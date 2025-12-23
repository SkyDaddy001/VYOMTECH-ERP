package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// MobileHandler handles mobile app-related HTTP requests
type MobileHandler struct {
	mobileService *services.MobileService
	logger        *logger.Logger
}

// NewMobileHandler creates a new MobileHandler
func NewMobileHandler(mobileService *services.MobileService, logger *logger.Logger) *MobileHandler {
	return &MobileHandler{
		mobileService: mobileService,
		logger:        logger,
	}
}

// ============================================================
// HELPER FUNCTIONS
// ============================================================

func getMobileTenantID(w http.ResponseWriter, r *http.Request) (string, bool) {
	ctx := r.Context()
	tenantID, ok := ctx.Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return "", false
	}
	return tenantID, true
}

func getMobileUserID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return 0, false
	}
	return userID, true
}

func generateUUID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixNano())
}

// ============================================================
// MOBILE APP HANDLERS
// ============================================================

// CreateApp handles POST /api/v1/mobile/apps
func (mh *MobileHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	var req models.CreateMobileAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	app := &models.MobileApp{
		ID:                  generateUUID(),
		TenantID:            tenantID,
		AppName:             req.AppName,
		AppType:             req.AppType,
		BundleIdentifier:    req.BundleIdentifier,
		Version:             req.Version,
		BuildNumber:         req.BuildNumber,
		Status:              "active",
		Description:         req.Description,
		APIKey:              generateUUID(),
		APISecret:           generateUUID(),
		MinSupportedVersion: req.MinSupportedVersion,
		MaxSupportedVersion: req.MaxSupportedVersion,
		StoreURLIOS:         req.StoreURLIOS,
		StoreURLAndroid:     req.StoreURLAndroid,
		SupportEmail:        req.SupportEmail,
		SupportPhone:        req.SupportPhone,
		PrivacyPolicyURL:    req.PrivacyPolicyURL,
		TermsOfServiceURL:   req.TermsOfServiceURL,
	}

	if err := mh.mobileService.CreateApp(ctx, app); err != nil {
		mh.logger.Error("Failed to create app", "error", err)
		http.Error(w, "failed to create app", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

// GetApp handles GET /api/v1/mobile/apps/{id}
func (mh *MobileHandler) GetApp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	appID := r.PathValue("id")
	if appID == "" {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	app, err := mh.mobileService.GetApp(ctx, appID, tenantID)
	if err != nil {
		http.Error(w, "app not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

// ListApps handles GET /api/v1/mobile/apps
func (mh *MobileHandler) ListApps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	apps, err := mh.mobileService.ListApps(ctx, tenantID, limit, offset)
	if err != nil {
		mh.logger.Error("Failed to list apps", "error", err)
		http.Error(w, "failed to list apps", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}

// ============================================================
// DEVICE HANDLERS
// ============================================================

// RegisterDevice handles POST /api/v1/mobile/devices
func (mh *MobileHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	var req models.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	appID := r.URL.Query().Get("app_id")
	if appID == "" {
		http.Error(w, "app_id query parameter required", http.StatusBadRequest)
		return
	}

	device := &models.MobileDevice{
		ID:                 generateUUID(),
		TenantID:           tenantID,
		UserID:             userID,
		AppID:              appID,
		DeviceID:           req.DeviceID,
		DeviceName:         req.DeviceName,
		DeviceModel:        req.DeviceModel,
		DeviceManufacturer: req.DeviceManufacturer,
		OSType:             req.OSType,
		OSVersion:          req.OSVersion,
		AppVersion:         req.AppVersion,
		AppBuild:           req.AppBuild,
		PushToken:          req.PushToken,
		BiometricEnabled:   false,
		BiometricType:      req.BiometricType,
		DeviceStatus:       "active",
	}

	if req.BiometricEnabled != nil {
		device.BiometricEnabled = *req.BiometricEnabled
	}

	if err := mh.mobileService.RegisterDevice(ctx, device); err != nil {
		mh.logger.Error("Failed to register device", "error", err)
		http.Error(w, "failed to register device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

// GetDevice handles GET /api/v1/mobile/devices/{id}
func (mh *MobileHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	deviceID := r.PathValue("id")
	if deviceID == "" {
		http.Error(w, "invalid device id", http.StatusBadRequest)
		return
	}

	device, err := mh.mobileService.GetDevice(ctx, deviceID, tenantID)
	if err != nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// ListUserDevices handles GET /api/v1/mobile/my-devices
func (mh *MobileHandler) ListUserDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	devices, err := mh.mobileService.ListUserDevices(ctx, userID, tenantID)
	if err != nil {
		mh.logger.Error("Failed to list user devices", "error", err)
		http.Error(w, "failed to list devices", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// ============================================================
// SESSION HANDLERS
// ============================================================

// CreateSession handles POST /api/v1/mobile/sessions
func (mh *MobileHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	var req models.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sessionToken := generateUUID()
	refreshToken := generateUUID()
	expiresAt := time.Now().Add(24 * time.Hour)
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	session := &models.MobileSession{
		ID:                    generateUUID(),
		TenantID:              tenantID,
		UserID:                userID,
		DeviceID:              req.DeviceID,
		AppID:                 req.AppID,
		SessionToken:          sessionToken,
		RefreshToken:          &refreshToken,
		TokenExpiresAt:        expiresAt,
		RefreshTokenExpiresAt: &refreshExpiresAt,
		SessionStatus:         "active",
		IPAddress:             req.IPAddress,
		UserAgent:             req.UserAgent,
		LoginMethod:           &req.LoginMethod,
		DeviceTrusted:         false,
		AppInForeground:       true,
	}

	if req.DeviceTrusted != nil {
		session.DeviceTrusted = *req.DeviceTrusted
	}

	if err := mh.mobileService.CreateSession(ctx, session); err != nil {
		mh.logger.Error("Failed to create session", "error", err)
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(session)
}

// ============================================================
// NOTIFICATION HANDLERS
// ============================================================

// SendNotification handles POST /api/v1/mobile/notifications
func (mh *MobileHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	appID := r.URL.Query().Get("app_id")
	if appID == "" {
		http.Error(w, "app_id query parameter required", http.StatusBadRequest)
		return
	}

	var req models.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	notification := &models.AppNotification{
		ID:               generateUUID(),
		TenantID:         tenantID,
		AppID:            appID,
		NotificationType: req.NotificationType,
		Title:            req.Title,
		Body:             req.Body,
		ActionURL:        req.ActionURL,
		CustomData:       req.CustomData,
		DeliveryStatus:   "pending",
		DeliveryAttempts: 0,
		Priority:         req.Priority,
	}

	if req.ScheduledTime != nil {
		notification.ScheduledTime = req.ScheduledTime
	} else {
		now := time.Now()
		notification.SendTime = &now
	}

	if err := mh.mobileService.SendNotification(ctx, notification); err != nil {
		mh.logger.Error("Failed to send notification", "error", err)
		http.Error(w, "failed to send notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

// GetUserNotifications handles GET /api/v1/mobile/notifications
func (mh *MobileHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	notifications, err := mh.mobileService.GetUserNotifications(ctx, userID, tenantID, limit, offset)
	if err != nil {
		mh.logger.Error("Failed to get notifications", "error", err)
		http.Error(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// MarkNotificationAsRead handles PUT /api/v1/mobile/notifications/{id}/read
func (mh *MobileHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	notificationID := r.PathValue("id")
	if notificationID == "" {
		http.Error(w, "invalid notification id", http.StatusBadRequest)
		return
	}

	if err := mh.mobileService.MarkNotificationAsRead(ctx, notificationID, tenantID); err != nil {
		mh.logger.Error("Failed to mark notification as read", "error", err)
		http.Error(w, "failed to mark as read", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "notification marked as read",
	})
}

// ============================================================
// FEATURE FLAG HANDLERS
// ============================================================

// GetAppFeatures handles GET /api/v1/mobile/features
func (mh *MobileHandler) GetAppFeatures(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	appID := r.URL.Query().Get("app_id")
	if appID == "" {
		http.Error(w, "app_id query parameter required", http.StatusBadRequest)
		return
	}

	features, err := mh.mobileService.GetAppFeatures(ctx, appID, tenantID)
	if err != nil {
		mh.logger.Error("Failed to get features", "error", err)
		http.Error(w, "failed to get features", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(features)
}

// ============================================================
// OFFLINE DATA HANDLERS
// ============================================================

// SyncOfflineData handles POST /api/v1/mobile/sync
func (mh *MobileHandler) SyncOfflineData(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	var req models.SyncOfflineDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deviceID := r.URL.Query().Get("device_id")
	appID := r.URL.Query().Get("app_id")
	if deviceID == "" || appID == "" {
		http.Error(w, "device_id and app_id query parameters required", http.StatusBadRequest)
		return
	}

	// Verify device and app exist for this tenant
	_, err := mh.mobileService.GetDevice(r.Context(), deviceID, tenantID)
	if err != nil {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}

	_, err = mh.mobileService.GetApp(r.Context(), appID, tenantID)
	if err != nil {
		http.Error(w, "app not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "offline data synced",
		"status":  "pending",
	})
}

// ============================================================
// CRASH REPORTING HANDLERS
// ============================================================

// ReportCrash handles POST /api/v1/mobile/crashes
func (mh *MobileHandler) ReportCrash(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	var req models.ReportCrashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deviceID := r.URL.Query().Get("device_id")
	appID := r.URL.Query().Get("app_id")
	if deviceID == "" || appID == "" {
		http.Error(w, "device_id and app_id query parameters required", http.StatusBadRequest)
		return
	}

	crash := &models.AppCrash{
		ID:                 generateUUID(),
		TenantID:           tenantID,
		AppID:              appID,
		DeviceID:           deviceID,
		UserID:             &userID,
		CrashTimestamp:     time.Now(),
		CrashType:          req.CrashType,
		ExceptionType:      req.ExceptionType,
		ExceptionMessage:   req.ExceptionMessage,
		StackTrace:         req.StackTrace,
		UserReport:         req.UserReport,
		MemoryUsedMB:       req.MemoryUsedMB,
		MemoryAvailableMB:  req.MemoryAvailableMB,
		BatteryLevel:       req.BatteryLevel,
		StorageAvailableMB: req.StorageAvailableMB,
		NetworkType:        req.NetworkType,
		CrashStatus:        "new",
	}

	if err := mh.mobileService.ReportCrash(ctx, crash); err != nil {
		mh.logger.Error("Failed to report crash", "error", err)
		http.Error(w, "failed to report crash", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(crash)
}

// ============================================================
// ANALYTICS HANDLERS
// ============================================================

// TrackEvent handles POST /api/v1/mobile/events
func (mh *MobileHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	var req models.TrackAnalyticsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	deviceID := r.URL.Query().Get("device_id")
	appID := r.URL.Query().Get("app_id")
	if deviceID == "" || appID == "" {
		http.Error(w, "device_id and app_id query parameters required", http.StatusBadRequest)
		return
	}

	event := &models.DeviceAnalytic{
		ID:             generateUUID(),
		TenantID:       tenantID,
		AppID:          appID,
		DeviceID:       deviceID,
		UserID:         userID,
		EventType:      req.EventType,
		EventName:      req.EventName,
		EventValue:     req.EventValue,
		ScreenName:     req.ScreenName,
		EventID:        generateUUID(),
		EventTimestamp: time.Now(),
	}

	if err := mh.mobileService.TrackEvent(ctx, event); err != nil {
		mh.logger.Error("Failed to track event", "error", err)
		http.Error(w, "failed to track event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// ============================================================
// APP UPDATE HANDLERS
// ============================================================

// GetLatestUpdate handles GET /api/v1/mobile/updates/latest
func (mh *MobileHandler) GetLatestUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}

	appID := r.URL.Query().Get("app_id")
	if appID == "" {
		http.Error(w, "app_id query parameter required", http.StatusBadRequest)
		return
	}

	update, err := mh.mobileService.GetLatestUpdate(ctx, appID, tenantID)
	if err != nil {
		http.Error(w, "no update available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(update)
}

// ============================================================
// SETTINGS HANDLERS
// ============================================================

// GetUserSettings handles GET /api/v1/mobile/settings
func (mh *MobileHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := getMobileTenantID(w, r)
	if !ok {
		return
	}
	userID, ok := getMobileUserID(w, r)
	if !ok {
		return
	}

	appID := r.URL.Query().Get("app_id")
	if appID == "" {
		http.Error(w, "app_id query parameter required", http.StatusBadRequest)
		return
	}

	settings, err := mh.mobileService.GetUserSettings(ctx, userID, appID, tenantID)
	if err != nil {
		mh.logger.Error("Failed to get settings", "error", err)
		http.Error(w, "failed to get settings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
