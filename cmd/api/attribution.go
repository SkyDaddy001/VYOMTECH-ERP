package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================================
// ATTRIBUTION EVENT INGESTION
// ============================================================

// AttributionEventPayload is the incoming request structure from all channels
type AttributionEventPayload struct {
	TenantID       string           `json:"tenant_id" binding:"required"`
	LeadID         string           `json:"lead_id" binding:"required"`
	TouchType      string           `json:"touch_type" binding:"required"` // "click" | "form" | "call" | "whatsapp" | "visit"
	Source         string           `json:"source" binding:"required"`     // "google_ads" | "facebook" | "phone"
	SubSource      *string          `json:"sub_source"`
	Medium         *string          `json:"medium"`
	Campaign       *string          `json:"campaign"`
	Project        *string          `json:"project"`
	AdID           *string          `json:"ad_id"`
	CreativeID     *string          `json:"creative_id"`
	Placement      *string          `json:"placement"`
	Keyword        *string          `json:"keyword"`
	LandingPage    *string          `json:"landing_page"`
	Referrer       *string          `json:"referrer"`
	UTMSource      *string          `json:"utm_source"`
	UTMMedium      *string          `json:"utm_medium"`
	UTMCampaign    *string          `json:"utm_campaign"`
	UTMContent     *string          `json:"utm_content"`
	UTMTerm        *string          `json:"utm_term"`
	Device         *string          `json:"device"` // "mobile" | "desktop"
	UserAgent      *string          `json:"user_agent"`
	IPAddress      *string          `json:"ip_address"`
	Country        *string          `json:"country"` // "AE" | "US"
	City           *string          `json:"city"`
	Timezone       *string          `json:"timezone"`
	CustomPayload  *json.RawMessage `json:"custom_payload"`
	OccurredAt     time.Time        `json:"occurred_at" binding:"required"`
	IdempotencyKey string           `json:"idempotency_key" binding:"required"` // Unique per channel
	SessionID      *string          `json:"session_id"`
}

// attributionEventHandler ingests attribution events from all channels
// This is the core endpoint that ALL channels call
func attributionEventHandler(c *gin.Context) {
	var payload AttributionEventPayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if err := validateAttributionEvent(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Idempotency check: ensure this exact event hasn't been ingested
	existingEvent, err := checkDuplicateEvent(c.Request.Context(), payload.TenantID, payload.IdempotencyKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "idempotency check failed"})
		return
	}

	if existingEvent != nil {
		// Event already processed, return success (idempotent)
		c.JSON(http.StatusOK, gin.H{
			"id":      existingEvent.ID,
			"status":  "already_ingested",
			"message": "This event has already been processed",
		})
		return
	}

	// Calculate touch order (sequence number for this lead)
	touchOrder, err := calculateTouchOrder(c.Request.Context(), payload.TenantID, payload.LeadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate touch order"})
		return
	}

	// Create the immutable attribution event
	event := &AttributionEvent{
		ID:             uuid.NewString(),
		TenantID:       payload.TenantID,
		LeadID:         payload.LeadID,
		TouchType:      payload.TouchType,
		Source:         payload.Source,
		SubSource:      payload.SubSource,
		Medium:         payload.Medium,
		Campaign:       payload.Campaign,
		Project:        payload.Project,
		AdID:           payload.AdID,
		CreativeID:     payload.CreativeID,
		Placement:      payload.Placement,
		Keyword:        payload.Keyword,
		LandingPage:    payload.LandingPage,
		Referrer:       payload.Referrer,
		UTMSource:      payload.UTMSource,
		UTMMedium:      payload.UTMMedium,
		UTMCampaign:    payload.UTMCampaign,
		UTMContent:     payload.UTMContent,
		UTMTerm:        payload.UTMTerm,
		Device:         payload.Device,
		UserAgent:      payload.UserAgent,
		IPAddress:      payload.IPAddress,
		Country:        payload.Country,
		City:           payload.City,
		Timezone:       payload.Timezone,
		CustomPayload:  payload.CustomPayload,
		TouchOrder:     touchOrder,
		OccurredAt:     payload.OccurredAt,
		IdempotencyKey: payload.IdempotencyKey,
		SessionID:      payload.SessionID,
		IngestedAt:     time.Now(),
		CreatedAt:      time.Now(),
	}

	// Calculate ingestion delay
	ingestionDelay := int(event.IngestedAt.Sub(payload.OccurredAt).Milliseconds())
	event.IngestionDelay = &ingestionDelay

	// Store the event (append-only)
	if err := storeAttributionEvent(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store event"})
		return
	}

	// Trigger snapshot recalculation (async or sync depending on load)
	// For now, we recalculate synchronously (can be moved to job queue)
	go recalculateLeadAttributionSnapshot(context.Background(), payload.TenantID, payload.LeadID)

	c.JSON(http.StatusCreated, gin.H{
		"id":          event.ID,
		"touch_order": event.TouchOrder,
		"status":      "ingested",
		"occurred_at": event.OccurredAt,
		"ingested_at": event.IngestedAt,
	})
}

