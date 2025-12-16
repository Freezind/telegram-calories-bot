# Implementation Plan: Calorie Estimation from Image

**Branch**: `002-calorie-image-estimate` | **Date**: 2025-12-15 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-calorie-image-estimate/spec.md`

## Summary

Implement a minimal Telegram bot that accepts a single food image upload via the /estimate command, analyzes it using Google Gemini Vision API, and returns a calorie estimate with confidence indicator. The bot operates statelessly with no persistence, provides inline buttons (Re-estimate, Cancel) for user control, and includes comprehensive testing via LLM-based tools with archived prompts and test reports.

## Technical Context

**Language/Version**: Golang 1.21+
**Primary Dependencies**: github.com/tucnak/telebot (v3), Google Gemini SDK for Go
**Storage**: N/A (stateless operation, in-memory session state only)
**Testing**: Go testing package, LLM-based conversational testing tool
**Target Platform**: Single server deployment (local or cloud)
**Project Type**: Single (CLI-style bot executable)
**Performance Goals**: /estimate response <1s, image analysis + response <8s, full flow <10s
**Constraints**: Stateless (no persistence), single image per request, reject multiple uploads
**Scale/Scope**: MVP - single bot command (/estimate), 3 user stories, minimal inline button interaction

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Quality-Controlled Vibe Coding
- ✅ Feature has clear acceptance criteria in spec.md (3 user stories with Given-When-Then scenarios)
- ✅ Code review gates planned (constitution mandates review before merge)
- ✅ Linting required: gofmt, golangci-lint
- ✅ Prompt archiving required: All vibe coding prompts → `prompts/`
- ✅ Test coverage required: 70% minimum for services

### Test-First with LLM Validation
- ✅ Acceptance tests required based on user scenarios
- ✅ LLM-based testing tool required for bot interaction flows
- ✅ Test reports required as deliverables (SC-008)
- ✅ Red-Green-Refactor workflow documented in constitution

### Dual Interface Architecture
- ✅ **NOT APPLICABLE** - This feature implements LUI only (slash commands + inline buttons)
- ✅ Mini App explicitly out of scope per spec requirements
- ✅ No architecture violation - feature is scoped to LUI interface only

### Fixed Technology Stack
- ✅ Golang 1.21+ confirmed
- ✅ github.com/tucnak/telebot (v3) confirmed
- ⚠️ **NEEDS CLARIFICATION**: Google Gemini SDK for Go - need to research official SDK vs REST client

### Deliverable-Driven Development
- ✅ Working Bot: /estimate command with Gemini Vision integration
- ✅ Clean Code: Go code with linting and review gates
- ✅ Vibe Coding Prompts: Archive in `prompts/002-calorie-image-estimate/`
- ✅ LLM Test Prompts: Archive in `prompts/002-calorie-image-estimate/test-prompts/`
- ✅ Test Reports: Generate and export test report (pass/fail)

**Gate Result**: ⚠️ CONDITIONAL PASS - Proceed to Phase 0 to resolve NEEDS CLARIFICATION items

## Project Structure

### Documentation (this feature)

```text
specs/002-calorie-image-estimate/
├── plan.md              # This file
├── research.md          # Phase 0: Gemini SDK research, LLM test tool selection
├── data-model.md        # Phase 1: Session state model (in-memory only)
├── quickstart.md        # Phase 1: Setup instructions (bot token, Gemini key, run)
├── contracts/           # Phase 1: Gemini API request/response contracts
│   └── gemini-vision.yaml
└── tasks.md             # Phase 2: NOT created by this command
```

### Source Code (repository root)

```text
src/
├── handlers/
│   └── estimate.go      # /estimate command handler, Re-estimate/Cancel button handlers
├── services/
│   ├── gemini.go        # Gemini Vision API client (calorie estimation)
│   └── session.go       # In-memory session state manager (user flow tracking)
├── models/
│   └── estimate.go      # EstimateResult, SessionState structs
└── main.go              # Bot entry point, env var loading, handler registration

tests/
├── integration/
│   └── estimate_flow_test.go  # LLM-based conversational flow tests
└── unit/
    ├── gemini_test.go         # Mock Gemini API responses
    └── session_test.go        # Session state transitions

prompts/002-calorie-image-estimate/
├── vibe-coding/               # All prompts used to generate code
│   ├── 001-initial-bot-structure.md
│   ├── 002-gemini-integration.md
│   └── [additional prompts as generated]
└── test-prompts/              # LLM test prompts
    ├── 001-happy-path-test.md
    ├── 002-edge-cases-test.md
    └── test-report.md         # Exported test results (pass/fail)
