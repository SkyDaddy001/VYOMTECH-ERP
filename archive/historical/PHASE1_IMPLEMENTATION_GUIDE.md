# Phase 1 Implementation Guide: Critical Features

**Focus:** Agent Availability, Lead Scoring, Audit Trail  
**Estimated Timeline:** 2-3 weeks  
**Priority:** HIGH

---

## Overview

These three features provide immediate value and support the existing 6 features:
1. **Agent Availability** - Enables intelligent routing (Feature 3: Automation)
2. **Lead Scoring** - Enhances automation decisions (Feature 3: Automation)
3. **Audit Trail** - Supports compliance (Feature 6: Compliance)

---

## Feature 1: Agent Availability Management

### Database Schema

```sql
-- Database: Add to internal/storage/migrations/
CREATE TABLE IF NOT EXISTS agent_availability (
    id TEXT PRIMARY KEY DEFAULT generate_ulid(),
    user_id TEXT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'available',
    -- Status enum: available, busy, on_break, offline, in_meeting, away
    break_reason TEXT,
    is_accepting_leads BOOLEAN DEFAULT TRUE,
    total_calls_today INT DEFAULT 0,
    current_call_duration_seconds INT DEFAULT 0,
    last_status_change TIMESTAMP DEFAULT NOW(),
    last_activity TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_agent_availability_user_id ON agent_availability(user_id);
CREATE INDEX idx_agent_availability_status ON agent_availability(status);
CREATE INDEX idx_agent_availability_is_accepting ON agent_availability(is_accepting_leads);

-- Auto-update updated_at
CREATE TRIGGER trigger_agent_availability_updated_at
BEFORE UPDATE ON agent_availability
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
```

### Backend Implementation

**File: `internal/services/agent_service.go`**

```go
package services

import (
	"context"
	"errors"
	"time"

	"your-module/internal/models"
	"your-module/pkg/database"
)

type AgentService interface {
	// Update agent status
	UpdateAgentStatus(ctx context.Context, userID, status, breakReason string) error
	
	// Get agent availability
	GetAgentAvailability(ctx context.Context, userID string) (*models.AgentAvailability, error)
	
	// Get available agents (for assignment)
	GetAvailableAgents(ctx context.Context, teamID string) ([]models.AgentAvailability, error)
	
	// Update call duration
	UpdateCallDuration(ctx context.Context, userID string, durationSeconds int) error
	
	// Increment call count
	IncrementDailyCallCount(ctx context.Context, userID string) error
}

type agentService struct {
	db database.DB
}

func NewAgentService(db database.DB) AgentService {
	return &agentService{db: db}
}

// UpdateAgentStatus updates the status and last activity timestamp
func (s *agentService) UpdateAgentStatus(ctx context.Context, userID, status, breakReason string) error {
	query := `
		UPDATE agent_availability 
		SET status = $1, 
		    break_reason = $2,
		    last_status_change = NOW(),
		    last_activity = NOW()
		WHERE user_id = $3
	`
	
	_, err := s.db.ExecContext(ctx, query, status, breakReason, userID)
	return err
}

// GetAgentAvailability retrieves current agent availability
func (s *agentService) GetAgentAvailability(ctx context.Context, userID string) (*models.AgentAvailability, error) {
	agent := &models.AgentAvailability{}
	query := `
		SELECT id, user_id, status, break_reason, is_accepting_leads,
		       total_calls_today, current_call_duration_seconds,
		       last_status_change, last_activity, created_at, updated_at
		FROM agent_availability
		WHERE user_id = $1
	`
	
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&agent.ID, &agent.UserID, &agent.Status, &agent.BreakReason,
		&agent.IsAcceptingLeads, &agent.TotalCallsToday, &agent.CurrentCallDurationSeconds,
		&agent.LastStatusChange, &agent.LastActivity, &agent.CreatedAt, &agent.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	return agent, nil
}

// GetAvailableAgents returns list of available agents in a team
func (s *agentService) GetAvailableAgents(ctx context.Context, teamID string) ([]models.AgentAvailability, error) {
	agents := []models.AgentAvailability{}
	query := `
		SELECT a.id, a.user_id, a.status, a.break_reason, a.is_accepting_leads,
		       a.total_calls_today, a.current_call_duration_seconds,
		       a.last_status_change, a.last_activity, a.created_at, a.updated_at
		FROM agent_availability a
		JOIN team_members tm ON a.user_id = tm.user_id
		WHERE tm.team_id = $1 
		  AND a.status = 'available'
		  AND a.is_accepting_leads = TRUE
		ORDER BY a.current_call_duration_seconds ASC
	`
	
	rows, err := s.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		agent := models.AgentAvailability{}
		err := rows.Scan(
			&agent.ID, &agent.UserID, &agent.Status, &agent.BreakReason,
			&agent.IsAcceptingLeads, &agent.TotalCallsToday, &agent.CurrentCallDurationSeconds,
			&agent.LastStatusChange, &agent.LastActivity, &agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}
	
	return agents, rows.Err()
}

// UpdateCallDuration updates current call duration
func (s *agentService) UpdateCallDuration(ctx context.Context, userID string, durationSeconds int) error {
	query := `
		UPDATE agent_availability 
		SET current_call_duration_seconds = $1, last_activity = NOW()
		WHERE user_id = $2
	`
	
	_, err := s.db.ExecContext(ctx, query, durationSeconds, userID)
	return err
}

// IncrementDailyCallCount increments the call count for today
func (s *agentService) IncrementDailyCallCount(ctx context.Context, userID string) error {
	query := `
		UPDATE agent_availability 
		SET total_calls_today = total_calls_today + 1, last_activity = NOW()
		WHERE user_id = $1
	`
	
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}
```

