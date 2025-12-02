// ========== SOURCE & PORTAL INTEGRATIONS ==========
export interface SourceIntegration {
  id?: string
  tenant_id?: string
  source_type: 'google' | 'meta' | 'portal' | 'website' | 'landing_page' | 'referral' | 'offline' | 'email' | 'sms' | 'whatsapp' | 'custom'
  source_name: string // 'Google Ads', 'Facebook', 'MagicBricks', 'Website', 'Broker Network', etc.
  sub_source?: string // 'SEM', '99acres_leads', 'Chat', 'Incoming Call', etc.
  api_key?: string // Encrypted
  is_active: boolean
  integration_status: 'configured' | 'syncing' | 'paused' | 'error'
  last_sync?: string
  sync_frequency: 'real_time' | 'hourly' | 'daily' | 'manual'
  webhook_url?: string
  lead_count?: number
  created_at?: string
}

export interface GoogleIntegration extends SourceIntegration {
  source_type: 'google'
  google_ads_account_id?: string
  google_analytics_property_id?: string
  conversion_tracking_enabled?: boolean
  auto_bidding_enabled?: boolean
}

export interface MetaIntegration extends SourceIntegration {
  source_type: 'meta'
  facebook_ad_account_id?: string
  instagram_business_account_id?: string
  pixel_id?: string
  conversion_api_enabled?: boolean
  form_leads_enabled?: boolean
}

export interface PortalIntegration extends SourceIntegration {
  source_type: 'portal'
  portal_name: 'magicbricks' | '99acres' | 'housing' | 'nobroker' | 'proptiger' | 'custom'
  portal_user_id?: string
  auto_lead_import?: boolean
  lead_assignment_logic?: 'round_robin' | 'custom' | 'manual'
}

export interface WebsiteIntegration extends SourceIntegration {
  source_type: 'website'
  website_url?: string
  form_tracking_enabled?: boolean
  chat_widget_id?: string
  call_tracking_enabled?: boolean
}

// Campaign Ideas - Planning Stage
export interface CampaignIdea {
  id?: string
  tenant_id?: string
  project_id: string
  project_name?: string
  title: string
  description?: string
  target_segment: 'residential' | 'commercial' | 'nri' | 'corporate' | 'investor' | 'all'
  budget_estimate: number
  expected_leads: number
  channels: ('email' | 'sms' | 'social' | 'website' | 'offline' | 'events' | 'broker_network' | 'google_ads' | 'facebook_ads' | 'portals' | 'landing_page' | 'whatsapp' | 'referral')[]
  integrations?: string[] // Integration IDs to use
  tags: string[] // Tag IDs or names
  status: 'draft' | 'approved' | 'in_execution' | 'completed' | 'archived'
  priority: 'low' | 'medium' | 'high' | 'critical'
  created_by?: string
  created_at?: string
  updated_at?: string
}

// Real Estate Marketing Campaign (Active)
export interface Campaign {
  id?: string
  tenant_id?: string
  idea_id?: string
  project_id: string
  project_name?: string
  name: string
  description?: string
  campaign_type: 'project_launch' | 'phase_launch' | 'grand_opening' | 'season_promotion' | 'inventory_clearance' | 'referral'
  status: 'draft' | 'active' | 'paused' | 'completed' | 'archived'
  start_date: string
  end_date: string
  budget: number
  spent?: number
  channels: ('email' | 'sms' | 'social' | 'website' | 'offline' | 'events' | 'broker_network' | 'google_ads' | 'facebook_ads' | 'portals' | 'landing_page' | 'whatsapp' | 'referral')[]
  target_segment: 'residential' | 'commercial' | 'nri' | 'corporate' | 'investor' | 'all'
  target_price_range_min: number
  target_price_range_max: number
  expected_leads: number
  expected_bookings: number
  // Integration sources
  active_integrations?: string[] // Integration IDs
  google_campaign_id?: string
  meta_campaign_id?: string
  portal_listing_ids?: Record<string, string> // portal_name -> listing_id
  landing_page_url?: string
  tags: string[] // Tag IDs
  created_at?: string
  updated_at?: string
}