// ============================================================
// ATTRIBUTION RULES ENGINE
// ============================================================

// AttributionModel defines how credit is distributed
type AttributionModel string

const (
	FirstTouch     AttributionModel = "first-touch"
	LastTouch      AttributionModel = "last-touch"
	PositionBased  AttributionModel = "position-based"  // 40-20-40
	LinearTouch    AttributionModel = "linear"          // Equal credit
	TimeDecay      AttributionModel = "time-decay"      // Recent touches weighted more
	CustomWeighted AttributionModel = "custom-weighted" // User-defined weights
)

// AttributionConfig is the rule configuration
type AttributionConfig struct {
	AttributionModel AttributionModel   `json:"attribution_model"`
	Weights          map[string]float64 `json:"weights,omitempty"`          // {first: 0.4, middle: 0.2, last: 0.4}
	Priorities       []SourcePriority   `json:"priorities,omitempty"`       // [{source: "paid", weight: 1.2}]
	ExcludedSources  []string           `json:"excluded_sources,omitempty"` // ["direct", "organic"]
	DecayDays        int                `json:"decay_days,omitempty"`       // For time-decay model
}

// SourcePriority defines relative importance of a source
type SourcePriority struct {
	Source string  `json:"source"`
	Weight float64 `json:"weight"`
}

// CalculateAttribution determines first-touch, last-touch, and assisted
// This function is called when generating reports or dashboards
func calculateAttribution(ctx context.Context, tenantID, leadID string, config *AttributionConfig) (*AttributionResult, error) {
	// Fetch all events for this lead (in order)
	events, err := getLeadAttributionEvents(ctx, tenantID, leadID)
	if err != nil {
		return nil, err
	}

	result := &AttributionResult{
		LeadID:       leadID,
		TotalTouches: len(events),
	}

	if len(events) == 0 {
		return result, nil
	}

	// Sort by touch order (should already be ordered, but ensure it)
	sort.Slice(events, func(i, j int) bool {
		return events[i].TouchOrder < events[j].TouchOrder
	})

	// First touch is always the first event
	firstEvent := events[0]
	result.FirstTouch = &TouchAttribution{
		EventID:    firstEvent.ID,
		OccurredAt: firstEvent.OccurredAt,
		Source:     firstEvent.Source,
		SubSource:  firstEvent.SubSource,
		Medium:     firstEvent.Medium,
		Campaign:   firstEvent.Campaign,
		Credit:     1.0,
	}

	// Last touch is always the last event
	lastEvent := events[len(events)-1]
	result.LastTouch = &TouchAttribution{
		EventID:    lastEvent.ID,
		OccurredAt: lastEvent.OccurredAt,
		Source:     lastEvent.Source,
		SubSource:  lastEvent.SubSource,
		Medium:     lastEvent.Medium,
		Campaign:   lastEvent.Campaign,
		Credit:     1.0,
	}

	// Calculate assisted touches (all except first and last)
	if len(events) > 2 {
		for _, event := range events[1 : len(events)-1] {
			weight := getSourcePriority(config.Priorities, event.Source)
			result.AssistedTouches = append(result.AssistedTouches, &TouchAttribution{
				EventID:    event.ID,
				OccurredAt: event.OccurredAt,
				Source:     event.Source,
				SubSource:  event.SubSource,
				Medium:     event.Medium,
				Campaign:   event.Campaign,
				Credit:     weight,
			})
		}
	}

	// Apply attribution model weights
	result.Credits = applyAttributionModel(events, config)

	// Time to conversion
	if len(events) > 0 {
		lastTouch := events[len(events)-1]
		firstTouch := events[0]
		daysToConversion := int(lastTouch.OccurredAt.Sub(firstTouch.OccurredAt).Hours() / 24)
		result.DaysToConversion = &daysToConversion
	}

	return result, nil
}

