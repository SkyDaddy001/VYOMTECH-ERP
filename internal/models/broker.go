package models

import (
	"database/sql"
	"time"
)

// ============================================================================
// BrokerProfile Model
// ============================================================================
type BrokerProfile struct {
	ID       int64 `gorm:"primaryKey" json:"id"`
	TenantID int64 `gorm:"index" json:"tenant_id"`

	// Broker Details
	BrokerName         string  `json:"broker_name"`
	BrokerEmail        *string `json:"broker_email"`
	BrokerPhone        *string `json:"broker_phone"`
	BrokerLicenseNo    *string `gorm:"uniqueIndex" json:"broker_license_no"`
	RERARegistrationNo *string `gorm:"uniqueIndex" json:"rera_registration_no"`
	PANNo              *string `json:"pan_no"`
	GSTNo              *string `json:"gst_no"`

	// Business Details
	BrokerFirmName     *string `json:"broker_firm_name"`
	FirmRegistrationNo *string `json:"firm_registration_no"`
	BusinessAddress    *string `json:"business_address"`
	BusinessCity       *string `json:"business_city"`
	BusinessState      *string `json:"business_state"`
	BusinessPostalCode *string `json:"business_postal_code"`

	// Financial Details
	BankAccountNo   *string `json:"bank_account_no"`
	BankName        *string `json:"bank_name"`
	IFSCCode        *string `json:"ifsc_code"`
	BeneficiaryName *string `json:"beneficiary_name"`

	// Status & Compliance
	Status           string     `gorm:"index" json:"status"` // active, inactive, suspended
	KYCVerified      bool       `json:"kyc_verified"`
	KYCVerifiedDate  *time.Time `json:"kyc_verified_date"`
	ComplianceStatus *string    `json:"compliance_status"` // pending, approved, rejected

	// Performance Tracking
	TotalBookings         int     `json:"total_bookings"`
	TotalCommissionEarned float64 `json:"total_commission_earned"`
	TotalCommissionPaid   float64 `json:"total_commission_paid"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	CommissionStructures []BrokerCommissionStructure `gorm:"foreignKey:BrokerID" json:"commission_structures,omitempty"`
	BookingLinks         []BrokerBookingLink         `gorm:"foreignKey:BrokerID" json:"booking_links,omitempty"`
	Payouts              []BrokerCommissionPayout    `gorm:"foreignKey:BrokerID" json:"payouts,omitempty"`
}

// TableName specifies the table name
func (BrokerProfile) TableName() string {
	return "broker_profile"
}

// ============================================================================
// BrokerCommissionStructure Model
// ============================================================================
type BrokerCommissionStructure struct {
	ID       int64 `gorm:"primaryKey" json:"id"`
	TenantID int64 `gorm:"index" json:"tenant_id"`
	BrokerID int64 `gorm:"index" json:"broker_id"`

	// Commission Type
	CommissionType string `json:"commission_type"` // percentage, fixed_amount, slab_based

	// For Percentage-based
	CommissionPercentage *float64 `json:"commission_percentage"`

	// For Fixed Amount
	FixedAmount *float64 `json:"fixed_amount"`

	// Slab-based Commission
	Slab1MaxAmount            *float64 `json:"slab_1_max_amount"`
	Slab1CommissionPercentage *float64 `json:"slab_1_commission_percentage"`
	Slab2MinAmount            *float64 `json:"slab_2_min_amount"`
	Slab2MaxAmount            *float64 `json:"slab_2_max_amount"`
	Slab2CommissionPercentage *float64 `json:"slab_2_commission_percentage"`
	Slab3MinAmount            *float64 `json:"slab_3_min_amount"`
	Slab3CommissionPercentage *float64 `json:"slab_3_commission_percentage"`

	// Applicability
	ApplicableFor *string  `json:"applicable_for"` // residential, commercial, all
	MinUnitPrice  *float64 `json:"min_unit_price"`
	MaxUnitPrice  *float64 `json:"max_unit_price"`

	// Validity
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTill *time.Time `json:"effective_till"`

	// Status
	Status string `json:"status"` // active, inactive

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Broker *BrokerProfile `gorm:"foreignKey:BrokerID" json:"broker,omitempty"`
}

// TableName specifies the table name
func (BrokerCommissionStructure) TableName() string {
	return "broker_commission_structure"
}

// ============================================================================
// BrokerBookingLink Model
// ============================================================================
type BrokerBookingLink struct {
	ID        int64 `gorm:"primaryKey" json:"id"`
	TenantID  int64 `gorm:"index" json:"tenant_id"`
	BrokerID  int64 `gorm:"index" json:"broker_id"`
	BookingID int64 `gorm:"index" json:"booking_id"`

	// Booking Details
	UnitPrice     *float64 `json:"unit_price"`
	BookingAmount *float64 `json:"booking_amount"`

	// Commission Calculation
	CommissionStructureID *int64   `json:"commission_structure_id"`
	CommissionPercentage  *float64 `json:"commission_percentage"`
	CommissionAmount      *float64 `json:"commission_amount"`

	// Status
	BookingStatus    string `json:"booking_status"`                 // active, cancelled, completed
	CommissionStatus string `gorm:"index" json:"commission_status"` // pending, approved, paid, cancelled

	// Dates
	BookingDate  *time.Time `json:"booking_date"`
	ApprovalDate *time.Time `json:"approval_date"`
	PaymentDate  *time.Time `json:"payment_date"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Broker              *BrokerProfile             `gorm:"foreignKey:BrokerID" json:"broker,omitempty"`
	CommissionStructure *BrokerCommissionStructure `gorm:"foreignKey:CommissionStructureID" json:"commission_structure,omitempty"`
}

