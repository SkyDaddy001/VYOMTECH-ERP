package services

import (
	"fmt"
	"time"

	"vyomtech-backend/internal/models"

	"gorm.io/gorm"
)

// ============================================================================
// BankFinancingService
// ============================================================================
type BankFinancingService struct {
	db *gorm.DB
}

// NewBankFinancingService creates new bank financing service
func NewBankFinancingService(db *gorm.DB) *BankFinancingService {
	return &BankFinancingService{db: db}
}

// ============================================================================
// Financing Management
// ============================================================================

// CreateFinancing creates new financing record
func (s *BankFinancingService) CreateFinancing(tenantID int64, req *models.CreateFinancingRequest, userID int64) (*models.BankFinancing, error) {
	financing := &models.BankFinancing{
		TenantID:         tenantID,
		BookingID:        req.BookingID,
		BankID:           req.BankID,
		LoanAmount:       req.LoanAmount,
		SanctionedAmount: req.SanctionedAmount,
		LoanType:         req.LoanType,
		InterestRate:     req.InterestRate,
		TenureMonths:     req.TenureMonths,
		ApplicationRefNo: req.ApplicationRefNo,
		Status:           "pending",
		CreatedBy:        &userID,
		CreatedAt:        time.Now(),
	}

	// Calculate EMI if tenure and interest rate provided
	if req.TenureMonths != nil && req.InterestRate != nil && *req.TenureMonths > 0 {
		emi := s.calculateEMI(req.LoanAmount, *req.InterestRate, *req.TenureMonths)
		financing.EMIAmount = &emi
	}

	if err := s.db.Create(financing).Error; err != nil {
		return nil, fmt.Errorf("failed to create financing: %w", err)
	}

	return financing, nil
}

// GetFinancingByID retrieves financing by ID
func (s *BankFinancingService) GetFinancingByID(tenantID, financingID int64) (*models.BankFinancing, error) {
	var financing models.BankFinancing
	if err := s.db.
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, financingID).
		Preload("Bank").
		Preload("Disbursements").
		Preload("NOCs").
		First(&financing).Error; err != nil {
		return nil, fmt.Errorf("financing not found: %w", err)
	}
	return &financing, nil
}

// GetFinancingByBookingID retrieves financing by booking ID
func (s *BankFinancingService) GetFinancingByBookingID(tenantID, bookingID int64) (*models.BankFinancing, error) {
	var financing models.BankFinancing
	if err := s.db.
		Where("tenant_id = ? AND booking_id = ? AND deleted_at IS NULL", tenantID, bookingID).
		Preload("Bank").
		First(&financing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get financing: %w", err)
	}
	return &financing, nil
}

// ListFinancing lists all financing records
func (s *BankFinancingService) ListFinancing(tenantID int64, filters map[string]interface{}, limit, offset int) ([]models.BankFinancing, int64, error) {
	var financings []models.BankFinancing
	var total int64

	query := s.db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	// Apply filters
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if loanType, ok := filters["loan_type"]; ok {
		query = query.Where("loan_type = ?", loanType)
	}
	if bankID, ok := filters["bank_id"]; ok {
		query = query.Where("bank_id = ?", bankID)
	}

	if err := query.Model(&models.BankFinancing{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Bank").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&financings).Error; err != nil {
		return nil, 0, err
	}

	return financings, total, nil
}

// UpdateFinancing updates financing record
func (s *BankFinancingService) UpdateFinancing(tenantID, financingID int64, req *models.UpdateFinancingRequest, userID int64) (*models.BankFinancing, error) {
	financing, err := s.GetFinancingByID(tenantID, financingID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"updated_by": userID,
		"updated_at": time.Now(),
	}

	if req.SanctionedAmount != nil {
		updates["sanctioned_amount"] = *req.SanctionedAmount
	}
	if req.LoanType != nil {
		updates["loan_type"] = *req.LoanType
	}
	if req.InterestRate != nil {
		updates["interest_rate"] = *req.InterestRate
	}
	if req.TenureMonths != nil {
		updates["tenure_months"] = *req.TenureMonths
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.ApprovalDate != nil {
		updates["approval_date"] = *req.ApprovalDate
	}
	if req.SanctionDate != nil {
		updates["sanction_date"] = *req.SanctionDate
	}

	if err := s.db.Model(financing).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update financing: %w", err)
	}

	return s.GetFinancingByID(tenantID, financingID)
}

// DeleteFinancing soft deletes financing record
func (s *BankFinancingService) DeleteFinancing(tenantID, financingID int64, userID int64) error {
	financing, err := s.GetFinancingByID(tenantID, financingID)
	if err != nil {
		return err
	}

	if err := s.db.Model(financing).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"updated_by": userID,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("failed to delete financing: %w", err)
	}

	return nil
}

