package routes

import (
	"net/http"
	"strings"

	"github.com/username/backend/controllers"
	"github.com/username/backend/middleware"
)

// Router handles all application routes
type Router struct {
	AuthController       *controllers.AuthController
	UserController      *controllers.UserController
	DepartmentController *controllers.DepartmentController
}

// NewRouter creates new router instance
func NewRouter(authController *controllers.AuthController, userController *controllers.UserController, departmentController *controllers.DepartmentController) *Router {
	return &Router{
		AuthController:       authController,
		UserController:      userController,
		DepartmentController: departmentController,
	}
}

// SetupRoutes configures all application routes
func (router *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/api/register", middleware.CORSMiddleware(middleware.LoggingMiddleware(router.AuthController.Register)))
	mux.HandleFunc("/api/login", middleware.CORSMiddleware(middleware.LoggingMiddleware(router.AuthController.Login)))

	// Protected routes (authentication required)
	mux.HandleFunc("/api/user", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.AuthController.GetCurrentUser))))
	mux.HandleFunc("/api/test", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.AuthController.TestConnection))))

	// User management routes (authentication required)
	mux.HandleFunc("/api/users", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.UserController.GetAllUsers))))
	mux.HandleFunc("/api/users/create", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.UserController.CreateUser))))
	mux.HandleFunc("/api/users/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleUserRoutes))))

	// Department management routes (authentication required)
	mux.HandleFunc("/api/departments", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleDepartmentRoutes))))
	mux.HandleFunc("/api/departments/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleDepartmentRoutes))))

	// Public test route (for testing backend connection without auth)
	mux.HandleFunc("/api/hello", middleware.CORSMiddleware(middleware.LoggingMiddleware(router.AuthController.TestConnection)))

	// Static file serving for uploaded images
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	return mux
}

// handleUserRoutes handles dynamic user routes based on HTTP method
func (router *Router) handleUserRoutes(w http.ResponseWriter, r *http.Request) {
	// Extract path after /api/users/
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	
	// If path is empty, return error
	if path == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success": false, "message": "User ID is required"}`))
		return
	}

	// Check if this is a password change request
	if strings.HasSuffix(path, "/password") {
		if r.Method == "PUT" {
			router.UserController.ChangePassword(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// Check if this is a profile update request with image upload
	if strings.HasSuffix(path, "/profile") {
		if r.Method == "PUT" {
			router.UserController.UpdateUserProfile(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// Route based on HTTP method
	switch r.Method {
	case "GET":
		router.UserController.GetUserByID(w, r)
	case "PUT":
		router.UserController.UpdateUser(w, r)
	case "DELETE":
		router.UserController.DeleteUser(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
	}
}

// handleDepartmentRoutes handles dynamic department routes based on HTTP method
func (router *Router) handleDepartmentRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/departments")

	// Route for getting all departments or creating new department
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.DepartmentController.GetAllDepartments(w, r)
		case "POST":
			router.DepartmentController.CreateDepartment(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// Routes for specific department operations (update, delete)
	switch r.Method {
	case "PUT":
		router.DepartmentController.UpdateDepartment(w, r)
	case "DELETE":
		router.DepartmentController.DeleteDepartment(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
	}
}