# Quickstart: Telegram Mini App Local Development

**Feature**: 003-miniapp-mvp
**Target**: Developers setting up local environment for Mini App development
**Estimated Time**: 20 minutes

## Overview

This guide walks you through setting up and running the Telegram Mini App MVP locally. You'll run two servers:
1. **Go HTTP server** (backend) on `localhost:8080`
2. **Vite dev server** (frontend) on `localhost:5173`

**What you'll have at the end**: A working Mini App accessible through Telegram, with CRUD operations for calorie logs stored in-memory.

---

## Prerequisites

Ensure you have the following installed:

- **Go 1.21+**: [Download](https://go.dev/dl/)
- **Node.js 18+**: [Download](https://nodejs.org/)
- **Telegram account**: To test the Mini App
- **Git**: For cloning the repository

**Verify installations**:
```bash
go version        # Should show 1.21 or higher
node --version    # Should show v18 or higher
npm --version     # Should show 8 or higher
```

---

## Step 1: Clone and Setup Repository

```bash
# Clone the repository
git clone https://github.com/freezind/telegram-calories-bot.git
cd telegram-calories-bot

# Checkout the Mini App branch
git checkout 003-miniapp-mvp
```

---

## Step 2: Backend Setup (Go)

### Install Go Dependencies

```bash
# From repo root
go mod download
go get github.com/rs/cors
go get github.com/google/uuid
```

### Configure Environment Variables

Create a `.env` file in the repo root:

```bash
# Telegram Bot Token (from BotFather)
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Gemini API Key (from Google AI Studio)
GEMINI_API_KEY=your_gemini_api_key_here

# Mini App URL (for local development)
MINIAPP_URL=http://localhost:5173
```

**How to get tokens**:
- **Telegram Bot Token**: Message [@BotFather](https://t.me/botfather) on Telegram, create a bot, copy the token
- **Gemini API Key**: Visit [Google AI Studio](https://makersuite.google.com/app/apikey) and generate a key

### Start the Backend Server

```bash
# Option 1: Run directly
go run cmd/miniapp/main.go

# Option 2: Use helper script (loads .env automatically)
chmod +x run-miniapp.sh
./run-miniapp.sh
```

**Expected output**:
```
Server starting on http://localhost:8080
CORS enabled for http://localhost:5173
```

**Test the backend**:
```bash
curl http://localhost:8080/api/health
# Should return: {"status":"ok"}
```

---

## Step 3: Frontend Setup (React + Vite)

### Install Node Dependencies

```bash
cd web
npm install
```

**Expected packages**:
- react, react-dom (UI library)
- @vitejs/plugin-react (Vite React plugin)
- @twa-dev/sdk (Telegram WebApp SDK)
- @headlessui/react (Modal dialogs)
- typescript (Type checking)

### Configure Vite

Verify `vite.config.ts` has the proxy configuration:

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
```

### Start the Frontend Dev Server

```bash
# From web/ directory
npm run dev
```

**Expected output**:
```
VITE v5.x.x  ready in 500 ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

**Test the frontend**:
Open `http://localhost:5173` in your browser. You should see the Mini App UI (though it will show "Unauthorized" because you're not in Telegram).

---

## Step 4: Configure Telegram Bot Menu Button

To make the Mini App accessible from Telegram, you need to configure the bot's menu button.

### Method 1: Using BotFather (Recommended)

1. Open Telegram and message [@BotFather](https://t.me/botfather)
2. Send `/mybots`
3. Select your bot
4. Click **Bot Settings** → **Menu Button**
5. Click **Edit Menu Button URL**
6. Enter the URL: `http://localhost:5173` (for local testing)
7. Enter the button text: `Open Calorie Logs` (or any name you prefer)

### Method 2: Using Telegram Bot API (Advanced)

```bash
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setChatMenuButton" \
  -H "Content-Type: application/json" \
  -d '{
    "menu_button": {
      "type": "web_app",
      "text": "Open Calorie Logs",
      "web_app": {
        "url": "http://localhost:5173"
      }
    }
  }'
```

**Note**: For local development, Telegram must be able to reach `localhost:5173`. This works if you're testing on the same machine. For mobile testing, you'll need to:
- Use ngrok or similar tunneling service
- Deploy to a public HTTPS URL (e.g., Cloudflare Pages, Vercel)

---

## Step 5: Test the Mini App

### Open the Mini App from Telegram

1. Open your Telegram bot chat
2. Click the **menu button** (bottom-left, next to attach icon)
3. Select **Open Calorie Logs**
4. The Mini App should open in Telegram's WebView

### Test CRUD Operations

**Create a Log**:
1. Click **"Add New Log"** button
2. Fill in the form:
   - Food Items: `Rice, Chicken, Salad`
   - Calories: `650`
   - Confidence: `High`
3. Click **Save**
4. Verify the log appears at the top of the table

**Edit a Log**:
1. Click **Edit** on a log entry
2. Modify the calories to `700`
3. Click **Save**
4. Verify the table updates immediately

**Delete a Log**:
1. Click **Delete** on a log entry
2. Confirm the deletion in the dialog
3. Verify the log is removed from the table

**Verify User Isolation**:
1. Test with a different Telegram account
2. Verify each user sees only their own logs

---

## Step 6: Development Workflow

### Backend Development

**Watch for changes** (requires `air` or similar):
```bash
# Install air (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

**Manual restart**:
```bash
# Stop the server (Ctrl+C)
# Restart
go run cmd/miniapp/main.go
```

**Run tests**:
```bash
go test ./tests/unit/... -v
go test ./tests/integration/... -v
```

### Frontend Development

Vite provides **hot module replacement (HMR)** automatically. Just save your changes and the browser updates instantly.

**Type checking**:
```bash
npm run tsc  # TypeScript type check
```

**Linting**:
```bash
npm run lint
```

**Build for production**:
```bash
npm run build
# Output: dist/ directory
```

---

## Troubleshooting

### Error: "Unauthorized: missing X-Telegram-Init-Data header"

**Cause**: Accessing the Mini App directly in a browser (not through Telegram).

**Solution**: Open the Mini App from Telegram's menu button. The WebApp SDK only works inside Telegram's WebView.

### Error: "CORS policy: No 'Access-Control-Allow-Origin' header"

**Cause**: Backend CORS not configured or frontend port mismatch.

**Solutions**:
1. Verify backend is running on `localhost:8080`
2. Verify frontend is running on `localhost:5173`
3. Check `cors.New()` in backend has `AllowedOrigins: []string{"http://localhost:5173"}`
4. Restart both servers

### Error: "Failed to fetch" or "Network request failed"

**Cause**: Backend server not running or Vite proxy misconfigured.

**Solutions**:
1. Verify backend is running: `curl http://localhost:8080/api/health`
2. Check Vite config has `proxy: { '/api': { target: 'http://localhost:8080' } }`
3. Check browser console for detailed error

### Empty Table / No Logs Appear

**Cause**: In-memory storage is lost on server restart.

**Expected Behavior**: This is normal for MVP. Data persists only while the server is running.

**Workaround**: Seed initial data in `cmd/miniapp/main.go` for testing:

```go
// In main() after creating storage
storage.CreateLog(context.Background(), &models.Log{
    ID: uuid.New().String(),
    UserID: 12345, // Your Telegram user ID
    FoodItems: []string{"Test Food"},
    Calories: 100,
    Confidence: "high",
    Timestamp: time.Now(),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
})
```

### Mobile Testing (ngrok setup)

If you want to test on mobile without deploying:

```bash
# Install ngrok (https://ngrok.com/)
ngrok http 5173

# Copy the HTTPS URL (e.g., https://abc123.ngrok.io)
# Update BotFather menu button URL to this ngrok URL
# Update backend CORS to allow ngrok origin
```

**Update backend CORS**:
```go
AllowedOrigins: []string{
    "http://localhost:5173",
    "https://abc123.ngrok.io", // Add ngrok URL
}
```

---

## Next Steps

Once you have the local environment working:

1. **Review the codebase**:
   - `cmd/miniapp/main.go` - Server setup and routing
   - `internal/storage/memory.go` - In-memory storage implementation
   - `internal/handlers/logs.go` - HTTP handlers
   - `web/src/components/` - React components

2. **Run manual tests**: Complete the checklist in `tests/integration/manual-checklist.md`

3. **Make changes**: Follow the development workflow above

4. **Archive prompts**: Save all vibe coding prompts verbatim in `prompts/003-miniapp-mvp/vibe-coding/`

5. **Deploy (optional)**: For production deployment, see Spec 004 (Cloudflare D1 + Workers)

---

## Quick Reference

**Ports**:
- Backend: `http://localhost:8080`
- Frontend: `http://localhost:5173`
- API endpoints: `/api/logs` (proxied through frontend)

**File Locations**:
- Backend: `cmd/miniapp/main.go`, `internal/`
- Frontend: `web/src/`
- Tests: `tests/unit/`, `tests/integration/`
- Docs: `specs/003-miniapp-mvp/`

**Common Commands**:
```bash
# Backend
go run cmd/miniapp/main.go
go test ./tests/unit/... -v

# Frontend (from web/ directory)
npm run dev
npm run build
npm run lint
```

---

## Support

If you encounter issues not covered here:
1. Check API contracts: `specs/003-miniapp-mvp/contracts/api.yaml`
2. Review research findings: `specs/003-miniapp-mvp/research.md`
3. Check browser console for frontend errors
4. Check server logs for backend errors

Happy coding! 🚀
