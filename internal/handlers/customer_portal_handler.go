package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
)

// CustomerPortalHandler handles customer portal operations
type CustomerPortalHandler struct {
	service *services.CustomerPortalService
}

// NewCustomerPortalHandler creates a new customer portal handler
func NewCustomerPortalHandler(service *services.CustomerPortalService) *CustomerPortalHandler {
	return &CustomerPortalHandler{service: service}
}

// Helper functions
func getPortalTenantID(r *http.Request) (int64, error) {
	tenantID := r.Context().Value("tenant_id")
	if tenantID == nil {
		return 0, errors.New("tenant_id not found in context")
	}
	return tenantID.(int64), nil
}

func getPortalUserID(r *http.Request) (int64, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0, errors.New("user_id not found in context")
	}
	return userID.(int64), nil
}

func respondWithPortalJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getPortalPagination(r *http.Request) (int, int) {
	limit := 10
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

// CreateProfile creates a customer profile
func (h *CustomerPortalHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.CreateCustomerPortalProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	profile, err := h.service.CreateCustomerProfile(tenantID, userID, &req)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusCreated, profile)
}

// GetProfile retrieves customer profile
func (h *CustomerPortalHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	profile, err := h.service.GetCustomerProfile(tenantID, userID)
	if err != nil {
		if err.Error() == "profile not found" {
			respondWithPortalJSON(w, http.StatusNotFound, map[string]string{"error": "Profile not found"})
			return
		}
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusOK, profile)
}

// UpdateProfile updates customer profile
func (h *CustomerPortalHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.UpdateCustomerPortalProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdateCustomerProfile(tenantID, userID, &req); err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

// GetNotifications retrieves notifications
func (h *CustomerPortalHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	limit, offset := getPortalPagination(r)

	notifications, total, err := h.service.GetNotifications(tenantID, userID, limit, offset)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPortalJSON(w, http.StatusOK, map[string]interface{}{
		"data":  notifications,
		"total": total,
	})
}

// MarkNotificationAsRead marks notification as read
func (h *CustomerPortalHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	notificationID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid notification ID"})
		return
	}

	if err := h.service.MarkNotificationAsRead(tenantID, notificationID); err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

// CreateConversation creates a support conversation
func (h *CustomerPortalHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.CreatePortalConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.ConversationType == "" {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "conversation_type is required"})
		return
	}

	conversation, err := h.service.CreateConversation(tenantID, userID, req.Subject, req.ConversationType)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusCreated, conversation)
}

// SendMessage sends a message
func (h *CustomerPortalHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	conversationID, err := strconv.ParseInt(r.PathValue("conversation_id"), 10, 64)
	if err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation ID"})
		return
	}

	var req models.SendPortalMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.MessageText == "" {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "message_text is required"})
		return
	}

	message, err := h.service.SendMessage(tenantID, conversationID, userID, req.MessageText, "customer")
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusCreated, message)
}

// GetConversationMessages retrieves conversation messages
func (h *CustomerPortalHandler) GetConversationMessages(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	conversationID, err := strconv.ParseInt(r.PathValue("conversation_id"), 10, 64)
	if err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid conversation ID"})
		return
	}

	limit, offset := getPortalPagination(r)

	messages, total, err := h.service.GetConversationMessages(tenantID, conversationID, limit, offset)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPortalJSON(w, http.StatusOK, map[string]interface{}{
		"data":  messages,
		"total": total,
	})
}

// UploadDocument uploads a document
func (h *CustomerPortalHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.CreatePortalDocumentUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.DocumentName == "" {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "document_name is required"})
		return
	}

	document, err := h.service.CreateDocumentUpload(tenantID, userID, &req)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusCreated, document)
}

// GetBookingTracking retrieves booking tracking
func (h *CustomerPortalHandler) GetBookingTracking(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	bookingID, err := strconv.ParseInt(r.PathValue("booking_id"), 10, 64)
	if err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid booking ID"})
		return
	}

	tracking, err := h.service.GetCustomerBookingTracking(tenantID, bookingID)
	if err != nil {
		if err.Error() == "booking tracking not found" {
			respondWithPortalJSON(w, http.StatusNotFound, map[string]string{"error": "Booking tracking not found"})
			return
		}
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusOK, tracking)
}

// GetPayments retrieves payment tracking
func (h *CustomerPortalHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	bookingID, err := strconv.ParseInt(r.PathValue("booking_id"), 10, 64)
	if err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid booking ID"})
		return
	}

	limit, offset := getPortalPagination(r)

	payments, total, err := h.service.GetPaymentTracking(tenantID, bookingID, limit, offset)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	respondWithPortalJSON(w, http.StatusOK, map[string]interface{}{
		"data":  payments,
		"total": total,
	})
}

// CreateFeedback creates feedback
func (h *CustomerPortalHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.CreatePortalFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.FeedbackType == "" || req.Message == "" {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "feedback_type and message are required"})
		return
	}

	feedback, err := h.service.CreateFeedback(tenantID, userID, &req)
	if err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusCreated, feedback)
}

// UpdatePreferences updates customer preferences
func (h *CustomerPortalHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	tenantID, err := getPortalTenantID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	userID, err := getPortalUserID(r)
	if err != nil {
		respondWithPortalJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	var req models.UpdatePortalPreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithPortalJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.service.UpdatePreferences(tenantID, userID, &req); err != nil {
		respondWithPortalJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithPortalJSON(w, http.StatusOK, map[string]string{"message": "Preferences updated successfully"})
}
