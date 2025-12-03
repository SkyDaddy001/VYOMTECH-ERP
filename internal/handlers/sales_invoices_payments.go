package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"vyomtech-backend/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ============================================================================
// SALES INVOICES HANDLERS
// ============================================================================

// CreateSalesInvoice creates a new sales invoice from order
func (h *SalesHandler) CreateSalesInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		SalesOrderID   string                   `json:"sales_order_id"`
		InvoiceDate    string                   `json:"invoice_date"`
		DueDate        *string                  `json:"due_date"`
		Items          []map[string]interface{} `json:"items"`
		DiscountAmount float64                  `json:"discount_amount"`
		CGSTAmount     float64                  `json:"cgst_amount"`
		SGSTAmount     float64                  `json:"sgst_amount"`
		IGSTAmount     float64                  `json:"igst_amount"`
		Notes          *string                  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	invoiceID := uuid.New().String()
	invoiceNumber := "INV-" + time.Now().Format("20060102") + "-" + invoiceID[:8]
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	invoiceDate, _ := time.Parse(time.RFC3339, req.InvoiceDate)
	var dueDate *time.Time
	if req.DueDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.DueDate)
		dueDate = &t
	}

	// Get order info
	var customerID string
	var orderSubtotal, orderDiscount, orderTax, orderTotal float64

	orderQuery := `
		SELECT customer_id, subtotal_amount, discount_amount, tax_amount, total_amount
		FROM sales_orders
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := h.DB.QueryRow(orderQuery, req.SalesOrderID, tenantID).Scan(
		&customerID, &orderSubtotal, &orderDiscount, &orderTax, &orderTotal)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Order not found")
		return
	}

	// Calculate totals
	var subtotal float64
	totalTax := req.CGSTAmount + req.SGSTAmount + req.IGSTAmount
	totalAmount := subtotal - req.DiscountAmount + totalTax

	// Insert invoice
	invoiceQuery := `
		INSERT INTO sales_invoices (
			id, tenant_id, invoice_number, customer_id, sales_order_id, invoice_date,
			due_date, subtotal_amount, discount_amount, cgst_amount, sgst_amount,
			igst_amount, total_tax, total_amount, payment_status, paid_amount,
			pending_amount, ar_posting_status, document_status, notes, created_by,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = h.DB.Exec(invoiceQuery,
		invoiceID, tenantID, invoiceNumber, customerID, req.SalesOrderID, invoiceDate,
		dueDate, subtotal, req.DiscountAmount, req.CGSTAmount, req.SGSTAmount,
		req.IGSTAmount, totalTax, totalAmount, "unpaid", 0.0,
		totalAmount, "not_posted", "draft", req.Notes, &userID,
		now, now)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	// Insert invoice items
	itemQuery := `
		INSERT INTO sales_invoice_items (
			id, tenant_id, invoice_id, line_number, description, hsn_code, quantity,
			unit_price, line_total, discount_percent, discount_amount, cgst_rate,
			cgst_amount, sgst_rate, sgst_amount, igst_rate, igst_amount,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for idx, item := range req.Items {
		itemID := uuid.New().String()
		h.DB.Exec(itemQuery,
			itemID, tenantID, invoiceID, idx+1,
			item["description"], item["hsn_code"], item["quantity"],
			item["unit_price"], item["line_total"], item["discount_percent"],
			item["discount_amount"], item["cgst_rate"], item["cgst_amount"],
			item["sgst_rate"], item["sgst_amount"], item["igst_rate"],
			item["igst_amount"], now, now)
	}

	// Update order status
	h.DB.Exec(`
		UPDATE sales_orders
		SET status = ?, invoiced_amount = invoiced_amount + ?, 
			pending_amount = total_amount - (invoiced_amount + ?), updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`, "partially_invoiced", totalAmount, totalAmount, now, req.SalesOrderID, tenantID)

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             invoiceID,
			"invoice_number": invoiceNumber,
			"customer_id":    customerID,
			"invoice_date":   invoiceDate,
			"total_amount":   totalAmount,
			"status":         "unpaid",
		},
	})
}

// GetSalesInvoice retrieves a specific invoice with items
func (h *SalesHandler) GetSalesInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	invoiceID := mux.Vars(r)["id"]

	query := `
		SELECT id, tenant_id, invoice_number, customer_id, sales_order_id, invoice_date,
			due_date, subtotal_amount, discount_amount, cgst_amount, sgst_amount,
			igst_amount, total_tax, total_amount, payment_status, paid_amount,
			pending_amount, ar_posting_status, gl_reference_number, document_status,
			notes, created_by, created_at, updated_at
		FROM sales_invoices
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var invoice models.SalesInvoice
	err := h.DB.QueryRow(query, invoiceID, tenantID).Scan(
		&invoice.ID, &invoice.TenantID, &invoice.InvoiceNumber, &invoice.CustomerID,
		&invoice.SalesOrderID, &invoice.InvoiceDate, &invoice.DueDate,
		&invoice.SubtotalAmount, &invoice.DiscountAmount, &invoice.CGSTAmount,
		&invoice.SGSTAmount, &invoice.IGSTAmount, &invoice.TotalTax, &invoice.TotalAmount,
		&invoice.PaymentStatus, &invoice.PaidAmount, &invoice.PendingAmount,
		&invoice.ARPostingStatus, &invoice.GLReferenceNumber, &invoice.DocumentStatus,
		&invoice.Notes, &invoice.CreatedBy, &invoice.CreatedAt, &invoice.UpdatedAt)

	if err == sql.ErrNoRows {
		h.respondError(w, http.StatusNotFound, "Invoice not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve invoice")
		return
	}

	// Get invoice items
	itemQuery := `
		SELECT id, tenant_id, invoice_id, line_number, description, hsn_code, quantity,
			unit_price, line_total, discount_percent, discount_amount, cgst_rate,
			cgst_amount, sgst_rate, sgst_amount, igst_rate, igst_amount,
			created_at, updated_at
		FROM sales_invoice_items
		WHERE invoice_id = ? AND tenant_id = ?
		ORDER BY line_number
	`

	rows, _ := h.DB.Query(itemQuery, invoiceID, tenantID)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var item models.SalesInvoiceItem
			rows.Scan(&item.ID, &item.TenantID, &item.InvoiceID, &item.LineNumber,
				&item.Description, &item.HSNCode, &item.Quantity, &item.UnitPrice,
				&item.LineTotal, &item.DiscountPercent, &item.DiscountAmount,
				&item.CGSTRate, &item.CGSTAmount, &item.SGSTRate, &item.SGSTAmount,
				&item.IGSTRate, &item.IGSTAmount, &item.CreatedAt, &item.UpdatedAt)
			invoice.Items = append(invoice.Items, item)
		}
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": invoice})
}

