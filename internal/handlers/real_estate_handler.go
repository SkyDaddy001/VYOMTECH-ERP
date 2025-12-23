package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gorilla/mux"
)

// RealEstateHandler handles all real estate related operations
type RealEstateHandler struct {
	DB          *sql.DB
	RBACService *services.RBACService
}

// NewRealEstateHandler creates a new real estate handler
func NewRealEstateHandler(db *sql.DB, rbacService *services.RBACService) *RealEstateHandler {
	return &RealEstateHandler{
		DB:          db,
		RBACService: rbacService,
	}
}

// ============================================
// PROJECT ENDPOINTS
// ============================================

// CreateProject creates a new property project
func (h *RealEstateHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenantID == "" {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found in context")
		return
	}

	// For real estate, create is combined with general property operations
	// Using a generic property.create permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenantID, userID, "properties.create"); err != nil {
		h.respondError(w, http.StatusForbidden, fmt.Sprintf("Permission denied: %s", err.Error()))
		return
	}

	var req models.CreatePropertyProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project := &models.PropertyProject{
		TenantID:           tenantID,
		ProjectName:        req.ProjectName,
		ProjectCode:        req.ProjectCode,
		Location:           req.Location,
		City:               req.City,
		State:              req.State,
		PostalCode:         req.PostalCode,
		TotalUnits:         req.TotalUnits,
		TotalArea:          req.TotalArea,
		ProjectType:        req.ProjectType,
		Status:             req.Status,
		LaunchDate:         req.LaunchDate,
		ExpectedCompletion: req.ExpectedCompletion,
		DeveloperName:      req.DeveloperName,
		ArchitectName:      req.ArchitectName,
	}

	query := `INSERT INTO property_projects 
		(tenant_id, project_name, project_code, location, city, state, postal_code, 
		 total_units, total_area, project_type, status, launch_date, expected_completion,
		 developer_name, architect_name) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at`

	err := h.DB.QueryRow(query,
		project.TenantID, project.ProjectName, project.ProjectCode, project.Location,
		project.City, project.State, project.PostalCode, project.TotalUnits,
		project.TotalArea, project.ProjectType, project.Status, project.LaunchDate,
		project.ExpectedCompletion, project.DeveloperName, project.ArchitectName,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	h.respondJSON(w, http.StatusCreated, project)
}

// GetProjects retrieves all projects for a tenant
func (h *RealEstateHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	query := `SELECT id, tenant_id, project_name, project_code, location, city, state, 
		postal_code, total_units, total_area, project_type, status, launch_date, 
		expected_completion, actual_completion, noc_status, noc_date, developer_name, 
		architect_name, created_at, updated_at, deleted_at, created_by
		FROM property_projects WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := h.DB.Query(query, tenantID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}
	defer rows.Close()

	projects := []models.PropertyProject{}
	for rows.Next() {
		var p models.PropertyProject
		if err := rows.Scan(&p.ID, &p.TenantID, &p.ProjectName, &p.ProjectCode, &p.Location,
			&p.City, &p.State, &p.PostalCode, &p.TotalUnits, &p.TotalArea, &p.ProjectType,
			&p.Status, &p.LaunchDate, &p.ExpectedCompletion, &p.ActualCompletion,
			&p.NOCStatus, &p.NOCDate, &p.DeveloperName, &p.ArchitectName,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy); err != nil {
			log.Printf("Error scanning project: %v", err)
			continue
		}
		projects = append(projects, p)
	}

	h.respondJSON(w, http.StatusOK, projects)
}

// ============================================
// PROPERTY UNIT ENDPOINTS
// ============================================

// CreateUnit creates a new property unit
func (h *RealEstateHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreatePropertyUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	unit := &models.PropertyUnit{
		TenantID:              tenantID,
		ProjectID:             req.ProjectID,
		BlockID:               req.BlockID,
		UnitNumber:            req.UnitNumber,
		Floor:                 req.Floor,
		UnitType:              req.UnitType,
		Facing:                req.Facing,
		CarpetArea:            req.CarpetArea,
		CarpetAreaWithBalcony: req.CarpetAreaWithBalcony,
		UtilityArea:           req.UtilityArea,
		PlinthArea:            req.PlinthArea,
		SBUA:                  req.SBUA,
		UDSSqft:               req.UDSSqft,
		Status:                "available",
	}

	query := `INSERT INTO property_units 
		(tenant_id, project_id, block_id, unit_number, floor, unit_type, facing,
		 carpet_area, carpet_area_with_balcony, utility_area, plinth_area, sbua, uds_sqft, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	err := h.DB.QueryRow(query,
		unit.TenantID, unit.ProjectID, unit.BlockID, unit.UnitNumber, unit.Floor,
		unit.UnitType, unit.Facing, unit.CarpetArea, unit.CarpetAreaWithBalcony,
		unit.UtilityArea, unit.PlinthArea, unit.SBUA, unit.UDSSqft, unit.Status,
	).Scan(&unit.ID, &unit.CreatedAt, &unit.UpdatedAt)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create unit")
		return
	}

	h.respondJSON(w, http.StatusCreated, unit)
}

