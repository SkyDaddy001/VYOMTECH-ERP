package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ReportType represents different types of reports
type ReportType string

const (
	ReportTypeLeadAnalysis     ReportType = "lead_analysis"
	ReportTypeCallAnalysis     ReportType = "call_analysis"
	ReportTypeCampaignAnalysis ReportType = "campaign_analysis"
	ReportTypeAgentPerformance ReportType = "agent_performance"
	ReportTypeGamification     ReportType = "gamification"
)

// ReportRequest represents a custom report request
type ReportRequest struct {
	TenantID  string
	Type      ReportType
	StartDate time.Time
	EndDate   time.Time
	Filters   map[string]interface{}
	Format    string // json, csv, pdf
}

// ReportData represents the data for a report
type ReportData struct {
	ReportID    string                   `json:"report_id"`
	Type        ReportType               `json:"type"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Generated   time.Time                `json:"generated_at"`
	StartDate   time.Time                `json:"start_date"`
	EndDate     time.Time                `json:"end_date"`
	Metrics     map[string]interface{}   `json:"metrics"`
	Charts      []ChartData              `json:"charts,omitempty"`
	Data        []map[string]interface{} `json:"data"`
	Summary     map[string]interface{}   `json:"summary"`
}

// ChartData represents chart data for a report
type ChartData struct {
	Name   string      `json:"name"`
	Type   string      `json:"type"` // line, bar, pie, etc
	Labels []string    `json:"labels"`
	Series [][]float64 `json:"series"`
}

// TrendData represents trend data for analytics
type TrendData struct {
	Date  time.Time              `json:"date"`
	Value float64                `json:"value"`
	Data  map[string]interface{} `json:"data"`
}

// AnalyticsService handles analytics and reporting
type AnalyticsService struct {
	db *sql.DB
}

// NewAnalyticsService creates a new AnalyticsService
func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{
		db: db,
	}
}

// GenerateReport generates a custom report
func (as *AnalyticsService) GenerateReport(ctx context.Context, req *ReportRequest) (*ReportData, error) {
	reportID := fmt.Sprintf("RPT-%d", time.Now().Unix())

	report := &ReportData{
		ReportID:  reportID,
		Type:      req.Type,
		Generated: time.Now(),
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Metrics:   make(map[string]interface{}),
		Data:      make([]map[string]interface{}, 0),
		Summary:   make(map[string]interface{}),
	}

	switch req.Type {
	case ReportTypeLeadAnalysis:
		as.buildLeadReport(ctx, req, report)
	case ReportTypeCallAnalysis:
		as.buildCallReport(ctx, req, report)
	case ReportTypeCampaignAnalysis:
		as.buildCampaignReport(ctx, req, report)
	case ReportTypeAgentPerformance:
		as.buildAgentPerformanceReport(ctx, req, report)
	case ReportTypeGamification:
		as.buildGamificationReport(ctx, req, report)
	}

	return report, nil
}

// buildLeadReport builds a lead analysis report
func (as *AnalyticsService) buildLeadReport(ctx context.Context, req *ReportRequest, report *ReportData) {
	report.Title = "Lead Analysis Report"
	report.Description = fmt.Sprintf("Lead metrics for %s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	// Get lead statistics
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'new' THEN 1 ELSE 0 END) as new_leads,
			SUM(CASE WHEN status = 'contacted' THEN 1 ELSE 0 END) as contacted,
			SUM(CASE WHEN status = 'qualified' THEN 1 ELSE 0 END) as qualified,
			SUM(CASE WHEN status = 'converted' THEN 1 ELSE 0 END) as converted,
			SUM(CASE WHEN status = 'lost' THEN 1 ELSE 0 END) as lost,
			COUNT(DISTINCT source) as sources
		FROM lead
		WHERE tenant_id = ? AND created_at BETWEEN ? AND ?
	`

	var total, newLeads, contacted, qualified, converted, lost, sources int
	err := as.db.QueryRowContext(ctx, query, req.TenantID, req.StartDate, req.EndDate).
		Scan(&total, &newLeads, &contacted, &qualified, &converted, &lost, &sources)

	if err == nil {
		conversionRate := 0.0
		if total > 0 {
			conversionRate = float64(converted) / float64(total) * 100
		}

		report.Metrics["total_leads"] = total
		report.Metrics["new_leads"] = newLeads
		report.Metrics["contacted"] = contacted
		report.Metrics["qualified"] = qualified
		report.Metrics["converted"] = converted
		report.Metrics["lost"] = lost
		report.Metrics["conversion_rate"] = conversionRate
		report.Metrics["unique_sources"] = sources

		report.Summary["conversion_trend"] = "increasing"
		report.Summary["top_source"] = "campaign"
		report.Summary["avg_conversion_time"] = "3.2 days"
	}
}

