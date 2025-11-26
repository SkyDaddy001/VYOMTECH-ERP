package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ============================================================================
// REPORTING AND ANALYTICS HANDLERS
// ============================================================================

// LeadFunnelAnalysis - GET /api/v1/sales/reports/funnel
func (sh *SalesHandler) LeadFunnelAnalysis(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			DATE_TRUNC('month', created_at)::DATE as month,
			COUNT(CASE WHEN status = 'new' THEN 1 END) as leads_new,
			COUNT(CASE WHEN status = 'contacted' THEN 1 END) as leads_contacted,
			COUNT(CASE WHEN status = 'qualified' THEN 1 END) as leads_qualified,
			COUNT(CASE WHEN status = 'negotiation' THEN 1 END) as leads_negotiation,
			COUNT(CASE WHEN status = 'converted' THEN 1 END) as leads_converted,
			COUNT(CASE WHEN status = 'lost' THEN 1 END) as leads_lost,
			COUNT(*) as total_leads,
			ROUND(CAST(COUNT(CASE WHEN status = 'converted' THEN 1 END) AS NUMERIC) / 
				  NULLIF(COUNT(*), 0) * 100, 2) as conversion_rate
		FROM sales_leads
		WHERE tenant_id = $1 AND deleted_at IS NULL
		GROUP BY DATE_TRUNC('month', created_at)::DATE
		ORDER BY month DESC
		LIMIT 12
	`

	rows, err := sh.DB.Query(query, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch funnel analysis")
		return
	}
	defer rows.Close()

	type FunnelData struct {
		Month            time.Time `json:"month"`
		LeadsNew         int       `json:"leads_new"`
		LeadsContacted   int       `json:"leads_contacted"`
		LeadsQualified   int       `json:"leads_qualified"`
		LeadsNegotiation int       `json:"leads_negotiation"`
		LeadsConverted   int       `json:"leads_converted"`
		LeadsLost        int       `json:"leads_lost"`
		TotalLeads       int       `json:"total_leads"`
		ConversionRate   float64   `json:"conversion_rate"`
	}

	data := []FunnelData{}

	for rows.Next() {
		var record FunnelData
		if err := rows.Scan(
			&record.Month, &record.LeadsNew, &record.LeadsContacted, &record.LeadsQualified,
			&record.LeadsNegotiation, &record.LeadsConverted, &record.LeadsLost,
			&record.TotalLeads, &record.ConversionRate,
		); err != nil {
			continue
		}
		data = append(data, record)
	}

	sh.respondJSON(w, http.StatusOK, data)
}

// LeadSourcePerformance - GET /api/v1/sales/reports/source-performance
func (sh *SalesHandler) LeadSourcePerformance(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			source as source_type,
			COUNT(id) as total_leads,
			COUNT(CASE WHEN converted_to_customer = true THEN 1 END) as converted_leads,
			ROUND(CAST(COUNT(CASE WHEN converted_to_customer = true THEN 1 END) AS NUMERIC) / 
				  NULLIF(COUNT(id), 0) * 100, 2) as conversion_rate,
			COUNT(CASE WHEN status = 'lost' THEN 1 END) as lost_leads,
			SUM(CASE WHEN probability > 0 THEN 1 ELSE 0 END) as qualified_leads
		FROM sales_leads
		WHERE tenant_id = $1 AND deleted_at IS NULL
		GROUP BY source
		ORDER BY total_leads DESC
	`

	rows, err := sh.DB.Query(query, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch source performance")
		return
	}
	defer rows.Close()

	type SourceData struct {
		SourceType     string  `json:"source_type"`
		TotalLeads     int     `json:"total_leads"`
		ConvertedLeads int     `json:"converted_leads"`
		ConversionRate float64 `json:"conversion_rate"`
		LostLeads      int     `json:"lost_leads"`
		QualifiedLeads int     `json:"qualified_leads"`
	}

	data := []SourceData{}

	for rows.Next() {
		var record SourceData
		if err := rows.Scan(
			&record.SourceType, &record.TotalLeads, &record.ConvertedLeads,
			&record.ConversionRate, &record.LostLeads, &record.QualifiedLeads,
		); err != nil {
			continue
		}
		data = append(data, record)
	}

	sh.respondJSON(w, http.StatusOK, data)
}

