# Multi-Tenant API Endpoint Testing Guide

This guide provides comprehensive instructions for testing all multi-tenant API endpoints using curl, Postman, or similar tools.

## Setup

### Prerequisites
- Backend running on http://localhost:8080
- PostgreSQL database with tenant tables
- Valid JWT token from login

### Get Authentication Token

```bash
# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password"
  }'

# Response will include:
# {
#   "token": "eyJhbGc...",
#   "user": {...},
#   "tenant": {...}
# }

# Save token for subsequent requests
TOKEN="your-jwt-token-here"
```

## Endpoint Testing

### 1. Get Current Tenant Info

**Endpoint:** `GET /api/v1/tenant`

**Purpose:** Get information about the user's current active tenant

```bash
curl -X GET http://localhost:8080/api/v1/tenant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"
```

**Expected Response (200 OK):**
```json
{
  "id": 1,
  "name": "ACME Corp",
  "domain": "acme.example.com",
  "description": "Main tenant",
  "max_users": 50,
  "status": "active",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Error Cases:**
- 401 Unauthorized: Invalid or missing token
- 404 Not Found: Tenant doesn't exist
- 403 Forbidden: User doesn't belong to tenant

---

### 2. Get Tenant User Count

**Endpoint:** `GET /api/v1/tenant/users/count`

**Purpose:** Get the count of users in the current tenant and max capacity

```bash
curl -X GET http://localhost:8080/api/v1/tenant/users/count \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**
```json
{
  "count": 15,
  "max_users": 50,
  "percentage": 30
}
```

**Troubleshooting:**
- If count is 0, verify users are added to tenant in database
- Check tenant_users table has entries

---

### 3. List All Tenants (Admin)

**Endpoint:** `GET /api/v1/tenants`

**Purpose:** List all tenants in the system (admin only)

```bash
# Without authentication (if public)
curl -X GET http://localhost:8080/api/v1/tenants

# With authentication
curl -X GET http://localhost:8080/api/v1/tenants \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "ACME Corp",
    "domain": "acme.example.com",
    "max_users": 50,
    "status": "active"
  },
  {
    "id": 2,
    "name": "TechStart Inc",
    "domain": "techstart.example.com",
    "max_users": 25,
    "status": "active"
  }
]
```

---

### 4. Get User's Tenants

**Endpoint:** `GET /api/v1/tenants`

**Purpose:** Get list of tenants the authenticated user belongs to

```bash
curl -X GET http://localhost:8080/api/v1/tenants \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "ACME Corp",
    "domain": "acme.example.com",
    "role": "admin",
    "joined_at": "2024-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "name": "TechStart Inc",
    "domain": "techstart.example.com",
    "role": "member",
    "joined_at": "2024-02-15T00:00:00Z"
  }
]
```

**Notes:**
- Response includes user's role in each tenant
- Shows join date for each tenant
- Only returns tenants user is member of

---

### 5. Switch Tenant

**Endpoint:** `POST /api/v1/tenants/{id}/switch`

**Purpose:** Switch the user's active tenant context

```bash
# Switch to tenant with ID 2
curl -X POST http://localhost:8080/api/v1/tenants/2/switch \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "message": "Tenant switched successfully",
  "tenant_id": 2,
  "tenant": {
    "id": 2,
    "name": "TechStart Inc",
    "domain": "techstart.example.com"
  }
}
```

**Error Cases:**
```bash
# Invalid tenant ID
# Response (404 Not Found):
# {
#   "error": "Tenant not found",
#   "code": "TENANT_NOT_FOUND"
# }

# User doesn't belong to tenant
# Response (403 Forbidden):
# {
#   "error": "Access denied to tenant",
#   "code": "ACCESS_DENIED"
# }
```

**Verification:**
After switching, verify by checking `/api/v1/tenant` again - should return the new tenant

---

### 6. Add Tenant Member

**Endpoint:** `POST /api/v1/tenants/{id}/members`

**Purpose:** Add a new member to a tenant

```bash
# Add user to tenant 1
curl -X POST http://localhost:8080/api/v1/tenants/1/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "role": "member"
  }'
```

**Request Parameters:**
- `email`: Email address of user to add (string, required)
- `role`: Role in tenant - "admin", "member", "viewer" (string, required)

**Expected Response (201 Created):**
```json
{
  "success": true,
  "message": "User added to tenant",
  "user_email": "newuser@example.com",
  "role": "member"
}
```

**Error Cases:**
```bash
# User already in tenant
# Response (400 Bad Request):
# {
#   "error": "User already exists in tenant",
#   "code": "USER_EXISTS"
# }

# User doesn't exist
# Response (404 Not Found):
# {
#   "error": "User not found",
#   "code": "USER_NOT_FOUND"
# }

# Insufficient permissions
# Response (403 Forbidden):
# {
#   "error": "Only tenant admins can add members",
#   "code": "INSUFFICIENT_PERMISSION"
# }
```

**Role Reference:**
- `admin`: Full access to tenant settings
- `member`: Standard access to tenant resources
- `viewer`: Read-only access

---

### 7. Remove Tenant Member

**Endpoint:** `DELETE /api/v1/tenants/{id}/members/{email}`

**Purpose:** Remove a member from a tenant