// ListSalesInvoices retrieves invoices with filtering
func (h *SalesHandler) ListSalesInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	status := r.URL.Query().Get("status")
	customerID := r.URL.Query().Get("customer_id")

	query := `
		SELECT id, tenant_id, invoice_number, customer_id, sales_order_id, invoice_date,
			due_date, subtotal_amount, discount_amount, cgst_amount, sgst_amount,
			igst_amount, total_tax, total_amount, payment_status, paid_amount,
			pending_amount, ar_posting_status, gl_reference_number, document_status,
			notes, created_by, created_at, updated_at
		FROM sales_invoices
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	args := []interface{}{tenantID}

	if status != "" {
		query += " AND payment_status = ?"
		args = append(args, status)
	}
	if customerID != "" {
		query += " AND customer_id = ?"
		args = append(args, customerID)
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve invoices")
		return
	}
	defer rows.Close()

	var invoices []models.SalesInvoice
	for rows.Next() {
		var inv models.SalesInvoice
		rows.Scan(&inv.ID, &inv.TenantID, &inv.InvoiceNumber, &inv.CustomerID,
			&inv.SalesOrderID, &inv.InvoiceDate, &inv.DueDate, &inv.SubtotalAmount,
			&inv.DiscountAmount, &inv.CGSTAmount, &inv.SGSTAmount, &inv.IGSTAmount,
			&inv.TotalTax, &inv.TotalAmount, &inv.PaymentStatus, &inv.PaidAmount,
			&inv.PendingAmount, &inv.ARPostingStatus, &inv.GLReferenceNumber,
			&inv.DocumentStatus, &inv.Notes, &inv.CreatedBy, &inv.CreatedAt, &inv.UpdatedAt)
		invoices = append(invoices, inv)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": invoices})
}

// ============================================================================
// SALES PAYMENTS HANDLERS
// ============================================================================

// CreateSalesPayment records a payment against invoice
func (h *SalesHandler) CreateSalesPayment(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		InvoiceID       string  `json:"invoice_id"`
		PaymentDate     string  `json:"payment_date"`
		PaymentAmount   float64 `json:"payment_amount"`
		PaymentMethod   string  `json:"payment_method"`
		ReferenceNumber string  `json:"reference_number"`
		Notes           *string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	paymentID := uuid.New().String()
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	paymentDate, _ := time.Parse(time.RFC3339, req.PaymentDate)

	// Validate payment method
	validMethods := map[string]bool{
		"cheque": true, "bank_transfer": true, "cash": true,
		"credit_card": true, "digital_payment": true,
	}
	if !validMethods[req.PaymentMethod] {
		h.respondError(w, http.StatusBadRequest, "Invalid payment method")
		return
	}

	// Get invoice details
	var totalAmount, paidAmount float64
	invoiceQuery := `
		SELECT total_amount, paid_amount FROM sales_invoices
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	err := h.DB.QueryRow(invoiceQuery, req.InvoiceID, tenantID).Scan(&totalAmount, &paidAmount)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Invoice not found")
		return
	}

	// Check if payment exceeds invoice amount
	if paidAmount+req.PaymentAmount > totalAmount {
		h.respondError(w, http.StatusBadRequest, "Payment amount exceeds invoice total")
		return
	}

	// Insert payment
	query := `
		INSERT INTO sales_payments (
			id, tenant_id, invoice_id, payment_date, payment_amount, payment_method,
			reference_number, payment_status, notes, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = h.DB.Exec(query,
		paymentID, tenantID, req.InvoiceID, paymentDate, req.PaymentAmount,
		req.PaymentMethod, req.ReferenceNumber, "processed", req.Notes, &userID, now, now)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create payment")
		return
	}

	// Update invoice payment status
	newPaidAmount := paidAmount + req.PaymentAmount
	newPendingAmount := totalAmount - newPaidAmount
	var newStatus string

	if newPendingAmount <= 0 {
		newStatus = "paid"
	} else if newPaidAmount > 0 {
		newStatus = "partially_paid"
	} else {
		newStatus = "unpaid"
	}

	h.DB.Exec(`
		UPDATE sales_invoices
		SET paid_amount = ?, pending_amount = ?, payment_status = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`, newPaidAmount, newPendingAmount, newStatus, now, req.InvoiceID, tenantID)

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":             paymentID,
			"invoice_id":     req.InvoiceID,
			"payment_amount": req.PaymentAmount,
			"payment_method": req.PaymentMethod,
			"invoice_paid":   newStatus == "paid",
			"pending_amount": newPendingAmount,
		},
	})
}

// ============================================================================
// DELIVERY NOTES HANDLERS
// ============================================================================

// CreateDeliveryNote creates a delivery note for an order
func (h *SalesHandler) CreateDeliveryNote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		OrderID           string  `json:"order_id"`
		DeliveryDate      string  `json:"delivery_date"`
		DeliveryLocation  string  `json:"delivery_location"`
		DeliveredBy       string  `json:"delivered_by"`
		VehicleNumber     string  `json:"vehicle_number"`
		DeliveredQuantity float64 `json:"delivered_quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	deliveryID := uuid.New().String()
	now := time.Now()
	deliveryDate, _ := time.Parse(time.RFC3339, req.DeliveryDate)

	query := `
		INSERT INTO sales_delivery_notes (
			id, tenant_id, order_id, delivery_date, delivery_location, delivered_by,
			vehicle_number, delivered_quantity, delivery_status, pod_received,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query,
		deliveryID, tenantID, req.OrderID, deliveryDate, req.DeliveryLocation,
		req.DeliveredBy, req.VehicleNumber, req.DeliveredQuantity, "in_transit", false,
		now, now)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create delivery note")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":              deliveryID,
			"order_id":        req.OrderID,
			"delivery_date":   deliveryDate,
			"delivery_status": "in_transit",
			"pod_received":    false,
		},
	})
}

// UpdateDeliveryPOD updates delivery proof of delivery
func (h *SalesHandler) UpdateDeliveryPOD(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	deliveryID := mux.Vars(r)["id"]

	var req struct {
		ReceiverName         string `json:"receiver_name"`
		ReceiverSignatureURL string `json:"receiver_signature_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	now := time.Now()

	query := `
		UPDATE sales_delivery_notes
		SET pod_received = ?, pod_date = ?, receiver_name = ?, receiver_signature_url = ?,
			delivery_status = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?
	`

	result, err := h.DB.Exec(query,
		true, now, req.ReceiverName, req.ReceiverSignatureURL,
		"delivered", now, deliveryID, tenantID)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update delivery")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "Delivery note not found")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Delivery POD updated",
	})
}

