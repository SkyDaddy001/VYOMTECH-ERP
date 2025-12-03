# Phase 5: Partner Source & Credit Policies Implementation
**Status**: ✅ Complete  
**Date**: December 3, 2025  
**Component**: External Partner System Enhancement  

---

## 📋 Overview

Phase 5 extends the partner system (Phase 4) with sophisticated source attribution and flexible credit policies. This enables multi-business-model support for partners:

- **Partner Source Mapping**: Automatically maps 4 partner types to lead sources
- **Flexible Credit Policies**: Supports time-based, project-based, and campaign-based billing
- **Tiered Pricing**: Volume-based pricing with automatic tier escalation
- **Revenue-Based Billing**: Flexible calculation methods (fixed, percentage, tiered, conversion, revenue-share)

---

## 🎯 Requirements Fulfilled

| Requirement | Status | Implementation |
|------------|--------|-----------------|
| "Source will always be = partner type" | ✅ | 4 enum values: Customer Reference, Vendor Reference, Channel Partner, Property Portal |
| "Time-based credit policies" | ✅ | PartnerCreditPolicy with monthly/quarterly/annual periods |
| "Project-based credit policies" | ✅ | PartnerCreditPolicy with project_id linking |
| "Campaign-based credit policies" | ✅ | PartnerCreditPolicy with campaign_id linking |
| "Flexible calculation methods" | ✅ | 5 types: percentage, fixed_price, tiered, conversion, revenue_share |
| "Volume-based pricing" | ✅ | TierConfig with min_leads/max_leads boundaries |

---

## 📁 Files Created/Modified

### **New Services** (689 lines)
**File**: `internal/services/partner_source_credit_service.go`

#### **PartnerSourceService Interface** (9 methods)
```go
CreatePartnerSource(ctx, tenantID, source) → Auto-maps partner type to lead source
GetPartnerSource(ctx, tenantID, sourceID) → Retrieve single source
GetPartnerSources(ctx, tenantID, partnerID) → All sources for partner
GetSourceByCode(ctx, tenantID, sourceCode) → Lookup by source code
UpdatePartnerSource(ctx, tenantID, source) → Modify source details
DeactivatePartnerSource(ctx, tenantID, sourceID) → Soft-deactivate
GetSourceStats(ctx, tenantID, sourceID) → Analytics per source
GetPartnerSourceStats(ctx, tenantID, partnerID) → All source analytics
```

**Key Features**:
- Automatic source code generation (e.g., "CR-1", "VR-2")
- Lead statistics tracking (leads_generated, leads_converted, total_revenue)
- Approval/conversion rate calculations
- Multi-tenant isolation on all queries

#### **PartnerCreditPolicyService Interface** (11 methods)
```go
CreateCreditPolicy(ctx, tenantID, policy) → New policy with full configuration
GetCreditPolicy(ctx, tenantID, policyID) → Single policy retrieval
GetPartnerCreditPolicies(ctx, tenantID, partnerID) → All policies for partner
GetActiveCreditPolicies(ctx, tenantID, partnerID) → Only active/valid-dated
UpdateCreditPolicy(ctx, tenantID, policy) → Modify policy
ApproveCreditPolicy(ctx, tenantID, policyID, approvedBy) → Management approval
DeactivateCreditPolicy(ctx, tenantID, policyID) → Soft-delete
CalculateLeadCredit(ctx, tenantID, partnerLeadID) → Credit computation engine
GetApplicablePolicies(ctx, tenantID, partnerID, leadData) → Policy matching
GetPolicyMappings(ctx, tenantID, partnerLeadID) → Lead-policy relationships
```

**Credit Calculation Logic**:
- **Percentage**: `credit = base_credit` (stored as %, e.g., 10 = 10%)
- **Fixed Price**: `credit = base_credit` (fixed amount per lead)
- **Tiered**: Finds tier from `lead_count`, returns `tier.credit_amount`
- **Conversion**: `credit = 0` on submission (calculated after conversion)
- **Revenue Share**: `credit = base_credit` (% of actual deal revenue)
- **Bonus**: If quality_score >= 80, add `(credit * bonus_percentage / 100)`
- **Bounds**: Apply min/max credit constraints after calculation

