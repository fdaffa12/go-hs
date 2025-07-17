package models

import (
	"database/sql"
	"fmt"
	"time"
)

// SewNumProcess represents a sewing numbering process record
type SewNumProcess struct {
	ID             int64      `json:"id" db:"id"`
	BuyerShortName string     `json:"buyer_short_name" db:"BUYER_SHORT_NAME"`
	StyleNo        string     `json:"style_no" db:"STYLE_NO"`
	NoProcess      int        `json:"no_process" db:"NO_PROCESS"`
	Category       string     `json:"category" db:"CATEGORY"`
	SubCategory    string     `json:"sub_category" db:"SUB_CATEGORY"`
	SmvProcGsd     *float64   `json:"smv_proc_gsd" db:"SMV_PROC_GSD"`
	SmvProcEst     float64    `json:"smv_proc_est" db:"SMV_PROC_EST"`
	ProcessNameEng *string    `json:"process_name_eng" db:"PROCESS_NAME_ENG"`
	ProcessNameInd *string    `json:"process_name_ind" db:"PROCESS_NAME_IND"`
	MachineCode    string     `json:"machine_code" db:"MACHINE_CODE"`
	LinkTo1        *int       `json:"link_to_1" db:"LINK_TO_1"`
	LinkTo2        *int       `json:"link_to_2" db:"LINK_TO_2"`
	LinkTo3        *int       `json:"link_to_3" db:"LINK_TO_3"`
	LinkTo4        *int       `json:"link_to_4" db:"LINK_TO_4"`
	LinkTo5        *int       `json:"link_to_5" db:"LINK_TO_5"`
	DeleteStatus   int        `json:"delete_status" db:"DELETE_STATUS"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

// SewNumProcessModel handles database operations for sew numbering process
type SewNumProcessModel struct {
	DB *sql.DB
}

// NewSewNumProcessModel creates a new SewNumProcessModel instance
func NewSewNumProcessModel(db *sql.DB) *SewNumProcessModel {
	return &SewNumProcessModel{DB: db}
}

// GetAll retrieves all sew numbering processes for a given style
func (m *SewNumProcessModel) GetAll(styleNo string) ([]*SewNumProcess, error) {
	query := `
		SELECT * FROM hs_sew_numprocess 
		WHERE STYLE_NO = ? 
		ORDER BY NO_PROCESS ASC
	`

	rows, err := m.DB.Query(query, styleNo)
	if err != nil {
		return nil, fmt.Errorf("failed to get sew numbering processes: %v", err)
	}
	defer rows.Close()

	var processes []*SewNumProcess
	for rows.Next() {
		process := &SewNumProcess{}
		err := rows.Scan(
			&process.ID,
			&process.BuyerShortName,
			&process.StyleNo,
			&process.NoProcess,
			&process.Category,
			&process.SubCategory,
			&process.SmvProcGsd,
			&process.SmvProcEst,
			&process.ProcessNameEng,
			&process.ProcessNameInd,
			&process.MachineCode,
			&process.LinkTo1,
			&process.LinkTo2,
			&process.LinkTo3,
			&process.LinkTo4,
			&process.LinkTo5,
			&process.DeleteStatus,
			&process.CreatedAt,
			&process.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sew numbering process: %v", err)
		}
		processes = append(processes, process)
	}

	return processes, nil
}

// Create adds a new sew numbering process
func (m *SewNumProcessModel) Create(process *SewNumProcess) error {
	query := `
		INSERT INTO hs_sew_numprocess (
			BUYER_SHORT_NAME, STYLE_NO, NO_PROCESS, CATEGORY, SUB_CATEGORY,
			SMV_PROC_GSD, SMV_PROC_EST, PROCESS_NAME_ENG, PROCESS_NAME_IND,
			MACHINE_CODE, LINK_TO_1, LINK_TO_2, LINK_TO_3, LINK_TO_4, LINK_TO_5,
			DELETE_STATUS, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err := m.DB.Exec(
		query,
		process.BuyerShortName,
		process.StyleNo,
		process.NoProcess,
		process.Category,
		process.SubCategory,
		process.SmvProcGsd,
		process.SmvProcEst,
		process.ProcessNameEng,
		process.ProcessNameInd,
		process.MachineCode,
		process.LinkTo1,
		process.LinkTo2,
		process.LinkTo3,
		process.LinkTo4,
		process.LinkTo5,
	)

	if err != nil {
		return fmt.Errorf("failed to create sew numbering process: %v", err)
	}

	return nil
}

// Update modifies an existing sew numbering process
func (m *SewNumProcessModel) Update(process *SewNumProcess) error {
	query := `
		UPDATE hs_sew_numprocess 
		SET 
			BUYER_SHORT_NAME = ?,
			STYLE_NO = ?,
			NO_PROCESS = ?,
			CATEGORY = ?,
			SUB_CATEGORY = ?,
			SMV_PROC_GSD = ?,
			SMV_PROC_EST = ?,
			PROCESS_NAME_ENG = ?,
			PROCESS_NAME_IND = ?,
			MACHINE_CODE = ?,
			LINK_TO_1 = ?,
			LINK_TO_2 = ?,
			LINK_TO_3 = ?,
			LINK_TO_4 = ?,
			LINK_TO_5 = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := m.DB.Exec(
		query,
		process.BuyerShortName,
		process.StyleNo,
		process.NoProcess,
		process.Category,
		process.SubCategory,
		process.SmvProcGsd,
		process.SmvProcEst,
		process.ProcessNameEng,
		process.ProcessNameInd,
		process.MachineCode,
		process.LinkTo1,
		process.LinkTo2,
		process.LinkTo3,
		process.LinkTo4,
		process.LinkTo5,
		process.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sew numbering process: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sew numbering process not found")
	}

	return nil
}

// Delete soft deletes a sew numbering process
func (m *SewNumProcessModel) Delete(id int64) error {
	query := `
		UPDATE hs_sew_numprocess 
		SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sew numbering process: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sew numbering process not found")
	}

	return nil
}

