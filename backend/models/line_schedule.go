package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// LineSchedule represents a line schedule record
type LineSchedule struct {
	RowID          int64      `json:"row_id" db:"ROWID"`
	IDRegistrasi   string     `json:"id_registrasi" db:"ID_REGISTRASI"`
	Date           time.Time  `json:"date" db:"DATE"`
	Factory        string     `json:"factory" db:"FACTORY"`
	Line           string     `json:"line" db:"LINE"`
	BuyerShortName string     `json:"buyer_short_name" db:"BUYER_SHORT_NAME"`
	StyleNo        string     `json:"style_no" db:"STYLE_NO"`
	StartDate      *time.Time `json:"start_date" db:"START_DATE"`
	NumberOfMP     int        `json:"number_of_mp" db:"NUMBER_OF_MP"`
	WorkingDay     int        `json:"working_day" db:"WORKING_DAY"`
	DeleteStatus   int        `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

// PaginatedLineSchedules represents paginated line schedules data
type PaginatedLineSchedules struct {
	LineSchedules []*LineSchedule `json:"line_schedules"`
	TotalItems    int            `json:"total_items"`
	TotalPages    int            `json:"total_pages"`
	CurrentPage   int            `json:"current_page"`
	PageSize      int            `json:"page_size"`
}

// LineScheduleRequest represents the request body for creating/updating a line schedule
type LineScheduleRequest struct {
	IDRegistrasi   string    `json:"id_registrasi"`
	Date           time.Time `json:"date"`
	Factory        string    `json:"factory"`
	Line           string    `json:"line"`
	BuyerShortName string    `json:"buyer_short_name"`
	StyleNo        string    `json:"style_no"`
	StartDate      *time.Time `json:"start_date"`
	NumberOfMP     int       `json:"number_of_mp"`
	WorkingDay     int       `json:"working_day"`
}

// UnmarshalJSON custom unmarshaler for LineScheduleRequest
func (r *LineScheduleRequest) UnmarshalJSON(data []byte) error {
	// Create an auxiliary struct with string fields for dates
	type Aux struct {
		IDRegistrasi   string `json:"id_registrasi"`
		Date           string `json:"date"`
		Factory        string `json:"factory"`
		Line           string `json:"line"`
		BuyerShortName string `json:"buyer_short_name"`
		StyleNo        string `json:"style_no"`
		StartDate      string `json:"start_date"`
		NumberOfMP     int    `json:"number_of_mp"`
		WorkingDay     int    `json:"working_day"`
	}

	// Parse into auxiliary struct
	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse dates
	date, err := time.Parse("2006-01-02", aux.Date)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}

	var startDate *time.Time
	if aux.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", aux.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format: %v", err)
		}
		startDate = &parsedStartDate
	}

	// Assign values
	r.IDRegistrasi = aux.IDRegistrasi
	r.Date = date
	r.Factory = aux.Factory
	r.Line = aux.Line
	r.BuyerShortName = aux.BuyerShortName
	r.StyleNo = aux.StyleNo
	r.StartDate = startDate
	r.NumberOfMP = aux.NumberOfMP
	r.WorkingDay = aux.WorkingDay

	return nil
}

// LineScheduleModel handles database operations for line schedules
type LineScheduleModel struct {
	DB *sql.DB
}

// NewLineScheduleModel creates a new LineScheduleModel instance
func NewLineScheduleModel(db *sql.DB) *LineScheduleModel {
	return &LineScheduleModel{DB: db}
}

// Create adds a new line schedule
func (m *LineScheduleModel) Create(schedule *LineScheduleRequest) (*LineSchedule, error) {
	query := `
		INSERT INTO hs_wsb_lineschedule (
			ID_REGISTRASI, DATE, FACTORY, LINE, BUYER_SHORT_NAME, 
			STYLE_NO, START_DATE, NUMBER_OF_MP, WORKING_DAY, DELETE_STATUS
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`

	result, err := m.DB.Exec(
		query,
		schedule.IDRegistrasi,
		schedule.Date,
		schedule.Factory,
		schedule.Line,
		schedule.BuyerShortName,
		schedule.StyleNo,
		schedule.StartDate,
		schedule.NumberOfMP,
		schedule.WorkingDay,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create line schedule: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return m.GetByID(id)
}

// GetByID retrieves a line schedule by its ID
func (m *LineScheduleModel) GetByID(id int64) (*LineSchedule, error) {
	schedule := &LineSchedule{}
	query := `
		SELECT * FROM hs_wsb_lineschedule 
		WHERE ROWID = ?
	`

	err := m.DB.QueryRow(query, id).Scan(
		&schedule.RowID,
		&schedule.IDRegistrasi,
		&schedule.Date,
		&schedule.Factory,
		&schedule.Line,
		&schedule.BuyerShortName,
		&schedule.StyleNo,
		&schedule.StartDate,
		&schedule.NumberOfMP,
		&schedule.WorkingDay,
		&schedule.DeleteStatus,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("line schedule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get line schedule: %v", err)
	}

	return schedule, nil
}

// GetAll retrieves all line schedules with pagination and search
func (m *LineScheduleModel) GetAll(page, pageSize int, filterStr string) (*PaginatedLineSchedules, error) {
	offset := (page - 1) * pageSize

	// Base query
	baseQuery := `FROM hs_wsb_lineschedule WHERE 1=1`
	args := []interface{}{}

	// Parse filters
	filters := make(map[string]string)
	if filterStr != "" {
		pairs := strings.Split(filterStr, "&")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				filters[kv[0]] = kv[1]
			}
		}
	}

	// Add factory filter if provided
	if factory, ok := filters["factory"]; ok && factory != "" {
		baseQuery += ` AND FACTORY = ?`
		args = append(args, factory)
	}

	// Add date filter if provided
	if date, ok := filters["date"]; ok && date != "" {
		baseQuery += ` AND DATE = ?`
		args = append(args, date)
	}

	countQuery := `SELECT COUNT(*) ` + baseQuery
	selectQuery := `SELECT * ` + baseQuery

	// Add pagination
	selectQuery = selectQuery + ` ORDER BY DATE DESC, FACTORY, LINE LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	// Get total count
	var totalItems int
	err := m.DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Calculate total pages
	totalPages := (totalItems + pageSize - 1) / pageSize

	// Get paginated data
	rows, err := m.DB.Query(selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get line schedules: %v", err)
	}
	defer rows.Close()

	var schedules []*LineSchedule
	for rows.Next() {
		schedule := &LineSchedule{}
		err := rows.Scan(
			&schedule.RowID,
			&schedule.IDRegistrasi,
			&schedule.Date,
			&schedule.Factory,
			&schedule.Line,
			&schedule.BuyerShortName,
			&schedule.StyleNo,
			&schedule.StartDate,
			&schedule.NumberOfMP,
			&schedule.WorkingDay,
			&schedule.DeleteStatus,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan line schedule: %v", err)
		}
		schedules = append(schedules, schedule)
	}

	return &PaginatedLineSchedules{
		LineSchedules: schedules,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  page,
		PageSize:     pageSize,
	}, nil
}

// Update modifies an existing line schedule
func (m *LineScheduleModel) Update(id int64, schedule *LineScheduleRequest) (*LineSchedule, error) {
	query := `
		UPDATE hs_wsb_lineschedule 
		SET 
			ID_REGISTRASI = ?,
			DATE = ?,
			FACTORY = ?,
			LINE = ?,
			BUYER_SHORT_NAME = ?,
			STYLE_NO = ?,
			START_DATE = ?,
			NUMBER_OF_MP = ?,
			WORKING_DAY = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE ROWID = ?
	`

	result, err := m.DB.Exec(
		query,
		schedule.IDRegistrasi,
		schedule.Date,
		schedule.Factory,
		schedule.Line,
		schedule.BuyerShortName,
		schedule.StyleNo,
		schedule.StartDate,
		schedule.NumberOfMP,
		schedule.WorkingDay,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update line schedule: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("line schedule not found")
	}

	return m.GetByID(id)
}

// Delete soft deletes a line schedule
func (m *LineScheduleModel) Delete(id int64) error {
	query := `
		UPDATE hs_wsb_lineschedule 
		SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP
		WHERE ROWID = ?
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete line schedule: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("line schedule not found")
	}

	return nil
}

// Activate reactivates a soft-deleted line schedule
func (m *LineScheduleModel) Activate(id int64) error {
	query := `
		UPDATE hs_wsb_lineschedule 
		SET DELETE_STATUS = 0, updated_at = CURRENT_TIMESTAMP
		WHERE ROWID = ?
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to activate line schedule: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("line schedule not found")
	}

	return nil
}

