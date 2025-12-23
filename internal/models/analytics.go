package models

import (
	"encoding/json"
	"time"
)

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

// ============================================
// Phase 3.1: Advanced Analytics Models
// ============================================

// AnalyticsDashboard represents a dashboard configuration
type AnalyticsDashboard struct {
	ID              int64            `json:"id"`
	TenantID        int64            `json:"tenant_id"`
	UserID          int64            `json:"user_id"`
	DashboardName   string           `json:"dashboard_name"`
	Description     *string          `json:"description"`
	DashboardType   string           `json:"dashboard_type"`
	LayoutConfig    *json.RawMessage `json:"layout_config"`
	IsPublic        bool             `json:"is_public"`
	IsDefault       bool             `json:"is_default"`
	WidgetIDs       *json.RawMessage `json:"widget_ids"`
	RefreshInterval int              `json:"refresh_interval"`
	DataDateRange   *string          `json:"data_date_range"`
	LastRefreshed   *time.Time       `json:"last_refreshed_at"`
	CreatedBy       *int64           `json:"created_by"`
	UpdatedBy       *int64           `json:"updated_by"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at"`
}

// AnalyticsWidget represents a widget for dashboards
type AnalyticsWidget struct {
	ID                int64            `json:"id"`
	TenantID          int64            `json:"tenant_id"`
	WidgetName        string           `json:"widget_name"`
	WidgetType        string           `json:"widget_type"`
	MetricSource      *string          `json:"metric_source"`
	MetricKey         *string          `json:"metric_key"`
	AggregationType   *string          `json:"aggregation_type"`
	TimePeriod        *string          `json:"time_period"`
	ComparisonEnabled bool             `json:"comparison_enabled"`
	ComparisonPeriod  *string          `json:"comparison_period"`
	DrillDownEnabled  bool             `json:"drill_down_enabled"`
	DrillDownTarget   *string          `json:"drill_down_target"`
	WidgetConfig      *json.RawMessage `json:"widget_config"`
	RefreshInterval   int              `json:"refresh_interval"`
	CreatedBy         *int64           `json:"created_by"`
	UpdatedBy         *int64           `json:"updated_by"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at"`
}

// AnalyticsKPI represents a key performance indicator
type AnalyticsKPI struct {
	ID                int64      `json:"id"`
	TenantID          int64      `json:"tenant_id"`
	KPIName           string     `json:"kpi_name"`
	Description       *string    `json:"description"`
	Category          string     `json:"category"`
	MetricSource      *string    `json:"metric_source"`
	CalculationMethod *string    `json:"calculation_method"`
	TargetValue       *float64   `json:"target_value"`
	CurrentValue      *float64   `json:"current_value"`
	PreviousValue     *float64   `json:"previous_value"`
	WarningThreshold  *float64   `json:"threshold_warning"`
	CriticalThreshold *float64   `json:"threshold_critical"`
	Status            string     `json:"status"`
	PercentageChange  *float64   `json:"percentage_change"`
	Trend             string     `json:"trend"`
	TimePeriod        string     `json:"time_period"`
	MeasurementDate   *time.Time `json:"measurement_date"`
	LastUpdated       *time.Time `json:"last_updated_at"`
	OwnerID           *int64     `json:"owner_id"`
	CreatedBy         *int64     `json:"created_by"`
	UpdatedBy         *int64     `json:"updated_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// AnalyticsReport represents a report definition
type AnalyticsReport struct {
	ID                int64            `json:"id"`
	TenantID          int64            `json:"tenant_id"`
	ReportName        string           `json:"report_name"`
	Description       *string          `json:"description"`
	ReportType        string           `json:"report_type"`
	Category          *string          `json:"category"`
	Template          *string          `json:"template"`
	DataSource        *string          `json:"data_source"`
	Filters           *json.RawMessage `json:"filters"`
	Grouping          *json.RawMessage `json:"grouping"`
	Sorting           *json.RawMessage `json:"sorting"`
	Metrics           *json.RawMessage `json:"metrics"`
	Dimensions        *json.RawMessage `json:"dimensions"`
	Format            string           `json:"format"`
	TimePeriod        string           `json:"time_period"`
	ScheduleFrequency *string          `json:"schedule_frequency"`
	ScheduleDay       *int             `json:"schedule_day"`
	ScheduleTime      *string          `json:"schedule_time"`
	IsPublic          bool             `json:"is_public"`
	IsScheduled       bool             `json:"is_scheduled"`
	LastGenerated     *time.Time       `json:"last_generated_at"`
	NextScheduled     *time.Time       `json:"next_scheduled_run"`
	RecipientEmails   *json.RawMessage `json:"recipient_emails"`
	CreatedBy         *int64           `json:"created_by"`
	UpdatedBy         *int64           `json:"updated_by"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at"`
}