// ListUnits retrieves all units for a project
func (h *RealEstateHandler) ListUnits(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	projectID := mux.Vars(r)["project_id"]

	query := `SELECT id, tenant_id, project_id, block_id, unit_number, floor, unit_type, facing,
		carpet_area, carpet_area_with_balcony, utility_area, plinth_area, sbua, uds_sqft, status,
		alloted_to, allotment_date, created_at, updated_at, deleted_at
		FROM property_units WHERE tenant_id = $1 AND project_id = $2 AND deleted_at IS NULL
		ORDER BY unit_number`

	rows, err := h.DB.Query(query, tenantID, projectID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch units")
		return
	}
	defer rows.Close()

	units := []models.PropertyUnit{}
	for rows.Next() {
		var u models.PropertyUnit
		if err := rows.Scan(&u.ID, &u.TenantID, &u.ProjectID, &u.BlockID, &u.UnitNumber,
			&u.Floor, &u.UnitType, &u.Facing, &u.CarpetArea, &u.CarpetAreaWithBalcony,
			&u.UtilityArea, &u.PlinthArea, &u.SBUA, &u.UDSSqft, &u.Status,
			&u.AllotedTo, &u.AllotmentDate, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			continue
		}
		units = append(units, u)
	}

	h.respondJSON(w, http.StatusOK, units)
}

// ============================================
// CUSTOMER BOOKING ENDPOINTS
// ============================================

// CreateBooking creates a new customer booking
func (h *RealEstateHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateCustomerBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking := &models.CustomerBooking{
		TenantID:                tenantID,
		UnitID:                  req.UnitID,
		CustomerID:              &req.CustomerID,
		BookingDate:             req.BookingDate,
		BookingStatus:           "active",
		RatePerSqft:             req.RatePerSqft,
		CompositeGuidelineValue: req.CompositeGuidelineValue,
		CarParkingType:          req.CarParkingType,
		ParkingLocation:         req.ParkingLocation,
	}

	// Generate booking reference
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM customer_bookings WHERE tenant_id = $1", tenantID).Scan(&count)
	booking.BookingReference = fmt.Sprintf("BKG-%s-%d", tenantID, count+1)

	query := `INSERT INTO customer_bookings 
		(tenant_id, unit_id, customer_id, booking_date, booking_reference, booking_status,
		 rate_per_sqft, composite_guideline_value, car_parking_type, parking_location)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	err := h.DB.QueryRow(query,
		booking.TenantID, booking.UnitID, booking.CustomerID, booking.BookingDate,
		booking.BookingReference, booking.BookingStatus, booking.RatePerSqft,
		booking.CompositeGuidelineValue, booking.CarParkingType, booking.ParkingLocation,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create booking")
		return
	}

	// Update unit status
	h.DB.Exec("UPDATE property_units SET status = $1 WHERE id = $2", "booked", booking.UnitID)

	h.respondJSON(w, http.StatusCreated, booking)
}

// GetBookings retrieves all bookings for a tenant
func (h *RealEstateHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	query := `SELECT id, tenant_id, unit_id, lead_id, customer_id, booking_date, booking_reference,
		booking_status, welcome_date, allotment_date, agreement_date, registration_date,
		handover_date, possession_date, rate_per_sqft, composite_guideline_value,
		car_parking_type, parking_location, created_at, updated_at, deleted_at
		FROM customer_bookings WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY booking_date DESC`

	rows, err := h.DB.Query(query, tenantID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	defer rows.Close()

	bookings := []models.CustomerBooking{}
	for rows.Next() {
		var b models.CustomerBooking
		if err := rows.Scan(&b.ID, &b.TenantID, &b.UnitID, &b.LeadID, &b.CustomerID,
			&b.BookingDate, &b.BookingReference, &b.BookingStatus, &b.WelcomeDate,
			&b.AllotmentDate, &b.AgreementDate, &b.RegistrationDate, &b.HandoverDate,
			&b.PossessionDate, &b.RatePerSqft, &b.CompositeGuidelineValue,
			&b.CarParkingType, &b.ParkingLocation, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt); err != nil {
			continue
		}
		bookings = append(bookings, b)
	}

	h.respondJSON(w, http.StatusOK, bookings)
}

// ============================================
// PAYMENT ENDPOINTS
// ============================================

// RecordPayment records a payment for a booking
func (h *RealEstateHandler) RecordPayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.CreateBookingPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment := &models.BookingPayment{
		TenantID:      tenantID,
		BookingID:     req.BookingID,
		PaymentDate:   req.PaymentDate,
		PaymentMode:   req.PaymentMode,
		PaidBy:        req.PaidBy,
		ReceiptNumber: req.ReceiptNumber,
		Towards:       req.Towards,
		Amount:        req.Amount,
		BankName:      req.BankName,
		TransactionID: req.TransactionID,
		Remarks:       req.Remarks,
		Status:        "cleared",
	}

	query := `INSERT INTO booking_payments 
		(tenant_id, booking_id, payment_date, payment_mode, paid_by, receipt_number,
		 towards, amount, bank_name, transaction_id, status, remarks)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	err := h.DB.QueryRow(query,
		payment.TenantID, payment.BookingID, payment.PaymentDate, payment.PaymentMode,
		payment.PaidBy, payment.ReceiptNumber, payment.Towards, payment.Amount,
		payment.BankName, payment.TransactionID, payment.Status, payment.Remarks,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to record payment")
		return
	}

	// Create ledger entry
	h.createLedgerEntry(payment.BookingID, "credit", fmt.Sprintf("Payment received: %s", payment.Towards), payment.Amount)

	h.respondJSON(w, http.StatusCreated, payment)
}

