package services

import (
	"fmt"
	"net/smtp"

	"vyomtech-backend/internal/config"
	"vyomtech-backend/pkg/logger"
)

type EmailService struct {
	config *config.EmailConfig
	logger *logger.Logger
}

func NewEmailService(cfg *config.EmailConfig, logger *logger.Logger) *EmailService {
	return &EmailService{
		config: cfg,
		logger: logger,
	}
}

func (s *EmailService) SendPasswordResetEmail(email, resetLink string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		<h2>Password Reset</h2>
		<p>Click the link below to reset your password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
	`, resetLink)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	return s.sendEmail(to, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	from := s.config.FromEmail
	password := s.config.SMTPPassword
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, to, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	if err != nil {
		s.logger.Error("Failed to send email", "error", err, "to", to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Email sent successfully", "to", to, "subject", subject)
	return nil
}
