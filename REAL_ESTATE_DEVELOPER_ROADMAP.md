# 🏗️ VYOMTECH ERP - Real Estate Developer Roadmap
## Complete Implementation Plan for Missing Features

**Document**: Next Steps for Complete Real Estate Developer ERP  
**Generated**: December 23, 2025  
**Status**: Ready for Implementation

---

## 📋 Executive Summary

The VYOMTECH ERP has **core real estate functionality** in place but requires **12 major feature additions** to be a complete, production-ready Real Estate Developer ERP. This document provides a phased implementation roadmap with specific database schemas, API endpoints, and frontend components needed.

**Current Status**: ✅ 65% Complete  
**Target Completion**: Phase 3 in 6-8 weeks  
**Development Effort**: ~400-500 hours

---

## 🎯 Phase 1: CRITICAL FEATURES (1-2 weeks)
### Must-have for functional real estate operations

---

### 1️⃣ **BANK LOAN & FINANCING MANAGEMENT** 
**Priority**: 🔴 CRITICAL  
**Effort**: 40 hours  
**Impact**: High - Essential for tracking buyer financing

#### What's Needed:
```
✅ Database: property_bank_financing table (partially exists)
✅ Database: property_disbursement_schedule table
❌ Backend: Bank financing service & handlers
❌ Frontend: Financing dashboard & approval workflow
❌ API: 12 endpoints for financing operations
```

#### Database Schema Addition:
```sql
-- Already exists but needs enhancement
-- migrations/022_project_management_system.sql

-- Key tables:
property_bank_financing
  - id, tenant_id, project_id, unit_id, customer_id
  - bank_name, sanctioned_amount, sanctioned_date
  - noc_received, noc_date, noc_document_url
  - disbursement_status, collection_status
  - created_at, updated_at

property_disbursement_schedule
  - id, financing_id, milestone_id
  - expected_disbursement_date, actual_disbursement_date
  - cheque_number, neft_reference, amount
  - status (pending, processed, rejected)
```

#### API Endpoints Needed:
```
POST   /api/v1/projects/:projectId/financing          Create financing record
GET    /api/v1/projects/:projectId/financing          List financing
GET    /api/v1/financing/:financingId                 Get details
PUT    /api/v1/financing/:financingId                 Update financing
POST   /api/v1/financing/:financingId/disbursement    Schedule disbursement
GET    /api/v1/financing/:financingId/disbursements   Get schedule
PATCH  /api/v1/financing/:financingId/noc             Upload NOC documents
POST   /api/v1/financing/:financingId/collection      Record collection
GET    /api/v1/financing/reports/bank-wise            Bank-wise summary
GET    /api/v1/financing/reports/outstanding          Outstanding amount by customer
```

#### Frontend Components:
```
Pages:
  - /dashboard/financing - Financing dashboard
  - /dashboard/financing/new - Add new financing
  - /dashboard/financing/:id - Financing details & tracking
  - /dashboard/financing/noc - NOC upload & management

Components:
  - FinancingForm - Create/edit financing
  - DisbursementSchedule - View & manage disbursements
  - NOCDocumentUpload - Upload & track NOC
  - FinancingReports - Bank-wise & customer-wise reports
  - DisbursementTracker - Track disbursement status
```

#### Implementation Steps:
1. Create `internal/models/bank_financing.go`
2. Create `internal/services/bank_financing_service.go` (100+ methods)
3. Create `internal/handlers/bank_financing_handler.go`
4. Register routes in `pkg/router/router.go`
5. Create frontend page components
6. Add API service methods in frontend

**Completion Criteria**:
- ✅ Create 5 financing records with different statuses
- ✅ Test disbursement schedule generation
- ✅ Verify NOC document upload
- ✅ Validate bank-wise reporting

---

### 2️⃣ **BROKERAGE & COMMISSION MANAGEMENT**
**Priority**: 🔴 CRITICAL  
**Effort**: 35 hours  
**Impact**: High - Essential for agent management & commissions

#### What's Needed:
```
❌ Database: broker_profile, broker_commission tables
❌ Backend: Broker service & handlers
❌ Frontend: Broker management & commission tracking
❌ API: 14 endpoints for broker operations
```

