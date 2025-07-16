package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/username/backend/models"
)

// StyleController handles style management operations
type StyleController struct {
	StyleModel *models.StyleModel
}

// NewStyleController creates new style controller instance
func NewStyleController(styleModel *models.StyleModel) *StyleController {
	return &StyleController{
		StyleModel: styleModel,
	}
}

// GetAllStyles handles getting all styles
func (sc *StyleController) GetAllStyles(w http.ResponseWriter, r *http.Request) {
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

	styles, err := sc.StyleModel.GetAll(page, pageSize, search)
	if err != nil {
		fmt.Printf("Error getting styles: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal mengambil data style: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Style berhasil diambil",
		Data:    styles,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateStyle handles creating a new style
func (sc *StyleController) CreateStyle(w http.ResponseWriter, r *http.Request) {
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

	var req models.StyleRequest
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
	if req.StyleNo == "" || req.BuyerShortName == "" || req.Unit == "" {
		response := Response{
			Success: false,
			Message: "Style No, Buyer, dan Unit wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate Unit value
	if req.Unit != "SET" && req.Unit != "PCS" {
		response := Response{
			Success: false,
			Message: "Unit harus SET atau PCS",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate T/B value if Unit is SET
	if req.Unit == "SET" {
		if req.TB == nil {
			response := Response{
				Success: false,
				Message: "T/B wajib diisi jika Unit adalah SET",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		tbValue := *req.TB
		if tbValue != "T1" && tbValue != "B1" {
			response := Response{
				Success: false,
				Message: "T/B harus T1 atau B1",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Check if style already exists
	existingStyle, _ := sc.StyleModel.GetByStyleNo(req.StyleNo)
	if existingStyle != nil {
		response := Response{
			Success: false,
			Message: "Style No sudah digunakan",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create style
	style, err := sc.StyleModel.Create(&req)
	if err != nil {
		fmt.Printf("Error creating style: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal membuat style: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Style berhasil dibuat",
		Data:    style,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateStyle handles updating style information
func (sc *StyleController) UpdateStyle(w http.ResponseWriter, r *http.Request) {
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

	// Extract style number from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/styles/")
	styleNo := path

	var req models.StyleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Format JSON tidak valid",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate Unit value
	if req.Unit != "SET" && req.Unit != "PCS" {
		response := Response{
			Success: false,
			Message: "Unit harus SET atau PCS",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate T/B value if Unit is SET
	if req.Unit == "SET" {
		if req.TB == nil {
			response := Response{
				Success: false,
				Message: "T/B wajib diisi jika Unit adalah SET",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		tbValue := *req.TB
		if tbValue != "T1" && tbValue != "B1" {
			response := Response{
				Success: false,
				Message: "T/B harus T1 atau B1",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Update style
	req.StyleNo = styleNo // Ensure we use the path parameter
	updatedStyle, err := sc.StyleModel.Update(styleNo, &req)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal memperbarui style: " + err.Error(),
		}
		if err.Error() == "style not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Style berhasil diperbarui",
		Data:    updatedStyle,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteStyle handles deleting a style
func (sc *StyleController) DeleteStyle(w http.ResponseWriter, r *http.Request) {
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

	// Extract style number from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/styles/")
	styleNo := path

	// Delete style
	err := sc.StyleModel.Delete(styleNo)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus style: " + err.Error(),
		}
		if err.Error() == "style not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Style berhasil dihapus",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDeleteStyle handles permanently deleting a style
func (sc *StyleController) HardDeleteStyle(w http.ResponseWriter, r *http.Request) {
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

	// Extract style number from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/styles/hard-delete/")
	styleNo := path

	// Delete style permanently
	err := sc.StyleModel.HardDelete(styleNo)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus style secara permanen: " + err.Error(),
		}
		if err.Error() == "style not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Style berhasil dihapus secara permanen",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ActivateStyle handles reactivating a soft-deleted style
func (sc *StyleController) ActivateStyle(w http.ResponseWriter, r *http.Request) {
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

	// Extract style number from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/styles/activate/")
	styleNo := path

	// Activate style
	err := sc.StyleModel.Activate(styleNo)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal mengaktifkan style: " + err.Error(),
		}
		if err.Error() == "style not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Style berhasil diaktifkan",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 