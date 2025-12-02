# 🎯 MISSING COMPONENTS - Implementation Roadmap

## 📊 Quick Status

| Category | Status | Count |
|----------|--------|-------|
| Frontend Dashboards | ✅ DONE | 4 new |
| Frontend Types | ✅ DONE | 4 new |
| Frontend Components | ❌ TODO | 4 folders |
| Backend Models | ❌ TODO | 2 files |
| Backend Handlers | ❌ TODO | 2 files |
| Backend Services | ❌ TODO | 4 files |
| Database Migrations | ❌ TODO | 4 files |
| API Routes | ❌ TODO | Integration |

---

## 🔴 CRITICAL PATH - Frontend Components Missing

### 1. Civil Engineering Components
**Location**: `frontend/components/modules/Civil/`

**Needed Files**:
```typescript
// SiteForm.tsx - Create/Edit sites
interface SiteFormProps {
  site?: Site
  onSubmit: (data: Partial<Site>) => Promise<void>
  onCancel: () => void
}

// SiteList.tsx - List all sites
interface SiteListProps {
  sites: Site[]
  loading: boolean
  onEdit: (site: Site) => void
  onDelete: (site: Site) => void
}

// IncidentForm.tsx - Report incidents
// IncidentList.tsx - View incidents
// ComplianceTracking.tsx - Track compliance
// PermitManagement.tsx - Manage permits
```

### 2. Construction Components
**Location**: `frontend/components/modules/Construction/`

**Needed Files**:
```typescript
// ProjectForm.tsx
// ProjectList.tsx
// BOQForm.tsx - Bill of Quantities
// BOQList.tsx
// ProgressForm.tsx
// ProgressList.tsx
// QualityControlForm.tsx
// QualityControlList.tsx
```

### 3. Gamification Components
**Location**: `frontend/components/modules/Gamification/`

**Needed Files**:
```typescript
// LeaderboardView.tsx
// BadgeGallery.tsx
// ChallengeCard.tsx
// AchievementTracker.tsx
// RewardsShop.tsx (Optional)
```

### 4. Scheduled Tasks Components
**Location**: `frontend/components/modules/ScheduledTasks/`

**Needed Files**:
```typescript
// TaskForm.tsx
// TaskList.tsx
// ExecutionHistoryView.tsx
// TemplateManager.tsx
// TaskNotifications.tsx
```

---

## 🔴 CRITICAL PATH - Backend Models Missing

### 1. `internal/models/civil.go`

**Structure**:
```go
package models

import "time"

// Site represents a construction site
type Site struct {
    ID              int64     `db:"id"`
    TenantID        string    `db:"tenant_id"`
    SiteName        string    `db:"site_name"`
    Location        string    `db:"location"`
    ProjectID       string    `db:"project_id"`
    SiteManager     string    `db:"site_manager"`
    StartDate       time.Time `db:"start_date"`
    ExpectedEndDate time.Time `db:"expected_end_date"`
    CurrentStatus   string    `db:"current_status"` // planning, active, paused, completed, closed
    SiteAreaSqm     float64   `db:"site_area_sqm"`
    WorkforceCount  int       `db:"workforce_count"`
    CreatedAt       time.Time `db:"created_at"`
    UpdatedAt       time.Time `db:"updated_at"`
}

type SafetyIncident struct {
    ID              int64     `db:"id"`
    SiteID          string    `db:"site_id"`
    IncidentType    string    `db:"incident_type"` // accident, near_miss, hazard, violation
    Severity        string    `db:"severity"` // low, medium, high, critical
    IncidentDate    time.Time `db:"incident_date"`
    Description     string    `db:"description"`
    ReportedBy      string    `db:"reported_by"`
    Status          string    `db:"status"` // open, investigating, resolved, closed
    IncidentNumber  string    `db:"incident_number"`
    CreatedAt       time.Time `db:"created_at"`
}

type Compliance struct {
    ID              int64     `db:"id"`
    SiteID          string    `db:"site_id"`
    ComplianceType  string    `db:"compliance_type"` // safety, environmental, labor, regulatory
    Requirement     string    `db:"requirement"`
    DueDate         time.Time `db:"due_date"`
    Status          string    `db:"status"` // compliant, non_compliant, in_progress, not_applicable
    LastAuditDate   time.Time `db:"last_audit_date"`
    AuditResult     string    `db:"audit_result"` // pass, fail, pending
    Notes           string    `db:"notes"`
    CreatedAt       time.Time `db:"created_at"`
}

type Permit struct {
    ID               int64     `db:"id"`
    SiteID           string    `db:"site_id"`
    PermitType       string    `db:"permit_type"`
    PermitNumber     string    `db:"permit_number"`
    IssuedDate       time.Time `db:"issued_date"`
    ExpiryDate       time.Time `db:"expiry_date"`
    IssuingAuthority string    `db:"issuing_authority"`
    Status           string    `db:"status"` // active, expired, cancelled, pending
    DocumentURL      string    `db:"document_url"`
    CreatedAt        time.Time `db:"created_at"`
}
```

