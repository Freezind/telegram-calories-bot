# Tasks: Calorie Estimation from Image

**Feature Branch**: `002-calorie-image-estimate`
**Input**: Design documents from `/specs/002-calorie-image-estimate/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/gemini-vision.yaml ✅

**Tests**: This feature explicitly requires test-first development per constitution Principle II and spec.md SC-005, SC-006, SC-007, SC-008. All test tasks are REQUIRED.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. Follow test-first workflow (Red-Green-Refactor).

## Format: `- [ ] [ID] [P?] [Story?] Description`

- **Checkbox**: Always start with `- [ ]`
- **[ID]**: Sequential task number (T001, T002, etc.)
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story label (US1, US2, US3) - only for story-specific tasks
- **File paths**: Always include exact file paths in descriptions

## Path Conventions

Per plan.md, this is a **single project** structure:
- Source code: `src/` at repository root
- Tests: `tests/` at repository root
- Prompts: `prompts/002-calorie-image-estimate/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure. These tasks set up the Go project and dependencies.

- [x] T001 Initialize Go module in repository root: `go mod init telegram-calories-bot`
- [x] T002 Create project directory structure per plan.md: `src/{handlers,services,models}`, `tests/{integration,unit}`, `prompts/002-calorie-image-estimate/{vibe-coding,test-prompts}`
- [x] T003 [P] Add telebot v3 dependency: `go get gopkg.in/telebot.v3`
- [x] T004 [P] Add Google Gemini SDK dependency: `go get google.golang.org/genai`
- [x] T005 [P] Add testing dependencies: `go get github.com/stretchr/testify`
- [x] T006 [P] Create .env.example file with TELEGRAM_BOT_TOKEN and GEMINI_API_KEY placeholders in repository root
- [x] T007 [P] Add .env to .gitignore (if not already present) in repository root
- [x] T008 [P] Install golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- [x] T009 Create .golangci.yml linter configuration file in repository root