#### Database Schema (New Migration):
```sql
-- New migration file: migrations/029_broker_management.sql

CREATE TABLE broker_profile (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  broker_name VARCHAR(255) NOT NULL,
  broker_code VARCHAR(50) UNIQUE NOT NULL,
  broker_type VARCHAR(50), -- individual, firm, corporate
  rera_number VARCHAR(100),
  pan_number VARCHAR(20),
  gst_number VARCHAR(20),
  email VARCHAR(100),
  phone VARCHAR(20),
  address TEXT,
  city VARCHAR(100),
  state VARCHAR(100),
  postal_code VARCHAR(20),
  bank_account_number VARCHAR(50),
  ifsc_code VARCHAR(20),
  is_active BOOLEAN DEFAULT true,
  commission_percentage DECIMAL(5,2),
  commission_type VARCHAR(50), -- percentage, flat, tiered
  payment_terms VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
  KEY idx_tenant (tenant_id),
  KEY idx_code (broker_code),
  UNIQUE KEY idx_rera (rera_number)
);

CREATE TABLE broker_commission_structure (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  broker_id VARCHAR(36) NOT NULL,
  project_id VARCHAR(36),
  commission_slab_start DECIMAL(18,2),
  commission_slab_end DECIMAL(18,2),
  commission_rate DECIMAL(5,2),
  effective_from DATE,
  effective_to DATE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
  FOREIGN KEY (broker_id) REFERENCES broker_profile(id) ON DELETE CASCADE,
  FOREIGN KEY (project_id) REFERENCES property_project(id) ON DELETE CASCADE,
  KEY idx_tenant (tenant_id),
  KEY idx_broker (broker_id)
);

CREATE TABLE broker_booking_link (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  broker_id VARCHAR(36) NOT NULL,
  commission_percentage DECIMAL(5,2),
  commission_amount DECIMAL(18,2),
  commission_due_date DATE,
  commission_paid_date DATE,
  payment_status VARCHAR(50), -- pending, partial, paid, cancelled
  payment_reference VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
  FOREIGN KEY (booking_id) REFERENCES property_booking(id) ON DELETE CASCADE,
  FOREIGN KEY (broker_id) REFERENCES broker_profile(id) ON DELETE CASCADE,
  KEY idx_tenant (tenant_id),
  KEY idx_booking (booking_id),
  KEY idx_broker (broker_id),
  KEY idx_status (payment_status)
);

CREATE TABLE broker_commission_payout (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  broker_id VARCHAR(36) NOT NULL,
  payout_period_start DATE,
  payout_period_end DATE,
  total_commission_amount DECIMAL(18,2),
  tds_percentage DECIMAL(5,2),
  tds_amount DECIMAL(18,2),
  net_payout DECIMAL(18,2),
  payout_date DATE,
  payment_mode VARCHAR(50), -- bank_transfer, cheque, neft, rtgs
  bank_reference VARCHAR(100),
  payment_status VARCHAR(50), -- pending, processed, failed
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
  FOREIGN KEY (broker_id) REFERENCES broker_profile(id) ON DELETE CASCADE,
  KEY idx_tenant (tenant_id),
  KEY idx_broker (broker_id),
  KEY idx_period (payout_period_start, payout_period_end)
);
```

#### API Endpoints:
```
-- Broker Profile Management
POST   /api/v1/brokers                      Create broker profile
GET    /api/v1/brokers                      List all brokers
GET    /api/v1/brokers/:brokerId            Get broker details
PUT    /api/v1/brokers/:brokerId            Update broker profile
DELETE /api/v1/brokers/:brokerId            Deactivate broker

-- Commission Structure
POST   /api/v1/brokers/:brokerId/commission-structure    Add commission slab
GET    /api/v1/brokers/:brokerId/commission-structure    Get slabs
PUT    /api/v1/brokers/:brokerId/commission-structure/:id Update slab

-- Commission Tracking
GET    /api/v1/brokers/:brokerId/commissions             Get all commissions
GET    /api/v1/brokers/:brokerId/payouts                 Get payout history
POST   /api/v1/brokers/:brokerId/payout                  Generate payout
GET    /api/v1/brokers/reports/top-performers            Top brokers report
GET    /api/v1/brokers/reports/commission-due            Outstanding commission
```

