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
	EmployeeController   *controllers.EmployeeController
}

// NewRouter creates new router instance
func NewRouter(authController *controllers.AuthController, userController *controllers.UserController, departmentController *controllers.DepartmentController, employeeController *controllers.EmployeeController) *Router {
	return &Router{
		AuthController:       authController,
		UserController:      userController,
		DepartmentController: departmentController,
		EmployeeController:   employeeController,
	}
}

// SetupRoutes configures all application routes
func (router *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("/api/register", middleware.CORSMiddleware(middleware.LoggingMiddleware(router.AuthController.Register)))
	mux.HandleFunc("/api/login", middleware.CORSMiddleware(middleware.LoggingMiddleware(router.AuthController.Login)))
	mux.HandleFunc("/api/available-employees", middleware.CORSMiddleware(middleware.LoggingMiddleware(router.AuthController.GetAvailableEmployees)))

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

	// Employee management routes (authentication required)
	mux.HandleFunc("/api/employees", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleEmployeeRoutes))))
	mux.HandleFunc("/api/employees/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleEmployeeRoutes))))

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
		w.Write([]byte(`{"success": false, "message": "NIK is required"}`))
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
		router.UserController.GetUserByNIK(w, r)
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

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete/") {
		router.DepartmentController.HardDeleteDepartment(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate/") {
		router.DepartmentController.ActivateDepartment(w, r)
		return
	}

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

// handleEmployeeRoutes handles dynamic employee routes based on HTTP method
func (router *Router) handleEmployeeRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/employees")

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete/") {
		router.EmployeeController.HardDeleteEmployee(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate/") {
		router.EmployeeController.ActivateEmployee(w, r)
		return
	}

	// Route for getting all employees or creating new employee
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.EmployeeController.GetAllEmployees(w, r)
		case "POST":
			router.EmployeeController.CreateEmployee(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// Routes for specific employee operations (update, delete)
	switch r.Method {
	case "PUT":
		router.EmployeeController.UpdateEmployee(w, r)
	case "DELETE":
		router.EmployeeController.DeleteEmployee(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
	}
}