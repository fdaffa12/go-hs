package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Holiday represents holiday model
type Holiday struct {
	HolidayID     int        `json:"holiday_id" db:"HOLIDAY_ID"`
	HolidayDate   string     `json:"holiday_date" db:"HOLIDAY_DATE"`
	HolidayName   string     `json:"holiday_name" db:"HOLIDAY_NAME"`
	HolidayType   string     `json:"holiday_type" db:"HOLIDAY_TYPE"`
	DeleteStatus  int        `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt     *time.Time `json:"created_at" db:"CREATED_AT"`
	UpdatedAt     *time.Time `json:"updated_at" db:"UPDATED_AT"`
}

// PaginatedHolidays represents paginated response
type PaginatedHolidays struct {
	Holidays     []*Holiday `json:"holidays"`
	TotalItems   int       `json:"total_items"`
	TotalPages   int       `json:"total_pages"`
	CurrentPage  int       `json:"current_page"`
	PageSize     int       `json:"page_size"`
}

// HolidayRequest represents request payload for holiday operations
type HolidayRequest struct {
	HolidayDate string `json:"holiday_date"`
	HolidayName string `json:"holiday_name"`
	HolidayType string `json:"holiday_type"`
}

// HolidayModel handles holiday database operations
type HolidayModel struct {
	DB *sql.DB
}

// NewHolidayModel creates new holiday model instance
func NewHolidayModel(db *sql.DB) *HolidayModel {
	return &HolidayModel{DB: db}
}

// Create creates a new holiday
func (m *HolidayModel) Create(holiday *HolidayRequest) (*Holiday, error) {
	query := `INSERT INTO hs_set_holidays (HOLIDAY_DATE, HOLIDAY_NAME, HOLIDAY_TYPE) VALUES (?, ?, ?)`
	fmt.Printf("Executing query: %s with values: %v, %v, %v\n", query, holiday.HolidayDate, holiday.HolidayName, holiday.HolidayType)

	result, err := m.DB.Exec(query, holiday.HolidayDate, holiday.HolidayName, holiday.HolidayType)
	if err != nil {
		return nil, fmt.Errorf("failed to create holiday: %v", err)
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}

	// Return created holiday
	return m.GetByID(int(id))
}

// GetByID gets holiday by ID
func (m *HolidayModel) GetByID(id int) (*Holiday, error) {
	query := `SELECT HOLIDAY_ID, HOLIDAY_DATE, HOLIDAY_NAME, HOLIDAY_TYPE, DELETE_STATUS, CREATED_AT, UPDATED_AT
              FROM hs_set_holidays
              WHERE HOLIDAY_ID = ?`
	fmt.Printf("Executing query: %s with value: %v\n", query, id)

	row := m.DB.QueryRow(query, id)

	var holiday Holiday
	err := row.Scan(&holiday.HolidayID, &holiday.HolidayDate, &holiday.HolidayName, &holiday.HolidayType, &holiday.DeleteStatus, &holiday.CreatedAt, &holiday.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("holiday not found")
		}
		return nil, fmt.Errorf("failed to get holiday: %v", err)
	}

	return &holiday, nil
}

// GetAll gets all holidays with pagination and search
func (m *HolidayModel) GetAll(page, pageSize int, search string) (*PaginatedHolidays, error) {
    // Base conditions for search
    searchCondition := ""
    searchArgs := []interface{}{}
    
    if search != "" {
        searchCondition = "WHERE (HOLIDAY_NAME LIKE ? OR HOLIDAY_TYPE LIKE ?)"
        searchPattern := "%" + search + "%"
        searchArgs = append(searchArgs, searchPattern, searchPattern)
    }

    // Get total count first with search condition
    var totalItems int
    countQuery := `SELECT COUNT(*) FROM hs_set_holidays ` + searchCondition
    err := m.DB.QueryRow(countQuery, searchArgs...).Scan(&totalItems)
    if err != nil {
        return nil, fmt.Errorf("failed to get total count: %v", err)
    }

    // Calculate total pages
    totalPages := (totalItems + pageSize - 1) / pageSize

    // Main query with pagination and search
    query := `
        SELECT HOLIDAY_ID, HOLIDAY_DATE, HOLIDAY_NAME, HOLIDAY_TYPE, DELETE_STATUS, CREATED_AT, UPDATED_AT 
        FROM hs_set_holidays 
        ` + searchCondition + `
        ORDER BY HOLIDAY_DATE DESC
        LIMIT ? OFFSET ?
    `
    // Add pagination parameters to search args
    searchArgs = append(searchArgs, pageSize, (page-1)*pageSize)
    
    rows, err := m.DB.Query(query, searchArgs...)
    if err != nil {
        return nil, fmt.Errorf("failed to get holidays: %v", err)
    }
    defer rows.Close()

    var holidays []*Holiday
    for rows.Next() {
        var holiday Holiday
        err := rows.Scan(
            &holiday.HolidayID,
            &holiday.HolidayDate,
            &holiday.HolidayName,
            &holiday.HolidayType,
            &holiday.DeleteStatus,
            &holiday.CreatedAt,
            &holiday.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan holiday: %v", err)
        }
        holidays = append(holidays, &holiday)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to iterate holidays: %v", err)
    }

    return &PaginatedHolidays{
        Holidays:    holidays,
        TotalItems:  totalItems,
        TotalPages:  totalPages,
        CurrentPage: page,
        PageSize:    pageSize,
    }, nil
}

// Update updates holiday information
func (m *HolidayModel) Update(id int, holiday *HolidayRequest) (*Holiday, error) {
	query := `UPDATE hs_set_holidays SET HOLIDAY_DATE = ?, HOLIDAY_NAME = ?, HOLIDAY_TYPE = ?, UPDATED_AT = CURRENT_TIMESTAMP WHERE HOLIDAY_ID = ?`
	result, err := m.DB.Exec(query, holiday.HolidayDate, holiday.HolidayName, holiday.HolidayType, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update holiday: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("holiday not found")
	}

	return m.GetByID(id)
}

// Delete soft deletes holiday by ID
func (m *HolidayModel) Delete(id int) error {
	query := `UPDATE hs_set_holidays SET DELETE_STATUS = 1, UPDATED_AT = CURRENT_TIMESTAMP WHERE HOLIDAY_ID = ?`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete holiday: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("holiday not found")
	}

	return nil
}

// Activate reactivates a soft-deleted holiday
func (m *HolidayModel) Activate(id int) error {
	query := `UPDATE hs_set_holidays SET DELETE_STATUS = 0, UPDATED_AT = CURRENT_TIMESTAMP WHERE HOLIDAY_ID = ?`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to activate holiday: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("holiday not found")
	}

	return nil
}

// HardDelete permanently deletes holiday from database
func (m *HolidayModel) HardDelete(id int) error {
    query := `DELETE FROM hs_set_holidays WHERE HOLIDAY_ID = ?`
    result, err := m.DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete holiday: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("holiday not found")
    }

    return nil
} 