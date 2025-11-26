package models

import (
	"time"
)

// ============================================
// PROJECT & PROPERTY MODELS
// ============================================

// PropertyProject represents a real estate project
type PropertyProject struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	ProjectName        string     `json:"project_name"`
	ProjectCode        string     `json:"project_code"`
	Location           string     `json:"location"`
	City               string     `json:"city"`
	State              string     `json:"state"`
	PostalCode         string     `json:"postal_code"`
	TotalUnits         int        `json:"total_units"`
	TotalArea          float64    `json:"total_area"`
	ProjectType        string     `json:"project_type"` // residential, commercial, mixed
	Status             string     `json:"status"`       // planning, under_construction, ready, sold_out
	LaunchDate         *time.Time `json:"launch_date"`
	ExpectedCompletion *time.Time `json:"expected_completion"`
	ActualCompletion   *time.Time `json:"actual_completion"`
	NOCStatus          string     `json:"noc_status"` // pending, approved, rejected
	NOCDate            *time.Time `json:"noc_date"`
	DeveloperName      string     `json:"developer_name"`
	ArchitectName      string     `json:"architect_name"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
	CreatedBy          *string    `json:"created_by"`
}

// PropertyBlock represents a block within a project
type PropertyBlock struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	ProjectID  string     `json:"project_id"`
	BlockName  string     `json:"block_name"`
	BlockCode  string     `json:"block_code"`
	WingName   string     `json:"wing_name"`
	TotalUnits int        `json:"total_units"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// PropertyUnit represents an individual property unit
type PropertyUnit struct {
	ID                    string     `json:"id"`
	TenantID              string     `json:"tenant_id"`
	ProjectID             string     `json:"project_id"`
	BlockID               string     `json:"block_id"`
	UnitNumber            string     `json:"unit_number"`
	Floor                 int        `json:"floor"`
	UnitType              string     `json:"unit_type"` // 1BHK, 2BHK, 3BHK, shop, office
	Facing                string     `json:"facing"`    // north, south, east, west, corner
	CarpetArea            float64    `json:"carpet_area"`
	CarpetAreaWithBalcony float64    `json:"carpet_area_with_balcony"`
	UtilityArea           float64    `json:"utility_area"`
	PlinthArea            float64    `json:"plinth_area"`
	SBUA                  float64    `json:"sbua"` // Super Built Up Area
	UDSSqft               float64    `json:"uds_sqft"`
	Status                string     `json:"status"` // available, booked, sold, reserved
	AllotedTo             string     `json:"alloted_to"`
	AllotmentDate         *time.Time `json:"allotment_date"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`
}

// ============================================
// COST SHEET MODELS
// ============================================

// UnitCostSheet represents the pricing details of a property unit
type UnitCostSheet struct {
	ID                      string     `json:"id"`
	TenantID                string     `json:"tenant_id"`
	UnitID                  string     `json:"unit_id"`
	RatePerSqft             float64    `json:"rate_per_sqft"`
	SBUARate                float64    `json:"sbua_rate"`
	BasePrice               float64    `json:"base_price"`
	FRC                     float64    `json:"frc"` // Floor Rise Charge
	CarParkingCost          float64    `json:"car_parking_cost"`
	PLC                     float64    `json:"plc"` // Parking + Lease Charge
	StatutoryCharges        float64    `json:"statutory_charges"`
	OtherCharges            float64    `json:"other_charges"`
	LegalCharges            float64    `json:"legal_charges"`
	ApartmentCostExcGovt    float64    `json:"apartment_cost_exc_govt"`
	ApartmentCostIncGovt    float64    `json:"apartment_cost_inc_govt"`
	CompositeGuidelineValue float64    `json:"composite_guideline_value"`
	ActualSoldPrice         float64    `json:"actual_sold_price"`
	CarParkingType          string     `json:"car_parking_type"` // covered, open, none
	ParkingLocation         string     `json:"parking_location"`
	EffectiveDate           *time.Time `json:"effective_date"`
	ValidityDate            *time.Time `json:"validity_date"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at"`
	CreatedBy               *string    `json:"created_by"`
}

// ============================================
// CUSTOMER BOOKING MODELS
// ============================================

