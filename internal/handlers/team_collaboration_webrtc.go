package handlers

import (
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== TEAM CHAT HANDLERS ====================

type TeamChatHandler struct {
	service *services.TeamChatService
}

func NewTeamChatHandler(db *gorm.DB) *TeamChatHandler {
	return &TeamChatHandler{
		service: services.NewTeamChatService(db),
	}
}

// CreateChannel creates a new team chat channel
func (h *TeamChatHandler) CreateChannel(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var req models.TeamChatChannel

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.service.CreateChannel(tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// SendMessage sends a message to a channel
func (h *TeamChatHandler) SendMessage(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	var req models.TeamChatSendMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := models.TeamChatMessage{
		ChannelID:   req.ChannelID,
		SenderID:    userID,
		MessageType: req.MessageType,
		MessageBody: req.MessageBody,
		FileURL:     req.FileURL,
		FileName:    req.FileName,
	}

	result, err := h.service.SendMessage(tenantID, msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.MessageResponse{
		MessageID: result.ID,
		Status:    "SENT",
		CreatedAt: result.CreatedAt,
	})
}

// GetChannelMessages retrieves messages from a channel
func (h *TeamChatHandler) GetChannelMessages(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	channelID := c.Param("channel_id")
	limit := 50
	offset := 0

	if l := c.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := c.Query("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	messages, err := h.service.GetChannelMessages(tenantID, channelID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}

// ==================== VOICE/VIDEO CALL HANDLERS ====================

type VoiceVideoCallHandler struct {
	service *services.WebRTCService
}

func NewVoiceVideoCallHandler(db *gorm.DB) *VoiceVideoCallHandler {
	return &VoiceVideoCallHandler{
		service: services.NewWebRTCService(db),
	}
}

// InitiateCall initiates a new WebRTC call
func (h *VoiceVideoCallHandler) InitiateCall(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	var req models.InitiateCallRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	call, err := h.service.InitiateCall(tenantID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.CallResponse{
		CallID:          call.ID,
		Status:          call.CallStatus,
		WebRTCRoomID:    call.WebRTCRoomID,
		SignalingServer: call.SignalingServer,
		CreatedAt:       call.CreatedAt,
	})
}

// UpdateCallStatus updates the status of a call
func (h *VoiceVideoCallHandler) UpdateCallStatus(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	callID := c.Param("call_id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateCallStatus(tenantID, callID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// EndCall ends a call
func (h *VoiceVideoCallHandler) EndCall(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	callID := c.Param("call_id")

	err := h.service.EndCall(tenantID, callID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ended"})
}

// ==================== MEETING ROOM HANDLERS ====================

type MeetingRoomHandler struct {
	service *services.MeetingRoomService
}

func NewMeetingRoomHandler(db *gorm.DB) *MeetingRoomHandler {
	return &MeetingRoomHandler{
		service: services.NewMeetingRoomService(db),
	}
}

// CreateMeetingRoom creates a new meeting room
func (h *MeetingRoomHandler) CreateMeetingRoom(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	var req models.MeetingRoom

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.OwnerID = userID
	room, err := h.service.CreateMeetingRoom(tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// GetRoomByCode retrieves a room by room code
func (h *MeetingRoomHandler) GetRoomByCode(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	roomCode := c.Param("room_code")

	room, err := h.service.GetRoomByCode(tenantID, roomCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// GrantRoomAccess grants access to a user for a room
func (h *MeetingRoomHandler) GrantRoomAccess(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var req struct {
		RoomID     string `json:"room_id" binding:"required"`
		UserID     string `json:"user_id" binding:"required"`
		AccessType string `json:"access_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, err := h.service.GrantRoomAccess(tenantID, req.RoomID, req.UserID, req.AccessType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, access)
}

// ==================== CALENDAR HANDLERS ====================

type CalendarHandler struct {
	service *services.CalendarService
}

func NewCalendarHandler(db *gorm.DB) *CalendarHandler {
	return &CalendarHandler{
		service: services.NewCalendarService(db),
	}
}

// CreateEvent creates a new calendar event
func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	var req models.CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(tenantID, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.EventResponse{
		EventID:         event.ID,
		Status:          event.Status,
		StartTime:       event.StartTime,
		EndTime:         event.EndTime,
		CreatedAt:       event.CreatedAt,
		InvitationsSent: 0,
	})
}

// GetUserEvents retrieves events for a user
func (h *CalendarHandler) GetUserEvents(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)

	if start.IsZero() {
		start = time.Now()
	}
	if end.IsZero() {
		end = start.AddDate(0, 1, 0)
	}

	events, err := h.service.GetUserEvents(tenantID, userID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events":     events,
		"start_date": start,
		"end_date":   end,
	})
}

// UpdateEventStatus updates event status
func (h *CalendarHandler) UpdateEventStatus(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	eventID := c.Param("event_id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateEventStatus(tenantID, eventID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// ==================== AUTO-DIALER HANDLERS ====================

type AutoDialerHandler struct {
	service *services.AutoDialerService
}

func NewAutoDialerHandler(db *gorm.DB) *AutoDialerHandler {
	return &AutoDialerHandler{
		service: services.NewAutoDialerService(db),
	}
}

// CreateCampaign creates a new dialer campaign
func (h *AutoDialerHandler) CreateCampaign(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	var req models.CreateDialerCampaignRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := h.service.CreateCampaign(tenantID, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.DialerCampaignResponse{
		CampaignID: campaign.ID,
		Status:     campaign.CampaignStatus,
		CreatedAt:  campaign.CreatedAt,
	})
}

// AddCallToQueue adds a call to the priority queue
func (h *AutoDialerHandler) AddCallToQueue(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	var req struct {
		CampaignID    string `json:"campaign_id" binding:"required"`
		PhoneNumber   string `json:"phone_number" binding:"required"`
		ContactName   string `json:"contact_name"`
		PriorityLevel int    `json:"priority_level" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	call, err := h.service.AddCallToQueue(tenantID, req.CampaignID, req.PhoneNumber, req.ContactName, req.PriorityLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"queue_call_id":  call.ID,
		"campaign_id":    call.CampaignID,
		"phone_number":   call.ContactPhoneNumber,
		"priority_level": call.PriorityLevel,
		"queue_status":   call.QueueStatus,
	})
}

// GetCampaignStats retrieves campaign statistics
func (h *AutoDialerHandler) GetCampaignStats(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	campaignID := c.Param("campaign_id")

	stats, err := h.service.GetCampaignStats(tenantID, campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ==================== WORK TRACKING HANDLERS ====================

type WorkTrackingHandler struct {
	service *services.WorkTrackingService
}

func NewWorkTrackingHandler(db *gorm.DB) *WorkTrackingHandler {
	return &WorkTrackingHandler{
		service: services.NewWorkTrackingService(db),
	}
}

// CreateWorkItem creates a new work item
func (h *WorkTrackingHandler) CreateWorkItem(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	var req models.CreateWorkItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.CreateWorkItem(tenantID, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.WorkItemResponse{
		WorkItemID: item.ID,
		Status:     item.Status,
		CreatedAt:  item.CreatedAt,
	})
}

// UpdateWorkItemStatus updates work item status
func (h *WorkTrackingHandler) UpdateWorkItemStatus(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	workItemID := c.Param("work_item_id")
	var req struct {
		Status             string `json:"status" binding:"required"`
		PercentageComplete int    `json:"percentage_complete"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateWorkItemStatus(tenantID, workItemID, req.Status, req.PercentageComplete)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// LogTime logs time on a work item
func (h *WorkTrackingHandler) LogTime(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")
	workItemID := c.Param("work_item_id")
	var req struct {
		TimeSpentMinutes int    `json:"time_spent_minutes" binding:"required"`
		Notes            string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.service.LogTimeOnWorkItem(tenantID, workItemID, userID, req.TimeSpentMinutes, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"log_id":             log.ID,
		"work_item_id":       log.WorkItemID,
		"time_spent_minutes": log.TimeSpentMin,
		"log_date":           log.LogDate,
	})
}

// ==================== HELPER FUNCTION ====================

// RegisterTeamCollaborationRoutes registers all team collaboration routes
func RegisterTeamCollaborationRoutes(router *gin.Engine, db *gorm.DB) {
	// Team Chat Routes
	chatHandler := NewTeamChatHandler(db)
	chat := router.Group("/api/v1/team-chat")
	{
		chat.POST("/channels", chatHandler.CreateChannel)
		chat.POST("/messages", chatHandler.SendMessage)
		chat.GET("/channels/:channel_id/messages", chatHandler.GetChannelMessages)
	}

	// Voice/Video Call Routes
	callHandler := NewVoiceVideoCallHandler(db)
	calls := router.Group("/api/v1/calls")
	{
		calls.POST("/initiate", callHandler.InitiateCall)
		calls.PUT("/:call_id/status", callHandler.UpdateCallStatus)
		calls.POST("/:call_id/end", callHandler.EndCall)
	}

	// Meeting Room Routes
	roomHandler := NewMeetingRoomHandler(db)
	rooms := router.Group("/api/v1/meeting-rooms")
	{
		rooms.POST("", roomHandler.CreateMeetingRoom)
		rooms.GET("/:room_code", roomHandler.GetRoomByCode)
		rooms.POST("/access", roomHandler.GrantRoomAccess)
	}

	// Calendar Routes
	calendarHandler := NewCalendarHandler(db)
	calendar := router.Group("/api/v1/calendar")
	{
		calendar.POST("/events", calendarHandler.CreateEvent)
		calendar.GET("/events", calendarHandler.GetUserEvents)
		calendar.PUT("/events/:event_id/status", calendarHandler.UpdateEventStatus)
	}

	// Auto-Dialer Routes
	dialerHandler := NewAutoDialerHandler(db)
	dialer := router.Group("/api/v1/dialer")
	{
		dialer.POST("/campaigns", dialerHandler.CreateCampaign)
		dialer.POST("/queue", dialerHandler.AddCallToQueue)
		dialer.GET("/campaigns/:campaign_id/stats", dialerHandler.GetCampaignStats)
	}

	// Work Tracking Routes
	workHandler := NewWorkTrackingHandler(db)
	work := router.Group("/api/v1/work")
	{
		work.POST("/items", workHandler.CreateWorkItem)
		work.PUT("/items/:work_item_id/status", workHandler.UpdateWorkItemStatus)
		work.POST("/items/:work_item_id/log-time", workHandler.LogTime)
	}
}
