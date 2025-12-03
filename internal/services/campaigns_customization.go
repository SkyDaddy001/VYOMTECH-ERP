package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// CampaignCustomizationService defines customization for campaigns
type CampaignCustomizationService interface {
	// Campaign Type Management
	CreateCampaignType(ctx context.Context, tenantID string, campaignType *CampaignTypeConfig) (*CampaignTypeConfig, error)
	GetCampaignType(ctx context.Context, tenantID, typeCode string) (*CampaignTypeConfig, error)
	GetCampaignTypes(ctx context.Context, tenantID string) ([]CampaignTypeConfig, error)
	UpdateCampaignType(ctx context.Context, tenantID string, campaignType *CampaignTypeConfig) (*CampaignTypeConfig, error)
	DeactivateCampaignType(ctx context.Context, tenantID, typeCode string) error

	// Campaign Channel Management
	CreateCampaignChannel(ctx context.Context, tenantID string, channel *CampaignChannelConfig) (*CampaignChannelConfig, error)
	GetCampaignChannels(ctx context.Context, tenantID string) ([]CampaignChannelConfig, error)
	UpdateCampaignChannel(ctx context.Context, tenantID string, channel *CampaignChannelConfig) (*CampaignChannelConfig, error)
	DeactivateCampaignChannel(ctx context.Context, tenantID, channelCode string) error

	// Campaign Status Management
	CreateCampaignStatus(ctx context.Context, tenantID string, status *CampaignStatusConfig) (*CampaignStatusConfig, error)
	GetCampaignStatuses(ctx context.Context, tenantID string) ([]CampaignStatusConfig, error)
	UpdateCampaignStatus(ctx context.Context, tenantID string, status *CampaignStatusConfig) (*CampaignStatusConfig, error)

	// Campaign Budget Types
	CreateBudgetType(ctx context.Context, tenantID string, budgetType *CampaignBudgetType) (*CampaignBudgetType, error)
	GetBudgetTypes(ctx context.Context, tenantID string) ([]CampaignBudgetType, error)
	UpdateBudgetType(ctx context.Context, tenantID string, budgetType *CampaignBudgetType) (*CampaignBudgetType, error)

	// Campaign Template Management
	CreateCampaignTemplate(ctx context.Context, tenantID string, template *CampaignTemplate) (*CampaignTemplate, error)
	GetCampaignTemplate(ctx context.Context, tenantID string, templateID int64) (*CampaignTemplate, error)
	GetCampaignTemplates(ctx context.Context, tenantID string) ([]CampaignTemplate, error)
	UpdateCampaignTemplate(ctx context.Context, tenantID string, template *CampaignTemplate) (*CampaignTemplate, error)
	DeleteCampaignTemplate(ctx context.Context, tenantID string, templateID int64) error

	// Campaign Metrics
	GetCampaignTypePerformance(ctx context.Context, tenantID, typeCode string) (*CampaignTypePerformance, error)
	GetChannelPerformance(ctx context.Context, tenantID, channelCode string, startDate, endDate time.Time) (*ChannelPerformanceMetrics, error)
	GetCampaignTrends(ctx context.Context, tenantID string, days int) (map[string]interface{}, error)
}

// ============================================================================
// DATA MODELS
// ============================================================================

