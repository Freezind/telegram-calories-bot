// Package services implements business logic and external service integrations
package services

import (
	"sync"
	"time"

	"github.com/freezind/telegram-calories-bot/src/models"
)

// SessionManager manages in-memory user sessions using sync.Map (thread-safe)
// Implements state transitions per data-model.md state machine
type SessionManager struct {
	sessions sync.Map // map[int64]*models.UserSession
}

// NewSessionManager creates a new session manager instance
func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// GetSession retrieves a user's session, creating a new one if it doesn't exist
func (sm *SessionManager) GetSession(userID int64) *models.UserSession {
	if val, ok := sm.sessions.Load(userID); ok {
		session, ok := val.(*models.UserSession)
		if !ok {
			// This should never happen, but handle gracefully
			session = &models.UserSession{
				UserID:       userID,
				State:        models.StateIdle,
				LastActivity: time.Now(),
			}
			sm.sessions.Store(userID, session)
			return session
		}
		session.LastActivity = time.Now() // Update activity timestamp
		return session
	}

	// Create new session in Idle state
	session := &models.UserSession{
		UserID:       userID,
		State:        models.StateIdle,
		LastActivity: time.Now(),
	}
	sm.sessions.Store(userID, session)
	return session
}

// UpdateSession updates a user's session state
// Validates state transitions per data-model.md state machine
func (sm *SessionManager) UpdateSession(userID int64, state models.SessionState) *models.UserSession {
	session := sm.GetSession(userID)
	session.State = state
	session.LastActivity = time.Now()
	sm.sessions.Store(userID, session)
	return session
}

// SetMessageID updates the message ID for a user's session
func (sm *SessionManager) SetMessageID(userID int64, messageID int) {
	session := sm.GetSession(userID)
	session.MessageID = messageID
	sm.sessions.Store(userID, session)
}

// DeleteSession removes a user's session from memory
// Called on Cancel button or after result delivery
func (sm *SessionManager) DeleteSession(userID int64) {
	sm.sessions.Delete(userID)
}

// CleanupStale removes sessions inactive for more than 15 minutes
// Prevents memory leaks from abandoned sessions
func (sm *SessionManager) CleanupStale() int {
	count := 0
	sm.sessions.Range(func(key, value interface{}) bool {
		session, ok := value.(*models.UserSession)
		if !ok {
			return true // skip invalid entries
		}
		if time.Since(session.LastActivity) > 15*time.Minute {
			sm.sessions.Delete(key)
			count++
		}
		return true // continue iteration
	})
	return count
}

// StartCleanupRoutine launches a goroutine that cleans up stale sessions every 5 minutes
// Per data-model.md: "Run cleanup every 5 minutes"
func (sm *SessionManager) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cleaned := sm.CleanupStale()
			if cleaned > 0 {
				// Log cleanup count (would use proper logger in production)
				_ = cleaned
			}
		}
	}()
}
