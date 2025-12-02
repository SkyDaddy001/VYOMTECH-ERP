# Multi-Source Lead Integration & UTM Tracking System

## Overview

This guide covers the comprehensive integration architecture for VYOM's real estate marketing platform, supporting all major paid advertising channels, real estate portals, organic sources, and owned channels.

---

## 1. Integration Architecture

### 1.1 Source Categories

#### **Paid Advertising Channels**
- **Google Ads** (SEM)
  - Search campaigns
  - Performance Max
  - Google Local Services Ads
  - Lead form extensions
  - Real-time conversion tracking

- **Meta Platform** (Social)
  - Facebook Lead Ads
  - Instagram Lead Forms
  - Facebook Conversion API
  - Pixel-based tracking
  - Video lead ads

#### **Real Estate Portals**
- **MagicBricks**
  - Lead form submissions
  - Inquiry tracking
  - Property page leads
  - Subvention leads

- **99 Acres**
  - Lead generation
  - Buyer inquiries
  - Phone inquiries
  - Direct form leads

- **Housing.com**
  - Buyer requirements matching
  - Project inquiry leads
  - Agent referral leads
  - Direct inbound leads

- **NoBroker**
  - Direct buyer leads
  - Self-curated leads
  - NoBroker Pro leads
  - Referral leads

- **PropTiger**
  - Property inquiry leads
  - Shortlist notifications
  - Call-based leads
  - Digital marketing leads

#### **Organic & Owned Channels**
- **Website**
  - Form submissions
  - Chat interactions
  - Phone inquiries
  - Google Analytics tracking

- **Landing Pages**
  - Campaign-specific landing pages
  - Property showcase pages
  - Phase launch pages
  - Promotional pages

- **Email Marketing**
  - Campaign responses
  - Link clicks
  - Newsletter signups
  - Re-engagement campaigns

- **SMS & WhatsApp**
  - Campaign responses
  - Click-to-call
  - Link tracking
  - Broadcast responses

#### **Referral & Offline**
- **Broker Network**
  - Channel partner leads
  - Individual agent leads
  - Network referrals
  - Cooperative leads

- **Offline Sources**
  - Site visits (STL)
  - Hoarding inquiries
  - Events & activations
  - Radio/Print inquiries
  - Incoming calls

---

## 2. Integration APIs & Methods

### 2.1 Google Ads Integration

```typescript
// Configuration
const googleIntegration: GoogleIntegration = {
  source_type: 'google',
  source_name: 'Google Ads - Residential',
  google_ads_account_id: 'xxx-xxx-xxxx',
  google_analytics_property_id: 'G-XXXXXXXXXX',
  conversion_tracking_enabled: true,
  auto_bidding_enabled: true,
  sync_frequency: 'real_time'
}

// Lead Sync Webhook
POST /api/v1/marketing/integrations/google/webhook
{
  "lead_id": "google_lead_123",
  "user_data": {
    "email": "customer@example.com",
    "phone": "+919876543210",
    "first_name": "Rajesh",
    "last_name": "Kumar"
  },
  "location_data": {
    "country_code": "IN",
    "state": "Tamil Nadu"
  },
  "form_id": "123456",
  "timestamp": "2025-11-29T10:30:00Z",
  "utm_parameters": {
    "utm_source": "google",
    "utm_medium": "cpc",
    "utm_campaign": "summer_residential_launch",
    "utm_content": "ad_v1",
    "utm_term": "2bhk_apartments_mumbai"
  }
}

// Get Leads from Google
GET /api/v1/marketing/integrations/google/leads?since=2025-11-29T00:00:00Z

// Track Conversions
POST /api/v1/marketing/integrations/google/conversion
{
  "lead_id": "google_lead_123",
  "conversion_type": "site_visit_scheduled",
  "conversion_value": 50000,
  "conversion_currency": "INR"
}
```

### 2.2 Meta (Facebook & Instagram) Integration

