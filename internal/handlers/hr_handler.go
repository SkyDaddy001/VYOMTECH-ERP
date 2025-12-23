package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vyomtech-backend/internal/constants"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/models"
	"vyomtech-backend/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HRHandler struct {
	Service     *services.HRService
	RBACService *services.RBACService
}

// NewHRHandler creates a new HR handler
func NewHRHandler(service *services.HRService, rbacService *services.RBACService) *HRHandler {
	return &HRHandler{
		Service:     service,
		RBACService: rbacService,
	}
}

// ============================================================================
// EMPLOYEE ENDPOINTS
// ============================================================================

// CreateEmployee - POST /api/v1/hr/employees
func (h *HRHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.EmployeeCreate); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	emp.ID = uuid.New().String()
	emp.Status = "active"

	if err := h.Service.CreateEmployee(tenant, &emp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create employee: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

// GetEmployee - GET /api/v1/hr/employees/{id}
func (h *HRHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	emp, err := h.Service.GetEmployee(tenant, employeeID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Employee not found: %s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

// ListEmployees - GET /api/v1/hr/employees
func (h *HRHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get("X-Tenant-ID")

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed > 0 {
			offset = parsed
		}
	}

	employees, total, err := h.Service.ListEmployees(tenant, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to list employees: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employees": employees,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateEmployee - PUT /api/v1/hr/employees/{id}
func (h *HRHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.EmployeeUpdate); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	employeeID := vars["id"]

	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	emp.ID = employeeID
	emp.TenantID = tenant

	if err := h.Service.UpdateEmployee(tenant, &emp); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to update employee: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

// DeleteEmployee - DELETE /api/v1/hr/employees/{id}
func (h *HRHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract user and tenant from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, `{"error": "User ID not found in context"}`, http.StatusUnauthorized)
		return
	}

	tenant, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok || tenant == "" {
		http.Error(w, `{"error": "Tenant ID not found in context"}`, http.StatusForbidden)
		return
	}

	// Verify permission
	if err := h.RBACService.VerifyPermission(r.Context(), tenant, userID, constants.EmployeeDelete); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Permission denied: %s"}`, err.Error()), http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	employeeID := vars["id"]

	if err := h.Service.DeleteEmployee(tenant, employeeID); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to delete employee: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Employee deleted successfully"})
}

// ============================================================================
// ATTENDANCE ENDPOINTS
// ============================================================================

// RecordAttendance - POST /api/v1/hr/attendance
func (h *HRHandler) RecordAttendance(w http.ResponseWriter, r *http.Request) {
	var att models.Attendance
	if err := json.NewDecoder(r.Body).Decode(&att); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("X-Tenant-ID")
	att.ID = uuid.New().String()

	if err := h.Service.RecordAttendance(tenant, &att); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to record attendance: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(att)
}

// GetAttendanceRecord - GET /api/v1/hr/attendance/{employee_id}/{date}
func (h *HRHandler) GetAttendanceRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employee_id"]
	dateStr := vars["date"]
	tenant := r.Header.Get("X-Tenant-ID")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
		return
	}

	att, err := h.Service.GetAttendanceRecord(tenant, employeeID, date)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Attendance record not found: %s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(att)
}

// ListEmployeeAttendance - GET /api/v1/hr/attendance/{employee_id}
func (h *HRHandler) ListEmployeeAttendance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employee_id"]
	tenant := r.Header.Get("X-Tenant-ID")

	fromDate := r.URL.Query().Get("from_date")
	toDate := r.URL.Query().Get("to_date")

	from, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		from = time.Now().AddDate(0, 0, -30)
	}

	to, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		to = time.Now()
	}

	records, err := h.Service.ListEmployeeAttendance(tenant, employeeID, from, to)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch attendance: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"records": records,
		"from":    from.Format("2006-01-02"),
		"to":      to.Format("2006-01-02"),
	})
}

// ============================================================================
// PAYROLL ENDPOINTS
// ============================================================================

