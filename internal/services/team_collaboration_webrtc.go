package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"vyomtech-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ==================== TEAM CHAT SERVICE ====================

type TeamChatService struct {
	db *gorm.DB
}

func NewTeamChatService(db *gorm.DB) *TeamChatService {
	return &TeamChatService{db: db}
}

// CreateChannel creates a new team chat channel
func (s *TeamChatService) CreateChannel(tenantID string, req models.TeamChatChannel) (*models.TeamChatChannel, error) {
	req.ID = uuid.New().String()
	req.TenantID = tenantID
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := s.db.Create(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

// SendMessage sends a message to a channel
func (s *TeamChatService) SendMessage(tenantID string, msg models.TeamChatMessage) (*models.TeamChatMessage, error) {
	msg.ID = uuid.New().String()
	msg.TenantID = tenantID
	msg.CreatedAt = time.Now()

	if err := s.db.Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetChannelMessages retrieves messages from a channel with pagination
func (s *TeamChatService) GetChannelMessages(tenantID, channelID string, limit, offset int) ([]models.TeamChatMessage, error) {
	var messages []models.TeamChatMessage
	err := s.db.Where("tenant_id = ? AND channel_id = ? AND deleted_at IS NULL", tenantID, channelID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// AddChannelMember adds a member to a channel
func (s *TeamChatService) AddChannelMember(tenantID, channelID, userID string, role string) (*models.TeamChatMember, error) {
	member := models.TeamChatMember{
		ID:        uuid.New().String(),
		ChannelID: channelID,
		TenantID:  tenantID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  time.Now(),
	}

	if err := s.db.Create(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetChannelMembers retrieves all members of a channel
func (s *TeamChatService) GetChannelMembers(tenantID, channelID string) ([]models.TeamChatMember, error) {
	var members []models.TeamChatMember
	err := s.db.Where("tenant_id = ? AND channel_id = ? AND left_at IS NULL", tenantID, channelID).
		Find(&members).Error
	return members, err
}

// EditMessage edits an existing message
func (s *TeamChatService) EditMessage(tenantID, messageID, newBody string) error {
	editedAt := time.Now()
	return s.db.Model(&models.TeamChatMessage{}).
		Where("tenant_id = ? AND id = ?", tenantID, messageID).
		Updates(map[string]interface{}{
			"message_body": newBody,
			"is_edited":    1,
			"edited_at":    editedAt,
		}).Error
}

// DeleteMessage soft-deletes a message
func (s *TeamChatService) DeleteMessage(tenantID, messageID string) error {
	deletedAt := time.Now()
	return s.db.Model(&models.TeamChatMessage{}).
		Where("tenant_id = ? AND id = ?", tenantID, messageID).
		Update("deleted_at", deletedAt).Error
}

// ==================== WEBRTC CALL SERVICE ====================

type WebRTCService struct {
	db *gorm.DB
}

func NewWebRTCService(db *gorm.DB) *WebRTCService {
	return &WebRTCService{db: db}
}

// InitiateCall initiates a new WebRTC call
func (s *WebRTCService) InitiateCall(tenantID, initiatorID string, req models.InitiateCallRequest) (*models.VoiceVideoCall, error) {
	call := models.VoiceVideoCall{
		ID:              uuid.New().String(),
		TenantID:        tenantID,
		CallType:        req.CallType,
		InitiatorID:     initiatorID,
		CallStatus:      "RINGING",
		CallDirection:   "OUTBOUND",
		IsAudioEnabled:  1,
		IsVideoEnabled:  1,
		WebRTCRoomID:    uuid.New().String(),
		SignalingServer: "wss://signaling.example.com",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if !req.IsAudioEnabled {
		call.IsAudioEnabled = 0
	}
	if !req.IsVideoEnabled {
		call.IsVideoEnabled = 0
	}

	// Set STUN/TURN servers
	stunServers := []string{
		"stun:stun.l.google.com:19302",
		"stun:stun1.l.google.com:19302",
	}
	call.STUNServers = datatypes.NewJSONType(stunServers)

	if err := s.db.Create(&call).Error; err != nil {
		return nil, err
	}

	// Add initiator as first participant
	s.AddCallParticipant(tenantID, call.ID, initiatorID, "JOINED")

	// Add other participants
	for _, participantID := range req.ParticipantIDs {
		if participantID != initiatorID {
			s.AddCallParticipant(tenantID, call.ID, participantID, "INVITED")
		}
	}

	return &call, nil
}

// AddCallParticipant adds a participant to a call
func (s *WebRTCService) AddCallParticipant(tenantID, callID, userID, status string) (*models.VoiceVideoCallParticipant, error) {
	participant := models.VoiceVideoCallParticipant{
		ID:                uuid.New().String(),
		CallID:            callID,
		TenantID:          tenantID,
		UserID:            userID,
		ParticipantStatus: status,
		CreatedAt:         time.Now(),
	}

	if status == "JOINED" {
		now := time.Now()
		participant.JoinedAt = &now
	}

	if err := s.db.Create(&participant).Error; err != nil {
		return nil, err
	}
	return &participant, nil
}

// UpdateCallStatus updates the call status
func (s *WebRTCService) UpdateCallStatus(tenantID, callID, status string) error {
	now := time.Now()
	update := map[string]interface{}{
		"call_status": status,
		"updated_at":  now,
	}

	switch status {
	case "IN_PROGRESS":
		update["started_at"] = now
	case "ENDED":
		update["ended_at"] = now
	}

	return s.db.Model(&models.VoiceVideoCall{}).
		Where("tenant_id = ? AND id = ?", tenantID, callID).
		Updates(update).Error
}

// EndCall ends a call and calculates duration
func (s *WebRTCService) EndCall(tenantID, callID string) error {
	var call models.VoiceVideoCall
	if err := s.db.First(&call, "tenant_id = ? AND id = ?", tenantID, callID).Error; err != nil {
		return err
	}

	now := time.Now()
	duration := 0
	if call.StartedAt != nil {
		duration = int(now.Sub(*call.StartedAt).Seconds())
	}

	return s.db.Model(&models.VoiceVideoCall{}).
		Where("tenant_id = ? AND id = ?", tenantID, callID).
		Updates(map[string]interface{}{
			"call_status":           "ENDED",
			"ended_at":              now,
			"call_duration_seconds": duration,
		}).Error
}

// GetActiveCall retrieves an active call
func (s *WebRTCService) GetActiveCall(tenantID, callID string) (*models.VoiceVideoCall, error) {
	var call models.VoiceVideoCall
	err := s.db.First(&call, "tenant_id = ? AND id = ? AND call_status IN ('RINGING', 'IN_PROGRESS')", tenantID, callID).Error
	return &call, err
}

// ==================== MEETING ROOM SERVICE ====================

type MeetingRoomService struct {
	db *gorm.DB
}

func NewMeetingRoomService(db *gorm.DB) *MeetingRoomService {
	return &MeetingRoomService{db: db}
}

// CreateMeetingRoom creates a new meeting room
func (s *MeetingRoomService) CreateMeetingRoom(tenantID string, room models.MeetingRoom) (*models.MeetingRoom, error) {
	room.ID = uuid.New().String()
	room.TenantID = tenantID
	room.RoomCode = generateRoomCode()
	room.RoomStatus = "AVAILABLE"
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	if err := s.db.Create(&room).Error; err != nil {
		return nil, err
	}

	// Add creator as owner
	s.GrantRoomAccess(tenantID, room.ID, room.OwnerID, "OWNER")

	return &room, nil
}

// GenerateRoomCode generates a unique room code
func generateRoomCode() string {
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// GrantRoomAccess grants access to a user for a room
func (s *MeetingRoomService) GrantRoomAccess(tenantID, roomID, userID, accessType string) (*models.MeetingRoomAccess, error) {
	access := models.MeetingRoomAccess{
		ID:         uuid.New().String(),
		RoomID:     roomID,
		TenantID:   tenantID,
		UserID:     userID,
		AccessType: accessType,
		IsActive:   1,
		CreatedAt:  time.Now(),
	}

	// Set permissions based on access type
	switch accessType {
	case "OWNER", "MODERATOR":
		access.CanMuteOthers = 1
		access.CanRemoveParticipants = 1
		access.CanRecord = 1
		access.CanShareScreen = 1
	case "PRESENTER":
		access.CanShareScreen = 1
		access.CanRecord = 1
	default:
		access.CanShareScreen = 1
	}

	if err := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&access).Error; err != nil {
		return nil, err
	}
	return &access, nil
}

// GetRoomByCode retrieves a room by room code
func (s *MeetingRoomService) GetRoomByCode(tenantID, roomCode string) (*models.MeetingRoom, error) {
	var room models.MeetingRoom
	err := s.db.First(&room, "tenant_id = ? AND room_code = ?", tenantID, roomCode).Error
	return &room, err
}

// UpdateRoomStatus updates room status and participant count
func (s *MeetingRoomService) UpdateRoomStatus(tenantID, roomID, status string) error {
	now := time.Now()
	update := map[string]interface{}{
		"room_status": status,
		"updated_at":  now,
	}

	if status == "ARCHIVED" {
		update["archived_at"] = now
	}

	return s.db.Model(&models.MeetingRoom{}).
		Where("tenant_id = ? AND id = ?", tenantID, roomID).
		Updates(update).Error
}

// ==================== CALENDAR SERVICE ====================

type CalendarService struct {
	db *gorm.DB
}

func NewCalendarService(db *gorm.DB) *CalendarService {
	return &CalendarService{db: db}
}

// CreateEvent creates a new calendar event
func (s *CalendarService) CreateEvent(tenantID string, req models.CreateEventRequest, creatorID string) (*models.CalendarEvent, error) {
	event := models.CalendarEvent{
		ID:               uuid.New().String(),
		TenantID:         tenantID,
		EventTitle:       req.EventTitle,
		EventDescription: req.EventDescription,
		EventType:        req.EventType,
		CreatorID:        creatorID,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		Location:         req.Location,
		LinkedRoomID:     req.LinkedRoomID,
		IsRecurring:      0,
		Status:           "SCHEDULED",
		IsBusy:           1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if req.IsRecurring {
		event.IsRecurring = 1
		event.RecurrencePattern = req.RecurrencePattern
	}

	if err := s.db.Create(&event).Error; err != nil {
		return nil, err
	}

	// Add attendees
	for _, attendeeID := range req.AttendeeIDs {
		s.AddAttendee(tenantID, event.ID, attendeeID, "INVITED", 0)
	}

	// Add creator as organizer
	s.AddAttendee(tenantID, event.ID, creatorID, "ACCEPTED", 1)

	return &event, nil
}

// AddAttendee adds an attendee to an event
func (s *CalendarService) AddAttendee(tenantID, eventID, userID, status string, isOrganizer int) (*models.CalendarAttendee, error) {
	attendee := models.CalendarAttendee{
		ID:               uuid.New().String(),
		EventID:          eventID,
		TenantID:         tenantID,
		UserID:           userID,
		AttendanceStatus: status,
		IsOrganizer:      isOrganizer,
		CreatedAt:        time.Now(),
	}

	if err := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&attendee).Error; err != nil {
		return nil, err
	}
	return &attendee, nil
}

// GetUserEvents retrieves events for a user within a date range
func (s *CalendarService) GetUserEvents(tenantID, userID string, startTime, endTime time.Time) ([]models.CalendarEvent, error) {
	var events []models.CalendarEvent
	err := s.db.Joins("JOIN calendar_attendee ON calendar_event.id = calendar_attendee.event_id").
		Where("calendar_event.tenant_id = ? AND calendar_attendee.user_id = ? AND calendar_event.start_time >= ? AND calendar_event.start_time <= ?",
			tenantID, userID, startTime, endTime).
		Find(&events).Error
	return events, err
}

// UpdateEventStatus updates the status of an event
func (s *CalendarService) UpdateEventStatus(tenantID, eventID, status string) error {
	now := time.Now()
	update := map[string]interface{}{
		"status":     status,
		"updated_at": now,
	}

	if status == "COMPLETED" {
		update["completed_date"] = now
	}

	return s.db.Model(&models.CalendarEvent{}).
		Where("tenant_id = ? AND id = ?", tenantID, eventID).
		Updates(update).Error
}

// ==================== AUTO-DIALER SERVICE ====================

type AutoDialerService struct {
	db *gorm.DB
}

func NewAutoDialerService(db *gorm.DB) *AutoDialerService {
	return &AutoDialerService{db: db}
}

// CreateCampaign creates a new dialer campaign
func (s *AutoDialerService) CreateCampaign(tenantID string, req models.CreateDialerCampaignRequest, createdBy string) (*models.DialerCampaign, error) {
	campaign := models.DialerCampaign{
		ID:                 uuid.New().String(),
		TenantID:           tenantID,
		CampaignName:       req.CampaignName,
		CampaignType:       req.CampaignType,
		CampaignStatus:     "DRAFT",
		Description:        req.Description,
		ScriptID:           req.ScriptID,
		DialStrategy:       req.DialStrategy,
		MaxConcurrentCalls: req.MaxConcurrentCalls,
		CallerIDNumber:     req.CallerIDNumber,
		ScheduledStartTime: &req.ScheduledStartTime,
		ScheduledEndTime:   &req.ScheduledEndTime,
		VoicemailDetection: 1,
		RecordingEnabled:   1,
		AMDEnabled:         1,
		CreatedBy:          createdBy,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.db.Create(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

// AddCallToQueue adds a call to the priority queue
func (s *AutoDialerService) AddCallToQueue(tenantID, campaignID, phoneNumber, contactName string, priority int) (*models.CallPriorityQueue, error) {
	queueCall := models.CallPriorityQueue{
		ID:                 uuid.New().String(),
		TenantID:           tenantID,
		CampaignID:         campaignID,
		ContactPhoneNumber: phoneNumber,
		ContactName:        contactName,
		PriorityLevel:      priority,
		QueueStatus:        "PENDING",
		CallAttemptCount:   0,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.db.Create(&queueCall).Error; err != nil {
		return nil, err
	}
	return &queueCall, nil
}

// GetNextCallInQueue retrieves the next call to make based on priority
func (s *AutoDialerService) GetNextCallInQueue(tenantID, campaignID string) (*models.CallPriorityQueue, error) {
	var call models.CallPriorityQueue
	err := s.db.Where("tenant_id = ? AND campaign_id = ? AND queue_status = ? AND next_call_time IS NULL OR next_call_time <= ?",
		tenantID, campaignID, "PENDING", time.Now()).
		Order("priority_level DESC, created_at ASC").
		First(&call).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &call, err
}

// AssignCallToAgent assigns a queued call to an agent
func (s *AutoDialerService) AssignCallToAgent(tenantID, queueCallID, agentID string) error {
	now := time.Now()
	return s.db.Model(&models.CallPriorityQueue{}).
		Where("tenant_id = ? AND id = ?", tenantID, queueCallID).
		Updates(map[string]interface{}{
			"assigned_agent_id": agentID,
			"assigned_at":       now,
			"queue_status":      "ASSIGNED",
		}).Error
}

// UpdateCallResult updates the result of a dialed call
func (s *AutoDialerService) UpdateCallResult(tenantID, queueCallID, callResult, notes string, duration int) error {
	now := time.Now()
	update := map[string]interface{}{
		"call_result":           callResult,
		"call_notes":            notes,
		"call_duration_seconds": duration,
		"last_call_time":        now,
		"updated_at":            now,
	}

	// Determine queue status based on result
	switch callResult {
	case "CONNECTED", "VOICEMAIL":
		update["queue_status"] = "COMPLETED"
	default:
		update["queue_status"] = "FAILED"
	}

	return s.db.Model(&models.CallPriorityQueue{}).
		Where("tenant_id = ? AND id = ?", tenantID, queueCallID).
		Updates(update).Error
}

// GetCampaignStats retrieves statistics for a campaign
func (s *AutoDialerService) GetCampaignStats(tenantID, campaignID string) (map[string]interface{}, error) {
	var campaign models.DialerCampaign
	if err := s.db.First(&campaign, "tenant_id = ? AND id = ?", tenantID, campaignID).Error; err != nil {
		return nil, err
	}

	var queueStats struct {
		Pending   int
		Assigned  int
		Completed int
		Failed    int
	}

	s.db.Model(&models.CallPriorityQueue{}).
		Where("tenant_id = ? AND campaign_id = ?", tenantID, campaignID).
		Select("SUM(CASE WHEN queue_status = 'PENDING' THEN 1 ELSE 0 END) as pending",
			"SUM(CASE WHEN queue_status = 'ASSIGNED' THEN 1 ELSE 0 END) as assigned",
			"SUM(CASE WHEN queue_status = 'COMPLETED' THEN 1 ELSE 0 END) as completed",
			"SUM(CASE WHEN queue_status = 'FAILED' THEN 1 ELSE 0 END) as failed").
		Scan(&queueStats)

	stats := map[string]interface{}{
		"campaign_id":     campaign.ID,
		"campaign_name":   campaign.CampaignName,
		"status":          campaign.CampaignStatus,
		"pending_calls":   queueStats.Pending,
		"assigned_calls":  queueStats.Assigned,
		"completed_calls": queueStats.Completed,
		"failed_calls":    queueStats.Failed,
		"total_contacts":  campaign.TotalContacts,
		"contacted_count": campaign.ContactedCount,
		"connected_count": campaign.ConnectedCount,
		"abandoned_count": campaign.AbandonedCount,
	}

	return stats, nil
}

// ==================== WORK TRACKING SERVICE ====================

type WorkTrackingService struct {
	db *gorm.DB
}

func NewWorkTrackingService(db *gorm.DB) *WorkTrackingService {
	return &WorkTrackingService{db: db}
}

// CreateWorkItem creates a new work item
func (s *WorkTrackingService) CreateWorkItem(tenantID string, req models.CreateWorkItemRequest, createdBy string) (*models.WorkItem, error) {
	item := models.WorkItem{
		ID:                 uuid.New().String(),
		TenantID:           tenantID,
		WorkTitle:          req.WorkTitle,
		WorkDescription:    req.WorkDescription,
		WorkType:           req.WorkType,
		Priority:           req.Priority,
		AssignedTo:         req.AssignedTo,
		CreatedBy:          createdBy,
		ParentItemID:       req.ParentItemID,
		DueDate:            &req.DueDate,
		EstimatedHours:     req.EstimatedHours,
		Status:             "TODO",
		PercentageComplete: 0,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateWorkItemStatus updates the status of a work item
func (s *WorkTrackingService) UpdateWorkItemStatus(tenantID, workItemID, status string, percentComplete int) error {
	update := map[string]interface{}{
		"status":              status,
		"percentage_complete": percentComplete,
		"updated_at":          time.Now(),
	}

	if status == "COMPLETED" {
		now := time.Now()
		update["completed_date"] = now
	}

	return s.db.Model(&models.WorkItem{}).
		Where("tenant_id = ? AND id = ?", tenantID, workItemID).
		Updates(update).Error
}

// LogTimeOnWorkItem logs time spent on a work item
func (s *WorkTrackingService) LogTimeOnWorkItem(tenantID, workItemID, userID string, timeSpentMin int, notes string) (*models.WorkItemTimeLog, error) {
	log := models.WorkItemTimeLog{
		ID:           uuid.New().String(),
		WorkItemID:   workItemID,
		TenantID:     tenantID,
		UserID:       userID,
		TimeSpentMin: timeSpentMin,
		LogDate:      time.Now(),
		LogNotes:     notes,
		CreatedAt:    time.Now(),
	}

	if err := s.db.Create(&log).Error; err != nil {
		return nil, err
	}

	// Update actual hours in work item
	var totalMinutes int64
	s.db.Model(&models.WorkItemTimeLog{}).
		Where("tenant_id = ? AND work_item_id = ?", tenantID, workItemID).
		Select("COALESCE(SUM(time_spent_minutes), 0)").
		Row().Scan(&totalMinutes)

	s.db.Model(&models.WorkItem{}).
		Where("tenant_id = ? AND id = ?", tenantID, workItemID).
		Update("actual_hours", float64(totalMinutes)/60.0)

	return &log, nil
}

// AddCommentToWorkItem adds a comment to a work item
func (s *WorkTrackingService) AddCommentToWorkItem(tenantID, workItemID, commenterID, commentText string) (*models.WorkItemComment, error) {
	comment := models.WorkItemComment{
		ID:          uuid.New().String(),
		WorkItemID:  workItemID,
		TenantID:    tenantID,
		CommenterID: commenterID,
		CommentText: commentText,
		CommentType: "COMMENT",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetUserWorkItems retrieves work items assigned to a user
func (s *WorkTrackingService) GetUserWorkItems(tenantID, userID string, status string) ([]models.WorkItem, error) {
	var items []models.WorkItem
	query := s.db.Where("tenant_id = ? AND assigned_to = ?", tenantID, userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("priority DESC, due_date ASC").Find(&items).Error
	return items, err
}
