package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// ============================================================================
// BankFinancingHandler
// ============================================================================
type BankFinancingHandler struct {
	financingService *services.BankFinancingService
}

// NewBankFinancingHandler creates new bank financing handler
func NewBankFinancingHandler(financingService *services.BankFinancingService) *BankFinancingHandler {
	return &BankFinancingHandler{
		financingService: financingService,
	}
}

// ============================================================================
// BANK FINANCING ENDPOINTS
// ============================================================================

// CreateBankFinancing POST /api/v1/financing
func (h *BankFinancingHandler) CreateBankFinancing(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req struct {
		BookingID              string  `json:"booking_id" binding:"required"`
		BankID                 string  `json:"bank_id" binding:"required"`
		LoanAmount             float64 `json:"loan_amount" binding:"required"`
		SanctionedAmount       float64 `json:"sanctioned_amount"`
		LoanType               string  `json:"loan_type"`
		InterestRate           float64 `json:"interest_rate"`
		TenureMonths           int     `json:"tenure_months"`
		ApplicationDate        string  `json:"application_date"`
		SanctionDate           string  `json:"sanction_date"`
		ApplicationRefNo       string  `json:"application_ref_no"`
		SanctionLetterURL      string  `json:"sanction_letter_url"`
		ExpectedCompletionDate string  `json:"expected_completion_date"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Convert pointers for nullable fields
	var bankID *string
	if req.BankID != "" {
		bankID = &req.BankID
	}

	var interestRate *float64
	if req.InterestRate > 0 {
		interestRate = &req.InterestRate
	}

	var tenureMonths *int
	if req.TenureMonths > 0 {
		tenureMonths = &req.TenureMonths
	}

	var appRefNo *string
	if req.ApplicationRefNo != "" {
		appRefNo = &req.ApplicationRefNo
	}

	var sanctionLetterURL *string
	if req.SanctionLetterURL != "" {
		sanctionLetterURL = &req.SanctionLetterURL
	}

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	financing := &models.BankFinancing{
		TenantID:               tenantID,
		BookingID:              req.BookingID,
		BankID:                 bankID,
		LoanAmount:             req.LoanAmount,
		SanctionedAmount:       req.SanctionedAmount,
		LoanType:               req.LoanType,
		InterestRate:           interestRate,
		TenureMonths:           tenureMonths,
		ApplicationRefNo:       appRefNo,
		SanctionLetterURL:      sanctionLetterURL,
		Status:                 "draft",
		CreatedBy:              userIDPtr,
		UpdatedBy:              userIDPtr,
	}

	created, err := h.financingService.CreateBankFinancing(c.Request.Context(), financing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create financing"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    created,
		"message": "Financing created successfully",
	})
}

// GetBankFinancing GET /api/v1/financing/:id
func (h *BankFinancingHandler) GetBankFinancing(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	financingID := c.Param("id")

	financing, err := h.financingService.GetBankFinancing(c.Request.Context(), tenantID, financingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch financing"})
		return
	}

	if financing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Financing not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    financing,
	})
}

// ListBankFinancing GET /api/v1/financing
func (h *BankFinancingHandler) ListBankFinancing(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	financings, err := h.financingService.ListBankFinancing(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list financing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    financings,
		"count":   len(financings),
	})
}

// ============================================================================
// BANK DISBURSEMENT ENDPOINTS
// ============================================================================

// CreateBankDisbursement POST /api/v1/financing/:id/disbursement
func (h *BankFinancingHandler) CreateBankDisbursement(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	financingID := c.Param("id")

	var req struct {
		DisbursementNumber  int     `json:"disbursement_number" binding:"required"`
		ScheduledAmount     float64 `json:"scheduled_amount" binding:"required"`
		MilestoneID         string  `json:"milestone_id"`
		MilestonePercentage int     `json:"milestone_percentage"`
		ScheduledDate       string  `json:"scheduled_date"`
		ClaimDocumentURL    string  `json:"claim_document_url"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var milestoneID *string
	if req.MilestoneID != "" {
		milestoneID = &req.MilestoneID
	}

	var milestonePercentage *int
	if req.MilestonePercentage > 0 {
		milestonePercentage = &req.MilestonePercentage
	}

	var claimDocURL *string
	if req.ClaimDocumentURL != "" {
		claimDocURL = &req.ClaimDocumentURL
	}

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	disbursement := &models.BankDisbursement{
		TenantID:             tenantID,
		FinancingID:          financingID,
		DisbursementNumber:   req.DisbursementNumber,
		ScheduledAmount:      req.ScheduledAmount,
		MilestoneID:          milestoneID,
		MilestonePercentage:  milestonePercentage,
		ClaimDocumentURL:     claimDocURL,
		Status:               "pending",
		CreatedBy:            userIDPtr,
		UpdatedBy:            userIDPtr,
	}

	created, err := h.financingService.CreateBankDisbursement(c.Request.Context(), disbursement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create disbursement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    created,
		"message": "Disbursement created successfully",
	})
}

// ============================================================================
// BANK NOC ENDPOINTS
// ============================================================================