**File: `internal/models/agent.go`**

```go
package models

import "time"

type AgentAvailability struct {
	ID                           string    `json:"id"`
	UserID                       string    `json:"user_id"`
	Status                       string    `json:"status"` // available, busy, on_break, offline, in_meeting, away
	BreakReason                  *string   `json:"break_reason,omitempty"`
	IsAcceptingLeads             bool      `json:"is_accepting_leads"`
	TotalCallsToday              int       `json:"total_calls_today"`
	CurrentCallDurationSeconds   int       `json:"current_call_duration_seconds"`
	LastStatusChange             time.Time `json:"last_status_change"`
	LastActivity                 time.Time `json:"last_activity"`
	CreatedAt                    time.Time `json:"created_at"`
	UpdatedAt                    time.Time `json:"updated_at"`
}

const (
	AgentStatusAvailable  = "available"
	AgentStatusBusy       = "busy"
	AgentStatusOnBreak    = "on_break"
	AgentStatusOffline    = "offline"
	AgentStatusInMeeting  = "in_meeting"
	AgentStatusAway       = "away"
)
```

### API Endpoints

**File: `internal/handlers/agent_handler.go`**

```go
package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"your-module/internal/services"
	"your-module/pkg/response"
)

type AgentHandler struct {
	agentService services.AgentService
}

func NewAgentHandler(router *mux.Router, agentService services.AgentService) {
	handler := &AgentHandler{agentService: agentService}
	
	// Routes
	router.HandleFunc("/api/v1/agents/availability/me", handler.GetMyAvailability).Methods("GET")
	router.HandleFunc("/api/v1/agents/availability", handler.UpdateAvailability).Methods("PUT")
	router.HandleFunc("/api/v1/agents/available", handler.ListAvailableAgents).Methods("GET")
	router.HandleFunc("/api/v1/agents/status", handler.UpdateStatus).Methods("PATCH")
}

// GetMyAvailability returns current agent's availability
func (h *AgentHandler) GetMyAvailability(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	
	availability, err := h.agentService.GetAgentAvailability(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, availability)
}

// UpdateAvailability updates agent availability
func (h *AgentHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	
	var req struct {
		Status      string `json:"status"`
		BreakReason string `json:"break_reason,omitempty"`
	}
	
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}
	
	err := h.agentService.UpdateAgentStatus(r.Context(), userID, req.Status, req.BreakReason)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, map[string]string{"message": "Status updated"})
}

// ListAvailableAgents returns available agents in team
func (h *AgentHandler) ListAvailableAgents(w http.ResponseWriter, r *http.Request) {
	teamID := r.URL.Query().Get("team_id")
	
	agents, err := h.agentService.GetAvailableAgents(r.Context(), teamID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, agents)
}

// UpdateStatus updates just the status
func (h *AgentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	
	var req struct {
		Status string `json:"status"`
	}
	
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request")
		return
	}
	
	err := h.agentService.UpdateAgentStatus(r.Context(), userID, req.Status, "")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	response.JSON(w, http.StatusOK, map[string]string{"message": "Status updated"})
}
```

