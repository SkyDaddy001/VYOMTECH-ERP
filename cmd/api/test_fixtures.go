package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================================
// TEST FIXTURES & MOCK DATA
// ============================================================
// Pre-configured test data for development and testing
// Use these instead of real API calls

// ============================================================
// MOCK DATA STRUCTURES
// ============================================================

type MockTestFixture struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
	CreatedAt   time.Time   `json:"created_at"`
	ExpiresAt   *time.Time  `json:"expires_at"`
}

// ============================================================
// GOOGLE ADS TEST DATA
// ============================================================

var GoogleAdsTestFixtures = map[string]interface{}{
	"account_connection": map[string]interface{}{
		"id":            "conn-google-001",
		"tenant_id":     "tenant-test-001",
		"platform":      "google_ads",
		"account_id":    "123-456-7890",
		"account_name":  "Mock Google Ads Account",
		"is_active":     true,
		"sync_enabled":  true,
		"sync_status":   "success",
		"access_token":  "access_test_google_12345",
		"refresh_token": "refresh_test_google_12345",
		"last_sync_at":  time.Now(),
	},
	"campaign": map[string]interface{}{
		"id":                   "camp-google-001",
		"tenant_id":            "tenant-test-001",
		"google_campaign_id":   "123456789",
		"google_customer_id":   "123-456-7890",
		"name":                 "Test Search Campaign",
		"status":               "ENABLED",
		"campaign_type":        "SEARCH",
		"daily_budget_micros":  int64(5000000000), // 5000 in base currency
		"targeting_dimensions": json.RawMessage(`{"keywords": ["test keyword"]}`),
		"project_id":           "proj-001",
		"created_at":           time.Now(),
		"updated_at":           time.Now(),
	},
	"campaigns_list": []map[string]interface{}{
		{
			"id":                  "camp-google-001",
			"name":                "Test Search Campaign",
			"status":              "ENABLED",
			"campaign_type":       "SEARCH",
			"daily_budget_micros": int64(5000000000),
		},
		{
			"id":                  "camp-google-002",
			"name":                "Test Display Campaign",
			"status":              "PAUSED",
			"campaign_type":       "DISPLAY",
			"daily_budget_micros": int64(3000000000),
		},
		{
			"id":                  "camp-google-003",
			"name":                "Test Shopping Campaign",
			"status":              "ENABLED",
			"campaign_type":       "SHOPPING",
			"daily_budget_micros": int64(2000000000),
		},
	},
	"metrics": map[string]interface{}{
		"id":                  "metrics-google-001",
		"tenant_id":           "tenant-test-001",
		"campaign_id":         "camp-google-001",
		"date":                time.Now().Format("2006-01-02"),
		"impressions":         int64(15000),
		"clicks":              int64(450),
		"ctr":                 float64(0.03),
		"average_cpc":         int64(50000),       // micros
		"cost_micros":         int64(22500000000), // 22500
		"conversions":         float64(45),
		"conversion_rate":     float64(0.10),
		"cost_per_conversion": int64(500000000), // 500
	},
	"metrics_week": []map[string]interface{}{
		{
			"date":        "2025-12-24",
			"impressions": int64(15000),
			"clicks":      int64(450),
			"cost_micros": int64(22500000000),
		},
		{
			"date":        "2025-12-23",
			"impressions": int64(14500),
			"clicks":      int64(435),
			"cost_micros": int64(21750000000),
		},
		{
			"date":        "2025-12-22",
			"impressions": int64(16000),
			"clicks":      int64(480),
			"cost_micros": int64(24000000000),
		},
	},
}

// ============================================================
// META ADS TEST DATA
// ============================================================