// applyAttributionModel distributes credit based on the chosen model
func applyAttributionModel(events []*AttributionEventRecord, config *AttributionConfig) map[string]*CreditDistribution {
	credits := make(map[string]*CreditDistribution)

	// Initialize credits for each event
	for _, event := range events {
		credits[event.ID] = &CreditDistribution{
			EventID:  event.ID,
			Source:   event.Source,
			Campaign: event.Campaign,
			Credit:   0.0,
		}
	}

	switch config.AttributionModel {
	case FirstTouch:
		// All credit to first touch
		if len(events) > 0 {
			credits[events[0].ID].Credit = 1.0
		}

	case LastTouch:
		// All credit to last touch
		if len(events) > 0 {
			credits[events[len(events)-1].ID].Credit = 1.0
		}

	case PositionBased:
		// 40-20-40: 40% first, 40% last, 20% divided among middle
		if len(events) == 1 {
			credits[events[0].ID].Credit = 1.0
		} else if len(events) == 2 {
			credits[events[0].ID].Credit = 0.5
			credits[events[1].ID].Credit = 0.5
		} else {
			credits[events[0].ID].Credit = 0.4
			credits[events[len(events)-1].ID].Credit = 0.4
			middleCredit := 0.2 / float64(len(events)-2)
			for i := 1; i < len(events)-1; i++ {
				credits[events[i].ID].Credit = middleCredit
			}
		}

	case LinearTouch:
		// Equal credit to all
		equalCredit := 1.0 / float64(len(events))
		for _, event := range events {
			credits[event.ID].Credit = equalCredit
		}

	case TimeDecay:
		// Recent touches weighted more
		decayDays := float64(config.DecayDays)
		if decayDays == 0 {
			decayDays = 7 // Default: 7-day decay
		}
		totalWeight := 0.0
		weights := make([]float64, len(events))

		for i, event := range events {
			daysSince := math.Max(0, float64(events[len(events)-1].OccurredAt.Sub(event.OccurredAt).Hours()/24))
			weight := math.Exp(-daysSince / decayDays) // Exponential decay
			weights[i] = weight
			totalWeight += weight
		}

		// Normalize
		for i, weight := range weights {
			credits[events[i].ID].Credit = weight / totalWeight
		}

	case CustomWeighted:
		// User-defined weights
		totalWeight := 0.0
		weights := make([]float64, len(events))

		for i, event := range events {
			priority := getSourcePriority(config.Priorities, event.Source)
			weights[i] = priority
			totalWeight += priority
		}

		// Normalize
		if totalWeight > 0 {
			for i, weight := range weights {
				credits[events[i].ID].Credit = weight / totalWeight
			}
		}
	}

	return credits
}

// ============================================================
// SNAPSHOT RECALCULATION (DERIVED FROM EVENTS)
// ============================================================

