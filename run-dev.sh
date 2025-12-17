#!/bin/bash

# Telegram Mini App Development Environment Startup Script
# This script starts both backend and frontend servers concurrently

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Telegram Mini App Development Environment${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    echo ""
    echo "Please create a .env file with the following variables:"
    echo "  TELEGRAM_BOT_TOKEN=your_bot_token_here"
    echo "  GEMINI_API_KEY=your_gemini_api_key_here"
    echo "  MINIAPP_URL=http://localhost:5173"
    echo ""
    echo "You can copy .env.example:"
    echo -e "${YELLOW}  cp .env.example .env${NC}"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go 1.21 or higher: https://go.dev/dl/"
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed${NC}"
    echo "Please install Node.js 18 or higher: https://nodejs.org/"
    exit 1
fi

# Load environment variables
export $(grep -v '^#' .env | xargs)

echo -e "${BLUE}Environment:${NC}"
echo "  Backend:  http://localhost:${PORT:-8080}"
echo "  Frontend: http://localhost:5173"
echo ""

# Install frontend dependencies if needed
if [ ! -d "web/node_modules" ]; then
    echo -e "${YELLOW}Installing frontend dependencies...${NC}"
    cd web && npm install && cd ..
    echo ""
fi

# Create log directory
mkdir -p logs

# Cleanup function to kill both servers on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping servers...${NC}"
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

# Start unified backend server (bot + HTTP API)
echo -e "${GREEN}Starting unified backend (bot + HTTP API)...${NC}"
go run cmd/unified/main.go > logs/backend.log 2>&1 &
BACKEND_PID=$!

# Wait a moment for backend to start
sleep 2

# Check if backend started successfully
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}Error: Backend server failed to start${NC}"
    echo "Check logs/backend.log for details"
    exit 1
fi

echo -e "${GREEN}✓ Backend started (PID: $BACKEND_PID)${NC}"

# Start frontend server
echo -e "${GREEN}Starting frontend server...${NC}"
cd web && npm run dev > ../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# Wait a moment for frontend to start
sleep 2

# Check if frontend started successfully
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo -e "${RED}Error: Frontend server failed to start${NC}"
    echo "Check logs/frontend.log for details"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

echo -e "${GREEN}✓ Frontend started (PID: $FRONTEND_PID)${NC}"
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Both servers are running!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Access the Mini App:${NC}"
echo "  1. Open Telegram"
echo "  2. Go to your bot chat"
echo "  3. Click the menu button (bottom-left)"
echo "  4. Select 'Open Calorie Logs'"
echo ""
echo -e "${YELLOW}Logs:${NC}"
echo "  Backend:  logs/backend.log"
echo "  Frontend: logs/frontend.log"
echo ""
echo -e "${YELLOW}To view logs in real-time:${NC}"
echo "  tail -f logs/backend.log"
echo "  tail -f logs/frontend.log"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop both servers${NC}"
echo ""

# Wait for both processes
wait $BACKEND_PID $FRONTEND_PID
