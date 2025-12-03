package services

import (
	"database/sql"
	"fmt"
	"time"
	"vyomtech-backend/internal/models"

	"github.com/google/uuid"
)

// ProjectManagementService provides project management functionality
type ProjectManagementService struct {
	DB *sql.DB
}

// NewProjectManagementService creates a new project management service instance
func NewProjectManagementService(db *sql.DB) *ProjectManagementService {
	return &ProjectManagementService{
		DB: db,
	}
}

// ============================================================
// PROPERTY CUSTOMER PROFILE OPERATIONS
// ============================================================

// CreateCustomerProfile creates a new customer profile
func (s *ProjectManagementService) CreateCustomerProfile(tenantID string, req *models.PropertyCustomerProfile) (*models.PropertyCustomerProfile, error) {
	req.ID = uuid.New().String()
	req.TenantID = tenantID
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	query := `INSERT INTO property_customer_profile 
		(id, tenant_id, customer_code, unit_id, first_name, middle_name, last_name, email, 
		 phone_primary, phone_secondary, alternate_phone, company_name, designation, 
		 pan_number, aadhar_number, pan_copy_url, aadhar_copy_url, poa_document_no, care_of,
		 communication_address_line1, communication_address_line2, communication_city, 
		 communication_state, communication_country, communication_zip,
		 permanent_address_line1, permanent_address_line2, permanent_city, 
		 permanent_state, permanent_country, permanent_zip,
		 profession, employer_name, employment_type, monthly_income, customer_type,
		 co_applicant_1_name, co_applicant_1_number, co_applicant_1_alternate_number, co_applicant_1_email,
		 co_applicant_1_communication_address, co_applicant_1_permanent_address, co_applicant_1_aadhar, co_applicant_1_pan,
		 co_applicant_1_care_of, co_applicant_1_relation,
		 co_applicant_2_name, co_applicant_2_number, co_applicant_2_alternate_number, co_applicant_2_email,
		 co_applicant_2_communication_address, co_applicant_2_permanent_address, co_applicant_2_aadhar, co_applicant_2_pan,
		 co_applicant_2_care_of, co_applicant_2_relation,
		 co_applicant_3_name, co_applicant_3_number, co_applicant_3_alternate_number, co_applicant_3_email,
		 co_applicant_3_communication_address, co_applicant_3_permanent_address, co_applicant_3_aadhar, co_applicant_3_pan,
		 co_applicant_3_care_of, co_applicant_3_relation,
		 loan_required, loan_amount, loan_sanction_date, bank_name, bank_branch, bank_contact_person, bank_contact_number,
		 connector_code_number, lead_id, sales_executive_id, sales_executive_name, sales_head_id, sales_head_name, 
		 booking_source, life_certificate, customer_status, booking_date, welcome_date, allotment_date, agreement_date,
		 registration_date, handover_date, noc_received_date, rate_per_sqft, composite_guideline_value, car_parking_type,
		 maintenance_charges, other_works_charges, corpus_charges, eb_deposit, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
		 $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38,
		 $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56, $57,
		 $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71, $72, $73, $74, $75, $76,
		 $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88, $89)
		RETURNING id, created_at, updated_at`

	err := s.DB.QueryRow(query,
		req.ID, req.TenantID, req.CustomerCode, req.UnitID, req.FirstName, req.MiddleName, req.LastName,
		req.Email, req.PhonePrimary, req.PhoneSecondary, req.AlternatePhone, req.CompanyName,
		req.Designation, req.PANNumber, req.AadharNumber, req.PANCopyURL, req.AadharCopyURL,
		req.POADocumentNo, req.CareOf, req.CommunicationAddressLine1, req.CommunicationAddressLine2,
		req.CommunicationCity, req.CommunicationState, req.CommunicationCountry, req.CommunicationZip,
		req.PermanentAddressLine1, req.PermanentAddressLine2, req.PermanentCity,
		req.PermanentState, req.PermanentCountry, req.PermanentZip,
		req.Profession, req.EmployerName, req.EmploymentType, req.MonthlyIncome, req.CustomerType,
		req.CoApplicant1Name, req.CoApplicant1Number, req.CoApplicant1AlternateNumber, req.CoApplicant1Email,
		req.CoApplicant1CommunicationAddress, req.CoApplicant1PermanentAddress, req.CoApplicant1Aadhar, req.CoApplicant1PAN,
		req.CoApplicant1CareOf, req.CoApplicant1Relation,
		req.CoApplicant2Name, req.CoApplicant2Number, req.CoApplicant2AlternateNumber, req.CoApplicant2Email,
		req.CoApplicant2CommunicationAddress, req.CoApplicant2PermanentAddress, req.CoApplicant2Aadhar, req.CoApplicant2PAN,
		req.CoApplicant2CareOf, req.CoApplicant2Relation,
		req.CoApplicant3Name, req.CoApplicant3Number, req.CoApplicant3AlternateNumber, req.CoApplicant3Email,
		req.CoApplicant3CommunicationAddress, req.CoApplicant3PermanentAddress, req.CoApplicant3Aadhar, req.CoApplicant3PAN,
		req.CoApplicant3CareOf, req.CoApplicant3Relation,
		req.LoanRequired, req.LoanAmount, req.LoanSanctionDate, req.BankName, req.BankBranch, req.BankContactPerson, req.BankContactNumber,
		req.ConnectorCodeNumber, req.LeadID, req.SalesExecutiveID, req.SalesExecutiveName, req.SalesHeadID, req.SalesHeadName,
		req.BookingSource, req.LifeCertificate, req.CustomerStatus, req.BookingDate, req.WelcomeDate, req.AllotmentDate, req.AgreementDate,
		req.RegistrationDate, req.HandoverDate, req.NOCReceivedDate, req.RatePerSqft, req.CompositeGuidelineValue, req.CarParkingType,
		req.MaintenanceCharges, req.OtherWorksCharges, req.CorpusCharges, req.EBDeposit, req.Notes, req.CreatedBy, req.CreatedAt, req.UpdatedAt,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer profile: %w", err)
	}

	return req, nil
}

