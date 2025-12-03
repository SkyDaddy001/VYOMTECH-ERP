# 🚀 External Partner System - Quick Reference

**Status**: ✅ Production Ready | **Files**: 4 production + 1 migration + 2 docs

---

## 📋 Quick Facts

| Aspect | Details |
|--------|---------|
| **Partner Types** | Portal, Channel Partner, Vendor, Customer |
| **Partner Statuses** | Pending, Active, Inactive, Suspended, Rejected |
| **Lead Submission Types** | New Lead, Referral, Import Batch |
| **Approval Model** | Management-approved only (NO auto-payouts) |
| **Commission Models** | Percentage, Fixed Price, Hybrid |
| **Payment Methods** | Bank Transfer, PayPal, Check, Wire |
| **Quality Score** | 0-100 (auto-calculated on submission) |
| **Tables** | 8 (partners, users, leads, credits, payouts, details, activities, audit) |
| **Service Methods** | 50+ methods across 3 services |
| **Multi-Tenant** | ✅ Yes - all queries tenant-scoped |

---

## 🔐 Core Workflow

### 1️⃣ Partner Registration
```
Partner applies → Status: pending
                    ↓
           Admin reviews docs
                    ↓
        Admin approves → Status: active
                    ↓
        Partner creates user account
                    ↓
       Partner can now submit leads
```

### 2️⃣ Lead Submission
```
Partner submits lead → Quality score auto-calculated
                    ↓
           Status: submitted
                    ↓
        Management reviews (sorted by quality)
                    ↓
    ✅ Approve → Lead ID assigned → Status: approved
    ❌ Reject → Reason tracked → Status: rejected
```

### 3️⃣ Lead-to-Credit Flow
```
Lead approved by management
                    ↓
    Lead converts to sale (in system)
                    ↓
Partner requests credit approval
                    ↓
Management approves → Balance updated
Status: approved
```

### 4️⃣ Payout Flow (KEY: MANAGEMENT APPROVAL REQUIRED)
```
Payout period ends (monthly)
                    ↓
    System generates batch (pending)
                    ↓
    Management reviews details
                    ↓
    ✅ Approve → Status: approved → Payment processed
    ⏱️ Partial → Approve some items → Status: partially_approved
    ❌ Reject → Status: rejected → Partner resubmits
```

---

## 📁 Files Overview

### Models: `internal/models/partner.go`
```go
Partner              // Organization
PartnerUser         // User account
PartnerLead         // Submitted lead
PartnerLeadCredit   // Credit request
PartnerPayout       // Payout batch
PartnerPayoutDetail // Payout line item
PartnerActivity     // Audit log
```

### Services

**`partner_service.go`** - Partner & user management (19 methods)
```go
CreatePartner, GetPartner, GetPartners, UpdatePartner
UpdatePartnerStatus, DeactivatePartner, SuspendPartner
CreatePartnerUser, GetPartnerUser, GetPartnerUsers
UpdatePartnerUser, UpdatePartnerUserPassword, DeactivatePartnerUser
GetPartnerStats, GetPartnerMonthlyStats
CalculateLeadQualityScore
```

**`partner_lead_service.go`** - Lead submission & approval (16 methods)
```go
SubmitPartnerLead, GetPartnerLead, GetPartnerLeads
UpdatePartnerLeadStatus
ApprovePartnerLead, RejectPartnerLead          // ⭐ Management
GetPendingLeadsForReview
GetLeadCredits, SubmitLeadCreditApprovalRequest
ApproveLeadCredit, RejectLeadCredit           // ⭐ Management
LogPartnerActivity, GetPartnerActivity
```

**`partner_payout_service.go`** - Payout management (18 methods)
```go
GeneratePayoutPeriod, CreatePayout, GetPayout, GetPayouts
GetPendingPayouts
ApprovePayout, RejectPayout                   // ⭐ Management (CRITICAL)
PartiallyApprovePayout
GetPayoutDetails, AddPayoutDetail
ApprovePayoutDetail, RejectPayoutDetail       // ⭐ Line-item approval
MarkPayoutAsPaid
GetPayoutStats
```

### Database: `migrations/022_external_partner_system.sql`
```sql
partners                   -- Org profiles
partner_users             -- User accounts
partner_leads             -- Lead submissions
partner_lead_credits      -- Credit requests
partner_payouts           -- Payout batches
partner_payout_details    -- Payout lines
partner_activities        -- Audit log
```