#### Frontend:
```
Pages:
  - /dashboard/brokers - Broker management
  - /dashboard/brokers/new - Add new broker
  - /dashboard/brokers/:id - Broker profile & performance
  - /dashboard/brokers/:id/commissions - Commission tracking
  - /dashboard/brokers/reports - Broker performance reports

Components:
  - BrokerForm - Create/edit broker
  - CommissionStructure - Manage commission slabs
  - CommissionTracker - View pending/paid commissions
  - BrokerReports - Performance analytics
  - PayoutGenerator - Create payout batches
```

#### Implementation Steps:
1. Create migration file `migrations/029_broker_management.sql`
2. Create models in `internal/models/broker.go`
3. Create `internal/services/broker_service.go`
4. Create `internal/handlers/broker_handler.go`
5. Register routes in router
6. Create frontend pages & components

---

### 3️⃣ **JOINT APPLICANTS & CO-OWNERS**
**Priority**: 🔴 CRITICAL  
**Effort**: 30 hours  
**Impact**: High - Most properties have multiple owners

#### What's Needed:
```
❌ Database: property_co_applicant table
❌ Enhanced: property_customer_profile (add co-applicant fields)
❌ Backend: Co-applicant service & validation
❌ Frontend: Co-applicant management UI
```

#### Database Schema:
```sql
-- New migration: migrations/030_joint_applicants.sql

CREATE TABLE property_co_applicant (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  applicant_number INT, -- 1 (primary), 2, 3, 4
  first_name VARCHAR(100) NOT NULL,
  middle_name VARCHAR(100),
  last_name VARCHAR(100),
  email VARCHAR(100),
  phone VARCHAR(20),
  date_of_birth DATE,
  pan_number VARCHAR(20),
  aadhar_number VARCHAR(20),
  passport_number VARCHAR(50),
  
  -- Address Details
  address_line_1 VARCHAR(255),
  address_line_2 VARCHAR(255),
  city VARCHAR(100),
  state VARCHAR(100),
  postal_code VARCHAR(20),
  country VARCHAR(100),
  
  -- Ownership Details
  ownership_percentage DECIMAL(5,2),
  is_nri BOOLEAN DEFAULT false,
  nri_country VARCHAR(100),
  
  -- Financial Details
  monthly_income DECIMAL(18,2),
  occupation VARCHAR(100),
  employer_name VARCHAR(255),
  
  -- Document URLs
  pan_document_url VARCHAR(500),
  aadhar_document_url VARCHAR(500),
  passport_document_url VARCHAR(500),
  
  -- Status
  kyc_status VARCHAR(50), -- pending, verified, rejected
  kyc_verified_by VARCHAR(36),
  kyc_verified_date TIMESTAMP,
  
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  
  FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
  FOREIGN KEY (booking_id) REFERENCES property_booking(id) ON DELETE CASCADE,
  KEY idx_tenant (tenant_id),
  KEY idx_booking (booking_id),
  UNIQUE KEY unique_applicant (booking_id, applicant_number)
);
```

#### API Endpoints:
```
POST   /api/v1/bookings/:bookingId/co-applicants           Add co-applicant
GET    /api/v1/bookings/:bookingId/co-applicants           List co-applicants
GET    /api/v1/co-applicants/:coApplicantId                Get details
PUT    /api/v1/co-applicants/:coApplicantId                Update co-applicant
DELETE /api/v1/co-applicants/:coApplicantId                Remove co-applicant
POST   /api/v1/co-applicants/:coApplicantId/kyc            Submit KYC
PATCH  /api/v1/co-applicants/:coApplicantId/kyc-verify     Verify KYC
POST   /api/v1/co-applicants/:coApplicantId/documents      Upload documents
GET    /api/v1/co-applicants/:coApplicantId/documents      List documents
```

---

## 🎯 Phase 2: IMPORTANT FEATURES (2-3 weeks)
### High-value features for business operations

---

### 4️⃣ **DOCUMENT MANAGEMENT SYSTEM**
**Priority**: 🟠 HIGH  
**Effort**: 45 hours  
**Impact**: High - Legal compliance & record keeping

