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
// SALES QUOTATIONS HANDLERS
// ============================================================================

// CreateSalesQuotation creates a new sales quotation
func (h *SalesHandler) CreateSalesQuotation(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		CustomerID         string                   `json:"customer_id"`
		QuotationDate      string                   `json:"quotation_date"`
		ValidUntil         *string                  `json:"valid_until"`
		Items              []map[string]interface{} `json:"items"`
		DiscountAmount     float64                  `json:"discount_amount"`
		Notes              *string                  `json:"notes"`
		TermsAndConditions *string                  `json:"terms_and_conditions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	quotationID := uuid.New().String()
	quotationNumber := "QT-" + time.Now().Format("20060102") + "-" + quotationID[:8]
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	quotationDate, _ := time.Parse(time.RFC3339, req.QuotationDate)
	var validUntil *time.Time
	if req.ValidUntil != nil {
		t, _ := time.Parse(time.RFC3339, *req.ValidUntil)
		validUntil = &t
	}

	// Calculate totals from items
	var subtotal, totalTax float64
	for _, item := range req.Items {
		if q, ok := item["quantity"].(float64); ok {
			if up, ok := item["unit_price"].(float64); ok {
				lineTotal := q * up
				subtotal += lineTotal
				if tr, ok := item["tax_rate"].(float64); ok {
					taxAmount := lineTotal * (tr / 100)
					totalTax += taxAmount
				}
			}
		}
	}

	totalAmount := subtotal - req.DiscountAmount + totalTax

	query := `
		INSERT INTO sales_quotations (
			id, tenant_id, quotation_number, customer_id, quotation_date, valid_until,
			subtotal_amount, discount_amount, tax_amount, total_amount, status,
			converted_to_order, notes, terms_and_conditions, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query,
		quotationID, tenantID, quotationNumber, req.CustomerID, quotationDate, validUntil,
		subtotal, req.DiscountAmount, totalTax, totalAmount, "draft",
		false, req.Notes, req.TermsAndConditions, &userID, now, now,
	)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create quotation")
		return
	}

	// Insert quotation items
	itemQuery := `
		INSERT INTO sales_quotation_items (
			id, tenant_id, quotation_id, line_number, description, product_service_code,
			quantity, unit_price, line_total, discount_percent, discount_amount,
			hsn_code, tax_rate, tax_amount, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for idx, item := range req.Items {
		itemID := uuid.New().String()
		quantity := item["quantity"].(float64)
		unitPrice := item["unit_price"].(float64)
		taxRate := item["tax_rate"].(float64)
		lineTotal := quantity * unitPrice
		taxAmount := lineTotal * (taxRate / 100)

		h.DB.Exec(itemQuery,
			itemID, tenantID, quotationID, idx+1,
			item["description"], item["product_service_code"],
			quantity, unitPrice, lineTotal,
			item["discount_percent"], item["discount_amount"],
			item["hsn_code"], taxRate, taxAmount,
			now, now,
		)
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":                 quotationID,
			"quotation_number":   quotationNumber,
			"customer_id":        req.CustomerID,
			"quotation_date":     quotationDate,
			"subtotal_amount":    subtotal,
			"discount_amount":    req.DiscountAmount,
			"tax_amount":         totalTax,
			"total_amount":       totalAmount,
			"status":             "draft",
			"converted_to_order": false,
		},
	})
}

// GetSalesQuotation retrieves a specific quotation with items
func (h *SalesHandler) GetSalesQuotation(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	quotationID := mux.Vars(r)["id"]

	query := `
		SELECT id, tenant_id, quotation_number, customer_id, quotation_date, valid_until,
			subtotal_amount, discount_amount, tax_amount, total_amount, status,
			converted_to_order, sales_order_id, notes, terms_and_conditions,
			created_by, created_at, updated_at
		FROM sales_quotations
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var quotation models.SalesQuotation
	err := h.DB.QueryRow(query, quotationID, tenantID).Scan(
		&quotation.ID, &quotation.TenantID, &quotation.QuotationNumber, &quotation.CustomerID,
		&quotation.QuotationDate, &quotation.ValidUntil, &quotation.SubtotalAmount,
		&quotation.DiscountAmount, &quotation.TaxAmount, &quotation.TotalAmount,
		&quotation.Status, &quotation.ConvertedToOrder, &quotation.SalesOrderID,
		&quotation.Notes, &quotation.TermsAndConditions, &quotation.CreatedBy,
		&quotation.CreatedAt, &quotation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		h.respondError(w, http.StatusNotFound, "Quotation not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve quotation")
		return
	}

	// Get quotation items
	itemQuery := `
		SELECT id, tenant_id, quotation_id, line_number, description, product_service_code,
			quantity, unit_price, line_total, discount_percent, discount_amount,
			hsn_code, tax_rate, tax_amount, created_at, updated_at
		FROM sales_quotation_items
		WHERE quotation_id = ? AND tenant_id = ?
		ORDER BY line_number
	`

	rows, _ := h.DB.Query(itemQuery, quotationID, tenantID)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var item models.SalesQuotationItem
			rows.Scan(&item.ID, &item.TenantID, &item.QuotationID, &item.LineNumber,
				&item.Description, &item.ProductServiceCode, &item.Quantity, &item.UnitPrice,
				&item.LineTotal, &item.DiscountPercent, &item.DiscountAmount,
				&item.HSNCode, &item.TaxRate, &item.TaxAmount, &item.CreatedAt, &item.UpdatedAt)
			quotation.Items = append(quotation.Items, item)
		}
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": quotation})
}

