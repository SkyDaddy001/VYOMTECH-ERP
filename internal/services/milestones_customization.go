package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// MilestoneCustomizationService defines customization for sales milestones
type MilestoneCustomizationService interface {
	// Milestone Type Management
	CreateMilestoneType(ctx context.Context, tenantID string, milestoneType *MilestoneTypeConfig) (*MilestoneTypeConfig, error)
	GetMilestoneType(ctx context.Context, tenantID, typeCode string) (*MilestoneTypeConfig, error)
	GetMilestoneTypes(ctx context.Context, tenantID string) ([]MilestoneTypeConfig, error)
	UpdateMilestoneType(ctx context.Context, tenantID string, milestoneType *MilestoneTypeConfig) (*MilestoneTypeConfig, error)
	DeactivateMilestoneType(ctx context.Context, tenantID, typeCode string) error
	DeleteMilestoneType(ctx context.Context, tenantID, typeCode string) error

	// Milestone Template Management
	CreateMilestoneTemplate(ctx context.Context, tenantID string, template *MilestoneTemplate) (*MilestoneTemplate, error)
	GetMilestoneTemplate(ctx context.Context, tenantID string, templateID int64) (*MilestoneTemplate, error)
	GetMilestoneTemplates(ctx context.Context, tenantID string, templateType string) ([]MilestoneTemplate, error)
	UpdateMilestoneTemplate(ctx context.Context, tenantID string, template *MilestoneTemplate) (*MilestoneTemplate, error)
	DeleteMilestoneTemplate(ctx context.Context, tenantID string, templateID int64) error

	// Lead Milestone Management
	CreateLeadMilestone(ctx context.Context, tenantID string, milestone *LeadMilestone) (*LeadMilestone, error)
	GetLeadMilestones(ctx context.Context, tenantID string, leadID int64) ([]LeadMilestone, error)
	UpdateLeadMilestoneStatus(ctx context.Context, tenantID string, milestoneID int64, status string) error
	GetMilestoneTimeline(ctx context.Context, tenantID string, leadID int64) ([]LeadMilestone, error)

	// Milestone Analytics
	GetMilestoneCompletionMetrics(ctx context.Context, tenantID string) (map[string]interface{}, error)
	GetMilestoneTimeTrends(ctx context.Context, tenantID string, milestoneType string, days int) ([]MilestoneTimeTrendData, error)
	GetMilestoneBottlenecks(ctx context.Context, tenantID string) ([]BottleneckAnalysis, error)
}

// ============================================================================
// DATA MODELS
// ============================================================================

// MilestoneTypeConfig defines a type of milestone
type MilestoneTypeConfig struct {
	ID               int64           `json:"id" db:"id"`
	TenantID         string          `json:"tenant_id" db:"tenant_id"`
	TypeCode         string          `json:"type_code" db:"type_code"`
	TypeName         string          `json:"type_name" db:"type_name"`
	Description      *string         `json:"description,omitempty" db:"description"`
	Icon             *string         `json:"icon,omitempty" db:"icon"`
	ColorHex         *string         `json:"color_hex,omitempty" db:"color_hex"`
	DisplayOrder     int             `json:"display_order" db:"display_order"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	IsRequired       bool            `json:"is_required" db:"is_required"`
	IsMandatory      bool            `json:"is_mandatory" db:"is_mandatory"`
	Category         string          `json:"category" db:"category"` // engagement, action, decision, completion
	TypicalDuration  *int            `json:"typical_duration_days,omitempty" db:"typical_duration_days"`
	SLADays          *int            `json:"sla_days,omitempty" db:"sla_days"`
	AllowsAttachment bool            `json:"allows_attachment" db:"allows_attachment"`
	AllowsLocation   bool            `json:"allows_location" db:"allows_location"`
	AllowsNotes      bool            `json:"allows_notes" db:"allows_notes"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedBy        *int64          `json:"created_by,omitempty" db:"created_by"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// MilestoneTemplate defines a milestone workflow template
type MilestoneTemplate struct {
	ID            int64         `json:"id" db:"id"`
	TenantID      string        `json:"tenant_id" db:"tenant_id"`
	TemplateName  string        `json:"template_name" db:"template_name"`
	TemplateType  string        `json:"template_type" db:"template_type"` // standard, campaign, project, custom
	Description   *string       `json:"description,omitempty" db:"description"`
	IsActive      bool          `json:"is_active" db:"is_active"`
	IsDefault     bool          `json:"is_default" db:"is_default"`
	Sequence      json.RawMessage `json:"sequence" db:"sequence"` // JSON array of milestone types in order
	EstimatedDays int           `json:"estimated_days" db:"estimated_days"`
	CreatedBy     *int64        `json:"created_by,omitempty" db:"created_by"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
}