// CustomerBooking represents a booking/reservation of a property unit
type CustomerBooking struct {
	ID                      string     `json:"id"`
	TenantID                string     `json:"tenant_id"`
	UnitID                  string     `json:"unit_id"`
	LeadID                  *string    `json:"lead_id"`
	CustomerID              *string    `json:"customer_id"`
	BookingDate             time.Time  `json:"booking_date"`
	BookingReference        string     `json:"booking_reference"`
	BookingStatus           string     `json:"booking_status"` // active, cancelled, completed
	WelcomeDate             *time.Time `json:"welcome_date"`
	AllotmentDate           *time.Time `json:"allotment_date"`
	AgreementDate           *time.Time `json:"agreement_date"`
	RegistrationDate        *time.Time `json:"registration_date"`
	HandoverDate            *time.Time `json:"handover_date"`
	PossessionDate          *time.Time `json:"possession_date"`
	RatePerSqft             float64    `json:"rate_per_sqft"`
	CompositeGuidelineValue float64    `json:"composite_guideline_value"`
	CarParkingType          string     `json:"car_parking_type"`
	ParkingLocation         string     `json:"parking_location"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	DeletedAt               *time.Time `json:"deleted_at"`
}

// CustomerDetails represents detailed customer information for a booking
type CustomerDetails struct {
	ID                       string     `json:"id"`
	TenantID                 string     `json:"tenant_id"`
	BookingID                string     `json:"booking_id"`
	PrimaryName              string     `json:"primary_name"`
	PrimaryPhone             string     `json:"primary_phone"`
	PrimaryAlternatePhone    string     `json:"primary_alternate_phone"`
	PrimaryEmail             string     `json:"primary_email"`
	PrimaryCommunicationAddr string     `json:"primary_communication_address"`
	PrimaryPermanentAddr     string     `json:"primary_permanent_address"`
	PrimaryAadharNo          string     `json:"primary_aadhar_no"`
	PrimaryPanNo             string     `json:"primary_pan_no"`
	CoApplicant1Name         string     `json:"coapplicant1_name"`
	CoApplicant1Phone        string     `json:"coapplicant1_phone"`
	CoApplicant1Email        string     `json:"coapplicant1_email"`
	CoApplicant1Aadhar       string     `json:"coapplicant1_aadhar_no"`
	CoApplicant1Pan          string     `json:"coapplicant1_pan_no"`
	CoApplicant1Relation     string     `json:"coapplicant1_relation"`
	CoApplicant2Name         string     `json:"coapplicant2_name"`
	CoApplicant2Phone        string     `json:"coapplicant2_phone"`
	CoApplicant2Email        string     `json:"coapplicant2_email"`
	CoApplicant2Aadhar       string     `json:"coapplicant2_aadhar_no"`
	CoApplicant2Pan          string     `json:"coapplicant2_pan_no"`
	CoApplicant2Relation     string     `json:"coapplicant2_relation"`
	POAHolderName            string     `json:"poa_holder_name"`
	POADocumentNo            string     `json:"poa_document_no"`
	LifeCertificateNo        string     `json:"life_certificate_no"`
	BankName                 string     `json:"bank_name"`
	LoanSanctionDate         *time.Time `json:"loan_sanction_date"`
	SanctionAmount           float64    `json:"sanction_amount"`
	SalesExecutiveName       string     `json:"sales_executive_name"`
	SalesHeadName            string     `json:"sales_head_name"`
	BookingSource            string     `json:"booking_source"` // direct, broker, partner, site_visit
	MaintenanceCharges       float64    `json:"maintenance_charges"`
	CorpusCharges            float64    `json:"corpus_charges"`
	EBDeposit                float64    `json:"eb_deposit"`
	NOCReceivedDate          *time.Time `json:"noc_received_date"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at"`
}

// ============================================
// PAYMENT MODELS
// ============================================

// BookingPayment represents a payment received for a booking
type BookingPayment struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	BookingID     string     `json:"booking_id"`
	PaymentDate   time.Time  `json:"payment_date"`
	PaymentMode   string     `json:"payment_mode"` // cash, cheque, transfer, neft, rtgs, demand_draft
	PaidBy        string     `json:"paid_by"`
	ReceiptNumber string     `json:"receipt_number"`
	ReceiptDate   *time.Time `json:"receipt_date"`
	Towards       string     `json:"towards"` // advance, booking, installment_1, balance, etc.
	Amount        float64    `json:"amount"`
	ChequeNumber  string     `json:"cheque_number"`
	ChequeDate    *time.Time `json:"cheque_date"`
	BankName      string     `json:"bank_name"`
	TransactionID string     `json:"transaction_id"`
	Status        string     `json:"status"` // pending, cleared, bounced, cancelled
	Remarks       string     `json:"remarks"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	CreatedBy     *string    `json:"created_by"`
}

// PaymentSchedule represents scheduled payments for a booking
type PaymentSchedule struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	BookingID      string     `json:"booking_id"`
	ScheduleName   string     `json:"schedule_name"`
	PaymentStage   string     `json:"payment_stage"` // booking, agreement, possession, handover
	PaymentPercent float64    `json:"payment_percentage"`
	PaymentAmount  float64    `json:"payment_amount"`
	DueDate        time.Time  `json:"due_date"`
	AmountPaid     float64    `json:"amount_paid"`
	Outstanding    float64    `json:"outstanding"`
	Status         string     `json:"status"` // pending, partial, completed, overdue
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

// ============================================
// LEDGER MODELS
// ============================================

// CustomerAccountLedger represents transaction records for a customer
type CustomerAccountLedger struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	BookingID       string     `json:"booking_id"`
	CustomerID      *string    `json:"customer_id"`
	TransactionDate time.Time  `json:"transaction_date"`
	TransactionType string     `json:"transaction_type"` // credit, debit, adjustment
	Description     string     `json:"description"`
	DebitAmount     float64    `json:"debit_amount"`
	CreditAmount    float64    `json:"credit_amount"`
	OpeningBalance  float64    `json:"opening_balance"`
	ClosingBalance  float64    `json:"closing_balance"`
	PaymentID       *string    `json:"payment_id"`
	ReferenceNum    string     `json:"reference_number"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