**Deliverable Checkpoint**: Archive all vibe coding prompts for Phase 1 in `prompts/002-calorie-image-estimate/vibe-coding/001-setup.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

### Foundation: Data Models (All Stories Depend On These)

- [x] T010 [P] Create SessionState enum in src/models/estimate.go (defines Idle, AwaitingImage, Processing states)
- [x] T011 [P] Create UserSession struct in src/models/estimate.go (tracks user flow state)
- [x] T012 [P] Create EstimateResult struct in src/models/estimate.go (holds calorie estimate output)

### Foundation: Core Services

- [x] T013 Create SessionManager service in src/services/session.go (manages sync.Map for user sessions, implements state transitions per data-model.md)
- [x] T014 [P] Create GeminiClient service in src/services/gemini.go (wraps Google Gemini SDK, handles API calls per contracts/gemini-vision.yaml)
- [x] T015 [P] Implement FormatResult function in src/models/estimate.go (fixed-format response per FR-006)

### Foundation: Bot Entry Point

- [x] T016 Create main.go in src/ (load env vars, initialize telebot, register handlers, start polling)
- [x] T017 Add environment variable validation in src/main.go (fail fast if TELEGRAM_BOT_TOKEN or GEMINI_API_KEY missing per contracts)
- [x] T018 [P] Implement session cleanup goroutine in src/services/session.go (cleanup stale sessions every 5 minutes per data-model.md)

**Deliverable Checkpoint**: Archive all vibe coding prompts for Phase 2 in `prompts/002-calorie-image-estimate/vibe-coding/002-foundation.md`

**Foundational Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Single Image Calorie Estimation (Priority: P1) 🎯 MVP

**Goal**: Implement core /estimate command flow: user sends command → uploads image → receives calorie estimate with inline buttons.

**Independent Test**: Send /estimate, upload food image, receive calorie estimate with confidence indicator and inline buttons (Re-estimate, Cancel).

**Acceptance Criteria** (from spec.md):
1. Bot responds to /estimate with prompt to upload one image
2. Bot analyzes uploaded food image and returns calorie estimate in fixed format
3. Result includes inline buttons: Re-estimate and Cancel

### Unit Tests for User Story 1 (Write FIRST - Red State)

> **RED**: Write these tests FIRST, ensure they FAIL before implementation

- [ ] T019 [P] [US1] Unit test for EstimateResult validation in tests/unit/estimate_test.go (test calories >0, confidence in {low,medium,high}, empty items for no food)
- [ ] T020 [P] [US1] Unit test for SessionManager state transitions in tests/unit/session_test.go (test Idle→AwaitingImage, AwaitingImage→Processing, Processing→Idle)
- [ ] T021 [P] [US1] Unit test for GeminiClient JSON parsing in tests/unit/gemini_test.go (mock Gemini API responses, test success + no-food + error cases)
- [ ] T022 [P] [US1] Unit test for image format validation in tests/unit/handlers_test.go (test JPEG/PNG/WebP accepted, others rejected per FR-003)
- [ ] T023 [P] [US1] Unit test for multiple image rejection in tests/unit/handlers_test.go (test FR-015: reject all images when multiple uploaded)

### Implementation for User Story 1 (Make Tests GREEN)

> **GREEN**: Implement minimum code to make tests pass

- [ ] T024 [US1] Implement /estimate command handler in src/handlers/estimate.go (transition Idle→AwaitingImage, send prompt "Please upload one food image")
- [ ] T025 [US1] Implement OnPhoto handler in src/handlers/estimate.go (validate state=AwaitingImage, check single image, download image)
- [ ] T026 [US1] Add image format validation in src/handlers/estimate.go (check MIME type is image/jpeg|png|webp per FR-003)
- [ ] T027 [US1] Add multiple image detection in src/handlers/estimate.go (check len(msg.Album) > 1, reject all with FR-015 message)
- [ ] T028 [US1] Integrate GeminiClient in src/handlers/estimate.go (call EstimateCalories with image bytes, handle response)
- [ ] T029 [US1] Create inline keyboard with Re-estimate and Cancel buttons in src/handlers/estimate.go (per FR-007)
- [ ] T030 [US1] Format and send calorie estimate result in src/handlers/estimate.go (use FormatResult, include inline buttons, transition Processing→Idle)
- [ ] T031 [US1] Add error handling for invalid images in src/handlers/estimate.go (corrupted file, invalid format → FR-013 message)
- [ ] T032 [US1] Add error handling for no food detected in src/handlers/estimate.go (empty FoodItems → FR-014 message)
- [ ] T033 [US1] Add error handling for Gemini API errors in src/handlers/estimate.go (timeout, rate limit → clear session, go Idle)

### Integration Tests for User Story 1 (Verify End-to-End)

> **Integration**: Test full user journey with real bot and Gemini API

- [ ] T034 [US1] Integration test for happy path in tests/integration/estimate_flow_test.go (simulate: /estimate → upload test image → verify result format + buttons)
- [ ] T035 [US1] Integration test for non-food image in tests/integration/estimate_flow_test.go (simulate: /estimate → upload landscape → verify FR-014 error)
- [ ] T036 [US1] Integration test for multiple images in tests/integration/estimate_flow_test.go (simulate: /estimate → upload 3 images → verify FR-015 rejection)
- [ ] T037 [US1] Integration test for invalid format in tests/integration/estimate_flow_test.go (simulate: /estimate → upload .gif → verify FR-013 error)

### Deliverables for User Story 1

- [ ] T038 [US1] Archive all vibe coding prompts for US1 in `prompts/002-calorie-image-estimate/vibe-coding/003-user-story-1.md`
- [ ] T039 [US1] Run gofmt on all US1 files: `gofmt -w src/handlers/estimate.go src/services/`
- [ ] T040 [US1] Run golangci-lint on US1 code: `golangci-lint run src/handlers/ src/services/ src/models/`
- [ ] T041 [US1] Run unit tests with coverage for US1: `go test ./tests/unit/... -cover` (must be ≥70% per constitution)

**User Story 1 Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Bot can accept /estimate, process single images, return calorie estimates with buttons. This is the **MVP**.

---

## Phase 4: User Story 2 - Re-estimation Flow (Priority: P2)

**Goal**: Implement Re-estimate button functionality to allow users to analyze another image without restarting flow.

**Independent Test**: Complete US1 flow, click "Re-estimate" button, upload new image, receive new estimate. Tests US1 is still working + adds re-estimation.

**Acceptance Criteria** (from spec.md):
1. Re-estimate button restarts image upload prompt
2. New image upload provides fresh calorie estimate

### Unit Tests for User Story 2 (Write FIRST - Red State)

> **RED**: Write these tests FIRST, ensure they FAIL before implementation

- [ ] T042 [P] [US2] Unit test for Re-estimate button handler in tests/unit/handlers_test.go (test callback transitions any state → AwaitingImage)
- [ ] T043 [P] [US2] Unit test for session state after Re-estimate in tests/unit/session_test.go (verify session not deleted, state reset to AwaitingImage)

### Implementation for User Story 2 (Make Tests GREEN)

> **GREEN**: Implement minimum code to make tests pass

- [ ] T044 [US2] Register callback handler for "re_estimate" button in src/main.go (bot.Handle(&btnReEstimate, handler))
- [ ] T045 [US2] Implement OnReEstimate callback handler in src/handlers/estimate.go (create/update session to AwaitingImage, send prompt "Please upload one food image")
- [ ] T046 [US2] Ensure Re-estimate works from any state in src/handlers/estimate.go (handle case where user clicks Re-estimate while in AwaitingImage or Processing)

### Integration Tests for User Story 2 (Verify End-to-End)

> **Integration**: Test re-estimation flow end-to-end

- [ ] T047 [US2] Integration test for Re-estimate button in tests/integration/estimate_flow_test.go (simulate: complete US1 → click Re-estimate → upload new image → verify new result)
- [ ] T048 [US2] Integration test for Re-estimate during AwaitingImage in tests/integration/estimate_flow_test.go (simulate: /estimate → click Re-estimate → verify prompt resent)

### Deliverables for User Story 2

- [ ] T049 [US2] Archive all vibe coding prompts for US2 in `prompts/002-calorie-image-estimate/vibe-coding/004-user-story-2.md`
- [ ] T050 [US2] Run gofmt and golangci-lint on US2 files
- [ ] T051 [US2] Run unit tests with coverage for US2: `go test ./tests/unit/... -cover`

**User Story 2 Checkpoint**: Re-estimation flow is now working. Users can iteratively estimate multiple images without typing /estimate again.

---

## Phase 5: User Story 3 - Cancellation Flow (Priority: P3)

**Goal**: Implement Cancel button functionality to allow users to exit estimation flow cleanly.

**Independent Test**: Start /estimate flow, click "Cancel" at any point, verify session cleared and cancellation acknowledged.

**Acceptance Criteria** (from spec.md):
1. Cancel button clears current interaction
2. Bot acknowledges cancellation
3. Session transitions to Idle state

### Unit Tests for User Story 3 (Write FIRST - Red State)

> **RED**: Write these tests FIRST, ensure they FAIL before implementation

- [ ] T052 [P] [US3] Unit test for Cancel button handler in tests/unit/handlers_test.go (test callback deletes session, sends confirmation)
- [ ] T053 [P] [US3] Unit test for Cancel from AwaitingImage state in tests/unit/session_test.go (verify session deleted)
- [ ] T054 [P] [US3] Unit test for Cancel from Processing state in tests/unit/session_test.go (verify session deleted even if Gemini call in progress)

### Implementation for User Story 3 (Make Tests GREEN)

> **GREEN**: Implement minimum code to make tests pass

- [ ] T055 [US3] Register callback handler for "cancel" button in src/main.go (bot.Handle(&btnCancel, handler))
- [ ] T056 [US3] Implement OnCancel callback handler in src/handlers/estimate.go (delete session from sync.Map, send "Estimation cancelled")
- [ ] T057 [US3] Ensure Cancel works from any state in src/handlers/estimate.go (AwaitingImage, Processing, even Idle)

### Integration Tests for User Story 3 (Verify End-to-End)

> **Integration**: Test cancellation flow end-to-end

- [ ] T058 [US3] Integration test for Cancel from AwaitingImage in tests/integration/estimate_flow_test.go (simulate: /estimate → Cancel → verify confirmation)
- [ ] T059 [US3] Integration test for Cancel after result in tests/integration/estimate_flow_test.go (simulate: complete US1 → Cancel → verify session cleared)

### Deliverables for User Story 3

- [ ] T060 [US3] Archive all vibe coding prompts for US3 in `prompts/002-calorie-image-estimate/vibe-coding/005-user-story-3.md`
- [ ] T061 [US3] Run gofmt and golangci-lint on US3 files
- [ ] T062 [US3] Run unit tests with coverage for US3: `go test ./tests/unit/... -cover`

**User Story 3 Checkpoint**: Cancellation flow is now complete. Users have full control over estimation flow with clean exit points.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Edge case handling, performance optimization, and final deliverables.

### Edge Cases (from spec.md)

- [ ] T063 [P] Handle text message after /estimate prompt in src/handlers/estimate.go (ignore or send reminder: "Please upload an image")
- [ ] T064 [P] Handle duplicate /estimate during active flow in src/handlers/estimate.go (ignore or warn: "Finish current estimation first")
- [ ] T065 [P] Handle extremely large images in src/handlers/estimate.go (check file size <20MB per Telegram limit)
- [ ] T066 [P] Add timeout for Gemini API calls in src/services/gemini.go (30 seconds per data-model.md, show timeout error to user)
- [ ] T067 [P] Add exponential backoff for Gemini rate limits in src/services/gemini.go (retry 429 errors with delay)

### Logging & Monitoring

- [ ] T068 [P] Add structured logging to all handlers in src/handlers/estimate.go (log UserID, operation type, timestamp per constitution)
- [ ] T069 [P] Add performance logging in src/handlers/estimate.go (measure time from /estimate → result, verify <10s per SC-001)
- [ ] T070 [P] Log Gemini API latency in src/services/gemini.go (measure time for generateContent call, verify <8s per SC-004)

### Final Testing & Deliverables

- [ ] T071 Create LLM test harness in tests/integration/llm_test_harness.go (custom Go tool per research.md Decision 2)
- [ ] T072 Generate LLM test prompts in `prompts/002-calorie-image-estimate/test-prompts/001-happy-path-test.md` (document prompts used for LLM testing)
- [ ] T073 Generate LLM test prompts in `prompts/002-calorie-image-estimate/test-prompts/002-edge-cases-test.md` (document edge case test prompts)
- [ ] T074 Run full integration test suite with LLM harness: `go test ./tests/integration/... -v -json > test-results.json`
- [ ] T075 Generate test report from test-results.json in `prompts/002-calorie-image-estimate/test-prompts/test-report.md` (must show 100% pass rate per SC-005)
- [ ] T076 Run final coverage check: `go test ./... -cover -coverprofile=coverage.out` (verify ≥70% per constitution)
- [ ] T077 Generate coverage HTML report: `go tool cover -html=coverage.out -o coverage.html`
- [ ] T078 Run final linting: `golangci-lint run` (must pass with zero warnings per SC-006)
- [ ] T079 Create bot via @BotFather with descriptive name per FR-017 (e.g., "CalorieScannerBot", "FoodEstimateBot")
- [ ] T080 Manual smoke test: Run bot locally, test all 3 user stories + edge cases per quickstart.md

### Final Deliverables Verification

Per constitution Principle V, verify all 5 deliverables:

- [ ] T081 **Deliverable 1**: Verify working bot - Bot runs, responds to /estimate, integrates with Gemini, has descriptive name
- [ ] T082 **Deliverable 2**: Verify clean code - All code passes gofmt and golangci-lint with zero warnings
- [ ] T083 **Deliverable 3**: Verify vibe coding prompts - All prompts archived in `prompts/002-calorie-image-estimate/vibe-coding/` (files 001-005)
- [ ] T084 **Deliverable 4**: Verify LLM test prompts - Test prompts archived in `prompts/002-calorie-image-estimate/test-prompts/` (files 001-002)
- [ ] T085 **Deliverable 5**: Verify test report - Test report generated with pass/fail results in `prompts/002-calorie-image-estimate/test-prompts/test-report.md`

---

## Task Summary

**Total Tasks**: 85
**Parallelizable Tasks**: 32 (marked with [P])

**Tasks by User Story**:
- **Setup (Phase 1)**: 9 tasks
- **Foundation (Phase 2)**: 9 tasks (blocking for all stories)
- **User Story 1 (MVP)**: 23 tasks (T019-T041)
- **User Story 2**: 10 tasks (T042-T051)
- **User Story 3**: 11 tasks (T052-T062)
- **Polish & Deliverables**: 23 tasks (T063-T085)

---

## Dependencies & Execution Order

### Critical Path (Must Execute Sequentially)

1. **Phase 1 (Setup)** → MUST complete before Phase 2
2. **Phase 2 (Foundation)** → MUST complete before ANY user story work
3. **User Stories (Phase 3-5)** → Can execute in parallel OR sequentially by priority

### User Story Dependencies

```
Foundation (T010-T018)
    ↓
    ├─→ User Story 1 (T019-T041) [P1] ← MVP
    │       ↓
    │       ├─→ User Story 2 (T042-T051) [P2] (depends on US1 handlers)
    │       └─→ User Story 3 (T052-T062) [P3] (depends on US1 handlers)
    └─→ Polish (T063-T085) (can start after US1 complete)
