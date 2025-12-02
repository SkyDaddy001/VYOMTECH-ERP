#!/bin/bash

# BOQ Service Setup Script
# Creates BOQ import/export functionality with INR currency support

set -e

PROJECT_DIR="/c/Users/Skydaddy/Desktop/VYOM - ERP"
SERVICES_DIR="$PROJECT_DIR/internal/services"
HANDLERS_DIR="$PROJECT_DIR/internal/handlers"

echo "=== BOQ Service Setup ==="
echo "Creating BOQ Service..."

# Create boq_service.go with sql.DB support
cat > "$SERVICES_DIR/boq_service.go" << 'EOF'
package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"multi-tenant-ai-callcenter/internal/models"
)

type BOQService struct {
	DB *sql.DB
}

// NewBOQService creates a new BOQ service
func NewBOQService(db *sql.DB) *BOQService {
	return &BOQService{DB: db}
}

// BOQImportResult contains results of BOQ import
type BOQImportResult struct {
	TotalRows      int
	SuccessCount   int
	FailureCount   int
	Errors         []string
	CreatedBOQs    int
	UpdatedBOQs    int
	DuplicateCount int
	TotalAmountINR float64 // Total amount in Indian Rupees
}

// ImportBOQFromExcel imports BOQ items from Excel file
func (s *BOQService) ImportBOQFromExcel(tenantID string, projectID int64, filePath string) (*BOQImportResult, error) {
	result := &BOQImportResult{
		Errors: []string{},
	}

	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return result, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	// Get first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return result, fmt.Errorf("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return result, fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(rows) < 2 {
		return result, fmt.Errorf("excel file must have header row and data rows")
	}

	headerMap := parseHeaders(rows[0])
	result.TotalRows = len(rows) - 1

	// Parse and insert BOQ items
	for idx := 1; idx < len(rows); idx++ {
		row := rows[idx]
		if len(row) == 0 || row[0] == "" {
			continue
		}

		boqItem, err := parseBOQRow(row, headerMap, tenantID, projectID)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: %v", idx+1, err))
			continue
		}

		// Check if BOQ already exists
		var id int64
		checkErr := s.DB.QueryRow(
			"SELECT id FROM bill_of_quantities WHERE tenant_id = ? AND project_id = ? AND item_description = ? AND deleted_at IS NULL",
			tenantID, projectID, boqItem.ItemDescription,
		).Scan(&id)

		if checkErr == nil {
			// Update existing
			_, err = s.DB.Exec(
				"UPDATE bill_of_quantities SET unit = ?, quantity = ?, unit_rate = ?, total_amount = ?, category = ?, status = ?, updated_at = NOW() WHERE id = ? AND tenant_id = ?",
				boqItem.Unit, boqItem.Quantity, boqItem.UnitRate, boqItem.TotalAmount, boqItem.Category, boqItem.Status, id, tenantID,
			)
			if err != nil {
				result.FailureCount++
				result.Errors = append(result.Errors, fmt.Sprintf("Row %d: failed to update: %v", idx+1, err))
				continue
			}
			result.UpdatedBOQs++
		} else if checkErr == sql.ErrNoRows {
			// Create new
			_, err = s.DB.Exec(
				"INSERT INTO bill_of_quantities (tenant_id, project_id, boq_number, item_description, unit, quantity, unit_rate, total_amount, category, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())",
				tenantID, projectID, boqItem.BOQNumber, boqItem.ItemDescription, boqItem.Unit, boqItem.Quantity, boqItem.UnitRate, boqItem.TotalAmount, boqItem.Category, boqItem.Status,
			)
			if err != nil {
				result.FailureCount++
				result.Errors = append(result.Errors, fmt.Sprintf("Row %d: failed to create: %v", idx+1, err))
				continue
			}
			result.CreatedBOQs++
		} else {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: database error: %v", idx+1, checkErr))
			continue
		}

		result.SuccessCount++
	}

	// Calculate total amount in INR
	var totalAmount sql.NullFloat64
	s.DB.QueryRow(
		"SELECT COALESCE(SUM(total_amount), 0) FROM bill_of_quantities WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL",
		tenantID, projectID,
	).Scan(&totalAmount)
	if totalAmount.Valid {
		result.TotalAmountINR = totalAmount.Float64
	}

	return result, nil
}