// AnalyticsReportInstance represents a generated report
type AnalyticsReportInstance struct {
	ID              int64            `json:"id"`
	TenantID        int64            `json:"tenant_id"`
	ReportID        int64            `json:"report_id"`
	ReportName      string           `json:"report_name"`
	ReportType      *string          `json:"report_type"`
	GenerationStart *time.Time       `json:"generation_start_time"`
	GenerationEnd   *time.Time       `json:"generation_end_time"`
	DurationSeconds *int             `json:"generation_duration_seconds"`
	Status          string           `json:"status"`
	FilePath        *string          `json:"file_path"`
	FileSize        *int64           `json:"file_size"`
	FileFormat      *string          `json:"file_format"`
	S3Bucket        *string          `json:"s3_bucket"`
	S3Key           *string          `json:"s3_key"`
	RecordCount     *int64           `json:"record_count"`
	DataFromDate    *time.Time       `json:"data_from_date"`
	DataToDate      *time.Time       `json:"data_to_date"`
	GeneratedBy     *int64           `json:"generated_by"`
	DownloadCount   int              `json:"download_count"`
	LastDownloaded  *time.Time       `json:"last_downloaded_at"`
	ErrorMessage    *string          `json:"error_message"`
	Metadata        *json.RawMessage `json:"metadata"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at"`
}

// AnalyticsAlert represents an analytics alert
type AnalyticsAlert struct {
	ID                   int64            `json:"id"`
	TenantID             int64            `json:"tenant_id"`
	AlertName            string           `json:"alert_name"`
	AlertType            string           `json:"alert_type"`
	MetricName           *string          `json:"metric_name"`
	Condition            string           `json:"condition"`
	ThresholdValue       *float64         `json:"threshold_value"`
	ThresholdLower       *float64         `json:"threshold_lower"`
	ThresholdUpper       *float64         `json:"threshold_upper"`
	Severity             string           `json:"severity"`
	Enabled              bool             `json:"enabled"`
	AlertFrequency       string           `json:"alert_frequency"`
	NotificationChannels *json.RawMessage `json:"notification_channels"`
	RecipientEmails      *json.RawMessage `json:"recipient_emails"`
	RecipientUserIDs     *json.RawMessage `json:"recipient_user_ids"`
	LastTriggered        *time.Time       `json:"last_triggered_at"`
	TriggerCount         int              `json:"trigger_count"`
	CreatedBy            *int64           `json:"created_by"`
	UpdatedBy            *int64           `json:"updated_by"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	DeletedAt            *time.Time       `json:"deleted_at"`
}

// AnalyticsUserPreferences stores user analytics preferences
type AnalyticsUserPreferences struct {
	ID                  int64            `json:"id"`
	TenantID            int64            `json:"tenant_id"`
	UserID              int64            `json:"user_id"`
	DefaultDashboardID  *int64           `json:"default_dashboard_id"`
	DefaultReportID     *int64           `json:"default_report_id"`
	PreferredTimeZone   *string          `json:"preferred_time_zone"`
	PreferredCurrency   *string          `json:"preferred_currency"`
	DateFormat          *string          `json:"date_format"`
	NumberFormat        *string          `json:"number_format"`
	AutoRefreshEnabled  bool             `json:"auto_refresh_enabled"`
	AutoRefreshInterval int              `json:"auto_refresh_interval"`
	EmailNotifications  bool             `json:"email_notifications_enabled"`
	ReportFrequency     *string          `json:"report_frequency"`
	ThemePreference     string           `json:"theme_preference"`
	ChartPreferences    *json.RawMessage `json:"chart_preferences"`
	SavedFilters        *json.RawMessage `json:"saved_filters"`
	ExportFormatDefault *string          `json:"export_format_default"`
	CreatedBy           *int64           `json:"created_by"`
	UpdatedBy           *int64           `json:"updated_by"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
	DeletedAt           *time.Time       `json:"deleted_at"`
}