```

### Independent Testing per Story

Each user story has explicit **Independent Test** criteria (from spec.md):

- **US1**: Send /estimate → upload image → receive estimate (completely independent)
- **US2**: Complete US1 → click Re-estimate → upload new image (tests US1 + US2)
- **US3**: Start any flow → click Cancel → verify cleared (tests any state + US3)

---

## Parallel Execution Opportunities

### Within Phase 1 (Setup)
```bash
# Can run in parallel:
T003, T004, T005 (different go get commands)
T006, T007, T008, T009 (different files)
```

### Within Phase 2 (Foundation)
```bash
# Can run in parallel:
T010, T011, T012 (different structs in same file)
T014, T015 (different files: gemini.go vs estimate.go)
```

### Within User Story 1
```bash
# Unit tests can run in parallel (different test files):
T019, T020, T021, T022, T023

# Integration tests can run in parallel (independent scenarios):
T034, T035, T036, T037
```

### Across User Stories (If Not Following Priority Order)
```bash
# If building all stories at once:
# User Story 2 and 3 CAN run in parallel after US1 foundation is complete
# (T044-T046 for US2) || (T055-T057 for US3)
```

---

## Implementation Strategy

### Recommended: Incremental Delivery by Priority

Follow spec.md user story priorities for maximum value delivery:

**Sprint 1: MVP (P1 - User Story 1)**
```bash
# Tasks: T001-T041
# Delivers: Core /estimate command with calorie estimation
# Independent Test: Send /estimate, upload image, receive result
# MVP Checkpoint: Bot has basic functionality, ready for user testing
```

**Sprint 2: Enhanced UX (P2 - User Story 2)**
```bash
# Tasks: T042-T051
# Delivers: Re-estimation without restarting flow
# Independent Test: Click Re-estimate button, upload new image
# Value Add: Improved user experience for iterative estimation
```

**Sprint 3: User Control (P3 - User Story 3)**
```bash
# Tasks: T052-T062
# Delivers: Cancellation flow
# Independent Test: Click Cancel at any point
# Value Add: Clean exit points, better user control
```

**Sprint 4: Production Ready (Polish)**
```bash
# Tasks: T063-T085
# Delivers: Edge case handling, logging, all 5 deliverables
# Final Checkpoint: Feature complete, constitution compliant, ready for merge
```

### Alternative: Parallel Development (If Team Has Multiple Developers)

```bash
# Developer 1: US1 (T019-T041)
# Developer 2: US2 (T042-T051) after US1 foundation (T024-T030) complete
# Developer 3: US3 (T052-T062) after US1 foundation (T024-T030) complete
# All: Merge and run Polish (T063-T085) together
```

---

## Success Criteria Mapping

Tasks mapped to spec.md Success Criteria:

| Success Criterion | Related Tasks | Verification |
|-------------------|---------------|--------------|
| SC-001: Full flow <10s | T069, T070 | Performance logging in T069-T070 |
| SC-002: 90% food images estimated | T021, T032, T034 | Gemini API testing in T021, T034 |
| SC-003: /estimate response <1s | T024, T069 | Command handler latency in T069 |
| SC-004: Analysis <8s | T028, T070 | Gemini API latency in T070 |
| SC-005: 100% test pass rate | T074, T075 | Test report generation in T075 |
| SC-006: Code passes linting | T040, T050, T061, T078 | Linting runs in multiple phases |
| SC-007: Prompts archived | T038, T049, T060, T083, T084 | Deliverables 3 & 4 verification |
| SC-008: Test report with pass/fail | T075, T085 | Deliverable 5 verification |

---

## Constitution Compliance Checklist

Per constitution, verify before marking feature complete:

### Principle I: Quality-Controlled Vibe Coding
- [ ] All code passes gofmt (T040, T050, T061, T078)
- [ ] All code passes golangci-lint (T040, T050, T061, T078)
- [ ] Test coverage ≥70% for services (T076)
- [ ] All prompts archived (T038, T049, T060, T083)

### Principle II: Test-First with LLM Validation
- [ ] Unit tests written before implementation (T019-T023, T042-T043, T052-T054)
- [ ] Integration tests validate user journeys (T034-T037, T047-T048, T058-T059)
- [ ] LLM test harness implemented (T071)
- [ ] Test report generated (T075)

### Principle III: Dual Interface Architecture
- [ ] N/A (LUI only per spec)

### Principle IV: Fixed Technology Stack
- [ ] Golang 1.21+ confirmed (T001)
- [ ] telebot v3 confirmed (T003)
- [ ] Gemini SDK confirmed (T004)

### Principle V: Deliverable-Driven Development
- [ ] Working bot verified (T081)
- [ ] Clean code verified (T082)
- [ ] Vibe prompts verified (T083)
- [ ] LLM test prompts verified (T084)
- [ ] Test report verified (T085)

---

## Quick Start for Implementation

1. **Clone and checkout branch**:
   ```bash
   git checkout 002-calorie-image-estimate
   ```

2. **Run Setup (Phase 1)**:
   ```bash
   # Execute T001-T009
   # Takes ~10 minutes
   ```

3. **Run Foundation (Phase 2)**:
   ```bash
   # Execute T010-T018
   # Takes ~2 hours
   # CHECKPOINT: Foundation ready
   ```

4. **Run MVP (User Story 1)**:
   ```bash
   # Execute T019-T041
   # Takes ~8 hours (test-first workflow)
   # CHECKPOINT: MVP ready for testing
   ```

5. **Test MVP**:
   ```bash
   go test ./tests/integration/... -v
   # Verify US1 works independently
   ```

6. **Continue with US2, US3, Polish** as time permits

---

## Notes

- **Test-First is Mandatory**: All unit tests (T019-T023, T042-T043, T052-T054) MUST be written and FAIL before implementation tasks
- **Prompt Archiving**: Archive prompts incrementally (T038, T049, T060) not at the end
- **Parallel Opportunities**: Tasks marked [P] can run concurrently if working with multiple developers or AI agents
- **MVP Scope**: User Story 1 alone (T001-T041) is a complete, shippable MVP
- **Constitution Gates**: T081-T085 verify all constitution requirements before merge

---

**Generated**: 2025-12-15
**Status**: Ready for implementation via `/speckit.implement`
**Next Command**: `/speckit.implement` to begin executing tasks in order
