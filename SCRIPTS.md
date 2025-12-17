# Run Scripts Guide

All helper scripts have been unified to use the **shared storage backend** (cmd/unified/main.go).

---

## Quick Start

### Option 1: Interactive Launcher (Recommended)
```bash
./run.sh
```

Shows a menu with all available options. Perfect for first-time users.

### Option 2: Direct Commands

**Full dev environment:**
```bash
./run-dev.sh
```

**Backend only:**
```bash
./run-unified.sh
```

**Frontend only:**
```bash
./run-frontend.sh
```

**With public tunnel:**
```bash
./run-with-tunnel.sh
```

---

## Script Reference

### Available Scripts

| Script | Purpose | What It Runs |
|--------|---------|--------------|
| `run.sh` | **Interactive launcher** | Shows menu, launches other scripts |
| `run-dev.sh` | **Full dev environment** | Unified backend + frontend |
| `run-unified.sh` | **Backend only** | Bot + HTTP API (shared storage) |
| `run-frontend.sh` | **Frontend only** | Vite dev server |
| `run-with-tunnel.sh` | **Remote testing** | Everything + cloudflared tunnel |
| `view-logs.sh` | **Log viewer** | Tail logs in real-time |
| `check-bot-config.sh` | **Bot config checker** | Verify Telegram bot settings |

---

## Detailed Usage

### 1. Interactive Launcher (`run.sh`)

**When to use:** First time, or when you're not sure what to run.

```bash
./run.sh
```

**Menu options:**
1. Full development environment (backend + frontend)
2. Unified backend only
3. Frontend only
4. Full setup + Cloudflare tunnel
5. View logs
q. Quit

**Example:**
```
$ ./run.sh

╔════════════════════════════════════════════════╗
║   Telegram Calorie Bot - Unified Launcher     ║
╚════════════════════════════════════════════════╝

What would you like to run?

  1) Full Development Environment (recommended)
     → ./run-dev.sh
  ...

Enter your choice [1-5, q]: 1
```

---

### 2. Full Dev Environment (`run-dev.sh`)

**When to use:** Normal development (default choice)

```bash
./run-dev.sh
```

**What it starts:**
- ✅ Unified backend on `localhost:8080` (bot + HTTP API)
- ✅ Frontend on `localhost:5173`
- ✅ Shared storage (bot logs visible in miniapp)

**Logs:**
- Backend: `logs/backend.log`
- Frontend: `logs/frontend.log`

**Stop:** Press `Ctrl+C`

---

### 3. Unified Backend Only (`run-unified.sh`)

**When to use:**
- Testing bot commands only
- Running backend separately from frontend
- Production deployment

```bash
./run-unified.sh
```

**What it starts:**
- ✅ Telegram bot (`/estimate`, `/start`, callbacks)
- ✅ HTTP API server (`/api/logs`)
- ✅ Shared in-memory storage

**Ports:**
- HTTP: `8080` (configurable via `PORT` env var)

**Stop:** Press `Ctrl+C`

---

### 4. Frontend Only (`run-frontend.sh`)

**When to use:**
- Frontend development (backend already running)
- Debugging frontend issues

```bash
./run-frontend.sh
```

**What it starts:**
- ✅ Vite dev server on `localhost:5173`
- ✅ Hot module replacement (HMR)
- ✅ API proxy to `localhost:8080`

**Prerequisites:**
- Backend must be running (use `./run-unified.sh` in another terminal)

**Stop:** Press `Ctrl+C`

---

### 5. Full Setup + Tunnel (`run-with-tunnel.sh`)

**When to use:**
- Testing on real Telegram app (mobile/desktop)
- Need public HTTPS URL
- Configuring bot menu button

```bash
./run-with-tunnel.sh
```

**What it starts:**
1. ✅ Frontend on `localhost:5173`
2. ✅ Cloudflare tunnel (public HTTPS URL)
3. ✅ Unified backend on `localhost:8080`

**Requirements:**
- `cloudflared` installed:
  ```bash
  # macOS
  brew install cloudflared

  # Linux
  wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64
  ```

**Output:**
```
Your Mini App URL (for Telegram):
https://abc-123-def.trycloudflare.com

Next Steps:
1. Copy the URL above
2. Configure Telegram bot menu button:
   curl -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setChatMenuButton" ...
```

**Logs:**
- Backend: `logs/backend.log`
- Frontend: `logs/frontend.log`
- Tunnel: `logs/tunnel.log`

**Stop:** Press `Ctrl+C` (kills all services)

---

### 6. View Logs (`view-logs.sh`)

**When to use:** Debugging, monitoring

```bash
./view-logs.sh
```

