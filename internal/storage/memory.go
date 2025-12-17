package storage

import (
	"errors"
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
		return []models.Log{}, nil
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
func (s *MemoryStorage) CreateLog(userID int64, log *models.Log) error {
	if err := log.Validate(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate UUID for the log
	log.ID = uuid.New().String()
	log.UserID = userID
	now := time.Now()
	log.CreatedAt = now
	log.UpdatedAt = now

	// Initialize user's log slice if it doesn't exist
	if _, exists := s.logs[userID]; !exists {
		s.logs[userID] = []models.Log{}
	}

	s.logs[userID] = append(s.logs[userID], *log)
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
	for i, log := range logs {
		if log.ID == logID {
			// Authorization check: verify log belongs to user
			if log.UserID != userID {
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
	log := &s.logs[userID][logIndex]
	if update.FoodItems != nil {
		log.FoodItems = *update.FoodItems
	}
	if update.Calories != nil {
		log.Calories = *update.Calories
	}
	if update.Confidence != nil {
		log.Confidence = *update.Confidence
	}
	if update.Timestamp != nil {
		log.Timestamp = *update.Timestamp
	}
	log.UpdatedAt = time.Now()

	// Validate after updates
	if err := log.Validate(); err != nil {
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
	for i, log := range logs {
		if log.ID == logID {
			// Authorization check: verify log belongs to user
			if log.UserID != userID {
				return errors.New("unauthorized: log does not belong to user")
			}
			// Remove log by slicing
			s.logs[userID] = append(logs[:i], logs[i+1:]...)
			return nil
		}
	}

	return errors.New("log not found")
}
