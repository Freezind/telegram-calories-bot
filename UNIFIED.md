# Unified Backend Architecture

## Overview

The bot backend (spec 002) and miniapp backend (spec 003) have been **unified into a single Go process** that runs both services concurrently with **shared in-memory storage**.

This means:
- ✅ Logs created via `/estimate` bot command are **immediately visible** in the miniapp
- ✅ Logs created/edited in the miniapp are visible to the bot's storage layer
- ✅ **One process** to run instead of two separate servers
- ✅ Simpler deployment and development workflow

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│         cmd/unified/main.go                     │
│  (Single Go process)                            │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────────┐      ┌──────────────────┐    │
│  │ Telegram Bot │      │ HTTP API Server  │    │
│  │ (Spec 002)   │      │ (Spec 003)       │    │
│  │              │      │                  │    │
│  │ /estimate    │      │ GET  /api/logs   │    │
│  │ /start       │      │ POST /api/logs   │    │
│  │ Photo upload │      │ PATCH /api/logs  │    │
│  │ Callbacks    │      │ DELETE /api/logs │    │
│  └──────┬───────┘      └────────┬─────────┘    │
│         │                       │              │
│         └───────┬───────────────┘              │
│                 ▼                              │
│      ┌─────────────────────┐                   │
│      │  Shared Storage     │                   │
│      │  (MemoryStorage)    │                   │
│      │  - CreateLog()      │                   │
│      │  - ListLogs()       │                   │
│      │  - UpdateLog()      │                   │
│      │  - DeleteLog()      │                   │
│      └─────────────────────┘                   │
└─────────────────────────────────────────────────┘
```

---

## Key Changes

### 1. Modified Files

**Backend (Go):**
- ✅ `cmd/unified/main.go` - New unified entrypoint (runs bot + HTTP server)
- ✅ `src/handlers/estimate.go` - EstimateHandler now accepts `storage` parameter and calls `CreateLog()` after successful estimation
- ✅ `run-unified.sh` - New script to run unified backend only
- ✅ `run-dev.sh` - Updated to use `cmd/unified/main.go` instead of `cmd/miniapp/main.go`

**Frontend Auth Fix (separate from unification):**
- ✅ `internal/auth/initdata.go` - Fixed JSON parsing bug (was looking for `id` param instead of parsing `user` JSON)
- ✅ `internal/middleware/auth.go` - Added debug logging
- ✅ `web/src/App.tsx` - Added guard for missing initData
- ✅ `web/src/main.tsx` - Enhanced logging

### 2. Code Mapping: EstimateResult → Log

When `/estimate` succeeds, the bot now creates a log entry:

```go
// src/handlers/estimate.go (after sending result to user)
logEntry := &internalmodels.Log{
    FoodItems:  result.FoodItems,  // []string from Gemini
    Calories:   result.Calories,   // int from Gemini
    Confidence: internalmodels.ConfidenceLevel(result.Confidence), // "low"/"medium"/"high"
    Timestamp:  time.Now(),
}

h.storage.CreateLog(userID, logEntry)
```

The storage automatically adds:
- `ID` (UUID)
- `UserID` (from Telegram sender ID)
- `CreatedAt` / `UpdatedAt` timestamps

### 3. Authentication

Each service maintains its own auth method:

**Bot (Telegram SDK):**
- UserID: `c.Sender().ID` (int64)
- Auth: Telegram bot token validates all incoming updates

**HTTP API (Miniapp):**
- UserID: Extracted from `X-Telegram-Init-Data` header
- Auth: `AuthMiddleware` parses initData and validates user

Both auth flows produce the **same UserID format** (int64), so storage works seamlessly.

---

## Running the Unified Backend

### Option 1: Full Dev Environment (Backend + Frontend)

```bash
./run-dev.sh
```

This starts:
1. Unified backend (bot + HTTP API) on `localhost:8080`
2. Vite frontend dev server on `localhost:5173`

Logs:
- `logs/backend.log`
- `logs/frontend.log`

### Option 2: Backend Only

```bash
./run-unified.sh
```

This starts only the unified backend. Useful if you:
- Already have frontend running separately
- Only testing bot commands
- Running in production

### Option 3: Manual (for debugging)

```bash
# Load environment variables
export $(grep -v '^#' .env | xargs)