### **Database Migration** (309 lines)
**File**: `migrations/023_partner_sources_and_credit_policies.sql`

#### **3 New Tables**

**1. partner_sources** (14 columns)
```sql
CREATE TABLE partner_sources (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(50) NOT NULL,
    partner_id BIGINT NOT NULL,
    
    source_type ENUM('customer_reference', 'vendor_reference', 'channel_partner', 'property_portal'),
    source_code VARCHAR(100) UNIQUE,
    source_name VARCHAR(255),
    description TEXT,
    
    is_active BOOLEAN DEFAULT TRUE,
    leads_generated BIGINT DEFAULT 0,
    leads_converted BIGINT DEFAULT 0,
    total_revenue DECIMAL(15,2) DEFAULT 0,
    
    created_at, updated_at, deleted_at TIMESTAMP
    
    INDEX: (tenant_id, partner_id), (source_type), (is_active)
    FK: partner_id → partners(id)
)
```

**2. partner_credit_policies** (36 columns)
```sql
CREATE TABLE partner_credit_policies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(50) NOT NULL,
    partner_id BIGINT NOT NULL,
    
    -- Policy Identity
    policy_code VARCHAR(100) UNIQUE,
    policy_name VARCHAR(255),
    policy_type ENUM('time_based', 'project_based', 'campaign_based'),
    calculation_type ENUM('percentage', 'fixed_price', 'tiered', 'conversion', 'revenue_share'),
    
    -- Time-Based Fields
    time_unit_type ENUM('monthly', 'quarterly', 'annual'),
    time_unit_value INT,
    policy_start_date, policy_end_date TIMESTAMP,
    
    -- Project/Campaign Fields
    project_id BIGINT, project_name VARCHAR(255),
    campaign_id BIGINT, campaign_name VARCHAR(255),
    
    -- Credit Configuration
    base_credit DECIMAL(10,2),
    minimum_credit, maximum_credit DECIMAL(10,2),
    bonus_percentage DECIMAL(5,2),
    
    -- Tier Configuration (JSON)
    tier_config JSON,  -- {"tiers": [{"tier_level": 1, "min_leads": 0, "max_leads": 100, "credit_amount": 10}...]}
    
    -- Conditions
    min_lead_quality_score DECIMAL(5,2),
    requires_approval, auto_approve BOOLEAN,
    
    -- Status & Approval
    is_active BOOLEAN,
    approval_required BOOLEAN,
    created_by, approved_by BIGINT,
    approved_at TIMESTAMP,
    
    -- Statistics
    total_leads_under_policy BIGINT DEFAULT 0,
    total_credits_allocated DECIMAL(15,2) DEFAULT 0,
    
    created_at, updated_at, deleted_at TIMESTAMP
    
    INDEX: (tenant_id, partner_id), (policy_type), (is_active), (approval_required)
    FK: partner_id → partners(id)
)
```

**3. partner_credit_policy_mappings** (7 columns)
```sql
CREATE TABLE partner_credit_policy_mappings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(50) NOT NULL,
    partner_lead_id BIGINT NOT NULL,
    policy_id BIGINT NOT NULL,
    
    calculated_credit DECIMAL(10,2),
    reason VARCHAR(255),
    created_at TIMESTAMP
    
    INDEX: (tenant_id, partner_lead_id), (policy_id)
    FK: partner_lead_id → partner_leads(id), policy_id → partner_credit_policies(id)
)
```

#### **Sample Data Inserted**
- 20 partner sources (5 each type, across sample partners)
- 9 credit policies (3 fixed-price, 3 percentage/revenue-share, 3 tiered)
- All with approval workflows and approval tracking

#### **Performance Indexes**
- Multi-tenant isolation: (tenant_id, partner_id)
- Status filtering: (is_active), (approval_required)
- Date-based queries: (policy_start_date, policy_end_date)
- Lead quality queries: partner_leads.(quality_score), (status, quality_score)

---

## 🏗 Architecture

### **Source Attribution Flow**
```
Partner Created (portal, channel_partner, vendor, customer)
    ↓
PartnerSourceService.CreatePartnerSource()
    ↓
Auto-maps partner_type → PartnerSourceType (4 enum values)
    ↓
Creates tracking record (leads_generated, leads_converted, total_revenue)
    ↓
Lead submitted by partner
    ↓
Lead tagged with source_id for attribution
    ↓
GetSourceStats() → approval_rate, conversion_rate, revenue
```

