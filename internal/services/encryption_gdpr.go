package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

// EncryptionService handles data encryption for sensitive fields
type EncryptionService struct {
	db        *sql.DB
	logger    *logger.Logger
	encKey    []byte // 32-byte AES-256 key
	nonceSize int
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService(db *sql.DB, log *logger.Logger, encryptionKey string) (*EncryptionService, error) {
	// Convert string key to bytes - should be 32 bytes for AES-256
	key := []byte(encryptionKey)
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}

	return &EncryptionService{
		db:        db,
		logger:    log,
		encKey:    key,
		nonceSize: 12, // 96-bit nonce for GCM
	}, nil
}

// EncryptField encrypts a sensitive field value
func (es *EncryptionService) EncryptField(fieldValue string) (string, error) {
	if fieldValue == "" {
		return "", nil
	}

	block, err := aes.NewCipher(es.encKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, es.nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(fieldValue), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptField decrypts a sensitive field value
func (es *EncryptionService) DecryptField(encryptedValue string) (string, error) {
	if encryptedValue == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(es.encKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce, ciphertext := data[:es.nonceSize], data[es.nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// StoreEncryptedField stores an encrypted field in the database
func (es *EncryptionService) StoreEncryptedField(ctx context.Context, tenantID, fieldName string, fieldValue string, resourceType string, resourceID int64) error {
	encryptedValue, err := es.EncryptField(fieldValue)
	if err != nil {
		es.logger.Error("Failed to encrypt field", "error", err, "field", fieldName)
		return err
	}

	query := `
		INSERT INTO data_encryption (tenant_id, field_name, field_value, resource_type, resource_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = es.db.ExecContext(ctx, query, tenantID, fieldName, encryptedValue, resourceType, resourceID, time.Now(), time.Now())
	if err != nil {
		es.logger.Error("Failed to store encrypted field", "error", err)
		return err
	}

	return nil
}

// GetEncryptedField retrieves and decrypts a field
func (es *EncryptionService) GetEncryptedField(ctx context.Context, tenantID, fieldName, resourceType string, resourceID int64) (string, error) {
	query := `
		SELECT field_value
		FROM data_encryption
		WHERE tenant_id = ? AND field_name = ? AND resource_type = ? AND resource_id = ?
	`

	var encryptedValue string
	err := es.db.QueryRowContext(ctx, query, tenantID, fieldName, resourceType, resourceID).Scan(&encryptedValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return es.DecryptField(encryptedValue)
}

// RotateEncryptionKey rotates the encryption key for all stored data
// This is a demonstration - in production, use key management services
func (es *EncryptionService) RotateEncryptionKey(ctx context.Context, tenantID string, newKey []byte) error {
	if len(newKey) != 32 {
		return errors.New("new encryption key must be 32 bytes for AES-256")
	}

	// Fetch all encrypted records
	query := `
		SELECT id, field_value
		FROM data_encryption
		WHERE tenant_id = ?
	`

	rows, err := es.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Re-encrypt all data with new key
	oldKey := es.encKey
	es.encKey = newKey
	newEncService := &EncryptionService{
		db:        es.db,
		logger:    es.logger,
		encKey:    newKey,
		nonceSize: es.nonceSize,
	}

	var updates []struct {
		id    int64
		value string
	}

	for rows.Next() {
		var id int64
		var encryptedValue string
		if err := rows.Scan(&id, &encryptedValue); err != nil {
			return err
		}

		// Decrypt with old key
		es.encKey = oldKey
		decrypted, err := es.DecryptField(encryptedValue)
		if err != nil {
			es.logger.Warn("Failed to decrypt field during rotation", "error", err, "id", id)
			continue
		}

		// Encrypt with new key
		reEncrypted, err := newEncService.EncryptField(decrypted)
		if err != nil {
			es.logger.Warn("Failed to re-encrypt field during rotation", "error", err, "id", id)
			continue
		}

		updates = append(updates, struct {
			id    int64
			value string
		}{id, reEncrypted})
	}

	// Update all records
	for _, update := range updates {
		updateQuery := `UPDATE data_encryption SET field_value = ?, updated_at = ? WHERE id = ?`
		_, err := es.db.ExecContext(ctx, updateQuery, update.value, time.Now(), update.id)
		if err != nil {
			es.logger.Error("Failed to update encrypted field", "error", err, "id", update.id)
		}
	}

	es.encKey = newKey
	es.logger.Info("Encryption key rotated successfully", "tenant", tenantID, "records_updated", len(updates))

	return nil
}

// GDPRService handles GDPR compliance requests
type GDPRService struct {
	db         *sql.DB
	logger     *logger.Logger
	encService *EncryptionService
}

// NewGDPRService creates a new GDPR service
func NewGDPRService(db *sql.DB, log *logger.Logger, encService *EncryptionService) *GDPRService {
	return &GDPRService{
		db:         db,
		logger:     log,
		encService: encService,
	}
}

// CreateDataAccessRequest creates a GDPR data access request
func (gs *GDPRService) CreateDataAccessRequest(ctx context.Context, tenantID string, userID int64) (*models.GDPRRequest, error) {
	query := `
		INSERT INTO gdpr_requests (tenant_id, user_id, type, status, reason, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL 30 DAY), ?, ?)
	`

	result, err := gs.db.ExecContext(ctx, query, tenantID, userID, "access", "pending", "Data access request", time.Now(), time.Now())
	if err != nil {
		gs.logger.Error("Failed to create data access request", "error", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	request := &models.GDPRRequest{
		ID:        id,
		TenantID:  tenantID,
		UserID:    userID,
		Type:      "access",
		Status:    "pending",
		Reason:    "Data access request",
		ExpiresAt: time.Now().AddDate(0, 0, 30),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return request, nil
}

// CreateDataDeletionRequest creates a GDPR data deletion request
func (gs *GDPRService) CreateDataDeletionRequest(ctx context.Context, tenantID string, userID int64, reason string) (*models.GDPRRequest, error) {
	query := `
		INSERT INTO gdpr_requests (tenant_id, user_id, type, status, reason, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL 30 DAY), ?, ?)
	`

	result, err := gs.db.ExecContext(ctx, query, tenantID, userID, "deletion", "pending", reason, time.Now(), time.Now())
	if err != nil {
		gs.logger.Error("Failed to create data deletion request", "error", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	request := &models.GDPRRequest{
		ID:        id,
		TenantID:  tenantID,
		UserID:    userID,
		Type:      "deletion",
		Status:    "pending",
		Reason:    reason,
		ExpiresAt: time.Now().AddDate(0, 0, 30),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return request, nil
}

// GetGDPRRequest retrieves a GDPR request
func (gs *GDPRService) GetGDPRRequest(ctx context.Context, tenantID string, requestID int64) (*models.GDPRRequest, error) {
	query := `
		SELECT id, tenant_id, user_id, type, status, reason, expires_at, created_at, updated_at
		FROM gdpr_requests
		WHERE id = ? AND tenant_id = ?
	`

	request := &models.GDPRRequest{}
	err := gs.db.QueryRowContext(ctx, query, requestID, tenantID).Scan(
		&request.ID, &request.TenantID, &request.UserID, &request.Type, &request.Status,
		&request.Reason, &request.ExpiresAt, &request.CreatedAt, &request.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("request not found")
		}
		return nil, err
	}

	return request, nil
}

// ExportUserData exports all data for a user (GDPR portability)
func (gs *GDPRService) ExportUserData(ctx context.Context, tenantID string, userID int64) (map[string]interface{}, error) {
	userData := make(map[string]interface{})

	// Export user profile
	userQuery := `SELECT id, email, role, tenant_id, created_at, updated_at FROM users WHERE id = ? AND tenant_id = ?`
	var id int64
	var email, role, tenantIDVal string
	var createdAt, updatedAt time.Time

	err := gs.db.QueryRowContext(ctx, userQuery, userID, tenantID).Scan(
		&id, &email, &role, &tenantIDVal, &createdAt, &updatedAt,
	)
	if err == nil {
		userData["user_profile"] = map[string]interface{}{
			"id":         id,
			"email":      email,
			"role":       role,
			"tenant_id":  tenantIDVal,
			"created_at": createdAt,
			"updated_at": updatedAt,
		}
	}

	// Export leads
	leadsQuery := `SELECT id, title, status, created_at FROM leads WHERE assigned_to = ? AND tenant_id = ? LIMIT 1000`
	rows, err := gs.db.QueryContext(ctx, leadsQuery, userID, tenantID)
	if err == nil {
		defer rows.Close()
		var leads []map[string]interface{}
		for rows.Next() {
			var leadID int64
			var title, status string
			var createdAt time.Time
			if rows.Scan(&leadID, &title, &status, &createdAt) == nil {
				leads = append(leads, map[string]interface{}{
					"id":         leadID,
					"title":      title,
					"status":     status,
					"created_at": createdAt,
				})
			}
		}
		userData["leads"] = leads
	}

	// Export calls
	callsQuery := `SELECT id, phone_number, duration, status, created_at FROM calls WHERE agent_id = ? AND tenant_id = ? LIMIT 1000`
	rows, err = gs.db.QueryContext(ctx, callsQuery, userID, tenantID)
	if err == nil {
		defer rows.Close()
		var calls []map[string]interface{}
		for rows.Next() {
			var callID, duration int64
			var phoneNumber, status string
			var createdAt time.Time
			if rows.Scan(&callID, &phoneNumber, &duration, &status, &createdAt) == nil {
				calls = append(calls, map[string]interface{}{
					"id":           callID,
					"phone_number": phoneNumber,
					"duration":     duration,
					"status":       status,
					"created_at":   createdAt,
				})
			}
		}
		userData["calls"] = calls
	}

	// Export consent records
	consentQuery := `SELECT type, given, consented_at FROM consent_records WHERE user_id = ? AND tenant_id = ?`
	rows, err = gs.db.QueryContext(ctx, consentQuery, userID, tenantID)
	if err == nil {
		defer rows.Close()
		var consents []map[string]interface{}
		for rows.Next() {
			var consentType string
			var given bool
			var consentedAt time.Time
			if rows.Scan(&consentType, &given, &consentedAt) == nil {
				consents = append(consents, map[string]interface{}{
					"type":         consentType,
					"given":        given,
					"consented_at": consentedAt,
				})
			}
		}
		userData["consent_records"] = consents
	}

	// Export audit logs
	auditQuery := `SELECT action, resource, details, created_at FROM audit_logs WHERE user_id = ? AND tenant_id = ? ORDER BY created_at DESC LIMIT 1000`
	rows, err = gs.db.QueryContext(ctx, auditQuery, userID, tenantID)
	if err == nil {
		defer rows.Close()
		var auditLogs []map[string]interface{}
		for rows.Next() {
			var action, resource, details string
			var createdAt time.Time
			if rows.Scan(&action, &resource, &details, &createdAt) == nil {
				var detailsMap map[string]interface{}
				json.Unmarshal([]byte(details), &detailsMap)
				auditLogs = append(auditLogs, map[string]interface{}{
					"action":     action,
					"resource":   resource,
					"details":    detailsMap,
					"created_at": createdAt,
				})
			}
		}
		userData["audit_logs"] = auditLogs
	}

	return userData, nil
}

// DeleteUserData permanently deletes all user data (GDPR right to be forgotten)
func (gs *GDPRService) DeleteUserData(ctx context.Context, tenantID string, userID int64) error {
	tx, err := gs.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// Delete in order of foreign key dependencies
	tables := []string{
		"audit_logs",
		"security_events",
		"data_encryption",
		"consent_records",
		"gdpr_requests",
		"user_roles",
		"calls",
		"leads",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s WHERE user_id = ? AND tenant_id = ?", table)
		_, err := tx.ExecContext(ctx, query, userID, tenantID)
		if err != nil {
			gs.logger.Warn("Failed to delete from table", "table", table, "error", err)
			// Continue with other tables
		}
	}

	// Finally, anonymize user record
	anonQuery := `
		UPDATE users
		SET email = CONCAT('deleted_', id, '@deleted.local'), 
		    password_hash = '', 
		    updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`
	_, err = tx.ExecContext(ctx, anonQuery, time.Now(), userID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to anonymize user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	gs.logger.Info("User data deleted successfully", "user_id", userID, "tenant_id", tenantID)

	return nil
}

// RecordConsent records user consent for data processing
func (gs *GDPRService) RecordConsent(ctx context.Context, tenantID string, userID int64, consentType string, given bool) error {
	query := `
		INSERT INTO consent_records (tenant_id, user_id, type, given, consented_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := gs.db.ExecContext(ctx, query, tenantID, userID, consentType, given, time.Now(), time.Now(), time.Now())
	if err != nil {
		gs.logger.Error("Failed to record consent", "error", err)
		return err
	}

	return nil
}

// GetUserConsents retrieves user consent records
func (gs *GDPRService) GetUserConsents(ctx context.Context, tenantID string, userID int64) ([]models.ConsentRecord, error) {
	query := `
		SELECT id, tenant_id, user_id, type, given, consented_at, expires_at, created_at, updated_at
		FROM consent_records
		WHERE tenant_id = ? AND user_id = ?
		ORDER BY consented_at DESC
	`

	rows, err := gs.db.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consents []models.ConsentRecord
	for rows.Next() {
		var consent models.ConsentRecord
		err := rows.Scan(&consent.ID, &consent.TenantID, &consent.UserID, &consent.Type, &consent.Given,
			&consent.ConsentedAt, &consent.ExpiresAt, &consent.CreatedAt, &consent.UpdatedAt)
		if err != nil {
			return nil, err
		}
		consents = append(consents, consent)
	}

	return consents, rows.Err()
}
