package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config holds test configuration from environment variables
type Config struct {
	// Mini App configuration
	MiniAppURL string

	// Test assets
	TestImagePath string

	// LLM judge
	GeminiAPIKey string

	// Optional
	TestTimeout time.Duration
}

func main() {
	log.Println("========================================")
	log.Println("LLM-Based Bot Testing Suite")
	log.Println("========================================")

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Validate test image exists
	if _, err := os.Stat(config.TestImagePath); os.IsNotExist(err) {
		log.Fatalf("Test image not found: %s", config.TestImagePath)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.TestTimeout)
	defer cancel()

	// Initialize components
	judge := NewGeminiJudge(config.GeminiAPIKey)
	botTester, err := NewBotTester(judge)
	if err != nil {
		log.Fatalf("Failed to create bot tester: %v", err)
	}
	miniappTester := NewMiniAppTester(config.MiniAppURL, judge)
	reporter := NewReporter()

	log.Println("Test configuration loaded successfully")
	log.Printf("  Mini App URL: %s", config.MiniAppURL)
	log.Printf("  Test Image: %s", config.TestImagePath)
	log.Printf("  Timeout: %s", config.TestTimeout)
	log.Println()

	// Track test start time
	testStart := time.Now()
	var results []TestResult

	// Execute scenarios sequentially
	log.Println("Starting test execution...")
	log.Println()

	// Scenario 1: /start welcome message
	log.Println("[S1] Testing /start command...")
	s1Result := botTester.TestStart(ctx)
	results = append(results, *s1Result)
	logScenarioResult(s1Result)

	// Scenario 2: /estimate + image upload
	log.Println("[S2] Testing /estimate + image upload...")
	s2Result := botTester.TestEstimate(ctx, config.TestImagePath)
	results = append(results, *s2Result)
	logScenarioResult(s2Result)

	// Scenario 3: Re-estimate button preservation
	log.Println("[S3] Testing Re-estimate button...")
	s3Result := botTester.TestReEstimate(ctx)
	results = append(results, *s3Result)
	logScenarioResult(s3Result)

	// Scenario 4: Cancel button preservation
	log.Println("[S4] Testing Cancel button...")
	s4Result := botTester.TestCancel(ctx)
	results = append(results, *s4Result)
	logScenarioResult(s4Result)

	// Scenario 5: Mini App page load
	log.Println("[S5] Testing Mini App page load...")
	s5Result := miniappTester.TestPageLoad(ctx)
	results = append(results, *s5Result)
	logScenarioResult(s5Result)

	// Track test end time
	testEnd := time.Now()

	log.Println()
	log.Println("========================================")
	log.Println("Test Execution Complete")
	log.Println("========================================")

	// Generate report
	log.Println("Generating test report...")
	reportContent := reporter.Generate(results, testStart, testEnd)

	// Ensure reports directory exists
	reportsDir := "reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		log.Fatalf("Failed to create reports directory: %v", err)
	}

	// Write report to file
	reportPath := filepath.Join(reportsDir, "004-test-report.md")
	if err := os.WriteFile(reportPath, reportContent, 0644); err != nil {
		log.Fatalf("Failed to write test report: %v", err)
	}

	log.Printf("Test report generated: %s", reportPath)

	// Archive LLM judge prompts
	log.Println("Archiving LLM judge prompts...")
	promptsPath := filepath.Join(reportsDir, "llm-judge-prompts.md")
	if err := archivePrompts(judge.GetPrompts(), promptsPath); err != nil {
		log.Fatalf("Failed to archive prompts: %v", err)
	}

	log.Printf("LLM judge prompts archived to %s", promptsPath)

	// Print summary
	log.Println()
	log.Println("========================================")
	log.Println("Test Summary")
	log.Println("========================================")

	passed := 0
	failed := 0
	for _, result := range results {
		if result.Verdict == "PASS" {
			passed++
		} else {
			failed++
		}
	}

	log.Printf("Total scenarios: %d", len(results))
	log.Printf("Passed: %d", passed)
	log.Printf("Failed: %d", failed)
	log.Printf("Duration: %s", testEnd.Sub(testStart).Round(time.Second))

	// Exit with appropriate code
	if failed > 0 {
		log.Println()
		log.Println("Result: FAIL (one or more scenarios failed)")
		os.Exit(1)
	}

	log.Println()
	log.Println("Result: PASS (all scenarios passed)")
	os.Exit(0)
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	config := &Config{}

	// Mini App configuration
	config.MiniAppURL = os.Getenv("MINIAPP_URL")
	if config.MiniAppURL == "" {
		return nil, fmt.Errorf("MINIAPP_URL environment variable is required")
	}

	// Test assets
	config.TestImagePath = os.Getenv("TEST_FOOD_IMAGE_PATH")
	if config.TestImagePath == "" {
		config.TestImagePath = "tests/fixtures/food.jpg" // Default
	}

	// LLM judge
	config.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	if config.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	// Optional timeout
	timeoutStr := os.Getenv("TEST_TIMEOUT_SECONDS")
	if timeoutStr == "" {
		config.TestTimeout = 120 * time.Second // Default 2 minutes
	} else {
		timeoutSec, err := strconv.Atoi(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid TEST_TIMEOUT_SECONDS: %w", err)
		}
		config.TestTimeout = time.Duration(timeoutSec) * time.Second
	}

	return config, nil
}

// logScenarioResult logs the result of a scenario
func logScenarioResult(result *TestResult) {
	if result.Verdict == "PASS" {
		log.Printf("  ✅ PASS - %s", result.Rationale)
	} else {
		log.Printf("  ❌ FAIL - %s", result.Rationale)
		if result.Error != "" {
			log.Printf("     Error: %s", result.Error)
		}
	}
	log.Println()
}

// archivePrompts appends all judge prompts to the specified file
func archivePrompts(prompts []string, filepath string) error {
	if len(prompts) == 0 {
		return nil
	}

	// Open file in append mode
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", filepath, err)
	}
	defer f.Close()

	// Add header for this test run
	header := fmt.Sprintf("\n\n## LLM-Based Bot Testing - %s\n\n", time.Now().Format(time.RFC3339))
	if _, err := f.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each prompt
	for i, prompt := range prompts {
		section := fmt.Sprintf("### Scenario %d Judge Prompt\n\n```\n%s\n```\n\n", i+1, prompt)
		if _, err := f.WriteString(section); err != nil {
			return fmt.Errorf("failed to write prompt %d: %w", i+1, err)
		}
	}

	return nil
}

// Environment variable documentation (for reference)
const envDocs = `
Required Environment Variables:
  MINIAPP_URL                - Deployed Mini App HTTPS URL
  GEMINI_API_KEY             - Gemini API key for LLM judge

Optional Environment Variables:
  TEST_FOOD_IMAGE_PATH       - Path to test food image (default: tests/fixtures/food.jpg)
  TEST_TIMEOUT_SECONDS       - Overall test timeout in seconds (default: 120)

Example:
  export MINIAPP_URL="https://your-app.railway.app"
  export GEMINI_API_KEY="your-gemini-api-key"

Note: You can copy .env.test.example to .env.test and fill in the values,
      then run: source .env.test && go run cmd/tester/*.go
`