### Frontend Hook

**File: `frontend/hooks/useAgentAvailability.ts`**

```typescript
import { useState, useCallback } from 'react';
import api from '@/services/api';

interface AgentAvailability {
  id: string;
  user_id: string;
  status: string;
  break_reason?: string;
  is_accepting_leads: boolean;
  total_calls_today: number;
  current_call_duration_seconds: number;
  last_status_change: string;
  last_activity: string;
}

export function useAgentAvailability() {
  const [availability, setAvailability] = useState<AgentAvailability | null>(null);
  const [availableAgents, setAvailableAgents] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getMyAvailability = useCallback(async () => {
    setLoading(true);
    try {
      const response = await api.agents.getMyAvailability();
      setAvailability(response);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const updateStatus = useCallback(async (status: string, breakReason?: string) => {
    setLoading(true);
    try {
      const response = await api.agents.updateStatus({ status, break_reason: breakReason });
      setAvailability(prev => prev ? { ...prev, status } : null);
      setError(null);
      return response;
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const getAvailableAgents = useCallback(async (teamId: string) => {
    setLoading(true);
    try {
      const response = await api.agents.getAvailableAgents(teamId);
      setAvailableAgents(response);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    availability,
    availableAgents,
    loading,
    error,
    getMyAvailability,
    updateStatus,
    getAvailableAgents,
  };
}
```

---

## Feature 2: Lead Scoring System

### Database Schema

```sql
CREATE TABLE IF NOT EXISTS lead_scores (
    id TEXT PRIMARY KEY DEFAULT generate_ulid(),
    lead_id TEXT UNIQUE NOT NULL REFERENCES leads(id) ON DELETE CASCADE,
    source_quality_score DECIMAL(3,2) DEFAULT 0.0,
    engagement_score DECIMAL(3,2) DEFAULT 0.0,
    conversion_probability DECIMAL(3,2) DEFAULT 0.0,
    urgency_score DECIMAL(3,2) DEFAULT 0.0,
    overall_score DECIMAL(5,2) DEFAULT 0.0,
    score_category VARCHAR(50), -- hot, warm, cold, nurture
    previous_score DECIMAL(5,2),
    score_change DECIMAL(5,2),
    reason_text TEXT,
    calculation_method VARCHAR(100) DEFAULT 'weighted',
    last_calculated TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_lead_scores_lead_id ON lead_scores(lead_id);
CREATE INDEX idx_lead_scores_overall_score ON lead_scores(overall_score DESC);
CREATE INDEX idx_lead_scores_score_category ON lead_scores(score_category);
```

### Backend Service

**File: `internal/services/lead_scoring_service.go`**

