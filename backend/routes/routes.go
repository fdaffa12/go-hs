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
	BuyerController     *controllers.BuyerController
	StyleController     *controllers.StyleController
	LineScheduleController *controllers.LineScheduleController
	HolidayController *controllers.HolidayController
	SewNumProcessController *controllers.SewNumProcessController
}

// NewRouter creates new router instance
func NewRouter(
	authController *controllers.AuthController,
	userController *controllers.UserController,
	departmentController *controllers.DepartmentController,
	employeeController *controllers.EmployeeController,
	buyerController *controllers.BuyerController,
	styleController *controllers.StyleController,
	lineScheduleController *controllers.LineScheduleController,
	holidayController *controllers.HolidayController,
	sewNumProcessController *controllers.SewNumProcessController,
) *Router {
	return &Router{
		AuthController:         authController,
		UserController:        userController,
		DepartmentController:  departmentController,
		EmployeeController:    employeeController,
		BuyerController:      buyerController,
		StyleController:      styleController,
		LineScheduleController: lineScheduleController,
		HolidayController:    holidayController,
		SewNumProcessController: sewNumProcessController,
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

	// Buyer management routes (authentication required)
	mux.HandleFunc("/api/buyers", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleBuyerRoutes))))
	mux.HandleFunc("/api/buyers/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleBuyerRoutes))))

	// Style management routes (authentication required)
	mux.HandleFunc("/api/styles", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleStyleRoutes))))
	mux.HandleFunc("/api/styles/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleStyleRoutes))))

	// Line Schedule management routes (authentication required)
	mux.HandleFunc("/api/line-schedules", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleLineScheduleRoutes))))
	mux.HandleFunc("/api/line-schedules/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleLineScheduleRoutes))))

	// Holiday management routes (authentication required)
	mux.HandleFunc("/api/holidays", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleHolidayRoutes))))
	mux.HandleFunc("/api/holidays/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleHolidayRoutes))))

	// Sew Numbering Process routes (authentication required)
	mux.HandleFunc("/api/sew-num-process", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleSewNumProcessRoutes))))
	mux.HandleFunc("/api/sew-num-process/", middleware.CORSMiddleware(middleware.LoggingMiddleware(middleware.AuthMiddleware(router.handleSewNumProcessRoutes))))

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

// handleBuyerRoutes handles dynamic buyer routes based on HTTP method
func (router *Router) handleBuyerRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/buyers")

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete/") {
		router.BuyerController.HardDeleteBuyer(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate/") {
		router.BuyerController.ActivateBuyer(w, r)
		return
	}

	// Route for getting all buyers or creating new buyer
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.BuyerController.GetAllBuyers(w, r)
		case "POST":
			router.BuyerController.CreateBuyer(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// Routes for specific buyer operations (update, delete)
	switch r.Method {
	case "PUT":
		router.BuyerController.UpdateBuyer(w, r)
	case "DELETE":
		router.BuyerController.DeleteBuyer(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
	}
}

// handleStyleRoutes handles dynamic style routes based on HTTP method
func (router *Router) handleStyleRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/styles")

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete/") {
		router.StyleController.HardDeleteStyle(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate/") {
		router.StyleController.ActivateStyle(w, r)
		return
	}

	// Route for getting all styles or creating new style
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.StyleController.GetAllStyles(w, r)
		case "POST":
			router.StyleController.CreateStyle(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// Routes for specific style operations (update, delete)
	switch r.Method {
	case "PUT":
		router.StyleController.UpdateStyle(w, r)
	case "DELETE":
		router.StyleController.DeleteStyle(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
	}
}

// handleLineScheduleRoutes handles dynamic line schedule routes based on HTTP method
func (router *Router) handleLineScheduleRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/line-schedules")

	// Handle calculate working days
	if strings.HasPrefix(path, "/calculate-working-days") {
		router.LineScheduleController.CalculateWorkingDays(w, r)
		return
	}

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete/") {
		router.LineScheduleController.HardDeleteLineSchedule(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate/") {
		router.LineScheduleController.ActivateLineSchedule(w, r)
		return
	}

	// Route for getting all line schedules or creating new line schedule
	// Also handle update and delete with query parameters
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.LineScheduleController.GetAllLineSchedules(w, r)
		case "POST":
			router.LineScheduleController.CreateLineSchedule(w, r)
		case "PUT":
			router.LineScheduleController.UpdateLineSchedule(w, r)
		case "DELETE":
			router.LineScheduleController.DeleteLineSchedule(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// If we get here, the path is not recognized
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"success": false, "message": "Route not found"}`))
}

// handleHolidayRoutes handles dynamic holiday routes based on HTTP method
func (router *Router) handleHolidayRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/holidays")

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete") {
		router.HolidayController.HardDeleteHoliday(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate") {
		router.HolidayController.ActivateHoliday(w, r)
		return
	}

	// Route for getting all holidays or creating new holiday
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.HolidayController.GetAllHolidays(w, r)
		case "POST":
			router.HolidayController.CreateHoliday(w, r)
		case "PUT":
			router.HolidayController.UpdateHoliday(w, r)
		case "DELETE":
			router.HolidayController.DeleteHoliday(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// If we get here, the path is not recognized
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"success": false, "message": "Route not found"}`))
}
// handleSewNumProcessRoutes handles dynamic sew numbering process routes based on HTTP method
func (router *Router) handleSewNumProcessRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/sew-num-process")

	// Handle hard delete
	if strings.HasPrefix(path, "/hard-delete/") {
		router.SewNumProcessController.HardDelete(w, r)
		return
	}

	// Handle activate
	if strings.HasPrefix(path, "/activate/") {
		router.SewNumProcessController.Activate(w, r)
		return
	}

	// Route for getting all processes or bulk operations
	if path == "" || path == "/" {
		switch r.Method {
		case "GET":
			router.SewNumProcessController.GetAllProcesses(w, r)
		case "POST":
			if strings.Contains(r.URL.Path, "bulk-delete") {
				router.SewNumProcessController.BulkDelete(w, r)
			} else {
				router.SewNumProcessController.BulkSave(w, r)
			}
		case "DELETE":
			router.SewNumProcessController.Delete(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"success": false, "message": "Method not allowed"}`))
		}
		return
	}

	// If we get here, the path is not recognized
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"success": false, "message": "Route not found"}`))
}
