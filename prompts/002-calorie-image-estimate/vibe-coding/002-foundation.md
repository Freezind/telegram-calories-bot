# Vibe Coding Prompts - Phase 2: Foundational

**Feature**: 002-calorie-image-estimate - Calorie Estimation from Image
**Phase**: Foundational (Blocking Prerequisites)
**Date**: 2025-12-16

## Context

Building core infrastructure for Telegram bot: data models, services (session management, Gemini integration), and bot entry point with validation.

## Prompts Used

### 1. Data Models (T010-T012)

**Prompt**: "Create Go data models for:
1. SessionState enum (Idle, AwaitingImage, Processing)
2. UserSession struct (UserID, State, LastActivity, MessageID)
3. EstimateResult struct (Calories, Confidence, FoodItems, Reasoning)
Include validation methods per data-model.md specification"

**Result**: Created `src/models/estimate.go` with all three types plus:
- `FormatResult()` function for fixed-format response (FR-006)
- `HasFood()` method to detect empty results (FR-014)
- `Validate()` method for confidence and calorie validation

### 2. Session Manager (T013, T018)

**Prompt**: "Create SessionManager service using sync.Map for thread-safe in-memory session storage. Implement:
- GetSession (create if missing)
- UpdateSession (state transitions)
- DeleteSession (on Cancel/completion)
- CleanupStale (remove sessions inactive >15 min)
- StartCleanupRoutine (goroutine running every 5 min)"

**Result**: Created `src/services/session.go` with full state management per data-model.md state machine

### 3. Gemini Client (T014)

**Prompt**: "Create GeminiClient service wrapping google.golang.org/genai SDK. Implement EstimateCalories method that:
1. Creates client with API key from GEMINI_API_KEY env var
2. Uses gemini-2.5-flash model with temperature 0.2
3. Sends structured JSON prompt + image bytes
4. Parses JSON response to EstimateResult
5. Validates result before returning"

**Result**: Created `src/services/gemini.go` per research.md Decision 1 & 3

**API Corrections Made**:
- Used `genai.NewClient(ctx, &genai.ClientConfig{...})` pattern
- Used `genai.NewPartFromText()` and `genai.NewPartFromBytes()`
- Used `client.Models.GenerateContent()` method
- Client doesn't have Close() method (removed defer)
- Response Part.Text is string not *string (fixed type assertion)

### 4. Bot Entry Point (T016-T017)

**Prompt**: "Create main.go that:
1. Validates TELEGRAM_BOT_TOKEN and GEMINI_API_KEY env vars (fail fast if missing)
2. Initializes SessionManager and GeminiClient
3. Starts session cleanup goroutine
4. Initializes telebot with LongPoller
5. Logs startup progress
6. Starts bot polling (handlers to be added in Phase 3)"

**Result**: Created `src/main.go` with environment validation and service initialization

## Technical Decisions

1. **Module Path**: Used `github.com/freezind/telegram-calories-bot` (existing module name)
2. **Gemini SDK**: google.golang.org/genai v1.39.0 (required Go 1.24+)
3. **Response Parsing**: Added markdown code block stripping (```json) for robustness
4. **Error Messages**: Included response text in parse errors for debugging

## Build Verification

```bash
cd src && go build -o /tmp/bot-test .
# BUILD SUCCESS
```

All Phase 2 code compiles with no errors.

## Files Created

- `src/models/estimate.go` (3 types + 3 methods)
- `src/services/session.go` (SessionManager)
- `src/services/gemini.go` (GeminiClient)
- `src/main.go` (bot entry point)

## Next Phase

Phase 3: User Story 1 - MVP Implementation
- Unit tests (5 test files)
- Handlers for /estimate command, photo upload, inline buttons
- Integration tests (4 test scenarios)
- Linting and deliverables

**Estimated effort**: Phase 3 is the largest phase (23 tasks) - implements full MVP functionality.
