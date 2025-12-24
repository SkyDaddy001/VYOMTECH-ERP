package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================
// ROI CALCULATION SERVICE
// ============================================================
// Purpose: Calculate TRUE ROI by combining:
// - Spend from Google Ads API
// - Spend from Meta Insights API
// - Revenue from ERP CRM (actual conversions, closed deals)
// - Attribution from ERP attribution events (which channel drove conversion)

// ROIMetrics is the core ROI data structure
type ROIMetrics struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Date      time.Time `json:"date"`
	DateRange string    `json:"date_range"` // "day" | "week" | "month" | "quarter" | "year"

	// Spend (from ad platforms)
	GoogleAdsCostMicros int64 `json:"google_ads_cost_micros"`
	MetaCostMicros      int64 `json:"meta_cost_micros"`
	TotalSpendMicros    int64 `json:"total_spend_micros"`

	// Attribution (from ERP)
	LeadsFromGoogleAds int64 `json:"leads_from_google_ads"`
	LeadsFromMeta      int64 `json:"leads_from_meta"`
	TotalLeads         int64 `json:"total_leads"`

	// Conversions (ERP CRM: qualified leads, site visits, closed deals)
	QualifiedLeads int64 `json:"qualified_leads"` // Lead moved to "qualified" stage
	SiteVisits     int64 `json:"site_visits"`     // Visitor actually visited property
	ClosedDeals    int64 `json:"closed_deals"`    // Actual sales

	// Revenue (from ERP sales data)
	TotalRevenueMicros int64 `json:"total_revenue_micros"`

	// Cost per metric
	CostPerLead  float64 `json:"cost_per_lead"`
	CostPerVisit float64 `json:"cost_per_visit"`
	CostPerDeal  float64 `json:"cost_per_deal"`

	// Attribution breakdown (ERP attribution)
	GoogleAdsRevenue int64 `json:"google_ads_revenue"`
	MetaRevenue      int64 `json:"meta_revenue"`

	// ROI (Return on Investment)
	GoogleAdsROI float64 `json:"google_ads_roi"` // (Revenue - Spend) / Spend
	MetaROI      float64 `json:"meta_roi"`
	OverallROI   float64 `json:"overall_roi"`

	// ROAS (Return on Ad Spend)
	GoogleAdsROAS float64 `json:"google_ads_roas"` // Revenue / Spend
	MetaROAS      float64 `json:"meta_roas"`
	OverallROAS   float64 `json:"overall_roas"`

	CalculatedAt time.Time `json:"calculated_at"`
}

// ============================================================
// ROI CALCULATION ENDPOINTS
// ============================================================

