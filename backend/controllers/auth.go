package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// RegisterRequest represents register request payload
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents login response payload
type LoginResponse struct {
	Token string                `json:"token"`
	User  *models.UserResponse `json:"user"`
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

	var req RegisterRequest
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
	existingUser, _ := ac.UserModel.GetByEmail(req.Email)
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
	existingUser, _ = ac.UserModel.GetByUsername(req.Username)
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
	userReq := &models.UserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := ac.UserModel.Create(userReq)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to create user: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Email)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to generate token",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := Response{
		Success: true,
		Message: "User registered successfully",
		Data: LoginResponse{
			Token: token,
			User:  user.ToResponse(),
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

	// Validate input
	if req.Email == "" || req.Password == "" {
		response := Response{
			Success: false,
			Message: "Email and password are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Find user by email
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
	token, err := middleware.GenerateJWT(user.ID, user.Email)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to generate token",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := Response{
		Success: true,
		Message: "Login successful",
		Data: LoginResponse{
			Token: token,
			User:  user.ToResponse(),
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

	// Get user ID from middleware
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		response := Response{
			Success: false,
			Message: "User ID not found in request",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get user from database
	user, err := ac.UserModel.GetByID(userID)
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