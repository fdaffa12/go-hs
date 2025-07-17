package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/username/backend/models"
)

// LineScheduleController handles line schedule operations
type LineScheduleController struct {
	LineScheduleModel *models.LineScheduleModel
}

// NewLineScheduleController creates a new line schedule controller instance
func NewLineScheduleController(lineScheduleModel *models.LineScheduleModel) *LineScheduleController {
	return &LineScheduleController{
		LineScheduleModel: lineScheduleModel,
	}
}

// GetAllLineSchedules handles getting all line schedules
func (lc *LineScheduleController) GetAllLineSchedules(w http.ResponseWriter, r *http.Request) {
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

	// Get filter parameters
	factory := r.URL.Query().Get("factory")
	date := r.URL.Query().Get("date")

	// Construct filter string
	var filters []string
	if factory != "" {
		filters = append(filters, fmt.Sprintf("factory=%s", factory))
	}
	if date != "" {
		filters = append(filters, fmt.Sprintf("date=%s", date))
	}
	filterStr := strings.Join(filters, "&")

	schedules, err := lc.LineScheduleModel.GetAll(page, pageSize, filterStr)
	if err != nil {
		fmt.Printf("Error getting line schedules: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to get line schedules: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Line schedules retrieved successfully",
		Data:    schedules,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateLineSchedule handles creating a new line schedule
func (lc *LineScheduleController) CreateLineSchedule(w http.ResponseWriter, r *http.Request) {
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
	fmt.Printf("Received create request body: %s\n", string(body))

	// Create a new reader with the body for further processing
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.LineScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Invalid request format: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate required fields
	if req.IDRegistrasi == "" || req.Factory == "" || req.Line == "" || req.BuyerShortName == "" || req.StyleNo == "" {
		response := Response{
			Success: false,
			Message: "All required fields must be provided",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate.Format("2006-01-02"))
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Invalid start date format: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	req.StartDate = &startDate

	date, err := time.Parse("2006-01-02", req.Date.Format("2006-01-02"))
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Invalid date format: %v", err),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	req.Date = date

	// Calculate working days
	workingDays, err := lc.LineScheduleModel.CalculateWorkingDays(*req.StartDate, req.Date)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to calculate working days: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update working days in request
	req.WorkingDay = workingDays

	schedule, err := lc.LineScheduleModel.Create(&req)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to create line schedule: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Line schedule created successfully",
		Data:    schedule,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateLineSchedule handles updating a line schedule
func (lc *LineScheduleController) UpdateLineSchedule(w http.ResponseWriter, r *http.Request) {
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

	// Get ID from query parameter
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid line schedule ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req models.LineScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid request format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate.Format("2006-01-02"))
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid start date format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	req.StartDate = &startDate

	// Parse date
	date, err := time.Parse("2006-01-02", req.Date.Format("2006-01-02"))
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid date format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	req.Date = date

	// Calculate working days
	workingDays, err := lc.LineScheduleModel.CalculateWorkingDays(*req.StartDate, req.Date)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to calculate working days: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	req.WorkingDay = workingDays

	schedule, err := lc.LineScheduleModel.Update(id, &req)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to update line schedule: %v", err),
		}
		if err.Error() == "line schedule not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Line schedule updated successfully",
		Data:    schedule,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteLineSchedule handles deleting a line schedule
func (lc *LineScheduleController) DeleteLineSchedule(w http.ResponseWriter, r *http.Request) {
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
			Message: "Invalid line schedule ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = lc.LineScheduleModel.Delete(id)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to delete line schedule: %v", err),
		}
		if err.Error() == "line schedule not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Line schedule deleted successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HardDeleteLineSchedule handles permanently deleting a line schedule
func (lc *LineScheduleController) HardDeleteLineSchedule(w http.ResponseWriter, r *http.Request) {
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
	path := strings.TrimPrefix(r.URL.Path, "/api/line-schedules/hard-delete/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid line schedule ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = lc.LineScheduleModel.HardDelete(id)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to permanently delete line schedule: %v", err),
		}
		if err.Error() == "line schedule not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Line schedule permanently deleted",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ActivateLineSchedule handles reactivating a soft-deleted line schedule
func (lc *LineScheduleController) ActivateLineSchedule(w http.ResponseWriter, r *http.Request) {
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
	path := strings.TrimPrefix(r.URL.Path, "/api/line-schedules/activate/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid line schedule ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = lc.LineScheduleModel.Activate(id)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to activate line schedule: %v", err),
		}
		if err.Error() == "line schedule not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Line schedule activated successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CalculateWorkingDays calculates working days between two dates
func (lc *LineScheduleController) CalculateWorkingDays(w http.ResponseWriter, r *http.Request) {
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

	// Parse start date and end date from request
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := Response{
			Success: false,
			Message: "Invalid request format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid start date format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		response := Response{
			Success: false,
			Message: "Invalid end date format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate working days
	workingDays, err := lc.LineScheduleModel.CalculateWorkingDays(startDate, endDate)
	if err != nil {
		response := Response{
			Success: false,
			Message: fmt.Sprintf("Failed to calculate working days: %v", err),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Success: true,
		Message: "Working days calculated successfully",
		Data: map[string]int{
			"working_days": workingDays,
		},
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
} 