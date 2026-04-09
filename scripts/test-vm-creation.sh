#!/bin/bash

# Test VM Creation After Fix
# Usage: ./scripts/test-vm-creation.sh

set -e

echo "🔍 Testing VM Creation After Column Fix"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:3000"

# Check if backend is running
if curl -s "${BASE_URL}/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend is not running${NC}"
    echo "Please start backend: cd backend && go run main.go"
    exit 1
fi

echo ""
echo "📝 Step 1: Login"
echo "----------------"

# Login
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}')

echo "$LOGIN_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$LOGIN_RESPONSE"

# Check if login successful
SUCCESS=$(echo "$LOGIN_RESPONSE" | grep -o '"success":true' || echo "")
if [ -z "$SUCCESS" ]; then
    echo -e "${RED}✗ Login failed${NC}"
    echo "Creating test user..."

    # Register
    REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/auth/register" \
      -H "Content-Type: application/json" \
      -d '{"username":"testuser","email":"test@test.com","password":"test12345"}')

    echo "$REGISTER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$REGISTER_RESPONSE"

    # Login again
    LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/auth/login" \
      -H "Content-Type: application/json" \
      -d '{"email":"test@test.com","password":"test12345"}')
fi

# Extract token
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ Failed to get token${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}"
echo ""

echo "📝 Step 2: Get Available Plans"
echo "-------------------------------"

# Get plans
PLANS_RESPONSE=$(curl -s "${BASE_URL}/api/v1/billing/plans" \
  -H "Authorization: Bearer $TOKEN")

echo "$PLANS_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$PLANS_RESPONSE"

# Extract plan ID
PLAN_ID=$(echo "$PLANS_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$PLAN_ID" ]; then
    echo -e "${RED}✗ No plans available${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Found plan ID: $PLAN_ID${NC}"
echo ""

echo "📝 Step 3: Create VM"
echo "--------------------"

VM_NAME="test-vm-$(date +%s)"

echo "Creating VM: $VM_NAME"
echo ""

VM_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/vms" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"name\": \"$VM_NAME\",
    \"hostname\": \"$VM_NAME\",
    \"plan_id\": $PLAN_ID
  }")

echo "$VM_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$VM_RESPONSE"

# Check for column error
COLUMN_ERROR=$(echo "$VM_RESPONSE" | grep -i "vm_id.*does not exist" || echo "")

if [ -n "$COLUMN_ERROR" ]; then
    echo ""
    echo -e "${RED}✗ VM creation failed with column error!${NC}"
    echo -e "${RED}  Error: $COLUMN_ERROR${NC}"
    echo ""
    echo "The fix didn't work. Please check:"
    echo "  1. backend/models/vm.go has 'column:vmid' tag"
    echo "  2. Backend server was restarted"
    exit 1
fi

# Check if successful
VM_SUCCESS=$(echo "$VM_RESPONSE" | grep -o '"success":true' || echo "")

if [ -n "$VM_SUCCESS" ]; then
    echo ""
    echo -e "${GREEN}✅ VM creation successful!${NC}"
    echo ""
    echo "🎉 The column fix is working correctly!"
else
    echo ""
    echo -e "${YELLOW}⚠️ VM creation returned non-success status${NC}"
    echo "Check the response above for error details"
fi

echo ""
echo "========================================="
echo "✅ Test complete!"
