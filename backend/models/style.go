package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Style represents style model
type Style struct {
	StyleNo        string     `json:"style_no" db:"STYLE_NO"`
	BuyerShortName string     `json:"buyer_short_name" db:"BUYER_SHORT_NAME"`
	Unit           string     `json:"unit" db:"UNIT"`
	TB             *string    `json:"t_b" db:"T_B"`
	SubCategory    *string    `json:"sub_category" db:"SUB_CATEGORY"`
	Fabric         *string    `json:"fabric" db:"FABRIC"`
	SmvAccum      *float64   `json:"smv_accum" db:"SMV_ACCUM"`
	EStyleNo      *string    `json:"e_style_no" db:"E_STYLE_NO"`
	DeleteStatus  int        `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}

// PaginatedStyles represents paginated response
type PaginatedStyles struct {
	Styles      []*Style `json:"styles"`
	TotalItems  int     `json:"total_items"`
	TotalPages  int     `json:"total_pages"`
	CurrentPage int     `json:"current_page"`
	PageSize    int     `json:"page_size"`
}

// StyleRequest represents request payload for style operations
type StyleRequest struct {
	StyleNo        string   `json:"style_no"`
	BuyerShortName string   `json:"buyer_short_name"`
	Unit           string   `json:"unit"`
	TB             *string  `json:"t_b"`
	SubCategory    *string  `json:"sub_category"`
	Fabric         *string  `json:"fabric"`
	SmvAccum      *float64 `json:"smv_accum"`
}

// StyleModel handles style database operations
type StyleModel struct {
	DB *sql.DB
}

// NewStyleModel creates new style model instance
func NewStyleModel(db *sql.DB) *StyleModel {
	return &StyleModel{DB: db}
}

// generateEStyleNo generates E-Style number based on style number and T/B value
func generateEStyleNo(styleNo string, tb *string) string {
	// Remove any existing T/B suffix from style number
	baseStyleNo := strings.TrimRight(strings.TrimRight(styleNo, ")"), "(T1)(B1")
	
	if tb != nil {
		return fmt.Sprintf("E-%s(%s)", baseStyleNo, *tb)
	}
	return fmt.Sprintf("E-%s", baseStyleNo)
}

// Create creates a new style
func (m *StyleModel) Create(style *StyleRequest) (*Style, error) {
	// Validate unit and T/B combination
	if style.Unit == "SET" && style.TB == nil {
		return nil, fmt.Errorf("T/B value is required when unit is SET")
	}
	if style.Unit == "PCS" && style.TB != nil {
		style.TB = nil // Clear T/B when unit is PCS
	}

	// Generate E-Style number
	eStyleNo := generateEStyleNo(style.StyleNo, style.TB)

	query := `INSERT INTO hs_ord_style (
		STYLE_NO, BUYER_SHORT_NAME, UNIT, T_B, SUB_CATEGORY, FABRIC, SMV_ACCUM, E_STYLE_NO, 
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := m.DB.Exec(query,
		style.StyleNo,
		style.BuyerShortName,
		style.Unit,
		style.TB,
		style.SubCategory,
		style.Fabric,
		style.SmvAccum,
		eStyleNo,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create style: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("no rows affected when creating style")
	}

	return m.GetByStyleNo(style.StyleNo)
}

// GetByStyleNo gets style by style number
func (m *StyleModel) GetByStyleNo(styleNo string) (*Style, error) {
	query := `SELECT STYLE_NO, BUYER_SHORT_NAME, UNIT, T_B, SUB_CATEGORY, FABRIC, 
              SMV_ACCUM, E_STYLE_NO, DELETE_STATUS, created_at, updated_at
              FROM hs_ord_style
              WHERE STYLE_NO = ?`

	var style Style
	err := m.DB.QueryRow(query, styleNo).Scan(
		&style.StyleNo,
		&style.BuyerShortName,
		&style.Unit,
		&style.TB,
		&style.SubCategory,
		&style.Fabric,
		&style.SmvAccum,
		&style.EStyleNo,
		&style.DeleteStatus,
		&style.CreatedAt,
		&style.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("style not found")
		}
		return nil, fmt.Errorf("failed to get style: %v", err)
	}

	return &style, nil
}

