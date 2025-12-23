package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// SiteVisitSchedule represents a scheduled site visit
type SiteVisitSchedule struct {
	ID            string         `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	TenantID      string         `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	LeadID        *string        `gorm:"column:lead_id;type:varchar(36)" json:"lead_id"`
	VisitorName   string         `gorm:"column:visitor_name;type:varchar(100);not null" json:"visitor_name"`
	VisitorPhone  *string        `gorm:"column:visitor_phone;type:varchar(20)" json:"visitor_phone"`
	VisitorEmail  *string        `gorm:"column:visitor_email;type:varchar(100)" json:"visitor_email"`
	ScheduledDate time.Time      `gorm:"column:scheduled_date;type:datetime;not null" json:"scheduled_date"`
	ScheduledBy   string         `gorm:"column:scheduled_by;type:varchar(36);not null" json:"scheduled_by"`
	Status        string         `gorm:"column:status;type:varchar(50);not null;default:'scheduled'" json:"status"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
}

// SiteVisitLog represents a logged site visit
type SiteVisitLog struct {
	ID               string       `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	TenantID         string       `gorm:"column:tenant_id;type:varchar(36);not null" json:"tenant_id"`
	VisitScheduleID  string       `gorm:"column:visit_schedule_id;type:varchar(36);not null" json:"visit_schedule_id"`
	CheckInTime      *time.Time   `gorm:"column:check_in_time;type:datetime" json:"check_in_time"`
	CheckOutTime     *time.Time   `gorm:"column:check_out_time;type:datetime" json:"check_out_time"`
	VisitedBy        string       `gorm:"column:visited_by;type:varchar(36);not null" json:"visited_by"`
	UnitsViewed      *UnitsViewed `gorm:"column:units_viewed;type:json" json:"units_viewed"`
	Feedback         *string      `gorm:"column:feedback;type:text" json:"feedback"`
	FollowUpRequired bool         `gorm:"column:follow_up_required;type:boolean;default:false" json:"follow_up_required"`
	NextFollowupDate *time.Time   `gorm:"column:next_followup_date;type:date" json:"next_followup_date"`
	CreatedAt        time.Time    `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time    `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;onUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
}

// UnitsViewed represents the units viewed during a site visit
type UnitsViewed []string

// Value implements the driver.Valuer interface for UnitsViewed
func (u UnitsViewed) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan implements the sql.Scanner interface for UnitsViewed
func (u *UnitsViewed) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("[]"), u)
	}
	return json.Unmarshal(bytes, u)
}

// VisitFeedback represents feedback from a site visit
type VisitFeedback struct {
	Rating      int    `json:"rating"`
	Comments    string `json:"comments"`
	Suggestions string `json:"suggestions"`
}

// DTOs for API requests and responses

// CreateSiteVisitScheduleRequest represents a request to create a site visit schedule
type CreateSiteVisitScheduleRequest struct {
	LeadID        *string   `json:"lead_id"`
	VisitorName   string    `json:"visitor_name" validate:"required"`
	VisitorPhone  *string   `json:"visitor_phone"`
	VisitorEmail  *string   `json:"visitor_email"`
	ScheduledDate time.Time `json:"scheduled_date" validate:"required"`
	ScheduledBy   string    `json:"scheduled_by" validate:"required"`
	Status        string    `json:"status" validate:"oneof=scheduled completed cancelled no_show"`
}

// UpdateSiteVisitScheduleRequest represents a request to update a site visit schedule
type UpdateSiteVisitScheduleRequest struct {
	VisitorName   *string    `json:"visitor_name"`
	VisitorPhone  *string    `json:"visitor_phone"`
	VisitorEmail  *string    `json:"visitor_email"`
	ScheduledDate *time.Time `json:"scheduled_date"`
	Status        *string    `json:"status" validate:"oneof=scheduled completed cancelled no_show"`
}

// SiteVisitScheduleResponse represents a response for site visit schedule
type SiteVisitScheduleResponse struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	LeadID        *string   `json:"lead_id"`
	VisitorName   string    `json:"visitor_name"`
	VisitorPhone  *string   `json:"visitor_phone"`
	VisitorEmail  *string   `json:"visitor_email"`
	ScheduledDate time.Time `json:"scheduled_date"`
	ScheduledBy   string    `json:"scheduled_by"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateSiteVisitLogRequest represents a request to create a site visit log
type CreateSiteVisitLogRequest struct {
	VisitScheduleID  string       `json:"visit_schedule_id" validate:"required"`
	CheckInTime      *time.Time   `json:"check_in_time"`
	CheckOutTime     *time.Time   `json:"check_out_time"`
	VisitedBy        string       `json:"visited_by" validate:"required"`
	UnitsViewed      *UnitsViewed `json:"units_viewed"`
	Feedback         *string      `json:"feedback"`
	FollowUpRequired bool         `json:"follow_up_required"`
	NextFollowupDate *time.Time   `json:"next_followup_date"`
}

// UpdateSiteVisitLogRequest represents a request to update a site visit log
type UpdateSiteVisitLogRequest struct {
	CheckInTime      *time.Time   `json:"check_in_time"`
	CheckOutTime     *time.Time   `json:"check_out_time"`
	UnitsViewed      *UnitsViewed `json:"units_viewed"`
	Feedback         *string      `json:"feedback"`
	FollowUpRequired *bool        `json:"follow_up_required"`
	NextFollowupDate *time.Time   `json:"next_followup_date"`
}

// SiteVisitLogResponse represents a response for site visit log
type SiteVisitLogResponse struct {
	ID               string       `json:"id"`
	TenantID         string       `json:"tenant_id"`
	VisitScheduleID  string       `json:"visit_schedule_id"`
	CheckInTime      *time.Time   `json:"check_in_time"`
	CheckOutTime     *time.Time   `json:"check_out_time"`
	VisitedBy        string       `json:"visited_by"`
	UnitsViewed      *UnitsViewed `json:"units_viewed"`
	Feedback         *string      `json:"feedback"`
	FollowUpRequired bool         `json:"follow_up_required"`
	NextFollowupDate *time.Time   `json:"next_followup_date"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// SiteVisitSummaryResponse provides a summary of site visits
type SiteVisitSummaryResponse struct {
	TotalScheduled   int64     `json:"total_scheduled"`
	TotalCompleted   int64     `json:"total_completed"`
	TotalCancelled   int64     `json:"total_cancelled"`
	TotalNoShow      int64     `json:"total_no_show"`
	UpcomingVisits   int64     `json:"upcoming_visits"`
	FollowUpRequired int64     `json:"follow_up_required"`
	LastUpdated      time.Time `json:"last_updated"`
}

// Validation helper methods

// ValidateStatus validates the status field
func (s *SiteVisitSchedule) ValidateStatus() bool {
	validStatuses := []string{"scheduled", "completed", "cancelled", "no_show"}
	for _, status := range validStatuses {
		if s.Status == status {
			return true
		}
	}
	return false
}

// IsCompleted checks if the visit is completed
func (s *SiteVisitSchedule) IsCompleted() bool {
	return s.Status == "completed"
}

// IsUpcoming checks if the visit is upcoming
func (s *SiteVisitSchedule) IsUpcoming() bool {
	return s.Status == "scheduled" && s.ScheduledDate.After(time.Now())
}

// RequiresFollowUp checks if follow-up is required
func (l *SiteVisitLog) RequiresFollowUp() bool {
	return l.FollowUpRequired
}
