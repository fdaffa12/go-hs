package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents user model
type User struct {
	NIK             string     `json:"nik" db:"NIK"`
	Name            string     `json:"name" db:"NAME"`
	Email           string     `json:"email" db:"EMAIL"`
	Password        string     `json:"-" db:"PASSWORD"` // "-" means don't include in JSON
	Level           int        `json:"level" db:"LEVEL"`
	ProfilePicture  *string    `json:"profile_picture" db:"PROFILE_PICTURE"`
	RememberToken   *string    `json:"remember_token" db:"remember_token"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
}

// UserRequest represents request payload for user operations
type UserRequest struct {
	NIK      string `json:"nik"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Level    int    `json:"level"`
}

// UserResponse represents response payload for user operations
type UserResponse struct {
	NIK            string  `json:"nik"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Level          int     `json:"level"`
	ProfilePicture *string `json:"profile_picture"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// UserModel handles user database operations
type UserModel struct {
	DB *sql.DB
}

// NewUserModel creates new user model instance
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// CheckNIKInEmployees checks if NIK exists in employees table
func (m *UserModel) CheckNIKInEmployees(nik string) (bool, error) {
	query := `SELECT COUNT(*) FROM hs_hrd_employee WHERE NIK = ? AND DELETE_STATUS = 0`
	var count int
	err := m.DB.QueryRow(query, nik).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check NIK in employees: %v", err)
	}
	return count > 0, nil
}

