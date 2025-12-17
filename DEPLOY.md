# Deployment Guide

## Option 1: Docker (Recommended)

### Quick Start with Docker Compose

```bash
# 1. Set environment variables in .env file
cp .env.example .env
# Edit .env and add your tokens

# 2. Build and run
docker-compose up -d

# 3. Check logs
docker-compose logs -f

# 4. Stop
docker-compose down
```

### Manual Docker Build

```bash
# Build image
docker build -t telegram-calories-bot .

# Run container
docker run -d \
  --name telegram-bot \
  -p 8080:8080 \
  -e TELEGRAM_BOT_TOKEN="your_token" \
  -e GEMINI_API_KEY="your_key" \
  -e PORT=8080 \
  --restart unless-stopped \
  telegram-calories-bot

# Check logs
docker logs -f telegram-bot

# Stop and remove
docker stop telegram-bot && docker rm telegram-bot
```

### Deploy to Railway with Docker

Railway automatically detects and uses the Dockerfile:

1. Push code to GitHub
2. Create new project in Railway
3. Connect GitHub repository
4. Set environment variables in Railway dashboard
5. Railway will automatically build and deploy using Dockerfile

## Option 2: Railway (Nixpacks)

### 1. Deploy to Railway

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template)

Or manually:

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login
railway login

# Initialize project
railway init

# Link to existing project (if already created)
railway link

# Deploy
railway up
```

### 2. Set Environment Variables

In Railway dashboard, add these environment variables:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
GEMINI_API_KEY=your_gemini_api_key_here
```

**Note:** `PORT` is automatically provided by Railway.

### 3. Verify Deployment

Railway will automatically:
- Build using `go build -o unified cmd/unified/main.go`
- Start using `./unified`
- Assign a public URL

Check logs in Railway dashboard to confirm:
```
[BOT] ✓ Telegram bot initialized
[HTTP] ✓ HTTP API server initialized
[HTTP] 🚀 HTTP server listening on http://localhost:XXXX
[BOT] 🚀 Telegram bot started
```

## Manual Build & Test Locally

```bash
# Build
go build -o unified cmd/unified/main.go

# Run (requires environment variables)
export TELEGRAM_BOT_TOKEN="your_token"
export GEMINI_API_KEY="your_key"
export PORT=8080
./unified
```

Or use the start script:

```bash
# Set environment variables first
export TELEGRAM_BOT_TOKEN="your_token"
export GEMINI_API_KEY="your_key"

# Run
./start.sh
```

## Architecture

The unified backend runs both:
1. **Telegram Bot** (LUI) - Handles /start, /estimate commands
2. **HTTP API Server** - Serves Mini App and CRUD endpoints

Both share the same in-memory storage for seamless integration.

## Health Check

Once deployed, test the health endpoint:

```bash
curl https://your-app.railway.app/health
# Expected: OK
```

## Troubleshooting

### Bot not responding
- Verify `TELEGRAM_BOT_TOKEN` is set correctly
- Check Railway logs for connection errors

### Mini App 401 errors
- Ensure Mini App URL in Telegram BotFather matches Railway URL
- Check CORS configuration includes your tunnel/production URL

### Port binding errors
- Railway automatically sets `PORT` - don't override it manually
- Ensure app binds to `0.0.0.0:$PORT`, not `localhost:$PORT`
