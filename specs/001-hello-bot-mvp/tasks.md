# Tasks: Hello Bot MVP

**Input**: Design documents from `/specs/001-hello-bot-mvp/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md

**Tests**: Unit test for greeter service is REQUIRED per spec (FR-005).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

Single project structure per plan.md:
- Source: `cmd/`, `internal/` at repository root
- Tests: `tests/` at repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create directory structure (cmd/bot, internal/handlers, internal/services, tests/unit)
- [x] T002 Initialize Go module with `go mod init github.com/freezind/telegram-calories-bot`
- [x] T003 Add Telebot dependency with `go get github.com/tucnak/telebot/v3`
- [x] T004 Create .env.example file documenting TELEGRAM_BOT_TOKEN requirement
- [x] T005 Create .gitignore with .env entry to prevent committing secrets

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core greeting service that MUST be complete before ANY user story handler can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Write failing test in tests/unit/greeter_test.go (test-first: verify "Hello! 👋" output)
- [x] T007 Create Greet() service function in internal/services/greeter.go (returns "Hello! 👋")
- [x] T008 Run test and verify it passes

**Checkpoint**: Foundation ready - greeter service tested and working. User story handlers can now be implemented in parallel.

---

## Phase 3: User Story 1 - Command Response (Priority: P1) 🎯 MVP

**Goal**: Bot responds to /start command with "Hello! 👋" and logs user ID

**Independent Test**: Send /start to bot via Telegram, verify exact response "Hello! 👋" and check logs for user ID

### Implementation for User Story 1

- [x] T009 [US1] Create cmd/bot/main.go with bot initialization (load TELEGRAM_BOT_TOKEN, create bot, start polling)
- [x] T010 [US1] Implement /start handler in internal/handlers/greeting.go (calls Greet(), logs user ID, sends response)
- [x] T011 [US1] Register /start handler in cmd/bot/main.go
- [x] T012 [US1] Test manually: set TELEGRAM_BOT_TOKEN, run bot, send /start, verify response and logs

**Checkpoint**: At this point, User Story 1 should be fully functional - bot handles /start command correctly

---

## Phase 4: User Story 2 - Text Message Response (Priority: P1)

**Goal**: Bot responds to plain text "hello" with "Hello! 👋" and logs user ID

**Independent Test**: Send "hello" text message to bot via Telegram, verify exact response "Hello! 👋" and check logs for user ID

### Implementation for User Story 2

- [x] T013 [US2] Implement OnText handler in internal/handlers/greeting.go (filters for "hello", calls Greet(), logs user ID)
- [x] T014 [US2] Register OnText handler in cmd/bot/main.go
- [X] T015 [US2] Test manually: send "hello", verify response and logs
- [X] T016 [US2] Test edge case: send "Hello" (capitalized), verify NO response (case-sensitive per spec)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - bot handles both /start and "hello"

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T017 Run `go fmt ./...` to format all Go code
- [x] T018 Run `go test ./... -v` and verify greeter test passes
- [x] T019 Verify go.mod and go.sum are committed
- [X] T020 Test concurrent users: send multiple messages quickly, verify all get responses
- [x] T021 Verify logs contain user IDs for all messages
- [x] T022 Review code against spec - confirm no out-of-scope features added (no database, no Mini App, no LLM, no deployment, no extra commands)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **User Story 2 (Phase 4)**: Depends on Foundational phase completion - Can run in parallel with US1 if different files
- **Polish (Phase 5)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on US2
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies on US1

**IMPORTANT**: Both user stories depend on the greeter service (Phase 2), so Phase 2 must complete first. However, once Phase 2 is done, US1 and US2 can proceed in parallel since they modify different parts of the same files or can be structured to avoid conflicts.

### Within Each User Story

- US1: main.go initialization → handler implementation → handler registration → manual test
- US2: handler implementation → handler registration → manual test → edge case test

### Parallel Opportunities

- Phase 1 tasks T001-T005: Some can run in parallel (directory creation, go mod init, dependency install)
- Phase 2: Test must be written first (T006), then implementation (T007), then verify (T008) - sequential
- **Phase 3 and Phase 4 can overlap**: If US1 developer works on main.go initialization and /start handler while US2 developer works on OnText handler, they can merge afterward. However, for a single developer, do US1 fully first (it initializes the bot), then add US2 handlers.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (greeter service + test)
3. Complete Phase 3: User Story 1 (/start handler)
4. **STOP and VALIDATE**: Test /start command independently
5. Bot is now minimally functional with /start support

### Full MVP (Both User Stories)

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Bot handles /start
3. Add User Story 2 → Test independently → Bot handles "hello"
4. Both stories work independently and together

### Parallel Team Strategy

With two developers:

1. Both complete Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (main.go + /start handler)
   - Developer B: User Story 2 (OnText handler)
3. Merge: Developer A's bot initialization + Developer B's OnText handler

---

## Notes

- **Test-First**: Phase 2 starts with failing test (T006) before implementing greeter service
- **Minimal Scope**: This is the absolute minimum to validate bot functionality. No database, Mini App, LLM, or deployment.
- **Manual Testing**: No automated integration tests for Telegram API interaction - manual testing via actual bot messages
- **File Paths Are Explicit**: Every task references exact file paths from plan.md structure
- **Both Stories Are P1**: They're equally critical for demonstrating bot functionality
- **Expected Total LOC**: ~80-100 lines across 4 files (very minimal)
- **Constitution Compliance**: Follows test-first (Phase 2), separation of concerns (service vs handler), and quality control (manual validation)

---

## Validation Checklist

Before considering this feature complete:

- [x] All tasks follow checklist format with IDs
- [x] Tasks organized by user story (Setup → Foundational → US1 → US2 → Polish)
- [x] Each user story has clear goal and independent test criteria
- [x] File paths are explicit in task descriptions
- [x] Dependencies clearly documented
- [x] Parallel opportunities identified
- [x] Test-first workflow enforced (Phase 2)
- [x] Scope strictly controlled (no out-of-scope features)