// GetPayments retrieves all payments for a booking
func (h *RealEstateHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	bookingID := mux.Vars(r)["booking_id"]

	query := `SELECT id, tenant_id, booking_id, payment_date, payment_mode, paid_by,
		receipt_number, receipt_date, towards, amount, cheque_number, cheque_date,
		bank_name, transaction_id, status, remarks, created_at, updated_at, deleted_at
		FROM booking_payments WHERE tenant_id = $1 AND booking_id = $2 AND deleted_at IS NULL
		ORDER BY payment_date DESC`

	rows, err := h.DB.Query(query, tenantID, bookingID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}
	defer rows.Close()

	payments := []models.BookingPayment{}
	for rows.Next() {
		var p models.BookingPayment
		if err := rows.Scan(&p.ID, &p.TenantID, &p.BookingID, &p.PaymentDate, &p.PaymentMode,
			&p.PaidBy, &p.ReceiptNumber, &p.ReceiptDate, &p.Towards, &p.Amount,
			&p.ChequeNumber, &p.ChequeDate, &p.BankName, &p.TransactionID, &p.Status,
			&p.Remarks, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			continue
		}
		payments = append(payments, p)
	}

	h.respondJSON(w, http.StatusOK, payments)
}

// ============================================
// MILESTONE TRACKING ENDPOINTS
// ============================================

