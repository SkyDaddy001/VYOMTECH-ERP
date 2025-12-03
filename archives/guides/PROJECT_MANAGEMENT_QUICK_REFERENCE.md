# Migration 022 - Quick Reference Guide

## 📚 What's New

### 7 New Tables + 2 Extended Tables

| Table | Purpose | Key Fields |
|-------|---------|-----------|
| `property_customer_profile` | Detailed customer KYC | PAN, Aadhar, income, employer, loan |
| `property_customer_unit_link` | Customer↔Unit mapping | booking_date, agreement_date, possession_date |
| `property_unit_area_statement` | Area breakup per unit | carpet_area, buildup_area, super_area |
| `property_payment_receipt` | Payment transactions | payment_mode, cheque#, transaction_id, gl_account_id |
| `property_project_milestone` | Construction phases | completion_status, budget_allocated, budget_spent |
| `property_project_activity` | Daily work logs | activity_type, activity_date, completion_percentage |
| `property_project_document` | Document management | document_type, approval_status, version_number |
| `property_project_summary` | KPI dashboard | units_sold, revenue_booked, gross_profit, margin_% |

---

## 🔄 Integration with Existing Tables

```
Migration 008 (Real Estate) ← EXTENDED BY → Migration 022

property_project (add: financial fields)
property_block  (unchanged)
property_unit   (add: area statements)
unit_cost_sheet (add: GST, charges, validity)
property_booking (linked via customer_unit_link)
installment (linked via payment_receipt)
```

---

## 🚀 API Endpoints to Build

### Customer Management
```
POST   /api/v1/customers/profiles               - Create customer
GET    /api/v1/customers/{id}                   - Get customer details
PUT    /api/v1/customers/{id}                   - Update customer
GET    /api/v1/customers?status=inquiry         - List customers by status
POST   /api/v1/customers/{id}/units/{unitId}    - Link customer to unit
```

### Payment Processing
```
POST   /api/v1/payments/receipts                - Create payment receipt
GET    /api/v1/payments/{id}                    - Get receipt details
GET    /api/v1/customers/{id}/payments          - Payment history
GET    /api/v1/units/{id}/payments              - Unit payment status
PUT    /api/v1/payments/{id}/status             - Update payment status
```

### Project Milestones
```
POST   /api/v1/projects/{id}/milestones         - Create milestone
GET    /api/v1/projects/{id}/milestones         - List milestones
PUT    /api/v1/milestones/{id}                  - Update milestone progress
GET    /api/v1/milestones/{id}/activities       - Get milestone activities
```

### Activity Logging
```
POST   /api/v1/projects/{id}/activities         - Log activity
GET    /api/v1/projects/{id}/activity-log       - Get activity timeline
GET    /api/v1/milestones/{id}/activities       - Milestone activities
```

### Documentation
```
POST   /api/v1/projects/{id}/documents          - Upload document
GET    /api/v1/projects/{id}/documents          - List documents
PUT    /api/v1/documents/{id}/approve           - Approve document
```

### Dashboard
```
GET    /api/v1/projects/{id}/summary            - Project KPI dashboard
GET    /api/v1/projects/{id}/summary/daily      - Daily summary
GET    /api/v1/projects?status=active           - Active projects list
GET    /api/v1/dashboard/pipeline               - Sales pipeline
```

---

## 📊 Key Queries Examples

### Get Customer with Units
```sql
SELECT c.*, cul.unit_id, pu.unit_number, pu.unit_type
FROM property_customer_profile c
JOIN property_customer_unit_link cul ON c.id = cul.customer_id
JOIN property_unit pu ON cul.unit_id = pu.id
WHERE c.tenant_id = ? AND c.customer_status = 'BOOKING_CONFIRMED';
```

### Payment Collection Summary
```sql
SELECT 
  SUM(pr.payment_amount) as total_collected,
  COUNT(*) as receipt_count,
  COUNT(DISTINCT pr.customer_id) as unique_customers
FROM property_payment_receipt pr
WHERE pr.tenant_id = ? 
  AND pr.payment_status IN ('RECEIVED', 'CLEARED')
  AND DATE(pr.payment_date) = CURDATE();
```

### Project Progress
```sql
SELECT 
  ppm.milestone_name,
  ppm.percentage_completion,
  ppm.completion_status,
  SUM(ppa.completion_percentage)/COUNT(*) as avg_activity_completion
FROM property_project_milestone ppm
LEFT JOIN property_project_activity ppa ON ppm.id = ppa.milestone_id
WHERE ppm.project_id = ?
GROUP BY ppm.id
ORDER BY ppm.start_date;
```

