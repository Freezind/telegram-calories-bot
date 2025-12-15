# Data Model: Calorie Estimation from Image

**Feature**: 002-calorie-image-estimate
**Created**: 2025-12-15
**Status**: Phase 1 Design

## Overview

This feature operates **statelessly** with **no persistent storage** (per FR-010, FR-011, FR-012). All data structures represent in-memory, ephemeral state that exists only during user interaction sessions.

## Core Data Structures

### SessionState (Enum)

Represents the current state of a user's /estimate flow.

```go
package models

type SessionState string

const (
    // User has no active estimation flow
    StateIdle SessionState = "idle"

    // User ran /estimate, bot is waiting for image upload
    StateAwaitingImage SessionState = "awaiting_image"

    // Bot is processing uploaded image via Gemini Vision API
    StateProcessing SessionState = "processing"
)
```

**Validation Rules**:
- State transitions must follow defined state machine (see State Machine section below)
- Invalid state transitions should be logged but not block user interaction
- State resets to `Idle` on Cancel button or after result delivery

---

### UserSession (Struct)

Tracks in-memory session state for a single user during /estimate flow.

```go
package models

import "time"

type UserSession struct {
    // Telegram user ID (from c.Sender().ID)
    UserID int64

    // Current flow state
    State SessionState

    // Unix timestamp of last user interaction (for cleanup)
    LastActivity time.Time

    // Message ID of last bot message (for editing/deletion)
    MessageID int
}
```

**Lifecycle**:
- Created when user sends `/estimate` command
- Updated on state transitions and user interactions
- Cleaned up automatically after timeout (e.g., 15 minutes of inactivity)
- Deleted immediately on Cancel button or result delivery

**Storage**: Stored in `sync.Map` keyed by `UserID` (thread-safe, in-memory only)

---

### EstimateResult (Struct)

Holds the calorie estimation output from Gemini Vision API.

```go
package models

type EstimateResult struct {
    // Total estimated calories (kcal)
    Calories int `json:"calories"`

    // Confidence level: "low", "medium", or "high"
    Confidence string `json:"confidence"`

    // List of detected food items (optional, for debugging/display)
    FoodItems []string `json:"items,omitempty"`

    // Brief reasoning from Gemini (optional, for transparency)
    Reasoning string `json:"reasoning,omitempty"`
}
```

**Validation Rules** (from research.md Gemini prompt template):
- `Calories`: Must be positive integer (>0)
- `Confidence`: Must be one of {"low", "medium", "high"}
- `FoodItems`: Empty array if no food detected (triggers FR-014 error handling)
- `Reasoning`: Optional but recommended for user trust

**Parsing**: Unmarshaled from Gemini Vision JSON response (see contracts/gemini-vision.yaml)

---

## State Machine

### State Transition Diagram

```
┌─────┐  /estimate   ┌────────────────┐  upload image   ┌────────────┐
│Idle │─────────────→│AwaitingImage   │────────────────→│Processing  │
└─────┘              └────────────────┘                 └────────────┘
   ↑                        ↑                                  │
   │  Cancel button         │  Re-estimate button              │
   │  (any state)           │  (from result)                   │
   └────────────────────────┴──────────────────────────────────┘
                                  │
                                  │ result sent
                                  ↓
                               ┌─────┐
                               │Idle │
                               └─────┘
```

### State Transition Table

| From State       | Event               | To State       | Action                                  |
|------------------|---------------------|----------------|-----------------------------------------|
| Idle             | /estimate command   | AwaitingImage  | Prompt user: "Please upload one image"  |
| AwaitingImage    | Single image upload | Processing     | Download image, call Gemini API         |
| AwaitingImage    | Multiple images     | AwaitingImage  | Reject all: "Please send exactly one image" (FR-015) |
| AwaitingImage    | Text message        | AwaitingImage  | Ignore or prompt: "Please upload an image" |
| AwaitingImage    | Cancel button       | Idle           | Send: "Estimation cancelled", clear session |
| Processing       | Gemini success      | Idle           | Send result + inline buttons, clear session |
| Processing       | Gemini error        | Idle           | Send error message (FR-013), clear session |
| Any state        | Re-estimate button  | AwaitingImage  | Prompt: "Please upload one image"       |
| Any state        | /estimate (duplicate) | Current state | Ignore or warn: "Finish current estimation first" |

### Validation Rules (from Functional Requirements)

**FR-003: Image Format Validation**
```go
// Supported MIME types
var supportedFormats = []string{
    "image/jpeg",
    "image/png",
    "image/webp",
}

func ValidateImageFormat(mimeType string) error {
    for _, format := range supportedFormats {
        if mimeType == format {
            return nil
        }
    }
    return fmt.Errorf("unsupported format %s, use JPEG/PNG/WebP", mimeType)
}
```

**FR-015: Multiple Image Detection**
```go
// In telebot v3, c.Message().Album contains media group
func IsMultipleImages(msg *telebot.Message) bool {
    return len(msg.Album) > 1
}
```

**FR-014: No Food Detected**
```go
// From EstimateResult
func (r *EstimateResult) HasFood() bool {
    return len(r.FoodItems) > 0 && r.Calories > 0
}

// Error response if no food
if !result.HasFood() {
    return errors.New("no food detected in image, please try another photo")
}
```

---

## Message Formatting

### Fixed-Format Result (FR-006)

Per spec requirement for "deterministic, fixed-format result":

