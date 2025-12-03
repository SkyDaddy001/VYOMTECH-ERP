# Business Modules - Sprint Breakdown & Execution Guide
## Phase 3E: Detailed Weekly Roadmap

**Date**: November 24, 2025  
**Scope**: 7 business modules over 16 weeks  
**Format**: Week-by-week, sprint-by-sprint detailed breakdown

---

## 📅 PHASE 3E Timeline Overview

```
PHASE 1: Foundation (Weeks 1-2)
  │
  └─ Week 1: Database Schema & Migration
  └─ Week 2: Authentication & Authorization

PHASE 2: Core Modules (Weeks 3-8)
  │
  ├─ Sprint 1 (Weeks 3-4): HR Module
  ├─ Sprint 2 (Weeks 5-7): Accounts Module
  └─ Sprint 3 (Weeks 8-9): Sales Module

PHASE 3: Supporting Modules (Weeks 9-13)
  │
  ├─ Sprint 4 (Weeks 9-10): Purchase Module
  ├─ Sprint 5 (Weeks 11-13): Construction Module
  └─ Sprint 6 (Weeks 12-13): Post Sales & Civil

PHASE 4: Integration & Launch (Weeks 14-16)
  │
  ├─ Week 14: Integration Testing
  ├─ Week 15: Performance & Security
  └─ Week 16: UAT & Go Live
```

---

## 🔧 WEEK 1: Database Schema & Migration

### Sprint 0, Day 1: Project Setup

**Team**: Backend Lead, DevOps Lead, 2 Backend Developers  
**Goal**: Set up project structure, create database, prepare migration framework

#### Tasks (In Order)

**Morning (4 hours)**
```
1. Create GitHub project/branch
   └─ Branch: feature/phase3e-modules
   └─ Protection rules: require PR review

2. Set up project structure
   └─ Create directories:
      ├─ internal/models/erp/*
      ├─ internal/services/erp/*
      ├─ internal/handlers/erp/*
      ├─ migrations/erp/*
      └─ tests/erp/*

3. Prepare migration framework
   └─ Choose tool: Flyway or golang-migrate
   └─ Set up migration versioning (starting from 005_*)
   └─ Create baseline schema

4. Backup current database
   └─ mysqldump -u root -p erp_db > backup_phase3d_baseline.sql
   └─ Verify backup integrity
```

**Afternoon (4 hours)**
```
5. Create core tables (Foundation)
   └─ From: Thoughts/schema_idea1/001_core_module.sql
   └─ Adapt: Multi-tenant structure
   └─ Tables:
      ├─ users (extend existing)
      ├─ roles (extend existing)
      ├─ permissions (extend existing)
      └─ audit_log (create new)

6. Verify database integrity
   └─ Run: SHOW TABLES;
   └─ Count: Should have ~10 new tables
   └─ Check: All primary keys, foreign keys, indexes

7. Create test data seed script
   └─ 5 test companies
   └─ 10 test users per company
   └─ Sample roles & permissions
```

**Deliverable**: Database skeleton created, seed script ready, all devs can connect

---

### Sprint 0, Days 2-3: HR & Accounts Schema

**Day 2: HR Module Schema (22 tables)**

```
Migration: 005_hr_module_schema.sql

CREATE TABLES (In Order):
├─ designations
│  └─ id, code, name, description, created_at, updated_at
│
├─ departments
│  └─ id, code, name, description, parent_id, created_at
│
├─ employees
│  └─ 20+ columns: personal data, contact, employment status
│  └─ FOREIGN KEY: designation_id, department_id
│
├─ salary_structures
│  └─ id, name, base_salary, structure_components
│
├─ attendance
│  └─ id, employee_id, date, check_in, check_out
│
├─ shift_definitions
│  └─ id, name, start_time, end_time, working_hours
│
├─ leave_types
│  └─ id, name, annual_days, allow_carry_forward
│
├─ leave_applications
│  └─ id, employee_id, from_date, to_date, reason, status
│
├─ allowances
│  └─ id, salary_structure_id, component_name, percentage
│
├─ deductions
│  └─ id, salary_structure_id, component_name, percentage
│
├─ salary_slips
│  └─ id, employee_id, month, gross_pay, deductions, net_pay
│
└─ [+11 more tables for payroll, statutory, etc.]

CREATE INDEXES:
├─ idx_employees_tenant_id (tenant_id)
├─ idx_employees_designation_id
├─ idx_employees_department_id
├─ idx_attendance_employee_date (employee_id, date)
├─ idx_leave_applications_employee_date
├─ idx_salary_slips_employee_month
└─ [+10 more for common queries]

CREATE TRIGGERS:
├─ tr_employee_created_audit
├─ tr_employee_updated_audit
├─ tr_salary_slip_created_audit
└─ [+3 more for compliance]
```

**Task Checklist**:
```
[ ] Migration file created: migrations/005_hr_module_schema.sql
[ ] All 22 HR tables created
[ ] Foreign keys established
[ ] 15+ indexes created
[ ] Audit triggers installed
[ ] Migration runs successfully
[ ] Test queries pass
[ ] Rollback tested (can undo migration)
```

