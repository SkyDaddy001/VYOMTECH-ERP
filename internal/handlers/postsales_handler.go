package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"vyomtech-backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PostSalesHandler manages post-sales operations
type PostSalesHandler struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewPostSalesHandler creates a new post-sales handler
func NewPostSalesHandler(db *sql.DB, logger *logger.Logger) *PostSalesHandler {
	return &PostSalesHandler{
		db:     db,
		logger: logger,
	}
}

// PaymentRequest represents payment creation request
type PaymentRequest struct {
	BookingID         string  `json:"booking_id"`
	PaymentScheduleID string  `json:"payment_schedule_id,omitempty"`
	CustomerName      string  `json:"customer_name"`
	UnitID            string  `json:"unit_id,omitempty"`
	ReceiptNo         string  `json:"receipt_no"`
	ReceivedOn        string  `json:"received_on"`
	ClearedOn         string  `json:"cleared_on,omitempty"`
	PaymentDate       string  `json:"payment_date"`
	PaymentMode       string  `json:"payment_mode"`
	PaidBy            string  `json:"paid_by"`
	Towards           string  `json:"towards"`
	Amount            float64 `json:"amount"`
}

// CreatePayment handles POST /api/v1/post-sales/payments
func (h *PostSalesHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	receivedOn, _ := time.Parse("2006-01-02", req.ReceivedOn)
	paymentDate, _ := time.Parse("2006-01-02", req.PaymentDate)

	query := `INSERT INTO payment (id, tenant_id, booking_id, payment_schedule_id, customer_name, 
		unit_id, receipt_no, received_on, payment_date, payment_mode, paid_by, towards, amount, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := h.db.Exec(query, id, tenantID, req.BookingID, req.PaymentScheduleID, req.CustomerName,
		req.UnitID, req.ReceiptNo, receivedOn, paymentDate, req.PaymentMode, req.PaidBy, req.Towards, req.Amount)

	if err != nil {
		h.logger.Error("Failed to create payment: %v", err)
		http.Error(w, "Failed to create payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Payment recorded successfully",
	})
}

// ListPayments handles GET /api/v1/post-sales/payments
func (h *PostSalesHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, tenant_id, booking_id, customer_name, receipt_no, received_on, 
		payment_date, payment_mode, paid_by, towards, amount, created_at 
		FROM payment WHERE tenant_id = ? ORDER BY payment_date DESC`

	rows, err := h.db.Query(query, tenantID)
	if err != nil {
		h.logger.Error("Failed to list payments: %v", err)
		http.Error(w, "Failed to fetch payments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var payments []map[string]interface{}
	for rows.Next() {
		var id, tenantID, bookingID, customerName, receiptNo, paymentMode, paidBy, towards string
		var amount float64
		var receivedOn, paymentDate, createdAt time.Time

		err := rows.Scan(&id, &tenantID, &bookingID, &customerName, &receiptNo, &receivedOn,
			&paymentDate, &paymentMode, &paidBy, &towards, &amount, &createdAt)
		if err != nil {
			continue
		}

		payments = append(payments, map[string]interface{}{
			"id":            id,
			"tenant_id":     tenantID,
			"booking_id":    bookingID,
			"customer_name": customerName,
			"receipt_no":    receiptNo,
			"received_on":   receivedOn.Format("2006-01-02"),
			"payment_date":  paymentDate.Format("2006-01-02"),
			"payment_mode":  paymentMode,
			"paid_by":       paidBy,
			"towards":       towards,
			"amount":        amount,
			"created_at":    createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    payments,
	})
}

// GetPayment handles GET /api/v1/post-sales/payments/:id
func (h *PostSalesHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	vars := mux.Vars(r)
	paymentID := vars["id"]

	query := `SELECT id, tenant_id, booking_id, customer_name, receipt_no, received_on, 
		payment_date, payment_mode, paid_by, towards, amount, created_at 
		FROM payment WHERE id = ? AND tenant_id = ?`

	var id, tid, bookingID, customerName, receiptNo, paymentMode, paidBy, towards string
	var amount float64
	var receivedOn, paymentDate, createdAt time.Time

	err := h.db.QueryRow(query, paymentID, tenantID).Scan(&id, &tid, &bookingID, &customerName,
		&receiptNo, &receivedOn, &paymentDate, &paymentMode, &paidBy, &towards, &amount, &createdAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("Failed to get payment: %v", err)
		http.Error(w, "Failed to fetch payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":            id,
			"tenant_id":     tid,
			"booking_id":    bookingID,
			"customer_name": customerName,
			"receipt_no":    receiptNo,
			"received_on":   receivedOn.Format("2006-01-02"),
			"payment_date":  paymentDate.Format("2006-01-02"),
			"payment_mode":  paymentMode,
			"paid_by":       paidBy,
			"towards":       towards,
			"amount":        amount,
			"created_at":    createdAt,
		},
	})
}

// UpdatePayment handles PUT /api/v1/post-sales/payments/:id
func (h *PostSalesHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	vars := mux.Vars(r)
	paymentID := vars["id"]

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE payment SET customer_name = ?, payment_mode = ?, paid_by = ?, 
		towards = ?, amount = ?, updated_at = NOW() WHERE id = ? AND tenant_id = ?`

	result, err := h.db.Exec(query, req.CustomerName, req.PaymentMode, req.PaidBy,
		req.Towards, req.Amount, paymentID, tenantID)

	if err != nil {
		h.logger.Error("Failed to update payment: %v", err)
		http.Error(w, "Failed to update payment", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Payment updated successfully",
	})
}

// DeletePayment handles DELETE /api/v1/post-sales/payments/:id
func (h *PostSalesHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	vars := mux.Vars(r)
	paymentID := vars["id"]

	query := `DELETE FROM payment WHERE id = ? AND tenant_id = ?`
	result, err := h.db.Exec(query, paymentID, tenantID)

	if err != nil {
		h.logger.Error("Failed to delete payment: %v", err)
		http.Error(w, "Failed to delete payment", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Payment deleted successfully",
	})
}

// PaymentScheduleRequest represents payment schedule creation request
type PaymentScheduleRequest struct {
	BookingID         string  `json:"booking_id"`
	PaymentStage      int     `json:"payment_stage"`
	ConstructionStage string  `json:"construction_stage"`
	ScheduledDate     string  `json:"scheduled_date"`
	AmountDue         float64 `json:"amount_due"`
	PaymentType       string  `json:"payment_type"`
	Status            string  `json:"status"`
}

// CreatePaymentSchedule handles POST /api/v1/post-sales/payment-schedules
func (h *PostSalesHandler) CreatePaymentSchedule(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req PaymentScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	scheduledDate, _ := time.Parse("2006-01-02", req.ScheduledDate)

	query := `INSERT INTO payment_schedule (id, tenant_id, booking_id, payment_stage, construction_stage, 
		scheduled_date, amount_due, payment_type, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := h.db.Exec(query, id, tenantID, req.BookingID, req.PaymentStage, req.ConstructionStage,
		scheduledDate, req.AmountDue, req.PaymentType, req.Status)

	if err != nil {
		h.logger.Error("Failed to create payment schedule: %v", err)
		http.Error(w, "Failed to create payment schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Payment schedule created successfully",
	})
}

// ListPaymentSchedules handles GET /api/v1/post-sales/payment-schedules
func (h *PostSalesHandler) ListPaymentSchedules(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	query := `SELECT id, tenant_id, booking_id, payment_stage, construction_stage, 
		scheduled_date, amount_due, payment_type, status, created_at 
		FROM payment_schedule WHERE tenant_id = ?`

	args := []interface{}{tenantID}
	if bookingID != "" {
		query += " AND booking_id = ?"
		args = append(args, bookingID)
	}
	query += " ORDER BY payment_stage ASC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		h.logger.Error("Failed to list payment schedules: %v", err)
		http.Error(w, "Failed to fetch payment schedules", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var schedules []map[string]interface{}
	for rows.Next() {
		var id, tenantID, bookingID, constructionStage, paymentType, status string
		var paymentStage int
		var amountDue float64
		var scheduledDate, createdAt time.Time

		err := rows.Scan(&id, &tenantID, &bookingID, &paymentStage, &constructionStage,
			&scheduledDate, &amountDue, &paymentType, &status, &createdAt)
		if err != nil {
			continue
		}

		schedules = append(schedules, map[string]interface{}{
			"id":                 id,
			"tenant_id":          tenantID,
			"booking_id":         bookingID,
			"payment_stage":      paymentStage,
			"construction_stage": constructionStage,
			"scheduled_date":     scheduledDate.Format("2006-01-02"),
			"amount_due":         amountDue,
			"payment_type":       paymentType,
			"status":             status,
			"created_at":         createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    schedules,
	})
}

// BankLoanRequest represents bank loan creation request
type BankLoanRequest struct {
	BookingID          string  `json:"booking_id"`
	BankName           string  `json:"bank_name"`
	ContactPerson      string  `json:"contact_person"`
	Phone              string  `json:"phone"`
	LoanSanctionDate   string  `json:"loan_sanction_date"`
	ConnectorCode      string  `json:"connector_code"`
	SanctionAmount     float64 `json:"sanction_amount"`
	DisbursedAmount    float64 `json:"disbursed_amount"`
	DisbursementDate   string  `json:"disbursement_date"`
	DisbursementStatus string  `json:"disbursement_status"`
}

// CreateBankLoan handles POST /api/v1/post-sales/bank-loans
func (h *PostSalesHandler) CreateBankLoan(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req BankLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	sanctionDate, _ := time.Parse("2006-01-02", req.LoanSanctionDate)

	query := `INSERT INTO bank_loan (id, tenant_id, booking_id, bank_name, contact_person, phone, 
		loan_sanction_date, connector_code, sanction_amount, disbursed_amount, disbursement_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := h.db.Exec(query, id, tenantID, req.BookingID, req.BankName, req.ContactPerson, req.Phone,
		sanctionDate, req.ConnectorCode, req.SanctionAmount, req.DisbursedAmount, req.DisbursementStatus)

	if err != nil {
		h.logger.Error("Failed to create bank loan: %v", err)
		http.Error(w, "Failed to create bank loan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Bank loan details created successfully",
	})
}

// ListBankLoans handles GET /api/v1/post-sales/bank-loans
func (h *PostSalesHandler) ListBankLoans(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, tenant_id, booking_id, bank_name, contact_person, phone, 
		loan_sanction_date, sanction_amount, disbursed_amount, disbursement_status, created_at 
		FROM bank_loan WHERE tenant_id = ? ORDER BY created_at DESC`

	rows, err := h.db.Query(query, tenantID)
	if err != nil {
		h.logger.Error("Failed to list bank loans: %v", err)
		http.Error(w, "Failed to fetch bank loans", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var loans []map[string]interface{}
	for rows.Next() {
		var id, tenantID, bookingID, bankName, contactPerson, phone, disbursementStatus string
		var sanctionAmount, disbursedAmount float64
		var sanctionDate, createdAt time.Time

		err := rows.Scan(&id, &tenantID, &bookingID, &bankName, &contactPerson, &phone,
			&sanctionDate, &sanctionAmount, &disbursedAmount, &disbursementStatus, &createdAt)
		if err != nil {
			continue
		}

		loans = append(loans, map[string]interface{}{
			"id":                  id,
			"tenant_id":           tenantID,
			"booking_id":          bookingID,
			"bank_name":           bankName,
			"contact_person":      contactPerson,
			"phone":               phone,
			"loan_sanction_date":  sanctionDate.Format("2006-01-02"),
			"sanction_amount":     sanctionAmount,
			"disbursed_amount":    disbursedAmount,
			"disbursement_status": disbursementStatus,
			"created_at":          createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    loans,
	})
}

// UpdateBankLoan handles PUT /api/v1/post-sales/bank-loans/:id
func (h *PostSalesHandler) UpdateBankLoan(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	vars := mux.Vars(r)
	loanID := vars["id"]

	var req BankLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `UPDATE bank_loan SET bank_name = ?, contact_person = ?, phone = ?, 
		sanction_amount = ?, disbursed_amount = ?, disbursement_status = ?, updated_at = NOW() 
		WHERE id = ? AND tenant_id = ?`

	result, err := h.db.Exec(query, req.BankName, req.ContactPerson, req.Phone,
		req.SanctionAmount, req.DisbursedAmount, req.DisbursementStatus, loanID, tenantID)

	if err != nil {
		h.logger.Error("Failed to update bank loan: %v", err)
		http.Error(w, "Failed to update bank loan", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Bank loan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Bank loan updated successfully",
	})
}

// RegistrationRequest represents registration details creation request
type RegistrationRequest struct {
	BookingID                 string  `json:"booking_id"`
	GSTApplicable             bool    `json:"gst_applicable"`
	GSTPercentage             float64 `json:"gst_percentage"`
	GSTCost                   float64 `json:"gst_cost"`
	ApartmentCostIncludingGST float64 `json:"apartment_cost_including_gst"`
	RegistrationType          string  `json:"registration_type"`
	RegistrationCost          float64 `json:"registration_cost"`
	NOCReceivedDate           string  `json:"noc_received_date"`
	Status                    string  `json:"status"`
}

// CreateRegistrationDetails handles POST /api/v1/post-sales/registrations
func (h *PostSalesHandler) CreateRegistrationDetails(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	query := `INSERT INTO registration_details (id, tenant_id, booking_id, gst_applicable, gst_percentage, 
		gst_cost, apartment_cost_including_gst, registration_type, registration_cost, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := h.db.Exec(query, id, tenantID, req.BookingID, req.GSTApplicable, req.GSTPercentage,
		req.GSTCost, req.ApartmentCostIncludingGST, req.RegistrationType, req.RegistrationCost, req.Status)

	if err != nil {
		h.logger.Error("Failed to create registration details: %v", err)
		http.Error(w, "Failed to create registration details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Registration details created successfully",
	})
}

// ListRegistrationDetails handles GET /api/v1/post-sales/registrations
func (h *PostSalesHandler) ListRegistrationDetails(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, tenant_id, booking_id, gst_applicable, gst_percentage, 
		gst_cost, apartment_cost_including_gst, registration_type, registration_cost, status, created_at 
		FROM registration_details WHERE tenant_id = ? ORDER BY created_at DESC`

	rows, err := h.db.Query(query, tenantID)
	if err != nil {
		h.logger.Error("Failed to list registration details: %v", err)
		http.Error(w, "Failed to fetch registration details", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var registrations []map[string]interface{}
	for rows.Next() {
		var id, tenantID, bookingID, registrationType, status string
		var gstApplicable bool
		var gstPercentage, gstCost, apartmentCostIncludingGST, registrationCost float64
		var createdAt time.Time

		err := rows.Scan(&id, &tenantID, &bookingID, &gstApplicable, &gstPercentage,
			&gstCost, &apartmentCostIncludingGST, &registrationType, &registrationCost, &status, &createdAt)
		if err != nil {
			continue
		}

		registrations = append(registrations, map[string]interface{}{
			"id":                           id,
			"tenant_id":                    tenantID,
			"booking_id":                   bookingID,
			"gst_applicable":               gstApplicable,
			"gst_percentage":               gstPercentage,
			"gst_cost":                     gstCost,
			"apartment_cost_including_gst": apartmentCostIncludingGST,
			"registration_type":            registrationType,
			"registration_cost":            registrationCost,
			"status":                       status,
			"created_at":                   createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    registrations,
	})
}

// AdditionalChargesRequest represents additional charges creation request
type AdditionalChargesRequest struct {
	BookingID         string  `json:"booking_id"`
	MaintenanceCharge float64 `json:"maintenance_charge"`
	CorpusCharge      float64 `json:"corpus_charge"`
	EBDeposit         float64 `json:"eb_deposit"`
	OtherWorksCharge  float64 `json:"other_works_charge"`
}

// CreateAdditionalCharges handles POST /api/v1/post-sales/additional-charges
func (h *PostSalesHandler) CreateAdditionalCharges(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	var req AdditionalChargesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	query := `INSERT INTO additional_charges (id, tenant_id, booking_id, maintenance_charge, 
		corpus_charge, eb_deposit, other_works_charge, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := h.db.Exec(query, id, tenantID, req.BookingID, req.MaintenanceCharge,
		req.CorpusCharge, req.EBDeposit, req.OtherWorksCharge)

	if err != nil {
		h.logger.Error("Failed to create additional charges: %v", err)
		http.Error(w, "Failed to create additional charges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Additional charges created successfully",
	})
}

// ListAdditionalCharges handles GET /api/v1/post-sales/additional-charges
func (h *PostSalesHandler) ListAdditionalCharges(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "X-Tenant-ID header required", http.StatusBadRequest)
		return
	}

	query := `SELECT id, tenant_id, booking_id, maintenance_charge, corpus_charge, eb_deposit, other_works_charge, created_at 
		FROM additional_charges WHERE tenant_id = ? ORDER BY created_at DESC`

	rows, err := h.db.Query(query, tenantID)
	if err != nil {
		h.logger.Error("Failed to list additional charges: %v", err)
		http.Error(w, "Failed to fetch additional charges", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var charges []map[string]interface{}
	for rows.Next() {
		var id, tenantID, bookingID string
		var maintenanceCharge, corpusCharge, ebDeposit, otherWorksCharge float64
		var createdAt time.Time

		err := rows.Scan(&id, &tenantID, &bookingID, &maintenanceCharge, &corpusCharge, &ebDeposit, &otherWorksCharge, &createdAt)
		if err != nil {
			continue
		}

		charges = append(charges, map[string]interface{}{
			"id":                 id,
			"tenant_id":          tenantID,
			"booking_id":         bookingID,
			"maintenance_charge": maintenanceCharge,
			"corpus_charge":      corpusCharge,
			"eb_deposit":         ebDeposit,
			"other_works_charge": otherWorksCharge,
			"created_at":         createdAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    charges,
	})
}
