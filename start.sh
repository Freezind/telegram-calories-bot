#!/bin/bash

# Start unified backend (bot + miniapp API server)
# Non-interactive startup script for Railway deployment

set -e  # Exit on error

echo "=========================================="
echo "Starting Unified Backend"
echo "=========================================="

# Verify required environment variables
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
  echo "❌ Error: TELEGRAM_BOT_TOKEN is not set"
  exit 1
fi

if [ -z "$GEMINI_API_KEY" ]; then
  echo "❌ Error: GEMINI_API_KEY is not set"
  exit 1
fi

# PORT is automatically provided by Railway
if [ -z "$PORT" ]; then
  export PORT=8080
  echo "ℹ️  PORT not set, using default: 8080"
fi

echo "✓ Environment variables validated"
echo "✓ Port: $PORT"
echo ""

# Build if binary doesn't exist
if [ ! -f "./unified" ]; then
  echo "Building unified backend..."
  go build -o unified cmd/unified/main.go
  echo "✓ Build complete"
  echo ""
fi

# Start the unified backend
echo "Starting unified backend..."
exec ./unified
