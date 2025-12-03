package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"vyomtech-backend/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// ============================================================================
// MILESTONE HANDLERS
// ============================================================================

// CreateLeadMilestone - POST /api/v1/sales/milestones/lead
func (sh *SalesHandler) CreateLeadMilestone(w http.ResponseWriter, r *http.Request) {
	var milestone models.SalesLeadMilestone

	if err := json.NewDecoder(r.Body).Decode(&milestone); err != nil {
		sh.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	milestone.ID = uuid.New().String()
	milestone.CreatedAt = time.Now()
	milestone.UpdatedAt = time.Now()

	query := `
		INSERT INTO sales_lead_milestones 
		(id, tenant_id, lead_id, milestone_type, milestone_date, milestone_time, notes, 
		 location_latitude, location_longitude, location_name, status_before, status_after, 
		 visited_by, duration_minutes, outcome, follow_up_date, follow_up_required, metadata, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, created_at
	`

	metadataJSON, _ := json.Marshal(milestone.Metadata)

	err := sh.DB.QueryRow(query,
		milestone.ID, milestone.TenantID, milestone.LeadID, milestone.MilestoneType,
		milestone.MilestoneDate, milestone.MilestoneTime, milestone.Notes,
		milestone.LocationLatitude, milestone.LocationLongitude, milestone.LocationName,
		milestone.StatusBefore, milestone.StatusAfter, milestone.VisitedBy,
		milestone.DurationMinutes, milestone.Outcome, milestone.FollowUpDate,
		milestone.FollowUpRequired, metadataJSON, milestone.CreatedBy,
		milestone.CreatedAt, milestone.UpdatedAt,
	).Scan(&milestone.ID, &milestone.CreatedAt)

	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to create milestone")
		return
	}

	sh.respondJSON(w, http.StatusCreated, milestone)
}

// GetLeadMilestones - GET /api/v1/sales/milestones/lead/{lead_id}
func (sh *SalesHandler) GetLeadMilestones(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["lead_id"]

	query := `
		SELECT id, tenant_id, lead_id, milestone_type, milestone_date, milestone_time, notes,
		       location_latitude, location_longitude, location_name, status_before, status_after,
		       visited_by, duration_minutes, outcome, follow_up_date, follow_up_required, metadata,
		       created_by, created_at, updated_at
		FROM sales_lead_milestones
		WHERE lead_id = $1 AND tenant_id = $2
		ORDER BY milestone_date DESC
	`

	rows, err := sh.DB.Query(query, leadID, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch milestones")
		return
	}
	defer rows.Close()

	milestones := []models.SalesLeadMilestone{}

	for rows.Next() {
		var milestone models.SalesLeadMilestone
		var metadataJSON []byte

		if err := rows.Scan(
			&milestone.ID, &milestone.TenantID, &milestone.LeadID, &milestone.MilestoneType,
			&milestone.MilestoneDate, &milestone.MilestoneTime, &milestone.Notes,
			&milestone.LocationLatitude, &milestone.LocationLongitude, &milestone.LocationName,
			&milestone.StatusBefore, &milestone.StatusAfter, &milestone.VisitedBy,
			&milestone.DurationMinutes, &milestone.Outcome, &milestone.FollowUpDate,
			&milestone.FollowUpRequired, &metadataJSON, &milestone.CreatedBy,
			&milestone.CreatedAt, &milestone.UpdatedAt,
		); err != nil {
			continue
		}

		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &milestone.Metadata)
		}

		milestones = append(milestones, milestone)
	}

	sh.respondJSON(w, http.StatusOK, milestones)
}

// ============================================================================
// ENGAGEMENT HANDLERS
// ============================================================================

// CreateLeadEngagement - POST /api/v1/sales/engagement
func (sh *SalesHandler) CreateLeadEngagement(w http.ResponseWriter, r *http.Request) {
	var engagement models.SalesLeadEngagement

	if err := json.NewDecoder(r.Body).Decode(&engagement); err != nil {
		sh.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	engagement.ID = uuid.New().String()
	if engagement.EngagementDate.IsZero() {
		engagement.EngagementDate = time.Now()
	}
	engagement.CreatedAt = time.Now()

	query := `
		INSERT INTO sales_lead_engagement 
		(id, tenant_id, lead_id, engagement_type, engagement_date, engagement_channel, 
		 subject, notes, status, response_received, response_date, response_notes, assigned_to, 
		 duration_seconds, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at
	`

	err := sh.DB.QueryRow(query,
		engagement.ID, engagement.TenantID, engagement.LeadID, engagement.EngagementType,
		engagement.EngagementDate, engagement.EngagementChannel, engagement.Subject,
		engagement.Notes, engagement.Status, engagement.ResponseReceived, engagement.ResponseDate,
		engagement.ResponseNotes, engagement.AssignedTo, engagement.DurationSeconds,
		engagement.CreatedBy, engagement.CreatedAt,
	).Scan(&engagement.ID, &engagement.CreatedAt)

	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to create engagement")
		return
	}

	sh.respondJSON(w, http.StatusCreated, engagement)
}

