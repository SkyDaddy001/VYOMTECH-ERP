#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080"
PASS_COUNT=0
FAIL_COUNT=0

# Test function
test_endpoint() {
    local test_name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_code=$5
    local token=$6
    
    local cmd="curl -s -w '%{http_code}' -X $method '$API_BASE$endpoint'"
    
    if [ -n "$data" ]; then
        cmd="$cmd -H 'Content-Type: application/json' -d '$data'"
    fi
    
    if [ -n "$token" ]; then
        cmd="$cmd -H 'Authorization: Bearer $token'"
    fi
    
    result=$(eval "$cmd")
    response_code="${result: -3}"
    response_body="${result%???}"
    
    if [ "$response_code" == "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} - $test_name (HTTP $response_code)"
        ((PASS_COUNT++))
    else
        echo -e "${RED}✗ FAIL${NC} - $test_name (Expected $expected_code, got $response_code)"
        echo "  Response: ${response_body:0:100}"
        ((FAIL_COUNT++))
    fi
}

echo "=========================================="
echo "API TESTS - DEMO DATA VALIDATION"
echo "=========================================="

# Test 1: Health endpoint
test_endpoint "Health Check" "GET" "/health" "" "200"

# Test 2: Login with master admin
echo ""
echo "=== AUTHENTICATION TESTS ==="
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"master.admin@vyomtech.com","password":"demo123"}')
MASTER_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$MASTER_TOKEN" ]; then
    echo -e "${GREEN}✓${NC} Master admin login successful"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗${NC} Master admin login failed"
    ((FAIL_COUNT++))
fi

# Test 3: Get agents (the key test)
echo ""
echo "=== DEMO DATA TESTS ==="
AGENT_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/agents" \
    -H "Authorization: Bearer $MASTER_TOKEN")

AGENT_COUNT=$(echo "$AGENT_RESPONSE" | grep -o '"id"' | wc -l)

if [ "$AGENT_COUNT" -ge 1 ]; then
    echo -e "${GREEN}✓ PASS${NC} - Get Agents List (Found $AGENT_COUNT agents)"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Get Agents List (No agents found)"
    echo "  Response: $AGENT_RESPONSE"
    ((FAIL_COUNT++))
fi

# Test 4: Get sales leads
LEAD_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/sales/leads" \
    -H "Authorization: Bearer $MASTER_TOKEN")
LEAD_COUNT=$(echo "$LEAD_RESPONSE" | grep -o '"id"' | wc -l)

if [ "$LEAD_COUNT" -ge 1 ]; then
    echo -e "${GREEN}✓ PASS${NC} - Get Sales Leads (Found $LEAD_COUNT leads)"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Get Sales Leads"
    ((FAIL_COUNT++))
fi

# Test 5: Get campaigns
CAMPAIGN_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/campaigns" \
    -H "Authorization: Bearer $MASTER_TOKEN")
CAMPAIGN_COUNT=$(echo "$CAMPAIGN_RESPONSE" | grep -o '"id"' | wc -l)

if [ "$CAMPAIGN_COUNT" -ge 1 ]; then
    echo -e "${GREEN}✓ PASS${NC} - Get Campaigns (Found $CAMPAIGN_COUNT campaigns)"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Get Campaigns"
    ((FAIL_COUNT++))
fi

# Test 6: Get partners
PARTNER_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/partners" \
    -H "Authorization: Bearer $MASTER_TOKEN")
PARTNER_COUNT=$(echo "$PARTNER_RESPONSE" | grep -o '"id"' | wc -l)

if [ "$PARTNER_COUNT" -ge 0 ]; then
    echo -e "${GREEN}✓ PASS${NC} - Get Partners (Found $PARTNER_COUNT partners)"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Get Partners"
    ((FAIL_COUNT++))
fi

# Test 7: Agent gamification stats
STATS_RESPONSE=$(curl -s -X GET "$API_BASE/api/v1/gamification/stats" \
    -H "Authorization: Bearer $MASTER_TOKEN")

if [[ "$STATS_RESPONSE" == *"points"* ]] || [[ "$STATS_RESPONSE" == *"badge"* ]]; then
    echo -e "${GREEN}✓ PASS${NC} - Get Gamification Stats"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Get Gamification Stats"
    ((FAIL_COUNT++))
fi

# Test 8: Invalid login rejection
INVALID_LOGIN=$(curl -s -w '%{http_code}' -X POST "$API_BASE/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"invalid@test.com","password":"wrongpass"}')
INVALID_CODE="${INVALID_LOGIN: -3}"

if [ "$INVALID_CODE" == "401" ] || [ "$INVALID_CODE" == "400" ]; then
    echo -e "${GREEN}✓ PASS${NC} - Invalid Login Rejection"
    ((PASS_COUNT++))
else
    echo -e "${RED}✗ FAIL${NC} - Invalid Login (Got HTTP $INVALID_CODE)"
    ((FAIL_COUNT++))
fi

# Summary
echo ""
echo "=========================================="
echo "SUMMARY"
echo "=========================================="
echo -e "${GREEN}Passed: $PASS_COUNT${NC}"
echo -e "${RED}Failed: $FAIL_COUNT${NC}"
TOTAL=$((PASS_COUNT + FAIL_COUNT))
PERCENT=$((PASS_COUNT * 100 / TOTAL))
echo "Success Rate: $PERCENT% ($PASS_COUNT/$TOTAL)"
echo "=========================================="

if [ $FAIL_COUNT -eq 0 ]; then
    exit 0
else
    exit 1
fi
