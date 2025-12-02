package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
)

// HRService handles HR and Payroll operations
type HRService struct {
	DB *sql.DB
}

// NewHRService creates a new HR service
func NewHRService(db *sql.DB) *HRService {
	return &HRService{DB: db}
}

// ============================================================================
// EMPLOYEE MANAGEMENT
// ============================================================================

// CreateEmployee creates a new employee record
func (s *HRService) CreateEmployee(tenantID string, emp *models.Employee) error {
	emp.TenantID = tenantID
	emp.CreatedAt = time.Now()
	emp.UpdatedAt = time.Now()

	query := `INSERT INTO employees (
		id, tenant_id, first_name, last_name, email, phone, date_of_birth, gender, nationality,
		address, city, state, country, postal_code, employee_id, designation, department, report_to,
		employment_type, joining_date, status, bank_account_number, bank_ifsc_code, bank_name,
		account_holder_name, base_salary, dearness_allowance, house_rent_allowance, special_allowance,
		conveyance_allowance, medical_allowance, other_allowances, epf_deduction, esi_deduction,
		professional_tax, income_tax, loan_deduction, advance_deduction, other_deductions,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query,
		emp.ID, emp.TenantID, emp.FirstName, emp.LastName, emp.Email, emp.Phone, emp.DateOfBirth, emp.Gender, emp.Nationality,
		emp.Address, emp.City, emp.State, emp.Country, emp.PostalCode, emp.EmployeeID, emp.Designation, emp.Department, emp.ReportTo,
		emp.EmploymentType, emp.JoiningDate, emp.Status, emp.BankAccountNumber, emp.BankIFSCCode, emp.BankName,
		emp.AccountHolderName, emp.BaseSalary, emp.DAAllowance, emp.HRAAllowance, emp.SpecialAllowance,
		emp.ConveyanceAllowance, emp.MedicalAllowance, emp.OtherAllowances, emp.EPFDeduction, emp.ESIDeduction,
		emp.ProfessionalTax, emp.IncomeTax, emp.LoanDeduction, emp.AdvanceDeduction, emp.OtherDeductions,
		emp.CreatedAt, emp.UpdatedAt,
	)

	return err
}

// GetEmployee retrieves an employee by ID
func (s *HRService) GetEmployee(tenantID, employeeID string) (*models.Employee, error) {
	var emp models.Employee
	query := `SELECT id, tenant_id, first_name, last_name, email, phone, date_of_birth, gender, nationality,
		address, city, state, country, postal_code, employee_id, designation, department, report_to,
		employment_type, joining_date, status, bank_account_number, bank_ifsc_code, bank_name,
		account_holder_name, base_salary, dearness_allowance, house_rent_allowance, special_allowance,
		conveyance_allowance, medical_allowance, other_allowances, epf_deduction, esi_deduction,
		professional_tax, income_tax, loan_deduction, advance_deduction, other_deductions,
		created_at, updated_at, deleted_at
		FROM employees WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, employeeID, tenantID).Scan(
		&emp.ID, &emp.TenantID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Phone, &emp.DateOfBirth, &emp.Gender, &emp.Nationality,
		&emp.Address, &emp.City, &emp.State, &emp.Country, &emp.PostalCode, &emp.EmployeeID, &emp.Designation, &emp.Department, &emp.ReportTo,
		&emp.EmploymentType, &emp.JoiningDate, &emp.Status, &emp.BankAccountNumber, &emp.BankIFSCCode, &emp.BankName,
		&emp.AccountHolderName, &emp.BaseSalary, &emp.DAAllowance, &emp.HRAAllowance, &emp.SpecialAllowance,
		&emp.ConveyanceAllowance, &emp.MedicalAllowance, &emp.OtherAllowances, &emp.EPFDeduction, &emp.ESIDeduction,
		&emp.ProfessionalTax, &emp.IncomeTax, &emp.LoanDeduction, &emp.AdvanceDeduction, &emp.OtherDeductions,
		&emp.CreatedAt, &emp.UpdatedAt, &emp.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found")
	}
	return &emp, err
}