// HardDelete permanently deletes a line schedule
func (m *LineScheduleModel) HardDelete(id int64) error {
	query := `DELETE FROM hs_wsb_lineschedule WHERE ROWID = ?`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete line schedule: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("line schedule not found")
	}

	return nil
} 

// CalculateWorkingDays calculates working days between start date and end date, excluding weekends and holidays
func (m *LineScheduleModel) CalculateWorkingDays(startDate, endDate time.Time) (int, error) {
    // Get holidays between start and end date
    query := `
        SELECT DATE_FORMAT(HOLIDAY_DATE, '%Y-%m-%d') as holiday_date
        FROM hs_set_holidays 
        WHERE DATE(HOLIDAY_DATE) BETWEEN ? AND ? 
        AND DELETE_STATUS = 0
    `
    
    // Format dates for SQL query
    startStr := startDate.Format("2006-01-02")
    endStr := endDate.Format("2006-01-02")
    
    fmt.Printf("Checking holidays between %s and %s\n", startStr, endStr)
    
    rows, err := m.DB.Query(query, startStr, endStr)
    if err != nil {
        return 0, fmt.Errorf("failed to get holidays: %v", err)
    }
    defer rows.Close()

    // Store holidays in a map for O(1) lookup
    holidays := make(map[string]bool)
    for rows.Next() {
        var holidayDate string
        if err := rows.Scan(&holidayDate); err != nil {
            return 0, fmt.Errorf("failed to scan holiday date: %v", err)
        }
        holidays[holidayDate] = true
        fmt.Printf("Found holiday on: %s\n", holidayDate)
    }

    workingDays := 0
    currentDate := startDate

    // Iterate through each day
    for !currentDate.After(endDate) {
        currentDateStr := currentDate.Format("2006-01-02")
        
        // Skip weekends (Saturday = 6, Sunday = 0)
        if currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday {
            // Check if it's not a holiday
            isHoliday := holidays[currentDateStr]
            if !isHoliday {
                workingDays++
                fmt.Printf("Adding working day: %s\n", currentDateStr)
            } else {
                fmt.Printf("Skipping holiday: %s\n", currentDateStr)
            }
        } else {
            fmt.Printf("Skipping weekend: %s\n", currentDateStr)
        }
        
        currentDate = currentDate.AddDate(0, 0, 1)
    }

    fmt.Printf("Total working days calculated: %d\n", workingDays)
    return workingDays, nil
} 