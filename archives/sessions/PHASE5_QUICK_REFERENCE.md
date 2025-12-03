# Phase 5: Quick Reference - Partner Sources & Credit Policies

## 📊 Service Interfaces Quick Lookup

### PartnerSourceService
```go
// Source Management
CreatePartnerSource(ctx, tenantID, source) → *PartnerSource
GetPartnerSource(ctx, tenantID, sourceID) → *PartnerSource
GetPartnerSources(ctx, tenantID, partnerID) → []PartnerSource
GetSourceByCode(ctx, tenantID, sourceCode) → *PartnerSource
UpdatePartnerSource(ctx, tenantID, source) → *PartnerSource
DeactivatePartnerSource(ctx, tenantID, sourceID) → error

// Statistics
GetSourceStats(ctx, tenantID, sourceID) → *PartnerSourceStats
GetPartnerSourceStats(ctx, tenantID, partnerID) → []PartnerSourceStats
```

### PartnerCreditPolicyService
```go
// Policy Management
CreateCreditPolicy(ctx, tenantID, policy) → *PartnerCreditPolicy
GetCreditPolicy(ctx, tenantID, policyID) → *PartnerCreditPolicy
GetPartnerCreditPolicies(ctx, tenantID, partnerID) → []PartnerCreditPolicy
GetActiveCreditPolicies(ctx, tenantID, partnerID) → []PartnerCreditPolicy
UpdateCreditPolicy(ctx, tenantID, policy) → *PartnerCreditPolicy
ApproveCreditPolicy(ctx, tenantID, policyID, approvedBy) → error
DeactivateCreditPolicy(ctx, tenantID, policyID) → error

// Credit Calculation
CalculateLeadCredit(ctx, tenantID, partnerLeadID) → (float64, []int64, error)
GetApplicablePolicies(ctx, tenantID, partnerID, leadData) → []PartnerCreditPolicy
GetPolicyMappings(ctx, tenantID, partnerLeadID) → []PartnerCreditPolicyMapping
```

---

## 🔗 4 Partner Source Types

```go
const (
    "customer_reference"    // B2B: Customers referring new leads
    "vendor_reference"      // Vendors/suppliers referring opportunities
    "channel_partner"       // Channel/reseller partners
    "property_portal"       // Real estate portals, aggregators
)
```

**Auto-Mapping**: When partner created with `partner_type`, auto-create corresponding PartnerSource

---

## 💰 5 Credit Calculation Methods

| Method | Formula | Example |
|--------|---------|---------|
| **percentage** | `credit = base_credit` (%) | 10% of lead value |
| **fixed_price** | `credit = base_credit` ($) | $10 per lead |
| **tiered** | Find tier by lead_count → credit_amount | $8 (1-50), $12 (51-200), $15 (201+) |
| **conversion** | Credit only on conversion | $100 per closed deal |
| **revenue_share** | `credit = deal_revenue × base_credit%` | 10% of $500 deal = $50 |

**All methods**: Apply bonus if quality_score >= 80, then apply min/max bounds

---

## 📅 3 Policy Types

| Type | Applies To | Example Use Case |
|------|-----------|------------------|
| **time_based** | Recurring periods | Monthly commission, $10/lead all month |
| **project_based** | Specific project | Q1 Office Leasing project, tiered pricing |
| **campaign_based** | Marketing campaign | "Summer Promo 2025", 15% revenue share |

---

## 🏆 Tiered Pricing Structure

```json
{
  "tiers": [
    {
      "tier_level": 1,
      "min_leads": 0,
      "max_leads": 100,
      "credit_amount": 8,
      "bonus_percent": 0
    },
    {
      "tier_level": 2,
      "min_leads": 101,
      "max_leads": 500,
      "credit_amount": 12,
      "bonus_percent": 5
    },
    {
      "tier_level": 3,
      "min_leads": 501,
      "max_leads": 999999,
      "credit_amount": 15,
      "bonus_percent": 10
    }
  ]
}
```

**Logic**: 
```
Lead count = 250
Find tier: min_leads <= 250 <= max_leads
Matches Tier 2: credit_amount = $12
Apply bonus 5%: $12.60 if quality_score >= 80
```

---

## 📊 Database Tables

### partner_sources (14 columns)
- **Key Fields**: source_type (enum), source_code (unique), leads_generated, leads_converted, total_revenue
- **Soft Delete**: deleted_at nullable timestamp
- **Indexes**: (tenant_id, partner_id), (source_type), (is_active)