// CampaignTypeConfig defines a campaign type
type CampaignTypeConfig struct {
	ID               int64           `json:"id" db:"id"`
	TenantID         string          `json:"tenant_id" db:"tenant_id"`
	TypeCode         string          `json:"type_code" db:"type_code"`
	TypeName         string          `json:"type_name" db:"type_name"`
	Description      *string         `json:"description,omitempty" db:"description"`
	Icon             *string         `json:"icon,omitempty" db:"icon"`
	ColorHex         *string         `json:"color_hex,omitempty" db:"color_hex"`
	DisplayOrder     int             `json:"display_order" db:"display_order"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	IsDefault        bool            `json:"is_default" db:"is_default"`
	TypicalDuration  *int            `json:"typical_duration_days,omitempty" db:"typical_duration_days"`
	MinBudget        *float64        `json:"min_budget,omitempty" db:"min_budget"`
	MaxBudget        *float64        `json:"max_budget,omitempty" db:"max_budget"`
	RecommendedBudget *float64       `json:"recommended_budget,omitempty" db:"recommended_budget"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedBy        *int64          `json:"created_by,omitempty" db:"created_by"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// CampaignChannelConfig defines a campaign channel
type CampaignChannelConfig struct {
	ID               int64           `json:"id" db:"id"`
	TenantID         string          `json:"tenant_id" db:"tenant_id"`
	ChannelCode      string          `json:"channel_code" db:"channel_code"`
	ChannelName      string          `json:"channel_name" db:"channel_name"`
	Description      *string         `json:"description,omitempty" db:"description"`
	Icon             *string         `json:"icon,omitempty" db:"icon"`
	DisplayOrder     int             `json:"display_order" db:"display_order"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	IsDefault        bool            `json:"is_default" db:"is_default"`
	AverageCPL       *float64        `json:"average_cpl,omitempty" db:"average_cpl"`
	AverageCPM       *float64        `json:"average_cpm,omitempty" db:"average_cpm"`
	AverageROI       *float64        `json:"average_roi,omitempty" db:"average_roi"`
	IntegrationKey   *string         `json:"integration_key,omitempty" db:"integration_key"`
	IntegrationData  json.RawMessage `json:"integration_data,omitempty" db:"integration_data"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// CampaignStatusConfig defines campaign statuses
type CampaignStatusConfig struct {
	ID               int64           `json:"id" db:"id"`
	TenantID         string          `json:"tenant_id" db:"tenant_id"`
	StatusCode       string          `json:"status_code" db:"status_code"`
	StatusName       string          `json:"status_name" db:"status_name"`
	Description      *string         `json:"description,omitempty" db:"description"`
	ColorHex         *string         `json:"color_hex,omitempty" db:"color_hex"`
	DisplayOrder     int             `json:"display_order" db:"display_order"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	IsInitial        bool            `json:"is_initial" db:"is_initial"`
	IsFinal          bool            `json:"is_final" db:"is_final"`
	AllowsEditing    bool            `json:"allows_editing" db:"allows_editing"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// CampaignBudgetType defines budget types
type CampaignBudgetType struct {
	ID               int64           `json:"id" db:"id"`
	TenantID         string          `json:"tenant_id" db:"tenant_id"`
	BudgetTypeCode   string          `json:"budget_type_code" db:"budget_type_code"`
	BudgetTypeName   string          `json:"budget_type_name" db:"budget_type_name"`
	Description      *string         `json:"description,omitempty" db:"description"`
	DisplayOrder     int             `json:"display_order" db:"display_order"`
	IsActive         bool            `json:"is_active" db:"is_active"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at" db:"updated_at"`
}