#### Database Schema:
```sql
-- migrations/031_document_management.sql

CREATE TABLE document_template (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  document_type VARCHAR(100),
  document_name VARCHAR(255),
  file_url VARCHAR(500),
  file_size INT,
  version INT DEFAULT 1,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  KEY idx_type (document_type)
);

CREATE TABLE booking_document (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  document_type VARCHAR(100), -- agreement, poa, receipt, kyc, etc.
  document_name VARCHAR(255),
  document_url VARCHAR(500),
  file_size INT,
  uploaded_by VARCHAR(36),
  execution_date DATE,
  execution_status VARCHAR(50), -- pending, executed, rejected
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (booking_id) REFERENCES property_booking(id),
  KEY idx_booking (booking_id),
  KEY idx_type (document_type)
);

CREATE TABLE document_signature (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_document_id VARCHAR(36) NOT NULL,
  signer_id VARCHAR(36) NOT NULL,
  signer_role VARCHAR(50), -- customer, broker, company_rep
  signature_date TIMESTAMP,
  signature_image_url VARCHAR(500),
  e_signature_reference VARCHAR(255),
  status VARCHAR(50), -- pending, signed, rejected
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (booking_document_id) REFERENCES booking_document(id),
  KEY idx_booking_doc (booking_document_id)
);
```

#### API Endpoints:
```
-- Document Template Management
GET    /api/v1/document-templates                Document templates list
GET    /api/v1/document-templates/:templateId   Get template
POST   /api/v1/document-templates                Create template

-- Document Handling
POST   /api/v1/bookings/:bookingId/documents                 Upload document
GET    /api/v1/bookings/:bookingId/documents                 List documents
GET    /api/v1/bookings/:bookingId/documents/:documentId     Get document
DELETE /api/v1/bookings/:bookingId/documents/:documentId     Delete document

-- E-Signature
POST   /api/v1/booking-documents/:docId/request-signature    Request signature
POST   /api/v1/booking-documents/:docId/sign                 Sign document
GET    /api/v1/booking-documents/:docId/signature-status     Get status
```

---

### 5️⃣ **POSSESSION & HANDOVER TRACKING**
**Priority**: 🟠 HIGH  
**Effort**: 35 hours  
**Impact**: High - Post-sales operations

#### Database Schema:
```sql
-- migrations/032_possession_handover.sql

CREATE TABLE possession_milestone (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  milestone_type VARCHAR(50), -- oc_obtained, unit_ready, possession_offered, possession_taken
  milestone_name VARCHAR(255),
  scheduled_date DATE,
  actual_date DATE,
  status VARCHAR(50), -- pending, completed, postponed
  notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (booking_id) REFERENCES property_booking(id),
  KEY idx_booking (booking_id),
  KEY idx_type (milestone_type)
);

CREATE TABLE possession_handover (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  booking_id VARCHAR(36) NOT NULL,
  handover_date DATE,
  handover_status VARCHAR(50), -- offered, accepted, completed
  property_status VARCHAR(50), -- ready, defects, remediated
  
  -- Documents
  oc_certificate_url VARCHAR(500),
  occupancy_letter_url VARCHAR(500),
  possession_letter_url VARCHAR(500),
  
  -- Inspection
  inspection_date DATE,
  inspection_notes TEXT,
  defect_list JSON, -- {defects: [{area, description, status}]}
  
  -- Handover Details
  handed_over_by VARCHAR(36),
  handed_over_to VARCHAR(36),
  handover_sign_off_date DATE,
  
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (booking_id) REFERENCES property_booking(id),
  KEY idx_booking (booking_id)
);

CREATE TABLE possession_defect (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  handover_id VARCHAR(36) NOT NULL,
  defect_area VARCHAR(100),
  defect_description TEXT,
  severity VARCHAR(50), -- critical, major, minor
  status VARCHAR(50), -- identified, assigned, in_progress, resolved
  assigned_to VARCHAR(36),
  resolution_date DATE,
  resolution_notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (handover_id) REFERENCES possession_handover(id),
  KEY idx_handover (handover_id)
);
```

