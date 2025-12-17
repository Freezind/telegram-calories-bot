#!/bin/bash

# Unified Telegram Bot + Mini App Backend
# Runs both bot and HTTP API server in a single Go process with shared storage

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Unified Telegram Bot + Mini App Backend${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    echo ""
    echo "Please create a .env file with the following variables:"
    echo "  TELEGRAM_BOT_TOKEN=your_bot_token_here"
    echo "  GEMINI_API_KEY=your_gemini_api_key_here"
    echo "  PORT=8080 (optional, defaults to 8080)"
    echo "  TUNNEL_URL=https://your-tunnel-url (optional, for Cloudflare Tunnel)"
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

# Load environment variables
export $(grep -v '^#' .env | xargs)

# Verify required variables
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo -e "${RED}Error: TELEGRAM_BOT_TOKEN not set in .env${NC}"
    exit 1
fi

if [ -z "$GEMINI_API_KEY" ]; then
    echo -e "${RED}Error: GEMINI_API_KEY not set in .env${NC}"
    exit 1
fi

echo -e "${BLUE}Configuration:${NC}"
echo "  HTTP Server: http://localhost:${PORT:-8080}"
echo "  Bot Token: ${TELEGRAM_BOT_TOKEN:0:10}..."
echo "  Gemini API: ${GEMINI_API_KEY:0:10}..."
if [ -n "$TUNNEL_URL" ]; then
    echo "  Tunnel URL: $TUNNEL_URL"
fi
echo ""

echo -e "${GREEN}Starting unified backend...${NC}"
echo -e "${YELLOW}This process runs BOTH:${NC}"
echo "  1. Telegram bot (/estimate command)"
echo "  2. HTTP API server (miniapp CRUD)"
echo "  3. Shared in-memory storage"
echo ""
echo -e "${YELLOW}Logs created via /estimate will appear in the miniapp!${NC}"
echo ""

# Run the unified backend
go run cmd/unified/main.go