// ============================================
// MILESTONE & TRACKING MODELS
// ============================================

// PropertyMilestone tracks important dates and campaign information for a booking
type PropertyMilestone struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	BookingID         string     `json:"booking_id"`
	CampaignID        *string    `json:"campaign_id"`
	CampaignName      string     `json:"campaign_name"`
	Source            string     `json:"source"` // direct, site_visit, broker, referral, digital, exhibition
	SubSource         string     `json:"subsource"`
	LeadGeneratedDate *time.Time `json:"lead_generated_date"`
	ReEngagedDate     *time.Time `json:"re_engaged_date"`
	SiteVisitDate     *time.Time `json:"site_visit_date"`
	ReVisitDate       *time.Time `json:"revisit_date"`
	BookingDate       *time.Time `json:"booking_date"`
	CancelledDate     *time.Time `json:"cancelled_date"`
	Status            string     `json:"status"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

// ============================================
// CONTROL SHEET MODEL
// ============================================

// ProjectControlSheet represents project configuration/attributes
type ProjectControlSheet struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	ProjectID      string     `json:"project_id"`
	AttributeName  string     `json:"attribute_name"`
	AttributeValue string     `json:"attribute_value"`
	AttributeType  string     `json:"attribute_type"` // text, number, date, percentage, boolean
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

// ============================================
// REQUEST/RESPONSE MODELS
// ============================================

// CreatePropertyProjectRequest for creating a new project
type CreatePropertyProjectRequest struct {
	ProjectName        string     `json:"project_name" validate:"required"`
	ProjectCode        string     `json:"project_code" validate:"required"`
	Location           string     `json:"location"`
	City               string     `json:"city"`
	State              string     `json:"state"`
	PostalCode         string     `json:"postal_code"`
	TotalUnits         int        `json:"total_units"`
	TotalArea          float64    `json:"total_area"`
	ProjectType        string     `json:"project_type"`
	Status             string     `json:"status"`
	LaunchDate         *time.Time `json:"launch_date"`
	ExpectedCompletion *time.Time `json:"expected_completion"`
	DeveloperName      string     `json:"developer_name"`
	ArchitectName      string     `json:"architect_name"`
}

// CreatePropertyUnitRequest for creating a new unit
type CreatePropertyUnitRequest struct {
	ProjectID             string  `json:"project_id" validate:"required"`
	BlockID               string  `json:"block_id" validate:"required"`
	UnitNumber            string  `json:"unit_number" validate:"required"`
	Floor                 int     `json:"floor"`
	UnitType              string  `json:"unit_type"`
	Facing                string  `json:"facing"`
	CarpetArea            float64 `json:"carpet_area"`
	CarpetAreaWithBalcony float64 `json:"carpet_area_with_balcony"`
	UtilityArea           float64 `json:"utility_area"`
	PlinthArea            float64 `json:"plinth_area"`
	SBUA                  float64 `json:"sbua"`
	UDSSqft               float64 `json:"uds_sqft"`
}

// CreateCustomerBookingRequest for creating a new booking
type CreateCustomerBookingRequest struct {
	UnitID                  string    `json:"unit_id" validate:"required"`
	CustomerID              string    `json:"customer_id"`
	BookingDate             time.Time `json:"booking_date" validate:"required"`
	RatePerSqft             float64   `json:"rate_per_sqft"`
	CompositeGuidelineValue float64   `json:"composite_guideline_value"`
	CarParkingType          string    `json:"car_parking_type"`
	ParkingLocation         string    `json:"parking_location"`
}

// CreateBookingPaymentRequest for recording a payment
type CreateBookingPaymentRequest struct {
	BookingID     string    `json:"booking_id" validate:"required"`
	PaymentDate   time.Time `json:"payment_date" validate:"required"`
	PaymentMode   string    `json:"payment_mode" validate:"required"`
	PaidBy        string    `json:"paid_by"`
	ReceiptNumber string    `json:"receipt_number" validate:"required"`
	Towards       string    `json:"towards"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	BankName      string    `json:"bank_name"`
	TransactionID string    `json:"transaction_id"`
	Remarks       string    `json:"remarks"`
}

// PropertyMilestoneRequest for tracking milestones
type PropertyMilestoneRequest struct {
	BookingID         string     `json:"booking_id" validate:"required"`
	CampaignName      string     `json:"campaign_name"`
	Source            string     `json:"source"`
	SubSource         string     `json:"subsource"`
	LeadGeneratedDate *time.Time `json:"lead_generated_date"`
	ReEngagedDate     *time.Time `json:"re_engaged_date"`
	SiteVisitDate     *time.Time `json:"site_visit_date"`
	ReVisitDate       *time.Time `json:"revisit_date"`
	BookingDate       *time.Time `json:"booking_date"`
	CancelledDate     *time.Time `json:"cancelled_date"`
	Notes             string     `json:"notes"`
}
