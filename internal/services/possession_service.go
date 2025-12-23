package services

import (
	"database/sql"
	"encoding/json"
	"vyomtech-backend/internal/models"
)

// PossessionService handles possession management operations
type PossessionService struct {
	db *sql.DB
}

// NewPossessionService creates a new possession service
func NewPossessionService(db *sql.DB) *PossessionService {
	return &PossessionService{db: db}
}

// CreatePossessionStatus creates a new possession status
func (s *PossessionService) CreatePossessionStatus(tenantID, bookingID int64, status, possessionType string, notes *string, createdBy *int64) (*models.PossessionStatus, error) {
	query := `
		INSERT INTO possession_statuses (tenant_id, booking_id, status, possession_type, notes, created_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, tenantID, bookingID, status, possessionType, notes, createdBy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetPossessionStatus(id)
}

// GetPossessionStatus retrieves a possession status by ID
func (s *PossessionService) GetPossessionStatus(id int64) (*models.PossessionStatus, error) {
	ps := &models.PossessionStatus{}
	query := `SELECT id, tenant_id, booking_id, status, possession_date, estimated_possession_date, possession_reason, possession_type, is_complete, completion_percentage, notes, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_statuses WHERE id = ? AND deleted_at IS NULL`
	err := s.db.QueryRow(query, id).Scan(&ps.ID, &ps.TenantID, &ps.BookingID, &ps.Status, &ps.PossessionDate, &ps.EstimatedPossessionDate, &ps.PossessionReason, &ps.PossessionType, &ps.IsComplete, &ps.CompletionPercentage, &ps.Notes, &ps.CreatedBy, &ps.UpdatedBy, &ps.CreatedAt, &ps.UpdatedAt, &ps.DeletedAt)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

// ListPossessionStatuses lists all possession statuses for a tenant
func (s *PossessionService) ListPossessionStatuses(tenantID int64, limit, offset int) ([]*models.PossessionStatus, int, error) {
	var statuses []*models.PossessionStatus
	countQuery := `SELECT COUNT(*) FROM possession_statuses WHERE tenant_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, booking_id, status, possession_date, estimated_possession_date, possession_reason, possession_type, is_complete, completion_percentage, notes, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_statuses WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		ps := &models.PossessionStatus{}
		err := rows.Scan(&ps.ID, &ps.TenantID, &ps.BookingID, &ps.Status, &ps.PossessionDate, &ps.EstimatedPossessionDate, &ps.PossessionReason, &ps.PossessionType, &ps.IsComplete, &ps.CompletionPercentage, &ps.Notes, &ps.CreatedBy, &ps.UpdatedBy, &ps.CreatedAt, &ps.UpdatedAt, &ps.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		statuses = append(statuses, ps)
	}

	return statuses, total, nil
}