// recalculateLeadAttributionSnapshot derives attribution snapshot from events
// This is triggered after every new event (async or batch)
func recalculateLeadAttributionSnapshot(ctx context.Context, tenantID, leadID string) error {
	// Fetch all events for this lead
	events, err := getLeadAttributionEvents(ctx, tenantID, leadID)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		// No events yet, snapshot stays empty
		return nil
	}

	// Sort by touch order
	sort.Slice(events, func(i, j int) bool {
		return events[i].TouchOrder < events[j].TouchOrder
	})

	// Calculate derived metrics
	snapshot := &LeadAttributionSnapshot{
		ID:                 uuid.NewString(),
		TenantID:           tenantID,
		LeadID:             leadID,
		TotalTouches:       len(events),
		UniqueSources:      countUniqueSources(events),
		UniqueCampaigns:    countUniqueCampaigns(events),
		PrimaryChannel:     findPrimaryChannel(events),
		TouchSequence:      buildTouchSequence(events),
		LastRecalculatedAt: time.Now(),
		CalculationVersion: 1,
	}

	// First touch
	firstEvent := events[0]
	snapshot.FirstTouchEventID = &firstEvent.ID
	snapshot.FirstTouchAt = &firstEvent.OccurredAt
	snapshot.FirstSource = &firstEvent.Source
	snapshot.FirstSubSource = firstEvent.SubSource
	snapshot.FirstMedium = firstEvent.Medium
	snapshot.FirstCampaign = firstEvent.Campaign
	snapshot.FirstProject = firstEvent.Project
	snapshot.FirstLandingPage = firstEvent.LandingPage
	snapshot.FirstDevice = firstEvent.Device

	// Last touch
	lastEvent := events[len(events)-1]
	snapshot.LastTouchEventID = &lastEvent.ID
	snapshot.LastTouchAt = &lastEvent.OccurredAt
	snapshot.LastSource = &lastEvent.Source
	snapshot.LastSubSource = lastEvent.SubSource
	snapshot.LastMedium = lastEvent.Medium
	snapshot.LastCampaign = lastEvent.Campaign
	snapshot.LastProject = lastEvent.Project
	snapshot.LastDevice = lastEvent.Device

	// Days to conversion
	if len(events) > 0 {
		daysToConv := int(lastEvent.OccurredAt.Sub(firstEvent.OccurredAt).Hours() / 24)
		snapshot.DaysToConversion = &daysToConv
	}

	// Assisted touches (all except first and last)
	if len(events) > 2 {
		assistedList := make([]map[string]interface{}, 0)
		for _, event := range events[1 : len(events)-1] {
			assisted := map[string]interface{}{
				"source":      event.Source,
				"touch_type":  event.TouchType,
				"occurred_at": event.OccurredAt,
				"campaign":    event.Campaign,
			}
			assistedList = append(assistedList, assisted)
		}
		if len(assistedList) > 0 {
			assistedJSON, _ := json.Marshal(assistedList)
			snapshot.AssistedTouches = (*json.RawMessage)(&assistedJSON)
		}
	}

	// Store or update snapshot
	return storeLeadAttributionSnapshot(ctx, snapshot)
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func validateAttributionEvent(p *AttributionEventPayload) error {
	validTouchTypes := map[string]bool{
		"click": true, "form": true, "call": true, "whatsapp": true, "visit": true, "email": true,
	}
	if !validTouchTypes[p.TouchType] {
		return fmt.Errorf("invalid touch_type: %s", p.TouchType)
	}

	if p.OccurredAt.IsZero() {
		return fmt.Errorf("occurred_at is required and must be a valid timestamp")
	}

	if p.IdempotencyKey == "" {
		return fmt.Errorf("idempotency_key is required")
	}

	return nil
}

func checkDuplicateEvent(ctx context.Context, tenantID, idempotencyKey string) (*AttributionEventRecord, error) {
	// Query: SELECT * FROM attribution_event WHERE tenant_id = ? AND idempotency_key = ? AND deleted_at IS NULL
	// This would be a database query in real implementation
	// For now, return nil (no duplicate found)
	return nil, nil
}