// ExportBOQToExcel exports BOQ items to Excel file
func (s *BOQService) ExportBOQToExcel(tenantID string, projectID int64, outputPath string) error {
	// Fetch BOQ items
	rows, err := s.DB.Query(
		"SELECT boq_number, item_description, unit, quantity, unit_rate, total_amount, category, status, created_at FROM bill_of_quantities WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL ORDER BY created_at",
		tenantID, projectID,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch BOQ items: %w", err)
	}
	defer rows.Close()

	// Create Excel file
	f := excelize.NewFile()
	headers := []string{"BOQ Number", "Item Description", "Unit", "Quantity", "Unit Rate (₹)", "Total Amount (₹)", "Category", "Status", "Created Date"}

	// Write headers
	for col, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+rune(col), 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Write data
	rowNum := 2
	var totalAmount float64
	for rows.Next() {
		var boqNumber, itemDesc, unit, category, status string
		var qty, rate, amount float64
		var createdAt time.Time

		err := rows.Scan(&boqNumber, &itemDesc, &unit, &qty, &rate, &amount, &category, &status, &createdAt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		values := []interface{}{boqNumber, itemDesc, unit, qty, rate, amount, category, status, createdAt.Format("2006-01-02")}
		for col, val := range values {
			cell := fmt.Sprintf("%c%d", 'A'+rune(col), rowNum)
			f.SetCellValue("Sheet1", cell, val)
		}

		totalAmount += amount
		rowNum++
	}

	// Add summary
	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum+1), "Total:")
	f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum+1), totalAmount)

	// Set column widths
	f.SetColWidth("Sheet1", "A", "I", 15)

	// Save file
	if err := f.SaveAs(outputPath); err != nil {
		return fmt.Errorf("failed to save excel file: %w", err)
	}

	return nil
}

// GetBOQItems retrieves BOQ items for a project
func (s *BOQService) GetBOQItems(tenantID string, projectID int64, limit int, offset int) ([]models.BillOfQuantities, int64, error) {
	// Count total
	var total int64
	err := s.DB.QueryRow(
		"SELECT COUNT(*) FROM bill_of_quantities WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL",
		tenantID, projectID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	rows, err := s.DB.Query(
		"SELECT id, tenant_id, project_id, boq_number, item_description, unit, quantity, unit_rate, total_amount, category, status, created_at, updated_at FROM bill_of_quantities WHERE tenant_id = ? AND project_id = ? AND deleted_at IS NULL LIMIT ? OFFSET ?",
		tenantID, projectID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var boqItems []models.BillOfQuantities
	for rows.Next() {
		var boq models.BillOfQuantities
		err := rows.Scan(&boq.ID, &boq.TenantID, &boq.ProjectID, &boq.BOQNumber, &boq.ItemDescription, &boq.Unit, &boq.Quantity, &boq.UnitRate, &boq.TotalAmount, &boq.Category, &boq.Status, &boq.CreatedAt, &boq.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning BOQ: %v", err)
			continue
		}
		boqItems = append(boqItems, boq)
	}

	return boqItems, total, nil
}

// UpdateBOQItem updates a BOQ item
func (s *BOQService) UpdateBOQItem(tenantID string, boqID int64, quantity, unitRate float64) error {
	totalAmount := quantity * unitRate
	_, err := s.DB.Exec(
		"UPDATE bill_of_quantities SET quantity = ?, unit_rate = ?, total_amount = ?, updated_at = NOW() WHERE id = ? AND tenant_id = ?",
		quantity, unitRate, totalAmount, boqID, tenantID,
	)
	return err
}

// DeleteBOQItem soft-deletes a BOQ item
func (s *BOQService) DeleteBOQItem(tenantID string, boqID int64) error {
	_, err := s.DB.Exec(
		"UPDATE bill_of_quantities SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?",
		boqID, tenantID,
	)
	return err
}

// Helper functions
func parseHeaders(headerRow []string) map[string]int {
	headerMap := make(map[string]int)
	for idx, header := range headerRow {
		headerMap[header] = idx
	}
	return headerMap
}

func parseBOQRow(row []string, headerMap map[string]int, tenantID string, projectID int64) (*models.BillOfQuantities, error) {
	boq := &models.BillOfQuantities{
		TenantID:  tenantID,
		ProjectID: projectID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "planned",
	}

	for header, idx := range headerMap {
		if idx >= len(row) || row[idx] == "" {
			continue
		}

		switch header {
		case "Item Description", "Description":
			boq.ItemDescription = row[idx]
		case "Unit":
			boq.Unit = row[idx]
		case "Quantity":
			if qty, err := strconv.ParseFloat(row[idx], 64); err == nil {
				boq.Quantity = qty
			}
		case "Unit Rate", "Rate":
			if rate, err := strconv.ParseFloat(row[idx], 64); err == nil {
				boq.UnitRate = rate
			}
		case "Total Amount", "Total":
			if total, err := strconv.ParseFloat(row[idx], 64); err == nil {
				boq.TotalAmount = total
			}
		case "Category":
			boq.Category = row[idx]
		case "Status":
			boq.Status = row[idx]
		case "BOQ Number", "Number":
			boq.BOQNumber = row[idx]
		}
	}

	if boq.ItemDescription == "" {
		return nil, fmt.Errorf("item description is required")
	}

	if boq.TotalAmount == 0 && boq.Quantity > 0 && boq.UnitRate > 0 {
		boq.TotalAmount = boq.Quantity * boq.UnitRate
	}

	return boq, nil
}
EOF

echo "✅ Created boq_service.go"

# Create boq_handler.go
cat > "$HANDLERS_DIR/boq_handler.go" << 'EOF'
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"multi-tenant-ai-callcenter/internal/services"
)

type BOQHandler struct {
	service *services.BOQService
}

func NewBOQHandler(service *services.BOQService) *BOQHandler {
	return &BOQHandler{service: service}
}

func RegisterBOQRoutes(router *mux.Router, service *services.BOQService) {
	handler := NewBOQHandler(service)
	router.HandleFunc("/api/v1/boq/import", handler.ImportBOQ).Methods("POST")
	router.HandleFunc("/api/v1/boq/export", handler.ExportBOQ).Methods("GET")
	router.HandleFunc("/api/v1/boq/list", handler.ListBOQItems).Methods("GET")
	router.HandleFunc("/api/v1/boq/update", handler.UpdateBOQItem).Methods("PUT")
	router.HandleFunc("/api/v1/boq/delete", handler.DeleteBOQItem).Methods("DELETE")
}

type ImportBOQResponse struct {
	TotalRows       int      `json:"total_rows"`
	SuccessCount    int      `json:"success_count"`
	FailureCount    int      `json:"failure_count"`
	CreatedBOQs     int      `json:"created_boqs"`
	UpdatedBOQs     int      `json:"updated_boqs"`
	DuplicateCount  int      `json:"duplicate_count"`
	TotalAmountINR  float64  `json:"total_amount_inr"`
	Currency        string   `json:"currency"`
	Errors          []string `json:"errors,omitempty"`
}

func (h *BOQHandler) ImportBOQ(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID header", http.StatusBadRequest)
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(header.Filename) != ".xlsx" && filepath.Ext(header.Filename) != ".xls" {
		http.Error(w, "Only .xlsx and .xls files are supported", http.StatusBadRequest)
		return
	}

	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fmt.Sprintf("boq_%d_%s.xlsx", projectID, time.Now().Format("20060102150405")))
	tempFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Failed to create temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempPath)

	_, err = io.Copy(tempFile, file)
	tempFile.Close()
	if err != nil {
		http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
		return
	}

	result, err := h.service.ImportBOQFromExcel(tenantID, projectID, tempPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Import failed: %v", err)})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ImportBOQResponse{
		TotalRows:      result.TotalRows,
		SuccessCount:   result.SuccessCount,
		FailureCount:   result.FailureCount,
		CreatedBOQs:    result.CreatedBOQs,
		UpdatedBOQs:    result.UpdatedBOQs,
		DuplicateCount: result.DuplicateCount,
		TotalAmountINR: result.TotalAmountINR,
		Currency:       "INR",
		Errors:         result.Errors,
	})
}

