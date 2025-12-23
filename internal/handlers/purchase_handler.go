package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/constants"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PurchaseHandler struct {
	DB          *sql.DB
	RBACService *services.RBACService
}

// NewPurchaseHandler creates a new purchase handler
func NewPurchaseHandler(db *sql.DB, rbacService *services.RBACService) *PurchaseHandler {
	return &PurchaseHandler{
		DB:          db,
		RBACService: rbacService,
	}
}

// ============================================================================
// VENDOR MANAGEMENT HANDLERS
// ============================================================================

// CreateVendor - POST /api/v1/purchase/vendors
func (h *PurchaseHandler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.VendorCreate); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	var vendor models.Vendor
	if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	vendor.ID = uuid.New().String()
	vendor.TenantID = tenant
	vendor.Status = "active"
	vendor.CreatedAt = time.Now()
	vendor.UpdatedAt = time.Now()

	query := `INSERT INTO vendors (id, tenant_id, name, email, phone, address, city, state, country, zip_code, vendor_type, payment_terms, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := h.DB.Exec(query, vendor.ID, vendor.TenantID, vendor.Name, vendor.Email, vendor.Phone, vendor.Address, vendor.City, vendor.State, vendor.Country, vendor.PostalCode, vendor.VendorType, vendor.PaymentTerms, vendor.Status, vendor.CreatedAt, vendor.UpdatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vendor)
}

// ListVendors - GET /api/v1/purchase/vendors
func (h *PurchaseHandler) ListVendors(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	if pageNum == 0 {
		pageNum = 1
	}
	if limitNum == 0 {
		limitNum = 20
	}

	offset := (pageNum - 1) * limitNum

	query := `SELECT id, tenant_id, name, email, phone, address, city, state, country, zip_code, vendor_type, payment_terms, status, created_at, updated_at 
	          FROM vendors WHERE tenant_id = ? AND status = 'active' LIMIT ? OFFSET ?`

	rows, err := h.DB.Query(query, tenant, limitNum, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vendors []models.Vendor
	for rows.Next() {
		var vendor models.Vendor
		if err := rows.Scan(&vendor.ID, &vendor.TenantID, &vendor.Name, &vendor.Email, &vendor.Phone, &vendor.Address, &vendor.City, &vendor.State, &vendor.Country, &vendor.PostalCode, &vendor.VendorType, &vendor.PaymentTerms, &vendor.Status, &vendor.CreatedAt, &vendor.UpdatedAt); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		vendors = append(vendors, vendor)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendors)
}

// GetVendor - GET /api/v1/purchase/vendors/{id}
func (h *PurchaseHandler) GetVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	query := `SELECT id, tenant_id, name, email, phone, address, city, state, country, zip_code, vendor_type, payment_terms, status, created_at, updated_at 
	          FROM vendors WHERE id = ? AND tenant_id = ?`

	var vendor models.Vendor
	err := h.DB.QueryRow(query, vendorID, tenant).Scan(&vendor.ID, &vendor.TenantID, &vendor.Name, &vendor.Email, &vendor.Phone, &vendor.Address, &vendor.City, &vendor.State, &vendor.Country, &vendor.PostalCode, &vendor.VendorType, &vendor.PaymentTerms, &vendor.Status, &vendor.CreatedAt, &vendor.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Vendor not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

// UpdateVendor - PUT /api/v1/purchase/vendors/{id}
func (h *PurchaseHandler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	var vendor models.Vendor
	if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	vendor.UpdatedAt = time.Now()

	query := `UPDATE vendors SET name = ?, email = ?, phone = ?, address = ?, city = ?, state = ?, country = ?, zip_code = ?, vendor_type = ?, payment_terms = ?, updated_at = ? 
	          WHERE id = ? AND tenant_id = ?`

	result, err := h.DB.Exec(query, vendor.Name, vendor.Email, vendor.Phone, vendor.Address, vendor.City, vendor.State, vendor.Country, vendor.PostalCode, vendor.VendorType, vendor.PaymentTerms, vendor.UpdatedAt, vendorID, tenant)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Vendor not found"}`, http.StatusNotFound)
		return
	}

	vendor.ID = vendorID
	vendor.TenantID = tenant

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vendor)
}