// TableName specifies the table name
func (BrokerBookingLink) TableName() string {
	return "broker_booking_link"
}

// ============================================================================
// BrokerCommissionPayout Model
// ============================================================================
type BrokerCommissionPayout struct {
	ID       int64 `gorm:"primaryKey" json:"id"`
	TenantID int64 `gorm:"index" json:"tenant_id"`
	BrokerID int64 `gorm:"index" json:"broker_id"`

	// Payout Period
	PayoutPeriodFrom  time.Time `gorm:"index" json:"payout_period_from"`
	PayoutPeriodTill  time.Time `gorm:"index" json:"payout_period_till"`
	PayoutReferenceNo *string   `gorm:"uniqueIndex" json:"payout_reference_no"`

	// Commission Calculation
	TotalCommission  float64  `json:"total_commission"`
	TDSAmount        *float64 `json:"tds_amount"`
	GSTAmount        *float64 `json:"gst_amount"`
	NetPayableAmount float64  `json:"net_payable_amount"`

	// Payment Details
	PaymentMode        *string    `json:"payment_mode"` // bank_transfer, cheque, cash, neft, rtgs
	PaymentDate        *time.Time `json:"payment_date"`
	PaymentReferenceNo *string    `json:"payment_reference_no"`
	BankName           *string    `json:"bank_name"`

	// Status
	Status       string     `gorm:"index" json:"status"` // pending, approved, paid, rejected, cancelled
	ApprovalBy   *int64     `json:"approval_by"`
	ApprovalDate *time.Time `json:"approval_date"`

	// Documents
	CalculationSheetURL *string `json:"calculation_sheet_url"`
	PaymentProofURL     *string `json:"payment_proof_url"`

	// Metadata
	CreatedBy *int64       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedBy *int64       `json:"updated_by"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at"`

	// Relations
	Broker *BrokerProfile `gorm:"foreignKey:BrokerID" json:"broker,omitempty"`
}

// TableName specifies the table name
func (BrokerCommissionPayout) TableName() string {
	return "broker_commission_payout"
}

// ============================================================================
// DTO Models for API Requests/Responses
// ============================================================================

// CreateBrokerRequest request to create new broker
type CreateBrokerRequest struct {
	BrokerName         string  `json:"broker_name" binding:"required"`
	BrokerEmail        *string `json:"broker_email"`
	BrokerPhone        *string `json:"broker_phone"`
	BrokerLicenseNo    *string `json:"broker_license_no"`
	RERARegistrationNo *string `json:"rera_registration_no"`
	PANNo              *string `json:"pan_no"`
	GSTNo              *string `json:"gst_no"`
	BrokerFirmName     *string `json:"broker_firm_name"`
	BusinessAddress    *string `json:"business_address"`
	BankAccountNo      *string `json:"bank_account_no"`
	BankName           *string `json:"bank_name"`
	IFSCCode           *string `json:"ifsc_code"`
}

// UpdateBrokerRequest request to update broker
type UpdateBrokerRequest struct {
	BrokerEmail   *string `json:"broker_email"`
	BrokerPhone   *string `json:"broker_phone"`
	Status        *string `json:"status"`
	KYCVerified   *bool   `json:"kyc_verified"`
	BankAccountNo *string `json:"bank_account_no"`
	BankName      *string `json:"bank_name"`
}

// CreateCommissionStructureRequest request to create commission structure
type CreateCommissionStructureRequest struct {
	BrokerID             int64    `json:"broker_id" binding:"required"`
	CommissionType       string   `json:"commission_type" binding:"required"`
	CommissionPercentage *float64 `json:"commission_percentage"`
	FixedAmount          *float64 `json:"fixed_amount"`
	ApplicableFor        *string  `json:"applicable_for"`
	MinUnitPrice         *float64 `json:"min_unit_price"`
	MaxUnitPrice         *float64 `json:"max_unit_price"`
	EffectiveFrom        string   `json:"effective_from" binding:"required"`
}

// CreateBookingLinkRequest request to link booking to broker
type CreateBookingLinkRequest struct {
	BrokerID              int64    `json:"broker_id" binding:"required"`
	BookingID             int64    `json:"booking_id" binding:"required"`
	CommissionStructureID int64    `json:"commission_structure_id" binding:"required"`
	CommissionPercentage  *float64 `json:"commission_percentage"`
	CommissionAmount      *float64 `json:"commission_amount"`
}

// CreatePayoutRequest request to create commission payout
type CreatePayoutRequest struct {
	BrokerID         int64    `json:"broker_id" binding:"required"`
	PayoutPeriodFrom string   `json:"payout_period_from" binding:"required"`
	PayoutPeriodTill string   `json:"payout_period_till" binding:"required"`
	TotalCommission  float64  `json:"total_commission" binding:"required"`
	NetPayableAmount float64  `json:"net_payable_amount" binding:"required"`
	TDSAmount        *float64 `json:"tds_amount"`
	GSTAmount        *float64 `json:"gst_amount"`
}

// BrokerPerformanceResponse response with broker performance metrics
type BrokerPerformanceResponse struct {
	BrokerID                int64   `json:"broker_id"`
	BrokerName              string  `json:"broker_name"`
	TotalBookings           int     `json:"total_bookings"`
	TotalCommissionEarned   float64 `json:"total_commission_earned"`
	TotalCommissionPaid     float64 `json:"total_commission_paid"`
	PendingPayout           float64 `json:"pending_payout"`
	AvgCommissionPerBooking float64 `json:"avg_commission_per_booking"`
}
