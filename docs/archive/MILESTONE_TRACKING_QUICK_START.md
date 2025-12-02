# MILESTONE TRACKING - QUICK REFERENCE GUIDE

## 🎯 Quick Start (5 Minutes)

### 1. Access Milestone Tracking
```
Navigate to: Dashboard → Sales → "Milestones & Tracking" tab
```

### 2. Track a Lead's Journey
```
1. Enter Lead ID (from Lead Management tab)
2. Click "Load" to fetch existing milestones
3. Click "New Milestone" to add event
4. Fill form and submit
```

### 3. Record Engagement
```
1. Same lead ID
2. Switch to "Engagement" tab
3. Click "New Engagement"
4. Select type (call, email, etc.)
5. Add details and submit
```

---

## 📅 Milestone Types & When to Use

| Type | Use Case | Example |
|------|----------|---------|
| **lead_generated** | When lead enters system | Imported from website form |
| **contacted** | First contact made | Initial email or call sent |
| **site_visit** | In-person visit to office/property | Team visited site on Nov 25 |
| **revisit** | Follow-up site visit | Second visit after negotiations |
| **demo** | Product/service demonstration | Showed 3D walkthrough |
| **proposal** | Formal proposal sent | Sent quotation doc |
| **negotiation** | Active negotiation phase | Discussing terms |
| **booking** | Customer confirmed booking | 500k booking confirmed |
| **cancellation** | Deal cancelled | Budget constraints |
| **reengaged** | Dead lead contacted again | Old lead got budget approval |

---

## 💬 Engagement Types & Channels

### Types
- **email_sent** → Track outgoing emails
- **call_made** → Phone conversations
- **message_sent** → SMS/WhatsApp messages
- **meeting_scheduled** → Calendar events
- **proposal_sent** → Document sharing
- **quote_sent** → Pricing documents

### Channels
- **email** → Email communications
- **phone** → Voice calls
- **sms** → Text messages
- **whatsapp** → WhatsApp messages
- **in_person** → Face-to-face meetings
- **video** → Video calls (Zoom, Teams, etc.)

---

## 📊 Reports Dashboard - All Views

### Dashboard Tab (KPI Overview)
**What to monitor**:
- Total Leads trend
- Conversion Rate % (target: 20%+)
- Active Customers count
- Outstanding Balance (AR health)
- Pending Follow-ups (action items)

**Action when high**:
- High pending follow-ups → Assign to team
- Low conversion rate → Review sales process
- High outstanding → Escalate collections

### Lead Funnel Tab (12-Month History)
**What to look for**:
- Funnel leakage (where leads drop off)
- Monthly trends (seasonal patterns)
- Conversion rate improvement

**Example analysis**:
- Month: Oct 2025
- New: 50 | Contacted: 40 | Qualified: 30 | Converted: 5
- Conv Rate: 10% (opportunity to improve)

### Source Performance Tab (Channel Analysis)
**What to evaluate**:
- Best performing source (highest conversion)
- Cost per lead by source
- Quality of leads by source

**Example**: 
- Google Ads: 100 leads, 25 converted (25% rate) ← Best
- Referral: 20 leads, 8 converted (40% rate) ← Highest rate
- Direct: 80 leads, 4 converted (5% rate) ← Needs attention

### Bookings Tab (Sales Analytics)
**What to track**:
- Confirmed bookings (revenue pipeline)
- Pending bookings (follow-up needed)
- Cancelled bookings (churn analysis)
- Average booking size (deal health)

---

## 💰 Financial Tracking - Account Ledger

### Transaction Types
```
INVOICE     → Debit amount (customer owes)
PAYMENT     → Credit amount (customer paid)
CREDIT_NOTE → Credit amount (refund/adjustment)
DEBIT_NOTE  → Debit amount (additional charge)
ADJUSTMENT  → Manual adjustment for corrections
```

### Balance Calculation
```
Formula: Running Balance = Previous Balance + Debits - Credits

Example:
Invoice for 100k         → Balance = 100,000 (customer owes)
Payment of 60k          → Balance = 40,000 (still owes)
Credit note of 10k      → Balance = 30,000 (after discount)
Payment of 30k          → Balance = 0 (settled)
```

### Checking Customer Health
```
Outstanding Balance < 50k         ← Good
Outstanding Balance 50-100k       ← Monitor
Outstanding Balance > 100k        ← High risk (follow-up needed)
Past due (>30 days)              ← Escalate
```

---

## 🎬 Common Workflows

### Workflow 1: New Lead to Booking (7 Steps)
```
Day 1:  Lead Generated    → Milestone with source tag
        Email Sent        → Engagement record
Day 2:  Call Made         → Engagement with notes
        Qualified         → Status milestone
Day 4:  Site Visit        → Milestone with location, outcome
        Follow-up Set     → Milestone with follow-up date
Day 7:  Demo Scheduled    → Meeting engagement
Day 14: Proposal Sent     → Engagement type
Day 21: Booking Confirmed → Milestone + Booking record
Day 22: Invoice Created   → Ledger entry (debit)
Day 30: Payment Received  → Ledger entry (credit)
```

