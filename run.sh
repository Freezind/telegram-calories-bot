#!/bin/bash

# Telegram Calorie Bot - Unified Launcher
# Smart script that detects your needs and runs the appropriate setup

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

clear
echo -e "${GREEN}╔════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   Telegram Calorie Bot - Unified Launcher     ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════╝${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  No .env file found${NC}"
    echo ""
    echo "First time setup? Let me help you create the .env file."
    echo ""

    if [ -f .env.example ]; then
        echo -e "${BLUE}Option 1: Copy from .env.example (recommended)${NC}"
        echo "  cp .env.example .env"
        echo ""
        echo -e "${BLUE}Option 2: Create manually with these variables:${NC}"
    else
        echo -e "${BLUE}Create .env file with these required variables:${NC}"
    fi

    echo "  TELEGRAM_BOT_TOKEN=your_bot_token_here"
    echo "  GEMINI_API_KEY=your_gemini_api_key_here"
    echo "  PORT=8080 (optional)"
    echo ""
    exit 1
fi

echo -e "${CYAN}What would you like to run?${NC}"
echo ""
echo "  1) Full Development Environment (recommended)"
echo "     - Unified backend (bot + HTTP API)"
echo "     - Frontend dev server"
echo "     - Shared storage (bot logs visible in miniapp)"
echo "     → ./run-dev.sh"
echo ""
echo "  2) Unified Backend Only"
echo "     - Bot + HTTP API in one process"
echo "     - No frontend server"
echo "     → ./run-unified.sh"
echo ""
echo "  3) Frontend Only"
echo "     - Vite dev server on localhost:5173"
echo "     - Requires backend running separately"
echo "     → ./run-frontend.sh"
echo ""
echo "  4) Full Setup + Cloudflare Tunnel (for remote testing)"
echo "     - Unified backend + frontend"
echo "     - Public HTTPS URL via cloudflared"
echo "     - Configure Telegram bot menu button"
echo "     → ./run-with-tunnel.sh"
echo ""
echo "  5) View Logs"
echo "     → ./view-logs.sh"
echo ""
echo "  q) Quit"
echo ""

read -p "Enter your choice [1-5, q]: " choice

case $choice in
    1)
        echo ""
        echo -e "${GREEN}Starting full development environment...${NC}"
        ./run-dev.sh
        ;;
    2)
        echo ""
        echo -e "${GREEN}Starting unified backend...${NC}"
        ./run-unified.sh
        ;;
    3)
        echo ""
        echo -e "${GREEN}Starting frontend only...${NC}"
        echo -e "${YELLOW}⚠️  Make sure the backend is running separately!${NC}"
        sleep 2
        ./run-frontend.sh
        ;;
    4)
        echo ""
        echo -e "${GREEN}Starting with Cloudflare tunnel...${NC}"
        ./run-with-tunnel.sh
        ;;
    5)
        echo ""
        ./view-logs.sh
        ;;
    q|Q)
        echo ""
        echo "Goodbye!"
        exit 0
        ;;
    *)
        echo ""
        echo -e "${RED}Invalid choice. Please run ./run.sh again.${NC}"
        exit 1
        ;;
esac
