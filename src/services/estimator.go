package services

import (
	"context"

	"github.com/freezind/telegram-calories-bot/src/models"
)

// Estimator is the interface for calorie estimation
type Estimator interface {
	// EstimateFromImage analyzes image bytes and returns calorie estimate
	EstimateFromImage(ctx context.Context, imageBytes []byte, mimeType string) (*models.EstimateResult, error)
}

// GeminiEstimator uses Gemini API for real estimation
type GeminiEstimator struct {
	client *GeminiClient
}

// NewGeminiEstimator creates a production estimator using Gemini
func NewGeminiEstimator(client *GeminiClient) Estimator {
	return &GeminiEstimator{client: client}
}

// EstimateFromImage estimates calories from image bytes using Gemini
func (e *GeminiEstimator) EstimateFromImage(ctx context.Context, imageBytes []byte, mimeType string) (*models.EstimateResult, error) {
	return e.client.EstimateCalories(ctx, imageBytes, mimeType)
}
