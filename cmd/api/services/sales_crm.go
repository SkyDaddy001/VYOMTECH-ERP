package services

import (
	"database/sql"
	"fmt"
	"time"

	"lms/cmd/api/models"
)

// ===== SALES & PRESALES REP SERVICES =====

type LeadService struct {
	db *sql.DB
}

func NewLeadService(db *sql.DB) *LeadService {
	return &LeadService{db: db}
}

func (s *LeadService) CreateLead(tenantID string, lead *models.Lead) (string, error) {
	lead.ID = generateID("lead")
	lead.TenantID = tenantID
	lead.CreatedAt = time.Now()
	lead.UpdatedAt = time.Now()
	lead.Status = models.LeadStatusNew

	query := `
		INSERT INTO leads (id, tenant_id, assigned_to_id, assigned_to_name, first_name, last_name, 
			company_name, email, phone, mobile_phone, website, industry, company_size, location, 
			city, state, country, postal_code, source, status, budget, currency, description, 
			rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, 
			$18, $19, $20, $21, $22, $23, $24, $25, $26)
	`

	_, err := s.db.Exec(query,
		lead.ID, lead.TenantID, lead.AssignedToID, lead.AssignedToName, lead.FirstName,
		lead.LastName, lead.CompanyName, lead.Email, lead.Phone, lead.MobilePhone,
		lead.Website, lead.Industry, lead.CompanySize, lead.Location, lead.City,
		lead.State, lead.Country, lead.PostalCode, lead.Source, lead.Status,
		lead.Budget, lead.Currency, lead.Description, lead.Rating, lead.CreatedAt, lead.UpdatedAt,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create lead: %w", err)
	}

	return lead.ID, nil
}

func (s *LeadService) GetLead(tenantID, leadID string) (*models.Lead, error) {
	lead := &models.Lead{}
	query := `SELECT id, tenant_id, assigned_to_id, assigned_to_name, first_name, last_name, 
		company_name, email, phone, mobile_phone, website, industry, company_size, location, 
		city, state, country, postal_code, source, status, budget, currency, description, 
		rating, created_at, updated_at, converted_at, last_contacted_at 
		FROM leads WHERE id = $1 AND tenant_id = $2`

	err := s.db.QueryRow(query, leadID, tenantID).Scan(
		&lead.ID, &lead.TenantID, &lead.AssignedToID, &lead.AssignedToName, &lead.FirstName,
		&lead.LastName, &lead.CompanyName, &lead.Email, &lead.Phone, &lead.MobilePhone,
		&lead.Website, &lead.Industry, &lead.CompanySize, &lead.Location, &lead.City,
		&lead.State, &lead.Country, &lead.PostalCode, &lead.Source, &lead.Status,
		&lead.Budget, &lead.Currency, &lead.Description, &lead.Rating, &lead.CreatedAt,
		&lead.UpdatedAt, &lead.ConvertedAt, &lead.LastContactedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lead not found")
		}
		return nil, fmt.Errorf("failed to fetch lead: %w", err)
	}

	return lead, nil
}

func (s *LeadService) ListLeads(tenantID string, filters map[string]interface{}, limit, offset int) ([]*models.Lead, int64, error) {
	query := `SELECT id, tenant_id, assigned_to_id, assigned_to_name, first_name, last_name, 
		company_name, email, phone, mobile_phone, website, industry, company_size, location, 
		city, state, country, postal_code, source, status, budget, currency, description, 
		rating, created_at, updated_at, converted_at, last_contacted_at 
		FROM leads WHERE tenant_id = $1`

	args := []interface{}{tenantID}
	argIndex := 2

	if status, ok := filters["status"].(string); ok && status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if source, ok := filters["source"].(string); ok && source != "" {
		query += fmt.Sprintf(" AND source = $%d", argIndex)
		args = append(args, source)
		argIndex++
	}

	if assignedToID, ok := filters["assigned_to_id"].(string); ok && assignedToID != "" {
		query += fmt.Sprintf(" AND assigned_to_id = $%d", argIndex)
		args = append(args, assignedToID)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as t", query)
	var total int64
	s.db.QueryRow(countQuery, args...).Scan(&total)

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list leads: %w", err)
	}
	defer rows.Close()

	var leads []*models.Lead
	for rows.Next() {
		lead := &models.Lead{}
		err := rows.Scan(
			&lead.ID, &lead.TenantID, &lead.AssignedToID, &lead.AssignedToName, &lead.FirstName,
			&lead.LastName, &lead.CompanyName, &lead.Email, &lead.Phone, &lead.MobilePhone,
			&lead.Website, &lead.Industry, &lead.CompanySize, &lead.Location, &lead.City,
			&lead.State, &lead.Country, &lead.PostalCode, &lead.Source, &lead.Status,
			&lead.Budget, &lead.Currency, &lead.Description, &lead.Rating, &lead.CreatedAt,
			&lead.UpdatedAt, &lead.ConvertedAt, &lead.LastContactedAt,
		)
		if err != nil {
			continue
		}
		leads = append(leads, lead)
	}

	return leads, total, nil
}

