# Feature Specification: Telegram Mini App MVP (Local In-Memory)

**Feature Branch**: `003-miniapp-mvp`
**Created**: 2025-12-16
**Status**: Draft
**Input**: User description: "Build Spec 003: Telegram Mini App MVP (local, in-memory, React) with a storage abstraction that is designed to migrate to Cloudflare D1 in Spec 004."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - View Calorie Log History (Priority: P1)

As a Telegram user, I want to view all my calorie estimation logs in a table so that I can track my food intake history over time.

**Why this priority**: This is the core value proposition of the Mini App - giving users visibility into their historical data. Without this, the Mini App provides no value.

**Independent Test**: Can be fully tested by creating mock log data in the in-memory storage, launching the Mini App, and verifying the table displays all columns correctly (Date/Time, Food Items, Calories, Confidence). Delivers immediate value: users can see their estimation history.

**Acceptance Scenarios**:

1. **Given** I have 5 calorie logs in the system, **When** I open the Mini App from Telegram bot menu, **Then** I see a table with 5 rows showing all log details
2. **Given** I have no logs in the system, **When** I open the Mini App, **Then** I see "No logs yet" empty state message
3. **Given** I have logs from the past week, **When** the table loads, **Then** logs are sorted newest-first (most recent at top)

---

### User Story 2 - Create New Log Entry (Priority: P2)

As a user, I want to manually add a new calorie log entry so that I can record meals that weren't estimated via the bot.

**Why this priority**: Enables users to maintain complete records even when they can't use the /estimate command. This makes the app more useful but isn't critical for MVP value.

**Independent Test**: Can be tested by clicking "Add New Log" button, filling the form (food items, calories, confidence), clicking Save, and verifying the new entry appears at the top of the table without page refresh. Delivers value: users can build comprehensive logs.

**Acceptance Scenarios**:

1. **Given** I'm viewing the log table, **When** I click "Add New Log" button, **Then** a form appears with fields for food items, calories, confidence, and optional date/time
2. **Given** I fill the form with valid data and click "Save", **When** the save completes, **Then** the new entry appears at the top of the table and the form is cleared/closed
3. **Given** I enter invalid data (negative calories), **When** I try to save, **Then** I see a validation error message

---

### User Story 3 - Edit Existing Log Entry (Priority: P3)

As a user, I want to edit calories, food items, or confidence of an existing log so that I can correct mistakes or update my estimates.

**Why this priority**: Nice-to-have for data accuracy but not critical for initial MVP. Users can work around this by deleting and re-creating entries.

**Independent Test**: Can be tested by clicking "Edit" on a log entry row, modifying editable fields, clicking Save, and verifying the updated values appear in the table immediately. Delivers value: users can maintain accurate records without re-creating entries.

**Acceptance Scenarios**:

1. **Given** I'm viewing a log entry, **When** I click "Edit" on that row, **Then** the row enters edit mode or a modal appears with editable fields
2. **Given** I modify the calories from 500 to 600 and click "Save", **When** the update completes, **Then** the table shows 600 kcal for that entry
3. **Given** I'm editing an entry, **When** I click "Cancel", **Then** the changes are discarded and the original values remain

---

### User Story 4 - Delete Log Entry (Priority: P3)

As a user, I want to delete a log entry with confirmation so that I can remove incorrect or unwanted logs.

**Why this priority**: Data hygiene feature but not critical for MVP. Users can tolerate a few incorrect entries initially.

**Independent Test**: Can be tested by clicking "Delete" on a log entry, confirming the action in the dialog, and verifying the entry is removed from the table without page refresh. Delivers value: users can maintain clean log history.

**Acceptance Scenarios**:

1. **Given** I'm viewing a log entry, **When** I click "Delete" on that row, **Then** a confirmation dialog appears asking "Are you sure you want to delete this log?"
2. **Given** the confirmation dialog is open and I click "Confirm", **When** the deletion completes, **Then** the log entry is removed from the table
3. **Given** the confirmation dialog is open and I click "Cancel", **When** I cancel, **Then** the entry remains in the table unchanged

---

### User Story 5 - User-Scoped Logs (Priority: P1)

