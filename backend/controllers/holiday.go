package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/username/backend/models"
)

// HolidayController handles holiday management operations
type HolidayController struct {
	HolidayModel *models.HolidayModel
}

// NewHolidayController creates new holiday controller instance
func NewHolidayController(holidayModel *models.HolidayModel) *HolidayController {
	return &HolidayController{
		HolidayModel: holidayModel,
	}
}

// GetAllHolidays handles getting all holidays
func (hc *HolidayController) GetAllHolidays(w http.ResponseWriter, r *http.Request) {
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

	holidays, err := hc.HolidayModel.GetAll(page, pageSize, search)
	if err != nil {
		fmt.Printf("Error getting holidays: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal mengambil data holiday: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Holiday berhasil diambil",
		Data:    holidays,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateHoliday handles creating a new holiday
func (hc *HolidayController) CreateHoliday(w http.ResponseWriter, r *http.Request) {
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

	var req models.HolidayRequest
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
	if req.HolidayDate == "" || req.HolidayName == "" || req.HolidayType == "" {
		response := Response{
			Success: false,
			Message: "Tanggal, nama, dan tipe holiday wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create holiday
	holiday, err := hc.HolidayModel.Create(&req)
	if err != nil {
		fmt.Printf("Error creating holiday: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal membuat holiday: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Holiday berhasil dibuat",
		Data:    holiday,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateHoliday handles updating holiday information
func (hc *HolidayController) UpdateHoliday(w http.ResponseWriter, r *http.Request) {
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

	// Extract holiday ID from URL path
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid holiday ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req models.HolidayRequest
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
	if req.HolidayDate == "" || req.HolidayName == "" || req.HolidayType == "" {
		response := Response{
			Success: false,
			Message: "Tanggal, nama, dan tipe holiday wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update holiday
	updatedHoliday, err := hc.HolidayModel.Update(id, &req)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal memperbarui holiday: " + err.Error(),
		}
		if err.Error() == "holiday not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Holiday berhasil diperbarui",
		Data:    updatedHoliday,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteHoliday handles soft deleting a holiday
func (hc *HolidayController) DeleteHoliday(w http.ResponseWriter, r *http.Request) {
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

	// Extract holiday ID from URL path
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid holiday ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete holiday
	err = hc.HolidayModel.Delete(id)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus holiday: " + err.Error(),
		}
		if err.Error() == "holiday not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Holiday berhasil dihapus",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDeleteHoliday handles permanently deleting a holiday
func (hc *HolidayController) HardDeleteHoliday(w http.ResponseWriter, r *http.Request) {
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

    // Extract holiday ID from URL path
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Invalid holiday ID",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    // Delete holiday permanently
    err = hc.HolidayModel.HardDelete(id)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Gagal menghapus holiday secara permanen: " + err.Error(),
        }
        if err.Error() == "holiday not found" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Holiday berhasil dihapus secara permanen",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// ActivateHoliday handles reactivating a soft-deleted holiday
func (hc *HolidayController) ActivateHoliday(w http.ResponseWriter, r *http.Request) {
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

    // Extract holiday ID from URL path
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Invalid holiday ID",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    // Activate holiday
    err = hc.HolidayModel.Activate(id)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Gagal mengaktifkan holiday: " + err.Error(),
        }
        if err.Error() == "holiday not found" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Holiday berhasil diaktifkan",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
} 