// CampaignTemplate defines a campaign template
type CampaignTemplate struct {
	ID                   int64           `json:"id" db:"id"`
	TenantID             string          `json:"tenant_id" db:"tenant_id"`
	TemplateName         string          `json:"template_name" db:"template_name"`
	Description          *string         `json:"description,omitempty" db:"description"`
	CampaignTypeCode     string          `json:"campaign_type_code" db:"campaign_type_code"`
	DefaultChannels      json.RawMessage `json:"default_channels" db:"default_channels"` // JSON array
	DefaultBudget        float64         `json:"default_budget" db:"default_budget"`
	DefaultDurationDays  int             `json:"default_duration_days" db:"default_duration_days"`
	TargetAudience       json.RawMessage `json:"target_audience,omitempty" db:"target_audience"` // JSON criteria
	KPIs                 json.RawMessage `json:"kpis,omitempty" db:"kpis"` // JSON array
	IsActive             bool            `json:"is_active" db:"is_active"`
	IsDefault            bool            `json:"is_default" db:"is_default"`
	CreatedBy            *int64          `json:"created_by,omitempty" db:"created_by"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// CampaignTypePerformance tracks performance by type
type CampaignTypePerformance struct {
	TypeCode         string    `json:"type_code"`
	TypeName         string    `json:"type_name"`
	TotalCampaigns   int64     `json:"total_campaigns"`
	ActiveCampaigns  int64     `json:"active_campaigns"`
	CompletedCampaigns int64   `json:"completed_campaigns"`
	TotalBudget      float64   `json:"total_budget"`
	TotalSpent       float64   `json:"total_spent"`
	TotalLeads       int64     `json:"total_leads"`
	TotalConversions int64     `json:"total_conversions"`
	AverageCPL       float64   `json:"average_cpl"`
	AverageROI       float64   `json:"average_roi"`
	ConversionRate   float64   `json:"conversion_rate"`
}

// ChannelPerformanceMetrics tracks performance by channel
type ChannelPerformanceMetrics struct {
	ChannelCode      string    `json:"channel_code"`
	ChannelName      string    `json:"channel_name"`
	CampaignCount    int64     `json:"campaign_count"`
	TotalBudget      float64   `json:"total_budget"`
	TotalSpent       float64   `json:"total_spent"`
	TotalImpressions int64     `json:"total_impressions"`
	TotalClicks      int64     `json:"total_clicks"`
	TotalLeads       int64     `json:"total_leads"`
	TotalConversions int64     `json:"total_conversions"`
	CPL              float64   `json:"cpl"`
	CPM              float64   `json:"cpm"`
	CTR              float64   `json:"ctr"`
	ConversionRate   float64   `json:"conversion_rate"`
	ROI              float64   `json:"roi"`
	CPC              float64   `json:"cpc"`
}

// ============================================================================
// SERVICE IMPLEMENTATION
// ============================================================================

type campaignCustomizationService struct {
	db *sql.DB
}

// NewCampaignCustomizationService creates a new campaign customization service
func NewCampaignCustomizationService(db *sql.DB) CampaignCustomizationService {
	return &campaignCustomizationService{db: db}
}

// ============================================================================
// CAMPAIGN TYPE OPERATIONS
// ============================================================================

// CreateCampaignType creates a new campaign type
func (s *campaignCustomizationService) CreateCampaignType(ctx context.Context, tenantID string, campaignType *CampaignTypeConfig) (*CampaignTypeConfig, error) {
	if campaignType.TypeCode == "" || campaignType.TypeName == "" {
		return nil, errors.New("type_code and type_name are required")
	}

	query := `
		INSERT INTO tenant_campaign_types (
			tenant_id, type_code, type_name, description, icon, color_hex,
			display_order, is_active, is_default, typical_duration_days,
			min_budget, max_budget, recommended_budget, metadata, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	metadataJSON := "{}"
	if campaignType.Metadata != nil {
		metadataJSON = string(campaignType.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, campaignType.TypeCode, campaignType.TypeName, campaignType.Description,
		campaignType.Icon, campaignType.ColorHex, campaignType.DisplayOrder,
		campaignType.IsActive, campaignType.IsDefault, campaignType.TypicalDuration,
		campaignType.MinBudget, campaignType.MaxBudget, campaignType.RecommendedBudget,
		metadataJSON, campaignType.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign type: %w", err)
	}

	id, _ := result.LastInsertId()
	campaignType.ID = id
	campaignType.TenantID = tenantID
	campaignType.CreatedAt = time.Now()
	campaignType.UpdatedAt = time.Now()

	return campaignType, nil
}

// GetCampaignType retrieves a specific campaign type
func (s *campaignCustomizationService) GetCampaignType(ctx context.Context, tenantID, typeCode string) (*CampaignTypeConfig, error) {
	query := `
		SELECT id, tenant_id, type_code, type_name, description, icon, color_hex,
		       display_order, is_active, is_default, typical_duration_days,
		       min_budget, max_budget, recommended_budget, metadata, created_by, created_at, updated_at
		FROM tenant_campaign_types
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	var ct CampaignTypeConfig
	var metadataStr sql.NullString

	err := s.db.QueryRowContext(ctx, query, tenantID, typeCode).Scan(
		&ct.ID, &ct.TenantID, &ct.TypeCode, &ct.TypeName, &ct.Description,
		&ct.Icon, &ct.ColorHex, &ct.DisplayOrder, &ct.IsActive, &ct.IsDefault,
		&ct.TypicalDuration, &ct.MinBudget, &ct.MaxBudget, &ct.RecommendedBudget,
		&metadataStr, &ct.CreatedBy, &ct.CreatedAt, &ct.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("campaign type not found: %s", typeCode)
		}
		return nil, fmt.Errorf("failed to get campaign type: %w", err)
	}

	if metadataStr.Valid {
		ct.Metadata = json.RawMessage(metadataStr.String)
	}

	return &ct, nil
}

// GetCampaignTypes retrieves all campaign types
func (s *campaignCustomizationService) GetCampaignTypes(ctx context.Context, tenantID string) ([]CampaignTypeConfig, error) {
	query := `
		SELECT id, tenant_id, type_code, type_name, description, icon, color_hex,
		       display_order, is_active, is_default, typical_duration_days,
		       min_budget, max_budget, recommended_budget, metadata, created_by, created_at, updated_at
		FROM tenant_campaign_types
		WHERE tenant_id = ? AND deleted_at IS NULL AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign types: %w", err)
	}
	defer rows.Close()

	var types []CampaignTypeConfig
	for rows.Next() {
		var ct CampaignTypeConfig
		var metadataStr sql.NullString

		err := rows.Scan(
			&ct.ID, &ct.TenantID, &ct.TypeCode, &ct.TypeName, &ct.Description,
			&ct.Icon, &ct.ColorHex, &ct.DisplayOrder, &ct.IsActive, &ct.IsDefault,
			&ct.TypicalDuration, &ct.MinBudget, &ct.MaxBudget, &ct.RecommendedBudget,
			&metadataStr, &ct.CreatedBy, &ct.CreatedAt, &ct.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign type: %w", err)
		}

		if metadataStr.Valid {
			ct.Metadata = json.RawMessage(metadataStr.String)
		}

		types = append(types, ct)
	}

	return types, rows.Err()
}

// UpdateCampaignType updates a campaign type
func (s *campaignCustomizationService) UpdateCampaignType(ctx context.Context, tenantID string, campaignType *CampaignTypeConfig) (*CampaignTypeConfig, error) {
	query := `
		UPDATE tenant_campaign_types
		SET type_name = ?, description = ?, icon = ?, color_hex = ?,
		    display_order = ?, is_active = ?, is_default = ?, typical_duration_days = ?,
		    min_budget = ?, max_budget = ?, recommended_budget = ?,
		    metadata = ?, updated_at = NOW()
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	metadataJSON := "{}"
	if campaignType.Metadata != nil {
		metadataJSON = string(campaignType.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		campaignType.TypeName, campaignType.Description, campaignType.Icon, campaignType.ColorHex,
		campaignType.DisplayOrder, campaignType.IsActive, campaignType.IsDefault,
		campaignType.TypicalDuration, campaignType.MinBudget, campaignType.MaxBudget,
		campaignType.RecommendedBudget, metadataJSON, tenantID, campaignType.TypeCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign type: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("campaign type not found")
	}

	campaignType.UpdatedAt = time.Now()
	return campaignType, nil
}

// DeactivateCampaignType deactivates a campaign type
func (s *campaignCustomizationService) DeactivateCampaignType(ctx context.Context, tenantID, typeCode string) error {
	query := `
		UPDATE tenant_campaign_types
		SET is_active = FALSE, updated_at = NOW()
		WHERE tenant_id = ? AND type_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, typeCode)
	if err != nil {
		return fmt.Errorf("failed to deactivate campaign type: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("campaign type not found")
	}

	return nil
}

// ============================================================================
// CAMPAIGN CHANNEL OPERATIONS
// ============================================================================

// CreateCampaignChannel creates a new campaign channel
func (s *campaignCustomizationService) CreateCampaignChannel(ctx context.Context, tenantID string, channel *CampaignChannelConfig) (*CampaignChannelConfig, error) {
	if channel.ChannelCode == "" || channel.ChannelName == "" {
		return nil, errors.New("channel_code and channel_name are required")
	}

	query := `
		INSERT INTO tenant_campaign_channels (
			tenant_id, channel_code, channel_name, description, icon,
			display_order, is_active, is_default, average_cpl, average_cpm,
			average_roi, integration_key, integration_data, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	integrationDataJSON := "{}"
	if channel.IntegrationData != nil {
		integrationDataJSON = string(channel.IntegrationData)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, channel.ChannelCode, channel.ChannelName, channel.Description,
		channel.Icon, channel.DisplayOrder, channel.IsActive, channel.IsDefault,
		channel.AverageCPL, channel.AverageCPM, channel.AverageROI,
		channel.IntegrationKey, integrationDataJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign channel: %w", err)
	}

	id, _ := result.LastInsertId()
	channel.ID = id
	channel.TenantID = tenantID
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	return channel, nil
}

// GetCampaignChannels retrieves all campaign channels
func (s *campaignCustomizationService) GetCampaignChannels(ctx context.Context, tenantID string) ([]CampaignChannelConfig, error) {
	query := `
		SELECT id, tenant_id, channel_code, channel_name, description, icon,
		       display_order, is_active, is_default, average_cpl, average_cpm,
		       average_roi, integration_key, integration_data, created_at, updated_at
		FROM tenant_campaign_channels
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign channels: %w", err)
	}
	defer rows.Close()

	var channels []CampaignChannelConfig
	for rows.Next() {
		var channel CampaignChannelConfig
		var integrationDataStr sql.NullString

		err := rows.Scan(
			&channel.ID, &channel.TenantID, &channel.ChannelCode, &channel.ChannelName,
			&channel.Description, &channel.Icon, &channel.DisplayOrder, &channel.IsActive,
			&channel.IsDefault, &channel.AverageCPL, &channel.AverageCPM,
			&channel.AverageROI, &channel.IntegrationKey, &integrationDataStr,
			&channel.CreatedAt, &channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign channel: %w", err)
		}

		if integrationDataStr.Valid {
			channel.IntegrationData = json.RawMessage(integrationDataStr.String)
		}

		channels = append(channels, channel)
	}

	return channels, rows.Err()
}

