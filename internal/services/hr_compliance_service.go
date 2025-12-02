package services

import (
	"database/sql"
	"fmt"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
)

// ============================================================================
// HR COMPLIANCE SERVICE
// ============================================================================
// Manages Labour Law compliance - ESI, EPF, PF, Gratuity, Bonus, Leaves, etc.

type HRComplianceService struct {
	DB *sql.DB
}

func NewHRComplianceService(db *sql.DB) *HRComplianceService {
	return &HRComplianceService{DB: db}
}

// CreateESICompliance creates ESI compliance record for an employee
// ESI: Employee State Insurance - 0.75% employee + 3.25% employer contribution
func (s *HRComplianceService) CreateESICompliance(
	tenantID, employeeID string,
	esiNumber, esiOffice string,
) (*models.ESICompliance, error) {

	esiState := "Maharashtra" // Will be extracted from employee address in actual implementation

	esiCompliance := &models.ESICompliance{
		ID:                       fmt.Sprintf("ESI-%s-%d", employeeID, time.Now().Unix()),
		TenantID:                 tenantID,
		EmployeeID:               employeeID,
		ESINumber:                esiNumber,
		ESIRegistrationDate:      &time.Time{},
		ESIOffice:                esiOffice,
		ESIState:                 esiState,
		EffectiveFrom:            time.Now(),
		EmployeeContributionRate: 0.75,
		EmployerContributionRate: 3.25,
		WageLimit:                21000,
		SickLeaveBalance:         0,
		IsCompliant:              true,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	query := `INSERT INTO esi_compliance 
		(id, tenant_id, employee_id, esi_number, esi_office, esi_state, 
		 effective_from, employee_contribution_rate, employer_contribution_rate, 
		 wage_limit, sick_leave_balance, is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		esiCompliance.ID, esiCompliance.TenantID, esiCompliance.EmployeeID,
		esiCompliance.ESINumber, esiCompliance.ESIOffice, esiCompliance.ESIState,
		esiCompliance.EffectiveFrom, esiCompliance.EmployeeContributionRate,
		esiCompliance.EmployerContributionRate, esiCompliance.WageLimit,
		esiCompliance.SickLeaveBalance, esiCompliance.IsCompliant,
		esiCompliance.CreatedAt, esiCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create ESI compliance: %w", err)
	}

	return esiCompliance, nil
}

// CreateEPFCompliance creates EPF/PF compliance record for an employee
// EPF: Employee Provident Fund - 12% employee + 12% employer contribution
func (s *HRComplianceService) CreateEPFCompliance(
	tenantID, employeeID string,
	pfAccountNumber, pfOffice string,
	pfExempt bool,
) (*models.EPFCompliance, error) {

	pfState := "Maharashtra" // Will be extracted from employee address

	epfCompliance := &models.EPFCompliance{
		ID:                       fmt.Sprintf("EPF-%s-%d", employeeID, time.Now().Unix()),
		TenantID:                 tenantID,
		EmployeeID:               employeeID,
		PFAccountNumber:          pfAccountNumber,
		PFOffice:                 pfOffice,
		PFState:                  pfState,
		EffectiveFrom:            time.Now(),
		EmployeeContributionRate: 12,
		EmployerContributionRate: 12,
		VPSContributionRate:      0,
		PFExempt:                 pfExempt,
		AccumulatedBalance:       0,
		CurrentBalance:           0,
		PartialWithdrawalAllowed: true,
		IsCompliant:              true,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	query := `INSERT INTO epf_compliance 
		(id, tenant_id, employee_id, pf_account_number, pf_office, pf_state, 
		 effective_from, employee_contribution_rate, employer_contribution_rate, 
		 pf_exempt, accumulated_balance, current_balance, partial_withdrawal_allowed, 
		 is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		epfCompliance.ID, epfCompliance.TenantID, epfCompliance.EmployeeID,
		epfCompliance.PFAccountNumber, epfCompliance.PFOffice, epfCompliance.PFState,
		epfCompliance.EffectiveFrom, epfCompliance.EmployeeContributionRate,
		epfCompliance.EmployerContributionRate, epfCompliance.PFExempt,
		epfCompliance.AccumulatedBalance, epfCompliance.CurrentBalance,
		epfCompliance.PartialWithdrawalAllowed, epfCompliance.IsCompliant,
		epfCompliance.CreatedAt, epfCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create EPF compliance: %w", err)
	}

	return epfCompliance, nil
}

// CreateProfessionalTaxCompliance creates Professional Tax (PT) compliance
// PT varies by state: e.g., Maharashtra, Tamil Nadu, etc.
func (s *HRComplianceService) CreateProfessionalTaxCompliance(
	tenantID, employeeID string,
	ptState string,
	monthlySalaryThreshold, ptAmount float64,
) (*models.ProfessionalTaxCompliance, error) {

	ptCompliance := &models.ProfessionalTaxCompliance{
		ID:                     fmt.Sprintf("PT-%s-%d", employeeID, time.Now().Unix()),
		TenantID:               tenantID,
		EmployeeID:             employeeID,
		PTRegistrationState:    ptState,
		MonthlySalaryThreshold: monthlySalaryThreshold,
		PTAmount:               ptAmount,
		ApplicableFrom:         time.Now(),
		IsCompliant:            true,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	query := `INSERT INTO professional_tax_compliance 
		(id, tenant_id, employee_id, pt_registration_state, monthly_salary_threshold, 
		 pt_amount, applicable_from, is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		ptCompliance.ID, ptCompliance.TenantID, ptCompliance.EmployeeID,
		ptCompliance.PTRegistrationState, ptCompliance.MonthlySalaryThreshold,
		ptCompliance.PTAmount, ptCompliance.ApplicableFrom, ptCompliance.IsCompliant,
		ptCompliance.CreatedAt, ptCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create PT compliance: %w", err)
	}

	return ptCompliance, nil
}

// CreateGratuityCompliance creates Gratuity compliance record
// Gratuity: 15 days' salary per year for first 5 years + 30 days for subsequent years
// Applicable only if service > 5 years
func (s *HRComplianceService) CreateGratuityCompliance(
	tenantID, employeeID string,
	joiningDate time.Time,
) (*models.GratuityCompliance, error) {

	gratuityCompliance := &models.GratuityCompliance{
		ID:               fmt.Sprintf("GRT-%s-%d", employeeID, time.Now().Unix()),
		TenantID:         tenantID,
		EmployeeID:       employeeID,
		JoiningDate:      joiningDate,
		GratuityEligible: false, // Will be set to true when service > 5 years
		YearsOfService:   0,
		FundType:         "Gratuity_Fund",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `INSERT INTO gratuity_compliance 
		(id, tenant_id, employee_id, joining_date, gratuity_eligible, 
		 years_of_service, fund_type, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		gratuityCompliance.ID, gratuityCompliance.TenantID, gratuityCompliance.EmployeeID,
		gratuityCompliance.JoiningDate, gratuityCompliance.GratuityEligible,
		gratuityCompliance.YearsOfService, gratuityCompliance.FundType,
		gratuityCompliance.CreatedAt, gratuityCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gratuity compliance: %w", err)
	}

	return gratuityCompliance, nil
}

// CheckAndUpdateGratuityEligibility checks if employee is eligible for gratuity (5+ years)
func (s *HRComplianceService) CheckAndUpdateGratuityEligibility(
	tenantID, employeeID string,
	lastDrawnSalary float64,
) error {

	query := `SELECT joining_date FROM gratuity_compliance 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL`

	var joiningDate time.Time
	err := s.DB.QueryRow(query, tenantID, employeeID).Scan(&joiningDate)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get gratuity record: %w", err)
	}

	// Calculate years of service
	yearsOfService := int(time.Since(joiningDate).Hours() / (24 * 365))
	gratuityEligible := yearsOfService >= 5

	var gratuityAccrued float64
	if gratuityEligible {
		// Calculate gratuity: 15 days' salary for first 5 years + 30 days for subsequent
		daysPerYear := float64(15)
		if yearsOfService > 5 {
			daysPerYear = float64(30)
		}
		gratuityAccrued = (lastDrawnSalary / 30) * daysPerYear * float64(yearsOfService)
	}

	updateQuery := `UPDATE gratuity_compliance 
		SET gratuity_eligible = ?, years_of_service = ?, gratuity_accrued = ?, 
		    last_drawn_salary = ?, updated_at = ? 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL`

	_, err = s.DB.Exec(updateQuery,
		gratuityEligible, yearsOfService, gratuityAccrued, lastDrawnSalary,
		time.Now(), tenantID, employeeID,
	)

	if err != nil {
		return fmt.Errorf("failed to update gratuity eligibility: %w", err)
	}

	return nil
}

// RecordBonusPayment records bonus payment for an employee
func (s *HRComplianceService) RecordBonusPayment(
	tenantID, employeeID string,
	bonusType string,
	bonusAmount float64,
	bonusYear int,
) (*models.BonusCompliance, error) {

	bonus := &models.BonusCompliance{
		ID:                        fmt.Sprintf("BONUS-%s-%d", employeeID, time.Now().Unix()),
		TenantID:                  tenantID,
		EmployeeID:                employeeID,
		BonusType:                 bonusType,
		Eligible:                  true,
		ApplicableSalaryComponent: "Basic_DA",
		BonusAmount:               bonusAmount,
		BonusYear:                 bonusYear,
		PaymentStatus:             "Paid",
		PaymentDate:               &time.Time{},
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}

	now := time.Now()
	bonus.PaymentDate = &now

	query := `INSERT INTO bonus_compliance 
		(id, tenant_id, employee_id, bonus_type, bonus_name, eligible, 
		 applicable_salary_component, bonus_amount, bonus_year, payment_date, 
		 payment_status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		bonus.ID, bonus.TenantID, bonus.EmployeeID, bonus.BonusType, bonus.BonusType,
		bonus.Eligible, bonus.ApplicableSalaryComponent, bonus.BonusAmount,
		bonus.BonusYear, bonus.PaymentDate, bonus.PaymentStatus,
		bonus.CreatedAt, bonus.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to record bonus payment: %w", err)
	}

	return bonus, nil
}

// InitializeLeaveCompliance sets up leave entitlements for a fiscal year
func (s *HRComplianceService) InitializeLeaveCompliance(
	tenantID, employeeID, leaveTypeID string,
	fiscalYearStart, fiscalYearEnd time.Time,
	annualEntitlement int,
) (*models.LeaveCompliance, error) {

	leaveCompliance := &models.LeaveCompliance{
		ID:                fmt.Sprintf("LEAVE-%s-%d", employeeID, time.Now().Unix()),
		TenantID:          tenantID,
		EmployeeID:        employeeID,
		LeaveTypeID:       leaveTypeID,
		FiscalYearFrom:    fiscalYearStart,
		FiscalYearTo:      fiscalYearEnd,
		AnnualEntitlement: annualEntitlement,
		OpeningBalance:    0,
		CarryForwardLimit: 5,
		Utilized:          0,
		Available:         annualEntitlement,
		EncashmentAllowed: true,
		IsCompliant:       true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	query := `INSERT INTO leave_compliance 
		(id, tenant_id, employee_id, leave_type_id, fiscal_year_from, fiscal_year_to, 
		 annual_entitlement, opening_balance, carry_forward_limit, utilized, available, 
		 encashment_allowed, is_compliant, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		leaveCompliance.ID, leaveCompliance.TenantID, leaveCompliance.EmployeeID,
		leaveCompliance.LeaveTypeID, leaveCompliance.FiscalYearFrom,
		leaveCompliance.FiscalYearTo, leaveCompliance.AnnualEntitlement,
		leaveCompliance.OpeningBalance, leaveCompliance.CarryForwardLimit,
		leaveCompliance.Utilized, leaveCompliance.Available,
		leaveCompliance.EncashmentAllowed, leaveCompliance.IsCompliant,
		leaveCompliance.CreatedAt, leaveCompliance.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize leave compliance: %w", err)
	}

	return leaveCompliance, nil
}

// LogComplianceAudit creates a compliance audit log entry
func (s *HRComplianceService) LogComplianceAudit(
	tenantID, employeeID string,
	complianceType, complianceItem string,
	isCompliant bool,
	violationFound string,
	severity string,
) (*models.StatutoryComplianceAudit, error) {

	audit := &models.StatutoryComplianceAudit{
		ID:               fmt.Sprintf("AUDIT-%s-%d", employeeID, time.Now().Unix()),
		TenantID:         tenantID,
		EmployeeID:       employeeID,
		ComplianceType:   complianceType,
		ComplianceItem:   complianceItem,
		CompliancePeriod: &time.Time{},
		IsCompliant:      isCompliant,
		ViolationFound:   violationFound,
		Severity:         severity,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	now := time.Now()
	audit.CompliancePeriod = &now

	query := `INSERT INTO statutory_compliance_audit 
		(id, tenant_id, employee_id, compliance_type, compliance_item, compliance_period, 
		 is_compliant, violation_found, severity, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		audit.ID, audit.TenantID, audit.EmployeeID, audit.ComplianceType,
		audit.ComplianceItem, audit.CompliancePeriod, audit.IsCompliant,
		audit.ViolationFound, audit.Severity, audit.CreatedAt, audit.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to log compliance audit: %w", err)
	}

	return audit, nil
}

// GetEmployeeComplianceStatus returns comprehensive compliance status for an employee
func (s *HRComplianceService) GetEmployeeComplianceStatus(tenantID, employeeID string) (*models.ComplianceCheckResult, error) {

	// Get employee name
	var employeeName string
	nameQuery := `SELECT CONCAT(first_name, ' ', last_name) FROM employees 
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`
	err := s.DB.QueryRow(nameQuery, employeeID, tenantID).Scan(&employeeName)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	result := &models.ComplianceCheckResult{
		EmployeeID:   employeeID,
		EmployeeName: employeeName,
	}

	// Check ESI Compliance
	var esiCompliant bool
	esiQuery := `SELECT is_compliant FROM esi_compliance 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL`
	err = s.DB.QueryRow(esiQuery, tenantID, employeeID).Scan(&esiCompliant)
	if err == nil {
		result.ESICompliance = &models.ComplianceStatus{
			IsCompliant: esiCompliant,
		}
	}

	// Check EPF Compliance
	var epfCompliant bool
	epfQuery := `SELECT is_compliant FROM epf_compliance 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL`
	err = s.DB.QueryRow(epfQuery, tenantID, employeeID).Scan(&epfCompliant)
	if err == nil {
		result.EPFCompliance = &models.ComplianceStatus{
			IsCompliant: epfCompliant,
		}
	}

	// Check PT Compliance
	var ptCompliant bool
	ptQuery := `SELECT is_compliant FROM professional_tax_compliance 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL`
	err = s.DB.QueryRow(ptQuery, tenantID, employeeID).Scan(&ptCompliant)
	if err == nil {
		result.PTCompliance = &models.ComplianceStatus{
			IsCompliant: ptCompliant,
		}
	}

	// Check Gratuity Compliance
	var gratuityCompliant bool
	gratuityQuery := `SELECT gratuity_eligible FROM gratuity_compliance 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL`
	err = s.DB.QueryRow(gratuityQuery, tenantID, employeeID).Scan(&gratuityCompliant)
	if err == nil {
		result.GratuityCompliance = &models.ComplianceStatus{
			IsCompliant: gratuityCompliant,
		}
	}

	// Check Leave Compliance
	var leaveCompliant bool
	leaveQuery := `SELECT is_compliant FROM leave_compliance 
		WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL ORDER BY fiscal_year_from DESC LIMIT 1`
	err = s.DB.QueryRow(leaveQuery, tenantID, employeeID).Scan(&leaveCompliant)
	if err == nil {
		result.LeaveCompliance = &models.ComplianceStatus{
			IsCompliant: leaveCompliant,
		}
	}

	// Determine overall status
	overallCompliant := true
	criticalViolations := 0
	if result.ESICompliance != nil && !result.ESICompliance.IsCompliant {
		overallCompliant = false
		criticalViolations++
	}
	if result.EPFCompliance != nil && !result.EPFCompliance.IsCompliant {
		overallCompliant = false
		criticalViolations++
	}

	if overallCompliant {
		result.OverallStatus = "Compliant"
	} else if criticalViolations == 0 {
		result.OverallStatus = "Action_Required"
	} else {
		result.OverallStatus = "Non-Compliant"
	}

	result.CriticalViolations = criticalViolations
	result.Recommendations = []string{
		"Ensure all statutory documents are filed on time",
		"Review ESI and EPF contributions regularly",
		"Maintain leave records and balance",
		"Track gratuity accruals",
	}

	return result, nil
}
