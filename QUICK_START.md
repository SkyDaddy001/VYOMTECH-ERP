# VYOM ERP - Quick Start Guide

## 🚀 Getting Started

### Prerequisites
- Go 1.24+
- Node.js 18+
- MySQL 8.0+
- Docker (optional)

### Installation

```bash
# Clone repository
git clone https://github.com/SkyDaddy001/VYOMTECH-ERP.git
cd VYOM-ERP

# Backend setup
go mod download
go build -o main ./cmd/main.go

# Frontend setup
cd frontend
npm install
npm run build

# Database setup
mysql -u root -p < migrations/001_initial_schema.sql
```

### Running Locally

```bash
# Terminal 1: Backend
go run cmd/main.go
# Server running at http://localhost:8080

# Terminal 2: Frontend
cd frontend
npm run dev
# UI available at http://localhost:3000
```

---

## 📚 Core Modules Overview

### 1. General Ledger (GL)
**Purpose:** Double-entry accounting system
- Chart of Accounts hierarchy
- Journal entries with automatic balancing
- Financial period management
- Audit trail

**Key Endpoints:**
```
GET    /api/v1/gl/chart-of-accounts
POST   /api/v1/gl/journal-entries
GET    /api/v1/gl/entries/{period}
POST   /api/v1/gl/post-entry
```

### 2. Accounts Payable (AP)
**Purpose:** Vendor management & payments
- Vendor master
- Purchase orders
- Invoice matching
- Payment scheduling

**Key Endpoints:**
```
POST   /api/v1/ap/vendors
GET    /api/v1/ap/invoices
POST   /api/v1/ap/schedule-payment
GET    /api/v1/ap/ageing-report
```

### 3. Accounts Receivable (AR)
**Purpose:** Customer invoicing & collections
- Customer master
- Sales invoicing
- Collection tracking
- Dunning management

**Key Endpoints:**
```
POST   /api/v1/ar/customers
POST   /api/v1/ar/invoices
GET    /api/v1/ar/collections
GET    /api/v1/ar/aging-report
```

### 4. HR & Payroll
**Purpose:** Employee management & salary processing
- Employee master with full details
- Salary components
- Payroll processing
- Statutory deductions (TDS)

**Key Endpoints:**
```
POST   /api/v1/hr/employees
POST   /api/v1/hr/payroll/process
GET    /api/v1/hr/payslips
POST   /api/v1/hr/attendance/mark
```

### 5. Leave Management
**Purpose:** Leave request & tracking
- Leave types (Annual, Casual, Sick, Maternity)
- Leave balance management
- Approval workflows
- Carry-forward rules

**Key Endpoints:**
```
POST   /api/v1/leave/request
GET    /api/v1/leave/balance
POST   /api/v1/leave/approve
GET    /api/v1/leave/history
```

### 6. Sales
**Purpose:** Sales pipeline & invoicing
- Opportunities management
- Sales invoicing
- Payment tracking
- Customer interaction

**Key Endpoints:**
```
POST   /api/v1/sales/opportunities
POST   /api/v1/sales/invoices
GET    /api/v1/sales/pipeline
GET    /api/v1/sales/forecast
```

### 7. Real Estate
**Purpose:** Project management (RERA compliant)
- Project setup with segregated accounts
- Collection tracking
- Fund utilization logging
- Borrowing management

**Key Endpoints:**
```
POST   /api/v1/realestate/projects
POST   /api/v1/realestate/collections
GET    /api/v1/realestate/fund-status
POST   /api/v1/realestate/reconciliation
```

### 8. Construction
**Purpose:** Material & contractor management
- BOQ (Bill of Quantities)
- Material tracking
- Contractor management
- Progress billing

**Key Endpoints:**
```
POST   /api/v1/construction/boq
POST   /api/v1/construction/material-orders
GET    /api/v1/construction/progress
POST   /api/v1/construction/contractor-bill
```

### 9. Purchase
**Purpose:** Vendor & material procurement
- Purchase orders
- Material receipt
- Vendor performance
- Invoice matching

**Key Endpoints:**
```
POST   /api/v1/purchase/orders
POST   /api/v1/purchase/goods-receipt
GET    /api/v1/purchase/vendor-performance
GET    /api/v1/purchase/pending-invoices
```

### 10. Compliance
**Purpose:** Regulatory compliance tracking
- RERA compliance
- Tax compliance (GST, TDS, ITR)
- Labour law compliance (ESI, EPF)
- Audit trails

**Key Endpoints:**
```
GET    /api/v1/compliance/rera-status
POST   /api/v1/compliance/tax-filing
GET    /api/v1/compliance/compliance-calendar
POST   /api/v1/compliance/audit-report
```

### 11. Dashboard
**Purpose:** Executive analytics & reporting
- Financial dashboard
- HR dashboard
- Sales dashboard
- Compliance dashboard

**Key Endpoints:**
```
GET    /api/v1/dashboard/financial/ratios
POST   /api/v1/dashboard/hr/payroll
GET    /api/v1/dashboard/sales/pipeline
GET    /api/v1/dashboard/compliance/health-score
```

