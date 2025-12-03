package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// LeadSourceCustomizationService defines customization for lead sources and subsources
type LeadSourceCustomizationService interface {
	// Lead Source Management
	CreateLeadSource(ctx context.Context, tenantID string, source *LeadSourceConfig) (*LeadSourceConfig, error)
	GetLeadSource(ctx context.Context, tenantID, sourceCode string) (*LeadSourceConfig, error)
	GetLeadSources(ctx context.Context, tenantID string, filter *LeadSourceFilter) ([]LeadSourceConfig, error)
	UpdateLeadSource(ctx context.Context, tenantID string, source *LeadSourceConfig) (*LeadSourceConfig, error)
	DeactivateLeadSource(ctx context.Context, tenantID, sourceCode string) error
	DeleteLeadSource(ctx context.Context, tenantID, sourceCode string) error

	// Lead SubSource Management
	CreateLeadSubSource(ctx context.Context, tenantID, sourceCode string, subsource *LeadSubSourceConfig) (*LeadSubSourceConfig, error)
	GetLeadSubSource(ctx context.Context, tenantID, sourceCode, subSourceCode string) (*LeadSubSourceConfig, error)
	GetLeadSubSources(ctx context.Context, tenantID, sourceCode string) ([]LeadSubSourceConfig, error)
	UpdateLeadSubSource(ctx context.Context, tenantID string, subsource *LeadSubSourceConfig) (*LeadSubSourceConfig, error)
	DeactivateLeadSubSource(ctx context.Context, tenantID, sourceCode, subSourceCode string) error
	DeleteLeadSubSource(ctx context.Context, tenantID, sourceCode, subSourceCode string) error

	// Channel Management
	CreateChannel(ctx context.Context, tenantID string, channel *ChannelConfig) (*ChannelConfig, error)
	GetChannels(ctx context.Context, tenantID string) ([]ChannelConfig, error)
	UpdateChannel(ctx context.Context, tenantID string, channel *ChannelConfig) (*ChannelConfig, error)

	// Lead Source Analytics
	GetSourcePerformance(ctx context.Context, tenantID, sourceCode string, startDate, endDate time.Time) (*SourcePerformanceMetrics, error)
	GetSubSourcePerformance(ctx context.Context, tenantID, sourceCode, subSourceCode string, startDate, endDate time.Time) (*SubSourcePerformanceMetrics, error)
	GetLeadSourceTrends(ctx context.Context, tenantID string, days int) (map[string]interface{}, error)
}

// ============================================================================
// DATA MODELS
// ============================================================================

// LeadSourceConfig defines a lead source
type LeadSourceConfig struct {
	ID               int64         `json:"id" db:"id"`
	TenantID         string        `json:"tenant_id" db:"tenant_id"`
	SourceCode       string        `json:"source_code" db:"source_code"`
	SourceName       string        `json:"source_name" db:"source_name"`
	SourceType       string        `json:"source_type" db:"source_type"` // website, email, phone, referral, event, social, direct, partner, import
	Description      *string       `json:"description,omitempty" db:"description"`
	Icon             *string       `json:"icon,omitempty" db:"icon"`
	ColorHex         *string       `json:"color_hex,omitempty" db:"color_hex"`
	DisplayOrder     int           `json:"display_order" db:"display_order"`
	IsActive         bool          `json:"is_active" db:"is_active"`
	IsDefault        bool          `json:"is_default" db:"is_default"`
	LeadsGenerated   int64         `json:"leads_generated" db:"leads_generated"`
	ConversionCount  int64         `json:"conversion_count" db:"conversion_count"`
	ConversionRate   float64       `json:"conversion_rate" db:"conversion_rate"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"` // Additional configuration
	CreatedBy        *int64        `json:"created_by,omitempty" db:"created_by"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
	SubSources       []LeadSubSourceConfig `json:"sub_sources,omitempty" db:"-"`
}

// LeadSubSourceConfig defines a subsource within a source
type LeadSubSourceConfig struct {
	ID               int64         `json:"id" db:"id"`
	TenantID         string        `json:"tenant_id" db:"tenant_id"`
	SourceID         int64         `json:"source_id" db:"source_id"`
	SourceCode       string        `json:"source_code" db:"source_code"`
	SubSourceCode    string        `json:"sub_source_code" db:"sub_source_code"`
	SubSourceName    string        `json:"sub_source_name" db:"sub_source_name"` // Google Ads, Facebook, LinkedIn, etc.
	Description      *string       `json:"description,omitempty" db:"description"`
	Icon             *string       `json:"icon,omitempty" db:"icon"`
	DisplayOrder     int           `json:"display_order" db:"display_order"`
	IsActive         bool          `json:"is_active" db:"is_active"`
	LeadsGenerated   int64         `json:"leads_generated" db:"leads_generated"`
	ConversionCount  int64         `json:"conversion_count" db:"conversion_count"`
	ConversionRate   float64       `json:"conversion_rate" db:"conversion_rate"`
	CostPerLead      float64       `json:"cost_per_lead" db:"cost_per_lead"`
	TotalCost        float64       `json:"total_cost" db:"total_cost"`
	LastActivityDate *time.Time    `json:"last_activity_date,omitempty" db:"last_activity_date"`
	Metadata         json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedBy        *int64        `json:"created_by,omitempty" db:"created_by"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
}

