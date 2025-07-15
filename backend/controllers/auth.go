package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/username/backend/middleware"
	"github.com/username/backend/models"
)

// AuthController handles authentication operations
type AuthController struct {
	UserModel *models.UserModel
}

// NewAuthController creates new auth controller instance
func NewAuthController(userModel *models.UserModel) *AuthController {
	return &AuthController{
		UserModel: userModel,
	}
}

// Response represents API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents login response payload
type LoginResponse struct {
	NIK            string  `json:"nik"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Level          int     `json:"level"`
	ProfilePicture *string `json:"profile_picture"`
	Token          string  `json:"token"`
}

// Register handles user registration
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("Error decoding request body: %v\n", err)
		response := Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Printf("Received registration request: %+v\n", req)

	// Validate input
	if req.NIK == "" || req.Name == "" || req.Email == "" || req.Password == "" {
		fmt.Printf("Validation failed: NIK=%s, Name=%s, Email=%s, Password=<redacted>\n", 
			req.NIK, req.Name, req.Email)
		response := Response{
			Success: false,
			Message: "NIK, name, email, and password are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if email already exists
	existingUser, err := ac.UserModel.GetByEmail(req.Email)
	if err != nil && err.Error() != "user not found" {
		fmt.Printf("Error checking existing email: %v\n", err)
	}
	if existingUser != nil {
		fmt.Printf("Email already exists: %s\n", req.Email)
		response := Response{
			Success: false,
			Message: "Email already registered",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if NIK already exists
	existingUser, err = ac.UserModel.GetByNIK(req.NIK)
	if err != nil && err.Error() != "user not found" {
		fmt.Printf("Error checking existing NIK: %v\n", err)
	}
	if existingUser != nil {
		fmt.Printf("NIK already exists: %s\n", req.NIK)
		response := Response{
			Success: false,
			Message: "NIK already taken",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if NIK exists in employees table
	exists, err := ac.UserModel.CheckNIKInEmployees(req.NIK)
	if err != nil {
		fmt.Printf("Error checking NIK in employees: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to validate NIK: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	if !exists {
		fmt.Printf("NIK not found in employees: %s\n", req.NIK)
		response := Response{
			Success: false,
			Message: "NIK not found in employees database",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create user
	user, err := ac.UserModel.Create(&req)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to create user: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.NIK, user.Email)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		response := Response{
			Success: false,
			Message: "Failed to generate token",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Printf("Successfully registered user: %s\n", user.NIK)
	response := Response{
		Success: true,
		Message: "User registered successfully",
		Data: LoginResponse{
			NIK:            user.NIK,
			Name:           user.Name,
			Email:          user.Email,
			Level:          user.Level,
			ProfilePicture: user.ProfilePicture,
			Token:          token,
		},
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
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

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid JSON format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get user by email
	user, err := ac.UserModel.GetByEmail(req.Email)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid email or password",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check password
	if !ac.UserModel.CheckPassword(user, req.Password) {
		response := Response{
			Success: false,
			Message: "Invalid email or password",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.NIK, user.Email)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to generate token",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Login successful",
		Data: LoginResponse{
			NIK:            user.NIK,
			Name:           user.Name,
			Email:          user.Email,
			Level:          user.Level,
			ProfilePicture: user.ProfilePicture,
			Token:          token,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns current authenticated user info
func (ac *AuthController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
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

	// Get user NIK from middleware
	userNIK := r.Header.Get("X-User-ID")
	if userNIK == "" {
		response := Response{
			Success: false,
			Message: "User NIK not found in request",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get user from database
	user, err := ac.UserModel.GetByNIK(userNIK)
	if err != nil {
		response := Response{
			Success: false,
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user.ToResponse(),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// TestConnection tests backend connection
func (ac *AuthController) TestConnection(w http.ResponseWriter, r *http.Request) {
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

	// Check if authenticated
	userEmail := r.Header.Get("X-User-Email")
	message := "Backend connection successful"
	if userEmail != "" {
		message = fmt.Sprintf("Backend connection successful for user: %s", userEmail)
	}

	response := Response{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"timestamp":     fmt.Sprintf("%v", strings.Replace(fmt.Sprintf("%v", r.Header.Get("Date")), " ", "T", 1)),
			"authenticated": userEmail != "",
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetAvailableEmployees returns list of employees that haven't registered yet
func (ac *AuthController) GetAvailableEmployees(w http.ResponseWriter, r *http.Request) {
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

	// Get employees that haven't registered
	employees, err := ac.UserModel.GetAvailableEmployees()
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to get available employees",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Available employees retrieved successfully",
		Data:    employees,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}