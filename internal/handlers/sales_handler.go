package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"multi-tenant-ai-callcenter/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SalesHandler struct {
	DB *sql.DB
}

// NewSalesHandler creates a new sales handler instance
func NewSalesHandler(db *sql.DB) *SalesHandler {
	return &SalesHandler{
		DB: db,
	}
}

// Helper function to respond with JSON
func (h *SalesHandler) respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

// Helper function to respond with error
func (h *SalesHandler) respondError(w http.ResponseWriter, code int, message string) {
	h.respondJSON(w, code, map[string]interface{}{"error": message})
}

// ============================================================================
// SALES LEADS HANDLERS
// ============================================================================

// CreateSalesLead creates a new sales lead
func (h *SalesHandler) CreateSalesLead(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		FirstName      string  `json:"first_name"`
		LastName       string  `json:"last_name"`
		Email          string  `json:"email"`
		Phone          string  `json:"phone"`
		CompanyName    string  `json:"company_name"`
		Industry       string  `json:"industry"`
		Source         string  `json:"source"`
		CampaignID     *string `json:"campaign_id"`
		AssignedTo     *string `json:"assigned_to"`
		NextActionDate *string `json:"next_action_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	leadID := uuid.New().String()
	leadCode := fmt.Sprintf("LEAD-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	query := `
		INSERT INTO sales_leads (
			id, tenant_id, lead_code, first_name, last_name, email, phone,
			company_name, industry, status, probability, source, campaign_id,
			assigned_to, assigned_date, converted_to_customer, next_action_date,
			created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var assignedDate *time.Time
	if req.AssignedTo != nil {
		assignedDate = &now
	}

	var nextActionDate *time.Time
	if req.NextActionDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.NextActionDate)
		nextActionDate = &t
	}

	_, err := h.DB.Exec(query,
		leadID, tenantID, leadCode, req.FirstName, req.LastName, req.Email, req.Phone,
		req.CompanyName, req.Industry, "new", 0.0, req.Source, req.CampaignID,
		req.AssignedTo, assignedDate, false, nextActionDate,
		userID, now, now,
	)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create lead")
		return
	}

	lead := models.SalesLead{
		ID:                  leadID,
		TenantID:            tenantID,
		LeadCode:            leadCode,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		Email:               req.Email,
		Phone:               req.Phone,
		CompanyName:         req.CompanyName,
		Industry:            req.Industry,
		Status:              "new",
		Probability:         0.0,
		Source:              req.Source,
		CampaignID:          req.CampaignID,
		AssignedTo:          req.AssignedTo,
		AssignedDate:        assignedDate,
		ConvertedToCustomer: false,
		CreatedBy:           &userID,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": lead})
}

// GetSalesLead retrieves a specific sales lead
func (h *SalesHandler) GetSalesLead(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	leadID := mux.Vars(r)["id"]

	query := `
		SELECT id, tenant_id, lead_code, first_name, last_name, email, phone,
			company_name, industry, status, probability, source, campaign_id,
			assigned_to, assigned_date, converted_to_customer, customer_id,
			next_action_date, next_action_notes, created_by, created_at, updated_at
		FROM sales_leads
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var lead models.SalesLead
	err := h.DB.QueryRow(query, leadID, tenantID).Scan(
		&lead.ID, &lead.TenantID, &lead.LeadCode, &lead.FirstName, &lead.LastName,
		&lead.Email, &lead.Phone, &lead.CompanyName, &lead.Industry, &lead.Status,
		&lead.Probability, &lead.Source, &lead.CampaignID, &lead.AssignedTo,
		&lead.AssignedDate, &lead.ConvertedToCustomer, &lead.CustomerID,
		&lead.NextActionDate, &lead.NextActionNotes, &lead.CreatedBy,
		&lead.CreatedAt, &lead.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		h.respondError(w, http.StatusNotFound, "Lead not found")
		return
	}

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve lead")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": lead})
}

// ListSalesLeads retrieves all sales leads
func (h *SalesHandler) ListSalesLeads(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	query := `
		SELECT id, tenant_id, lead_code, first_name, last_name, email, phone,
			company_name, industry, status, probability, source, campaign_id,
			assigned_to, assigned_date, converted_to_customer, customer_id,
			next_action_date, next_action_notes, created_by, created_at, updated_at
		FROM sales_leads
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := h.DB.Query(query, tenantID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve leads")
		return
	}
	defer rows.Close()

	var leads []models.SalesLead
	for rows.Next() {
		var lead models.SalesLead
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.LeadCode, &lead.FirstName, &lead.LastName,
			&lead.Email, &lead.Phone, &lead.CompanyName, &lead.Industry, &lead.Status,
			&lead.Probability, &lead.Source, &lead.CampaignID, &lead.AssignedTo,
			&lead.AssignedDate, &lead.ConvertedToCustomer, &lead.CustomerID,
			&lead.NextActionDate, &lead.NextActionNotes, &lead.CreatedBy,
			&lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			continue
		}
		leads = append(leads, lead)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": leads})
}

