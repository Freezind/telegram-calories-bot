#!/bin/bash

# Telegram Mini App Frontend Startup Script
# This script starts ONLY the Vite development server
# NOTE: You must run the backend separately (./run-unified.sh or ./run-dev.sh)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting Telegram Mini App Frontend...${NC}"
echo ""

# Check if node is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed${NC}"
    echo "Please install Node.js 18 or higher: https://nodejs.org/"
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: npm is not installed${NC}"
    echo "Please install npm (comes with Node.js): https://nodejs.org/"
    exit 1
fi

# Change to web directory
if [ ! -d "web" ]; then
    echo -e "${RED}Error: web directory not found${NC}"
    echo "Please run this script from the repository root"
    exit 1
fi

cd web

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}node_modules not found, installing dependencies...${NC}"
    npm install
    echo ""
fi

echo -e "${GREEN}Starting Vite development server...${NC}"
echo -e "${YELLOW}Frontend will be available at: http://localhost:5173${NC}"
echo -e "${YELLOW}API requests will be proxied to: http://localhost:8080${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
echo ""

# Start the Vite dev server
npm run dev
