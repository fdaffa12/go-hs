package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Buyer represents buyer model
type Buyer struct {
	ShortName     string     `json:"short_name" db:"BUYER_SHORT_NAME"`
	LongName      string     `json:"long_name" db:"BUYER_LONG_NAME"`
	DeleteStatus  int        `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}

// PaginatedBuyers represents paginated response
type PaginatedBuyers struct {
	Buyers       []*Buyer `json:"buyers"`
	TotalItems   int     `json:"total_items"`
	TotalPages   int     `json:"total_pages"`
	CurrentPage  int     `json:"current_page"`
	PageSize     int     `json:"page_size"`
}

// BuyerRequest represents request payload for buyer operations
type BuyerRequest struct {
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}

// BuyerModel handles buyer database operations
type BuyerModel struct {
	DB *sql.DB
}

// NewBuyerModel creates new buyer model instance
func NewBuyerModel(db *sql.DB) *BuyerModel {
	return &BuyerModel{DB: db}
}

// Create creates a new buyer
func (m *BuyerModel) Create(buyer *BuyerRequest) (*Buyer, error) {
	query := `INSERT INTO hs_ord_buyer (BUYER_SHORT_NAME, BUYER_LONG_NAME) VALUES (?, ?)`
	fmt.Printf("Executing query: %s with values: %v, %v\n", query, buyer.ShortName, buyer.LongName)

	result, err := m.DB.Exec(query, buyer.ShortName, buyer.LongName)
	if err != nil {
		return nil, fmt.Errorf("failed to create buyer: %v", err)
	}

	// Check if any rows were affected
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("no rows affected when creating buyer")
	}

	// Return created buyer
	return m.GetByShortName(buyer.ShortName)
}

// GetByShortName gets buyer by short name
func (m *BuyerModel) GetByShortName(shortName string) (*Buyer, error) {
	query := `SELECT BUYER_SHORT_NAME, BUYER_LONG_NAME, DELETE_STATUS, created_at, updated_at
              FROM hs_ord_buyer
              WHERE BUYER_SHORT_NAME = ?`
	fmt.Printf("Executing query: %s with value: %v\n", query, shortName)

	row := m.DB.QueryRow(query, shortName)

	var buyer Buyer
	err := row.Scan(&buyer.ShortName, &buyer.LongName, &buyer.DeleteStatus, &buyer.CreatedAt, &buyer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("buyer not found")
		}
		return nil, fmt.Errorf("failed to get buyer: %v", err)
	}

	return &buyer, nil
}

// GetAll gets all buyers with pagination and search
func (m *BuyerModel) GetAll(page, pageSize int, search string) (*PaginatedBuyers, error) {
    // Base conditions for search
    searchCondition := ""
    searchArgs := []interface{}{}
    
    if search != "" {
        searchCondition = "WHERE (BUYER_SHORT_NAME LIKE ? OR BUYER_LONG_NAME LIKE ?)"
        searchPattern := "%" + search + "%"
        searchArgs = append(searchArgs, searchPattern, searchPattern)
    }

    // Get total count first with search condition
    var totalItems int
    countQuery := `SELECT COUNT(*) FROM hs_ord_buyer ` + searchCondition
    err := m.DB.QueryRow(countQuery, searchArgs...).Scan(&totalItems)
    if err != nil {
        return nil, fmt.Errorf("failed to get total count: %v", err)
    }

    // Calculate total pages
    totalPages := (totalItems + pageSize - 1) / pageSize

    // Main query with pagination and search
    query := `
        SELECT BUYER_SHORT_NAME, BUYER_LONG_NAME, DELETE_STATUS, created_at, updated_at 
        FROM hs_ord_buyer 
        ` + searchCondition + `
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `
    // Add pagination parameters to search args
    searchArgs = append(searchArgs, pageSize, (page-1)*pageSize)
    
    rows, err := m.DB.Query(query, searchArgs...)
    if err != nil {
        return nil, fmt.Errorf("failed to get buyers: %v", err)
    }
    defer rows.Close()

    var buyers []*Buyer
    for rows.Next() {
        var buyer Buyer
        err := rows.Scan(
            &buyer.ShortName,
            &buyer.LongName,
            &buyer.DeleteStatus,
            &buyer.CreatedAt,
            &buyer.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan buyer: %v", err)
        }
        buyers = append(buyers, &buyer)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to iterate buyers: %v", err)
    }

    return &PaginatedBuyers{
        Buyers:      buyers,
        TotalItems:  totalItems,
        TotalPages:  totalPages,
        CurrentPage: page,
        PageSize:    pageSize,
    }, nil
}

// Update updates buyer information
func (m *BuyerModel) Update(shortName string, buyer *BuyerRequest) (*Buyer, error) {
	query := `UPDATE hs_ord_buyer SET BUYER_LONG_NAME = ?, updated_at = CURRENT_TIMESTAMP WHERE BUYER_SHORT_NAME = ?`
	result, err := m.DB.Exec(query, buyer.LongName, shortName)
	if err != nil {
		return nil, fmt.Errorf("failed to update buyer: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("buyer not found")
	}

	return m.GetByShortName(shortName)
}

// Delete soft deletes buyer by short name
func (m *BuyerModel) Delete(shortName string) error {
	query := `UPDATE hs_ord_buyer SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP WHERE BUYER_SHORT_NAME = ?`
	result, err := m.DB.Exec(query, shortName)
	if err != nil {
		return fmt.Errorf("failed to delete buyer: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("buyer not found")
	}

	return nil
}

// Activate reactivates a soft-deleted buyer
func (m *BuyerModel) Activate(shortName string) error {
	query := `UPDATE hs_ord_buyer SET DELETE_STATUS = 0, updated_at = CURRENT_TIMESTAMP WHERE BUYER_SHORT_NAME = ?`
	result, err := m.DB.Exec(query, shortName)
	if err != nil {
		return fmt.Errorf("failed to activate buyer: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("buyer not found")
	}

	return nil
}

// HardDelete permanently deletes buyer from database
func (m *BuyerModel) HardDelete(shortName string) error {
    query := `DELETE FROM hs_ord_buyer WHERE BUYER_SHORT_NAME = ?`
    result, err := m.DB.Exec(query, shortName)
    if err != nil {
        return fmt.Errorf("failed to delete buyer: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("buyer not found")
    }

    return nil
} 