### partner_credit_policies (36 columns)
- **Policy ID**: policy_code (unique per tenant), policy_name
- **Type**: policy_type (time/project/campaign), calculation_type (5 options)
- **Configuration**: base_credit, min/max_credit, bonus_percentage, tier_config (JSON)
- **Conditions**: min_lead_quality_score, requires_approval, auto_approve
- **Approval**: approval_required, approved_by, approved_at
- **Stats**: total_leads_under_policy, total_credits_allocated
- **Indexes**: (tenant_id, partner_id), (policy_type), (is_active), (approval_required)

### partner_credit_policy_mappings (7 columns)
- **Link**: partner_lead_id → policy_id
- **Result**: calculated_credit, reason (audit trail)
- **Indexes**: (tenant_id, partner_lead_id), (policy_id)

---

## 🎯 Credit Calculation Example Walkthrough

**Scenario**: Partner submits 150th lead (from 200-lead tier)
- Quality score: 85 (high quality)
- Policy: Tiered with 3 tiers (1-100: $8, 101-500: $12, 501+: $15)
- Bonus percentage: 5%
- Min/max: $5-$25

**Calculation**:
1. Find tier: lead_count (150) in 101-500 range → Tier 2
2. Base credit: $12
3. Apply bonus: quality_score (85) >= 80 → +5% = $12 × 1.05 = $12.60
4. Verify bounds: $12.60 between $5-$25 ✓
5. **Final credit: $12.60**

---

## 🚀 Integration Checklist

- [ ] **PartnerSourceService initialized** in main.go/dependency injection
- [ ] **PartnerCreditPolicyService initialized** in main.go/dependency injection
- [ ] **Migration 023 executed** on database
- [ ] **Update PartnerLeadService**: Call CalculateLeadCredit() on lead approval
- [ ] **Update PartnerPayoutService**: Include source/policy info in payout calculation
- [ ] **Create API routes**:
  - POST /api/v1/partners/:id/sources
  - GET /api/v1/partners/:id/sources
  - GET /api/v1/partners/:id/sources/:sourceId/stats
  - POST /api/v1/partners/:id/credit-policies
  - GET /api/v1/partners/:id/credit-policies
  - POST /api/v1/credit-policies/:id/approve
- [ ] **Create frontend components**:
  - Source tracking dashboard
  - Credit policy builder UI
  - Tiered pricing visualizer

---

## 📁 Files Reference

| File | Lines | Purpose |
|------|-------|---------|
| `internal/services/partner_source_credit_service.go` | 689 | 2 services, 20 methods |
| `migrations/023_partner_sources_and_credit_policies.sql` | 309 | 3 tables, indexes, sample data |
| `internal/models/partner.go` | +150 | 8 new model definitions |
| **Total Phase 5** | **1,148** | Complete partner enhancement |

---

## ⚡ Common Operations

### Create Monthly Fixed-Price Policy
```go
policy := &models.PartnerCreditPolicy{
    PartnerID: partnerID,
    PolicyCode: "MONTHLY-FP",
    PolicyName: "Monthly Fixed Price",
    PolicyType: "time_based",
    CalculationType: "fixed_price",
    TimeUnitType: "monthly",
    BaseCredit: 10,
    BonusPercentage: 5,
    MinLeadQualityScore: 60,
    AutoApprove: true,
}
policy, _ := creditPolicySvc.CreateCreditPolicy(ctx, tenantID, policy)
```

### Calculate Lead Credit
```go
credit, policyIDs, _ := creditPolicySvc.CalculateLeadCredit(ctx, tenantID, leadID)
// credit: final amount, policyIDs: applied policies
```

### Get All Partner Sources with Stats
```go
sources, _ := sourceSvc.GetPartnerSources(ctx, tenantID, partnerID)
for _, src := range sources {
    stats, _ := sourceSvc.GetSourceStats(ctx, tenantID, src.ID)
    // Use stats for dashboard
}
```

### Approve New Policy (Management)
```go
creditPolicySvc.ApproveCreditPolicy(ctx, tenantID, policyID, managedBy)
```

---

## 🔐 Multi-Tenant Safety

✅ All queries filter by `tenant_id`
✅ All operations start with TenantID validation
✅ Soft deletes used throughout (deleted_at nullable)
✅ Audit trail: CreatedBy, ApprovedBy timestamps
✅ No cross-tenant data leakage possible

---

## 📈 Typical Flow

```
Partner Registered
    ↓
Auto-create PartnerSource (maps partner_type → source_type)
    ↓
Create PartnerCreditPolicy (time/project/campaign-based)
    ↓
Policy requires approval
    ↓
Management approves policy
    ↓
Partner submits lead
    ↓
On lead approval: CalculateLeadCredit()
    ↓
Create PartnerCreditPolicyMapping (audit trail)
    ↓
Credit allocated to lead
    ↓
Payout generated with source + policy info
    ↓
Management approves payout
```

---

**Status**: ✅ Ready for API handler implementation
