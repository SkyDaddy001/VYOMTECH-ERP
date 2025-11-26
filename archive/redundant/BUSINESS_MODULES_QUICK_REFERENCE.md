# Business Modules - Quick Technical Reference
## Phase 3E Development Handbook

**Date**: November 24, 2025  
**Audience**: Developers, Architects, DevOps  
**Version**: 1.0

---

## 🗂️ Quick Module Summary

### Module Overview Table

| Module | Priority | DB Tables | Endpoints | Frontend Screens | Integrates With |
|--------|----------|-----------|-----------|-----------------|-----------------|
| **HR & Payroll** | HIGH | 22 | 45 | 12 | Accounts, Gamification |
| **Accounts (GL)** | HIGH | 20 | 40 | 15 | All modules |
| **Sales** | HIGH | 18 | 35 | 10 | Accounts, Call Center, Gamification |
| **Purchase** | HIGH | 16 | 30 | 8 | Accounts, Inventory |
| **Post Sales** | MEDIUM | 12 | 25 | 8 | Sales, Accounts |
| **Construction** | MEDIUM | 20 | 40 | 12 | Accounts, Purchase, HR |
| **Civil** | MEDIUM | 12 | 20 | 6 | Construction, Accounts |
| **TOTAL** | - | 130 | 235+ | 71 | 15+ integration points |

---

## 🗄️ Database Schema Quick Reference

### Core Tables by Module

#### HR Module (22 tables)
```sql
-- Core HR
employees, designations, departments, salary_structures
-- Attendance
attendance, shift_definitions, shift_assignments
-- Leave
leave_types, leave_applications, leaves_taken
-- Payroll
salary_slips, allowances, deductions, payroll_runs
-- Compliance
statutory_details, bank_details
-- Support
employee_documents, employee_skills, performance_reviews, employee_audit_log
```

#### Accounts Module (20 tables)
```sql
-- GL Management
gl_masters, gl_hierarchy, cost_centers, account_settings
-- Journal Posting
journal_entries, journal_line_items
-- Banking
bank_accounts, bank_reconciliation, cheques
-- Billing
invoices, invoice_items, payments, payment_terms
-- Reporting
financial_reports, report_templates, balance_sheet_data, profit_loss_data
-- Compliance
tax_rates, account_audit_log
-- Banking Integration
bank_feeds
```

#### Sales Module (18 tables)
```sql
-- CRM
customers, customer_contacts, customer_addresses, customer_preferences
-- Opportunities
opportunities, opportunity_stages, opportunity_line_items
-- Quotes
quotations, quote_line_items
-- Orders
sales_orders, order_line_items, order_fulfillment
-- Commission
sales_representatives, commission_structures, commission_calculations
-- History
customer_interaction_history, sales_targets, sales_audit_log
```

#### Purchase Module (16 tables)
```sql
-- Vendor Management
vendors, vendor_contacts, vendor_addresses
-- Requisition & PO
purchase_requisitions, purchase_orders, po_line_items
-- Receipt
goods_receipts, receipt_line_items
-- Quality
quality_inspections
-- Invoice & Payment
vendor_invoices, invoice_line_items
-- Payment
payments
-- Audit
vendor_performance_metrics, purchase_approvals, purchase_audit_log
```

#### Construction Module (20 tables)
```sql
-- Project Planning
construction_projects, project_phases, work_breakdown_structure
-- Tasks
tasks, task_dependencies, task_resources
-- BOQ
boq_master, boq_line_items
-- Rates
material_rates, labour_rates, equipment_rates
-- Tracking
daily_reports, material_usage, labour_usage, equipment_usage
-- Quality
quality_inspections, defect_register, non_conformance_reports
-- Support
project_milestones, construction_audit_log
```

#### Post Sales Module (12 tables)
```sql
-- Service Tickets
service_tickets, ticket_categories, ticket_assignments, ticket_comments
-- Warranty
warranty_records, warranty_claims, warranty_claim_documents
-- Support
support_tickets, support_communications
-- Knowledge
knowledge_base_articles, faq_items, satisfaction_surveys
```

