#!/bin/bash

# TerasVPS - Quick Start Script for Testing Auth Fix
# Usage: ./scripts/test-auth.sh

set -e

echo "🔐 TerasVPS Authentication Fix - Quick Test"
echo "=============================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if backend is running
echo "📡 Checking backend server..."
if pgrep -f "teras-vps-api" > /dev/null; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${YELLOW}⚠ Backend is not running. Starting...${NC}"
    cd backend
    go build -o teras-vps-api .
    ./teras-vps-api &
    cd ..
    sleep 2
fi

# Check if frontend is running
echo "🌐 Checking frontend dev server..."
if pgrep -f "astro dev" > /dev/null || pgrep -f "astro" > /dev/null; then
    echo -e "${GREEN}✓ Frontend is running${NC}"
else
    echo -e "${YELLOW}⚠ Frontend is not running. Please start it manually:${NC}"
    echo "   cd frontend && npm run dev"
fi

echo ""
echo "📝 Test Steps:"
echo "==============="
echo "1. Open browser: http://localhost:4321/auth/login"
echo "2. Login with your credentials"
echo "3. Check DevTools → Network tab"
echo "4. Verify token is sent with requests"
echo "5. Navigate to /dashboard/billing"
echo "6. Check that /api/v1/billing/plans returns 200 (not 401)"
echo ""

echo "🔍 Quick Token Check:"
echo "===================="
echo "To manually test if token is being sent:"
echo ""
echo "  # Get token from localStorage after login"
echo "  TOKEN=$(cat << 'EOF'
  <run this in browser console: localStorage.getItem('token')>
  EOF
  )"
echo ""
echo "  # Test API with token"
echo "  curl -H \"Authorization: Bearer \$TOKEN\" http://localhost:3000/api/v1/billing/plans"
echo ""

echo -e "${GREEN}✨ Auth fix is ready!${NC}"
echo "For detailed information, see: AUTH_FIX_SUMMARY.md"
