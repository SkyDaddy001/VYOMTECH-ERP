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

func (s *AgentService) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	var agent models.Agent
	var skillsJSON sql.NullString

	err := s.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, agent_code, first_name, last_name, email, phone, status, agent_type, skills, available, created_at, updated_at
		FROM agent WHERE id = ?
	`, agentID).Scan(
		&agent.ID, &agent.TenantID, &agent.AgentCode, &agent.FirstName, &agent.LastName,
		&agent.Email, &agent.Phone, &agent.Status, &agent.AgentType, &skillsJSON,
		&agent.Available, &agent.CreatedAt, &agent.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("agent not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Parse skills JSON
	if skillsJSON.Valid && skillsJSON.String != "" {
		err = json.Unmarshal([]byte(skillsJSON.String), &agent.Skills)
		if err != nil {
			s.logger.Warn("Failed to unmarshal agent skills", "error", err)
		}
	}

	return &agent, nil
}

func (s *AgentService) GetAgentsByTenant(ctx context.Context, tenantID string) ([]models.Agent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, agent_code, first_name, last_name, email, phone, status, agent_type, skills, available, created_at, updated_at
		FROM agent WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var skillsJSON sql.NullString
		var id, createdAt, updatedAt sql.NullString

		err = rows.Scan(
			&id, &agent.TenantID, &agent.AgentCode, &agent.FirstName, &agent.LastName,
			&agent.Email, &agent.Phone, &agent.Status, &agent.AgentType, &skillsJSON,
			&agent.Available, &createdAt, &updatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan agent row", "error", err)
			continue
		}

		if id.Valid {
			agent.ID = id.String
		}

		if skillsJSON.Valid && skillsJSON.String != "" {
			err = json.Unmarshal([]byte(skillsJSON.String), &agent.Skills)
			if err != nil {
				s.logger.Warn("Failed to unmarshal agent skills", "error", err)
			}
		}

		agents = append(agents, agent)
	}

	return agents, rows.Err()
}

func (s *AgentService) UpdateAgentAvailability(ctx context.Context, agentID string, availability string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET available = CASE WHEN ? = 'available' THEN TRUE ELSE FALSE END, updated_at = NOW() 
		WHERE id = ?
	`, availability, agentID)

	if err != nil {
		return fmt.Errorf("failed to update agent availability: %w", err)
	}

	s.logger.Info("Agent availability updated", "availability", availability)
	return nil
}

func (s *AgentService) UpdateAgentStatus(ctx context.Context, agentID string, status string) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent 
		SET status = ?, updated_at = NOW() 
		WHERE id = ?
	`, status, agentID)

	if err != nil {
		return fmt.Errorf("failed to update agent status: %w", err)
	}

	s.logger.Info("Agent status updated", "status", status)
	return nil
}

func (s *AgentService) GetAvailableAgents(ctx context.Context, tenantID string) ([]models.Agent, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tenant_id, agent_code, first_name, last_name, email, phone, status, agent_type, skills, available, created_at, updated_at
		FROM agent 
		WHERE tenant_id = ? AND status = 'available' AND available = TRUE
		ORDER BY created_at DESC
	`, tenantID)

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var skillsJSON sql.NullString

		err = rows.Scan(
			&agent.ID, &agent.TenantID, &agent.AgentCode, &agent.FirstName, &agent.LastName,
			&agent.Email, &agent.Phone, &agent.Status, &agent.AgentType, &skillsJSON,
			&agent.Available, &agent.CreatedAt, &agent.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("Failed to scan agent row", "error", err)
			continue
		}

		if skillsJSON.Valid && skillsJSON.String != "" {
			err = json.Unmarshal([]byte(skillsJSON.String), &agent.Skills)
			if err != nil {
				s.logger.Warn("Failed to unmarshal agent skills", "error", err)
			}
		}

		agents = append(agents, agent)
	}

	return agents, rows.Err()
}

func (s *AgentService) IncrementAgentCallCount(ctx context.Context, agentID string) error {
	// This feature is not supported in the current schema
	// Keeping for backward compatibility
	return nil
}

func (s *AgentService) DecrementAgentCallCount(ctx context.Context, agentID string) error {
	// This feature is not supported in the current schema
	// Keeping for backward compatibility
	return nil
}

func (s *AgentService) GetAgentStats(ctx context.Context, tenantID string) (*models.AgentStats, error) {
	var stats models.AgentStats

	// Get counts for different availability statuses
	err := s.db.QueryRowContext(ctx, `
		SELECT 
			SUM(CASE WHEN status = 'available' THEN 1 ELSE 0 END) as online_agents,
			SUM(CASE WHEN status = 'busy' THEN 1 ELSE 0 END) as busy_agents,
			COUNT(*) as total_agents
		FROM agent WHERE tenant_id = ?
	`, tenantID).Scan(&stats.OnlineAgents, &stats.BusyAgents, &stats.TotalAgents)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get agent stats: %w", err)
	}

	return &stats, nil
}

func (s *AgentService) UpdateAgentStats(ctx context.Context, agentID string, avgHandleTime float64, satisfactionScore float64) error {
	// This feature is not supported in the current schema
	// Keeping for backward compatibility
	return nil
}