---

## 🔑 Key Constants

### Partner Types
```go
"portal"           // White-label portal
"channel_partner"  // Channel/reseller
"vendor"          // Lead supplier
"customer"        // Direct customer
```

### Partner Statuses
```go
"pending"    // Awaiting approval
"active"     // Approved & can submit
"inactive"   // Deactivated by admin
"suspended"  // Suspended for violations
"rejected"   // Rejected during approval
```

### Lead Statuses
```go
"submitted"    // Initial state
"under_review" // Being reviewed
"approved"     // Accepted by management
"rejected"     // Rejected by management
"converted"    // Became customer
```

### Payout Statuses
```go
"pending"            // Generated, awaiting review
"approved"           // Approved by management
"rejected"           // Rejected by management
"partially_approved" // Some items approved
"paid"              // Payment processed
```

### Credit Statuses
```go
"pending_approval" // Awaiting management
"approved"        // Approved by management
"rejected"        // Rejected by management
"paid"           // Included in paid payout
```

---

## 💡 Usage Examples

### Initialize Services
```go
import "vyomtech/internal/services"

db := // ... database connection
partnerSvc := services.NewPartnerService(db)
leadSvc := services.NewPartnerLeadService(db, partnerSvc)
payoutSvc := services.NewPartnerPayoutService(db, partnerSvc, leadSvc)
```

### Create Partner (Registration)
```go
partner := &models.Partner{
    OrganizationName: "Lead Portal Inc",
    PartnerType:      models.PartnerTypePortal,
    ContactEmail:     "contact@leadportal.com",
    ContactPerson:    "John Doe",
    BankingDetails: models.BankingDetails{
        BankName:      "Bank of America",
        AccountNumber: "****1234",
        PaymentMethod: "bank_transfer",
    },
    CommissionPercentage: 15.00,
    CreatedBy:            userId,
}

created, err := partnerSvc.CreatePartner(ctx, tenantID, partner)
```

### Approve Partner (Admin)
```go
err := partnerSvc.UpdatePartnerStatus(
    ctx, 
    tenantID, 
    partnerID, 
    models.PartnerStatusActive, 
    "", // no rejection reason
    adminUserID,
)
```

### Submit Lead (Partner)
```go
lead := &models.PartnerLead{
    PartnerID:      partnerID,
    SubmissionType: "new_lead",
    LeadData: models.LeadData{
        FirstName: "Jane",
        LastName:  "Smith",
        Email:     "jane@example.com",
        Phone:     "+1-555-0123",
        Company:   "ABC Corp",
        JobTitle:  "Manager",
        Budget:    "high",
        Timeline:  30,
    },
    SubmittedBy: partnerUserID,
}

submitted, err := leadSvc.SubmitPartnerLead(ctx, tenantID, lead)
// Quality score auto-calculated
// Status: "submitted"
```

### Approve Lead (Management)
```go
// First create actual lead
actualLead := &models.Lead{
    TenantID:  tenantID,
    Name:      "Jane Smith",
    Email:     "jane@example.com",
    Phone:     "+1-555-0123",
    Status:    "new",
    Source:    "partner_referral",
}

newLead, _ := leadService.CreateLead(ctx, tenantID, actualLead)

// Then approve partner lead
err := leadSvc.ApprovePartnerLead(
    ctx,
    tenantID,
    partnerLeadID,
    managerID,          // Who approved
    newLead.ID,         // Link to actual lead
)
```

### Approve Credit (Management - REQUIRED)
```go
credit := &models.PartnerLeadCredit{
    PartnerLeadID:    partnerLeadID,
    PartnerID:        partnerID,
    CreditAmount:     150.00,
    CalculationType:  "percentage", // or "fixed_price"
}

submitted, _ := leadSvc.SubmitLeadCreditApprovalRequest(ctx, tenantID, credit)

// Management must approve
err := leadSvc.ApproveLeadCredit(
    ctx,
    tenantID,
    creditID,
    managerID,  // Who approved
)
// Now partner.available_balance increases
```

### Generate Payout Batch
```go
payout, err := payoutSvc.GeneratePayoutPeriod(
    ctx,
    tenantID,
    partnerID,
    startDate,  // 2025-01-01
    endDate,    // 2025-01-31
)
// Automatically aggregates approved leads
// Status: "pending"
```

