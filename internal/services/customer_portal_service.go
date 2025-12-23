package services

import (
	"database/sql"
	"errors"
	"time"

	"vyomtech-backend/internal/models"
)

// CustomerPortalService handles customer portal operations
type CustomerPortalService struct {
	db *sql.DB
}

// NewCustomerPortalService creates a new customer portal service
func NewCustomerPortalService(db *sql.DB) *CustomerPortalService {
	return &CustomerPortalService{db: db}
}

// CreateCustomerProfile creates a new customer profile
func (s *CustomerPortalService) CreateCustomerProfile(tenantID, userID int64, req *models.CreateCustomerPortalProfileRequest) (*models.CustomerProfile, error) {
	profile := &models.CustomerProfile{
		TenantID:                    tenantID,
		UserID:                      userID,
		PhoneNumber:                 req.PhoneNumber,
		EmailAddress:                req.EmailAddress,
		DateOfBirth:                 req.DateOfBirth,
		Gender:                      req.Gender,
		CurrentAddress:              req.CurrentAddress,
		City:                        req.City,
		State:                       req.State,
		PostalCode:                  req.PostalCode,
		Country:                     req.Country,
		CommunicationPreference:     req.CommunicationPreference,
		LanguagePreference:          req.LanguagePreference,
		Timezone:                    req.Timezone,
		NotificationEnabled:         true,
		EmailUpdatesEnabled:         true,
		PushNotificationsEnabled:    true,
		ProfileCompletionPercentage: 30,
		CreatedAt:                   time.Now(),
		UpdatedAt:                   time.Now(),
	}

	query := `INSERT INTO customer_profiles (tenant_id, user_id, phone_number, email_address, date_of_birth, gender, current_address, city, state, postal_code, country, communication_preference, language_preference, timezone, notification_enabled, email_updates_enabled, push_notifications_enabled, profile_completion_percentage, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := s.db.QueryRow(query,
		profile.TenantID, profile.UserID, profile.PhoneNumber, profile.EmailAddress,
		profile.DateOfBirth, profile.Gender, profile.CurrentAddress, profile.City,
		profile.State, profile.PostalCode, profile.Country, profile.CommunicationPreference,
		profile.LanguagePreference, profile.Timezone, profile.NotificationEnabled,
		profile.EmailUpdatesEnabled, profile.PushNotificationsEnabled,
		profile.ProfileCompletionPercentage, profile.CreatedAt, profile.UpdatedAt).Scan(&profile.ID)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

// GetCustomerProfile retrieves a customer profile
func (s *CustomerPortalService) GetCustomerProfile(tenantID, userID int64) (*models.CustomerProfile, error) {
	profile := &models.CustomerProfile{}
	query := `SELECT id, tenant_id, user_id, booking_id, phone_number, alternate_phone, email_address, alternate_email, date_of_birth, gender, current_address, permanent_address, city, state, postal_code, country, id_proof_type, id_proof_number, id_proof_file_url, s3_bucket, s3_key, communication_preference, language_preference, timezone, notification_enabled, email_updates_enabled, sms_updates_enabled, push_notifications_enabled, profile_completion_percentage, is_verified, verified_at, verification_notes, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM customer_profiles WHERE tenant_id = ? AND user_id = ? AND deleted_at IS NULL`

	err := s.db.QueryRow(query, tenantID, userID).Scan(
		&profile.ID, &profile.TenantID, &profile.UserID, &profile.BookingID, &profile.PhoneNumber,
		&profile.AlternatePhone, &profile.EmailAddress, &profile.AlternateEmail, &profile.DateOfBirth,
		&profile.Gender, &profile.CurrentAddress, &profile.PermanentAddress, &profile.City, &profile.State,
		&profile.PostalCode, &profile.Country, &profile.IDProofType, &profile.IDProofNumber,
		&profile.IDProofFileURL, &profile.S3Bucket, &profile.S3Key, &profile.CommunicationPreference,
		&profile.LanguagePreference, &profile.Timezone, &profile.NotificationEnabled,
		&profile.EmailUpdatesEnabled, &profile.SMSUpdatesEnabled, &profile.PushNotificationsEnabled,
		&profile.ProfileCompletionPercentage, &profile.IsVerified, &profile.VerifiedAt,
		&profile.VerificationNotes, &profile.Metadata, &profile.CreatedBy, &profile.UpdatedBy,
		&profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("profile not found")
	}
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateCustomerProfile updates a customer profile
func (s *CustomerPortalService) UpdateCustomerProfile(tenantID, userID int64, req *models.UpdateCustomerPortalProfileRequest) error {
	query := `UPDATE customer_profiles SET `
	args := []interface{}{}
	fieldCount := 0

	if req.PhoneNumber != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "phone_number = ?"
		args = append(args, *req.PhoneNumber)
		fieldCount++
	}
	if req.AlternatePhone != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "alternate_phone = ?"
		args = append(args, *req.AlternatePhone)
		fieldCount++
	}
	if req.EmailAddress != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "email_address = ?"
		args = append(args, *req.EmailAddress)
		fieldCount++
	}
	if req.CurrentAddress != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "current_address = ?"
		args = append(args, *req.CurrentAddress)
		fieldCount++
	}
	if req.City != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "city = ?"
		args = append(args, *req.City)
		fieldCount++
	}
	if req.State != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "state = ?"
		args = append(args, *req.State)
		fieldCount++
	}
	if req.PostalCode != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "postal_code = ?"
		args = append(args, *req.PostalCode)
		fieldCount++
	}
	if req.LanguagePreference != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "language_preference = ?"
		args = append(args, *req.LanguagePreference)
		fieldCount++
	}
	if req.Timezone != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "timezone = ?"
		args = append(args, *req.Timezone)
		fieldCount++
	}

	if fieldCount == 0 {
		return errors.New("no fields to update")
	}

	query += " updated_at = ? WHERE tenant_id = ? AND user_id = ?"
	args = append(args, time.Now(), tenantID, userID)

	_, err := s.db.Exec(query, args...)
	return err
}

// CreateNotification creates a notification
func (s *CustomerPortalService) CreateNotification(tenantID, userID int64, notificationType, title, message string) (*models.CustomerNotification, error) {
	notif := &models.CustomerNotification{
		TenantID:         tenantID,
		UserID:           userID,
		NotificationType: notificationType,
		Title:            title,
		Message:          message,
		Category:         "status_update",
		Priority:         "normal",
		IsRead:           false,
		DeliveryStatus:   "pending",
		InAppSent:        true,
		InAppSentAt:      now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `INSERT INTO customer_notifications (tenant_id, user_id, notification_type, title, message, category, priority, is_read, delivery_status, in_app_sent, in_app_sent_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := s.db.QueryRow(query,
		notif.TenantID, notif.UserID, notif.NotificationType, notif.Title, notif.Message,
		notif.Category, notif.Priority, notif.IsRead, notif.DeliveryStatus, notif.InAppSent,
		notif.InAppSentAt, notif.CreatedAt, notif.UpdatedAt).Scan(&notif.ID)

	if err != nil {
		return nil, err
	}

	return notif, nil
}

// GetNotifications retrieves notifications for a user
func (s *CustomerPortalService) GetNotifications(tenantID, userID int64, limit, offset int) ([]*models.CustomerNotification, int, error) {
	notifications := []*models.CustomerNotification{}

	countQuery := `SELECT COUNT(*) FROM customer_notifications WHERE tenant_id = ? AND user_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, tenantID, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, user_id, notification_type, title, message, category, priority, is_read, is_archived, delivery_status, created_at, updated_at FROM customer_notifications WHERE tenant_id = ? AND user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, tenantID, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		notif := &models.CustomerNotification{}
		err := rows.Scan(
			&notif.ID, &notif.TenantID, &notif.UserID, &notif.NotificationType, &notif.Title,
			&notif.Message, &notif.Category, &notif.Priority, &notif.IsRead, &notif.IsArchived,
			&notif.DeliveryStatus, &notif.CreatedAt, &notif.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, total, rows.Err()
}

// MarkNotificationAsRead marks a notification as read
func (s *CustomerPortalService) MarkNotificationAsRead(tenantID, notificationID int64) error {
	query := `UPDATE customer_notifications SET is_read = TRUE, read_at = ? WHERE id = ? AND tenant_id = ?`
	_, err := s.db.Exec(query, time.Now(), notificationID, tenantID)
	return err
}

// CreateConversation creates a support conversation
func (s *CustomerPortalService) CreateConversation(tenantID, customerUserID int64, subject *string, conversationType string) (*models.CustomerConversation, error) {
	conv := &models.CustomerConversation{
		TenantID:         tenantID,
		CustomerUserID:   customerUserID,
		Subject:          subject,
		ConversationType: conversationType,
		Status:           "open",
		Priority:         "normal",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `INSERT INTO customer_conversations (tenant_id, customer_user_id, subject, conversation_type, status, priority, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := s.db.QueryRow(query,
		conv.TenantID, conv.CustomerUserID, conv.Subject, conv.ConversationType,
		conv.Status, conv.Priority, conv.CreatedAt, conv.UpdatedAt).Scan(&conv.ID)

	if err != nil {
		return nil, err
	}

	return conv, nil
}

// SendMessage sends a message in a conversation
func (s *CustomerPortalService) SendMessage(tenantID, conversationID, senderUserID int64, messageText, senderType string) (*models.CustomerMessage, error) {
	msg := &models.CustomerMessage{
		TenantID:       tenantID,
		ConversationID: conversationID,
		SenderUserID:   senderUserID,
		SenderType:     senderType,
		MessageText:    messageText,
		MessageType:    "text",
		IsRead:         false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO customer_messages (tenant_id, conversation_id, sender_user_id, sender_type, message_text, message_type, is_read, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := s.db.QueryRow(query,
		msg.TenantID, msg.ConversationID, msg.SenderUserID, msg.SenderType,
		msg.MessageText, msg.MessageType, msg.IsRead, msg.CreatedAt, msg.UpdatedAt).Scan(&msg.ID)

	if err != nil {
		return nil, err
	}

	// Update conversation last message
	updateQuery := `UPDATE customer_conversations SET last_message_at = ?, last_message_from = ?, updated_at = ? WHERE id = ?`
	s.db.Exec(updateQuery, time.Now(), senderType, time.Now(), conversationID)

	return msg, nil
}

// GetConversationMessages retrieves messages from a conversation
func (s *CustomerPortalService) GetConversationMessages(tenantID, conversationID int64, limit, offset int) ([]*models.CustomerMessage, int, error) {
	messages := []*models.CustomerMessage{}

	countQuery := `SELECT COUNT(*) FROM customer_messages WHERE tenant_id = ? AND conversation_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, tenantID, conversationID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, conversation_id, sender_user_id, sender_type, message_text, message_type, is_read, created_at FROM customer_messages WHERE tenant_id = ? AND conversation_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, tenantID, conversationID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		msg := &models.CustomerMessage{}
		err := rows.Scan(
			&msg.ID, &msg.TenantID, &msg.ConversationID, &msg.SenderUserID, &msg.SenderType,
			&msg.MessageText, &msg.MessageType, &msg.IsRead, &msg.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, msg)
	}

	return messages, total, rows.Err()
}

// CreateDocumentUpload creates a document upload record
func (s *CustomerPortalService) CreateDocumentUpload(tenantID, userID int64, req *models.CreatePortalDocumentUploadRequest) (*models.CustomerDocumentUpload, error) {
	doc := &models.CustomerDocumentUpload{
		TenantID:                 tenantID,
		UserID:                   userID,
		DocumentType:             &req.DocumentType,
		DocumentName:             req.DocumentName,
		DocumentDescription:      req.DocumentDescription,
		UploadStatus:             "uploading",
		UploadProgressPercentage: 0,
		VerificationStatus:       "pending",
		IsRequiredDocument:       req.IsRequiredDocument,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	query := `INSERT INTO customer_document_uploads (tenant_id, user_id, document_type, document_name, document_description, upload_status, upload_progress_percentage, verification_status, is_required_document, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := s.db.QueryRow(query,
		doc.TenantID, doc.UserID, doc.DocumentType, doc.DocumentName, doc.DocumentDescription,
		doc.UploadStatus, doc.UploadProgressPercentage, doc.VerificationStatus, doc.IsRequiredDocument,
		doc.CreatedAt, doc.UpdatedAt).Scan(&doc.ID)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

// UpdateDocumentUploadStatus updates document upload status
func (s *CustomerPortalService) UpdateDocumentUploadStatus(tenantID, documentID int64, status string, progress float64) error {
	query := `UPDATE customer_document_uploads SET upload_status = ?, upload_progress_percentage = ?, updated_at = ? WHERE id = ? AND tenant_id = ?`
	_, err := s.db.Exec(query, status, progress, time.Now(), documentID, tenantID)
	return err
}

// GetCustomerBookingTracking retrieves booking tracking
func (s *CustomerPortalService) GetCustomerBookingTracking(tenantID, bookingID int64) (*models.CustomerBookingTracking, error) {
	tracking := &models.CustomerBookingTracking{}
	query := `SELECT id, tenant_id, booking_id, user_id, current_status, status_changed_at, property_name, property_location, booking_date, possession_date, estimated_handover_date, total_amount, amount_paid, amount_pending, payment_percentage, required_documents_count, uploaded_documents_count, verified_documents_count, last_update_at, updated_at FROM customer_booking_tracking WHERE tenant_id = ? AND booking_id = ? AND deleted_at IS NULL`

	err := s.db.QueryRow(query, tenantID, bookingID).Scan(
		&tracking.ID, &tracking.TenantID, &tracking.BookingID, &tracking.UserID, &tracking.CurrentStatus,
		&tracking.StatusChangedAt, &tracking.PropertyName, &tracking.PropertyLocation, &tracking.BookingDate,
		&tracking.PossessionDate, &tracking.EstimatedHandoverDate, &tracking.TotalAmount, &tracking.AmountPaid,
		&tracking.AmountPending, &tracking.PaymentPercentage, &tracking.RequiredDocumentsCount,
		&tracking.UploadedDocumentsCount, &tracking.VerifiedDocumentsCount, &tracking.LastUpdateAt,
		&tracking.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("booking tracking not found")
	}
	if err != nil {
		return nil, err
	}

	return tracking, nil
}

// GetPaymentTracking retrieves payment tracking
func (s *CustomerPortalService) GetPaymentTracking(tenantID, bookingID int64, limit, offset int) ([]*models.CustomerPaymentTracking, int, error) {
	payments := []*models.CustomerPaymentTracking{}

	countQuery := `SELECT COUNT(*) FROM customer_payment_tracking WHERE tenant_id = ? AND booking_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, tenantID, bookingID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, booking_id, user_id, invoice_number, invoice_date, due_date, invoice_amount, total_amount, payment_status, payment_date, amount_paid FROM customer_payment_tracking WHERE tenant_id = ? AND booking_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := s.db.Query(query, tenantID, bookingID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		payment := &models.CustomerPaymentTracking{}
		err := rows.Scan(
			&payment.ID, &payment.TenantID, &payment.BookingID, &payment.UserID, &payment.InvoiceNumber,
			&payment.InvoiceDate, &payment.DueDate, &payment.InvoiceAmount, &payment.TotalAmount,
			&payment.PaymentStatus, &payment.PaymentDate, &payment.AmountPaid)
		if err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	return payments, total, rows.Err()
}

// CreateFeedback creates customer feedback
func (s *CustomerPortalService) CreateFeedback(tenantID, userID int64, req *models.CreatePortalFeedbackRequest) (*models.CustomerFeedback, error) {
	feedback := &models.CustomerFeedback{
		TenantID:            tenantID,
		UserID:              userID,
		FeedbackType:        req.FeedbackType,
		Subject:             req.Subject,
		Message:             req.Message,
		OverallRating:       req.OverallRating,
		ServiceRating:       req.ServiceRating,
		CommunicationRating: req.CommunicationRating,
		DocumentationRating: req.DocumentationRating,
		FeedbackStatus:      "open",
		AttachmentsCount:    0,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	query := `INSERT INTO customer_feedback (tenant_id, user_id, feedback_type, subject, message, overall_rating, service_rating, communication_rating, documentation_rating, feedback_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	err := s.db.QueryRow(query,
		feedback.TenantID, feedback.UserID, feedback.FeedbackType, feedback.Subject, feedback.Message,
		feedback.OverallRating, feedback.ServiceRating, feedback.CommunicationRating, feedback.DocumentationRating,
		feedback.FeedbackStatus, feedback.CreatedAt, feedback.UpdatedAt).Scan(&feedback.ID)

	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// UpdatePreferences updates customer preferences
func (s *CustomerPortalService) UpdatePreferences(tenantID, userID int64, req *models.UpdatePortalPreferencesRequest) error {
	query := `UPDATE customer_preferences SET `
	args := []interface{}{}
	fieldCount := 0

	if req.EmailNotifications != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "email_notifications = ?"
		args = append(args, *req.EmailNotifications)
		fieldCount++
	}
	if req.SMSNotifications != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "sms_notifications = ?"
		args = append(args, *req.SMSNotifications)
		fieldCount++
	}
	if req.PushNotifications != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "push_notifications = ?"
		args = append(args, *req.PushNotifications)
		fieldCount++
	}
	if req.ReceiveBookingUpdates != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "receive_booking_updates = ?"
		args = append(args, *req.ReceiveBookingUpdates)
		fieldCount++
	}
	if req.ReceivePaymentReminders != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "receive_payment_reminders = ?"
		args = append(args, *req.ReceivePaymentReminders)
		fieldCount++
	}
	if req.PreferredLanguage != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "preferred_language = ?"
		args = append(args, *req.PreferredLanguage)
		fieldCount++
	}
	if req.ThemePreference != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "theme_preference = ?"
		args = append(args, *req.ThemePreference)
		fieldCount++
	}

	if fieldCount == 0 {
		return errors.New("no fields to update")
	}

	query += " updated_at = ? WHERE tenant_id = ? AND user_id = ?"
	args = append(args, time.Now(), tenantID, userID)

	_, err := s.db.Exec(query, args...)
	return err
}

// LogActivity logs customer activity
func (s *CustomerPortalService) LogActivity(tenantID, userID int64, activityType, entityType string, entityID int64) error {
	query := `INSERT INTO customer_activity_log (tenant_id, user_id, activity_type, entity_type, entity_id, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, tenantID, userID, activityType, entityType, entityID, time.Now())
	return err
}