```go
func FormatResult(result *EstimateResult) string {
    return fmt.Sprintf(
        "🍽️ Calorie Estimate\n\n"+
        "Estimated Calories: %d kcal\n"+
        "Confidence: %s\n\n"+
        "Detected Items: %s",
        result.Calories,
        strings.Title(result.Confidence),
        strings.Join(result.FoodItems, ", "),
    )
}
```

**Example Output**:
```
🍽️ Calorie Estimate

Estimated Calories: 650 kcal
Confidence: High

Detected Items: Grilled chicken breast, Steamed broccoli, Brown rice
```

---

## Session Cleanup Strategy

Since there's no persistence, sessions must be cleaned up to prevent memory leaks:

```go
// Cleanup sessions inactive for >15 minutes
func CleanupStaleSessions(sessionMap *sync.Map) {
    sessionMap.Range(func(key, value interface{}) bool {
        session := value.(*UserSession)
        if time.Since(session.LastActivity) > 15*time.Minute {
            sessionMap.Delete(key)
        }
        return true // continue iteration
    })
}

// Run cleanup every 5 minutes
go func() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        CleanupStaleSessions(sessions)
    }
}()
```

---

## Inline Keyboard Buttons

Per FR-007, FR-008, FR-009:

```go
import "gopkg.in/telebot.v3"

var estimateKeyboard = &telebot.ReplyMarkup{}
var (
    btnReEstimate = estimateKeyboard.Data("🔄 Re-estimate", "re_estimate")
    btnCancel     = estimateKeyboard.Data("❌ Cancel", "cancel")
)

estimateKeyboard.Inline(
    estimateKeyboard.Row(btnReEstimate, btnCancel),
)
```

**Button Handlers**:
- `re_estimate`: Transition to `AwaitingImage`, prompt for new image
- `cancel`: Transition to `Idle`, send cancellation confirmation

---

## Data Flow Summary

1. **User sends /estimate**
   - Create `UserSession{State: AwaitingImage}`
   - Store in `sync.Map[UserID]`
   - Send prompt message

2. **User uploads image**
   - Validate format (FR-003)
   - Check for multiple images (FR-015)
   - Update session: `State = Processing`
   - Download image to memory
   - Send to Gemini API (see contracts/gemini-vision.yaml)

3. **Gemini returns result**
   - Unmarshal JSON → `EstimateResult`
   - Validate `HasFood()` (FR-014)
   - Format result (FR-006)
   - Send with inline buttons (FR-007)
   - Delete session from `sync.Map`

4. **User clicks Re-estimate**
   - Create new `UserSession{State: AwaitingImage}`
   - Prompt for new image

5. **User clicks Cancel**
   - Delete session from `sync.Map`
   - Send cancellation message

---

## Error Handling

All error cases must provide clear messages (FR-013):

| Error Condition            | User Message                                      | Technical Action        |
|----------------------------|---------------------------------------------------|-------------------------|
| Invalid image format       | "Invalid format. Please upload JPEG, PNG, or WebP" | Reject, stay in `AwaitingImage` |
| Multiple images uploaded   | "Please send exactly one image (not multiple)"    | Reject all, stay in `AwaitingImage` |
| No food detected           | "No food detected in image. Try another photo?"   | Clear session, go `Idle` |
| Gemini API timeout         | "Analysis timed out. Please try again."           | Clear session, go `Idle` |
| Gemini API rate limit      | "Service busy. Please wait a moment."             | Clear session, go `Idle` |
| Invalid API key (startup)  | Bot fails to start with error log                 | Exit process            |

---

## Performance Considerations

**Success Criteria Alignment**:
- **SC-003**: Bot responds to `/estimate` in <1s → Minimal session creation overhead
- **SC-004**: Image analysis + response in <8s → Gemini API call must be fast (use `gemini-2.0-flash` from research.md)
- **SC-001**: Full flow in <10s → Total latency budget: 1s (command) + 8s (analysis) + 1s (buffer)

**Memory Management**:
- No image persistence (FR-011) → Images held in memory only during Gemini API call
- Session cleanup every 5 minutes → Prevents unbounded `sync.Map` growth
- No historical data → Zero storage footprint

---

## Testing Checklist

Based on this data model, tests must cover:

- [ ] **State transitions**: All paths in state machine diagram
- [ ] **Format validation**: JPEG/PNG/WebP accepted, others rejected
- [ ] **Multiple image rejection**: Array of photos triggers FR-015 error
- [ ] **No food handling**: Empty `FoodItems` triggers FR-014 error
- [ ] **Session cleanup**: Stale sessions removed after timeout
- [ ] **Thread safety**: Concurrent `/estimate` from different users
- [ ] **Button handlers**: Re-estimate and Cancel transition to correct states
- [ ] **Message formatting**: Fixed format matches FR-006 example

See `tests/unit/session_test.go` for state machine tests, `tests/integration/estimate_flow_test.go` for end-to-end scenarios.

---

## References

- **Spec**: [spec.md](./spec.md) - Functional requirements (FR-001 to FR-017)
- **Research**: [research.md](./research.md) - Session state pattern (Decision 5)
- **Contracts**: [contracts/gemini-vision.yaml](./contracts/gemini-vision.yaml) - API request/response format
- **Constitution**: [.specify/memory/constitution.md](../../.specify/memory/constitution.md) - No persistence principle (II)
