package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/freezind/telegram-calories-bot/internal/middleware"
	"github.com/freezind/telegram-calories-bot/internal/models"
	"github.com/freezind/telegram-calories-bot/internal/storage"
)

// LogsHandler handles log-related HTTP requests
type LogsHandler struct {
	storage storage.LogStorage
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(storage storage.LogStorage) *LogsHandler {
	return &LogsHandler{storage: storage}
}

// ListLogs handles GET /api/logs
func (h *LogsHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	// Extract userID from context (added by AuthMiddleware)
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: user ID not found in context", http.StatusUnauthorized)
		return
	}

	// Fetch logs for user
	logs, err := h.storage.ListLogs(userID)
	if err != nil {
		http.Error(w, "Failed to fetch logs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return logs as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateLog handles POST /api/logs
func (h *LogsHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	// Extract userID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: user ID not found in context", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var log models.Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create log (storage will generate ID and timestamps)
	if err := h.storage.CreateLog(userID, &log); err != nil {
		http.Error(w, "Failed to create log: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Return created log as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(log); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateLog handles PATCH /api/logs/:id
func (h *LogsHandler) UpdateLog(w http.ResponseWriter, r *http.Request) {
	// Extract userID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: user ID not found in context", http.StatusUnauthorized)
		return
	}

	// Extract log ID from URL path
	// Path format: /api/logs/{id}
	path := r.URL.Path
	logID := path[len("/api/logs/"):]
	if logID == "" {
		http.Error(w, "Log ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var update models.LogUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update log
	if err := h.storage.UpdateLog(userID, logID, &update); err != nil {
		if err.Error() == "log not found" {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: log does not belong to user" {
			http.Error(w, "Unauthorized: log does not belong to user", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to update log: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch updated log to return
	logs, err := h.storage.ListLogs(userID)
	if err != nil {
		http.Error(w, "Failed to fetch updated log", http.StatusInternalServerError)
		return
	}

	// Find the updated log
	var updatedLog *models.Log
	for _, log := range logs {
		if log.ID == logID {
			updatedLog = &log
			break
		}
	}

	if updatedLog == nil {
		http.Error(w, "Updated log not found", http.StatusInternalServerError)
		return
	}

	// Return updated log as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedLog); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteLog handles DELETE /api/logs/:id
func (h *LogsHandler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	// Extract userID from context
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized: user ID not found in context", http.StatusUnauthorized)
		return
	}

	// Extract log ID from URL path
	path := r.URL.Path
	logID := path[len("/api/logs/"):]
	if logID == "" {
		http.Error(w, "Log ID is required", http.StatusBadRequest)
		return
	}

	// Delete log
	if err := h.storage.DeleteLog(userID, logID); err != nil {
		if err.Error() == "log not found" {
			http.Error(w, "Log not found", http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized: log does not belong to user" {
			http.Error(w, "Unauthorized: log does not belong to user", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to delete log: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}
