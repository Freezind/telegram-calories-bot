# Feature Specification: Telegram Mini App MVP (Local, In-Memory, Telegram-Only)

**Feature Branch**: `003-miniapp-mvp`
**Created**: 2025-12-16
**Status**: Draft

**Scope Note**:
This Mini App is intended to be used **only within the Telegram WebApp environment**.
Direct access from a normal browser is **not supported** in Spec 003.
If Telegram `initData` is missing or cannot be parsed, API requests **MUST be rejected with 401 Unauthorized**.

---

## Clarifications

### Session 2025-12-16

- Q: What fields should the Log entity contain? → A: ID, UserID, FoodItems (array), Calories, Confidence, Timestamp, CreatedAt, UpdatedAt
- Q: How should the Edit functionality work in the UI? → A: Modal dialog with editable fields
- Q: What are the maximum length limits for food item validation? → A: 1000 characters per food item, 10 items max in array
- Q: Which HTTP framework should the backend use? → A: net/http (Go standard library)
- Q: How should the system handle concurrent edits to the same log entry? → A: Last write wins (no conflict detection)

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 – View Calorie Log History (Priority: P1)

As a Telegram user, I want to view all my calorie estimation logs in a table so that I can track my food intake history over time.

**Why this priority**:
This is the core value of the Mini App. Without visibility into existing logs, the Mini App provides no meaningful functionality.

**Independent Test**:
Seed in-memory storage with mock data, open the Mini App from Telegram, and verify logs are displayed correctly and sorted newest-first.

**Acceptance Scenarios**:

1. **Given** I have 5 calorie logs, **When** I open the Mini App, **Then** I see a table with 5 rows showing all log details
2. **Given** I have no logs, **When** I open the Mini App, **Then** I see a “No logs yet” empty state
3. **Given** logs exist across different timestamps, **When** the table loads, **Then** logs are sorted newest-first

---

### User Story 2 – Create New Log Entry (Priority: P1)

As a user, I want to manually add a new calorie log entry so that I can record meals not captured by bot estimation.

**Why this priority**:
Without the ability to create logs, the Mini App would always remain empty during demos and local testing.

**Independent Test**:
Use the UI to add a new log and verify it appears immediately at the top of the table without a page refresh.

**Acceptance Scenarios**:

1. **Given** I am viewing the log table, **When** I click “Add New Log”, **Then** a form appears
2. **Given** I submit valid data, **When** the save completes, **Then** the new entry appears at the top of the table
3. **Given** I enter invalid data (e.g. negative calories), **When** I save, **Then** I see a validation error

---

### User Story 3 – Edit Existing Log Entry (Priority: P3)

As a user, I want to edit an existing log so that I can correct mistakes.

**Acceptance Scenarios**:

1. **Given** I view a log entry, **When** I click "Edit", **Then** a modal dialog opens with editable fields (FoodItems, Calories, Confidence)
2. **Given** I modify fields and click "Save", **When** the update completes, **Then** the modal closes and the table reflects the updated values
3. **Given** the modal is open and I click "Cancel", **Then** the modal closes and no changes are applied

---

### User Story 4 – Delete Log Entry (Priority: P3)

As a user, I want to delete a log entry with confirmation.

**Acceptance Scenarios**:

1. **Given** I click “Delete”, **Then** a confirmation dialog appears
2. **Given** I confirm deletion, **Then** the entry is removed from the table
3. **Given** I cancel deletion, **Then** the entry remains unchanged

---

### User Story 5 – User-Scoped Logs (Priority: P1)

As a user, I want to see **only my own logs**, so that my data is private.

**Acceptance Scenarios**:

1. **Given** User A has 3 logs and User B has 5 logs, **When** User A opens the Mini App, **Then** only User A’s logs are visible
2. **Given** User B creates a log, **When** User A reloads, **Then** User A does not see it
3. **Given** a request is made, **Then** the server derives user identity from Telegram initData and scopes all operations accordingly

---

## Edge Cases

- Logs with 0 calories are allowed
- Food item text exceeding 1000 characters per item is rejected with validation error
- Food item arrays exceeding 10 items are rejected with validation error
- Concurrent edits to the same log entry use last write wins strategy (no conflict detection)
- Server restart clears all data (expected in-memory behavior)
- Missing or unparsable initData → 401 Unauthorized

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Mini App MUST display logs in a table/list
- **FR-002**: System MUST show a “No logs yet” empty state
- **FR-003**: Logs MUST be sorted newest-first
- **FR-004**: Users MUST be able to create new log entries
- **FR-005**: Users MUST be able to edit existing log entries via modal dialog (editable fields: FoodItems, Calories, Confidence)
- **FR-006**: Users MUST be able to delete log entries with confirmation dialog
- **FR-007**: System MUST expose CRUD HTTP endpoints under `/api/logs` using net/http standard library
- **FR-008**: API payloads MUST use JSON
- **FR-009**: System MUST define a storage abstraction interface (LogStorage)
- **FR-010**: Data MUST be stored in memory only (Spec 003) using sync.Map or mutex-protected map
- **FR-011**: User identity MUST be derived from Telegram initData
- **FR-012**: Data MUST be scoped per authenticated user

---

## Data Model

### Log Entity

```go
type Log struct {
    ID         string    `json:"id"`         // UUID v4
    UserID     int64     `json:"userID"`     // Telegram user ID from initData
    FoodItems  []string  `json:"foodItems"`  // Array of food item names
    Calories   int       `json:"calories"`   // Estimated kcal (≥0)
    Confidence string    `json:"confidence"` // "low" | "medium" | "high"
    Timestamp  time.Time `json:"timestamp"`  // When food was consumed
    CreatedAt  time.Time `json:"createdAt"`  // Record creation time
    UpdatedAt  time.Time `json:"updatedAt"`  // Last modification time
}
```

**Field Constraints:**
- `ID`: Unique identifier, generated as UUID v4
- `UserID`: Must match authenticated Telegram user from initData
- `FoodItems`: Non-empty array of strings; max 10 items; each item max 1000 characters
- `Calories`: Integer ≥ 0
- `Confidence`: Enum value from ["low", "medium", "high"]
- `Timestamp`: Defaults to current time if not provided; used for newest-first sorting
- `CreatedAt`: Set once on creation
- `UpdatedAt`: Updated on every modification

---

## Success Criteria

- Logs render correctly and update without page reload
- CRUD operations complete successfully
- Users cannot access other users’ logs
- System behaves deterministically for demo usage

---

## Non-Goals

- Persistence or cloud deployment
- initData signature verification
- Image uploads or LLM integration
- Advanced UI features

---

## Next Steps

1. Run `/speckit.plan`
2. Run `/speckit.tasks`
3. Implement Spec 003