// GetCustomerProfile retrieves a customer profile by ID
func (s *ProjectManagementService) GetCustomerProfile(tenantID, customerID string) (*models.PropertyCustomerProfile, error) {
	query := `SELECT id, tenant_id, customer_code, unit_id, first_name, middle_name, last_name, 
		email, phone_primary, phone_secondary, alternate_phone, company_name, designation,
		pan_number, aadhar_number, pan_copy_url, aadhar_copy_url, poa_document_no, care_of,
		communication_address_line1, communication_address_line2, communication_city, communication_state, communication_country, communication_zip,
		permanent_address_line1, permanent_address_line2, permanent_city, permanent_state, permanent_country, permanent_zip,
		profession, employer_name, employment_type, monthly_income, customer_type,
		co_applicant_1_name, co_applicant_1_number, co_applicant_1_alternate_number, co_applicant_1_email,
		co_applicant_1_communication_address, co_applicant_1_permanent_address, co_applicant_1_aadhar, co_applicant_1_pan, co_applicant_1_care_of, co_applicant_1_relation,
		co_applicant_2_name, co_applicant_2_number, co_applicant_2_alternate_number, co_applicant_2_email,
		co_applicant_2_communication_address, co_applicant_2_permanent_address, co_applicant_2_aadhar, co_applicant_2_pan, co_applicant_2_care_of, co_applicant_2_relation,
		co_applicant_3_name, co_applicant_3_number, co_applicant_3_alternate_number, co_applicant_3_email,
		co_applicant_3_communication_address, co_applicant_3_permanent_address, co_applicant_3_aadhar, co_applicant_3_pan, co_applicant_3_care_of, co_applicant_3_relation,
		loan_required, loan_amount, loan_sanction_date, bank_name, bank_branch, bank_contact_person, bank_contact_number,
		connector_code_number, lead_id, sales_executive_id, sales_executive_name, sales_head_id, sales_head_name, booking_source, life_certificate,
		customer_status, booking_date, welcome_date, allotment_date, agreement_date, registration_date, handover_date, noc_received_date,
		rate_per_sqft, composite_guideline_value, car_parking_type, maintenance_charges, other_works_charges, corpus_charges, eb_deposit,
		notes, created_by, created_at, updated_at, deleted_at
		FROM property_customer_profile 
		WHERE id = $1 AND tenant_id = $2`

	customer := &models.PropertyCustomerProfile{}
	err := s.DB.QueryRow(query, customerID, tenantID).Scan(
		&customer.ID, &customer.TenantID, &customer.CustomerCode, &customer.UnitID,
		&customer.FirstName, &customer.MiddleName, &customer.LastName, &customer.Email,
		&customer.PhonePrimary, &customer.PhoneSecondary, &customer.AlternatePhone,
		&customer.CompanyName, &customer.Designation, &customer.PANNumber, &customer.AadharNumber,
		&customer.PANCopyURL, &customer.AadharCopyURL, &customer.POADocumentNo, &customer.CareOf,
		&customer.CommunicationAddressLine1, &customer.CommunicationAddressLine2, &customer.CommunicationCity, &customer.CommunicationState, &customer.CommunicationCountry, &customer.CommunicationZip,
		&customer.PermanentAddressLine1, &customer.PermanentAddressLine2, &customer.PermanentCity, &customer.PermanentState, &customer.PermanentCountry, &customer.PermanentZip,
		&customer.Profession, &customer.EmployerName, &customer.EmploymentType, &customer.MonthlyIncome, &customer.CustomerType,
		&customer.CoApplicant1Name, &customer.CoApplicant1Number, &customer.CoApplicant1AlternateNumber, &customer.CoApplicant1Email,
		&customer.CoApplicant1CommunicationAddress, &customer.CoApplicant1PermanentAddress, &customer.CoApplicant1Aadhar, &customer.CoApplicant1PAN, &customer.CoApplicant1CareOf, &customer.CoApplicant1Relation,
		&customer.CoApplicant2Name, &customer.CoApplicant2Number, &customer.CoApplicant2AlternateNumber, &customer.CoApplicant2Email,
		&customer.CoApplicant2CommunicationAddress, &customer.CoApplicant2PermanentAddress, &customer.CoApplicant2Aadhar, &customer.CoApplicant2PAN, &customer.CoApplicant2CareOf, &customer.CoApplicant2Relation,
		&customer.CoApplicant3Name, &customer.CoApplicant3Number, &customer.CoApplicant3AlternateNumber, &customer.CoApplicant3Email,
		&customer.CoApplicant3CommunicationAddress, &customer.CoApplicant3PermanentAddress, &customer.CoApplicant3Aadhar, &customer.CoApplicant3PAN, &customer.CoApplicant3CareOf, &customer.CoApplicant3Relation,
		&customer.LoanRequired, &customer.LoanAmount, &customer.LoanSanctionDate, &customer.BankName, &customer.BankBranch, &customer.BankContactPerson, &customer.BankContactNumber,
		&customer.ConnectorCodeNumber, &customer.LeadID, &customer.SalesExecutiveID, &customer.SalesExecutiveName, &customer.SalesHeadID, &customer.SalesHeadName, &customer.BookingSource, &customer.LifeCertificate,
		&customer.CustomerStatus, &customer.BookingDate, &customer.WelcomeDate, &customer.AllotmentDate, &customer.AgreementDate, &customer.RegistrationDate, &customer.HandoverDate, &customer.NOCReceivedDate,
		&customer.RatePerSqft, &customer.CompositeGuidelineValue, &customer.CarParkingType, &customer.MaintenanceCharges, &customer.OtherWorksCharges, &customer.CorpusCharges, &customer.EBDeposit,
		&customer.Notes, &customer.CreatedBy, &customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get customer profile: %w", err)
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}

	return customer, nil
}