// LeadMilestone represents a milestone achieved by a lead
type LeadMilestone struct {
	ID               int64         `json:"id" db:"id"`
	TenantID         string        `json:"tenant_id" db:"tenant_id"`
	LeadID           int64         `json:"lead_id" db:"lead_id"`
	MilestoneTypeID  int64         `json:"milestone_type_id" db:"milestone_type_id"`
	TypeCode         string        `json:"type_code" db:"type_code"`
	TypeName         string        `json:"type_name" db:"type_name"`
	AchievedDate     time.Time     `json:"achieved_date" db:"achieved_date"`
	AchievedTime     *time.Time    `json:"achieved_time,omitempty" db:"achieved_time"`
	Status           string        `json:"status" db:"status"` // pending, in_progress, completed, skipped, failed
	DaysFromPrevious *int          `json:"days_from_previous,omitempty" db:"days_from_previous"`
	Notes            *string       `json:"notes,omitempty" db:"notes"`
	LocationLatitude *float64      `json:"location_latitude,omitempty" db:"location_latitude"`
	LocationLongitude *float64     `json:"location_longitude,omitempty" db:"location_longitude"`
	LocationName     *string       `json:"location_name,omitempty" db:"location_name"`
	DurationMinutes  *int          `json:"duration_minutes,omitempty" db:"duration_minutes"`
	Outcome          *string       `json:"outcome,omitempty" db:"outcome"` // positive, neutral, negative
	FollowUpDate     *time.Time    `json:"follow_up_date,omitempty" db:"follow_up_date"`
	FollowUpRequired bool          `json:"follow_up_required" db:"follow_up_required"`
	DocumentURLs     json.RawMessage `json:"document_urls,omitempty" db:"document_urls"` // JSON array
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CompletedBy      *int64        `json:"completed_by,omitempty" db:"completed_by"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
}

// MilestoneTimeTrendData represents trend data for milestones
type MilestoneTimeTrendData struct {
	Date              time.Time `json:"date"`
	MilestoneType     string    `json:"milestone_type"`
	MilestoneCount    int64     `json:"milestone_count"`
	AverageTimeToHit  float64   `json:"average_time_to_hit"`
	CompletionRate    float64   `json:"completion_rate"`
	FollowUpRequired  int64     `json:"follow_up_required"`
}

// BottleneckAnalysis identifies slow milestone transitions
type BottleneckAnalysis struct {
	FromMilestone    string    `json:"from_milestone"`
	ToMilestone      string    `json:"to_milestone"`
	AverageDays      float64   `json:"average_days"`
	MedianDays       float64   `json:"median_days"`
	LeadsAffected    int64     `json:"leads_affected"`
	CompletionRate   float64   `json:"completion_rate"`
	SLABreachPercent float64   `json:"sla_breach_percent"`
}

// ============================================================================
// SERVICE IMPLEMENTATION
// ============================================================================

type milestoneCustomizationService struct {
	db *sql.DB
}

// NewMilestoneCustomizationService creates a new milestone customization service
func NewMilestoneCustomizationService(db *sql.DB) MilestoneCustomizationService {
	return &milestoneCustomizationService{db: db}
}

// ============================================================================
// MILESTONE TYPE CRUD OPERATIONS
// ============================================================================

// CreateMilestoneType creates a new milestone type
func (s *milestoneCustomizationService) CreateMilestoneType(ctx context.Context, tenantID string, milestoneType *MilestoneTypeConfig) (*MilestoneTypeConfig, error) {
	if milestoneType.TypeCode == "" || milestoneType.TypeName == "" {
		return nil, errors.New("type_code and type_name are required")
	}

	query := `
		INSERT INTO tenant_milestone_types (
			tenant_id, type_code, type_name, description, icon, color_hex,
			display_order, is_active, is_required, is_mandatory, category,
			typical_duration_days, sla_days, allows_attachment, allows_location,
			allows_notes, metadata, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	metadataJSON := "{}"
	if milestoneType.Metadata != nil {
		metadataJSON = string(milestoneType.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, milestoneType.TypeCode, milestoneType.TypeName, milestoneType.Description,
		milestoneType.Icon, milestoneType.ColorHex, milestoneType.DisplayOrder,
		milestoneType.IsActive, milestoneType.IsRequired, milestoneType.IsMandatory,
		milestoneType.Category, milestoneType.TypicalDuration, milestoneType.SLADays,
		milestoneType.AllowsAttachment, milestoneType.AllowsLocation,
		milestoneType.AllowsNotes, metadataJSON, milestoneType.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create milestone type: %w", err)
	}

	id, _ := result.LastInsertId()
	milestoneType.ID = id
	milestoneType.TenantID = tenantID
	milestoneType.CreatedAt = time.Now()
	milestoneType.UpdatedAt = time.Now()

	return milestoneType, nil
}

// GetMilestoneType retrieves a specific milestone type
func (s *milestoneCustomizationService) GetMilestoneType(ctx context.Context, tenantID, typeCode string) (*MilestoneTypeConfig, error) {
	query := `
		SELECT id, tenant_id, type_code, type_name, description, icon, color_hex,
		       display_order, is_active, is_required, is_mandatory, category,
		       typical_duration_days, sla_days, allows_attachment, allows_location,
		       allows_notes, metadata, created_by, created_at, updated_at
		FROM tenant_milestone_types
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	var milestoneType MilestoneTypeConfig
	var metadataStr sql.NullString

	err := s.db.QueryRowContext(ctx, query, tenantID, typeCode).Scan(
		&milestoneType.ID, &milestoneType.TenantID, &milestoneType.TypeCode,
		&milestoneType.TypeName, &milestoneType.Description, &milestoneType.Icon,
		&milestoneType.ColorHex, &milestoneType.DisplayOrder, &milestoneType.IsActive,
		&milestoneType.IsRequired, &milestoneType.IsMandatory, &milestoneType.Category,
		&milestoneType.TypicalDuration, &milestoneType.SLADays,
		&milestoneType.AllowsAttachment, &milestoneType.AllowsLocation,
		&milestoneType.AllowsNotes, &metadataStr, &milestoneType.CreatedBy,
		&milestoneType.CreatedAt, &milestoneType.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("milestone type not found: %s", typeCode)
		}
		return nil, fmt.Errorf("failed to get milestone type: %w", err)
	}

	if metadataStr.Valid {
		milestoneType.Metadata = json.RawMessage(metadataStr.String)
	}

	return &milestoneType, nil
}

// GetMilestoneTypes retrieves all milestone types
func (s *milestoneCustomizationService) GetMilestoneTypes(ctx context.Context, tenantID string) ([]MilestoneTypeConfig, error) {
	query := `
		SELECT id, tenant_id, type_code, type_name, description, icon, color_hex,
		       display_order, is_active, is_required, is_mandatory, category,
		       typical_duration_days, sla_days, allows_attachment, allows_location,
		       allows_notes, metadata, created_by, created_at, updated_at
		FROM tenant_milestone_types
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY display_order ASC, type_name ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone types: %w", err)
	}
	defer rows.Close()

	var types []MilestoneTypeConfig
	for rows.Next() {
		var mt MilestoneTypeConfig
		var metadataStr sql.NullString

		err := rows.Scan(
			&mt.ID, &mt.TenantID, &mt.TypeCode, &mt.TypeName, &mt.Description,
			&mt.Icon, &mt.ColorHex, &mt.DisplayOrder, &mt.IsActive,
			&mt.IsRequired, &mt.IsMandatory, &mt.Category, &mt.TypicalDuration,
			&mt.SLADays, &mt.AllowsAttachment, &mt.AllowsLocation,
			&mt.AllowsNotes, &metadataStr, &mt.CreatedBy, &mt.CreatedAt, &mt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone type: %w", err)
		}

		if metadataStr.Valid {
			mt.Metadata = json.RawMessage(metadataStr.String)
		}

		types = append(types, mt)
	}

	return types, rows.Err()
}

// UpdateMilestoneType updates a milestone type
func (s *milestoneCustomizationService) UpdateMilestoneType(ctx context.Context, tenantID string, milestoneType *MilestoneTypeConfig) (*MilestoneTypeConfig, error) {
	query := `
		UPDATE tenant_milestone_types
		SET type_name = ?, description = ?, icon = ?, color_hex = ?,
		    display_order = ?, is_active = ?, is_required = ?, is_mandatory = ?,
		    category = ?, typical_duration_days = ?, sla_days = ?,
		    allows_attachment = ?, allows_location = ?, allows_notes = ?,
		    metadata = ?, updated_at = NOW()
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	metadataJSON := "{}"
	if milestoneType.Metadata != nil {
		metadataJSON = string(milestoneType.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		milestoneType.TypeName, milestoneType.Description, milestoneType.Icon,
		milestoneType.ColorHex, milestoneType.DisplayOrder, milestoneType.IsActive,
		milestoneType.IsRequired, milestoneType.IsMandatory, milestoneType.Category,
		milestoneType.TypicalDuration, milestoneType.SLADays,
		milestoneType.AllowsAttachment, milestoneType.AllowsLocation,
		milestoneType.AllowsNotes, metadataJSON, tenantID, milestoneType.TypeCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update milestone type: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("milestone type not found")
	}

	milestoneType.UpdatedAt = time.Now()
	return milestoneType, nil
}

// DeactivateMilestoneType deactivates a milestone type
func (s *milestoneCustomizationService) DeactivateMilestoneType(ctx context.Context, tenantID, typeCode string) error {
	query := `
		UPDATE tenant_milestone_types
		SET is_active = FALSE, updated_at = NOW()
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, typeCode)
	if err != nil {
		return fmt.Errorf("failed to deactivate milestone type: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("milestone type not found")
	}

	return nil
}

// DeleteMilestoneType soft deletes a milestone type
func (s *milestoneCustomizationService) DeleteMilestoneType(ctx context.Context, tenantID, typeCode string) error {
	query := `
		UPDATE tenant_milestone_types
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, typeCode)
	if err != nil {
		return fmt.Errorf("failed to delete milestone type: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("milestone type not found")
	}

	return nil
}

// ============================================================================
// MILESTONE TEMPLATE OPERATIONS
// ============================================================================

// CreateMilestoneTemplate creates a new milestone template
func (s *milestoneCustomizationService) CreateMilestoneTemplate(ctx context.Context, tenantID string, template *MilestoneTemplate) (*MilestoneTemplate, error) {
	if template.TemplateName == "" || template.TemplateType == "" {
		return nil, errors.New("template_name and template_type are required")
	}

	query := `
		INSERT INTO tenant_milestone_templates (
			tenant_id, template_name, template_type, description,
			is_active, is_default, sequence, estimated_days, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	sequenceJSON := "[]"
	if template.Sequence != nil {
		sequenceJSON = string(template.Sequence)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, template.TemplateName, template.TemplateType, template.Description,
		template.IsActive, template.IsDefault, sequenceJSON, template.EstimatedDays,
		template.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create milestone template: %w", err)
	}

	id, _ := result.LastInsertId()
	template.ID = id
	template.TenantID = tenantID
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	return template, nil
}

// GetMilestoneTemplate retrieves a specific milestone template
func (s *milestoneCustomizationService) GetMilestoneTemplate(ctx context.Context, tenantID string, templateID int64) (*MilestoneTemplate, error) {
	query := `
		SELECT id, tenant_id, template_name, template_type, description,
		       is_active, is_default, sequence, estimated_days, created_by, created_at, updated_at
		FROM tenant_milestone_templates
		WHERE tenant_id = ? AND id = ?
	`

	var template MilestoneTemplate
	var sequenceStr sql.NullString

	err := s.db.QueryRowContext(ctx, query, tenantID, templateID).Scan(
		&template.ID, &template.TenantID, &template.TemplateName, &template.TemplateType,
		&template.Description, &template.IsActive, &template.IsDefault,
		&sequenceStr, &template.EstimatedDays, &template.CreatedBy,
		&template.CreatedAt, &template.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("milestone template not found")
		}
		return nil, fmt.Errorf("failed to get milestone template: %w", err)
	}

	if sequenceStr.Valid {
		template.Sequence = json.RawMessage(sequenceStr.String)
	}

	return &template, nil
}

// GetMilestoneTemplates retrieves milestone templates by type
func (s *milestoneCustomizationService) GetMilestoneTemplates(ctx context.Context, tenantID string, templateType string) ([]MilestoneTemplate, error) {
	query := `
		SELECT id, tenant_id, template_name, template_type, description,
		       is_active, is_default, sequence, estimated_days, created_by, created_at, updated_at
		FROM tenant_milestone_templates
		WHERE tenant_id = ? AND template_type = ? AND is_active = TRUE
		ORDER BY is_default DESC, template_name ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, templateType)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone templates: %w", err)
	}
	defer rows.Close()

	var templates []MilestoneTemplate
	for rows.Next() {
		var template MilestoneTemplate
		var sequenceStr sql.NullString

		err := rows.Scan(
			&template.ID, &template.TenantID, &template.TemplateName, &template.TemplateType,
			&template.Description, &template.IsActive, &template.IsDefault,
			&sequenceStr, &template.EstimatedDays, &template.CreatedBy,
			&template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone template: %w", err)
		}

		if sequenceStr.Valid {
			template.Sequence = json.RawMessage(sequenceStr.String)
		}

		templates = append(templates, template)
	}

	return templates, rows.Err()
}

// UpdateMilestoneTemplate updates a milestone template
func (s *milestoneCustomizationService) UpdateMilestoneTemplate(ctx context.Context, tenantID string, template *MilestoneTemplate) (*MilestoneTemplate, error) {
	query := `
		UPDATE tenant_milestone_templates
		SET template_name = ?, template_type = ?, description = ?,
		    is_active = ?, is_default = ?, sequence = ?, estimated_days = ?,
		    updated_at = NOW()
		WHERE tenant_id = ? AND id = ?
	`

	sequenceJSON := "[]"
	if template.Sequence != nil {
		sequenceJSON = string(template.Sequence)
	}

	result, err := s.db.ExecContext(ctx, query,
		template.TemplateName, template.TemplateType, template.Description,
		template.IsActive, template.IsDefault, sequenceJSON, template.EstimatedDays,
		tenantID, template.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update milestone template: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("milestone template not found")
	}

	template.UpdatedAt = time.Now()
	return template, nil
}

// DeleteMilestoneTemplate deletes a milestone template
func (s *milestoneCustomizationService) DeleteMilestoneTemplate(ctx context.Context, tenantID string, templateID int64) error {
	query := `DELETE FROM tenant_milestone_templates WHERE tenant_id = ? AND id = ?`

	result, err := s.db.ExecContext(ctx, query, tenantID, templateID)
	if err != nil {
		return fmt.Errorf("failed to delete milestone template: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("milestone template not found")
	}

	return nil
}

// ============================================================================
// LEAD MILESTONE OPERATIONS
// ============================================================================

// CreateLeadMilestone records a milestone for a lead
func (s *milestoneCustomizationService) CreateLeadMilestone(ctx context.Context, tenantID string, milestone *LeadMilestone) (*LeadMilestone, error) {
	if milestone.LeadID == 0 || milestone.MilestoneTypeID == 0 {
		return nil, errors.New("lead_id and milestone_type_id are required")
	}

	query := `
		INSERT INTO lead_milestones (
			tenant_id, lead_id, milestone_type_id, type_code, type_name,
			achieved_date, achieved_time, status, days_from_previous, notes,
			location_latitude, location_longitude, location_name, duration_minutes,
			outcome, follow_up_date, follow_up_required, document_urls, metadata,
			completed_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	docUrlsJSON := "[]"
	if milestone.DocumentURLs != nil {
		docUrlsJSON = string(milestone.DocumentURLs)
	}

	metadataJSON := "{}"
	if milestone.Metadata != nil {
		metadataJSON = string(milestone.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, milestone.LeadID, milestone.MilestoneTypeID, milestone.TypeCode,
		milestone.TypeName, milestone.AchievedDate, milestone.AchievedTime, milestone.Status,
		milestone.DaysFromPrevious, milestone.Notes, milestone.LocationLatitude,
		milestone.LocationLongitude, milestone.LocationName, milestone.DurationMinutes,
		milestone.Outcome, milestone.FollowUpDate, milestone.FollowUpRequired,
		docUrlsJSON, metadataJSON, milestone.CompletedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead milestone: %w", err)
	}

	id, _ := result.LastInsertId()
	milestone.ID = id
	milestone.TenantID = tenantID
	milestone.CreatedAt = time.Now()
	milestone.UpdatedAt = time.Now()

	return milestone, nil
}

// GetLeadMilestones retrieves all milestones for a lead
func (s *milestoneCustomizationService) GetLeadMilestones(ctx context.Context, tenantID string, leadID int64) ([]LeadMilestone, error) {
	query := `
		SELECT id, tenant_id, lead_id, milestone_type_id, type_code, type_name,
		       achieved_date, achieved_time, status, days_from_previous, notes,
		       location_latitude, location_longitude, location_name, duration_minutes,
		       outcome, follow_up_date, follow_up_required, document_urls, metadata,
		       completed_by, created_at, updated_at
		FROM lead_milestones
		WHERE tenant_id = ? AND lead_id = ?
		ORDER BY achieved_date ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, leadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead milestones: %w", err)
	}
	defer rows.Close()

	var milestones []LeadMilestone
	for rows.Next() {
		var m LeadMilestone
		var docUrlsStr, metadataStr sql.NullString

		err := rows.Scan(
			&m.ID, &m.TenantID, &m.LeadID, &m.MilestoneTypeID, &m.TypeCode,
			&m.TypeName, &m.AchievedDate, &m.AchievedTime, &m.Status,
			&m.DaysFromPrevious, &m.Notes, &m.LocationLatitude,
			&m.LocationLongitude, &m.LocationName, &m.DurationMinutes,
			&m.Outcome, &m.FollowUpDate, &m.FollowUpRequired, &docUrlsStr,
			&metadataStr, &m.CompletedBy, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead milestone: %w", err)
		}

		if docUrlsStr.Valid {
			m.DocumentURLs = json.RawMessage(docUrlsStr.String)
		}

		if metadataStr.Valid {
			m.Metadata = json.RawMessage(metadataStr.String)
		}

		milestones = append(milestones, m)
	}

	return milestones, rows.Err()
}

// UpdateLeadMilestoneStatus updates the status of a lead milestone
func (s *milestoneCustomizationService) UpdateLeadMilestoneStatus(ctx context.Context, tenantID string, milestoneID int64, status string) error {
	query := `
		UPDATE lead_milestones
		SET status = ?, updated_at = NOW()
		WHERE tenant_id = ? AND id = ?
	`

	result, err := s.db.ExecContext(ctx, query, status, tenantID, milestoneID)
	if err != nil {
		return fmt.Errorf("failed to update milestone status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("milestone not found")
	}

	return nil
}

// GetMilestoneTimeline retrieves ordered milestone timeline for a lead
func (s *milestoneCustomizationService) GetMilestoneTimeline(ctx context.Context, tenantID string, leadID int64) ([]LeadMilestone, error) {
	return s.GetLeadMilestones(ctx, tenantID, leadID)
}

// ============================================================================
// ANALYTICS
// ============================================================================

// GetMilestoneCompletionMetrics retrieves milestone completion metrics
func (s *milestoneCustomizationService) GetMilestoneCompletionMetrics(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	query := `
		SELECT
			tmt.type_code,
			tmt.type_name,
			COUNT(lm.id) as total_milestones,
			SUM(CASE WHEN lm.status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN lm.status IN ('pending', 'in_progress') THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN lm.status = 'skipped' THEN 1 ELSE 0 END) as skipped,
			ROUND(SUM(CASE WHEN lm.status = 'completed' THEN 1 ELSE 0 END) / COUNT(lm.id) * 100, 2) as completion_rate
		FROM tenant_milestone_types tmt
		LEFT JOIN lead_milestones lm ON tmt.id = lm.milestone_type_id
		WHERE tmt.tenant_id = ?
		GROUP BY tmt.type_code, tmt.type_name
		ORDER BY completion_rate DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone metrics: %w", err)
	}
	defer rows.Close()

	metrics := make(map[string]interface{})
	details := make([]map[string]interface{}, 0)

	for rows.Next() {
		var typeCode, typeName string
		var total, completed, inProgress, skipped int64
		var completionRate float64

		err := rows.Scan(&typeCode, &typeName, &total, &completed, &inProgress, &skipped, &completionRate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan milestone metric: %w", err)
		}

		detail := map[string]interface{}{
			"type_code":        typeCode,
			"type_name":        typeName,
			"total":            total,
			"completed":        completed,
			"in_progress":      inProgress,
			"skipped":          skipped,
			"completion_rate":  completionRate,
		}

		details = append(details, detail)
	}

	metrics["details"] = details
	metrics["timestamp"] = time.Now()
	return metrics, rows.Err()
}

// GetMilestoneTimeTrends retrieves time-based milestone trends
func (s *milestoneCustomizationService) GetMilestoneTimeTrends(ctx context.Context, tenantID string, milestoneType string, days int) ([]MilestoneTimeTrendData, error) {
	query := `
		SELECT
			DATE(lm.achieved_date) as date,
			tmt.type_code,
			COUNT(lm.id) as milestone_count,
			ROUND(AVG(DATEDIFF(lm.achieved_date, l.created_at)), 2) as average_time_to_hit,
			ROUND(SUM(CASE WHEN lm.status = 'completed' THEN 1 ELSE 0 END) / COUNT(lm.id) * 100, 2) as completion_rate,
			SUM(CASE WHEN lm.follow_up_required = TRUE THEN 1 ELSE 0 END) as follow_up_required
		FROM lead_milestones lm
		JOIN tenant_milestone_types tmt ON lm.milestone_type_id = tmt.id
		JOIN lead l ON lm.lead_id = l.id
		WHERE lm.tenant_id = ? AND tmt.type_code = ? AND lm.achieved_date >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(lm.achieved_date), tmt.type_code
		ORDER BY date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, milestoneType, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get milestone trends: %w", err)
	}
	defer rows.Close()

	var trends []MilestoneTimeTrendData
	for rows.Next() {
		var data MilestoneTimeTrendData
		var followUp int64

		err := rows.Scan(&data.Date, &data.MilestoneType, &data.MilestoneCount,
			&data.AverageTimeToHit, &data.CompletionRate, &followUp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trend: %w", err)
		}

		data.FollowUpRequired = followUp
		trends = append(trends, data)
	}

	return trends, rows.Err()
}

// GetMilestoneBottlenecks identifies milestone transition bottlenecks
func (s *milestoneCustomizationService) GetMilestoneBottlenecks(ctx context.Context, tenantID string) ([]BottleneckAnalysis, error) {
	query := `
		SELECT
			tmt1.type_name as from_milestone,
			tmt2.type_name as to_milestone,
			ROUND(AVG(DATEDIFF(lm2.achieved_date, lm1.achieved_date)), 2) as average_days,
			ROUND(MEDIAN(DATEDIFF(lm2.achieved_date, lm1.achieved_date)), 2) as median_days,
			COUNT(DISTINCT lm1.lead_id) as leads_affected,
			ROUND(COUNT(DISTINCT CASE WHEN lm2.status = 'completed' THEN lm1.lead_id END) / 
				COUNT(DISTINCT lm1.lead_id) * 100, 2) as completion_rate,
			ROUND(SUM(CASE WHEN DATEDIFF(lm2.achieved_date, lm1.achieved_date) > tmt2.sla_days THEN 1 ELSE 0 END) /
				COUNT(lm1.id) * 100, 2) as sla_breach_percent
		FROM lead_milestones lm1
		JOIN lead_milestones lm2 ON lm1.lead_id = lm2.lead_id AND lm2.achieved_date > lm1.achieved_date
		JOIN tenant_milestone_types tmt1 ON lm1.milestone_type_id = tmt1.id
		JOIN tenant_milestone_types tmt2 ON lm2.milestone_type_id = tmt2.id
		WHERE lm1.tenant_id = ? AND lm1.milestone_type_id < lm2.milestone_type_id
		GROUP BY lm1.milestone_type_id, lm2.milestone_type_id
		HAVING average_days > 7
		ORDER BY average_days DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bottlenecks: %w", err)
	}
	defer rows.Close()

	var bottlenecks []BottleneckAnalysis
	for rows.Next() {
		var ba BottleneckAnalysis
		err := rows.Scan(&ba.FromMilestone, &ba.ToMilestone, &ba.AverageDays,
			&ba.MedianDays, &ba.LeadsAffected, &ba.CompletionRate, &ba.SLABreachPercent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bottleneck: %w", err)
		}
		bottlenecks = append(bottlenecks, ba)
	}

	return bottlenecks, rows.Err()
}
