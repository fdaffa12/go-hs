package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/username/backend/models"
)

// BuyerController handles buyer management operations
type BuyerController struct {
	BuyerModel *models.BuyerModel
}

// NewBuyerController creates new buyer controller instance
func NewBuyerController(buyerModel *models.BuyerModel) *BuyerController {
	return &BuyerController{
		BuyerModel: buyerModel,
	}
}

// GetAllBuyers handles getting all buyers
func (bc *BuyerController) GetAllBuyers(w http.ResponseWriter, r *http.Request) {
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

	buyers, err := bc.BuyerModel.GetAll(page, pageSize, search)
	if err != nil {
		fmt.Printf("Error getting buyers: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal mengambil data buyer: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Buyer berhasil diambil",
		Data:    buyers,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateBuyer handles creating a new buyer
func (bc *BuyerController) CreateBuyer(w http.ResponseWriter, r *http.Request) {
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

	var req models.BuyerRequest
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
	if req.ShortName == "" || req.LongName == "" {
		response := Response{
			Success: false,
			Message: "Nama singkat dan nama panjang buyer wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if buyer already exists
	existingBuyer, _ := bc.BuyerModel.GetByShortName(req.ShortName)
	if existingBuyer != nil {
		response := Response{
			Success: false,
			Message: "Kode buyer sudah digunakan",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create buyer
	buyer, err := bc.BuyerModel.Create(&req)
	if err != nil {
		fmt.Printf("Error creating buyer: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Gagal membuat buyer: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Buyer berhasil dibuat",
		Data:    buyer,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateBuyer handles updating buyer information
func (bc *BuyerController) UpdateBuyer(w http.ResponseWriter, r *http.Request) {
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

	// Extract buyer short name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/buyers/")
	shortName := path

	var req models.BuyerRequest
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
			Message: "Nama panjang buyer wajib diisi",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update buyer
	req.ShortName = shortName // Ensure we use the path parameter
	updatedBuyer, err := bc.BuyerModel.Update(shortName, &req)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal memperbarui buyer: " + err.Error(),
		}
		if err.Error() == "buyer not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Buyer berhasil diperbarui",
		Data:    updatedBuyer,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteBuyer handles deleting a buyer
func (bc *BuyerController) DeleteBuyer(w http.ResponseWriter, r *http.Request) {
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

	// Extract buyer short name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/buyers/")
	shortName := path

	// Delete buyer
	err := bc.BuyerModel.Delete(shortName)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Gagal menghapus buyer: " + err.Error(),
		}
		if err.Error() == "buyer not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Buyer berhasil dihapus",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDeleteBuyer handles permanently deleting a buyer
func (bc *BuyerController) HardDeleteBuyer(w http.ResponseWriter, r *http.Request) {
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

    // Extract buyer short name from URL path
    path := strings.TrimPrefix(r.URL.Path, "/api/buyers/hard-delete/")
    shortName := path

    // Delete buyer permanently
    err := bc.BuyerModel.HardDelete(shortName)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Gagal menghapus buyer secara permanen: " + err.Error(),
        }
        if err.Error() == "buyer not found" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Buyer berhasil dihapus secara permanen",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// ActivateBuyer handles reactivating a soft-deleted buyer
func (bc *BuyerController) ActivateBuyer(w http.ResponseWriter, r *http.Request) {
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

    // Extract buyer short name from URL path
    path := strings.TrimPrefix(r.URL.Path, "/api/buyers/activate/")
    shortName := path

    // Activate buyer
    err := bc.BuyerModel.Activate(shortName)
    if err != nil {
        response := Response{
            Success: false,
            Message: "Gagal mengaktifkan buyer: " + err.Error(),
        }
        if err.Error() == "buyer not found" {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Buyer berhasil diaktifkan",
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
} 