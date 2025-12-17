# Implementation Plan: Telegram Mini App MVP (Local In-Memory)

**Branch**: `003-miniapp-mvp` | **Date**: 2025-12-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-miniapp-mvp/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Create a React-based Telegram Mini App providing CRUD operations for calorie log management. The backend uses Go with in-memory storage and a storage abstraction layer designed for future Cloudflare D1 migration (Spec 004). The Mini App is accessible only within Telegram WebApp environment, with user authentication via initData.

**Primary Requirement**: Enable users to view, create, edit, and delete calorie logs through a table UI accessible from Telegram bot menu.

**Technical Approach**:
- Frontend: React 18+ with Vite, TypeScript, Telegram WebApp SDK
- Backend: Go 1.21+ with net/http standard library
- Storage: In-memory (sync.Map) implementing LogStorage interface
- Authentication: Telegram initData extraction (basic userID validation, no signature verification in MVP)

## Technical Context

**Language/Version**: Go 1.21+ (backend), JavaScript/TypeScript (React 18+ frontend)
**Primary Dependencies**:
- Backend: net/http (stdlib), github.com/google/uuid, github.com/rs/cors
- Frontend: React 18+, Vite, @twa-dev/sdk, axios or fetch

**Storage**: In-memory only (sync.Map or mutex-protected map); LogStorage interface for future D1 migration
**Testing**: Go testing (unit tests for storage/handlers), manual testing checklist (integration)
**Target Platform**: Local development (localhost:8080 backend, localhost:5173 frontend)
**Project Type**: Web application (React frontend + Go backend)

**Constraints**:
- **Demo MVP scope**: Focus on functional correctness, simplicity, and demo viability
- **Performance**: Not measured or optimized (out of scope for demo)
- **Telegram WebApp only**: No direct browser access supported
- **No persistent storage**: Data lost on server restart (expected in-memory behavior)
- **Authentication**: User identity derived EXCLUSIVELY from X-Telegram-Init-Data header (MUST NOT accept userID via query params or request body)
- **User isolation**: All operations scoped by authenticated Telegram user
- **Validation**: Food items max 10 items, 1000 chars each; calories ≥ 0

**Scale/Scope**:
- Local development demo only
- No production deployment
- Manual testing (no automated UI tests)
- ~5 React components, 4 HTTP endpoints

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Quality-Controlled Vibe Coding
✅ **PASS** - Spec has clear acceptance criteria for all user stories (US-001 to US-005)
✅ **PASS** - Manual testing checklist defined (7 scenarios)
✅ **PASS** - Unit tests planned for Go storage and handlers
✅ **PASS** - Linting: gofmt, golangci-lint for Go; ESLint for React (standard tools)

### Principle II: Test-First with LLM Validation
✅ **PASS** - Manual testing checklist defined as primary validation
✅ **PASS** - Unit tests planned (TestMemoryStorage_*, TestHandlers_*)
⚠️ **PARTIAL** - No LLM-based bot integration testing (deferred: Mini App is separate from bot LUI)

**Justification**: Mini App operates independently from bot commands. LLM testing applies to bot conversation flows (Spec 002), not GUI CRUD operations. Manual testing sufficient for MVP.

### Principle III: Dual Interface Architecture
✅ **PASS** - Mini App is second interface alongside bot LUI (Spec 002)
✅ **PASS** - Shares data layer concept (storage abstraction) but operates independently
✅ **PASS** - No shared state between Mini App and bot handlers

**Note**: Spec 003 establishes storage abstraction. Future integration (Spec 004+) will connect bot /estimate results to Mini App logs via shared D1 database.

### Principle IV: Fixed Technology Stack
✅ **PASS** - Backend uses Golang 1.21+ per constitution
⚠️ **NEW DEPENDENCY** - React 18+ for frontend (not previously used)

**Justification**: Telegram Mini Apps require web frontend. React is industry-standard for this use case. Constitution allows justified dependencies when necessary for requirements.

### Principle V: Deliverable-Driven Development
✅ **PASS** - Core deliverables defined in spec:
1. Working Mini App accessible via Telegram
2. React codebase (modular components)
3. Go HTTP API with storage abstraction
4. Storage interface documented
5. Local development setup guide (quickstart.md)

**Note**: Prompt archiving handled by constitution (see ./prompts.md at repo root). No feature-specific prompt directories.

### Error Handling & Logging
✅ **PASS** - API errors must be logged with context (user ID, operation, error value)
✅ **PASS** - Validation errors returned to client with descriptive messages
✅ **PASS** - 401 Unauthorized for missing/invalid initData

### Authentication Requirements
✅ **PASS** - User identity MUST be derived exclusively from X-Telegram-Init-Data header
❌ **FORBIDDEN** - Backend MUST NOT accept userID via query params or request body
✅ **PASS** - All operations scoped by authenticated user from initData

**GATE RESULT**: ✅ **PASSED** - Proceed to Phase 0

## Project Structure

### Documentation (this feature)

```text
specs/003-miniapp-mvp/
├── spec.md              # Feature specification (already exists)
├── plan.md              # This file (/speckit.plan output)
├── data-model.md        # Phase 1 output (Log entity, storage interface)
├── quickstart.md        # Phase 1 output (local dev setup) - REQUIRED
├── research.md          # Phase 0 output (optional, non-blocking)
├── contracts/           # Phase 1 output (optional, non-blocking)
│   └── api.yaml         # OpenAPI spec for HTTP endpoints
└── tasks.md             # Phase 2 output (/speckit.tasks - not yet created)
```

**Documentation Focus**: Prioritize quickstart.md for developer onboarding. Research and API contracts are supplementary.

### Source Code (repository root)

