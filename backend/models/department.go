package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Department represents department model
type Department struct {
	ShortName     string     `json:"short_name" db:"DEPT_SHORT_NAME"`
	LongName      string     `json:"long_name" db:"DEPT_LONG_NAME"`
	DeleteStatus  bool       `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}

// DepartmentRequest represents request payload for department operations
type DepartmentRequest struct {
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}

// DepartmentModel handles department database operations
type DepartmentModel struct {
	DB *sql.DB
}

// NewDepartmentModel creates new department model instance
func NewDepartmentModel(db *sql.DB) *DepartmentModel {
	return &DepartmentModel{DB: db}
}

// Create creates a new department
func (m *DepartmentModel) Create(dept *DepartmentRequest) (*Department, error) {
	query := `INSERT INTO hs_mst_departement (DEPT_SHORT_NAME, DEPT_LONG_NAME) VALUES (?, ?)`
	fmt.Printf("Executing query: %s with values: %v, %v\n", query, dept.ShortName, dept.LongName) // Add logging

	result, err := m.DB.Exec(query, dept.ShortName, dept.LongName)
	if err != nil {
		return nil, fmt.Errorf("failed to create department: %v", err)
	}

	// Check if any rows were affected
	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("no rows affected when creating department")
	}

	// Return created department
	return m.GetByShortName(dept.ShortName)
}

// GetByShortName gets department by short name
func (m *DepartmentModel) GetByShortName(shortName string) (*Department, error) {
	query := `SELECT DEPT_SHORT_NAME, DEPT_LONG_NAME, DELETE_STATUS, created_at, updated_at
              FROM hs_mst_departement
              WHERE DEPT_SHORT_NAME = ?`
	fmt.Printf("Executing query: %s with value: %v\n", query, shortName) // Add logging

	row := m.DB.QueryRow(query, shortName)

	var dept Department
	err := row.Scan(&dept.ShortName, &dept.LongName, &dept.DeleteStatus, &dept.CreatedAt, &dept.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %v", err)
	}

	return &dept, nil
}

// GetAll gets all hs_mst_departement
func (m *DepartmentModel) GetAll() ([]*Department, error) {
    query := `
        SELECT DEPT_SHORT_NAME, DEPT_LONG_NAME, DELETE_STATUS, created_at, updated_at 
        FROM hs_mst_departement 
        ORDER BY created_at DESC
    `
    rows, err := m.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to get hs_mst_departement: %v", err)
    }
    defer rows.Close()

    var departments []*Department
    for rows.Next() {
        var dept Department
        err := rows.Scan(
            &dept.ShortName,
            &dept.LongName,
            &dept.DeleteStatus,
            &dept.CreatedAt,
            &dept.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan department: %v", err)
        }
        departments = append(departments, &dept)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("failed to iterate departments: %v", err)
    }

    return departments, nil
}

// Update updates department information
func (m *DepartmentModel) Update(shortName string, dept *DepartmentRequest) (*Department, error) {
	query := `UPDATE hs_mst_departement SET DEPT_LONG_NAME = ?, updated_at = CURRENT_TIMESTAMP WHERE DEPT_SHORT_NAME = ?`
	result, err := m.DB.Exec(query, dept.LongName, shortName)
	if err != nil {
		return nil, fmt.Errorf("failed to update department: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("department not found")
	}

	return m.GetByShortName(shortName)
}

// Delete soft deletes department by short name
func (m *DepartmentModel) Delete(shortName string) error {
	query := `UPDATE hs_mst_departement SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP WHERE DEPT_SHORT_NAME = ?`
	result, err := m.DB.Exec(query, shortName)
	if err != nil {
		return fmt.Errorf("failed to delete department: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}

	return nil
}

// Activate reactivates a soft-deleted department
func (m *DepartmentModel) Activate(shortName string) error {
	query := `UPDATE hs_mst_departement SET DELETE_STATUS = 0, updated_at = CURRENT_TIMESTAMP WHERE DEPT_SHORT_NAME = ?`
	result, err := m.DB.Exec(query, shortName)
	if err != nil {
		return fmt.Errorf("failed to activate department: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}

	return nil
}

// HardDelete permanently deletes department from database
func (m *DepartmentModel) HardDelete(shortName string) error {
    query := `DELETE FROM hs_mst_departement WHERE DEPT_SHORT_NAME = ?`
    result, err := m.DB.Exec(query, shortName)
    if err != nil {
        return fmt.Errorf("failed to delete department: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("department not found")
    }

    return nil
} 