// ============================================================================
// CREDIT NOTES HANDLERS
// ============================================================================

// CreateCreditNote creates a credit note for an invoice
func (h *SalesHandler) CreateCreditNote(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		InvoiceID string                   `json:"invoice_id"`
		Reason    string                   `json:"reason"`
		Items     []map[string]interface{} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	creditNoteID := uuid.New().String()
	creditNoteNumber := "CN-" + time.Now().Format("20060102") + "-" + creditNoteID[:8]
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	// Calculate total
	var totalAmount float64
	for _, item := range req.Items {
		if lt, ok := item["line_total"].(float64); ok {
			totalAmount += lt
		}
	}

	query := `
		INSERT INTO sales_credit_notes (
			id, tenant_id, invoice_id, credit_note_number, credit_note_date, reason,
			total_amount, status, ar_posting_status, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query,
		creditNoteID, tenantID, req.InvoiceID, creditNoteNumber, now, req.Reason,
		totalAmount, "draft", "not_posted", &userID, now, now)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create credit note")
		return
	}

	// Insert credit note items
	itemQuery := `
		INSERT INTO sales_credit_note_items (
			id, tenant_id, credit_note_id, line_number, description, quantity,
			unit_price, line_total, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for idx, item := range req.Items {
		itemID := uuid.New().String()
		h.DB.Exec(itemQuery,
			itemID, tenantID, creditNoteID, idx+1,
			item["description"], item["quantity"], item["unit_price"],
			item["line_total"], now)
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":                 creditNoteID,
			"credit_note_number": creditNoteNumber,
			"invoice_id":         req.InvoiceID,
			"total_amount":       totalAmount,
			"reason":             req.Reason,
		},
	})
}

