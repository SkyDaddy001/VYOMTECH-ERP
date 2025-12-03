package services

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"vyomtech-backend/pkg/logger"
)

type PasswordResetService struct {
	db           *sql.DB
	emailService *EmailService
	logger       *logger.Logger
}

func NewPasswordResetService(db *sql.DB, emailService *EmailService, logger *logger.Logger) *PasswordResetService {
	return &PasswordResetService{
		db:           db,
		emailService: emailService,
		logger:       logger,
	}
}

func (s *PasswordResetService) RequestPasswordReset(email string) error {
	// Check if user exists
	var userID int
	err := s.db.QueryRow("SELECT id FROM user WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Generate reset token
	token := make([]byte, 32)
	rand.Read(token)
	resetToken := fmt.Sprintf("%x", token)

	// Store token with expiration
	_, err = s.db.Exec(`
        INSERT INTO password_reset_tokens (user_id, token, expires_at)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE token = ?, expires_at = ?`,
		userID, resetToken, time.Now().Add(1*time.Hour), resetToken, time.Now().Add(1*time.Hour))
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send reset email
	resetLink := fmt.Sprintf("https://app.example.com/reset-password?token=%s", resetToken)
	err = s.emailService.SendPasswordResetEmail(email, resetLink)
	if err != nil {
		s.logger.Error("Failed to send password reset email", "error", err)
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	s.logger.Info("Password reset requested", "user_id", userID)
	return nil
}

func (s *PasswordResetService) ResetPassword(token, newPassword string) error {
	// Verify token
	var userID int
	var expiresAt time.Time
	err := s.db.QueryRow(`
        SELECT user_id, expires_at FROM password_reset_tokens
        WHERE token = ? AND expires_at > NOW()`, token).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invalid or expired token")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	_, err = s.db.Exec("UPDATE user SET password_hash = ? WHERE id = ?", string(hashedPassword), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete used token
	s.db.Exec("DELETE FROM password_reset_tokens WHERE token = ?", token)

	s.logger.Info("Password reset successfully", "user_id", userID)
	return nil
}