// UpdatePossessionStatus updates a possession status
func (s *PossessionService) UpdatePossessionStatus(id int64, updates map[string]interface{}, updatedBy *int64) (*models.PossessionStatus, error) {
	query := `UPDATE possession_statuses SET `
	args := []interface{}{}
	count := 0

	for key, value := range updates {
		if count > 0 {
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		count++
	}

	query += ", updated_by = ? WHERE id = ?"
	args = append(args, updatedBy, id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetPossessionStatus(id)
}

// DeletePossessionStatus soft deletes a possession status
func (s *PossessionService) DeletePossessionStatus(id int64) error {
	query := `UPDATE possession_statuses SET deleted_at = NOW() WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// CreatePossessionDocument creates a new possession document
func (s *PossessionService) CreatePossessionDocument(tenantID, possessionID int64, docType, docName string, isMandatory bool, metadata *models.JSONMap, uploadedBy *int64) (*models.PossessionDocument, error) {
	metadataJSON, _ := json.Marshal(metadata)
	query := `
		INSERT INTO possession_documents (tenant_id, possession_id, document_type, document_name, is_mandatory, metadata, uploaded_by, uploaded_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`
	result, err := s.db.Exec(query, tenantID, possessionID, docType, docName, isMandatory, metadataJSON, uploadedBy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetPossessionDocument(id)
}

// GetPossessionDocument retrieves a possession document by ID
func (s *PossessionService) GetPossessionDocument(id int64) (*models.PossessionDocument, error) {
	pd := &models.PossessionDocument{}
	var metadataJSON []byte
	query := `SELECT id, tenant_id, possession_id, document_type, document_name, document_url, file_name, file_size, file_format, s3_bucket, s3_key, document_status, verification_notes, verified_by, verified_at, is_mandatory, uploaded_by, uploaded_at, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_documents WHERE id = ? AND deleted_at IS NULL`
	err := s.db.QueryRow(query, id).Scan(&pd.ID, &pd.TenantID, &pd.PossessionID, &pd.DocumentType, &pd.DocumentName, &pd.DocumentURL, &pd.FileName, &pd.FileSize, &pd.FileFormat, &pd.S3Bucket, &pd.S3Key, &pd.DocumentStatus, &pd.VerificationNotes, &pd.VerifiedBy, &pd.VerifiedAt, &pd.IsMandatory, &pd.UploadedBy, &pd.UploadedAt, &metadataJSON, &pd.CreatedBy, &pd.UpdatedBy, &pd.CreatedAt, &pd.UpdatedAt, &pd.DeletedAt)
	if err != nil {
		return nil, err
	}

	if metadataJSON != nil {
		var metadata models.JSONMap
		if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
			pd.Metadata = &metadata
		}
	}

	return pd, nil
}

// ListPossessionDocuments lists documents for a possession
func (s *PossessionService) ListPossessionDocuments(possessionID int64, limit, offset int) ([]*models.PossessionDocument, int, error) {
	var documents []*models.PossessionDocument
	countQuery := `SELECT COUNT(*) FROM possession_documents WHERE possession_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, possessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, possession_id, document_type, document_name, document_url, file_name, file_size, file_format, s3_bucket, s3_key, document_status, verification_notes, verified_by, verified_at, is_mandatory, uploaded_by, uploaded_at, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_documents WHERE possession_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, possessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		pd := &models.PossessionDocument{}
		var metadataJSON []byte
		err := rows.Scan(&pd.ID, &pd.TenantID, &pd.PossessionID, &pd.DocumentType, &pd.DocumentName, &pd.DocumentURL, &pd.FileName, &pd.FileSize, &pd.FileFormat, &pd.S3Bucket, &pd.S3Key, &pd.DocumentStatus, &pd.VerificationNotes, &pd.VerifiedBy, &pd.VerifiedAt, &pd.IsMandatory, &pd.UploadedBy, &pd.UploadedAt, &metadataJSON, &pd.CreatedBy, &pd.UpdatedBy, &pd.CreatedAt, &pd.UpdatedAt, &pd.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		if metadataJSON != nil {
			var metadata models.JSONMap
			if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
				pd.Metadata = &metadata
			}
		}
		documents = append(documents, pd)
	}

	return documents, total, nil
}

// UpdatePossessionDocument updates a possession document
func (s *PossessionService) UpdatePossessionDocument(id int64, updates map[string]interface{}, updatedBy *int64) (*models.PossessionDocument, error) {
	query := `UPDATE possession_documents SET `
	args := []interface{}{}
	count := 0

	for key, value := range updates {
		if count > 0 {
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		count++
	}

	query += ", updated_by = ? WHERE id = ?"
	args = append(args, updatedBy, id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetPossessionDocument(id)
}

// VerifyPossessionDocument verifies a possession document
func (s *PossessionService) VerifyPossessionDocument(id int64, status string, verificationNotes *string, verifiedBy *int64) (*models.PossessionDocument, error) {
	query := `UPDATE possession_documents SET document_status = ?, verification_notes = ?, verified_by = ?, verified_at = NOW() WHERE id = ?`
	_, err := s.db.Exec(query, status, verificationNotes, verifiedBy, id)
	if err != nil {
		return nil, err
	}
	return s.GetPossessionDocument(id)
}

// CreatePossessionRegistration creates a new registration
func (s *PossessionService) CreatePossessionRegistration(tenantID, possessionID int64, regType string, createdBy *int64) (*models.PossessionRegistration, error) {
	query := `
		INSERT INTO possession_registrations (tenant_id, possession_id, registration_type, registration_status, created_by)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, tenantID, possessionID, regType, "pending", createdBy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetPossessionRegistration(id)
}

// GetPossessionRegistration retrieves a registration by ID
func (s *PossessionService) GetPossessionRegistration(id int64) (*models.PossessionRegistration, error) {
	pr := &models.PossessionRegistration{}
	query := `SELECT id, tenant_id, possession_id, registration_type, registration_number, registration_office, registration_date, registration_status, amount_paid, amount_pending, payment_mode, reference_number, submission_date, expected_completion_date, actual_completion_date, remarks, approved_by, approved_at, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_registrations WHERE id = ? AND deleted_at IS NULL`
	err := s.db.QueryRow(query, id).Scan(&pr.ID, &pr.TenantID, &pr.PossessionID, &pr.RegistrationType, &pr.RegistrationNumber, &pr.RegistrationOffice, &pr.RegistrationDate, &pr.RegistrationStatus, &pr.AmountPaid, &pr.AmountPending, &pr.PaymentMode, &pr.ReferenceNumber, &pr.SubmissionDate, &pr.ExpectedCompletionDate, &pr.ActualCompletionDate, &pr.Remarks, &pr.ApprovedBy, &pr.ApprovedAt, &pr.CreatedBy, &pr.UpdatedBy, &pr.CreatedAt, &pr.UpdatedAt, &pr.DeletedAt)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

// ListPossessionRegistrations lists registrations for a possession
func (s *PossessionService) ListPossessionRegistrations(possessionID int64, limit, offset int) ([]*models.PossessionRegistration, int, error) {
	var registrations []*models.PossessionRegistration
	countQuery := `SELECT COUNT(*) FROM possession_registrations WHERE possession_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, possessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, possession_id, registration_type, registration_number, registration_office, registration_date, registration_status, amount_paid, amount_pending, payment_mode, reference_number, submission_date, expected_completion_date, actual_completion_date, remarks, approved_by, approved_at, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_registrations WHERE possession_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, possessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		pr := &models.PossessionRegistration{}
		err := rows.Scan(&pr.ID, &pr.TenantID, &pr.PossessionID, &pr.RegistrationType, &pr.RegistrationNumber, &pr.RegistrationOffice, &pr.RegistrationDate, &pr.RegistrationStatus, &pr.AmountPaid, &pr.AmountPending, &pr.PaymentMode, &pr.ReferenceNumber, &pr.SubmissionDate, &pr.ExpectedCompletionDate, &pr.ActualCompletionDate, &pr.Remarks, &pr.ApprovedBy, &pr.ApprovedAt, &pr.CreatedBy, &pr.UpdatedBy, &pr.CreatedAt, &pr.UpdatedAt, &pr.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		registrations = append(registrations, pr)
	}

	return registrations, total, nil
}

// UpdatePossessionRegistration updates a registration
func (s *PossessionService) UpdatePossessionRegistration(id int64, updates map[string]interface{}, updatedBy *int64) (*models.PossessionRegistration, error) {
	query := `UPDATE possession_registrations SET `
	args := []interface{}{}
	count := 0

	for key, value := range updates {
		if count > 0 {
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		count++
	}

	query += ", updated_by = ? WHERE id = ?"
	args = append(args, updatedBy, id)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return s.GetPossessionRegistration(id)
}

// ApproveRegistration approves a registration
func (s *PossessionService) ApproveRegistration(id int64, approvalStatus string, approvalNotes *string, approvedBy *int64) (*models.PossessionRegistration, error) {
	query := `UPDATE possession_registrations SET registration_status = ?, remarks = ?, approved_by = ?, approved_at = NOW() WHERE id = ?`
	_, err := s.db.Exec(query, approvalStatus, approvalNotes, approvedBy, id)
	if err != nil {
		return nil, err
	}
	return s.GetPossessionRegistration(id)
}

// CreatePossessionCertificate creates a new certificate
func (s *PossessionService) CreatePossessionCertificate(tenantID, possessionID int64, certType string, createdBy *int64) (*models.PossessionCertificate, error) {
	query := `
		INSERT INTO possession_certificates (tenant_id, possession_id, certificate_type, certificate_status, verification_status, created_by)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, tenantID, possessionID, certType, "pending", "pending", createdBy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetPossessionCertificate(id)
}

// GetPossessionCertificate retrieves a certificate by ID
func (s *PossessionService) GetPossessionCertificate(id int64) (*models.PossessionCertificate, error) {
	pc := &models.PossessionCertificate{}
	var metadataJSON []byte
	query := `SELECT id, tenant_id, possession_id, certificate_type, certificate_number, issuing_authority, issue_date, validity_date, certificate_url, file_name, file_size, s3_bucket, s3_key, certificate_status, verification_status, verified_by, verified_at, verification_notes, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_certificates WHERE id = ? AND deleted_at IS NULL`
	err := s.db.QueryRow(query, id).Scan(&pc.ID, &pc.TenantID, &pc.PossessionID, &pc.CertificateType, &pc.CertificateNumber, &pc.IssuingAuthority, &pc.IssueDate, &pc.ValidityDate, &pc.CertificateURL, &pc.FileName, &pc.FileSize, &pc.S3Bucket, &pc.S3Key, &pc.CertificateStatus, &pc.VerificationStatus, &pc.VerifiedBy, &pc.VerifiedAt, &pc.VerificationNotes, &metadataJSON, &pc.CreatedBy, &pc.UpdatedBy, &pc.CreatedAt, &pc.UpdatedAt, &pc.DeletedAt)
	if err != nil {
		return nil, err
	}

	if metadataJSON != nil {
		var metadata models.JSONMap
		if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
			pc.Metadata = &metadata
		}
	}

	return pc, nil
}

// ListPossessionCertificates lists certificates for a possession
func (s *PossessionService) ListPossessionCertificates(possessionID int64, limit, offset int) ([]*models.PossessionCertificate, int, error) {
	var certificates []*models.PossessionCertificate
	countQuery := `SELECT COUNT(*) FROM possession_certificates WHERE possession_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, possessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, possession_id, certificate_type, certificate_number, issuing_authority, issue_date, validity_date, certificate_url, file_name, file_size, s3_bucket, s3_key, certificate_status, verification_status, verified_by, verified_at, verification_notes, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_certificates WHERE possession_id = ? AND deleted_at IS NULL ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, possessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		pc := &models.PossessionCertificate{}
		var metadataJSON []byte
		err := rows.Scan(&pc.ID, &pc.TenantID, &pc.PossessionID, &pc.CertificateType, &pc.CertificateNumber, &pc.IssuingAuthority, &pc.IssueDate, &pc.ValidityDate, &pc.CertificateURL, &pc.FileName, &pc.FileSize, &pc.S3Bucket, &pc.S3Key, &pc.CertificateStatus, &pc.VerificationStatus, &pc.VerifiedBy, &pc.VerifiedAt, &pc.VerificationNotes, &metadataJSON, &pc.CreatedBy, &pc.UpdatedBy, &pc.CreatedAt, &pc.UpdatedAt, &pc.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		if metadataJSON != nil {
			var metadata models.JSONMap
			if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
				pc.Metadata = &metadata
			}
		}
		certificates = append(certificates, pc)
	}

	return certificates, total, nil
}

// VerifyCertificate verifies a certificate
func (s *PossessionService) VerifyCertificate(id int64, verificationStatus string, verifiedBy *int64) (*models.PossessionCertificate, error) {
	query := `UPDATE possession_certificates SET verification_status = ?, verified_by = ?, verified_at = NOW() WHERE id = ?`
	_, err := s.db.Exec(query, verificationStatus, verifiedBy, id)
	if err != nil {
		return nil, err
	}
	return s.GetPossessionCertificate(id)
}

// CreatePossessionApproval creates a new approval
func (s *PossessionService) CreatePossessionApproval(tenantID, possessionID, approverID int64, approvalType string, sequenceOrder int, createdBy *int64) (*models.PossessionApproval, error) {
	query := `
		INSERT INTO possession_approvals (tenant_id, possession_id, approval_type, approver_id, approval_status, sequence_order, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query, tenantID, possessionID, approvalType, approverID, "pending", sequenceOrder, createdBy)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetPossessionApproval(id)
}

// GetPossessionApproval retrieves an approval by ID
func (s *PossessionService) GetPossessionApproval(id int64) (*models.PossessionApproval, error) {
	pa := &models.PossessionApproval{}
	var metadataJSON []byte
	query := `SELECT id, tenant_id, possession_id, approval_type, approver_id, approver_role, approval_status, approval_notes, conditional_remarks, approval_date, valid_from, valid_till, sequence_order, is_final_approval, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_approvals WHERE id = ? AND deleted_at IS NULL`
	err := s.db.QueryRow(query, id).Scan(&pa.ID, &pa.TenantID, &pa.PossessionID, &pa.ApprovalType, &pa.ApproverID, &pa.ApproverRole, &pa.ApprovalStatus, &pa.ApprovalNotes, &pa.ConditionalRemarks, &pa.ApprovalDate, &pa.ValidFrom, &pa.ValidTill, &pa.SequenceOrder, &pa.IsFinalApproval, &metadataJSON, &pa.CreatedBy, &pa.UpdatedBy, &pa.CreatedAt, &pa.UpdatedAt, &pa.DeletedAt)
	if err != nil {
		return nil, err
	}

	if metadataJSON != nil {
		var metadata models.JSONMap
		if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
			pa.Metadata = &metadata
		}
	}

	return pa, nil
}

// ListPossessionApprovals lists approvals for a possession
func (s *PossessionService) ListPossessionApprovals(possessionID int64, limit, offset int) ([]*models.PossessionApproval, int, error) {
	var approvals []*models.PossessionApproval
	countQuery := `SELECT COUNT(*) FROM possession_approvals WHERE possession_id = ? AND deleted_at IS NULL`
	var total int
	err := s.db.QueryRow(countQuery, possessionID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, possession_id, approval_type, approver_id, approver_role, approval_status, approval_notes, conditional_remarks, approval_date, valid_from, valid_till, sequence_order, is_final_approval, metadata, created_by, updated_by, created_at, updated_at, deleted_at FROM possession_approvals WHERE possession_id = ? AND deleted_at IS NULL ORDER BY sequence_order ASC, created_at DESC LIMIT ? OFFSET ?`
	rows, err := s.db.Query(query, possessionID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		pa := &models.PossessionApproval{}
		var metadataJSON []byte
		err := rows.Scan(&pa.ID, &pa.TenantID, &pa.PossessionID, &pa.ApprovalType, &pa.ApproverID, &pa.ApproverRole, &pa.ApprovalStatus, &pa.ApprovalNotes, &pa.ConditionalRemarks, &pa.ApprovalDate, &pa.ValidFrom, &pa.ValidTill, &pa.SequenceOrder, &pa.IsFinalApproval, &metadataJSON, &pa.CreatedBy, &pa.UpdatedBy, &pa.CreatedAt, &pa.UpdatedAt, &pa.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		if metadataJSON != nil {
			var metadata models.JSONMap
			if err := json.Unmarshal(metadataJSON, &metadata); err == nil {
				pa.Metadata = &metadata
			}
		}
		approvals = append(approvals, pa)
	}

	return approvals, total, nil
}

// ApprovePossession approves a possession
func (s *PossessionService) ApprovePossession(id int64, approvalStatus string, approvalNotes *string, approvedBy *int64, isFinal bool) (*models.PossessionApproval, error) {
	query := `UPDATE possession_approvals SET approval_status = ?, approval_notes = ?, approval_date = NOW(), is_final_approval = ? WHERE id = ?`
	_, err := s.db.Exec(query, approvalStatus, approvalNotes, isFinal, id)
	if err != nil {
		return nil, err
	}
	return s.GetPossessionApproval(id)
}

// LogPossessionAction logs an audit action
func (s *PossessionService) LogPossessionAction(tenantID, possessionID int64, action, entityType string, entityID int64, oldValue, newValue, changeReason *string, performedBy *int64, userIP, userAgent *string) error {
	query := `
		INSERT INTO possession_audit_log (tenant_id, possession_id, action, entity_type, entity_id, old_value, new_value, change_reason, performed_by, user_ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, tenantID, possessionID, action, entityType, entityID, oldValue, newValue, changeReason, performedBy, userIP, userAgent)
	return err
}

// GetPossessionSummary retrieves a summary of possession status
func (s *PossessionService) GetPossessionSummary(possessionID int64) (*models.PossessionSummaryResponse, error) {
	ps, err := s.GetPossessionStatus(possessionID)
	if err != nil {
		return nil, err
	}

	documents, docTotal, _ := s.ListPossessionDocuments(possessionID, 1000, 0)
	_, regTotal, _ := s.ListPossessionRegistrations(possessionID, 1000, 0)
	_, certTotal, _ := s.ListPossessionCertificates(possessionID, 1000, 0)
	approvals, _, _ := s.ListPossessionApprovals(possessionID, 1000, 0)

	verifiedDocs := 0
	pendingDocs := 0
	for _, doc := range documents {
		switch doc.DocumentStatus {
		case "verified", "approved":
			verifiedDocs++
		case "pending":
			pendingDocs++
		}
	}

	pendingApprovals := 0
	for _, approval := range approvals {
		if approval.ApprovalStatus == "pending" {
			pendingApprovals++
		}
	}

	return &models.PossessionSummaryResponse{
		PossessionID:            ps.ID,
		BookingID:               ps.BookingID,
		Status:                  ps.Status,
		PossessionType:          ps.PossessionType,
		CompletionPercentage:    ps.CompletionPercentage,
		PossessionDate:          ps.PossessionDate,
		EstimatedPossessionDate: ps.EstimatedPossessionDate,
		TotalDocuments:          docTotal,
		VerifiedDocuments:       verifiedDocs,
		PendingDocuments:        pendingDocs,
		Registrations:           regTotal,
		Certificates:            certTotal,
		PendingApprovals:        pendingApprovals,
		CreatedAt:               ps.CreatedAt,
		UpdatedAt:               ps.UpdatedAt,
	}, nil
}