// UpdateSalesLead updates an existing sales lead
func (h *SalesHandler) UpdateSalesLead(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Status      *string `json:"status"`
		Probability *int    `json:"probability"`
		AssignedTo  *string `json:"assigned_to"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateParts := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Status != nil {
		updateParts = append(updateParts, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *req.Status)
		argIdx++
	}
	if req.Probability != nil {
		updateParts = append(updateParts, fmt.Sprintf("probability = $%d", argIdx))
		args = append(args, *req.Probability)
		argIdx++
	}
	if req.AssignedTo != nil {
		updateParts = append(updateParts, fmt.Sprintf("assigned_to = $%d", argIdx))
		args = append(args, *req.AssignedTo)
		argIdx++
	}

	if len(updateParts) == 0 {
		h.respondError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	args = append(args, id, tenantID)

	query := fmt.Sprintf(
		"UPDATE sales_leads SET %s, updated_at = NOW() WHERE id = $%d AND tenant_id = $%d",
		strings.Join(updateParts, ", "),
		argIdx,
		argIdx+1,
	)

	result, err := h.DB.Exec(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update lead")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "Lead not found")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Lead updated"})
}

// DeleteSalesLead deletes a sales lead (soft delete)
func (h *SalesHandler) DeleteSalesLead(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	query := "UPDATE sales_leads SET deleted_at = NOW() WHERE id = $1 AND tenant_id = $2"
	result, err := h.DB.Exec(query, id, tenantID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete lead")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "Lead not found")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Lead deleted"})
}

// ============================================================================
// SALES CUSTOMERS HANDLERS
// ============================================================================

// CreateSalesCustomer creates a new sales customer
func (h *SalesHandler) CreateSalesCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		CustomerName       string  `json:"customer_name"`
		BusinessName       string  `json:"business_name"`
		BusinessType       string  `json:"business_type"`
		Industry           string  `json:"industry"`
		PrimaryContactName string  `json:"primary_contact_name"`
		PrimaryEmail       string  `json:"primary_email"`
		PrimaryPhone       string  `json:"primary_phone"`
		BillingAddress     string  `json:"billing_address"`
		BillingCity        string  `json:"billing_city"`
		BillingState       string  `json:"billing_state"`
		BillingCountry     string  `json:"billing_country"`
		BillingZip         string  `json:"billing_zip"`
		ShippingAddress    string  `json:"shipping_address"`
		ShippingCity       string  `json:"shipping_city"`
		ShippingState      string  `json:"shipping_state"`
		ShippingCountry    string  `json:"shipping_country"`
		ShippingZip        string  `json:"shipping_zip"`
		PANNumber          string  `json:"pan_number"`
		GSTNumber          string  `json:"gst_number"`
		CreditLimit        float64 `json:"credit_limit"`
		CreditDays         int     `json:"credit_days"`
		PaymentTerms       string  `json:"payment_terms"`
		CustomerCategory   string  `json:"customer_category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customerID := uuid.New().String()
	customerCode := fmt.Sprintf("CUST-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	query := `
		INSERT INTO sales_customers (
			id, tenant_id, customer_code, customer_name, business_name, business_type,
			industry, primary_contact_name, primary_email, primary_phone,
			billing_address, billing_city, billing_state, billing_country, billing_zip,
			shipping_address, shipping_city, shipping_state, shipping_country, shipping_zip,
			pan_number, gst_number, credit_limit, credit_days, payment_terms,
			customer_category, status, current_balance, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query,
		customerID, tenantID, customerCode, req.CustomerName, req.BusinessName, req.BusinessType,
		req.Industry, req.PrimaryContactName, req.PrimaryEmail, req.PrimaryPhone,
		req.BillingAddress, req.BillingCity, req.BillingState, req.BillingCountry, req.BillingZip,
		req.ShippingAddress, req.ShippingCity, req.ShippingState, req.ShippingCountry, req.ShippingZip,
		req.PANNumber, req.GSTNumber, req.CreditLimit, req.CreditDays, req.PaymentTerms,
		req.CustomerCategory, "active", 0.0, &userID, now, now,
	)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create customer")
		return
	}

	customer := models.SalesCustomer{
		ID:                 customerID,
		TenantID:           tenantID,
		CustomerCode:       customerCode,
		CustomerName:       req.CustomerName,
		BusinessName:       req.BusinessName,
		BusinessType:       req.BusinessType,
		Industry:           req.Industry,
		PrimaryContactName: req.PrimaryContactName,
		PrimaryEmail:       req.PrimaryEmail,
		PrimaryPhone:       req.PrimaryPhone,
		BillingAddress:     req.BillingAddress,
		BillingCity:        req.BillingCity,
		BillingState:       req.BillingState,
		BillingCountry:     req.BillingCountry,
		BillingZip:         req.BillingZip,
		ShippingAddress:    req.ShippingAddress,
		ShippingCity:       req.ShippingCity,
		ShippingState:      req.ShippingState,
		ShippingCountry:    req.ShippingCountry,
		ShippingZip:        req.ShippingZip,
		PANNumber:          req.PANNumber,
		GSTNumber:          req.GSTNumber,
		CreditLimit:        req.CreditLimit,
		CreditDays:         req.CreditDays,
		PaymentTerms:       req.PaymentTerms,
		CustomerCategory:   req.CustomerCategory,
		Status:             "active",
		CurrentBalance:     0.0,
		CreatedBy:          &userID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{"success": true, "data": customer})
}

