package models

import "time"

// ==================== ANALYTICS MODELS ====================

// AnalyticsEvent represents a tracked event in the system
type AnalyticsEvent struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	EventType string    `db:"event_type" json:"event_type"`
	EventData string    `db:"event_data" json:"event_data"` // JSON
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ConversionFunnel tracks lead conversion through stages
type ConversionFunnel struct {
	ID             int64      `db:"id" json:"id"`
	TenantID       string     `db:"tenant_id" json:"tenant_id"`
	LeadID         int64      `db:"lead_id" json:"lead_id"`
	CampaignID     int64      `db:"campaign_id" json:"campaign_id"`
	Stage          string     `db:"stage" json:"stage"`
	ConversionRate float64    `db:"conversion_rate" json:"conversion_rate"`
	TimeInStage    int64      `db:"time_in_stage" json:"time_in_stage"` // seconds
	EnteredAt      time.Time  `db:"entered_at" json:"entered_at"`
	ExitedAt       *time.Time `db:"exited_at" json:"exited_at"`
	ConvertedAt    *time.Time `db:"converted_at" json:"converted_at"`
}

// AgentMetrics tracks individual agent performance
type AgentMetrics struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        string    `db:"tenant_id" json:"tenant_id"`
	AgentID         int64     `db:"agent_id" json:"agent_id"`
	MetricDate      time.Time `db:"metric_date" json:"metric_date"`
	CallsHandled    int64     `db:"calls_handled" json:"calls_handled"`
	AverageCallTime int64     `db:"average_call_time" json:"average_call_time"` // seconds
	LeadsConverted  int64     `db:"leads_converted" json:"leads_converted"`
	ConversionRate  float64   `db:"conversion_rate" json:"conversion_rate"`
	CustomerRating  float64   `db:"customer_rating" json:"customer_rating"`
	TasksCompleted  int64     `db:"tasks_completed" json:"tasks_completed"`
	AvailableTime   int64     `db:"available_time" json:"available_time"` // seconds
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// CampaignMetrics tracks campaign performance
type CampaignMetrics struct {
	ID               int64     `db:"id" json:"id"`
	TenantID         string    `db:"tenant_id" json:"tenant_id"`
	CampaignID       int64     `db:"campaign_id" json:"campaign_id"`
	MetricDate       time.Time `db:"metric_date" json:"metric_date"`
	LeadsGenerated   int64     `db:"leads_generated" json:"leads_generated"`
	LeadsContacted   int64     `db:"leads_contacted" json:"leads_contacted"`
	LeadsConverted   int64     `db:"leads_converted" json:"leads_converted"`
	ConversionRate   float64   `db:"conversion_rate" json:"conversion_rate"`
	AverageLeadValue float64   `db:"average_lead_value" json:"average_lead_value"`
	TotalRevenue     float64   `db:"total_revenue" json:"total_revenue"`
	ROI              float64   `db:"roi" json:"roi"`
	CostPerLead      float64   `db:"cost_per_lead" json:"cost_per_lead"`
	Cost             float64   `db:"cost" json:"cost"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// DailyReport aggregated daily statistics
type DailyReport struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        string    `db:"tenant_id" json:"tenant_id"`
	ReportDate      time.Time `db:"report_date" json:"report_date"`
	TotalCalls      int64     `db:"total_calls" json:"total_calls"`
	TotalLeads      int64     `db:"total_leads" json:"total_leads"`
	ConvertedLeads  int64     `db:"converted_leads" json:"converted_leads"`
	ConversionRate  float64   `db:"conversion_rate" json:"conversion_rate"`
	AverageCallTime int64     `db:"average_call_time" json:"average_call_time"`
	TasksCompleted  int64     `db:"tasks_completed" json:"tasks_completed"`
	ActiveAgents    int64     `db:"active_agents" json:"active_agents"`
	TotalRevenue    float64   `db:"total_revenue" json:"total_revenue"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// CustomReport user-defined report configuration
type CustomReport struct {
	ID          int64      `db:"id" json:"id"`
	TenantID    string     `db:"tenant_id" json:"tenant_id"`
	CreatedBy   int64      `db:"created_by" json:"created_by"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	Metrics     string     `db:"metrics" json:"metrics"`       // JSON array
	Filters     string     `db:"filters" json:"filters"`       // JSON object
	Schedule    string     `db:"schedule" json:"schedule"`     // cron expression
	Format      string     `db:"format" json:"format"`         // csv, pdf, json
	Recipients  string     `db:"recipients" json:"recipients"` // JSON array of emails
	Enabled     bool       `db:"enabled" json:"enabled"`
	LastRunAt   *time.Time `db:"last_run_at" json:"last_run_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// ReportExecution tracks report generation
type ReportExecution struct {
	ID          int64      `db:"id" json:"id"`
	TenantID    string     `db:"tenant_id" json:"tenant_id"`
	ReportID    int64      `db:"report_id" json:"report_id"`
	Status      string     `db:"status" json:"status"` // pending, processing, completed, failed
	FileURL     *string    `db:"file_url" json:"file_url"`
	Error       *string    `db:"error" json:"error"`
	StartedAt   time.Time  `db:"started_at" json:"started_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
}

// AnalyticsQuery represents a user analytics search
type AnalyticsQuery struct {
	StartDate  time.Time              `json:"start_date"`
	EndDate    time.Time              `json:"end_date"`
	MetricType string                 `json:"metric_type"` // agent, campaign, funnel, custom
	Filters    map[string]interface{} `json:"filters"`
	GroupBy    string                 `json:"group_by"` // date, agent, campaign
	Limit      int                    `json:"limit"`
	Offset     int                    `json:"offset"`
}

// AnalyticsResponse represents aggregated analytics data
type AnalyticsResponse struct {
	Period          string       `json:"period"`
	StartDate       time.Time    `json:"start_date"`
	EndDate         time.Time    `json:"end_date"`
	Data            interface{}  `json:"data"`
	Summary         SummaryStats `json:"summary"`
	Trends          []TrendPoint `json:"trends"`
	Comparison      interface{}  `json:"comparison,omitempty"`
	Recommendations []string     `json:"recommendations,omitempty"`
}

// SummaryStats contains aggregated metrics
type SummaryStats struct {
	TotalCalls       int64   `json:"total_calls"`
	TotalLeads       int64   `json:"total_leads"`
	ConvertedLeads   int64   `json:"converted_leads"`
	ConversionRate   float64 `json:"conversion_rate"`
	AverageCallTime  int64   `json:"average_call_time"`
	AverageLeadValue float64 `json:"average_lead_value"`
	TotalRevenue     float64 `json:"total_revenue"`
	ROI              float64 `json:"roi"`
}

// TrendPoint represents a single data point in a trend
type TrendPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Label     string    `json:"label"`
	Change    float64   `json:"change"` // percent change from previous
}

// PerformanceComparison compares metrics across periods
type PerformanceComparison struct {
	Metric     string  `json:"metric"`
	Current    float64 `json:"current"`
	Previous   float64 `json:"previous"`
	Change     float64 `json:"change"`      // absolute
	ChangeRate float64 `json:"change_rate"` // percentage
	Status     string  `json:"status"`      // up, down, flat
}

// FunnelStage represents a stage in conversion funnel
type FunnelStage struct {
	Name        string  `json:"name"`
	Count       int64   `json:"count"`
	Rate        float64 `json:"rate"`         // percentage of previous stage
	AverageTime int64   `json:"average_time"` // seconds
	Completed   int64   `json:"completed"`
	Dropped     int64   `json:"dropped"`
}

// FunnelAnalysis represents full funnel analysis
type FunnelAnalysis struct {
	CampaignID  int64         `json:"campaign_id"`
	Stages      []FunnelStage `json:"stages"`
	TotalEntry  int64         `json:"total_entry"`
	TotalExit   int64         `json:"total_exit"`
	OverallRate float64       `json:"overall_rate"`
}

// DashboardWidget represents a widget on the analytics dashboard
type DashboardWidget struct {
	ID         int64     `db:"id" json:"id"`
	TenantID   string    `db:"tenant_id" json:"tenant_id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	Title      string    `db:"title" json:"title"`
	WidgetType string    `db:"widget_type" json:"widget_type"` // chart, table, metric, etc
	Config     string    `db:"config" json:"config"`           // JSON configuration
	Position   int       `db:"position" json:"position"`
	Size       string    `db:"size" json:"size"` // small, medium, large
	Enabled    bool      `db:"enabled" json:"enabled"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
