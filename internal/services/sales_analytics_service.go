package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"vyomtech-backend/internal/models"
)

// SalesAnalyticsService handles business logic for sales analytics and reporting
type SalesAnalyticsService struct {
	db     *sql.DB
	logger *log.Logger
}

// NewSalesAnalyticsService creates a new sales analytics service instance
func NewSalesAnalyticsService(db *sql.DB, logger *log.Logger) *SalesAnalyticsService {
	return &SalesAnalyticsService{
		db:     db,
		logger: logger,
	}
}

// GetDB returns the database connection for direct queries
func (s *SalesAnalyticsService) GetDB() *sql.DB {
	return s.db
}

// GetMonthlySalesAnalysis retrieves monthly sales summary for a tenant
func (s *SalesAnalyticsService) GetMonthlySalesAnalysis(ctx context.Context, tenantID string, startDate, endDate *time.Time) ([]*models.SalesMonthlySummary, error) {
	query := `
		SELECT 
			b.tenant_id,
			DATE_FORMAT(b.booking_date, '%Y-%m-01') as financial_month,
			QUARTER(b.booking_date) as financial_quarter,
			YEAR(b.booking_date) as financial_year,
			COUNT(DISTINCT CASE WHEN b.status = 'active' THEN b.id END) as units_sold,
			COUNT(DISTINCT CASE WHEN u.status = 'available' THEN u.id END) as units_unsold,
			COUNT(DISTINCT u.id) as total_units,
			SUM(u.uds_per_sqft) as total_uds,
			SUM(CASE WHEN u.status = 'available' THEN u.uds_per_sqft ELSE 0 END) as unsold_uds,
			SUM(CASE WHEN b.status = 'active' THEN u.uds_per_sqft ELSE 0 END) as sold_uds,
			SUM(u.sbua) as total_sbua,
			SUM(CASE WHEN u.status = 'available' THEN u.sbua ELSE 0 END) as unsold_sbua,
			SUM(CASE WHEN b.status = 'active' THEN u.sbua ELSE 0 END) as sold_sbua,
			SUM(CASE WHEN b.status = 'active' THEN ucs.apartment_cost_excluding_govt ELSE 0 END) as sold_value,
			SUM(CASE WHEN b.status = 'active' THEN rd.gst_cost ELSE 0 END) as gst_total,
			SUM(CASE WHEN b.status = 'active' THEN (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) ELSE 0 END) as sold_value_with_gst,
			SUM(COALESCE(p.amount, 0)) as collections_done,
			SUM(CASE WHEN b.status = 'active' THEN (ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) ELSE 0 END) - SUM(COALESCE(p.amount, 0)) as pending_due
		FROM booking b
		LEFT JOIN property_unit u ON b.unit_id = u.id
		LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
		LEFT JOIN registration_details rd ON b.id = rd.booking_id
		LEFT JOIN payment p ON b.id = p.booking_id
		WHERE b.tenant_id = ?
	`

	args := []interface{}{tenantID}

	if startDate != nil && endDate != nil {
		query += ` AND b.booking_date BETWEEN ? AND ?`
		args = append(args, startDate, endDate)
	}

	query += ` GROUP BY b.tenant_id, DATE_FORMAT(b.booking_date, '%Y-%m-01'), QUARTER(b.booking_date), YEAR(b.booking_date)
		ORDER BY financial_month DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.logger.Printf("Error querying monthly sales analysis: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []*models.SalesMonthlySummary
	for rows.Next() {
		record := &models.SalesMonthlySummary{}
		var monthStr string
		var quarterInt int

		err := rows.Scan(
			&record.TenantID,
			&monthStr,
			&quarterInt,
			&record.FinancialYear,
			&record.UnitsSold,
			&record.UnitsUnsold,
			&record.TotalUnits,
			&record.TotalUDS,
			&record.UnsoldUDS,
			&record.SoldUDS,
			&record.TotalSBUA,
			&record.UnsoldSBUA,
			&record.SoldSBUA,
			nil, // sold_value (not in model)
			nil, // gst_total (not in model)
			&record.SoldValueWithGST,
			&record.CollectionsDone,
			&record.PendingDue,
		)
		if err != nil {
			s.logger.Printf("Error scanning monthly sales record: %v", err)
			return nil, err
		}

		// Parse month string to time.Time
		if monthTime, err := time.Parse("2006-01-02", monthStr); err == nil {
			record.FinancialMonth = monthTime
		}

		record.FinancialQuarter = fmt.Sprintf("Q%d", quarterInt)
		results = append(results, record)
	}

	return results, rows.Err()
}

// GetCollectionReport retrieves collection report for a quarter
func (s *SalesAnalyticsService) GetCollectionReport(ctx context.Context, tenantID string, quarter int, year string) ([]*models.CollectionReport, error) {
	query := `
		SELECT 
			p.tenant_id,
			DATE_FORMAT(p.payment_date, '%Y-%m-01') as collection_month,
			QUARTER(p.payment_date) as financial_quarter,
			YEAR(p.payment_date) as financial_year,
			SUM(p.amount) as overall_collections,
			SUM(CASE WHEN p.towards LIKE '%apartment%' THEN p.amount ELSE 0 END) as apartment_cost,
			SUM(CASE WHEN p.towards LIKE '%gst%' THEN p.amount ELSE 0 END) as gst_collected,
			SUM(CASE WHEN p.towards LIKE '%tds%' THEN p.amount ELSE 0 END) as tds_collected,
			SUM(CASE WHEN p.towards NOT IN ('apartment', 'gst', 'tds') THEN p.amount ELSE 0 END) as others_collected
		FROM payment p
		WHERE p.tenant_id = ? AND QUARTER(p.payment_date) = ? AND YEAR(p.payment_date) = ?
		GROUP BY p.tenant_id, DATE_FORMAT(p.payment_date, '%Y-%m-01'), QUARTER(p.payment_date), YEAR(p.payment_date)
		ORDER BY collection_month DESC
	`

	yearInt := 0
	fmt.Sscanf(year, "%d", &yearInt)

	rows, err := s.db.QueryContext(ctx, query, tenantID, quarter, yearInt)
	if err != nil {
		s.logger.Printf("Error querying collection report: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []*models.CollectionReport
	for rows.Next() {
		record := &models.CollectionReport{}
		var monthStr string
		var quarterInt int

		err := rows.Scan(
			&record.TenantID,
			&monthStr,
			&quarterInt,
			&record.FinancialYear,
			&record.OverallCollections,
			&record.ApartmentCost,
			&record.GSTCollected,
			&record.TDSCollected,
			&record.OthersCollected,
		)
		if err != nil {
			s.logger.Printf("Error scanning collection record: %v", err)
			return nil, err
		}

		// Parse month string to time.Time
		if monthTime, err := time.Parse("2006-01-02", monthStr); err == nil {
			record.CollectionMonth = monthTime
		}

		record.FinancialQuarter = fmt.Sprintf("Q%d", quarterInt)
		results = append(results, record)
	}

	return results, rows.Err()
}

// GetBankOwnPaymentAnalysis retrieves bank vs own payment analysis
func (s *SalesAnalyticsService) GetBankOwnPaymentAnalysis(ctx context.Context, tenantID string) ([]*models.BankOwnPaymentAnalysis, error) {
	query := `
		SELECT
			b.id,
			b.tenant_id,
			b.id,
			sl.first_name,
			bl.bank_name,
			bl.sanction_amount,
			COALESCE(bl.disbursed_amount, 0),
			bl.sanction_amount - COALESCE(bl.disbursed_amount, 0),
			SUM(CASE WHEN p.payment_mode IN ('cash', 'cheque', 'bank_transfer') THEN p.amount ELSE 0 END),
			(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) - COALESCE(bl.disbursed_amount, 0) - SUM(CASE WHEN p.payment_mode IN ('cash', 'cheque', 'bank_transfer') THEN p.amount ELSE 0 END),
			u.block,
			ucs.apartment_cost_excluding_govt,
			ps.payment_stage,
			ROUND((ps.payment_stage / 13) * 100, 2),
			ps.amount_due,
			COALESCE(SUM(p.amount), 0),
			CURDATE()
		FROM booking b
		LEFT JOIN property_unit u ON b.unit_id = u.id
		LEFT JOIN sales_lead sl ON b.lead_id = sl.id
		LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
		LEFT JOIN registration_details rd ON b.id = rd.booking_id
		LEFT JOIN bank_loan bl ON b.id = bl.booking_id
		LEFT JOIN payment_schedule ps ON b.id = ps.booking_id
		LEFT JOIN payment p ON b.id = p.booking_id
		WHERE b.tenant_id = ?
		GROUP BY b.id, b.tenant_id, sl.first_name, bl.bank_name, bl.sanction_amount, COALESCE(bl.disbursed_amount, 0), ucs.apartment_cost_excluding_govt, COALESCE(rd.gst_cost, 0), u.block, ps.payment_stage, ps.amount_due
	`

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		s.logger.Printf("Error querying bank payment analysis: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []*models.BankOwnPaymentAnalysis
	for rows.Next() {
		record := &models.BankOwnPaymentAnalysis{}
		var asOnDateStr string

		err := rows.Scan(
			&record.ID,
			&record.TenantID,
			&record.BookingID,
			&record.CustomerName,
			&record.BankName,
			&record.BankSanctioned,
			&record.BankLoanDisbursed,
			&record.LoanAvailableForDisbursement,
			&record.CustomerOwnPaid,
			&record.CustomerOwnDue,
			&record.Block,
			&record.ApartmentCost,
			&record.PaymentStage,
			&record.StagePercentage,
			&record.StageDue,
			&record.StageReceived,
			&asOnDateStr,
		)
		if err != nil {
			s.logger.Printf("Error scanning bank payment record: %v", err)
			return nil, err
		}

		// Parse date
		if asOnDate, err := time.Parse("2006-01-02", asOnDateStr); err == nil {
			record.AsOnDate = asOnDate
		}

		results = append(results, record)
	}

	return results, rows.Err()
}

// GetDashboardSummary retrieves overall sales dashboard summary
func (s *SalesAnalyticsService) GetDashboardSummary(ctx context.Context, tenantID string) (*models.SalesDashboardSummary, error) {
	query := `
		SELECT
			u.tenant_id,
			CURDATE(),
			COUNT(DISTINCT u.id),
			COUNT(DISTINCT CASE WHEN b.status = 'active' THEN u.id END),
			COUNT(DISTINCT CASE WHEN u.status = 'available' THEN u.id END),
			SUM(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)),
			SUM(COALESCE(p.amount, 0)),
			SUM(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0)) - SUM(COALESCE(p.amount, 0)),
			SUM(u.uds_per_sqft),
			SUM(u.sbua),
			COUNT(DISTINCT CASE WHEN b.agreement_date IS NOT NULL THEN b.id END),
			COUNT(DISTINCT CASE WHEN b.agreement_date IS NULL THEN b.id END),
			ROUND((SUM(COALESCE(p.amount, 0)) / SUM(ucs.apartment_cost_excluding_govt + COALESCE(rd.gst_cost, 0))) * 100, 2),
			ROUND((COUNT(DISTINCT CASE WHEN b.status = 'active' THEN u.id END) / COUNT(DISTINCT u.id)) * 100, 2)
		FROM property_unit u
		LEFT JOIN booking b ON u.id = b.unit_id
		LEFT JOIN unit_cost_sheet ucs ON u.id = ucs.unit_id
		LEFT JOIN registration_details rd ON b.id = rd.booking_id
		LEFT JOIN payment p ON b.id = p.booking_id
		WHERE u.tenant_id = ?
		GROUP BY u.tenant_id
	`

	summary := &models.SalesDashboardSummary{}
	var reportDateStr string

	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(
		&summary.TenantID,
		&reportDateStr,
		&summary.TotalUnits,
		&summary.SoldUnits,
		&summary.UnsoldUnits,
		&summary.TotalValue,
		&summary.TotalCollected,
		&summary.TotalPending,
		&summary.TotalUDS,
		&summary.TotalSBUA,
		&summary.AgreementsSigned,
		&summary.AgreementsPending,
		&summary.CollectionPercentage,
		&summary.OccupancyPercentage,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data found for tenant")
		}
		s.logger.Printf("Error querying dashboard summary: %v", err)
		return nil, err
	}

	// Parse report date
	if reportDate, err := time.Parse("2006-01-02", reportDateStr); err == nil {
		summary.ReportDate = reportDate
	}

	return summary, nil
}
