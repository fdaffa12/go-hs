package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/username/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles user management operations
type UserController struct {
	UserModel *models.UserModel
}

// NewUserController creates new user controller instance
func NewUserController(userModel *models.UserModel) *UserController {
	return &UserController{
		UserModel: userModel,
	}
}

// UpdateUserRequest represents update user request payload
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ChangePasswordRequest represents change password request payload
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// GetAllUsers handles getting all users (admin only)
func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	users, err := uc.UserModel.GetAll()
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to get users",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response format
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	response := Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    userResponses,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserByID handles getting user by ID
func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userID, err := strconv.Atoi(path)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := uc.UserModel.GetByID(userID)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user.ToResponse(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateUser handles creating a new user
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		response := Response{
			Success: false,
			Message: "Username, email, and password are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if email already exists
	existingUser, _ := uc.UserModel.GetByEmail(req.Email)
	if existingUser != nil {
		response := Response{
			Success: false,
			Message: "Email already registered",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if username already exists
	existingUser, _ = uc.UserModel.GetByUsername(req.Username)
	if existingUser != nil {
		response := Response{
			Success: false,
			Message: "Username already taken",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create user
	user, err := uc.UserModel.Create(&req)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to create user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "User created successfully",
		Data:    user.ToResponse(),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateUser handles updating user information
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userID, err := strconv.Atoi(path)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if user exists
	_, err = uc.UserModel.GetByID(userID)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" {
		response := Response{
			Success: false,
			Message: "Username and email are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if email already exists (excluding current user)
	existingUser, _ := uc.UserModel.GetByEmail(req.Email)
	if existingUser != nil && existingUser.ID != userID {
		response := Response{
			Success: false,
			Message: "Email already registered",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if username already exists (excluding current user)
	existingUser, _ = uc.UserModel.GetByUsername(req.Username)
	if existingUser != nil && existingUser.ID != userID {
		response := Response{
			Success: false,
			Message: "Username already taken",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update user
	userReq := &models.UserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: "", // Password not updated in this endpoint
	}

	updatedUser, err := uc.UserModel.Update(userID, userReq)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to update user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "User updated successfully",
		Data:    updatedUser.ToResponse(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteUser handles deleting a user
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	userID, err := strconv.Atoi(path)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if user exists
	_, err = uc.UserModel.GetByID(userID)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete user
	err = uc.UserModel.Delete(userID)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to delete user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "User deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ChangePassword handles changing user password
func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	path = strings.TrimSuffix(path, "/password")
	userID, err := strconv.Atoi(path)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get current user
	currentUser, err := uc.UserModel.GetByID(userID)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		response := Response{
			Success: false,
			Message: "Current password and new password are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(req.NewPassword) < 6 {
		response := Response{
			Success: false,
			Message: "New password must be at least 6 characters",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.CurrentPassword))
	if err != nil {
		response := Response{
			Success: false,
			Message: "Current password is incorrect",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update password
	err = uc.UserModel.UpdatePassword(userID, req.NewPassword)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to update password",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := Response{
		Success: true,
		Message: "Password changed successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateUserProfile handles updating user profile with image upload
func (uc *UserController) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "PUT" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract user ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	path = strings.TrimSuffix(path, "/profile")
	userID, err := strconv.Atoi(path)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if user exists
	_, err = uc.UserModel.GetByID(userID)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to parse form data",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get form values
	username := r.FormValue("username")
	email := r.FormValue("email")
	removeProfilePicture := r.FormValue("remove_profile_picture")

	// Validate input
	if username == "" || email == "" {
		response := Response{
			Success: false,
			Message: "Username and email are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle profile picture removal
	if removeProfilePicture == "true" {
		// Get current user to get existing profile picture path
		currentUser, err := uc.UserModel.GetByID(userID)
		if err == nil && currentUser.ProfilePicture != nil && *currentUser.ProfilePicture != "" {
			// Delete existing file
			oldFilePath := "." + *currentUser.ProfilePicture
			os.Remove(oldFilePath)
		}

		// Update user with empty profile picture
		updatedUser, err := uc.UserModel.UpdateWithProfilePicture(userID, username, email, "")
		if err != nil {
			response := Response{
				Success: false,
				Message: "Failed to remove profile picture",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := Response{
			Success: true,
			Message: "Profile picture removed successfully",
			Data:    updatedUser.ToResponse(),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle file upload
	var profilePictureURL string
	file, header, err := r.FormFile("profile_picture")
	if err == nil {
		defer file.Close()

		// Validate file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
		}

		contentType := header.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			response := Response{
				Success: false,
				Message: "Invalid file type. Only JPEG, PNG, and GIF are allowed",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Validate file size (2MB)
		if header.Size > 2*1024*1024 {
			response := Response{
				Success: false,
				Message: "File size too large. Maximum 2MB allowed",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Create uploads directory if it doesn't exist
		uploadsDir := "uploads/profile_pictures"
		err = os.MkdirAll(uploadsDir, 0755)
		if err != nil {
			response := Response{
				Success: false,
				Message: "Failed to create upload directory",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
		filePath := filepath.Join(uploadsDir, filename)

		// Create the file
		dst, err := os.Create(filePath)
		if err != nil {
			response := Response{
				Success: false,
				Message: "Failed to create file",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		defer dst.Close()

		// Copy file content
		_, err = io.Copy(dst, file)
		if err != nil {
			response := Response{
				Success: false,
				Message: "Failed to save file",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Set profile picture URL
		profilePictureURL = fmt.Sprintf("/uploads/profile_pictures/%s", filename)
	}

	// Update user in database
	var updatedUser *models.User
	if profilePictureURL != "" {
		updatedUser, err = uc.UserModel.UpdateWithProfilePicture(userID, username, email, profilePictureURL)
	} else {
		// Update without changing profile picture
		userReq := &models.UserRequest{
			Username: username,
			Email:    email,
			Password: "", // Password not updated in this endpoint
		}
		updatedUser, err = uc.UserModel.Update(userID, userReq)
	}

	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to update user profile",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Profile updated successfully",
		Data:    updatedUser.ToResponse(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}