# Run with live logs
go run cmd/unified/main.go
```

---

## Environment Variables

Required in `.env`:
```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
GEMINI_API_KEY=your_gemini_api_key_here
PORT=8080                           # Optional, defaults to 8080
TUNNEL_URL=https://your-tunnel-url  # Optional, for Cloudflare Tunnel CORS
```

---

## Verification Steps

### 1. Check Backend Logs

When unified backend starts, you should see:
```
[STORAGE] ✓ Shared MemoryStorage initialized
[BOT] ✓ Telegram bot initialized (session cleanup routine started)
[HTTP] ✓ HTTP API server initialized
[HTTP] 🚀 HTTP server listening on http://localhost:8080
[BOT] 🚀 Telegram bot started
```

### 2. Test Bot → Miniapp Integration

1. Open Telegram and send `/estimate` to your bot
2. Upload a food image
3. Bot responds with calorie estimate
4. Backend logs: `[HANDLER] ✓ Log entry saved for user 123456: 450 kcal, 2 items`
5. Open the miniapp (Menu → "Open Calorie Logs")
6. **You should see the log entry appear immediately!**

### 3. Test Miniapp → Storage

1. In miniapp, click "Add New Log"
2. Fill in food items and calories
3. Save
4. Backend logs: `[HTTP PATCH] /api/logs completed in 2ms`

---

## Troubleshooting

### Bot works but logs don't appear in miniapp

**Check**: Backend logs should show:
```
[HANDLER] ✓ Log entry saved for user 123456: 450 kcal, 2 items
```

If you see:
```
[HANDLER ERROR] Failed to save log entry for user 123456: ...
```

Check that the EstimateResult validation passes (food items non-empty, confidence valid).

### Miniapp shows 401 errors

See the auth debugging guide in `TROUBLESHOOTING.md`. The recent fix should resolve initData parsing issues.

### "address already in use" error

Another process is using port 8080. Either:
- Kill existing process: `lsof -ti:8080 | xargs kill`
- Change port: `export PORT=8081` before running

---

## Migration Notes

### Old Setup (Separate Processes) - DEPRECATED

**Before unification (old scripts removed):**
```bash
# Terminal 1: Bot only (run-bot.sh - REMOVED)
go run src/main.go

# Terminal 2: HTTP API only (run-miniapp.sh - REMOVED)
go run cmd/miniapp/main.go

# Terminal 3: Frontend
cd web && npm run dev
```

**Issue**: Bot logs and miniapp logs were in **separate storage instances** (both in-memory, but different maps). `/estimate` logs were invisible to miniapp.

**Note**: The old helper scripts `run-bot.sh` and `run-miniapp.sh` have been **removed** to prevent confusion.

### New Setup (Unified)

**After unification:**
```bash
# Terminal 1: Unified backend (bot + HTTP API)
./run-unified.sh

# Terminal 2: Frontend
cd web && npm run dev
```

**Or use the all-in-one script:**
```bash
./run-dev.sh   # Starts both in background
```

---

## Next Steps

1. ✅ Test end-to-end: `/estimate` → miniapp displays log
2. ✅ Verify auth fixes work (see `TROUBLESHOOTING.md`)
3. ⬜ Deploy unified backend to production (single Dockerfile/service)
4. ⬜ Add persistent storage (D1/SQLite) to replace in-memory map
5. ⬜ Add WebSocket/SSE for real-time log sync (optional)

---

## Files Reference

**Unified Backend:**
- `cmd/unified/main.go` - Main entrypoint
- `src/handlers/estimate.go` - Bot handlers with storage injection
- `internal/handlers/logs.go` - HTTP API handlers
- `internal/storage/memory.go` - Shared storage implementation

**Scripts:**
- `run.sh` - Interactive launcher (recommended)
- `run-unified.sh` - Run backend only
- `run-dev.sh` - Run backend + frontend
- `run-with-tunnel.sh` - Run with Cloudflare tunnel

**Docs:**
- `UNIFIED.md` - This file
- `TROUBLESHOOTING.md` - Auth debugging guide
- `specs/003-miniapp-mvp/quickstart.md` - Original miniapp setup guide