var MetaAdsTestFixtures = map[string]interface{}{
	"account_connection": map[string]interface{}{
		"id":            "conn-meta-001",
		"tenant_id":     "tenant-test-001",
		"platform":      "meta",
		"account_id":    "act_1234567890",
		"account_name":  "Mock Meta Ads Account",
		"is_active":     true,
		"sync_enabled":  true,
		"sync_status":   "success",
		"access_token":  "access_test_meta_12345",
		"refresh_token": "refresh_test_meta_12345",
		"last_sync_at":  time.Now(),
	},
	"campaign": map[string]interface{}{
		"id":                   "camp-meta-001",
		"tenant_id":            "tenant-test-001",
		"meta_campaign_id":     "23456789012345",
		"meta_ad_account_id":   "act_1234567890",
		"name":                 "Test Meta Campaign",
		"status":               "ACTIVE",
		"objective":            "LINK_CLICKS",
		"daily_budget_micros":  int64(5000000000),
		"targeting_dimensions": json.RawMessage(`{"interests": ["technology"]}`),
		"project_id":           "proj-001",
		"created_at":           time.Now(),
		"updated_at":           time.Now(),
	},
	"ad_set": map[string]interface{}{
		"id":                   "adset-meta-001",
		"campaign_id":          "camp-meta-001",
		"meta_ad_set_id":       "34567890123456",
		"name":                 "Test Ad Set",
		"status":               "ACTIVE",
		"daily_budget_micros":  int64(2500000000),
		"targeting_dimensions": json.RawMessage(`{"placements": ["feed", "stories"]}`),
		"created_at":           time.Now(),
		"updated_at":           time.Now(),
	},
	"campaigns_list": []map[string]interface{}{
		{
			"id":                  "camp-meta-001",
			"name":                "Test Meta Campaign",
			"status":              "ACTIVE",
			"objective":           "LINK_CLICKS",
			"daily_budget_micros": int64(5000000000),
		},
		{
			"id":                  "camp-meta-002",
			"name":                "Retargeting Campaign",
			"status":              "ACTIVE",
			"objective":           "CONVERSIONS",
			"daily_budget_micros": int64(3000000000),
		},
	},
	"metrics": map[string]interface{}{
		"id":              "metrics-meta-001",
		"tenant_id":       "tenant-test-001",
		"campaign_id":     "camp-meta-001",
		"date":            time.Now().Format("2006-01-02"),
		"impressions":     int64(20000),
		"clicks":          int64(600),
		"ctr":             float64(0.03),
		"average_cpc":     int64(40000),
		"cost_micros":     int64(24000000000), // 24000
		"spend":           int64(24000000000),
		"web_conversions": float64(60),
		"reach":           int64(18000),
		"frequency":       float64(1.11),
	},
}

// ============================================================
// ROI CALCULATION TEST DATA
// ============================================================

var ROICalculationTestFixtures = map[string]interface{}{
	"daily_roi": map[string]interface{}{
		"id":                      "roi-2025-12-24",
		"tenant_id":               "tenant-test-001",
		"date":                    "2025-12-24",
		"date_range":              "day",
		"total_spend_micros":      int64(46500000000), // Google + Meta
		"total_leads":             int64(105),         // From attribution
		"qualified_leads":         int64(78),          // Filtered leads
		"conversions":             int64(25),          // Actual closed deals
		"total_revenue":           int64(250000000),   // 250000
		"roas":                    float64(5.38),      // 250000 / 46500
		"roi_percentage":          float64(438),       // (250000 - 46500) / 46500 * 100
		"cost_per_lead":           int64(443),
		"cost_per_qualified_lead": int64(596),
		"cost_per_conversion":     int64(1860000000), // 1860
		"attribution_breakdown": map[string]interface{}{
			"google_ads": map[string]interface{}{
				"spend_micros":     int64(22500000000),
				"attributed_leads": int64(52),
				"conversions":      int64(14),
				"revenue":          int64(140000000),
			},
			"meta_ads": map[string]interface{}{
				"spend_micros":     int64(24000000000),
				"attributed_leads": int64(53),
				"conversions":      int64(11),
				"revenue":          int64(110000000),
			},
		},
	},
	"weekly_roi": map[string]interface{}{
		"id":                 "roi-2025-w51",
		"tenant_id":          "tenant-test-001",
		"date":               "2025-12-22",
		"date_range":         "week",
		"total_spend_micros": int64(330000000000), // 330000
		"total_leads":        int64(750),
		"conversions":        int64(180),
		"total_revenue":      int64(1800000000), // 1800000
		"roas":               float64(5.45),
		"roi_percentage":     float64(445),
	},
	"campaign_roi": map[string]interface{}{
		"id":                  "roi-camp-google-001",
		"tenant_id":           "tenant-test-001",
		"campaign_id":         "camp-google-001",
		"campaign_name":       "Test Search Campaign",
		"platform":            "google_ads",
		"date_range":          "day",
		"date":                "2025-12-24",
		"spend_micros":        int64(22500000000),
		"attributed_leads":    int64(52),
		"conversions":         int64(14),
		"revenue":             int64(140000000),
		"roas":                float64(6.22),
		"roi_percentage":      float64(522),
		"cost_per_lead":       int64(433),
		"cost_per_conversion": int64(1607000000),
	},
}

// ============================================================
// ATTRIBUTION TEST DATA
// ============================================================

