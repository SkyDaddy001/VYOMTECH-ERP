package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type TenantService struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewTenantService(db *sql.DB, logger *logger.Logger) *TenantService {
	return &TenantService{
		db:     db,
		logger: logger,
	}
}

// CreateTenant creates a new tenant with default configuration
func (s *TenantService) CreateTenant(ctx context.Context, name, domain string) (*models.Tenant, error) {
	if name == "" {
		return nil, errors.New("tenant name is required")
	}

	tenantID := uuid.New().String()

	// Create tenant
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO tenant (id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at) 
		 VALUES (?, ?, ?, 'active', 100, 50, 1000.00, NOW(), NOW())`,
		tenantID, name, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	tenant := &models.Tenant{
		ID:                 tenantID,
		Name:               name,
		Domain:             domain,
		Status:             "active",
		MaxUsers:           100,
		MaxConcurrentCalls: 50,
		AIBudgetMonthly:    1000.00,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	s.logger.Info("Tenant created successfully", "tenant_id", tenantID, "name", name)
	return tenant, nil
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at FROM tenant WHERE id = ?",
		tenantID).Scan(
		&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers,
		&tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &tenant, nil
}

// GetTenantByDomain retrieves a tenant by domain
func (s *TenantService) GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := s.db.QueryRowContext(ctx,
		"SELECT id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at FROM tenant WHERE domain = ?",
		domain).Scan(
		&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers,
		&tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &tenant, nil
}

// UpdateTenantStatus updates the status of a tenant
func (s *TenantService) UpdateTenantStatus(ctx context.Context, tenantID, status string) error {
	validStatuses := map[string]bool{"active": true, "inactive": true, "suspended": true}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	result, err := s.db.ExecContext(ctx,
		"UPDATE tenant SET status = ?, updated_at = NOW() WHERE id = ?",
		status, tenantID)
	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("tenant not found")
	}

	s.logger.Info("Tenant status updated", "tenant_id", tenantID, "status", status)
	return nil
}

// ListTenants returns all tenants (admin only)
func (s *TenantService) ListTenants(ctx context.Context) ([]*models.Tenant, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, name, domain, status, max_users, max_concurrent_calls, ai_budget_monthly, created_at, updated_at FROM tenant ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var tenants []*models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers,
			&tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		tenants = append(tenants, &tenant)
	}

	return tenants, rows.Err()
}

// GetTenantUserCount returns the number of users in a tenant
func (s *TenantService) GetTenantUserCount(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM user WHERE tenant_id = ?", tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	return count, nil
}

// CanAddUser checks if a tenant can add more users
func (s *TenantService) CanAddUser(ctx context.Context, tenantID string) (bool, error) {
	tenant, err := s.GetTenant(ctx, tenantID)
	if err != nil {
		return false, err
	}

	count, err := s.GetTenantUserCount(ctx, tenantID)
	if err != nil {
		return false, err
	}

	return count < tenant.MaxUsers, nil
}

// DeleteTenant (soft delete - sets status to inactive)
func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE tenant SET status = 'inactive', updated_at = NOW() WHERE id = ?",
		tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("tenant not found")
	}

	s.logger.Info("Tenant deleted (soft delete)", "tenant_id", tenantID)
	return nil
}

// SwitchUserTenant updates the user's current active tenant
func (s *TenantService) SwitchUserTenant(ctx context.Context, userID int64, tenantID string) error {
	// Verify user is a member of the tenant
	isMember, err := s.UserIsTenantMember(ctx, userID, tenantID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("user is not a member of this tenant")
	}

	// Update user's current tenant
	result, err := s.db.ExecContext(ctx,
		"UPDATE `user` SET current_tenant_id = ?, updated_at = NOW() WHERE id = ?",
		tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to switch tenant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}

	s.logger.Info("User switched tenant", "user_id", userID, "tenant_id", tenantID)
	return nil
}

// AddTenantMember adds a user to a tenant with a specified role
func (s *TenantService) AddTenantMember(ctx context.Context, tenantID string, userID int64, email, role string) (*models.TenantMember, error) {
	// Validate role
	validRoles := map[string]bool{"admin": true, "member": true, "viewer": true}
	if !validRoles[role] {
		return nil, errors.New("invalid role")
	}

	memberID := uuid.New().String()

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO tenant_members (id, tenant_id, user_id, email, role, created_at, updated_at) 
		 VALUES (?, ?, ?, ?, ?, NOW(), NOW())`,
		memberID, tenantID, userID, email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to add tenant member: %w", err)
	}

	member := &models.TenantMember{
		ID:        memberID,
		TenantID:  tenantID,
		UserID:    userID,
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.logger.Info("Tenant member added", "tenant_id", tenantID, "user_id", userID, "role", role)
	return member, nil
}

