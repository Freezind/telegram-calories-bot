package models

import (
	"errors"
	"strings"
	"time"
)

// ConfidenceLevel represents the AI's confidence in calorie estimation
type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceLow    ConfidenceLevel = "low"
)

// Log represents a calorie log entry
type Log struct {
	ID         string          `json:"id"`
	UserID     int64           `json:"userId"`
	FoodItems  []string        `json:"foodItems"`
	Calories   int             `json:"calories"`
	Confidence ConfidenceLevel `json:"confidence"`
	Timestamp  time.Time       `json:"timestamp"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// LogUpdate represents partial updates to a log entry
type LogUpdate struct {
	FoodItems  *[]string        `json:"foodItems,omitempty"`
	Calories   *int             `json:"calories,omitempty"`
	Confidence *ConfidenceLevel `json:"confidence,omitempty"`
	Timestamp  *time.Time       `json:"timestamp,omitempty"`
}

// Validate performs validation on a Log instance
func (l *Log) Validate() error {
	// Calories must be non-negative
	if l.Calories < 0 {
		return errors.New("calories must be non-negative")
	}

	// Confidence must be valid enum value
	if l.Confidence != ConfidenceHigh && l.Confidence != ConfidenceMedium && l.Confidence != ConfidenceLow {
		return errors.New("confidence must be one of: high, medium, low")
	}

	// Food items constraints
	if len(l.FoodItems) == 0 {
		return errors.New("food items cannot be empty")
	}
	if len(l.FoodItems) > 10 {
		return errors.New("food items cannot exceed 10 items")
	}

	// Check total length of all food items
	totalLength := 0
	for _, item := range l.FoodItems {
		item = strings.TrimSpace(item)
		if item == "" {
			return errors.New("food items cannot contain empty strings")
		}
		totalLength += len(item)
	}
	if totalLength > 1000 {
		return errors.New("total food items text cannot exceed 1000 characters")
	}

	return nil
}
