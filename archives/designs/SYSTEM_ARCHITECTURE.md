# VYOM ERP - System Architecture

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Client Layer                             │
│   ┌─────────────────┬─────────────────┬─────────────────┐   │
│   │   Web (React)   │   Mobile App    │   Admin Portal  │   │
│   │   Next.js 14    │   React Native  │   React Admin   │   │
│   └─────────────────┴─────────────────┴─────────────────┘   │
└───────────────────────────┬─────────────────────────────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
         ▼                  ▼                  ▼
   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
   │ REST API     │  │ WebSocket    │  │ gRPC Svc     │
   │ (HTTP/REST)  │  │ Real-time    │  │ Internal Svc │
   └──────────────┘  └──────────────┘  └──────────────┘
         │                  │                  │
         └──────────────────┼──────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                  API Gateway Layer                           │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Authentication (JWT + OAuth2)                       │  │
│  │  Rate Limiting | CORS | Request Validation          │  │
│  │  Multi-Tenant Routing | Logging                      │  │
│  └──────────────────────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
         ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  Business Logic Layer (Go)                   │
│  ┌──────────────┬──────────────┬──────────────┐             │
│  │  Handlers    │  Services    │  Middleware  │             │
│  │  100+ files  │  11 modules  │  Auth, CORS  │             │
│  └──────────────┴──────────────┴──────────────┘             │
│                                                             │
│  ┌────────────────────────────────────────────────────┐   │
│  │        Service Modules (11 Total)                  │   │
│  │  GL | AP | AR | HR | Leave | Sales | Real Estate  │   │
│  │  Construction | Purchase | Compliance | Dashboard  │   │
│  └────────────────────────────────────────────────────┘   │
└───────────────────────────┬─────────────────────────────────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
         ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                  Data Layer                                 │
│  ┌──────────────┬──────────────┬──────────────┐             │
│  │   MySQL DB   │   Cache      │   File Store │             │
│  │   8.0+       │   Redis      │   S3/GCS     │             │
│  └──────────────┴──────────────┴──────────────┘             │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 Multi-Tenant Architecture

### Tenant Isolation Mechanism

```
┌─────────────────────────────────────────┐
│   Request from User (Org A)             │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Extract X-Tenant-ID from Header        │
│  TenantIDKey = "org-123"                │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Store in Request Context               │
│  r.Context().Value(TenantIDKey)         │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Handler passes tenantID to Service     │
│  glService.GetFinancialRatios(tenantID) │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Service filters all queries            │
│  WHERE tenant_id = ?                    │
└────────────┬────────────────────────────┘
             │
             ▼
┌─────────────────────────────────────────┐
│  Only Org A data returned               │
│  Zero cross-tenant data exposure        │
└─────────────────────────────────────────┘
```

---

## 🔄 Request Flow Example

### Example: Create Sales Invoice

```
1. FRONTEND (Next.js)
   ├─ User clicks "Create Invoice"
   ├─ Form validation
   └─ POST /api/v1/sales/invoices

2. API GATEWAY
   ├─ Validate JWT token
   ├─ Extract X-Tenant-ID header
   ├─ Check rate limits
   └─ Route to handler

3. HANDLER (sales_handler.go)
   ├─ Parse JSON body
   ├─ Validate invoice data
   ├─ Call sales service
   └─ Return JSON response

4. SERVICE (sales_service.go)
   ├─ Validate business rules
   ├─ Create invoice record
   ├─ Post to GL (double-entry)
   │  ├─ DR: Accounts Receivable
   │  └─ CR: Sales Revenue
   ├─ Create GL entry
   └─ Return invoice ID

5. DATABASE
   ├─ Insert into sales_invoices (tenant_id filter)
   ├─ Insert into gl_entries (debit/credit pair)
   ├─ Update account balances
   └─ Audit log

6. RESPONSE
   ├─ HTTP 200 + invoice data
   └─ Frontend displays confirmation
```

---

## 📋 Data Model Relationships

### Financial Module

