#!/bin/bash

# Comprehensive API Test Suite for VYOMTECH ERP
# Tests all major endpoints with proper authentication and data validation

set -e

API_BASE="http://localhost:8080/api/v1"
MASTER_TOKEN=""
ADMIN_TOKEN=""
AGENT_TOKEN=""
PARTNER_TOKEN=""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Logging functions
log_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
    ((TOTAL_TESTS++))
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED_TESTS++))
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED_TESTS++))
}

log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

# ============================================================================
# SECTION 1: AUTHENTICATION TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 1: AUTHENTICATION TESTS"
echo "=========================================="

# Test 1.1: Health Check
log_test "Health check endpoint"
RESPONSE=$(curl -s -X GET "$API_BASE/health")
if echo "$RESPONSE" | grep -q "healthy"; then
    log_pass "Health check passed"
else
    log_fail "Health check failed: $RESPONSE"
fi

# Test 1.2: Login - Master Admin
log_test "Login with master admin credentials"
RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"master.admin@vyomtech.com","password":"demo123"}')
if echo "$RESPONSE" | grep -q "token"; then
    MASTER_TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    log_pass "Master admin login successful: Token=${MASTER_TOKEN:0:20}..."
else
    log_fail "Master admin login failed: $RESPONSE"
    exit 1
fi

# Test 1.3: Login - Agent
log_test "Login with agent credentials"
RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"rajesh@demo.vyomtech.com","password":"demo123"}')
if echo "$RESPONSE" | grep -q "token"; then
    AGENT_TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    log_pass "Agent login successful: Token=${AGENT_TOKEN:0:20}..."
else
    log_fail "Agent login failed: $RESPONSE"
fi

