# Quickstart: Calorie Estimation Bot

**Feature**: 002-calorie-image-estimate
**Status**: Phase 1 Design
**Last Updated**: 2025-12-15

## Overview

This guide will help you set up and run the Telegram calorie estimation bot locally. The bot accepts food images via the `/estimate` command and returns calorie estimates using Google Gemini Vision AI.

---

## Prerequisites

### Required Software

- **Go 1.21+** ([Download](https://go.dev/dl/))
- **Git** (for cloning repository)
- **Text editor** (VS Code, GoLand, vim, etc.)

### Required Credentials

1. **Telegram Bot Token**
   - Open Telegram and search for `@BotFather`
   - Send `/newbot` and follow prompts to create a bot
   - Copy the token (format: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`)
   - Choose a descriptive bot name (e.g., "CalorieScannerBot", "FoodEstimateBot")

2. **Google Gemini API Key**
   - Go to [Google AI Studio](https://ai.google.dev/)
   - Sign in with Google account
   - Click "Get API Key" → "Create API Key"
   - Copy the key (format: `AIzaSy...`)
   - Free tier: 1500 requests/day, 15 requests/minute (sufficient for MVP)

---

## Setup Instructions

### 1. Clone Repository

```bash
git clone <repository-url>
cd telegram-calories-bot
git checkout 002-calorie-image-estimate
```

### 2. Create Environment File

**CRITICAL**: Never commit secrets to Git!

Create `.env` file in project root:

```bash
# .env (DO NOT COMMIT THIS FILE)
TELEGRAM_BOT_TOKEN=your_bot_token_here
GEMINI_API_KEY=your_gemini_api_key_here
```

Verify `.env` is in `.gitignore`:

```bash
cat .gitignore | grep .env
# Should output: .env
```

### 3. Install Dependencies

```bash
# Download Go modules
go mod download

# Install linting tools (required by constitution)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 4. Verify Installation

```bash
# Check Go version (must be 1.21+)
go version

# Check linter installation
golangci-lint --version

# Verify environment variables (should NOT print actual values)
go run -tags=dev src/main.go --check-env
```

---

## Running the Bot

### Development Mode (Local Polling)

```bash
# Run with auto-reload (recommended for development)
go run src/main.go
```

Expected output:

```
[INFO] Starting Calorie Estimation Bot...
[INFO] Bot name: @YourBotName
[INFO] Gemini API key validated
[INFO] Polling for updates...
```

Bot is now running! Open Telegram and:

1. Search for your bot (`@YourBotName`)
2. Send `/estimate`
3. Upload a food image
4. Receive calorie estimate with Re-estimate/Cancel buttons

### Stopping the Bot

Press `Ctrl+C` in terminal. Bot will gracefully shut down.

---

## Testing the Bot

### Manual Testing Flow

**Happy Path** (User Story 1 - P1):

1. Send `/estimate` → Bot responds: "Please upload one food image"
2. Upload food image (JPEG/PNG/WebP) → Bot responds with:
   ```
   🍽️ Calorie Estimate

   Estimated Calories: 650 kcal
   Confidence: High

   Detected Items: Grilled chicken breast, Steamed broccoli, Brown rice

   [Re-estimate] [Cancel]
   ```

**Re-estimation** (User Story 2 - P2):

1. Complete happy path above
2. Click "Re-estimate" button → Bot prompts: "Please upload one food image"
3. Upload new image → Receive new estimate

**Cancellation** (User Story 3 - P3):

1. Send `/estimate` → Bot prompts for image
2. Click "Cancel" → Bot responds: "Estimation cancelled"

### Edge Case Testing

Test these scenarios manually (from spec.md Edge Cases):

| Scenario | Expected Behavior |
|----------|------------------|
| Upload landscape photo | "No food detected in image. Try another photo?" |
| Upload multiple images | "Please send exactly one image (not multiple)" |
| Send text after `/estimate` | Bot ignores or prompts: "Please upload an image" |
| Upload extremely large image | Bot accepts if <20MB (Telegram limit) |
| Upload corrupted image | "Invalid image format. Please upload JPEG/PNG/WebP" |

---

## Running Automated Tests

### Unit Tests

```bash
# Run all unit tests
go test ./tests/unit/... -v

# Run with coverage report (must be ≥70% per constitution)
go test ./tests/unit/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # View coverage report in browser
```

Expected output:

```
=== RUN   TestSessionStateTransitions
--- PASS: TestSessionStateTransitions (0.01s)
=== RUN   TestEstimateResultValidation
--- PASS: TestEstimateResultValidation (0.00s)
=== RUN   TestGeminiAPIResponseParsing
--- PASS: TestGeminiAPIResponseParsing (0.00s)
PASS
coverage: 75.3% of statements
ok      telegram-calories-bot/tests/unit        0.123s
```

### Integration Tests (LLM-Based)

**Prerequisites**: Bot token + Gemini API key in `.env`

```bash
# Run integration tests (simulates user conversations)
go test ./tests/integration/... -v -timeout 60s

# Generate test report (deliverable requirement)
go test ./tests/integration/... -v -json > test-results.json
go run tests/tools/report-generator.go test-results.json > prompts/002-calorie-image-estimate/test-prompts/test-report.md
```

Expected output:

```
=== RUN   TestEstimateFlowHappyPath
[LLM] Simulating user: /estimate
[BOT] Response: Please upload one food image
[LLM] Simulating user: [uploads test-food-1.jpg]
[BOT] Response: 🍽️ Calorie Estimate...
[LLM] ✅ PASS: Bot returned calorie estimate with confidence
--- PASS: TestEstimateFlowHappyPath (5.23s)

=== RUN   TestEstimateFlowMultipleImages
[LLM] Simulating user: /estimate
[LLM] Simulating user: [uploads 3 images]
[BOT] Response: Please send exactly one image (not multiple)
[LLM] ✅ PASS: Bot rejected multiple images (FR-015)
--- PASS: TestEstimateFlowMultipleImages (1.45s)

PASS
All 8 scenarios passed (100% success rate - SC-005 met)
```

---

## Code Quality Checks

### Linting (Required Before Commit)

```bash
# Format code (must run before commit)
gofmt -w .

# Run linter (must pass with zero warnings)
golangci-lint run

# Expected output:
# (no output = all checks passed)
```

If linter fails:

```bash
# View detailed errors
golangci-lint run --verbose

# Auto-fix simple issues
golangci-lint run --fix
```

### Pre-Commit Checklist

Before committing code, verify:

- [ ] `go test ./... -cover` passes with ≥70% coverage
- [ ] `golangci-lint run` passes with zero warnings
- [ ] `gofmt -w .` applied (no formatting changes when re-run)
- [ ] No secrets in code (grep for "AIza" or bot token patterns)
- [ ] All vibe coding prompts archived in `prompts/002-calorie-image-estimate/vibe-coding/`

---

## Deliverables Checklist

Per constitution Principle V and spec.md SC-007/SC-008, before marking feature complete:

### 1. Working Bot ✅
- [ ] Bot runs and responds to `/estimate` command
- [ ] Bot integrates with Gemini Vision API
- [ ] Inline buttons (Re-estimate, Cancel) work
- [ ] Bot has descriptive name (e.g., "CalorieScannerBot")

### 2. Clean Code ✅
- [ ] All code passes `gofmt` formatting
- [ ] `golangci-lint run` passes with zero warnings
- [ ] Code follows project structure (src/handlers, src/services, src/models)
- [ ] Comments explain non-obvious logic

### 3. Vibe Coding Prompts ✅
- [ ] All prompts used to generate code archived in:
  ```
  prompts/002-calorie-image-estimate/vibe-coding/
  ├── 001-initial-bot-structure.md
  ├── 002-gemini-integration.md
  ├── 003-session-management.md
  └── [additional prompts as generated]
  ```

### 4. LLM Test Prompts ✅
- [ ] Test prompts archived in:
  ```
  prompts/002-calorie-image-estimate/test-prompts/
  ├── 001-happy-path-test.md
  ├── 002-edge-cases-test.md
  └── test-report.md
  ```

### 5. Test Report ✅
- [ ] Test report generated with pass/fail results
- [ ] Report shows 100% pass rate for defined scenarios (SC-005)
- [ ] Report includes:
  - Total scenarios tested (8 expected)
  - Pass/fail status for each scenario
  - Execution time
  - Any failures with error details

---

## Troubleshooting

### Bot Not Starting

**Error**: `panic: API key not valid`

**Solution**: Verify `GEMINI_API_KEY` in `.env` is correct. Test with:

```bash
curl -H "x-goog-api-key: $GEMINI_API_KEY" \
  https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash
# Should return model info, not 403 error
```

---

**Error**: `telegram: unauthorized (401)`

**Solution**: Verify `TELEGRAM_BOT_TOKEN` in `.env` is correct. Get new token from @BotFather if needed.

---

### Bot Not Responding to Commands

**Issue**: Send `/estimate` but bot doesn't respond

**Debug Steps**:

1. Check bot logs for errors
2. Verify bot is running (check terminal output)
3. Test with simple `/start` command
4. Check BotFather settings: `/mybots` → Select bot → Bot Settings → Group Privacy (should be OFF for private chats)

---

### Gemini API Errors

**Error**: `RESOURCE_EXHAUSTED (429)`

**Solution**: Hit rate limit (15 req/min). Wait 1 minute and retry. For production, implement exponential backoff.

---

**Error**: `INVALID_ARGUMENT (400)` on image upload

**Solution**: Image format not supported or corrupted. Verify:

```bash
file your-image.jpg
# Should output: JPEG image data, ... (not "data" or unknown)
```

---

### Tests Failing

**Issue**: Unit tests pass but integration tests fail

**Debug Steps**:

1. Verify `.env` file exists with valid credentials
2. Check network connection (integration tests call real APIs)
3. Run single test in verbose mode:
   ```bash
   go test -v ./tests/integration/... -run TestEstimateFlowHappyPath
   ```
4. Check Gemini API quota: [Google AI Studio Dashboard](https://ai.google.dev/)

---

## Performance Benchmarks

Per spec.md Success Criteria, monitor these metrics during testing:

| Metric | Target | Measurement |
|--------|--------|-------------|
| `/estimate` response time | <1 second (SC-003) | Time from command send to prompt received |
| Image analysis + response | <8 seconds (SC-004) | Time from image upload to result received |
| Full flow (command → result) | <10 seconds (SC-001) | End-to-end user experience |

**How to Measure**:

```bash
# Add timestamps to bot logs
# In src/main.go:
log.Printf("[%s] /estimate received", time.Now().Format(time.RFC3339))
log.Printf("[%s] Gemini API call started", time.Now().Format(time.RFC3339))
log.Printf("[%s] Result sent to user", time.Now().Format(time.RFC3339))
```

---

## Next Steps

After completing this quickstart:

1. **Read**: [data-model.md](./data-model.md) - Understand session state machine and data structures
2. **Read**: [contracts/gemini-vision.yaml](./contracts/gemini-vision.yaml) - API request/response format
3. **Implement**: Follow tasks in [tasks.md](./tasks.md) (generated via `/speckit.tasks`)
4. **Deploy** (out of scope for MVP): See deployment guide (future feature)

---

## Support & Resources

- **Spec**: [spec.md](./spec.md) - Full feature specification with requirements
- **Plan**: [plan.md](./plan.md) - Implementation plan and architecture
- **Research**: [research.md](./research.md) - Technical decisions and alternatives
- **Telegram Bot API**: https://core.telegram.org/bots/api
- **Telebot v3 Docs**: https://github.com/tucnak/telebot
- **Gemini API Docs**: https://ai.google.dev/gemini-api/docs

---

## FAQ

**Q: Can I use a different LLM instead of Gemini?**

A: No. Per user constraints and spec.md assumptions, only Google Gemini Vision via ai.dev is permitted for this feature.

---

**Q: Where are images stored?**

A: Nowhere. Per FR-010, FR-011, FR-012, the bot operates statelessly with no persistence. Images are held in memory only during Gemini API calls.

---

**Q: Can I add history/logging of past estimates?**

A: No. Persistence is explicitly out of scope per spec requirements. This feature is stateless by design.

---

**Q: How do I change the bot's response format?**

A: Edit the `FormatResult()` function in `src/models/estimate.go`. Maintain the fixed format per FR-006 (must be deterministic).

---

**Q: Tests pass but coverage is <70%, what now?**

A: Add more unit tests for uncovered code paths. Focus on:
- Error handling branches
- Edge cases (multiple images, no food, invalid format)
- State machine transitions

Use `go tool cover -html=coverage.out` to identify uncovered lines.

---

**Q: How do I contribute code?**

A: Follow constitution Pre-Implementation Checklist:
1. User scenarios written in spec.md ✅ (already done)
2. Acceptance criteria defined ✅ (already done)
3. Constitution compliance verified ✅ (see plan.md)
4. Technical approach documented ✅ (research.md)
5. Test stubs created and failing (write tests first!)
6. Implement feature (make tests pass)
7. Archive prompts in `prompts/002-calorie-image-estimate/vibe-coding/`
8. Run linting and coverage checks
9. Submit PR with constitution compliance checklist

---

**Last Updated**: 2025-12-15
**Version**: 1.0.0
**Feature**: 002-calorie-image-estimate
