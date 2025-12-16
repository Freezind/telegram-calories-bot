package storage

import "github.com/freezind/telegram-calories-bot/internal/models"

// LogStorage defines the interface for log data persistence
type LogStorage interface {
	// ListLogs retrieves all logs for a given user, sorted by Timestamp descending
	ListLogs(userID int64) ([]models.Log, error)

	// CreateLog creates a new log entry for a user
	CreateLog(userID int64, log *models.Log) error

	// UpdateLog updates an existing log entry
	// Returns error if log not found or user is not authorized
	UpdateLog(userID int64, logID string, update *models.LogUpdate) error

	// DeleteLog deletes a log entry
	// Returns error if log not found or user is not authorized
	DeleteLog(userID int64, logID string) error
}