// DeleteVendor - DELETE /api/v1/purchase/vendors/{id}
func (h *PurchaseHandler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vendorID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	query := `UPDATE vendors SET status = 'inactive' WHERE id = ? AND tenant_id = ?`
	result, err := h.DB.Exec(query, vendorID, tenant)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Vendor not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// ============================================================================
// PURCHASE ORDER HANDLERS
// ============================================================================

// CreatePurchaseOrder - POST /api/v1/purchase/orders
func (h *PurchaseHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var po models.PurchaseOrder
	if err := json.NewDecoder(r.Body).Decode(&po); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("X-Tenant-ID")
	po.ID = uuid.New().String()
	po.TenantID = tenant
	po.PONumber = fmt.Sprintf("PO-%d-%s", time.Now().Unix(), po.ID[:8])
	po.Status = "draft"
	po.CreatedAt = time.Now()
	po.UpdatedAt = time.Now()

	query := `INSERT INTO purchase_orders (id, tenant_id, po_number, vendor_id, total_amount, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := h.DB.Exec(query, po.ID, po.TenantID, po.PONumber, po.VendorID, po.TotalAmount, po.Status, po.CreatedAt, po.UpdatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(po)
}

// ListPurchaseOrders - GET /api/v1/purchase/orders
func (h *PurchaseHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	if pageNum == 0 {
		pageNum = 1
	}
	if limitNum == 0 {
		limitNum = 20
	}

	offset := (pageNum - 1) * limitNum

	query := `SELECT id, tenant_id, po_number, vendor_id, total_amount, status, created_at, updated_at 
	          FROM purchase_orders WHERE tenant_id = ? LIMIT ? OFFSET ?`

	rows, err := h.DB.Query(query, tenant, limitNum, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.PurchaseOrder
	for rows.Next() {
		var po models.PurchaseOrder
		if err := rows.Scan(&po.ID, &po.TenantID, &po.PONumber, &po.VendorID, &po.TotalAmount, &po.Status, &po.CreatedAt, &po.UpdatedAt); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		orders = append(orders, po)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetPurchaseOrder - GET /api/v1/purchase/orders/{id}
func (h *PurchaseHandler) GetPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	poID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	query := `SELECT id, tenant_id, po_number, vendor_id, total_amount, status, created_at, updated_at 
	          FROM purchase_orders WHERE id = ? AND tenant_id = ?`

	var po models.PurchaseOrder
	err := h.DB.QueryRow(query, poID, tenant).Scan(&po.ID, &po.TenantID, &po.PONumber, &po.VendorID, &po.TotalAmount, &po.Status, &po.CreatedAt, &po.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Purchase order not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(po)
}

// UpdatePurchaseOrder - PUT /api/v1/purchase/orders/{id}
func (h *PurchaseHandler) UpdatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	poID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	var po models.PurchaseOrder
	if err := json.NewDecoder(r.Body).Decode(&po); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	po.UpdatedAt = time.Now()

	query := `UPDATE purchase_orders SET vendor_id = ?, total_amount = ?, status = ?, updated_at = ? 
	          WHERE id = ? AND tenant_id = ?`

	result, err := h.DB.Exec(query, po.VendorID, po.TotalAmount, po.Status, po.UpdatedAt, poID, tenant)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Purchase order not found"}`, http.StatusNotFound)
		return
	}

	po.ID = poID
	po.TenantID = tenant

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(po)
}

// ============================================================================
// GOODS RECEIPT NOTE (GRN) HANDLERS
// ============================================================================

// CreateGRN - POST /api/v1/purchase/grn
func (h *PurchaseHandler) CreateGRN(w http.ResponseWriter, r *http.Request) {
	var grn models.GoodsReceipt
	if err := json.NewDecoder(r.Body).Decode(&grn); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("X-Tenant-ID")
	grn.ID = uuid.New().String()
	grn.TenantID = tenant
	grn.GRNNumber = fmt.Sprintf("GRN-%d-%s", time.Now().Unix(), grn.ID[:8])
	grn.Status = "pending"
	grn.CreatedAt = time.Now()
	grn.UpdatedAt = time.Now()

	query := `INSERT INTO goods_receipts (id, tenant_id, grn_number, po_id, received_date, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := h.DB.Exec(query, grn.ID, grn.TenantID, grn.GRNNumber, grn.POID, grn.ReceiptDate, grn.Status, grn.CreatedAt, grn.UpdatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(grn)
}

// ListGRNs - GET /api/v1/purchase/grn
func (h *PurchaseHandler) ListGRNs(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")

	query := `SELECT id, tenant_id, grn_number, po_id, received_date, status, created_at, updated_at 
	          FROM goods_receipts WHERE tenant_id = ?`

	rows, err := h.DB.Query(query, tenant)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var grns []models.GoodsReceipt
	for rows.Next() {
		var grn models.GoodsReceipt
		if err := rows.Scan(&grn.ID, &grn.TenantID, &grn.GRNNumber, &grn.POID, &grn.ReceiptDate, &grn.Status, &grn.CreatedAt, &grn.UpdatedAt); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		grns = append(grns, grn)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grns)
}

// ============================================================================
// NOTE: Invoice functionality is part of the Billing Module
// See internal/handlers/billing_handler.go and internal/models/billing.go
// ============================================================================