---

## 🔐 Authentication

### Login
```bash
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {...}
}
```

### Using Token
```bash
Header: Authorization: Bearer {token}
Header: X-Tenant-ID: {tenant-id}
```

---

## 📊 Dashboard Access

### Financial Dashboard
```bash
GET /api/v1/dashboard/financial/ratios
Headers:
  - Authorization: Bearer {token}
  - X-Tenant-ID: {tenant-id}

Response: {
  "timestamp": "2024-12-02T10:30:00Z",
  "data": {
    "current_ratio": 1.5,
    "quick_ratio": 1.2,
    "debt_to_equity": 0.8,
    "roe": 0.25,
    "roa": 0.15
  }
}
```

### HR Dashboard
```bash
POST /api/v1/dashboard/hr/payroll
Headers:
  - Authorization: Bearer {token}
  - X-Tenant-ID: {tenant-id}

Body: {
  "payroll_month": "2024-12"
}

Response: {
  "timestamp": "2024-12-02T10:30:00Z",
  "data": {
    "total_salary": 5000000,
    "by_department": {...}
  }
}
```

### Sales Dashboard
```bash
GET /api/v1/dashboard/sales/overview
Headers:
  - Authorization: Bearer {token}
  - X-Tenant-ID: {tenant-id}

Response: {
  "timestamp": "2024-12-02T10:30:00Z",
  "data": {
    "ytd_revenue": 10000000,
    "pipeline_value": 5000000,
    "conversion_rate": 0.35
  }
}
```

### Compliance Dashboard
```bash
GET /api/v1/dashboard/compliance/health-score
Headers:
  - Authorization: Bearer {token}
  - X-Tenant-ID: {tenant-id}

Response: {
  "timestamp": "2024-12-02T10:30:00Z",
  "data": {
    "health_score": 95,
    "risk_level": "Green",
    "critical_items": 0
  }
}
```

---

## 🗄️ Database Schema

### Key Tables

**Chart of Accounts**
```sql
SELECT * FROM chart_of_accounts 
WHERE tenant_id = 'org-123' AND deleted_at IS NULL;
```

**GL Entries**
```sql
SELECT * FROM gl_entries 
WHERE tenant_id = 'org-123' AND entry_date BETWEEN '2024-01-01' AND '2024-12-31';
```

**Employees**
```sql
SELECT * FROM employees 
WHERE tenant_id = 'org-123' AND status = 'Active';
```

**Sales Invoices**
```sql
SELECT * FROM sales_invoices 
WHERE tenant_id = 'org-123' AND invoice_date >= DATE_SUB(NOW(), INTERVAL 90 DAY);
```

---

## 🚀 Deployment Options

### Docker
```bash
docker build -t vyom-erp:latest .
docker run -p 8080:8080 -e DB_HOST=mysql vyom-erp:latest
```

### Kubernetes
```bash
kubectl apply -f k8s/
kubectl port-forward svc/vyom-api 8080:8080
```

### AWS
```bash
# ECR
aws ecr create-repository --repository-name vyom-erp
docker tag vyom-erp:latest {aws-account}.dkr.ecr.ap-south-1.amazonaws.com/vyom-erp:latest
docker push {aws-account}.dkr.ecr.ap-south-1.amazonaws.com/vyom-erp:latest

# ECS/EKS deployment
aws ecs create-service --cluster vyom-prod --service-name vyom-api --task-definition vyom:1
```

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run specific module tests
go test ./internal/services -v

# Integration tests
go test ./internal/handlers -v

# Coverage report
go test -cover ./...
```

---

## 📝 API Documentation

For complete API documentation, see:
- `/docs/API_REFERENCE.md` - All 176+ endpoints
- `/docs/SYSTEM_ARCHITECTURE.md` - Technical design
- `/docs/archive/` - Detailed phase documentation

---

## 🐛 Troubleshooting

### Build Issues
```bash
# Clear Go cache
go clean -cache

# Rebuild
go build -o main ./cmd/main.go
```

### Database Connection
```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1;"

# Reset database
mysql -u root -p < migrations/001_initial_schema.sql
```

### Multi-Tenant Issues
```bash
# Ensure X-Tenant-ID header is present
curl -H "X-Tenant-ID: tenant-123" http://localhost:8080/api/v1/dashboard/financial/ratios
```

---

## 📞 Support & Documentation

- **GitHub:** github.com/SkyDaddy001/VYOMTECH-ERP
- **Documentation:** See `/docs/` directory
- **Archived Docs:** See `/docs/archive/` for detailed phase documentation

---

## ✅ Checklist for Production

- [ ] Database backups configured
- [ ] SSL/TLS certificates installed
- [ ] API keys securely stored
- [ ] Monitoring & alerts set up
- [ ] Log aggregation configured
- [ ] Rate limiting enabled
- [ ] Multi-tenant isolation verified
- [ ] Compliance requirements met

---

**Status:** ✅ Production Ready
**Last Updated:** December 2, 2025