// ChannelConfig defines a channel for lead categorization
type ChannelConfig struct {
	ID           int64         `json:"id" db:"id"`
	TenantID     string        `json:"tenant_id" db:"tenant_id"`
	ChannelCode  string        `json:"channel_code" db:"channel_code"`
	ChannelName  string        `json:"channel_name" db:"channel_name"` // Direct, Organic, Paid, Referral, etc.
	Description  *string       `json:"description,omitempty" db:"description"`
	DisplayOrder int           `json:"display_order" db:"display_order"`
	IsActive     bool          `json:"is_active" db:"is_active"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}

// LeadSourceFilter for filtering sources
type LeadSourceFilter struct {
	SourceType string
	IsActive   *bool
	SearchText string
	Limit      int
	Offset     int
}

// SourcePerformanceMetrics tracks performance for a source
type SourcePerformanceMetrics struct {
	SourceCode       string    `json:"source_code"`
	SourceName       string    `json:"source_name"`
	LeadsGenerated   int64     `json:"leads_generated"`
	LeadsContacted   int64     `json:"leads_contacted"`
	LeadsQualified   int64     `json:"leads_qualified"`
	LeadsConverted   int64     `json:"leads_converted"`
	LeadsLost        int64     `json:"leads_lost"`
	ConversionRate   float64   `json:"conversion_rate"`
	QualificationRate float64   `json:"qualification_rate"`
	AverageDaysToClose float64  `json:"average_days_to_close"`
	TotalValue       float64   `json:"total_value"`
	CostPerLead      float64   `json:"cost_per_lead"`
	ROI              float64   `json:"roi"`
	SubSources       []SubSourcePerformanceMetrics `json:"sub_sources,omitempty"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// SubSourcePerformanceMetrics tracks performance for a subsource
type SubSourcePerformanceMetrics struct {
	SubSourceCode    string    `json:"sub_source_code"`
	SubSourceName    string    `json:"sub_source_name"`
	LeadsGenerated   int64     `json:"leads_generated"`
	LeadsContacted   int64     `json:"leads_contacted"`
	LeadsQualified   int64     `json:"leads_qualified"`
	LeadsConverted   int64     `json:"leads_converted"`
	LeadsLost        int64     `json:"leads_lost"`
	ConversionRate   float64   `json:"conversion_rate"`
	TotalCost        float64   `json:"total_cost"`
	CostPerLead      float64   `json:"cost_per_lead"`
	CostPerConversion float64   `json:"cost_per_conversion"`
	Revenue          float64   `json:"revenue"`
	ROI              float64   `json:"roi"`
}

// ============================================================================
// SERVICE IMPLEMENTATION
// ============================================================================

type leadSourceCustomizationService struct {
	db *sql.DB
}

// NewLeadSourceCustomizationService creates a new lead source customization service
func NewLeadSourceCustomizationService(db *sql.DB) LeadSourceCustomizationService {
	return &leadSourceCustomizationService{db: db}
}

// ============================================================================
// LEAD SOURCE CRUD OPERATIONS
// ============================================================================

// CreateLeadSource creates a new lead source configuration
func (s *leadSourceCustomizationService) CreateLeadSource(ctx context.Context, tenantID string, source *LeadSourceConfig) (*LeadSourceConfig, error) {
	if source.SourceCode == "" || source.SourceName == "" || source.SourceType == "" {
		return nil, errors.New("source_code, source_name, and source_type are required")
	}

	query := `
		INSERT INTO tenant_lead_sources (
			tenant_id, source_code, source_name, source_type, description, icon, 
			color_hex, display_order, is_active, is_default, leads_generated,
			conversion_count, conversion_rate, metadata, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	metadataJSON := "{}"
	if source.Metadata != nil {
		metadataJSON = string(source.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		tenantID, source.SourceCode, source.SourceName, source.SourceType,
		source.Description, source.Icon, source.ColorHex, source.DisplayOrder,
		source.IsActive, source.IsDefault, source.LeadsGenerated,
		source.ConversionCount, source.ConversionRate, metadataJSON,
		source.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead source: %w", err)
	}

	id, _ := result.LastInsertId()
	source.ID = id
	source.TenantID = tenantID
	source.CreatedAt = time.Now()
	source.UpdatedAt = time.Now()

	return source, nil
}

// GetLeadSource retrieves a specific lead source
func (s *leadSourceCustomizationService) GetLeadSource(ctx context.Context, tenantID, sourceCode string) (*LeadSourceConfig, error) {
	query := `
		SELECT id, tenant_id, source_code, source_name, source_type, description, 
		       icon, color_hex, display_order, is_active, is_default, leads_generated,
		       conversion_count, conversion_rate, metadata, created_by, created_at, updated_at
		FROM tenant_lead_sources
		WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL
	`

	var source LeadSourceConfig
	var metadataStr sql.NullString

	err := s.db.QueryRowContext(ctx, query, tenantID, sourceCode).Scan(
		&source.ID, &source.TenantID, &source.SourceCode, &source.SourceName,
		&source.SourceType, &source.Description, &source.Icon, &source.ColorHex,
		&source.DisplayOrder, &source.IsActive, &source.IsDefault,
		&source.LeadsGenerated, &source.ConversionCount, &source.ConversionRate,
		&metadataStr, &source.CreatedBy, &source.CreatedAt, &source.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lead source not found: %s", sourceCode)
		}
		return nil, fmt.Errorf("failed to get lead source: %w", err)
	}

	if metadataStr.Valid {
		source.Metadata = json.RawMessage(metadataStr.String)
	}

	return &source, nil
}

// GetLeadSources retrieves all lead sources with optional filtering
func (s *leadSourceCustomizationService) GetLeadSources(ctx context.Context, tenantID string, filter *LeadSourceFilter) ([]LeadSourceConfig, error) {
	query := `
		SELECT id, tenant_id, source_code, source_name, source_type, description, 
		       icon, color_hex, display_order, is_active, is_default, leads_generated,
		       conversion_count, conversion_rate, metadata, created_by, created_at, updated_at
		FROM tenant_lead_sources
		WHERE tenant_id = ? AND deleted_at IS NULL
	`

	args := []interface{}{tenantID}

	if filter != nil {
		if filter.SourceType != "" {
			query += " AND source_type = ?"
			args = append(args, filter.SourceType)
		}
		if filter.IsActive != nil {
			query += " AND is_active = ?"
			args = append(args, *filter.IsActive)
		}
		if filter.SearchText != "" {
			query += " AND (source_name LIKE ? OR source_code LIKE ?)"
			searchTerm := "%" + filter.SearchText + "%"
			args = append(args, searchTerm, searchTerm)
		}

		query += " ORDER BY display_order ASC, source_name ASC"

		if filter.Limit > 0 {
			query += " LIMIT ? OFFSET ?"
			args = append(args, filter.Limit, filter.Offset)
		}
	} else {
		query += " ORDER BY display_order ASC"
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead sources: %w", err)
	}
	defer rows.Close()

	var sources []LeadSourceConfig
	for rows.Next() {
		var source LeadSourceConfig
		var metadataStr sql.NullString

		err := rows.Scan(
			&source.ID, &source.TenantID, &source.SourceCode, &source.SourceName,
			&source.SourceType, &source.Description, &source.Icon, &source.ColorHex,
			&source.DisplayOrder, &source.IsActive, &source.IsDefault,
			&source.LeadsGenerated, &source.ConversionCount, &source.ConversionRate,
			&metadataStr, &source.CreatedBy, &source.CreatedAt, &source.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead source: %w", err)
		}

		if metadataStr.Valid {
			source.Metadata = json.RawMessage(metadataStr.String)
		}

		sources = append(sources, source)
	}

	return sources, rows.Err()
}

// UpdateLeadSource updates an existing lead source
func (s *leadSourceCustomizationService) UpdateLeadSource(ctx context.Context, tenantID string, source *LeadSourceConfig) (*LeadSourceConfig, error) {
	query := `
		UPDATE tenant_lead_sources
		SET source_name = ?, source_type = ?, description = ?, icon = ?, 
		    color_hex = ?, display_order = ?, is_active = ?, is_default = ?,
		    metadata = ?, updated_at = NOW()
		WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL
	`

	metadataJSON := "{}"
	if source.Metadata != nil {
		metadataJSON = string(source.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		source.SourceName, source.SourceType, source.Description, source.Icon,
		source.ColorHex, source.DisplayOrder, source.IsActive, source.IsDefault,
		metadataJSON, tenantID, source.SourceCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update lead source: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("lead source not found")
	}

	source.UpdatedAt = time.Now()
	return source, nil
}

// DeactivateLeadSource deactivates a lead source
func (s *leadSourceCustomizationService) DeactivateLeadSource(ctx context.Context, tenantID, sourceCode string) error {
	query := `
		UPDATE tenant_lead_sources
		SET is_active = FALSE, updated_at = NOW()
		WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, sourceCode)
	if err != nil {
		return fmt.Errorf("failed to deactivate lead source: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("lead source not found")
	}

	return nil
}

// DeleteLeadSource soft deletes a lead source
func (s *leadSourceCustomizationService) DeleteLeadSource(ctx context.Context, tenantID, sourceCode string) error {
	query := `
		UPDATE tenant_lead_sources
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, sourceCode)
	if err != nil {
		return fmt.Errorf("failed to delete lead source: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("lead source not found")
	}

	return nil
}

// ============================================================================
// LEAD SUBSOURCE CRUD OPERATIONS
// ============================================================================

// CreateLeadSubSource creates a new lead subsource
func (s *leadSourceCustomizationService) CreateLeadSubSource(ctx context.Context, tenantID, sourceCode string, subsource *LeadSubSourceConfig) (*LeadSubSourceConfig, error) {
	if subsource.SubSourceCode == "" || subsource.SubSourceName == "" {
		return nil, errors.New("sub_source_code and sub_source_name are required")
	}

	// Get source ID
	var sourceID int64
	query := `SELECT id FROM tenant_lead_sources WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL`
	err := s.db.QueryRowContext(ctx, query, tenantID, sourceCode).Scan(&sourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lead source not found: %s", sourceCode)
		}
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	insertQuery := `
		INSERT INTO tenant_lead_subsources (
			tenant_id, source_id, source_code, sub_source_code, sub_source_name,
			description, icon, display_order, is_active, leads_generated,
			conversion_count, conversion_rate, cost_per_lead, total_cost,
			last_activity_date, metadata, created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	metadataJSON := "{}"
	if subsource.Metadata != nil {
		metadataJSON = string(subsource.Metadata)
	}

	result, err := s.db.ExecContext(ctx, insertQuery,
		tenantID, sourceID, sourceCode, subsource.SubSourceCode, subsource.SubSourceName,
		subsource.Description, subsource.Icon, subsource.DisplayOrder, subsource.IsActive,
		subsource.LeadsGenerated, subsource.ConversionCount, subsource.ConversionRate,
		subsource.CostPerLead, subsource.TotalCost, subsource.LastActivityDate,
		metadataJSON, subsource.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lead subsource: %w", err)
	}

	id, _ := result.LastInsertId()
	subsource.ID = id
	subsource.TenantID = tenantID
	subsource.SourceCode = sourceCode
	subsource.SourceID = sourceID
	subsource.CreatedAt = time.Now()
	subsource.UpdatedAt = time.Now()

	return subsource, nil
}

// GetLeadSubSource retrieves a specific subsource
func (s *leadSourceCustomizationService) GetLeadSubSource(ctx context.Context, tenantID, sourceCode, subSourceCode string) (*LeadSubSourceConfig, error) {
	query := `
		SELECT id, tenant_id, source_id, source_code, sub_source_code, sub_source_name,
		       description, icon, display_order, is_active, leads_generated,
		       conversion_count, conversion_rate, cost_per_lead, total_cost,
		       last_activity_date, metadata, created_by, created_at, updated_at
		FROM tenant_lead_subsources
		WHERE tenant_id = ? AND source_code = ? AND sub_source_code = ? AND deleted_at IS NULL
	`

	var subsource LeadSubSourceConfig
	var metadataStr sql.NullString

	err := s.db.QueryRowContext(ctx, query, tenantID, sourceCode, subSourceCode).Scan(
		&subsource.ID, &subsource.TenantID, &subsource.SourceID, &subsource.SourceCode,
		&subsource.SubSourceCode, &subsource.SubSourceName, &subsource.Description,
		&subsource.Icon, &subsource.DisplayOrder, &subsource.IsActive,
		&subsource.LeadsGenerated, &subsource.ConversionCount, &subsource.ConversionRate,
		&subsource.CostPerLead, &subsource.TotalCost, &subsource.LastActivityDate,
		&metadataStr, &subsource.CreatedBy, &subsource.CreatedAt, &subsource.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lead subsource not found: %s/%s", sourceCode, subSourceCode)
		}
		return nil, fmt.Errorf("failed to get lead subsource: %w", err)
	}

	if metadataStr.Valid {
		subsource.Metadata = json.RawMessage(metadataStr.String)
	}

	return &subsource, nil
}

// GetLeadSubSources retrieves all subsources for a source
func (s *leadSourceCustomizationService) GetLeadSubSources(ctx context.Context, tenantID, sourceCode string) ([]LeadSubSourceConfig, error) {
	query := `
		SELECT id, tenant_id, source_id, source_code, sub_source_code, sub_source_name,
		       description, icon, display_order, is_active, leads_generated,
		       conversion_count, conversion_rate, cost_per_lead, total_cost,
		       last_activity_date, metadata, created_by, created_at, updated_at
		FROM tenant_lead_subsources
		WHERE tenant_id = ? AND source_code = ? AND deleted_at IS NULL
		ORDER BY display_order ASC, sub_source_name ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, sourceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead subsources: %w", err)
	}
	defer rows.Close()

	var subsources []LeadSubSourceConfig
	for rows.Next() {
		var subsource LeadSubSourceConfig
		var metadataStr sql.NullString

		err := rows.Scan(
			&subsource.ID, &subsource.TenantID, &subsource.SourceID, &subsource.SourceCode,
			&subsource.SubSourceCode, &subsource.SubSourceName, &subsource.Description,
			&subsource.Icon, &subsource.DisplayOrder, &subsource.IsActive,
			&subsource.LeadsGenerated, &subsource.ConversionCount, &subsource.ConversionRate,
			&subsource.CostPerLead, &subsource.TotalCost, &subsource.LastActivityDate,
			&metadataStr, &subsource.CreatedBy, &subsource.CreatedAt, &subsource.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead subsource: %w", err)
		}

		if metadataStr.Valid {
			subsource.Metadata = json.RawMessage(metadataStr.String)
		}

		subsources = append(subsources, subsource)
	}

	return subsources, rows.Err()
}

// UpdateLeadSubSource updates an existing subsource
func (s *leadSourceCustomizationService) UpdateLeadSubSource(ctx context.Context, tenantID string, subsource *LeadSubSourceConfig) (*LeadSubSourceConfig, error) {
	query := `
		UPDATE tenant_lead_subsources
		SET sub_source_name = ?, description = ?, icon = ?, display_order = ?,
		    is_active = ?, cost_per_lead = ?, total_cost = ?,
		    last_activity_date = ?, metadata = ?, updated_at = NOW()
		WHERE tenant_id = ? AND sub_source_code = ? AND deleted_at IS NULL
	`

	metadataJSON := "{}"
	if subsource.Metadata != nil {
		metadataJSON = string(subsource.Metadata)
	}

	result, err := s.db.ExecContext(ctx, query,
		subsource.SubSourceName, subsource.Description, subsource.Icon,
		subsource.DisplayOrder, subsource.IsActive, subsource.CostPerLead,
		subsource.TotalCost, subsource.LastActivityDate, metadataJSON,
		tenantID, subsource.SubSourceCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update lead subsource: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("lead subsource not found")
	}

	subsource.UpdatedAt = time.Now()
	return subsource, nil
}

// DeactivateLeadSubSource deactivates a subsource
func (s *leadSourceCustomizationService) DeactivateLeadSubSource(ctx context.Context, tenantID, sourceCode, subSourceCode string) error {
	query := `
		UPDATE tenant_lead_subsources
		SET is_active = FALSE, updated_at = NOW()
		WHERE tenant_id = ? AND source_code = ? AND sub_source_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, sourceCode, subSourceCode)
	if err != nil {
		return fmt.Errorf("failed to deactivate lead subsource: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("lead subsource not found")
	}

	return nil
}

// DeleteLeadSubSource soft deletes a subsource
func (s *leadSourceCustomizationService) DeleteLeadSubSource(ctx context.Context, tenantID, sourceCode, subSourceCode string) error {
	query := `
		UPDATE tenant_lead_subsources
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE tenant_id = ? AND source_code = ? AND sub_source_code = ? AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tenantID, sourceCode, subSourceCode)
	if err != nil {
		return fmt.Errorf("failed to delete lead subsource: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("lead subsource not found")
	}

	return nil
}

// ============================================================================
// CHANNEL MANAGEMENT
// ============================================================================

// CreateChannel creates a new channel
func (s *leadSourceCustomizationService) CreateChannel(ctx context.Context, tenantID string, channel *ChannelConfig) (*ChannelConfig, error) {
	if channel.ChannelCode == "" || channel.ChannelName == "" {
		return nil, errors.New("channel_code and channel_name are required")
	}

	query := `
		INSERT INTO tenant_lead_channels (
			tenant_id, channel_code, channel_name, description,
			display_order, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := s.db.ExecContext(ctx, query,
		tenantID, channel.ChannelCode, channel.ChannelName,
		channel.Description, channel.DisplayOrder, channel.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	id, _ := result.LastInsertId()
	channel.ID = id
	channel.TenantID = tenantID
	channel.CreatedAt = time.Now()
	channel.UpdatedAt = time.Now()

	return channel, nil
}

// GetChannels retrieves all channels for a tenant
func (s *leadSourceCustomizationService) GetChannels(ctx context.Context, tenantID string) ([]ChannelConfig, error) {
	query := `
		SELECT id, tenant_id, channel_code, channel_name, description,
		       display_order, is_active, created_at, updated_at
		FROM tenant_lead_channels
		WHERE tenant_id = ? AND is_active = TRUE
		ORDER BY display_order ASC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	defer rows.Close()

	var channels []ChannelConfig
	for rows.Next() {
		var channel ChannelConfig
		err := rows.Scan(
			&channel.ID, &channel.TenantID, &channel.ChannelCode, &channel.ChannelName,
			&channel.Description, &channel.DisplayOrder, &channel.IsActive,
			&channel.CreatedAt, &channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channels = append(channels, channel)
	}

	return channels, rows.Err()
}

// UpdateChannel updates an existing channel
func (s *leadSourceCustomizationService) UpdateChannel(ctx context.Context, tenantID string, channel *ChannelConfig) (*ChannelConfig, error) {
	query := `
		UPDATE tenant_lead_channels
		SET channel_name = ?, description = ?, display_order = ?,
		    is_active = ?, updated_at = NOW()
		WHERE tenant_id = ? AND channel_code = ?
	`

	result, err := s.db.ExecContext(ctx, query,
		channel.ChannelName, channel.Description, channel.DisplayOrder,
		channel.IsActive, tenantID, channel.ChannelCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update channel: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("channel not found")
	}

	channel.UpdatedAt = time.Now()
	return channel, nil
}

// ============================================================================
// ANALYTICS
// ============================================================================

// GetSourcePerformance retrieves performance metrics for a source
func (s *leadSourceCustomizationService) GetSourcePerformance(ctx context.Context, tenantID, sourceCode string, startDate, endDate time.Time) (*SourcePerformanceMetrics, error) {
	query := `
		SELECT
			ls.source_code,
			ls.source_name,
			COUNT(DISTINCT l.id) as leads_generated,
			SUM(CASE WHEN l.status IN ('contacted', 'qualified', 'converted') THEN 1 ELSE 0 END) as leads_contacted,
			SUM(CASE WHEN l.status = 'qualified' THEN 1 ELSE 0 END) as leads_qualified,
			SUM(CASE WHEN l.status = 'converted' THEN 1 ELSE 0 END) as leads_converted,
			SUM(CASE WHEN l.status = 'lost' THEN 1 ELSE 0 END) as leads_lost,
			ROUND(SUM(CASE WHEN l.status = 'converted' THEN 1 ELSE 0 END) / COUNT(DISTINCT l.id) * 100, 2) as conversion_rate,
			ROUND(SUM(CASE WHEN l.status = 'qualified' THEN 1 ELSE 0 END) / COUNT(DISTINCT l.id) * 100, 2) as qualification_rate,
			ROUND(AVG(DATEDIFF(l.converted_date, l.created_at)), 2) as average_days_to_close,
			COALESCE(SUM(l.expected_value), 0) as total_value
		FROM tenant_lead_sources ls
		LEFT JOIN lead l ON l.tenant_id = ls.tenant_id AND l.source = ls.source_code
		WHERE ls.tenant_id = ? AND ls.source_code = ?
		AND l.created_at BETWEEN ? AND ?
		GROUP BY ls.source_code, ls.source_name
	`

	var metrics SourcePerformanceMetrics
	err := s.db.QueryRowContext(ctx, query, tenantID, sourceCode, startDate, endDate).Scan(
		&metrics.SourceCode, &metrics.SourceName, &metrics.LeadsGenerated,
		&metrics.LeadsContacted, &metrics.LeadsQualified, &metrics.LeadsConverted,
		&metrics.LeadsLost, &metrics.ConversionRate, &metrics.QualificationRate,
		&metrics.AverageDaysToClose, &metrics.TotalValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data available for source: %s", sourceCode)
		}
		return nil, fmt.Errorf("failed to get source performance: %w", err)
	}

	if metrics.LeadsGenerated > 0 {
		metrics.CostPerLead = metrics.TotalValue / float64(metrics.LeadsGenerated)
		metrics.ROI = ((metrics.TotalValue - metrics.TotalValue) / metrics.TotalValue) * 100 // Simplified
	}

	metrics.UpdatedAt = time.Now()
	return &metrics, nil
}

// GetSubSourcePerformance retrieves performance metrics for a subsource
func (s *leadSourceCustomizationService) GetSubSourcePerformance(ctx context.Context, tenantID, sourceCode, subSourceCode string, startDate, endDate time.Time) (*SubSourcePerformanceMetrics, error) {
	query := `
		SELECT
			lss.sub_source_code,
			lss.sub_source_name,
			COUNT(DISTINCT l.id) as leads_generated,
			SUM(CASE WHEN l.status IN ('contacted', 'qualified', 'converted') THEN 1 ELSE 0 END) as leads_contacted,
			SUM(CASE WHEN l.status = 'qualified' THEN 1 ELSE 0 END) as leads_qualified,
			SUM(CASE WHEN l.status = 'converted' THEN 1 ELSE 0 END) as leads_converted,
			SUM(CASE WHEN l.status = 'lost' THEN 1 ELSE 0 END) as leads_lost,
			ROUND(SUM(CASE WHEN l.status = 'converted' THEN 1 ELSE 0 END) / COUNT(DISTINCT l.id) * 100, 2) as conversion_rate,
			COALESCE(lss.total_cost, 0) as total_cost,
			COALESCE(lss.cost_per_lead, 0) as cost_per_lead,
			COALESCE(SUM(l.expected_value), 0) as revenue
		FROM tenant_lead_subsources lss
		LEFT JOIN lead l ON l.tenant_id = lss.tenant_id AND l.subsource = lss.sub_source_code
		WHERE lss.tenant_id = ? AND lss.source_code = ? AND lss.sub_source_code = ?
		AND l.created_at BETWEEN ? AND ?
		GROUP BY lss.sub_source_code, lss.sub_source_name
	`

	var metrics SubSourcePerformanceMetrics
	err := s.db.QueryRowContext(ctx, query, tenantID, sourceCode, subSourceCode, startDate, endDate).Scan(
		&metrics.SubSourceCode, &metrics.SubSourceName, &metrics.LeadsGenerated,
		&metrics.LeadsContacted, &metrics.LeadsQualified, &metrics.LeadsConverted,
		&metrics.LeadsLost, &metrics.ConversionRate, &metrics.TotalCost,
		&metrics.CostPerLead, &metrics.Revenue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data available for subsource: %s/%s", sourceCode, subSourceCode)
		}
		return nil, fmt.Errorf("failed to get subsource performance: %w", err)
	}

	if metrics.LeadsConverted > 0 {
		metrics.CostPerConversion = metrics.TotalCost / float64(metrics.LeadsConverted)
	}

	if metrics.Revenue > 0 {
		metrics.ROI = ((metrics.Revenue - metrics.TotalCost) / metrics.TotalCost) * 100
	}

	return &metrics, nil
}

// GetLeadSourceTrends retrieves trend data for lead sources
func (s *leadSourceCustomizationService) GetLeadSourceTrends(ctx context.Context, tenantID string, days int) (map[string]interface{}, error) {
	query := `
		SELECT
			DATE(l.created_at) as date,
			ls.source_code,
			ls.source_name,
			COUNT(l.id) as leads_generated,
			SUM(CASE WHEN l.status = 'converted' THEN 1 ELSE 0 END) as conversions,
			SUM(l.expected_value) as revenue
		FROM lead l
		LEFT JOIN tenant_lead_sources ls ON l.source = ls.source_code
		WHERE l.tenant_id = ? AND l.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY DATE(l.created_at), ls.source_code, ls.source_name
		ORDER BY date DESC
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get lead trends: %w", err)
	}
	defer rows.Close()

	trends := make(map[string]interface{})
	dailyData := make([]map[string]interface{}, 0)

	for rows.Next() {
		var date time.Time
		var sourceCode, sourceName string
		var leadsGenerated, conversions int64
		var revenue sql.NullFloat64

		err := rows.Scan(&date, &sourceCode, &sourceName, &leadsGenerated, &conversions, &revenue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trend: %w", err)
		}

		dailyEntry := map[string]interface{}{
			"date":            date,
			"source_code":     sourceCode,
			"source_name":     sourceName,
			"leads_generated": leadsGenerated,
			"conversions":     conversions,
		}

		if revenue.Valid {
			dailyEntry["revenue"] = revenue.Float64
		}

		dailyData = append(dailyData, dailyEntry)
	}

	trends["daily_data"] = dailyData
	trends["period_days"] = days
	return trends, rows.Err()
}
