package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/freezind/telegram-calories-bot/src/models"
	"google.golang.org/genai"
)

// GeminiClient wraps Google Gemini SDK for calorie estimation
// Handles API calls per contracts/gemini-vision.yaml
type GeminiClient struct {
	apiKey string
	model  string
}

// NewGeminiClient creates a new Gemini client instance
// Validates API key exists in environment
func NewGeminiClient() (*GeminiClient, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	return &GeminiClient{
		apiKey: apiKey,
		model:  "gemini-2.5-flash", // Fast, cost-effective model per research.md
	}, nil
}

// EstimateCalories analyzes a food image and returns calorie estimate
// Uses structured JSON prompt per contracts/gemini-vision.yaml
func (gc *GeminiClient) EstimateCalories(ctx context.Context, imageBytes []byte, mimeType string) (*models.EstimateResult, error) {
	// Create client with timeout (30 seconds per data-model.md)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Initialize client per research.md Decision 1
	client, err := genai.NewClient(ctxWithTimeout, &genai.ClientConfig{
		APIKey:  gc.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	// Structured prompt per research.md Decision 3 and contracts/gemini-vision.yaml
	prompt := `You are a nutrition analysis assistant. Analyze this food image and estimate total calories.

Output ONLY valid JSON with this exact structure:
{
  "calories": <number>,
  "confidence": "low|medium|high",
  "items": ["food1", "food2", ...],
  "reasoning": "brief explanation"
}

Confidence levels:
- high: Common foods, clear portions visible
- medium: Some foods recognizable, portions estimated
- low: Unclear foods or portions, or non-food image

If no food detected, return:
{"calories": 0, "confidence": "low", "items": [], "reasoning": "No food detected"}

Example (grilled chicken with vegetables):
{"calories": 450, "confidence": "high", "items": ["Grilled chicken breast (200g)", "Steamed broccoli (100g)", "Brown rice (150g)"], "reasoning": "Standard portions for grilled chicken plate"}`

	// Create multimodal content: prompt + image (per research.md)
	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
		genai.NewPartFromBytes(imageBytes, mimeType), // Supports JPEG, PNG, WebP
	}

	content := []*genai.Content{{
		Parts: parts,
		Role:  genai.RoleUser,
	}}

	// Generate content using Gemini 2.5 Flash (optimized for speed)
	response, err := client.Models.GenerateContent(ctxWithTimeout, gc.model, content, &genai.GenerateContentConfig{
		Temperature: floatPtr(0.2), // Low temperature for deterministic output
	})
	if err != nil {
		return nil, fmt.Errorf("gemini API call failed: %w", err)
	}

	// Parse response
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from Gemini API")
	}

	// Extract text from first part
	textPart := response.Candidates[0].Content.Parts[0].Text
	if textPart == "" {
		return nil, fmt.Errorf("unexpected empty response from Gemini API")
	}

	// Clean JSON response (remove markdown code blocks if present)
	jsonText := strings.TrimSpace(textPart)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	// Unmarshal JSON to EstimateResult
	var result models.EstimateResult
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini JSON response: %w (response: %s)", err, jsonText)
	}

	// Validate result per data-model.md
	if err := result.Validate(); err != nil {
		return nil, fmt.Errorf("invalid result from Gemini: %w", err)
	}

	return &result, nil
}

// floatPtr returns a pointer to a float32 value
func floatPtr(f float32) *float32 {
	return &f
}
