package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================================
// META ADS INTEGRATION (Facebook, Instagram)
// ============================================================
// Purpose: Control Meta campaigns, sync metrics to ERP
// Key difference from Google: Meta uses campaign → ad_set → ad hierarchy
// ERP is system of record, Meta is execution engine
// Meta attribution is probabilistic, ERP attribution is deterministic

// ============================================================
// DATA STRUCTURES (matching Prisma schema)
// ============================================================

type MetaAdsCampaign struct {
	ID                   string          `json:"id"`
	TenantID             string          `json:"tenant_id"`
	MetaCampaignID       string          `json:"meta_campaign_id"`
	MetaAdAccountID      string          `json:"meta_ad_account_id"`
	Name                 string          `json:"name"`
	Status               string          `json:"status"`    // "ACTIVE" | "PAUSED" | "DELETED"
	Objective            string          `json:"objective"` // "LINK_CLICKS" | "CONVERSIONS" | "LEAD_GENERATION" | "PAGE_LIKES" etc
	DailyBudgetMicros    *int64          `json:"daily_budget_micros"`
	TotalBudgetMicros    *int64          `json:"total_budget_micros"`
	LifetimeBudgetMicros *int64          `json:"lifetime_budget_micros"`
	TargetingDimensions  json.RawMessage `json:"targeting_dimensions"`
	ProjectID            string          `json:"project_id"`
	CustomMetadata       json.RawMessage `json:"custom_metadata"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type MetaAdsAdSet struct {
	ID                   string          `json:"id"`
	CampaignID           string          `json:"campaign_id"`
	MetaAdSetID          string          `json:"meta_ad_set_id"`
	Name                 string          `json:"name"`
	Status               string          `json:"status"` // "ACTIVE" | "PAUSED"
	DailyBudgetMicros    *int64          `json:"daily_budget_micros"`
	LifetimeBudgetMicros *int64          `json:"lifetime_budget_micros"`
	TargetingDimensions  json.RawMessage `json:"targeting_dimensions"`
	CreativeIds          json.RawMessage `json:"creative_ids"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type MetaAdsMetrics struct {
	ID                 string          `json:"id"`
	TenantID           string          `json:"tenant_id"`
	CampaignID         string          `json:"campaign_id"`
	Date               time.Time       `json:"date"`
	Impressions        int64           `json:"impressions"`
	Clicks             int64           `json:"clicks"`
	CTR                float64         `json:"ctr"`
	AverageCpc         int64           `json:"average_cpc"`
	CostMicros         int64           `json:"cost_micros"`
	Spend              int64           `json:"spend"`
	WebConversions     float64         `json:"web_conversions"` // Meta's conversion tracking
	WebConversionValue float64         `json:"web_conversion_value"`
	ConversionRate     float64         `json:"conversion_rate"`
	Actions            json.RawMessage `json:"actions"`
	Likes              int64           `json:"likes"`
	Comments           int64           `json:"comments"`
	Shares             int64           `json:"shares"`
	Reach              int64           `json:"reach"`
	Frequency          *float64        `json:"frequency"`
	SyncedAt           time.Time       `json:"synced_at"`
}

// CreateMetaCampaignRequest is incoming request
type CreateMetaCampaignRequest struct {
	AccountConnectionID  string          `json:"account_connection_id" binding:"required"`
	Name                 string          `json:"name" binding:"required"`
	Objective            string          `json:"objective" binding:"required"` // LINK_CLICKS, CONVERSIONS, LEAD_GENERATION, PAGE_LIKES, etc
	DailyBudgetMicros    *int64          `json:"daily_budget_micros"`
	LifetimeBudgetMicros *int64          `json:"lifetime_budget_micros"`
	TargetingDimensions  json.RawMessage `json:"targeting_dimensions"`
	ProjectID            string          `json:"project_id"`
	CustomMetadata       json.RawMessage `json:"custom_metadata"`
}

type CreateMetaAdSetRequest struct {
	CampaignID           string          `json:"campaign_id" binding:"required"`
	Name                 string          `json:"name" binding:"required"`
	DailyBudgetMicros    *int64          `json:"daily_budget_micros"`
	LifetimeBudgetMicros *int64          `json:"lifetime_budget_micros"`
	TargetingDimensions  json.RawMessage `json:"targeting_dimensions"`
	CreativeIds          json.RawMessage `json:"creative_ids"`
}

// ============================================================
// META ADS CONTROL ENDPOINTS (ERP → Meta)
// ============================================================

// createMetaAdsCampaign creates a campaign in Meta
// ERP command: Create campaign with objective, budget, targeting
// Meta execution: Create via Graph API
// Store: Save campaign info + mapping
func createMetaAdsCampaign(c *gin.Context) {
	var req CreateMetaCampaignRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant_id required"})
		return
	}

	// Validate objective
	validObjectives := map[string]bool{
		"LINK_CLICKS":     true,
		"CONVERSIONS":     true,
		"LEAD_GENERATION": true,
		"PAGE_LIKES":      true,
		"VIDEO_VIEWS":     true,
		"REACH":           true,
		"FREQUENCY":       true,
		"APP_INSTALLS":    true,
		"BRAND_AWARENESS": true,
	}
	if !validObjectives[req.Objective] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid objective"})
		return
	}

	// Get account connection
	conn, err := getAdAccountConnection(c.Request.Context(), tenantID, req.AccountConnectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account connection not found"})
		return
	}

	if conn.Platform != "meta" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "connection is not for Meta"})
		return
	}

	if !conn.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account connection is not active"})
		return
	}

	// Call Meta API to create campaign
	metaCampaignID, err := createMetaAdsCampaignViaAPI(c.Request.Context(), conn, &req)
	if err != nil {
		log.Printf("Meta API error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create campaign in Meta"})
		return
	}

	// Save to database
	campaign := &MetaAdsCampaign{
		ID:                   uuid.NewString(),
		TenantID:             tenantID,
		MetaCampaignID:       metaCampaignID,
		MetaAdAccountID:      conn.AccountID,
		Name:                 req.Name,
		Status:               "ACTIVE",
		Objective:            req.Objective,
		DailyBudgetMicros:    req.DailyBudgetMicros,
		LifetimeBudgetMicros: req.LifetimeBudgetMicros,
		TargetingDimensions:  req.TargetingDimensions,
		ProjectID:            req.ProjectID,
		CustomMetadata:       req.CustomMetadata,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := storeMetaAdsCampaign(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save campaign"})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// createMetaAdsAdSet creates an ad set within a campaign
// Meta hierarchy: Campaign → AdSet → Ads
// Ad set contains: budget, targeting, schedule, creatives
func createMetaAdsAdSet(c *gin.Context) {
	var req CreateMetaAdSetRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenant_id")

	// Get campaign
	campaign, err := getMetaAdsCampaign(c.Request.Context(), tenantID, req.CampaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}

	// Get account connection
	conn, err := getAdAccountConnectionByMetaID(c.Request.Context(), campaign.MetaAdAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account not found"})
		return
	}

	// Call Meta API to create ad set
	metaAdSetID, err := createMetaAdsAdSetViaAPI(c.Request.Context(), conn, campaign, &req)
	if err != nil {
		log.Printf("Meta API error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ad set"})
		return
	}

	// Save to database
	adSet := &MetaAdsAdSet{
		ID:                   uuid.NewString(),
		CampaignID:           campaign.ID,
		MetaAdSetID:          metaAdSetID,
		Name:                 req.Name,
		Status:               "ACTIVE",
		DailyBudgetMicros:    req.DailyBudgetMicros,
		LifetimeBudgetMicros: req.LifetimeBudgetMicros,
		TargetingDimensions:  req.TargetingDimensions,
		CreativeIds:          req.CreativeIds,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := storeMetaAdsAdSet(c.Request.Context(), adSet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save ad set"})
		return
	}

	c.JSON(http.StatusCreated, adSet)
}

// pauseMetaAdsCampaign pauses campaign
func pauseMetaAdsCampaign(c *gin.Context) {
	campaignID := c.Param("campaign_id")
	tenantID := c.GetString("tenant_id")

	campaign, err := getMetaAdsCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}

	conn, err := getAdAccountConnectionByMetaID(c.Request.Context(), campaign.MetaAdAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account not found"})
		return
	}

	if err := pauseMetaAdsCampaignViaAPI(c.Request.Context(), conn, campaign.MetaCampaignID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to pause campaign"})
		return
	}

	campaign.Status = "PAUSED"
	campaign.UpdatedAt = time.Now()
	if err := updateMetaAdsCampaignStatus(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// updateMetaAdsBudget updates budget for campaign or ad set
func updateMetaAdsBudget(c *gin.Context) {
	campaignID := c.Param("campaign_id")
	tenantID := c.GetString("tenant_id")

	var req struct {
		BudgetType   string `json:"budget_type" binding:"required"` // "daily" | "lifetime"
		BudgetMicros int64  `json:"budget_micros" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := getMetaAdsCampaign(c.Request.Context(), tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}

	conn, err := getAdAccountConnectionByMetaID(c.Request.Context(), campaign.MetaAdAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account not found"})
		return
	}

	if err := updateMetaAdsBudgetViaAPI(c.Request.Context(), conn, campaign.MetaCampaignID, req.BudgetType, req.BudgetMicros); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update budget"})
		return
	}

	if req.BudgetType == "daily" {
		campaign.DailyBudgetMicros = &req.BudgetMicros
	} else {
		campaign.LifetimeBudgetMicros = &req.BudgetMicros
	}
	campaign.UpdatedAt = time.Now()
	if err := updateMetaAdsCampaignStatus(c.Request.Context(), campaign); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update locally"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// getMetaAdsCampaigns gets all campaigns
func getMetaAdsCampaigns(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	campaigns, err := getMetaAdsCampaignsByTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch campaigns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":     len(campaigns),
		"campaigns": campaigns,
	})
}

// ============================================================
// META ADS SYNC ENDPOINTS (Meta → ERP)
// ============================================================

// syncMetaAdsMetrics pulls metrics from Meta API
// Called daily via scheduler
func syncMetaAdsMetrics(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	// Get all active Meta accounts
	accounts, err := getAdAccountConnectionsByPlatform(c.Request.Context(), tenantID, "meta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})
		return
	}

	syncResults := make([]gin.H, 0)

	for _, account := range accounts {
		if !account.SyncEnabled {
			continue
		}

		account.SyncStatus = "syncing"
		updateAdAccountConnectionStatus(c.Request.Context(), account)

		// Get campaigns for this account
		campaigns, err := getMetaAdsCampaignsByMetaID(c.Request.Context(), tenantID, account.AccountID)
		if err != nil {
			account.SyncStatus = "failed"
			account.LastSyncError = fmt.Sprintf("Failed to fetch campaigns: %v", err)
			updateAdAccountConnectionStatus(c.Request.Context(), account)
			continue
		}

		for _, campaign := range campaigns {
			// Pull metrics for last 7 days
			endDate := time.Now()
			startDate := endDate.AddDate(0, 0, -7)

			metrics, err := fetchMetaAdsMetricsViaAPI(c.Request.Context(), account, campaign, startDate, endDate)
			if err != nil {
				log.Printf("Failed to sync metrics for campaign %s: %v", campaign.ID, err)
				continue
			}

			// Store metrics
			for _, metric := range metrics {
				metric.TenantID = tenantID
				metric.CampaignID = campaign.ID
				metric.SyncedAt = time.Now()
				if err := storeMetaAdsMetrics(c.Request.Context(), metric); err != nil {
					log.Printf("Failed to store metric: %v", err)
				}
			}

			campaign.UpdatedAt = time.Now()
			updateMetaAdsCampaignStatus(c.Request.Context(), campaign)
		}

		// Update sync status to success
		now := time.Now()
		account.SyncStatus = "success"
		account.LastSyncAt = &now
		account.LastSyncError = ""
		updateAdAccountConnectionStatus(c.Request.Context(), account)

		syncResults = append(syncResults, gin.H{
			"account_id": account.AccountID,
			"status":     "success",
			"campaigns":  len(campaigns),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"synced_at": time.Now(),
		"results":   syncResults,
	})
}

// ============================================================
// HELPER FUNCTIONS (Database)
// ============================================================

func storeMetaAdsCampaign(ctx context.Context, campaign *MetaAdsCampaign) error {
	// TODO: INSERT INTO meta_ads_campaign (...) VALUES (...)
	return nil
}

func getMetaAdsCampaign(ctx context.Context, tenantID, campaignID string) (*MetaAdsCampaign, error) {
	// TODO: SELECT * FROM meta_ads_campaign WHERE tenant_id = $1 AND id = $2
	return nil, nil
}

func getMetaAdsCampaignsByTenant(ctx context.Context, tenantID string) ([]*MetaAdsCampaign, error) {
	// TODO: SELECT * FROM meta_ads_campaign WHERE tenant_id = $1
	return nil, nil
}

func getMetaAdsCampaignsByMetaID(ctx context.Context, tenantID, metaAdAccountID string) ([]*MetaAdsCampaign, error) {
	// TODO: SELECT * FROM meta_ads_campaign WHERE tenant_id = $1 AND meta_ad_account_id = $2
	return nil, nil
}

func updateMetaAdsCampaignStatus(ctx context.Context, campaign *MetaAdsCampaign) error {
	// TODO: UPDATE meta_ads_campaign SET status = $1, updated_at = $2 WHERE id = $3
	return nil
}

func storeMetaAdsAdSet(ctx context.Context, adSet *MetaAdsAdSet) error {
	// TODO: INSERT INTO meta_ads_ad_set (...) VALUES (...)
	return nil
}

func storeMetaAdsMetrics(ctx context.Context, metrics *MetaAdsMetrics) error {
	// TODO: INSERT INTO meta_ads_metrics (...) VALUES (...) ON DUPLICATE KEY UPDATE
	return nil
}

func getAdAccountConnectionByMetaID(ctx context.Context, metaAdAccountID string) (*AdAccountConnection, error) {
	// TODO: SELECT * FROM ad_account_connection WHERE account_id = $1 AND platform = 'meta'
	return nil, nil
}

// ============================================================
// HELPER FUNCTIONS (Meta Marketing API calls)
// ============================================================

func createMetaAdsCampaignViaAPI(ctx context.Context, conn *AdAccountConnection, req *CreateMetaCampaignRequest) (string, error) {
	// TODO: Use Facebook Marketing API
	// POST /act_{AD_ACCOUNT_ID}/campaigns
	// Return campaign ID
	return "", nil
}

func createMetaAdsAdSetViaAPI(ctx context.Context, conn *AdAccountConnection, campaign *MetaAdsCampaign, req *CreateMetaAdSetRequest) (string, error) {
	// TODO: Use Facebook Marketing API
	// POST /act_{AD_ACCOUNT_ID}/adsets
	// Return ad set ID
	return "", nil
}

func pauseMetaAdsCampaignViaAPI(ctx context.Context, conn *AdAccountConnection, metaCampaignID string) error {
	// TODO: Use Facebook Marketing API
	// POST /{CAMPAIGN_ID}?status=PAUSED
	return nil
}

func updateMetaAdsBudgetViaAPI(ctx context.Context, conn *AdAccountConnection, metaCampaignID, budgetType string, budgetMicros int64) error {
	// TODO: Use Facebook Marketing API
	// POST /{CAMPAIGN_ID} with daily_budget or lifetime_budget
	return nil
}

func fetchMetaAdsMetricsViaAPI(ctx context.Context, conn *AdAccountConnection, campaign *MetaAdsCampaign, startDate, endDate time.Time) ([]*MetaAdsMetrics, error) {
	// TODO: Use Facebook Insights API
	// GET /{CAMPAIGN_ID}/insights?fields=impressions,clicks,spend,actions,etc
	// Return slice of metrics
	return nil, nil
}

// ============================================================
// REGISTER ROUTES
// ============================================================

func registerMetaAdsRoutes(router *gin.Engine) {
	// Control endpoints (ERP → Meta)
	router.POST("/api/v1/meta-ads/campaigns", createMetaAdsCampaign)
	router.GET("/api/v1/meta-ads/campaigns", getMetaAdsCampaigns)
	router.POST("/api/v1/meta-ads/campaigns/:campaign_id/pause", pauseMetaAdsCampaign)
	router.PATCH("/api/v1/meta-ads/campaigns/:campaign_id/budget", updateMetaAdsBudget)
	router.POST("/api/v1/meta-ads/adsets", createMetaAdsAdSet)

	// Sync endpoints (Meta → ERP)
	router.POST("/api/v1/meta-ads/sync/metrics", syncMetaAdsMetrics)
}