#### Civil Module (12 tables)
```sql
-- Site Management
civil_sites, site_contacts, site_amenities
-- Contractors
contractors, contractor_rates, contractor_agreements
-- Compliance
safety_incidents, safety_checklist_items, compliance_audits
-- Permits & Environment
permits, waste_management_records, environmental_data
```

---

## 🔌 Integration Points Reference

### GL Integration Pattern
Every transaction module posts to GL using this pattern:

```go
// Transaction occurs in source module
transaction := CreateTransaction(...)

// Generate GL entries
entries := transaction.ToGLEntries()  
// Returns: Debit Account, Credit Account, Amount, Narration

// Post to GL
gl.PostEntries(entries)

// Example: HR Payroll → Accounts
// Entry 1: DR Salary Expense | CR Bank Account
// Entry 2: DR Tax Expense | CR Tax Payable
```

### API Response Structure (Consistent Across All Modules)
```json
{
  "success": true,
  "status": 200,
  "message": "Operation successful",
  "data": {
    "id": "ulid_here",
    "created_at": "2025-11-24T10:00:00Z",
    ...
  },
  "meta": {
    "tenant_id": "tenant_123",
    "company_id": "company_456",
    "user_id": "user_789"
  },
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

### Error Response Structure
```json
{
  "success": false,
  "status": 400,
  "message": "Validation failed",
  "errors": [
    {
      "field": "amount",
      "error": "must be positive"
    }
  ],
  "request_id": "req_123456"
}
```

---

## 🔄 Data Flow Examples

### Flow 1: Sales Order → Invoice → GL Entry
```
Sales Module
├─ Create Sales Order (SO)
│  ├─ Customer ID ✓
│  ├─ Order Items ✓
│  └─ Total Amount ✓
│
├─ Mark Order as "Fulfilled"
│  └─ Trigger Invoice Creation
│
└─ Auto-Generate Invoice
   ├─ Invoice ID: INV-2025-001
   ├─ Amount: $5,000
   └─ Link to SO

     ↓

Accounts Module (GL)
├─ Receive Invoice Event
├─ Generate GL Entry:
│  ├─ DR: Accounts Receivable ($5,000)
│  └─ CR: Revenue ($5,000)
├─ Post Entry
└─ Update GL Balance
```

### Flow 2: Purchase Order → GRN → Invoice → Payment
```
Purchase Module
├─ Create PO
│  └─ Link to Vendor
│
├─ Receive Goods
│  └─ Create GRN
│     ├─ QC Inspection
│     └─ Accept/Reject
│
├─ Receive Vendor Invoice
│  └─ Match: PO vs GRN vs Invoice
│
└─ Approve for Payment

     ↓

Accounts Module
├─ Record GL Entry:
│  ├─ DR: Inventory
│  └─ CR: Accounts Payable
│
├─ Approve Payment
│  └─ Generate Check/Transfer
│
└─ Update GL:
    ├─ DR: Accounts Payable
    └─ CR: Bank
```

### Flow 3: HR Payroll → GL → Bank Transfer
```
HR Module
├─ Run Monthly Payroll
│  ├─ Calculate Salaries
│  ├─ Apply Deductions
│  └─ Generate Salary Slips
│
└─ Approve Payroll Run

     ↓

Accounts Module
├─ Receive Payroll Event
├─ Generate GL Entries:
│  ├─ DR: Salary Expense
│  ├─ DR: Tax Expense
│  ├─ CR: Bank (net salary)
│  ├─ CR: Tax Payable
│  └─ CR: PF Payable
│
└─ Post Entries

     ↓

Bank Transfer
├─ Funds transferred to employees
└─ Bank reconciliation updates GL
```

---

## 🛠️ Development Tools & Setup

### Required Development Stack
```
Backend:
├─ Go 1.25.4
├─ Go Modules (dependency management)
├─ Gorilla Mux (routing - existing)
├─ GORM (ORM for database operations)
├─ sqlc (SQL code generation)
└─ Testing: testify, mockery

Frontend:
├─ Next.js 16.0.3
├─ React 19.2.0
├─ TypeScript 5.3
├─ Tailwind CSS (styling - existing)
├─ Zustand (state management - existing)
├─ React Query (data fetching)
├─ React Hook Form (form handling)
└─ Testing: Jest, React Testing Library

