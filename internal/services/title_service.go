package services

import (
	"database/sql"
	"errors"
	"time"

	"vyomtech-backend/internal/models"
)

// TitleService handles title clearance operations
type TitleService struct {
	db *sql.DB
}

// NewTitleService creates a new title service
func NewTitleService(db *sql.DB) *TitleService {
	return &TitleService{db: db}
}

// CreateTitleClearance creates a new title clearance
func (ts *TitleService) CreateTitleClearance(tenantID, bookingID int64, clearanceType string, req *models.CreateTitleClearanceRequest) (*models.TitleClearance, error) {
	clearance := &models.TitleClearance{
		TenantID:             tenantID,
		BookingID:            bookingID,
		PropertyID:           req.PropertyID,
		Status:               "pending",
		ClearanceType:        clearanceType,
		TargetClearanceDate:  req.TargetClearanceDate,
		Priority:             req.Priority,
		Notes:                req.Notes,
		ClearancePercentage:  0,
		IssuesCount:          0,
		ResolvedIssuesCount:  0,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	query := `INSERT INTO title_clearances 
		(tenant_id, booking_id, property_id, status, clearance_type, target_clearance_date, 
		 priority, notes, clearance_percentage, issues_count, resolved_issues_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	err := ts.db.QueryRow(query,
		clearance.TenantID, clearance.BookingID, clearance.PropertyID, clearance.Status,
		clearance.ClearanceType, clearance.TargetClearanceDate, clearance.Priority, clearance.Notes,
		clearance.ClearancePercentage, clearance.IssuesCount, clearance.ResolvedIssuesCount,
		clearance.CreatedAt, clearance.UpdatedAt).Scan(&clearance.ID)

	if err != nil {
		return nil, err
	}

	return clearance, nil
}

// GetTitleClearance retrieves a title clearance by ID
func (ts *TitleService) GetTitleClearance(tenantID, clearanceID int64) (*models.TitleClearance, error) {
	clearance := &models.TitleClearance{}

	query := `SELECT id, tenant_id, booking_id, property_id, status, clearance_type, 
		start_date, target_clearance_date, actual_clearance_date, clearance_percentage, 
		issues_count, resolved_issues_count, priority, notes, created_by, updated_by, 
		created_at, updated_at, deleted_at
		FROM title_clearances
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := ts.db.QueryRow(query, clearanceID, tenantID).Scan(
		&clearance.ID, &clearance.TenantID, &clearance.BookingID, &clearance.PropertyID,
		&clearance.Status, &clearance.ClearanceType, &clearance.StartDate,
		&clearance.TargetClearanceDate, &clearance.ActualClearanceDate,
		&clearance.ClearancePercentage, &clearance.IssuesCount, &clearance.ResolvedIssuesCount,
		&clearance.Priority, &clearance.Notes, &clearance.CreatedBy, &clearance.UpdatedBy,
		&clearance.CreatedAt, &clearance.UpdatedAt, &clearance.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("clearance not found")
	}
	if err != nil {
		return nil, err
	}

	return clearance, nil
}

// UpdateTitleClearance updates a title clearance
func (ts *TitleService) UpdateTitleClearance(tenantID, clearanceID int64, req *models.UpdateTitleClearanceRequest) error {
	query := `UPDATE title_clearances SET `
	args := []interface{}{}
	fieldCount := 0

	if req.Status != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "status = ?"
		args = append(args, *req.Status)
		fieldCount++
	}
	if req.TargetClearanceDate != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "target_clearance_date = ?"
		args = append(args, *req.TargetClearanceDate)
		fieldCount++
	}
	if req.ClearancePercentage != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "clearance_percentage = ?"
		args = append(args, *req.ClearancePercentage)
		fieldCount++
	}
	if req.Priority != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "priority = ?"
		args = append(args, *req.Priority)
		fieldCount++
	}
	if req.Notes != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "notes = ?"
		args = append(args, *req.Notes)
		fieldCount++
	}

	if fieldCount == 0 {
		return errors.New("no fields to update")
	}

	query += " updated_at = ? WHERE id = ? AND tenant_id = ?"
	args = append(args, time.Now(), clearanceID, tenantID)

	_, err := ts.db.Exec(query, args...)
	return err
}

