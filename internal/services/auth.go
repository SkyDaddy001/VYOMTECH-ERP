package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/auth"
	"vyomtech-backend/pkg/logger"
)

type AuthService struct {
	db         *sql.DB
	jwtManager *auth.JWTManager
	logger     *logger.Logger
}

func NewAuthService(db *sql.DB, jwtManager *auth.JWTManager, logger *logger.Logger) *AuthService {
	return &AuthService{
		db:         db,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, role, tenantID string) (*models.User, error) {
	// Check if user already exists
	var existingID int
	err := s.db.QueryRowContext(ctx, "SELECT id FROM user WHERE email = ?", email).Scan(&existingID)
	if err == nil {
		return nil, errors.New("user already exists")
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// If no tenant ID provided, create a default tenant for the user
	if tenantID == "" {
		// Generate a simple tenant ID (email-based)
		tenantID = fmt.Sprintf("tenant-%s", email)

		// Create tenant if it doesn't exist
		var existingTenant string
		err := s.db.QueryRowContext(ctx, "SELECT id FROM tenant WHERE id = ?", tenantID).Scan(&existingTenant)
		if err == sql.ErrNoRows {
			// Tenant doesn't exist, create it
			_, err = s.db.ExecContext(ctx,
				"INSERT INTO tenant (id, name, status, max_users, max_concurrent_calls, ai_budget_monthly) VALUES (?, ?, 'active', 100, 50, 1000.00)",
				tenantID, email)
			if err != nil {
				return nil, fmt.Errorf("failed to create tenant: %w", err)
			}
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Default role to 'admin' if not provided (first user of tenant should be admin)
	if role == "" {
		role = "admin"
	}

	// Insert user
	result, err := s.db.ExecContext(ctx,
		"INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
		email, string(hashedPassword), role, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	user := &models.User{
		ID:        int(userID),
		Email:     email,
		Role:      role,
		TenantID:  tenantID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.logger.WithUser(int(userID)).Info("User registered successfully")
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx,
		"SELECT id, email, password_hash, role, tenant_id FROM user WHERE email = ?",
		email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.TenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Role, user.TenantID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.WithUser(user.ID).Info("User logged in successfully")
	return token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	userIDFloat, ok := (*claims)["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}
	userID := int(userIDFloat)

	var user models.User
	err = s.db.QueryRow(
		"SELECT id, email, role, tenant_id FROM user WHERE id = ?",
		userID).Scan(&user.ID, &user.Email, &user.Role, &user.TenantID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	// Get current password hash
	var currentHash string
	err := s.db.QueryRowContext(ctx, "SELECT password_hash FROM user WHERE id = ?", userID).Scan(&currentHash)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	_, err = s.db.ExecContext(ctx,
		"UPDATE user SET password_hash = ?, updated_at = NOW() WHERE id = ?",
		string(newHash), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.WithUser(userID).Info("Password changed successfully")
	return nil
}

func (s *AuthService) GenerateToken(userID int, email, role, tenantID string) (string, error) {
	return s.jwtManager.GenerateToken(userID, email, role, tenantID)
}