Database:
├─ MySQL 8.0.44
├─ Redis 7.0+ (caching)
├─ Flyway (migrations - consider switching to Go migrate)
└─ MySQL Workbench (ER diagrams)

DevOps:
├─ Docker & Docker Compose
├─ GitHub Actions (CI/CD)
├─ kubectl (k8s)
└─ Prometheus + Grafana (monitoring)
```

### IDE Extensions (VS Code)
```
Backend Development:
├─ Go (golang.go)
├─ Database Client (cweijan.vscode-database-client2)
├─ REST Client (humao.rest-client)
├─ Thunder Client (rangav.vscode-thunder-client)
└─ Go Doc (ms-vscode.Go)

Frontend Development:
├─ ES7+ React/Redux/React-Native (dsznajder.es7-react-js-snippets)
├─ Tailwind CSS IntelliSense (bradlc.vscode-tailwindcss)
├─ TypeScript Vue Plugin (vue.vscode-typescript-vue-plugin)
└─ Prettier (esbenp.prettier-vscode)

General:
├─ GitLens (eamodio.gitlens)
├─ Docker (ms-azuretools.vscode-docker)
├─ SQLTools (mtxr.sqltools)
└─ Thunder Client (rangav.vscode-thunder-client)
```

### Local Environment Setup Checklist
```
Database:
[ ] MySQL running on localhost:3306
[ ] Database created: erp_db
[ ] Test user with privileges
[ ] Sample test data loaded

Backend:
[ ] Go 1.25.4 installed
[ ] Go modules initialized
[ ] Dependencies downloaded (go mod download)
[ ] Backend running on localhost:8080

Frontend:
[ ] Node.js 18+ installed
[ ] npm/yarn dependencies installed
[ ] Frontend running on localhost:3000

Integration:
[ ] Backend → Database connection verified
[ ] Frontend → Backend API connectivity verified
[ ] Logs accessible from all components
```

---

## 📝 API Endpoint Structure

### Naming Convention
```
/api/v2/{module}/{resource}
/api/v2/{module}/{resource}/{id}
/api/v2/{module}/{resource}/{id}/{action}

Examples:
GET    /api/v2/hr/employees                 # List employees
GET    /api/v2/hr/employees/{id}            # Get employee
POST   /api/v2/hr/employees                 # Create employee
PUT    /api/v2/hr/employees/{id}            # Update employee
DELETE /api/v2/hr/employees/{id}            # Delete employee

GET    /api/v2/hr/attendance               # List attendance
POST   /api/v2/hr/attendance/{id}/approve  # Approve attendance

GET    /api/v2/accounts/gl-masters         # List GL accounts
POST   /api/v2/accounts/journal-entries    # Create journal entry

