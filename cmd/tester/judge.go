package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genai"
)

// GeminiJudge evaluates test scenarios using Gemini LLM
type GeminiJudge struct {
	apiKey  string
	model   string
	prompts []string // Archive all prompts used during testing
}

// JudgeVerdict represents the structured output from LLM judge
type JudgeVerdict struct {
	Verdict   string `json:"verdict"`   // "PASS" or "FAIL"
	Rationale string `json:"rationale"` // Brief explanation
}

// NewGeminiJudge creates a new Gemini judge instance
func NewGeminiJudge(apiKey string) *GeminiJudge {
	return &GeminiJudge{
		apiKey:  apiKey,
		model:   "gemini-2.5-flash", // Fast model for deterministic evaluation
		prompts: []string{},
	}
}

// Evaluate evaluates a test scenario using Gemini LLM judge
func (gj *GeminiJudge) Evaluate(ctx context.Context, scenarioName, expectedBehavior string, evidence []Evidence) (JudgeVerdict, error) {
	// Build evidence JSON
	evidenceJSON, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		return JudgeVerdict{}, fmt.Errorf("failed to marshal evidence: %w", err)
	}

	// Format judge prompt
	prompt := fmt.Sprintf(`You are a test evaluator for a Telegram bot testing system.

**Scenario:** %s
**Expected Behavior:** %s

**Captured Evidence:**
%s

Evaluate whether the captured evidence demonstrates the expected behavior.

Output ONLY valid JSON:
{
  "verdict": "PASS" or "FAIL",
  "rationale": "brief explanation (1-2 sentences)"
}

Rules:
- PASS if evidence clearly matches expected behavior
- FAIL if evidence contradicts expected behavior or is missing critical elements
- Be strict: ambiguous evidence = FAIL
- Do not output anything other than the JSON object`, scenarioName, expectedBehavior, string(evidenceJSON))

	// Archive prompt
	gj.prompts = append(gj.prompts, prompt)

	// Create Gemini client with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctxWithTimeout, &genai.ClientConfig{
		APIKey:  gj.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return JudgeVerdict{}, fmt.Errorf("failed to create genai client: %w", err)
	}

	// Generate content with low temperature for deterministic output
	parts := []*genai.Part{
		genai.NewPartFromText(prompt),
	}

	content := []*genai.Content{{
		Parts: parts,
		Role:  genai.RoleUser,
	}}

	response, err := client.Models.GenerateContent(ctxWithTimeout, gj.model, content, &genai.GenerateContentConfig{
		Temperature: floatPtr(0.0), // Deterministic evaluation
	})
	if err != nil {
		return JudgeVerdict{}, fmt.Errorf("gemini API call failed: %w", err)
	}

	// Parse response
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return JudgeVerdict{}, fmt.Errorf("no response from Gemini API")
	}

	// Extract text from first part
	textPart := response.Candidates[0].Content.Parts[0].Text
	if textPart == "" {
		return JudgeVerdict{}, fmt.Errorf("unexpected empty response from Gemini API")
	}

	// Clean JSON response (remove markdown code blocks if present)
	jsonText := strings.TrimSpace(textPart)
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	// Unmarshal JSON to JudgeVerdict
	var verdict JudgeVerdict
	if err := json.Unmarshal([]byte(jsonText), &verdict); err != nil {
		return JudgeVerdict{}, fmt.Errorf("failed to parse Gemini JSON response: %w (response: %s)", err, jsonText)
	}

	// Validate verdict
	if verdict.Verdict != "PASS" && verdict.Verdict != "FAIL" {
		return JudgeVerdict{}, fmt.Errorf("invalid verdict from Gemini: %s (expected PASS or FAIL)", verdict.Verdict)
	}

	return verdict, nil
}

// GetPrompts returns all prompts used during testing
func (gj *GeminiJudge) GetPrompts() []string {
	return gj.prompts
}

// floatPtr returns a pointer to a float32 value
func floatPtr(f float32) *float32 {
	return &f
}