```

**Structure Decision**: Single project structure (Option 1) selected. This is a simple bot with one command, no web frontend, no persistence layer. All logic fits in `src/` with standard Go package layout (handlers, services, models). Test separation (unit vs integration) follows constitution requirements.

## Complexity Tracking

> This section is empty because there are no constitution violations requiring justification.

**Rationale**: Feature is intentionally minimal (MVP scope), uses mandated technology stack (Go + telebot), follows test-first workflow, and produces all required deliverables. No complexity additions beyond constitution baseline.

## Phase 0: Outline & Research

### Research Tasks

Based on Technical Context NEEDS CLARIFICATION items, Phase 0 must resolve:

1. **Google Gemini SDK for Go**
   - **Question**: Does Google provide an official Gemini SDK for Go, or should we use REST API with standard HTTP client?
   - **Research Target**: Google AI Studio documentation, google-generativeai Go packages
   - **Decision Criteria**: Prefer official SDK if stable; fallback to REST if no Go SDK exists
   - **Output**: Document chosen approach in research.md with installation/setup instructions

2. **LLM-Based Testing Tool Selection**
   - **Question**: Which LLM-based testing tool should validate conversational flow (simulate /estimate → image → result)?
   - **Research Target**: Available tools for testing chat bots (e.g., Botium, custom GPT-4 test script, Claude-based test harness)
   - **Decision Criteria**: Must simulate Telegram user interactions, support image uploads in test scenarios, export pass/fail reports
   - **Output**: Document selected tool in research.md with example test case

3. **Gemini Vision Prompt Engineering**
   - **Question**: What prompt structure yields best calorie estimation results from Gemini Vision?
   - **Research Target**: Gemini Vision API documentation, prompt templates for food recognition + calorie estimation
   - **Decision Criteria**: Deterministic output format, confidence scoring, handling non-food images
   - **Output**: Document prompt template in research.md with example input/output

4. **Telegram Image Handling with telebot v3**
   - **Question**: How to download images from Telegram via telebot, detect multiple uploads, and pass to Gemini API?
   - **Research Target**: telebot v3 documentation for Photo handling, GetFile API
   - **Decision Criteria**: Must support JPEG/PNG/WebP, detect multiple images in single message
   - **Output**: Document code patterns in research.md with snippet examples

5. **Session State Management (In-Memory)**
   - **Question**: How to track user flow state (awaiting image, processing, showing results) without persistence?
   - **Research Target**: Go concurrency patterns (sync.Map, context.Context), telebot session patterns
   - **Decision Criteria**: Thread-safe, simple cleanup, no external dependencies
   - **Output**: Document state machine design in research.md

### Expected Outputs

**Artifact**: `specs/002-calorie-image-estimate/research.md`

**Contents**:
- Decision: Gemini SDK approach (official SDK vs REST client)
- Decision: LLM testing tool (specific tool + setup steps)
- Decision: Gemini Vision prompt template (with examples)
- Decision: Image download pattern (telebot code snippet)
- Decision: Session state pattern (Go code snippet)
- Rationale for each decision
- Alternatives considered and why rejected

## Phase 1: Design & Contracts

### Prerequisites
- `research.md` complete with all NEEDS CLARIFICATION items resolved

### Data Model

**Artifact**: `specs/002-calorie-image-estimate/data-model.md`

**Contents** (based on spec Key Entities):

```go
// SessionState tracks user flow for /estimate command (in-memory only)
type SessionState string

const (
    StateIdle           SessionState = "idle"
    StateAwaitingImage  SessionState = "awaiting_image"
    StateProcessing     SessionState = "processing"
)

// EstimateResult holds calorie analysis output
type EstimateResult struct {
    Calories   int     // Estimated total calories
    Confidence string  // "low", "medium", "high" or percentage (TBD in research)
    FoodItems  []string // Detected food items (optional)
}

// UserSession represents in-memory session (NOT persisted)
type UserSession struct {
    UserID    int64        // Telegram user ID
    State     SessionState
    MessageID int          // Last bot message for editing
}
```

**Validation Rules** (from functional requirements):
- FR-003: Accept JPEG, PNG, WebP only
- FR-015: Reject if multiple images detected in single message
- FR-014: Handle "no food detected" case (return error or low-confidence result)

**State Transitions**:
- Idle → AwaitingImage (on /estimate command)
- AwaitingImage → Processing (on image upload)
- Processing → Idle (after result sent)
- AwaitingImage → Idle (on Cancel button)
- Any state → AwaitingImage (on Re-estimate button)

### API Contracts

**Artifact**: `specs/002-calorie-image-estimate/contracts/gemini-vision.yaml`

**Contents**: OpenAPI-style contract for Gemini Vision API interaction

```yaml
# Gemini Vision API Contract (simplified)
endpoint: POST https://generativelanguage.googleapis.com/v1/models/gemini-pro-vision:generateContent
authentication: API Key (header: x-goog-api-key from env GEMINI_API_KEY)

request:
  contents:
    - parts:
        - text: "[Prompt template from research.md]"
        - inline_data:
            mime_type: "image/jpeg" # or image/png, image/webp
            data: "[base64 encoded image]"

response:
  candidates:
    - content:
        parts:
          - text: "Estimated Calories: 450 kcal | Confidence: High"
      finishReason: "STOP"

error_cases:
  - RESOURCE_EXHAUSTED: Rate limit exceeded (return error to user: "Service busy, try again")
  - INVALID_ARGUMENT: Invalid image format (return error: "Invalid image, please upload JPEG/PNG/WebP")
  - PERMISSION_DENIED: Invalid API key (fail fast on startup)