// ListSalesQuotations retrieves quotations with filtering
func (h *SalesHandler) ListSalesQuotations(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	status := r.URL.Query().Get("status")
	customerID := r.URL.Query().Get("customer_id")

	query := `
		SELECT id, tenant_id, quotation_number, customer_id, quotation_date, valid_until,
			subtotal_amount, discount_amount, tax_amount, total_amount, status,
			converted_to_order, sales_order_id, notes, terms_and_conditions,
			created_by, created_at, updated_at
		FROM sales_quotations
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	args := []interface{}{tenantID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if customerID != "" {
		query += " AND customer_id = ?"
		args = append(args, customerID)
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve quotations")
		return
	}
	defer rows.Close()

	var quotations []models.SalesQuotation
	for rows.Next() {
		var q models.SalesQuotation
		rows.Scan(&q.ID, &q.TenantID, &q.QuotationNumber, &q.CustomerID, &q.QuotationDate,
			&q.ValidUntil, &q.SubtotalAmount, &q.DiscountAmount, &q.TaxAmount, &q.TotalAmount,
			&q.Status, &q.ConvertedToOrder, &q.SalesOrderID, &q.Notes, &q.TermsAndConditions,
			&q.CreatedBy, &q.CreatedAt, &q.UpdatedAt)
		quotations = append(quotations, q)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": quotations})
}

// ============================================================================
// SALES ORDERS HANDLERS
// ============================================================================

// CreateSalesOrder creates a new sales order
func (h *SalesHandler) CreateSalesOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	var req struct {
		CustomerID           string                   `json:"customer_id"`
		OrderDate            string                   `json:"order_date"`
		RequiredByDate       *string                  `json:"required_by_date"`
		DeliveryLocation     string                   `json:"delivery_location"`
		DeliveryInstructions *string                  `json:"delivery_instructions"`
		Items                []map[string]interface{} `json:"items"`
		DiscountAmount       float64                  `json:"discount_amount"`
		Notes                *string                  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	orderID := uuid.New().String()
	orderNumber := "SO-" + time.Now().Format("20060102") + "-" + orderID[:8]
	now := time.Now()
	userID := r.Header.Get("X-User-ID")

	orderDate, _ := time.Parse(time.RFC3339, req.OrderDate)
	var requiredByDate *time.Time
	if req.RequiredByDate != nil {
		t, _ := time.Parse(time.RFC3339, *req.RequiredByDate)
		requiredByDate = &t
	}

	// Calculate totals
	var subtotal, totalTax float64
	for _, item := range req.Items {
		if q, ok := item["quantity"].(float64); ok {
			if up, ok := item["unit_price"].(float64); ok {
				lineTotal := q * up
				subtotal += lineTotal
				if tr, ok := item["tax_rate"].(float64); ok {
					taxAmount := lineTotal * (tr / 100)
					totalTax += taxAmount
				}
			}
		}
	}

	totalAmount := subtotal - req.DiscountAmount + totalTax

	query := `
		INSERT INTO sales_orders (
			id, tenant_id, order_number, customer_id, order_date, required_by_date,
			delivery_location, delivery_instructions, subtotal_amount, discount_amount,
			tax_amount, total_amount, invoiced_amount, pending_amount, status,
			notes, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := h.DB.Exec(query,
		orderID, tenantID, orderNumber, req.CustomerID, orderDate, requiredByDate,
		req.DeliveryLocation, req.DeliveryInstructions, subtotal, req.DiscountAmount,
		totalTax, totalAmount, 0.0, totalAmount, "draft",
		req.Notes, &userID, now, now,
	)

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	// Insert order items
	itemQuery := `
		INSERT INTO sales_order_items (
			id, tenant_id, order_id, line_number, description, product_service_code,
			ordered_quantity, invoiced_quantity, unit_price, line_total,
			discount_percent, discount_amount, hsn_code, tax_rate, tax_amount,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	for idx, item := range req.Items {
		itemID := uuid.New().String()
		quantity := item["quantity"].(float64)
		unitPrice := item["unit_price"].(float64)
		taxRate := item["tax_rate"].(float64)
		lineTotal := quantity * unitPrice
		taxAmount := lineTotal * (taxRate / 100)

		h.DB.Exec(itemQuery,
			itemID, tenantID, orderID, idx+1,
			item["description"], item["product_service_code"],
			quantity, 0.0, unitPrice, lineTotal,
			item["discount_percent"], item["discount_amount"],
			item["hsn_code"], taxRate, taxAmount,
			now, now,
		)
	}

	h.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":              orderID,
			"order_number":    orderNumber,
			"customer_id":     req.CustomerID,
			"order_date":      orderDate,
			"subtotal_amount": subtotal,
			"discount_amount": req.DiscountAmount,
			"tax_amount":      totalTax,
			"total_amount":    totalAmount,
			"status":          "draft",
		},
	})
}

// GetSalesOrder retrieves a specific order with items
func (h *SalesHandler) GetSalesOrder(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	orderID := mux.Vars(r)["id"]

	query := `
		SELECT id, tenant_id, order_number, customer_id, quotation_id, order_date,
			required_by_date, delivery_location, delivery_instructions,
			subtotal_amount, discount_amount, tax_amount, total_amount, invoiced_amount,
			pending_amount, status, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var order models.SalesOrder
	err := h.DB.QueryRow(query, orderID, tenantID).Scan(
		&order.ID, &order.TenantID, &order.OrderNumber, &order.CustomerID, &order.QuotationID,
		&order.OrderDate, &order.RequiredByDate, &order.DeliveryLocation,
		&order.DeliveryInstructions, &order.SubtotalAmount, &order.DiscountAmount,
		&order.TaxAmount, &order.TotalAmount, &order.InvoicedAmount, &order.PendingAmount,
		&order.Status, &order.Notes, &order.CreatedBy, &order.CreatedAt, &order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		h.respondError(w, http.StatusNotFound, "Order not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve order")
		return
	}

	// Get order items
	itemQuery := `
		SELECT id, tenant_id, order_id, line_number, description, product_service_code,
			ordered_quantity, invoiced_quantity, unit_price, line_total,
			discount_percent, discount_amount, hsn_code, tax_rate, tax_amount,
			created_at, updated_at
		FROM sales_order_items
		WHERE order_id = ? AND tenant_id = ?
		ORDER BY line_number
	`

	rows, _ := h.DB.Query(itemQuery, orderID, tenantID)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var item models.SalesOrderItem
			rows.Scan(&item.ID, &item.TenantID, &item.OrderID, &item.LineNumber,
				&item.Description, &item.ProductServiceCode, &item.OrderedQuantity,
				&item.InvoicedQuantity, &item.UnitPrice, &item.LineTotal,
				&item.DiscountPercent, &item.DiscountAmount, &item.HSNCode,
				&item.TaxRate, &item.TaxAmount, &item.CreatedAt, &item.UpdatedAt)
			order.Items = append(order.Items, item)
		}
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": order})
}

