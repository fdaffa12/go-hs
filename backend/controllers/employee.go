package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/username/backend/models"
)

// EmployeeController handles employee management operations
type EmployeeController struct {
	EmployeeModel *models.EmployeeModel
}

// NewEmployeeController creates new employee controller instance
func NewEmployeeController(employeeModel *models.EmployeeModel) *EmployeeController {
	return &EmployeeController{
		EmployeeModel: employeeModel,
	}
}

// GetAllEmployees handles getting all employees
func (ec *EmployeeController) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
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

	result, err := ec.EmployeeModel.GetAll(page, pageSize, search)
	if err != nil {
		fmt.Printf("Error getting employees: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal mengambil data karyawan: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Data karyawan berhasil diambil",
		Data:    result,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateEmployee handles creating a new employee
func (ec *EmployeeController) CreateEmployee(w http.ResponseWriter, r *http.Request) {
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

	var req models.EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		response := Response{
			Success: false,
			Message: "Format JSON tidak valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if req.NIK == "" || req.Name == "" || req.DeptShortName == "" || req.EnteranceDate == "" || req.Title == "" {
		response := Response{
			Success: false,
			Message: "Semua field wajib diisi kecuali RFID ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create employee
	emp, err := ec.EmployeeModel.Create(&req)
	if err != nil {
		fmt.Printf("Error creating employee: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal membuat karyawan: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Karyawan berhasil dibuat",
		Data:    emp,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateEmployee handles updating employee information
func (ec *EmployeeController) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
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

	// Extract NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/employees/")
	nik := path

	var req models.EmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Format JSON tidak valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if req.Name == "" || req.DeptShortName == "" || req.EnteranceDate == "" || req.Title == "" {
		response := Response{
			Success: false,
			Message: "Semua field wajib diisi kecuali RFID ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update employee
	req.NIK = nik // Ensure we use the path parameter
	updatedEmp, err := ec.EmployeeModel.Update(nik, &req)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal memperbarui karyawan: " + err.Error(),
		}
		if err.Error() == "employee not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Karyawan berhasil diperbarui",
		Data:    updatedEmp,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteEmployee handles soft deleting an employee
func (ec *EmployeeController) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
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

	// Extract NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/employees/")
	nik := path

	// Delete employee
	err := ec.EmployeeModel.Delete(nik)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus karyawan: " + err.Error(),
		}
		if err.Error() == "employee not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Karyawan berhasil dinonaktifkan",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDeleteEmployee handles permanently deleting an employee
func (ec *EmployeeController) HardDeleteEmployee(w http.ResponseWriter, r *http.Request) {
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

	// Extract NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/employees/hard-delete/")
	nik := path

	// Delete employee permanently
	err := ec.EmployeeModel.HardDelete(nik)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus karyawan secara permanen: " + err.Error(),
		}
		if err.Error() == "employee not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Karyawan berhasil dihapus secara permanen",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ActivateEmployee handles reactivating a soft-deleted employee
func (ec *EmployeeController) ActivateEmployee(w http.ResponseWriter, r *http.Request) {
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

	// Extract NIK from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/employees/activate/")
	nik := path

	// Activate employee
	err := ec.EmployeeModel.Activate(nik)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal mengaktifkan karyawan: " + err.Error(),
		}
		if err.Error() == "employee not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Karyawan berhasil diaktifkan",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 