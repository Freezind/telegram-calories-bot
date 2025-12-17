#!/bin/bash

# View logs from backend and frontend in real-time

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Telegram Mini App - Live Logs${NC}"
echo ""
echo "Checking for running servers..."
echo ""

# Check if servers are running
BACKEND_PID=$(pgrep -f "go run cmd/miniapp/main.go" | head -1)
FRONTEND_PID=$(pgrep -f "npm run dev" | head -1)

if [ -n "$BACKEND_PID" ]; then
    echo -e "${GREEN}✓ Backend running (PID: $BACKEND_PID)${NC}"
else
    echo -e "${YELLOW}⚠ Backend not running${NC}"
fi

if [ -n "$FRONTEND_PID" ]; then
    echo -e "${GREEN}✓ Frontend running (PID: $FRONTEND_PID)${NC}"
else
    echo -e "${YELLOW}⚠ Frontend not running${NC}"
fi

echo ""
echo -e "${BLUE}Select which logs to view:${NC}"
echo "  1) Backend logs"
echo "  2) Frontend logs"
echo "  3) Both (split view)"
echo "  4) All logs in one stream"
echo ""
read -p "Choice (1-4): " choice

case $choice in
    1)
        if [ -f logs/backend.log ]; then
            tail -f logs/backend.log
        else
            echo "No backend log file found. Backend might be running in a terminal."
            echo "Check the terminal where you started the backend."
        fi
        ;;
    2)
        if [ -f logs/frontend.log ]; then
            tail -f logs/frontend.log
        else
            echo "No frontend log file found. Frontend might be running in a terminal."
            echo "Check the terminal where you started the frontend."
        fi
        ;;
    3)
        if command -v tmux &> /dev/null; then
            tmux new-session \; \
              send-keys 'tail -f logs/backend.log' C-m \; \
              split-window -v \; \
              send-keys 'tail -f logs/frontend.log' C-m \;
        else
            echo "tmux not installed. Showing both logs in one stream..."
            tail -f logs/backend.log logs/frontend.log
        fi
        ;;
    4)
        if [ -f logs/backend.log ] && [ -f logs/frontend.log ]; then
            tail -f logs/backend.log logs/frontend.log
        else
            echo "Log files not found."
            echo ""
            echo "The servers might be running in terminals."
            echo "Or use ./run-dev.sh to start with file logging."
        fi
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac
