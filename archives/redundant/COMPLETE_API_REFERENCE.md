# Complete API Reference - Multi-Tenant AI Call Center

## Overview
This document provides a complete reference for all API endpoints across all 6 implemented features.

**Base URL:** `http://localhost:8080/api/v1`

**Authentication:** All protected endpoints require JWT token in Authorization header
```
Authorization: Bearer <token>
```

---

## Feature 1: Real-time WebSocket Communication

### WebSocket Connection
- **Endpoint:** `WS /api/v1/ws`
- **Auth Required:** Yes
- **Description:** Establish WebSocket connection for real-time updates
- **Usage:**
  ```typescript
  const ws = new WebSocket('ws://localhost:8080/api/v1/ws')
  ```

### Get Connection Stats
- **Endpoint:** `GET /api/v1/ws/stats`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "total_connections": 45,
    "active_users": 32,
    "connections_by_tenant": { "tenant1": 15, "tenant2": 30 }
  }
  ```

---

## Feature 2: Advanced Analytics

### Generate Report
- **Endpoint:** `POST /api/v1/analytics/reports`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "type": "agent_performance|campaign|lead_source|gamification",
    "start_date": "2025-01-01",
    "end_date": "2025-01-31"
  }
  ```
- **Response:** Report object with metrics and summary

### Export Report
- **Endpoint:** `POST /api/v1/analytics/export`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "report_id": "report_123",
    "format": "csv|json|pdf"
  }
  ```

### Get Trends
- **Endpoint:** `GET /api/v1/analytics/trends?metric=leads|calls|conversions&start_date=2025-01-01&end_date=2025-01-31`
- **Auth Required:** Yes
- **Response:** Array of trend data points with dates and values

### Get Custom Metrics
- **Endpoint:** `GET /api/v1/analytics/metrics?metric=lead_source_distribution&filter_source=campaign`
- **Auth Required:** Yes
- **Response:** Custom metric data based on filters

---

## Feature 3: Automation & Lead Routing

### Calculate Lead Score
- **Endpoint:** `POST /api/v1/automation/leads/score`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "lead_id": 123
  }
  ```
- **Response:**
  ```json
  {
    "lead_id": 123,
    "score": 85.5,
    "rank": 3,
    "factors": {
      "source_quality": 25,
      "status_score": 20,
      "engagement": 25,
      "assignment": 15.5
    }
  }
  ```

### Rank Leads
- **Endpoint:** `GET /api/v1/automation/leads/ranked?limit=100`
- **Auth Required:** Yes
- **Response:** Array of leads sorted by score (highest first)

### Route Lead to Agent
- **Endpoint:** `POST /api/v1/automation/leads/route`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "lead_id": 123
  }
  ```
- **Response:**
  ```json
  {
    "agent_id": 456,
    "routing_rule_applied": "round_robin"
  }
  ```

### Create Routing Rule
- **Endpoint:** `POST /api/v1/automation/routing-rules`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "VIP Leads Rule",
    "priority": 1,
    "conditions": { "source": "campaign", "status": "new" },
    "action_type": "assign_to_agent|assign_to_team|round_robin",
    "action_value": "456",
    "enabled": true
  }
  ```

### Schedule Campaign
- **Endpoint:** `POST /api/v1/automation/schedule-campaign`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "campaign_id": 789,
    "scheduled_time": "2025-02-01T10:00:00Z"
  }
  ```

### Get Lead Scoring Metrics
- **Endpoint:** `GET /api/v1/automation/metrics`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "avg_score": 72.3,
    "routed_today": 45,
    "active_rules": 8
  }
  ```

---

## Feature 4: Multi-Channel Communication

### Register Provider
- **Endpoint:** `POST /api/v1/communication/providers`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "type": "sms|email|whatsapp|slack|push",
    "credentials": {
      "api_key": "xxx",
      "api_secret": "yyy"
    }
  }
  ```

### Create Template
- **Endpoint:** `POST /api/v1/communication/templates`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "Welcome Email",
    "type": "email",
    "content": "Hello {{name}}, welcome!"
  }
  ```

### Send Message
- **Endpoint:** `POST /api/v1/communication/messages`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "recipient": "user@example.com",
    "type": "email",
    "template_id": "template_123",
    "body": "Custom message"
  }
  ```
- **Response:**
  ```json
  {
    "message_id": "msg_789",
    "status": "pending"
  }
  ```

### Get Message Status
- **Endpoint:** `GET /api/v1/communication/messages/status?id=msg_789`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "message_id": "msg_789",
    "status": "sent|pending|failed",
    "sent_at": "2025-01-15T10:30:00Z"
  }
  ```