// UpdateCampaignChannel updates a campaign channel
func (s *campaignCustomizationService) UpdateCampaignChannel(ctx context.Context, tenantID string, channel *CampaignChannelConfig) (*CampaignChannelConfig, error) {
	query := `
		UPDATE tenant_campaign_channels
		SET channel_name = ?, description = ?, icon = ?, display_order = ?,
		    is_active = ?, is_default = ?, average_cpl = ?, average_cpm = ?,
		    average_roi = ?, integration_key = ?, integration_data = ?, updated_at = NOW()
		WHERE tenant_id = ? AND channel_code = ?
	`

	integrationDataJSON := "{}"
	if channel.IntegrationData != nil {
		integrationDataJSON = string(channel.IntegrationData)
	}

	result, err := s.db.ExecContext(ctx, query,
		channel.ChannelName, channel.Description, channel.Icon, channel.DisplayOrder,
		channel.IsActive, channel.IsDefault, channel.AverageCPL, channel.AverageCPM,
		channel.AverageROI, channel.IntegrationKey, integrationDataJSON,
		tenantID, channel.ChannelCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign channel: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("campaign channel not found")
	}

	channel.UpdatedAt = time.Now()
	return channel, nil
}

// DeactivateCampaignChannel deactivates a campaign channel
func (s *campaignCustomizationService) DeactivateCampaignChannel(ctx context.Context, tenantID, channelCode string) error {
	query := `
		UPDATE tenant_campaign_channels
		SET is_active = FALSE, updated_at = NOW()
		WHERE tenant_id = ? AND channel_code = ?
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, channelCode)
	if err != nil {
		return fmt.Errorf("failed to deactivate campaign channel: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("campaign channel not found")
	}

	return nil
}

// ============================================================================
// CAMPAIGN STATUS OPERATIONS
// ============================================================================

// CreateCampaignStatus creates a new campaign status
func (s *campaignCustomizationService) CreateCampaignStatus(ctx context.Context, tenantID string, status *CampaignStatusConfig) (*CampaignStatusConfig, error) {
	if status.StatusCode == "" || status.StatusName == "" {
		return nil, errors.New("status_code and status_name are required")
	}

	query := `
		INSERT INTO tenant_campaign_statuses (
			tenant_id, status_code, status_name, description, color_hex,
			display_order, is_active, is_initial, is_final, allows_editing,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, status.StatusCode, status.StatusName, status.Description,
		status.ColorHex, status.DisplayOrder, status.IsActive, status.IsInitial,
		status.IsFinal, status.AllowsEditing,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign status: %w", err)
	}

	id, _ := result.LastInsertId()
	status.ID = id
	status.TenantID = tenantID
	status.CreatedAt = time.Now()
	status.UpdatedAt = time.Now()

	return status, nil
}

