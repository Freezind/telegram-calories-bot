#!/bin/bash
# LLM-Based Bot Testing Runner
# This script runs the automated LLM-based test suite for the Telegram bot

set -e

echo "========================================"
echo "LLM-Based Bot Testing Suite"
echo "========================================"
echo ""

# Check if .env.test exists
if [ ! -f ".env.test" ]; then
    echo "Error: .env.test file not found"
    echo ""
    echo "Please create .env.test from .env.test.example:"
    echo "  cp .env.test.example .env.test"
    echo "  # Edit .env.test and fill in the required values"
    echo ""
    exit 1
fi

# Load environment variables
echo "Loading environment variables from .env.test..."
set -a
source .env.test
set +a

# Validate required environment variables
required_vars=("TELEGRAM_BOT_TOKEN" "TELEGRAM_TEST_CHAT_ID" "MINIAPP_URL" "GEMINI_API_KEY")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -gt 0 ]; then
    echo "Error: Missing required environment variables:"
    for var in "${missing_vars[@]}"; do
        echo "  - $var"
    done
    echo ""
    echo "Please update .env.test with all required values"
    exit 1
fi

# Check if test image exists
if [ ! -f "${TEST_FOOD_IMAGE_PATH:-tests/fixtures/food.jpg}" ]; then
    echo "Warning: Test image not found at ${TEST_FOOD_IMAGE_PATH:-tests/fixtures/food.jpg}"
    echo "Please add a food image to tests/fixtures/food.jpg"
    echo "See tests/fixtures/README.md for details"
    echo ""
fi

# Install Playwright browsers if not already installed
if [ ! -d "$HOME/.cache/ms-playwright" ]; then
    echo "Installing Playwright browsers (first-time setup)..."
    go run github.com/playwright-community/playwright-go/cmd/playwright@latest install chromium
    echo ""
fi

# Run tests
echo "Running LLM-based test suite..."
echo ""

go run cmd/tester/*.go

# Exit with the same code as the test runner
exit $?