// ListTitleClearances lists title clearances with pagination
func (ts *TitleService) ListTitleClearances(tenantID int64, limit, offset int) ([]*models.TitleClearance, int, error) {
	clearances := []*models.TitleClearance{}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM title_clearances WHERE tenant_id = ? AND deleted_at IS NULL`
	var total int
	err := ts.db.QueryRow(countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, booking_id, property_id, status, clearance_type, 
		start_date, target_clearance_date, actual_clearance_date, clearance_percentage, 
		issues_count, resolved_issues_count, priority, notes, created_by, updated_by, 
		created_at, updated_at, deleted_at
		FROM title_clearances
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := ts.db.Query(query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		clearance := &models.TitleClearance{}
		err := rows.Scan(
			&clearance.ID, &clearance.TenantID, &clearance.BookingID, &clearance.PropertyID,
			&clearance.Status, &clearance.ClearanceType, &clearance.StartDate,
			&clearance.TargetClearanceDate, &clearance.ActualClearanceDate,
			&clearance.ClearancePercentage, &clearance.IssuesCount, &clearance.ResolvedIssuesCount,
			&clearance.Priority, &clearance.Notes, &clearance.CreatedBy, &clearance.UpdatedBy,
			&clearance.CreatedAt, &clearance.UpdatedAt, &clearance.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		clearances = append(clearances, clearance)
	}

	return clearances, total, rows.Err()
}

// CreateTitleIssue creates a new title issue
func (ts *TitleService) CreateTitleIssue(tenantID, clearanceID int64, req *models.CreateTitleIssueRequest) (*models.TitleIssue, error) {
	issue := &models.TitleIssue{
		TenantID:        tenantID,
		ClearanceID:     clearanceID,
		IssueType:       req.IssueType,
		IssueTitle:      req.IssueTitle,
		IssueDescription: req.IssueDescription,
		Severity:        req.Severity,
		Status:          "open",
		SourceDocument:  req.SourceDocument,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	query := `INSERT INTO title_issues 
		(tenant_id, clearance_id, issue_type, issue_title, issue_description, severity, status, 
		 source_document, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	err := ts.db.QueryRow(query,
		issue.TenantID, issue.ClearanceID, issue.IssueType, issue.IssueTitle,
		issue.IssueDescription, issue.Severity, issue.Status, issue.SourceDocument,
		issue.CreatedAt, issue.UpdatedAt).Scan(&issue.ID)

	if err != nil {
		return nil, err
	}

	// Update issue count in clearance
	ts.updateClearanceIssueCount(tenantID, clearanceID)

	return issue, nil
}

// GetTitleIssue retrieves a title issue by ID
func (ts *TitleService) GetTitleIssue(tenantID, issueID int64) (*models.TitleIssue, error) {
	issue := &models.TitleIssue{}

	query := `SELECT id, tenant_id, clearance_id, issue_type, issue_title, issue_description, 
		severity, status, reported_date, source_document, affected_parties, resolution_notes, 
		resolved_date, resolved_by, resolution_method, metadata, created_by, updated_by, 
		created_at, updated_at, deleted_at
		FROM title_issues
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := ts.db.QueryRow(query, issueID, tenantID).Scan(
		&issue.ID, &issue.TenantID, &issue.ClearanceID, &issue.IssueType, &issue.IssueTitle,
		&issue.IssueDescription, &issue.Severity, &issue.Status, &issue.ReportedDate,
		&issue.SourceDocument, &issue.AffectedParties, &issue.ResolutionNotes,
		&issue.ResolvedDate, &issue.ResolvedBy, &issue.ResolutionMethod, &issue.Metadata,
		&issue.CreatedBy, &issue.UpdatedBy, &issue.CreatedAt, &issue.UpdatedAt, &issue.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("issue not found")
	}
	if err != nil {
		return nil, err
	}

	return issue, nil
}

// UpdateTitleIssue updates a title issue
func (ts *TitleService) UpdateTitleIssue(tenantID, issueID int64, req *models.UpdateTitleIssueRequest) error {
	query := `UPDATE title_issues SET `
	args := []interface{}{}
	fieldCount := 0

	if req.Status != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "status = ?"
		args = append(args, *req.Status)
		fieldCount++
	}
	if req.Severity != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "severity = ?"
		args = append(args, *req.Severity)
		fieldCount++
	}
	if req.ResolutionNotes != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "resolution_notes = ?"
		args = append(args, *req.ResolutionNotes)
		fieldCount++
	}
	if req.ResolutionMethod != nil {
		if fieldCount > 0 {
			query += ", "
		}
		query += "resolution_method = ?"
		args = append(args, *req.ResolutionMethod)
		fieldCount++
	}

	if fieldCount == 0 {
		return errors.New("no fields to update")
	}

	query += " updated_at = ? WHERE id = ? AND tenant_id = ?"
	args = append(args, time.Now(), issueID, tenantID)

	_, err := ts.db.Exec(query, args...)
	return err
}

// ResolveTitleIssue resolves a title issue
func (ts *TitleService) ResolveTitleIssue(tenantID, issueID, userID int64, req *models.ResolveTitleIssueRequest) error {
	now := time.Now()
	query := `UPDATE title_issues SET status = ?, resolution_method = ?, resolved_date = ?, 
		resolved_by = ?, updated_at = ? WHERE id = ? AND tenant_id = ?`

	_, err := ts.db.Exec(query, req.Status, req.ResolutionMethod, now, userID, now, issueID, tenantID)
	return err
}

// ListTitleIssues lists title issues for a clearance
func (ts *TitleService) ListTitleIssues(tenantID, clearanceID int64, limit, offset int) ([]*models.TitleIssue, int, error) {
	issues := []*models.TitleIssue{}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM title_issues WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`
	var total int
	err := ts.db.QueryRow(countQuery, tenantID, clearanceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, clearance_id, issue_type, issue_title, issue_description, 
		severity, status, reported_date, source_document, affected_parties, resolution_notes, 
		resolved_date, resolved_by, resolution_method, metadata, created_by, updated_by, 
		created_at, updated_at, deleted_at
		FROM title_issues
		WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL
		ORDER BY severity DESC, created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := ts.db.Query(query, tenantID, clearanceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		issue := &models.TitleIssue{}
		err := rows.Scan(
			&issue.ID, &issue.TenantID, &issue.ClearanceID, &issue.IssueType, &issue.IssueTitle,
			&issue.IssueDescription, &issue.Severity, &issue.Status, &issue.ReportedDate,
			&issue.SourceDocument, &issue.AffectedParties, &issue.ResolutionNotes,
			&issue.ResolvedDate, &issue.ResolvedBy, &issue.ResolutionMethod, &issue.Metadata,
			&issue.CreatedBy, &issue.UpdatedBy, &issue.CreatedAt, &issue.UpdatedAt, &issue.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		issues = append(issues, issue)
	}

	return issues, total, rows.Err()
}

// CreateTitleSearchReport creates a new search report
func (ts *TitleService) CreateTitleSearchReport(tenantID, clearanceID int64, req *models.CreateSearchReportRequest) (*models.TitleSearchReport, error) {
	report := &models.TitleSearchReport{
		TenantID:       tenantID,
		ClearanceID:    clearanceID,
		SearchType:     req.SearchType,
		SearchAuthority: req.SearchAuthority,
		SearchReferenceNum: req.SearchReferenceNum,
		SearchCost:     *req.SearchCost,
		SearchStatus:   "pending",
		SearchDate:     now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO title_search_reports 
		(tenant_id, clearance_id, search_type, search_authority, search_reference_number, 
		 search_cost, search_status, search_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	err := ts.db.QueryRow(query,
		report.TenantID, report.ClearanceID, report.SearchType, report.SearchAuthority,
		report.SearchReferenceNum, report.SearchCost, report.SearchStatus, report.SearchDate,
		report.CreatedAt, report.UpdatedAt).Scan(&report.ID)

	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetTitleSearchReport retrieves a search report
func (ts *TitleService) GetTitleSearchReport(tenantID, reportID int64) (*models.TitleSearchReport, error) {
	report := &models.TitleSearchReport{}

	query := `SELECT id, tenant_id, clearance_id, search_type, search_date, search_authority, 
		search_reference_number, report_url, report_file_name, report_file_size, s3_bucket, s3_key, 
		encumbrances_found, search_status, verified_by, verified_at, verification_notes, search_cost, 
		metadata, created_by, updated_by, created_at, updated_at, deleted_at
		FROM title_search_reports
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := ts.db.QueryRow(query, reportID, tenantID).Scan(
		&report.ID, &report.TenantID, &report.ClearanceID, &report.SearchType, &report.SearchDate,
		&report.SearchAuthority, &report.SearchReferenceNum, &report.ReportURL, &report.ReportFileName,
		&report.ReportFileSize, &report.S3Bucket, &report.S3Key, &report.EncumbrancesFound,
		&report.SearchStatus, &report.VerifiedBy, &report.VerifiedAt, &report.VerificationNotes,
		&report.SearchCost, &report.Metadata, &report.CreatedBy, &report.UpdatedBy,
		&report.CreatedAt, &report.UpdatedAt, &report.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("search report not found")
	}
	if err != nil {
		return nil, err
	}

	return report, nil
}

// VerifySearchReport verifies a search report
func (ts *TitleService) VerifySearchReport(tenantID, reportID, userID int64, req *models.VerifySearchReportRequest) error {
	now := time.Now()
	query := `UPDATE title_search_reports SET search_status = ?, encumbrances_found = ?, 
		verification_notes = ?, verified_by = ?, verified_at = ?, updated_at = ? 
		WHERE id = ? AND tenant_id = ?`

	_, err := ts.db.Exec(query, req.SearchStatus, req.EncumbrancesFound, req.VerificationNotes,
		userID, now, now, reportID, tenantID)
	return err
}

// ListTitleSearchReports lists search reports
func (ts *TitleService) ListTitleSearchReports(tenantID, clearanceID int64, limit, offset int) ([]*models.TitleSearchReport, int, error) {
	reports := []*models.TitleSearchReport{}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM title_search_reports WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`
	var total int
	err := ts.db.QueryRow(countQuery, tenantID, clearanceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, clearance_id, search_type, search_date, search_authority, 
		search_reference_number, report_url, report_file_name, report_file_size, s3_bucket, s3_key, 
		encumbrances_found, search_status, verified_by, verified_at, verification_notes, search_cost, 
		metadata, created_by, updated_by, created_at, updated_at, deleted_at
		FROM title_search_reports
		WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := ts.db.Query(query, tenantID, clearanceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		report := &models.TitleSearchReport{}
		err := rows.Scan(
			&report.ID, &report.TenantID, &report.ClearanceID, &report.SearchType, &report.SearchDate,
			&report.SearchAuthority, &report.SearchReferenceNum, &report.ReportURL, &report.ReportFileName,
			&report.ReportFileSize, &report.S3Bucket, &report.S3Key, &report.EncumbrancesFound,
			&report.SearchStatus, &report.VerifiedBy, &report.VerifiedAt, &report.VerificationNotes,
			&report.SearchCost, &report.Metadata, &report.CreatedBy, &report.UpdatedBy,
			&report.CreatedAt, &report.UpdatedAt, &report.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		reports = append(reports, report)
	}

	return reports, total, rows.Err()
}

// CreateTitleLegalOpinion creates a legal opinion
func (ts *TitleService) CreateTitleLegalOpinion(tenantID, clearanceID int64, req *models.CreateLegalOpinionRequest) (*models.TitleLegalOpinion, error) {
	opinion := &models.TitleLegalOpinion{
		TenantID:            tenantID,
		ClearanceID:         clearanceID,
		OpinionType:         req.OpinionType,
		ExpertName:          req.ExpertName,
		ExpertOrganization:  req.ExpertOrganization,
		ExpertLicenseNumber: req.ExpertLicenseNumber,
		OpinionStatus:       "pending",
		Cost:                0,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	query := `INSERT INTO title_legal_opinions 
		(tenant_id, clearance_id, opinion_type, expert_name, expert_organization, 
		 expert_license_number, opinion_status, cost, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	err := ts.db.QueryRow(query,
		opinion.TenantID, opinion.ClearanceID, opinion.OpinionType, opinion.ExpertName,
		opinion.ExpertOrganization, opinion.ExpertLicenseNumber, opinion.OpinionStatus,
		opinion.Cost, opinion.CreatedAt, opinion.UpdatedAt).Scan(&opinion.ID)

	if err != nil {
		return nil, err
	}

	return opinion, nil
}

// GetTitleLegalOpinion retrieves a legal opinion
func (ts *TitleService) GetTitleLegalOpinion(tenantID, opinionID int64) (*models.TitleLegalOpinion, error) {
	opinion := &models.TitleLegalOpinion{}

	query := `SELECT id, tenant_id, clearance_id, opinion_type, expert_name, expert_organization, 
		expert_license_number, opinion_date, opinion_status, opinion_url, opinion_file_name, file_size, 
		s3_bucket, s3_key, opinion_summary, recommendations, risk_assessment, review_by_lawyer, 
		review_notes, reviewed_at, cost, metadata, created_by, updated_by, created_at, updated_at, deleted_at
		FROM title_legal_opinions
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	err := ts.db.QueryRow(query, opinionID, tenantID).Scan(
		&opinion.ID, &opinion.TenantID, &opinion.ClearanceID, &opinion.OpinionType, &opinion.ExpertName,
		&opinion.ExpertOrganization, &opinion.ExpertLicenseNumber, &opinion.OpinionDate, &opinion.OpinionStatus,
		&opinion.OpinionURL, &opinion.OpinionFileName, &opinion.FileSize, &opinion.S3Bucket, &opinion.S3Key,
		&opinion.OpinionSummary, &opinion.Recommendations, &opinion.RiskAssessment, &opinion.ReviewByLawyer,
		&opinion.ReviewNotes, &opinion.ReviewedAt, &opinion.Cost, &opinion.Metadata, &opinion.CreatedBy,
		&opinion.UpdatedBy, &opinion.CreatedAt, &opinion.UpdatedAt, &opinion.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("legal opinion not found")
	}
	if err != nil {
		return nil, err
	}

	return opinion, nil
}

// ReviewLegalOpinion reviews a legal opinion
func (ts *TitleService) ReviewLegalOpinion(tenantID, opinionID, userID int64, req *models.ReviewLegalOpinionRequest) error {
	now := time.Now()
	query := `UPDATE title_legal_opinions SET opinion_status = ?, risk_assessment = ?, 
		review_notes = ?, review_by_lawyer = ?, reviewed_at = ?, updated_at = ? 
		WHERE id = ? AND tenant_id = ?`

	_, err := ts.db.Exec(query, req.OpinionStatus, req.RiskAssessment, req.ReviewNotes,
		userID, now, now, opinionID, tenantID)
	return err
}

// ListTitleLegalOpinions lists legal opinions
func (ts *TitleService) ListTitleLegalOpinions(tenantID, clearanceID int64, limit, offset int) ([]*models.TitleLegalOpinion, int, error) {
	opinions := []*models.TitleLegalOpinion{}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM title_legal_opinions WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`
	var total int
	err := ts.db.QueryRow(countQuery, tenantID, clearanceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, clearance_id, opinion_type, expert_name, expert_organization, 
		expert_license_number, opinion_date, opinion_status, opinion_url, opinion_file_name, file_size, 
		s3_bucket, s3_key, opinion_summary, recommendations, risk_assessment, review_by_lawyer, 
		review_notes, reviewed_at, cost, metadata, created_by, updated_by, created_at, updated_at, deleted_at
		FROM title_legal_opinions
		WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := ts.db.Query(query, tenantID, clearanceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		opinion := &models.TitleLegalOpinion{}
		err := rows.Scan(
			&opinion.ID, &opinion.TenantID, &opinion.ClearanceID, &opinion.OpinionType, &opinion.ExpertName,
			&opinion.ExpertOrganization, &opinion.ExpertLicenseNumber, &opinion.OpinionDate, &opinion.OpinionStatus,
			&opinion.OpinionURL, &opinion.OpinionFileName, &opinion.FileSize, &opinion.S3Bucket, &opinion.S3Key,
			&opinion.OpinionSummary, &opinion.Recommendations, &opinion.RiskAssessment, &opinion.ReviewByLawyer,
			&opinion.ReviewNotes, &opinion.ReviewedAt, &opinion.Cost, &opinion.Metadata, &opinion.CreatedBy,
			&opinion.UpdatedBy, &opinion.CreatedAt, &opinion.UpdatedAt, &opinion.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		opinions = append(opinions, opinion)
	}

	return opinions, total, rows.Err()
}

// CreateVerificationChecklist creates a checklist item
func (ts *TitleService) CreateVerificationChecklist(tenantID, clearanceID int64, itemName string) (*models.TitleVerificationChecklist, error) {
	item := &models.TitleVerificationChecklist{
		TenantID:    tenantID,
		ClearanceID: clearanceID,
		ItemName:    itemName,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	query := `INSERT INTO title_verification_checklists 
		(tenant_id, clearance_id, item_name, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id`

	err := ts.db.QueryRow(query,
		item.TenantID, item.ClearanceID, item.ItemName, item.Status,
		item.CreatedAt, item.UpdatedAt).Scan(&item.ID)

	if err != nil {
		return nil, err
	}

	return item, nil
}

// VerifyChecklistItem marks a checklist item as verified
func (ts *TitleService) VerifyChecklistItem(tenantID, checklistID, userID int64, verificationNotes *string) error {
	now := time.Now()
	query := `UPDATE title_verification_checklists SET status = ?, verified_by = ?, 
		verified_at = ?, verification_notes = ?, updated_at = ? WHERE id = ? AND tenant_id = ?`

	_, err := ts.db.Exec(query, "verified", userID, now, verificationNotes, now, checklistID, tenantID)
	return err
}

// ListVerificationChecklists lists checklist items
func (ts *TitleService) ListVerificationChecklists(tenantID, clearanceID int64, limit, offset int) ([]*models.TitleVerificationChecklist, int, error) {
	items := []*models.TitleVerificationChecklist{}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM title_verification_checklists WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`
	var total int
	err := ts.db.QueryRow(countQuery, tenantID, clearanceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, clearance_id, item_name, item_category, description, is_mandatory, 
		status, verified_by, verified_at, verification_notes, sequence_order, metadata, created_by, 
		updated_by, created_at, updated_at, deleted_at
		FROM title_verification_checklists
		WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL
		ORDER BY sequence_order ASC
		LIMIT ? OFFSET ?`

	rows, err := ts.db.Query(query, tenantID, clearanceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.TitleVerificationChecklist{}
		err := rows.Scan(
			&item.ID, &item.TenantID, &item.ClearanceID, &item.ItemName, &item.ItemCategory,
			&item.Description, &item.IsMandatory, &item.Status, &item.VerifiedBy, &item.VerifiedAt,
			&item.VerificationNotes, &item.SequenceOrder, &item.Metadata, &item.CreatedBy,
			&item.UpdatedBy, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, rows.Err()
}

// CreateClearanceApproval creates an approval
func (ts *TitleService) CreateClearanceApproval(tenantID, clearanceID int64, req *models.CreateClearanceApprovalRequest) (*models.TitleClearanceApproval, error) {
	approval := &models.TitleClearanceApproval{
		TenantID:        tenantID,
		ClearanceID:     clearanceID,
		ApprovalType:    req.ApprovalType,
		ApproverID:      req.ApproverID,
		ApprovalStatus:  "pending",
		SequenceOrder:   req.SequenceOrder,
		IsFinalApproval: req.IsFinalApproval,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	query := `INSERT INTO title_clearance_approvals 
		(tenant_id, clearance_id, approval_type, approver_id, approval_status, 
		 sequence_order, is_final_approval, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id`

	err := ts.db.QueryRow(query,
		approval.TenantID, approval.ClearanceID, approval.ApprovalType, approval.ApproverID,
		approval.ApprovalStatus, approval.SequenceOrder, approval.IsFinalApproval,
		approval.CreatedAt, approval.UpdatedAt).Scan(&approval.ID)

	if err != nil {
		return nil, err
	}

	return approval, nil
}

// ApproveClearance approves a clearance
func (ts *TitleService) ApproveClearance(tenantID, approvalID int64, req *models.ApproveClearanceRequest) error {
	now := time.Now()
	query := `UPDATE title_clearance_approvals SET approval_status = ?, approval_notes = ?, 
		conditional_requirements = ?, is_final_approval = ?, approval_date = ?, updated_at = ? 
		WHERE id = ? AND tenant_id = ?`

	_, err := ts.db.Exec(query, req.ApprovalStatus, req.ApprovalNotes,
		req.ConditionalRequirements, req.IsFinalApproval, now, now, approvalID, tenantID)
	return err
}

// ListClearanceApprovals lists approvals
func (ts *TitleService) ListClearanceApprovals(tenantID, clearanceID int64, limit, offset int) ([]*models.TitleClearanceApproval, int, error) {
	approvals := []*models.TitleClearanceApproval{}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM title_clearance_approvals WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`
	var total int
	err := ts.db.QueryRow(countQuery, tenantID, clearanceID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, tenant_id, clearance_id, approval_type, approver_id, approver_role, 
		approval_status, approval_notes, conditional_requirements, approval_date, valid_from, 
		valid_till, sequence_order, is_final_approval, metadata, created_by, updated_by, 
		created_at, updated_at, deleted_at
		FROM title_clearance_approvals
		WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL
		ORDER BY sequence_order ASC
		LIMIT ? OFFSET ?`

	rows, err := ts.db.Query(query, tenantID, clearanceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		approval := &models.TitleClearanceApproval{}
		err := rows.Scan(
			&approval.ID, &approval.TenantID, &approval.ClearanceID, &approval.ApprovalType,
			&approval.ApproverID, &approval.ApproverRole, &approval.ApprovalStatus, &approval.ApprovalNotes,
			&approval.ConditionalRequirements, &approval.ApprovalDate, &approval.ValidFrom, &approval.ValidTill,
			&approval.SequenceOrder, &approval.IsFinalApproval, &approval.Metadata, &approval.CreatedBy,
			&approval.UpdatedBy, &approval.CreatedAt, &approval.UpdatedAt, &approval.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		approvals = append(approvals, approval)
	}

	return approvals, total, rows.Err()
}

// GetClearanceSummary returns a summary of the clearance status
func (ts *TitleService) GetClearanceSummary(tenantID, clearanceID int64) (*models.TitleClearanceSummaryResponse, error) {
	clearance, err := ts.GetTitleClearance(tenantID, clearanceID)
	if err != nil {
		return nil, err
	}

	// Count issues
	var openIssues, resolvedIssues int
	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_issues WHERE tenant_id = ? AND clearance_id = ? AND status = 'open' AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&openIssues)
	if err != nil {
		return nil, err
	}

	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_issues WHERE tenant_id = ? AND clearance_id = ? AND status = 'resolved' AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&resolvedIssues)
	if err != nil {
		return nil, err
	}

	// Count search reports
	var searchReportCount int
	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_search_reports WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&searchReportCount)
	if err != nil {
		return nil, err
	}

	// Count legal opinions
	var legalOpinionCount int
	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_legal_opinions WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&legalOpinionCount)
	if err != nil {
		return nil, err
	}

	// Count pending approvals
	var pendingApprovals int
	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_clearance_approvals WHERE tenant_id = ? AND clearance_id = ? AND approval_status = 'pending' AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&pendingApprovals)
	if err != nil {
		return nil, err
	}

	// Count checklist items
	var checklistItems, verifiedItems int
	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_verification_checklists WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&checklistItems)
	if err != nil {
		return nil, err
	}

	err = ts.db.QueryRow(`SELECT COUNT(*) FROM title_verification_checklists WHERE tenant_id = ? AND clearance_id = ? AND status = 'verified' AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&verifiedItems)
	if err != nil {
		return nil, err
	}

	return &models.TitleClearanceSummaryResponse{
		ClearanceID:            clearance.ID,
		BookingID:              clearance.BookingID,
		Status:                 clearance.Status,
		ClearanceType:          clearance.ClearanceType,
		ClearancePercentage:    clearance.ClearancePercentage,
		TargetClearanceDate:    clearance.TargetClearanceDate,
		ActualClearanceDate:    clearance.ActualClearanceDate,
		TotalIssues:            openIssues + resolvedIssues,
		OpenIssues:             openIssues,
		ResolvedIssues:         resolvedIssues,
		SearchReports:          searchReportCount,
		LegalOpinions:          legalOpinionCount,
		PendingApprovals:       pendingApprovals,
		ChecklistItems:         checklistItems,
		VerifiedChecklistItems: verifiedItems,
		CreatedAt:              clearance.CreatedAt,
		UpdatedAt:              clearance.UpdatedAt,
	}, nil
}

// Helper functions

func (ts *TitleService) updateClearanceIssueCount(tenantID, clearanceID int64) error {
	var count int
	err := ts.db.QueryRow(`SELECT COUNT(*) FROM title_issues WHERE tenant_id = ? AND clearance_id = ? AND deleted_at IS NULL`, tenantID, clearanceID).Scan(&count)
	if err != nil {
		return err
	}

	_, err = ts.db.Exec(`UPDATE title_clearances SET issues_count = ? WHERE id = ? AND tenant_id = ?`, count, clearanceID, tenantID)
	return err
}

func now() *time.Time {
	t := time.Now()
	return &t
}
