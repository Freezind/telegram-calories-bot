package main

import (
	"encoding/json"
	"time"
)

// Evidence represents captured test evidence from a scenario execution
type Evidence struct {
	Type      string                 `json:"type"`      // "bot_message", "callback_response", "http_response", "error"
	Timestamp time.Time              `json:"timestamp"` // When evidence was captured
	Data      map[string]interface{} `json:"data"`      // Flexible key-value storage for scenario-specific data
}

// TestResult represents the outcome of a single test scenario
type TestResult struct {
	ScenarioID   string    `json:"scenario_id"`   // e.g., "S1", "S2", "S3", "S4", "S5"
	ScenarioName string    `json:"scenario_name"` // Human-readable scenario name
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Evidence     []Evidence `json:"evidence"`      // All evidence collected during scenario
	Verdict      string    `json:"verdict"`       // "PASS" or "FAIL"
	Rationale    string    `json:"rationale"`     // LLM judge's explanation
	Error        string    `json:"error,omitempty"` // Error message if scenario execution failed
}

// Duration returns the scenario execution time
func (tr *TestResult) Duration() time.Duration {
	return tr.EndTime.Sub(tr.StartTime)
}

// AddEvidence appends evidence to the test result
func (tr *TestResult) AddEvidence(evidenceType string, data map[string]interface{}) {
	tr.Evidence = append(tr.Evidence, Evidence{
		Type:      evidenceType,
		Timestamp: time.Now(),
		Data:      data,
	})
}

// SetError marks the test as failed with an error message
func (tr *TestResult) SetError(err error) {
	tr.Verdict = "FAIL"
	tr.Error = err.Error()
	tr.Rationale = "Test execution error: " + err.Error()
}

// MarshalJSON converts evidence data to JSON string for embedding in reports
func (e *Evidence) MarshalJSON() ([]byte, error) {
	type Alias Evidence
	return json.Marshal(&struct {
		*Alias
		Timestamp string `json:"timestamp"`
	}{
		Alias:     (*Alias)(e),
		Timestamp: e.Timestamp.Format(time.RFC3339),
	})
}

// NewTestResult creates a new test result for a scenario
func NewTestResult(scenarioID, scenarioName string) *TestResult {
	return &TestResult{
		ScenarioID:   scenarioID,
		ScenarioName: scenarioName,
		StartTime:    time.Now(),
		Evidence:     []Evidence{},
	}
}

// Complete marks the test result as complete and sets the end time
func (tr *TestResult) Complete() {
	tr.EndTime = time.Now()
}
