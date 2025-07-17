package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/username/backend/models"
)

// SewNumProcessController handles sew numbering process operations
type SewNumProcessController struct {
	SewNumProcessModel *models.SewNumProcessModel
}

// NewSewNumProcessController creates a new sew numbering process controller instance
func NewSewNumProcessController(sewNumProcessModel *models.SewNumProcessModel) *SewNumProcessController {
	return &SewNumProcessController{
		SewNumProcessModel: sewNumProcessModel,
	}
}

// GetAllProcesses handles getting all processes for a style
func (c *SewNumProcessController) GetAllProcesses(w http.ResponseWriter, r *http.Request) {
    // Set CORS headers
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    w.Header().Set("Access-Control-Allow-Credentials", "true")

    // Handle preflight requests
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

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

    styleNo := r.URL.Query().Get("style_no")
    if styleNo == "" {
        response := Response{
            Success: false,
            Message: "Style number is required",
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(response)
        return
    }

    processes, err := c.SewNumProcessModel.GetAll(styleNo)
    if err != nil {
        response := Response{
            Success: false,
            Message: fmt.Sprintf("Failed to get processes: %v", err),
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(response)
        return
    }

    response := Response{
        Success: true,
        Message: "Processes retrieved successfully",
        Data:    processes,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

// BulkSave handles bulk create/update operations
func (c *SewNumProcessController) BulkSave(w http.ResponseWriter, r *http.Request) {
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

	// Read and log the request body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Failed to read request body",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Printf("Received bulk save request body: %s\n", string(body))

	// Create a new reader with the body for further processing
	r.Body = io.NopCloser(strings.NewReader(string(body)))

	var processes []*models.SewNumProcess
	if err := json.NewDecoder(r.Body).Decode(&processes); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Invalid request format: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.SewNumProcessModel.BulkSave(processes); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to save processes: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Processes saved successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete handles soft deleting a process
func (c *SewNumProcessController) Delete(w http.ResponseWriter, r *http.Request) {
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

	// Get ID from query parameter
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid process ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.SewNumProcessModel.Delete(id); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to delete process: %v", err),
		}
		if err.Error() == "process not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Process deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// BulkDelete handles bulk soft delete operations
func (c *SewNumProcessController) BulkDelete(w http.ResponseWriter, r *http.Request) {
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

	var ids []int64
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Invalid request format: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.SewNumProcessModel.BulkDelete(ids); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to delete processes: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Processes deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Activate handles reactivating a soft-deleted process
func (c *SewNumProcessController) Activate(w http.ResponseWriter, r *http.Request) {
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

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/sew-num-process/activate/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid process ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.SewNumProcessModel.Activate(id); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to activate process: %v", err),
		}
		if err.Error() == "process not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Process activated successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDelete handles permanently deleting a process
func (c *SewNumProcessController) HardDelete(w http.ResponseWriter, r *http.Request) {
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

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/sew-num-process/hard-delete/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid process ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := c.SewNumProcessModel.HardDelete(id); err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to permanently delete process: %v", err),
		}
		if err.Error() == "process not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Process permanently deleted",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 