// GetAll gets all styles with pagination and search
func (m *StyleModel) GetAll(page, pageSize int, search string) (*PaginatedStyles, error) {
	// Base conditions for search
	searchCondition := ""
	searchArgs := []interface{}{}

	if search != "" {
		searchCondition = `WHERE (STYLE_NO LIKE ? OR BUYER_SHORT_NAME LIKE ? 
                          OR SUB_CATEGORY LIKE ? OR FABRIC LIKE ?)`
		searchPattern := "%" + search + "%"
		searchArgs = append(searchArgs, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Get total count first with search condition
	var totalItems int
	countQuery := `SELECT COUNT(*) FROM hs_ord_style ` + searchCondition
	err := m.DB.QueryRow(countQuery, searchArgs...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Calculate total pages
	totalPages := (totalItems + pageSize - 1) / pageSize

	// Main query with pagination and search
	query := `
		SELECT STYLE_NO, BUYER_SHORT_NAME, UNIT, T_B, SUB_CATEGORY, FABRIC, 
		       SMV_ACCUM, E_STYLE_NO, DELETE_STATUS, created_at, updated_at 
		FROM hs_ord_style 
		` + searchCondition + `
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	// Add pagination parameters to search args
	searchArgs = append(searchArgs, pageSize, (page-1)*pageSize)

	rows, err := m.DB.Query(query, searchArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get styles: %v", err)
	}
	defer rows.Close()

	var styles []*Style
	for rows.Next() {
		var style Style
		err := rows.Scan(
			&style.StyleNo,
			&style.BuyerShortName,
			&style.Unit,
			&style.TB,
			&style.SubCategory,
			&style.Fabric,
			&style.SmvAccum,
			&style.EStyleNo,
			&style.DeleteStatus,
			&style.CreatedAt,
			&style.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan style: %v", err)
		}
		styles = append(styles, &style)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate styles: %v", err)
	}

	return &PaginatedStyles{
		Styles:      styles,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}

// Update updates style information
func (m *StyleModel) Update(styleNo string, style *StyleRequest) (*Style, error) {
	// Validate unit and T/B combination
	if style.Unit == "SET" && style.TB == nil {
		return nil, fmt.Errorf("T/B value is required when unit is SET")
	}
	if style.Unit == "PCS" && style.TB != nil {
		style.TB = nil // Clear T/B when unit is PCS
	}

	// Generate new E-Style number
	eStyleNo := generateEStyleNo(style.StyleNo, style.TB)

	query := `UPDATE hs_ord_style SET 
		BUYER_SHORT_NAME = ?, 
		UNIT = ?, 
		T_B = ?, 
		SUB_CATEGORY = ?, 
		FABRIC = ?, 
		SMV_ACCUM = ?,
		E_STYLE_NO = ?,
		updated_at = CURRENT_TIMESTAMP 
		WHERE STYLE_NO = ?`

	result, err := m.DB.Exec(query,
		style.BuyerShortName,
		style.Unit,
		style.TB,
		style.SubCategory,
		style.Fabric,
		style.SmvAccum,
		eStyleNo,
		styleNo,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update style: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("style not found")
	}

	return m.GetByStyleNo(styleNo)
}

// Delete soft deletes style by style number
func (m *StyleModel) Delete(styleNo string) error {
	query := `UPDATE hs_ord_style SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP 
              WHERE STYLE_NO = ?`
	result, err := m.DB.Exec(query, styleNo)
	if err != nil {
		return fmt.Errorf("failed to delete style: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("style not found")
	}

	return nil
}

// Activate reactivates a soft-deleted style
func (m *StyleModel) Activate(styleNo string) error {
	query := `UPDATE hs_ord_style SET DELETE_STATUS = 0, updated_at = CURRENT_TIMESTAMP 
              WHERE STYLE_NO = ?`
	result, err := m.DB.Exec(query, styleNo)
	if err != nil {
		return fmt.Errorf("failed to activate style: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("style not found")
	}

	return nil
}

// HardDelete permanently deletes style from database
func (m *StyleModel) HardDelete(styleNo string) error {
	query := `DELETE FROM hs_ord_style WHERE STYLE_NO = ?`
	result, err := m.DB.Exec(query, styleNo)
	if err != nil {
		return fmt.Errorf("failed to delete style: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("style not found")
	}

	return nil
} 