As a user, I want to see only my own logs so that my data is private and isolated from other users.

**Why this priority**: Fundamental security and privacy requirement. Without this, the app is unusable in a multi-user scenario.

**Independent Test**: Can be tested with two different Telegram accounts. Each user opens the Mini App and verifies they only see logs created by their own user ID. Delivers value: users trust their data is private.

**Acceptance Scenarios**:

1. **Given** User A has 3 logs and User B has 5 logs, **When** User A opens the Mini App, **Then** User A sees only their 3 logs
2. **Given** User B creates a new log, **When** User A refreshes the Mini App, **Then** User A does not see User B's new log
3. **Given** the app extracts userID from Telegram initData, **When** API calls are made, **Then** all requests include the authenticated userID

---

### Edge Cases

- What happens when a user creates a log with 0 calories? (System should allow, as some foods like water have 0 kcal)
- What happens when a user enters 10,000+ characters in the food items field? (System should enforce max length validation)
- What happens if the Go server restarts while user is viewing logs? (In-memory data is lost, table shows empty state on next load - expected MVP behavior)
- What happens when Telegram initData is missing or invalid? (System should reject API calls with 401 Unauthorized)
- What happens if two users try to create logs at the exact same timestamp? (System should handle with unique UUIDs, no collision risk)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Mini App MUST display logs in a table/list format with columns: Date/Time, Food Items, Calories (kcal), Confidence (Low/Medium/High), Actions
- **FR-002**: System MUST show "No logs yet" empty state message when user has zero logs
- **FR-003**: System MUST sort logs in descending chronological order (newest first) by timestamp
- **FR-004**: Users MUST be able to create new log entries via a form with fields: food items (string array), calories (integer), confidence (Low/Medium/High), optional timestamp (defaults to now)
- **FR-005**: Users MUST be able to edit existing log entries (food items, calories, confidence fields only)
- **FR-006**: Users MUST be able to delete log entries with a confirmation dialog
- **FR-007**: System MUST provide HTTP API endpoints:
  - `GET /api/logs?userID={id}` - List all logs for user
  - `POST /api/logs` - Create new log
  - `PATCH /api/logs/:id` - Update existing log
  - `DELETE /api/logs/:id` - Delete log
- **FR-008**: System MUST use JSON format for all API request/response bodies
- **FR-009**: System MUST define a `LogStorage` interface with methods: ListLogs, CreateLog, UpdateLog, DeleteLog
- **FR-010**: System MUST implement in-memory storage using sync.Map or mutex-protected map
- **FR-011**: System MUST extract Telegram user ID from initData for authentication
- **FR-012**: System MUST scope all logs per user (no cross-user data access)
- **FR-013**: System MUST validate log data: calories ≥ 0, confidence in ["low", "medium", "high"], food items non-empty
- **FR-014**: System MUST enable CORS for local development (frontend on :5173, backend on :8080)
- **FR-015**: React app MUST use Telegram WebApp SDK (@twa-dev/sdk) to access initData

### Key Entities *(include if feature involves data)*

- **Log**: Represents a calorie estimation record with attributes: ID (UUID), UserID (Telegram user ID), FoodItems (array of strings), Calories (integer kcal), Confidence (low/medium/high), Timestamp (when the food was consumed), CreatedAt (record creation time), UpdatedAt (last modification time)
- **LogStorage Interface**: Abstraction layer for data persistence operations, designed to support both in-memory (Spec 003) and Cloudflare D1 (Spec 004) implementations
- **UserSession**: Telegram user context containing user ID extracted from initData, used to scope API operations

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-009**: Mini App displays empty table with "No logs yet" message on first launch (0 logs)
- **SC-010**: User can create a new log entry via form, and the entry appears in the table within 200ms without page refresh
- **SC-011**: Created logs appear at the top of the table (newest-first sorting works correctly)
- **SC-012**: User can edit an existing log's calories/food items/confidence, and changes appear immediately in the table
- **SC-013**: User can delete a log with confirmation dialog, and the entry is removed from the table immediately
- **SC-014**: All CRUD operations (Create, Read, Update, Delete) work without requiring page refresh (SPA behavior)
- **SC-015**: Storage interface is defined in Go and implemented for in-memory storage, designed to be swappable with D1 implementation in Spec 004
- **SC-016**: Logs are scoped per user via Telegram initData userID - users can only see and modify their own logs
- **SC-017**: Table renders within 500ms when displaying up to 100 logs
- **SC-018**: Manual testing checklist (7 scenarios) completes with 100% pass rate