// BulkSave handles bulk create/update operations for sew numbering processes
func (m *SewNumProcessModel) BulkSave(processes []*SewNumProcess) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Prepare create statement
	createStmt, err := tx.Prepare(`
		INSERT INTO hs_sew_numprocess (
			BUYER_SHORT_NAME, STYLE_NO, NO_PROCESS, CATEGORY, SUB_CATEGORY,
			SMV_PROC_GSD, SMV_PROC_EST, PROCESS_NAME_ENG, PROCESS_NAME_IND,
			MACHINE_CODE, LINK_TO_1, LINK_TO_2, LINK_TO_3, LINK_TO_4, LINK_TO_5,
			DELETE_STATUS, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare create statement: %v", err)
	}
	defer createStmt.Close()

	// Prepare update statement
	updateStmt, err := tx.Prepare(`
		UPDATE hs_sew_numprocess 
		SET 
			BUYER_SHORT_NAME = ?,
			STYLE_NO = ?,
			NO_PROCESS = ?,
			CATEGORY = ?,
			SUB_CATEGORY = ?,
			SMV_PROC_GSD = ?,
			SMV_PROC_EST = ?,
			PROCESS_NAME_ENG = ?,
			PROCESS_NAME_IND = ?,
			MACHINE_CODE = ?,
			LINK_TO_1 = ?,
			LINK_TO_2 = ?,
			LINK_TO_3 = ?,
			LINK_TO_4 = ?,
			LINK_TO_5 = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer updateStmt.Close()

	for _, process := range processes {
		if process.ID > 0 {
			// Update existing process
			_, err = updateStmt.Exec(
				process.BuyerShortName,
				process.StyleNo,
				process.NoProcess,
				process.Category,
				process.SubCategory,
				process.SmvProcGsd,
				process.SmvProcEst,
				process.ProcessNameEng,
				process.ProcessNameInd,
				process.MachineCode,
				process.LinkTo1,
				process.LinkTo2,
				process.LinkTo3,
				process.LinkTo4,
				process.LinkTo5,
				process.ID,
			)
		} else {
			// Create new process
			_, err = createStmt.Exec(
				process.BuyerShortName,
				process.StyleNo,
				process.NoProcess,
				process.Category,
				process.SubCategory,
				process.SmvProcGsd,
				process.SmvProcEst,
				process.ProcessNameEng,
				process.ProcessNameInd,
				process.MachineCode,
				process.LinkTo1,
				process.LinkTo2,
				process.LinkTo3,
				process.LinkTo4,
				process.LinkTo5,
			)
		}

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// BulkDelete handles bulk soft delete operations
func (m *SewNumProcessModel) BulkDelete(ids []int64) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		UPDATE hs_sew_numprocess 
		SET DELETE_STATUS = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, id := range ids {
		_, err = stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete process: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Activate reactivates a soft-deleted sew numbering process
func (m *SewNumProcessModel) Activate(id int64) error {
	query := `
		UPDATE hs_sew_numprocess 
		SET DELETE_STATUS = 0, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to activate sew numbering process: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sew numbering process not found")
	}

	return nil
}

// HardDelete permanently deletes a sew numbering process
func (m *SewNumProcessModel) HardDelete(id int64) error {
	query := `DELETE FROM hs_sew_numprocess WHERE id = ?`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete sew numbering process: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sew numbering process not found")
	}

	return nil
} 