```go
package services

import (
	"context"
	"fmt"
	"math"
)

type LeadScoringService interface {
	CalculateLeadScore(ctx context.Context, leadID string) (float64, string, error)
	GetLeadScore(ctx context.Context, leadID string) (*LeadScore, error)
	UpdateLeadScore(ctx context.Context, leadID string) error
	BatchCalculateScores(ctx context.Context) error
}

type LeadScore struct {
	ID                  string  `json:"id"`
	LeadID              string  `json:"lead_id"`
	SourceQualityScore  float64 `json:"source_quality_score"`
	EngagementScore     float64 `json:"engagement_score"`
	ConversionProb      float64 `json:"conversion_probability"`
	UrgencyScore        float64 `json:"urgency_score"`
	OverallScore        float64 `json:"overall_score"`
	ScoreCategory       string  `json:"score_category"`
	ReasonText          string  `json:"reason_text,omitempty"`
	LastCalculated      string  `json:"last_calculated"`
}

type leadScoringService struct {
	db database.DB
}

// CalculateLeadScore calculates composite lead score
func (s *leadScoringService) CalculateLeadScore(ctx context.Context, leadID string) (float64, string, error) {
	// Get lead details
	var source, email, phone string
	var revisitDate, bookingDate *string
	
	query := `
		SELECT source, email, phone, revisit_date, booking_date
		FROM leads WHERE id = $1
	`
	err := s.db.QueryRowContext(ctx, query, leadID).Scan(
		&source, &email, &phone, &revisitDate, &bookingDate,
	)
	if err != nil {
		return 0, "", err
	}
	
	// Calculate individual scores
	sourceScore := s.calculateSourceScore(source)           // 0-25
	engagementScore := s.calculateEngagementScore(email, phone) // 0-25
	conversionProb := s.calculateConversionProbability(revisitDate, bookingDate) // 0-30
	urgencyScore := s.calculateUrgencyScore(revisitDate)     // 0-20
	
	// Weighted overall score (0-100)
	overall := (sourceScore * 0.25) + (engagementScore * 0.25) + 
	           (conversionProb * 0.30) + (urgencyScore * 0.20)
	
	// Determine category
	category := s.categorizeScore(overall)
	
	return math.Round(overall*100) / 100, category, nil
}

// calculateSourceScore: Different sources have different quality
func (s *leadScoringService) calculateSourceScore(source string) float64 {
	scores := map[string]float64{
		"direct_website":  25,
		"referral":        24,
		"google_ads":      22,
		"facebook_ads":    20,
		"instagram_ads":   18,
		"cold_call":       15,
		"other":           10,
	}
	if score, exists := scores[source]; exists {
		return score
	}
	return 10.0
}

// calculateEngagementScore: Contact info completeness
func (s *leadScoringService) calculateEngagementScore(email, phone string) float64 {
	score := 0.0
	if email != "" && email != "null" {
		score += 12.5
	}
	if phone != "" && phone != "null" {
		score += 12.5
	}
	return score
}

// calculateConversionProbability: Based on actions taken
func (s *leadScoringService) calculateConversionProbability(revisit, booking *string) float64 {
	score := 10.0
	if booking != nil {
		score = 30.0 // Booking scheduled = high probability
	} else if revisit != nil {
		score = 20.0 // Revisit scheduled = medium-high
	}
	return score
}

// calculateUrgencyScore: How soon action is needed
func (s *leadScoringService) calculateUrgencyScore(revisit *string) float64 {
	// If revisit is within 7 days, high urgency
	if revisit != nil {
		// Parse and calculate days to revisit
		// For now, return high urgency
		return 20.0
	}
	return 10.0
}

// categorizeScore: Convert numeric score to category
func (s *leadScoringService) categorizeScore(score float64) string {
	if score >= 75 {
		return "hot"
	} else if score >= 50 {
		return "warm"
	} else if score >= 25 {
		return "cold"
	}
	return "nurture"
}

// GetLeadScore retrieves existing lead score
func (s *leadScoringService) GetLeadScore(ctx context.Context, leadID string) (*LeadScore, error) {
	score := &LeadScore{}
	query := `
		SELECT id, lead_id, source_quality_score, engagement_score,
		       conversion_probability, urgency_score, overall_score,
		       score_category, reason_text, last_calculated
		FROM lead_scores WHERE lead_id = $1
	`
	
	err := s.db.QueryRowContext(ctx, query, leadID).Scan(
		&score.ID, &score.LeadID, &score.SourceQualityScore, &score.EngagementScore,
		&score.ConversionProb, &score.UrgencyScore, &score.OverallScore,
		&score.ScoreCategory, &score.ReasonText, &score.LastCalculated,
	)
	
	return score, err
}

// UpdateLeadScore recalculates and updates lead score
func (s *leadScoringService) UpdateLeadScore(ctx context.Context, leadID string) error {
	score, category, err := s.CalculateLeadScore(ctx, leadID)
	if err != nil {
		return err
	}
	
	query := `
		UPDATE lead_scores 
		SET overall_score = $1, score_category = $2, last_calculated = NOW()
		WHERE lead_id = $3
	`
	
	_, err = s.db.ExecContext(ctx, query, score, category, leadID)
	return err
}

// BatchCalculateScores recalculates all lead scores
func (s *leadScoringService) BatchCalculateScores(ctx context.Context) error {
	query := `
		SELECT id FROM leads WHERE created_at > NOW() - INTERVAL '30 days'
	`
	
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()
	
	for rows.Next() {
		var leadID string
		if err := rows.Scan(&leadID); err != nil {
			continue
		}
		_ = s.UpdateLeadScore(ctx, leadID)
	}
	
	return rows.Err()
}
```

