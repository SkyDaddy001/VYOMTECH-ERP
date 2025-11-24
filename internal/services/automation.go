package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// LeadScore represents a lead's score and ranking
type LeadScore struct {
	LeadID    int64              `json:"lead_id"`
	Score     float64            `json:"score"`
	Rank      int                `json:"rank"`
	Factors   map[string]float64 `json:"factors"`
	CreatedAt time.Time          `json:"created_at"`
}

// RoutingRule defines how leads should be routed to agents
type RoutingRule struct {
	ID          int64                  `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Name        string                 `json:"name"`
	Priority    int                    `json:"priority"`
	Conditions  map[string]interface{} `json:"conditions"`
	ActionType  string                 `json:"action_type"` // assign_to_agent, assign_to_team, round_robin
	ActionValue string                 `json:"action_value"`
	Enabled     bool                   `json:"enabled"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// WorkflowTask represents a task in an automation workflow
type WorkflowTask struct {
	ID         int64                  `json:"id"`
	WorkflowID int64                  `json:"workflow_id"`
	Trigger    string                 `json:"trigger"` // lead_created, call_ended, etc
	Action     string                 `json:"action"`  // send_email, create_call, etc
	Config     map[string]interface{} `json:"config"`
	Order      int                    `json:"order"`
	Enabled    bool                   `json:"enabled"`
	CreatedAt  time.Time              `json:"created_at"`
}

// ScheduledCampaign represents a campaign scheduled for automation
type ScheduledCampaign struct {
	ID            int64      `json:"id"`
	CampaignID    int64      `json:"campaign_id"`
	TenantID      string     `json:"tenant_id"`
	ScheduledTime time.Time  `json:"scheduled_time"`
	Status        string     `json:"status"` // pending, running, completed, failed
	ExecutedAt    *time.Time `json:"executed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// AutomationService handles lead scoring, routing, and workflows
type AutomationService struct {
	db *sql.DB
}

// NewAutomationService creates a new AutomationService
func NewAutomationService(db *sql.DB) *AutomationService {
	return &AutomationService{
		db: db,
	}
}

// CalculateLeadScore calculates the quality score for a lead
// Scoring factors: lead source quality, email validity, phone validity, engagement, conversion likelihood
func (as *AutomationService) CalculateLeadScore(ctx context.Context, tenantID string, leadID int64) (*LeadScore, error) {
	score := &LeadScore{
		LeadID:  leadID,
		Score:   0,
		Factors: make(map[string]float64),
	}

	// Query lead data
	query := `
		SELECT source, status, assigned_agent_id, created_at 
		FROM lead 
		WHERE id = ? AND tenant_id = ?
	`
	var source, status string
	var agentID sql.NullInt64
	var createdAt time.Time

	err := as.db.QueryRowContext(ctx, query, leadID, tenantID).
		Scan(&source, &status, &agentID, &createdAt)
	if err != nil {
		return nil, err
	}

	// Source quality scoring (0-25 points)
	sourceScore := as.getSourceQualityScore(source)
	score.Factors["source_quality"] = sourceScore

	// Status score (0-25 points)
	statusScore := as.getStatusScore(status)
	score.Factors["status_score"] = statusScore

	// Engagement score (0-25 points) - based on days since creation
	daysOld := time.Since(createdAt).Hours() / 24
	engagementScore := as.calculateEngagementScore(daysOld)
	score.Factors["engagement"] = engagementScore

	// Assignment score (0-25 points) - unassigned leads get bonus points
	assignmentScore := 0.0
	if !agentID.Valid {
		assignmentScore = 25.0 // Unassigned leads get highest score
	}
	score.Factors["assignment"] = assignmentScore

	// Calculate total score
	for _, value := range score.Factors {
		score.Score += value
	}

	// Normalize to 0-100
	score.Score = (score.Score / 100) * 100
	score.CreatedAt = time.Now()

	return score, nil
}

// getSourceQualityScore returns quality score for lead source
func (as *AutomationService) getSourceQualityScore(source string) float64 {
	scoreMap := map[string]float64{
		"campaign": 25.0,
		"referral": 20.0,
		"organic":  18.0,
		"paid_ads": 15.0,
		"import":   10.0,
		"manual":   5.0,
	}
	if score, ok := scoreMap[source]; ok {
		return score
	}
	return 5.0
}

// getStatusScore returns score based on lead status
func (as *AutomationService) getStatusScore(status string) float64 {
	scoreMap := map[string]float64{
		"new":       25.0,
		"contacted": 20.0,
		"qualified": 15.0,
		"converted": 10.0,
		"lost":      0.0,
	}
	if score, ok := scoreMap[status]; ok {
		return score
	}
	return 10.0
}

// calculateEngagementScore scores based on lead age
func (as *AutomationService) calculateEngagementScore(daysOld float64) float64 {
	if daysOld <= 1 {
		return 25.0 // Very recent - highest priority
	} else if daysOld <= 7 {
		return 20.0 // Recent
	} else if daysOld <= 30 {
		return 10.0 // Moderate age
	} else {
		return 0.0 // Stale lead
	}
}

// RankLeads ranks all leads by score for a tenant
func (as *AutomationService) RankLeads(ctx context.Context, tenantID string, limit int) ([]LeadScore, error) {
	scores := make([]LeadScore, 0)

	// In a real implementation, scores would be cached/precomputed
	// For now, get all leads and score them
	query := `SELECT id FROM lead WHERE tenant_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := as.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rank := 1
	for rows.Next() {
		var leadID int64
		if err := rows.Scan(&leadID); err != nil {
			continue
		}

		score, err := as.CalculateLeadScore(ctx, tenantID, leadID)
		if err == nil {
			score.Rank = rank
			scores = append(scores, *score)
			rank++
		}
	}

	return scores, nil
}

// RouteLeadToAgent routes a lead to an agent based on rules
func (as *AutomationService) RouteLeadToAgent(ctx context.Context, tenantID string, leadID int64) (int64, error) {
	// Get applicable routing rules
	query := `
		SELECT id, action_type, action_value FROM routing_rule
		WHERE tenant_id = ? AND enabled = true
		ORDER BY priority DESC
		LIMIT 1
	`
	var ruleID int64
	var actionType, actionValue string

	err := as.db.QueryRowContext(ctx, query, tenantID).Scan(&ruleID, &actionType, &actionValue)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	// If no rule found, use simple round-robin
	if err == sql.ErrNoRows {
		return as.routeByRoundRobin(ctx, tenantID)
	}

	// Apply routing rule based on action type
	switch actionType {
	case "assign_to_agent":
		// Parse agent ID from action value
		var agentID int64
		fmt.Sscanf(actionValue, "%d", &agentID)
		return agentID, nil

	case "assign_to_team":
		// Get agent with lowest workload from team
		return as.getAvailableAgentInTeam(ctx, tenantID, actionValue)

	case "round_robin":
		return as.routeByRoundRobin(ctx, tenantID)

	default:
		return as.routeByRoundRobin(ctx, tenantID)
	}
}

// routeByRoundRobin routes lead using round-robin method
func (as *AutomationService) routeByRoundRobin(ctx context.Context, tenantID string) (int64, error) {
	// Get agent with least assigned leads
	query := `
		SELECT a.id FROM agent a
		LEFT JOIN lead l ON a.id = l.assigned_agent_id
		WHERE a.tenant_id = ?
		GROUP BY a.id
		ORDER BY COUNT(l.id) ASC
		LIMIT 1
	`
	var agentID int64
	err := as.db.QueryRowContext(ctx, query, tenantID).Scan(&agentID)
	return agentID, err
}

// getAvailableAgentInTeam gets an available agent from a team
func (as *AutomationService) getAvailableAgentInTeam(ctx context.Context, tenantID string, _ string) (int64, error) {
	query := `
		SELECT id FROM agent
		WHERE tenant_id = ? AND availability = 'available'
		ORDER BY active_calls ASC
		LIMIT 1
	`
	var agentID int64
	err := as.db.QueryRowContext(ctx, query, tenantID).Scan(&agentID)
	return agentID, err
}

// CreateRoutingRule creates a new routing rule
func (as *AutomationService) CreateRoutingRule(ctx context.Context, rule *RoutingRule) error {
	query := `
		INSERT INTO routing_rule (tenant_id, name, priority, conditions, action_type, action_value, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := as.db.ExecContext(ctx, query,
		rule.TenantID, rule.Name, rule.Priority, fmt.Sprintf("%v", rule.Conditions),
		rule.ActionType, rule.ActionValue, rule.Enabled)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	rule.ID = id
	return err
}

// CreateWorkflow creates a new automation workflow
func (as *AutomationService) CreateWorkflow(ctx context.Context, tenantID string, name string, trigger string) (int64, error) {
	query := `
		INSERT INTO automation_workflow (tenant_id, name, trigger, enabled, created_at)
		VALUES (?, ?, ?, ?, NOW())
	`

	result, err := as.db.ExecContext(ctx, query, tenantID, name, trigger, true)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// ScheduleCampaign schedules a campaign for future execution
func (as *AutomationService) ScheduleCampaign(ctx context.Context, campaignID int64, tenantID string, scheduledTime time.Time) error {
	query := `
		INSERT INTO scheduled_campaign (campaign_id, tenant_id, scheduled_time, status, created_at)
		VALUES (?, ?, ?, 'pending', NOW())
	`

	_, err := as.db.ExecContext(ctx, query, campaignID, tenantID, scheduledTime)
	return err
}

// ExecutePendingCampaigns executes all campaigns scheduled for now
func (as *AutomationService) ExecutePendingCampaigns(ctx context.Context) error {
	query := `
		SELECT id, campaign_id FROM scheduled_campaign
		WHERE status = 'pending' AND scheduled_time <= NOW()
	`

	rows, err := as.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	updateQuery := `UPDATE scheduled_campaign SET status = 'completed', executed_at = NOW() WHERE id = ?`

	for rows.Next() {
		var schedID, campaignID int64
		if err := rows.Scan(&schedID, &campaignID); err != nil {
			continue
		}

		// Execute campaign (in real implementation, trigger campaign service)
		as.db.ExecContext(ctx, updateQuery, schedID)
	}

	return nil
}

// GetLeadScoringMetrics returns metrics for lead scoring
func (as *AutomationService) GetLeadScoringMetrics(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// Average lead score
	query := `SELECT AVG(score) FROM lead_score WHERE tenant_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)`
	var avgScore sql.NullFloat64
	as.db.QueryRowContext(ctx, query, tenantID).Scan(&avgScore)
	metrics["avg_score"] = avgScore.Float64

	// Leads routed today
	query = `SELECT COUNT(*) FROM lead WHERE tenant_id = ? AND assigned_agent_id IS NOT NULL AND DATE(created_at) = CURDATE()`
	var routed int
	as.db.QueryRowContext(ctx, query, tenantID).Scan(&routed)
	metrics["routed_today"] = routed

	// Active rules
	query = `SELECT COUNT(*) FROM routing_rule WHERE tenant_id = ? AND enabled = true`
	var activeRules int
	as.db.QueryRowContext(ctx, query, tenantID).Scan(&activeRules)
	metrics["active_rules"] = activeRules

	return metrics, nil
}