```typescript
// Configuration
const metaIntegration: MetaIntegration = {
  source_type: 'meta',
  source_name: 'Facebook & Instagram Ads',
  facebook_ad_account_id: 'act_123456789',
  instagram_business_account_id: '17841406338772046',
  pixel_id: '1234567890',
  conversion_api_enabled: true,
  form_leads_enabled: true,
  sync_frequency: 'real_time'
}

// Lead Sync Webhook
POST /api/v1/marketing/integrations/meta/webhook
{
  "form_id": "fb_lead_123",
  "ad_id": "123456789",
  "campaign_id": "987654321",
  "adset_id": "555555555",
  "created_time": "2025-11-29T10:30:00Z",
  "field_data": [
    { "name": "email", "value": "customer@example.com" },
    { "name": "phone_number", "value": "+919876543210" },
    { "name": "first_name", "value": "Priya" },
    { "name": "last_name", "value": "Singh" },
    { "name": "property_interested", "value": "2BHK" },
    { "name": "budget", "value": "50-75L" }
  ]
}

// Conversion Tracking (Server-to-Server)
POST /api/v1/marketing/integrations/meta/conversions
{
  "data": [{
    "event_name": "Purchase",  // or Lead, ViewContent, etc.
    "event_time": 1630000000,
    "user_data": {
      "em": "customer_email_hash",
      "ph": "customer_phone_hash"
    },
    "custom_data": {
      "value": 50000,
      "currency": "INR",
      "content_name": "2BHK_Project_X",
      "content_type": "property"
    }
  }]
}

// Get Page Insights
GET /api/v1/marketing/integrations/meta/insights?metric=impressions,clicks,leads
```

### 2.3 Portal Integration (MagicBricks, 99Acres, etc.)

```typescript
// Configuration
const portalIntegration: PortalIntegration = {
  source_type: 'portal',
  source_name: 'MagicBricks Integration',
  portal_name: 'magicbricks',
  portal_user_id: 'mb_user_123',
  auto_lead_import: true,
  lead_assignment_logic: 'round_robin',
  sync_frequency: 'real_time'
}

// Portal Lead Webhook (MagicBricks)
POST /api/v1/marketing/integrations/portal/webhook/magicbricks
{
  "lead_id": "mb_lead_123",
  "property_id": "mb_prop_456",
  "project_name": "Project Paradise",
  "lead_name": "Amit Verma",
  "lead_email": "amit@example.com",
  "lead_phone": "9876543210",
  "property_type": "2BHK",
  "budget_range": "30-50L",
  "lead_type": "inquiry",
  "created_at": "2025-11-29T10:30:00Z"
}

// Lead Export from Portal
GET /api/v1/marketing/integrations/portal/leads?portal=magicbricks&since=2025-11-29T00:00:00Z

// Portal Listing Management
POST /api/v1/marketing/integrations/portal/listings
{
  "portal_name": "magicbricks",
  "project_id": "proj_123",
  "property_listings": [
    {
      "property_id": "prop_1",
      "title": "2BHK Apartment",
      "price": 5000000,
      "area": 1000,
      "description": "Luxury 2BHK apartment in premium location"
    }
  ]
}
```

### 2.4 Website & Landing Page Integration

```typescript
// Configuration
const websiteIntegration: WebsiteIntegration = {
  source_type: 'website',
  source_name: 'Main Website',
  website_url: 'https://skylivingproperty.com',
  form_tracking_enabled: true,
  chat_widget_id: 'drift_widget_123',
  call_tracking_enabled: true,
  sync_frequency: 'real_time'
}

// Website Form Submission
POST /api/v1/marketing/integrations/website/form-submission
{
  "form_id": "contact_form_1",
  "form_name": "Property Inquiry",
  "form_page": "https://skylivingproperty.com/projects/project-1",
  "user_data": {
    "email": "visitor@example.com",
    "phone": "9876543210",
    "name": "Ravi Kumar",
    "property_interested": "3BHK",
    "budget": "75-100L"
  },
  "utm_parameters": {
    "utm_source": "google",
    "utm_medium": "organic",
    "utm_campaign": "organic_search"
  },
  "timestamp": "2025-11-29T10:30:00Z"
}

// Chat Lead Integration (Drift, Intercom, etc.)
POST /api/v1/marketing/integrations/website/chat-lead
{
  "chat_id": "chat_123",
  "visitor_email": "visitor@example.com",
  "visitor_name": "Sarah Johnson",
  "conversation_summary": "Interested in 2BHK apartments",
  "has_contact_info": true,
  "quality_score": 8
}

// Phone Call Integration (CallTrack, Twilio, etc.)
POST /api/v1/marketing/integrations/website/call-lead
{
  "call_id": "call_123",
  "caller_number": "+919876543210",
  "call_duration": 180,
  "recording_url": "https://...",
  "call_type": "inbound",
  "tracking_number": "+919123456789",
  "timestamp": "2025-11-29T10:30:00Z"
}
```