var AttributionTestFixtures = map[string]interface{}{
	"event": map[string]interface{}{
		"id":              "event-001",
		"tenant_id":       "tenant-test-001",
		"lead_id":         "lead-001",
		"event_type":      "page_view",
		"touch_point":     "google_ads",
		"source":          "google_ads",
		"medium":          "cpc",
		"campaign":        "Test Search Campaign",
		"utm_source":      "google",
		"utm_medium":      "cpc",
		"utm_campaign":    "test_campaign",
		"utm_content":     "ad_1",
		"utm_term":        "test keyword",
		"referrer":        "google.com",
		"landing_page":    "/products",
		"device":          "mobile",
		"user_agent":      "Mozilla/5.0",
		"ip_address":      "192.168.1.1",
		"session_id":      "session-001",
		"sequence_number": 1,
		"timestamp":       time.Now(),
	},
	"events_sequence": []map[string]interface{}{
		{
			"sequence_number": 1,
			"event_type":      "page_view",
			"source":          "google_ads",
			"campaign":        "Test Campaign",
			"timestamp":       time.Now(),
		},
		{
			"sequence_number": 2,
			"event_type":      "add_to_cart",
			"source":          "direct",
			"timestamp":       time.Now().Add(2 * time.Hour),
		},
		{
			"sequence_number": 3,
			"event_type":      "purchase",
			"source":          "direct",
			"revenue":         int64(50000000), // 50000
			"timestamp":       time.Now().Add(24 * time.Hour),
		},
	},
	"attribution_model": map[string]interface{}{
		"first_touch": map[string]interface{}{
			"source":   "google_ads",
			"campaign": "Test Campaign",
			"weight":   float64(0.5), // 50% credit
		},
		"last_touch": map[string]interface{}{
			"source":   "direct",
			"campaign": "direct",
			"weight":   float64(0.5), // 50% credit
		},
		"linear": map[string]interface{}{
			"google_ads": float64(0.33),
			"direct":     float64(0.67),
		},
	},
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

// GetGoogleAdsTestData returns a copy of Google Ads test fixture
func GetGoogleAdsTestData(key string) interface{} {
	if data, ok := GoogleAdsTestFixtures[key]; ok {
		return data
	}
	return nil
}

// GetMetaAdsTestData returns a copy of Meta Ads test fixture
func GetMetaAdsTestData(key string) interface{} {
	if data, ok := MetaAdsTestFixtures[key]; ok {
		return data
	}
	return nil
}

// GetROITestData returns a copy of ROI calculation test fixture
func GetROITestData(key string) interface{} {
	if data, ok := ROICalculationTestFixtures[key]; ok {
		return data
	}
	return nil
}

// GetAttributionTestData returns a copy of attribution test fixture
func GetAttributionTestData(key string) interface{} {
	if data, ok := AttributionTestFixtures[key]; ok {
		return data
	}
	return nil
}

// ============================================================
// TEST DATA RETRIEVAL ENDPOINTS
// ============================================================

// RegisterTestDataRoutes adds test data retrieval routes
func RegisterTestDataRoutes(router *gin.Engine) {
	testData := router.Group("/test/data")
	{
		// Google Ads test data
		testData.GET("/google-ads/:key", func(c *gin.Context) {
			key := c.Param("key")
			data := GetGoogleAdsTestData(key)
			if data == nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "test data not found",
					"key":   key,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"key":  key,
				"data": data,
			})
		})

		// Meta Ads test data
		testData.GET("/meta-ads/:key", func(c *gin.Context) {
			key := c.Param("key")
			data := GetMetaAdsTestData(key)
			if data == nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "test data not found",
					"key":   key,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"key":  key,
				"data": data,
			})
		})

		// ROI test data
		testData.GET("/roi/:key", func(c *gin.Context) {
			key := c.Param("key")
			data := GetROITestData(key)
			if data == nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "test data not found",
					"key":   key,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"key":  key,
				"data": data,
			})
		})

		// Attribution test data
		testData.GET("/attribution/:key", func(c *gin.Context) {
			key := c.Param("key")
			data := GetAttributionTestData(key)
			if data == nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "test data not found",
					"key":   key,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"key":  key,
				"data": data,
			})
		})

		// All fixtures summary
		testData.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"google_ads_keys": []string{
					"account_connection", "campaign", "campaigns_list", "metrics", "metrics_week",
				},
				"meta_ads_keys": []string{
					"account_connection", "campaign", "ad_set", "campaigns_list", "metrics",
				},
				"roi_keys": []string{
					"daily_roi", "weekly_roi", "campaign_roi",
				},
				"attribution_keys": []string{
					"event", "events_sequence", "attribution_model",
				},
				"usage": "GET /test/data/{type}/{key} where type is google-ads, meta-ads, roi, or attribution",
			})
		})
	}
}
