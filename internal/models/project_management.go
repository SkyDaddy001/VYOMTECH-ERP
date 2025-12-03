package models

import (
	"database/sql"
	"time"
)

// ============================================================
// MODELS FOR EXTENDED PROPERTY PROJECT SYSTEM
// Extends migration 008 (real_estate.sql) with enhanced customer management,
// detailed cost sheets, and construction milestone tracking
// ============================================================

// PropertyCustomerProfile extends customer data with comprehensive KYC support
// Covers primary applicant (15 fields), up to 3 co-applicants (30 fields), booking dates (6),
// banking/loan (7), sales tracking (5), maintenance charges (4), compliance docs (2), plus metadata (7)
type PropertyCustomerProfile struct {
	// Primary Keys & Identification
	ID           string `gorm:"primaryKey" json:"id"`
	TenantID     string `json:"tenant_id"`
	CustomerCode string `json:"customer_code"`
	UnitID       string `json:"unit_id"`

	// PRIMARY APPLICANT - Personal Information
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	PhonePrimary   string `json:"phone_primary"`
	PhoneSecondary string `json:"phone_secondary"`
	AlternatePhone string `json:"alternate_phone"`
	CompanyName    string `json:"company_name"`
	Designation    string `json:"designation"`
	PANNumber      string `json:"pan_number"`
	AadharNumber   string `json:"aadhar_number"`
	PANCopyURL     string `json:"pan_copy_url"`
	AadharCopyURL  string `json:"aadhar_copy_url"`
	POADocumentNo  string `json:"poa_document_no"`
	CareOf         string `json:"care_of"`

	// PRIMARY APPLICANT - Communication Address
	CommunicationAddressLine1 string `json:"communication_address_line1"`
	CommunicationAddressLine2 string `json:"communication_address_line2"`
	CommunicationCity         string `json:"communication_city"`
	CommunicationState        string `json:"communication_state"`
	CommunicationCountry      string `json:"communication_country"`
	CommunicationZip          string `json:"communication_zip"`

	// PRIMARY APPLICANT - Permanent Address
	PermanentAddressLine1 string `json:"permanent_address_line1"`
	PermanentAddressLine2 string `json:"permanent_address_line2"`
	PermanentCity         string `json:"permanent_city"`
	PermanentState        string `json:"permanent_state"`
	PermanentCountry      string `json:"permanent_country"`
	PermanentZip          string `json:"permanent_zip"`

	// PRIMARY APPLICANT - Employment & Financial
	Profession     string  `json:"profession"`
	EmployerName   string  `json:"employer_name"`
	EmploymentType string  `json:"employment_type"`
	MonthlyIncome  float64 `json:"monthly_income"`

	// CO-APPLICANT 1 - Full Details
	CoApplicant1Name                 string `json:"co_applicant_1_name"`
	CoApplicant1Number               string `json:"co_applicant_1_number"`
	CoApplicant1AlternateNumber      string `json:"co_applicant_1_alternate_number"`
	CoApplicant1Email                string `json:"co_applicant_1_email"`
	CoApplicant1CommunicationAddress string `json:"co_applicant_1_communication_address"`
	CoApplicant1PermanentAddress     string `json:"co_applicant_1_permanent_address"`
	CoApplicant1Aadhar               string `json:"co_applicant_1_aadhar"`
	CoApplicant1PAN                  string `json:"co_applicant_1_pan"`
	CoApplicant1CareOf               string `json:"co_applicant_1_care_of"`
	CoApplicant1Relation             string `json:"co_applicant_1_relation"`

	// CO-APPLICANT 2 - Full Details
	CoApplicant2Name                 string `json:"co_applicant_2_name"`
	CoApplicant2Number               string `json:"co_applicant_2_number"`
	CoApplicant2AlternateNumber      string `json:"co_applicant_2_alternate_number"`
	CoApplicant2Email                string `json:"co_applicant_2_email"`
	CoApplicant2CommunicationAddress string `json:"co_applicant_2_communication_address"`
	CoApplicant2PermanentAddress     string `json:"co_applicant_2_permanent_address"`
	CoApplicant2Aadhar               string `json:"co_applicant_2_aadhar"`
	CoApplicant2PAN                  string `json:"co_applicant_2_pan"`
	CoApplicant2CareOf               string `json:"co_applicant_2_care_of"`
	CoApplicant2Relation             string `json:"co_applicant_2_relation"`

	// CO-APPLICANT 3 - Full Details (Bonus support for HUF/Family structures)
	CoApplicant3Name                 string `json:"co_applicant_3_name"`
	CoApplicant3Number               string `json:"co_applicant_3_number"`
	CoApplicant3AlternateNumber      string `json:"co_applicant_3_alternate_number"`
	CoApplicant3Email                string `json:"co_applicant_3_email"`
	CoApplicant3CommunicationAddress string `json:"co_applicant_3_communication_address"`
	CoApplicant3PermanentAddress     string `json:"co_applicant_3_permanent_address"`
	CoApplicant3Aadhar               string `json:"co_applicant_3_aadhar"`
	CoApplicant3PAN                  string `json:"co_applicant_3_pan"`
	CoApplicant3CareOf               string `json:"co_applicant_3_care_of"`
	CoApplicant3Relation             string `json:"co_applicant_3_relation"`

	// BOOKING LIFECYCLE DATES
	BookingDate      *time.Time `json:"booking_date"`
	WelcomeDate      *time.Time `json:"welcome_date"`
	AllotmentDate    *time.Time `json:"allotment_date"`
	AgreementDate    *time.Time `json:"agreement_date"`
	RegistrationDate *time.Time `json:"registration_date"`
	HandoverDate     *time.Time `json:"handover_date"`
	NOCReceivedDate  *time.Time `json:"noc_received_date"`

	// PROPERTY & RATE DETAILS
	RatePerSqft             float64 `json:"rate_per_sqft"`
	CompositeGuidelineValue float64 `json:"composite_guideline_value"`
	CarParkingType          string  `json:"car_parking_type"`

	// FINANCING & LOAN DETAILS
	LoanRequired      bool       `json:"loan_required"`
	LoanAmount        float64    `json:"loan_amount"`
	LoanSanctionDate  *time.Time `json:"loan_sanction_date"`
	BankName          string     `json:"bank_name"`
	BankBranch        string     `json:"bank_branch"`
	BankContactPerson string     `json:"bank_contact_person"`
	BankContactNumber string     `json:"bank_contact_number"`

	// SALES & CRM TRACKING
	ConnectorCodeNumber string `json:"connector_code_number"`
	LeadID              string `json:"lead_id"`
	SalesExecutiveID    string `json:"sales_executive_id"`
	SalesExecutiveName  string `json:"sales_executive_name"`
	SalesHeadID         string `json:"sales_head_id"`
	SalesHeadName       string `json:"sales_head_name"`
	BookingSource       string `json:"booking_source"`

	// MAINTENANCE & CHARGES
	MaintenanceCharges float64 `json:"maintenance_charges"`
	OtherWorksCharges  float64 `json:"other_works_charges"`
	CorpusCharges      float64 `json:"corpus_charges"`
	EBDeposit          float64 `json:"eb_deposit"`

	// COMPLIANCE & DOCUMENTS
	LifeCertificate string `json:"life_certificate"`

	// STATUS & METADATA
	CustomerType   string     `json:"customer_type"`
	CustomerStatus string     `json:"customer_status"`
	Notes          string     `json:"notes"`
	CreatedBy      string     `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

// PropertyCustomerUnitLink links customers to units and bookings
type PropertyCustomerUnitLink struct {
	ID              string     `gorm:"primaryKey" json:"id"`
	TenantID        string     `json:"tenant_id"`
	CustomerID      string     `json:"customer_id"`
	UnitID          string     `json:"unit_id"`
	BookingID       string     `json:"booking_id"`
	BookingStatus   string     `json:"booking_status"`
	BookingDate     *time.Time `json:"booking_date"`
	AgreementDate   *time.Time `json:"agreement_date"`
	PossessionDate  *time.Time `json:"possession_date"`
	HandoverDate    *time.Time `json:"handover_date"`
	PrimaryCustomer bool       `json:"primary_customer"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PropertyUnitAreaStatement represents the complete Area Statement for a unit
// Includes all area measurements (carpet, plinth, SBUA, balcony, utility, etc.) with RERA compliance
// Supports multiple area types and UDS calculation
type PropertyUnitAreaStatement struct {
	ID        string `gorm:"primaryKey" json:"id"`
	TenantID  string `json:"tenant_id"`
	ProjectID string `json:"project_id"`
	BlockID   string `json:"block_id"`
	UnitID    string `json:"unit_id"`
	// Unit Identification
	AptNo    string `json:"apt_no"`    // Apartment Number
	Floor    string `json:"floor"`     // Floor Number
	UnitType string `json:"unit_type"` // 1BHK, 2BHK, 3BHK, etc.
	Facing   string `json:"facing"`    // NORTH, SOUTH, EAST, WEST, etc.
	// RERA Carpet Area
	RERACarPetAreaSqft float64 `json:"rera_carpet_area_sqft"` // RERA Carpet Area in Sq.Ft.
	RERACarPetAreaSqm  float64 `json:"rera_carpet_area_sqm"`  // RERA Carpet Area in Sq.Mtrs.
	// Carpet Area with Balcony & Utility
	CarPetAreaWithBalconySqft float64 `json:"carpet_area_with_balcony_sqft"` // With Balcony and Utility in Sq.Ft.
	CarPetAreaWithBalconySqm  float64 `json:"carpet_area_with_balcony_sqm"`  // With Balcony and Utility in Sq.Mtrs.
	// Plinth Area
	PlinthAreaSqft float64 `json:"plinth_area_sqft"` // Plinth Area in Sq.Ft.
	PlinthAreaSqm  float64 `json:"plinth_area_sqm"`  // Plinth Area in Sq.Mtrs.
	// Super Built-Up Area (SBUA)
	SBUASqft float64 `json:"sbua_sqft"` // SBUA in Sq.Ft.
	SBUASqm  float64 `json:"sbua_sqm"`  // SBUA in Sq.Mtrs.
	// Undivided Share (UDS)
	UDSPerSqft   float64 `json:"uds_per_sqft"`   // UDS per Sq.Ft.
	UDSTotalSqft float64 `json:"uds_total_sqft"` // Total UDS in Sq.Ft.
	// Additional Areas
	BalconyAreaSqft float64 `json:"balcony_area_sqft"`
	BalconyAreaSqm  float64 `json:"balcony_area_sqm"`
	UtilityAreaSqft float64 `json:"utility_area_sqft"`
	UtilityAreaSqm  float64 `json:"utility_area_sqm"`
	GardenAreaSqft  float64 `json:"garden_area_sqft"`
	GardenAreaSqm   float64 `json:"garden_area_sqm"`
	TerraceAreaSqft float64 `json:"terrace_area_sqft"`
	TerraceAreaSqm  float64 `json:"terrace_area_sqm"`
	ParkingAreaSqft float64 `json:"parking_area_sqft"`
	ParkingAreaSqm  float64 `json:"parking_area_sqm"`
	CommonAreaSqft  float64 `json:"common_area_sqft"`
	CommonAreaSqm   float64 `json:"common_area_sqm"`
	// Ownership & Allocation
	AlotedTo             string  `json:"alloted_to"`            // Primary allottee name
	KeyHolder            string  `json:"key_holder"`            // Key holder details
	PercentageAllocation float64 `json:"percentage_allocation"` // Percentage share of common area
	// NOC & Compliance
	NOCTaken       string     `json:"noc_taken"` // YES, NO, PENDING, NA
	NOCDate        *time.Time `json:"noc_date"`
	NOCDocumentURL string     `json:"noc_document_url"`
	// Other Details
	AreaType    string    `json:"area_type"` // CARPET_AREA, BUILTUP_AREA, SUPER_AREA, etc.
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectCostConfiguration manages project-wise charge definitions
// Allows each project to define custom OTHER_CHARGES (e.g., CMWSSB, Water Tax, etc.)
// Supports both per-sqft and lump-sum charge calculation
type ProjectCostConfiguration struct {
	ID                    string    `gorm:"primaryKey" json:"id"`
	TenantID              string    `json:"tenant_id"`
	ProjectID             string    `json:"project_id"`
	ConfigName            string    `json:"config_name"`   // e.g., CMWSSB, Water Tax, Electricity Deposit
	ConfigType            string    `json:"config_type"`   // OTHER_CHARGE_1, OTHER_CHARGE_2, etc.
	ChargeType            string    `json:"charge_type"`   // PER_SQFT or LUMPSUM
	ChargeAmount          float64   `json:"charge_amount"` // Amount per sqft or lump-sum amount
	DisplayOrder          int       `json:"display_order"`
	IsMandatory           bool      `json:"is_mandatory"`             // Is this charge mandatory for all units
	ApplicableForUnitType string    `json:"applicable_for_unit_type"` // Comma-separated unit types or null for all
	Description           string    `json:"description"`
	Active                bool      `json:"active"`
	CreatedBy             string    `json:"created_by"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// PropertyPaymentReceipt tracks detailed payment transactions
type PropertyPaymentReceipt struct {
	ID            string `gorm:"primaryKey" json:"id"`
	TenantID      string `json:"tenant_id"`
	CustomerID    string `json:"customer_id"`
	CustomerName  string `json:"customer_name"`
	UnitID        string `json:"unit_id"`
	InstallmentID string `json:"installment_id"`

	// Receipt Details
	ReceiptNumber string     `json:"receipt_number"`
	ReceiptDate   *time.Time `json:"receipt_date"`
	PaymentDate   *time.Time `json:"payment_date"`

	// Payment Information
	PaymentMode          string  `json:"payment_mode"` // Cheque, NEFT, RTGS, Online, DD, Cash
	PaymentAmount        float64 `json:"payment_amount"`
	InstallmentAmountDue float64 `json:"installment_amount_due"`
	ShortfallAmount      float64 `json:"shortfall_amount"`
	ExcessAmount         float64 `json:"excess_amount"`
	PaymentStatus        string  `json:"payment_status"` // Pending, Received, Processed, Cleared, Bounced, Cancelled

	// Payment Mode Details (For Cheque)
	BankName     string     `json:"bank_name"`
	ChequeNumber string     `json:"cheque_number"`
	ChequeDate   *time.Time `json:"cheque_date"`

	// Payment Mode Details (For NEFT/RTGS)
	TransactionID string `json:"transaction_id"`
	AccountNumber string `json:"account_number"`

	// Receipt Details
	TowardsDescription    string `json:"towards_description"`      // APARTMENT_COST, MAINTENANCE, CORPUS, etc.
	ReceivedInBankAccount string `json:"received_in_bank_account"` // Bank account where payment received

	// Accounting
	GLAccountID string `json:"gl_account_id"`
	PaidBy      string `json:"paid_by"` // Customer, Agent, other

	// Metadata
	Remarks   string       `json:"remarks"`
	CreatedBy string       `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// PropertyProjectMilestone tracks construction progress milestones
type PropertyProjectMilestone struct {
	ID                    string     `gorm:"primaryKey" json:"id"`
	TenantID              string     `json:"tenant_id"`
	ProjectID             string     `json:"project_id"`
	BlockID               string     `json:"block_id"`
	MilestoneName         string     `json:"milestone_name"`
	MilestoneType         string     `json:"milestone_type"`
	MilestoneDescription  string     `json:"milestone_description"`
	StartDate             *time.Time `json:"start_date"`
	PlannedCompletionDate *time.Time `json:"planned_completion_date"`
	ActualCompletionDate  *time.Time `json:"actual_completion_date"`
	CompletionStatus      string     `json:"completion_status"`
	PercentageCompletion  int        `json:"percentage_completion"`
	ResponsiblePartyID    string     `json:"responsible_party_id"`
	ResponsiblePartyType  string     `json:"responsible_party_type"`
	BudgetAllocated       float64    `json:"budget_allocated"`
	BudgetSpent           float64    `json:"budget_spent"`
	BudgetVariance        float64    `json:"budget_variance"`
	Priority              int        `json:"priority"`
	Notes                 string     `json:"notes"`
	DocumentsURL          string     `json:"documents_url"`
	CreatedBy             string     `json:"created_by"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// PropertyProjectActivity logs daily project activities
type PropertyProjectActivity struct {
	ID                   string     `gorm:"primaryKey" json:"id"`
	TenantID             string     `json:"tenant_id"`
	ProjectID            string     `json:"project_id"`
	MilestoneID          string     `json:"milestone_id"`
	ActivityType         string     `json:"activity_type"`
	ActivityDescription  string     `json:"activity_description"`
	ActivityDate         *time.Time `json:"activity_date"`
	ActivityTime         *time.Time `json:"activity_time"`
	AssignedTo           string     `json:"assigned_to"`
	Status               string     `json:"status"`
	CompletionDate       *time.Time `json:"completion_date"`
	CompletionPercentage int        `json:"completion_percentage"`
	AttachmentsURL       string     `json:"attachments_url"`
	Notes                string     `json:"notes"`
	CreatedBy            string     `json:"created_by"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// PropertyProjectDocument manages project documentation
type PropertyProjectDocument struct {
	ID                  string     `gorm:"primaryKey" json:"id"`
	TenantID            string     `json:"tenant_id"`
	ProjectID           string     `json:"project_id"`
	DocumentType        string     `json:"document_type"`
	DocumentName        string     `json:"document_name"`
	DocumentDescription string     `json:"document_description"`
	DocumentURL         string     `json:"document_url"`
	FileSizeBytes       int        `json:"file_size_bytes"`
	FileType            string     `json:"file_type"`
	UploadedBy          string     `json:"uploaded_by"`
	UploadDate          *time.Time `json:"upload_date"`
	ExpiryDate          *time.Time `json:"expiry_date"`
	IsActive            bool       `json:"is_active"`
	ApprovalStatus      string     `json:"approval_status"`
	ApprovedBy          string     `json:"approved_by"`
	ApprovalDate        *time.Time `json:"approval_date"`
	VersionNumber       int        `json:"version_number"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// PropertyProjectSummary provides KPI dashboard data
type PropertyProjectSummary struct {
	ID                          string     `gorm:"primaryKey" json:"id"`
	TenantID                    string     `json:"tenant_id"`
	ProjectID                   string     `json:"project_id"`
	SummaryDate                 *time.Time `json:"summary_date"`
	TotalUnits                  int        `json:"total_units"`
	UnitsAvailable              int        `json:"units_available"`
	UnitsReserved               int        `json:"units_reserved"`
	UnitsSold                   int        `json:"units_sold"`
	UnitsHandedOver             int        `json:"units_handed_over"`
	TotalRevenueBooked          float64    `json:"total_revenue_booked"`
	TotalRevenueReceived        float64    `json:"total_revenue_received"`
	TotalConstructionCost       float64    `json:"total_construction_cost"`
	TotalCostIncurred           float64    `json:"total_cost_incurred"`
	GrossProfit                 float64    `json:"gross_profit"`
	MarginPercentage            float64    `json:"margin_percentage"`
	ProjectCompletionPercentage int        `json:"project_completion_percentage"`
	ExpectedCompletionDate      *time.Time `json:"expected_completion_date"`
	NumberOfActiveMilestones    int        `json:"number_of_active_milestones"`
	NumberOfDelayedMilestones   int        `json:"number_of_delayed_milestones"`
	CustomerSatisfactionScore   float32    `json:"customer_satisfaction_score"`
	LastUpdatedBy               string     `json:"last_updated_by"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

// ============================================================
// API REQUEST/RESPONSE TYPES
// ============================================================

type CreateCustomerProfileRequest struct {
	// Customer Identification
	CustomerCode string `json:"customer_code" binding:"required"`
	UnitID       string `json:"unit_id"`

	// PRIMARY APPLICANT - Personal Information
	FirstName      string `json:"first_name" binding:"required"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	PhonePrimary   string `json:"phone_primary"`
	PhoneSecondary string `json:"phone_secondary"`
	AlternatePhone string `json:"alternate_phone"`
	CompanyName    string `json:"company_name"`
	Designation    string `json:"designation"`
	PANNumber      string `json:"pan_number"`
	AadharNumber   string `json:"aadhar_number"`
	POADocumentNo  string `json:"poa_document_no"`
	CareOf         string `json:"care_of"`

	// PRIMARY APPLICANT - Communication Address
	CommunicationAddressLine1 string `json:"communication_address_line1"`
	CommunicationAddressLine2 string `json:"communication_address_line2"`
	CommunicationCity         string `json:"communication_city"`
	CommunicationState        string `json:"communication_state"`
	CommunicationCountry      string `json:"communication_country"`
	CommunicationZip          string `json:"communication_zip"`

	// PRIMARY APPLICANT - Permanent Address
	PermanentAddressLine1 string `json:"permanent_address_line1"`
	PermanentAddressLine2 string `json:"permanent_address_line2"`
	PermanentCity         string `json:"permanent_city"`
	PermanentState        string `json:"permanent_state"`
	PermanentCountry      string `json:"permanent_country"`
	PermanentZip          string `json:"permanent_zip"`

	// PRIMARY APPLICANT - Employment & Financial
	Profession     string  `json:"profession"`
	EmployerName   string  `json:"employer_name"`
	EmploymentType string  `json:"employment_type"`
	MonthlyIncome  float64 `json:"monthly_income"`

	// CO-APPLICANT 1
	CoApplicant1Name                 string `json:"co_applicant_1_name"`
	CoApplicant1Number               string `json:"co_applicant_1_number"`
	CoApplicant1AlternateNumber      string `json:"co_applicant_1_alternate_number"`
	CoApplicant1Email                string `json:"co_applicant_1_email"`
	CoApplicant1CommunicationAddress string `json:"co_applicant_1_communication_address"`
	CoApplicant1PermanentAddress     string `json:"co_applicant_1_permanent_address"`
	CoApplicant1Aadhar               string `json:"co_applicant_1_aadhar"`
	CoApplicant1PAN                  string `json:"co_applicant_1_pan"`
	CoApplicant1CareOf               string `json:"co_applicant_1_care_of"`
	CoApplicant1Relation             string `json:"co_applicant_1_relation"`

	// CO-APPLICANT 2
	CoApplicant2Name                 string `json:"co_applicant_2_name"`
	CoApplicant2Number               string `json:"co_applicant_2_number"`
	CoApplicant2AlternateNumber      string `json:"co_applicant_2_alternate_number"`
	CoApplicant2Email                string `json:"co_applicant_2_email"`
	CoApplicant2CommunicationAddress string `json:"co_applicant_2_communication_address"`
	CoApplicant2PermanentAddress     string `json:"co_applicant_2_permanent_address"`
	CoApplicant2Aadhar               string `json:"co_applicant_2_aadhar"`
	CoApplicant2PAN                  string `json:"co_applicant_2_pan"`
	CoApplicant2CareOf               string `json:"co_applicant_2_care_of"`
	CoApplicant2Relation             string `json:"co_applicant_2_relation"`

	// CO-APPLICANT 3
	CoApplicant3Name                 string `json:"co_applicant_3_name"`
	CoApplicant3Number               string `json:"co_applicant_3_number"`
	CoApplicant3AlternateNumber      string `json:"co_applicant_3_alternate_number"`
	CoApplicant3Email                string `json:"co_applicant_3_email"`
	CoApplicant3CommunicationAddress string `json:"co_applicant_3_communication_address"`
	CoApplicant3PermanentAddress     string `json:"co_applicant_3_permanent_address"`
	CoApplicant3Aadhar               string `json:"co_applicant_3_aadhar"`
	CoApplicant3PAN                  string `json:"co_applicant_3_pan"`
	CoApplicant3CareOf               string `json:"co_applicant_3_care_of"`
	CoApplicant3Relation             string `json:"co_applicant_3_relation"`

	// BOOKING DATES
	BookingDate      string `json:"booking_date"` // YYYY-MM-DD format
	WelcomeDate      string `json:"welcome_date"`
	AllotmentDate    string `json:"allotment_date"`
	AgreementDate    string `json:"agreement_date"`
	RegistrationDate string `json:"registration_date"`
	HandoverDate     string `json:"handover_date"`
	NOCReceivedDate  string `json:"noc_received_date"`

	// PROPERTY DETAILS
	RatePerSqft             float64 `json:"rate_per_sqft"`
	CompositeGuidelineValue float64 `json:"composite_guideline_value"`
	CarParkingType          string  `json:"car_parking_type"`

	// FINANCING
	LoanRequired      bool    `json:"loan_required"`
	LoanAmount        float64 `json:"loan_amount"`
	LoanSanctionDate  string  `json:"loan_sanction_date"` // YYYY-MM-DD format
	BankName          string  `json:"bank_name"`
	BankBranch        string  `json:"bank_branch"`
	BankContactPerson string  `json:"bank_contact_person"`
	BankContactNumber string  `json:"bank_contact_number"`

	// SALES & CRM
	ConnectorCodeNumber string `json:"connector_code_number"`
	LeadID              string `json:"lead_id"`
	SalesExecutiveID    string `json:"sales_executive_id"`
	SalesExecutiveName  string `json:"sales_executive_name"`
	SalesHeadID         string `json:"sales_head_id"`
	SalesHeadName       string `json:"sales_head_name"`
	BookingSource       string `json:"booking_source"`

	// CHARGES
	MaintenanceCharges float64 `json:"maintenance_charges"`
	OtherWorksCharges  float64 `json:"other_works_charges"`
	CorpusCharges      float64 `json:"corpus_charges"`
	EBDeposit          float64 `json:"eb_deposit"`

	// COMPLIANCE
	LifeCertificate string `json:"life_certificate"`

	// STATUS
	CustomerType   string `json:"customer_type"`
	CustomerStatus string `json:"customer_status"`
	Notes          string `json:"notes"`
}

type LinkCustomerToUnitRequest struct {
	CustomerID      string     `json:"customer_id" binding:"required"`
	UnitID          string     `json:"unit_id" binding:"required"`
	BookingID       string     `json:"booking_id"`
	BookingStatus   string     `json:"booking_status"`
	BookingDate     *time.Time `json:"booking_date"`
	AgreementDate   *time.Time `json:"agreement_date"`
	PrimaryCustomer bool       `json:"primary_customer"`
}

type CreatePaymentReceiptRequest struct {
	// Core References
	CustomerID    string `json:"customer_id" binding:"required"`
	CustomerName  string `json:"customer_name"`
	UnitID        string `json:"unit_id" binding:"required"`
	InstallmentID string `json:"installment_id"`

	// Receipt Details
	PaymentDate   *time.Time `json:"payment_date" binding:"required"`
	PaymentMode   string     `json:"payment_mode" binding:"required"` // Cheque/NEFT/RTGS/Online/DD/Cash
	PaymentAmount float64    `json:"payment_amount" binding:"required"`

	// Payment Mode Details
	BankName      string     `json:"bank_name"`
	ChequeNumber  string     `json:"cheque_number"`
	ChequeDate    *time.Time `json:"cheque_date"`
	TransactionID string     `json:"transaction_id"`
	AccountNumber string     `json:"account_number"`

	// Receipt Details
	TowardsDescription    string `json:"towards_description"`      // APARTMENT_COST, MAINTENANCE, CORPUS, etc.
	ReceivedInBankAccount string `json:"received_in_bank_account"` // Bank account where payment received
	PaidBy                string `json:"paid_by"`                  // Customer, Agent, etc.

	// Additional
	Remarks string `json:"remarks"`
}

type CreateMilestoneRequest struct {
	ProjectID             string     `json:"project_id" binding:"required"`
	BlockID               string     `json:"block_id"`
	MilestoneName         string     `json:"milestone_name" binding:"required"`
	MilestoneType         string     `json:"milestone_type"`
	MilestoneDescription  string     `json:"milestone_description"`
	StartDate             *time.Time `json:"start_date"`
	PlannedCompletionDate *time.Time `json:"planned_completion_date"`
	ResponsiblePartyID    string     `json:"responsible_party_id"`
	ResponsiblePartyType  string     `json:"responsible_party_type"`
	BudgetAllocated       float64    `json:"budget_allocated"`
	Priority              int        `json:"priority"`
}

type UpdateMilestoneRequest struct {
	MilestoneID          string     `json:"milestone_id" binding:"required"`
	PercentageCompletion int        `json:"percentage_completion"`
	CompletionStatus     string     `json:"completion_status"`
	ActualCompletionDate *time.Time `json:"actual_completion_date"`
	BudgetSpent          float64    `json:"budget_spent"`
	Notes                string     `json:"notes"`
}

type CreateProjectActivityRequest struct {
	ProjectID           string     `json:"project_id" binding:"required"`
	MilestoneID         string     `json:"milestone_id"`
	ActivityType        string     `json:"activity_type" binding:"required"`
	ActivityDescription string     `json:"activity_description" binding:"required"`
	ActivityDate        *time.Time `json:"activity_date"`
	ActivityTime        *time.Time `json:"activity_time"`
	AssignedTo          string     `json:"assigned_to"`
	Notes               string     `json:"notes"`
}

// CreateAreaStatementRequest for creating/updating area statement with complete area breakup
type CreateAreaStatementRequest struct {
	ProjectID string `json:"project_id" binding:"required"`
	BlockID   string `json:"block_id"`
	UnitID    string `json:"unit_id" binding:"required"`
	// Unit Identification
	AptNo    string `json:"apt_no"`    // Apartment Number
	Floor    string `json:"floor"`     // Floor Number
	UnitType string `json:"unit_type"` // 1BHK, 2BHK, 3BHK, etc.
	Facing   string `json:"facing"`    // NORTH, SOUTH, EAST, WEST, etc.
	// RERA Carpet Area
	RERACarPetAreaSqft float64 `json:"rera_carpet_area_sqft"` // RERA Carpet Area in Sq.Ft.
	RERACarPetAreaSqm  float64 `json:"rera_carpet_area_sqm"`  // RERA Carpet Area in Sq.Mtrs.
	// Carpet Area with Balcony & Utility
	CarPetAreaWithBalconySqft float64 `json:"carpet_area_with_balcony_sqft"` // With Balcony and Utility in Sq.Ft.
	CarPetAreaWithBalconySqm  float64 `json:"carpet_area_with_balcony_sqm"`  // With Balcony and Utility in Sq.Mtrs.
	// Plinth Area
	PlinthAreaSqft float64 `json:"plinth_area_sqft"` // Plinth Area in Sq.Ft.
	PlinthAreaSqm  float64 `json:"plinth_area_sqm"`  // Plinth Area in Sq.Mtrs.
	// Super Built-Up Area (SBUA)
	SBUASqft float64 `json:"sbua_sqft"` // SBUA in Sq.Ft.
	SBUASqm  float64 `json:"sbua_sqm"`  // SBUA in Sq.Mtrs.
	// Undivided Share (UDS)
	UDSPerSqft   float64 `json:"uds_per_sqft"`   // UDS per Sq.Ft.
	UDSTotalSqft float64 `json:"uds_total_sqft"` // Total UDS in Sq.Ft.
	// Additional Areas
	BalconyAreaSqft float64 `json:"balcony_area_sqft"`
	BalconyAreaSqm  float64 `json:"balcony_area_sqm"`
	UtilityAreaSqft float64 `json:"utility_area_sqft"`
	UtilityAreaSqm  float64 `json:"utility_area_sqm"`
	GardenAreaSqft  float64 `json:"garden_area_sqft"`
	GardenAreaSqm   float64 `json:"garden_area_sqm"`
	TerraceAreaSqft float64 `json:"terrace_area_sqft"`
	TerraceAreaSqm  float64 `json:"terrace_area_sqm"`
	ParkingAreaSqft float64 `json:"parking_area_sqft"`
	ParkingAreaSqm  float64 `json:"parking_area_sqm"`
	CommonAreaSqft  float64 `json:"common_area_sqft"`
	CommonAreaSqm   float64 `json:"common_area_sqm"`
	// Ownership & Allocation
	AlotedTo             string  `json:"alloted_to"`            // Primary allottee name
	KeyHolder            string  `json:"key_holder"`            // Key holder details
	PercentageAllocation float64 `json:"percentage_allocation"` // Percentage share of common area
	// NOC & Compliance
	NOCTaken       string `json:"noc_taken"` // YES, NO, PENDING, NA
	NOCDate        string `json:"noc_date"`
	NOCDocumentURL string `json:"noc_document_url"`
	// Other Details
	AreaType    string `json:"area_type"` // CARPET_AREA, BUILTUP_AREA, SUPER_AREA, etc.
	Description string `json:"description"`
}

type UpdateCostSheetRequest struct {
	UnitID                       string     `json:"unit_id" binding:"required"`
	BlockName                    string     `json:"block_name"`
	SBUA                         float64    `json:"sbua"` // Super Built-Up Area
	RatePerSqft                  float64    `json:"rate_per_sqft"`
	CarParkingCost               float64    `json:"car_parking_cost"`
	PLC                          float64    `json:"plc"`                       // PLC charges if any
	StatutoryApprovalCharge      float64    `json:"statutory_approval_charge"` // Local body, infrastructure
	LegalDocumentationCharge     float64    `json:"legal_documentation_charge"`
	AmenitiesEquipmentCharge     float64    `json:"amenities_equipment_charge"`
	OtherCharges1                float64    `json:"other_charges_1"`
	OtherCharges1Name            string     `json:"other_charges_1_name"` // e.g., CMWSSB
	OtherCharges1Type            string     `json:"other_charges_1_type"` // PER_SQFT or LUMPSUM
	OtherCharges2                float64    `json:"other_charges_2"`
	OtherCharges2Name            string     `json:"other_charges_2_name"`
	OtherCharges2Type            string     `json:"other_charges_2_type"`
	OtherCharges3                float64    `json:"other_charges_3"`
	OtherCharges3Name            string     `json:"other_charges_3_name"`
	OtherCharges3Type            string     `json:"other_charges_3_type"`
	OtherCharges4                float64    `json:"other_charges_4"`
	OtherCharges4Name            string     `json:"other_charges_4_name"`
	OtherCharges4Type            string     `json:"other_charges_4_type"`
	OtherCharges5                float64    `json:"other_charges_5"`
	OtherCharges5Name            string     `json:"other_charges_5_name"`
	OtherCharges5Type            string     `json:"other_charges_5_type"`
	ApartmentCostExcludingGovt   float64    `json:"apartment_cost_excluding_govt"`
	ActualSoldPriceExcludingGovt float64    `json:"actual_sold_price_excluding_govt"`
	GSTApplicable                bool       `json:"gst_applicable"`
	GSTPercentage                float64    `json:"gst_percentage"`
	ClubMembership               float64    `json:"club_membership"`
	RegistrationCharge           float64    `json:"registration_charge"`
	EffectiveDate                *time.Time `json:"effective_date"`
	ValidUntil                   *time.Time `json:"valid_until"`
}

// CreateProjectCostConfigRequest for setting up project-wise charge definitions
// ChargeType: PER_SQFT (multiply by SBUA) or LUMPSUM (fixed amount)
type CreateProjectCostConfigRequest struct {
	ProjectID             string  `json:"project_id" binding:"required"`
	ConfigName            string  `json:"config_name" binding:"required"`   // e.g., CMWSSB
	ConfigType            string  `json:"config_type" binding:"required"`   // OTHER_CHARGE_1, OTHER_CHARGE_2, etc.
	ChargeType            string  `json:"charge_type" binding:"required"`   // PER_SQFT or LUMPSUM
	ChargeAmount          float64 `json:"charge_amount" binding:"required"` // Amount per sqft or lump-sum
	DisplayOrder          int     `json:"display_order"`
	IsMandatory           bool    `json:"is_mandatory"`
	ApplicableForUnitType string  `json:"applicable_for_unit_type"` // Comma-separated: 1BHK,2BHK,3BHK or null for all
	Description           string  `json:"description"`
}

// PropertyBankFinancing tracks banker's sanction, disbursement, and collection details per unit
type PropertyBankFinancing struct {
	ID         string `gorm:"primaryKey" json:"id"`
	TenantID   string `json:"tenant_id"`
	ProjectID  string `json:"project_id"`
	BlockID    string `json:"block_id"`
	UnitID     string `json:"unit_id"`
	CustomerID string `json:"customer_id"`
	// Unit & Apartment Details
	AptNo         string  `json:"apt_no"`
	BlockName     string  `json:"block_name"`
	ApartmentCost float64 `json:"apartment_cost"` // Base apartment cost
	// Bank Details
	BankName          string `json:"bank_name"`
	BankerReferenceNo string `json:"banker_reference_no"`
	// Sanction Details
	SanctionedAmount float64    `json:"sanctioned_amount"` // Total sanctioned loan amount
	SanctionedDate   *time.Time `json:"sanctioned_date"`
	// Disbursement Tracking
	TotalDisbursedAmount  float64    `json:"total_disbursed_amount"` // Total disbursed till date
	DisbursementStatus    string     `json:"disbursement_status"`    // PENDING, PARTIAL, COMPLETED
	LastDisbursementDate  *time.Time `json:"last_disbursement_date"`
	RemainingDisbursement float64    `json:"remaining_disbursement"` // Amount yet to be disbursed
	// Collection from Unit Owner
	TotalCollectionFromUnit float64 `json:"total_collection_from_unit"` // Non-bank collections
	CollectionStatus        string  `json:"collection_status"`          // PENDING, PARTIAL, COMPLETED
	// Financial Summary
	TotalCommitment   float64 `json:"total_commitment"`   // Total project cost commitment
	OutstandingAmount float64 `json:"outstanding_amount"` // Amount still to be collected
	// Compliance & Documentation
	NOCRequired     bool       `json:"noc_required"`
	NOCReceived     bool       `json:"noc_received"`
	NOCDate         *time.Time `json:"noc_date"`
	DocumentsStatus string     `json:"documents_status"` // COMPLETE, PENDING, SUBMITTED_FOR_VERIFICATION
	DocumentsURL    string     `json:"documents_url"`
	// Metadata
	Remarks   string    `json:"remarks"`
	Active    bool      `json:"active"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PropertyDisbursementSchedule tracks expected vs actual disbursements linked to milestones
type PropertyDisbursementSchedule struct {
	ID          string `gorm:"primaryKey" json:"id"`
	TenantID    string `json:"tenant_id"`
	FinancingID string `json:"financing_id"`
	UnitID      string `json:"unit_id"`
	CustomerID  string `json:"customer_id"`
	// Disbursement Details
	DisbursementNo             int        `json:"disbursement_no"` // Sequential number (1st, 2nd, 3rd, etc.)
	ExpectedDisbursementDate   *time.Time `json:"expected_disbursement_date"`
	ActualDisbursementDate     *time.Time `json:"actual_disbursement_date"`
	ExpectedDisbursementAmount float64    `json:"expected_disbursement_amount"`
	ActualDisbursementAmount   float64    `json:"actual_disbursement_amount"`
	DisbursementPercentage     float64    `json:"disbursement_percentage"` // % of sanctioned amount
	// Milestone Linkage
	LinkedMilestoneID string `json:"linked_milestone_id"`
	MilestoneStage    string `json:"milestone_stage"` // FOUNDATION, STRUCTURE, FINISHING, HANDOVER
	// Bank Documentation
	ChequeNo           string `json:"cheque_no"`
	BankReferenceNo    string `json:"bank_reference_no"`
	DisbursementStatus string `json:"disbursement_status"` // PENDING, CLEARED, FAILED, CANCELLED
	NEFTRefID          string `json:"neft_ref_id"`
	// Variance Tracking
	VarianceDays   int     `json:"variance_days"`   // Expected vs actual date
	VarianceAmount float64 `json:"variance_amount"` // Expected vs actual amount
	VarianceReason string  `json:"variance_reason"`
	// Metadata
	Remarks   string    `json:"remarks"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PropertyPaymentStage maps installment stages to percentage of cost and collections
type PropertyPaymentStage struct {
	ID         string `gorm:"primaryKey" json:"id"`
	TenantID   string `json:"tenant_id"`
	ProjectID  string `json:"project_id"`
	UnitID     string `json:"unit_id"`
	CustomerID string `json:"customer_id"`
	// Stage Definition
	StageName        string `json:"stage_name"` // BOOKING, FOUNDATION, STRUCTURE, FINISHING, HANDOVER
	StageNumber      int    `json:"stage_number"`
	StageDescription string `json:"stage_description"`
	// Percentage & Cost
	StagePercentage float64 `json:"stage_percentage"` // % of apartment cost due at this stage
	StageDueAmount  float64 `json:"stage_due_amount"` // Calculated from apartment_cost * percentage
	ApartmentCost   float64 `json:"apartment_cost"`
	// Collection Details
	AmountDue        float64 `json:"amount_due"`
	AmountReceived   float64 `json:"amount_received"`
	AmountPending    float64 `json:"amount_pending"`
	CollectionStatus string  `json:"collection_status"` // PENDING, PARTIAL, COMPLETED, OVERDUE
	// Timeline
	DueDate                *time.Time `json:"due_date"`
	ExpectedCollectionDate *time.Time `json:"expected_collection_date"`
	ActualCollectionDate   *time.Time `json:"actual_collection_date"`
	DaysOverdue            int        `json:"days_overdue"` // Days past due date
	// Payment Details
	PaymentReceivedDate *time.Time `json:"payment_received_date"`
	PaymentMode         string     `json:"payment_mode"` // CASH, CHEQUE, NEFT, RTGS, ONLINE, DD
	ReferenceNo         string     `json:"reference_no"` // Cheque no, transaction ID, etc.
	// Variance
	VarianceDays   int     `json:"variance_days"`   // Difference from due date
	VarianceAmount float64 `json:"variance_amount"` // Shortfall or excess
	// Metadata
	Remarks   string    `json:"remarks"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ============================================================
// REQUEST TYPES FOR BANK FINANCING OPERATIONS
// ============================================================

// CreateBankFinancingRequest for initiating bank financing tracking
type CreateBankFinancingRequest struct {
	ProjectID         string  `json:"project_id" binding:"required"`
	BlockID           string  `json:"block_id"`
	UnitID            string  `json:"unit_id" binding:"required"`
	CustomerID        string  `json:"customer_id" binding:"required"`
	AptNo             string  `json:"apt_no"`
	BlockName         string  `json:"block_name"`
	ApartmentCost     float64 `json:"apartment_cost" binding:"required"`
	BankName          string  `json:"bank_name" binding:"required"`
	BankerReferenceNo string  `json:"banker_reference_no"`
	SanctionedAmount  float64 `json:"sanctioned_amount" binding:"required"`
	SanctionedDate    string  `json:"sanctioned_date"`
	TotalCommitment   float64 `json:"total_commitment"`
}

// CreateDisbursementScheduleRequest for adding disbursement schedule
type CreateDisbursementScheduleRequest struct {
	FinancingID                string  `json:"financing_id" binding:"required"`
	UnitID                     string  `json:"unit_id" binding:"required"`
	CustomerID                 string  `json:"customer_id" binding:"required"`
	DisbursementNo             int     `json:"disbursement_no" binding:"required"`
	ExpectedDisbursementDate   string  `json:"expected_disbursement_date"`
	ExpectedDisbursementAmount float64 `json:"expected_disbursement_amount" binding:"required"`
	DisbursementPercentage     float64 `json:"disbursement_percentage"`
	LinkedMilestoneID          string  `json:"linked_milestone_id"`
	MilestoneStage             string  `json:"milestone_stage"`
}

// UpdateDisbursementRequest for recording actual disbursement
type UpdateDisbursementRequest struct {
	DisbursementID           string  `json:"disbursement_id" binding:"required"`
	ActualDisbursementDate   string  `json:"actual_disbursement_date"`
	ActualDisbursementAmount float64 `json:"actual_disbursement_amount"`
	DisbursementStatus       string  `json:"disbursement_status"`
	ChequeNo                 string  `json:"cheque_no"`
	BankReferenceNo          string  `json:"bank_reference_no"`
	NEFTRefID                string  `json:"neft_ref_id"`
}

// CreatePaymentStageRequest for setting up payment stages
type CreatePaymentStageRequest struct {
	ProjectID        string  `json:"project_id" binding:"required"`
	UnitID           string  `json:"unit_id" binding:"required"`
	CustomerID       string  `json:"customer_id" binding:"required"`
	StageName        string  `json:"stage_name" binding:"required"`
	StageNumber      int     `json:"stage_number" binding:"required"`
	StageDescription string  `json:"stage_description"`
	StagePercentage  float64 `json:"stage_percentage" binding:"required"`
	ApartmentCost    float64 `json:"apartment_cost" binding:"required"`
	DueDate          string  `json:"due_date"`
}

// UpdatePaymentStageRequest for recording collections
type UpdatePaymentStageRequest struct {
	PaymentStageID      string  `json:"payment_stage_id" binding:"required"`
	AmountReceived      float64 `json:"amount_received" binding:"required"`
	PaymentReceivedDate string  `json:"payment_received_date"`
	PaymentMode         string  `json:"payment_mode"`
	ReferenceNo         string  `json:"reference_no"`
	CollectionStatus    string  `json:"collection_status"`
}

// JSON marshaling helper - removed custom methods on json.RawMessage
// Use standard JSON marshaling instead