// GeneratePayroll - POST /api/v1/hr/payroll/generate
func (h *HRHandler) GeneratePayroll(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID   string `json:"employee_id"`
		PayrollMonth string `json:"payroll_month"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("X-Tenant-ID")

	payrollMonth, err := time.Parse("2006-01-02", req.PayrollMonth)
	if err != nil {
		payrollMonth, _ = time.Parse("2006-01", req.PayrollMonth)
	}

	payroll, err := h.Service.CalculateAndCreatePayroll(tenant, req.EmployeeID, payrollMonth)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to generate payroll: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payroll)
}

// GetPayrollRecord - GET /api/v1/hr/payroll/{id}
func (h *HRHandler) GetPayrollRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	payrollID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	payroll, err := h.Service.GetPayrollRecord(tenant, payrollID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Payroll record not found: %s"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payroll)
}

// ListPayrollRecords - GET /api/v1/hr/payroll/{employee_id}
func (h *HRHandler) ListPayrollRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employee_id"]
	tenant := r.Header.Get("X-Tenant-ID")

	records, err := h.Service.ListPayrollRecords(tenant, employeeID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch payroll records: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"payroll_records": records,
		"total":           len(records),
	})
}

// ============================================================================
// LEAVE ENDPOINTS
// ============================================================================

// RequestLeave - POST /api/v1/hr/leaves
func (h *HRHandler) RequestLeave(w http.ResponseWriter, r *http.Request) {
	var leave models.LeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&leave); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	tenant := r.Header.Get("X-Tenant-ID")
	leave.ID = uuid.New().String()

	if err := h.Service.RequestLeave(tenant, &leave); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to request leave: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(leave)
}

// ApproveLeave - POST /api/v1/hr/leaves/{id}/approve
func (h *HRHandler) ApproveLeave(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leaveID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	var req struct {
		ApprovedBy string `json:"approved_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.Service.ApproveLeave(tenant, leaveID, req.ApprovedBy); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to approve leave: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Leave approved successfully"})
}

// RejectLeave - POST /api/v1/hr/leaves/{id}/reject
func (h *HRHandler) RejectLeave(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leaveID := vars["id"]
	tenant := r.Header.Get("X-Tenant-ID")

	var req struct {
		RejectionReason string `json:"rejection_reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := h.Service.RejectLeave(tenant, leaveID, req.RejectionReason); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to reject leave: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Leave rejected successfully"})
}

// GetLeaveBalance - GET /api/v1/hr/leave-balance/{employee_id}
func (h *HRHandler) GetLeaveBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["employee_id"]
	tenant := r.Header.Get("X-Tenant-ID")

	balance, err := h.Service.GetLeaveBalance(tenant, employeeID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to fetch leave balance: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employee_id": employeeID,
		"balance":     balance,
	})
}

// ============================================================================
// ROUTE REGISTRATION
// ============================================================================

// RegisterHRRoutes registers all HR routes
func RegisterHRRoutes(r *mux.Router, hrService *services.HRService, rbacService *services.RBACService) {
	handler := NewHRHandler(hrService, rbacService)

	// Employee routes
	r.HandleFunc("/api/v1/hr/employees", handler.CreateEmployee).Methods("POST")
	r.HandleFunc("/api/v1/hr/employees", handler.ListEmployees).Methods("GET")
	r.HandleFunc("/api/v1/hr/employees/{id}", handler.GetEmployee).Methods("GET")
	r.HandleFunc("/api/v1/hr/employees/{id}", handler.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/api/v1/hr/employees/{id}", handler.DeleteEmployee).Methods("DELETE")

	// Attendance routes
	r.HandleFunc("/api/v1/hr/attendance", handler.RecordAttendance).Methods("POST")
	r.HandleFunc("/api/v1/hr/attendance/{employee_id}/{date}", handler.GetAttendanceRecord).Methods("GET")
	r.HandleFunc("/api/v1/hr/attendance/{employee_id}", handler.ListEmployeeAttendance).Methods("GET")

	// Payroll routes
	r.HandleFunc("/api/v1/hr/payroll/generate", handler.GeneratePayroll).Methods("POST")
	r.HandleFunc("/api/v1/hr/payroll/{id}", handler.GetPayrollRecord).Methods("GET")
	r.HandleFunc("/api/v1/hr/payroll/{employee_id}", handler.ListPayrollRecords).Methods("GET")

	// Leave routes
	r.HandleFunc("/api/v1/hr/leaves", handler.RequestLeave).Methods("POST")
	r.HandleFunc("/api/v1/hr/leaves/{id}/approve", handler.ApproveLeave).Methods("POST")
	r.HandleFunc("/api/v1/hr/leaves/{id}/reject", handler.RejectLeave).Methods("POST")
	r.HandleFunc("/api/v1/hr/leave-balance/{employee_id}", handler.GetLeaveBalance).Methods("GET")
}
