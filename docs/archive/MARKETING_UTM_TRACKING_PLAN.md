# Real Estate Marketing Module - Idea Tracking & UTM Analytics Plan

## Overview
Enhanced marketing module with:
- Campaign Ideas → Execution tracking
- UTM parameter generation & tracking
- Lead source attribution
- Tag-based organization
- Multi-channel analytics

---

## 1. Core Entities

### Campaign Ideas (Planning Stage)
```typescript
interface CampaignIdea {
  id: string
  project_id: string
  title: string
  description: string
  target_segment: 'residential' | 'commercial' | 'nri' | 'corporate' | 'investor'
  budget_estimate: number
  expected_leads: number
  channels: string[] // email, sms, social, website, offline, broker, events
  tags: string[] // seasonal, urgent, vip, niche, bulk
  status: 'draft' | 'approved' | 'in_execution' | 'completed' | 'archived'
  priority: 'low' | 'medium' | 'high' | 'critical'
  created_by: string
  created_at: string
  updated_at: string
}
```

### UTM Tracking
```typescript
interface UTMTracker {
  id: string
  campaign_id: string
  utm_source: string // 'google', 'facebook', 'email', 'direct', etc.
  utm_medium: string // 'cpc', 'organic', 'social', 'email', 'referral'
  utm_campaign: string // campaign name
  utm_content: string // ad variant/version
  utm_term: string // keyword (if applicable)
  full_url: string // complete tracking URL
  short_url?: string // shortened URL
  qr_code?: string // base64 encoded QR code
  created_at: string
}
```

### Lead Source Tracking
```typescript
interface LeadSourceTracking {
  id: string
  lead_id: string
  utm_source: string
  utm_medium: string
  utm_campaign: string
  utm_content?: string
  utm_term?: string
  referrer_url?: string
  landing_page?: string
  device_type: 'mobile' | 'tablet' | 'desktop'
  browser?: string
  ip_address?: string
  tracked_at: string
}
```

### Marketing Tags
```typescript
interface MarketingTag {
  id: string
  name: string // 'seasonal_promo', 'vip_outreach', 'first_time_buyer', etc.
  category: 'campaign_type' | 'lead_quality' | 'priority' | 'channel' | 'custom'
  color: string // hex color for UI
  description?: string
  is_system: boolean
  created_at: string
}
```

### Campaign Performance Analytics
```typescript
interface CampaignAnalytics {
  campaign_id: string
  total_clicks: number
  total_impressions?: number
  ctr: number // click-through rate
  leads_generated: number
  leads_qualified: number
  bookings_converted: number
  conversion_rate: number
  cost_per_lead: number
  cost_per_booking: number
  revenue_generated: number
  roi: number
  by_source: {
    [key: string]: SourceMetrics
  }
  tracked_at: string
}

interface SourceMetrics {
  clicks: number
  leads: number
  bookings: number
  revenue: number
}
```

---

## 2. Campaign Workflow

```
IDEA CREATION
    ↓
IDEA APPROVAL (Management review)
    ↓
UTM SETUP & TRACKING LINK GENERATION
    ↓
CAMPAIGN EXECUTION
    ↓
LEAD CAPTURE & SOURCE ATTRIBUTION
    ↓
ANALYTICS & REPORTING
    ↓
CAMPAIGN COMPLETION
```

---

## 3. Core Features

### A. Idea Management
- [ ] Create campaign ideas with budget, target segment, channels
- [ ] Assign tags for organization
- [ ] Priority/status tracking
- [ ] Bulk edit ideas
- [ ] Archive old ideas
- [ ] Approval workflow

### B. UTM Generation & Management
- [ ] Auto-generate UTM parameters from campaign details
- [ ] Create multiple variants with different content (A/B testing)
- [ ] Generate QR codes for offline campaigns
- [ ] Create short URLs for sharing
- [ ] Preview tracking URLs
- [ ] Copy to clipboard functionality
- [ ] Export UTM manifest

### C. Lead Attribution
- [ ] Automatic source attribution via UTM params
- [ ] Manual source override
- [ ] Lead source analytics dashboard
- [ ] Cross-channel attribution
- [ ] First-touch vs Last-touch models

### D. Tags & Organization
- [ ] System tags (seasonal, vip, urgent, bulk, niche)
- [ ] Custom tags
- [ ] Multi-tag filtering
- [ ] Tag-based reporting
- [ ] Tag auto-assignment based on campaign type