// Create creates a new user
func (m *UserModel) Create(user *UserRequest) (*User, error) {
	// Check if NIK exists in employees table
	exists, err := m.CheckNIKInEmployees(user.NIK)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("NIK not found in employees database")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Insert user into database
	query := `INSERT INTO hs_wsb_user (NIK, NAME, EMAIL, PASSWORD, LEVEL) VALUES (?, ?, ?, ?, ?)`
	_, err = m.DB.Exec(query, user.NIK, user.Name, user.Email, string(hashedPassword), user.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Return created user
	return m.GetByNIK(user.NIK)
}

// GetByNIK gets user by NIK
func (m *UserModel) GetByNIK(nik string) (*User, error) {
	query := `SELECT NIK, NAME, EMAIL, PASSWORD, LEVEL, PROFILE_PICTURE, remember_token, created_at, updated_at FROM hs_wsb_user WHERE NIK = ?`
	row := m.DB.QueryRow(query, nik)

	var user User
	err := row.Scan(&user.NIK, &user.Name, &user.Email, &user.Password, &user.Level, &user.ProfilePicture, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// GetByEmail gets user by email
func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT NIK, NAME, EMAIL, PASSWORD, LEVEL, PROFILE_PICTURE, remember_token, created_at, updated_at FROM hs_wsb_user WHERE EMAIL = ?`
	row := m.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.NIK, &user.Name, &user.Email, &user.Password, &user.Level, &user.ProfilePicture, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// Update updates user information
func (m *UserModel) Update(nik string, user *UserRequest) (*User, error) {
    // First check if email exists for another user
    if user.Email != "" {
        existingUser, err := m.GetByEmail(user.Email)
        if err == nil && existingUser != nil && existingUser.NIK != nik {
            return nil, fmt.Errorf("email already registered")
        }
    }

    var query string
    var args []interface{}

    if user.Password != "" {
        // Hash new password if provided
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            return nil, fmt.Errorf("failed to hash password: %v", err)
        }
        query = `UPDATE hs_wsb_user SET NAME = ?, EMAIL = ?, PASSWORD = ?, LEVEL = ?, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
        args = []interface{}{user.Name, user.Email, string(hashedPassword), user.Level, nik}
    } else {
        // Update without password
        query = `UPDATE hs_wsb_user SET NAME = ?, EMAIL = ?, LEVEL = ?, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
        args = []interface{}{user.Name, user.Email, user.Level, nik}
    }

    result, err := m.DB.Exec(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to update user: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return nil, fmt.Errorf("no user was updated")
    }

    return m.GetByNIK(nik)
}

// UpdateWithProfilePicture updates user with profile picture
func (m *UserModel) UpdateWithProfilePicture(nik string, name, email, profilePicture string) (*User, error) {
	// Debug log
	fmt.Printf("UpdateWithProfilePicture called with: NIK=%s, name=%s, email=%s, profilePicture=%s\n", 
		nik, name, email, profilePicture)

	// Check if user exists first
	existingUser, err := m.GetByNIK(nik)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Debug log
	fmt.Printf("Existing user found: %+v\n", existingUser)

	// Check if email exists for another user
	if email != "" && email != existingUser.Email {
		existingUserByEmail, err := m.GetByEmail(email)
		if err == nil && existingUserByEmail != nil && existingUserByEmail.NIK != nik {
			return nil, fmt.Errorf("email already registered")
		}
	}

	// Check if any changes are needed
	needsUpdate := false
	if name != existingUser.Name || email != existingUser.Email {
		needsUpdate = true
	}

	// For profile picture, compare with existing value considering nil cases
	currentProfilePicture := ""
	if existingUser.ProfilePicture != nil {
		currentProfilePicture = *existingUser.ProfilePicture
	}
	if profilePicture != currentProfilePicture {
		needsUpdate = true
	}

	// Debug log update status
	fmt.Printf("Update needed: %v (name changed: %v, email changed: %v, profile picture changed: %v)\n",
		needsUpdate,
		name != existingUser.Name,
		email != existingUser.Email,
		profilePicture != currentProfilePicture)

	// If no changes needed, return existing user
	if !needsUpdate {
		return existingUser, nil
	}

	// Build query and args based on what's being updated
	var query string
	var args []interface{}

	if profilePicture != "" {
		// Update with new profile picture
		query = `UPDATE hs_wsb_user SET NAME = ?, EMAIL = ?, PROFILE_PICTURE = ?, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
		args = []interface{}{name, email, profilePicture, nik}
	} else {
		// Remove profile picture or update without profile picture
		query = `UPDATE hs_wsb_user SET NAME = ?, EMAIL = ?, PROFILE_PICTURE = NULL, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
		args = []interface{}{name, email, nik}
	}

	// Debug log
	fmt.Printf("Executing query: %s with args: %v\n", query, args)

	result, err := m.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	// Debug log
	fmt.Printf("Rows affected: %d\n", rowsAffected)

	// Return updated user
	return m.GetByNIK(nik)
}

// Delete deletes user by NIK
func (m *UserModel) Delete(nik string) error {
	query := `DELETE FROM hs_wsb_user WHERE NIK = ?`
	_, err := m.DB.Exec(query, nik)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

// CheckPassword checks if provided password matches user's password
func (m *UserModel) CheckPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// PaginatedUsers represents paginated response for users
type PaginatedUsers struct {
	Users       []*User `json:"users"`
	TotalItems  int     `json:"total_items"`
	TotalPages  int     `json:"total_pages"`
	CurrentPage int     `json:"current_page"`
	PageSize    int     `json:"page_size"`
}

// GetAll gets all users (for admin purposes)
func (m *UserModel) GetAll(page, pageSize int, search string) (*PaginatedUsers, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Base conditions for search
	searchCondition := ""
	searchArgs := []interface{}{}

	if search != "" {
		searchCondition = "WHERE (NIK LIKE ? OR NAME LIKE ? OR EMAIL LIKE ?)"
		searchPattern := "%" + search + "%"
		searchArgs = append(searchArgs, searchPattern, searchPattern, searchPattern)
	}

	// Get total count first with search condition
	countQuery := `SELECT COUNT(*) FROM hs_wsb_user ` + searchCondition
	var totalItems int
	err := m.DB.QueryRow(countQuery, searchArgs...).Scan(&totalItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}

	// Calculate total pages
	totalPages := (totalItems + pageSize - 1) / pageSize

	// Main query with pagination and search
	query := `
		SELECT NIK, NAME, EMAIL, LEVEL, PROFILE_PICTURE, remember_token, created_at, updated_at 
		FROM hs_wsb_user 
		` + searchCondition + `
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	// Add pagination parameters to search args
	searchArgs = append(searchArgs, pageSize, offset)

	rows, err := m.DB.Query(query, searchArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.NIK,
			&user.Name,
			&user.Email,
			&user.Level,
			&user.ProfilePicture,
			&user.RememberToken,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %v", err)
	}

	return &PaginatedUsers{
		Users:       users,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}

// UpdatePassword updates user's password
func (m *UserModel) UpdatePassword(nik string, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update password in database
	query := `UPDATE hs_wsb_user SET PASSWORD = ?, updated_at = CURRENT_TIMESTAMP WHERE NIK = ?`
	_, err = m.DB.Exec(query, string(hashedPassword), nik)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	createdAt := ""
	updatedAt := ""
	
	if u.CreatedAt != nil {
		createdAt = u.CreatedAt.Format("2006-01-02 15:04:05")
	}
	if u.UpdatedAt != nil {
		updatedAt = u.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	return &UserResponse{
		NIK:            u.NIK,
		Name:           u.Name,
		Email:          u.Email,
		Level:          u.Level,
		ProfilePicture: u.ProfilePicture,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

// GetAvailableEmployees gets list of employees that haven't registered yet
func (m *UserModel) GetAvailableEmployees() ([]*Employee, error) {
	query := `
		SELECT e.NIK, e.NAME as Name 
		FROM hs_hrd_employee e 
		LEFT JOIN hs_wsb_user u ON e.NIK = u.NIK 
		WHERE u.NIK IS NULL AND e.DELETE_STATUS = 0
		ORDER BY e.NIK
	`
	
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get available employees: %v", err)
	}
	defer rows.Close()

	var employees []*Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.NIK, &emp.Name); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %v", err)
		}
		employees = append(employees, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employees: %v", err)
	}

	return employees, nil
}