# Implementation Tasks: Telegram Mini App MVP

**Feature**: 003-miniapp-mvp
**Branch**: `003-miniapp-mvp`
**Created**: 2025-12-16
**Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md)

## Overview

This document organizes implementation tasks by user story to enable independent, incremental delivery. Each user story phase represents a complete, testable increment of functionality.

**Key Principles**:
- Demo MVP scope: Functional correctness, simplicity, no performance optimization
- User identity EXCLUSIVELY from X-Telegram-Init-Data header (NO query params/body)
- Manual testing only (no automated UI tests)
- In-memory storage with LogStorage abstraction for future D1 migration

---

## Task Summary

- **Total Tasks**: 62
- **Parallelizable Tasks**: 28 (marked with [P])
- **User Stories**: 5 (US1-US5)
- **Phases**: 8 (Setup → 5 User Stories → Polish)

---

## Phase 1: Project Setup (T001-T010)

**Goal**: Initialize project structure, install dependencies, configure development environment.

- [x] T001 Create Go module and project structure (cmd/miniapp/, internal/, web/)
- [x] T002 [P] Install Go dependencies (github.com/google/uuid, github.com/rs/cors)
- [x] T003 [P] Initialize React project with Vite and TypeScript in web/
- [x] T004 [P] Install frontend dependencies (@twa-dev/sdk, @headlessui/react)
- [x] T005 [P] Create .env.example file with TELEGRAM_BOT_TOKEN, GEMINI_API_KEY, MINIAPP_URL
- [x] T006 [P] Configure Vite proxy for /api/* → localhost:8080 in web/vite.config.ts
- [x] T007 [P] Set up ESLint and TypeScript config in web/
- [x] T008 [P] Create README.md or update existing with Mini App setup instructions
- [x] T009 [P] Initialize Git ignore patterns for node_modules/, dist/, .env
- [x] T010 Run `go mod tidy` and `npm install` to verify setup

---

## Phase 2: Foundational Layer (T011-T022)

**Goal**: Implement core backend infrastructure required by ALL user stories (blocking prerequisites).

### Data Model & Storage Interface

- [x] T011 [P] Create Log entity in internal/models/log.go with 8 fields (ID, UserID, FoodItems, Calories, Confidence, Timestamp, CreatedAt, UpdatedAt)
- [x] T012 [P] Implement Log.Validate() method with validation rules (calories ≥0, confidence enum, food items constraints)
- [x] T013 [P] Create LogUpdate struct in internal/models/log.go for partial updates
- [x] T014 [P] Define LogStorage interface in internal/storage/interface.go with 4 methods (ListLogs, CreateLog, UpdateLog, DeleteLog)

### In-Memory Storage

- [x] T015 [P] Implement MemoryStorage struct in internal/storage/memory.go with sync.RWMutex and map[int64][]Log
- [x] T016 [P] Implement MemoryStorage.ListLogs() sorted by Timestamp descending
- [x] T017 [P] Implement MemoryStorage.CreateLog() with UUID generation
- [x] T018 [P] Implement MemoryStorage.UpdateLog() with userID authorization check
- [x] T019 [P] Implement MemoryStorage.DeleteLog() with userID authorization check

### Authentication & Middleware

- [x] T020 [P] Create TelegramUser struct and ParseInitData() function in internal/auth/initdata.go (extracts userID from X-Telegram-Init-Data header)
- [x] T021 [P] Implement AuthMiddleware in internal/middleware/auth.go (validates initData, adds userID to context, returns 401 if missing/invalid)
- [x] T022 Create HTTP server in cmd/miniapp/main.go with CORS config for localhost:5173

---

## Phase 3: User Story 1 – View Calorie Log History (T023-T032) [P1]

**Story Goal**: Users can view all their calorie logs in a table, sorted newest-first, with empty state handling.

**Independent Test**: Seed in-memory storage with mock data, open Mini App from Telegram, verify logs display correctly and are sorted newest-first.

### Backend API

- [x] T023 [US1] Implement GET /api/logs handler in internal/handlers/logs.go (extracts userID from context, calls storage.ListLogs(), returns JSON)
- [x] T024 [US1] Register GET /api/logs route with AuthMiddleware in cmd/miniapp/main.go

### Frontend Components

- [x] T025 [P] [US1] Create LogTable component in web/src/components/LogTable.tsx (displays logs in table with columns: Date/Time, Food Items, Calories, Confidence, Actions)
- [x] T026 [P] [US1] Implement empty state UI in LogTable ("No logs yet" message when logs array is empty)
- [x] T027 [P] [US1] Create API client function fetchLogs() in web/src/api/logs.ts (includes X-Telegram-Init-Data header from WebApp.initData)
- [x] T028 [P] [US1] Create App.tsx in web/src/ as main component (fetches logs on mount, passes to LogTable)
- [x] T029 [P] [US1] Add Telegram WebApp SDK initialization in web/src/main.tsx
- [x] T030 [P] [US1] Style LogTable with basic CSS (responsive, mobile-friendly)

### Integration & Testing

- [ ] T031 [US1] Test US1: Seed storage with 5 logs, verify table displays all rows sorted newest-first
- [ ] T032 [US1] Test US1: Start with empty storage, verify "No logs yet" message appears

**Story 1 Complete**: ✅ Users can view their log history

---

## Phase 4: User Story 2 – Create New Log Entry (T033-T042) [P1]

**Story Goal**: Users can manually add new calorie log entries via a form modal.

**Independent Test**: Use UI to add a new log, verify it appears immediately at the top of the table without page refresh.

### Backend API

- [ ] T033 [US2] Implement POST /api/logs handler in internal/handlers/logs.go (parses JSON body, validates, generates ID/timestamps, calls storage.CreateLog())
- [ ] T034 [US2] Register POST /api/logs route with AuthMiddleware in cmd/miniapp/main.go

### Frontend Components

- [ ] T035 [P] [US2] Create LogForm component in web/src/components/LogForm.tsx using @headlessui/react Dialog (form fields: foodItems, calories, confidence, optional timestamp)
- [ ] T036 [P] [US2] Implement form state management in LogForm with useState (controlled inputs)
- [ ] T037 [P] [US2] Add client-side validation in LogForm (negative calories error, empty food items error, max 10 items/1000 chars)
- [ ] T038 [P] [US2] Create API client function createLog() in web/src/api/logs.ts (POST request with X-Telegram-Init-Data header)
- [ ] T039 [P] [US2] Add "Add New Log" button in App.tsx (opens LogForm modal)
- [ ] T040 [P] [US2] Implement form submission in LogForm (calls createLog(), refreshes log list on success, closes modal)

### Integration & Testing

- [ ] T041 [US2] Test US2: Submit valid log data, verify new entry appears at top of table without page refresh
- [ ] T042 [US2] Test US2: Submit invalid data (negative calories), verify validation error appears in modal

**Story 2 Complete**: ✅ Users can create new log entries

---

## Phase 5: User Story 5 – User-Scoped Logs (T043-T047) [P1]

**Story Goal**: Users see only their own logs (data isolation via Telegram initData authentication).

**Independent Test**: Test with two different Telegram accounts, verify each user sees only their own logs.

### Backend Security

- [ ] T043 [US5] Add userID scope enforcement in all storage methods (verify in MemoryStorage.UpdateLog and DeleteLog that log.UserID matches authenticated userID)
- [ ] T044 [US5] Add logging for authentication events in AuthMiddleware (log userID extraction, log 401 errors with context)

### Frontend

- [ ] T045 [P] [US5] Ensure all API calls in web/src/api/logs.ts include X-Telegram-Init-Data header from WebApp.initData
- [ ] T046 [P] [US5] Add error handling in App.tsx for 401 Unauthorized (display "Please open from Telegram" message)

### Integration & Testing

- [ ] T047 [US5] Test US5: Use two Telegram accounts (User A, User B), verify User A sees only their 3 logs and User B sees only their 5 logs

**Story 5 Complete**: ✅ User data is isolated and private

---

## Phase 6: User Story 3 – Edit Existing Log Entry (T048-T054) [P3]

**Story Goal**: Users can edit existing log entries (food items, calories, confidence) via modal dialog.

**Independent Test**: Click Edit on a log, modify calories, save, verify table updates immediately.

### Backend API

- [ ] T048 [US3] Implement PATCH /api/logs/:id handler in internal/handlers/logs.go (parses log ID from URL, parses LogUpdate from body, calls storage.UpdateLog())
- [ ] T049 [US3] Register PATCH /api/logs/:id route with AuthMiddleware in cmd/miniapp/main.go

### Frontend Components

- [ ] T050 [P] [US3] Extend LogForm component to support edit mode (populate form with initialData prop, change modal title to "Edit Log")
- [ ] T051 [P] [US3] Create API client function updateLog() in web/src/api/logs.ts (PATCH request with partial updates)
- [ ] T052 [P] [US3] Add Edit button to each row in LogTable (opens LogForm modal in edit mode with log data)
- [ ] T053 [P] [US3] Implement edit submission in LogForm (calls updateLog(), refreshes log list, closes modal)

### Integration & Testing

- [ ] T054 [US3] Test US3: Edit a log's calories from 500 to 600, verify table shows updated value and modal closes

**Story 3 Complete**: ✅ Users can edit log entries

---

## Phase 7: User Story 4 – Delete Log Entry (T055-T060) [P3]

**Story Goal**: Users can delete log entries with confirmation dialog.

**Independent Test**: Click Delete, confirm, verify entry is removed from table.

### Backend API

- [ ] T055 [US4] Implement DELETE /api/logs/:id handler in internal/handlers/logs.go (parses log ID, calls storage.DeleteLog(), returns 204 No Content)
- [ ] T056 [US4] Register DELETE /api/logs/:id route with AuthMiddleware in cmd/miniapp/main.go

### Frontend Components

- [ ] T057 [P] [US4] Create DeleteConfirm component in web/src/components/DeleteConfirm.tsx using @headlessui/react Dialog (confirmation message, Confirm/Cancel buttons)
- [ ] T058 [P] [US4] Create API client function deleteLog() in web/src/api/logs.ts (DELETE request)
- [ ] T059 [P] [US4] Add Delete button to each row in LogTable (opens DeleteConfirm dialog)
- [ ] T060 [P] [US4] Implement delete confirmation in DeleteConfirm (calls deleteLog(), refreshes log list, closes dialog)

### Integration & Testing

- [ ] T061 [US4] Test US4: Click Delete on a log, confirm, verify log is removed from table immediately

**Story 4 Complete**: ✅ Users can delete log entries

---

## Phase 8: Polish & Final Integration (T062)

**Goal**: Final testing, documentation, and demo readiness.

- [ ] T062 Perform end-to-end manual testing per quickstart.md checklist (all 5 user stories, all edge cases)

---

## Dependencies Between User Stories

```
Setup (Phase 1)
  ↓
Foundational (Phase 2) ← BLOCKING for all user stories
  ↓
┌─────────────┬─────────────┬─────────────┐
│   US1 (P1)  │   US2 (P1)  │   US5 (P1)  │  ← Can implement in parallel
│   View Logs │ Create Logs │ User Scoping│
└─────────────┴─────────────┴─────────────┘
  ↓
┌─────────────┬─────────────┐
│   US3 (P3)  │   US4 (P3)  │  ← Can implement in parallel (depend on US1, US2)
│  Edit Logs  │ Delete Logs │
└─────────────┴─────────────┘
  ↓
Polish (Phase 8)
```

**Critical Path**:
1. Setup → Foundational (MUST complete first)
2. US1 + US2 + US5 (P1 stories, can parallelize)
3. US3 + US4 (P3 stories, depend on US1/US2 completing)

---

## Parallel Execution Opportunities

### After Setup (Phase 1) Complete:

**Parallel Track 1 - Backend Foundation**:
- T011-T022 (Data model, storage, auth middleware) - All parallelizable

**Parallel Track 2 - Frontend Foundation** (can start immediately):
- T006-T007 (Vite config, ESLint setup)

### After Foundational (Phase 2) Complete:

**Parallel User Stories (P1)**:
- **Track A**: US1 (View Logs) - T023-T032
- **Track B**: US2 (Create Logs) - T033-T042
- **Track C**: US5 (User Scoping) - T043-T047

All three can be developed simultaneously by different developers.

### After P1 Stories Complete:

**Parallel User Stories (P3)**:
- **Track D**: US3 (Edit Logs) - T048-T054
- **Track E**: US4 (Delete Logs) - T055-T061

---

## MVP Scope Recommendation

**Minimum Viable Demo**: Complete Phase 1 (Setup) + Phase 2 (Foundational) + Phase 3 (US1: View Logs) + Phase 4 (US2: Create Logs)

This provides:
- ✅ Users can view their log history
- ✅ Users can manually add new logs
- ✅ Basic CRUD functionality for demo
- ✅ Telegram WebApp integration
- ✅ User authentication via initData

**Incremental Delivery**:
1. **MVP**: US1 + US2 (View + Create)
2. **Phase 2**: Add US5 (User Scoping) for multi-user demos
3. **Phase 3**: Add US3 + US4 (Edit + Delete) for complete CRUD

---

## Task Completion Checklist

When completing each task:

1. ✅ Implement the functionality per task description
2. ✅ Verify file path is correct (use absolute paths from repo root)
3. ✅ Run linting (gofmt for Go, ESLint for TypeScript)
4. ✅ Test locally (manual verification as per user story acceptance scenarios)
5. ✅ Mark task as complete in this file: `- [x]`

---

## Implementation Strategy

### Recommended Order:

1. **Week 1**: Setup + Foundational (T001-T022)
   - Get local dev environment working
   - Backend API skeleton ready
   - Frontend scaffold ready

2. **Week 2**: MVP (US1 + US2) (T023-T042)
   - View logs functionality
   - Create logs functionality
   - Basic demo ready

3. **Week 3**: User Isolation (US5) + Edit/Delete (US3, US4) (T043-T061)
   - Multi-user testing
   - Complete CRUD operations
   - Full feature set

4. **Week 4**: Polish & Testing (T062)
   - Manual testing all scenarios
   - Bug fixes
   - Demo preparation

---

## Notes

- **No automated tests**: Manual testing only per constitution (demo MVP scope)
- **Authentication**: User identity derived EXCLUSIVELY from X-Telegram-Init-Data header (backend MUST NOT accept userID from query params or request body)
- **Performance**: Not optimized (out of scope for demo)
- **Storage**: In-memory only, data lost on restart (expected behavior)
- **Prompt Archiving**: Use ./prompts.md at repo root per constitution (no feature-specific prompt directories)

---

**Status**: Ready for implementation
**Last Updated**: 2025-12-16