#### API Endpoints:
```
-- Milestones
POST   /api/v1/bookings/:bookingId/milestones              Create milestone
GET    /api/v1/bookings/:bookingId/milestones              List milestones
PUT    /api/v1/milestones/:milestoneId                     Update milestone

-- Handover
POST   /api/v1/bookings/:bookingId/handover                Create handover
GET    /api/v1/bookings/:bookingId/handover                Get handover details
PUT    /api/v1/handover/:handoverId                        Update handover
POST   /api/v1/handover/:handoverId/accept                 Accept possession

-- Defects
POST   /api/v1/handover/:handoverId/defects                Report defect
GET    /api/v1/handover/:handoverId/defects                List defects
PUT    /api/v1/defects/:defectId                           Update defect status

-- Reports
GET    /api/v1/projects/:projectId/handover-status         Handover status report
GET    /api/v1/projects/:projectId/defects-pending         Pending defects report
```

---

### 6️⃣ **TITLE CLEARANCE & LEGAL COMPLIANCE**
**Priority**: 🟠 HIGH  
**Effort**: 30 hours  
**Impact**: Medium - Legal compliance essential

#### Database Schema:
```sql
-- migrations/033_title_compliance.sql

CREATE TABLE title_clearance (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  project_id VARCHAR(36) NOT NULL,
  title_issue_type VARCHAR(100), -- encumbrance, objection, dispute, unclear
  issue_description TEXT,
  status VARCHAR(50), -- open, in_review, resolved, escalated
  clearance_date DATE,
  assigned_to VARCHAR(36),
  resolution_notes TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (project_id) REFERENCES property_project(id),
  KEY idx_project (project_id),
  KEY idx_status (status)
);

CREATE TABLE rera_compliance_tracker (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  project_id VARCHAR(36) NOT NULL,
  rera_registration_number VARCHAR(100),
  registration_date DATE,
  expiry_date DATE,
  status VARCHAR(50), -- active, expired, suspended, cancelled
  
  -- Required Documents
  promoter_id_proof VARCHAR(500),
  approvals_documents VARCHAR(500),
  site_plan_url VARCHAR(500),
  
  -- Compliance Status
  disclosure_requirements_met BOOLEAN,
  maintenance_account_opened BOOLEAN,
  insurance_policy_active BOOLEAN,
  
  last_audit_date DATE,
  audit_findings TEXT,
  
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (project_id) REFERENCES property_project(id),
  KEY idx_project (project_id)
);
```

---

### 7️⃣ **CUSTOMER PORTAL & SELF-SERVICE**
**Priority**: 🟠 HIGH  
**Effort**: 50 hours  
**Impact**: High - Customer engagement & reduced support load

#### Frontend Pages Needed:
```
/customer-portal/
  ├── dashboard
  │   ├── booking-status
  │   ├── payment-tracking
  │   ├── milestone-timeline
  │   └── documents
  ├── payments
  │   ├── payment-history
  │   ├── outstanding-dues
  │   └── payment-schedule
  ├── documents
  │   ├── my-documents
  │   ├── signatures-pending
  │   └── download-receipt
  ├── support
  │   ├── raise-complaint
  │   ├── ticket-history
  │   └── faq
  └── profile
      ├── my-details
      ├── co-applicants
      └── bank-details
```

---

## 🎯 Phase 3: ENHANCEMENT FEATURES (2-3 weeks)
### Nice-to-have features for advanced operations

---

### 8️⃣ **SITE VISIT & LEAD MANAGEMENT**
**Priority**: 🟡 MEDIUM  
**Effort**: 30 hours  

#### Database Schema:
```sql
-- migrations/034_site_visit_management.sql

CREATE TABLE site_visit_schedule (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  lead_id VARCHAR(36),
  visitor_name VARCHAR(100),
  visitor_phone VARCHAR(20),
  visitor_email VARCHAR(100),
  scheduled_date DATETIME,
  scheduled_by VARCHAR(36),
  status VARCHAR(50), -- scheduled, completed, cancelled, no_show
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  KEY idx_lead (lead_id),
  KEY idx_date (scheduled_date)
);

CREATE TABLE site_visit_log (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  visit_schedule_id VARCHAR(36),
  check_in_time DATETIME,
  check_out_time DATETIME,
  visited_by VARCHAR(36),
  units_viewed JSON, -- list of unit IDs viewed
  feedback TEXT,
  follow_up_required BOOLEAN,
  next_followup_date DATE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (visit_schedule_id) REFERENCES site_visit_schedule(id),
  KEY idx_schedule (visit_schedule_id)
);
```