## Technical Architecture *(optional)*

### Frontend Stack
- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite (fast dev server, hot reload)
- **State Management**: useState/useEffect (no Redux for MVP)
- **HTTP Client**: fetch() or axios for API calls
- **Telegram SDK**: @twa-dev/sdk for WebApp integration
- **Styling**: Minimal CSS (optional: Tailwind CSS or plain CSS)

### Backend Stack
- **Language**: Golang 1.21+ (per Constitution Principle IV)
- **HTTP Framework**: net/http standard library or github.com/gin-gonic/gin
- **Storage**: In-memory using sync.Map or mutex-protected map[int64][]Log
- **UUID**: github.com/google/uuid for log ID generation
- **CORS**: Middleware to enable cross-origin requests from :5173 to :8080

### Project Structure
```
telegram-calories-bot/
├── cmd/
│   └── miniapp/           # Mini App HTTP server entry point
│       └── main.go
├── internal/
│   ├── storage/
│   │   ├── interface.go   # LogStorage interface definition
│   │   └── memory.go      # In-memory implementation
│   ├── handlers/
│   │   └── logs.go        # HTTP handlers for /api/logs endpoints
│   └── models/
│       └── log.go         # Log data model and validation
├── web/                   # React frontend
│   ├── src/
│   │   ├── App.tsx
│   │   ├── components/
│   │   │   ├── LogTable.tsx
│   │   │   ├── LogForm.tsx
│   │   │   └── DeleteConfirm.tsx
│   │   └── api/
│   │       └── logs.ts    # API client functions
│   ├── package.json
│   ├── vite.config.ts
│   └── index.html
├── tests/
│   ├── unit/
│   │   ├── storage_test.go
│   │   └── handlers_test.go
│   └── integration/
│       └── manual-checklist.md
└── docs/
    └── miniapp-local-dev.md  # Local development instructions
```

### API Specification

#### GET /api/logs
**Query Parameters**:
- `userID` (int64, required): Telegram user ID from initData

**Response** (200 OK):
```json
{
  "logs": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "userID": 12345,
      "foodItems": ["Rice", "Chicken"],
      "calories": 650,
      "confidence": "high",
      "timestamp": "2025-12-16T10:30:00Z",
      "createdAt": "2025-12-16T10:30:00Z",
      "updatedAt": "2025-12-16T10:30:00Z"
    }
  ]
}
```

#### POST /api/logs
**Request Body**:
```json
{
  "userID": 12345,
  "foodItems": ["Salad", "Bread"],
  "calories": 350,
  "confidence": "medium",
  "timestamp": "2025-12-16T12:00:00Z"
}
```

**Response** (201 Created):
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "userID": 12345,
  "foodItems": ["Salad", "Bread"],
  "calories": 350,
  "confidence": "medium",
  "timestamp": "2025-12-16T12:00:00Z",
  "createdAt": "2025-12-16T12:00:05Z",
  "updatedAt": "2025-12-16T12:00:05Z"
}
```

#### PATCH /api/logs/:id
**URL Parameters**:
- `id` (string, UUID): Log ID

**Query Parameters**:
- `userID` (int64, required): Telegram user ID for authorization

**Request Body** (partial update):
```json
{
  "calories": 400,
  "confidence": "high"
}
```

**Response** (200 OK):
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "userID": 12345,
  "foodItems": ["Salad", "Bread"],
  "calories": 400,
  "confidence": "high",
  "timestamp": "2025-12-16T12:00:00Z",
  "createdAt": "2025-12-16T12:00:05Z",
  "updatedAt": "2025-12-16T12:10:00Z"
}
```

#### DELETE /api/logs/:id
**URL Parameters**:
- `id` (string, UUID): Log ID

**Query Parameters**:
- `userID` (int64, required): Telegram user ID for authorization

**Response** (204 No Content)

### Data Model