// ============================================================
// PROPERTY UNIT AREA STATEMENT OPERATIONS
// ============================================================

// CreateAreaStatement creates a new area statement
func (s *ProjectManagementService) CreateAreaStatement(tenantID string, req *models.CreateAreaStatementRequest) (*models.PropertyUnitAreaStatement, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `INSERT INTO property_unit_area_statement 
		(id, tenant_id, project_id, block_id, unit_id, apt_no, floor, unit_type, facing,
		 rera_carpet_area_sqft, rera_carpet_area_sqm, carpet_area_with_balcony_sqft, carpet_area_with_balcony_sqm,
		 plinth_area_sqft, plinth_area_sqm, sbua_sqft, sbua_sqm, uds_per_sqft, uds_total_sqft,
		 balcony_area_sqft, balcony_area_sqm, utility_area_sqft, utility_area_sqm,
		 garden_area_sqft, garden_area_sqm, terrace_area_sqft, terrace_area_sqm,
		 parking_area_sqft, parking_area_sqm, common_area_sqft, common_area_sqm,
		 alloted_to, key_holder, percentage_allocation, noc_taken, noc_date, noc_document_url,
		 area_type, description, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
		 $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42)
		RETURNING id, created_at, updated_at`

	nocDate := (*time.Time)(nil)
	if req.NOCDate != "" {
		t, err := time.Parse("2006-01-02", req.NOCDate)
		if err == nil {
			nocDate = &t
		}
	}

	result := &models.PropertyUnitAreaStatement{
		ID:                        id,
		TenantID:                  tenantID,
		ProjectID:                 req.ProjectID,
		BlockID:                   req.BlockID,
		UnitID:                    req.UnitID,
		AptNo:                     req.AptNo,
		Floor:                     req.Floor,
		UnitType:                  req.UnitType,
		Facing:                    req.Facing,
		RERACarPetAreaSqft:        req.RERACarPetAreaSqft,
		RERACarPetAreaSqm:         req.RERACarPetAreaSqm,
		CarPetAreaWithBalconySqft: req.CarPetAreaWithBalconySqft,
		CarPetAreaWithBalconySqm:  req.CarPetAreaWithBalconySqm,
		PlinthAreaSqft:            req.PlinthAreaSqft,
		PlinthAreaSqm:             req.PlinthAreaSqm,
		SBUASqft:                  req.SBUASqft,
		SBUASqm:                   req.SBUASqm,
		UDSPerSqft:                req.UDSPerSqft,
		UDSTotalSqft:              req.UDSTotalSqft,
		BalconyAreaSqft:           req.BalconyAreaSqft,
		BalconyAreaSqm:            req.BalconyAreaSqm,
		UtilityAreaSqft:           req.UtilityAreaSqft,
		UtilityAreaSqm:            req.UtilityAreaSqm,
		GardenAreaSqft:            req.GardenAreaSqft,
		GardenAreaSqm:             req.GardenAreaSqm,
		TerraceAreaSqft:           req.TerraceAreaSqft,
		TerraceAreaSqm:            req.TerraceAreaSqm,
		ParkingAreaSqft:           req.ParkingAreaSqft,
		ParkingAreaSqm:            req.ParkingAreaSqm,
		CommonAreaSqft:            req.CommonAreaSqft,
		CommonAreaSqm:             req.CommonAreaSqm,
		AlotedTo:                  req.AlotedTo,
		KeyHolder:                 req.KeyHolder,
		PercentageAllocation:      req.PercentageAllocation,
		NOCTaken:                  req.NOCTaken,
		NOCDate:                   nocDate,
		NOCDocumentURL:            req.NOCDocumentURL,
		AreaType:                  req.AreaType,
		Description:               req.Description,
		Active:                    true,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	err := s.DB.QueryRow(query,
		id, tenantID, req.ProjectID, req.BlockID, req.UnitID, req.AptNo, req.Floor, req.UnitType, req.Facing,
		req.RERACarPetAreaSqft, req.RERACarPetAreaSqm, req.CarPetAreaWithBalconySqft, req.CarPetAreaWithBalconySqm,
		req.PlinthAreaSqft, req.PlinthAreaSqm, req.SBUASqft, req.SBUASqm, req.UDSPerSqft, req.UDSTotalSqft,
		req.BalconyAreaSqft, req.BalconyAreaSqm, req.UtilityAreaSqft, req.UtilityAreaSqm,
		req.GardenAreaSqft, req.GardenAreaSqm, req.TerraceAreaSqft, req.TerraceAreaSqm,
		req.ParkingAreaSqft, req.ParkingAreaSqm, req.CommonAreaSqft, req.CommonAreaSqm,
		req.AlotedTo, req.KeyHolder, req.PercentageAllocation, req.NOCTaken, nocDate, req.NOCDocumentURL,
		req.AreaType, req.Description, true, now, now,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create area statement: %w", err)
	}

	return result, nil
}

// ============================================================
// COST SHEET OPERATIONS
// ============================================================

// UpdateCostSheet updates or creates a cost sheet for a unit
func (s *ProjectManagementService) UpdateCostSheet(tenantID string, req *models.UpdateCostSheetRequest) error {
	now := time.Now()

	query := `INSERT INTO unit_cost_sheet 
		(tenant_id, unit_id, block_name, sbua, rate_per_sqft, car_parking_cost, plc,
		 statutory_approval_charge, legal_documentation_charge, amenities_equipment_charge,
		 other_charges_1, other_charges_1_name, other_charges_1_type,
		 other_charges_2, other_charges_2_name, other_charges_2_type,
		 other_charges_3, other_charges_3_name, other_charges_3_type,
		 other_charges_4, other_charges_4_name, other_charges_4_type,
		 other_charges_5, other_charges_5_name, other_charges_5_type,
		 apartment_cost_excluding_govt, actual_sold_price_excluding_govt,
		 gst_applicable, gst_percentage, gst_amount, grand_total, club_membership,
		 registration_charge, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		 $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
		 $33, $34, $35)
		ON DUPLICATE KEY UPDATE
		sbua = $4, rate_per_sqft = $5, car_parking_cost = $6, plc = $7,
		other_charges_1 = $11, other_charges_1_name = $12, other_charges_1_type = $13,
		other_charges_2 = $14, other_charges_2_name = $15, other_charges_2_type = $16,
		other_charges_3 = $17, other_charges_3_name = $18, other_charges_3_type = $19,
		other_charges_4 = $20, other_charges_4_name = $21, other_charges_4_type = $22,
		other_charges_5 = $23, other_charges_5_name = $24, other_charges_5_type = $25,
		gst_applicable = $28, gst_percentage = $29, gst_amount = $30, grand_total = $31,
		updated_at = $35`

	_, err := s.DB.Exec(query,
		tenantID, req.UnitID, req.BlockName, req.SBUA, req.RatePerSqft, req.CarParkingCost, req.PLC,
		req.StatutoryApprovalCharge, req.LegalDocumentationCharge, req.AmenitiesEquipmentCharge,
		req.OtherCharges1, req.OtherCharges1Name, req.OtherCharges1Type,
		req.OtherCharges2, req.OtherCharges2Name, req.OtherCharges2Type,
		req.OtherCharges3, req.OtherCharges3Name, req.OtherCharges3Type,
		req.OtherCharges4, req.OtherCharges4Name, req.OtherCharges4Type,
		req.OtherCharges5, req.OtherCharges5Name, req.OtherCharges5Type,
		req.ApartmentCostExcludingGovt, req.ActualSoldPriceExcludingGovt,
		req.GSTApplicable, req.GSTPercentage,
		req.ClubMembership, req.RegistrationCharge, now, now,
	)

	if err != nil {
		return fmt.Errorf("failed to update cost sheet: %w", err)
	}

	return nil
}

// ============================================================
// PROJECT COST CONFIGURATION OPERATIONS
// ============================================================

// CreateProjectCostConfiguration creates a project cost configuration
func (s *ProjectManagementService) CreateProjectCostConfiguration(tenantID string, req *models.CreateProjectCostConfigRequest) (*models.ProjectCostConfiguration, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `INSERT INTO project_cost_configuration
		(id, tenant_id, project_id, config_name, config_type, charge_type, charge_amount,
		 display_order, is_mandatory, applicable_for_unit_type, description, active,
		 created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	config := &models.ProjectCostConfiguration{
		ID:                    id,
		TenantID:              tenantID,
		ProjectID:             req.ProjectID,
		ConfigName:            req.ConfigName,
		ConfigType:            req.ConfigType,
		ChargeType:            req.ChargeType,
		ChargeAmount:          req.ChargeAmount,
		DisplayOrder:          req.DisplayOrder,
		IsMandatory:           req.IsMandatory,
		ApplicableForUnitType: req.ApplicableForUnitType,
		Description:           req.Description,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	err := s.DB.QueryRow(query,
		id, tenantID, req.ProjectID, req.ConfigName, req.ConfigType, req.ChargeType, req.ChargeAmount,
		req.DisplayOrder, req.IsMandatory, req.ApplicableForUnitType, req.Description, true, now, now,
	).Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create cost configuration: %w", err)
	}

	return config, nil
}

// ============================================================
// BANK FINANCING OPERATIONS
// ============================================================

// CreateBankFinancing creates a bank financing record
func (s *ProjectManagementService) CreateBankFinancing(tenantID string, req *models.CreateBankFinancingRequest) (*models.PropertyBankFinancing, error) {
	id := uuid.New().String()
	now := time.Now()

	sanctionedDate := (*time.Time)(nil)
	if req.SanctionedDate != "" {
		t, err := time.Parse("2006-01-02", req.SanctionedDate)
		if err == nil {
			sanctionedDate = &t
		}
	}

	query := `INSERT INTO property_bank_financing
		(id, tenant_id, project_id, block_id, unit_id, customer_id, apt_no, block_name,
		 apartment_cost, bank_name, banker_reference_no, sanctioned_amount, sanctioned_date,
		 total_disbursed_amount, disbursement_status, remaining_disbursement,
		 total_collection_from_unit, collection_status, total_commitment, outstanding_amount,
		 noc_required, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
		 $17, $18, $19, $20, $21, $22, $23, $24)
		RETURNING id, created_at, updated_at`

	financing := &models.PropertyBankFinancing{
		ID:                    id,
		TenantID:              tenantID,
		ProjectID:             req.ProjectID,
		BlockID:               req.BlockID,
		UnitID:                req.UnitID,
		CustomerID:            req.CustomerID,
		AptNo:                 req.AptNo,
		BlockName:             req.BlockName,
		ApartmentCost:         req.ApartmentCost,
		BankName:              req.BankName,
		BankerReferenceNo:     req.BankerReferenceNo,
		SanctionedAmount:      req.SanctionedAmount,
		SanctionedDate:        sanctionedDate,
		TotalDisbursedAmount:  0,
		DisbursementStatus:    "PENDING",
		RemainingDisbursement: req.SanctionedAmount,
		CollectionStatus:      "PENDING",
		TotalCommitment:       req.TotalCommitment,
		OutstandingAmount:     req.TotalCommitment,
		NOCRequired:           false,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	err := s.DB.QueryRow(query,
		id, tenantID, req.ProjectID, req.BlockID, req.UnitID, req.CustomerID,
		req.AptNo, req.BlockName, req.ApartmentCost, req.BankName, req.BankerReferenceNo,
		req.SanctionedAmount, sanctionedDate, 0, "PENDING", req.SanctionedAmount,
		0, "PENDING", req.TotalCommitment, req.TotalCommitment, false, true, now, now,
	).Scan(&financing.ID, &financing.CreatedAt, &financing.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create bank financing: %w", err)
	}

	return financing, nil
}

// ============================================================
// DISBURSEMENT SCHEDULE OPERATIONS
// ============================================================

// CreateDisbursementSchedule creates a disbursement schedule
func (s *ProjectManagementService) CreateDisbursementSchedule(tenantID string, req *models.CreateDisbursementScheduleRequest) (*models.PropertyDisbursementSchedule, error) {
	id := uuid.New().String()
	now := time.Now()

	expectedDate := (*time.Time)(nil)
	if req.ExpectedDisbursementDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpectedDisbursementDate)
		if err == nil {
			expectedDate = &t
		}
	}

	query := `INSERT INTO property_disbursement_schedule
		(id, tenant_id, financing_id, unit_id, customer_id, disbursement_no,
		 expected_disbursement_date, expected_disbursement_amount, disbursement_percentage,
		 linked_milestone_id, milestone_stage, disbursement_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	schedule := &models.PropertyDisbursementSchedule{
		ID:                         id,
		TenantID:                   tenantID,
		FinancingID:                req.FinancingID,
		UnitID:                     req.UnitID,
		CustomerID:                 req.CustomerID,
		DisbursementNo:             req.DisbursementNo,
		ExpectedDisbursementDate:   expectedDate,
		ExpectedDisbursementAmount: req.ExpectedDisbursementAmount,
		DisbursementPercentage:     req.DisbursementPercentage,
		LinkedMilestoneID:          req.LinkedMilestoneID,
		MilestoneStage:             req.MilestoneStage,
		DisbursementStatus:         "PENDING",
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}

	err := s.DB.QueryRow(query,
		id, tenantID, req.FinancingID, req.UnitID, req.CustomerID, req.DisbursementNo,
		expectedDate, req.ExpectedDisbursementAmount, req.DisbursementPercentage,
		req.LinkedMilestoneID, req.MilestoneStage, "PENDING", now, now,
	).Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create disbursement schedule: %w", err)
	}

	return schedule, nil
}

// ============================================================
// PAYMENT STAGE OPERATIONS
// ============================================================

// CreatePaymentStage creates a payment stage
func (s *ProjectManagementService) CreatePaymentStage(tenantID string, req *models.CreatePaymentStageRequest) (*models.PropertyPaymentStage, error) {
	id := uuid.New().String()
	now := time.Now()

	dueDate := (*time.Time)(nil)
	if req.DueDate != "" {
		t, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			dueDate = &t
		}
	}

	// Calculate stage due amount
	stageDueAmount := req.ApartmentCost * (req.StagePercentage / 100)

	query := `INSERT INTO property_payment_stage
		(id, tenant_id, project_id, unit_id, customer_id, stage_name, stage_number,
		 stage_percentage, stage_due_amount, apartment_cost, amount_due, amount_pending,
		 collection_status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at`

	stage := &models.PropertyPaymentStage{
		ID:               id,
		TenantID:         tenantID,
		ProjectID:        req.ProjectID,
		UnitID:           req.UnitID,
		CustomerID:       req.CustomerID,
		StageName:        req.StageName,
		StageNumber:      req.StageNumber,
		StagePercentage:  req.StagePercentage,
		StageDueAmount:   stageDueAmount,
		ApartmentCost:    req.ApartmentCost,
		AmountDue:        stageDueAmount,
		AmountPending:    stageDueAmount,
		CollectionStatus: "PENDING",
		DueDate:          dueDate,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err := s.DB.QueryRow(query,
		id, tenantID, req.ProjectID, req.UnitID, req.CustomerID, req.StageName, req.StageNumber,
		req.StagePercentage, stageDueAmount, req.ApartmentCost, stageDueAmount, stageDueAmount,
		"PENDING", dueDate, now, now,
	).Scan(&stage.ID, &stage.CreatedAt, &stage.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment stage: %w", err)
	}

	return stage, nil
}

// UpdatePaymentStageCollection records collection for a payment stage
func (s *ProjectManagementService) UpdatePaymentStageCollection(tenantID, stageID string, req *models.UpdatePaymentStageRequest) error {
	paymentDate := (*time.Time)(nil)
	if req.PaymentReceivedDate != "" {
		t, err := time.Parse("2006-01-02", req.PaymentReceivedDate)
		if err == nil {
			paymentDate = &t
		}
	}

	now := time.Now()

	query := `UPDATE property_payment_stage 
		SET amount_received = $1, amount_pending = GREATEST(0, amount_due - $1),
		 payment_received_date = $2, payment_mode = $3, reference_no = $4,
		 collection_status = $5, updated_at = $6
		WHERE id = $7 AND tenant_id = $8`

	_, err := s.DB.Exec(query,
		req.AmountReceived, paymentDate, req.PaymentMode, req.ReferenceNo,
		req.CollectionStatus, now, stageID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment stage collection: %w", err)
	}

	return nil
}

// ============================================================
// LIST/GET OPERATIONS
// ============================================================

// ListCustomerProfiles retrieves all customer profiles for a tenant
func (s *ProjectManagementService) ListCustomerProfiles(tenantID string, limit, offset int) ([]models.PropertyCustomerProfile, int, error) {
	query := `SELECT id, tenant_id, customer_code, unit_id, first_name, middle_name, last_name, email, phone_primary, phone_secondary,
		company_name, designation, pan_number, aadhar_number, customer_type, customer_status, booking_date, agreement_date,
		rate_per_sqft, car_parking_type, loan_required, sales_executive_id, sales_executive_name, created_at, updated_at
		FROM property_customer_profile 
		WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := s.DB.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []models.PropertyCustomerProfile
	for rows.Next() {
		var c models.PropertyCustomerProfile
		if err := rows.Scan(&c.ID, &c.TenantID, &c.CustomerCode, &c.UnitID, &c.FirstName, &c.MiddleName, &c.LastName,
			&c.Email, &c.PhonePrimary, &c.PhoneSecondary, &c.CompanyName, &c.Designation, &c.PANNumber, &c.AadharNumber,
			&c.CustomerType, &c.CustomerStatus, &c.BookingDate, &c.AgreementDate, &c.RatePerSqft, &c.CarParkingType,
			&c.LoanRequired, &c.SalesExecutiveID, &c.SalesExecutiveName, &c.CreatedAt, &c.UpdatedAt); err != nil {
			continue
		}
		customers = append(customers, c)
	}

	// Get count
	var count int
	countQuery := `SELECT COUNT(*) FROM property_customer_profile WHERE tenant_id = $1`
	s.DB.QueryRow(countQuery, tenantID).Scan(&count)

	return customers, count, nil
}

// GetAreaStatement retrieves a specific area statement
func (s *ProjectManagementService) GetAreaStatement(tenantID, unitID string) (*models.PropertyUnitAreaStatement, error) {
	query := `SELECT id, tenant_id, project_id, unit_id, apt_no, floor, unit_type, facing,
		rera_carpet_area_sqft, sbua_sqft, alloted_to, key_holder, percentage_allocation,
		noc_taken, active, created_at, updated_at
		FROM property_unit_area_statement 
		WHERE unit_id = $1 AND tenant_id = $2 AND active = 1`

	statement := &models.PropertyUnitAreaStatement{}
	err := s.DB.QueryRow(query, unitID, tenantID).Scan(
		&statement.ID, &statement.TenantID, &statement.ProjectID, &statement.UnitID,
		&statement.AptNo, &statement.Floor, &statement.UnitType, &statement.Facing,
		&statement.RERACarPetAreaSqft, &statement.SBUASqft, &statement.AlotedTo,
		&statement.KeyHolder, &statement.PercentageAllocation, &statement.NOCTaken,
		&statement.Active, &statement.CreatedAt, &statement.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("area statement not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get area statement: %w", err)
	}

	return statement, nil
}

// ListAreaStatements retrieves area statements for a project
func (s *ProjectManagementService) ListAreaStatements(tenantID, projectID string, limit, offset int) ([]models.PropertyUnitAreaStatement, int, error) {
	query := `SELECT id, tenant_id, project_id, unit_id, apt_no, floor, unit_type, facing,
		rera_carpet_area_sqft, sbua_sqft, alloted_to, created_at
		FROM property_unit_area_statement 
		WHERE tenant_id = $1 AND project_id = $2 AND active = 1
		ORDER BY apt_no DESC LIMIT $3 OFFSET $4`

	rows, err := s.DB.Query(query, tenantID, projectID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list area statements: %w", err)
	}
	defer rows.Close()

	var statements []models.PropertyUnitAreaStatement
	for rows.Next() {
		var s models.PropertyUnitAreaStatement
		if err := rows.Scan(&s.ID, &s.TenantID, &s.ProjectID, &s.UnitID, &s.AptNo, &s.Floor,
			&s.UnitType, &s.Facing, &s.RERACarPetAreaSqft, &s.SBUASqft, &s.AlotedTo, &s.CreatedAt); err != nil {
			continue
		}
		statements = append(statements, s)
	}

	var count int
	countQuery := `SELECT COUNT(*) FROM property_unit_area_statement WHERE tenant_id = $1 AND project_id = $2 AND active = 1`
	s.DB.QueryRow(countQuery, tenantID, projectID).Scan(&count)

	return statements, count, nil
}

// GetBankFinancing retrieves bank financing details
func (s *ProjectManagementService) GetBankFinancing(tenantID, financingID string) (*models.PropertyBankFinancing, error) {
	query := `SELECT id, tenant_id, project_id, unit_id, customer_id, apt_no, apartment_cost,
		sanctioned_amount, total_disbursed_amount, remaining_disbursement, total_collection_from_unit,
		disbursement_status, collection_status, outstanding_amount, noc_received, created_at, updated_at
		FROM property_bank_financing 
		WHERE id = $1 AND tenant_id = $2`

	financing := &models.PropertyBankFinancing{}
	err := s.DB.QueryRow(query, financingID, tenantID).Scan(
		&financing.ID, &financing.TenantID, &financing.ProjectID, &financing.UnitID, &financing.CustomerID,
		&financing.AptNo, &financing.ApartmentCost, &financing.SanctionedAmount, &financing.TotalDisbursedAmount,
		&financing.RemainingDisbursement, &financing.TotalCollectionFromUnit, &financing.DisbursementStatus,
		&financing.CollectionStatus, &financing.OutstandingAmount, &financing.NOCReceived,
		&financing.CreatedAt, &financing.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("financing record not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get financing: %w", err)
	}

	return financing, nil
}

// ListPaymentStages retrieves payment stages for a unit
func (s *ProjectManagementService) ListPaymentStages(tenantID, unitID string) ([]models.PropertyPaymentStage, error) {
	query := `SELECT id, tenant_id, stage_name, stage_number, stage_percentage, stage_due_amount,
		amount_due, amount_received, amount_pending, collection_status, due_date
		FROM property_payment_stage 
		WHERE tenant_id = $1 AND unit_id = $2
		ORDER BY stage_number ASC`

	rows, err := s.DB.Query(query, tenantID, unitID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment stages: %w", err)
	}
	defer rows.Close()

	var stages []models.PropertyPaymentStage
	for rows.Next() {
		var ps models.PropertyPaymentStage
		if err := rows.Scan(&ps.ID, &ps.TenantID, &ps.StageName, &ps.StageNumber, &ps.StagePercentage,
			&ps.StageDueAmount, &ps.AmountDue, &ps.AmountReceived, &ps.AmountPending, &ps.CollectionStatus, &ps.DueDate); err != nil {
			continue
		}
		stages = append(stages, ps)
	}

	return stages, nil
}

// ============================================================
// UPDATE/DELETE OPERATIONS
// ============================================================

// UpdateCustomerProfile updates customer profile
func (s *ProjectManagementService) UpdateCustomerProfile(tenantID, customerID string, req *models.PropertyCustomerProfile) error {
	now := time.Now()

	query := `UPDATE property_customer_profile 
		SET first_name = $1, middle_name = $2, last_name = $3, email = $4, phone_primary = $5,
		 phone_secondary = $6, alternate_phone = $7, company_name = $8, designation = $9,
		 pan_number = $10, aadhar_number = $11, pan_copy_url = $12, aadhar_copy_url = $13,
		 poa_document_no = $14, care_of = $15, communication_address_line1 = $16, communication_address_line2 = $17,
		 communication_city = $18, communication_state = $19, communication_country = $20, communication_zip = $21,
		 permanent_address_line1 = $22, permanent_address_line2 = $23, permanent_city = $24, permanent_state = $25,
		 permanent_country = $26, permanent_zip = $27, profession = $28, employer_name = $29, employment_type = $30,
		 monthly_income = $31, customer_type = $32,
		 co_applicant_1_name = $33, co_applicant_1_number = $34, co_applicant_1_alternate_number = $35, co_applicant_1_email = $36,
		 co_applicant_1_communication_address = $37, co_applicant_1_permanent_address = $38, co_applicant_1_aadhar = $39, co_applicant_1_pan = $40,
		 co_applicant_1_care_of = $41, co_applicant_1_relation = $42,
		 co_applicant_2_name = $43, co_applicant_2_number = $44, co_applicant_2_alternate_number = $45, co_applicant_2_email = $46,
		 co_applicant_2_communication_address = $47, co_applicant_2_permanent_address = $48, co_applicant_2_aadhar = $49, co_applicant_2_pan = $50,
		 co_applicant_2_care_of = $51, co_applicant_2_relation = $52,
		 co_applicant_3_name = $53, co_applicant_3_number = $54, co_applicant_3_alternate_number = $55, co_applicant_3_email = $56,
		 co_applicant_3_communication_address = $57, co_applicant_3_permanent_address = $58, co_applicant_3_aadhar = $59, co_applicant_3_pan = $60,
		 co_applicant_3_care_of = $61, co_applicant_3_relation = $62,
		 loan_required = $63, loan_amount = $64, loan_sanction_date = $65, bank_name = $66, bank_branch = $67,
		 bank_contact_person = $68, bank_contact_number = $69, connector_code_number = $70, lead_id = $71,
		 sales_executive_id = $72, sales_executive_name = $73, sales_head_id = $74, sales_head_name = $75,
		 booking_source = $76, life_certificate = $77, customer_status = $78, booking_date = $79, welcome_date = $80,
		 allotment_date = $81, agreement_date = $82, registration_date = $83, handover_date = $84, noc_received_date = $85,
		 rate_per_sqft = $86, composite_guideline_value = $87, car_parking_type = $88, maintenance_charges = $89,
		 other_works_charges = $90, corpus_charges = $91, eb_deposit = $92, notes = $93, updated_at = $94
		WHERE id = $95 AND tenant_id = $96`

	_, err := s.DB.Exec(query,
		req.FirstName, req.MiddleName, req.LastName, req.Email, req.PhonePrimary,
		req.PhoneSecondary, req.AlternatePhone, req.CompanyName, req.Designation,
		req.PANNumber, req.AadharNumber, req.PANCopyURL, req.AadharCopyURL,
		req.POADocumentNo, req.CareOf, req.CommunicationAddressLine1, req.CommunicationAddressLine2,
		req.CommunicationCity, req.CommunicationState, req.CommunicationCountry, req.CommunicationZip,
		req.PermanentAddressLine1, req.PermanentAddressLine2, req.PermanentCity, req.PermanentState,
		req.PermanentCountry, req.PermanentZip, req.Profession, req.EmployerName, req.EmploymentType,
		req.MonthlyIncome, req.CustomerType,
		req.CoApplicant1Name, req.CoApplicant1Number, req.CoApplicant1AlternateNumber, req.CoApplicant1Email,
		req.CoApplicant1CommunicationAddress, req.CoApplicant1PermanentAddress, req.CoApplicant1Aadhar, req.CoApplicant1PAN,
		req.CoApplicant1CareOf, req.CoApplicant1Relation,
		req.CoApplicant2Name, req.CoApplicant2Number, req.CoApplicant2AlternateNumber, req.CoApplicant2Email,
		req.CoApplicant2CommunicationAddress, req.CoApplicant2PermanentAddress, req.CoApplicant2Aadhar, req.CoApplicant2PAN,
		req.CoApplicant2CareOf, req.CoApplicant2Relation,
		req.CoApplicant3Name, req.CoApplicant3Number, req.CoApplicant3AlternateNumber, req.CoApplicant3Email,
		req.CoApplicant3CommunicationAddress, req.CoApplicant3PermanentAddress, req.CoApplicant3Aadhar, req.CoApplicant3PAN,
		req.CoApplicant3CareOf, req.CoApplicant3Relation,
		req.LoanRequired, req.LoanAmount, req.LoanSanctionDate, req.BankName, req.BankBranch,
		req.BankContactPerson, req.BankContactNumber, req.ConnectorCodeNumber, req.LeadID,
		req.SalesExecutiveID, req.SalesExecutiveName, req.SalesHeadID, req.SalesHeadName,
		req.BookingSource, req.LifeCertificate, req.CustomerStatus, req.BookingDate, req.WelcomeDate,
		req.AllotmentDate, req.AgreementDate, req.RegistrationDate, req.HandoverDate, req.NOCReceivedDate,
		req.RatePerSqft, req.CompositeGuidelineValue, req.CarParkingType, req.MaintenanceCharges,
		req.OtherWorksCharges, req.CorpusCharges, req.EBDeposit, req.Notes, now, customerID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update customer profile: %w", err)
	}

	return nil
}

// UpdateBankFinancing updates bank financing details
func (s *ProjectManagementService) UpdateBankFinancing(tenantID, financingID string, req *models.CreateBankFinancingRequest) error {
	now := time.Now()

	query := `UPDATE property_bank_financing 
		SET bank_name = $1, sanctioned_amount = $2, total_commitment = $3, updated_at = $4
		WHERE id = $5 AND tenant_id = $6`

	_, err := s.DB.Exec(query,
		req.BankName, req.SanctionedAmount, req.TotalCommitment, now, financingID, tenantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update bank financing: %w", err)
	}

	return nil
}

// DeleteAreaStatement soft deletes an area statement
func (s *ProjectManagementService) DeleteAreaStatement(tenantID, unitID string) error {
	now := time.Now()

	query := `UPDATE property_unit_area_statement SET active = 0, updated_at = $1
		WHERE unit_id = $2 AND tenant_id = $3`

	_, err := s.DB.Exec(query, now, unitID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete area statement: %w", err)
	}

	return nil
}

// ============================================================
// CALCULATION & SUMMARY OPERATIONS
// ============================================================

// CalculateCostBreakdown calculates the complete cost breakdown
func (s *ProjectManagementService) CalculateCostBreakdown(tenantID, unitID string) (map[string]float64, error) {
	query := `SELECT sbua, rate_per_sqft, car_parking_cost, plc, statutory_approval_charge,
		legal_documentation_charge, amenities_equipment_charge, other_charges_1, other_charges_2,
		other_charges_3, other_charges_4, other_charges_5, gst_percentage, grand_total
		FROM unit_cost_sheet WHERE unit_id = $1 AND tenant_id = $2`

	var sbua, ratePerSqft, carParkingCost, plc, statutoryCharge, legalCharge, amenitiesCharge sql.NullFloat64
	var oc1, oc2, oc3, oc4, oc5, gstPct, grandTotal sql.NullFloat64

	err := s.DB.QueryRow(query, unitID, tenantID).Scan(
		&sbua, &ratePerSqft, &carParkingCost, &plc, &statutoryCharge, &legalCharge, &amenitiesCharge,
		&oc1, &oc2, &oc3, &oc4, &oc5, &gstPct, &grandTotal,
	)

	if err == sql.ErrNoRows {
		return map[string]float64{}, fmt.Errorf("cost sheet not found")
	}
	if err != nil {
		return map[string]float64{}, fmt.Errorf("failed to calculate cost breakdown: %w", err)
	}

	breakdown := map[string]float64{
		"sbua":              sbua.Float64,
		"rate_per_sqft":     ratePerSqft.Float64,
		"apartment_cost":    sbua.Float64 * ratePerSqft.Float64,
		"car_parking_cost":  carParkingCost.Float64,
		"statutory_charges": statutoryCharge.Float64,
		"legal_charges":     legalCharge.Float64,
		"amenities_charges": amenitiesCharge.Float64,
		"other_charges_1":   oc1.Float64,
		"other_charges_2":   oc2.Float64,
		"other_charges_3":   oc3.Float64,
		"other_charges_4":   oc4.Float64,
		"other_charges_5":   oc5.Float64,
		"gst_percentage":    gstPct.Float64,
		"grand_total":       grandTotal.Float64,
	}

	return breakdown, nil
}

// GetProjectSummary gets comprehensive project summary
func (s *ProjectManagementService) GetProjectSummary(tenantID, projectID string) (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Total units
	var totalUnits int
	s.DB.QueryRow(`SELECT COUNT(*) FROM property_unit_area_statement 
		WHERE tenant_id = $1 AND project_id = $2 AND active = 1`, tenantID, projectID).Scan(&totalUnits)

	// Total sanctioned amount
	var totalSanctioned, totalDisbursed, totalCollected sql.NullFloat64
	s.DB.QueryRow(`SELECT SUM(sanctioned_amount), SUM(total_disbursed_amount), SUM(total_collection_from_unit)
		FROM property_bank_financing WHERE tenant_id = $1 AND project_id = $2`,
		tenantID, projectID).Scan(&totalSanctioned, &totalDisbursed, &totalCollected)

	// Collection status breakdown
	var pendingCount, partialCount, completedCount int
	s.DB.QueryRow(`SELECT COUNT(*) FROM property_payment_stage 
		WHERE tenant_id = $1 AND project_id = $2 AND collection_status = 'PENDING'`,
		tenantID, projectID).Scan(&pendingCount)
	s.DB.QueryRow(`SELECT COUNT(*) FROM property_payment_stage 
		WHERE tenant_id = $1 AND project_id = $2 AND collection_status = 'PARTIAL'`,
		tenantID, projectID).Scan(&partialCount)
	s.DB.QueryRow(`SELECT COUNT(*) FROM property_payment_stage 
		WHERE tenant_id = $1 AND project_id = $2 AND collection_status = 'COMPLETED'`,
		tenantID, projectID).Scan(&completedCount)

	summary["total_units"] = totalUnits
	summary["total_sanctioned_amount"] = totalSanctioned.Float64
	summary["total_disbursed_amount"] = totalDisbursed.Float64
	summary["total_collected_amount"] = totalCollected.Float64
	summary["payment_status"] = map[string]int{
		"pending":   pendingCount,
		"partial":   partialCount,
		"completed": completedCount,
	}

	return summary, nil
}