func calculateTouchOrder(ctx context.Context, tenantID, leadID string) (int, error) {
	// Query: SELECT MAX(touch_order) FROM attribution_event WHERE tenant_id = ? AND lead_id = ?
	// In real implementation, this queries the DB
	// For MVP, return incrementing order
	return 1, nil
}

func storeAttributionEvent(ctx context.Context, event *AttributionEvent) error {
	// INSERT INTO attribution_event (id, tenant_id, lead_id, touch_type, source, ...) VALUES (...)
	// This is a database operation
	log.Printf("Storing attribution event: %s for lead: %s", event.ID, event.LeadID)
	return nil
}

func getLeadAttributionEvents(ctx context.Context, tenantID, leadID string) ([]*AttributionEventRecord, error) {
	// Query: SELECT * FROM attribution_event WHERE tenant_id = ? AND lead_id = ? AND deleted_at IS NULL ORDER BY touch_order
	// This is a database operation
	return make([]*AttributionEventRecord, 0), nil
}

func getSourcePriority(priorities []SourcePriority, source string) float64 {
	for _, p := range priorities {
		if p.Source == source {
			return p.Weight
		}
	}
	return 1.0 // Default: equal weight
}

func countUniqueSources(events []*AttributionEventRecord) int {
	sources := make(map[string]bool)
	for _, event := range events {
		sources[event.Source] = true
	}
	return len(sources)
}

func countUniqueCampaigns(events []*AttributionEventRecord) int {
	campaigns := make(map[string]bool)
	for _, event := range events {
		if event.Campaign != nil {
			campaigns[*event.Campaign] = true
		}
	}
	return len(campaigns)
}

func findPrimaryChannel(events []*AttributionEventRecord) string {
	channelCount := make(map[string]int)
	for _, event := range events {
		channelCount[event.Source]++
	}

	maxCount := 0
	primaryChannel := ""
	for channel, count := range channelCount {
		if count > maxCount {
			maxCount = count
			primaryChannel = channel
		}
	}
	return primaryChannel
}

func buildTouchSequence(events []*AttributionEventRecord) string {
	sequence := ""
	for i, event := range events {
		if i > 0 {
			sequence += ","
		}
		sequence += event.TouchType
	}
	return sequence
}

func storeLeadAttributionSnapshot(ctx context.Context, snapshot *LeadAttributionSnapshot) error {
	// INSERT OR UPDATE lead_attribution_snapshot
	log.Printf("Storing snapshot for lead: %s", snapshot.LeadID)
	return nil
}

// ============================================================
// DATA STRUCTURES
// ============================================================

type AttributionEvent struct {
	ID             string
	TenantID       string
	LeadID         string
	TouchType      string
	Source         string
	SubSource      *string
	Medium         *string
	Campaign       *string
	Project        *string
	AdID           *string
	CreativeID     *string
	Placement      *string
	Keyword        *string
	LandingPage    *string
	Referrer       *string
	UTMSource      *string
	UTMMedium      *string
	UTMCampaign    *string
	UTMContent     *string
	UTMTerm        *string
	Device         *string
	UserAgent      *string
	IPAddress      *string
	Country        *string
	City           *string
	Timezone       *string
	CustomPayload  *json.RawMessage
	TouchOrder     int
	OccurredAt     time.Time
	IngestionDelay *int
	IdempotencyKey string
	SessionID      *string
	DeletedAt      *time.Time
	IngestedAt     time.Time
	CreatedAt      time.Time
}

type AttributionEventRecord struct {
	ID          string
	TenantID    string
	LeadID      string
	TouchType   string
	Source      string
	SubSource   *string
	Medium      *string
	Campaign    *string
	Project     *string
	Device      *string
	OccurredAt  time.Time
	LandingPage *string
	TouchOrder  int
}

