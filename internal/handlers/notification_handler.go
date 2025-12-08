package handlers

import (
	"vyomtech-backend/internal/middleware"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"

	"github.com/gorilla/mux"
)

// NotificationHandler handles notification management endpoints
type NotificationHandler struct {
	notificationService services.NotificationService
	log                 *logger.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService services.NotificationService, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		log:                 log,
	}
}

// RegisterRoutes registers notification routes
func (h *NotificationHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/notifications", h.CreateNotification).Methods(http.MethodPost)
	router.HandleFunc("/notifications", h.ListNotifications).Methods(http.MethodGet)
	router.HandleFunc("/notifications/{id}", h.GetNotification).Methods(http.MethodGet)
	router.HandleFunc("/notifications/{id}", h.DeleteNotification).Methods(http.MethodDelete)
	router.HandleFunc("/notifications/{id}/read", h.MarkAsRead).Methods(http.MethodPost)
	router.HandleFunc("/notifications/{id}/archive", h.ArchiveNotification).Methods(http.MethodPost)
	router.HandleFunc("/notifications/user/{userID}/unread", h.GetUnreadNotifications).Methods(http.MethodGet)
	router.HandleFunc("/notifications/stats", h.GetNotificationStats).Methods(http.MethodGet)
	router.HandleFunc("/notifications/preferences", h.GetPreferences).Methods(http.MethodGet)
	router.HandleFunc("/notifications/preferences", h.UpdatePreferences).Methods(http.MethodPut)
}

// CreateNotification creates a new notification
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	var notif services.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	notif.TenantID = tenantID
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if ok && notif.UserID == 0 {
		notif.UserID = userID
	}

	created, err := h.notificationService.CreateNotification(ctx, tenantID, &notif)
	if err != nil {
		h.log.Error("Failed to create notification", "error", err)
		http.Error(w, "failed to create notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetNotification retrieves a specific notification
func (h *NotificationHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	notifID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid notification id", http.StatusBadRequest)
		return
	}

	notif, err := h.notificationService.GetNotification(ctx, tenantID, notifID)
	if err != nil {
		h.log.Error("Failed to get notification", "error", err)
		http.Error(w, "notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notif)
}

// ListNotifications lists notifications for current user
func (h *NotificationHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	limit := int64(20)
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := int64(0)
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.ParseInt(o, 10, 64); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	notifs, err := h.notificationService.GetUserNotifications(ctx, tenantID, userID, int(limit), int(offset))
	if err != nil {
		h.log.Error("Failed to list notifications", "error", err)
		http.Error(w, "failed to list notifications", http.StatusInternalServerError)
		return
	}

	if notifs == nil {
		notifs = []services.Notification{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(notifs),
		"notifs": notifs,
	})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	notifID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid notification id", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.DeleteNotification(ctx, tenantID, notifID); err != nil {
		h.log.Error("Failed to delete notification", "error", err)
		http.Error(w, "failed to delete notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "notification deleted"})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	notifID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid notification id", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.MarkAsRead(ctx, tenantID, notifID); err != nil {
		h.log.Error("Failed to mark notification as read", "error", err)
		http.Error(w, "failed to update notification", http.StatusInternalServerError)
		return
	}

	// Fetch updated notification
	notif, err := h.notificationService.GetNotification(ctx, tenantID, notifID)
	if err != nil {
		h.log.Error("Failed to fetch notification", "error", err)
		http.Error(w, "failed to fetch notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notif)
}

// ArchiveNotification archives a notification
func (h *NotificationHandler) ArchiveNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	notifID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid notification id", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.ArchiveNotification(ctx, tenantID, notifID); err != nil {
		h.log.Error("Failed to archive notification", "error", err)
		http.Error(w, "failed to update notification", http.StatusInternalServerError)
		return
	}

	// Fetch updated notification
	notif, err := h.notificationService.GetNotification(ctx, tenantID, notifID)
	if err != nil {
		h.log.Error("Failed to fetch notification", "error", err)
		http.Error(w, "failed to fetch notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notif)
}

// GetUnreadNotifications gets unread notifications for a user
func (h *NotificationHandler) GetUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userID"], 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	notifs, err := h.notificationService.GetUserNotifications(ctx, tenantID, userID, 100, 0)
	if err != nil {
		h.log.Error("Failed to get unread notifications", "error", err)
		http.Error(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	// Filter for unread only
	unread := []services.Notification{}
	for _, notif := range notifs {
		if !notif.IsRead {
			unread = append(unread, notif)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count":  len(unread),
		"notifs": unread,
	})
}

// GetNotificationStats gets notification statistics
func (h *NotificationHandler) GetNotificationStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		userID = 0
	}

	stats, err := h.notificationService.GetNotificationStats(ctx, tenantID, userID)
	if err != nil {
		h.log.Error("Failed to get notification stats", "error", err)
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetPreferences gets notification preferences for current user
func (h *NotificationHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	prefs, err := h.notificationService.GetPreferences(ctx, tenantID, userID)
	if err != nil {
		h.log.Error("Failed to get preferences", "error", err)
		http.Error(w, "failed to get preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prefs)
}

// UpdatePreferences updates notification preferences for current user
func (h *NotificationHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant_id required", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	var prefs services.NotificationPreferences
	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.notificationService.UpdatePreferences(ctx, tenantID, userID, &prefs); err != nil {
		h.log.Error("Failed to update preferences", "error", err)
		http.Error(w, "failed to update preferences", http.StatusInternalServerError)
		return
	}

	// Fetch updated preferences
	updated, err := h.notificationService.GetPreferences(ctx, tenantID, userID)
	if err != nil {
		h.log.Error("Failed to fetch preferences", "error", err)
		http.Error(w, "failed to fetch preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
