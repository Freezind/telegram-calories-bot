// Package models defines data structures for calorie estimation
package models

import (
	"fmt"
	"strings"
	"time"
)

// SessionState represents the current state of a user's /estimate flow
type SessionState string

const (
	// StateIdle indicates user has no active estimation flow
	StateIdle SessionState = "idle"

	// StateAwaitingImage indicates bot is waiting for image upload
	StateAwaitingImage SessionState = "awaiting_image"

	// StateProcessing indicates bot is processing uploaded image via Gemini
	StateProcessing SessionState = "processing"
)

// UserSession tracks in-memory session state for a single user during /estimate flow
// Stored in sync.Map keyed by UserID (thread-safe, in-memory only)
type UserSession struct {
	// UserID is the Telegram user ID (from c.Sender().ID)
	UserID int64

	// State is the current flow state
	State SessionState

	// LastActivity tracks the last user interaction (for cleanup after 15 min timeout)
	LastActivity time.Time

	// MessageID is the ID of the last bot message (for editing/deletion)
	MessageID int
}

// EstimateResult holds the calorie estimation output from Gemini Vision API
type EstimateResult struct {
	// Calories is the total estimated calories (kcal)
	Calories int `json:"calories"`

	// Confidence level: "low", "medium", or "high"
	Confidence string `json:"confidence"`

	// FoodItems contains detected food items (empty array if no food)
	FoodItems []string `json:"items,omitempty"`

	// Reasoning provides brief explanation from Gemini (optional)
	Reasoning string `json:"reasoning,omitempty"`
}

// FormatResult formats an EstimateResult into fixed-format response per FR-006
// Returns deterministic, consistent message structure
func FormatResult(result *EstimateResult) string {
	itemsList := "None detected"
	if len(result.FoodItems) > 0 {
		itemsList = strings.Join(result.FoodItems, ", ")
	}

	// Capitalize first letter of confidence
	confidence := result.Confidence
	if len(confidence) > 0 {
		confidence = strings.ToUpper(string(confidence[0])) + strings.ToLower(confidence[1:])
	}

	return fmt.Sprintf(
		"🍽️ Calorie Estimate\n\n"+
			"Estimated Calories: %d kcal\n"+
			"Confidence: %s\n\n"+
			"Detected Items: %s",
		result.Calories,
		confidence,
		itemsList,
	)
}

// HasFood returns true if the result contains recognized food items
// Used to trigger FR-014 error handling when no food detected
func (r *EstimateResult) HasFood() bool {
	return len(r.FoodItems) > 0 && r.Calories > 0
}

// Validate checks if EstimateResult fields meet validation rules from data-model.md
func (r *EstimateResult) Validate() error {
	if r.Calories < 0 {
		return fmt.Errorf("calories must be non-negative, got %d", r.Calories)
	}

	validConfidence := map[string]bool{"low": true, "medium": true, "high": true}
	if !validConfidence[strings.ToLower(r.Confidence)] {
		return fmt.Errorf("confidence must be low/medium/high, got %s", r.Confidence)
	}

	return nil
}