### 2. `internal/models/construction.go`

**Structure**:
```go
package models

import "time"

type ConstructionProject struct {
    ID                     int64     `db:"id"`
    TenantID               string    `db:"tenant_id"`
    ProjectName            string    `db:"project_name"`
    ProjectCode            string    `db:"project_code"`
    Location               string    `db:"location"`
    Client                 string    `db:"client"`
    ContractValue          float64   `db:"contract_value"`
    StartDate              time.Time `db:"start_date"`
    ExpectedCompletion     time.Time `db:"expected_completion"`
    CurrentProgressPercent int       `db:"current_progress_percentage"`
    Status                 string    `db:"status"` // planning, active, suspended, completed, on_hold
    ProjectManager         string    `db:"project_manager"`
    CreatedAt              time.Time `db:"created_at"`
}

type BillOfQuantities struct {
    ID              int64     `db:"id"`
    ProjectID       string    `db:"project_id"`
    BOQNumber       string    `db:"boq_number"`
    ItemDescription string    `db:"item_description"`
    Unit            string    `db:"unit"`
    Quantity        float64   `db:"quantity"`
    UnitRate        float64   `db:"unit_rate"`
    TotalAmount     float64   `db:"total_amount"`
    Category        string    `db:"category"` // civil, structural, electrical, plumbing, finishing, other
    Status          string    `db:"status"` // planned, in_progress, completed, on_hold
    CreatedAt       time.Time `db:"created_at"`
}

type ProgressTracking struct {
    ID                int64     `db:"id"`
    ProjectID         string    `db:"project_id"`
    Date              time.Time `db:"date"`
    ActivityDesc      string    `db:"activity_description"`
    QuantityCompleted float64   `db:"quantity_completed"`
    Unit              string    `db:"unit"`
    PercentComplete   int       `db:"percentage_complete"`
    WorkforceDeployed int       `db:"workforce_deployed"`
    Notes             string    `db:"notes"`
    PhotoURL          string    `db:"photo_url"`
    CreatedAt         time.Time `db:"created_at"`
}

type QualityControl struct {
    ID               int64     `db:"id"`
    ProjectID        string    `db:"project_id"`
    BOQItemID        string    `db:"boq_item_id"`
    InspectionDate   time.Time `db:"inspection_date"`
    InspectorName    string    `db:"inspector_name"`
    QualityStatus    string    `db:"quality_status"` // passed, failed, partial, pending
    Observations     string    `db:"observations"`
    CorrectiveAction string    `db:"corrective_actions"`
    FollowUpDate     time.Time `db:"follow_up_date"`
    CreatedAt        time.Time `db:"created_at"`
}
```

---

## 🔴 CRITICAL PATH - Backend Handlers Missing

### 1. `internal/handlers/civil_handler.go`

**Endpoint Structure** (Reference: purchase_handler.go):
```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type CivilHandler struct {
    service *services.CivilService
}

// GET /api/v1/civil/sites
func (h *CivilHandler) GetSites(c *gin.Context) {
    // List all sites for tenant
    // Include filtering, pagination
}

// POST /api/v1/civil/sites
func (h *CivilHandler) CreateSite(c *gin.Context) {
    // Create new site
    // Validate input, check permissions
}

// PUT /api/v1/civil/sites/:id
func (h *CivilHandler) UpdateSite(c *gin.Context) {
    // Update site
}

// DELETE /api/v1/civil/sites/:id
func (h *CivilHandler) DeleteSite(c *gin.Context) {
    // Delete site
}

// POST /api/v1/civil/incidents
func (h *CivilHandler) ReportIncident(c *gin.Context) {
    // Create safety incident
}

// GET /api/v1/civil/incidents/:siteId
func (h *CivilHandler) GetIncidents(c *gin.Context) {
    // List incidents for site
}

// Similar endpoints for:
// - Compliance tracking
// - Permit management
// - Dashboard metrics
```

### 2. `internal/handlers/construction_handler.go`

**Similar structure with endpoints**:
- GetProjects, CreateProject, UpdateProject, DeleteProject
- GetBOQ, CreateBOQItem, UpdateBOQItem
- LogProgress, GetProgress
- QualityInspection, GetQC

---

## 🔴 CRITICAL PATH - Service Layer Missing

### 1. `internal/services/civil_service.go`

**Interface**:
```go
package services

type CivilService struct {
    db *sql.DB
}

func NewCivilService(db *sql.DB) *CivilService {
    return &CivilService{db: db}
}

// Site operations
func (s *CivilService) GetSites(ctx context.Context, tenantID string) ([]models.Site, error)
func (s *CivilService) CreateSite(ctx context.Context, site *models.Site) error
func (s *CivilService) UpdateSite(ctx context.Context, site *models.Site) error
func (s *CivilService) DeleteSite(ctx context.Context, siteID string) error

// Incident operations
func (s *CivilService) ReportIncident(ctx context.Context, incident *models.SafetyIncident) error
func (s *CivilService) GetIncidents(ctx context.Context, siteID string) ([]models.SafetyIncident, error)

// Compliance operations
func (s *CivilService) GetCompliance(ctx context.Context, siteID string) ([]models.Compliance, error)
func (s *CivilService) UpdateComplianceStatus(ctx context.Context, complianceID, status string) error

// Permit operations
func (s *CivilService) GetPermits(ctx context.Context, siteID string) ([]models.Permit, error)
```