type LeadAttributionSnapshot struct {
	ID                 string
	TenantID           string
	LeadID             string
	FirstTouchEventID  *string
	FirstTouchAt       *time.Time
	FirstSource        *string
	FirstSubSource     *string
	FirstMedium        *string
	FirstCampaign      *string
	FirstProject       *string
	FirstLandingPage   *string
	FirstDevice        *string
	LastTouchEventID   *string
	LastTouchAt        *time.Time
	LastSource         *string
	LastSubSource      *string
	LastMedium         *string
	LastCampaign       *string
	LastProject        *string
	LastDevice         *string
	TotalTouches       int
	UniqueSources      int
	UniqueCampaigns    int
	PrimaryChannel     string
	DaysToConversion   *int
	TouchSequence      string
	AssistedTouches    *json.RawMessage
	LastRecalculatedAt time.Time
	CalculationVersion int
}

type AttributionResult struct {
	LeadID           string
	TotalTouches     int
	FirstTouch       *TouchAttribution
	LastTouch        *TouchAttribution
	AssistedTouches  []*TouchAttribution
	Credits          map[string]*CreditDistribution
	DaysToConversion *int
}

type TouchAttribution struct {
	EventID    string
	OccurredAt time.Time
	Source     string
	SubSource  *string
	Medium     *string
	Campaign   *string
	Credit     float64
}

type CreditDistribution struct {
	EventID  string
	Source   string
	Campaign *string
	Credit   float64
}

// ============================================================
// API ENDPOINTS
// ============================================================

// registerAttributionRoutes registers all attribution endpoints
func registerAttributionRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	// Core attribution ingestion (all channels call this)
	api.POST("/attribution/events", attributionEventHandler)

	// Analytics queries (internal use)
	api.GET("/attribution/snapshot/:leadId", getAttributionSnapshotHandler)
	api.GET("/attribution/events/:leadId", getLeadAttributionEventsHandler)
	api.GET("/attribution/report/:leadId", getAttributionReportHandler)

	// Rules management (admin only)
	api.POST("/attribution/rules", createAttributionRuleHandler)
	api.GET("/attribution/rules", listAttributionRulesHandler)
	api.PUT("/attribution/rules/:ruleId", updateAttributionRuleHandler)

	// Snapshot recalculation (maintenance)
	api.POST("/attribution/snapshot/recalculate", recalculateSnapshotsHandler)

	log.Println("✅ Attribution routes registered")
}

// getAttributionSnapshotHandler returns the derived snapshot for a lead
func getAttributionSnapshotHandler(c *gin.Context) {
	leadID := c.Param("leadId")
	tenantID := c.GetString("tenantID")

	// Query snapshot from DB
	snapshot, err := getLeadSnapshotFromDB(c.Request.Context(), tenantID, leadID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "snapshot not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch snapshot"})
		}
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

// getLeadAttributionEventsHandler returns all events for a lead
func getLeadAttributionEventsHandler(c *gin.Context) {
	leadID := c.Param("leadId")
	tenantID := c.GetString("tenantID")

	events, err := getLeadAttributionEvents(c.Request.Context(), tenantID, leadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"lead_id": leadID,
		"count":   len(events),
		"events":  events,
	})
}