POST   /api/v2/sales/orders/{id}/fulfill   # Fulfill order
```

### HTTP Status Codes
```
200 OK                  - Successful GET/PUT
201 Created             - Successful POST
204 No Content          - Successful DELETE
400 Bad Request         - Validation error
401 Unauthorized        - No/invalid auth
403 Forbidden           - Lacks permissions
404 Not Found           - Resource doesn't exist
409 Conflict            - State violation (e.g., amount negative)
422 Unprocessable       - Business logic violation
429 Too Many Requests   - Rate limited
500 Server Error        - Unexpected error
503 Service Unavailable - Maintenance/outage
```

### Required Request Headers
```
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
X-Company-ID: {company_id}  # Optional, defaults to user's company
X-Request-ID: {uuid}         # For tracing
Content-Type: application/json
```

### Query Parameters (Standard Across All Modules)
```
GET /api/v2/{module}/{resource}
├─ limit=20              (default: 20, max: 100)
├─ offset=0              (pagination)
├─ sort_by=created_at    (which field to sort)
├─ order=asc             (asc|desc)
├─ filter[status]=active (dynamic filters)
├─ include=relations     (include related data)
└─ fields=id,name        (return only specified fields)
```

---

## 🔐 Security Checklist per Module

### Authentication & Authorization
```
[ ] All endpoints require JWT token
[ ] Token validation on each request
[ ] Tenant isolation enforced
[ ] Company-level isolation enforced
[ ] RBAC permission check before each operation
```

### Data Protection
```
[ ] All user passwords hashed (bcrypt)
[ ] Sensitive data encrypted (salary, bank details)
[ ] SQL injection protection (parameterized queries)
[ ] XSS protection (input sanitization)
[ ] CSRF tokens for POST/PUT/DELETE
```

### Audit & Compliance
```
[ ] All transactional operations logged
[ ] Audit table for each entity
[ ] Change tracking (old vs new value)
[ ] User identification on each change
[ ] Timestamp on each operation
```

### API Security
```
[ ] Rate limiting (100 req/min per user)
[ ] Request size limits
[ ] API key rotation support
[ ] HTTPS/TLS everywhere
[ ] CORS properly configured
```

---

## 📊 Reporting Endpoints (Finance/BI)

### HR Module Reports
```
GET /api/v2/hr/reports/salary-summary       # Monthly salary summary
GET /api/v2/hr/reports/attendance-analysis  # Attendance trends
GET /api/v2/hr/reports/payroll-reconcile    # Payroll vs GL reconciliation
GET /api/v2/hr/reports/compliance           # Tax, PF, ESI compliance
GET /api/v2/hr/reports/employee-roster      # Current employee roster
```

### Accounts Module Reports
```
GET /api/v2/accounts/reports/trial-balance           # Trial balance
GET /api/v2/accounts/reports/profit-loss/{period}    # P&L statement
GET /api/v2/accounts/reports/balance-sheet/{period}  # Balance sheet
GET /api/v2/accounts/reports/cash-flow/{period}      # Cash flow
GET /api/v2/accounts/reports/aging                   # AR/AP aging
```

### Sales Module Reports
```
GET /api/v2/sales/reports/pipeline          # Sales pipeline summary
GET /api/v2/sales/reports/forecast          # Revenue forecast
GET /api/v2/sales/reports/performance       # Sales rep performance
GET /api/v2/sales/reports/customer-analysis # Customer analytics
```

### Purchase Module Reports
```
GET /api/v2/purchase/reports/outstanding-po    # Outstanding POs
GET /api/v2/purchase/reports/vendor-performance # Vendor KPIs
GET /api/v2/purchase/reports/spend-analysis     # Spend analytics
```

### Construction Module Reports
```
GET /api/v2/construction/reports/progress       # Project progress
GET /api/v2/construction/reports/budget-variance # Budget vs actual
GET /api/v2/construction/reports/resource-usage # Resource utilization
GET /api/v2/construction/reports/quality        # Quality metrics
```

---

## 🚀 Performance Optimization Techniques

### Database Performance
```
Indexing Strategy:
├─ Foreign keys: Always index
├─ Tenant/Company filter columns: Always index
├─ Status/Type fields: Always index
├─ Date range queries: Composite indexes
└─ Search fields: Full-text indexes

Query Optimization:
├─ Use SELECT specific_columns (not *)
├─ Filter early with WHERE
├─ Join on indexed columns
├─ Use LIMIT for pagination
├─ Avoid N+1 queries (use JOINs)

Example - Bad Query:
FOR EACH employee:
  SELECT salary FROM salaries WHERE employee_id = ?
(N+1 problem)

Example - Good Query:
SELECT e.id, s.salary 
FROM employees e 
JOIN salaries s ON e.id = s.employee_id
WHERE e.tenant_id = ?
```

### Caching Strategy
```
Module Cache Layers:
├─ Level 1: In-memory (Go cache, 1MB)
│  ├─ User roles/permissions (5 min TTL)
│  ├─ GL account list (30 min TTL)
│  └─ Employee designation list (1 hour TTL)
│
├─ Level 2: Redis (cluster, 2GB)
│  ├─ User sessions (24 hours)
│  ├─ Frequently accessed reports (1 hour)
│  └─ GL account balances (30 min)
│
└─ Level 3: Browser cache (Frontend)
    ├─ Static assets (1 week)
    ├─ API responses (5 min)
    └─ Immutable data (master data)
```

### API Response Optimization
```
Pagination Strategy:
├─ Default limit: 20 records
├─ Max limit: 100 records
├─ Cursor-based for large datasets
└─ Example: ?limit=50&offset=0