// ListEmployees retrieves employees with pagination
func (s *HRService) ListEmployees(tenantID string, limit, offset int) ([]models.Employee, int, error) {
	var employees []models.Employee

	// Get total count
	countQuery := `SELECT COUNT(*) FROM employees WHERE tenant_id = ? AND deleted_at IS NULL`
	var total int
	s.DB.QueryRow(countQuery, tenantID).Scan(&total)

	query := `SELECT id, tenant_id, first_name, last_name, email, phone, date_of_birth, gender, nationality,
		address, city, state, country, postal_code, employee_id, designation, department, report_to,
		employment_type, joining_date, status, bank_account_number, bank_ifsc_code, bank_name,
		account_holder_name, base_salary, dearness_allowance, house_rent_allowance, special_allowance,
		conveyance_allowance, medical_allowance, other_allowances, epf_deduction, esi_deduction,
		professional_tax, income_tax, loan_deduction, advance_deduction, other_deductions,
		created_at, updated_at, deleted_at
		FROM employees WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := s.DB.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.ID, &emp.TenantID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Phone, &emp.DateOfBirth, &emp.Gender, &emp.Nationality,
			&emp.Address, &emp.City, &emp.State, &emp.Country, &emp.PostalCode, &emp.EmployeeID, &emp.Designation, &emp.Department, &emp.ReportTo,
			&emp.EmploymentType, &emp.JoiningDate, &emp.Status, &emp.BankAccountNumber, &emp.BankIFSCCode, &emp.BankName,
			&emp.AccountHolderName, &emp.BaseSalary, &emp.DAAllowance, &emp.HRAAllowance, &emp.SpecialAllowance,
			&emp.ConveyanceAllowance, &emp.MedicalAllowance, &emp.OtherAllowances, &emp.EPFDeduction, &emp.ESIDeduction,
			&emp.ProfessionalTax, &emp.IncomeTax, &emp.LoanDeduction, &emp.AdvanceDeduction, &emp.OtherDeductions,
			&emp.CreatedAt, &emp.UpdatedAt, &emp.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning employee: %v", err)
			continue
		}
		employees = append(employees, emp)
	}

	return employees, total, rows.Err()
}

// UpdateEmployee updates an employee record
func (s *HRService) UpdateEmployee(tenantID string, emp *models.Employee) error {
	emp.UpdatedAt = time.Now()

	query := `UPDATE employees SET 
		first_name = ?, last_name = ?, email = ?, phone = ?, gender = ?, 
		address = ?, city = ?, state = ?, country = ?, postal_code = ?,
		designation = ?, department = ?, report_to = ?, status = ?,
		bank_account_number = ?, bank_ifsc_code = ?, bank_name = ?, account_holder_name = ?,
		base_salary = ?, dearness_allowance = ?, house_rent_allowance = ?, special_allowance = ?,
		conveyance_allowance = ?, medical_allowance = ?, other_allowances = ?,
		epf_deduction = ?, esi_deduction = ?, professional_tax = ?, income_tax = ?,
		loan_deduction = ?, advance_deduction = ?, other_deductions = ?, updated_at = ?
		WHERE id = ? AND tenant_id = ?`

	_, err := s.DB.Exec(query,
		emp.FirstName, emp.LastName, emp.Email, emp.Phone, emp.Gender,
		emp.Address, emp.City, emp.State, emp.Country, emp.PostalCode,
		emp.Designation, emp.Department, emp.ReportTo, emp.Status,
		emp.BankAccountNumber, emp.BankIFSCCode, emp.BankName, emp.AccountHolderName,
		emp.BaseSalary, emp.DAAllowance, emp.HRAAllowance, emp.SpecialAllowance,
		emp.ConveyanceAllowance, emp.MedicalAllowance, emp.OtherAllowances,
		emp.EPFDeduction, emp.ESIDeduction, emp.ProfessionalTax, emp.IncomeTax,
		emp.LoanDeduction, emp.AdvanceDeduction, emp.OtherDeductions, emp.UpdatedAt,
		emp.ID, tenantID,
	)

	return err
}

// DeleteEmployee soft deletes an employee
func (s *HRService) DeleteEmployee(tenantID, employeeID string) error {
	query := `UPDATE employees SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?`
	_, err := s.DB.Exec(query, employeeID, tenantID)
	return err
}

// ============================================================================
// ATTENDANCE MANAGEMENT
// ============================================================================

// RecordAttendance records employee attendance
func (s *HRService) RecordAttendance(tenantID string, att *models.Attendance) error {
	att.TenantID = tenantID
	att.CreatedAt = time.Now()
	att.UpdatedAt = time.Now()

	query := `INSERT INTO attendance (id, tenant_id, employee_id, attendance_date, check_in_time, check_out_time, working_hours, status, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query, att.ID, att.TenantID, att.EmployeeID, att.AttendanceDate, att.CheckInTime, att.CheckOutTime, att.WorkingHours, att.Status, att.Notes, att.CreatedAt, att.UpdatedAt)
	return err
}