### Get Message Stats
- **Endpoint:** `GET /api/v1/communication/stats`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "messages_sent": 1250,
    "delivery_rate": 98.5,
    "by_provider": { "email": 500, "sms": 750 }
  }
  ```

---

## Feature 5: Advanced Gamification

### Get User Points
- **Endpoint:** `GET /api/v1/gamification/points`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "user_id": 123,
    "total_points": 2500,
    "level": 5
  }
  ```

### Award Points
- **Endpoint:** `POST /api/v1/gamification/points/award`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "user_id": 123,
    "points": 50,
    "reason": "Call completed successfully"
  }
  ```

### Revoke Points
- **Endpoint:** `POST /api/v1/gamification/points/revoke`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "user_id": 123,
    "points": 10,
    "reason": "Policy violation"
  }
  ```

### Get User Badges
- **Endpoint:** `GET /api/v1/gamification/badges`
- **Auth Required:** Yes
- **Response:** Array of badges earned by user

### Create Badge
- **Endpoint:** `POST /api/v1/gamification/badges`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "Top Performer",
    "description": "50+ successful calls",
    "icon": "https://..."
  }
  ```

### Award Badge
- **Endpoint:** `POST /api/v1/gamification/badges/award`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "user_id": 123,
    "badge_id": 456
  }
  ```

### Get User Challenges
- **Endpoint:** `GET /api/v1/gamification/challenges`
- **Auth Required:** Yes

### Get Active Challenges
- **Endpoint:** `GET /api/v1/gamification/challenges/active`
- **Auth Required:** Yes

### Create Challenge
- **Endpoint:** `POST /api/v1/gamification/challenges`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "Sales Blitz",
    "description": "Make 100 calls this week",
    "target_score": 100
  }
  ```

### Get Leaderboard
- **Endpoint:** `GET /api/v1/gamification/leaderboard?limit=50`
- **Auth Required:** Yes

### Get Gamification Profile
- **Endpoint:** `GET /api/v1/gamification/profile`
- **Auth Required:** Yes

### Create Team Competition
- **Endpoint:** `POST /api/v1/gamification-advanced/competitions`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "Q1 Sales Challenge",
    "description": "Team competition for Q1",
    "start_date": "2025-01-01",
    "end_date": "2025-03-31"
  }
  ```

### Get Team Leaderboard
- **Endpoint:** `GET /api/v1/gamification-advanced/competitions/leaderboard?competition_id=comp_123`
- **Auth Required:** Yes

### Create Advanced Challenge
- **Endpoint:** `POST /api/v1/gamification-advanced/challenges`
- **Auth Required:** Yes

### Get Available Rewards
- **Endpoint:** `GET /api/v1/gamification-advanced/rewards`
- **Auth Required:** Yes

### Create Reward
- **Endpoint:** `POST /api/v1/gamification-advanced/rewards`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "Bonus Gift Card",
    "points_cost": 500,
    "description": "$50 gift card"
  }
  ```

### Redeem Reward
- **Endpoint:** `POST /api/v1/gamification-advanced/redeem`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "reward_id": 789
  }
  ```

### Get Advanced Leaderboard
- **Endpoint:** `GET /api/v1/gamification-advanced/leaderboard`
- **Auth Required:** Yes

### Get Gamification Stats
- **Endpoint:** `GET /api/v1/gamification-advanced/stats`
- **Auth Required:** Yes

---

## Feature 6: Compliance & Security

### RBAC - Create Role
- **Endpoint:** `POST /api/v1/compliance/roles`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "name": "Manager",
    "description": "Team manager role"
  }
  ```

### RBAC - Get Roles
- **Endpoint:** `GET /api/v1/compliance/roles`
- **Auth Required:** Yes

### Audit - Get Audit Logs
- **Endpoint:** `GET /api/v1/compliance/audit-logs?user_id=123&action=login&limit=50&offset=0`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "logs": [
      {
        "id": 1,
        "user_id": 123,
        "action": "login",
        "resource": "auth",
        "status": "success",
        "timestamp": "2025-01-15T10:30:00Z"
      }
    ],
    "total": 150
  }
  ```