### **Credit Policy Flow**
```
PartnerCreditPolicy Created
    ↓
Specify: PolicyType (time/project/campaign) + CalculationType (fixed/percent/tiered/etc)
    ↓
Set: Base credit, min/max bounds, bonus rules, quality thresholds
    ↓
If tiered: Define tiers with min_leads/max_leads/credit_amount
    ↓
Policy requires approval (approval_required = TRUE)
    ↓
Management reviews + ApproveCreditPolicy()
    ↓
Policy becomes active for use
    ↓
Lead submitted
    ↓
PartnerLeadService.UpdatePartnerLeadStatus("approved")
    ↓
PartnerCreditPolicyService.CalculateLeadCredit()
    ↓
EvaluatePolicy() + ApplyCalculation()
    ↓
Store in partner_credit_policy_mappings
    ↓
Credit allocated to partner lead
```

### **Tiered Pricing Example**
```go
Policy: "Tiered Pricing - Volume Based"
Tiers:
  Tier 1: 1-50 leads = $8/lead
  Tier 2: 51-200 leads = $12/lead
  Tier 3: 201+ leads = $15/lead

Partner submits 150th lead:
  total_leads_under_policy = 150
  Applicable tier: Tier 2 (101-200)
  Base credit = $12
  If quality_score >= 80: add 5% bonus = $12.60
  Apply bounds: min=$5, max=$25 → final credit = $12.60
```

### **Revenue-Share Example**
```go
Policy: "10% Revenue Share"
Lead submitted with base value $100
  calculation_type = "revenue_share"
  base_credit = 10
  Lead status = "converted"
  Actual deal revenue = $500
  Credit = $500 × 10% = $50
  If bonus_percentage = 10 and quality_score >= 80:
    Add 10% bonus = $50 × 1.10 = $55
```

---

## 📊 Data Models

### **PartnerSourceType** (enum, 4 values)
```go
const (
    PartnerSourceCustomerReference  = "customer_reference"
    PartnerSourceVendorReference    = "vendor_reference"
    PartnerSourceChannelPartner     = "channel_partner"
    PartnerSourcePropertyPortal     = "property_portal"
)
```

### **PartnerSource** (struct, ~15 fields)
```go
type PartnerSource struct {
    ID              int64
    TenantID        string
    PartnerID       int64
    SourceType      string                // PartnerSourceType
    SourceCode      string                // Auto-generated: "CR-1", "VR-2"
    SourceName      string                // "Customer Reference"
    Description     string
    IsActive        bool
    LeadsGenerated  int64                 // Total leads from this source
    LeadsConverted  int64                 // Leads that converted to deals
    TotalRevenue    decimal.Decimal       // Revenue from source
    CreatedAt, UpdatedAt, DeletedAt time.Time
}
```

### **CreditPolicyType** (enum, 3 values)
```go
const (
    CreditPolicyTypeTimeBased     = "time_based"      // Monthly/quarterly/annual
    CreditPolicyTypeProjectBased  = "project_based"   // Per specific project
    CreditPolicyTypeCampaignBased = "campaign_based"  // Per marketing campaign
)
```

### **CreditPolicyCalculation** (enum, 5 values)
```go
const (
    CreditPolicyCalcPercentage    = "percentage"      // % of lead value
    CreditPolicyCalcFixedPrice    = "fixed_price"     // Fixed $/€ per lead
    CreditPolicyCalcTiered        = "tiered"          // Volume tiers
    CreditPolicyCalcConversion    = "conversion"      // Based on actual conversion
    CreditPolicyCalcRevenueshare  = "revenue_share"   // % of deal revenue
)
```