### 2.5 Email Marketing Integration

```typescript
// Configuration
const emailIntegration: SourceIntegration = {
  source_type: 'email',
  source_name: 'MailChimp Newsletter',
  api_key: 'encrypted_key_123',
  sync_frequency: 'daily'
}

// Email Campaign Link Tracking
POST /api/v1/marketing/integrations/email/click
{
  "email": "subscriber@example.com",
  "campaign_id": "campaign_123",
  "link_id": "link_456",
  "link_url": "https://skylivingproperty.com/project",
  "timestamp": "2025-11-29T10:30:00Z",
  "user_agent": "Mozilla/5.0...",
  "ip_address": "123.45.67.89"
}

// Email Signup/Unsubscribe
POST /api/v1/marketing/integrations/email/subscription-change
{
  "email": "subscriber@example.com",
  "event": "subscribe", // or "unsubscribe"
  "source_campaign": "weekly_updates",
  "timestamp": "2025-11-29T10:30:00Z"
}
```

### 2.6 SMS & WhatsApp Integration

```typescript
// Configuration
const smsIntegration: SourceIntegration = {
  source_type: 'sms',
  source_name: 'SMS Blast - Campaigns',
  api_key: 'encrypted_key_456',
  sync_frequency: 'real_time'
}

// SMS Response Tracking
POST /api/v1/marketing/integrations/sms/response
{
  "phone": "+919876543210",
  "message_content": "Interested in the 2BHK apartments. Please share details.",
  "inbound_message_id": "sms_123",
  "outbound_campaign_id": "campaign_789",
  "timestamp": "2025-11-29T10:30:00Z"
}

// WhatsApp Form Response
POST /api/v1/marketing/integrations/whatsapp/response
{
  "sender_number": "+919876543210",
  "message": "Yes, I'm interested in project details",
  "campaign_id": "whatsapp_campaign_1",
  "media": [
    { "type": "image", "url": "https://..." }
  ],
  "timestamp": "2025-11-29T10:30:00Z"
}
```

---

## 3. UTM Parameter Strategy

### 3.1 Standard UTM Parameters

```
https://skylivingproperty.com/projects/project-1?
  utm_source=google
  &utm_medium=cpc
  &utm_campaign=summer_2024_residential
  &utm_content=ad_version_1
  &utm_term=2bhk_apartments_mumbai
```

### 3.2 Real Estate Specific Parameters

**Extended UTM Mapping:**
- `utm_source`: Channel origin (google, facebook, magicbricks, website, email, etc.)
- `utm_medium`: Interaction type (cpc, organic, social, form, email, sms, referral, offline)
- `utm_campaign`: Campaign name (project_launch, phase_2_launch, seasonal_promo)
- `utm_content`: Variant ID (ad_v1, email_subject_a, hero_image_1)
- `utm_term`: Property/keyword focus (2bhk, luxury_apartments, mumbai_location)

**Portal-Specific UTM:**
```
Portal Lead URL Pattern:
https://magicbricks.com/project/project-123?
  utm_source=magicbricks
  &utm_medium=cpc
  &utm_campaign=sky_living_campaign_2024
  &utm_content=featured_listing
  &utm_term=residential_2bhk

Portal Lead Attribution:
- Source: magicbricks
- Medium: cpc (Cost Per Click)
- Property ID: From portal form
- Lead capture date: Automatic
```

### 3.3 UTM Generation Examples

**Google Ads Campaign:**
```
https://skylivingproperty.com/projects/paradise-heights?
utm_source=google
&utm_medium=cpc
&utm_campaign=paradise_heights_summer_2024
&utm_content=search_ad_responsive
&utm_term=luxury_apartments_pune
```

**Facebook Lead Ad:**
```
https://skylivingproperty.com/projects/paradise-heights?
utm_source=facebook
&utm_medium=social
&utm_campaign=paradise_heights_facebook
&utm_content=lead_form_2024
&utm_term=first_time_buyer
```

**MagicBricks Portal:**
```
https://magicbricks.com/project/paradise-heights-pune?
utm_source=magicbricks
&utm_medium=cpc
&utm_campaign=portal_lead_gen_2024
&utm_content=featured_listing_tier1
&utm_term=2bhk_pune
```

**Email Campaign:**
```
https://skylivingproperty.com/projects/paradise-heights?
utm_source=email
&utm_medium=email
&utm_campaign=monthly_newsletter_nov_2024
&utm_content=project_spotlight
&utm_term=newsletter_subscriber
```