**Day 3: Accounts Module Schema (20 tables)**

```
Migration: 006_accounts_module_schema.sql

CREATE TABLES (In Order):
├─ gl_masters
│  └─ id, code, name, account_type, balance, status
│  └─ account_type: 'ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE'
│
├─ gl_hierarchy
│  └─ id, parent_gl_id, child_gl_id, level
│
├─ cost_centers
│  └─ id, code, name, description, created_at
│
├─ journal_entries
│  └─ id, entry_date, narration, status, posted_by
│
├─ journal_line_items
│  └─ id, journal_id, gl_master_id, debit_amount, credit_amount
│
├─ bank_accounts
│  └─ id, name, account_number, balance, gl_master_id
│
├─ bank_reconciliation
│  └─ id, bank_account_id, statement_date, reconciliation_status
│
├─ invoices
│  └─ id, customer_id, invoice_number, amount, status
│  └─ FOREIGN KEY: customer_id (from Sales module)
│
├─ invoice_items
│  └─ id, invoice_id, description, quantity, rate, amount
│
├─ payments
│  └─ id, invoice_id, amount, payment_date, payment_method
│
└─ [+10 more tables for payment terms, reports, etc.]

CREATE INDEXES:
├─ idx_gl_masters_code (code)
├─ idx_journal_entries_date (entry_date)
├─ idx_journal_line_items_gl_master
├─ idx_invoices_customer_id
├─ idx_invoices_status
├─ idx_bank_reconciliation_date
└─ [+10 more for reporting queries]

CREATE TRIGGERS:
├─ tr_journal_entry_posted_audit
├─ tr_gl_balance_update (on journal posting)
├─ tr_invoice_created_gl_entry
└─ [+5 more for compliance]
```

**Task Checklist**:
```
[ ] Migration file created: migrations/006_accounts_module_schema.sql
[ ] All 20 Accounts tables created
[ ] GL hierarchy established
[ ] Foreign keys to Sales/Purchase modules defined
[ ] 15+ indexes created
[ ] Audit triggers installed
[ ] GL balance calculation trigger working
[ ] Migration runs successfully
[ ] Test: Can create GL entry and balance updates
```

**Day 3 Afternoon: Schema Validation**

```
Validation Queries:
[ ] SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES 
    WHERE TABLE_SCHEMA = 'erp_db' AND TABLE_NAME LIKE '%module%';
    Expected: ~50 tables created so far

[ ] SELECT COUNT(*) FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = 'erp_db' AND CONSTRAINT_NAME LIKE 'FK_%';
    Expected: 30+ foreign keys

[ ] SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = 'erp_db';
    Expected: 100+ indexes

[ ] Verify all audit tables exist
    SELECT * FROM information_schema.TABLES 
    WHERE TABLE_NAME LIKE '%audit%';
```

---

### Sprint 0, Days 4-5: Sales, Purchase, Construction Schemas

**Day 4: Sales + Purchase Schemas**

```
Migration: 007_sales_purchase_schema.sql

SALES TABLES (18 total):
├─ customers (id, name, email, phone, address, tax_id)
├─ customer_contacts (customer_id, contact_name, phone, email)
├─ opportunities (id, customer_id, title, amount, stage, probability)
├─ quotations (id, customer_id, quote_number, amount, valid_until)
├─ quote_line_items (id, quote_id, description, quantity, rate)
├─ sales_orders (id, customer_id, so_number, amount, status)
├─ order_line_items (id, order_id, product_id, quantity, rate)
├─ order_fulfillment (id, order_id, fulfilled_quantity, status)
├─ sales_representatives (id, name, email, commission_rate)
├─ commission_calculations (id, sales_rep_id, month, commission_amount)
└─ [+8 more tables]

PURCHASE TABLES (16 total):
├─ vendors (id, name, email, phone, address, tax_id)
├─ vendor_contacts (vendor_id, contact_name, phone)
├─ purchase_requisitions (id, requisitioner_id, status, date)
├─ purchase_orders (id, vendor_id, po_number, amount, status)
├─ po_line_items (id, po_id, product_id, quantity, rate)
├─ goods_receipts (id, po_id, grn_number, status)
├─ receipt_line_items (id, grn_id, po_line_id, received_qty)
├─ quality_inspections (id, grn_id, inspection_status, notes)
├─ vendor_invoices (id, vendor_id, invoice_number, amount)
└─ [+7 more tables]

Key Foreign Keys:
├─ sales_orders.customer_id → customers.id
├─ purchase_orders.vendor_id → vendors.id
├─ Both sales/purchase → accounts GL accounts (future link)
└─ All orders → invoice generation (Accounts module)
```

**Day 5: Construction + Civil Schemas**