Field Selection:
├─ Allow clients to specify fields
├─ Example: ?fields=id,name,email
└─ Reduces payload by 40-60%

Compression:
├─ gzip all responses
├─ brotli for modern clients
└─ Reduces bandwidth by 70%+

Response Caching:
├─ Cache-Control: public, max-age=300
├─ ETag for conditional requests
└─ 304 Not Modified for cached responses
```

---

## 🧪 Testing Strategy

### Unit Testing (Per Endpoint)
```go
Example: Testing HR Employee Creation

func TestCreateEmployee(t *testing.T) {
    // Arrange
    mockDB := NewMockDB()
    service := NewEmployeeService(mockDB)
    
    // Act
    employee, err := service.CreateEmployee(&Employee{
        Name: "John Doe",
        Email: "john@example.com",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, employee)
    assert.Equal(t, "john@example.com", employee.Email)
}

Target: 85%+ coverage per module
```

### Integration Testing (Module-to-Module)
```
Test: Sales Order → GL Entry Creation

Steps:
1. Create sales order via Sales API
2. Mark order as fulfilled
3. Verify GL entries created in Accounts module
4. Check account balances updated
5. Verify audit trail recorded

Expected: 
├─ AR account increased by order amount
├─ Revenue account increased
└─ Audit log contains both entries
```

### End-to-End Testing
```
Test Scenario: Monthly Payroll to GL Reconciliation

1. HR Module: Create salary structure
2. HR Module: Mark attendance complete
3. HR Module: Run monthly payroll
4. Accounts Module: Verify GL entries posted
5. Accounts Module: Run trial balance
6. Verify: HR payroll totals = GL salary expense
```

### Performance Testing
```
Load Testing Targets:
├─ 500 concurrent users
├─ API response time < 200ms (p95)
├─ Database queries < 100ms (p95)
├─ No errors under load
└─ Memory stable after 1 hour

Tools: k6, locust, JMeter
```

---

## 🔄 Module Integration Checklist

### Pre-Integration Requirements
```
Each module must have:
[ ] 85%+ test coverage
[ ] API documentation (OpenAPI spec)
[ ] Database schema finalized
[ ] Error codes documented
[ ] Audit table implemented
[ ] Performance tested
[ ] Security audit passed
```

### Integration Steps
```
1. Schema Integration
   [ ] Create foreign keys between modules
   [ ] Add composite indexes
   [ ] Test data consistency

2. API Integration
   [ ] All endpoints using standard response
   [ ] Error codes documented
   [ ] Request/response validation

3. Data Flow Integration
   [ ] Define GL posting rules
   [ ] Test transaction flow
   [ ] Verify data consistency

4. Testing Integration
   [ ] End-to-end workflow tests
   [ ] Performance tests with load
   [ ] Security penetration tests

5. Deployment Integration
   [ ] All modules deploy together
   [ ] Database migrations in order
   [ ] Rollback procedures documented
```

---

## 📋 Module Dependency Graph

```
Core Dependencies:
├─ Accounts Module (GL)
│  └─ Required by: ALL other modules
│     (every transaction posts to GL)
│
├─ HR Module
│  ├─ Inputs to: Accounts (salary expense)
│  └─ Optional input from: Gamification (bonuses)
│
├─ Sales Module
│  ├─ Inputs to: Accounts (AR, Revenue)
│  ├─ Inputs to: Purchase (stock check)
│  └─ Triggers: Post Sales (service tickets)
│
├─ Purchase Module
│  ├─ Inputs to: Accounts (AP, Expense)
│  ├─ Link to: Inventory (stock)
│  └─ Optional: Gamification (rebates)
│
├─ Construction Module
│  ├─ Consumes: Purchase (materials)
│  ├─ Inputs to: Accounts (project costs)
│  ├─ Inputs to: HR (labor tracking)
│  └─ Generates: Post Sales (maintenance contracts)
│
├─ Post Sales Module
│  ├─ Triggered by: Sales (order completion)
│  ├─ Inputs to: Accounts (service revenue)
│  └─ Optional: Gamification (customer satisfaction)
│
└─ Civil Module
   ├─ Consumes: Construction (project data)
   ├─ Consumes: HR (worker data)
   └─ Inputs to: Accounts (site costs)

Recommendation: Deploy in this order:
1. Accounts (foundational)
2. HR, Sales, Purchase (core business)
3. Construction, Civil, Post Sales (domain-specific)
```

---

## 🐛 Debugging Tips

### Common Issues & Solutions

**Issue 1: GL Entry Not Posted After Transaction**
```
Debug Steps:
1. Check transaction status is "completed"
2. Verify GL posting trigger is enabled
3. Check GL posting queue for errors
4. Review GL posting logs
5. Manually trigger GL posting

Query:
SELECT * FROM gl_posting_queue 
WHERE status = 'failed' 
ORDER BY created_at DESC;
```

**Issue 2: Employee Salary Calculation Incorrect**
```
Debug Steps:
1. Check salary structure assigned to employee
2. Verify attendance is marked complete
3. Check allowances and deductions configured
4. Trace calculation step-by-step
5. Compare to salary slip

Query:
SELECT e.id, ss.*, a.*, d.* 
FROM employees e
JOIN salary_structures ss ON e.salary_structure_id = ss.id
LEFT JOIN allowances a ON ss.id = a.salary_structure_id
LEFT JOIN deductions d ON ss.id = d.salary_structure_id
WHERE e.id = ?;
```

**Issue 3: Purchase Order Not Converting to Invoice**
```
Debug Steps:
1. Verify GRN received and accepted
2. Check invoice matching rules
3. Verify GL posting rules
4. Check purchase module permissions

Query:
SELECT po.*, grn.*, vi.* 
FROM purchase_orders po
LEFT JOIN goods_receipts grn ON po.id = grn.po_id
LEFT JOIN vendor_invoices vi ON po.id = vi.po_id
WHERE po.id = ?;
```

**Issue 4: Performance Degradation**
```
Debug Steps:
1. Check query execution plans
2. Verify indexes are used
3. Check cache hit rates
4. Monitor database connections
5. Review slow query log

Query:
SHOW ENGINE INNODB STATUS\G
SELECT * FROM performance_schema.events_statements_summary_by_digest 
LIMIT 10;
```

---

## 📞 Support & Escalation

### Support Matrix
```
Issue Type          | First Level | Escalation
─────────────────────────────────────────────────
API Error           | Dev Team    | Backend Lead
GL Posting Issue    | Finance Dev | Accounts Lead
Performance         | DevOps      | Architect
Data Corruption     | DBA         | Engineering Manager
Security Breach     | Security    | CTO
Data Loss           | DBA         | VP Engineering
```

### Key Contacts
```
HR Module Lead:           [To be assigned]
Accounts Module Lead:     [To be assigned]
Sales Module Lead:        [To be assigned]
Purchase Module Lead:     [To be assigned]
Frontend Lead:            [To be assigned]
DevOps Lead:              [To be assigned]
QA Lead:                  [To be assigned]
```

---

## 📚 Additional Resources

### Documentation References
- Phase 3E Implementation Plan: BUSINESS_MODULES_IMPLEMENTATION_PLAN.md
- Schema Designs: Thoughts/schema_idea1/ (22 SQL files)
- Existing Codebase: Internal/models, services, handlers
- API Standards: COMPLETE_API_REFERENCE.md

### Tools & Commands

**Database Commands**
```bash
# Backup database
mysqldump -u root -p erp_db > backup_$(date +%s).sql

# Restore database
mysql -u root -p erp_db < backup_file.sql

# Run migrations
go run cmd/migrate/main.go up

# Check schema
mysql> SHOW TABLES LIKE 'employee%';
```

**API Testing**
```bash
# Test HR employee creation
curl -X POST http://localhost:8080/api/v2/hr/employees \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# List employees
curl -X GET "http://localhost:8080/api/v2/hr/employees?limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

**Git Commands**
```bash
# Create feature branch
git checkout -b feature/hr-payroll-module

# Commit changes
git commit -m "feat(hr): implement payroll calculation"

# Create pull request
git push origin feature/hr-payroll-module
# Then create PR on GitHub
```

---

**Generated**: November 24, 2025  
**Last Updated**: November 24, 2025  
**Status**: Ready for Team Distribution