// getAttributionReportHandler returns first-touch, last-touch, and assisted attribution
func getAttributionReportHandler(c *gin.Context) {
	leadID := c.Param("leadId")
	tenantID := c.GetString("tenantID")

	// Get active attribution rule
	rule, err := getActiveAttributionRule(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch attribution rule"})
		return
	}

	// Calculate attribution
	result, err := calculateAttribution(c.Request.Context(), tenantID, leadID, rule.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate attribution"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// createAttributionRuleHandler creates a new attribution rule
func createAttributionRuleHandler(c *gin.Context) {
	var req struct {
		Name        string            `json:"name" binding:"required"`
		Description string            `json:"description"`
		RuleConfig  AttributionConfig `json:"rule_config" binding:"required"`
		IsActive    bool              `json:"is_active"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := c.GetString("tenantID")
	userID := c.GetString("userID")

	rule := &struct {
		ID        string
		TenantID  string
		Name      string
		Config    *AttributionConfig
		IsActive  bool
		CreatedBy string
		CreatedAt time.Time
	}{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Name:      req.Name,
		Config:    &req.RuleConfig,
		IsActive:  req.IsActive,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}

	// Store rule in DB
	if err := storeAttributionRule(c.Request.Context(), rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule"})
		return
	}

	// Trigger snapshot recalculation for all leads
	go recalculateAllSnapshots(context.Background(), tenantID)

	c.JSON(http.StatusCreated, rule)
}

// listAttributionRulesHandler lists all rules for a tenant
func listAttributionRulesHandler(c *gin.Context) {
	tenantID := c.GetString("tenantID")

	rules, err := listAttributionRulesFromDB(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch rules"})
		return
	}

	ruleSlice := rules.([]interface{})
	c.JSON(http.StatusOK, gin.H{
		"count": len(ruleSlice),
		"rules": ruleSlice,
	})
}

// updateAttributionRuleHandler updates an attribution rule
func updateAttributionRuleHandler(c *gin.Context) {
	ruleID := c.Param("ruleId")
	tenantID := c.GetString("tenantID")

	var req struct {
		Name       string            `json:"name"`
		RuleConfig AttributionConfig `json:"rule_config"`
		IsActive   bool              `json:"is_active"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update rule in DB
	if err := updateAttributionRuleInDB(c.Request.Context(), tenantID, ruleID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rule"})
		return
	}

	// Trigger snapshot recalculation (rule changed, snapshots need update)
	go recalculateAllSnapshots(context.Background(), tenantID)

	c.JSON(http.StatusOK, gin.H{"message": "rule updated, snapshots will be recalculated"})
}

// recalculateSnapshotsHandler manually triggers snapshot recalculation
func recalculateSnapshotsHandler(c *gin.Context) {
	tenantID := c.GetString("tenantID")

	go recalculateAllSnapshots(context.Background(), tenantID)

	c.JSON(http.StatusAccepted, gin.H{"message": "snapshot recalculation queued"})
}

// ============================================================
// DATABASE HELPERS (Stubs - implement with real queries)
// ============================================================

func getLeadSnapshotFromDB(ctx context.Context, tenantID, leadID string) (*LeadAttributionSnapshot, error) {
	// Query: SELECT * FROM lead_attribution_snapshot WHERE tenant_id = ? AND lead_id = ?
	return &LeadAttributionSnapshot{}, nil
}

func storeAttributionRule(ctx context.Context, rule interface{}) error {
	// INSERT INTO attribution_rule (...)
	return nil
}

func getActiveAttributionRule(ctx context.Context, tenantID string) (*struct {
	Config *AttributionConfig
}, error) {
	// Query active rule, return with decoded config
	return &struct {
		Config *AttributionConfig
	}{
		Config: &AttributionConfig{
			AttributionModel: FirstTouch,
		},
	}, nil
}

func listAttributionRulesFromDB(ctx context.Context, tenantID string) (interface{}, error) {
	// Query: SELECT * FROM attribution_rule WHERE tenant_id = ? ORDER BY priority DESC
	return []interface{}{}, nil
}

func updateAttributionRuleInDB(ctx context.Context, tenantID, ruleID string, req interface{}) error {
	// UPDATE attribution_rule SET ... WHERE tenant_id = ? AND id = ?
	return nil
}

func recalculateAllSnapshots(ctx context.Context, tenantID string) error {
	// Query all leads for tenant
	// For each lead, call recalculateLeadAttributionSnapshot
	log.Printf("Recalculating snapshots for tenant: %s", tenantID)
	return nil
}
