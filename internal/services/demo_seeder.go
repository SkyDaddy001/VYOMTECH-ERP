package services

import (
	"context"
	"database/sql"
	"fmt"
	"vyomtech-backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// DemoUser represents a demo user to be seeded
type DemoUser struct {
	Email    string
	Password string
	Role     string
	TenantID string
}

// DemoUserSeeder handles seeding demo users into the database
type DemoUserSeeder struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewDemoUserSeeder creates a new demo user seeder
func NewDemoUserSeeder(db *sql.DB, logger *logger.Logger) *DemoUserSeeder {
	return &DemoUserSeeder{
		db:     db,
		logger: logger,
	}
}

// SeedDemoUsers creates demo users for testing
func (s *DemoUserSeeder) SeedDemoUsers(ctx context.Context) error {
	const demoTenantID = "demo-tenant"
	const demoTenantName = "Demo Tenant"

	// First, ensure demo tenant exists
	if err := s.ensureTenant(ctx, demoTenantID, demoTenantName); err != nil {
		return fmt.Errorf("failed to ensure demo tenant: %w", err)
	}

	// Demo user credentials
	demoUsers := []DemoUser{
		{
			Email:    "demo@vyomtech.com",
			Password: "DemoPass@123",
			Role:     "admin",
			TenantID: demoTenantID,
		},
		{
			Email:    "agent@vyomtech.com",
			Password: "AgentPass@123",
			Role:     "agent",
			TenantID: demoTenantID,
		},
		{
			Email:    "manager@vyomtech.com",
			Password: "ManagerPass@123",
			Role:     "manager",
			TenantID: demoTenantID,
		},
		{
			Email:    "sales@vyomtech.com",
			Password: "SalesPass@123",
			Role:     "sales",
			TenantID: demoTenantID,
		},
		{
			Email:    "hr@vyomtech.com",
			Password: "HRPass@123",
			Role:     "hr",
			TenantID: demoTenantID,
		},
	}

	// Seed each demo user
	for _, user := range demoUsers {
		if err := s.seedUser(ctx, user); err != nil {
			s.logger.Warn("Failed to seed user", "email", user.Email, "error", err)
			// Continue seeding other users even if one fails
		}
	}

	s.logger.Info("Demo users seeded successfully")
	return nil
}

// ensureTenant creates the demo tenant if it doesn't exist
func (s *DemoUserSeeder) ensureTenant(ctx context.Context, tenantID, tenantName string) error {
	// Check if tenant already exists
	var existingID string
	err := s.db.QueryRowContext(ctx, "SELECT id FROM tenant WHERE id = ?", tenantID).Scan(&existingID)
	if err == nil {
		// Tenant already exists
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("database error: %w", err)
	}

	// Tenant doesn't exist, create it
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO tenant (id, name, status, max_users, max_concurrent_calls, ai_budget_monthly) 
		 VALUES (?, ?, 'active', 100, 50, 1000.00)`,
		tenantID, tenantName)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	s.logger.Info("Demo tenant created", "tenant_id", tenantID)
	return nil
}

// seedUser creates a demo user if it doesn't exist, or updates it if it does
func (s *DemoUserSeeder) seedUser(ctx context.Context, user DemoUser) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Check if user already exists
	var existingID int
	err = s.db.QueryRowContext(ctx, "SELECT id FROM user WHERE email = ?", user.Email).Scan(&existingID)
	if err == nil {
		// User exists, update the password and role
		_, err = s.db.ExecContext(ctx,
			"UPDATE user SET password_hash = ?, role = ?, updated_at = NOW() WHERE email = ?",
			string(hashedPassword), user.Role, user.Email)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
		s.logger.Info("Demo user updated", "email", user.Email, "role", user.Role)
		return nil
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("database error: %w", err)
	}

	// User doesn't exist, create it
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, NOW(), NOW())`,
		user.Email, string(hashedPassword), user.Role, user.TenantID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("Demo user created", "email", user.Email, "role", user.Role)
	return nil
}