// GetLeadEngagements - GET /api/v1/sales/engagement/{lead_id}
func (sh *SalesHandler) GetLeadEngagements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["lead_id"]

	query := `
		SELECT id, tenant_id, lead_id, engagement_type, engagement_date, engagement_channel,
		       subject, notes, status, response_received, response_date, response_notes,
		       assigned_to, duration_seconds, created_by, created_at
		FROM sales_lead_engagement
		WHERE lead_id = $1 AND tenant_id = $2
		ORDER BY engagement_date DESC
		LIMIT 100
	`

	rows, err := sh.DB.Query(query, leadID, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch engagements")
		return
	}
	defer rows.Close()

	engagements := []models.SalesLeadEngagement{}

	for rows.Next() {
		var engagement models.SalesLeadEngagement

		if err := rows.Scan(
			&engagement.ID, &engagement.TenantID, &engagement.LeadID, &engagement.EngagementType,
			&engagement.EngagementDate, &engagement.EngagementChannel, &engagement.Subject,
			&engagement.Notes, &engagement.Status, &engagement.ResponseReceived,
			&engagement.ResponseDate, &engagement.ResponseNotes, &engagement.AssignedTo,
			&engagement.DurationSeconds, &engagement.CreatedBy, &engagement.CreatedAt,
		); err != nil {
			continue
		}

		engagements = append(engagements, engagement)
	}

	sh.respondJSON(w, http.StatusOK, engagements)
}

// ============================================================================
// BOOKING HANDLERS
// ============================================================================

// CreateBooking - POST /api/v1/sales/bookings
func (sh *SalesHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.SalesBooking

	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		sh.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	booking.ID = uuid.New().String()
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	query := `
		INSERT INTO sales_bookings 
		(id, tenant_id, customer_id, lead_id, booking_code, booking_date, booking_amount,
		 unit_type, unit_count, units_booked, units_available, delivery_date, status,
		 cancellation_date, cancellation_reason, cancellation_refund_amount, notes, created_by,
		 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id, created_at
	`

	err := sh.DB.QueryRow(query,
		booking.ID, booking.TenantID, booking.CustomerID, booking.LeadID, booking.BookingCode,
		booking.BookingDate, booking.BookingAmount, booking.UnitType, booking.UnitCount,
		booking.UnitsBooked, booking.UnitsAvailable, booking.DeliveryDate, booking.Status,
		booking.CancellationDate, booking.CancellationReason, booking.CancellationRefundAmount,
		booking.Notes, booking.CreatedBy, booking.CreatedAt, booking.UpdatedAt,
	).Scan(&booking.ID, &booking.CreatedAt)

	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to create booking")
		return
	}

	sh.respondJSON(w, http.StatusCreated, booking)
}