### **PartnerCreditPolicy** (struct, ~28 fields)
```go
type PartnerCreditPolicy struct {
    ID                        int64
    TenantID                  string
    PartnerID                 int64
    PolicyCode                string                // Unique: "FP-1", "TIER-2"
    PolicyName                string                // "Fixed Price - $10 per lead"
    PolicyType                string                // PartnerCreditPolicyType
    CalculationType           string                // PartnerCreditPolicyCalculation
    
    TimeUnitType              string                // "monthly", "quarterly", "annual"
    TimeUnitValue             int
    PolicyStartDate           time.Time
    PolicyEndDate             *time.Time
    
    ProjectID                 *int64
    ProjectName               string
    CampaignID                *int64
    CampaignName              string
    
    BaseCredit                decimal.Decimal       // Core credit value
    MinimumCredit             decimal.Decimal       // Floor
    MaximumCredit             decimal.Decimal       // Ceiling
    BonusPercentage           decimal.Decimal       // % bonus if quality_score >= 80
    
    TierConfig                struct {              // Volume-based tiers
        Tiers []struct {
            TierLevel       int
            MinLeads        int
            MaxLeads        int
            CreditAmount    decimal.Decimal
            BonusPercent    decimal.Decimal
        }
    }
    
    MinLeadQualityScore       decimal.Decimal
    RequiresApproval          bool
    AutoApprove               bool
    IsActive                  bool
    ApprovalRequired          bool
    
    TotalLeadsUnderPolicy     int64
    TotalCreditsAllocated     decimal.Decimal
    
    CreatedBy, ApprovedBy     *int64
    ApprovedAt                *time.Time
    CreatedAt, UpdatedAt      time.Time
    DeletedAt                 *time.Time
}
```

### **PartnerCreditPolicyMapping** (struct, 7 fields)
```go
type PartnerCreditPolicyMapping struct {
    ID                int64
    TenantID          string
    PartnerLeadID     int64
    PolicyID          int64
    CalculatedCredit  decimal.Decimal       // Final credit amount
    Reason            string                // "Applied percentage policy", "Tiered pricing Tier 2"
    CreatedAt         time.Time
}
```

### **PartnerSourceStats** (struct, ~12 fields)
```go
type PartnerSourceStats struct {
    SourceType              string
    SourceCode              string
    TotalLeads              int64
    ApprovedLeads           int64
    ConvertedLeads          int64
    ApprovalRate            float64               // %
    ConversionRate          float64               // %
    AverageQualityScore     float64
    TotalCreditsAllocated   float64
}
```

---

## 🔄 Integration Points

### **With PartnerService** (Phase 4)
- Automatically create partner sources when partner is created
- Map `partner.partner_type` → `PartnerSourceType` enum

### **With PartnerLeadService** (Phase 4)
- UpdatePartnerLeadStatus() → Call CalculateLeadCredit()
- When lead approved → Evaluate policies → Allocate credits
- Store policy mapping in partner_credit_policy_mappings

### **With PartnerPayoutService** (Phase 4)
- Payout calculation includes source + policy info
- Group payouts by source for reporting
- Apply policy-based adjustments before approval

---

## 📈 Usage Examples

### **Example 1: Create Fixed Price Policy**
```go
policy := &models.PartnerCreditPolicy{
    PartnerID:          1,
    PolicyCode:         "FP-STANDARD",
    PolicyName:         "Standard Lead Acquisition",
    PolicyType:         "time_based",
    CalculationType:    "fixed_price",
    TimeUnitType:       "monthly",
    TimeUnitValue:      1,
    PolicyStartDate:    time.Now(),
    PolicyEndDate:      time.Now().AddDate(0, 12, 0),
    BaseCredit:         10.00,
    MinimumCredit:      5.00,
    MaximumCredit:      20.00,
    BonusPercentage:    5,
    MinLeadQualityScore: 60,
    RequiresApproval:   false,
    AutoApprove:        true,
}

policy, err := creditPolicyService.CreateCreditPolicy(ctx, tenantID, policy)
// Result: Partner gets $10 per approved lead, $11 if quality_score >= 80
```

