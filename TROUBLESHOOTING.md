# Troubleshooting Guide: Telegram Mini App

## Common Issues and Solutions

### Error: "Failed to fetch logs"

This error typically indicates that the frontend cannot connect to the backend server.

#### **Solution 1: Verify Backend is Running**

1. Check if the backend server is running:
   ```bash
   curl http://localhost:8080/health
   ```
   Expected output: `OK`

2. If no output, start the backend:
   ```bash
   ./run-miniapp.sh
   ```

#### **Solution 2: Check CORS Configuration**

Make sure the backend allows requests from `http://localhost:5173`:

1. Open `cmd/miniapp/main.go`
2. Verify the CORS settings include:
   ```go
   AllowedOrigins: []string{"http://localhost:5173"}
   ```

#### **Solution 3: Verify Vite Proxy Configuration**

1. Open `web/vite.config.ts`
2. Verify the proxy settings:
   ```typescript
   proxy: {
     '/api': {
       target: 'http://localhost:8080',
       changeOrigin: true,
     }
   }
   ```

---

### Error: "Unauthorized: missing X-Telegram-Init-Data header"

This error occurs when accessing the Mini App outside of Telegram.

#### **Solution: Open from Telegram**

The Mini App **must** be opened through Telegram:

1. Open your Telegram bot chat
2. Click the menu button (bottom-left, hamburger icon)
3. Select "Open Calorie Logs" or your configured menu button text

**Why?** The app requires `window.Telegram.WebApp.initData` which is only available inside Telegram's WebView.

#### **For Development Testing:**

If you need to test outside Telegram (not recommended), you can temporarily modify `web/src/api/logs.ts`:

```typescript
function getInitData(): string {
  // DEVELOPMENT ONLY - Remove for production
  return window.Telegram?.WebApp?.initData || 'user=%7B%22id%22%3A123456789%7D';
}
```

**⚠️ WARNING:** Remove this mock data before deploying to production!

---

### Error: "Cannot connect to backend server"

#### **Solution: Run Both Servers**

Use the development script to run both servers together:

```bash
./run-dev.sh
```

Or run them separately in two terminals:

**Terminal 1 (Backend):**
```bash
./run-miniapp.sh
```

**Terminal 2 (Frontend):**
```bash
./run-frontend.sh
```

---

### Frontend Shows Blank Page

#### **Solution 1: Rebuild Frontend**

```bash
cd web
npm run build
npm run dev
```

#### **Solution 2: Clear Browser Cache**

1. Open browser DevTools (F12)
2. Right-click the refresh button
3. Select "Empty Cache and Hard Reload"

---

### Backend Build Errors

#### **Solution: Update Dependencies**

```bash
go mod tidy
go mod download
```

#### **Common Dependency Issues:**

If you see `module not found` errors:

```bash
go get github.com/rs/cors
go get github.com/google/uuid
```

---

### Frontend Build Errors

#### **Solution: Reinstall Dependencies**

```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

---

### Port Already in Use

#### **Backend Port (8080) in Use:**

1. Find the process:
   ```bash
   lsof -ti:8080
   ```

2. Kill the process:
   ```bash
   kill -9 $(lsof -ti:8080)
   ```

3. Or change the port in `.env`:
   ```
   PORT=8081
   ```

#### **Frontend Port (5173) in Use:**

1. Kill the process:
   ```bash
   kill -9 $(lsof -ti:5173)
   ```

2. Or modify `web/vite.config.ts` to use a different port

---

### Logs Are Not Persisting

**This is expected behavior!** The MVP uses in-memory storage.

- Data is lost when the backend server restarts
- This is by design for local development
- Persistent storage will be added in Spec 004 (Cloudflare D1)

#### **For Testing: Seed Initial Data**

Add test data in `cmd/miniapp/main.go` after creating the storage:

```go
store := storage.NewMemoryStorage()

// Seed test data (DEVELOPMENT ONLY)
store.CreateLog(123456789, &models.Log{
    FoodItems:  []string{"Test Food"},
    Calories:   100,
    Confidence: "high",
    Timestamp:  time.Now(),
})
```

---

### Browser Console Errors

#### **"Telegram is not defined"**

This means you're accessing the app outside Telegram. See "Unauthorized" section above.

#### **"Failed to load resource: net::ERR_CONNECTION_REFUSED"**

The backend server is not running. Start it with `./run-miniapp.sh`

#### **CORS errors**

Check CORS configuration in `cmd/miniapp/main.go` and ensure it includes your frontend URL.

---

## Debugging Tips

### View Backend Logs

```bash
# If using run-dev.sh
tail -f logs/backend.log

# If running manually
go run cmd/miniapp/main.go
```

### View Frontend Logs

```bash
# If using run-dev.sh
tail -f logs/frontend.log

# Browser console (F12 → Console tab)
```

### Test API Endpoints Manually

```bash
# Health check
curl http://localhost:8080/health

# List logs (requires valid initData)
curl -H "X-Telegram-Init-Data: user=%7B%22id%22%3A123%7D" \
     http://localhost:8080/api/logs

# Create log
curl -X POST http://localhost:8080/api/logs \
  -H "X-Telegram-Init-Data: user=%7B%22id%22%3A123%7D" \
  -H "Content-Type: application/json" \
  -d '{
    "foodItems": ["Pizza"],
    "calories": 500,
    "confidence": "high",
    "timestamp": "2025-12-16T12:00:00Z"
  }'
```

---

## Getting Help

If you're still experiencing issues:

1. Check the quickstart guide: `specs/003-miniapp-mvp/quickstart.md`
2. Review the implementation plan: `specs/003-miniapp-mvp/plan.md`
3. Check backend logs: `logs/backend.log`
4. Check frontend console: Browser DevTools → Console
5. Verify all environment variables are set in `.env`

---

## Quick Start Checklist

- [ ] `.env` file created and populated
- [ ] Backend server running on port 8080
- [ ] Frontend server running on port 5173
- [ ] Telegram bot menu button configured
- [ ] Accessing app through Telegram (not direct browser)
- [ ] Both servers show no errors in logs

If all items are checked and you still have issues, check the error messages in the browser console and backend logs for specific details.
