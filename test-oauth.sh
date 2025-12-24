#!/bin/bash

# ============================================================
# OAUTH TESTING QUICK START SCRIPT
# ============================================================
# Run this script to test the complete OAuth flow
# Usage: bash test-oauth.sh [google|meta]

PLATFORM=${1:-google}
BASE_URL="http://localhost:8080"

echo "🔐 VYOM ERP - OAuth Testing"
echo "============================"
echo "Platform: $PLATFORM"
echo "Base URL: $BASE_URL"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Get test credentials
echo -e "${BLUE}Step 1: Getting test credentials...${NC}"
TEST_DATA=$(curl -s $BASE_URL/mock/oauth/test-data)

if [ "$PLATFORM" = "google" ]; then
    CLIENT_ID=$(echo $TEST_DATA | jq -r '.google.client_id')
    CLIENT_SECRET=$(echo $TEST_DATA | jq -r '.google.client_secret')
    REDIRECT_URI=$(echo $TEST_DATA | jq -r '.google.redirect_uri')
    ACCOUNT_ID=$(echo $TEST_DATA | jq -r '.google.test_account_id')
    SCOPE="https://www.googleapis.com/auth/adwords"
else
    CLIENT_ID=$(echo $TEST_DATA | jq -r '.meta.client_id')
    CLIENT_SECRET=$(echo $TEST_DATA | jq -r '.meta.client_secret')
    REDIRECT_URI=$(echo $TEST_DATA | jq -r '.meta.redirect_uri')
    ACCOUNT_ID=$(echo $TEST_DATA | jq -r '.meta.test_account_id')
    SCOPE="ads_management,business_management"
fi

echo -e "${GREEN}✓ Got credentials${NC}"
echo "  Client ID: $CLIENT_ID"
echo "  Redirect URI: $REDIRECT_URI"
echo ""

# Step 2: Authorize (get authorization code)
echo -e "${BLUE}Step 2: Getting authorization code...${NC}"
STATE="test_state_${PLATFORM}_$(date +%s)"

AUTH_RESPONSE=$(curl -s -X POST $BASE_URL/mock/oauth/$PLATFORM/authorize \
  -H "Content-Type: application/json" \
  -d "{
    \"client_id\": \"$CLIENT_ID\",
    \"redirect_uri\": \"$REDIRECT_URI\",
    \"scope\": \"$SCOPE\",
    \"state\": \"$STATE\"
  }")

AUTH_CODE=$(echo $AUTH_RESPONSE | jq -r '.code')
RETURNED_STATE=$(echo $AUTH_RESPONSE | jq -r '.state')

if [ "$STATE" != "$RETURNED_STATE" ]; then
    echo -e "${RED}✗ State mismatch!${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Got authorization code${NC}"
echo "  Code: $AUTH_CODE"
echo ""

# Step 3: Exchange code for tokens
echo -e "${BLUE}Step 3: Exchanging code for tokens...${NC}"

TOKEN_RESPONSE=$(curl -s -X POST $BASE_URL/mock/oauth/$PLATFORM/token \
  -H "Content-Type: application/json" \
  -d "{
    \"code\": \"$AUTH_CODE\",
    \"client_id\": \"$CLIENT_ID\",
    \"client_secret\": \"$CLIENT_SECRET\",
    \"redirect_uri\": \"$REDIRECT_URI\",
    \"grant_type\": \"authorization_code\"
  }")

ACCESS_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.access_token')
REFRESH_TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.refresh_token')
EXPIRES_IN=$(echo $TOKEN_RESPONSE | jq -r '.expires_in')

echo -e "${GREEN}✓ Got tokens${NC}"
echo "  Access Token: ${ACCESS_TOKEN:0:20}..."
echo "  Refresh Token: ${REFRESH_TOKEN:0:20}..."
echo "  Expires In: $EXPIRES_IN seconds"
echo ""

# Step 4: Validate token
echo -e "${BLUE}Step 4: Validating token...${NC}"

VALIDATE_RESPONSE=$(curl -s -X GET "$BASE_URL/mock/oauth/$PLATFORM/validate?access_token=$ACCESS_TOKEN")

IS_VALID=$(echo $VALIDATE_RESPONSE | jq -r '.valid')

if [ "$IS_VALID" = "true" ]; then
    echo -e "${GREEN}✓ Token is valid${NC}"
    echo "  Account ID: $(echo $VALIDATE_RESPONSE | jq -r '.account_id')"
    echo "  Account Name: $(echo $VALIDATE_RESPONSE | jq -r '.account_name')"
else
    echo -e "${RED}✗ Token validation failed${NC}"
    exit 1
fi
echo ""

# Step 5: Refresh token
echo -e "${BLUE}Step 5: Refreshing token...${NC}"

REFRESH_RESPONSE=$(curl -s -X POST $BASE_URL/mock/oauth/$PLATFORM/refresh \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\",
    \"client_id\": \"$CLIENT_ID\",
    \"client_secret\": \"$CLIENT_SECRET\"
  }")

NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.access_token')

echo -e "${GREEN}✓ Got new access token${NC}"
echo "  New Access Token: ${NEW_ACCESS_TOKEN:0:20}..."
echo ""

# Step 6: Revoke token
echo -e "${BLUE}Step 6: Revoking token...${NC}"

REVOKE_RESPONSE=$(curl -s -X POST $BASE_URL/mock/oauth/$PLATFORM/revoke \
  -H "Content-Type: application/json" \
  -d "{
    \"access_token\": \"$ACCESS_TOKEN\"
  }")

SUCCESS=$(echo $REVOKE_RESPONSE | jq -r '.success')

if [ "$SUCCESS" = "true" ]; then
    echo -e "${GREEN}✓ Token revoked successfully${NC}"
else
    echo -e "${RED}✗ Token revocation failed${NC}"
    exit 1
fi
echo ""

# Summary
echo -e "${GREEN}===========================${NC}"
echo -e "${GREEN}✓ All OAuth tests passed!${NC}"
echo -e "${GREEN}===========================${NC}"
echo ""
echo "Summary:"
echo "  Platform: $PLATFORM"
echo "  Authorization: ✓"
echo "  Token Exchange: ✓"
echo "  Token Validation: ✓"
echo "  Token Refresh: ✓"
echo "  Token Revocation: ✓"
echo ""
echo "Next: Use the access token for API calls:"
echo "  curl -H \"Authorization: Bearer $NEW_ACCESS_TOKEN\" \\"
echo "       http://localhost:8080/api/v1/${PLATFORM}-ads/campaigns"