// TrackMilestone tracks important dates and campaign info
func (h *RealEstateHandler) TrackMilestone(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)

	var req models.PropertyMilestoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	milestone := &models.PropertyMilestone{
		TenantID:          tenantID,
		BookingID:         req.BookingID,
		CampaignName:      req.CampaignName,
		Source:            req.Source,
		SubSource:         req.SubSource,
		LeadGeneratedDate: req.LeadGeneratedDate,
		ReEngagedDate:     req.ReEngagedDate,
		SiteVisitDate:     req.SiteVisitDate,
		ReVisitDate:       req.ReVisitDate,
		BookingDate:       req.BookingDate,
		CancelledDate:     req.CancelledDate,
		Notes:             req.Notes,
		Status:            "active",
	}

	query := `INSERT INTO property_milestones 
		(tenant_id, booking_id, campaign_name, source, subsource, lead_generated_date,
		 re_engaged_date, site_visit_date, revisit_date, booking_date, cancelled_date, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	err := h.DB.QueryRow(query,
		milestone.TenantID, milestone.BookingID, milestone.CampaignName, milestone.Source,
		milestone.SubSource, milestone.LeadGeneratedDate, milestone.ReEngagedDate,
		milestone.SiteVisitDate, milestone.ReVisitDate, milestone.BookingDate,
		milestone.CancelledDate, milestone.Status, milestone.Notes,
	).Scan(&milestone.ID, &milestone.CreatedAt, &milestone.UpdatedAt)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to track milestone")
		return
	}

	h.respondJSON(w, http.StatusCreated, milestone)
}

// GetMilestones retrieves milestones for a booking
func (h *RealEstateHandler) GetMilestones(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	bookingID := mux.Vars(r)["booking_id"]

	query := `SELECT id, tenant_id, booking_id, campaign_id, campaign_name, source, subsource,
		lead_generated_date, re_engaged_date, site_visit_date, revisit_date, booking_date,
		cancelled_date, status, notes, created_at, updated_at, deleted_at
		FROM property_milestones WHERE tenant_id = $1 AND booking_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := h.DB.Query(query, tenantID, bookingID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch milestones")
		return
	}
	defer rows.Close()

	milestones := []models.PropertyMilestone{}
	for rows.Next() {
		var m models.PropertyMilestone
		if err := rows.Scan(&m.ID, &m.TenantID, &m.BookingID, &m.CampaignID, &m.CampaignName,
			&m.Source, &m.SubSource, &m.LeadGeneratedDate, &m.ReEngagedDate,
			&m.SiteVisitDate, &m.ReVisitDate, &m.BookingDate, &m.CancelledDate,
			&m.Status, &m.Notes, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt); err != nil {
			continue
		}
		milestones = append(milestones, m)
	}

	h.respondJSON(w, http.StatusOK, milestones)
}

// ============================================
// LEDGER ENDPOINTS
// ============================================

// GetAccountLedger retrieves account ledger for a booking
func (h *RealEstateHandler) GetAccountLedger(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Context().Value(middleware.TenantIDKey).(string)
	bookingID := mux.Vars(r)["booking_id"]

	query := `SELECT id, tenant_id, booking_id, customer_id, transaction_date, transaction_type,
		description, debit_amount, credit_amount, opening_balance, closing_balance, payment_id,
		reference_number, created_at, updated_at, deleted_at
		FROM customer_account_ledgers WHERE tenant_id = $1 AND booking_id = $2 AND deleted_at IS NULL
		ORDER BY transaction_date ASC`

	rows, err := h.DB.Query(query, tenantID, bookingID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to fetch ledger")
		return
	}
	defer rows.Close()

	ledgers := []models.CustomerAccountLedger{}
	for rows.Next() {
		var l models.CustomerAccountLedger
		if err := rows.Scan(&l.ID, &l.TenantID, &l.BookingID, &l.CustomerID, &l.TransactionDate,
			&l.TransactionType, &l.Description, &l.DebitAmount, &l.CreditAmount,
			&l.OpeningBalance, &l.ClosingBalance, &l.PaymentID, &l.ReferenceNum,
			&l.CreatedAt, &l.UpdatedAt, &l.DeletedAt); err != nil {
			continue
		}
		ledgers = append(ledgers, l)
	}

	h.respondJSON(w, http.StatusOK, ledgers)
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func (h *RealEstateHandler) createLedgerEntry(bookingID string, txnType string, description string, amount float64) error {
	var openingBalance float64
	var closingBalance float64

	// Get last balance
	h.DB.QueryRow(`SELECT COALESCE(closing_balance, 0) FROM customer_account_ledgers 
		WHERE booking_id = $1 ORDER BY transaction_date DESC LIMIT 1`, bookingID).Scan(&openingBalance)

	if txnType == "credit" {
		closingBalance = openingBalance + amount
	} else {
		closingBalance = openingBalance - amount
	}

	query := `INSERT INTO customer_account_ledgers 
		(tenant_id, booking_id, transaction_date, transaction_type, description,
		 debit_amount, credit_amount, opening_balance, closing_balance)
		VALUES ($1, $2, NOW(), $3, $4, $5, $6, $7, $8)`

	debit := 0.0
	credit := 0.0
	if txnType == "credit" {
		credit = amount
	} else {
		debit = amount
	}

	_, err := h.DB.Exec(query, "", bookingID, txnType, description, debit, credit, openingBalance, closingBalance)
	return err
}

func (h *RealEstateHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *RealEstateHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	h.respondJSON(w, statusCode, map[string]string{"error": message})
}
