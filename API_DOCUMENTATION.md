# VYOM ERP - Comprehensive API Documentation

**Last Updated:** December 3, 2025  
**API Version:** v1  
**Base URL:** `http://localhost:8080/api/v1`

---

## 📋 Table of Contents

1. [Authentication](#authentication)
2. [Multi-Tenancy](#multi-tenancy)
3. [Common Response Formats](#common-response-formats)
4. [API Endpoints](#api-endpoints)
5. [Error Handling](#error-handling)
6. [Rate Limiting](#rate-limiting)

---

## Authentication

### JWT Token Authentication

All API requests require a valid JWT token in the `Authorization` header:

```http
Authorization: Bearer <jwt_token>
```

### OAuth2 Support

The system supports OAuth2 authentication with the following flow:
- Client credentials flow for service-to-service authentication
- Authorization code flow for user authentication

### Token Refresh

```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token"
}
```

**Response:**
```json
{
  "access_token": "new_jwt_token",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

---

## Multi-Tenancy

### Tenant Identification

Every request must include the tenant ID in the request header:

```http
X-Tenant-ID: tenant_uuid
```

The tenant context is automatically applied to:
- Database queries (WHERE tenant_id = X)
- Authorization checks
- Audit logging
- Resource isolation

### Tenant Isolation

- Data is logically isolated at the database level using `tenant_id` columns
- Middleware automatically filters all queries by tenant
- Cross-tenant data access is strictly prevented

---

## Common Response Formats

### Success Response (2xx)

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "example",
    "created_at": "2025-12-03T10:30:00Z"
  },
  "message": "Operation completed successfully"
}
```

### Error Response (4xx/5xx)

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "email",
        "message": "Invalid email format"
      }
    ]
  },
  "timestamp": "2025-12-03T10:30:00Z"
}
```

### Pagination Response

```json
{
  "success": true,
  "data": [
    { "id": "1", "name": "item1" },
    { "id": "2", "name": "item2" }
  ],
  "pagination": {
    "page": 1,
    "per_page": 50,
    "total": 150,
    "total_pages": 3
  }
}
```

---

## API Endpoints

### Authentication Endpoints

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user_uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "admin"
    },
    "access_token": "jwt_token",
    "refresh_token": "refresh_token"
  }
}
```

#### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <jwt_token>
```

#### Get Current User
```http
GET /api/v1/auth/me
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### User Management Endpoints

#### Create User
```http
POST /api/v1/users
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "email": "newuser@example.com",
  "name": "New User",
  "role": "manager",
  "department": "Sales"
}
```

#### Get All Users
```http
GET /api/v1/users?page=1&per_page=50
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get User by ID
```http
GET /api/v1/users/{user_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Update User
```http
PUT /api/v1/users/{user_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "name": "Updated Name",
  "role": "senior_manager"
}
```

#### Delete User
```http
DELETE /api/v1/users/{user_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Customer Management Endpoints

#### Create Customer
```http
POST /api/v1/customers
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "name": "ABC Corporation",
  "email": "contact@abc.com",
  "phone": "+1-800-123-4567",
  "address": "123 Main St, City, State",
  "city": "San Francisco",
  "state": "CA",
  "country": "USA",
  "postal_code": "94105",
  "credit_limit": 50000,
  "kyc_verified": true,
  "pan": "ABCDE1234F",
  "gst": "18AABCU1234H1Z0"
}
```

#### Get All Customers
```http
GET /api/v1/customers?page=1&per_page=50&search=name
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Customer by ID
```http
GET /api/v1/customers/{customer_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Update Customer
```http
PUT /api/v1/customers/{customer_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "name": "Updated Corp Name",
  "credit_limit": 75000
}
```

#### Delete Customer
```http
DELETE /api/v1/customers/{customer_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Sales Endpoints

#### Create Sales Order
```http
POST /api/v1/sales
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "customer_id": "customer_uuid",
  "order_date": "2025-12-03",
  "delivery_date": "2025-12-10",
  "items": [
    {
      "product_id": "product_uuid",
      "quantity": 10,
      "unit_price": 100.00,
      "tax_rate": 18
    }
  ],
  "notes": "Delivery instructions here"
}
```

**Response:** `201 Created`

#### Get All Sales Orders
```http
GET /api/v1/sales?page=1&per_page=50&status=pending
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Sales Order by ID
```http
GET /api/v1/sales/{sales_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Update Sales Order
```http
PUT /api/v1/sales/{sales_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "status": "shipped",
  "delivery_date": "2025-12-15"
}
```

#### Delete Sales Order
```http
DELETE /api/v1/sales/{sales_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Purchase Endpoints

#### Create Purchase Order
```http
POST /api/v1/purchase
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "vendor_id": "vendor_uuid",
  "order_date": "2025-12-03",
  "delivery_date": "2025-12-15",
  "items": [
    {
      "product_id": "product_uuid",
      "quantity": 50,
      "unit_cost": 50.00,
      "tax_rate": 18
    }
  ]
}
```

#### Get All Purchase Orders
```http
GET /api/v1/purchase?page=1&per_page=50
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Purchase Order by ID
```http
GET /api/v1/purchase/{purchase_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Update Purchase Order
```http
PUT /api/v1/purchase/{purchase_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Inventory Management Endpoints

#### Get Inventory
```http
GET /api/v1/inventory?page=1&per_page=50
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Stock by Product
```http
GET /api/v1/inventory/{product_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

**Response:**
```json
{
  "success": true,
  "data": {
    "product_id": "product_uuid",
    "quantity_on_hand": 500,
    "reorder_level": 100,
    "reorder_quantity": 200,
    "unit_cost": 50.00,
    "warehouse_locations": [
      {
        "warehouse_id": "warehouse_uuid",
        "quantity": 300
      }
    ]
  }
}
```

#### Update Stock
```http
POST /api/v1/inventory/{product_id}/adjust
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "quantity_adjustment": 50,
  "reason": "physical_count_adjustment",
  "notes": "Quarterly inventory count"
}
```

---

### Accounting Endpoints

#### Get General Ledger
```http
GET /api/v1/accounting/ledger?account_code=1000&from_date=2025-12-01&to_date=2025-12-31
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

**Response:**
```json
{
  "success": true,
  "data": {
    "account": {
      "code": "1000",
      "name": "Cash",
      "type": "asset"
    },
    "entries": [
      {
        "date": "2025-12-03",
        "description": "Sales deposit",
        "debit": 1000.00,
        "credit": 0,
        "balance": 5000.00
      }
    ]
  }
}
```

#### Create Journal Entry
```http
POST /api/v1/accounting/journal
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "entry_date": "2025-12-03",
  "description": "Monthly rent payment",
  "lines": [
    {
      "account_code": "6100",
      "description": "Rent expense",
      "debit": 5000.00
    },
    {
      "account_code": "1000",
      "description": "Cash",
      "credit": 5000.00
    }
  ]
}
```

#### Get Trial Balance
```http
GET /api/v1/accounting/trial-balance?as_of_date=2025-12-31
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Post to General Ledger
```http
POST /api/v1/accounting/post-ledger
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "journal_batch_id": "batch_uuid",
  "posting_date": "2025-12-03"
}
```

---

### Banking & Reconciliation Endpoints

#### Get Bank Accounts
```http
GET /api/v1/banking/accounts
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Create Bank Reconciliation
```http
POST /api/v1/banking/reconciliation
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "bank_account_id": "account_uuid",
  "statement_date": "2025-12-03",
  "statement_balance": 50000.00,
  "items": [
    {
      "check_number": "CHK001",
      "amount": 1000.00,
      "status": "cleared"
    }
  ]
}
```

#### Get Reconciliation Status
```http
GET /api/v1/banking/reconciliation/{reconciliation_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### HR & Payroll Endpoints

#### Get Employees
```http
GET /api/v1/hr/employees?page=1&per_page=50&department=Sales
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Create Employee
```http
POST /api/v1/hr/employees
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@company.com",
  "phone": "+1-800-123-4567",
  "date_of_birth": "1990-01-15",
  "date_of_joining": "2025-01-01",
  "department": "Sales",
  "designation": "Sales Manager",
  "salary": 50000,
  "employment_type": "full_time"
}
```

#### Get Salary Slip
```http
GET /api/v1/hr/salary-slip/{employee_id}/{year}/{month}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Generate Payroll
```http
POST /api/v1/hr/payroll/generate
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "month": 12,
  "year": 2025,
  "department_id": "department_uuid"
}
```

---

### Project Management Endpoints

#### Create Project
```http
POST /api/v1/projects
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "name": "Website Redesign",
  "description": "Complete redesign of company website",
  "start_date": "2025-12-01",
  "end_date": "2026-03-31",
  "budget": 50000,
  "client_id": "client_uuid",
  "project_manager_id": "user_uuid",
  "status": "active"
}
```

#### Get All Projects
```http
GET /api/v1/projects?page=1&per_page=50&status=active
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Project by ID
```http
GET /api/v1/projects/{project_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Update Project
```http
PUT /api/v1/projects/{project_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "status": "completed",
  "end_date": "2025-12-10"
}
```

#### Create Task
```http
POST /api/v1/projects/{project_id}/tasks
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "title": "Design mockups",
  "description": "Create UI mockups",
  "assigned_to": "user_uuid",
  "due_date": "2025-12-15",
  "priority": "high",
  "status": "in_progress"
}
```

#### Get Project Tasks
```http
GET /api/v1/projects/{project_id}/tasks?status=in_progress
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Call Center & Communication Endpoints

#### Get Call Records
```http
GET /api/v1/call-center/calls?from_date=2025-12-01&to_date=2025-12-31
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Initiate Click-to-Call
```http
POST /api/v1/call-center/click-to-call
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "customer_id": "customer_uuid",
  "phone_number": "+1-800-123-4567",
  "agent_id": "agent_uuid"
}
```

#### Get Call Metrics
```http
GET /api/v1/call-center/metrics?period=daily&date=2025-12-03
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_calls": 250,
    "answered_calls": 200,
    "missed_calls": 50,
    "average_duration": "4m 30s",
    "agents": [
      {
        "agent_id": "agent_uuid",
        "name": "John Smith",
        "calls_handled": 45,
        "avg_duration": "5m 12s"
      }
    ]
  }
}
```

#### Send Message (Multi-channel)
```http
POST /api/v1/communication/send
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "recipient_id": "customer_uuid",
  "channel": "sms|email|whatsapp|slack",
  "subject": "Order Update",
  "message": "Your order has been shipped",
  "template_id": "optional_template_uuid"
}
```

---

### Reports & Analytics Endpoints

#### Get Sales Report
```http
GET /api/v1/reports/sales?from_date=2025-12-01&to_date=2025-12-31&group_by=month
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Financial Report
```http
GET /api/v1/reports/financial?report_type=income_statement&period=monthly&date=2025-12-31
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get Dashboard Metrics
```http
GET /api/v1/analytics/dashboard?period=current_month
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Workflow & Automation Endpoints

#### Create Workflow
```http
POST /api/v1/workflows
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "name": "New Customer Onboarding",
  "description": "Automated onboarding process",
  "trigger": "customer_created",
  "steps": [
    {
      "type": "send_email",
      "template": "welcome_email"
    },
    {
      "type": "create_task",
      "assigned_to": "sales_manager"
    }
  ]
}
```

#### Get Workflows
```http
GET /api/v1/workflows?page=1&per_page=50
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Execute Workflow
```http
POST /api/v1/workflows/{workflow_id}/execute
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "trigger_data": {
    "customer_id": "customer_uuid"
  }
}
```

#### Get Workflow Executions
```http
GET /api/v1/workflows/{workflow_id}/executions
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### RBAC & Permissions Endpoints

#### Get Roles
```http
GET /api/v1/roles
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Create Role
```http
POST /api/v1/roles
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "name": "Sales Manager",
  "description": "Manager for sales team",
  "permissions": [
    "sales_view",
    "sales_create",
    "sales_update",
    "reports_view"
  ]
}
```

#### Assign Role to User
```http
POST /api/v1/users/{user_id}/roles/{role_id}
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Get User Permissions
```http
GET /api/v1/users/{user_id}/permissions
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

### Settings & Configuration Endpoints

#### Get Tenant Settings
```http
GET /api/v1/settings
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

#### Update Tenant Settings
```http
PUT /api/v1/settings
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "company_name": "VYOM Corp",
  "fiscal_year_start": "01-04",
  "currency": "INR",
  "timezone": "IST",
  "gst_applicable": true
}
```

#### Get Feature Flags
```http
GET /api/v1/features
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
```

---

## Error Handling

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or missing authentication token |
| `FORBIDDEN` | 403 | Insufficient permissions for the resource |
| `NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request parameters |
| `CONFLICT` | 409 | Resource already exists |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

### Error Response Example

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "email",
        "message": "Email is required"
      },
      {
        "field": "age",
        "message": "Age must be at least 18"
      }
    ]
  },
  "timestamp": "2025-12-03T10:30:00Z"
}
```

---

## Rate Limiting

### Rate Limits

- **Standard requests:** 1000 requests per hour per user
- **Authentication requests:** 10 requests per minute per IP
- **Bulk operations:** 100 requests per minute per user
- **Reports/Analytics:** 50 requests per hour per user

### Rate Limit Headers

All responses include rate limit information:

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1638534600
```