---

### 9️⃣ **ADVANCED REPORTING & ANALYTICS**
**Priority**: 🟡 MEDIUM  
**Effort**: 40 hours  

#### Reports to Create:
```
1. Sales Dashboard
   - Monthly bookings & pipeline
   - Revenue realization
   - Sales team performance
   - Broker performance

2. Collections Report
   - Payment received vs due
   - Customer-wise outstanding
   - Milestone-wise collections
   - Aging analysis

3. Broker Performance
   - Top brokers by commission
   - Commission due & paid
   - Booking attribution
   - Commission timeline

4. Project Health Report
   - Units sold/booked/unsold
   - Occupancy rate
   - Revenue status
   - Cost vs budget

5. Financial Reports
   - GL posting status
   - Customer receivables
   - Broker payables
   - Tax reporting (TDS, GST)
```

---

### 🔟 **BULK OPERATIONS & BATCH PROCESSING**
**Priority**: 🟡 MEDIUM  
**Effort**: 25 hours  

#### Features:
```
- Bulk booking/allotment upload
- Bulk payment receipt generation
- Batch email/SMS campaigns
- Bulk milestone update
- Excel import/export
```

---

### 1️⃣1️⃣ **INVENTORY & AVAILABILITY MANAGEMENT**
**Priority**: 🟡 MEDIUM  
**Effort**: 30 hours  

#### Database Schema:
```sql
-- migrations/035_inventory_management.sql

CREATE TABLE unit_availability (
  id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  unit_id VARCHAR(36) NOT NULL,
  
  status VARCHAR(50), -- available, booked, sold, reserved, blocked
  status_reason VARCHAR(255),
  status_date TIMESTAMP,
  
  -- Availability Window
  available_from DATE,
  available_until DATE,
  
  -- Block Details
  blocked_by VARCHAR(36),
  block_reason VARCHAR(100),
  block_until DATE,
  
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenant(id),
  FOREIGN KEY (unit_id) REFERENCES property_unit(id),
  KEY idx_unit (unit_id),
  KEY idx_status (status)
);
```

---

### 1️⃣2️⃣ **TAX PLANNING & COMPLIANCE TOOLS**
**Priority**: 🟡 MEDIUM  
**Effort**: 35 hours  

#### Calculators:
```
- TDS Calculator (on rent, professional fees, etc.)
- Stamp Duty Calculator (by state/city)
- GST on Construction
- Income Tax Slab Calculator
- Investment Benefit Calculator
```

---

## 📊 IMPLEMENTATION CHECKLIST

### Phase 1 (Weeks 1-2) - CRITICAL
```
Week 1:
[ ] Bank Financing Module
  [ ] Database migration
  [ ] Models & services
  [ ] API endpoints
  [ ] Basic UI
  
[ ] Brokerage Management
  [ ] Database migration
  [ ] Broker profile & commissions
  [ ] Commission structure setup
  [ ] Basic UI

Week 2:
[ ] Joint Applicants
  [ ] Database schema
  [ ] Co-applicant service
  [ ] KYC verification
  [ ] UI forms

[ ] Document Management (Phase 1)
  [ ] Document upload system
  [ ] Template management
  [ ] Basic document tracking
```

### Phase 2 (Weeks 3-5) - HIGH PRIORITY
```
Week 3:
[ ] Possession & Handover
  [ ] Database schema
  [ ] Milestone tracking service
  [ ] Handover workflow
  [ ] UI dashboard

[ ] Title Clearance
  [ ] Database schema
  [ ] Compliance tracker
  [ ] Issue tracking service

Week 4:
[ ] Document Management (Phase 2)
  [ ] E-signature integration
  [ ] Document versioning
  [ ] Digital archive
  [ ] Advanced permissions

Week 5:
[ ] Customer Portal (MVP)
  [ ] Booking status dashboard
  [ ] Payment tracker
  [ ] Document download
  [ ] Basic support ticket
```