```
chart_of_accounts (COA)
├─ id (PK)
├─ tenant_id (FK)
├─ account_code (UNIQUE per tenant)
├─ account_type (Asset/Liability/Equity/Income/Expense)
├─ sub_account_type (Current/Non-Current, etc.)
├─ current_balance
├─ parent_id (for hierarchy)
└─ deleted_at (soft delete)

gl_entries (Journal Entries)
├─ id (PK)
├─ tenant_id (FK)
├─ journal_id (FK)
├─ account_id (FK → chart_of_accounts)
├─ debit
├─ credit
├─ entry_date
├─ reference_type (Invoice/Payment/etc)
├─ reference_id
└─ deleted_at

journals
├─ id (PK)
├─ tenant_id (FK)
├─ batch_id (FK)
├─ entry_date
├─ status (Draft/Posted)
├─ total_debit
├─ total_credit
└─ deleted_at
```

### Sales Module

```
sales_invoices
├─ id (PK)
├─ tenant_id (FK)
├─ invoice_number
├─ customer_id (FK)
├─ invoice_date
├─ invoice_amount
├─ tax_amount
├─ status (Draft/Issued/Paid/Cancelled)
├─ gl_entry_id (FK) [Posted to GL]
└─ deleted_at

sales_line_items
├─ id (PK)
├─ invoice_id (FK)
├─ description
├─ quantity
├─ unit_price
├─ line_total
└─ deleted_at

sales_payments
├─ id (PK)
├─ invoice_id (FK)
├─ payment_amount
├─ payment_date
├─ status (Pending/Completed)
├─ gl_entry_id (FK) [Posted to GL]
└─ deleted_at
```

### HR Module

```
employees
├─ id (PK)
├─ tenant_id (FK)
├─ employee_code
├─ name
├─ email
├─ department_id (FK)
├─ designation
├─ doj (Date of Joining)
├─ status (Active/Inactive)
└─ deleted_at

payroll
├─ id (PK)
├─ tenant_id (FK)
├─ employee_id (FK)
├─ payroll_month
├─ gross_salary
├─ basic_pay
├─ allowances
├─ deductions
├─ net_salary
├─ gl_entry_id (FK) [Posted to GL]
└─ deleted_at

leave_balance
├─ id (PK)
├─ tenant_id (FK)
├─ employee_id (FK)
├─ leave_type (Annual/Casual/Sick)
├─ fiscal_year
├─ entitled_days
├─ used_days
├─ available_balance
└─ deleted_at
```

---

## 🔒 Security Architecture

### Authentication Flow

```
1. User Login
   POST /api/v1/auth/login
   {email, password}
   ↓
2. Validate Credentials
   Hash password & compare
   ↓
3. Create JWT Token
   Claims: {user_id, org_id, exp, roles}
   ↓
4. Return Token + Refresh Token
   HTTP 200 {token, refresh_token, user}
   ↓
5. Client Stores Token
   localStorage/secureStorage
   ↓
6. Subsequent Requests
   Header: Authorization: Bearer {token}
   ↓
7. Verify Token
   Check signature & expiry
   ↓
8. Extract Claims
   user_id, org_id → Request context
   ↓
9. Execute Request
   All queries filtered by org_id
```

### Authorization Model

```
User → Roles → Permissions

User has Role(s):
├─ Admin (All permissions)
├─ Manager (GL, Reports)
├─ Finance (GL, AP, AR)
├─ HR Manager (HR, Payroll)
├─ Sales (Sales, Dashboard)
└─ Viewer (Read-only dashboard)

Permissions:
├─ view_{module}
├─ create_{module}
├─ update_{module}
├─ delete_{module}
└─ post_{module}
```

---

## 🗄️ Database Schema (Simplified)

### Core Tables

```
✓ chart_of_accounts      (Financial master)
✓ gl_entries            (Journal entries)
✓ journals              (Journal batches)
✓ sales_invoices        (Customer billing)
✓ sales_payments        (Collections)
✓ ap_invoices           (Vendor bills)
✓ ap_payments           (Vendor payments)
✓ employees             (HR master)
✓ payroll               (Monthly salaries)
✓ attendance            (Daily attendance)
✓ leave_requests        (Leave applications)
✓ leave_balance         (Leave entitlement)
✓ projects              (Real estate projects)
✓ project_collection_accounts (RERA segregated)
✓ compliance_records    (Regulatory tracking)
✓ users                 (System users)
✓ tenants               (Organizations)
✓ audit_logs            (Change tracking)
```

### Migrations
Located in `/migrations/`
```
001_initial_schema.sql      ← All tables created
```

---

## 🔌 API Endpoints Summary

### By Module