// ============================================================================
// Disbursement Management
// ============================================================================

// CreateDisbursement creates disbursement schedule
func (s *BankFinancingService) CreateDisbursement(tenantID int64, req *models.CreateDisbursementRequest, userID int64) (*models.BankDisbursement, error) {
	// Verify financing exists
	if _, err := s.GetFinancingByID(tenantID, req.FinancingID); err != nil {
		return nil, err
	}

	// Get next disbursement number
	var maxNumber int
	s.db.
		Where("tenant_id = ? AND financing_id = ? AND deleted_at IS NULL", tenantID, req.FinancingID).
		Model(&models.BankDisbursement{}).
		Select("COALESCE(MAX(disbursement_number), 0)").
		Scan(&maxNumber)

	disbursement := &models.BankDisbursement{
		TenantID:            tenantID,
		FinancingID:         req.FinancingID,
		DisbursementNumber:  maxNumber + 1,
		ScheduledAmount:     req.ScheduledAmount,
		MilestoneID:         req.MilestoneID,
		MilestonePercentage: req.MilestonePercentage,
		ScheduledDate:       req.ScheduledDate,
		Status:              "pending",
		CreatedBy:           &userID,
		CreatedAt:           time.Now(),
	}

	if err := s.db.Create(disbursement).Error; err != nil {
		return nil, fmt.Errorf("failed to create disbursement: %w", err)
	}

	return disbursement, nil
}

// GetDisbursement retrieves disbursement by ID
func (s *BankFinancingService) GetDisbursement(tenantID, disbursementID int64) (*models.BankDisbursement, error) {
	var disbursement models.BankDisbursement
	if err := s.db.
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, disbursementID).
		First(&disbursement).Error; err != nil {
		return nil, fmt.Errorf("disbursement not found: %w", err)
	}
	return &disbursement, nil
}

// ListDisbursements lists disbursements for financing
func (s *BankFinancingService) ListDisbursements(tenantID, financingID int64) ([]models.BankDisbursement, error) {
	var disbursements []models.BankDisbursement
	if err := s.db.
		Where("tenant_id = ? AND financing_id = ? AND deleted_at IS NULL", tenantID, financingID).
		Order("disbursement_number ASC").
		Find(&disbursements).Error; err != nil {
		return nil, err
	}
	return disbursements, nil
}

// UpdateDisbursementStatus updates disbursement status
func (s *BankFinancingService) UpdateDisbursementStatus(tenantID, disbursementID int64, status string, actualAmount *float64, userID int64) (*models.BankDisbursement, error) {
	updates := map[string]interface{}{
		"status":     status,
		"updated_by": userID,
		"updated_at": time.Now(),
	}

	if actualAmount != nil {
		updates["actual_amount"] = *actualAmount
	}

	if status == "credited" {
		updates["actual_date"] = time.Now()
	}

	disbursement, err := s.GetDisbursement(tenantID, disbursementID)
	if err != nil {
		return nil, err
	}

	if err := s.db.Model(disbursement).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetDisbursement(tenantID, disbursementID)
}

// ============================================================================
// NOC Management
// ============================================================================

// CreateNOC creates NOC record
func (s *BankFinancingService) CreateNOC(tenantID int64, req *models.CreateNOCRequest, userID int64) (*models.BankNOC, error) {
	if _, err := s.GetFinancingByID(tenantID, req.FinancingID); err != nil {
		return nil, err
	}

	noc := &models.BankNOC{
		TenantID:       tenantID,
		FinancingID:    req.FinancingID,
		NOCType:        req.NOCType,
		NOCRequestDate: req.NOCRequestDate,
		NOCAmount:      req.NOCAmount,
		Status:         "requested",
		CreatedBy:      &userID,
		CreatedAt:      time.Now(),
	}

	if err := s.db.Create(noc).Error; err != nil {
		return nil, fmt.Errorf("failed to create NOC: %w", err)
	}

	return noc, nil
}

// GetNOC retrieves NOC by ID
func (s *BankFinancingService) GetNOC(tenantID, nocID int64) (*models.BankNOC, error) {
	var noc models.BankNOC
	if err := s.db.
		Where("tenant_id = ? AND id = ? AND deleted_at IS NULL", tenantID, nocID).
		First(&noc).Error; err != nil {
		return nil, fmt.Errorf("NOC not found: %w", err)
	}
	return &noc, nil
}

