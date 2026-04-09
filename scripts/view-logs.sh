#!/bin/bash

# TerasVPS - View API Logs in Real-Time
# This script helps you monitor API calls with beautiful formatting

set -e

echo "🔍 TerasVPS API Log Monitor"
echo "==========================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if backend is running
if pgrep -f "teras-vps-api" > /dev/null || pgrep -f "go run main.go" > /dev/null; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${YELLOW}⚠ Backend is not running${NC}"
    echo ""
    echo "Starting backend server..."
    cd "$(dirname "$0")/../backend"
    go run main.go &
    sleep 3
fi

echo ""
echo -e "${BLUE}📊 Watching API logs...${NC}"
echo -e "${BLUE}Press Ctrl+C to stop watching${NC}"
echo ""
echo "═══════════════════════════════════════════════════════════"
echo "💡 TIP: All API calls are already logged in the backend terminal!"
echo "   You don't need this script - just watch your backend console"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Optional: If you want to tail a log file (if you redirect output to file)
LOG_FILE="../logs/api.log"

if [ -f "$LOG_FILE" ]; then
    echo -e "${GREEN}📄 Tailing log file: $LOG_FILE${NC}"
    echo ""
    tail -f "$LOG_FILE"
else
    echo -e "${YELLOW}⚠ No log file found at: $LOG_FILE${NC}"
    echo ""
    echo "Logs are shown directly in the backend terminal where you ran:"
    echo "  cd backend && go run main.go"
    echo ""
    echo "To save logs to a file, run backend like this:"
    echo "  cd backend && go run main.go 2>&1 | tee ../logs/api.log"
    echo ""
    echo -e "${GREEN}✨ Your backend is already logging beautifully! Just watch the console.${NC}"
fi
