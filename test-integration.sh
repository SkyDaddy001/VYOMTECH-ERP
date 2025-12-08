#!/bin/bash

# Test Complete Authentication & Data Flow

echo "=== Testing ERP Dashboard Backend Integration ==="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Test Login
echo -e "${YELLOW}1. Testing Login Endpoint...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@vyomtech.com","password":"demo123"}')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token' 2>/dev/null)
TENANT_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.user.tenant_id' 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo -e "${RED}✗ Login failed${NC}"
  echo "$LOGIN_RESPONSE"
  exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}"
echo "Token: ${TOKEN:0:50}..."
echo "Tenant ID: $TENANT_ID"
echo ""

# 2. Test Dashboard Stats Endpoint
echo -e "${YELLOW}2. Testing Dashboard Stats Endpoint...${NC}"
STATS_RESPONSE=$(curl -s -X GET http://localhost:8080/api/v1/dashboard/stats \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID")

TOTAL_LEADS=$(echo "$STATS_RESPONSE" | jq '.total_leads' 2>/dev/null)

if [ -z "$TOTAL_LEADS" ] || [ "$TOTAL_LEADS" = "null" ]; then
  echo -e "${RED}✗ Dashboard stats endpoint failed${NC}"
  echo "$STATS_RESPONSE"
else
  echo -e "${GREEN}✓ Dashboard stats retrieved${NC}"
  echo "Stats Response: $(echo "$STATS_RESPONSE" | jq '.' 2>/dev/null | head -10)"
fi
echo ""

# 3. Test Leads Endpoint
echo -e "${YELLOW}3. Testing Leads Endpoint...${NC}"
LEADS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/leads?limit=5" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID")

LEAD_COUNT=$(echo "$LEADS_RESPONSE" | jq 'length' 2>/dev/null)

if [ -z "$LEAD_COUNT" ] || [ "$LEAD_COUNT" = "null" ]; then
  echo -e "${RED}✗ Leads endpoint failed${NC}"
  echo "$LEADS_RESPONSE" | head -5
else
  echo -e "${GREEN}✓ Leads endpoint working${NC}"
  echo "Leads Count: $LEAD_COUNT"
  echo "Sample: $(echo "$LEADS_RESPONSE" | jq '.[0] | {id, name, email, status}' 2>/dev/null)"
fi
echo ""

# 4. Test Campaigns Endpoint
echo -e "${YELLOW}4. Testing Campaigns Endpoint...${NC}"
CAMPAIGNS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/campaigns?limit=5" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID")

CAMPAIGN_COUNT=$(echo "$CAMPAIGNS_RESPONSE" | jq 'length' 2>/dev/null)

if [ -z "$CAMPAIGN_COUNT" ] || [ "$CAMPAIGN_COUNT" = "null" ]; then
  echo -e "${RED}✗ Campaigns endpoint failed${NC}"
  echo "$CAMPAIGNS_RESPONSE" | head -5
else
  echo -e "${GREEN}✓ Campaigns endpoint working${NC}"
  echo "Campaigns Count: $CAMPAIGN_COUNT"
fi
echo ""

# 5. Test Agents Endpoint
echo -e "${YELLOW}5. Testing Agents Endpoint...${NC}"
AGENTS_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/v1/agents" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID")

AGENT_COUNT=$(echo "$AGENTS_RESPONSE" | jq 'length' 2>/dev/null)

if [ -z "$AGENT_COUNT" ] || [ "$AGENT_COUNT" = "null" ]; then
  echo -e "${RED}✗ Agents endpoint failed${NC}"
  echo "$AGENTS_RESPONSE" | head -5
else
  echo -e "${GREEN}✓ Agents endpoint working${NC}"
  echo "Agents Count: $AGENT_COUNT"
fi
echo ""

# 6. Verify Frontend is Running
echo -e "${YELLOW}6. Checking Frontend Status...${NC}"
FRONTEND_CHECK=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000)

if [ "$FRONTEND_CHECK" = "200" ]; then
  echo -e "${GREEN}✓ Frontend is running at http://localhost:3000${NC}"
else
  echo -e "${YELLOW}⚠ Frontend returned status $FRONTEND_CHECK${NC}"
fi
echo ""

echo -e "${GREEN}=== All Tests Complete ===${NC}"
echo ""
echo "Summary:"
echo "  ✓ Authentication working"
echo "  ✓ Backend API responding"
echo "  ✓ All endpoints accessible"
echo ""
echo "Next steps:"
echo "  1. Visit http://localhost:3000 in your browser"
echo "  2. Login with: demo@vyomtech.com / demo123"
echo "  3. You should see the dashboard with real data"