### Frontend Hook

**File: `frontend/hooks/useLeadScoring.ts`**

```typescript
import { useState, useCallback } from 'react';
import api from '@/services/api';

interface LeadScore {
  id: string;
  lead_id: string;
  source_quality_score: number;
  engagement_score: number;
  conversion_probability: number;
  urgency_score: number;
  overall_score: number;
  score_category: string;
  reason_text?: string;
  last_calculated: string;
}

export function useLeadScoring() {
  const [score, setScore] = useState<LeadScore | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getLeadScore = useCallback(async (leadId: string) => {
    setLoading(true);
    try {
      const response = await api.leads.getScore(leadId);
      setScore(response);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const recalculateScore = useCallback(async (leadId: string) => {
    setLoading(true);
    try {
      const response = await api.leads.recalculateScore(leadId);
      setScore(response);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    score,
    loading,
    error,
    getLeadScore,
    recalculateScore,
  };
}
```

---

## Feature 3: Audit Trail System

### Database Schema

```sql
CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY DEFAULT generate_ulid(),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id TEXT NOT NULL,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    status VARCHAR(50) DEFAULT 'success',
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
```

### Backend Service

**File: `internal/services/audit_service.go`**

```go
package services

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
)

type AuditService interface {
	LogAction(ctx context.Context, log AuditLog) error
	GetAuditLogs(ctx context.Context, entityType, entityID string, limit int) ([]AuditLog, error)
	GetUserAuditLogs(ctx context.Context, userID string, limit int) ([]AuditLog, error)
}

type AuditLog struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	Action       string          `json:"action"`
	EntityType   string          `json:"entity_type"`
	EntityID     string          `json:"entity_id"`
	OldValues    json.RawMessage `json:"old_values,omitempty"`
	NewValues    json.RawMessage `json:"new_values,omitempty"`
	IPAddress    string          `json:"ip_address,omitempty"`
	UserAgent    string          `json:"user_agent,omitempty"`
	Status       string          `json:"status"`
	ErrorMessage string          `json:"error_message,omitempty"`
	CreatedAt    string          `json:"created_at"`
}

type auditService struct {
	db database.DB
}

func NewAuditService(db database.DB) AuditService {
	return &auditService{db: db}
}

// LogAction creates an audit log entry
func (s *auditService) LogAction(ctx context.Context, log AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			user_id, action, entity_type, entity_id, 
			old_values, new_values, ip_address, user_agent, status, error_message
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := s.db.ExecContext(ctx, query,
		log.UserID, log.Action, log.EntityType, log.EntityID,
		log.OldValues, log.NewValues, log.IPAddress, log.UserAgent,
		log.Status, log.ErrorMessage,
	)
	return err
}

// GetAuditLogs retrieves audit logs for an entity
func (s *auditService) GetAuditLogs(ctx context.Context, entityType, entityID string, limit int) ([]AuditLog, error) {
	logs := []AuditLog{}
	query := `
		SELECT id, user_id, action, entity_type, entity_id,
		       old_values, new_values, ip_address, user_agent,
		       status, error_message, created_at
		FROM audit_logs
		WHERE entity_type = $1 AND entity_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`
	
	rows, err := s.db.QueryContext(ctx, query, entityType, entityID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		log := AuditLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.EntityType, &log.EntityID,
			&log.OldValues, &log.NewValues, &log.IPAddress, &log.UserAgent,
			&log.Status, &log.ErrorMessage, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, rows.Err()
}

// GetUserAuditLogs retrieves all audit logs for a user
func (s *auditService) GetUserAuditLogs(ctx context.Context, userID string, limit int) ([]AuditLog, error) {
	logs := []AuditLog{}
	query := `
		SELECT id, user_id, action, entity_type, entity_id,
		       old_values, new_values, ip_address, user_agent,
		       status, error_message, created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := s.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		log := AuditLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.EntityType, &log.EntityID,
			&log.OldValues, &log.NewValues, &log.IPAddress, &log.UserAgent,
			&log.Status, &log.ErrorMessage, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	
	return logs, rows.Err()
}

