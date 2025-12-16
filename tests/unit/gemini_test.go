package unit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/freezind/telegram-calories-bot/src/models"
	"github.com/stretchr/testify/assert"
)

// T021: Unit test for GeminiClient JSON parsing
// Tests: mock Gemini API responses, test success + no-food + error cases

func TestGeminiJSONParsing_Success(t *testing.T) {
	tests := []struct {
		name     string
		jsonResp string
		expected models.EstimateResult
	}{
		{
			name: "high confidence with food",
			jsonResp: `{
				"calories": 650,
				"confidence": "high",
				"items": ["Grilled chicken breast (200g)", "Steamed broccoli (100g)", "Brown rice (150g)"],
				"reasoning": "Standard portions for grilled chicken plate"
			}`,
			expected: models.EstimateResult{
				Calories:   650,
				Confidence: "high",
				FoodItems:  []string{"Grilled chicken breast (200g)", "Steamed broccoli (100g)", "Brown rice (150g)"},
				Reasoning:  "Standard portions for grilled chicken plate",
			},
		},
		{
			name: "medium confidence",
			jsonResp: `{
				"calories": 300,
				"confidence": "medium",
				"items": ["Mixed salad", "Dressing"],
				"reasoning": "Portion sizes estimated"
			}`,
			expected: models.EstimateResult{
				Calories:   300,
				Confidence: "medium",
				FoodItems:  []string{"Mixed salad", "Dressing"},
				Reasoning:  "Portion sizes estimated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result models.EstimateResult
			err := json.Unmarshal([]byte(tt.jsonResp), &result)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Calories, result.Calories)
			assert.Equal(t, tt.expected.Confidence, result.Confidence)
			assert.Equal(t, tt.expected.FoodItems, result.FoodItems)
			assert.Equal(t, tt.expected.Reasoning, result.Reasoning)

			// Should pass validation
			assert.NoError(t, result.Validate())
		})
	}
}

func TestGeminiJSONParsing_NoFood(t *testing.T) {
	jsonResp := `{
		"calories": 0,
		"confidence": "low",
		"items": [],
		"reasoning": "No food detected"
	}`

	var result models.EstimateResult
	err := json.Unmarshal([]byte(jsonResp), &result)

	assert.NoError(t, err)
	assert.Equal(t, 0, result.Calories)
	assert.Equal(t, "low", result.Confidence)
	assert.Empty(t, result.FoodItems)
	assert.False(t, result.HasFood(), "Should not have food")
}

func TestGeminiJSONParsing_WithMarkdownCodeBlocks(t *testing.T) {
	// Test response wrapped in markdown code blocks (common Gemini behavior)
	jsonResp := "```json\n" + `{
		"calories": 450,
		"confidence": "high",
		"items": ["Pizza slice"],
		"reasoning": "One slice visible"
	}` + "\n```"

	// Simulate cleaning (as done in gemini.go)
	cleaned := strings.TrimSpace(jsonResp)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var result models.EstimateResult
	err := json.Unmarshal([]byte(cleaned), &result)

	assert.NoError(t, err)
	assert.Equal(t, 450, result.Calories)
	assert.Equal(t, "high", result.Confidence)
}

func TestGeminiJSONParsing_InvalidJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonResp string
	}{
		{
			name:     "malformed JSON",
			jsonResp: `{"calories": 450, "confidence": "high"`,
		},
		{
			name:     "not JSON at all",
			jsonResp: "This is just text, not JSON",
		},
		{
			name:     "empty response",
			jsonResp: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result models.EstimateResult
			err := json.Unmarshal([]byte(tt.jsonResp), &result)
			assert.Error(t, err, "Should fail to parse invalid JSON")
		})
	}
}

func TestGeminiJSONParsing_ValidationFailures(t *testing.T) {
	tests := []struct {
		name     string
		jsonResp string
		errMsg   string
	}{
		{
			name: "negative calories",
			jsonResp: `{
				"calories": -100,
				"confidence": "high",
				"items": []
			}`,
			errMsg: "calories must be non-negative",
		},
		{
			name: "invalid confidence",
			jsonResp: `{
				"calories": 100,
				"confidence": "super_high",
				"items": []
			}`,
			errMsg: "confidence must be low/medium/high",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result models.EstimateResult
			err := json.Unmarshal([]byte(tt.jsonResp), &result)
			assert.NoError(t, err, "JSON should parse")

			// But validation should fail
			err = result.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