// BookingSummary - GET /api/v1/sales/reports/bookings
func (sh *SalesHandler) BookingSummary(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			status,
			COUNT(id) as booking_count,
			SUM(booking_amount) as total_booking_amount,
			SUM(units_booked) as total_units_booked,
			AVG(booking_amount) as avg_booking_amount
		FROM sales_bookings
		WHERE tenant_id = $1 AND deleted_at IS NULL
		GROUP BY status
		ORDER BY booking_count DESC
	`

	rows, err := sh.DB.Query(query, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch booking summary")
		return
	}
	defer rows.Close()

	type BookingData struct {
		Status             string  `json:"status"`
		BookingCount       int     `json:"booking_count"`
		TotalBookingAmount float64 `json:"total_booking_amount"`
		TotalUnitsBooked   int     `json:"total_units_booked"`
		AvgBookingAmount   float64 `json:"avg_booking_amount"`
	}

	data := []BookingData{}

	for rows.Next() {
		var record BookingData
		if err := rows.Scan(
			&record.Status, &record.BookingCount, &record.TotalBookingAmount,
			&record.TotalUnitsBooked, &record.AvgBookingAmount,
		); err != nil {
			continue
		}
		data = append(data, record)
	}

	sh.respondJSON(w, http.StatusOK, data)
}

// CustomerLedgerSummary - GET /api/v1/sales/reports/customer-ledger/{customer_id}
func (sh *SalesHandler) CustomerLedgerSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID := vars["customer_id"]

	query := `
		SELECT 
			SUM(CASE WHEN transaction_type = 'invoice' THEN debit_amount ELSE 0 END) as total_invoiced,
			SUM(CASE WHEN transaction_type = 'payment' THEN credit_amount ELSE 0 END) as total_paid,
			SUM(CASE WHEN transaction_type = 'credit_note' THEN credit_amount ELSE 0 END) as total_credit_notes,
			(SUM(CASE WHEN transaction_type = 'invoice' THEN debit_amount ELSE 0 END) - 
			 SUM(CASE WHEN transaction_type = 'payment' THEN credit_amount ELSE 0 END) -
			 SUM(CASE WHEN transaction_type = 'credit_note' THEN credit_amount ELSE 0 END)) as outstanding_balance,
			COUNT(CASE WHEN transaction_type = 'invoice' THEN 1 END) as total_transactions
		FROM sales_account_ledgers
		WHERE customer_id = $1 AND tenant_id = $2 AND status = 'active'
	`

	var record struct {
		TotalInvoiced      float64 `json:"total_invoiced"`
		TotalPaid          float64 `json:"total_paid"`
		TotalCreditNotes   float64 `json:"total_credit_notes"`
		OutstandingBalance float64 `json:"outstanding_balance"`
		TotalTransactions  int     `json:"total_transactions"`
	}

	err := sh.DB.QueryRow(query, customerID, r.Header.Get("X-Tenant-ID")).Scan(
		&record.TotalInvoiced, &record.TotalPaid, &record.TotalCreditNotes,
		&record.OutstandingBalance, &record.TotalTransactions,
	)

	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch ledger summary")
		return
	}

	sh.respondJSON(w, http.StatusOK, record)
}

// MilestoneTimeline - GET /api/v1/sales/reports/milestone-timeline/{lead_id}
func (sh *SalesHandler) MilestoneTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["lead_id"]

	query := `
		SELECT 
			slm.id,
			slm.lead_id,
			sl.first_name,
			sl.last_name,
			slm.milestone_type,
			slm.milestone_date,
			slm.milestone_time,
			slm.notes,
			slm.status_before,
			slm.status_after,
			LEAD(slm.milestone_date) OVER (PARTITION BY slm.lead_id ORDER BY slm.milestone_date) as next_milestone_date,
			(LEAD(slm.milestone_date) OVER (PARTITION BY slm.lead_id ORDER BY slm.milestone_date) - slm.milestone_date)::INT as days_to_next_milestone
		FROM sales_lead_milestones slm
		JOIN sales_leads sl ON slm.lead_id = sl.id
		WHERE slm.lead_id = $1 AND slm.tenant_id = $2
		ORDER BY slm.milestone_date
	`

	rows, err := sh.DB.Query(query, leadID, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch timeline")
		return
	}
	defer rows.Close()

	type TimelineData struct {
		ID                  string     `json:"id"`
		LeadID              string     `json:"lead_id"`
		FirstName           string     `json:"first_name"`
		LastName            string     `json:"last_name"`
		MilestoneType       string     `json:"milestone_type"`
		MilestoneDate       time.Time  `json:"milestone_date"`
		MilestoneTime       *string    `json:"milestone_time"`
		Notes               *string    `json:"notes"`
		StatusBefore        *string    `json:"status_before"`
		StatusAfter         *string    `json:"status_after"`
		NextMilestoneDate   *time.Time `json:"next_milestone_date"`
		DaysToNextMilestone *int       `json:"days_to_next_milestone"`
	}

	data := []TimelineData{}

	for rows.Next() {
		var record TimelineData
		if err := rows.Scan(
			&record.ID, &record.LeadID, &record.FirstName, &record.LastName,
			&record.MilestoneType, &record.MilestoneDate, &record.MilestoneTime,
			&record.Notes, &record.StatusBefore, &record.StatusAfter,
			&record.NextMilestoneDate, &record.DaysToNextMilestone,
		); err != nil {
			continue
		}
		data = append(data, record)
	}

	sh.respondJSON(w, http.StatusOK, data)
}

// LeadEngagementStats - GET /api/v1/sales/reports/engagement-stats/{lead_id}
func (sh *SalesHandler) LeadEngagementStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leadID := vars["lead_id"]

	query := `
		SELECT 
			engagement_type,
			COUNT(*) as count,
			SUM(CASE WHEN response_received = true THEN 1 ELSE 0 END) as responses_received,
			AVG(duration_seconds) as avg_duration,
			MAX(engagement_date) as last_engagement
		FROM sales_lead_engagement
		WHERE lead_id = $1 AND tenant_id = $2
		GROUP BY engagement_type
		ORDER BY count DESC
	`

	rows, err := sh.DB.Query(query, leadID, r.Header.Get("X-Tenant-ID"))
	if err != nil {
		sh.respondError(w, http.StatusInternalServerError, "Failed to fetch engagement stats")
		return
	}
	defer rows.Close()

	type EngagementStats struct {
		EngagementType    string     `json:"engagement_type"`
		Count             int        `json:"count"`
		ResponsesReceived int        `json:"responses_received"`
		AvgDuration       *float64   `json:"avg_duration"`
		LastEngagement    *time.Time `json:"last_engagement"`
	}

	data := []EngagementStats{}

	for rows.Next() {
		var record EngagementStats
		if err := rows.Scan(
			&record.EngagementType, &record.Count, &record.ResponsesReceived,
			&record.AvgDuration, &record.LastEngagement,
		); err != nil {
			continue
		}
		data = append(data, record)
	}

	sh.respondJSON(w, http.StatusOK, data)
}

// DashboardMetrics - GET /api/v1/sales/reports/dashboard
func (sh *SalesHandler) DashboardMetrics(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")

	type DashboardMetrics struct {
		TotalLeads          int     `json:"total_leads"`
		NewLeads            int     `json:"new_leads"`
		QualifiedLeads      int     `json:"qualified_leads"`
		ConvertedLeads      int     `json:"converted_leads"`
		ConversionRate      float64 `json:"conversion_rate"`
		ActiveCustomers     int     `json:"active_customers"`
		TotalBookings       int     `json:"total_bookings"`
		BookedAmount        float64 `json:"booked_amount"`
		OutstandingBalance  float64 `json:"outstanding_balance"`
		EngagementThisMonth int     `json:"engagement_this_month"`
		PendingFollowUps    int     `json:"pending_follow_ups"`
	}

	metrics := DashboardMetrics{}

	// Total and converted leads
	sh.DB.QueryRow(`
		SELECT 
			COUNT(*),
			COUNT(CASE WHEN status = 'new' THEN 1 END),
			COUNT(CASE WHEN converted_to_customer = true THEN 1 END),
			ROUND(CAST(COUNT(CASE WHEN converted_to_customer = true THEN 1 END) AS NUMERIC) / 
				  NULLIF(COUNT(*), 0) * 100, 2)
		FROM sales_leads WHERE tenant_id = $1 AND deleted_at IS NULL
	`, tenantID).Scan(&metrics.TotalLeads, &metrics.NewLeads, &metrics.ConvertedLeads, &metrics.ConversionRate)

	// Qualified and active customers
	sh.DB.QueryRow(`
		SELECT 
			COUNT(CASE WHEN status = 'qualified' OR probability > 50 THEN 1 END),
			COUNT(CASE WHEN status = 'active' THEN 1 END)
		FROM sales_leads sl LEFT JOIN sales_customers sc ON sl.customer_id = sc.id 
		WHERE sl.tenant_id = $1 AND sl.deleted_at IS NULL
	`, tenantID).Scan(&metrics.QualifiedLeads, &metrics.ActiveCustomers)

	// Bookings
	sh.DB.QueryRow(`
		SELECT 
			COUNT(*),
			SUM(booking_amount)
		FROM sales_bookings WHERE tenant_id = $1 AND deleted_at IS NULL AND status = 'confirmed'
	`, tenantID).Scan(&metrics.TotalBookings, &metrics.BookedAmount)

	// Outstanding balance
	sh.DB.QueryRow(`
		SELECT COALESCE(SUM(
			(CASE WHEN transaction_type = 'invoice' THEN debit_amount ELSE 0 END) - 
			(CASE WHEN transaction_type = 'payment' THEN credit_amount ELSE 0 END) -
			(CASE WHEN transaction_type = 'credit_note' THEN credit_amount ELSE 0 END)
		), 0)
		FROM sales_account_ledgers WHERE tenant_id = $1 AND status = 'active'
	`, tenantID).Scan(&metrics.OutstandingBalance)

	// Engagement this month
	sh.DB.QueryRow(`
		SELECT COUNT(*) FROM sales_lead_engagement 
		WHERE tenant_id = $1 AND DATE_TRUNC('month', engagement_date) = DATE_TRUNC('month', CURRENT_DATE)
	`, tenantID).Scan(&metrics.EngagementThisMonth)

	// Pending follow-ups
	sh.DB.QueryRow(`
		SELECT COUNT(*) FROM sales_lead_milestones 
		WHERE tenant_id = $1 AND follow_up_required = true AND follow_up_date <= CURRENT_DATE
	`, tenantID).Scan(&metrics.PendingFollowUps)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}