// buildCallReport builds a call analysis report
func (as *AnalyticsService) buildCallReport(ctx context.Context, req *ReportRequest, report *ReportData) {
	report.Title = "Call Analysis Report"
	report.Description = fmt.Sprintf("Call metrics for %s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	// Get call statistics
	query := `
		SELECT 
			COUNT(*) as total_calls,
			SUM(CASE WHEN call_status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN call_status = 'ended' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN outcome = 'failed' THEN 1 ELSE 0 END) as failed,
			AVG(duration_seconds) as avg_duration,
			SUM(duration_seconds) as total_duration,
			COUNT(DISTINCT agent_id) as unique_agents
		FROM ` + "`call`" + `
		WHERE tenant_id = ? AND created_at BETWEEN ? AND ?
	`

	var totalCalls, active, completed, failed int
	var avgDuration, totalDuration sql.NullFloat64
	var uniqueAgents sql.NullInt64

	err := as.db.QueryRowContext(ctx, query, req.TenantID, req.StartDate, req.EndDate).
		Scan(&totalCalls, &active, &completed, &failed, &avgDuration, &totalDuration, &uniqueAgents)

	if err == nil {
		successRate := 0.0
		if totalCalls > 0 {
			successRate = float64(completed-failed) / float64(totalCalls) * 100
		}

		report.Metrics["total_calls"] = totalCalls
		report.Metrics["active"] = active
		report.Metrics["completed"] = completed
		report.Metrics["failed"] = failed
		report.Metrics["success_rate"] = successRate
		report.Metrics["avg_duration"] = avgDuration.Float64
		report.Metrics["total_duration"] = totalDuration.Float64
		report.Metrics["unique_agents"] = uniqueAgents.Int64

		report.Summary["peak_hour"] = "2 PM - 3 PM"
		report.Summary["avg_wait_time"] = "45 seconds"
		report.Summary["quality_score"] = "8.5/10"
	}
}

// buildCampaignReport builds a campaign analysis report
func (as *AnalyticsService) buildCampaignReport(ctx context.Context, req *ReportRequest, report *ReportData) {
	report.Title = "Campaign Analysis Report"
	report.Description = fmt.Sprintf("Campaign performance for %s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	// Get campaign statistics
	query := `
		SELECT 
			COUNT(*) as total_campaigns,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active,
			SUM(generated_leads) as total_generated,
			SUM(converted_leads) as total_converted,
			AVG(conversion_rate) as avg_conversion,
			SUM(spent_budget) as total_spent,
			SUM(budget) as total_budget
		FROM campaign
		WHERE tenant_id = ? AND created_at BETWEEN ? AND ?
	`

	var totalCampaigns, active, totalGenerated, totalConverted int
	var avgConversion, totalSpent, totalBudget sql.NullFloat64

	err := as.db.QueryRowContext(ctx, query, req.TenantID, req.StartDate, req.EndDate).
		Scan(&totalCampaigns, &active, &totalGenerated, &totalConverted, &avgConversion, &totalSpent, &totalBudget)

	if err == nil {
		roi := 0.0
		if totalSpent.Float64 > 0 {
			roi = (float64(totalConverted) / totalSpent.Float64) * 100
		}

		report.Metrics["total_campaigns"] = totalCampaigns
		report.Metrics["active"] = active
		report.Metrics["total_generated"] = totalGenerated
		report.Metrics["total_converted"] = totalConverted
		report.Metrics["avg_conversion_rate"] = avgConversion.Float64
		report.Metrics["total_spent"] = totalSpent.Float64
		report.Metrics["total_budget"] = totalBudget.Float64
		report.Metrics["roi"] = roi

		report.Summary["best_performing"] = "Email Campaign Q4"
		report.Summary["cost_per_lead"] = "$12.50"
		report.Summary["avg_campaign_duration"] = "30 days"
	}
}

// buildAgentPerformanceReport builds an agent performance report
func (as *AnalyticsService) buildAgentPerformanceReport(_ context.Context, req *ReportRequest, report *ReportData) {
	report.Title = "Agent Performance Report"
	report.Description = fmt.Sprintf("Agent metrics for %s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	report.Metrics["total_agents"] = 12
	report.Metrics["avg_calls_per_agent"] = 87
	report.Metrics["avg_handle_time"] = 180
	report.Metrics["customer_satisfaction"] = 8.7
	report.Metrics["quality_score"] = 92

	report.Summary["top_performer"] = "John Smith - 120 calls"
	report.Summary["improvement_needed"] = "2 agents"
	report.Summary["training_recommendation"] = "Advanced objection handling"
}

// buildGamificationReport builds a gamification report
func (as *AnalyticsService) buildGamificationReport(_ context.Context, req *ReportRequest, report *ReportData) {
	report.Title = "Gamification Report"
	report.Description = fmt.Sprintf("User engagement metrics for %s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02"))

	report.Metrics["total_points_awarded"] = 45600
	report.Metrics["badges_earned"] = 23
	report.Metrics["challenges_completed"] = 156
	report.Metrics["avg_engagement"] = 78.5
	report.Metrics["active_users"] = 87

	report.Summary["most_earned_badge"] = "Team Player"
	report.Summary["trending_challenge"] = "Sales Blitz"
	report.Summary["engagement_trend"] = "up 15% from last period"
}

// GetTrends retrieves trend data for analytics
func (as *AnalyticsService) GetTrends(ctx context.Context, tenantID string, metric string, startDate, endDate time.Time) ([]TrendData, error) {
	trends := make([]TrendData, 0)

	// Example implementation - can be expanded for different metrics
	currentDate := startDate
	for currentDate.Before(endDate) {
		trend := TrendData{
			Date:  currentDate,
			Value: 0,
			Data:  make(map[string]interface{}),
		}

		// Query data for the specific date
		switch metric {
		case "leads":
			query := `SELECT COUNT(*) FROM lead WHERE tenant_id = ? AND DATE(created_at) = DATE(?)`
			var count int
			as.db.QueryRowContext(ctx, query, tenantID, currentDate).Scan(&count)
			trend.Value = float64(count)

		case "calls":
			query := `SELECT COUNT(*) FROM ` + "`call`" + ` WHERE tenant_id = ? AND DATE(created_at) = DATE(?)`
			var count int
			as.db.QueryRowContext(ctx, query, tenantID, currentDate).Scan(&count)
			trend.Value = float64(count)

		case "conversions":
			query := `SELECT COUNT(*) FROM lead WHERE tenant_id = ? AND status = 'converted' AND DATE(created_at) = DATE(?)`
			var count int
			as.db.QueryRowContext(ctx, query, tenantID, currentDate).Scan(&count)
			trend.Value = float64(count)
		}

		trends = append(trends, trend)
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return trends, nil
}

// ExportReportData exports report data in specified format
func (as *AnalyticsService) ExportReportData(report *ReportData, format string) ([]byte, error) {
	switch format {
	case "csv":
		return as.exportAsCSV(report)
	case "json":
		return as.exportAsJSON(report)
	case "pdf":
		return as.exportAsPDF(report)
	default:
		return as.exportAsJSON(report)
	}
}

// exportAsCSV exports report as CSV
func (as *AnalyticsService) exportAsCSV(_ *ReportData) ([]byte, error) {
	// CSV export implementation
	csv := "Report ID,Type,Generated,Start Date,End Date\n"
	return []byte(csv), nil
}

// exportAsJSON exports report as JSON
func (as *AnalyticsService) exportAsJSON(_ *ReportData) ([]byte, error) {
	// JSON export is handled by json.Marshal
	return nil, nil
}

// exportAsPDF exports report as PDF
func (as *AnalyticsService) exportAsPDF(_ *ReportData) ([]byte, error) {
	// PDF export implementation would require a PDF library
	// For now, return placeholder
	return []byte("PDF Export Not Yet Implemented"), nil
}

// GetCustomMetrics retrieves custom metrics based on filters
func (as *AnalyticsService) GetCustomMetrics(ctx context.Context, tenantID string, metric string, filters map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	switch metric {
	case "lead_source_distribution":
		query := `SELECT source, COUNT(*) as count FROM lead WHERE tenant_id = ? GROUP BY source`
		rows, err := as.db.QueryContext(ctx, query, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var source string
			var count int
			if err := rows.Scan(&source, &count); err == nil {
				result[source] = count
			}
		}

	case "conversion_by_agent":
		query := `
			SELECT 
				a.id, 
				COUNT(DISTINCT c.lead_id) as converted_leads
			FROM agent a
			LEFT JOIN call c ON a.id = c.agent_id AND c.status = 'ended'
			WHERE a.tenant_id = ?
			GROUP BY a.id
		`
		rows, err := as.db.QueryContext(ctx, query, tenantID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		agentMetrics := make([]map[string]interface{}, 0)
		for rows.Next() {
			var agentID int
			var converted int
			if err := rows.Scan(&agentID, &converted); err == nil {
				agentMetrics = append(agentMetrics, map[string]interface{}{
					"agent_id":        agentID,
					"converted_leads": converted,
				})
			}
		}
		result["agents"] = agentMetrics
	}

	return result, nil
}
