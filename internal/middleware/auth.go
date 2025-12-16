package middleware

import (
	"context"
	"net/http"
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
		// Extract X-Telegram-Init-Data header
		initData := r.Header.Get("X-Telegram-Init-Data")
		if initData == "" {
			http.Error(w, "Unauthorized: X-Telegram-Init-Data header missing", http.StatusUnauthorized)
			return
		}

		// Parse initData to extract user information
		user, err := auth.ParseInitData(initData)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid initData - "+err.Error(), http.StatusUnauthorized)
			return
		}

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