// ListSalesOrders retrieves orders with filtering
func (h *SalesHandler) ListSalesOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	status := r.URL.Query().Get("status")
	customerID := r.URL.Query().Get("customer_id")

	query := `
		SELECT id, tenant_id, order_number, customer_id, quotation_id, order_date,
			required_by_date, delivery_location, delivery_instructions,
			subtotal_amount, discount_amount, tax_amount, total_amount, invoiced_amount,
			pending_amount, status, notes, created_by, created_at, updated_at
		FROM sales_orders
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	args := []interface{}{tenantID}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if customerID != "" {
		query += " AND customer_id = ?"
		args = append(args, customerID)
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve orders")
		return
	}
	defer rows.Close()

	var orders []models.SalesOrder
	for rows.Next() {
		var o models.SalesOrder
		rows.Scan(&o.ID, &o.TenantID, &o.OrderNumber, &o.CustomerID, &o.QuotationID,
			&o.OrderDate, &o.RequiredByDate, &o.DeliveryLocation, &o.DeliveryInstructions,
			&o.SubtotalAmount, &o.DiscountAmount, &o.TaxAmount, &o.TotalAmount,
			&o.InvoicedAmount, &o.PendingAmount, &o.Status, &o.Notes,
			&o.CreatedBy, &o.CreatedAt, &o.UpdatedAt)
		orders = append(orders, o)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": orders})
}

// UpdateSalesOrderStatus updates order status
func (h *SalesHandler) UpdateSalesOrderStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		h.respondError(w, http.StatusBadRequest, "X-Tenant-ID header required")
		return
	}

	orderID := mux.Vars(r)["id"]

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	validStatuses := map[string]bool{
		"draft": true, "confirmed": true, "partially_invoiced": true,
		"invoiced": true, "partially_delivered": true, "delivered": true, "cancelled": true,
	}

	if !validStatuses[req.Status] {
		h.respondError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	query := `
		UPDATE sales_orders
		SET status = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := h.DB.Exec(query, req.Status, time.Now(), orderID, tenantID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.respondError(w, http.StatusNotFound, "Order not found")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{"success": true, "message": "Order status updated"})
}