```bash
# Remove user from tenant 1
curl -X DELETE http://localhost:8080/api/v1/tenants/1/members/user@example.com \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**
```json
{
  "success": true,
  "message": "User removed from tenant",
  "user_email": "user@example.com"
}
```

**Error Cases:**
```bash
# User not in tenant
# Response (404 Not Found):
# {
#   "error": "User not found in tenant",
#   "code": "USER_NOT_IN_TENANT"
# }

# Last admin in tenant
# Response (400 Bad Request):
# {
#   "error": "Cannot remove last admin from tenant",
#   "code": "CANNOT_REMOVE_LAST_ADMIN"
# }

# Insufficient permissions
# Response (403 Forbidden):
# {
#   "error": "Only tenant admins can remove members",
#   "code": "INSUFFICIENT_PERMISSION"
# }
```

---

## Complete Workflow Test

This test covers a complete workflow of tenant operations:

```bash
#!/bin/bash

# Set up variables
API_URL="http://localhost:8080"
USER1_EMAIL="user1@example.com"
USER1_PASS="password123"
USER2_EMAIL="user2@example.com"

echo "=== Multi-Tenant API Workflow Test ==="

# 1. Login as User1
echo -e "\n1. Logging in User1..."
LOGIN_RESPONSE=$(curl -s -X POST $API_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER1_EMAIL\",
    \"password\": \"$USER1_PASS\"
  }")

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "Token: $TOKEN"

# 2. Get current tenant
echo -e "\n2. Getting current tenant..."
curl -s -X GET $API_URL/api/v1/tenant \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 3. Get user's tenants
echo -e "\n3. Getting user's tenants..."
curl -s -X GET $API_URL/api/v1/tenants \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 4. Get tenant user count
echo -e "\n4. Getting tenant user count..."
curl -s -X GET $API_URL/api/v1/tenant/users/count \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# 5. Add member (assume user has admin role)
echo -e "\n5. Adding new member to tenant..."
curl -s -X POST $API_URL/api/v1/tenants/1/members \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$USER2_EMAIL\",
    \"role\": \"member\"
  }" | jq '.'

# 6. List members again
echo -e "\n6. Checking tenant user count after adding member..."
curl -s -X GET $API_URL/api/v1/tenant/users/count \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo -e "\n=== Test Complete ==="
```

---

## Postman Collection

Here's a Postman collection you can import:

```json
{
  "info": {
    "name": "Multi-Tenant API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Get Current Tenant",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenant",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenant"]
        }
      }
    },
    {
      "name": "Get Tenant User Count",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenant/users/count",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenant", "users", "count"]
        }
      }
    },
    {
      "name": "List All Tenants",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenants",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenants"]
        }
      }
    },
    {
      "name": "Get User's Tenants",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenants",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenants"]
        }
      }
    },
    {
      "name": "Switch Tenant",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          },
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{}"
        },
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenants/{{tenantId}}/switch",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenants", "{{tenantId}}", "switch"]
        }
      }
    },
    {
      "name": "Add Tenant Member",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          },
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\"email\": \"user@example.com\", \"role\": \"member\"}"
        },
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenants/{{tenantId}}/members",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenants", "{{tenantId}}", "members"]
        }
      }
    },
    {
      "name": "Remove Tenant Member",
      "request": {
        "method": "DELETE",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/api/v1/tenants/{{tenantId}}/members/{{email}}",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "tenants", "{{tenantId}}", "members", "{{email}}"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080"
    },
    {
      "key": "token",
      "value": ""
    },
    {
      "key": "tenantId",
      "value": "1"
    },
    {
      "key": "email",
      "value": "user@example.com"
    }
  ]
}
```

---

## Performance Testing

### Load Testing with Apache Bench

```bash
# Test tenant listing performance
ab -n 100 -c 10 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenant

# Results should show:
# - Response time < 100ms
# - No failed requests
```

### Load Testing with wrk

```bash
# Install wrk (macOS)
brew install wrk

# Test tenant switching
wrk -t 4 -c 100 -d 30s \
  -s tenant_test.lua \
  http://localhost:8080/api/v1/tenants/1/switch
```

---

## Monitoring & Debugging

### Enable Debug Logging

Set environment variable before running:
```bash
export LOG_LEVEL=debug
```

### Check API Response Times

```bash
# Using curl with timing
curl -w "\nResponse time: %{time_total}s\n" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/tenant
```

### Database Queries

Monitor active queries:
```sql
-- PostgreSQL
SELECT * FROM pg_stat_activity WHERE query LIKE '%tenant%';
```

---

## Troubleshooting

### Common Issues

**401 Unauthorized**
```bash
# Token is invalid or expired
# Solution: Get new token via login
```

**404 Not Found**
```bash
# Tenant or user doesn't exist
# Check: 
#   - Tenant ID exists in database
#   - User is member of tenant
```

**403 Forbidden**
```bash
# User doesn't have permission
# Check:
#   - User's role in tenant
#   - User belongs to tenant
```

**500 Internal Server Error**
```bash
# Backend error
# Check:
#   - Server logs
#   - Database connection
#   - Database migrations applied
```

---

## Best Practices

1. **Always validate responses** - Check status codes and error messages
2. **Use environment variables** - Set token and URLs in Postman/scripts
3. **Test error cases** - Don't just test happy paths
4. **Monitor performance** - Track response times
5. **Log everything** - Keep records of test results
6. **Test with realistic data** - Use actual user counts and tenant sizes
7. **Automate tests** - Use scripts for repeated testing