### Retry-After Header

When rate limited, the response includes:

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 60

{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please retry after 60 seconds"
  }
}
```

---

## Webhook Events

The API supports webhooks for real-time event notifications:

### Supported Events

- `customer.created`
- `customer.updated`
- `customer.deleted`
- `sales.created`
- `sales.updated`
- `sales.shipped`
- `purchase.created`
- `purchase.received`
- `payment.received`
- `invoice.generated`
- `project.completed`
- `task.completed`

### Webhook Registration

```http
POST /api/v1/webhooks
Authorization: Bearer <jwt_token>
X-Tenant-ID: tenant_uuid
Content-Type: application/json

{
  "url": "https://your-app.com/webhooks/events",
  "events": ["sales.created", "payment.received"],
  "active": true
}
```

### Webhook Payload Example

```json
{
  "id": "webhook_event_uuid",
  "event": "sales.created",
  "timestamp": "2025-12-03T10:30:00Z",
  "data": {
    "sales_id": "sales_uuid",
    "customer_id": "customer_uuid",
    "total_amount": 10000.00,
    "items": [...]
  }
}
```

---

## Best Practices

### 1. Always Include Tenant ID
```http
X-Tenant-ID: tenant_uuid
```

### 2. Use Pagination for Large Datasets
```http
GET /api/v1/customers?page=1&per_page=50
```

### 3. Implement Retry Logic
- Implement exponential backoff for 5xx errors
- Retry maximum 3 times
- Don't retry 4xx errors

### 4. Cache Responses When Possible
- Use ETags for conditional requests
- Cache list endpoints for 5 minutes
- Invalidate cache on updates

### 5. Handle Time Zones
- Always use ISO 8601 format for dates
- Include timezone in API responses
- Consider user's timezone preference

### 6. Secure Token Storage
- Never store tokens in localStorage (use httpOnly cookies)
- Implement token rotation
- Clear tokens on logout

### 7. Implement Proper Logging
- Log all API requests/responses
- Never log sensitive data (passwords, PAN, etc.)
- Use structured logging (JSON format)

---

## SDK & Client Libraries

### JavaScript/TypeScript

```typescript
import { VyomERP } from '@vyomtech/erp-sdk';

