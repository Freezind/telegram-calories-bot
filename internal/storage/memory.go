package storage

import (
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/freezind/telegram-calories-bot/internal/models"
)

// MemoryStorage implements LogStorage using in-memory map
type MemoryStorage struct {
	mu   sync.RWMutex
	logs map[int64][]models.Log
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		logs: make(map[int64][]models.Log),
	}
}

// ListLogs retrieves all logs for a user, sorted by Timestamp descending
func (s *MemoryStorage) ListLogs(userID int64) ([]models.Log, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logs, exists := s.logs[userID]
	if !exists {
		// First-time user: no prior storage record exists
		log.Printf("[STORAGE] User %d has no storage record yet (first access) - returning empty list", userID)
		return []models.Log{}, nil
	}

	if len(logs) == 0 {
		log.Printf("[STORAGE] User %d has 0 logs", userID)
	} else {
		log.Printf("[STORAGE] User %d has %d log(s)", userID, len(logs))
	}

	// Create a copy to avoid external mutations
	result := make([]models.Log, len(logs))
	copy(result, logs)

	// Sort by Timestamp descending (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	return result, nil
}

// CreateLog creates a new log entry
func (s *MemoryStorage) CreateLog(userID int64, logEntry *models.Log) error {
	if err := logEntry.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate UUID for the log
	logEntry.ID = uuid.New().String()
	logEntry.UserID = userID
	now := time.Now()
	logEntry.CreatedAt = now
	logEntry.UpdatedAt = now

	// Initialize user's log slice if it doesn't exist
	if _, exists := s.logs[userID]; !exists {
		log.Printf("[STORAGE] Initializing storage for user %d (first log creation)", userID)
		s.logs[userID] = []models.Log{}
	}

	s.logs[userID] = append(s.logs[userID], *logEntry)
	log.Printf("[STORAGE] Created log for user %d (total logs: %d)", userID, len(s.logs[userID]))
	return nil
}

// UpdateLog updates an existing log entry
func (s *MemoryStorage) UpdateLog(userID int64, logID string, update *models.LogUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logs, exists := s.logs[userID]
	if !exists {
		return errors.New("log not found")
	}

	// Find the log by ID
	logIndex := -1
	for i, logEntry := range logs {
		if logEntry.ID == logID {
			// Authorization check: verify log belongs to user
			if logEntry.UserID != userID {
				return errors.New("unauthorized: log does not belong to user")
			}
			logIndex = i
			break
		}
	}

	if logIndex == -1 {
		return errors.New("log not found")
	}

	// Apply updates
	logEntry := &s.logs[userID][logIndex]
	if update.FoodItems != nil {
		logEntry.FoodItems = *update.FoodItems
	}
	if update.Calories != nil {
		logEntry.Calories = *update.Calories
	}
	if update.Confidence != nil {
		logEntry.Confidence = *update.Confidence
	}
	if update.Timestamp != nil {
		logEntry.Timestamp = *update.Timestamp
	}
	logEntry.UpdatedAt = time.Now()

	// Validate after updates
	if err := logEntry.Validate(); err != nil {
		return err
	}

	return nil
}

// DeleteLog deletes a log entry
func (s *MemoryStorage) DeleteLog(userID int64, logID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	logs, exists := s.logs[userID]
	if !exists {
		return errors.New("log not found")
	}

	// Find and remove the log
	for i, logEntry := range logs {
		if logEntry.ID == logID {
			// Authorization check: verify log belongs to user
			if logEntry.UserID != userID {
				return errors.New("unauthorized: log does not belong to user")
			}
			// Remove log by slicing
			s.logs[userID] = append(logs[:i], logs[i+1:]...)
			return nil
		}
	}

	return errors.New("log not found")
}