// ============================================================================
// PERFORMANCE METRICS HANDLERS
// ============================================================================

// GetSalespersonMetrics retrieves performance metrics for a salesperson
func (h *SalesHandler) GetSalespersonMetrics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	salespersonID := mux.Vars(r)["id"]
	metricMonth := r.URL.Query().Get("month")

	query := `
		SELECT id, tenant_id, salesperson_id, metric_month, leads_generated,
			leads_qualified, conversion_rate, total_orders, total_order_value,
			average_order_value, target_amount, actual_amount, achievement_percent,
			commission_rate, commission_amount, created_at, updated_at
		FROM sales_performance_metrics
		WHERE tenant_id = ? AND salesperson_id = ? AND metric_month = ?
	`

	var metrics models.SalesPerformanceMetrics
	err := h.DB.QueryRow(query, tenantID, salespersonID, metricMonth).Scan(
		&metrics.ID, &metrics.TenantID, &metrics.SalespersonID, &metrics.MetricMonth,
		&metrics.LeadsGenerated, &metrics.LeadsQualified, &metrics.ConversionRate,
		&metrics.TotalOrders, &metrics.TotalOrderValue, &metrics.AverageOrderValue,
		&metrics.TargetAmount, &metrics.ActualAmount, &metrics.AchievementPercent,
		&metrics.CommissionRate, &metrics.CommissionAmount, &metrics.CreatedAt, &metrics.UpdatedAt)

	if err == sql.ErrNoRows {
		h.respondError(w, http.StatusNotFound, "Metrics not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve metrics")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": metrics})
}