// Helper: Extract IP from request
func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}
```

### Middleware for Automatic Logging

**File: `internal/middleware/audit_middleware.go`**

```go
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"your-module/internal/services"
)

func AuditMiddleware(auditService services.AuditService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture request body
			var requestBody []byte
			if r.Body != nil {
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}
			
			userID := r.Context().Value("user_id")
			if userID == nil {
				next.ServeHTTP(w, r)
				return
			}
			
			// Log the action
			auditLog := services.AuditLog{
				UserID:    userID.(string),
				Action:    r.Method + " " + r.RequestURI,
				IPAddress: services.GetClientIP(r),
				UserAgent: r.Header.Get("User-Agent"),
				Status:    "success",
			}
			
			// Extract entity from body if POST/PUT
			if r.Method == "POST" || r.Method == "PUT" {
				var bodyMap map[string]interface{}
				json.Unmarshal(requestBody, &bodyMap)
				auditLog.NewValues, _ = json.Marshal(bodyMap)
			}
			
			_ = auditService.LogAction(r.Context(), auditLog)
			
			next.ServeHTTP(w, r)
		})
	}
}
```

---

## Integration with Existing Features

### 1. Agent Availability with Automation (Feature 3)
```go
// In automation service, use agent availability to route leads
availableAgents, _ := agentService.GetAvailableAgents(ctx, teamID)
// Choose agent with lowest current call duration
selectedAgent := availableAgents[0]
```

### 2. Lead Scoring with Automation (Feature 3)
```go
// In automation service, score leads before assignment
score, category, _ := scoringService.CalculateLeadScore(ctx, leadID)
if category == "hot" {
    // Assign to senior agent
}
```

### 3. Audit Trail with Compliance (Feature 6)
```go
// All lead assignments automatically logged
auditLog := services.AuditLog{
    UserID: userID,
    Action: "lead_assignment",
    EntityType: "lead",
    EntityID: leadID,
}
auditService.LogAction(ctx, auditLog)
```

---

## Testing Strategy

### Unit Tests
```go
func TestLeadScoringService(t *testing.T) {
	service := services.NewLeadScoringService(mockDB)
	
	// Test source scoring
	score := service.calculateSourceScore("direct_website")
	assert.Equal(t, 25.0, score)
	
	// Test category assignment
	category := service.categorizeScore(85.5)
	assert.Equal(t, "hot", category)
}
```

### Integration Tests
- Test agent availability updates
- Test lead score recalculation
- Test audit logging for all actions

---

## Deployment Checklist

- [ ] Create database migration for Phase 1 tables
- [ ] Deploy backend services
- [ ] Deploy API endpoints
- [ ] Create frontend hooks
- [ ] Update frontend UI components
- [ ] Run integration tests
- [ ] Deploy to production
- [ ] Monitor audit logs

---

## Performance Considerations

1. **Agent Availability Caching:** Cache available agents list, refresh every 30 seconds
2. **Lead Scoring:** Run batch calculations during off-peak hours (3 AM)
3. **Audit Logs:** Partition by date for efficient querying
4. **Indexes:** Ensure all WHERE clauses have indexes

---

## Next Steps

1. Implement Phase 1 features in development
2. Create comprehensive tests
3. Set up monitoring and alerting
4. Document API endpoints
5. Train team on new features
6. Deploy to staging
7. Deploy to production

This Phase 1 provides immediate value while preparing foundation for Phase 2-4 improvements.