### E. Analytics & Reporting
- [ ] Campaign performance dashboard
- [ ] UTM source performance breakdown
- [ ] Lead quality by source
- [ ] Conversion funnel by channel
- [ ] Cost per lead/booking trends
- [ ] ROI calculation
- [ ] Heatmap visualization (if possible)

---

## 4. Database Schema

### New Tables
```sql
-- Campaign Ideas
CREATE TABLE campaign_ideas (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  project_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  target_segment VARCHAR(50),
  budget_estimate DECIMAL(12,2),
  expected_leads INT,
  channels TEXT[], -- JSON array
  tags TEXT[], -- JSON array
  status VARCHAR(20),
  priority VARCHAR(20),
  created_by UUID,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (project_id) REFERENCES projects(id),
  FOREIGN KEY (created_by) REFERENCES users(id)
);

-- UTM Trackers
CREATE TABLE utm_trackers (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  campaign_id UUID NOT NULL,
  utm_source VARCHAR(100) NOT NULL,
  utm_medium VARCHAR(100) NOT NULL,
  utm_campaign VARCHAR(255) NOT NULL,
  utm_content VARCHAR(255),
  utm_term VARCHAR(255),
  full_url TEXT,
  short_url VARCHAR(255),
  qr_code LONGTEXT,
  created_at TIMESTAMP,
  FOREIGN KEY (campaign_id) REFERENCES campaign_ideas(id)
);

-- Lead Source Tracking
CREATE TABLE lead_source_tracking (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  lead_id UUID NOT NULL,
  utm_source VARCHAR(100),
  utm_medium VARCHAR(100),
  utm_campaign VARCHAR(255),
  utm_content VARCHAR(255),
  utm_term VARCHAR(255),
  referrer_url TEXT,
  landing_page VARCHAR(255),
  device_type VARCHAR(20),
  browser VARCHAR(100),
  ip_address VARCHAR(45),
  tracked_at TIMESTAMP,
  FOREIGN KEY (lead_id) REFERENCES leads(id)
);

-- Marketing Tags
CREATE TABLE marketing_tags (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  name VARCHAR(100) NOT NULL UNIQUE,
  category VARCHAR(50),
  color VARCHAR(7),
  description TEXT,
  is_system BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP
);

-- Campaign Analytics
CREATE TABLE campaign_analytics (
  id UUID PRIMARY KEY,
  tenant_id UUID NOT NULL,
  campaign_id UUID NOT NULL,
  total_clicks INT DEFAULT 0,
  total_impressions INT DEFAULT 0,
  leads_generated INT DEFAULT 0,
  leads_qualified INT DEFAULT 0,
  bookings_converted INT DEFAULT 0,
  revenue_generated DECIMAL(12,2) DEFAULT 0,
  analytics_json JSONB, -- Detailed metrics by source
  tracked_at TIMESTAMP,
  FOREIGN KEY (campaign_id) REFERENCES campaign_ideas(id)
);
```

---

## 5. Frontend Components

### Pages
- `/dashboard/marketing/ideas` - Campaign idea management
- `/dashboard/marketing/campaigns` - Active campaigns with UTM tracking
- `/dashboard/marketing/leads` - Lead source attribution
- `/dashboard/marketing/analytics` - Multi-channel analytics

### Components
- `IdeaForm.tsx` - Create/edit campaign ideas
- `IdeaList.tsx` - List all ideas with filters
- `UTMGenerator.tsx` - Generate UTM parameters and tracking links
- `UTMTrackingDashboard.tsx` - View all UTM links for a campaign
- `LeadSourceAttributionForm.tsx` - Manual source assignment
- `AnalyticsDashboard.tsx` - Multi-chart analytics view
- `TagManager.tsx` - Create and manage tags
- `CampaignAnalyticsCard.tsx` - Individual campaign stats

---

## 6. API Endpoints

### Campaign Ideas
```
POST   /api/v1/marketing/ideas              - Create idea
GET    /api/v1/marketing/ideas              - List ideas (with filters)
GET    /api/v1/marketing/ideas/{id}         - Get idea details
PUT    /api/v1/marketing/ideas/{id}         - Update idea
DELETE /api/v1/marketing/ideas/{id}         - Delete idea
PATCH  /api/v1/marketing/ideas/{id}/status  - Update status
PATCH  /api/v1/marketing/ideas/{id}/tags    - Update tags
```