// RemoveTenantMember removes a user from a tenant
func (s *TenantService) RemoveTenantMember(ctx context.Context, tenantID string, userID int64) error {
	// Check if user is the last admin
	adminCount, err := s.GetTenantAdminCount(ctx, tenantID)
	if err != nil {
		return err
	}

	// Get user's role
	isMember, role, err := s.GetTenantMemberRole(ctx, userID, tenantID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("user is not a member of this tenant")
	}

	// Prevent removing the last admin
	if role == "admin" && adminCount <= 1 {
		return errors.New("cannot remove the last admin from the tenant")
	}

	result, err := s.db.ExecContext(ctx,
		"DELETE FROM tenant_members WHERE tenant_id = ? AND user_id = ?",
		tenantID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove tenant member: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("tenant member not found")
	}

	s.logger.Info("Tenant member removed", "tenant_id", tenantID, "user_id", userID)
	return nil
}

// GetTenantMembers returns all members of a tenant
func (s *TenantService) GetTenantMembers(ctx context.Context, tenantID string) ([]*models.TenantMember, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, tenant_id, user_id, email, role, created_at, updated_at 
		 FROM tenant_members WHERE tenant_id = ? ORDER BY created_at DESC`,
		tenantID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var members []*models.TenantMember
	for rows.Next() {
		var member models.TenantMember
		err := rows.Scan(&member.ID, &member.TenantID, &member.UserID, &member.Email,
			&member.Role, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		members = append(members, &member)
	}

	return members, rows.Err()
}

// UserIsTenantMember checks if a user is a member of a tenant
func (s *TenantService) UserIsTenantMember(ctx context.Context, userID int64, tenantID string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tenant_members WHERE user_id = ? AND tenant_id = ?",
		userID, tenantID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}
	return count > 0, nil
}

// GetTenantMemberRole returns the role of a user in a tenant
func (s *TenantService) GetTenantMemberRole(ctx context.Context, userID int64, tenantID string) (bool, string, error) {
	var role string
	err := s.db.QueryRowContext(ctx,
		"SELECT role FROM tenant_members WHERE user_id = ? AND tenant_id = ?",
		userID, tenantID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", fmt.Errorf("database error: %w", err)
	}
	return true, role, nil
}

// UserIsTenantAdmin checks if a user is an admin of a tenant
func (s *TenantService) UserIsTenantAdmin(ctx context.Context, userID int64, tenantID string) (bool, error) {
	isMember, role, err := s.GetTenantMemberRole(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}
	return isMember && role == "admin", nil
}

// GetTenantAdminCount returns the number of admins in a tenant
func (s *TenantService) GetTenantAdminCount(ctx context.Context, tenantID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM tenant_members WHERE tenant_id = ? AND role = 'admin'",
		tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	return count, nil
}

// GetUserTenants returns all tenants a user is a member of
func (s *TenantService) GetUserTenants(ctx context.Context, userID int64) ([]*models.Tenant, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT DISTINCT t.id, t.name, t.domain, t.status, t.max_users, t.max_concurrent_calls, 
		 t.ai_budget_monthly, t.created_at, t.updated_at
		 FROM tenant t
		 INNER JOIN tenant_members tm ON t.id = tm.tenant_id
		 WHERE tm.user_id = ?
		 ORDER BY t.created_at DESC`,
		userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var tenants []*models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Status, &tenant.MaxUsers,
			&tenant.MaxConcurrentCalls, &tenant.AIBudgetMonthly, &tenant.CreatedAt, &tenant.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		tenants = append(tenants, &tenant)
	}

	return tenants, rows.Err()
}