### Phase 3 (Weeks 6-8) - ENHANCEMENT
```
Week 6:
[ ] Site Visit Management
[ ] Inventory Management
[ ] Advanced Reporting

Week 7:
[ ] Bulk Operations
[ ] Tax Calculators
[ ] API documentation

Week 8:
[ ] Integration testing
[ ] Performance optimization
[ ] User documentation
[ ] Go-live preparation
```

---

## 🛠️ TECHNICAL DETAILS

### Database Migrations to Create:
```
029_broker_management.sql           (4 tables)
030_joint_applicants.sql             (1 table)
031_document_management.sql           (3 tables)
032_possession_handover.sql           (3 tables)
033_title_compliance.sql              (2 tables)
034_site_visit_management.sql         (2 tables)
035_inventory_management.sql          (1 table)
```

### Backend Services to Create:
```
BankFinancingService      (15 methods)
BrokerService             (18 methods)
CoApplicantService        (10 methods)
DocumentService           (12 methods)
PossessionService         (14 methods)
ComplianceService         (8 methods)
SiteVisitService          (10 methods)
ReportingService          (20+ methods)
```

### API Endpoints to Implement:
```
Total: 120+ endpoints across all features
- Bank Financing: 8 endpoints
- Brokerage: 12 endpoints
- Co-Applicants: 8 endpoints
- Documents: 10 endpoints
- Possession: 15 endpoints
- Compliance: 8 endpoints
- Site Visits: 8 endpoints
- Reports: 20+ endpoints
- Customer Portal: 15 endpoints
```

### Frontend Pages to Create:
```
- 15+ new dashboard pages
- 25+ components
- 30+ API hooks
- Customer portal (8 pages)
- Admin panels for all features
```

---

## 🎯 SUCCESS CRITERIA

**Phase 1 Complete When**:
- ✅ Bank financing can be created, tracked, and disbursed
- ✅ Broker profiles and commissions fully operational
- ✅ Co-applicants added with KYC verification
- ✅ Documents uploaded and stored with audit trail

**Phase 2 Complete When**:
- ✅ Possession workflow from OC to handover tracked
- ✅ Title clearance issues managed
- ✅ Customer portal live with payment tracking
- ✅ Document e-signature integrated

**Phase 3 Complete When**:
- ✅ Advanced reports generated automatically
- ✅ Bulk operations reduce manual effort by 70%
- ✅ Tax compliance fully automated
- ✅ System handles 100+ concurrent users

---

## 📈 EFFORT ESTIMATION SUMMARY

| Phase | Feature Count | Backend | Frontend | Database | Total Hours |
|-------|--------------|---------|----------|----------|-------------|
| Phase 1 | 4 | 100 | 80 | 20 | **200 hrs** |
| Phase 2 | 3 | 85 | 75 | 15 | **175 hrs** |
| Phase 3 | 5 | 70 | 60 | 10 | **140 hrs** |
| **TOTAL** | **12** | **255** | **215** | **45** | **515 hrs** |

**Timeline**: 6-8 weeks (assuming 40 hrs/week development)  
**Team Size**: 2-3 developers recommended

---

## 🚀 NEXT IMMEDIATE ACTIONS

1. **Create migration files** for Phase 1 features
2. **Set up database** with new schemas
3. **Start Bank Financing module** (highest priority)
4. **Parallel track**: Brokerage & Co-applicants
5. **Weekly reviews** of progress against roadmap

---

## 📝 REVISION HISTORY

| Date | Version | Status | Notes |
|------|---------|--------|-------|
| 2025-12-23 | 1.0 | Final | Real Estate Developer ERP roadmap created |

---

**Document Prepared By**: AI Development Assistant  
**For**: VYOMTECH ERP Real Estate Module  
**Status**: Ready for Development Team Implementation

---

## 📚 Related Documents

- [PROJECT_SPEEDUP_GUIDE.md](PROJECT_SPEEDUP_GUIDE.md) - System Architecture
- [PRODUCTION_READINESS_REPORT.md](PRODUCTION_READINESS_REPORT.md) - Current Status
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Development Quick Start
- [SYSTEM_AUDIT_CHECKLIST.md](SYSTEM_AUDIT_CHECKLIST.md) - Verification Guide
