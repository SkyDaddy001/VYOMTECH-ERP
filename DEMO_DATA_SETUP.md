# Vyomtech Demo Data & 30-Day Auto-Reset Setup

## Overview

The Vyomtech ERP system includes a comprehensive demo tenant (`demo_vyomtech_001`) with sample data across all modules that **automatically resets every 30 days**. This enables:

- **Instant Testing**: Use pre-configured credentials to test the full platform
- **Multi-Partner Testing**: Test all 4 partner types (portal, channel, vendor, customer)
- **Real-World Scenarios**: Demo includes leads, campaigns, projects, agents, compliance data
- **Fresh Data Every Month**: Auto-reset ensures clean demo state

---

## Demo Credentials

### Primary Access Points

| Role | Email | Password | Partner Type |
|------|-------|----------|--------------|
| **Portal Admin** | `demo@vyomtech.com` | `demo123` | Portal |
| **Channel Manager** | `channel@demo.vyomtech.com` | `demo123` | Channel Partner |
| **Vendor Lead** | `vendor@demo.vyomtech.com` | `demo123` | Vendor |
| **Customer Manager** | `customer@demo.vyomtech.com` | `demo123` | Customer |

### Agent Credentials

| Agent Name | Email | Password | Performance Score |
|------------|-------|----------|-------------------|
| Rajesh Kumar | `rajesh@demo.vyomtech.com` | `demo123` | 95% |
| Priya Singh | `priya@demo.vyomtech.com` | `demo123` | 92% |
| Arun Patel | `arun@demo.vyomtech.com` | `demo123` | 88% |
| Neha Sharma | `neha@demo.vyomtech.com` | `demo123` | 90% |

---

## Demo Data Contents

### Partners (4 Total)

```
1. Vyomtech Portal Demo (portal)
   - Code: DEMO_PORTAL
   - Status: Active

2. Demo Channel Partner (channel_partner)
   - Code: DEMO_CHANNEL
   - Status: Active

3. Demo Vendor Solutions (vendor)
   - Code: DEMO_VENDOR
   - Status: Active

4. Demo Customer Account (customer)
   - Code: DEMO_CUSTOMER
   - Status: Active
```

### Leads (5 Total)

- **High Value Residential Project** - ₹50 Lakh (Mumbai)
- **Commercial Space Inquiry** - ₹35 Lakh (Bangalore)
- **Plot Purchase Interest** - ₹20 Lakh (Delhi)
- **Rental Inquiry** - ₹15 Lakh (Hyderabad)
- **Apartment Pre-booking** - ₹40 Lakh (Pune)

### Campaigns (4 Total)

1. Summer Residential Drive 2025 (Email) - ₹5 Lakh Budget
2. Commercial Real Estate Expo (Event) - ₹10 Lakh Budget
3. Digital Marketing Campaign (Social) - ₹3 Lakh Budget
4. Corporate Bulk Purchase (Direct Sales) - ₹8 Lakh Budget

### Projects (4 Total)

1. **Skyrise Towers Mumbai** (Residential)
   - Cost: ₹50 Cr | Status: Active
   - Progress: 65% (Foundation Work)

2. **Tech Park Bangalore** (Commercial)
   - Cost: ₹100 Cr | Status: Active
   - Progress: 45% (Construction Phase)

3. **Plot Development Delhi** (Plot Development)
   - Cost: ₹20 Cr | Status: Planning

4. **Green Spaces Pune** (Mixed Use)
   - Cost: ₹75 Cr | Status: Active

### Additional Data

- **4 Demo Agents** with performance scores (88-95%)
- **4 Tasks** with different priorities and statuses
- **4 Calls** (mix of inbound/outbound, completed)
- **Gamification Data** (points, badges, leaderboards)
- **Compliance Records** (RERA, FEMA, GST, Environmental)
- **Progress Tracking** (project milestones)

---

## Implementation Details

### Files Created

```
migrations/025_vyomtech_demo_data.sql
├─ Demo tenant setup
├─ 4 sample partners + users
├─ Sample agents (4)
├─ Sample leads (5)
├─ Sample campaigns (4)
├─ Sample projects (4)
├─ Sample tasks
├─ Gamification data
├─ Compliance records
└─ Progress tracking

internal/services/demo_reset_service.go
├─ 30-day reset scheduler
├─ Auto-clear old data
├─ Auto-reload fresh data
└─ Transaction management

frontend/app/demo-credentials.tsx
├─ React component
├─ Display on login page
├─ 4 credential blocks
└─ Auto-reset notice

scripts/reset-demo-data.sh
├─ Manual reset script
├─ Backup-aware deletion
├─ Fresh data reload
└─ Cron-schedulable
```

### How Auto-Reset Works

1. **Service Initialization**: `DemoResetService` starts when app boots
2. **First Reset**: Runs immediately on startup if needed
3. **Scheduled Resets**: Every 30 days automatically
4. **Transaction Safe**: Uses database transactions to ensure consistency
5. **Table Cleanup**: Deletes from 18+ related tables
6. **Fresh Reload**: Inserts pristine demo data

### Reset Schedule