### Approve Payout (Management - CRITICAL)
```go
// Get pending payouts for review
pendings, _ := payoutSvc.GetPendingPayouts(ctx, tenantID, 50, 0)

// Review details
details, _ := payoutSvc.GetPayoutDetails(ctx, tenantID, payoutID, 100, 0)

// Management approves payout
err := payoutSvc.ApprovePayout(
    ctx,
    tenantID,
    payoutID,
    5000.00,    // Approved amount (can be less than total)
    managerID,  // Who approved
)
// Status: "approved"
// partner.withdrawn_amount updated
// partner.available_balance updated
```

### Mark Payment Processed
```go
err := payoutSvc.MarkPayoutAsPaid(
    ctx,
    tenantID,
    payoutID,
    time.Now(),
    "TXN-2025-001", // Reference number
)
// Status: "paid"
```

---

## 📊 Dashboard Queries

### Partner Stats
```go
stats, _ := partnerSvc.GetPartnerStats(ctx, tenantID, partnerID)
// Returns:
//   total_leads_submitted, approved_leads, rejected_leads
//   conversion_rate, approval_rate, total_earnings
//   available_balance, pending_payout, average_lead_quality
```

### Pending Leads for Review
```go
leads, _ := leadSvc.GetPendingLeadsForReview(ctx, tenantID, 50, 0)
// Returns leads sorted by quality_score DESC
// For management dashboard
```

### Pending Payouts for Approval
```go
payouts, _ := payoutSvc.GetPendingPayouts(ctx, tenantID, 50, 0)
// Returns payouts with status = "pending"
// For management dashboard
```

### Partner Activity Log
```go
activities, _ := leadSvc.GetPartnerActivity(ctx, tenantID, partnerID, 50, 0)
// Returns all actions (lead_submitted, lead_approved, payout_requested, etc)
```

---

## 🔒 Security Checklist

- ✅ Always include `tenantID` in all queries
- ✅ Check user permissions before operations
- ✅ Verify `approvedBy` user is manager/admin
- ✅ Audit log all approvals/rejections
- ✅ Validate banking details before payout
- ✅ Never skip management approval for payouts
- ✅ Sanitize all user input
- ✅ Use context timeouts to prevent hanging

---

## 📈 Performance Tips

1. **Index for common queries**:
   - `(tenant_id, status)` on all tables
   - `(quality_score DESC)` on partner_leads
   - `(created_at DESC)` on all tables

2. **Use pagination**: Always include limit/offset

3. **Pre-calculate stats**: Store in partner row

4. **Archive old payouts**: Move to archive table after 1 year

5. **Monitor approval times**: Track approval lag

---

## 🚨 Critical Points

### ⭐ PAYOUTS ARE MANAGEMENT-APPROVED ONLY
- ✅ Partner can REQUEST payout
- ❌ Partner CANNOT withdraw automatically
- ✅ Management MUST explicitly approve
- ✅ All approvals LOGGED with timestamp & approver
- ❌ Partial payments requires re-approval

### ⭐ ALL CREDITS NEED APPROVAL
- ✅ Partner can request credit
- ❌ Credit NOT automatically added to balance
- ✅ Management MUST approve each credit
- ✅ Only approved credits increase balance

### ⭐ MULTI-TENANT ISOLATION
- Every query MUST filter by tenant_id
- Partner data never visible across tenants
- Payouts isolated per tenant

---

## 📞 Troubleshooting

**Q: Partner can't submit leads**  
A: Check partner status = "active" and user role = "lead_manager" or "admin"

**Q: Lead quality score is 0**  
A: Check lead_data completeness - missing required fields reduces score

**Q: Payout stuck in pending**  
A: Management hasn't approved yet - check pending payout dashboard

**Q: Balance not updating after approval**  
A: Ensure LeadCredit was approved - that's what updates balance

**Q: Can't create partner**  
A: Check organization_name & contact_email are provided

---

## 📚 Full Documentation

- **Complete Details**: See `EXTERNAL_PARTNER_SYSTEM_COMPLETE.md`
- **Delivery Summary**: See `EXTERNAL_PARTNER_SYSTEM_DELIVERY.md`
- **Code Comments**: See each service file for detailed method docs

---

**Status**: ✅ Ready for API Handler Development  
**Next**: Create REST endpoints for all services  
**Then**: Build partner portal UI & management dashboard