### UTM Tracking
```
POST   /api/v1/marketing/utm-trackers              - Create UTM tracker
GET    /api/v1/marketing/utm-trackers              - List trackers
GET    /api/v1/marketing/campaigns/{id}/utm        - Get campaign UTMs
POST   /api/v1/marketing/utm-trackers/{id}/qr      - Generate QR code
POST   /api/v1/marketing/utm-trackers/{id}/shorten - Shorten URL
```

### Lead Attribution
```
POST   /api/v1/marketing/lead-source/{id}         - Set lead source
GET    /api/v1/marketing/lead-sources              - List attributions
GET    /api/v1/marketing/lead-sources/by-campaign  - Sources by campaign
```

### Tags
```
GET    /api/v1/marketing/tags                     - List all tags
POST   /api/v1/marketing/tags                     - Create tag
PUT    /api/v1/marketing/tags/{id}                - Update tag
DELETE /api/v1/marketing/tags/{id}                - Delete tag
```

### Analytics
```
GET    /api/v1/marketing/analytics/{campaign_id}  - Campaign analytics
GET    /api/v1/marketing/analytics/by-source      - Analytics by source
GET    /api/v1/marketing/roi-report               - ROI calculations
```

---

## 7. Implementation Phases

### Phase 1: Core Types & Services
- [ ] Define all TypeScript interfaces
- [ ] Create API service methods
- [ ] Setup service layer

### Phase 2: Campaign Ideas
- [ ] IdeaForm & IdeaList components
- [ ] CRUD operations
- [ ] Tag assignment
- [ ] Status workflow

### Phase 3: UTM Tracking
- [ ] UTM parameter generation logic
- [ ] QR code generation
- [ ] URL shortening integration
- [ ] UTMTrackingDashboard component

### Phase 4: Lead Attribution
- [ ] Lead source form
- [ ] Automatic UTM attribution
- [ ] Manual override
- [ ] Source tracking component

### Phase 5: Analytics
- [ ] Campaign analytics aggregation
- [ ] Multi-source analytics
- [ ] ROI calculations
- [ ] Visualization dashboard

### Phase 6: Integration & Testing
- [ ] Connect all components
- [ ] End-to-end testing
- [ ] Performance optimization
- [ ] Mobile responsiveness

---

## 8. UTM Best Practices

### Standard Parameters
- `utm_source` - Where the traffic came from (google, facebook, email, direct)
- `utm_medium` - How they clicked (cpc, organic, email, social, referral)
- `utm_campaign` - Campaign name (project_launch_q4, referral_bonus)
- `utm_content` - Which ad/email variant (ad_v1, email_subject_a)
- `utm_term` - Keywords (only for paid search)

### Example URLs
```
https://project.com/property?utm_source=google&utm_medium=cpc&utm_campaign=project_launch&utm_content=ad_v1
https://project.com/property?utm_source=facebook&utm_medium=social&utm_campaign=retargeting&utm_content=carousel_ad
https://project.com/property?utm_source=newsletter&utm_medium=email&utm_campaign=weekly&utm_content=feature_property
```

---

## 9. Key Metrics to Track

1. **Acquisition Metrics**
   - Leads by source/medium/campaign
   - Cost per lead
   - Lead quality score

2. **Conversion Metrics**
   - Site visits to leads
   - Leads to qualified
   - Qualified to booking
   - Booking to completion

3. **Channel Performance**
   - Cost per booking by channel
   - Time to conversion by source
   - Customer lifetime value by source

4. **Campaign ROI**
   - Campaign spend
   - Revenue generated
   - ROI percentage
   - Break-even analysis

---

## 10. Tags Strategy

### System Tags
- `seasonal_promo` - Seasonal promotions
- `vip_outreach` - VIP customer targeting
- `first_time_buyer` - New buyers
- `bulk_inquiry` - Group inquiries
- `urgent` - High priority
- `niche_market` - Specialized targeting

### Dynamic Tags
- Auto-created based on campaign properties
- User-defined custom tags
- Tag-based filtering & reporting

---

## Implementation Ready!

This plan ensures:
✅ Complete idea-to-execution tracking
✅ UTM-based lead attribution
✅ Multi-channel analytics capability
✅ Tag-based organization
✅ ROI & performance measurement
✅ Scalable architecture
✅ Real estate specific metrics
