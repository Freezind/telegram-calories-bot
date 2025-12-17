package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/freezind/telegram-calories-bot/internal/auth"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserIDKey is the context key for storing authenticated user ID
	UserIDKey contextKey = "userID"
)

// AuthMiddleware validates Telegram initData and adds userID to request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for DEV_FAKE_USER_ID environment variable (dev/testing only)
		if devUserIDStr := os.Getenv("DEV_FAKE_USER_ID"); devUserIDStr != "" {
			devUserID, err := strconv.ParseInt(devUserIDStr, 10, 64)
			if err != nil {
				log.Printf("[AUTH] ⚠️  DEV FALLBACK: Invalid DEV_FAKE_USER_ID value: %s (error: %v)", devUserIDStr, err)
			} else {
				log.Printf("[AUTH] 🔧 DEV FALLBACK: Using DEV_FAKE_USER_ID=%d for %s %s (initData bypassed)", devUserID, r.Method, r.URL.Path)
				ctx := context.WithValue(r.Context(), UserIDKey, devUserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Extract X-Telegram-Init-Data header
		initData := r.Header.Get("X-Telegram-Init-Data")

		// Debug logging: header presence and length
		if initData == "" {
			log.Printf("[AUTH] ❌ X-Telegram-Init-Data header missing for %s %s (no DEV_FAKE_USER_ID set)", r.Method, r.URL.Path)
			http.Error(w, "Unauthorized: X-Telegram-Init-Data header missing", http.StatusUnauthorized)
			return
		}
		log.Printf("[AUTH] ✓ X-Telegram-Init-Data header present (length: %d) for %s %s", len(initData), r.Method, r.URL.Path)

		// Parse initData to extract user information
		user, err := auth.ParseInitData(initData)
		if err != nil {
			log.Printf("[AUTH] ❌ Failed to parse initData for %s %s: %v", r.Method, r.URL.Path, err)
			http.Error(w, "Unauthorized: Invalid initData - "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Debug logging: successful auth
		log.Printf("[AUTH] ✓ User authenticated: ID=%d, Username=%s for %s %s", user.ID, user.Username, r.Method, r.URL.Path)

		// Add userID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the authenticated user ID from request context
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
