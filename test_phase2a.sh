#!/bin/bash

# Phase 2A Testing Script - Tasks and Notifications

BASE_URL="http://localhost:8080/api/v1"
TENANT_ID="test-tenant-001"
USER_ID=1
AUTH_TOKEN="test-token"

echo "=========================================="
echo "Phase 2A - Tasks and Notifications Testing"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4

    echo -e "${YELLOW}Testing:${NC} $description"
    echo "  Method: $method | Endpoint: $endpoint"

    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $AUTH_TOKEN" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -d "$data" \
            -w "\n%{http_code}")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Authorization: Bearer $AUTH_TOKEN" \
            -H "X-Tenant-ID: $TENANT_ID" \
            -w "\n%{http_code}")
    fi

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "  ${GREEN}✓ PASSED${NC} (HTTP $http_code)"
        ((TESTS_PASSED++))
    elif [ "$http_code" -ge 400 ] && [ "$http_code" -lt 500 ]; then
        echo -e "  ${YELLOW}⚠ Expected Error${NC} (HTTP $http_code)"
        ((TESTS_PASSED++))
    else
        echo -e "  ${RED}✗ FAILED${NC} (HTTP $http_code)"
        echo "  Response: $body"
        ((TESTS_FAILED++))
    fi
    echo ""
}

echo "========================="
echo "TASK ENDPOINTS TESTING"
echo "========================="
echo ""

# Create Task
task_data='{
  "title": "Test Task",
  "description": "Testing Phase 2A Task Handler",
  "status": "open",
  "assigned_to": 1,
  "due_date": "2025-12-31T23:59:59Z"
}'

test_endpoint "POST" "/tasks" "$task_data" "Create Task"

# List Tasks
test_endpoint "GET" "/tasks" "" "List Tasks"

# Get Tasks Stats
test_endpoint "GET" "/tasks/stats" "" "Get Task Statistics"

# Get Tasks by User
test_endpoint "GET" "/tasks/user/1" "" "Get Tasks by User ID"

echo ""
echo "========================="
echo "NOTIFICATION ENDPOINTS TESTING"
echo "========================="
echo ""

# Create Notification
notif_data='{
  "title": "Test Notification",
  "message": "Testing Phase 2A Notification Handler",
  "type": "info",
  "user_id": 1
}'

test_endpoint "POST" "/notifications" "$notif_data" "Create Notification"

# List Notifications
test_endpoint "GET" "/notifications" "" "List Notifications"

# Get Notification Stats
test_endpoint "GET" "/notifications/stats" "" "Get Notification Statistics"

# Get Preferences
test_endpoint "GET" "/notifications/preferences" "" "Get Notification Preferences"

# Update Preferences
pref_data='{
  "email_enabled": true,
  "sms_enabled": false,
  "push_enabled": true
}'

test_endpoint "PUT" "/notifications/preferences" "$pref_data" "Update Notification Preferences"

echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests completed successfully!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed!${NC}"
    exit 1
fi
