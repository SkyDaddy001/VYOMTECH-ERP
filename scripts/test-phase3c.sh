#!/bin/bash

# Phase 3C Integration Test Script
# Tests all API endpoints for the Modular Monetization System

BASE_URL="http://localhost:8080/api/v1"
BEARER_TOKEN=""  # Set after login

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================"
echo "Phase 3C Integration Test Suite"
echo "================================"
echo ""

# Helper function to make requests
test_endpoint() {
  local method=$1
  local endpoint=$2
  local data=$3
  local description=$4

  echo -e "${YELLOW}Testing: $description${NC}"
  
  if [ "$method" = "GET" ]; then
    curl -s -X GET "$BASE_URL$endpoint" \
      -H "Authorization: Bearer $BEARER_TOKEN" \
      -H "Content-Type: application/json"
  elif [ "$method" = "POST" ]; then
    curl -s -X POST "$BASE_URL$endpoint" \
      -H "Authorization: Bearer $BEARER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "$data"
  elif [ "$method" = "PUT" ]; then
    curl -s -X PUT "$BASE_URL$endpoint" \
      -H "Authorization: Bearer $BEARER_TOKEN" \
      -H "Content-Type: application/json" \
      -d "$data"
  elif [ "$method" = "DELETE" ]; then
    curl -s -X DELETE "$BASE_URL$endpoint" \
      -H "Authorization: Bearer $BEARER_TOKEN" \
      -H "Content-Type: application/json"
  fi
  
  echo ""
  echo "---"
  echo ""
}

# ==================== MODULE TESTS ====================
echo -e "${GREEN}MODULE TESTS${NC}"
echo ""

test_endpoint "GET" "/modules" "" "List all modules"

test_endpoint "POST" "/modules/register" '{
  "id": "test-module-1",
  "name": "Test Module",
  "description": "A test module",
  "category": "analytics",
  "version": "1.0.0",
  "pricing_model": "per_user",
  "base_cost": 10.0,
  "cost_per_user": 5.0,
  "cost_per_project": 0,
  "cost_per_company": 0,
  "max_users": 100,
  "is_core": false,
  "requires_approval": false,
  "trial_days_allowed": 14
}' "Register a new module"

test_endpoint "POST" "/modules/subscribe" '{
  "module_id": "test-module-1",
  "scope_level": "tenant",
  "scope_id": "tenant-123"
}' "Subscribe to a module"

test_endpoint "GET" "/modules/subscriptions" "" "List module subscriptions"

test_endpoint "GET" "/modules/usage" "" "Get module usage metrics"

test_endpoint "PUT" "/modules/toggle" '{
  "subscription_id": "sub-123",
  "enabled": false
}' "Toggle module subscription"

echo ""

# ==================== COMPANY TESTS ====================
echo -e "${GREEN}COMPANY TESTS${NC}"
echo ""

test_endpoint "POST" "/companies" '{
  "name": "Test Company",
  "description": "A test company",
  "status": "active",
  "industry_type": "Technology",
  "employee_count": 50,
  "website": "https://test-company.com",
  "max_projects": 10,
  "max_users": 100,
  "billing_email": "billing@test-company.com",
  "billing_address": "123 Main St, City, State 12345"
}' "Create a company"

test_endpoint "GET" "/companies" "" "List all companies"

test_endpoint "GET" "/companies/company-1" "" "Get company details"

test_endpoint "PUT" "/companies/company-1" '{
  "description": "Updated description",
  "max_projects": 20
}' "Update company"

test_endpoint "POST" "/companies/company-1/projects" '{
  "name": "Test Project",
  "description": "A test project",
  "status": "active",
  "budget_allocated": 5000.0
}' "Create a project"

test_endpoint "GET" "/companies/company-1/projects" "" "List company projects"

test_endpoint "GET" "/companies/company-1/members" "" "Get company members"

test_endpoint "POST" "/companies/company-1/members" '{
  "user_id": "user-123",
  "role": "manager",
  "department": "Engineering"
}' "Add member to company"

echo ""

# ==================== BILLING TESTS ====================
echo -e "${GREEN}BILLING TESTS${NC}"
echo ""

test_endpoint "POST" "/billing/plans" '{
  "name": "Starter Plan",
  "description": "Basic plan for small teams",
  "price": 99.0,
  "billing_cycle": "monthly",
  "included_modules": ["gamification", "analytics"]
}' "Create pricing plan"

test_endpoint "GET" "/billing/plans" "" "List pricing plans"

test_endpoint "POST" "/billing/subscribe" '{
  "plan_id": "plan-1",
  "billing_cycle": "monthly"
}' "Subscribe to plan"

test_endpoint "POST" "/billing/usage" '{
  "metric_type": "api_calls",
  "value": 1000,
  "timestamp": "'$(date -u +'%Y-%m-%dT%H:%M:%SZ')'"
}' "Record usage metrics"

test_endpoint "GET" "/billing/usage" "" "Get usage metrics"

test_endpoint "GET" "/billing/invoices" "" "List invoices"

test_endpoint "GET" "/billing/charges" "" "Calculate monthly charges"

echo ""

# ==================== SUMMARY ====================
echo -e "${GREEN}================================"
echo "Test Suite Complete"
echo "================================${NC}"
echo ""
echo "Note: This test script requires:"
echo "1. Backend running on http://localhost:8080"
echo "2. Valid authentication token in BEARER_TOKEN"
echo "3. Database seeded with test data"
echo ""