### Unit Revenue Analysis
```sql
SELECT 
  pu.unit_number,
  ucs.grand_total as unit_price,
  SUM(COALESCE(pr.payment_amount, 0)) as collected,
  ucs.grand_total - SUM(COALESCE(pr.payment_amount, 0)) as pending
FROM property_unit pu
LEFT JOIN unit_cost_sheet ucs ON pu.id = ucs.unit_id
LEFT JOIN property_payment_receipt pr ON pu.id = pr.unit_id
WHERE pu.project_id = ?
GROUP BY pu.id;
```

---

## 🔐 Multi-Tenancy Enforcement

Every table has:
- `tenant_id` field (foreign key to `tenant(id)`)
- Index on `(tenant_id, ...)` for query optimization
- Data isolation by design

**Example service pattern**:
```go
// SECURE: Tenant data isolation
func (s *Service) GetCustomer(ctx context.Context, tenantID, customerID string) (*PropertyCustomerProfile, error) {
    var customer PropertyCustomerProfile
    return customer, s.db.WithContext(ctx).
        Where("tenant_id = ? AND id = ?", tenantID, customerID).
        First(&customer).
        Error
}
```

---

## 💰 GL Accounting Integration

Payment receipts connect to GL automatically:

```go
// When creating payment receipt:
receipt := &PropertyPaymentReceipt{
    CustomerID: customerID,
    PaymentAmount: 100000,
    GLAccountID: "AR-001",  // Accounts Receivable
}

// Service posts to GL:
glEntry := &GLPosting{
    AccountID: "AR-001",
    DebitAmount: 100000,
    Reference: receipt.ReceiptNumber,
}
```

---

## 📈 KPI Dashboard Fields

```
property_project_summary tracks:
✓ Units sold / Units available (inventory)
✓ Revenue booked / Revenue received (cash flow)
✓ Construction cost / Cost incurred (budget control)
✓ Gross profit / Margin percentage (profitability)
✓ Project completion percentage (progress)
✓ Active milestones / Delayed milestones (timeline)
✓ Customer satisfaction score (quality)
```

---

## 🛠️ Development Workflow

### Step 1: Create Service Layer
```
File: internal/services/project_management.go
- CustomerProfileService (CRUD)
- PaymentReceiptService (create, reconcile)
- MilestoneService (create, update progress)
- ProjectActivityService (log activities)
- SummaryService (generate KPIs)
```

### Step 2: Create Handlers
```
File: internal/handlers/project_management.go
- CustomerHandlers
- PaymentHandlers
- MilestoneHandlers
- ActivityHandlers
- DocumentHandlers
- DashboardHandlers
```

### Step 3: Register Routes
```
router.POST("/api/v1/customers/profiles", handlers.CreateCustomer)
router.GET("/api/v1/customers/:id", handlers.GetCustomer)
router.POST("/api/v1/payments/receipts", handlers.CreatePayment)
// ... more routes
```

### Step 4: Update Docker
```
Add to docker-compose.yml:
volumes:
  - ./migrations/022_project_management_extensions.sql:/docker-entrypoint-initdb.d/22-project-management.sql
```

---

## ✅ Validation Checklist

- [ ] Migration 022 SQL reviewed
- [ ] Models file matches schema
- [ ] Service layer created
- [ ] Handlers implemented
- [ ] Routes registered
- [ ] Docker-compose updated
- [ ] Unit tests written
- [ ] Integration tests running
- [ ] GL posting tested
- [ ] Multi-tenant isolation verified
- [ ] API documentation ready

---

## 📞 Related Migrations

| Migration | Module | Status |
|-----------|--------|--------|
| 001 | Foundation (Tenant, User, Auth) | ✅ Existing |
| 005 | GL Accounting | ✅ Existing |
| 008 | Real Estate | ✅ Existing |
| 022 | **Project Management** | ✅ **New** |

---

## 🎯 Success Metrics

After implementation, you should have:

✅ Create customers with full KYC details  
✅ Link customers to units (individual, joint, corporate)  
✅ Track payment receipts with GL posting  
✅ Monitor construction milestones with budget  
✅ Log daily project activities  
✅ Generate project KPI dashboard  
✅ Multi-tenant data isolation  
✅ Complete audit trails  
✅ Production-ready performance  

---

**Version**: 1.0  
**Created**: December 3, 2025  
**Status**: Ready for Service/Handler Implementation