func (s *LeadService) UpdateLeadStatus(tenantID, leadID string, status models.LeadStatus) error {
	query := `UPDATE leads SET status = $1, updated_at = $2 WHERE id = $3 AND tenant_id = $4`

	if status == models.LeadStatusConverted {
		query = `UPDATE leads SET status = $1, updated_at = $2, converted_at = $3 WHERE id = $4 AND tenant_id = $5`
		_, err := s.db.Exec(query, status, time.Now(), time.Now(), leadID, tenantID)
		return err
	}

	_, err := s.db.Exec(query, status, time.Now(), leadID, tenantID)
	return err
}

// ===== OPPORTUNITY SERVICE =====

type OpportunityService struct {
	db *sql.DB
}

func NewOpportunityService(db *sql.DB) *OpportunityService {
	return &OpportunityService{db: db}
}

func (s *OpportunityService) CreateOpportunity(tenantID string, opp *models.Opportunity) (string, error) {
	opp.ID = generateID("oppt")
	opp.TenantID = tenantID
	opp.CreatedAt = time.Now()
	opp.UpdatedAt = time.Now()
	if opp.Stage == "" {
		opp.Stage = models.OpportunityStageLead
	}
	if opp.Probability == 0 {
		opp.Probability = 30
	}
	opp.ExpectedRevenue = opp.Amount * float64(opp.Probability) / 100

	query := `
		INSERT INTO opportunities (id, tenant_id, lead_id, account_id, assigned_to_id, assigned_to_name, 
			name, description, stage, amount, currency, close_date, expected_revenue, probability, 
			source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := s.db.Exec(query,
		opp.ID, opp.TenantID, opp.LeadID, opp.AccountID, opp.AssignedToID, opp.AssignedToName,
		opp.Name, opp.Description, opp.Stage, opp.Amount, opp.Currency, opp.CloseDate,
		opp.ExpectedRevenue, opp.Probability, opp.Source, opp.CreatedAt, opp.UpdatedAt,
	)

	return opp.ID, err
}

func (s *OpportunityService) GetOpportunity(tenantID, oppID string) (*models.Opportunity, error) {
	opp := &models.Opportunity{}
	query := `SELECT id, tenant_id, lead_id, account_id, assigned_to_id, assigned_to_name, 
		name, description, stage, amount, currency, close_date, expected_revenue, probability, 
		source, competitor_info, next_action, next_action_date, created_at, updated_at, won_at, lost_at, lost_reason
		FROM opportunities WHERE id = $1 AND tenant_id = $2`

	err := s.db.QueryRow(query, oppID, tenantID).Scan(
		&opp.ID, &opp.TenantID, &opp.LeadID, &opp.AccountID, &opp.AssignedToID, &opp.AssignedToName,
		&opp.Name, &opp.Description, &opp.Stage, &opp.Amount, &opp.Currency, &opp.CloseDate,
		&opp.ExpectedRevenue, &opp.Probability, &opp.Source, &opp.CompetitorInfo, &opp.NextAction,
		&opp.NextActionDate, &opp.CreatedAt, &opp.UpdatedAt, &opp.WonAt, &opp.LostAt, &opp.LostReason,
	)

	return opp, err
}

func (s *OpportunityService) UpdateOpportunityStage(tenantID, oppID string, stage models.OpportunityStage) error {
	now := time.Now()
	query := `UPDATE opportunities SET stage = $1, updated_at = $2`
	args := []interface{}{stage, now}

	if stage == models.OpportunityStageClosedWon {
		query += `, won_at = $3`
		args = append(args, now)
	} else if stage == models.OpportunityStageClosedLost {
		query += `, lost_at = $3`
		args = append(args, now)
	}

	query += ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1) + ` AND tenant_id = $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, oppID, tenantID)

	_, err := s.db.Exec(query, args...)
	return err
}