### 2. `internal/services/construction_service.go`

Similar structure for construction operations

---

## 🔴 CRITICAL PATH - Database Migrations Missing

### 1. `internal/migrations/002_civil_schema.sql`

```sql
-- Civil Engineering Module Schema

CREATE TABLE IF NOT EXISTS sites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    site_name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    project_id VARCHAR(36),
    site_manager VARCHAR(255),
    start_date DATE,
    expected_end_date DATE,
    current_status VARCHAR(50) DEFAULT 'active',
    site_area_sqm DECIMAL(12,2),
    workforce_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS safety_incidents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    site_id BIGINT,
    incident_type VARCHAR(50),
    severity VARCHAR(50),
    incident_date DATETIME,
    description TEXT,
    reported_by VARCHAR(255),
    status VARCHAR(50) DEFAULT 'open',
    incident_number VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    KEY idx_site (site_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS compliance_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    site_id BIGINT,
    compliance_type VARCHAR(50),
    requirement VARCHAR(255),
    due_date DATE,
    status VARCHAR(50),
    last_audit_date DATE,
    audit_result VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS permits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    site_id BIGINT,
    permit_type VARCHAR(100),
    permit_number VARCHAR(100) UNIQUE,
    issued_date DATE,
    expiry_date DATE,
    issuing_authority VARCHAR(255),
    status VARCHAR(50),
    document_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 2. `internal/migrations/003_construction_schema.sql`

```sql
-- Construction Module Schema

CREATE TABLE IF NOT EXISTS construction_projects (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    project_name VARCHAR(255),
    project_code VARCHAR(100) UNIQUE,
    location VARCHAR(255),
    client VARCHAR(255),
    contract_value DECIMAL(15,2),
    start_date DATE,
    expected_completion DATE,
    current_progress_percentage INT DEFAULT 0,
    status VARCHAR(50) DEFAULT 'planning',
    project_manager VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS bill_of_quantities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    project_id BIGINT,
    boq_number VARCHAR(100),
    item_description VARCHAR(500),
    unit VARCHAR(50),
    quantity DECIMAL(12,2),
    unit_rate DECIMAL(12,2),
    total_amount DECIMAL(15,2),
    category VARCHAR(50),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    KEY idx_project (project_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS progress_tracking (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    project_id BIGINT,
    date DATE,
    activity_description TEXT,
    quantity_completed DECIMAL(12,2),
    unit VARCHAR(50),
    percentage_complete INT,
    workforce_deployed INT,
    notes TEXT,
    photo_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS quality_control (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    project_id BIGINT,
    boq_item_id BIGINT,
    inspection_date DATE,
    inspector_name VARCHAR(255),
    quality_status VARCHAR(50),
    observations TEXT,
    corrective_actions TEXT,
    follow_up_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY idx_tenant (tenant_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 📋 Implementation Checklist

### Phase 1: Backend Models & Services (4 Hours)
- [ ] Create `civil.go` (50 lines)
- [ ] Create `construction.go` (50 lines)
- [ ] Create `civil_service.go` (150 lines)
- [ ] Create `construction_service.go` (150 lines)

### Phase 2: Backend Handlers (4 Hours)
- [ ] Create `civil_handler.go` (200 lines)
- [ ] Create `construction_handler.go` (200 lines)
- [ ] Register routes in `cmd/main.go`

### Phase 3: Database Migrations (2 Hours)
- [ ] Create `002_civil_schema.sql`
- [ ] Create `003_construction_schema.sql`
- [ ] Run migrations

### Phase 4: Frontend Components (6 Hours)
- [ ] Create Civil/ folder with 6 components
- [ ] Create Construction/ folder with 8 components
- [ ] Create Gamification/ folder with 5 components
- [ ] Create ScheduledTasks/ folder with 5 components

### Phase 5: Integration (4 Hours)
- [ ] Connect frontend to backend APIs
- [ ] Replace mock data with real data
- [ ] Test all CRUD operations
- [ ] Verify multi-tenant isolation

### Phase 6: Testing (4 Hours)
- [ ] Unit tests for services
- [ ] Integration tests for handlers
- [ ] E2E tests for workflows
- [ ] Performance testing

---

## ⏱️ Total Estimated Time: 24 Hours

**Breakdown**:
- Backend: 14 hours
- Frontend: 6 hours
- Testing & Integration: 4 hours

---

**Priority**: CRITICAL
**Next Action**: Start with backend models & services
**Recommended Approach**: Follow Phase2B patterns from archive
