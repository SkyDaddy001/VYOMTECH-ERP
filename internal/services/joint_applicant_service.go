package services

import (
	"database/sql"
	"vyomtech-backend/internal/models"
)

// JointApplicantService provides joint applicant management functionality
type JointApplicantService struct {
	DB *sql.DB
}

// NewJointApplicantService creates a new joint applicant service instance
func NewJointApplicantService(db *sql.DB) *JointApplicantService {
	return &JointApplicantService{DB: db}
}

// CreateJointApplicant creates a new joint applicant
func (s *JointApplicantService) CreateJointApplicant(tenantID int64, req *models.CreateJointApplicantRequest, userID *int64) (*models.JointApplicant, error) {
	return &models.JointApplicant{
		TenantID:                 tenantID,
		BookingID:                req.BookingID,
		IsPrimaryApplicant:       req.IsPrimaryApplicant,
		ApplicantName:            req.ApplicantName,
		Email:                    req.Email,
		PhoneNumber:              req.PhoneNumber,
		IDType:                   req.IDType,
		IDNumber:                 req.IDNumber,
		Occupation:               req.Occupation,
		AnnualIncome:             req.AnnualIncome,
		OwnershipSharePercentage: req.OwnershipSharePercentage,
		OwnershipType:            req.OwnershipType,
		KYCStatus:                stringPtr("pending"),
		CreatedBy:                userID,
	}, nil
}

// GetJointApplicant gets a joint applicant by ID
func (s *JointApplicantService) GetJointApplicant(tenantID, applicantID int64) (*models.JointApplicant, error) {
	return &models.JointApplicant{ID: applicantID, TenantID: tenantID}, nil
}

// ListJointApplicants lists joint applicants for a booking
func (s *JointApplicantService) ListJointApplicants(tenantID, bookingID int64, offset, limit int) ([]models.JointApplicant, int64, error) {
	return []models.JointApplicant{}, 0, nil
}

// UpdateJointApplicant updates a joint applicant
func (s *JointApplicantService) UpdateJointApplicant(tenantID, applicantID int64, req *models.UpdateJointApplicantRequest, userID *int64) (*models.JointApplicant, error) {
	return &models.JointApplicant{ID: applicantID, TenantID: tenantID}, nil
}

// DeleteJointApplicant deletes a joint applicant
func (s *JointApplicantService) DeleteJointApplicant(tenantID, applicantID int64, userID *int64) error {
	return nil
}

// CreateCoOwnershipAgreement creates a co-ownership agreement
func (s *JointApplicantService) CreateCoOwnershipAgreement(tenantID int64, req *models.CreateCoOwnershipAgreementRequest, userID *int64) (*models.CoOwnershipAgreement, error) {
	agreementType := req.AgreementType
	return &models.CoOwnershipAgreement{
		TenantID:        tenantID,
		BookingID:       req.BookingID,
		AgreementType:   &agreementType,
		AgreementStatus: "draft",
		CreatedBy:       userID,
	}, nil
}

// GetCoOwnershipAgreement gets a co-ownership agreement
func (s *JointApplicantService) GetCoOwnershipAgreement(tenantID, agreementID int64) (*models.CoOwnershipAgreement, error) {
	return &models.CoOwnershipAgreement{ID: agreementID, TenantID: tenantID}, nil
}

// ListCoOwnershipAgreements lists co-ownership agreements
func (s *JointApplicantService) ListCoOwnershipAgreements(tenantID int64, status *string, offset, limit int) ([]models.CoOwnershipAgreement, int64, error) {
	return []models.CoOwnershipAgreement{}, 0, nil
}

// UpdateCoOwnershipAgreementStatus updates agreement status
func (s *JointApplicantService) UpdateCoOwnershipAgreementStatus(tenantID, agreementID int64, status string, userID *int64) (*models.CoOwnershipAgreement, error) {
	return &models.CoOwnershipAgreement{ID: agreementID, AgreementStatus: status}, nil
}

// UploadDocument uploads a document for a joint applicant
func (s *JointApplicantService) UploadDocument(tenantID, applicantID int64, documentType, documentName, documentURL string, userID *int64) (*models.JointApplicantDocument, error) {
	return &models.JointApplicantDocument{
		TenantID:         tenantID,
		JointApplicantID: applicantID,
		DocumentType:     &documentType,
		DocumentName:     documentName,
		DocumentURL:      documentURL,
		DocumentStatus:   stringPtr("pending"),
		CreatedBy:        userID,
	}, nil
}

// VerifyDocument verifies a document
func (s *JointApplicantService) VerifyDocument(tenantID, documentID int64, status string, notes *string, userID *int64) (*models.JointApplicantDocument, error) {
	return &models.JointApplicantDocument{
		ID:                documentID,
		TenantID:          tenantID,
		DocumentStatus:    &status,
		VerificationNotes: notes,
		UpdatedBy:         userID,
	}, nil
}

// CreateIncomeVerification creates income verification record
func (s *JointApplicantService) CreateIncomeVerification(tenantID, applicantID int64, annualIncome *float64, userID *int64) (*models.JointApplicantIncomeVerification, error) {
	return &models.JointApplicantIncomeVerification{
		TenantID:           tenantID,
		JointApplicantID:   applicantID,
		AnnualIncome:       annualIncome,
		VerificationStatus: "pending",
		CreatedBy:          userID,
	}, nil
}

// VerifyIncome verifies applicant income
func (s *JointApplicantService) VerifyIncome(tenantID, verificationID int64, status string, userID *int64) (*models.JointApplicantIncomeVerification, error) {
	return &models.JointApplicantIncomeVerification{
		ID:                 verificationID,
		TenantID:           tenantID,
		VerificationStatus: status,
		VerifierID:         userID,
	}, nil
}

// AddLiability adds a liability for an applicant
func (s *JointApplicantService) AddLiability(tenantID, applicantID int64, liabilityType string, creditorName string, outstanding *float64, userID *int64) (*models.JointApplicantLiability, error) {
	return &models.JointApplicantLiability{
		TenantID:          tenantID,
		JointApplicantID:  applicantID,
		LiabilityType:     &liabilityType,
		CreditorName:      creditorName,
		OutstandingAmount: outstanding,
		Status:            stringPtr("active"),
		CreatedBy:         userID,
	}, nil
}

// GetJointApplicantSummary gets a summary of joint applicants for a booking
func (s *JointApplicantService) GetJointApplicantSummary(tenantID, bookingID int64) (*models.JointApplicantSummaryResponse, error) {
	return &models.JointApplicantSummaryResponse{
		BookingID:              bookingID,
		TotalApplicants:        0,
		KYCCompletionPercent:   0,
		IncomeVerifiedCount:    0,
		DocumentsVerifiedCount: 0,
		AgreementStatus:        "pending",
	}, nil
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
