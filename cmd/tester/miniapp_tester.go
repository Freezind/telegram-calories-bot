package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MiniAppTester tests Mini App HTTP handlers using httptest
type MiniAppTester struct {
	baseURL string
	judge   *GeminiJudge
}

// NewMiniAppTester creates a new Mini App tester instance
func NewMiniAppTester(baseURL string, judge *GeminiJudge) *MiniAppTester {
	return &MiniAppTester{
		baseURL: baseURL,
		judge:   judge,
	}
}

// TestPageLoad tests the Mini App page load (HTTP GET expecting status 200)
func (mt *MiniAppTester) TestPageLoad(ctx context.Context) *TestResult {
	result := NewTestResult("S5", "Mini App Page Load")

	// HTTP GET to Mini App URL
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(mt.baseURL)
	if err != nil {
		result.SetError(fmt.Errorf("failed to GET Mini App URL: %w", err))
		result.Complete()
		return result
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.SetError(fmt.Errorf("failed to read response body: %w", err))
		result.Complete()
		return result
	}

	bodyText := string(bodyBytes)

	// Check for expected key texts in HTML
	keyTexts := []string{"Calorie Log", "Add New Log", "No logs yet"}
	foundTexts := []string{}
	for _, text := range keyTexts {
		if strings.Contains(bodyText, text) {
			foundTexts = append(foundTexts, text)
		}
	}

	// Capture evidence
	result.AddEvidence("page_load", map[string]interface{}{
		"url":          mt.baseURL,
		"status_code":  resp.StatusCode,
		"expected_texts": keyTexts,
		"found_texts":    foundTexts,
		"body_snippet":   truncateString(bodyText, 500),
	})

	// Check status code
	if resp.StatusCode != http.StatusOK {
		result.SetError(fmt.Errorf("expected status 200, got %d", resp.StatusCode))
		result.Complete()
		return result
	}

	// Evaluate with LLM judge
	expectedBehavior := "Mini App page loads with HTTP 200 status and HTML contains at least one of the key texts: 'Calorie Log', 'Add New Log', or 'No logs yet'"
	verdict, err := mt.judge.Evaluate(ctx, result.ScenarioName, expectedBehavior, result.Evidence)
	if err != nil {
		result.SetError(fmt.Errorf("LLM judge evaluation failed: %w", err))
		result.Complete()
		return result
	}

	result.Verdict = verdict.Verdict
	result.Rationale = verdict.Rationale
	result.Complete()
	return result
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "... (truncated)"
}
