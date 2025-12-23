package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// DocumentHandler handles document HTTP requests
type DocumentHandler struct {
	Service *services.DocumentService
}

// NewDocumentHandler creates a new document handler
func NewDocumentHandler(service *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{Service: service}
}

// CreateDocument creates a new document
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	userID := getDocUserID(r)

	var req models.CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	document, err := h.Service.CreateDocument(tenantID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, document)
}

// GetDocument gets a document by ID
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	documentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	document, err := h.Service.GetDocument(tenantID, documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, document)
}

// ListDocuments lists documents for an entity
func (h *DocumentHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	entityType := r.URL.Query().Get("entity_type")
	entityIDStr := r.URL.Query().Get("entity_id")
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	offset := getDocOffsetFromQuery(r)
	limit := getDocLimitFromQuery(r)

	documents, total, err := h.Service.ListDocuments(tenantID, entityType, entityID, offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithDocJSON(w, http.StatusOK, documents, total)
}

// UpdateDocument updates a document
func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	userID := getDocUserID(r)
	vars := mux.Vars(r)
	documentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req models.UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	document, err := h.Service.UpdateDocument(tenantID, documentID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, document)
}

// DeleteDocument deletes a document
func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	userID := getDocUserID(r)
	vars := mux.Vars(r)
	documentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	err = h.Service.DeleteDocument(tenantID, documentID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// VerifyDocument verifies a document
func (h *DocumentHandler) VerifyDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	userID := getDocUserID(r)
	vars := mux.Vars(r)
	documentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var req models.VerifyDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	document, err := h.Service.VerifyDocument(tenantID, documentID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, document)
}

// GetDocumentsByType gets documents by type
func (h *DocumentHandler) GetDocumentsByType(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	docTypeID, err := strconv.ParseInt(vars["type_id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid document type ID")
		return
	}

	offset := getDocOffsetFromQuery(r)
	limit := getDocLimitFromQuery(r)

	documents, total, err := h.Service.GetDocumentsByType(tenantID, docTypeID, offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithDocJSON(w, http.StatusOK, documents, total)
}

// CheckDocumentExpiry checks for expiring documents
func (h *DocumentHandler) CheckDocumentExpiry(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		d, _ := strconv.Atoi(daysStr)
		days = d
	}

	documents, err := h.Service.CheckDocumentExpiry(tenantID, days)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, documents)
}

// CreateDocumentCollection creates a collection
func (h *DocumentHandler) CreateDocumentCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	userID := getDocUserID(r)

	var req models.CreateDocumentCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	collection, err := h.Service.CreateDocumentCollection(tenantID, &req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, collection)
}

// GetDocumentCollection gets a collection
func (h *DocumentHandler) GetDocumentCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	collectionID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid collection ID")
		return
	}

	collection, err := h.Service.GetDocumentCollection(tenantID, collectionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, collection)
}

// ListDocumentCollections lists collections
func (h *DocumentHandler) ListDocumentCollections(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	entityType := r.URL.Query().Get("entity_type")
	entityIDStr := r.URL.Query().Get("entity_id")
	entityID, err := strconv.ParseInt(entityIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid entity ID")
		return
	}

	collections, err := h.Service.ListDocumentCollections(tenantID, entityType, entityID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, collections)
}

// AddDocumentToCollection adds a document to collection
func (h *DocumentHandler) AddDocumentToCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	collectionID, _ := strconv.ParseInt(vars["collection_id"], 10, 64)

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	docID, _ := req["document_id"].(float64)
	isMandatory, _ := req["is_mandatory"].(bool)

	item, err := h.Service.AddDocumentToCollection(tenantID, collectionID, int64(docID), isMandatory)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, item)
}

// RemoveDocumentFromCollection removes a document from collection
func (h *DocumentHandler) RemoveDocumentFromCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	collectionID, _ := strconv.ParseInt(vars["collection_id"], 10, 64)
	documentID, _ := strconv.ParseInt(vars["document_id"], 10, 64)

	err := h.Service.RemoveDocumentFromCollection(tenantID, collectionID, documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// UpdateCollectionStatus updates collection status
func (h *DocumentHandler) UpdateCollectionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	collectionID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	status, ok := req["status"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Status field required")
		return
	}

	collection, err := h.Service.UpdateCollectionStatus(tenantID, collectionID, status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, collection)
}

// ShareDocument shares a document
func (h *DocumentHandler) ShareDocument(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	sharedBy := getDocUserID(r)
	vars := mux.Vars(r)
	documentID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := req["user_id"].(float64)
	permission, _ := req["permission"].(string)

	share, err := h.Service.ShareDocument(tenantID, documentID, int64(userID), permission, sharedBy)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, share)
}

// GetDocumentShares gets document shares
func (h *DocumentHandler) GetDocumentShares(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	documentID, _ := strconv.ParseInt(vars["id"], 10, 64)

	shares, err := h.Service.GetDocumentShares(tenantID, documentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, shares)
}

// RevokeDocumentShare revokes sharing
func (h *DocumentHandler) RevokeDocumentShare(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	shareID, _ := strconv.ParseInt(vars["share_id"], 10, 64)

	err := h.Service.RevokeDocumentShare(tenantID, shareID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// GetDocumentSummary gets document summary
func (h *DocumentHandler) GetDocumentSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	entityType := r.URL.Query().Get("entity_type")
	entityIDStr := r.URL.Query().Get("entity_id")
	entityID, _ := strconv.ParseInt(entityIDStr, 10, 64)

	summary, err := h.Service.GetDocumentSummary(tenantID, entityType, entityID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, summary)
}

// CreateDocumentTemplate creates a template
func (h *DocumentHandler) CreateDocumentTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	userID := getDocUserID(r)

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	template, err := h.Service.CreateDocumentTemplate(tenantID, req["name"], req["content"], req["type"], userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, template)
}

// ListDocumentTemplates lists templates
func (h *DocumentHandler) ListDocumentTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)

	templates, err := h.Service.ListDocumentTemplates(tenantID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, templates)
}

// ListDocumentCategories lists categories
func (h *DocumentHandler) ListDocumentCategories(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)

	categories, err := h.Service.ListDocumentCategories(tenantID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}

// ListDocumentTypes lists types for category
func (h *DocumentHandler) ListDocumentTypes(w http.ResponseWriter, r *http.Request) {
	tenantID := getDocTenantID(r)
	vars := mux.Vars(r)
	categoryID, _ := strconv.ParseInt(vars["category_id"], 10, 64)

	types, err := h.Service.ListDocumentTypes(tenantID, categoryID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, types)
}

// Helper functions
func getDocTenantID(r *http.Request) int64 {
	tenantIDStr, ok := r.Context().Value("tenant_id").(string)
	if ok {
		tenantID, _ := strconv.ParseInt(tenantIDStr, 10, 64)
		return tenantID
	}
	if tenantID, ok := r.Context().Value("tenant_id").(int64); ok {
		return tenantID
	}
	return 0
}

func getDocUserID(r *http.Request) *int64 {
	if userID, ok := r.Context().Value("user_id").(int64); ok {
		return &userID
	}
	return nil
}

func respondWithDocJSON(w http.ResponseWriter, code int, payload interface{}, total int64) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func getDocOffsetFromQuery(r *http.Request) int {
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		return 0
	}
	offset, _ := strconv.Atoi(offsetStr)
	return offset
}

func getDocLimitFromQuery(r *http.Request) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return 10
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit > 100 {
		limit = 100
	}
	return limit
}
