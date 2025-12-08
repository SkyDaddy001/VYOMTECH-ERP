package models

import "time"

// SalesMonthlySummary represents monthly sales summary data
type SalesMonthlySummary struct {
	ID                string    `db:"id" json:"id"`
	TenantID          string    `db:"tenant_id" json:"tenant_id"`
	FinancialMonth    time.Time `db:"financial_month" json:"financial_month"`
	FinancialQuarter  string    `db:"financial_quarter" json:"financial_quarter"`
	FinancialYear     string    `db:"financial_year" json:"financial_year"`
	UnitsSold         int       `db:"units_sold" json:"units_sold"`
	UnitsUnsold       int       `db:"units_unsold" json:"units_unsold"`
	TotalUnits        int       `db:"total_units" json:"total_units"`
	TotalUDS          float64   `db:"total_uds" json:"total_uds"`
	UnsoldUDS         float64   `db:"unsold_uds" json:"unsold_uds"`
	SoldUDS           float64   `db:"sold_uds" json:"sold_uds"`
	TotalSBUA         float64   `db:"total_sbua" json:"total_sbua"`
	UnsoldSBUA        float64   `db:"unsold_sbua" json:"unsold_sbua"`
	SoldSBUA          float64   `db:"sold_sbua" json:"sold_sbua"`
	SoldValueWithGST  float64   `db:"sold_value_with_gst" json:"sold_value_with_gst"`
	CollectionsDone   float64   `db:"collections_done" json:"collections_done"`
	PendingDue        float64   `db:"pending_due" json:"pending_due"`
	UnsoldReceivables float64   `db:"unsold_receivables" json:"unsold_receivables"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// SoldUnitsTracking represents individual sold unit details
type SoldUnitsTracking struct {
	ID                   string    `db:"id" json:"id"`
	TenantID             string    `db:"tenant_id" json:"tenant_id"`
	UnitID               string    `db:"unit_id" json:"unit_id"`
	BookingID            string    `db:"booking_id" json:"booking_id"`
	Block                string    `db:"block" json:"block"`
	ApartmentNo          string    `db:"apartment_no" json:"apartment_no"`
	SoldMonth            time.Time `db:"sold_month" json:"sold_month"`
	SoldValueExclTax     float64   `db:"sold_value_excl_tax" json:"sold_value_excl_tax"`
	ApartmentCostWithGST float64   `db:"apartment_cost_with_gst" json:"apartment_cost_with_gst"`
	TDSAmount            float64   `db:"tds_amount" json:"tds_amount"`
	GSTAmount            float64   `db:"gst_amount" json:"gst_amount"`
	RegistrationCost     float64   `db:"registration_cost" json:"registration_cost"`
	MODTAmount           float64   `db:"modt_amount" json:"modt_amount"`
	OtherWorks           float64   `db:"other_works" json:"other_works"`
	MaintenanceCharge    float64   `db:"maintenance_charge" json:"maintenance_charge"`
	CorpusAmount         float64   `db:"corpus_amount" json:"corpus_amount"`
	TotalReceivable      float64   `db:"total_receivable" json:"total_receivable"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

// CollectionReport represents collection data by quarter
type CollectionReport struct {
	ID                      string    `db:"id" json:"id"`
	TenantID                string    `db:"tenant_id" json:"tenant_id"`
	FinancialQuarter        string    `db:"financial_quarter" json:"financial_quarter"`
	FinancialYear           string    `db:"financial_year" json:"financial_year"`
	CollectionMonth         time.Time `db:"collection_month" json:"collection_month"`
	OverallCollections      float64   `db:"overall_collections" json:"overall_collections"`
	ApartmentCost           float64   `db:"apartment_cost" json:"apartment_cost"`
	GSTCollected            float64   `db:"gst_collected" json:"gst_collected"`
	TDSCollected            float64   `db:"tds_collected" json:"tds_collected"`
	OthersCollected         float64   `db:"others_collected" json:"others_collected"`
	TotalQuarterCollections float64   `db:"total_quarter_collections" json:"total_quarter_collections"`
	QuarterApartmentCost    float64   `db:"quarter_apartment_cost" json:"quarter_apartment_cost"`
	QuarterGST              float64   `db:"quarter_gst" json:"quarter_gst"`
	QuarterTDS              float64   `db:"quarter_tds" json:"quarter_tds"`
	QuarterOthers           float64   `db:"quarter_others" json:"quarter_others"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
}

// SalesConsiderationReport represents agreement signed/pending data
type SalesConsiderationReport struct {
	ID                string    `db:"id" json:"id"`
	TenantID          string    `db:"tenant_id" json:"tenant_id"`
	AgreementMonth    time.Time `db:"agreement_month" json:"agreement_month"`
	FinancialQuarter  string    `db:"financial_quarter" json:"financial_quarter"`
	FinancialYear     string    `db:"financial_year" json:"financial_year"`
	AgreementsSigned  int       `db:"agreements_signed" json:"agreements_signed"`
	AgreementsPending int       `db:"agreements_pending" json:"agreements_pending"`
	UnitsCount        int       `db:"units_count" json:"units_count"`
	TotalUDS          float64   `db:"total_uds" json:"total_uds"`
	TotalSBUA         float64   `db:"total_sbua" json:"total_sbua"`
	SoldValueWithGST  float64   `db:"sold_value_with_gst" json:"sold_value_with_gst"`
	GSTAmount         float64   `db:"gst_amount" json:"gst_amount"`
	CollectionsDone   float64   `db:"collections_done" json:"collections_done"`
	PendingDue        float64   `db:"pending_due" json:"pending_due"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// BankOwnPaymentAnalysis represents bank loan disbursement vs customer own contributions
type BankOwnPaymentAnalysis struct {
	ID                           string    `db:"id" json:"id"`
	TenantID                     string    `db:"tenant_id" json:"tenant_id"`
	BookingID                    string    `db:"booking_id" json:"booking_id"`
	CustomerName                 string    `db:"customer_name" json:"customer_name"`
	BankName                     string    `db:"bank_name" json:"bank_name"`
	BankSanctioned               float64   `db:"bank_sanctioned" json:"bank_sanctioned"`
	BankLoanDisbursed            float64   `db:"bank_loan_disbursed" json:"bank_loan_disbursed"`
	LoanAvailableForDisbursement float64   `db:"loan_available_for_disbursement" json:"loan_available_for_disbursement"`
	CustomerOwnPaid              float64   `db:"customer_own_paid" json:"customer_own_paid"`
	CustomerOwnDue               float64   `db:"customer_own_due" json:"customer_own_due"`
	Block                        string    `db:"block" json:"block"`
	ApartmentCost                float64   `db:"apartment_cost_excluding_govt" json:"apartment_cost_excluding_govt"`
	PaymentStage                 int       `db:"payment_stage" json:"payment_stage"`
	StagePercentage              float64   `db:"stage_percentage" json:"stage_percentage"`
	StageDue                     float64   `db:"stage_due" json:"stage_due"`
	StageReceived                float64   `db:"stage_received" json:"stage_received"`
	AsOnDate                     time.Time `db:"as_on_date" json:"as_on_date"`
}

// PendingAgreement represents units with pending agreements
type PendingAgreement struct {
	ID                    string    `db:"id" json:"id"`
	TenantID              string    `db:"tenant_id" json:"tenant_id"`
	UnitID                string    `db:"unit_id" json:"unit_id"`
	Block                 string    `db:"block" json:"block"`
	ApartmentNo           string    `db:"apartment_no" json:"apartment_no"`
	AgreementPendingSince time.Time `db:"agreement_pending_since" json:"agreement_pending_since"`
	UnitsCount            int       `db:"units_count" json:"units_count"`
	TotalUDS              float64   `db:"total_uds" json:"total_uds"`
	TotalSBUA             float64   `db:"total_sbua" json:"total_sbua"`
	SoldValueWithGST      float64   `db:"sold_value_with_gst" json:"sold_value_with_gst"`
	GSTAmount             float64   `db:"gst_amount" json:"gst_amount"`
	CollectionsDone       float64   `db:"collections_done" json:"collections_done"`
	PendingDue            float64   `db:"pending_due" json:"pending_due"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}

// SalesDashboardSummary represents cached aggregated data for dashboard
type SalesDashboardSummary struct {
	ID                   string    `db:"id" json:"id"`
	TenantID             string    `db:"tenant_id" json:"tenant_id"`
	ReportDate           time.Time `db:"report_date" json:"report_date"`
	TotalUnits           int       `db:"total_units" json:"total_units"`
	SoldUnits            int       `db:"sold_units" json:"sold_units"`
	UnsoldUnits          int       `db:"unsold_units" json:"unsold_units"`
	TotalValue           float64   `db:"total_value" json:"total_value"`
	TotalCollected       float64   `db:"total_collected" json:"total_collected"`
	TotalPending         float64   `db:"total_pending" json:"total_pending"`
	TotalUDS             float64   `db:"total_uds" json:"total_uds"`
	TotalSBUA            float64   `db:"total_sbua" json:"total_sbua"`
	AgreementsSigned     int       `db:"agreements_signed" json:"agreements_signed"`
	AgreementsPending    int       `db:"agreements_pending" json:"agreements_pending"`
	CollectionPercentage float64   `db:"collection_percentage" json:"collection_percentage"`
	OccupancyPercentage  float64   `db:"occupancy_percentage" json:"occupancy_percentage"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

// Analytics Request/Response DTOs

// MonthlySalesRequest is request to get monthly sales data
type MonthlySalesRequest struct {
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	FinancialYear string    `json:"financial_year"`
}

// CollectionReportRequest is request for collection report
type CollectionReportRequest struct {
	FinancialQuarter string `json:"financial_quarter"`
	FinancialYear    string `json:"financial_year"`
}

// SalesAnalyticsResponse provides comprehensive sales data
type SalesAnalyticsResponse struct {
	TenantID             string                 `json:"tenant_id"`
	MonthlySummary       []SalesMonthlySummary  `json:"monthly_summary"`
	SoldUnits            []SoldUnitsTracking    `json:"sold_units"`
	DashboardSummary     *SalesDashboardSummary `json:"dashboard_summary"`
	CollectionPercentage float64                `json:"collection_percentage"`
	OccupancyPercentage  float64                `json:"occupancy_percentage"`
	GeneratedAt          time.Time              `json:"generated_at"`
}

// CollectionReportResponse provides collection breakdown
type CollectionReportResponse struct {
	TenantID                string             `json:"tenant_id"`
	FinancialQuarter        string             `json:"financial_quarter"`
	FinancialYear           string             `json:"financial_year"`
	MonthlyBreakup          []CollectionReport `json:"monthly_breakup"`
	TotalQuarterCollections float64            `json:"total_quarter_collections"`
	QuarterApartmentCost    float64            `json:"quarter_apartment_cost"`
	QuarterGST              float64            `json:"quarter_gst"`
	QuarterTDS              float64            `json:"quarter_tds"`
	QuarterOthers           float64            `json:"quarter_others"`
	GeneratedAt             time.Time          `json:"generated_at"`
}

// BankOwnPaymentResponse provides bank vs own payment analysis
type BankOwnPaymentResponse struct {
	TenantID            string                   `json:"tenant_id"`
	PaymentAnalysis     []BankOwnPaymentAnalysis `json:"payment_analysis"`
	TotalBankSanctioned float64                  `json:"total_bank_sanctioned"`
	TotalBankDisbursed  float64                  `json:"total_bank_disbursed"`
	TotalBankDue        float64                  `json:"total_bank_due"`
	TotalOwnPaid        float64                  `json:"total_own_paid"`
	TotalOwnDue         float64                  `json:"total_own_due"`
	GeneratedAt         time.Time                `json:"generated_at"`
}

// SalesConsiderationResponse provides agreement data
type SalesConsiderationResponse struct {
	TenantID                string                     `json:"tenant_id"`
	AgreementsSigned        int                        `json:"agreements_signed"`
	AgreementsPending       int                        `json:"agreements_pending"`
	SignedConsolidation     []SalesConsiderationReport `json:"signed_consolidation"`
	PendingAgreements       []PendingAgreement         `json:"pending_agreements"`
	TotalPendingValue       float64                    `json:"total_pending_value"`
	TotalPendingCollections float64                    `json:"total_pending_collections"`
	GeneratedAt             time.Time                  `json:"generated_at"`
}