func (s *OpportunityService) GetPipelineValue(tenantID string) (map[models.OpportunityStage]float64, error) {
	pipeline := make(map[models.OpportunityStage]float64)
	query := `SELECT stage, SUM(expected_revenue) FROM opportunities 
		WHERE tenant_id = $1 AND stage NOT IN ($2, $3) 
		GROUP BY stage`

	rows, err := s.db.Query(query, tenantID, models.OpportunityStageClosedLost, models.OpportunityStageClosedWon)
	if err != nil {
		return pipeline, err
	}
	defer rows.Close()

	for rows.Next() {
		var stage models.OpportunityStage
		var value float64
		rows.Scan(&stage, &value)
		pipeline[stage] = value
	}

	return pipeline, nil
}

// ===== ACTIVITY SERVICE =====

type ActivityService struct {
	db *sql.DB
}

func NewActivityService(db *sql.DB) *ActivityService {
	return &ActivityService{db: db}
}

func (s *ActivityService) CreateActivity(tenantID string, activity *models.Activity) (string, error) {
	activity.ID = generateID("activity")
	activity.TenantID = tenantID
	activity.CreatedAt = time.Now()
	activity.UpdatedAt = time.Now()
	if activity.Status == "" {
		activity.Status = "pending"
	}

	query := `
		INSERT INTO activities (id, tenant_id, lead_id, opportunity_id, account_id, activity_type, 
			subject, description, created_by_id, created_by_name, assigned_to_id, assigned_to_name, 
			activity_date, due_date, status, priority, duration, outcome, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := s.db.Exec(query,
		activity.ID, activity.TenantID, activity.LeadID, activity.OpportunityID, activity.AccountID,
		activity.ActivityType, activity.Subject, activity.Description, activity.CreatedByID,
		activity.CreatedByName, activity.AssignedToID, activity.AssignedToName, activity.ActivityDate,
		activity.DueDate, activity.Status, activity.Priority, activity.Duration, activity.Outcome,
		activity.CreatedAt, activity.UpdatedAt,
	)

	return activity.ID, err
}

func (s *ActivityService) GetActivitiesByLead(tenantID, leadID string, limit, offset int) ([]*models.Activity, int64, error) {
	var total int64
	s.db.QueryRow(`SELECT COUNT(*) FROM activities WHERE tenant_id = $1 AND lead_id = $2`, tenantID, leadID).Scan(&total)

	query := `SELECT id, tenant_id, lead_id, opportunity_id, account_id, activity_type, 
		subject, description, created_by_id, created_by_name, assigned_to_id, assigned_to_name, 
		activity_date, due_date, status, priority, duration, outcome, created_at, updated_at
		FROM activities WHERE tenant_id = $1 AND lead_id = $2 ORDER BY activity_date DESC LIMIT $3 OFFSET $4`

	rows, err := s.db.Query(query, tenantID, leadID, limit, offset)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	var activities []*models.Activity
	for rows.Next() {
		activity := &models.Activity{}
		err := rows.Scan(
			&activity.ID, &activity.TenantID, &activity.LeadID, &activity.OpportunityID, &activity.AccountID,
			&activity.ActivityType, &activity.Subject, &activity.Description, &activity.CreatedByID,
			&activity.CreatedByName, &activity.AssignedToID, &activity.AssignedToName, &activity.ActivityDate,
			&activity.DueDate, &activity.Status, &activity.Priority, &activity.Duration, &activity.Outcome,
			&activity.CreatedAt, &activity.UpdatedAt,
		)
		if err == nil {
			activities = append(activities, activity)
		}
	}

	return activities, total, nil
}

// ===== CRM ACCOUNT SERVICE =====

type AccountService struct {
	db *sql.DB
}

func NewAccountService(db *sql.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(tenantID string, account *models.Account) (string, error) {
	account.ID = generateID("account")
	account.TenantID = tenantID
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	query := `
		INSERT INTO accounts (id, tenant_id, name, website, industry, company_size, billing_address, 
			shipping_address, phone, email, account_manager_id, annual_revenue, employees, rating, type, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := s.db.Exec(query,
		account.ID, account.TenantID, account.Name, account.Website, account.Industry, account.CompanySize,
		account.BillingAddress, account.ShippingAddress, account.Phone, account.Email, account.AccountManager,
		account.AnnualRevenue, account.Employees, account.Rating, account.Type, account.Description,
		account.CreatedAt, account.UpdatedAt,
	)

	return account.ID, err
}

func (s *AccountService) GetAccount(tenantID, accountID string) (*models.Account, error) {
	account := &models.Account{}
	query := `SELECT id, tenant_id, name, website, industry, company_size, billing_address, 
		shipping_address, phone, email, account_manager_id, annual_revenue, employees, rating, type, description, created_at, updated_at
		FROM accounts WHERE id = $1 AND tenant_id = $2`

	err := s.db.QueryRow(query, accountID, tenantID).Scan(
		&account.ID, &account.TenantID, &account.Name, &account.Website, &account.Industry, &account.CompanySize,
		&account.BillingAddress, &account.ShippingAddress, &account.Phone, &account.Email, &account.AccountManager,
		&account.AnnualRevenue, &account.Employees, &account.Rating, &account.Type, &account.Description,
		&account.CreatedAt, &account.UpdatedAt,
	)

	return account, err
}

// ===== CONTACT SERVICE =====

type ContactService struct {
	db *sql.DB
}

func NewContactService(db *sql.DB) *ContactService {
	return &ContactService{db: db}
}

func (s *ContactService) CreateContact(tenantID string, contact *models.Contact) (string, error) {
	contact.ID = generateID("contact")
	contact.TenantID = tenantID
	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()

	query := `
		INSERT INTO contacts (id, tenant_id, account_id, first_name, last_name, title, department, 
			email, phone, mobile_phone, role, is_primary, linkedin, twitter, preferred_language, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := s.db.Exec(query,
		contact.ID, contact.TenantID, contact.AccountID, contact.FirstName, contact.LastName,
		contact.Title, contact.Department, contact.Email, contact.Phone, contact.MobilePhone,
		contact.Role, contact.IsPrimary, contact.LinkedIn, contact.Twitter, contact.PreferredLanguage,
		contact.CreatedAt, contact.UpdatedAt,
	)

	return contact.ID, err
}

func (s *ContactService) GetContactsByAccount(tenantID, accountID string) ([]*models.Contact, error) {
	query := `SELECT id, tenant_id, account_id, first_name, last_name, title, department, 
		email, phone, mobile_phone, role, is_primary, linkedin, twitter, preferred_language, created_at, updated_at
		FROM contacts WHERE tenant_id = $1 AND account_id = $2 ORDER BY is_primary DESC, created_at DESC`

	rows, err := s.db.Query(query, tenantID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		contact := &models.Contact{}
		err := rows.Scan(
			&contact.ID, &contact.TenantID, &contact.AccountID, &contact.FirstName, &contact.LastName,
			&contact.Title, &contact.Department, &contact.Email, &contact.Phone, &contact.MobilePhone,
			&contact.Role, &contact.IsPrimary, &contact.LinkedIn, &contact.Twitter, &contact.PreferredLanguage,
			&contact.CreatedAt, &contact.UpdatedAt,
		)
		if err == nil {
			contacts = append(contacts, contact)
		}
	}

	return contacts, nil
}

// ===== INTERACTION SERVICE =====

type InteractionService struct {
	db *sql.DB
}

func NewInteractionService(db *sql.DB) *InteractionService {
	return &InteractionService{db: db}
}

func (s *InteractionService) CreateInteraction(tenantID string, interaction *models.Interaction) (string, error) {
	interaction.ID = generateID("interaction")
	interaction.TenantID = tenantID
	interaction.CreatedAt = time.Now()
	interaction.UpdatedAt = time.Now()

	query := `
		INSERT INTO interactions (id, tenant_id, account_id, contact_id, interaction_type, channel, 
			subject, message, status, priority, assigned_to_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := s.db.Exec(query,
		interaction.ID, interaction.TenantID, interaction.AccountID, interaction.ContactID,
		interaction.InteractionType, interaction.Channel, interaction.Subject, interaction.Message,
		interaction.Status, interaction.Priority, interaction.AssignedToID, interaction.CreatedAt,
		interaction.UpdatedAt,
	)

	return interaction.ID, err
}

func (s *InteractionService) GetAccountInteractions(tenantID, accountID string, limit, offset int) ([]*models.Interaction, int64, error) {
	var total int64
	s.db.QueryRow(`SELECT COUNT(*) FROM interactions WHERE tenant_id = $1 AND account_id = $2`, tenantID, accountID).Scan(&total)

	query := `SELECT id, tenant_id, account_id, contact_id, interaction_type, channel, 
		subject, message, status, priority, assigned_to_id, created_at, updated_at, resolved_at
		FROM interactions WHERE tenant_id = $1 AND account_id = $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`

	rows, err := s.db.Query(query, tenantID, accountID, limit, offset)
	if err != nil {
		return nil, total, err
	}
	defer rows.Close()

	var interactions []*models.Interaction
	for rows.Next() {
		interaction := &models.Interaction{}
		err := rows.Scan(
			&interaction.ID, &interaction.TenantID, &interaction.AccountID, &interaction.ContactID,
			&interaction.InteractionType, &interaction.Channel, &interaction.Subject, &interaction.Message,
			&interaction.Status, &interaction.Priority, &interaction.AssignedToID, &interaction.CreatedAt,
			&interaction.UpdatedAt, &interaction.ResolvedAt,
		)
		if err == nil {
			interactions = append(interactions, interaction)
		}
	}

	return interactions, total, nil
}

// ===== SALES ANALYTICS SERVICE =====

type SalesAnalyticsService struct {
	db *sql.DB
}

func NewSalesAnalyticsService(db *sql.DB) *SalesAnalyticsService {
	return &SalesAnalyticsService{db: db}
}

func (s *SalesAnalyticsService) GetSalesDashboard(tenantID string) (*models.SalesDashboard, error) {
	dashboard := &models.SalesDashboard{
		LeadsBySource:        make(map[string]int),
		OpportunitiesByStage: make(map[string]int),
		TopSalesReps:         make([]models.RepMetrics, 0),
	}

	// Total leads and qualified leads
	s.db.QueryRow(`SELECT COUNT(*) FROM leads WHERE tenant_id = $1 AND status != $2`,
		tenantID, models.LeadStatusLost).Scan(&dashboard.TotalLeads)
	s.db.QueryRow(`SELECT COUNT(*) FROM leads WHERE tenant_id = $1 AND status = $2`,
		tenantID, models.LeadStatusConverted).Scan(&dashboard.ConversionRate)

	if dashboard.TotalLeads > 0 {
		dashboard.ConversionRate = (float64(dashboard.ConversionRate) / float64(dashboard.TotalLeads)) * 100
	}

	// Opportunities
	s.db.QueryRow(`SELECT COUNT(*) FROM opportunities WHERE tenant_id = $1 AND stage NOT IN ($2, $3)`,
		tenantID, models.OpportunityStageClosedWon, models.OpportunityStageClosedLost).
		Scan(&dashboard.TotalOpportunities)

	s.db.QueryRow(`SELECT SUM(expected_revenue) FROM opportunities WHERE tenant_id = $1 AND stage NOT IN ($2, $3)`,
		tenantID, models.OpportunityStageClosedWon, models.OpportunityStageClosedLost).
		Scan(&dashboard.PipelineValue)

	// Won deals
	s.db.QueryRow(`SELECT COUNT(*), SUM(amount) FROM opportunities WHERE tenant_id = $1 AND stage = $2 AND EXTRACT(YEAR FROM won_at) = EXTRACT(YEAR FROM NOW())`,
		tenantID, models.OpportunityStageClosedWon).Scan(&dashboard.WonDeals, &dashboard.WonValue)

	// Lost deals
	s.db.QueryRow(`SELECT COUNT(*) FROM opportunities WHERE tenant_id = $1 AND stage = $2`,
		tenantID, models.OpportunityStageClosedLost).Scan(&dashboard.LostDeals)

	// Average deal size
	if dashboard.WonDeals > 0 && dashboard.WonValue > 0 {
		dashboard.AvgDealSize = dashboard.WonValue / float64(dashboard.WonDeals)
	}

	return dashboard, nil
}

// Helper function
func generateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