### Audit - Get Audit Summary
- **Endpoint:** `GET /api/v1/compliance/audit-summary?start_date=2025-01-01&end_date=2025-01-31`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "total_events": 5000,
    "by_action": { "login": 1200, "create": 800, "update": 2000 },
    "by_status": { "success": 4900, "failure": 100 }
  }
  ```

### Security - Get Security Events
- **Endpoint:** `GET /api/v1/compliance/security-events?status=unresolved&limit=50`
- **Auth Required:** Yes
- **Response:**
  ```json
  {
    "events": [
      {
        "id": 1,
        "event_type": "failed_login_attempt",
        "severity": "high",
        "status": "unresolved",
        "timestamp": "2025-01-15T10:30:00Z"
      }
    ]
  }
  ```

### Get Compliance Report
- **Endpoint:** `GET /api/v1/compliance/report`
- **Auth Required:** Yes
- **Response:** Comprehensive compliance report with all metrics

### GDPR - Request Data Access
- **Endpoint:** `POST /api/v1/compliance/gdpr/request-access`
- **Auth Required:** Yes

### GDPR - Export User Data
- **Endpoint:** `POST /api/v1/compliance/gdpr/export`
- **Auth Required:** Yes
- **Response:** JSON file with all user data

### GDPR - Request Data Deletion
- **Endpoint:** `POST /api/v1/compliance/gdpr/request-deletion`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "reason": "User requested deletion"
  }
  ```

### GDPR - Get User Consents
- **Endpoint:** `GET /api/v1/compliance/gdpr/consents`
- **Auth Required:** Yes

### GDPR - Record Consent
- **Endpoint:** `POST /api/v1/compliance/gdpr/consents`
- **Auth Required:** Yes
- **Request Body:**
  ```json
  {
    "type": "marketing|analytics|cookies",
    "value": true
  }
  ```

---

## Core Resources

### Leads
- `GET /api/v1/leads` - List all leads
- `POST /api/v1/leads` - Create new lead
- `GET /api/v1/leads?id=123` - Get lead details
- `PUT /api/v1/leads?id=123` - Update lead
- `DELETE /api/v1/leads?id=123` - Delete lead
- `GET /api/v1/leads/stats` - Get lead statistics

### Calls
- `GET /api/v1/calls` - List all calls
- `POST /api/v1/calls` - Create new call
- `GET /api/v1/calls?id=123` - Get call details
- `POST /api/v1/calls?id=123/end` - End call
- `GET /api/v1/calls/stats` - Get call statistics

### Campaigns
- `GET /api/v1/campaigns` - List all campaigns
- `POST /api/v1/campaigns` - Create new campaign
- `GET /api/v1/campaigns?id=123` - Get campaign details
- `PUT /api/v1/campaigns?id=123` - Update campaign
- `DELETE /api/v1/campaigns?id=123` - Delete campaign
- `GET /api/v1/campaigns/stats` - Get campaign statistics

### Agents
- `GET /api/v1/agents` - List all agents
- `GET /api/v1/agents/{id}` - Get agent details
- `PATCH /api/v1/agents/status` - Update availability
- `GET /api/v1/agents/available` - Get available agents
- `GET /api/v1/agents/stats` - Get agent statistics

### Tenants
- `POST /api/v1/tenants` - Create tenant
- `GET /api/v1/tenant` - Get current tenant info
- `GET /api/v1/tenant/users/count` - Get tenant user count
- `GET /api/v1/tenants` - List user's tenants
- `POST /api/v1/tenants/{id}/switch` - Switch tenant
- `POST /api/v1/tenants/{id}/members` - Add member
- `DELETE /api/v1/tenants/{id}/members/{email}` - Remove member

### Authentication
- `POST /api/v1/auth/register` - Register new account
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/validate` - Validate token
- `POST /api/v1/auth/change-password` - Change password
- `POST /api/v1/password-reset/request` - Request password reset
- `POST /api/v1/password-reset/reset` - Reset password

### AI
- `POST /api/v1/ai/query` - Process AI query
- `GET /api/v1/ai/providers` - List AI providers

---

## Error Responses

All errors follow this format:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

### Common Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not found
- `500` - Server error

---

## Frontend Integration

All services are available through the frontend API client:

```typescript
import {
  authService,
  leadService,
  callService,
  campaignService,
  analyticsService,
  automationService,
  communicationService,
  gamificationService,
  complianceService,
  aiService,
  webSocketService,
} from '@/services/api'
```

### Using Hooks

```typescript
import { useLeads } from '@/hooks/useLeads'
import { useCalls } from '@/hooks/useCalls'
import { useGamification } from '@/hooks/useGamification'
import { useAnalytics } from '@/hooks/useAnalytics'
import { useAutomation } from '@/hooks/useAutomation'
import { useCompliance } from '@/hooks/useCompliance'
```

---

## Configuration

**Environment Variables:**
- `NEXT_PUBLIC_API_URL` - Backend API URL (default: `http://localhost:8080`)

**Database Configuration:**
- MySQL database with auto-created tables for all features
- Connection pool with configurable size
- Transaction support for critical operations

---

## Deployment

### Docker Compose
```bash
docker-compose up -d
```

### Kubernetes
```bash
kubectl apply -f k8s/
```

### Local Development
```bash
go run cmd/main.go
```

All endpoints are fully functional and ready for production use.