// GetBookings - GET /api/v1/sales/bookings
func (sh *SalesHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, tenant_id, customer_id, lead_id, booking_code, booking_date, booking_amount,
		       unit_type, unit_count, units_booked, units_available, delivery_date, status,
		       cancellation_date, cancellation_reason, cancellation_refund_amount, notes,
		       created_by, created_at, updated_at
		FROM sales_bookings
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY booking_date DESC
		LIMIT 100
	`

	rows, err := sh.DB.Query(query, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch bookings")
		return
	}
	defer rows.Close()

	bookings := []models.SalesBooking{}

	for rows.Next() {
		var booking models.SalesBooking

		if err := rows.Scan(
			&booking.ID, &booking.TenantID, &booking.CustomerID, &booking.LeadID,
			&booking.BookingCode, &booking.BookingDate, &booking.BookingAmount,
			&booking.UnitType, &booking.UnitCount, &booking.UnitsBooked,
			&booking.UnitsAvailable, &booking.DeliveryDate, &booking.Status,
			&booking.CancellationDate, &booking.CancellationReason,
			&booking.CancellationRefundAmount, &booking.Notes, &booking.CreatedBy,
			&booking.CreatedAt, &booking.UpdatedAt,
		); err != nil {
			continue
		}

		bookings = append(bookings, booking)
	}

	sh.respondJSON(w, http.StatusOK, bookings)
}

// ============================================================================
// ACCOUNT LEDGER HANDLERS
// ============================================================================

// CreateLedgerEntry - POST /api/v1/sales/ledger
func (sh *SalesHandler) CreateLedgerEntry(w http.ResponseWriter, r *http.Request) {
	var ledger models.SalesAccountLedger

	if err := json.NewDecoder(r.Body).Decode(&ledger); err != nil {
		sh.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ledger.ID = uuid.New().String()
	ledger.CreatedAt = time.Now()
	ledger.UpdatedAt = time.Now()

	query := `
		INSERT INTO sales_account_ledgers 
		(id, tenant_id, customer_id, ledger_code, ledger_date, transaction_type,
		 reference_document_type, reference_document_id, debit_amount, credit_amount,
		 balance_after, description, remarks, status, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING id, created_at
	`

	err := sh.DB.QueryRow(query,
		ledger.ID, ledger.TenantID, ledger.CustomerID, ledger.LedgerCode, ledger.LedgerDate,
		ledger.TransactionType, ledger.ReferenceDocumentType, ledger.ReferenceDocumentID,
		ledger.DebitAmount, ledger.CreditAmount, ledger.BalanceAfter, ledger.Description,
		ledger.Remarks, ledger.Status, ledger.CreatedBy, ledger.CreatedAt, ledger.UpdatedAt,
	).Scan(&ledger.ID, &ledger.CreatedAt)

	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to create ledger entry")
		return
	}

	sh.respondJSON(w, http.StatusCreated, ledger)
}

// GetCustomerLedger - GET /api/v1/sales/ledger/{customer_id}
func (sh *SalesHandler) GetCustomerLedger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customer_id"]

	query := `
		SELECT id, tenant_id, customer_id, ledger_code, ledger_date, transaction_type,
		       reference_document_type, reference_document_id, debit_amount, credit_amount,
		       balance_after, description, remarks, status, created_by, created_at, updated_at
		FROM sales_account_ledgers
		WHERE customer_id = $1 AND tenant_id = $2 AND status = 'active'
		ORDER BY ledger_date DESC
		LIMIT 500
	`

	rows, err := sh.DB.Query(query, customerID, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch ledger")
		return
	}
	defer rows.Close()

	ledgers := []models.SalesAccountLedger{}

	for rows.Next() {
		var ledger models.SalesAccountLedger

		if err := rows.Scan(
			&ledger.ID, &ledger.TenantID, &ledger.CustomerID, &ledger.LedgerCode,
			&ledger.LedgerDate, &ledger.TransactionType, &ledger.ReferenceDocumentType,
			&ledger.ReferenceDocumentID, &ledger.DebitAmount, &ledger.CreditAmount,
			&ledger.BalanceAfter, &ledger.Description, &ledger.Remarks, &ledger.Status,
			&ledger.CreatedBy, &ledger.CreatedAt, &ledger.UpdatedAt,
		); err != nil {
			continue
		}

		ledgers = append(ledgers, ledger)
	}

	sh.respondJSON(w, http.StatusOK, ledgers)
}

// ============================================================================
// CAMPAIGN HANDLERS
// ============================================================================

// CreateCampaign - POST /api/v1/sales/campaigns
func (sh *SalesHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var campaign models.SalesCampaign

	if err := json.NewDecoder(r.Body).Decode(&campaign); err != nil {
		sh.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	campaign.ID = uuid.New().String()
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()

	query := `
		INSERT INTO sales_campaigns 
		(id, tenant_id, campaign_code, campaign_name, campaign_type, description,
		 start_date, end_date, budget, expected_roi, assigned_to, status, created_by,
		 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at
	`

	err := sh.DB.QueryRow(query,
		campaign.ID, campaign.TenantID, campaign.CampaignCode, campaign.CampaignName,
		campaign.CampaignType, campaign.Description, campaign.StartDate, campaign.EndDate,
		campaign.Budget, campaign.ExpectedROI, campaign.AssignedTo, campaign.Status,
		campaign.CreatedBy, campaign.CreatedAt, campaign.UpdatedAt,
	).Scan(&campaign.ID, &campaign.CreatedAt)

	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to create campaign")
		return
	}

	sh.respondJSON(w, http.StatusCreated, campaign)
}

// GetCampaigns - GET /api/v1/sales/campaigns
func (sh *SalesHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, tenant_id, campaign_code, campaign_name, campaign_type, description,
		       start_date, end_date, budget, expected_roi, assigned_to, status, created_by,
		       created_at, updated_at
		FROM sales_campaigns
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY start_date DESC
	`

	rows, err := sh.DB.Query(query, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch campaigns")
		return
	}
	defer rows.Close()

	campaigns := []models.SalesCampaign{}

	for rows.Next() {
		var campaign models.SalesCampaign

		if err := rows.Scan(
			&campaign.ID, &campaign.TenantID, &campaign.CampaignCode, &campaign.CampaignName,
			&campaign.CampaignType, &campaign.Description, &campaign.StartDate,
			&campaign.EndDate, &campaign.Budget, &campaign.ExpectedROI, &campaign.AssignedTo,
			&campaign.Status, &campaign.CreatedBy, &campaign.CreatedAt, &campaign.UpdatedAt,
		); err != nil {
			continue
		}

		campaigns = append(campaigns, campaign)
	}

	sh.respondJSON(w, http.StatusOK, campaigns)
}
