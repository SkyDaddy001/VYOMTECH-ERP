package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

type AgentService struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAgentService(db *sql.DB, logger *logger.Logger) *AgentService {
	return &AgentService{
		db:     db,
		logger: logger,
	}
}

func (s *AgentService) GetAgent(ctx context.Context, agentID int) (*models.Agent, error) {
	var agent models.Agent
	var skillsJSON string

	err := s.db.QueryRowContext(ctx, `
		SELECT user_id, status, availability, skills, max_concurrent_calls, 
		       current_calls, total_calls, avg_handle_time, satisfaction_score, 
		       last_active, tenant_id
		FROM agent WHERE user_id = ?
	`, agentID).Scan(
		&agent.UserID, &agent.Status, &agent.Availability, &skillsJSON,
		&agent.MaxConcurrentCalls, &agent.CurrentCalls, &agent.TotalCalls,
		&agent.AvgHandleTime, &agent.SatisfactionScore, &agent.LastActive, &agent.TenantID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Parse skills JSON
	if skillsJSON != "" {
		err = json.Unmarshal([]byte(skillsJSON), &agent.Skills)
		if err != nil {
			s.logger.Warn("Failed to unmarshal agent skills", "error", err)
		}
	}

	return &agent, nil
}

func (s *AgentService) GetAgentsByTenant(ctx context.Context, tenantID string) ([]models.Agent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, status, availability, skills, max_concurrent_calls,
		       current_calls, total_calls, avg_handle_time, satisfaction_score,
		       last_active, tenant_id
		FROM agent WHERE tenant_id = ?
		ORDER BY last_active DESC
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var skillsJSON string

		err = rows.Scan(
			&agent.UserID, &agent.Status, &agent.Availability, &skillsJSON,
			&agent.MaxConcurrentCalls, &agent.CurrentCalls, &agent.TotalCalls,
			&agent.AvgHandleTime, &agent.SatisfactionScore, &agent.LastActive, &agent.TenantID,
		)
		if err != nil {
			s.logger.Error("Failed to scan agent row", "error", err)
			continue
		}

		if skillsJSON != "" {
			err = json.Unmarshal([]byte(skillsJSON), &agent.Skills)
			if err != nil {
				s.logger.Warn("Failed to unmarshal agent skills", "error", err)
			}
		}

		agents = append(agents, agent)
	}

	return agents, rows.Err()
}

func (s *AgentService) UpdateAgentAvailability(ctx context.Context, agentID int, availability string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET availability = ?, last_active = NOW() 
		WHERE user_id = ?
	`, availability, agentID)

	if err != nil {
		return fmt.Errorf("failed to update agent availability: %w", err)
	}

	s.logger.WithUser(agentID).Info("Agent availability updated", "availability", availability)
	return nil
}

func (s *AgentService) UpdateAgentStatus(ctx context.Context, agentID int, status string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET status = ?, updated_at = NOW() 
		WHERE user_id = ?
	`, status, agentID)

	if err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}

	s.logger.WithUser(agentID).Info("Agent status updated", "status", status)
	return nil
}

func (s *AgentService) GetAvailableAgents(ctx context.Context, tenantID string) ([]models.Agent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, status, availability, skills, max_concurrent_calls,
		       current_calls, total_calls, avg_handle_time, satisfaction_score,
		       last_active, tenant_id
		FROM agent 
		WHERE tenant_id = ? AND status = 'active' AND availability = 'online'
		      AND current_calls < max_concurrent_calls
		ORDER BY satisfaction_score DESC
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var skillsJSON string

		err = rows.Scan(
			&agent.UserID, &agent.Status, &agent.Availability, &skillsJSON,
			&agent.MaxConcurrentCalls, &agent.CurrentCalls, &agent.TotalCalls,
			&agent.AvgHandleTime, &agent.SatisfactionScore, &agent.LastActive, &agent.TenantID,
		)
		if err != nil {
			s.logger.Error("Failed to scan agent row", "error", err)
			continue
		}

		if skillsJSON != "" {
			err = json.Unmarshal([]byte(skillsJSON), &agent.Skills)
			if err != nil {
				s.logger.Warn("Failed to unmarshal agent skills", "error", err)
			}
		}

		agents = append(agents, agent)
	}

	return agents, rows.Err()
}

func (s *AgentService) IncrementAgentCallCount(ctx context.Context, agentID int) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET current_calls = current_calls + 1, total_calls = total_calls + 1
		WHERE user_id = ?
	`, agentID)

	if err != nil {
		return fmt.Errorf("failed to increment agent call count: %w", err)
	}

	return nil
}

func (s *AgentService) DecrementAgentCallCount(ctx context.Context, agentID int) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET current_calls = GREATEST(0, current_calls - 1)
		WHERE user_id = ?
	`, agentID)

	if err != nil {
		return fmt.Errorf("failed to decrement agent call count: %w", err)
	}

	return nil
}

func (s *AgentService) GetAgentStats(ctx context.Context, tenantID string) (*models.AgentStats, error) {
	var stats models.AgentStats

	// Get counts for different availability statuses
	err := s.db.QueryRowContext(ctx, `
		SELECT 
			SUM(CASE WHEN availability = 'online' THEN 1 ELSE 0 END) as online_agents,
			SUM(CASE WHEN availability = 'busy' THEN 1 ELSE 0 END) as busy_agents,
			COUNT(*) as total_agents
		FROM agent WHERE tenant_id = ? AND status = 'active'
	`, tenantID).Scan(&stats.OnlineAgents, &stats.BusyAgents, &stats.TotalAgents)

	if err != nil {
		return nil, fmt.Errorf("failed to get agent stats: %w", err)
	}

	return &stats, nil
}

func (s *AgentService) UpdateAgentStats(ctx context.Context, agentID int, avgHandleTime float64, satisfactionScore float64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET avg_handle_time = ?, satisfaction_score = ?, last_active = NOW()
		WHERE user_id = ?
	`, avgHandleTime, satisfactionScore, agentID)

	if err != nil {
		return fmt.Errorf("failed to update agent stats: %w", err)
	}

	return nil
}