```
Day 1:  Demo data loaded
Day 30: Auto-reset triggered
Day 60: Auto-reset triggered
Day 90: Auto-reset triggered
... (every 30 days)
```

---

## Usage

### Login with Demo Account

1. Navigate to login page
2. See "Try the Demo" section with credentials
3. Select any credential pair
4. Use password: `demo123`
5. Explore the full platform

### Manual Reset

To manually reset demo data (emergency/testing):

```bash
# Linux/macOS
bash /path/to/scripts/reset-demo-data.sh

# Docker
docker exec callcenter-mysql bash /scripts/reset-demo-data.sh
```

### Cron Scheduling (Optional)

Add to crontab to run reset every 30 days:

```bash
# Reset on 1st of every month at midnight
0 0 1 * * /path/to/scripts/reset-demo-data.sh >> /var/log/vyomtech-demo-reset.log 2>&1

# Reset every 30 days from now
0 0 */30 * * /path/to/scripts/reset-demo-data.sh >> /var/log/vyomtech-demo-reset.log 2>&1
```

---

## Displaying on Login Page

### Add to Login Component

```tsx
// pages/login.tsx or components/LoginForm.tsx

import { DemoCredentials } from '@/app/demo-credentials';

export default function LoginPage() {
  return (
    <div className="login-container">
      <LoginForm />
      <DemoCredentials />  {/* Add this line */}
    </div>
  );
}
```

### Styling Notes

- Uses Tailwind CSS
- Blue-themed design
- Responsive (works on mobile/tablet)
- Accessible (proper contrast ratios)
- Auto-hides for non-demo environments

---

## Database Schema

### Demo Tenant Structure

```sql
-- Tenant
demo_vyomtech_001

-- Partners & Users
partners (id: 1-4)
partner_users (id: 1-4)
partner_leads (5 records)
partner_payouts (0, available for testing)

-- Agents & Leads
agent (4 records)
lead (5 records)

-- Campaigns
campaign (4 records)
campaign_recipient (auto-populated)

-- Projects
construction_projects (4 records)
progress_tracking (8 records)

-- Other
task (4 records)
call (4 records)
compliance_records (4 records)
gamification_points_history (4 records)
gamification_badges (4 records)
```

---

## Testing Scenarios

### 1. Portal Admin Testing
```
Login: demo@vyomtech.com / demo123
Test:  - View all leads
       - Create campaigns
       - Manage agents
       - View projects
```

### 2. Partner Multi-Tenancy
```
Login as each partner type
Verify:  - Tenant isolation
         - Data filtering
         - Correct permissions
```

### 3. Lead Conversion Flow
```
1. View demo leads
2. Create tasks for follow-up
3. Track progress
4. Monitor gamification points
```

### 4. Compliance & Audit
```
View RERA compliance records
View GST registration status
Check environmental clearances
Review audit logs
```

---

## Auto-Reset Details

### Deleted Tables (18)

On reset, data is deleted from:

```
partner_lead_credits
partner_leads
partner_payouts
partner_payout_details
partner_activities
partner_users
partners
gamification_points_history
gamification_badges
progress_tracking
compliance_records
task
call
campaign_recipient
campaign
lead
agent
construction_projects
```

### Reloaded Data

Fresh data is reloaded with:

- **Current timestamps**: All timestamps set to NOW()
- **Reset IDs**: Auto-increment IDs start fresh
- **Same structure**: Identical data relationships
- **No duplicates**: Uses INSERT IGNORE for safety

---

## Troubleshooting

### Demo Data Not Showing

1. Check if migration 025 was loaded:
   ```sql
   SELECT COUNT(*) FROM partners WHERE tenant_id = 'demo_vyomtech_001';
   ```

2. Verify demo tenant exists:
   ```sql
   SELECT * FROM tenants WHERE id = 'demo_vyomtech_001';
   ```

3. Check app logs for reset errors:
   ```bash
   docker logs callcenter-app | grep DemoReset
   ```

### Reset Not Triggering

1. Verify service is running:
   ```go
   // In main.go
   demoService := services.NewDemoResetService(db)
   demoService.StartScheduler()
   ```

2. Check system time (important for scheduler)

3. Review database transaction logs

### Cannot Login with Demo Credentials

1. Verify password hash matches:
   ```bash
   # bcrypt hash of "demo123"
   $2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm
   ```

2. Check if partner_users table has the credentials

3. Verify tenant ID matches in request header

---

## Future Enhancements

- [ ] UI dashboard showing last reset time
- [ ] Manual reset button in admin panel
- [ ] Email notification when reset completes
- [ ] Configurable reset interval (30 days → custom)
- [ ] Demo data templates for different industries
- [ ] Reset history audit trail
- [ ] Ability to add custom demo scenarios

---

## Support

For issues or questions:

1. Check app logs: `docker logs callcenter-app`
2. Check database: `SELECT * FROM tenants WHERE id = 'demo_vyomtech_001'`
3. Manual reset: `bash scripts/reset-demo-data.sh`

---

**Last Updated**: December 3, 2025
**Demo Tenant ID**: `demo_vyomtech_001`
**Reset Interval**: 30 days
**Status**: ✅ Production Ready