**SMS Campaign:**
```
https://bit.ly/phts24?
utm_source=sms
&utm_medium=sms
&utm_campaign=sms_blast_nov_2024
&utm_content=exclusive_offer
&utm_term=existing_leads
```

---

## 4. Lead Attribution Model

### 4.1 First-Touch Attribution
Attributes conversion to the first source that brought the lead.

```typescript
// Example: Lead entered through Google Ads, converted via Email
Attribution → Google Ads (first source)
```

### 4.2 Last-Touch Attribution
Attributes conversion to the last source before booking.

```typescript
// Example: Lead entered through Google, clicked Email, then visited website
Attribution → Website (last source)
```

### 4.3 Multi-Touch Attribution
Distributes credit across multiple touchpoints.

```typescript
// Example: 30% Google + 50% Email + 20% Portal
Distribution model:
- First Touch: 30%
- Last Touch: 50%
- Middle Touches: 20%
```

### 4.4 Custom Attribution Logic

```typescript
// Priority-based attribution
Priority Order:
1. Paid (Google, Facebook, Portals)
2. Owned (Email, SMS, Website)
3. Referral
4. Offline

const attributeLead = (touchpoints: Touchpoint[]) => {
  const lastPaid = touchpoints
    .reverse()
    .find(t => ['google', 'facebook', 'portal'].includes(t.source))
  
  if (lastPaid) return lastPaid.source
  
  return touchpoints[touchpoints.length - 1].source
}
```

---

## 5. Real Estate Portal Sources

### 5.1 Portal Data Mapping

| Portal | Lead Type | Fields | Sync Method |
|--------|-----------|--------|------------|
| **MagicBricks** | Form inquiry | Name, Email, Phone, Property, Budget | Webhook |
| **99 Acres** | Buyer match | Name, Phone, Preferences | API pull |
| **Housing.com** | Lead form | Email, Phone, Location, Budget | Webhook |
| **NoBroker** | Direct inquiry | Name, Phone, Preferences | API pull |
| **PropTiger** | Property inquiry | Email, Phone, Shortlist, Budget | Webhook |

### 5.2 Portal Lead Processing

```typescript
interface PortalLead {
  portal_id: string
  portal_name: 'magicbricks' | '99acres' | 'housing' | 'nobroker' | 'proptiger'
  lead_id: string
  prospect_name: string
  email: string
  phone: string
  property_interest: {
    type: string
    budget_range: string
    location: string
  }
  source_tracking: {
    utm_source: string // portal name
    utm_medium: string // 'cpc' for portals
    utm_campaign: string // platform campaign
    referrer_url: string // portal URL
  }
  created_at: string
  imported_at: string
}

// Lead import workflow
1. Receive lead from portal webhook/API
2. Extract prospect information
3. Set utm_source = portal_name
4. Check for duplicates (email/phone match)
5. Assign to sales team (round-robin or custom logic)
6. Create follow-up task
7. Log in LeadSourceTracking
```

---

## 6. Lead Source Tagging

### 6.1 Automatic Tagging

```typescript
const AUTO_TAGS = {
  'google': ['paid', 'sem', 'search', 'high_intent'],
  'facebook': ['paid', 'social', 'display', 'awareness'],
  'instagram': ['paid', 'social', 'visual', 'engagement'],
  'magicbricks': ['portal', 'real_estate', 'property_portal'],
  'website': ['organic', 'direct', 'owned', 'first_party_data'],
  'email': ['owned', 'repeat_engagement', 'nurture'],
  'referral': ['word_of_mouth', 'organic', 'trusted_source'],
  'offline': ['offline', 'high_touch', 'local']
}

// Apply tags when lead is created
const tagLead = (lead: Lead, source: string) => {
  const tags = AUTO_TAGS[source] || ['unknown']
  lead.tags = [...new Set([...lead.tags, ...tags])]
}
```

### 6.2 Manual Tags

Users can add custom tags:
- **Lead Quality**: hot, warm, cold, junk
- **Property Interest**: luxury, affordable, mixed
- **Priority**: vip, regular, follow_up
- **Stage**: inquiry, negotiation, lost, converted
- **Campaign**: seasonal, promotional, launch, referral

---

## 7. Analytics & Reporting

### 7.1 Source Performance Dashboard