// ============================================
// DTO Models for Request/Response Handling
// ============================================

// CreateAnalyticsDashboardRequest represents a request to create a dashboard
type CreateAnalyticsDashboardRequest struct {
	DashboardName   string           `json:"dashboard_name" binding:"required"`
	Description     *string          `json:"description"`
	DashboardType   string           `json:"dashboard_type"`
	LayoutConfig    *json.RawMessage `json:"layout_config"`
	IsPublic        *bool            `json:"is_public"`
	IsDefault       *bool            `json:"is_default"`
	RefreshInterval *int             `json:"refresh_interval"`
}

// CreateAnalyticsReportRequest represents a request to create a report
type CreateAnalyticsReportRequest struct {
	ReportName        string   `json:"report_name" binding:"required"`
	Description       *string  `json:"description"`
	ReportType        string   `json:"report_type" binding:"required"`
	Category          *string  `json:"category"`
	DataSource        *string  `json:"data_source"`
	Format            string   `json:"format"`
	TimePeriod        string   `json:"time_period"`
	ScheduleFrequency *string  `json:"schedule_frequency"`
	IsPublic          *bool    `json:"is_public"`
	IsScheduled       *bool    `json:"is_scheduled"`
	RecipientEmails   []string `json:"recipient_emails"`
}

// CreateAnalyticsAlertRequest represents a request to create an alert
type CreateAnalyticsAlertRequest struct {
	AlertName       string   `json:"alert_name" binding:"required"`
	AlertType       string   `json:"alert_type" binding:"required"`
	MetricName      *string  `json:"metric_name"`
	Condition       string   `json:"condition" binding:"required"`
	ThresholdValue  *float64 `json:"threshold_value"`
	Severity        string   `json:"severity"`
	Enabled         *bool    `json:"enabled"`
	AlertFrequency  string   `json:"alert_frequency"`
	RecipientEmails []string `json:"recipient_emails"`
}

// GenerateAnalyticsReportRequest represents a request to generate a report
type GenerateAnalyticsReportRequest struct {
	ReportID       int64      `json:"report_id" binding:"required"`
	DataFromDate   *time.Time `json:"data_from_date"`
	DataToDate     *time.Time `json:"data_to_date"`
	Format         *string    `json:"format"`
	EmailOn        bool       `json:"email_on_complete"`
	RecipientEmail *string    `json:"recipient_email"`
}

// AnalyticsDashboardResponse represents a dashboard response
type AnalyticsDashboardResponse struct {
	ID            int64      `json:"id"`
	DashboardName string     `json:"dashboard_name"`
	Description   *string    `json:"description"`
	DashboardType string     `json:"dashboard_type"`
	IsPublic      bool       `json:"is_public"`
	IsDefault     bool       `json:"is_default"`
	WidgetCount   int        `json:"widget_count"`
	LastRefreshed *time.Time `json:"last_refreshed_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AnalyticsReportResponse represents a report response
type AnalyticsReportResponse struct {
	ID            int64      `json:"id"`
	ReportName    string     `json:"report_name"`
	ReportType    string     `json:"report_type"`
	Format        string     `json:"format"`
	IsScheduled   bool       `json:"is_scheduled"`
	LastGenerated *time.Time `json:"last_generated_at"`
	NextScheduled *time.Time `json:"next_scheduled_run"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// KPIDashboardResponse represents a KPI dashboard overview
type KPIDashboardResponse struct {
	TotalKPIs     int             `json:"total_kpis"`
	OnTrackCount  int             `json:"on_track_count"`
	AtRiskCount   int             `json:"at_risk_count"`
	CriticalCount int             `json:"critical_count"`
	KPIs          []*AnalyticsKPI `json:"kpis"`
	LastUpdated   *time.Time      `json:"last_updated"`
}