// CreateBankNOC POST /api/v1/financing/:id/noc
func (h *BankFinancingHandler) CreateBankNOC(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	financingID := c.Param("id")

	var req struct {
		NOCType        string  `json:"noc_type" binding:"required"`
		NOCRequestDate string  `json:"noc_request_date"`
		NOCAmount      float64 `json:"noc_amount"`
		ValidTillDate  string  `json:"valid_till_date"`
		Remarks        string  `json:"remarks"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var nocAmount *float64
	if req.NOCAmount > 0 {
		nocAmount = &req.NOCAmount
	}

	var remarks *string
	if req.Remarks != "" {
		remarks = &req.Remarks
	}

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	noc := &models.BankNOC{
		TenantID:       tenantID,
		FinancingID:    financingID,
		NOCType:        req.NOCType,
		NOCAmount:      nocAmount,
		Remarks:        remarks,
		Status:         "requested",
		CreatedBy:      userIDPtr,
		UpdatedBy:      userIDPtr,
	}

	created, err := h.financingService.CreateBankNOC(c.Request.Context(), noc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create NOC"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    created,
		"message": "NOC created successfully",
	})
}

// ============================================================================
// BANK COLLECTION ENDPOINTS
// ============================================================================

// CreateBankCollection POST /api/v1/financing/:id/collection
func (h *BankFinancingHandler) CreateBankCollection(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	financingID := c.Param("id")

	var req struct {
		CollectionType     string  `json:"collection_type" binding:"required"`
		CollectionAmount   float64 `json:"collection_amount" binding:"required"`
		CollectionDate     string  `json:"collection_date"`
		PaymentMode        string  `json:"payment_mode"`
		PaymentReferenceNo string  `json:"payment_reference_no"`
		EMIMonth           string  `json:"emi_month"`
		EMINumber          int     `json:"emi_number"`
		PrincipalAmount    float64 `json:"principal_amount"`
		InterestAmount     float64 `json:"interest_amount"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var paymentMode *string
	if req.PaymentMode != "" {
		paymentMode = &req.PaymentMode
	}

	var paymentRefNo *string
	if req.PaymentReferenceNo != "" {
		paymentRefNo = &req.PaymentReferenceNo
	}

	var emiMonth *string
	if req.EMIMonth != "" {
		emiMonth = &req.EMIMonth
	}

	var emiNumber *int
	if req.EMINumber > 0 {
		emiNumber = &req.EMINumber
	}

	var principalAmount *float64
	if req.PrincipalAmount > 0 {
		principalAmount = &req.PrincipalAmount
	}

	var interestAmount *float64
	if req.InterestAmount > 0 {
		interestAmount = &req.InterestAmount
	}

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	collection := &models.BankCollectionTracking{
		TenantID:            tenantID,
		FinancingID:         financingID,
		CollectionType:      req.CollectionType,
		CollectionAmount:    req.CollectionAmount,
		PaymentMode:         paymentMode,
		PaymentReferenceNo:  paymentRefNo,
		EMIMonth:            emiMonth,
		EMINumber:           emiNumber,
		PrincipalAmount:     principalAmount,
		InterestAmount:      interestAmount,
		Status:              "recorded",
		CreatedBy:           userIDPtr,
		UpdatedBy:           userIDPtr,
	}

	created, err := h.financingService.CreateBankCollection(c.Request.Context(), collection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record collection"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    created,
		"message": "Collection recorded successfully",
	})
}

// ============================================================================
// BANK MASTER ENDPOINTS
// ============================================================================

// CreateBank POST /api/v1/banks
func (h *BankFinancingHandler) CreateBank(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	var req struct {
		BankName                 string `json:"bank_name" binding:"required"`
		BranchName               string `json:"branch_name" binding:"required"`
		IFSCCode                 string `json:"ifsc_code" binding:"required"`
		BranchContact            string `json:"branch_contact"`
		BranchEmail              string `json:"branch_email"`
		RelationshipManagerName  string `json:"relationship_manager_name"`
		RelationshipManagerPhone string `json:"relationship_manager_phone"`
		RelationshipManagerEmail string `json:"relationship_manager_email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var ifscCode *string
	if req.IFSCCode != "" {
		ifscCode = &req.IFSCCode
	}

	var branchContact *string
	if req.BranchContact != "" {
		branchContact = &req.BranchContact
	}

	var branchEmail *string
	if req.BranchEmail != "" {
		branchEmail = &req.BranchEmail
	}

	var rmName *string
	if req.RelationshipManagerName != "" {
		rmName = &req.RelationshipManagerName
	}

	var rmPhone *string
	if req.RelationshipManagerPhone != "" {
		rmPhone = &req.RelationshipManagerPhone
	}

	var rmEmail *string
	if req.RelationshipManagerEmail != "" {
		rmEmail = &req.RelationshipManagerEmail
	}

	bank := &models.Bank{
		TenantID:                   tenantID,
		BankName:                   req.BankName,
		BranchName:                 req.BranchName,
		IFSCCode:                   ifscCode,
		BranchContact:              branchContact,
		BranchEmail:                branchEmail,
		RelationshipManagerName:    rmName,
		RelationshipManagerPhone:   rmPhone,
		RelationshipManagerEmail:   rmEmail,
		Status:                     "active",
	}

	created, err := h.financingService.CreateBank(c.Request.Context(), bank)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bank"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    created,
		"message": "Bank created successfully",
	})
}
		RelationshipManagerEmail: req.RelationshipManagerEmail,
		Status:                   "active",
	}

	created, err := h.financingService.CreateBank(c.Request.Context(), bank)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bank"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    created,
		"message": "Bank created successfully",
	})
}

// ListBanks GET /api/v1/banks
func (h *BankFinancingHandler) ListBanks(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	banks, err := h.financingService.ListBanks(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list banks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    banks,
		"count":   len(banks),
	})
}
