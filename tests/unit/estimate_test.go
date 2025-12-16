package unit

import (
	"testing"

	"github.com/freezind/telegram-calories-bot/src/models"
	"github.com/stretchr/testify/assert"
)

// T019: Unit test for EstimateResult validation
// Tests: calories >0, confidence in {low,medium,high}, empty items for no food

func TestEstimateResult_Validate_ValidResults(t *testing.T) {
	tests := []struct {
		name   string
		result models.EstimateResult
	}{
		{
			name: "valid high confidence",
			result: models.EstimateResult{
				Calories:   450,
				Confidence: "high",
				FoodItems:  []string{"Chicken", "Rice"},
			},
		},
		{
			name: "valid medium confidence",
			result: models.EstimateResult{
				Calories:   300,
				Confidence: "medium",
				FoodItems:  []string{"Salad"},
			},
		},
		{
			name: "valid low confidence",
			result: models.EstimateResult{
				Calories:   100,
				Confidence: "low",
				FoodItems:  []string{},
			},
		},
		{
			name: "zero calories with low confidence (no food)",
			result: models.EstimateResult{
				Calories:   0,
				Confidence: "low",
				FoodItems:  []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.Validate()
			assert.NoError(t, err, "Valid result should not return error")
		})
	}
}

func TestEstimateResult_Validate_InvalidResults(t *testing.T) {
	tests := []struct {
		name        string
		result      models.EstimateResult
		expectedErr string
	}{
		{
			name: "negative calories",
			result: models.EstimateResult{
				Calories:   -100,
				Confidence: "high",
			},
			expectedErr: "calories must be non-negative",
		},
		{
			name: "invalid confidence level",
			result: models.EstimateResult{
				Calories:   450,
				Confidence: "very_high",
			},
			expectedErr: "confidence must be low/medium/high",
		},
		{
			name: "empty confidence",
			result: models.EstimateResult{
				Calories:   450,
				Confidence: "",
			},
			expectedErr: "confidence must be low/medium/high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.Validate()
			assert.Error(t, err, "Invalid result should return error")
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestEstimateResult_HasFood(t *testing.T) {
	tests := []struct {
		name     string
		result   models.EstimateResult
		expected bool
	}{
		{
			name: "has food - with items and calories",
			result: models.EstimateResult{
				Calories:  450,
				FoodItems: []string{"Chicken", "Rice"},
			},
			expected: true,
		},
		{
			name: "no food - empty items",
			result: models.EstimateResult{
				Calories:  0,
				FoodItems: []string{},
			},
			expected: false,
		},
		{
			name: "no food - zero calories with items",
			result: models.EstimateResult{
				Calories:  0,
				FoodItems: []string{"Something"},
			},
			expected: false,
		},
		{
			name: "no food - calories but empty items",
			result: models.EstimateResult{
				Calories:  100,
				FoodItems: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.result.HasFood())
		})
	}
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name     string
		result   models.EstimateResult
		contains []string
	}{
		{
			name: "format with food items",
			result: models.EstimateResult{
				Calories:   450,
				Confidence: "high",
				FoodItems:  []string{"Chicken", "Rice", "Vegetables"},
			},
			contains: []string{"450 kcal", "High", "Chicken", "Rice", "Vegetables"},
		},
		{
			name: "format without food items",
			result: models.EstimateResult{
				Calories:   0,
				Confidence: "low",
				FoodItems:  []string{},
			},
			contains: []string{"0 kcal", "Low", "None detected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := models.FormatResult(&tt.result)
			for _, expected := range tt.contains {
				assert.Contains(t, formatted, expected)
			}
		})
	}
}
