package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Employee represents employee model
type Employee struct {
	NIK           string     `json:"nik" db:"NIK"`
	Name          string     `json:"name" db:"NAME"`
	DeptShortName string     `json:"dept_short_name" db:"DEPT_SHORT_NAME"`
	EnteranceDate time.Time  `json:"enterance_date" db:"ENTERANCE_DATE"`
	Title         string     `json:"title" db:"TITLE"`
	Gender        int        `json:"gender" db:"GENDER"`
	RfidID        *string    `json:"rfid_id" db:"RFID_ID"`
	WorkingArea   int        `json:"working_area" db:"WORKING_AREA"`
	DeleteStatus  bool       `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
	// Additional fields for joined data
	DeptLongName string `json:"dept_long_name" db:"DEPT_LONG_NAME"`
}

// EmployeeRequest represents request payload for employee operations
type EmployeeRequest struct {
	NIK           string `json:"nik"`
	Name          string `json:"name"`
	DeptShortName string `json:"dept_short_name"`
	EnteranceDate string `json:"enterance_date"` // Format: YYYY-MM-DD
	Title         string `json:"title"`
	Gender        int    `json:"gender"`
	RfidID        string `json:"rfid_id"`
	WorkingArea   int    `json:"working_area"`
}

// EmployeeModel handles employee database operations
type EmployeeModel struct {
	DB *sql.DB
}

// NewEmployeeModel creates new employee model instance
func NewEmployeeModel(db *sql.DB) *EmployeeModel {
	return &EmployeeModel{DB: db}
}

// Create creates a new employee
func (m *EmployeeModel) Create(emp *EmployeeRequest) (*Employee, error) {
	query := `
		INSERT INTO hs_hrd_employee (
			NIK, NAME, DEPT_SHORT_NAME, ENTERANCE_DATE, 
			TITLE, GENDER, RFID_ID, WORKING_AREA, DELETE_STATUS
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0)
	`

	result, err := m.DB.Exec(
		query,
		emp.NIK,
		emp.Name,
		emp.DeptShortName,
		emp.EnteranceDate,
		emp.Title,
		emp.Gender,
		sql.NullString{String: emp.RfidID, Valid: emp.RfidID != ""},
		emp.WorkingArea,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		return nil, fmt.Errorf("no rows affected when creating employee")
	}

	// Get the created employee
	return m.GetByNIK(emp.NIK)
}

// GetByNIK gets employee by NIK with department name
func (m *EmployeeModel) GetByNIK(nik string) (*Employee, error) {
	if nik == "" {
		return nil, fmt.Errorf("NIK cannot be empty")
	}

	query := `
		SELECT 
			e.NIK, e.NAME, e.DEPT_SHORT_NAME, e.ENTERANCE_DATE, 
			e.TITLE, e.GENDER, e.RFID_ID, e.WORKING_AREA, 
			e.DELETE_STATUS, e.created_at, e.updated_at,
			COALESCE(d.DEPT_LONG_NAME, '') as DEPT_LONG_NAME
		FROM hs_hrd_employee e
		LEFT JOIN departments d ON e.DEPT_SHORT_NAME = d.DEPT_SHORT_NAME
		WHERE e.NIK = ?
	`

	var emp Employee
	var rfidID sql.NullString
	err := m.DB.QueryRow(query, nik).Scan(
		&emp.NIK,
		&emp.Name,
		&emp.DeptShortName,
		&emp.EnteranceDate,
		&emp.Title,
		&emp.Gender,
		&rfidID,
		&emp.WorkingArea,
		&emp.DeleteStatus,
		&emp.CreatedAt,
		&emp.UpdatedAt,
		&emp.DeptLongName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %v", err)
	}

	if rfidID.Valid {
		emp.RfidID = &rfidID.String
	}

	return &emp, nil
}

// GetAll gets all hs_hrd_employee with department names
func (m *EmployeeModel) GetAll() ([]*Employee, error) {
	query := `
		SELECT 
			e.NIK, e.NAME, e.DEPT_SHORT_NAME, e.ENTERANCE_DATE, 
			e.TITLE, e.GENDER, e.RFID_ID, e.WORKING_AREA, 
			e.DELETE_STATUS, e.created_at, e.updated_at,
			COALESCE(d.DEPT_LONG_NAME, '') as DEPT_LONG_NAME
		FROM hs_hrd_employee e
		LEFT JOIN departments d ON e.DEPT_SHORT_NAME = d.DEPT_SHORT_NAME
		ORDER BY e.created_at DESC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get hs_hrd_employee: %v", err)
	}
	defer rows.Close()

	var employees []*Employee
	for rows.Next() {
		var emp Employee
		var rfidID sql.NullString
		var deptLongName sql.NullString
		err := rows.Scan(
			&emp.NIK,
			&emp.Name,
			&emp.DeptShortName,
			&emp.EnteranceDate,
			&emp.Title,
			&emp.Gender,
			&rfidID,
			&emp.WorkingArea,
			&emp.DeleteStatus,
			&emp.CreatedAt,
			&emp.UpdatedAt,
			&deptLongName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %v", err)
		}

		if rfidID.Valid {
			emp.RfidID = &rfidID.String
		}

		emp.DeptLongName = deptLongName.String // Will be empty string if NULL

		employees = append(employees, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate hs_hrd_employee: %v", err)
	}

	return employees, nil
}

// Update updates employee information
func (m *EmployeeModel) Update(nik string, emp *EmployeeRequest) (*Employee, error) {
	query := `
		UPDATE hs_hrd_employee 
		SET 
			NAME = ?, 
			DEPT_SHORT_NAME = ?, 
			ENTERANCE_DATE = ?,
			TITLE = ?,
			GENDER = ?,
			RFID_ID = ?,
			WORKING_AREA = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE NIK = ?
	`

	result, err := m.DB.Exec(
		query,
		emp.Name,
		emp.DeptShortName,
		emp.EnteranceDate,
		emp.Title,
		emp.Gender,
		sql.NullString{String: emp.RfidID, Valid: emp.RfidID != ""},
		emp.WorkingArea,
		nik,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update employee: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("employee not found")
	}

	return m.GetByNIK(nik)
}

// Delete soft deletes employee by NIK
func (m *EmployeeModel) Delete(nik string) error {
	query := `UPDATE hs_hrd_employee SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
	result, err := m.DB.Exec(query, nik)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// HardDelete permanently deletes employee from database
func (m *EmployeeModel) HardDelete(nik string) error {
	query := `DELETE FROM hs_hrd_employee WHERE NIK = ?`
	result, err := m.DB.Exec(query, nik)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// Activate reactivates a soft-deleted employee
func (m *EmployeeModel) Activate(nik string) error {
	query := `UPDATE hs_hrd_employee SET DELETE_STATUS = 0, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
	result, err := m.DB.Exec(query, nik)
	if err != nil {
		return fmt.Errorf("failed to activate employee: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
} 