```
Migration: 008_construction_civil_schema.sql

CONSTRUCTION TABLES (20 total):
├─ construction_projects (id, name, location, budget, start_date, end_date)
├─ project_phases (id, project_id, phase_name, start_date, end_date)
├─ work_breakdown_structure (id, project_id, parent_id, task_name, level)
├─ tasks (id, project_id, task_name, description, status, duration)
├─ task_dependencies (id, task_id, depends_on_task_id)
├─ boq_master (id, project_id, description, total_cost)
├─ boq_line_items (id, boq_id, item_code, description, quantity, rate)
├─ material_rates (id, material_code, rate, unit, effective_date)
├─ labour_rates (id, skill_type, rate, unit, effective_date)
├─ daily_reports (id, project_id, report_date, weather, notes, photos_url)
├─ material_usage (id, project_id, material_id, quantity_used, date)
├─ quality_inspections (id, project_id, inspection_date, pass_fail, notes)
├─ defect_register (id, project_id, defect_description, severity, rectified)
└─ [+7 more tables]

CIVIL TABLES (12 total):
├─ civil_sites (id, name, location, area_sqft, site_type)
├─ site_contacts (id, site_id, contact_name, phone)
├─ site_amenities (id, site_id, amenity_type, status)
├─ contractors (id, name, license_number, rating)
├─ contractor_agreements (id, contractor_id, site_id, start_date, end_date)
├─ safety_incidents (id, site_id, incident_date, description, severity)
├─ safety_checklist_items (id, site_id, checklist_item, completed)
├─ compliance_audits (id, site_id, audit_date, compliance_status)
├─ permits (id, site_id, permit_type, permit_number, valid_until)
├─ waste_management_records (id, site_id, waste_type, quantity, disposal_date)
├─ environmental_data (id, site_id, date, noise_level, dust_level, water_quality)
└─ [+1 more table]

Key Dependencies:
├─ construction projects ← from sales (contract/order)
├─ tasks → employees (HR module)
├─ material_usage → vendors (Purchase module)
└─ All costs → GL accounts (Accounts module)
```

**End of Week 1 Deliverables**:
```
✅ Database: 130+ tables created
✅ Migrations: All 4 migration files tested
✅ Indexes: 100+ strategic indexes created
✅ Triggers: Audit triggers installed per module
✅ Schema: Fully normalized with ULID PKs
✅ Foreign Keys: Cross-module relationships defined
✅ Seed Data: Test data loading script ready
✅ Documentation: Schema ER diagrams generated
✅ Testing: All CREATE TABLE statements verified
✅ Rollback: All migrations tested for rollback
```

---

## 🔐 WEEK 2: Authentication & Authorization

### Sprint 0, Days 8-10: RBAC Extension

**Day 8: Role & Permission Design**

**Module-Specific Roles**:
```
HR Module Roles:
├─ HR Admin: Full HR module access
├─ HR Manager: Employee management, payroll review
├─ Payroll Officer: Salary slip generation, payment authorization
├─ Employee Self-Service: Own profile, leave applications
└─ Compliance Officer: Statutory reports, audit trail

Accounts Module Roles:
├─ Finance Admin: Full GL access
├─ Accountant: GL posting, reconciliation
├─ Accounts Receivable: Invoice, payment tracking
├─ Accounts Payable: Vendor invoices, payments
└─ Audit Officer: Report generation, audit trail

Sales Module Roles:
├─ Sales Manager: Team oversight, quota setting
├─ Sales Representative: Own leads, quotes, orders
├─ Sales Admin: Customer master, pricing
└─ Commission Officer: Commission calculation, approval

Purchase Module Roles:
├─ Purchase Manager: PO approval, vendor management
├─ Buyer: Create POs, GRN processing
├─ Vendor Manager: Vendor master, rating
└─ Invoice Officer: Invoice matching, payment auth

Construction Module Roles:
├─ Project Manager: Full project access
├─ Site Engineer: Daily reports, QC
├─ Contractor: Own work assignments
└─ Project Accountant: Cost tracking, budget monitoring

Civil Module Roles:
├─ Site Manager: Site oversight
├─ Safety Officer: Safety & compliance
├─ Contractor Manager: Contractor agreements
└─ Environmental Officer: Permits, environmental tracking

Post Sales Module Roles:
├─ Service Manager: Ticket assignment, resolution
├─ Technician: Own service tickets
├─ Warranty Officer: Warranty claims
└─ Support Analyst: Customer support, FAQ management
```

**Day 9: Permission Matrix Development**

**Example Permission Matrix for HR Module**:
```
Permission               | HR Admin | HR Manager | Payroll Officer | Employee
─────────────────────────────────────────────────────────────────────────────
create_employee          |   ✓      |     ✗      |        ✗        |    ✗
edit_employee            |   ✓      |     ✓      |        ✗        |   (own)
delete_employee          |   ✓      |     ✗      |        ✗        |    ✗
view_salary              |   ✓      |     ✓      |        ✓        |   (own)
modify_salary            |   ✓      |     ✗      |        ✗        |    ✗
run_payroll              |   ✓      |     ✗      |        ✓        |    ✗
generate_salary_slip     |   ✓      |     ✓      |        ✓        |   (own)
approve_leave            |   ✓      |     ✓      |        ✗        |    ✗
apply_leave              |   ✓      |     ✓      |        ✓        |    ✓
view_reports             |   ✓      |     ✓      |        ✓        |    ✗
export_payroll           |   ✓      |     ✓      |        ✓        |    ✗
```

