package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/username/backend/models"
)

// DepartmentController handles department management operations
type DepartmentController struct {
	DepartmentModel *models.DepartmentModel
}

// NewDepartmentController creates new department controller instance
func NewDepartmentController(departmentModel *models.DepartmentModel) *DepartmentController {
	return &DepartmentController{
		DepartmentModel: departmentModel,
	}
}

// GetAllDepartments handles getting all departments
func (dc *DepartmentController) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
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

	departments, err := dc.DepartmentModel.GetAll(page, pageSize, search)
	if err != nil {
		fmt.Printf("Error getting departments: %v\n", err) // Add logging
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal mengambil data departemen: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Departemen berhasil diambil",
		Data:    departments,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateDepartment handles creating a new department
func (dc *DepartmentController) CreateDepartment(w http.ResponseWriter, r *http.Request) {
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

	var req models.DepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request: %v\n", err) // Add logging
		response := Response{
			Success: false,
			Message: "Format JSON tidak valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate input
	if req.ShortName == "" || req.LongName == "" {
		response := Response{
			Success: false,
			Message: "Nama singkat dan nama panjang departemen wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if department already exists
	existingDept, _ := dc.DepartmentModel.GetByShortName(req.ShortName)
	if existingDept != nil {
		response := Response{
			Success: false,
			Message: "Kode departemen sudah digunakan",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create department
	dept, err := dc.DepartmentModel.Create(&req)
	if err != nil {
		fmt.Printf("Error creating department: %v\n", err) // Add logging
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal membuat departemen: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Departemen berhasil dibuat",
		Data:    dept,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateDepartment handles updating department information
func (dc *DepartmentController) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
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

	// Extract department short name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/departments/")
	shortName := path

	var req models.DepartmentRequest
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
	if req.LongName == "" {
		response := Response{
			Success: false,
			Message: "Nama panjang departemen wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update department
	req.ShortName = shortName // Ensure we use the path parameter
	updatedDept, err := dc.DepartmentModel.Update(shortName, &req)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal memperbarui departemen: " + err.Error(),
		}
		if err.Error() == "department not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Departemen berhasil diperbarui",
		Data:    updatedDept,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteDepartment handles deleting a department
func (dc *DepartmentController) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
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

	// Extract department short name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/departments/")
	shortName := path

	// Delete department
	err := dc.DepartmentModel.Delete(shortName)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus departemen: " + err.Error(),
		}
		if err.Error() == "department not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Departemen berhasil dihapus",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDeleteDepartment handles permanently deleting a department
func (dc *DepartmentController) HardDeleteDepartment(w http.ResponseWriter, r *http.Request) {
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

    // Extract department short name from URL path
    path := strings.TrimPrefix(r.URL.Path, "/api/departments/hard-delete/")
    shortName := path

    // Delete department permanently
    err := dc.DepartmentModel.HardDelete(shortName)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Gagal menghapus departemen secara permanen: " + err.Error(),
        }
        if err.Error() == "department not found" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Departemen berhasil dihapus secara permanen",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// ActivateDepartment handles reactivating a soft-deleted department
func (dc *DepartmentController) ActivateDepartment(w http.ResponseWriter, r *http.Request) {
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

    // Extract department short name from URL path
    path := strings.TrimPrefix(r.URL.Path, "/api/departments/activate/")
    shortName := path

    // Activate department
    err := dc.DepartmentModel.Activate(shortName)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Gagal mengaktifkan departemen: " + err.Error(),
        }
        if err.Error() == "department not found" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Departemen berhasil diaktifkan",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
} 