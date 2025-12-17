# Data Model: Telegram Mini App MVP

**Feature**: 003-miniapp-mvp
**Created**: 2025-12-16
**Status**: Draft

## Overview

This document defines the data structures and storage contracts for the Telegram Mini App MVP. The design prioritizes a clean storage abstraction that can migrate from in-memory (Spec 003) to Cloudflare D1 (Spec 004) with minimal changes to business logic.

## Entities

### Log

Represents a single calorie estimation record created by a user.

**Go Type Definition**:

```go
package models

import "time"

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

**Field Specifications**:

| Field | Type | Constraints | Default | Notes |
|-------|------|-------------|---------|-------|
| `ID` | string | UUID v4 format, unique | Generated on create | Primary identifier |
| `UserID` | int64 | Must match authenticated Telegram user | From initData | Foreign key to Telegram user |
| `FoodItems` | []string | Non-empty array; max 10 items; each item max 1000 chars | Required | User-entered or LLM-detected |
| `Calories` | int | ≥ 0 | Required | Total estimated kcal |
| `Confidence` | string | Enum: "low", "medium", "high" | Required | Estimation confidence level |
| `Timestamp` | time.Time | Valid RFC3339 timestamp | `time.Now()` if not provided | When food was consumed |
| `CreatedAt` | time.Time | Valid RFC3339 timestamp | `time.Now()` on create | Immutable after creation |
| `UpdatedAt` | time.Time | Valid RFC3339 timestamp | `time.Now()` on create/update | Updated on every modification |

**Validation Rules**:

```go
func (l *Log) Validate() error {
    // ID validation
    if l.ID == "" {
        return fmt.Errorf("ID is required")
    }
    if _, err := uuid.Parse(l.ID); err != nil {
        return fmt.Errorf("ID must be valid UUID: %w", err)
    }

    // UserID validation
    if l.UserID <= 0 {
        return fmt.Errorf("UserID must be positive")
    }

    // FoodItems validation
    if len(l.FoodItems) == 0 {
        return fmt.Errorf("FoodItems cannot be empty")
    }
    if len(l.FoodItems) > 10 {
        return fmt.Errorf("FoodItems cannot exceed 10 items (got %d)", len(l.FoodItems))
    }
    for i, item := range l.FoodItems {
        if len(item) > 1000 {
            return fmt.Errorf("FoodItems[%d] exceeds 1000 characters (got %d)", i, len(item))
        }
    }

    // Calories validation
    if l.Calories < 0 {
        return fmt.Errorf("Calories must be non-negative (got %d)", l.Calories)
    }

    // Confidence validation
    validConfidence := map[string]bool{"low": true, "medium": true, "high": true}
    if !validConfidence[l.Confidence] {
        return fmt.Errorf("Confidence must be low/medium/high (got %s)", l.Confidence)
    }

    return nil
}
```

**Lifecycle**:

1. **Create**: ID generated as UUID v4; CreatedAt and UpdatedAt set to current time; Timestamp defaults to current time if not provided
2. **Read**: No modifications
3. **Update**: Only FoodItems, Calories, Confidence can be modified; UpdatedAt set to current time; ID, UserID, Timestamp, CreatedAt are immutable
4. **Delete**: Hard delete (no soft delete in MVP)

---

### LogUpdate

Represents a partial update to an existing Log entry. Used for PATCH operations.

**Go Type Definition**:

```go
type LogUpdate struct {
    FoodItems  *[]string `json:"foodItems,omitempty"`
    Calories   *int      `json:"calories,omitempty"`
    Confidence *string   `json:"confidence,omitempty"`
}
```

**Usage**:
- Pointers allow distinguishing between "field not provided" (nil) vs "field set to zero value" (e.g., `&[]string{}`)
- Only non-nil fields are applied to the target Log
- All validation rules from Log entity apply to provided fields

**Example JSON**:

```json
{
  "calories": 500,
  "confidence": "high"
}
```

This updates only Calories and Confidence, leaving FoodItems unchanged.

---

## Storage Abstraction

### LogStorage Interface

Defines the contract for log persistence operations. Designed to support both in-memory (Spec 003) and Cloudflare D1 (Spec 004) implementations.

**Go Interface Definition**:

```go
package storage

import (
    "context"
    "github.com/freezind/telegram-calories-bot/internal/models"
)