```go
// Log represents a calorie estimation record
type Log struct {
    ID         string    `json:"id"`         // UUID v4
    UserID     int64     `json:"userID"`     // Telegram user ID
    FoodItems  []string  `json:"foodItems"`  // Detected food items
    Calories   int       `json:"calories"`   // Estimated kcal
    Confidence string    `json:"confidence"` // "low" | "medium" | "high"
    Timestamp  time.Time `json:"timestamp"`  // When food was consumed
    CreatedAt  time.Time `json:"createdAt"`  // Record creation time
    UpdatedAt  time.Time `json:"updatedAt"`  // Last update time
}

// LogUpdate represents partial update fields
type LogUpdate struct {
    FoodItems  *[]string `json:"foodItems,omitempty"`
    Calories   *int      `json:"calories,omitempty"`
    Confidence *string   `json:"confidence,omitempty"`
}

// LogStorage defines the storage abstraction interface
type LogStorage interface {
    // ListLogs returns all logs for a user, sorted by timestamp descending
    ListLogs(ctx context.Context, userID int64) ([]Log, error)

    // CreateLog adds a new log entry
    CreateLog(ctx context.Context, userID int64, log *Log) error

    // UpdateLog modifies an existing log (only provided fields updated)
    UpdateLog(ctx context.Context, userID int64, logID string, updates *LogUpdate) error

    // DeleteLog removes a log by ID
    DeleteLog(ctx context.Context, userID int64, logID string) error
}
```

## Non-Goals (Deferred to Spec 004)

- Cloudflare D1 database integration
- Deployment to Cloudflare Workers or Pages
- Production-ready authentication/authorization (full initData crypto verification)
- Image storage or retrieval from previous /estimate results
- Real-time sync between bot /estimate command and Mini App
- Persistent storage (data survives server restart)
- User management, admin features, or multi-user admin panel
- Advanced filtering/sorting/search capabilities
- Export logs to CSV/PDF

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Telegram initData validation complexity | Medium | Start with basic userID extraction, defer full HMAC-SHA256 signature verification to Spec 004 |
| Storage abstraction too coupled to in-memory implementation | High | Design interface from D1 perspective first (study D1 docs), implement in-memory as a conforming adapter |
| React state management becomes complex with CRUD operations | Medium | Use simple useState for form/table state, refactor to Context API only if prop drilling becomes unmanageable |
| CORS issues in local development | Low | Configure Vite proxy (`/api/*` → `http://localhost:8080`) and Go CORS middleware early in setup |
| In-memory data loss confuses users during testing | Low | Document clearly in README that data is lost on server restart (expected MVP behavior) |

## Dependencies

### External Libraries
- **Frontend**:
  - React 18+
  - Vite
  - @twa-dev/sdk (Telegram WebApp SDK)
  - axios or native fetch (HTTP client)
- **Backend**:
  - github.com/google/uuid (UUID generation)
  - github.com/gin-gonic/gin (optional HTTP framework)
  - github.com/rs/cors (CORS middleware)

### Prerequisites
- Spec 002 completed and merged (bot /estimate flow works)
- Telegram bot token configured in .env
- Node.js 18+ for React development
- Go 1.21+ for backend
- Telegram bot configured with menu button pointing to Mini App URL

## Deliverables

1. **Working Mini App**: Accessible via Telegram bot menu button with functional CRUD operations for calorie logs
2. **React Codebase**: Modular components (LogTable, LogForm, DeleteConfirm), clean separation of UI and API logic
3. **Go HTTP API**: RESTful endpoints (GET/POST/PATCH/DELETE /api/logs) with storage abstraction
4. **Storage Interface**: Documented `LogStorage` interface designed for D1 migration in Spec 004
5. **Local Development README**: Step-by-step instructions to run frontend (Vite dev server) + backend (Go HTTP server)
6. **Vibe Coding Prompts Archive**: All prompts used for implementation stored verbatim in `prompts/003-miniapp-mvp/vibe-coding/` (per Constitution v1.1.0)
7. **Test Reports**: Manual testing checklist results with pass/fail status for all 7 scenarios

## Testing Strategy

