package auth

import (
	"errors"
	"net/url"
	"strconv"
)

// TelegramUser represents a user extracted from Telegram initData
type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// ParseInitData extracts user information from Telegram WebApp initData
// The initData is a URL-encoded query string containing user info
func ParseInitData(initData string) (*TelegramUser, error) {
	if initData == "" {
		return nil, errors.New("initData is empty")
	}

	// Parse the URL-encoded query string
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, errors.New("invalid initData format")
	}

	// Extract user data (Telegram WebApp includes a "user" field as JSON string)
	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, errors.New("user data not found in initData")
	}

	// For demo MVP, we'll do simple parsing
	// In production, you would:
	// 1. Validate the hash to ensure data authenticity
	// 2. Parse the JSON user object properly
	// For now, we'll extract just the user ID from the query string

	// Alternative: Try to get id directly if it's in the query
	userIDStr := values.Get("id")
	if userIDStr == "" {
		// Try parsing from user JSON - simplified approach
		// In real implementation, unmarshal JSON properly
		return nil, errors.New("user ID not found in initData")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	return &TelegramUser{
		ID:        userID,
		FirstName: values.Get("first_name"),
		LastName:  values.Get("last_name"),
		Username:  values.Get("username"),
	}, nil
}
