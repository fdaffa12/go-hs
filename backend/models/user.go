package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents user model
type User struct {
	ID             int       `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	Password       string    `json:"-" db:"password"` // "-" means don't include in JSON
	ProfilePicture *string   `json:"profile_picture" db:"profile_picture"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// UserRequest represents request payload for user operations
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse represents response payload for user operations
type UserResponse struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
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

// Create creates a new user
func (m *UserModel) Create(user *UserRequest) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Insert user into database
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := m.DB.Exec(query, user.Username, user.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Get inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}

	// Return created user
	return m.GetByID(int(id))
}

// GetByID gets user by ID
func (m *UserModel) GetByID(id int) (*User, error) {
	query := `SELECT id, username, email, password, profile_picture, created_at, updated_at FROM users WHERE id = ?`
	row := m.DB.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)
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
	query := `SELECT id, username, email, password, profile_picture, created_at, updated_at FROM users WHERE email = ?`
	row := m.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// GetByUsername gets user by username
func (m *UserModel) GetByUsername(username string) (*User, error) {
	query := `SELECT id, username, email, password, profile_picture, created_at, updated_at FROM users WHERE username = ?`
	row := m.DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return &user, nil
}

// Update updates user information
func (m *UserModel) Update(id int, user *UserRequest) (*User, error) {
	query := `UPDATE users SET username = ?, email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.DB.Exec(query, user.Username, user.Email, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return m.GetByID(id)
}

// UpdateWithProfilePicture updates user with profile picture
func (m *UserModel) UpdateWithProfilePicture(id int, username, email, profilePicture string) (*User, error) {
	query := `UPDATE users SET username = ?, email = ?, profile_picture = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.DB.Exec(query, username, email, profilePicture, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user with profile picture: %v", err)
	}

	return m.GetByID(id)
}

// Delete deletes user by ID
func (m *UserModel) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := m.DB.Exec(query, id)
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

// GetAll gets all users (for admin purposes)
func (m *UserModel) GetAll() ([]*User, error) {
	query := `SELECT id, username, email, profile_picture, created_at, updated_at FROM users ORDER BY created_at DESC`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate users: %v", err)
	}

	return users, nil
}

// UpdatePassword updates user's password
func (m *UserModel) UpdatePassword(id int, newPassword string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Update password in database
	query := `UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = m.DB.Exec(query, string(hashedPassword), id)
	if err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return nil
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		ProfilePicture: u.ProfilePicture,
		CreatedAt:      u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      u.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}