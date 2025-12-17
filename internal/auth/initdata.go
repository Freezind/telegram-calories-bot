package auth

import (
	"encoding/json"
	"errors"
	"net/url"
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

	// Parse the JSON user object
	var user TelegramUser
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil, errors.New("failed to parse user JSON: " + err.Error())
	}

	// Validate that we got a user ID
	if user.ID == 0 {
		return nil, errors.New("user ID is zero or missing in parsed JSON")
	}

	// For demo MVP, we skip hash validation
	// In production, you would validate the hash to ensure data authenticity

	return &user, nil
}