// calculateDailyROI calculates ROI for a given day
func calculateDailyROI(c *gin.Context) {
	var req struct {
		Date string `json:"date" binding:"required"` // "2025-12-24"
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	metrics, err := calculateROIForDate(c.Request.Context(), tenantID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate ROI"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// calculatePeriodROI calculates ROI for a date range
func calculatePeriodROI(c *gin.Context) {
	var req struct {
		StartDate string `json:"start_date" binding:"required"` // "2025-12-01"
		EndDate   string `json:"end_date" binding:"required"`   // "2025-12-31"
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	metrics, err := calculateROIForPeriod(c.Request.Context(), tenantID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate ROI"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// getRiskDashboard shows which channels are performing well vs poorly
func getRiskDashboard(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	// Get last 30 days ROI by channel
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	googleMetrics, err := getChannelROI(c.Request.Context(), tenantID, "google_ads", startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Google ROI"})
		return
	}

	metaMetrics, err := getChannelROI(c.Request.Context(), tenantID, "meta", startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Meta ROI"})
		return
	}

	// Classify performance
	googlePerformance := classifyPerformance(googleMetrics)
	metaPerformance := classifyPerformance(metaMetrics)

	c.JSON(http.StatusOK, gin.H{
		"period": "last_30_days",
		"google_ads": gin.H{
			"metrics":        googleMetrics,
			"performance":    googlePerformance,
			"recommendation": performanceRecommendation(googlePerformance),
		},
		"meta": gin.H{
			"metrics":        metaMetrics,
			"performance":    metaPerformance,
			"recommendation": performanceRecommendation(metaPerformance),
		},
	})
}

// budgetOptimization suggests how to allocate budget for maximum ROI
func budgetOptimization(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	// Get last 90 days data
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -90)

	googleMetrics, _ := getChannelROI(c.Request.Context(), tenantID, "google_ads", startDate, endDate)
	metaMetrics, _ := getChannelROI(c.Request.Context(), tenantID, "meta", startDate, endDate)

	// Calculate budget allocation
	googleROAS := googleMetrics.OverallROAS
	metaROAS := metaMetrics.OverallROAS

	totalROAS := googleROAS + metaROAS
	if totalROAS == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "insufficient_data",
			"message": "Not enough data to optimize budget",
		})
		return
	}

	// Allocate budget proportionally to ROAS
	googleAllocation := (googleROAS / totalROAS) * 100
	metaAllocation := (metaROAS / totalROAS) * 100

	// Get current budgets
	googleBudget := getTotalGoogleAdsBudget(c.Request.Context(), tenantID, startDate, endDate)
	metaBudget := getTotalMetaAdsBudget(c.Request.Context(), tenantID, startDate, endDate)
	totalBudget := googleBudget + metaBudget

	recommendedGoogleBudget := int64(float64(totalBudget) * (googleAllocation / 100))
	recommendedMetaBudget := int64(float64(totalBudget) * (metaAllocation / 100))

	c.JSON(http.StatusOK, gin.H{
		"period": "last_90_days",
		"current": gin.H{
			"google_ads_budget": googleBudget,
			"meta_budget":       metaBudget,
			"total_budget":      totalBudget,
		},
		"recommended": gin.H{
			"google_ads_budget": recommendedGoogleBudget,
			"google_allocation": fmt.Sprintf("%.1f%%", googleAllocation),
			"meta_budget":       recommendedMetaBudget,
			"meta_allocation":   fmt.Sprintf("%.1f%%", metaAllocation),
			"reason":            "Allocation proportional to ROAS for maximum ROI",
		},
		"metrics": gin.H{
			"google_ads_roas": googleMetrics.OverallROAS,
			"meta_roas":       metaMetrics.OverallROAS,
		},
	})
}

// ============================================================
// ROI CALCULATION LOGIC
// ============================================================

// calculateROIForDate calculates ROI for a single day
func calculateROIForDate(ctx context.Context, tenantID string, date time.Time) (*ROIMetrics, error) {
	metrics := &ROIMetrics{
		TenantID:     tenantID,
		Date:         date,
		DateRange:    "day",
		CalculatedAt: time.Now(),
	}

	// Step 1: Get spend from Google Ads metrics
	googleSpend, err := getGoogleAdsSpendForDate(ctx, tenantID, date)
	if err != nil {
		log.Printf("Failed to get Google spend: %v", err)
	}
	metrics.GoogleAdsCostMicros = googleSpend

	// Step 2: Get spend from Meta Insights
	metaSpend, err := getMetaAdsSpendForDate(ctx, tenantID, date)
	if err != nil {
		log.Printf("Failed to get Meta spend: %v", err)
	}
	metrics.MetaCostMicros = metaSpend
	metrics.TotalSpendMicros = googleSpend + metaSpend

	// Step 3: Get attribution data from ERP
	// Count leads attributed to each channel
	googleLeads, err := getLeadsAttributedToChannel(ctx, tenantID, "google_ads", date, date)
	if err != nil {
		log.Printf("Failed to get Google leads: %v", err)
	}
	metrics.LeadsFromGoogleAds = googleLeads

	metaLeads, err := getLeadsAttributedToChannel(ctx, tenantID, "meta", date, date)
	if err != nil {
		log.Printf("Failed to get Meta leads: %v", err)
	}
	metrics.LeadsFromMeta = metaLeads
	metrics.TotalLeads = googleLeads + metaLeads

	// Step 4: Get conversions (from CRM)
	// Qualified leads: moved to qualified stage
	qualifiedFromGoogle, _, _ := getConversionsAttributedToChannel(ctx, tenantID, "google_ads", "qualified_lead", date, date)
	qualifiedFromMeta, _, _ := getConversionsAttributedToChannel(ctx, tenantID, "meta", "qualified_lead", date, date)
	metrics.QualifiedLeads = qualifiedFromGoogle + qualifiedFromMeta

	// Site visits: scheduled/completed property visit
	visitsFromGoogle, _, _ := getConversionsAttributedToChannel(ctx, tenantID, "google_ads", "site_visit", date, date)
	visitsFromMeta, _, _ := getConversionsAttributedToChannel(ctx, tenantID, "meta", "site_visit", date, date)
	metrics.SiteVisits = visitsFromGoogle + visitsFromMeta

	// Closed deals: actual sale
	dealsFromGoogle, revenueGoogle, _ := getConversionsAttributedToChannel(ctx, tenantID, "google_ads", "closed_deal", date, date)
	dealsFromMeta, revenueMeta, _ := getConversionsAttributedToChannel(ctx, tenantID, "meta", "closed_deal", date, date)
	metrics.ClosedDeals = dealsFromGoogle + dealsFromMeta
	metrics.TotalRevenueMicros = revenueGoogle + revenueMeta
	metrics.GoogleAdsRevenue = revenueGoogle
	metrics.MetaRevenue = revenueMeta

	// Step 5: Calculate KPIs
	if metrics.TotalLeads > 0 {
		metrics.CostPerLead = float64(metrics.TotalSpendMicros) / float64(metrics.TotalLeads) / 1_000_000
	}
	if metrics.SiteVisits > 0 {
		metrics.CostPerVisit = float64(metrics.TotalSpendMicros) / float64(metrics.SiteVisits) / 1_000_000
	}
	if metrics.ClosedDeals > 0 {
		metrics.CostPerDeal = float64(metrics.TotalSpendMicros) / float64(metrics.ClosedDeals) / 1_000_000
	}

	// Step 6: Calculate ROI and ROAS
	if metrics.GoogleAdsCostMicros > 0 {
		googleProfit := metrics.GoogleAdsRevenue - metrics.GoogleAdsCostMicros
		metrics.GoogleAdsROI = (float64(googleProfit) / float64(metrics.GoogleAdsCostMicros)) * 100
		metrics.GoogleAdsROAS = float64(metrics.GoogleAdsRevenue) / float64(metrics.GoogleAdsCostMicros)
	}

	if metrics.MetaCostMicros > 0 {
		metaProfit := metrics.MetaRevenue - metrics.MetaCostMicros
		metrics.MetaROI = (float64(metaProfit) / float64(metrics.MetaCostMicros)) * 100
		metrics.MetaROAS = float64(metrics.MetaRevenue) / float64(metrics.MetaCostMicros)
	}

	if metrics.TotalSpendMicros > 0 {
		totalProfit := metrics.TotalRevenueMicros - metrics.TotalSpendMicros
		metrics.OverallROI = (float64(totalProfit) / float64(metrics.TotalSpendMicros)) * 100
		metrics.OverallROAS = float64(metrics.TotalRevenueMicros) / float64(metrics.TotalSpendMicros)
	}

	return metrics, nil
}

// calculateROIForPeriod calculates ROI for a date range
func calculateROIForPeriod(ctx context.Context, tenantID string, startDate, endDate time.Time) (*ROIMetrics, error) {
	aggregatedMetrics := &ROIMetrics{
		TenantID:     tenantID,
		Date:         startDate,
		DateRange:    "custom",
		CalculatedAt: time.Now(),
	}

	// Sum up daily metrics
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		daily, err := calculateROIForDate(ctx, tenantID, currentDate)
		if err != nil {
			log.Printf("Failed to calculate daily ROI for %s: %v", currentDate, err)
			currentDate = currentDate.AddDate(0, 0, 1)
			continue
		}

		aggregatedMetrics.GoogleAdsCostMicros += daily.GoogleAdsCostMicros
		aggregatedMetrics.MetaCostMicros += daily.MetaCostMicros
		aggregatedMetrics.LeadsFromGoogleAds += daily.LeadsFromGoogleAds
		aggregatedMetrics.LeadsFromMeta += daily.LeadsFromMeta
		aggregatedMetrics.QualifiedLeads += daily.QualifiedLeads
		aggregatedMetrics.SiteVisits += daily.SiteVisits
		aggregatedMetrics.ClosedDeals += daily.ClosedDeals
		aggregatedMetrics.TotalRevenueMicros += daily.TotalRevenueMicros
		aggregatedMetrics.GoogleAdsRevenue += daily.GoogleAdsRevenue
		aggregatedMetrics.MetaRevenue += daily.MetaRevenue

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	aggregatedMetrics.TotalSpendMicros = aggregatedMetrics.GoogleAdsCostMicros + aggregatedMetrics.MetaCostMicros
	aggregatedMetrics.TotalLeads = aggregatedMetrics.LeadsFromGoogleAds + aggregatedMetrics.LeadsFromMeta

	// Recalculate KPIs
	if aggregatedMetrics.TotalLeads > 0 {
		aggregatedMetrics.CostPerLead = float64(aggregatedMetrics.TotalSpendMicros) / float64(aggregatedMetrics.TotalLeads) / 1_000_000
	}
	if aggregatedMetrics.SiteVisits > 0 {
		aggregatedMetrics.CostPerVisit = float64(aggregatedMetrics.TotalSpendMicros) / float64(aggregatedMetrics.SiteVisits) / 1_000_000
	}
	if aggregatedMetrics.ClosedDeals > 0 {
		aggregatedMetrics.CostPerDeal = float64(aggregatedMetrics.TotalSpendMicros) / float64(aggregatedMetrics.ClosedDeals) / 1_000_000
	}

	// Recalculate ROI and ROAS
	if aggregatedMetrics.GoogleAdsCostMicros > 0 {
		googleProfit := aggregatedMetrics.GoogleAdsRevenue - aggregatedMetrics.GoogleAdsCostMicros
		aggregatedMetrics.GoogleAdsROI = (float64(googleProfit) / float64(aggregatedMetrics.GoogleAdsCostMicros)) * 100
		aggregatedMetrics.GoogleAdsROAS = float64(aggregatedMetrics.GoogleAdsRevenue) / float64(aggregatedMetrics.GoogleAdsCostMicros)
	}

	if aggregatedMetrics.MetaCostMicros > 0 {
		metaProfit := aggregatedMetrics.MetaRevenue - aggregatedMetrics.MetaCostMicros
		aggregatedMetrics.MetaROI = (float64(metaProfit) / float64(aggregatedMetrics.MetaCostMicros)) * 100
		aggregatedMetrics.MetaROAS = float64(aggregatedMetrics.MetaRevenue) / float64(aggregatedMetrics.MetaCostMicros)
	}

	if aggregatedMetrics.TotalSpendMicros > 0 {
		totalProfit := aggregatedMetrics.TotalRevenueMicros - aggregatedMetrics.TotalSpendMicros
		aggregatedMetrics.OverallROI = (float64(totalProfit) / float64(aggregatedMetrics.TotalSpendMicros)) * 100
		aggregatedMetrics.OverallROAS = float64(aggregatedMetrics.TotalRevenueMicros) / float64(aggregatedMetrics.TotalSpendMicros)
	}

	return aggregatedMetrics, nil
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func classifyPerformance(metrics *ROIMetrics) string {
	if metrics.OverallROAS >= 3.0 {
		return "excellent"
	} else if metrics.OverallROAS >= 2.0 {
		return "good"
	} else if metrics.OverallROAS >= 1.0 {
		return "acceptable"
	} else {
		return "poor"
	}
}

func performanceRecommendation(performance string) string {
	recommendations := map[string]string{
		"excellent":  "Channel is performing well. Consider increasing budget allocation.",
		"good":       "Channel is performing well. Maintain current budget or increase slightly.",
		"acceptable": "Channel is profitable but has room for optimization. Review targeting and creative quality.",
		"poor":       "Channel is unprofitable. Review campaign setup, targeting, and consider pausing underperforming ads.",
	}
	return recommendations[performance]
}

// Database query stubs (implement with actual DB calls)
func getChannelROI(ctx context.Context, tenantID, channel string, startDate, endDate time.Time) (*ROIMetrics, error) {
	// TODO: Query roi_metrics table for channel
	return nil, nil
}

func getGoogleAdsSpendForDate(ctx context.Context, tenantID string, date time.Time) (int64, error) {
	// TODO: Query google_ads_metrics and SUM(cost_micros) for date
	return 0, nil
}

func getMetaAdsSpendForDate(ctx context.Context, tenantID string, date time.Time) (int64, error) {
	// TODO: Query meta_ads_metrics and SUM(cost_micros) for date
	return 0, nil
}

func getLeadsAttributedToChannel(ctx context.Context, tenantID, channel string, startDate, endDate time.Time) (int64, error) {
	// TODO: Query attribution_event WHERE source = $channel
	return 0, nil
}

func getConversionsAttributedToChannel(ctx context.Context, tenantID, channel, conversionType string, startDate, endDate time.Time) (int64, int64, error) {
	// TODO: Query lead_attribution_snapshot filtered by channel
	// Return count and revenue
	return 0, 0, nil
}

func getTotalGoogleAdsBudget(ctx context.Context, tenantID string, startDate, endDate time.Time) int64 {
	// TODO: Query google_ads_metrics and SUM daily_budget for period
	return 0
}

func getTotalMetaAdsBudget(ctx context.Context, tenantID string, startDate, endDate time.Time) int64 {
	// TODO: Query meta_ads_metrics and SUM daily_budget for period
	return 0
}

// ============================================================
// REGISTER ROUTES
// ============================================================

func registerROIRoutes(router *gin.Engine) {
	router.POST("/api/v1/roi/daily", calculateDailyROI)
	router.POST("/api/v1/roi/period", calculatePeriodROI)
	router.GET("/api/v1/roi/dashboard", getRiskDashboard)
	router.GET("/api/v1/roi/optimize-budget", budgetOptimization)
}
