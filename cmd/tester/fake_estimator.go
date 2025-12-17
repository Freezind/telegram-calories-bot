package main

import (
	"context"

	"github.com/freezind/telegram-calories-bot/src/models"
)

// FakeEstimator returns deterministic estimation results for testing
type FakeEstimator struct{}

// NewFakeEstimator creates a fake estimator for testing
func NewFakeEstimator() *FakeEstimator {
	return &FakeEstimator{}
}

// EstimateFromImage returns a deterministic fake estimate
func (f *FakeEstimator) EstimateFromImage(ctx context.Context, imageBytes []byte, mimeType string) (*models.EstimateResult, error) {
	// Return a deterministic structured estimate for stable testing
	return &models.EstimateResult{
		FoodItems:  []string{"Rice", "Chicken"},
		Calories:   500,
		Confidence: "high",
		Reasoning:  "Test estimation for integration testing",
	}, nil
}