```text
# Web application structure (React frontend + Go backend)

cmd/
└── miniapp/             # Mini App HTTP server entry point
    └── main.go          # Server startup, routing, CORS setup

internal/
├── storage/
│   ├── interface.go     # LogStorage interface definition
│   └── memory.go        # In-memory implementation (sync.Map)
├── handlers/
│   └── logs.go          # HTTP handlers for /api/logs endpoints
└── models/
    └── log.go           # Log entity, validation, JSON serialization

web/                     # React frontend (NEW directory)
├── src/
│   ├── App.tsx          # Main app component
│   ├── components/
│   │   ├── LogTable.tsx       # Table display component
│   │   ├── LogForm.tsx        # Create/edit form (modal)
│   │   └── DeleteConfirm.tsx  # Confirmation dialog
│   ├── api/
│   │   └── logs.ts      # API client functions (fetch/axios)
│   └── main.tsx         # React entry point
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts       # Vite config with /api/* proxy

tests/
├── unit/
│   ├── storage_test.go  # In-memory storage tests
│   └── handlers_test.go # HTTP handler tests
└── integration/
    └── manual-checklist.md  # Manual testing procedures

docs/
└── miniapp-local-dev.md # Local development setup instructions

prompts.md               # Prompt archive at repo root (Constitution v1.1.0)
```

**Structure Decision**: Web application structure selected due to React frontend + Go backend. Frontend lives in `/web` directory (standard React+Vite layout). Backend extends existing `/cmd` and `/internal` structure from Spec 002. Separation allows independent hot-reload for frontend (Vite) and backend (Go).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| React dependency (new tech) | Telegram Mini Apps require web UI; React is industry standard for component-based SPAs | Plain HTML/JS insufficient for CRUD state management and SPA behavior; Telegram WebApp SDK has React examples |

**Note**: This is not a true violation - Constitution allows justified dependencies. Documented here for transparency.

---

## Phase 0: Research & Unknowns

**Status**: ✅ COMPLETE

**Unknowns resolved**:
1. ✅ Telegram WebApp initData structure and parsing in Go → Use URL-encoded query string from `X-Telegram-Init-Data` header
2. ✅ CORS configuration for localhost cross-origin → Use github.com/rs/cors with explicit localhost:5173 origin
3. ✅ Vite proxy setup for /api/* requests → Use vite.config.ts server.proxy with target: localhost:8080
4. ✅ sync.Map vs mutex-protected map → Use map[int64][]Log with sync.RWMutex (better for CRUD operations)
5. ✅ Best practices for modal dialog implementation in React → Use @headlessui/react Dialog component

**Deliverable**: [research.md](./research.md) with detailed findings (optional, non-blocking)

---

## Phase 1: Design & Contracts

**Status**: ✅ COMPLETE

**Core Artifacts (Required)**:

1. **Data Model** ([data-model.md](./data-model.md))
   - Log entity with 8 fields (ID, UserID, FoodItems, Calories, Confidence, Timestamp, CreatedAt, UpdatedAt)
   - Field constraints and validation rules
   - LogUpdate partial update type
   - LogStorage interface (abstraction for in-memory + future D1)
   - MemoryStorage implementation strategy

2. **Quickstart Guide** ([quickstart.md](./quickstart.md)) - **REQUIRED**
   - Step-by-step local development setup
   - Backend and frontend installation instructions
   - Telegram bot menu button configuration
   - CRUD testing procedures
   - Troubleshooting guide

3. **Agent Context Update** (CLAUDE.md)
   - Added React 18+ frontend technology
   - Added in-memory storage with LogStorage interface
   - Added web application project type

**Supplementary Artifacts (Optional)**:

4. **API Contracts** ([contracts/api.yaml](./contracts/api.yaml)) - Non-blocking
   - OpenAPI 3.0.3 specification
   - 4 endpoints: GET /api/logs, POST /api/logs, PATCH /api/logs/:id, DELETE /api/logs/:id
   - Authentication via X-Telegram-Init-Data header (NO query params or body userID)
   - Complete request/response schemas with examples
   - Error handling (400, 401, 404)

---

## Constitution Re-Check (Post-Design)

**Status**: ✅ PASSED

All constitution principles remain satisfied after detailed design:

### Updated Assessments:
- ✅ **Quality Gates**: Data model includes complete validation rules, clear error handling
- ✅ **Dual Interface**: LogStorage abstraction confirmed, ready for D1 migration in Spec 004
- ✅ **Deliverables**: Core artifacts generated (data-model.md, quickstart.md); supplementary docs optional
- ✅ **Dependencies Justified**: React, @headlessui/react, rs/cors researched and documented
- ✅ **Authentication**: User identity exclusively from X-Telegram-Init-Data header (no query params/body)
- ✅ **Simplicity Focus**: Demo MVP scope, no performance optimization, functional correctness prioritized

**No new violations introduced.**

---

## Planning Summary

**Phase 0 & 1 Status**: ✅ COMPLETE

**Next Step**: Run `/speckit.tasks` to generate Phase 2 implementation tasks

**Key Decisions Made**:
1. **Technology stack**: React + TypeScript, Go + net/http, rs/cors, headlessui
2. **Storage**: RWMutex-protected map with LogStorage interface (D1-ready)
3. **Authentication**: User identity EXCLUSIVELY from X-Telegram-Init-Data header (forbidden: query params, request body)
4. **API design**: RESTful with partial PATCH updates, 204 for DELETE
5. **UX pattern**: Headless UI modals for Create/Edit forms
6. **Scope**: Demo MVP - functional correctness, simplicity, no performance metrics
7. **Documentation**: Quickstart guide required, research/contracts optional

**Implementation Ready**: Core design complete, authentication constraints clear, development guide available.
