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
	Name  string `json:"name"`
	Email string `json:"email"`
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

	// Get pagination parameters from query string
	page := 1
	pageSize := 10 // default page size

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	// Get search parameter
	search := r.URL.Query().Get("search")

	users, err := uc.UserModel.GetAll(page, pageSize, search)
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
	userResponses := make([]*models.UserResponse, len(users.Users))
	for i, user := range users.Users {
		userResponses[i] = user.ToResponse()
	}

	response := Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data: map[string]interface{}{
			"users":        userResponses,
			"total_items":  users.TotalItems,
			"total_pages":  users.TotalPages,
			"current_page": users.CurrentPage,
			"page_size":    users.PageSize,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetUserByNIK handles getting user by NIK
func (uc *UserController) GetUserByNIK(w http.ResponseWriter, r *http.Request) {
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

	// Extract user NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if path == "" {
		response := Response{
			Success: false,
			Message: "Invalid NIK",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := uc.UserModel.GetByNIK(path)
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
	if req.NIK == "" || req.Name == "" || req.Email == "" || req.Password == "" {
		response := Response{
			Success: false,
			Message: "NIK, name, email, and password are required",
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

	// Check if NIK already exists
	existingUser, _ = uc.UserModel.GetByNIK(req.NIK)
	if existingUser != nil {
		response := Response{
			Success: false,
			Message: "NIK already registered",
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

	// Extract user NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if path == "" {
		response := Response{
			Success: false,
			Message: "Invalid NIK",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if user exists
	_, err := uc.UserModel.GetByNIK(path)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
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
	if req.Name == "" || req.Email == "" {
		response := Response{
			Success: false,
			Message: "Name and email are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if email already exists (excluding current user)
	existingUser, _ := uc.UserModel.GetByEmail(req.Email)
	if existingUser != nil && existingUser.NIK != path {
		response := Response{
			Success: false,
			Message: "Email already registered",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update user
	userReq := &models.UserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Level:    req.Level,
		Password: req.Password, // Will be handled by model if provided
	}

	updatedUser, err := uc.UserModel.Update(path, userReq)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to update user: %v", err),
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

	// Extract user NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if path == "" {
		response := Response{
			Success: false,
			Message: "Invalid NIK",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if user exists
	_, err := uc.UserModel.GetByNIK(path)
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
	err = uc.UserModel.Delete(path)
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

	// Extract user NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	path = strings.TrimSuffix(path, "/password")
	if path == "" {
		response := Response{
			Success: false,
			Message: "Invalid NIK",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get current user
	currentUser, err := uc.UserModel.GetByNIK(path)
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
	if !uc.UserModel.CheckPassword(currentUser, req.CurrentPassword) {
		response := Response{
			Success: false,
			Message: "Current password is incorrect",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update password
	err = uc.UserModel.UpdatePassword(path, req.NewPassword)
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

	// Debug log
	fmt.Printf("UpdateUserProfile called with method: %s, URL: %s\n", r.Method, r.URL.Path)

	if r.Method != "PUT" {
		response := Response{
			Success: false,
			Message: "Method not allowed",
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Extract user NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	path = strings.TrimSuffix(path, "/profile")
	
	// Debug log
	fmt.Printf("Extracted NIK from path: %s\n", path)

	if path == "" {
		response := Response{
			Success: false,
			Message: "Invalid NIK",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if user exists
	currentUser, err := uc.UserModel.GetByNIK(path)
	if err != nil {
		// Debug log
		fmt.Printf("Error getting user by NIK: %v\n", err)
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Debug log
	fmt.Printf("Current user found: %+v\n", currentUser)

	// Parse multipart form with larger size limit
	err = r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		// Debug log
		fmt.Printf("Error parsing multipart form: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to parse form data: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get form values
	name := r.FormValue("name")
	email := r.FormValue("email")
	removeProfilePicture := r.FormValue("remove_profile_picture") == "true"

	// Debug log
	fmt.Printf("Form values - name: %s, email: %s, removeProfilePicture: %v\n", 
		name, email, removeProfilePicture)

	// Validate input
	if name == "" || email == "" {
		response := Response{
			Success: false,
			Message: "Name and email are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var profilePictureURL string

	// Handle profile picture removal
	if removeProfilePicture {
		// Debug log
		fmt.Printf("Removing profile picture for user: %s\n", path)
		
		// If user has an existing profile picture, delete the file
		if currentUser.ProfilePicture != nil && *currentUser.ProfilePicture != "" {
			oldFilePath := "." + *currentUser.ProfilePicture
			// Debug log
			fmt.Printf("Deleting old profile picture: %s\n", oldFilePath)
			os.Remove(oldFilePath)
		}
		profilePictureURL = "" // This will trigger profile picture removal in the model
	} else {
		// Handle new file upload if provided
		file, header, err := r.FormFile("profile_picture")
		if err == nil {
			defer file.Close()

			// Debug log
			fmt.Printf("New file upload detected - filename: %s, size: %d, content-type: %s\n",
				header.Filename, header.Size, header.Header.Get("Content-Type"))

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

			// Delete existing profile picture if exists
			if currentUser.ProfilePicture != nil && *currentUser.ProfilePicture != "" {
				oldFilePath := "." + *currentUser.ProfilePicture
				os.Remove(oldFilePath)
			}

			// Generate unique filename
			ext := filepath.Ext(header.Filename)
			filename := fmt.Sprintf("%s_%d%s", path, time.Now().Unix(), ext)
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
			
			// Debug log
			fmt.Printf("File saved successfully at: %s\n", filePath)
			fmt.Printf("Profile picture URL set to: %s\n", profilePictureURL)
		} else {
			// Debug log
			fmt.Printf("No new file uploaded: %v\n", err)
			
			// If no new file and not removing, keep existing profile picture
			if currentUser.ProfilePicture != nil {
				profilePictureURL = *currentUser.ProfilePicture
			}
		}
	}

	// Debug log before database update
	fmt.Printf("Updating user profile - NIK: %s, name: %s, email: %s, profilePictureURL: %s\n",
		path, name, email, profilePictureURL)

	// Update user in database
	updatedUser, err := uc.UserModel.UpdateWithProfilePicture(path, name, email, profilePictureURL)
	if err != nil {
		// Debug log
		fmt.Printf("Error updating user profile: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to update user profile: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Debug log success
	fmt.Printf("Successfully updated user profile: %+v\n", updatedUser)

	response := Response{
		Success: true,
		Message: "Profile updated successfully",
		Data:    updatedUser.ToResponse(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}