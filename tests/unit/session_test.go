package unit

import (
	"testing"
	"time"

	"github.com/freezind/telegram-calories-bot/src/models"
	"github.com/freezind/telegram-calories-bot/src/services"
	"github.com/stretchr/testify/assert"
)

// T020: Unit test for SessionManager state transitions
// Tests: Idle→AwaitingImage, AwaitingImage→Processing, Processing→Idle

func TestSessionManager_GetSession_CreatesNewSession(t *testing.T) {
	sm := services.NewSessionManager()
	userID := int64(12345)

	session := sm.GetSession(userID)

	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, models.StateIdle, session.State)
	assert.WithinDuration(t, time.Now(), session.LastActivity, time.Second)
}

func TestSessionManager_GetSession_ReturnsExistingSession(t *testing.T) {
	sm := services.NewSessionManager()
	userID := int64(12345)

	// Create initial session
	session1 := sm.GetSession(userID)
	initialTime := session1.LastActivity

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Get session again
	session2 := sm.GetSession(userID)

	assert.Equal(t, session1.UserID, session2.UserID)
	assert.NotEqual(t, initialTime, session2.LastActivity, "LastActivity should be updated")
}

func TestSessionManager_UpdateSession_StateTransitions(t *testing.T) {
	sm := services.NewSessionManager()
	userID := int64(12345)

	// Idle → AwaitingImage
	session := sm.UpdateSession(userID, models.StateAwaitingImage)
	assert.Equal(t, models.StateAwaitingImage, session.State)

	// AwaitingImage → Processing
	session = sm.UpdateSession(userID, models.StateProcessing)
	assert.Equal(t, models.StateProcessing, session.State)

	// Processing → Idle
	session = sm.UpdateSession(userID, models.StateIdle)
	assert.Equal(t, models.StateIdle, session.State)
}

func TestSessionManager_SetMessageID(t *testing.T) {
	sm := services.NewSessionManager()
	userID := int64(12345)
	messageID := 999

	sm.SetMessageID(userID, messageID)
	session := sm.GetSession(userID)

	assert.Equal(t, messageID, session.MessageID)
}

func TestSessionManager_DeleteSession(t *testing.T) {
	sm := services.NewSessionManager()
	userID := int64(12345)

	// Create session
	sm.GetSession(userID)

	// Delete it
	sm.DeleteSession(userID)

	// Getting again should create new session
	newSession := sm.GetSession(userID)
	assert.Equal(t, models.StateIdle, newSession.State)
	assert.Equal(t, 0, newSession.MessageID, "New session should have zero message ID")
}

func TestSessionManager_CleanupStale(t *testing.T) {
	sm := services.NewSessionManager()

	// Create multiple sessions
	user1 := int64(111)
	user2 := int64(222)
	user3 := int64(333)

	sm.GetSession(user1)
	sm.GetSession(user2)
	sm.GetSession(user3)

	// Manually set one session to be stale (>15 min old)
	session := sm.GetSession(user2)
	session.LastActivity = time.Now().Add(-20 * time.Minute)

	// Run cleanup
	cleaned := sm.CleanupStale()

	// Should have cleaned 1 session
	assert.Equal(t, 1, cleaned)

	// Verify user2's session was cleaned (will create new one)
	newSession := sm.GetSession(user2)
	assert.Equal(t, models.StateIdle, newSession.State)
	assert.WithinDuration(t, time.Now(), newSession.LastActivity, time.Second)
}

func TestSessionManager_ConcurrentAccess(t *testing.T) {
	sm := services.NewSessionManager()
	userID := int64(12345)
	iterations := 100

	// Concurrent reads and writes
	done := make(chan bool)

	go func() {
		for i := 0; i < iterations; i++ {
			sm.GetSession(userID)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			sm.UpdateSession(userID, models.StateAwaitingImage)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Should not crash and should have valid session
	session := sm.GetSession(userID)
	assert.NotNil(t, session)
}