type LogStorage interface {
    // ListLogs returns all logs for a user, sorted by Timestamp descending (newest first)
    // Returns empty slice if user has no logs
    ListLogs(ctx context.Context, userID int64) ([]models.Log, error)

    // CreateLog adds a new log entry
    // Returns error if log with same ID already exists or validation fails
    CreateLog(ctx context.Context, log *models.Log) error

    // UpdateLog modifies an existing log (only provided fields in updates are changed)
    // Returns error if log not found or userID mismatch (prevents cross-user modification)
    UpdateLog(ctx context.Context, userID int64, logID string, updates *models.LogUpdate) error

    // DeleteLog removes a log by ID
    // Returns error if log not found or userID mismatch (prevents cross-user deletion)
    DeleteLog(ctx context.Context, userID int64, logID string) error
}
```

**Design Principles**:

1. **User Scoping**: All operations (except CreateLog) require userID to enforce data isolation
2. **Context-Aware**: All methods accept `context.Context` for cancellation, timeouts, and future tracing
3. **Error Transparency**: Methods return errors for not-found, validation failures, and authorization mismatches
4. **Immutability**: Interface does not expose batch operations or transactions (can be added in Spec 004 if needed)

**Error Semantics**:

- `ErrLogNotFound`: Log ID does not exist
- `ErrUnauthorized`: UserID mismatch (user attempting to access another user's log)
- `ErrValidation`: Log or LogUpdate failed validation
- Other errors: Storage-layer failures (e.g., D1 connection errors in Spec 004)

---

### In-Memory Implementation (Spec 003)

**Implementation Strategy**:

```go
package storage

import (
    "context"
    "fmt"
    "sort"
    "sync"
    "github.com/freezind/telegram-calories-bot/internal/models"
)

type MemoryStorage struct {
    mu   sync.RWMutex
    logs map[int64][]models.Log  // userID -> logs array
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        logs: make(map[int64][]models.Log),
    }
}

// ListLogs implementation
func (m *MemoryStorage) ListLogs(ctx context.Context, userID int64) ([]models.Log, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    userLogs, exists := m.logs[userID]
    if !exists {
        return []models.Log{}, nil
    }

    // Deep copy to prevent external mutation
    result := make([]models.Log, len(userLogs))
    copy(result, userLogs)

    // Sort newest-first by Timestamp
    sort.Slice(result, func(i, j int) bool {
        return result[i].Timestamp.After(result[j].Timestamp)
    })

    return result, nil
}

// Additional methods follow similar pattern with RWMutex protection
```

**Concurrency Strategy**:
- Use `sync.RWMutex` for read-heavy workload (ListLogs frequent, writes occasional)
- Deep copy on reads to prevent external mutation of internal state
- Last-write-wins for concurrent updates (no optimistic locking in MVP)

**Why RWMutex over sync.Map**:
- `sync.Map` optimized for key-stable, write-once-read-many patterns
- Our use case has dynamic keys (new users) and frequent slice modifications
- `RWMutex` with `map[int64][]Log` provides simpler debugging and clearer semantics

---

## Data Flow

### Create Flow

```
User Form (React)
  → POST /api/logs {userID, foodItems, calories, confidence, timestamp?}
    → Handler validates JSON
      → Handler calls CreateLog(ctx, log)
        → MemoryStorage validates log.Validate()
          → MemoryStorage appends to logs[userID]
            ← Returns nil (success)
          ← Returns log with ID/timestamps
        ← Returns 201 Created {log}
      ← UI appends log to table state
```

### Update Flow

```
User clicks Edit → Modal opens
  → User modifies fields → Clicks Save
    → PATCH /api/logs/:id {calories?, foodItems?, confidence?}
      → Handler validates JSON
        → Handler calls UpdateLog(ctx, userID, logID, updates)
          → MemoryStorage finds log by ID
            → Validates userID match
              → Applies non-nil fields from updates
                → Updates UpdatedAt timestamp
                  ← Returns nil (success)
                ← Returns 200 OK {updatedLog}
              ← UI updates log in table state
```

### Delete Flow

```
User clicks Delete → Confirmation dialog
  → User confirms
    → DELETE /api/logs/:id?userID={userID}
      → Handler calls DeleteLog(ctx, userID, logID)
        → MemoryStorage finds log by ID
          → Validates userID match
            → Removes from logs[userID]
              ← Returns nil (success)
            ← Returns 204 No Content
          ← UI removes log from table state
```

---

## Migration Path to D1 (Spec 004)

**Changes Required**:

1. **Implement `D1Storage` struct** conforming to `LogStorage` interface
2. **SQL Schema**:
   ```sql
   CREATE TABLE logs (
       id TEXT PRIMARY KEY,
       user_id INTEGER NOT NULL,
       food_items TEXT NOT NULL,  -- JSON array
       calories INTEGER NOT NULL,
       confidence TEXT NOT NULL,
       timestamp TEXT NOT NULL,   -- ISO8601
       created_at TEXT NOT NULL,
       updated_at TEXT NOT NULL,
       INDEX idx_user_timestamp (user_id, timestamp DESC)
   );
   ```
3. **Swap implementation** in `cmd/miniapp/main.go`:
   ```go
   // Spec 003
   storage := storage.NewMemoryStorage()

   // Spec 004
   storage := storage.NewD1Storage(d1Client)
   ```

**Zero business logic changes** required in handlers or models.

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-12-16 | Initial data model with Log entity, LogStorage interface, MemoryStorage design |