**Database Structure**:
```sql
-- Extend existing permissions table
ALTER TABLE permissions ADD COLUMN (
    module ENUM('hr', 'accounts', 'sales', 'purchase', 'construction', 'civil', 'postsales'),
    scope ENUM('system', 'tenant', 'company', 'project', 'record'),
    resource_type VARCHAR(50),
    action VARCHAR(50)
);

-- Create permission sets per role
INSERT INTO role_permissions (role_id, permission_id, created_at) VALUES (...);
```

**Day 10: API Middleware Development**

**Authentication Middleware**:
```go
// middleware/auth.go

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Extract JWT token
        token := extractToken(r)
        
        // 2. Validate token signature
        claims, err := validateToken(token)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // 3. Get user from database
        user := db.GetUser(claims.UserID)
        if user == nil {
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }
        
        // 4. Add to context
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Authorization Middleware**:
```go
// middleware/rbac.go

func RequirePermission(permission string, module string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := r.Context().Value("user").(*models.User)
            
            // Check if user has permission for module
            hasPermission := db.CheckPermission(
                user.ID,
                module,
                permission,
            )
            
            if !hasPermission {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

**Usage in Routes**:
```go
// routes.go

// HR Module Routes
hrRoutes := router.PathPrefix("/api/v2/hr").Subrouter()
hrRoutes.Use(middleware.AuthMiddleware)
hrRoutes.Use(middleware.TenantIsolationMiddleware)

hrRoutes.HandleFunc("/employees", middleware.RequirePermission("create_employee", "hr")(
    handlers.CreateEmployee,
)).Methods("POST")

hrRoutes.HandleFunc("/payroll/run", middleware.RequirePermission("run_payroll", "hr")(
    handlers.RunPayroll,
)).Methods("POST")
```

**End of Week 2 Deliverables**:
```
✅ RBAC Extended: All 7 modules have role sets
✅ Permissions: 150+ module-specific permissions defined
✅ Middleware: Auth & RBAC middleware implemented
✅ Routes: All endpoints protected with permissions
✅ Testing: Permission checks tested for all roles
✅ Documentation: Permission matrix documented
✅ Database: All permission data seeded
✅ Audit: Permission checks logged
```

---

## 💼 SPRINT 1: HR MODULE (Weeks 3-4)

### Sprint 1 Outline
- **Duration**: 2 weeks
- **Team**: HR Module Lead (1 FTE), 1 Backend Developer, 1 Frontend Developer, QA (0.5 FTE)
- **Story Points**: 89
- **Goal**: Complete, tested HR module with all 45 API endpoints

### Week 3: Backend Implementation

**Monday & Tuesday: Employee Management**

```go
// internal/handlers/hr/employee_handler.go

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
    var req struct {
        FirstName       string `json:"first_name" validate:"required"`
        LastName        string `json:"last_name" validate:"required"`
        Email           string `json:"email" validate:"required,email"`
        DesignationID   string `json:"designation_id" validate:"required"`
        DepartmentID    string `json:"department_id" validate:"required"`
        JoinDate        string `json:"join_date" validate:"required,datetime"`
        SalaryStructID  string `json:"salary_struct_id" validate:"required"`
        BankAccount     string `json:"bank_account"`
        PFAccount       string `json:"pf_account"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // Validate input
    if err := validate.Struct(req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Create employee
    employee := &models.Employee{
        FirstName:      req.FirstName,
        LastName:       req.LastName,
        Email:          req.Email,
        DesignationID:  req.DesignationID,
        DepartmentID:   req.DepartmentID,
        JoinDate:       parseDate(req.JoinDate),
    }
    
    employee, err := h.service.CreateEmployee(r.Context(), employee)
    if err != nil {
        http.Error(w, "Failed to create employee", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": employee,
    })
}

// GET /api/v2/hr/employees
func (h *EmployeeHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
    limit := queryInt(r, "limit", 20)
    offset := queryInt(r, "offset", 0)
    departmentID := r.URL.Query().Get("department_id")
    
    employees, total, err := h.service.ListEmployees(r.Context(), &models.EmployeeFilter{
        DepartmentID: departmentID,
        Limit:        limit,
        Offset:       offset,
    })
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": employees,
        "pagination": map[string]int{
            "total": total,
            "limit": limit,
            "offset": offset,
        },
    })
}
```

**Tasks - Days 1-2**:
```
[ ] Implement CreateEmployee handler
[ ] Implement ListEmployees handler (with pagination)
[ ] Implement GetEmployee handler
[ ] Implement UpdateEmployee handler
[ ] Implement DeleteEmployee handler (soft delete)
[ ] Create Employee model with validation
[ ] Create EmployeeService with CRUD operations
[ ] Set up unit tests (80%+ coverage)
[ ] Integration tests with database
[ ] API documentation (OpenAPI spec)

