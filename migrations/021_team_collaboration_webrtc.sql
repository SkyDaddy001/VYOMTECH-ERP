-- Migration: 021_team_collaboration_webrtc.sql
-- Description: Team collaboration, WebRTC calls, meeting rooms, calendar, auto-dialers
-- Date: 2025-12-03

-- ==================== TEAM CHAT SYSTEM ====================

-- Team chat channels/groups
CREATE TABLE IF NOT EXISTS `team_chat_channel` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    channel_name VARCHAR(100) NOT NULL,
    channel_type ENUM('DIRECT', 'GROUP', 'ANNOUNCEMENT', 'DEPARTMENT', 'PROJECT') NOT NULL,
    description TEXT,
    avatar_url VARCHAR(500),
    is_archived TINYINT(1) DEFAULT 0,
    is_private TINYINT(1) DEFAULT 0,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_type (tenant_id, channel_type),
    INDEX idx_archived (is_archived),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Team chat channels/groups';

-- Chat channel members
CREATE TABLE IF NOT EXISTS `team_chat_member` (
    id CHAR(36) PRIMARY KEY,
    channel_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    role ENUM('OWNER', 'MODERATOR', 'MEMBER') DEFAULT 'MEMBER',
    is_muted TINYINT(1) DEFAULT 0,
    last_read_message_id CHAR(36),
    last_read_at TIMESTAMP NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP NULL,
    INDEX idx_channel_user (channel_id, user_id),
    INDEX idx_tenant_user (tenant_id, user_id),
    FOREIGN KEY (channel_id) REFERENCES team_chat_channel(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Team chat channel members';

-- Chat messages
CREATE TABLE IF NOT EXISTS `team_chat_message` (
    id CHAR(36) PRIMARY KEY,
    channel_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    sender_id CHAR(36) NOT NULL,
    message_type ENUM('TEXT', 'IMAGE', 'FILE', 'VIDEO', 'LINK', 'MENTION', 'SYSTEM') DEFAULT 'TEXT',
    message_body TEXT NOT NULL,
    message_html TEXT,
    file_url VARCHAR(500),
    file_name VARCHAR(255),
    file_size_bytes INT,
    file_mime_type VARCHAR(100),
    mentions JSON,
    reactions JSON,
    is_edited TINYINT(1) DEFAULT 0,
    is_pinned TINYINT(1) DEFAULT 0,
    edited_at TIMESTAMP NULL,
    replied_to_message_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_channel_message (channel_id, created_at DESC),
    INDEX idx_sender (sender_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (channel_id) REFERENCES team_chat_channel(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (replied_to_message_id) REFERENCES team_chat_message(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Team chat messages';

-- Message reactions
CREATE TABLE IF NOT EXISTS `team_chat_reaction` (
    id CHAR(36) PRIMARY KEY,
    message_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    reaction_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_message_user_reaction (message_id, user_id, reaction_type),
    FOREIGN KEY (message_id) REFERENCES team_chat_message(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Message reactions/emojis';

-- ==================== WEBRTC VOICE/VIDEO CALLS ====================

-- Voice and video call sessions
CREATE TABLE IF NOT EXISTS `voice_video_call` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    call_type ENUM('ONE_TO_ONE', 'GROUP', 'CONFERENCE') NOT NULL,
    initiator_id CHAR(36) NOT NULL,
    call_status ENUM('RINGING', 'ACCEPTED', 'IN_PROGRESS', 'MISSED', 'REJECTED', 'ENDED') DEFAULT 'RINGING',
    call_direction ENUM('INBOUND', 'OUTBOUND') DEFAULT 'OUTBOUND',
    is_audio_enabled TINYINT(1) DEFAULT 1,
    is_video_enabled TINYINT(1) DEFAULT 1,
    is_screen_shared TINYINT(1) DEFAULT 0,
    is_recording TINYINT(1) DEFAULT 0,
    recording_url VARCHAR(500),
    call_duration_seconds INT DEFAULT 0,
    started_at TIMESTAMP NULL,
    ended_at TIMESTAMP NULL,
    missed_reason VARCHAR(100),
    webrtc_room_id VARCHAR(100),
    signaling_server VARCHAR(255),
    stun_servers JSON,
    turn_servers JSON,
    ice_candidates JSON,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, call_status),
    INDEX idx_initiator (initiator_id),
    INDEX idx_created_at (created_at DESC),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='WebRTC voice and video calls';

-- Call participants
CREATE TABLE IF NOT EXISTS `voice_video_call_participant` (
    id CHAR(36) PRIMARY KEY,
    call_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    participant_status ENUM('INVITED', 'RINGING', 'JOINED', 'ON_HOLD', 'LEFT') DEFAULT 'INVITED',
    is_audio_muted TINYINT(1) DEFAULT 0,
    is_video_off TINYINT(1) DEFAULT 0,
    joined_at TIMESTAMP NULL,
    left_at TIMESTAMP NULL,
    duration_seconds INT DEFAULT 0,
    audio_quality_score INT,
    video_quality_score INT,
    packet_loss_percent DECIMAL(5,2),
    latency_ms INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_call_user (call_id, user_id),
    INDEX idx_call_status (call_id, participant_status),
    FOREIGN KEY (call_id) REFERENCES voice_video_call(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call participants and their status';

-- ==================== MEETING ROOMS ====================

-- Virtual meeting rooms
CREATE TABLE IF NOT EXISTS `meeting_room` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    room_name VARCHAR(150) NOT NULL,
    room_code VARCHAR(50) UNIQUE,
    description TEXT,
    room_avatar_url VARCHAR(500),
    max_participants INT DEFAULT 100,
    current_participants INT DEFAULT 0,
    room_status ENUM('AVAILABLE', 'IN_USE', 'MAINTENANCE', 'ARCHIVED') DEFAULT 'AVAILABLE',
    is_password_protected TINYINT(1) DEFAULT 0,
    password_hash VARCHAR(255),
    room_type ENUM('PERMANENT', 'TEMPORARY', 'RECURRING') DEFAULT 'PERMANENT',
    owner_id CHAR(36),
    allow_recording TINYINT(1) DEFAULT 1,
    allow_screen_share TINYINT(1) DEFAULT 1,
    allow_chat TINYINT(1) DEFAULT 1,
    webrtc_config JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    archived_at TIMESTAMP NULL,
    INDEX idx_tenant_status (tenant_id, room_status),
    INDEX idx_room_code (room_code),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Virtual meeting rooms';

-- Meeting room access control
CREATE TABLE IF NOT EXISTS `meeting_room_access` (
    id CHAR(36) PRIMARY KEY,
    room_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36),
    user_role VARCHAR(50),
    access_type ENUM('OWNER', 'MODERATOR', 'PRESENTER', 'PARTICIPANT', 'VIEWER') DEFAULT 'PARTICIPANT',
    can_mute_others TINYINT(1) DEFAULT 0,
    can_remove_participants TINYINT(1) DEFAULT 0,
    can_record TINYINT(1) DEFAULT 0,
    can_share_screen TINYINT(1) DEFAULT 1,
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_room_user (room_id, user_id),
    INDEX idx_room_access (room_id, access_type),
    FOREIGN KEY (room_id) REFERENCES meeting_room(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Meeting room access and permissions';

-- ==================== CALENDAR & APPOINTMENTS ====================

-- Calendar events/appointments
CREATE TABLE IF NOT EXISTS `calendar_event` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    event_title VARCHAR(200) NOT NULL,
    event_description TEXT,
    event_type ENUM('MEETING', 'CALL', 'TASK', 'REMINDER', 'APPOINTMENT', 'CONFERENCE') NOT NULL,
    creator_id CHAR(36) NOT NULL,
    assigned_to CHAR(36),
    linked_room_id CHAR(36),
    linked_call_id CHAR(36),
    location VARCHAR(255),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration_minutes INT,
    timezone VARCHAR(50),
    is_all_day TINYINT(1) DEFAULT 0,
    reminder_minutes INT DEFAULT 15,
    status ENUM('SCHEDULED', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED', 'RESCHEDULED') DEFAULT 'SCHEDULED',
    is_recurring TINYINT(1) DEFAULT 0,
    recurrence_pattern VARCHAR(100),
    recurrence_end_date TIMESTAMP NULL,
    is_busy TINYINT(1) DEFAULT 1,
    is_private TINYINT(1) DEFAULT 0,
    calendar_id CHAR(36),
    color_code VARCHAR(20),
    attachments JSON,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_time (tenant_id, start_time, end_time),
    INDEX idx_assigned_to (assigned_to),
    INDEX idx_creator (creator_id),
    INDEX idx_status (status),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (linked_room_id) REFERENCES meeting_room(id) ON DELETE SET NULL,
    FOREIGN KEY (linked_call_id) REFERENCES voice_video_call(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Calendar events and appointments';

-- Event attendees
CREATE TABLE IF NOT EXISTS `calendar_attendee` (
    id CHAR(36) PRIMARY KEY,
    event_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    attendance_status ENUM('INVITED', 'ACCEPTED', 'DECLINED', 'TENTATIVE', 'NO_RESPONSE') DEFAULT 'INVITED',
    reminder_sent TINYINT(1) DEFAULT 0,
    is_organizer TINYINT(1) DEFAULT 0,
    responded_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_event_user (event_id, user_id),
    INDEX idx_event_status (event_id, attendance_status),
    INDEX idx_user_status (user_id, attendance_status),
    FOREIGN KEY (event_id) REFERENCES calendar_event(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Calendar event attendees';

-- ==================== AUTO-DIALER WITH PRIORITY QUEUE ====================

-- Dialer campaigns
CREATE TABLE IF NOT EXISTS `dialer_campaign` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    campaign_name VARCHAR(150) NOT NULL,
    campaign_type ENUM('OUTBOUND', 'PREVIEW', 'PREDICTIVE', 'PROGRESSIVE') NOT NULL,
    campaign_status ENUM('DRAFT', 'SCHEDULED', 'ACTIVE', 'PAUSED', 'COMPLETED', 'CANCELLED') DEFAULT 'DRAFT',
    description TEXT,
    script_id CHAR(36),
    call_list_id CHAR(36),
    dial_strategy ENUM('SEQUENTIAL', 'RANDOM', 'PRIORITY_BASED', 'SKILL_BASED') DEFAULT 'SEQUENTIAL',
    max_concurrent_calls INT DEFAULT 10,
    max_retries INT DEFAULT 3,
    retry_interval_minutes INT DEFAULT 60,
    abandoned_call_threshold_percent DECIMAL(5,2) DEFAULT 3.0,
    do_not_call_list_id CHAR(36),
    caller_id_number VARCHAR(20),
    voicemail_detection TINYINT(1) DEFAULT 1,
    recording_enabled TINYINT(1) DEFAULT 1,
    amd_enabled TINYINT(1) DEFAULT 1,
    scheduled_start_time TIMESTAMP NULL,
    scheduled_end_time TIMESTAMP NULL,
    actual_start_time TIMESTAMP NULL,
    actual_end_time TIMESTAMP NULL,
    total_contacts INT DEFAULT 0,
    contacted_count INT DEFAULT 0,
    connected_count INT DEFAULT 0,
    failed_count INT DEFAULT 0,
    abandoned_count INT DEFAULT 0,
    created_by CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, campaign_status),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Dialer campaigns for auto-calling';

-- Call priority queue
CREATE TABLE IF NOT EXISTS `call_priority_queue` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    campaign_id CHAR(36) NOT NULL,
    contact_phone_number VARCHAR(20) NOT NULL,
    contact_name VARCHAR(100),
    contact_id CHAR(36),
    lead_id CHAR(36),
    priority_level INT DEFAULT 0,
    priority_reason VARCHAR(100),
    queue_status ENUM('PENDING', 'ASSIGNED', 'CALLING', 'COMPLETED', 'FAILED', 'RESCHEDULED') DEFAULT 'PENDING',
    assigned_agent_id CHAR(36),
    assigned_at TIMESTAMP NULL,
    call_attempt_count INT DEFAULT 0,
    last_call_time TIMESTAMP NULL,
    next_call_time TIMESTAMP NULL,
    call_result ENUM('CONNECTED', 'VOICEMAIL', 'DISCONNECTED', 'NO_ANSWER', 'BUSY', 'INVALID', 'DO_NOT_CALL') DEFAULT 'NO_ANSWER',
    call_notes TEXT,
    call_duration_seconds INT DEFAULT 0,
    is_callback TINYINT(1) DEFAULT 0,
    callback_requested_time TIMESTAMP NULL,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_campaign_priority (campaign_id, priority_level DESC),
    INDEX idx_campaign_status (campaign_id, queue_status),
    INDEX idx_agent_status (assigned_agent_id, queue_status),
    INDEX idx_next_call (next_call_time),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES dialer_campaign(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Call priority queue for auto-dialer';

-- Dialer scripts/call flows
CREATE TABLE IF NOT EXISTS `dialer_script` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    script_name VARCHAR(150) NOT NULL,
    script_description TEXT,
    script_type ENUM('GREETING', 'QUALIFICATION', 'OBJECTION_HANDLING', 'CLOSING', 'FOLLOW_UP') NOT NULL,
    script_content JSON,
    voice_guidance_audio_url VARCHAR(500),
    created_by CHAR(36),
    is_active TINYINT(1) DEFAULT 1,
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_type (tenant_id, script_type),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Dialer scripts and call flows';

-- ==================== WORK TRACKING & COMPLETION ====================

-- Work items/tasks
CREATE TABLE IF NOT EXISTS `work_item` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    work_title VARCHAR(255) NOT NULL,
    work_description TEXT,
    work_type ENUM('TASK', 'BUG', 'FEATURE', 'IMPROVEMENT', 'DOCUMENTATION') NOT NULL,
    status ENUM('TODO', 'IN_PROGRESS', 'IN_REVIEW', 'BLOCKED', 'COMPLETED', 'CANCELLED') DEFAULT 'TODO',
    priority ENUM('CRITICAL', 'HIGH', 'MEDIUM', 'LOW') DEFAULT 'MEDIUM',
    assigned_to CHAR(36),
    created_by CHAR(36) NOT NULL,
    parent_item_id CHAR(36),
    estimated_hours DECIMAL(8,2),
    actual_hours DECIMAL(8,2),
    due_date TIMESTAMP NULL,
    completed_date TIMESTAMP NULL,
    percentage_complete INT DEFAULT 0,
    tags JSON,
    attachments JSON,
    linked_chat_channel_id CHAR(36),
    linked_call_id CHAR(36),
    linked_event_id CHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_assigned_to (assigned_to),
    INDEX idx_due_date (due_date),
    INDEX idx_priority (priority),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_item_id) REFERENCES work_item(id) ON DELETE SET NULL,
    FOREIGN KEY (linked_chat_channel_id) REFERENCES team_chat_channel(id) ON DELETE SET NULL,
    FOREIGN KEY (linked_call_id) REFERENCES voice_video_call(id) ON DELETE SET NULL,
    FOREIGN KEY (linked_event_id) REFERENCES calendar_event(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Work items and tasks';

-- Work item comments/updates
CREATE TABLE IF NOT EXISTS `work_item_comment` (
    id CHAR(36) PRIMARY KEY,
    work_item_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    commenter_id CHAR(36) NOT NULL,
    comment_text TEXT NOT NULL,
    comment_type ENUM('COMMENT', 'STATUS_UPDATE', 'ATTACHMENT', 'MENTION') NOT NULL,
    mentions JSON,
    attachments JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_work_item (work_item_id),
    INDEX idx_commenter (commenter_id),
    FOREIGN KEY (work_item_id) REFERENCES work_item(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Work item comments and updates';

-- Work item time tracking
CREATE TABLE IF NOT EXISTS `work_item_time_log` (
    id CHAR(36) PRIMARY KEY,
    work_item_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    time_spent_minutes INT NOT NULL,
    log_date DATE NOT NULL,
    log_notes TEXT,
    is_billable TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_work_item_date (work_item_id, log_date),
    INDEX idx_user_date (user_id, log_date),
    FOREIGN KEY (work_item_id) REFERENCES work_item(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Work item time tracking and logging';

-- ==================== REAL-TIME NOTIFICATIONS ====================

-- Notifications for WebSocket delivery
CREATE TABLE IF NOT EXISTS `user_notification` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    notification_type ENUM('MESSAGE', 'CALL', 'TASK', 'EVENT', 'MENTION', 'SYSTEM') NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    reference_id CHAR(36),
    reference_type VARCHAR(50),
    is_read TINYINT(1) DEFAULT 0,
    read_at TIMESTAMP NULL,
    action_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    INDEX idx_user_read (user_id, is_read),
    INDEX idx_created_at (created_at DESC),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Real-time user notifications';

-- ==================== USER PRESENCE & STATUS ====================

-- User online status
CREATE TABLE IF NOT EXISTS `user_presence` (
    id CHAR(36) PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    status ENUM('ONLINE', 'AWAY', 'BUSY', 'DND', 'OFFLINE') DEFAULT 'OFFLINE',
    current_activity VARCHAR(255),
    last_seen TIMESTAMP,
    session_id VARCHAR(100),
    device_type VARCHAR(50),
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_tenant_user (tenant_id, user_id),
    INDEX idx_status (status),
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='User online presence and status';

-- ==================== INDEXES FOR PERFORMANCE ====================

CREATE INDEX idx_team_chat_channel_search ON team_chat_channel(tenant_id, channel_name);
CREATE INDEX idx_voice_call_timing ON voice_video_call(started_at, ended_at);
CREATE INDEX idx_calendar_date_range ON calendar_event(start_time, end_time);
CREATE INDEX idx_dialer_queue_priority ON call_priority_queue(campaign_id, priority_level, queue_status);
CREATE INDEX idx_work_item_search ON work_item(tenant_id, work_title);