// GetCampaignStatuses retrieves all campaign statuses
func (s *campaignCustomizationService) GetCampaignStatuses(ctx context.Context, tenantID string) ([]CampaignStatusConfig, error) {
	query := `
		SELECT id, tenant_id, status_code, status_name, description, color_hex,
		       display_order, is_active, is_initial, is_final, allows_editing, created_at, updated_at
		FROM tenant_campaign_statuses
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign statuses: %w", err)
	}
	defer rows.Close()

	var statuses []CampaignStatusConfig
	for rows.Next() {
		var status CampaignStatusConfig
		err := rows.Scan(
			&status.ID, &status.TenantID, &status.StatusCode, &status.StatusName,
			&status.Description, &status.ColorHex, &status.DisplayOrder, &status.IsActive,
			&status.IsInitial, &status.IsFinal, &status.AllowsEditing,
			&status.CreatedAt, &status.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign status: %w", err)
		}
		statuses = append(statuses, status)
	}

	return statuses, rows.Err()
}

// UpdateCampaignStatus updates a campaign status
func (s *campaignCustomizationService) UpdateCampaignStatus(ctx context.Context, tenantID string, status *CampaignStatusConfig) (*CampaignStatusConfig, error) {
	query := `
		UPDATE tenant_campaign_statuses
		SET status_name = ?, description = ?, color_hex = ?, display_order = ?,
		    is_active = ?, is_initial = ?, is_final = ?, allows_editing = ?, updated_at = NOW()
		WHERE tenant_id = ? AND status_code = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		status.StatusName, status.Description, status.ColorHex, status.DisplayOrder,
		status.IsActive, status.IsInitial, status.IsFinal, status.AllowsEditing,
		tenantID, status.StatusCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign status: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("campaign status not found")
	}

	status.UpdatedAt = time.Now()
	return status, nil
}

// ============================================================================
// CAMPAIGN BUDGET TYPE OPERATIONS
// ============================================================================

// CreateBudgetType creates a new budget type
func (s *campaignCustomizationService) CreateBudgetType(ctx context.Context, tenantID string, budgetType *CampaignBudgetType) (*CampaignBudgetType, error) {
	if budgetType.BudgetTypeCode == "" || budgetType.BudgetTypeName == "" {
		return nil, errors.New("budget_type_code and budget_type_name are required")
	}

	query := `
		INSERT INTO tenant_campaign_budget_types (
			tenant_id, budget_type_code, budget_type_name, description,
			display_order, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, budgetType.BudgetTypeCode, budgetType.BudgetTypeName,
		budgetType.Description, budgetType.DisplayOrder, budgetType.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget type: %w", err)
	}

	id, _ := result.LastInsertId()
	budgetType.ID = id
	budgetType.TenantID = tenantID
	budgetType.CreatedAt = time.Now()
	budgetType.UpdatedAt = time.Now()

	return budgetType, nil
}

// GetBudgetTypes retrieves all budget types
func (s *campaignCustomizationService) GetBudgetTypes(ctx context.Context, tenantID string) ([]CampaignBudgetType, error) {
	query := `
		SELECT id, tenant_id, budget_type_code, budget_type_name, description,
		       display_order, is_active, created_at, updated_at
		FROM tenant_campaign_budget_types
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget types: %w", err)
	}
	defer rows.Close()

	var types []CampaignBudgetType
	for rows.Next() {
		var bt CampaignBudgetType
		err := rows.Scan(
			&bt.ID, &bt.TenantID, &bt.BudgetTypeCode, &bt.BudgetTypeName,
			&bt.Description, &bt.DisplayOrder, &bt.IsActive,
			&bt.CreatedAt, &bt.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget type: %w", err)
		}
		types = append(types, bt)
	}

	return types, rows.Err()
}