// GetSalesCustomer retrieves a specific sales customer
func (h *SalesHandler) GetSalesCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	customerID := mux.Vars(r)["id"]

	query := `
		SELECT id, tenant_id, customer_code, customer_name, business_name, business_type,
			industry, primary_contact_name, primary_email, primary_phone,
			billing_address, billing_city, billing_state, billing_country, billing_zip,
			shipping_address, shipping_city, shipping_state, shipping_country, shipping_zip,
			pan_number, gst_number, credit_limit, credit_days, payment_terms,
			customer_category, status, current_balance, created_by, created_at, updated_at
		FROM sales_customers
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var customer models.SalesCustomer
	err := h.DB.QueryRow(query, customerID, tenantID).Scan(
		&customer.ID, &customer.TenantID, &customer.CustomerCode, &customer.CustomerName,
		&customer.BusinessName, &customer.BusinessType, &customer.Industry,
		&customer.PrimaryContactName, &customer.PrimaryEmail, &customer.PrimaryPhone,
		&customer.BillingAddress, &customer.BillingCity, &customer.BillingState,
		&customer.BillingCountry, &customer.BillingZip,
		&customer.ShippingAddress, &customer.ShippingCity, &customer.ShippingState,
		&customer.ShippingCountry, &customer.ShippingZip,
		&customer.PANNumber, &customer.GSTNumber, &customer.CreditLimit, &customer.CreditDays,
		&customer.PaymentTerms, &customer.CustomerCategory, &customer.Status,
		&customer.CurrentBalance, &customer.CreatedBy, &customer.CreatedAt, &customer.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		h.respondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve customer")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": customer})
}

// ListSalesCustomers retrieves all sales customers
func (h *SalesHandler) ListSalesCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	query := `
		SELECT id, tenant_id, customer_code, customer_name, business_name, business_type,
			industry, primary_contact_name, primary_email, primary_phone,
			billing_address, billing_city, billing_state, billing_country, billing_zip,
			shipping_address, shipping_city, shipping_state, shipping_country, shipping_zip,
			pan_number, gst_number, credit_limit, credit_days, payment_terms,
			customer_category, status, current_balance, created_by, created_at, updated_at
		FROM sales_customers
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := h.DB.Query(query, tenantID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve customers")
		return
	}
	defer rows.Close()

	var customers []models.SalesCustomer
	for rows.Next() {
		var customer models.SalesCustomer
		err := rows.Scan(
			&customer.ID, &customer.TenantID, &customer.CustomerCode, &customer.CustomerName,
			&customer.BusinessName, &customer.BusinessType, &customer.Industry,
			&customer.PrimaryContactName, &customer.PrimaryEmail, &customer.PrimaryPhone,
			&customer.BillingAddress, &customer.BillingCity, &customer.BillingState,
			&customer.BillingCountry, &customer.BillingZip,
			&customer.ShippingAddress, &customer.ShippingCity, &customer.ShippingState,
			&customer.ShippingCountry, &customer.ShippingZip,
			&customer.PANNumber, &customer.GSTNumber, &customer.CreditLimit, &customer.CreditDays,
			&customer.PaymentTerms, &customer.CustomerCategory, &customer.Status,
			&customer.CurrentBalance, &customer.CreatedBy, &customer.CreatedAt, &customer.UpdatedAt,
		)
		if err != nil {
			continue
		}
		customers = append(customers, customer)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": customers})
}

// UpdateSalesCustomer updates a sales customer
func (h *SalesHandler) UpdateSalesCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	customerID := mux.Vars(r)["id"]

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates := []string{}
	args := []interface{}{}

	allowedFields := map[string]bool{
		"customer_name": true, "business_name": true, "primary_email": true,
		"primary_phone": true, "billing_address": true, "billing_city": true,
		"billing_state": true, "gst_number": true, "credit_limit": true,
		"credit_days": true, "status": true, "current_balance": true,
	}

	for field, allowed := range allowedFields {
		if allowed && req[field] != nil {
			updates = append(updates, field+" = ?")
			args = append(args, req[field])
		}
	}

	if len(updates) == 0 {
		h.respondError(w, http.StatusBadRequest, "No valid fields to update")
		return
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, customerID, tenantID)

	query := fmt.Sprintf(`
		UPDATE sales_customers
		SET %s
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, strings.Join(updates, ", "))

	result, err := h.DB.Exec(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update customer")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Customer updated"})
}
