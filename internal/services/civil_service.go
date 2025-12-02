package services

import (
	"database/sql"
	"fmt"
)

// CivilService provides civil engineering management functionality
type CivilService struct {
	DB *sql.DB
}

// NewCivilService creates a new civil service instance
func NewCivilService(db *sql.DB) *CivilService {
	return &CivilService{
		DB: db,
	}
}

// CivilDashboardMetrics represents civil dashboard metrics
type CivilDashboardMetrics struct {
	TotalSites          int     `json:"total_sites"`
	ActiveSites         int     `json:"active_sites"`
	TotalIncidents      int     `json:"total_incidents"`
	CriticalIncidents   int     `json:"critical_incidents"`
	ComplianceScore     float64 `json:"compliance_score"`
	PermitsExpiringSoon int     `json:"permits_expiring_soon"`
}

// GetDashboardMetrics retrieves civil dashboard metrics for a tenant
func (s *CivilService) GetDashboardMetrics(tenantID string) (*CivilDashboardMetrics, error) {
	metrics := &CivilDashboardMetrics{}

	// Get total sites
	err := s.DB.QueryRow(
		"SELECT COUNT(*) FROM sites WHERE tenant_id = ? AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.TotalSites)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count total sites: %w", err)
	}

	// Get active sites
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM sites WHERE tenant_id = ? AND current_status = 'active' AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.ActiveSites)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count active sites: %w", err)
	}

	// Get total incidents
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM safety_incidents WHERE tenant_id = ? AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.TotalIncidents)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count incidents: %w", err)
	}

	// Get critical incidents
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM safety_incidents WHERE tenant_id = ? AND severity = 'critical' AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.CriticalIncidents)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count critical incidents: %w", err)
	}

	// Get compliance score (percentage of compliant records)
	var complianceRecords, compliantRecords int
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM compliance_records WHERE tenant_id = ? AND deleted_at IS NULL",
		tenantID,
	).Scan(&complianceRecords)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count compliance records: %w", err)
	}

	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM compliance_records WHERE tenant_id = ? AND status = 'compliant' AND deleted_at IS NULL",
		tenantID,
	).Scan(&compliantRecords)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count compliant records: %w", err)
	}

	if complianceRecords > 0 {
		metrics.ComplianceScore = (float64(compliantRecords) / float64(complianceRecords)) * 100
	} else {
		metrics.ComplianceScore = 100.0
	}

	// Get permits expiring soon (within 30 days)
	err = s.DB.QueryRow(
		"SELECT COUNT(*) FROM permits WHERE tenant_id = ? AND status = 'active' AND expiry_date BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 30 DAY) AND deleted_at IS NULL",
		tenantID,
	).Scan(&metrics.PermitsExpiringSoon)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to count expiring permits: %w", err)
	}

	return metrics, nil
}