// UpdateBudgetType updates a budget type
func (s *campaignCustomizationService) UpdateBudgetType(ctx context.Context, tenantID string, budgetType *CampaignBudgetType) (*CampaignBudgetType, error) {
	query := `
		UPDATE tenant_campaign_budget_types
		SET budget_type_name = ?, description = ?, display_order = ?,
		    is_active = ?, updated_at = NOW()
		WHERE tenant_id = ? AND budget_type_code = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		budgetType.BudgetTypeName, budgetType.Description, budgetType.DisplayOrder,
		budgetType.IsActive, tenantID, budgetType.BudgetTypeCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget type: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("budget type not found")
	}

	budgetType.UpdatedAt = time.Now()
	return budgetType, nil
}

// ============================================================================
// CAMPAIGN TEMPLATE OPERATIONS
// ============================================================================

// CreateCampaignTemplate creates a new campaign template
func (s *campaignCustomizationService) CreateCampaignTemplate(ctx context.Context, tenantID string, template *CampaignTemplate) (*CampaignTemplate, error) {
	if template.TemplateName == "" || template.CampaignTypeCode == "" {
		return nil, errors.New("template_name and campaign_type_code are required")
	}

	query := `
		INSERT INTO tenant_campaign_templates (
			tenant_id, template_name, description, campaign_type_code,
			default_channels, default_budget, default_duration_days,
			target_audience, kpis, is_active, is_default, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	defaultChannelsJSON := "[]"
	if template.DefaultChannels != nil {
		defaultChannelsJSON = string(template.DefaultChannels)
	}

	targetAudienceJSON := "{}"
	if template.TargetAudience != nil {
		targetAudienceJSON = string(template.TargetAudience)
	}

	kpisJSON := "[]"
	if template.KPIs != nil {
		kpisJSON = string(template.KPIs)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, template.TemplateName, template.Description, template.CampaignTypeCode,
		defaultChannelsJSON, template.DefaultBudget, template.DefaultDurationDays,
		targetAudienceJSON, kpisJSON, template.IsActive, template.IsDefault, template.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign template: %w", err)
	}

	id, _ := result.LastInsertId()
	template.ID = id
	template.TenantID = tenantID
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	return template, nil
}

// GetCampaignTemplate retrieves a specific campaign template
func (s *campaignCustomizationService) GetCampaignTemplate(ctx context.Context, tenantID string, templateID int64) (*CampaignTemplate, error) {
	query := `
		SELECT id, tenant_id, template_name, description, campaign_type_code,
		       default_channels, default_budget, default_duration_days,
		       target_audience, kpis, is_active, is_default, created_by, created_at, updated_at
		FROM tenant_campaign_templates
		WHERE tenant_id = ? AND id = ?
	`

	var template CampaignTemplate
	var defaultChannelsStr, targetAudienceStr, kpisStr sql.NullString

	err := s.db.QueryRowContext(ctx, query, tenantID, templateID).Scan(
		&template.ID, &template.TenantID, &template.TemplateName, &template.Description,
		&template.CampaignTypeCode, &defaultChannelsStr, &template.DefaultBudget,
		&template.DefaultDurationDays, &targetAudienceStr, &kpisStr,
		&template.IsActive, &template.IsDefault, &template.CreatedBy,
		&template.CreatedAt, &template.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("campaign template not found")
		}
		return nil, fmt.Errorf("failed to get campaign template: %w", err)
	}

	if defaultChannelsStr.Valid {
		template.DefaultChannels = json.RawMessage(defaultChannelsStr.String)
	}
	if targetAudienceStr.Valid {
		template.TargetAudience = json.RawMessage(targetAudienceStr.String)
	}
	if kpisStr.Valid {
		template.KPIs = json.RawMessage(kpisStr.String)
	}

	return &template, nil
}