// UTM Tracking Parameters
export interface UTMTracker {
  id?: string
  tenant_id?: string
  campaign_id: string
  integration_id?: string // Associated integration
  utm_source: string // 'google', 'facebook', 'instagram', 'email', 'direct', 'broker', 'magicbricks', 'website', 'landing_page', etc.
  utm_medium: string // 'cpc', 'organic', 'social', 'email', 'referral', 'organic_social', 'form', 'chat'
  utm_campaign: string // Campaign name
  utm_content: string // Ad variant/version (ad_v1, email_subject_a, etc.)
  utm_term?: string // Keyword (for search campaigns)
  full_url: string // Complete tracking URL
  short_url?: string // Shortened URL
  qr_code?: string // Base64 encoded QR code or URL
  // Source-specific tracking
  source_type?: 'google_ads' | 'facebook' | 'instagram' | 'portal' | 'email' | 'sms' | 'whatsapp' | 'website_form' | 'landing_page' | 'landing_page_form' | 'organic' | 'referral' | 'offline'
  portal_name?: string // 'magicbricks', '99acres', 'housing', 'nobroker', etc.
  landing_page_id?: string // For landing page campaigns
  email_campaign_id?: string // For email integrations
  phone_tracking_number?: string // For offline tracking
  clicks?: number
  impressions?: number
  conversions?: number
  created_at?: string
  updated_at?: string
}

// Lead Source Attribution
export interface LeadSourceTracking {
  id?: string
  tenant_id?: string
  lead_id: string
  utm_source?: string
  utm_medium?: string
  utm_campaign?: string
  utm_content?: string
  utm_term?: string
  referrer_url?: string
  landing_page?: string
  device_type?: 'mobile' | 'tablet' | 'desktop'
  browser?: string
  ip_address?: string
  tracked_at?: string
}

// Marketing Tags for Organization
export interface MarketingTag {
  id?: string
  tenant_id?: string
  name: string // 'seasonal_promo', 'vip_outreach', 'first_time_buyer'
  category: 'campaign_type' | 'lead_quality' | 'priority' | 'channel' | 'custom'
  color?: string // Hex color for UI
  description?: string
  is_system?: boolean // System tags vs user-created
  created_at?: string
}

// Real Estate Lead - Property Inquiries
export interface Lead {
  id?: string
  tenant_id?: string
  project_id: string
  project_name?: string
  prospect_name: string
  email: string
  phone: string
  alternate_phone?: string
  property_type_interest: 'residential' | 'commercial' | 'mixed_use' | 'any'
  budget_range_min: number
  budget_range_max: number
  required_bhk?: number
  super_area_range_min?: number
  super_area_range_max?: number
  source: 'website' | 'referral' | 'campaign' | 'broker' | 'social' | 'event' | 'walk_in' | 'cold_call' | 'other'
  lead_stage: 'inquiry' | 'interested' | 'site_visit_scheduled' | 'site_visit_done' | 'qualified' | 'negotiation' | 'lost' | 'converted'
  lead_quality: 'hot' | 'warm' | 'cold'
  assigned_to?: string // sales executive
  assigned_to_name?: string
  campaign_id?: string
  last_contact_date?: string
  follow_up_date?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

// Real Estate Marketing Performance
export interface MarketingCampaignPerformance {
  id?: string
  tenant_id?: string
  campaign_id: string
  leads_generated: number
  site_visits_scheduled: number
  site_visits_completed: number
  bookings_converted: number
  conversion_rate: number
  value_conversion: number // Total booking value
  roi: number
  cost_per_lead: number
  cost_per_booking: number
  by_utm_source?: Record<string, SourceMetrics>
  by_tag?: Record<string, SourceMetrics>
  created_at?: string
  updated_at?: string
}

export interface SourceMetrics {
  clicks: number
  leads: number
  bookings: number
  revenue: number
  cost_per_lead: number
  cost_per_booking: number
}

export interface MarketingMetrics {
  total_campaigns: number
  active_campaigns: number
  draft_ideas: number
  total_leads: number
  qualified_leads: number
  site_visits_scheduled: number
  site_visits_completed: number
  bookings_from_leads: number
  conversion_rate_inquiry_to_booking: number // %
  average_lead_quality_score: number
  total_budget: number
  total_spent: number
  avg_roi: number
  avg_cost_per_lead: number
  avg_cost_per_booking: number
  // UTM specific
  top_utm_source: string
  top_utm_medium: string
  top_utm_campaign: string
  utm_sources_count: number
  // Tag specific
  campaigns_by_tag_count: Record<string, number>
  top_performing_tag: string
}