```typescript
// Metrics by source
interface SourceAnalytics {
  source_name: string
  leads_count: number
  qualified_leads: number
  site_visits_scheduled: number
  site_visits_completed: number
  bookings_count: number
  revenue: number
  cost_per_lead: number
  cost_per_booking: number
  roi: number
  conversion_rate: number
}

// Example metrics
{
  "google": {
    leads: 245,
    qualified: 180,
    bookings: 18,
    revenue: "9000000", // ₹90L
    cost_per_lead: 4081,
    cost_per_booking: 45000,
    roi: 1.8,
    conversion_rate: 0.073
  },
  "magicbricks": {
    leads: 156,
    qualified: 98,
    bookings: 8,
    revenue: "4000000", // ₹40L
    cost_per_lead: 2051,
    cost_per_booking: 25000,
    roi: 2.3,
    conversion_rate: 0.051
  },
  "website": {
    leads: 89,
    qualified: 67,
    bookings: 6,
    revenue: "3000000", // ₹30L
    cost_per_lead: 0,
    cost_per_booking: 0,
    roi: 100,
    conversion_rate: 0.067
  }
}
```

### 7.2 Multi-Channel Attribution Report

```typescript
// Sample monthly report
{
  "period": "2025-11-01 to 2025-11-30",
  "total_leads": 890,
  "total_qualified": 645,
  "total_bookings": 42,
  "total_revenue": "21000000", // ₹2.1Cr
  "by_source": {
    "google": { leads: 245, bookings: 18, revenue: "9000000" },
    "facebook": { leads: 156, bookings: 9, revenue: "4500000" },
    "magicbricks": { leads: 156, bookings: 8, revenue: "4000000" },
    "website": { leads: 89, bookings: 6, revenue: "3000000" },
    "email": { leads: 67, bookings: 3, revenue: "1500000" },
    "other": { leads: 177, bookings: 2, revenue: "1000000" }
  }
}
```

---

## 8. Implementation Checklist

### Phase 1: Core Setup
- [ ] Database schema for SourceIntegration, UTMTracker, LeadSourceTracking
- [ ] API endpoints for integration CRUD operations
- [ ] IntegrationForm & IntegrationList components
- [ ] UTMGenerator component
- [ ] Encryption for API keys

### Phase 2: Paid Channels
- [ ] Google Ads webhook integration
- [ ] Meta (Facebook/Instagram) webhook integration
- [ ] Real-time lead sync
- [ ] Conversion tracking

### Phase 3: Portal Integration
- [ ] MagicBricks API integration
- [ ] 99 Acres API integration
- [ ] Housing.com webhook
- [ ] NoBroker API integration
- [ ] PropTiger webhook

### Phase 4: Organic Sources
- [ ] Website form tracking
- [ ] Landing page analytics
- [ ] Chat widget integration
- [ ] Call tracking integration

### Phase 5: Owned Channels
- [ ] Email platform integration (MailChimp, SendGrid, etc.)
- [ ] SMS provider integration (Twilio, MSG91, etc.)
- [ ] WhatsApp Business API integration
- [ ] Link click tracking

### Phase 6: Analytics & Reporting
- [ ] Source analytics dashboard
- [ ] Multi-touch attribution model
- [ ] ROI calculation engine
- [ ] Custom report builder

---

## 9. Best Practices

1. **Always track source of truth**: UTM parameters should match your integration records
2. **Use consistent naming**: Keep campaign names standardized across platforms
3. **Monitor data quality**: Check for duplicate leads and invalid data
4. **Set up alerts**: Monitor integration health and sync failures
5. **Regular audits**: Monthly source performance reviews
6. **Test webhooks**: Verify all webhook integrations work before going live
7. **Document custom logic**: Keep attribution rules documented
8. **Privacy first**: Ensure compliance with data privacy regulations (GDPR, CCPA)

---

## 10. Troubleshooting

### Common Issues

**Issue**: Leads from Google not syncing
- Check API credentials
- Verify webhook URL is accessible
- Check Google Ads campaign settings
- Monitor error logs

**Issue**: Duplicate leads from portals
- Enable duplicate detection
- Check email/phone matching logic
- Investigate portal API for duplicates

**Issue**: Missing UTM parameters
- Verify tracking URLs are correct
- Check form submission handlers
- Monitor UTM parsing logic

**Issue**: Low attribution accuracy
- Review touchpoint collection
- Adjust attribution model
- Check timezone settings

---

This comprehensive guide ensures full integration coverage across all major lead sources, enabling complete visibility into marketing performance and accurate lead attribution across all channels.
