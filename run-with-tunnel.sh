#!/bin/bash

# Unified Telegram Bot + Mini App with Cloudflared Tunnel
# This script runs the unified backend (bot + HTTP API) + frontend with a public tunnel

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Telegram Mini App + Cloudflared Tunnel${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if cloudflared is installed
if ! command -v cloudflared &> /dev/null; then
    echo -e "${RED}Error: cloudflared is not installed${NC}"
    echo ""
    echo "Install cloudflared:"
    echo "  macOS:   brew install cloudflared"
    echo "  Linux:   https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/"
    echo ""
    exit 1
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    echo "Copy .env.example to .env and configure it"
    exit 1
fi

# Load environment variables
export $(grep -v '^#' .env | xargs)

echo -e "${BLUE}Step 1: Starting Frontend Dev Server${NC}"
echo ""

# Start frontend in background
cd web
npm run dev > ../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

sleep 3

if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo -e "${RED}Error: Frontend failed to start${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Frontend started on http://localhost:5173${NC}"
echo ""

echo -e "${BLUE}Step 2: Starting Cloudflared Tunnel${NC}"
echo ""

# Start cloudflared tunnel
cloudflared tunnel --url http://localhost:5173 > logs/tunnel.log 2>&1 &
TUNNEL_PID=$!

echo -e "${YELLOW}Waiting for tunnel to start...${NC}"
sleep 5

# Extract tunnel URL from logs
TUNNEL_URL=$(grep -oP 'https://[a-zA-Z0-9-]+\.trycloudflare\.com' logs/tunnel.log | head -1)

if [ -z "$TUNNEL_URL" ]; then
    echo -e "${RED}Error: Could not get tunnel URL${NC}"
    echo "Check logs/tunnel.log for details"
    kill $FRONTEND_PID $TUNNEL_PID 2>/dev/null
    exit 1
fi

echo -e "${GREEN}✓ Tunnel started: ${TUNNEL_URL}${NC}"
echo ""

# Export tunnel URL for backend
export TUNNEL_URL

echo -e "${BLUE}Step 3: Starting Unified Backend (Bot + HTTP API)${NC}"
echo ""

# Start unified backend
go run cmd/unified/main.go > logs/backend.log 2>&1 &
BACKEND_PID=$!

sleep 2

if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}Error: Backend failed to start${NC}"
    kill $FRONTEND_PID $TUNNEL_PID 2>/dev/null
    exit 1
fi

echo -e "${GREEN}✓ Backend started on http://localhost:8080${NC}"
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}All services running!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

echo -e "${BLUE}Your Mini App URL (for Telegram):${NC}"
echo -e "${YELLOW}${TUNNEL_URL}${NC}"
echo ""

echo -e "${BLUE}Next Steps:${NC}"
echo "1. Copy the URL above"
echo "2. Configure Telegram bot menu button:"
echo ""
echo -e "${YELLOW}curl -X POST \"https://api.telegram.org/bot\${TELEGRAM_BOT_TOKEN}/setChatMenuButton\" \\${NC}"
echo -e "${YELLOW}  -H \"Content-Type: application/json\" \\${NC}"
echo -e "${YELLOW}  -d '{${NC}"
echo -e "${YELLOW}    \"menu_button\": {${NC}"
echo -e "${YELLOW}      \"type\": \"web_app\",${NC}"
echo -e "${YELLOW}      \"text\": \"Open Calorie Logs\",${NC}"
echo -e "${YELLOW}      \"web_app\": {${NC}"
echo -e "${YELLOW}        \"url\": \"${TUNNEL_URL}\"${NC}"
echo -e "${YELLOW}      }${NC}"
echo -e "${YELLOW}    }${NC}"
echo -e "${YELLOW}  }'${NC}"
echo ""
echo "3. Open Telegram → Your bot → Menu button"
echo ""

echo -e "${BLUE}Logs:${NC}"
echo "  Frontend: logs/frontend.log"
echo "  Backend:  logs/backend.log"
echo "  Tunnel:   logs/tunnel.log"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping all services...${NC}"
    kill $FRONTEND_PID $BACKEND_PID $TUNNEL_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

# Wait for processes
wait $FRONTEND_PID $BACKEND_PID $TUNNEL_PID