// GetAttendanceRecord retrieves a specific attendance record
func (s *HRService) GetAttendanceRecord(tenantID, employeeID string, date time.Time) (*models.Attendance, error) {
	var att models.Attendance
	query := `SELECT id, tenant_id, employee_id, attendance_date, check_in_time, check_out_time, working_hours, status, notes, created_at, updated_at, deleted_at
		FROM attendance WHERE tenant_id = ? AND employee_id = ? AND DATE(attendance_date) = DATE(?) AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, tenantID, employeeID, date).Scan(
		&att.ID, &att.TenantID, &att.EmployeeID, &att.AttendanceDate, &att.CheckInTime, &att.CheckOutTime, &att.WorkingHours, &att.Status, &att.Notes, &att.CreatedAt, &att.UpdatedAt, &att.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("attendance record not found")
	}
	return &att, err
}

// ListEmployeeAttendance retrieves attendance records for an employee
func (s *HRService) ListEmployeeAttendance(tenantID, employeeID string, fromDate, toDate time.Time) ([]models.Attendance, error) {
	var records []models.Attendance

	query := `SELECT id, tenant_id, employee_id, attendance_date, check_in_time, check_out_time, working_hours, status, notes, created_at, updated_at, deleted_at
		FROM attendance WHERE tenant_id = ? AND employee_id = ? AND attendance_date BETWEEN ? AND ? AND deleted_at IS NULL
		ORDER BY attendance_date DESC`

	rows, err := s.DB.Query(query, tenantID, employeeID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var att models.Attendance
		err := rows.Scan(
			&att.ID, &att.TenantID, &att.EmployeeID, &att.AttendanceDate, &att.CheckInTime, &att.CheckOutTime, &att.WorkingHours, &att.Status, &att.Notes, &att.CreatedAt, &att.UpdatedAt, &att.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning attendance: %v", err)
			continue
		}
		records = append(records, att)
	}

	return records, rows.Err()
}

// ============================================================================
// PAYROLL MANAGEMENT
// ============================================================================

// CalculateAndCreatePayroll calculates and creates payroll for an employee
func (s *HRService) CalculateAndCreatePayroll(tenantID, employeeID string, payrollMonth time.Time) (*models.PayrollRecord, error) {
	// Get employee details
	emp, err := s.GetEmployee(tenantID, employeeID)
	if err != nil {
		return nil, err
	}

	payroll := &models.PayrollRecord{
		TenantID:      tenantID,
		EmployeeID:    employeeID,
		PayrollMonth:  payrollMonth,
		PayrollStatus: "generated",

		// Earnings
		BasicSalary:      emp.BaseSalary,
		DAAllowance:      emp.DAAllowance,
		HRAAllowance:     emp.HRAAllowance,
		SpecialAllowance: emp.SpecialAllowance,
		ConveyanceAllow:  emp.ConveyanceAllowance,
		MedicalAllow:     emp.MedicalAllowance,
		OtherAllowances:  emp.OtherAllowances,

		// Deductions
		EPFDeduction:     emp.EPFDeduction,
		ESIDeduction:     emp.ESIDeduction,
		ProfessionalTax:  emp.ProfessionalTax,
		IncomeTax:        emp.IncomeTax,
		LoanDeduction:    emp.LoanDeduction,
		AdvanceDeduction: emp.AdvanceDeduction,
		OtherDeductions:  emp.OtherDeductions,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate totals
	payroll.TotalEarnings = payroll.BasicSalary + payroll.DAAllowance + payroll.HRAAllowance + payroll.SpecialAllowance + payroll.ConveyanceAllow + payroll.MedicalAllow + payroll.OtherAllowances
	payroll.TotalDeductions = payroll.EPFDeduction + payroll.ESIDeduction + payroll.ProfessionalTax + payroll.IncomeTax + payroll.LoanDeduction + payroll.AdvanceDeduction + payroll.OtherDeductions
	payroll.NetSalary = payroll.TotalEarnings - payroll.TotalDeductions

	// Insert into database
	query := `INSERT INTO payroll (id, tenant_id, employee_id, payroll_month, payroll_status,
		basic_salary, dearness_allowance, house_rent_allowance, special_allowance,
		conveyance_allowance, medical_allowance, other_allowances, total_earnings,
		epf_deduction, esi_deduction, professional_tax, income_tax, loan_deduction, advance_deduction, other_deductions, total_deductions,
		net_salary, working_days, leave_days, paid_days, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	id := fmt.Sprintf("%s-%s-%d", employeeID, payrollMonth.Format("2006-01"), time.Now().UnixNano())
	_, err = s.DB.Exec(query,
		id, tenantID, employeeID, payrollMonth, payroll.PayrollStatus,
		payroll.BasicSalary, payroll.DAAllowance, payroll.HRAAllowance, payroll.SpecialAllowance,
		payroll.ConveyanceAllow, payroll.MedicalAllow, payroll.OtherAllowances, payroll.TotalEarnings,
		payroll.EPFDeduction, payroll.ESIDeduction, payroll.ProfessionalTax, payroll.IncomeTax, payroll.LoanDeduction, payroll.AdvanceDeduction, payroll.OtherDeductions, payroll.TotalDeductions,
		payroll.NetSalary, payroll.WorkingDays, payroll.LeaveDays, payroll.PaidDays, payroll.CreatedAt, payroll.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	payroll.ID = id
	return payroll, nil
}

// GetPayrollRecord retrieves a payroll record
func (s *HRService) GetPayrollRecord(tenantID, payrollID string) (*models.PayrollRecord, error) {
	var payroll models.PayrollRecord
	query := `SELECT id, tenant_id, employee_id, payroll_month, payroll_status,
		basic_salary, dearness_allowance, house_rent_allowance, special_allowance,
		conveyance_allowance, medical_allowance, other_allowances, total_earnings,
		epf_deduction, esi_deduction, professional_tax, income_tax, loan_deduction, advance_deduction, other_deductions, total_deductions,
		net_salary, working_days, leave_days, paid_days, notes, created_at, updated_at, processed_at, deleted_at
		FROM payroll WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := s.DB.QueryRow(query, payrollID, tenantID).Scan(
		&payroll.ID, &payroll.TenantID, &payroll.EmployeeID, &payroll.PayrollMonth, &payroll.PayrollStatus,
		&payroll.BasicSalary, &payroll.DAAllowance, &payroll.HRAAllowance, &payroll.SpecialAllowance,
		&payroll.ConveyanceAllow, &payroll.MedicalAllow, &payroll.OtherAllowances, &payroll.TotalEarnings,
		&payroll.EPFDeduction, &payroll.ESIDeduction, &payroll.ProfessionalTax, &payroll.IncomeTax, &payroll.LoanDeduction, &payroll.AdvanceDeduction, &payroll.OtherDeductions, &payroll.TotalDeductions,
		&payroll.NetSalary, &payroll.WorkingDays, &payroll.LeaveDays, &payroll.PaidDays, &payroll.Notes, &payroll.CreatedAt, &payroll.UpdatedAt, &payroll.ProcessedAt, &payroll.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payroll record not found")
	}
	return &payroll, err
}

// ListPayrollRecords retrieves payroll records for an employee
func (s *HRService) ListPayrollRecords(tenantID, employeeID string) ([]models.PayrollRecord, error) {
	var records []models.PayrollRecord

	query := `SELECT id, tenant_id, employee_id, payroll_month, payroll_status,
		basic_salary, dearness_allowance, house_rent_allowance, special_allowance,
		conveyance_allowance, medical_allowance, other_allowances, total_earnings,
		epf_deduction, esi_deduction, professional_tax, income_tax, loan_deduction, advance_deduction, other_deductions, total_deductions,
		net_salary, working_days, leave_days, paid_days, notes, created_at, updated_at, processed_at, deleted_at
		FROM payroll WHERE tenant_id = ? AND employee_id = ? AND deleted_at IS NULL
		ORDER BY payroll_month DESC`

	rows, err := s.DB.Query(query, tenantID, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payroll models.PayrollRecord
		err := rows.Scan(
			&payroll.ID, &payroll.TenantID, &payroll.EmployeeID, &payroll.PayrollMonth, &payroll.PayrollStatus,
			&payroll.BasicSalary, &payroll.DAAllowance, &payroll.HRAAllowance, &payroll.SpecialAllowance,
			&payroll.ConveyanceAllow, &payroll.MedicalAllow, &payroll.OtherAllowances, &payroll.TotalEarnings,
			&payroll.EPFDeduction, &payroll.ESIDeduction, &payroll.ProfessionalTax, &payroll.IncomeTax, &payroll.LoanDeduction, &payroll.AdvanceDeduction, &payroll.OtherDeductions, &payroll.TotalDeductions,
			&payroll.NetSalary, &payroll.WorkingDays, &payroll.LeaveDays, &payroll.PaidDays, &payroll.Notes, &payroll.CreatedAt, &payroll.UpdatedAt, &payroll.ProcessedAt, &payroll.DeletedAt,
		)
		if err != nil {
			log.Printf("Error scanning payroll: %v", err)
			continue
		}
		records = append(records, payroll)
	}

	return records, rows.Err()
}

// ============================================================================
// LEAVE MANAGEMENT
// ============================================================================

// RequestLeave creates a leave request
func (s *HRService) RequestLeave(tenantID string, leave *models.LeaveRequest) error {
	leave.TenantID = tenantID
	leave.Status = "pending"
	leave.CreatedAt = time.Now()
	leave.UpdatedAt = time.Now()

	query := `INSERT INTO leave_requests (id, tenant_id, employee_id, leave_type_id, from_date, to_date, number_of_days, reason, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.Exec(query, leave.ID, leave.TenantID, leave.EmployeeID, leave.LeaveTypeID, leave.FromDate, leave.ToDate, leave.NumberOfDays, leave.Reason, leave.Status, leave.CreatedAt, leave.UpdatedAt)
	return err
}

// ApproveLeave approves a leave request
func (s *HRService) ApproveLeave(tenantID, leaveID, approvedBy string) error {
	query := `UPDATE leave_requests SET status = 'approved', approved_by = ?, approval_date = NOW(), updated_at = NOW()
		WHERE id = ? AND tenant_id = ?`
	_, err := s.DB.Exec(query, approvedBy, leaveID, tenantID)
	return err
}

// RejectLeave rejects a leave request
func (s *HRService) RejectLeave(tenantID, leaveID, reason string) error {
	query := `UPDATE leave_requests SET status = 'rejected', rejection_reason = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ?`
	_, err := s.DB.Exec(query, reason, leaveID, tenantID)
	return err
}

// GetLeaveBalance retrieves leave balance for an employee
func (s *HRService) GetLeaveBalance(tenantID, employeeID string) (map[string]int, error) {
	balance := make(map[string]int)

	query := `SELECT lt.leave_type_name, lt.annual_entitlement, COUNT(lr.id) as used_leaves
		FROM leave_types lt
		LEFT JOIN leave_requests lr ON lr.leave_type_id = lt.id AND lr.employee_id = ? AND lr.status = 'approved'
		WHERE lt.tenant_id = ?
		GROUP BY lt.id, lt.leave_type_name, lt.annual_entitlement`

	rows, err := s.DB.Query(query, employeeID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var leaveType string
		var entitlement, used int
		if err := rows.Scan(&leaveType, &entitlement, &used); err != nil {
			log.Printf("Error scanning leave balance: %v", err)
			continue
		}
		balance[leaveType] = entitlement - used
	}

	return balance, rows.Err()
}

// ============================================================================
// PAYROLL TO GL INTEGRATION
// ============================================================================

// PostPayrollToGL posts payroll entries to the General Ledger
// This creates debit/credit entries for:
// - Salary Expense (Debit)
// - Salary Payable (Credit)
// - Tax Payable (Credit)
// - Deduction Payables (Credit)
func (s *HRService) PostPayrollToGL(tenantID, payrollID string, glService *GLService, postedBy string) (string, error) {
	// Get payroll record
	payroll, err := s.GetPayrollRecord(tenantID, payrollID)
	if err != nil {
		return "", fmt.Errorf("failed to get payroll record: %w", err)
	}

	// Create journal entry header for payroll posting
	journalEntry := &models.JournalEntry{
		ID:              fmt.Sprintf("JE-HR-PAYROLL-%s", payrollID),
		TenantID:        tenantID,
		EntryDate:       time.Now(),
		ReferenceNumber: &payrollID,
		ReferenceType:   "HR_Payroll",
		ReferenceID:     &payrollID,
		Description:     fmt.Sprintf("Salary accrual for %s", payroll.PayrollMonth.Format("Jan 2006")),
		Amount:          payroll.TotalEarnings,
		Narration:       fmt.Sprintf("Monthly salary expense for employee %s", payroll.EmployeeID),
		EntryStatus:     "Draft",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Create journal entry in GL service
	if err := glService.CreateJournalEntry(tenantID, journalEntry); err != nil {
		return "", fmt.Errorf("failed to create journal entry: %w", err)
	}

	// Add debit line: Salary Expense account
	expenseDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-EXP", payrollID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-SALARY-EXPENSE", // Should be configured per tenant
		DebitAmount:    payroll.TotalEarnings,
		CreditAmount:   0,
		Description:    "Salary expense",
		LineNumber:     1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(expenseDetail); err != nil {
		return "", fmt.Errorf("failed to add expense detail: %w", err)
	}

	// Add credit line: Salary Payable account
	payableDetail := &models.JournalEntryDetail{
		ID:             fmt.Sprintf("JED-%s-PAY", payrollID),
		TenantID:       tenantID,
		JournalEntryID: journalEntry.ID,
		AccountID:      "ACC-SALARY-PAYABLE", // Should be configured per tenant
		DebitAmount:    0,
		CreditAmount:   payroll.NetSalary, // Net salary = Earnings - All deductions (employee receives this)
		Description:    "Net salary payable to employee",
		LineNumber:     2,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := glService.AddJournalEntryDetail(payableDetail); err != nil {
		return "", fmt.Errorf("failed to add payable detail: %w", err)
	}

	// Add credit line: Tax Payable if any (company owes government)
	if payroll.ProfessionalTax > 0 || payroll.IncomeTax > 0 {
		taxPayable := payroll.ProfessionalTax + payroll.IncomeTax
		taxDetail := &models.JournalEntryDetail{
			ID:             fmt.Sprintf("JED-%s-TAX", payrollID),
			TenantID:       tenantID,
			JournalEntryID: journalEntry.ID,
			AccountID:      "ACC-TAX-PAYABLE", // Should be configured per tenant
			DebitAmount:    0,
			CreditAmount:   taxPayable,
			Description:    "Income tax payable to government",
			LineNumber:     3,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := glService.AddJournalEntryDetail(taxDetail); err != nil {
			return "", fmt.Errorf("failed to add tax detail: %w", err)
		}
	}

	// Add credit line: EPF Payable if any (company owes EPF authority)
	if payroll.EPFDeduction > 0 {
		epfDetail := &models.JournalEntryDetail{
			ID:             fmt.Sprintf("JED-%s-EPF", payrollID),
			TenantID:       tenantID,
			JournalEntryID: journalEntry.ID,
			AccountID:      "ACC-EPF-PAYABLE", // Should be configured per tenant
			DebitAmount:    0,
			CreditAmount:   payroll.EPFDeduction,
			Description:    "EPF deduction payable to EPF authority",
			LineNumber:     4,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := glService.AddJournalEntryDetail(epfDetail); err != nil {
			return "", fmt.Errorf("failed to add EPF detail: %w", err)
		}
	}

	// Add credit line: ESI Payable if any (company owes ESI authority)
	if payroll.ESIDeduction > 0 {
		esiDetail := &models.JournalEntryDetail{
			ID:             fmt.Sprintf("JED-%s-ESI", payrollID),
			TenantID:       tenantID,
			JournalEntryID: journalEntry.ID,
			AccountID:      "ACC-ESI-PAYABLE", // Should be configured per tenant
			DebitAmount:    0,
			CreditAmount:   payroll.ESIDeduction,
			Description:    "ESI deduction payable to ESI authority",
			LineNumber:     5,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := glService.AddJournalEntryDetail(esiDetail); err != nil {
			return "", fmt.Errorf("failed to add ESI detail: %w", err)
		}
	}

	// Post the journal entry (validates debit=credit balance)
	if err := glService.PostJournalEntry(tenantID, journalEntry.ID, postedBy); err != nil {
		return "", fmt.Errorf("failed to post journal entry: %w", err)
	}

	// Update payroll status to indicate GL posting
	updateQuery := `UPDATE payroll SET payroll_status = 'posted_to_gl', updated_at = ? 
		WHERE id = ? AND tenant_id = ?`
	_, err = s.DB.Exec(updateQuery, time.Now(), payrollID, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to update payroll status: %w", err)
	}

	return journalEntry.ID, nil
}
