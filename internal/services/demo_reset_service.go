package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"vyomtech-backend/pkg/logger"
)

// DemoResetService handles automatic reset of demo data every 30 days
type DemoResetService struct {
	db     *sql.DB
	logger *logger.Logger
}

const (
	DemoTenantID  = "demo_vyomtech_001"
	ResetInterval = 30 * 24 * time.Hour // 30 days
)

// NewDemoResetService creates a new demo reset service
func NewDemoResetService(db *sql.DB, log *logger.Logger) *DemoResetService {
	return &DemoResetService{
		db:     db,
		logger: log,
	}
}

// StartScheduler starts the 30-day demo data reset scheduler
func (s *DemoResetService) StartScheduler() {
	go func() {
		ticker := time.NewTicker(ResetInterval)
		defer ticker.Stop()

		// First reset on startup if needed
		if err := s.ResetDemoData(); err != nil {
			s.logger.Error("[DemoReset] Initial reset failed", "error", err)
		}

		// Then reset every 30 days
		for range ticker.C {
			if err := s.ResetDemoData(); err != nil {
				s.logger.Error("[DemoReset] Scheduled reset failed", "error", err)
			}
		}
	}()

	s.logger.Info("[DemoReset] Scheduler started", "interval", ResetInterval.String())
}

// ResetDemoData clears and reloads all demo tenant data
func (s *DemoResetService) ResetDemoData() error {
	s.logger.Info("[DemoReset] Starting reset of demo data", "tenant_id", DemoTenantID)

	ctx := context.Background()

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("[DemoReset] Error starting transaction", "error", err)
		return err
	}

	// Clear all demo data
	if err := s.clearDemoData(ctx, tx); err != nil {
		tx.Rollback()
		s.logger.Error("[DemoReset] Error clearing demo data", "error", err)
		return err
	}

	// Reload fresh demo data
	if err := s.reloadDemoData(ctx, tx); err != nil {
		tx.Rollback()
		s.logger.Error("[DemoReset] Error reloading demo data", "error", err)
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("[DemoReset] Error committing transaction", "error", err)
		return err
	}

	s.logger.Info("[DemoReset] ✓ Demo data reset completed successfully", "timestamp", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

// tableExists checks if a table exists in the database
func (s *DemoResetService) tableExists(ctx context.Context, tableName string) bool {
	query := `SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`
	var name string
	err := s.db.QueryRowContext(ctx, query, tableName).Scan(&name)
	return err == nil
}

// clearDemoData removes all demo tenant data using SQL
func (s *DemoResetService) clearDemoData(ctx context.Context, tx *sql.Tx) error {
	// Delete in reverse dependency order
	tables := []string{
		"partner_lead_credits",
		"partner_leads",
		"partner_payouts",
		"partner_payout_details",
		"partner_activities",
		"partner_users",
		"partners",
		"gamification_points_history",
		"gamification_badges",
		"progress_tracking",
		"compliance_records",
		"task",
		"call",
		"campaign_recipient",
		"campaign",
		"lead",
		"agent",
		"construction_projects",
	}

	for _, table := range tables {
		// Check if table exists before trying to delete
		if !s.tableExists(ctx, table) {
			s.logger.Debug("[DemoReset] Table does not exist, skipping", "table", table)
			continue
		}

		// Use backticks for reserved keywords like 'call' and 'lead'
		query := fmt.Sprintf("DELETE FROM `%s` WHERE tenant_id = ?", table)
		result, err := tx.ExecContext(ctx, query, DemoTenantID)
		if err != nil {
			// Log but continue - might be other issues
			s.logger.Warn("[DemoReset] Failed to delete from table", "table", table, "error", err)
		} else {
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected > 0 {
				s.logger.Debug("[DemoReset] Deleted rows", "table", table, "count", rowsAffected)
			}
		}
	}

	s.logger.Info("[DemoReset] Cleared demo data", "tenant_id", DemoTenantID)
	return nil
}

// reloadDemoData inserts fresh demo data using SQL
func (s *DemoResetService) reloadDemoData(ctx context.Context, tx *sql.Tx) error {
	now := time.Now()

	// Insert demo agents (4 agents) - only if agent table exists
	if s.tableExists(ctx, "agent") {
		agentQueries := []struct {
			name  string
			email string
			phone string
			score float64
		}{
			{"Rajesh Kumar", "rajesh@demo.vyomtech.com", "+91-9876543210", 95.0},
			{"Priya Singh", "priya@demo.vyomtech.com", "+91-8765432109", 92.0},
			{"Arun Patel", "arun@demo.vyomtech.com", "+91-7654321098", 88.0},
			{"Neha Sharma", "neha@demo.vyomtech.com", "+91-6543210987", 90.0},
		}

		for i, agent := range agentQueries {
			agentCode := fmt.Sprintf("AGENT%03d", i+1)
			skillsJSON := `["Customer Support","Sales"]`
			_, err := tx.ExecContext(ctx,
				`INSERT IGNORE INTO agent (id, tenant_id, agent_code, first_name, last_name, email, phone, status, agent_type, skills, available, created_at, updated_at) 
				 VALUES (UUID(), ?, ?, ?, ?, ?, ?, 'available', 'inbound', ?, TRUE, ?, ?)`,
				DemoTenantID, agentCode, agent.name, "", agent.email, agent.phone, skillsJSON, now, now)
			if err != nil {
				s.logger.Warn("[DemoReset] Failed to insert agent", "email", agent.email, "error", err)
			}
		}
		s.logger.Info("[DemoReset] Inserted demo agents", "count", len(agentQueries))
	}

	// Insert demo leads (5 leads) - only if sales_lead table exists
	if s.tableExists(ctx, "sales_lead") {
		leadQueries := []struct {
			title    string
			value    float64
			contact  string
			email    string
			phone    string
			propType string
			location string
		}{
			{"High Value Residential Project", 5000000, "Amit Kumar", "amit@example.com", "+91-9876543211", "residential", "Mumbai, Maharashtra"},
			{"Commercial Space Inquiry", 3500000, "Sneha Desai", "sneha@example.com", "+91-9876543212", "commercial", "Bangalore, Karnataka"},
			{"Plot Purchase Interest", 2000000, "Vikram Singh", "vikram@example.com", "+91-9876543213", "plot", "Delhi, Delhi"},
			{"Rental Inquiry", 1500000, "Meera Nair", "meera@example.com", "+91-9876543214", "rental", "Hyderabad, Telangana"},
			{"Apartment Pre-booking", 4000000, "Rohan Gupta", "rohan@example.com", "+91-9876543215", "apartment", "Pune, Maharashtra"},
		}

		for i, lead := range leadQueries {
			leadCode := fmt.Sprintf("LEAD%03d", i+1)
			_, err := tx.ExecContext(ctx,
				`INSERT IGNORE INTO sales_lead (id, tenant_id, lead_code, first_name, last_name, email, phone, company_name, status, source, created_by, created_at, updated_at) 
				 VALUES (UUID(), ?, ?, ?, ?, ?, ?, ?, 'new', 'Direct', 1, ?, ?)`,
				DemoTenantID, leadCode, lead.contact, "", lead.email, lead.phone, lead.propType, now, now)
			if err != nil {
				s.logger.Warn("[DemoReset] Failed to insert lead", "title", lead.title, "error", err)
			}
		}
		s.logger.Info("[DemoReset] Inserted demo leads", "count", len(leadQueries))
	}

	// Insert demo campaigns (4 campaigns) - only if campaign table exists
	if s.tableExists(ctx, "campaign") {
		campaignQueries := []struct {
			name     string
			cType    string
			budget   float64
			audience string
		}{
			{"Summer Residential Drive 2025", "email", 500000, "Homebuyers"},
			{"Commercial Real Estate Expo", "event", 1000000, "Business Owners"},
			{"Digital Marketing Campaign", "social_media", 300000, "Young Professionals"},
			{"Corporate Bulk Purchase", "direct_sales", 800000, "Corporate Houses"},
		}

		for _, campaign := range campaignQueries {
			_, err := tx.ExecContext(ctx,
				`INSERT IGNORE INTO campaign (id, tenant_id, campaign_name, campaign_type, description, status, start_date, target_leads, budget, assigned_agents, created_by, created_at, updated_at) 
				 VALUES (UUID(), ?, ?, ?, ?, 'planning', DATE_ADD(NOW(), INTERVAL 1 MONTH), 500, ?, 2, 1, ?, ?)`,
				DemoTenantID, campaign.name, campaign.cType, campaign.audience, campaign.budget, now, now)
			if err != nil {
				s.logger.Warn("[DemoReset] Failed to insert campaign", "name", campaign.name, "error", err)
			}
		}
		s.logger.Info("[DemoReset] Inserted demo campaigns", "count", len(campaignQueries))
	}

	// Insert demo construction projects (4 projects) - only if construction_projects table exists
	if s.tableExists(ctx, "construction_projects") {
		projectQueries := []struct {
			name      string
			pType     string
			location  string
			cost      float64
			startDate string
			endDate   string
		}{
			{"Skyrise Towers Mumbai", "residential", "Mumbai, Maharashtra", 50000000, "2024-01-15", "2026-12-31"},
			{"Tech Park Bangalore", "commercial", "Bangalore, Karnataka", 100000000, "2023-06-01", "2025-12-31"},
			{"Plot Development Delhi", "plot_development", "Delhi, Delhi", 20000000, "2025-03-01", "2026-03-31"},
			{"Green Spaces Pune", "mixed_use", "Pune, Maharashtra", 75000000, "2024-08-01", "2027-08-31"},
		}

		for i, project := range projectQueries {
			projCode := fmt.Sprintf("PROJ%03d", i+1)
			_, err := tx.ExecContext(ctx,
				`INSERT IGNORE INTO construction_projects (tenant_id, project_name, project_code, location, contract_value, start_date, expected_completion, status, created_at, updated_at) 
				 VALUES (?, ?, ?, ?, ?, ?, ?, 'planning', ?, ?)`,
				DemoTenantID, project.name, projCode, project.location, project.cost, project.startDate, project.endDate, now, now)
			if err != nil {
				s.logger.Warn("[DemoReset] Failed to insert project", "name", project.name, "error", err)
			}
		}
		s.logger.Info("[DemoReset] Inserted demo projects", "count", len(projectQueries))
	}

	return nil
}
