#!/bin/bash

# Check Telegram Bot Configuration

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Checking Telegram Bot Configuration${NC}"
echo ""

# Load .env
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    exit 1
fi

export $(grep -v '^#' .env | xargs)

if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo -e "${RED}Error: TELEGRAM_BOT_TOKEN not set in .env${NC}"
    exit 1
fi

echo "Bot Token: ${TELEGRAM_BOT_TOKEN:0:10}...${TELEGRAM_BOT_TOKEN: -5}"
echo ""

echo "Fetching bot info..."
curl -s "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getMe" | jq '.'
echo ""

echo "Fetching menu button configuration..."
MENU=$(curl -s "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/getChatMenuButton")
echo "$MENU" | jq '.'

# Check if menu button is configured
if echo "$MENU" | grep -q '"type":"web_app"'; then
    echo ""
    echo -e "${GREEN}✓ Menu button is configured!${NC}"
    URL=$(echo "$MENU" | jq -r '.result.menu_button.web_app.url')
    echo -e "${GREEN}URL: $URL${NC}"
else
    echo ""
    echo -e "${RED}✗ Menu button is NOT configured or is default!${NC}"
    echo ""
    echo "To configure it, run:"
    echo ""
    echo -e "${YELLOW}curl -X POST \"https://api.telegram.org/bot\${TELEGRAM_BOT_TOKEN}/setChatMenuButton\" \\${NC}"
    echo -e "${YELLOW}  -H \"Content-Type: application/json\" \\${NC}"
    echo -e "${YELLOW}  -d '{${NC}"
    echo -e "${YELLOW}    \"menu_button\": {${NC}"
    echo -e "${YELLOW}      \"type\": \"web_app\",${NC}"
    echo -e "${YELLOW}      \"text\": \"Open Calorie Logs\",${NC}"
    echo -e "${YELLOW}      \"web_app\": {${NC}"
    echo -e "${YELLOW}        \"url\": \"YOUR_TUNNEL_URL\"${NC}"
    echo -e "${YELLOW}      }${NC}"
    echo -e "${YELLOW}    }${NC}"
    echo -e "${YELLOW}  }'${NC}"
fi
