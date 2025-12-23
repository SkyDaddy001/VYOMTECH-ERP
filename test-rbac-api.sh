#!/bin/bash

# ============================================================
# VYOMTECH ERP - RBAC API TEST SUITE
# Phase 3.4: Permission Management APIs
# ============================================================

# Configuration
API_BASE_URL="http://localhost:8080/api/v1"
TENANT_ID="test-tenant-001"
ADMIN_TOKEN="your-admin-token-here"
MANAGER_TOKEN="your-manager-token-here"
USER_TOKEN="your-user-token-here"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
test_endpoint() {
	local name=$1
	local method=$2
	local url=$3
	local token=$4
	local data=$5
	local expected_status=$6

	TESTS_RUN=$((TESTS_RUN + 1))

	echo -e "\n${BLUE}Test $TESTS_RUN: $name${NC}"
	echo "Method: $method"
	echo "URL: $url"

	# Build curl command
	local curl_cmd="curl -s -w '\n%{http_code}' -X $method"
	curl_cmd="$curl_cmd -H 'Content-Type: application/json'"
	curl_cmd="$curl_cmd -H 'Authorization: Bearer $token'"
	curl_cmd="$curl_cmd -H 'X-Tenant-ID: $TENANT_ID'"

	if [ -n "$data" ]; then
		curl_cmd="$curl_cmd -d '$data'"
	fi

	curl_cmd="$curl_cmd '$url'"

	# Execute and capture response
	response=$(eval $curl_cmd)
	status=$(echo "$response" | tail -n1)
	body=$(echo "$response" | head -n-1)

	echo "Status: $status (expected: $expected_status)"
	echo "Response: $body"

	if [ "$status" -eq "$expected_status" ]; then
		echo -e "${GREEN}✓ PASSED${NC}"
		TESTS_PASSED=$((TESTS_PASSED + 1))
	else
		echo -e "${RED}✗ FAILED${NC}"
		TESTS_FAILED=$((TESTS_FAILED + 1))
	fi
}

# ============================================================
# TEST SUITE
# ============================================================

echo "============================================================"
echo "VYOMTECH ERP - RBAC API TEST SUITE"
echo "============================================================"

# 1. List Roles (GET /api/v1/rbac/roles)
echo -e "\n\n${BLUE}=== TEST SET 1: LIST OPERATIONS ===${NC}"

test_endpoint \
	"List all roles (authenticated user)" \
	"GET" \
	"$API_BASE_URL/rbac/roles" \
	"$MANAGER_TOKEN" \
	"" \
	200

# 2. List Permissions (GET /api/v1/rbac/permissions)
test_endpoint \
	"List all permissions" \
	"GET" \
	"$API_BASE_URL/rbac/permissions" \
	"$MANAGER_TOKEN" \
	"" \
	200

# 3. Get Specific Role (GET /api/v1/rbac/roles/:id)
test_endpoint \
	"Get specific role with permissions" \
	"GET" \
	"$API_BASE_URL/rbac/roles/role-id-001" \
	"$MANAGER_TOKEN" \
	"" \
	200

# ============================================================
# TEST SET 2: CREATE OPERATIONS (ADMIN ONLY)
# ============================================================

echo -e "\n\n${BLUE}=== TEST SET 2: CREATE OPERATIONS ===${NC}"

# 4. Create Role (admin only)
test_endpoint \
	"Create new role (admin access)" \
	"POST" \
	"$API_BASE_URL/rbac/roles" \
	"$ADMIN_TOKEN" \
	'{"role_name":"Sales Manager","description":"Sales team manager role","is_active":true}' \
	201

# 5. Create Role (non-admin denied)
test_endpoint \
	"Create new role (non-admin denied)" \
	"POST" \
	"$API_BASE_URL/rbac/roles" \
	"$MANAGER_TOKEN" \
	'{"role_name":"Unauthorized Role","description":"Should fail","is_active":true}' \
	403

# 6. Create Role without authorization header
test_endpoint \
	"Create role without auth (should fail)" \
	"POST" \
	"$API_BASE_URL/rbac/roles" \
	"" \
	'{"role_name":"No Auth Role","description":"No auth","is_active":true}' \
	401

# ============================================================
# TEST SET 3: ASSIGN PERMISSIONS (ADMIN ONLY)
# ============================================================

echo -e "\n\n${BLUE}=== TEST SET 3: ASSIGN PERMISSIONS ===${NC}"

# 7. Assign permissions to role (admin)
test_endpoint \
	"Assign permissions to role (admin)" \
	"PUT" \
	"$API_BASE_URL/rbac/roles/role-id-001/permissions" \
	"$ADMIN_TOKEN" \
	'{"permission_ids":["perm-001","perm-002","perm-003"]}' \
	200

# 8. Assign permissions (non-admin denied)
test_endpoint \
	"Assign permissions (non-admin denied)" \
	"PUT" \
	"$API_BASE_URL/rbac/roles/role-id-001/permissions" \
	"$MANAGER_TOKEN" \
	'{"permission_ids":["perm-001"]}' \
	403

# ============================================================
# TEST SET 4: DELETE OPERATIONS (ADMIN ONLY)
# ============================================================

echo -e "\n\n${BLUE}=== TEST SET 4: DELETE OPERATIONS ===${NC}"

# 9. Delete role (admin)
test_endpoint \
	"Delete role (admin)" \
	"DELETE" \
	"$API_BASE_URL/rbac/roles/role-id-temp" \
	"$ADMIN_TOKEN" \
	"" \
	200

# 10. Delete role (non-admin denied)
test_endpoint \
	"Delete role (non-admin denied)" \
	"DELETE" \
	"$API_BASE_URL/rbac/roles/role-id-001" \
	"$MANAGER_TOKEN" \
	"" \
	403

# ============================================================
# TEST SET 5: ERROR HANDLING
# ============================================================

echo -e "\n\n${BLUE}=== TEST SET 5: ERROR HANDLING ===${NC}"

# 11. Create role with missing required field
test_endpoint \
	"Create role with missing role_name" \
	"POST" \
	"$API_BASE_URL/rbac/roles" \
	"$ADMIN_TOKEN" \
	'{"description":"Missing role name","is_active":true}' \
	400

# 12. Get non-existent role
test_endpoint \
	"Get non-existent role" \
	"GET" \
	"$API_BASE_URL/rbac/roles/nonexistent-id" \
	"$MANAGER_TOKEN" \
	"" \
	404

# 13. Assign non-existent permission
test_endpoint \
	"Assign non-existent permission" \
	"PUT" \
	"$API_BASE_URL/rbac/roles/role-id-001/permissions" \
	"$ADMIN_TOKEN" \
	'{"permission_ids":["nonexistent-perm"]}' \
	400

# ============================================================
# TEST SUMMARY
# ============================================================

echo -e "\n\n${BLUE}============================================================${NC}"
echo "TEST SUMMARY"
echo -e "${BLUE}============================================================${NC}"
echo "Total Tests Run: $TESTS_RUN"
echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Failed: ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
	echo -e "\n${GREEN}All tests passed!${NC}"
	exit 0
else
	echo -e "\n${RED}Some tests failed!${NC}"
	exit 1
fi