// ListNOCs lists NOCs for financing
func (s *BankFinancingService) ListNOCs(tenantID, financingID int64) ([]models.BankNOC, error) {
	var nocs []models.BankNOC
	if err := s.db.
		Where("tenant_id = ? AND financing_id = ? AND deleted_at IS NULL", tenantID, financingID).
		Order("created_at DESC").
		Find(&nocs).Error; err != nil {
		return nil, err
	}
	return nocs, nil
}

// ============================================================================
// Collection Management
// ============================================================================

// CreateCollection creates collection record
func (s *BankFinancingService) CreateCollection(tenantID int64, req *models.CreateCollectionRequest, userID int64) (*models.BankCollectionTracking, error) {
	if _, err := s.GetFinancingByID(tenantID, req.FinancingID); err != nil {
		return nil, err
	}

	collection := &models.BankCollectionTracking{
		TenantID:         tenantID,
		FinancingID:      req.FinancingID,
		CollectionType:   req.CollectionType,
		CollectionAmount: req.CollectionAmount,
		CollectionDate:   req.CollectionDate,
		PaymentMode:      req.PaymentMode,
		EMIMonth:         req.EMIMonth,
		EMINumber:        req.EMINumber,
		Status:           "pending",
		CreatedBy:        &userID,
		CreatedAt:        time.Now(),
	}

	if err := s.db.Create(collection).Error; err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return collection, nil
}

// ListCollections lists collections for financing
func (s *BankFinancingService) ListCollections(tenantID, financingID int64) ([]models.BankCollectionTracking, error) {
	var collections []models.BankCollectionTracking
	if err := s.db.
		Where("tenant_id = ? AND financing_id = ? AND deleted_at IS NULL", tenantID, financingID).
		Order("collection_date DESC").
		Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// calculateEMI calculates EMI (Equated Monthly Installment)
// Formula: EMI = P * r * (1 + r)^n / ((1 + r)^n - 1)
// where P = principal, r = monthly rate, n = tenure in months
func (s *BankFinancingService) calculateEMI(principal float64, annualRate float64, tenureMonths int) float64 {
	monthlyRate := annualRate / 12 / 100
	if monthlyRate == 0 {
		return principal / float64(tenureMonths)
	}

	numerator := principal * monthlyRate * func() float64 {
		base := 1 + monthlyRate
		result := 1.0
		for i := 0; i < tenureMonths; i++ {
			result *= base
		}
		return result
	}()

	denominator := func() float64 {
		base := 1 + monthlyRate
		result := 1.0
		for i := 0; i < tenureMonths; i++ {
			result *= base
		}
		return result - 1
	}()

	return numerator / denominator
}

// GetFinancingSummary gets comprehensive financing summary
func (s *BankFinancingService) GetFinancingSummary(tenantID, financingID int64) (map[string]interface{}, error) {
	financing, err := s.GetFinancingByID(tenantID, financingID)
	if err != nil {
		return nil, err
	}

	// Get total disbursed
	var totalDisbursed float64
	s.db.
		Where("tenant_id = ? AND financing_id = ? AND status != 'cancelled' AND deleted_at IS NULL", tenantID, financingID).
		Model(&models.BankDisbursement{}).
		Select("COALESCE(SUM(actual_amount), 0)").
		Scan(&totalDisbursed)

	// Get total collected
	var totalCollected float64
	s.db.
		Where("tenant_id = ? AND financing_id = ? AND status = 'credited' AND deleted_at IS NULL", tenantID, financingID).
		Model(&models.BankCollectionTracking{}).
		Select("COALESCE(SUM(collection_amount), 0)").
		Scan(&totalCollected)

	// Get pending disbursements
	var pendingDisbursements int64
	s.db.
		Where("tenant_id = ? AND financing_id = ? AND status = 'pending' AND deleted_at IS NULL", tenantID, financingID).
		Model(&models.BankDisbursement{}).
		Count(&pendingDisbursements)

	return map[string]interface{}{
		"financing_id":            financing.ID,
		"loan_amount":             financing.LoanAmount,
		"sanctioned_amount":       financing.SanctionedAmount,
		"total_disbursed":         totalDisbursed,
		"total_collected":         totalCollected,
		"outstanding_amount":      financing.OutstandingAmount,
		"disbursement_percentage": (totalDisbursed / financing.SanctionedAmount) * 100,
		"pending_disbursements":   pendingDisbursements,
		"status":                  financing.Status,
	}, nil
}