func (h *BOQHandler) ExportBOQ(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID header", http.StatusBadRequest)
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fmt.Sprintf("boq_export_%d_%s.xlsx", projectID, time.Now().Format("20060102150405")))
	defer os.Remove(tempPath)

	err = h.service.ExportBOQToExcel(tenantID, projectID, tempPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Export failed: %v", err), http.StatusInternalServerError)
		return
	}

	fileData, err := os.ReadFile(tempPath)
	if err != nil {
		http.Error(w, "Failed to read exported file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=BOQ_Project_%d.xlsx", projectID))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.WriteHeader(http.StatusOK)
	w.Write(fileData)
}

func (h *BOQHandler) ListBOQItems(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID header", http.StatusBadRequest)
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 500 {
			limit = parsedLimit
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	boqItems, total, err := h.service.GetBOQItems(tenantID, projectID, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch BOQ items: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":  boqItems,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *BOQHandler) UpdateBOQItem(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID header", http.StatusBadRequest)
		return
	}

	boqIDStr := r.URL.Query().Get("id")
	boqID, err := strconv.ParseInt(boqIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid BOQ ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Quantity float64 `json:"quantity"`
		UnitRate float64 `json:"unit_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateBOQItem(tenantID, boqID, req.Quantity, req.UnitRate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "BOQ item updated successfully"})
}

func (h *BOQHandler) DeleteBOQItem(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Missing X-Tenant-ID header", http.StatusBadRequest)
		return
	}

	boqIDStr := r.URL.Query().Get("id")
	boqID, err := strconv.ParseInt(boqIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid BOQ ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteBOQItem(tenantID, boqID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "BOQ item deleted successfully"})
}
EOF

echo "✅ Created boq_handler.go"

# Build
cd "$PROJECT_DIR"
echo "Building..."
go build -o bin/main cmd/main.go 2>&1

if [ $? -eq 0 ]; then
	echo "✅ Build successful!"
	ls -lh bin/main
else
	echo "❌ Build failed"
	exit 1
fi

echo "=== BOQ Setup Complete ==="
echo "Endpoints created:"
echo "  POST   /api/v1/boq/import?project_id=X    - Upload Excel BOQ file (INR)"
echo "  GET    /api/v1/boq/export?project_id=X    - Download Excel BOQ file (INR)"
echo "  GET    /api/v1/boq/list?project_id=X      - List BOQ items with INR values"
echo "  PUT    /api/v1/boq/update?id=X            - Update BOQ item (INR)"
echo "  DELETE /api/v1/boq/delete?id=X            - Delete BOQ item"
echo ""
echo "All monetary values in INR (Indian Rupees) ₹"
