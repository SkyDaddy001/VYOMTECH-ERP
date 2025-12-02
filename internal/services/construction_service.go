package services

import (
	"database/sql"
	"fmt"
	"time"
)

// ConstructionService provides construction management functionality
type ConstructionService struct {
	DB *sql.DB
}

// NewConstructionService creates a new construction service instance
func NewConstructionService(db *sql.DB) *ConstructionService {
	return &ConstructionService{
		DB: db,
	}
}

// ConstructionDashboardMetrics represents construction dashboard metrics
type ConstructionDashboardMetrics struct {
	TotalProjects      int     `json:"total_projects"`
	ActiveProjects     int     `json:"active_projects"`
	CompletedProjects  int     `json:"completed_projects"`
	AverageProgress    float64 `json:"average_progress"`
	OnScheduleProjects int     `json:"on_schedule_projects"`
	DelayedProjects    int     `json:"delayed_projects"`
	TotalContractValue float64 `json:"total_contract_value"`
}

// GetDashboardMetrics retrieves construction dashboard metrics for a tenant
func (s *ConstructionService) GetDashboardMetrics(tenantID string) (*ConstructionDashboardMetrics, error) {
	metrics := &ConstructionDashboardMetrics{}

	// Get total projects
	err := s.DB.QueryRow(
		"SELECT COUNT(*) FROM construction_projects WHERE tenant_id = ? AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.TotalProjects)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count total projects: %w", err)
	}

	// Get active projects
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM construction_projects WHERE tenant_id = ? AND status = 'active' AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.ActiveProjects)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count active projects: %w", err)
	}

	// Get completed projects
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM construction_projects WHERE tenant_id = ? AND status = 'completed' AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.CompletedProjects)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count completed projects: %w", err)
	}

	// Get average progress percentage
	err = s.DB.QueryRow(
		"SELECT COALESCE(AVG(current_progress_percent), 0) FROM construction_projects WHERE tenant_id = ? AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.AverageProgress)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to calculate average progress: %w", err)
	}

	// Get total contract value
	err = s.DB.QueryRow(
		"SELECT COALESCE(SUM(contract_value), 0) FROM construction_projects WHERE tenant_id = ? AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.TotalContractValue)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to sum contract value: %w", err)
	}

	// Calculate on-schedule projects (expected_completion >= now or current_progress_percent >= (days_elapsed/total_days)*100)
	today := time.Now()
	rows, err := s.DB.Query(
		`SELECT COUNT(*) FROM construction_projects 
		 WHERE tenant_id = ? AND status IN ('active', 'planning') AND deleted_at IS NULL
		 AND expected_completion >= ? 
		 AND (current_progress_percent >= (DATEDIFF(?, start_date) / DATEDIFF(expected_completion, start_date) * 100) OR current_progress_percent >= 0)`,
		tenantID, today, today,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate on-schedule projects: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&metrics.OnScheduleProjects); err != nil {
			return nil, fmt.Errorf("failed to scan on-schedule projects: %w", err)
		}
	}

	// Calculate delayed projects (expected_completion < now and status != completed)
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM construction_projects WHERE tenant_id = ? AND expected_completion < ? AND status != 'completed' AND deleted_at IS NULL",
		tenantID, today,
	).Scan(&metrics.DelayedProjects)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count delayed projects: %w", err)
	}

	return metrics, nil
}
