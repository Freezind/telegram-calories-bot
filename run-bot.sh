#!/bin/bash

# Load environment variables from .env file and run the bot
# Usage: ./run-bot.sh

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found"
    echo "Please create a .env file with the following content:"
    echo ""
    echo "TELEGRAM_BOT_TOKEN=your_bot_token_here"
    echo "GEMINI_API_KEY=your_gemini_api_key_here"
    echo ""
    exit 1
fi

# Load environment variables from .env
echo "Loading environment variables from .env..."
export $(grep -v '^#' .env | xargs)

# Verify required variables are set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "Error: TELEGRAM_BOT_TOKEN not set in .env"
    exit 1
fi

if [ -z "$GEMINI_API_KEY" ]; then
    echo "Error: GEMINI_API_KEY not set in .env"
    exit 1
fi

echo "Environment variables loaded successfully"
echo "Starting Telegram Calories Bot..."
echo ""

# Run the bot
go run src/main.go
