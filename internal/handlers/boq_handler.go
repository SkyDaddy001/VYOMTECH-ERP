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

	"vyomtech-backend/internal/services"
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