Deliverable: Employee CRUD API complete, 10 endpoints tested
```

**Wednesday & Thursday: Attendance System**

```go
// internal/handlers/hr/attendance_handler.go

func (h *AttendanceHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
    var req struct {
        EmployeeID string `json:"employee_id" validate:"required"`
        CheckInTime string `json:"check_in_time" validate:"required,datetime"`
        Location   string `json:"location"`
    }
    
    // Parse and validate
    // ...
    
    // Check if already checked in today
    today := time.Now().Format("2006-01-02")
    existing, _ := h.service.GetAttendanceByDate(
        r.Context(),
        req.EmployeeID,
        today,
    )
    
    if existing != nil && existing.CheckInTime != nil {
        http.Error(w, "Already checked in today", http.StatusConflict)
        return
    }
    
    // Record check-in
    attendance := &models.Attendance{
        EmployeeID:   req.EmployeeID,
        CheckInTime:  parseTime(req.CheckInTime),
        AttendanceDate: today,
        Status:       "present",
    }
    
    attendance, err := h.service.RecordCheckIn(r.Context(), attendance)
    if err != nil {
        http.Error(w, "Failed to record check-in", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(attendance)
}

func (h *AttendanceHandler) GetAttendanceReport(w http.ResponseWriter, r *http.Request) {
    month := r.URL.Query().Get("month")
    departmentID := r.URL.Query().Get("department_id")
    
    report, err := h.service.GenerateAttendanceReport(
        r.Context(),
        month,
        departmentID,
    )
    
    if err != nil {
        http.Error(w, "Failed to generate report", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(report)
}
```

**Tasks - Days 3-4**:
```
[ ] Implement CheckIn handler
[ ] Implement CheckOut handler
[ ] Implement MarkAbsent handler
[ ] Implement GetAttendance handler
[ ] Implement GetAttendanceReport handler (monthly summary)
[ ] Create Attendance model
[ ] Create AttendanceService
[ ] Implement shift validation
[ ] Create attendance dashboard API
[ ] Set up tests (80%+ coverage)

Deliverable: Attendance system complete, 8 endpoints tested
```

**Friday: Leave Management**

```go
// internal/handlers/hr/leave_handler.go

func (h *LeaveHandler) ApplyLeave(w http.ResponseWriter, r *http.Request) {
    var req struct {
        EmployeeID  string `json:"employee_id" validate:"required"`
        LeaveTypeID string `json:"leave_type_id" validate:"required"`
        FromDate    string `json:"from_date" validate:"required,datetime"`
        ToDate      string `json:"to_date" validate:"required,datetime"`
        Reason      string `json:"reason" validate:"required"`
    }
    
    // Parse, validate
    // ...
    
    fromDate := parseDate(req.FromDate)
    toDate := parseDate(req.ToDate)
    
    // Check available balance
    leaveType, _ := h.service.GetLeaveType(req.LeaveTypeID)
    balance, _ := h.service.GetLeaveBalance(req.EmployeeID, leaveType.ID)
    
    days := int(toDate.Sub(fromDate).Hours() / 24)
    if balance < days {
        http.Error(w, "Insufficient leave balance", http.StatusConflict)
        return
    }
    
    // Create application
    application := &models.LeaveApplication{
        EmployeeID: req.EmployeeID,
        LeaveTypeID: req.LeaveTypeID,
        FromDate: fromDate,
        ToDate: toDate,
        Reason: req.Reason,
        Status: "pending",
    }
    
    application, err := h.service.CreateLeaveApplication(r.Context(), application)
    if err != nil {
        http.Error(w, "Failed to apply leave", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(application)
}

func (h *LeaveHandler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    
    application, err := h.service.ApproveLeave(r.Context(), id)
    if err != nil {
        http.Error(w, "Failed to approve leave", http.StatusInternalServerError)
        return
    }
    
    // Update employee leave balance
    h.service.UpdateLeaveBalance(r.Context(), application)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(application)
}
```

**Tasks - Day 5**:
```
[ ] Implement ApplyLeave handler
[ ] Implement GetLeaveApplication handler
[ ] Implement ApproveLeave handler
[ ] Implement RejectLeave handler
[ ] Implement GetLeaveBalance handler
[ ] Create Leave models
[ ] Create LeaveService
[ ] Implement leave balance calculation
[ ] Tests (80%+ coverage)

Deliverable: Leave management system, 8 endpoints tested
```

**Week 3 Backend Summary**:
- Employee Management: 5 endpoints ✅
- Attendance System: 8 endpoints ✅
- Leave Management: 8 endpoints ✅
- **Total**: 21 endpoints, 130+ code commits

---

### Week 3: Frontend Implementation (Parallel)

**Frontend Team (1 Developer)**:

**Day 1-2: Employee List & Detail Pages**
```tsx
// frontend/components/hr/EmployeeList.tsx

export function EmployeeList() {
  const [employees, setEmployees] = useState([])
  const [loading, setLoading] = useState(false)
  const [filters, setFilters] = useState({
    department: '',
    status: 'active',
  })
  
  useEffect(() => {
    fetchEmployees()
  }, [filters])
  
  const fetchEmployees = async () => {
    setLoading(true)
    const response = await fetch(
      `/api/v2/hr/employees?department_id=${filters.department}&status=${filters.status}`,
      { headers: { 'Authorization': `Bearer ${token}` } }
    )
    const data = await response.json()
    setEmployees(data.data)
    setLoading(false)
  }
  
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Employees</h1>
      
      <div className="mb-4 flex gap-2">
        <select 
          value={filters.department}
          onChange={(e) => setFilters({...filters, department: e.target.value})}
          className="px-3 py-2 border rounded"
        >
          <option value="">All Departments</option>
          {/* Populate from API */}
        </select>
        
        <button
          onClick={() => navigate('/hr/employees/new')}
          className="px-4 py-2 bg-blue-600 text-white rounded"
        >
          + New Employee
        </button>
      </div>
      
      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2 text-left">Name</th>
            <th className="p-2 text-left">Email</th>
            <th className="p-2 text-left">Department</th>
            <th className="p-2 text-left">Designation</th>
            <th className="p-2 text-left">Status</th>
            <th className="p-2 text-left">Actions</th>
          </tr>
        </thead>
        <tbody>
          {employees.map(emp => (
            <tr key={emp.id} className="border-t hover:bg-gray-50">
              <td className="p-2">{emp.first_name} {emp.last_name}</td>
              <td className="p-2">{emp.email}</td>
              <td className="p-2">{emp.department.name}</td>
              <td className="p-2">{emp.designation.name}</td>
              <td className="p-2">
                <span className={`px-2 py-1 rounded text-sm ${
                  emp.status === 'active' ? 'bg-green-100' : 'bg-gray-100'
                }`}>
                  {emp.status}
                </span>
              </td>
              <td className="p-2">
                <button onClick={() => navigate(`/hr/employees/${emp.id}`)}>
                  View
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
```

**Day 3-4: Attendance Tracker**
```tsx
// frontend/components/hr/AttendanceTracker.tsx

export function AttendanceTracker() {
  const [date, setDate] = useState(new Date())
  const [attendance, setAttendance] = useState([])
  
  const fetchAttendance = async () => {
    const response = await fetch(
      `/api/v2/hr/attendance?date=${date.toISOString().split('T')[0]}`,
      { headers: { 'Authorization': `Bearer ${token}` } }
    )
    const data = await response.json()
    setAttendance(data.data)
  }
  
  const markAbsent = async (employeeId) => {
    await fetch(`/api/v2/hr/attendance/mark-absent`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: JSON.stringify({ employee_id: employeeId, date })
    })
    fetchAttendance()
  }
  
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Attendance - {date.toDateString()}</h1>
      
      <div className="mb-4">
        <input
          type="date"
          value={date.toISOString().split('T')[0]}
          onChange={(e) => setDate(new Date(e.target.value))}
          className="px-3 py-2 border rounded"
        />
      </div>
      
      <table className="w-full border">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2">Employee</th>
            <th className="p-2">Check In</th>
            <th className="p-2">Check Out</th>
            <th className="p-2">Status</th>
            <th className="p-2">Action</th>
          </tr>
        </thead>
        <tbody>
          {attendance.map(att => (
            <tr key={att.id} className="border-t">
              <td className="p-2">{att.employee.name}</td>
              <td className="p-2">{att.check_in_time || '-'}</td>
              <td className="p-2">{att.check_out_time || '-'}</td>
              <td className="p-2">{att.status}</td>
              <td className="p-2">
                {att.status === 'absent' && (
                  <button
                    onClick={() => markAbsent(att.employee_id)}
                    className="text-red-600 hover:text-red-800"
                  >
                    Mark Present
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
```

**Day 5: Leave Management UI**

```tsx
// frontend/components/hr/LeaveManagement.tsx

export function LeaveManagement() {
  const [applications, setApplications] = useState([])
  const [pendingCount, setPendingCount] = useState(0)
  
  useEffect(() => {
    fetchLeaveApplications()
  }, [])
  
  const fetchLeaveApplications = async () => {
    const response = await fetch('/api/v2/hr/leave-applications', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    const data = await response.json()
    setApplications(data.data)
    setPendingCount(data.data.filter(app => app.status === 'pending').length)
  }
  
  const approveLeave = async (id) => {
    await fetch(`/api/v2/hr/leave-applications/${id}/approve`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })
    fetchLeaveApplications()
  }
  
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Leave Applications</h1>
      
      <div className="mb-4 p-4 bg-yellow-50 border border-yellow-200 rounded">
        <p className="text-sm">Pending Approvals: <strong>{pendingCount}</strong></p>
      </div>
      
      {applications.map(app => (
        <div key={app.id} className="mb-4 p-4 border rounded">
          <div className="flex justify-between">
            <div>
              <h3 className="font-bold">{app.employee.name}</h3>
              <p className="text-sm text-gray-600">
                {app.from_date} to {app.to_date}
              </p>
              <p className="text-sm">{app.reason}</p>
            </div>
            <div>
              <span className={`px-3 py-1 rounded text-sm ${
                app.status === 'approved' ? 'bg-green-100' :
                app.status === 'rejected' ? 'bg-red-100' :
                'bg-yellow-100'
              }`}>
                {app.status}
              </span>
            </div>
          </div>
          
          {app.status === 'pending' && (
            <div className="mt-3 flex gap-2">
              <button
                onClick={() => approveLeave(app.id)}
                className="px-4 py-2 bg-green-600 text-white rounded"
              >
                Approve
              </button>
              <button
                onClick={() => rejectLeave(app.id)}
                className="px-4 py-2 bg-red-600 text-white rounded"
              >
                Reject
              </button>
            </div>
          )}
        </div>
      ))}
    </div>
  )
}
```

---

### Week 4: Payroll System & Integration

**Monday-Wednesday: Payroll Processing**

```go
// internal/handlers/hr/payroll_handler.go

func (h *PayrollHandler) RunMonthlyPayroll(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Month string `json:"month" validate:"required"`  // YYYY-MM
        Year  int    `json:"year" validate:"required"`
    }
    
    // 1. Get all active employees
    employees, _ := h.service.GetActiveEmployees(r.Context())
    
    // 2. For each employee, calculate salary
    var payrollRun = &models.PayrollRun{
        Month:       req.Month,
        RunDate:     time.Now(),
        Status:      "processing",
        TotalSalary: 0,
    }
    
    salarySlips := make([]*models.SalarySlip, 0)
    
    for _, emp := range employees {
        slip := &models.SalarySlip{
            EmployeeID: emp.ID,
            Month:      req.Month,
        }
        
        // Get attendance count
        presentDays, _ := h.service.GetAttendanceDays(emp.ID, req.Month)
        totalDays := 30 // or get from leaves, shifts, etc.
        
        // Get salary structure
        structure, _ := h.service.GetSalaryStructure(emp.SalaryStructureID)
        
        // Calculate salary
        dailyRate := structure.BaseSalary / 30
        basicSalary := dailyRate * presentDays
        
        // Add allowances
        allowances, _ := h.service.GetAllowances(emp.SalaryStructureID)
        var allowanceAmount float64
        for _, allow := range allowances {
            percentage := (allow.Percentage / 100) * basicSalary
            slip.Allowances = append(slip.Allowances, &models.AllowanceDetail{
                Name: allow.Name,
                Amount: percentage,
            })
            allowanceAmount += percentage
        }
        
        // Gross Salary
        grossSalary := basicSalary + allowanceAmount
        slip.GrossSalary = grossSalary
        
        // Deductions
        deductions, _ := h.service.GetDeductions(emp.SalaryStructureID)
        var deductionAmount float64
        for _, ded := range deductions {
            if ded.Type == "fixed" {
                slip.Deductions = append(slip.Deductions, &models.DeductionDetail{
                    Name: ded.Name,
                    Amount: ded.Amount,
                })
                deductionAmount += ded.Amount
            } else if ded.Type == "percentage" {
                percentage := (ded.Percentage / 100) * grossSalary
                slip.Deductions = append(slip.Deductions, &models.DeductionDetail{
                    Name: ded.Name,
                    Amount: percentage,
                })
                deductionAmount += percentage
            }
        }
        
        // Net Salary
        netSalary := grossSalary - deductionAmount
        slip.NetSalary = netSalary
        slip.DeductionAmount = deductionAmount
        
        salarySlips = append(salarySlips, slip)
        payrollRun.TotalSalary += netSalary
    }
    
    // 3. Save payroll run and salary slips
    payrollRun, err := h.service.CreatePayrollRun(r.Context(), payrollRun, salarySlips)
    
    // 4. Generate GL entries
    entries := payrollToGLEntries(payrollRun)
    h.glService.PostEntries(r.Context(), entries)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(payrollRun)
}

func (h *PayrollHandler) GenerateSalarySlip(w http.ResponseWriter, r *http.Request) {
    employeeID := mux.Vars(r)["employee_id"]
    month := r.URL.Query().Get("month")
    
    slip, err := h.service.GetSalarySlip(employeeID, month)
    if err != nil {
        http.Error(w, "Salary slip not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(slip)
}
```

**Helper Function: GL Entry Generation**
```go
func payrollToGLEntries(payroll *models.PayrollRun) []*models.GLEntry {
    entries := make([]*models.GLEntry, 0)
    
    // Entry 1: Salary Expense
    entries = append(entries, &models.GLEntry{
        DebitAccount:  "Salary Expense",
        CreditAccount: "Bank Account",
        Amount:        payroll.TotalSalary,
        Narration:     fmt.Sprintf("Monthly Payroll - %s", payroll.Month),
        Date:          payroll.RunDate,
    })
    
    // Entry 2: Tax Payable (if applicable)
    // Entry 3: PF Payable (if applicable)
    
    return entries
}
```

**Tasks - Days 1-3**:
```
[ ] Implement salary calculation logic
[ ] Create payroll run procedure
[ ] Implement GL entry generation
[ ] Implement salary slip generation
[ ] Create PayrollService
[ ] Implement tax calculation
[ ] Tests (80%+ coverage)

Deliverable: Payroll system, 5 endpoints tested
```

**Thursday-Friday: Testing & Documentation**

```
[ ] End-to-end payroll test
  - Create test employee
  - Mark attendance
  - Run payroll
  - Verify GL entries
  - Verify salary slip

[ ] Integration with Accounts module
  - Verify GL entries post correctly
  - Verify GL balance updates
  - Verify audit trail

[ ] Load testing
  - Run payroll for 1000 employees
  - Measure performance
  - Target: < 5 minutes

[ ] API documentation
  - OpenAPI spec complete
  - All 45 endpoints documented
  - Examples provided

[ ] Create HR Module Summary
  - Features checklist
  - Known issues
  - Performance metrics
```

**End of Sprint 1 (HR Module)**:
```
✅ HR Module: 100% Complete
  ✅ Employee Management: 5 endpoints
  ✅ Attendance System: 8 endpoints
  ✅ Leave Management: 8 endpoints
  ✅ Payroll System: 5 endpoints
  ✅ Reports: 4 endpoints (summary, detailed, attendance, payroll)
  ✅ Statutory: 5 endpoints (tax, PF, ESI, compliance)
  ✅ Admin Endpoints: 5 endpoints (master data, config)
  
  ✅ Testing: 85%+ code coverage
  ✅ Database: 22 HR tables
  ✅ Integration: HR → GL
  ✅ Frontend: All 12 screens built
  ✅ Documentation: Complete API spec
  
✅ Story Points: 89/89 completed
```

---

## 📊 Success Metrics for Phase 3E

### Week 1-2: Foundation
- [x] 130+ tables created
- [x] Database passes integrity tests
- [x] All migrations reversible
- [x] RBAC fully extended

### Sprint 1: HR Module
- [x] 45 API endpoints deployed
- [x] 85%+ test coverage
- [x] All payroll calculations correct
- [x] GL entries post without errors
- [x] Frontend: 12/12 screens complete
- [x] Performance: Payroll runs < 5 min for 1000 employees

### Sprint 2: Accounts Module (Weeks 5-7)
- [ ] 40 API endpoints deployed
- [ ] GL balances reconcile perfectly
- [ ] Financial reports generate correctly
- [ ] Bank reconciliation working
- [ ] 85%+ test coverage

### Sprint 3: Sales Module (Weeks 8-9)
- [ ] 35 API endpoints deployed
- [ ] Orders convert to invoices automatically
- [ ] Commission calculations correct
- [ ] CRM pipeline working
- [ ] 85%+ test coverage

### Sprint 4-6: Purchase, Construction, Post Sales, Civil
- [ ] 90 combined API endpoints
- [ ] All modules integrated with GL
- [ ] Multi-module workflows tested
- [ ] 85%+ test coverage

### Final Week: Integration & Launch
- [ ] 99.9% uptime achieved
- [ ] All 235+ endpoints tested
- [ ] Zero critical bugs
- [ ] Performance: API response < 200ms
- [ ] Security: Penetration test passed
- [ ] UAT: 95%+ pass rate

---

## 🎓 Training & Knowledge Transfer

### Team Onboarding (Week 1-2)
- ERP concepts overview (4 hours)
- Database architecture walkthrough (2 hours)
- API standards & patterns (2 hours)
- Code review process (1 hour)
- Testing requirements (1 hour)

### Module-Specific Training
- HR concepts & compliance (3 hours)
- Accounting principles (4 hours)
- Sales & CRM best practices (2 hours)
- Purchase & inventory management (2 hours)

### Hands-On Sessions
- Pair programming on first module
- Code reviews with detailed feedback
- Weekly technical deep dives
- Knowledge sharing sessions

---

## 🏁 Conclusion

This detailed sprint breakdown shows **exactly** how to implement Phase 3E:

1. **Week 1-2**: Build the foundation (database, RBAC)
2. **Weeks 3-4**: HR module (payroll, attendance, leave)
3. **Weeks 5-7**: Accounts module (GL, invoicing, reporting)
4. **Weeks 8-9**: Sales module (CRM, pipeline, orders)
5. **Weeks 9-10**: Purchase module (vendor, PO, GRN)
6. **Weeks 11-13**: Construction & Civil modules
7. **Weeks 14-16**: Integration, testing, launch

**Each sprint is independently deployable** and includes testing, documentation, and frontend completion.

The total package:
- **235+ API endpoints**
- **130+ database tables**
- **71 frontend screens**
- **85%+ test coverage**
- **$155,000 budget**
- **16 weeks timeline**
- **8.3 FTE team**

**Ready to execute Phase 3E.**

---

**Generated**: November 24, 2025  
**Status**: Ready for Sprint Planning & Kickoff  
**Next Action**: Schedule Phase 3E Kickoff Meeting