Shows a menu to tail logs in real-time.

---

### 7. Check Bot Config (`check-bot-config.sh`)

**When to use:** Verify Telegram bot settings

```bash
./check-bot-config.sh
```

Checks:
- Bot token validity
- Current menu button configuration
- Bot info (username, name)

---

## Environment Variables

All scripts expect a `.env` file:

```bash
# Required
TELEGRAM_BOT_TOKEN=your_bot_token_here
GEMINI_API_KEY=your_gemini_api_key_here

# Optional
PORT=8080                          # HTTP server port
TUNNEL_URL=https://your-tunnel     # For CORS (auto-set by run-with-tunnel.sh)
```

**Setup:**
```bash
cp .env.example .env
# Edit .env with your values
```

---

## Common Workflows

### Workflow 1: Local Development

**Terminal 1:**
```bash
./run-dev.sh
```

**What you get:**
- Bot running (test with `/estimate` in Telegram)
- Miniapp accessible at `http://localhost:5173`
- Shared storage (logs sync between bot and miniapp)

---

### Workflow 2: Backend Dev + Frontend Dev (Separate)

**Terminal 1 (backend):**
```bash
./run-unified.sh
```

**Terminal 2 (frontend):**
```bash
./run-frontend.sh
```

**Use case:** Frontend changes with HMR, backend stays running.

---

### Workflow 3: Remote Testing (Real Telegram App)

**Terminal 1:**
```bash
./run-with-tunnel.sh
```

**Steps:**
1. Copy the tunnel URL from output
2. Run the provided `curl` command to set bot menu button
3. Open Telegram → Your bot → Menu button
4. Test the miniapp on your phone/desktop

---

### Workflow 4: Bot Testing Only (No Miniapp)

**Terminal 1:**
```bash
./run-unified.sh
```

**Use case:** Test `/estimate` command, verify bot responses, debug Gemini API.

---

## Troubleshooting

### "Address already in use" (port 8080)

**Solution 1:** Kill existing process
```bash
lsof -ti:8080 | xargs kill
```

**Solution 2:** Change port
```bash
export PORT=8081
./run-unified.sh
```

### "cloudflared: command not found"

Install cloudflared:
```bash
# macOS
brew install cloudflared

# Linux
wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64
chmod +x cloudflared-linux-amd64
sudo mv cloudflared-linux-amd64 /usr/local/bin/cloudflared
```

### "Frontend failed to start"

Check if `node_modules` exists:
```bash
cd web && npm install && cd ..
```

### Backend logs show errors

View logs:
```bash
tail -f logs/backend.log
```

Common issues:
- Missing `TELEGRAM_BOT_TOKEN` or `GEMINI_API_KEY` in `.env`
- Invalid bot token
- Network issues (Gemini API)

---

## Script Architecture

```
run.sh (interactive launcher)
  ├── run-dev.sh (full dev env)
  │   ├── cmd/unified/main.go (backend)
  │   └── web/npm run dev (frontend)
  │
  ├── run-unified.sh (backend only)
  │   └── cmd/unified/main.go
  │
  ├── run-frontend.sh (frontend only)
  │   └── web/npm run dev
  │
  └── run-with-tunnel.sh (full + tunnel)
      ├── web/npm run dev (frontend)
      ├── cloudflared tunnel (public URL)
      └── cmd/unified/main.go (backend)
```

---

## Migration from Old Scripts

### Before (Separate Storage) ❌

The old scripts `run-bot.sh` and `run-miniapp.sh` have been **removed** because they ran separate storage instances:

```bash
# OLD (removed):
./run-bot.sh          # Bot with separate storage
./run-miniapp.sh      # HTTP API with separate storage
./run-frontend.sh     # Frontend
```

**Problem:** Bot and miniapp had **separate storage**. Logs created via `/estimate` were invisible to miniapp.

### After (Unified Storage) ✅

```bash
# NEW (unified backend):
./run-dev.sh          # Bot + HTTP API + Frontend (shared storage)
```

**Solution:** One backend process, shared storage, logs sync automatically.

---

## Next Steps

1. Run `./run.sh` to get started
2. Choose option 1 (full dev environment)
3. Test `/estimate` in Telegram
4. Open miniapp and verify logs appear
5. Read `UNIFIED.md` for architecture details

---

## Quick Reference

| Task | Command |
|------|---------|
| Start everything | `./run-dev.sh` |
| Backend only | `./run-unified.sh` |
| Frontend only | `./run-frontend.sh` |
| Remote testing | `./run-with-tunnel.sh` |
| View logs | `./view-logs.sh` |
| Interactive menu | `./run.sh` |