const client = new VyomERP({
  baseUrl: 'http://localhost:8080/api/v1',
  accessToken: 'your_jwt_token',
  tenantId: 'your_tenant_id'
});

// Create customer
const customer = await client.customers.create({
  name: 'ABC Corp',
  email: 'contact@abc.com'
});

// Get sales orders
const orders = await client.sales.list({
  page: 1,
  per_page: 50,
  status: 'pending'
});
```

### Python

```python
from vyomtech import VyomERP

client = VyomERP(
    base_url='http://localhost:8080/api/v1',
    access_token='your_jwt_token',
    tenant_id='your_tenant_id'
)

# Create customer
customer = client.customers.create(
    name='ABC Corp',
    email='contact@abc.com'
)

# Get sales orders
orders = client.sales.list(page=1, per_page=50, status='pending')
```

---

## Support & Resources

- **Documentation:** https://docs.vyomtech.com
- **API Status:** https://status.vyomtech.com
- **Support Email:** support@vyomtech.com
- **Community Forum:** https://community.vyomtech.com
- **GitHub Issues:** https://github.com/vyomtech/erp-sdk/issues

---

## Version History

| Version | Released | Changes |
|---------|----------|---------|
| v1.5.0 | 2025-12-03 | Multi-channel communication, WebRTC support |
| v1.4.0 | 2025-11-15 | Fixed assets & depreciation, cost centers |
| v1.3.0 | 2025-10-20 | Bank reconciliation, payment enhancements |
| v1.2.0 | 2025-09-10 | Click-to-call system, AI call center |
| v1.1.0 | 2025-08-01 | RBAC, compliance, tax features |
| v1.0.0 | 2025-07-01 | Initial release |

---

**Last Updated:** December 3, 2025  
**Maintained By:** VYOM Tech Development Team