### Manual Testing Checklist
1. **Empty State**: Open Mini App → Verify "No logs yet" message appears
2. **Create Log**: Click "Add New Log", fill form, save → Verify new entry appears at top of table
3. **View Logs**: Check table displays all columns correctly (Date/Time, Food Items, Calories, Confidence, Actions)
4. **Edit Log**: Click "Edit" on existing entry, modify calories, save → Verify updated value appears
5. **Delete Log**: Click "Delete", confirm dialog → Verify entry removed from table
6. **User Isolation**: Test with two different Telegram accounts → Verify each user sees only their own logs
7. **Page Refresh**: Reload app, restart Go server → Verify in-memory data lost (expected MVP behavior)

### Unit Tests (Go)
- `TestMemoryStorage_CreateLog` - Verify log creation with valid data
- `TestMemoryStorage_ListLogs` - Verify logs sorted newest-first
- `TestMemoryStorage_UpdateLog` - Verify partial updates work correctly
- `TestMemoryStorage_DeleteLog` - Verify deletion by ID
- `TestLogValidation` - Verify validation rules (calories ≥ 0, confidence enum, food items non-empty)
- `TestHandlers_GetLogs` - Verify HTTP handler with mock storage
- `TestHandlers_CreateLog` - Verify POST handler validation and creation
- `TestHandlers_UserIsolation` - Verify users cannot access other users' logs

### Integration Tests (Optional for MVP)
- End-to-end API tests with actual HTTP server
- React component tests with React Testing Library

## Constitution Compliance Checklist

- ✅ **Principle I - Quality-Controlled Vibe Coding**: Clear acceptance criteria for all 5 user stories with Given-When-Then scenarios
- ✅ **Principle II - Test-First**: Manual testing checklist defined (7 scenarios), unit tests planned (8 tests)
- ✅ **Principle III - Dual Interface**: Mini App is second interface alongside LUI bot, shares data layer but operates independently
- ✅ **Principle IV - Fixed Tech Stack**: Golang 1.21+ for backend confirmed, React for frontend (new addition, justified for Mini App requirement)
- ✅ **Principle V - Deliverable-Driven**: All 7 deliverables listed above
- ✅ **Prompt Archiving Standards**: Verbatim prompt archiving required in `prompts/003-miniapp-mvp/vibe-coding/` per Constitution v1.1.0
- ✅ **Error Handling & Logging**: All API errors logged to console with context (user ID, operation, error value)
- ✅ **Dual Interface Architecture**: Mini App and LUI share storage layer (LogStorage interface) but have independent UIs

## Timeline and Phases

**Phase 1: Foundation** (T001-T015)
- Set up React project with Vite and TypeScript
- Define `LogStorage` interface in Go
- Implement in-memory storage with sync.Map
- Create basic HTTP server with CORS
- Set up Vite proxy for /api/* requests

**Phase 2: API Endpoints** (T016-T030)
- Implement GET /api/logs with userID filtering
- Implement POST /api/logs with validation
- Implement PATCH /api/logs/:id with partial updates
- Implement DELETE /api/logs/:id with authorization
- Add unit tests for storage and handlers

**Phase 3: Frontend CRUD** (T031-T050)
- Create LogTable component with responsive design
- Create LogForm component for Create mode
- Create LogForm component for Edit mode (inline or modal)
- Implement DeleteConfirm dialog component
- Wire up API calls with error handling

**Phase 4: Telegram Integration** (T051-T060)
- Add @twa-dev/sdk to React project
- Extract userID from Telegram initData
- Include userID in all API requests (query param or header)
- Test user isolation with two accounts
- Add basic initData validation in Go

**Phase 5: Testing & Documentation** (T061-T075)
- Write unit tests for in-memory storage (8 tests)
- Write unit tests for HTTP handlers
- Complete manual testing checklist (7 scenarios)
- Write local development README with setup steps
- Archive all vibe coding prompts verbatim

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-12-16 | Initial spec created with 5 user stories, 15 functional requirements, storage abstraction design |

---

**Next Steps**:
1. Run `/speckit.plan` to generate detailed implementation plan
2. Run `/speckit.tasks` to create actionable task list with dependencies
3. Run `/speckit.implement` to execute implementation