// CalculateAndUpdateMetrics calculates and updates monthly metrics
func (h *SalesHandler) CalculateAndUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		SalespersonID string  `json:"salesperson_id"`
		MetricMonth   string  `json:"metric_month"`
		TargetAmount  float64 `json:"target_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Calculate actual metrics from orders
	var leadsGenerated, leadsQualified, totalOrders int
	var totalOrderValue float64

	leadQuery := `SELECT COUNT(*) FROM sales_leads WHERE tenant_id = ? AND assigned_to = ?`
	h.DB.QueryRow(leadQuery, tenantID, req.SalespersonID).Scan(&leadsGenerated)

	orderQuery := `
		SELECT COUNT(*), COALESCE(SUM(total_amount), 0)
		FROM sales_orders
		WHERE tenant_id = ? AND DATE_FORMAT(created_at, '%Y-%m') = ? AND deleted_at IS NULL
	`
	h.DB.QueryRow(orderQuery, tenantID, req.MetricMonth).Scan(&totalOrders, &totalOrderValue)

	conversionRate := 0.0
	if leadsGenerated > 0 {
		conversionRate = float64(totalOrders) / float64(leadsGenerated) * 100
	}

	averageOrderValue := 0.0
	if totalOrders > 0 {
		averageOrderValue = totalOrderValue / float64(totalOrders)
	}

	achievementPercent := 0.0
	if req.TargetAmount > 0 {
		achievementPercent = (totalOrderValue / req.TargetAmount) * 100
	}

	commissionRate := 0.0
	if achievementPercent >= 100 {
		commissionRate = 5.0
	} else if achievementPercent >= 80 {
		commissionRate = 3.0
	}

	commissionAmount := totalOrderValue * (commissionRate / 100)

	metricsID := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO sales_performance_metrics (
			id, tenant_id, salesperson_id, metric_month, leads_generated,
			leads_qualified, conversion_rate, total_orders, total_order_value,
			average_order_value, target_amount, actual_amount, achievement_percent,
			commission_rate, commission_amount, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query,
		metricsID, tenantID, req.SalespersonID, req.MetricMonth, leadsGenerated,
		leadsQualified, conversionRate, totalOrders, totalOrderValue, averageOrderValue,
		req.TargetAmount, totalOrderValue, achievementPercent, commissionRate,
		commissionAmount, now, now)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to calculate metrics")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"salesperson_id":      req.SalespersonID,
			"metric_month":        req.MetricMonth,
			"total_orders":        totalOrders,
			"total_order_value":   totalOrderValue,
			"achievement_percent": achievementPercent,
			"commission_amount":   commissionAmount,
		},
	})
}