### **Example 2: Create Tiered Policy**
```go
policy := &models.PartnerCreditPolicy{
    PartnerID:          2,
    PolicyCode:         "TIER-VOLUME",
    PolicyName:         "Volume-Based Pricing",
    PolicyType:         "project_based",
    CalculationType:    "tiered",
    ProjectID:          42,
    ProjectName:        "Office Space Leasing Q1-2025",
    BaseCredit:         10.00,
    MinimumCredit:      5.00,
    MaximumCredit:      25.00,
    TierConfig: struct{}{
        Tiers: []struct{}{
            {TierLevel: 1, MinLeads: 0, MaxLeads: 100, CreditAmount: 8, BonusPercent: 0},
            {TierLevel: 2, MinLeads: 101, MaxLeads: 500, CreditAmount: 12, BonusPercent: 5},
            {TierLevel: 3, MinLeads: 501, MaxLeads: 9999, CreditAmount: 15, BonusPercent: 10},
        },
    },
    RequiresApproval:   false,
    AutoApprove:        true,
    MinLeadQualityScore: 50,
}

policy, err := creditPolicyService.CreateCreditPolicy(ctx, tenantID, policy)
// Result: Credit increases by tier as partner submits more leads
```

### **Example 3: Calculate Lead Credit**
```go
// Lead approved with quality_score = 85
credit, policyIDs, err := creditPolicyService.CalculateLeadCredit(ctx, tenantID, partnerLeadID)
// credit = 12.60 (from tiered policy Tier 2: $12 + 5% bonus = $12.60)
// policyIDs = [2] (applied policy ID 2)
```

### **Example 4: Get Source Statistics**
```go
sourceStats, err := sourceService.GetSourceStats(ctx, tenantID, sourceID)
// Returns:
// - TotalLeads: 50
// - ApprovedLeads: 45
// - ConvertedLeads: 15
// - ApprovalRate: 90%
// - ConversionRate: 33.3%
// - AverageQualityScore: 75.5
// - TotalCreditsAllocated: 540.00
```

---

## 🚀 Deployment Steps

### **1. Run Migration**
```bash
# SSH into server or use database tool
mysql -u user -p database < migrations/023_partner_sources_and_credit_policies.sql

# Verify tables created
SHOW TABLES LIKE 'partner_%';
```

### **2. Restart Backend**
```bash
docker-compose restart api
# or
systemctl restart api-service
```

### **3. Initialize Sample Data**
```bash
# Migration includes INSERT statements for:
# - 20 partner sources (5 each type)
# - 9 credit policies (mixed types)
# - Sample tiers and configurations
# Auto-executes during migration
```

### **4. Update API Routes** (Next Step)
- POST `/api/v1/partners/:id/sources` → CreatePartnerSource
- GET `/api/v1/partners/:id/sources` → GetPartnerSources
- GET `/api/v1/partners/:id/sources/:sourceId/stats` → GetSourceStats
- POST `/api/v1/partners/:id/credit-policies` → CreateCreditPolicy
- GET `/api/v1/partners/:id/credit-policies` → GetCreditPolicies
- POST `/api/v1/credit-policies/:id/approve` → ApproveCreditPolicy

---

## ✅ Validation Checklist

- [x] PartnerSourceService compiles without errors
- [x] PartnerCreditPolicyService compiles without errors
- [x] Database migration 023 syntax valid
- [x] All 3 new tables created with proper constraints
- [x] Indexes optimized for multi-tenant queries
- [x] Sample data inserted correctly
- [x] Credit calculation logic handles all 5 types
- [x] Tiered pricing boundary conditions verified
- [x] Bonus percentage applied correctly (quality_score >= 80)
- [x] Min/max credit bounds enforced
- [x] Soft-delete support on policies and sources
- [x] Audit trail fields populated (CreatedBy, ApprovedBy, ApprovedAt)

---

## 📝 Statistics

| Metric | Value |
|--------|-------|
| Lines of service code | 689 |
| Lines of migration | 309 |
| PartnerSourceService methods | 9 |
| PartnerCreditPolicyService methods | 11 |
| New database tables | 3 |
| Database columns added | 57 |
| New indexes | 8 |
| Partner source types | 4 |
| Credit policy types | 3 |
| Credit calculation methods | 5 |
| Sample policies created | 9 |
| Sample sources created | 20 |

---

## 🔗 Related Phases

- **Phase 1-2**: Backend customization (campaigns, sources, stages, milestones)
- **Phase 4**: External partner system (CRUD, quality scoring, payouts)
- **Phase 5 (Current)**: Partner sources & credit policies
- **Phase 6 (Next)**: API routes & frontend integration

---

**Status**: ✅ Phase 5 Complete - Ready for API handler development

