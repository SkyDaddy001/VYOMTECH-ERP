#!/bin/bash

# VYOMTECH API Testing Script
# Tests all endpoints and verifies responses

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

API_URL="http://localhost:8080"
FRONTEND_URL="http://localhost:3000"
TENANT_ID="default"
JWT_TOKEN=""
USER_ID=""

# Demo credentials
DEMO_EMAIL="demo@vyomtech.com"
DEMO_PASSWORD="DemoPass@123"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}VYOMTECH API Testing Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to make API call and display result
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${YELLOW}Testing: $description${NC}"
    echo -e "  ${BLUE}$method $API_URL$endpoint${NC}"
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -d "$data" \
            "$API_URL$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" -lt 400 ]; then
        echo -e "  ${GREEN}✅ Status: $http_code${NC}"
        echo -e "  Response: $body" | head -c 200
        echo ""
    else
        echo -e "  ${RED}❌ Status: $http_code${NC}"
        echo -e "  Response: $body" | head -c 200
        echo ""
    fi
    
    echo ""
}

# Test 1: Login
echo -e "${YELLOW}Step 1: Testing Authentication${NC}"
echo "=================="
login_data="{\"email\":\"$DEMO_EMAIL\",\"password\":\"$DEMO_PASSWORD\"}"
login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    "$API_URL/auth/login" \
    -d "$login_data")

echo -e "Login Response: $login_response"

JWT_TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
USER_ID=$(echo "$login_response" | grep -o '"user_id":"[^"]*' | cut -d'"' -f4)

if [ -z "$JWT_TOKEN" ]; then
    echo -e "${RED}❌ Failed to get JWT token${NC}"
    exit 1
else
    echo -e "${GREEN}✅ JWT Token obtained${NC}"
    echo "  Token (first 50 chars): ${JWT_TOKEN:0:50}..."
    echo ""
fi

# Test 2: Dashboard Analytics
echo -e "${YELLOW}Step 2: Testing Dashboard Analytics${NC}"
echo "=================="
test_endpoint "GET" "/api/v1/dashboard/analytics" "" "Get Dashboard Analytics"

# Test 3: Agents List
echo -e "${YELLOW}Step 3: Testing Agents Module${NC}"
echo "=================="
test_endpoint "GET" "/api/v1/agents" "" "List Agents"

# Test 4: Sales Data
echo -e "${YELLOW}Step 4: Testing Sales Module${NC}"
echo "=================="
test_endpoint "GET" "/api/v1/sales/orders" "" "List Sales Orders"

# Test 5: Accounting Data
echo -e "${YELLOW}Step 5: Testing Accounting Module${NC}"
echo "=================="
test_endpoint "GET" "/api/v1/accounts/chart-of-accounts" "" "Get Chart of Accounts"

# Test 6: HR Data
echo -e "${YELLOW}Step 6: Testing HR Module${NC}"
echo "=================="
test_endpoint "GET" "/api/v1/hr/employees" "" "List Employees"

# Test 7: Construction Data
echo -e "${YELLOW}Step 7: Testing Construction Module${NC}"
echo "=================="
test_endpoint "GET" "/api/v1/construction/projects" "" "List Construction Projects"

# Test 8: User Profile
echo -e "${YELLOW}Step 8: Testing User Profile${NC}"
echo "=================="
if [ -n "$USER_ID" ]; then
    test_endpoint "GET" "/api/v1/users/$USER_ID" "" "Get User Profile"
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ API Testing Complete${NC}"
echo -e "${GREEN}========================================${NC}"
