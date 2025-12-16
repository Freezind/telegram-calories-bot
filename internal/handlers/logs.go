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
