package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Reporter generates markdown test reports
type Reporter struct{}

// NewReporter creates a new reporter instance
func NewReporter() *Reporter {
	return &Reporter{}
}

// Generate creates a self-contained markdown report from test results
func (r *Reporter) Generate(results []TestResult, startTime, endTime time.Time) []byte {
	var sb strings.Builder

	// Header
	sb.WriteString("# Test Report: LLM-Based Bot Testing\n\n")
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n", time.Now().UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("**Test Duration:** %s\n", endTime.Sub(startTime).Round(time.Second)))
	sb.WriteString(fmt.Sprintf("**Total Scenarios:** %d\n", len(results)))

	// Count pass/fail
	passed := 0
	failed := 0
	for _, result := range results {
		if result.Verdict == "PASS" {
			passed++
		} else {
			failed++
		}
	}

	sb.WriteString(fmt.Sprintf("**Passed:** %d\n", passed))
	sb.WriteString(fmt.Sprintf("**Failed:** %d\n\n", failed))
	sb.WriteString("---\n\n")

	// Summary table
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Scenario | Verdict | Rationale |\n")
	sb.WriteString("|----------|---------|-----------||\n")

	for _, result := range results {
		verdict := result.Verdict
		if verdict == "PASS" {
			verdict = "✅ PASS"
		} else {
			verdict = "❌ FAIL"
		}
		rationale := strings.ReplaceAll(result.Rationale, "\n", " ")
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", result.ScenarioName, verdict, rationale))
	}

	sb.WriteString("\n---\n\n")

	// Scenario details
	sb.WriteString("## Scenario Details\n\n")

	for _, result := range results {
		sb.WriteString(fmt.Sprintf("### %s: %s\n\n", result.ScenarioID, result.ScenarioName))
		sb.WriteString(fmt.Sprintf("**Duration:** %s\n\n", result.Duration().Round(time.Millisecond)))

		// Evidence
		if len(result.Evidence) > 0 {
			sb.WriteString("**Captured Evidence:**\n\n")
			for i, evidence := range result.Evidence {
				evidenceJSON, _ := json.MarshalIndent(evidence, "", "  ")
				sb.WriteString(fmt.Sprintf("Evidence %d (%s):\n```json\n%s\n```\n\n", i+1, evidence.Type, string(evidenceJSON)))
			}
		}

		// Error details if present
		if result.Error != "" {
			sb.WriteString("**Error:**\n```\n")
			sb.WriteString(result.Error)
			sb.WriteString("\n```\n\n")
		}

		// Verdict
		verdict := result.Verdict
		if verdict == "PASS" {
			sb.WriteString("**LLM Judge Verdict:** ✅ PASS\n")
		} else {
			sb.WriteString("**LLM Judge Verdict:** ❌ FAIL\n")
		}
		sb.WriteString(fmt.Sprintf("**Rationale:** %s\n\n", result.Rationale))

		sb.WriteString("---\n\n")
	}

	// Footer
	sb.WriteString("## Test Completion\n\n")
	sb.WriteString(fmt.Sprintf("**Test run completed at:** %s\n", endTime.UTC().Format(time.RFC3339)))

	if failed == 0 {
		sb.WriteString(fmt.Sprintf("**Overall result:** ✅ PASS (all %d scenarios passed)\n", passed))
	} else {
		sb.WriteString(fmt.Sprintf("**Overall result:** ❌ FAIL (%d scenario(s) failed)\n", failed))
	}

	return []byte(sb.String())
}