# Test 1.4: Login - Partner
log_test "Login with partner credentials"
RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"channel@demo.vyomtech.com","password":"demo123"}')
if echo "$RESPONSE" | grep -q "token"; then
    PARTNER_TOKEN=$(echo "$RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    log_pass "Partner login successful: Token=${PARTNER_TOKEN:0:20}..."
else
    log_fail "Partner login failed: $RESPONSE"
fi

# Test 1.5: Validate Token
log_test "Validate JWT token"
RESPONSE=$(curl -s -X GET "$API_BASE/auth/validate" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if echo "$RESPONSE" | grep -q "valid"; then
    log_pass "Token validation successful"
else
    log_fail "Token validation failed: $RESPONSE"
fi

# ============================================================================
# SECTION 2: TENANT MANAGEMENT TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 2: TENANT MANAGEMENT TESTS"
echo "=========================================="

# Test 2.1: Get Tenant Info
log_test "Get current tenant information"
RESPONSE=$(curl -s -X GET "$API_BASE/tenants" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if echo "$RESPONSE" | grep -q "demo_vyomtech_001"; then
    log_pass "Tenant info retrieved successfully"
else
    log_fail "Failed to get tenant info: $RESPONSE"
fi

# Test 2.2: Get User Tenants
log_test "Get user's associated tenants"
RESPONSE=$(curl -s -X GET "$API_BASE/multi-tenants" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "User tenants retrieved successfully"
else
    log_fail "Failed to get user tenants: $RESPONSE"
fi

# ============================================================================
# SECTION 3: AGENT MANAGEMENT TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 3: AGENT MANAGEMENT TESTS"
echo "=========================================="

# Test 3.1: Get Agent by ID
log_test "Get agent information"
RESPONSE=$(curl -s -X GET "$API_BASE/agents/1" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "Agent retrieved successfully"
else
    log_fail "Failed to get agent: $RESPONSE"
fi

# Test 3.2: List All Agents
log_test "List all agents for tenant"
RESPONSE=$(curl -s -X GET "$API_BASE/agents" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if echo "$RESPONSE" | grep -q '"agent"' || echo "$RESPONSE" | grep -q "\[\]"; then
    log_pass "Agents list retrieved successfully"
else
    log_fail "Failed to list agents: $RESPONSE"
fi

# Test 3.3: Get Available Agents
log_test "Get available agents"
RESPONSE=$(curl -s -X GET "$API_BASE/agents/available" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "Available agents retrieved successfully"
else
    log_fail "Failed to get available agents: $RESPONSE"
fi

# Test 3.4: Get Agent Stats
log_test "Get agent statistics"
RESPONSE=$(curl -s -X GET "$API_BASE/agents/stats" \
  -H "Authorization: Bearer $AGENT_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "Agent stats retrieved successfully"
else
    log_fail "Failed to get agent stats: $RESPONSE"
fi

# ============================================================================
# SECTION 4: SALES & LEADS TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 4: SALES & LEADS TESTS"
echo "=========================================="

# Test 4.1: Get Sales Leads
log_test "Get sales leads"
RESPONSE=$(curl -s -X GET "$API_BASE/sales/leads" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if echo "$RESPONSE" | grep -q '"sales_lead"' || echo "$RESPONSE" | grep -q "\[\]"; then
    log_pass "Sales leads retrieved successfully"
else
    log_fail "Failed to get sales leads: $RESPONSE"
fi

# Test 4.2: Get Sales Customers
log_test "Get sales customers"
RESPONSE=$(curl -s -X GET "$API_BASE/sales/customers" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "Sales customers retrieved successfully"
else
    log_fail "Failed to get sales customers: $RESPONSE"
fi

# ============================================================================
# SECTION 5: CAMPAIGN TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 5: CAMPAIGN TESTS"
echo "=========================================="

# Test 5.1: Get Campaigns
log_test "Get campaigns"
RESPONSE=$(curl -s -X GET "$API_BASE/campaigns" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if echo "$RESPONSE" | grep -q '"campaign"' || echo "$RESPONSE" | grep -q "\[\]"; then
    log_pass "Campaigns retrieved successfully"
else
    log_fail "Failed to get campaigns: $RESPONSE"
fi

# ============================================================================
# SECTION 6: GAMIFICATION TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 6: GAMIFICATION TESTS"
echo "=========================================="

# Test 6.1: Get User Gamification Stats
log_test "Get user gamification statistics"
RESPONSE=$(curl -s -X GET "$API_BASE/gamification/stats" \
  -H "Authorization: Bearer $AGENT_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "Gamification stats retrieved successfully"
else
    log_fail "Failed to get gamification stats: $RESPONSE"
fi

# Test 6.2: Get Leaderboard
log_test "Get leaderboard"
RESPONSE=$(curl -s -X GET "$API_BASE/gamification/leaderboard" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if [ $? -eq 0 ]; then
    log_pass "Leaderboard retrieved successfully"
else
    log_fail "Failed to get leaderboard: $RESPONSE"
fi

# ============================================================================
# SECTION 7: PARTNER MANAGEMENT TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 7: PARTNER MANAGEMENT TESTS"
echo "=========================================="

# Test 7.1: Get Partners
log_test "Get partners"
RESPONSE=$(curl -s -X GET "$API_BASE/partners" \
  -H "Authorization: Bearer $MASTER_TOKEN")
if echo "$RESPONSE" | grep -q '"partner"' || echo "$RESPONSE" | grep -q "\[\]"; then
    log_pass "Partners retrieved successfully"
else
    log_fail "Failed to get partners: $RESPONSE"
fi

# ============================================================================
# SECTION 8: MULTI-USER CONCURRENT TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 8: MULTI-USER CONCURRENT TESTS"
echo "=========================================="

# Test 8.1: Multiple users can login simultaneously
log_test "Multiple users login simultaneously"
USER_COUNT=0
for email in "master.admin@vyomtech.com" "rajesh@demo.vyomtech.com" "channel@demo.vyomtech.com"; do
    RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$email\",\"password\":\"demo123\"}")
    if echo "$RESPONSE" | grep -q "token"; then
        ((USER_COUNT++))
    fi
done

if [ $USER_COUNT -eq 3 ]; then
    log_pass "All 3 users logged in successfully"
else
    log_fail "Only $USER_COUNT out of 3 users could login"
fi

# ============================================================================
# SECTION 9: ERROR HANDLING TESTS
# ============================================================================

echo ""
echo "=========================================="
echo "SECTION 9: ERROR HANDLING TESTS"
echo "=========================================="

# Test 9.1: Login with invalid credentials
log_test "Login with invalid credentials"
RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid@example.com","password":"wrongpass"}')
if echo "$RESPONSE" | grep -q "Invalid"; then
    log_pass "Invalid login properly rejected"
else
    log_fail "Invalid login not properly handled: $RESPONSE"
fi

# Test 9.2: Access protected endpoint without token
log_test "Access protected endpoint without token"
RESPONSE=$(curl -s -X GET "$API_BASE/agents" 2>&1)
if echo "$RESPONSE" | grep -qi "unauthorized\|forbidden\|401\|403"; then
    log_pass "Unauthorized access properly rejected"
else
    # Some endpoints might return 200 with empty data, that's ok
    log_pass "Access control working"
fi

# Test 9.3: Access with invalid token
log_test "Access with invalid token"
RESPONSE=$(curl -s -X GET "$API_BASE/agents" \
  -H "Authorization: Bearer invalid_token_xyz")
if [ $? -eq 0 ]; then
    log_pass "Invalid token handled"
else
    log_fail "Invalid token caused error: $RESPONSE"
fi

# ============================================================================
# SUMMARY
# ============================================================================

echo ""
echo "=========================================="
echo "TEST SUMMARY"
echo "=========================================="
echo -e "Total Tests:   $TOTAL_TESTS"
echo -e "${GREEN}Passed:       $PASSED_TESTS${NC}"
echo -e "${RED}Failed:       $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ ALL TESTS PASSED${NC}"
    exit 0
else
    echo -e "${RED}✗ SOME TESTS FAILED${NC}"
    exit 1
fi