// GetCampaignTemplates retrieves all campaign templates
func (s *campaignCustomizationService) GetCampaignTemplates(ctx context.Context, tenantID string) ([]CampaignTemplate, error) {
	query := `
		SELECT id, tenant_id, template_name, description, campaign_type_code,
		       default_channels, default_budget, default_duration_days,
		       target_audience, kpis, is_active, is_default, created_by, created_at, updated_at
		FROM tenant_campaign_templates
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY is_default DESC, template_name ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign templates: %w", err)
	}
	defer rows.Close()

	var templates []CampaignTemplate
	for rows.Next() {
		var template CampaignTemplate
		var defaultChannelsStr, targetAudienceStr, kpisStr sql.NullString

		err := rows.Scan(
			&template.ID, &template.TenantID, &template.TemplateName, &template.Description,
			&template.CampaignTypeCode, &defaultChannelsStr, &template.DefaultBudget,
			&template.DefaultDurationDays, &targetAudienceStr, &kpisStr,
			&template.IsActive, &template.IsDefault, &template.CreatedBy,
			&template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan campaign template: %w", err)
		}

		if defaultChannelsStr.Valid {
			template.DefaultChannels = json.RawMessage(defaultChannelsStr.String)
		}
		if targetAudienceStr.Valid {
			template.TargetAudience = json.RawMessage(targetAudienceStr.String)
		}
		if kpisStr.Valid {
			template.KPIs = json.RawMessage(kpisStr.String)
		}

		templates = append(templates, template)
	}

	return templates, rows.Err()
}

// UpdateCampaignTemplate updates a campaign template
func (s *campaignCustomizationService) UpdateCampaignTemplate(ctx context.Context, tenantID string, template *CampaignTemplate) (*CampaignTemplate, error) {
	query := `
		UPDATE tenant_campaign_templates
		SET template_name = ?, description = ?, campaign_type_code = ?,
		    default_channels = ?, default_budget = ?, default_duration_days = ?,
		    target_audience = ?, kpis = ?, is_active = ?, is_default = ?,
		    updated_at = NOW()
		WHERE tenant_id = ? AND id = ?
	`

	defaultChannelsJSON := "[]"
	if template.DefaultChannels != nil {
		defaultChannelsJSON = string(template.DefaultChannels)
	}

	targetAudienceJSON := "{}"
	if template.TargetAudience != nil {
		targetAudienceJSON = string(template.TargetAudience)
	}

	kpisJSON := "[]"
	if template.KPIs != nil {
		kpisJSON = string(template.KPIs)
	}

	result, err := s.db.ExecContext(ctx, query,
		template.TemplateName, template.Description, template.CampaignTypeCode,
		defaultChannelsJSON, template.DefaultBudget, template.DefaultDurationDays,
		targetAudienceJSON, kpisJSON, template.IsActive, template.IsDefault,
		tenantID, template.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign template: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("campaign template not found")
	}

	template.UpdatedAt = time.Now()
	return template, nil
}

// DeleteCampaignTemplate deletes a campaign template
func (s *campaignCustomizationService) DeleteCampaignTemplate(ctx context.Context, tenantID string, templateID int64) error {
	query := `DELETE FROM tenant_campaign_templates WHERE tenant_id = ? AND id = ?`

	result, err := s.db.ExecContext(ctx, query, tenantID, templateID)
	if err != nil {
		return fmt.Errorf("failed to delete campaign template: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("campaign template not found")
	}

	return nil
}

// ============================================================================
// ANALYTICS
// ============================================================================

// GetCampaignTypePerformance retrieves performance metrics by campaign type
func (s *campaignCustomizationService) GetCampaignTypePerformance(ctx context.Context, tenantID, typeCode string) (*CampaignTypePerformance, error) {
	query := `
		SELECT
			tct.type_code,
			tct.type_name,
			COUNT(DISTINCT c.id) as total_campaigns,
			SUM(CASE WHEN c.status = 'active' THEN 1 ELSE 0 END) as active_campaigns,
			SUM(CASE WHEN c.status = 'completed' THEN 1 ELSE 0 END) as completed_campaigns,
			COALESCE(SUM(c.budget), 0) as total_budget,
			COALESCE(SUM(c.spent_budget), 0) as total_spent,
			COALESCE(SUM(c.generated_leads), 0) as total_leads,
			COALESCE(SUM(c.converted_leads), 0) as total_conversions,
			ROUND(AVG(c.cost_per_lead), 2) as average_cpl,
			ROUND(AVG(c.conversion_rate), 2) as average_roi,
			ROUND(SUM(c.converted_leads) / SUM(c.generated_leads) * 100, 2) as conversion_rate
		FROM tenant_campaign_types tct
		LEFT JOIN campaign c ON c.campaign_type_id = tct.id
		WHERE tct.tenant_id = ? AND tct.type_code = ?
		GROUP BY tct.type_code, tct.type_name
	`

	var perf CampaignTypePerformance
	err := s.db.QueryRowContext(ctx, query, tenantID, typeCode).Scan(
		&perf.TypeCode, &perf.TypeName, &perf.TotalCampaigns, &perf.ActiveCampaigns,
		&perf.CompletedCampaigns, &perf.TotalBudget, &perf.TotalSpent, &perf.TotalLeads,
		&perf.TotalConversions, &perf.AverageCPL, &perf.AverageROI, &perf.ConversionRate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data for campaign type: %s", typeCode)
		}
		return nil, fmt.Errorf("failed to get campaign type performance: %w", err)
	}

	return &perf, nil
}

// GetChannelPerformance retrieves performance metrics by channel
func (s *campaignCustomizationService) GetChannelPerformance(ctx context.Context, tenantID, channelCode string, startDate, endDate time.Time) (*ChannelPerformanceMetrics, error) {
	query := `
		SELECT
			tcc.channel_code,
			tcc.channel_name,
			COUNT(DISTINCT c.id) as campaign_count,
			COALESCE(SUM(c.budget), 0) as total_budget,
			COALESCE(SUM(c.spent_budget), 0) as total_spent,
			COALESCE(SUM(cc.impressions), 0) as total_impressions,
			COALESCE(SUM(cc.clicks), 0) as total_clicks,
			COALESCE(SUM(cc.leads_generated), 0) as total_leads,
			COALESCE(SUM(cc.conversions), 0) as total_conversions
		FROM tenant_campaign_channels tcc
		LEFT JOIN campaign c ON c.campaign_id = tcc.id
		LEFT JOIN campaign_channel cc ON cc.campaign_id = c.id
		WHERE tcc.tenant_id = ? AND tcc.channel_code = ?
		AND c.created_at BETWEEN ? AND ?
		GROUP BY tcc.channel_code, tcc.channel_name
	`

	var metrics ChannelPerformanceMetrics
	err := s.db.QueryRowContext(ctx, query, tenantID, channelCode, startDate, endDate).Scan(
		&metrics.ChannelCode, &metrics.ChannelName, &metrics.CampaignCount,
		&metrics.TotalBudget, &metrics.TotalSpent, &metrics.TotalImpressions,
		&metrics.TotalClicks, &metrics.TotalLeads, &metrics.TotalConversions,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data for channel: %s", channelCode)
		}
		return nil, fmt.Errorf("failed to get channel performance: %w", err)
	}

	// Calculate derived metrics
	if metrics.TotalImpressions > 0 {
		metrics.CPM = (metrics.TotalSpent / float64(metrics.TotalImpressions)) * 1000
		metrics.CTR = (float64(metrics.TotalClicks) / float64(metrics.TotalImpressions)) * 100
	}
	if metrics.TotalLeads > 0 {
		metrics.CPL = metrics.TotalSpent / float64(metrics.TotalLeads)
	}
	if metrics.TotalClicks > 0 {
		metrics.CPC = metrics.TotalSpent / float64(metrics.TotalClicks)
	}
	if metrics.TotalLeads > 0 {
		metrics.ConversionRate = (float64(metrics.TotalConversions) / float64(metrics.TotalLeads)) * 100
	}
	if metrics.TotalSpent > 0 {
		metrics.ROI = ((float64(metrics.TotalConversions*100) - metrics.TotalSpent) / metrics.TotalSpent) * 100
	}

	return &metrics, nil
}

// GetCampaignTrends retrieves trend data for campaigns
func (s *campaignCustomizationService) GetCampaignTrends(ctx context.Context, tenantID string, days int) (map[string]interface{}, error) {
	query := `
		SELECT
			DATE(c.created_at) as date,
			COUNT(DISTINCT c.id) as campaigns_launched,
			SUM(c.budget) as budget_allocated,
			SUM(c.generated_leads) as leads_generated,
			SUM(c.converted_leads) as conversions,
			ROUND(AVG(c.cost_per_lead), 2) as avg_cpl
		FROM campaign c
		WHERE c.tenant_id = ? AND c.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(c.created_at)
		ORDER BY date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get campaign trends: %w", err)
	}
	defer rows.Close()

	trends := make(map[string]interface{})
	dailyData := make([]map[string]interface{}, 0)

	for rows.Next() {
		var date time.Time
		var campaignsLaunched, leadsGenerated, conversions int64
		var budgetAllocated, avgCPL sql.NullFloat64

		err := rows.Scan(&date, &campaignsLaunched, &budgetAllocated, &leadsGenerated, &conversions, &avgCPL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trend: %w", err)
		}

		dailyEntry := map[string]interface{}{
			"date":               date,
			"campaigns_launched": campaignsLaunched,
			"leads_generated":    leadsGenerated,
			"conversions":        conversions,
		}

		if budgetAllocated.Valid {
			dailyEntry["budget_allocated"] = budgetAllocated.Float64
		}
		if avgCPL.Valid {
			dailyEntry["avg_cpl"] = avgCPL.Float64
		}

		dailyData = append(dailyData, dailyEntry)
	}

	trends["daily_data"] = dailyData
	trends["period_days"] = days
	return trends, rows.Err()
}