### Workflow 2: Campaign Performance Tracking
```
Step 1: Create Campaign
        Name: "Q4 2025 Email Blast"
        Type: Email
        Budget: ₹50,000
        Expected ROI: 300%

Step 2: Link Leads to Campaign
        When creating leads, set campaign_id

Step 3: Track in Reports
        Go to Lead Funnel → See conversion rate for campaign

Step 4: Analyze Results
        50 leads → 10 converted → 20% rate → ₹150k revenue
        ROI = (150k - 50k) / 50k × 100 = 200%
```

### Workflow 3: Cancellation Tracking
```
Step 1: Record Cancellation Milestone
        milestone_type: "cancellation"
        notes: "Budget constraints - client company downsizing"

Step 2: Update Booking Status
        status: "cancelled"
        cancellation_reason: "Budget constraints"
        cancellation_refund_amount: 500000

Step 3: Create Ledger Entry
        transaction_type: "credit_note"
        credit_amount: 500000 (full refund)
        remarks: "Full refund for cancelled booking"

Step 4: Track for Analysis
        View Bookings Report → See cancellation trend
```

---

## 📈 KPI Formulas (Dashboard Calculates Automatically)

```
Conversion Rate = (Converted Leads / Total Leads) × 100
Lead Quality % = (Qualified Leads / Total Leads) × 100
Customer LTV = Total Invoiced / Customer Count
Payment Collection % = (Total Paid / Total Invoiced) × 100
Booking Rate = (Total Bookings / Converted Leads) × 100
Days to Conversion = AVG(Booking Date - Lead Created Date)
Engagement Frequency = Total Engagements / Active Leads
Outstanding AR = Σ(Debit - Credit) from Ledger
```

---

## 🔍 Data Entry Best Practices

### When Creating Milestone
✅ Always set milestone_date  
✅ Use specific milestone_type (not generic)  
✅ Add detailed notes for context  
✅ Set follow_up_required = true if action needed  
✅ Record location for site visits  
✅ Capture outcome (positive/neutral/negative)  

### When Recording Engagement
✅ Always set engagement_date  
✅ Be specific in subject line  
✅ Add outcome in notes  
✅ Mark response_received when applicable  
✅ Set engagement_channel (not just type)  
✅ Track call/meeting duration  

### When Creating Booking
✅ Set unique booking_code  
✅ Record exact booking_date  
✅ Specify unit type and count  
✅ Set delivery_date for visibility  
✅ Link to correct customer_id  
✅ Include notes on special terms  

### When Creating Ledger Entry
✅ Use sequential ledger_code  
✅ Set transaction_date (don't leave blank)  
✅ Ensure debit/credit adds up  
✅ Link to source document (invoice, payment, etc.)  
✅ Add description for clarity  
✅ Verify running balance  

---

## ⚠️ Common Mistakes to Avoid

❌ Using same milestone multiple times without dates  
   → Use one milestone per event with specific date

❌ Not linking leads to campaigns  
   → Always set campaign_id when creating leads

❌ Leaving engagement subject blank  
   → Always document what was discussed

❌ Incorrect ledger balances  
   → Verify: running balance = previous + debits - credits

❌ Not recording follow-ups  
   → Always check follow_up_required and set follow_up_date

❌ Mixing engagement types  
   → "call_made" for calls, "email_sent" for emails, etc.

❌ Soft-deleted leads still showing in reports  
   → System auto-filters, but verify in queries

---

## 📞 Troubleshooting

### Milestones not appearing?
1. Verify lead_id exists in database
2. Check X-Tenant-ID header matches user's tenant
3. Ensure dates are in ISO format (YYYY-MM-DD)
4. Try refreshing browser cache

### Reports showing no data?
1. Verify leads have milestones/engagements
2. Check date range includes current month
3. Ensure milestones have follow-up_date set
4. Try different report period

### Ledger balance wrong?
1. Verify transaction_type is correct
2. Check debit/credit amounts are positive
3. Ensure no duplicate entries
4. Recalculate manually

### Dashboard metrics zero?
1. Ensure at least one lead exists
2. Check leads have created_at date
3. Verify leads have status field
4. Try refreshing dashboard tab

---

## 🚀 Pro Tips

### Tip 1: Use Milestones for Pipeline Visibility
Create milestone for every status change → Clear pipeline visibility

### Tip 2: Set Follow-ups Regularly
Enable follow_up_required on site visits → Auto-remind team

### Tip 3: Tag Campaigns on All Leads
Always set campaign_id → Track marketing ROI accurately

### Tip 4: Monitor Daily Engagements
Check "Engagement This Month" widget → Track team activity

### Tip 5: Review Funnel Monthly
Monthly funnel review → Identify bottlenecks → Improve process

### Tip 6: Track Visit Outcomes
Record positive/neutral/negative → Data for lead scoring

### Tip 7: Keep Ledger Current
Update ledger same day as transaction → Accurate AR reporting

---

## 📞 Support Contacts

**Feature Questions**: Review MILESTONE_TRACKING_COMPLETE.md  
**API Questions**: Check COMPLETE_API_REFERENCE.md  
**Database Issues**: Contact DBA team  
**Performance Issues**: Monitor query performance on dashboard  

---

**Last Updated**: November 25, 2025  
**Version**: 1.0.0  
**Status**: Production Ready ✅