```

### Quickstart Guide

**Artifact**: `specs/002-calorie-image-estimate/quickstart.md`

**Contents**:

```markdown
# Quickstart: Calorie Estimation Bot

## Prerequisites
- Go 1.21+
- Telegram Bot Token (from @BotFather)
- Google Gemini API Key (from ai.google.dev)

## Setup

1. Clone repository and navigate to project root
2. Create `.env` file (NEVER commit this):
   ```
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   GEMINI_API_KEY=your_gemini_key_here
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run linting (must pass before running):
   ```bash
   gofmt -w .
   golangci-lint run
   ```

## Running the Bot

```bash
go run src/main.go
```

Bot will start polling. Test with:
1. Send `/estimate` to bot
2. Upload a food image
3. Receive calorie estimate with Re-estimate/Cancel buttons

## Running Tests

```bash
# Unit tests
go test ./tests/unit/... -v

# Integration tests (requires bot token + Gemini key in .env)
go test ./tests/integration/... -v
```

## Deliverables Checklist

Before marking feature complete, verify:
- [ ] Bot runs and responds to /estimate
- [ ] Code passes gofmt and golangci-lint
- [ ] All prompts archived in `prompts/002-calorie-image-estimate/`
- [ ] Test report exported in `prompts/002-calorie-image-estimate/test-prompts/test-report.md`
- [ ] Tests pass with 70%+ coverage
```

### Agent Context Update

After generating design artifacts, run:

```bash
.specify/scripts/bash/update-agent-context.sh claude
```

This will update `CLAUDE.md` with:
- Golang 1.21+ (already present from 001-hello-bot-mvp)
- Google Gemini SDK/client (new addition)
- LLM testing tool (new addition)

## Post-Phase 1 Constitution Re-Check

### Quality-Controlled Vibe Coding
- ✅ No changes from initial check

### Test-First with LLM Validation
- ✅ Integration test structure defined in data-model.md
- ✅ LLM testing tool selected in research.md

### Dual Interface Architecture
- ✅ Not applicable (LUI only)

### Fixed Technology Stack
- ✅ Gemini SDK/client approach documented in research.md
- ✅ No new frameworks added (telebot already approved)

### Deliverable-Driven Development
- ✅ All deliverables mapped to artifacts:
  - Working Bot → quickstart.md run instructions
  - Clean Code → linting requirements in quickstart.md
  - Vibe Coding Prompts → prompts/002-calorie-image-estimate/vibe-coding/
  - LLM Test Prompts → prompts/002-calorie-image-estimate/test-prompts/
  - Test Reports → test-report.md

**Gate Result**: ✅ PASS - Ready for Phase 2 (task generation via `/speckit.tasks`)

## Scope Enforcement

### Features Explicitly OUT of Scope

Per user input constraints, the following are **FORBIDDEN** and must be rejected during implementation:

1. **Storage/Persistence**:
   - ❌ No database (PostgreSQL, SQLite, etc.)
   - ❌ No key-value stores (Redis, Cloudflare KV)
   - ❌ No file-based history/logs
   - ❌ No Cloudflare Workers storage (DO/D1/R2)

2. **GUI/Mini App**:
   - ❌ No Telegram Mini App
   - ❌ No web interface for history viewing
   - ❌ No CRUD operations for historical data

3. **Infrastructure**:
   - ❌ No deployment automation (CI/CD pipelines)
   - ❌ No webhook infrastructure beyond basic telebot polling
   - ❌ No load balancers, orchestration, etc.

4. **Alternative Vision Providers**:
   - ❌ No OpenAI GPT-4 Vision
   - ❌ No Anthropic Claude Vision
   - ❌ No other LLM vision services
   - ✅ **ONLY** Google Gemini Vision via ai.dev

5. **Abstractions/Over-Engineering**:
   - ❌ No repository pattern (no persistence layer needed)
   - ❌ No dependency injection frameworks
   - ❌ No plugin systems or extensibility hooks
   - ✅ **ONLY** minimal code to satisfy functional requirements

### Validation Checkpoints

During implementation (Phase 2), if any task introduces:
- New storage dependency → **REJECT** (violates FR-010, FR-011, FR-012)
- Non-Gemini vision API → **REJECT** (violates user constraint)
- Abstraction layers not required for MVP → **FLAG** for review

## Next Steps

This planning command stops here. To proceed:

1. **Resolve Research** (if running manually):
   ```bash
   # Research tasks will be executed automatically as part of Phase 0
   ```

2. **Generate Tasks**:
   ```bash
   /speckit.tasks
   ```
   This will create `tasks.md` with executable implementation steps based on this plan.

3. **Begin Implementation**:
   ```bash
   /speckit.implement
   ```
   This will execute tasks in dependency order, following test-first workflow.

## Plan Metadata

**Created**: 2025-12-15
**Status**: Ready for Phase 0 Research
**Blocking Issues**: None (NEEDS CLARIFICATION items moved to Phase 0 research tasks)
**Constitution Version**: 1.0.0