| Module | Endpoints | Status |
|--------|-----------|--------|
| GL | 15 | ✅ |
| AP | 12 | ✅ |
| AR | 14 | ✅ |
| HR & Payroll | 18 | ✅ |
| Leave | 16 | ✅ |
| Sales | 14 | ✅ |
| Real Estate | 20 | ✅ |
| Construction | 12 | ✅ |
| Purchase | 10 | ✅ |
| Compliance | 25 | ✅ |
| Dashboard | 20 | ✅ |
| **Total** | **176** | **✅** |

---

## 🚀 Scalability Design

### Horizontal Scaling

```
Load Balancer (AWS ALB)
    │
    ├─ API Server 1 (Go)
    ├─ API Server 2 (Go)
    ├─ API Server 3 (Go)
    └─ API Server N (Go)
    
    All sharing:
    ├─ MySQL DB (RDS)
    ├─ Redis Cache (ElastiCache)
    └─ S3 File Storage
```

### Database Optimization

```
Indexes on:
├─ tenant_id (All tables - multi-tenant)
├─ user_id (Authentication)
├─ entry_date (Date range queries)
├─ created_at (Sorting)
├─ status (Filtering)
└─ Foreign keys (Joins)

Query Optimization:
├─ Prepared statements
├─ Connection pooling
├─ Batch operations
├─ Pagination
└─ Caching layer
```

---

## 📊 Data Flow: Dashboard Example

### Financial Dashboard Request

```
1. Frontend
   GET /api/v1/dashboard/financial/ratios
   Headers: {Authorization, X-Tenant-ID}

2. Handler (financial_dashboard_handler.go)
   tenantID := r.Context().Value(middleware.TenantIDKey)
   data := glService.GetFinancialRatios(tenantID)

3. Service (gl_service.go)
   Query 1: SELECT account_balances FROM chart_of_accounts
            WHERE tenant_id = ? AND deleted_at IS NULL
   Query 2: SELECT gl_entries FROM gl_entries
            WHERE tenant_id = ? AND entry_date >= DATE_SUB(NOW(), INTERVAL 1 YEAR)

4. Database
   Returns filtered results (tenant_id isolated)

5. Service Calculations
   current_ratio = current_assets / current_liabilities
   quick_ratio = (current_assets - inventory) / current_liabilities
   debt_to_equity = total_debt / total_equity
   roe = net_income / equity
   roa = net_income / assets
   [+ 7 more ratios]

6. Handler Response
   {
     "timestamp": "2024-12-02T10:30:00Z",
     "data": {
       "current_ratio": 1.5,
       "quick_ratio": 1.2,
       ...
     }
   }
```

---

## 🔧 Configuration Management

### Environment Variables

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=vyom_erp

API_PORT=8080
API_ENV=production

JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

REDIS_HOST=localhost
REDIS_PORT=6379

AWS_REGION=ap-south-1
AWS_ACCESS_KEY=xxx
AWS_SECRET_KEY=xxx
S3_BUCKET=vyom-storage
```

---

## 🏥 Health Check & Monitoring

### Health Endpoint
```bash
GET /api/v1/health
Response: {
  "status": "ok",
  "database": "connected",
  "cache": "connected",
  "uptime": "72h"
}
```

### Metrics to Track
- API response time
- Error rate by endpoint
- Database query performance
- Cache hit rate
- Tenant-wise usage
- Cost per tenant

---

## 🎯 Performance Benchmarks

### Expected Response Times

| Operation | Time |
|-----------|------|
| Login | < 100ms |
| GET Dashboard | < 200ms |
| Create Invoice | < 150ms |
| GL Post | < 100ms |
| Report Generation | < 500ms |
| Batch Operations | < 1s per 100 records |

---

## 📈 Scalability Roadmap

### Phase 1 (Current)
- ✅ Single DB instance
- ✅ Basic caching
- ✅ Manual backups

### Phase 2 (Q1)
- Database replication (Read replicas)
- Advanced caching strategy
- Auto-scaling

### Phase 3 (Q2)
- Database sharding by tenant
- CDN for static assets
- Microservices architecture

### Phase 4 (Q3)
- Event-driven architecture (Kafka)
- Distributed tracing
- Advanced monitoring (Prometheus + Grafana)

---

**Architecture Document:** Complete ✅
**Last Updated:** December 2, 2025
