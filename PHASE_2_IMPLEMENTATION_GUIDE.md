# Phase 2: Handler-Level Permission Checks Implementation Guide

## Overview
Phase 1 (Route-level permission middleware) is COMPLETE. All routes are now protected with role-based access control.

Phase 2 adds granular permission checks at the handler level to enforce specific operations.

## Pattern to Apply

### Step 1: Modify Handler Constructor
```go
type SalesHandler struct {
	DB           *sql.DB
	RBACService  *services.RBACService  // ADD THIS
}

// Update constructor
func NewSalesHandler(db *sql.DB, rbacService *services.RBACService) *SalesHandler {
	return &SalesHandler{
		DB:          db,
		RBACService: rbacService,
	}
}
```

### Step 2: Add Permission Check at Start of Each CRUD Operation
```go
func (h *SalesHandler) CreateSalesLead(w http.ResponseWriter, r *http.Request) {
	// ADD THESE LINES AT START
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok {
		h.respondError(w, http.StatusForbidden, "Tenant ID not found in context")
		return
	}

	// VERIFY SPECIFIC PERMISSION
	if err := h.RBACService.VerifyPermission(r.Context(), tenantID, userID, "leads.create"); err != nil {
		h.respondError(w, http.StatusForbidden, "Permission denied: lacks leads.create")
		return
	}

	// Rest of handler logic...
}
```

## Files to Modify (in order of priority)

### Tier 1: Critical Business Operations (P1)
- `internal/handlers/sales_handler.go` (10 CRUD ops)
  - CreateSalesLead → "leads.create"
  - UpdateSalesLead → "leads.update"
  - DeleteSalesLead → "leads.delete"
  - etc.

- `internal/handlers/gl_handler.go` (8 CRUD ops)
  - CreateAccount → "accounts.create"
  - PostJournalEntry → "entries.post"
  - ClosePeriod → "period.close"
  - etc.

- `internal/handlers/hr_handler.go` (12 CRUD ops)
  - CreateEmployee → "employees.create"
  - UpdateEmployee → "employees.update"
  - GeneratePayroll → "payroll.generate"
  - etc.

### Tier 2: Domain Operations (P2)
- `internal/handlers/real_estate_handler.go` (8 CRUD ops)
- `internal/handlers/construction_handler.go` (6 CRUD ops)
- `internal/handlers/civil_handler.go` (5 CRUD ops)
- `internal/handlers/purchase_handler.go` (7 CRUD ops)

### Tier 3: Support Operations (P3)
- `internal/handlers/billing_handler.go` (4 CRUD ops)
- `internal/handlers/company_handler.go` (4 CRUD ops)
- `internal/handlers/task_handler.go` (4 CRUD ops)
- Compliance handlers (RERA, HR, Tax)
- Dashboard handlers

## Update Router to Pass RBACService

### In `pkg/router/router.go`, update all handler creations:

```go
// Example - update all similar lines
salesHandler := handlers.NewSalesHandler(salesService.DB, rbacService)  // ADD rbacService

// For handler registrations, update signatures:
func RegisterSalesRoutes(r *mux.Router, salesService *services.SalesService, rbacService *services.RBACService) {
	handler := handlers.NewSalesHandler(salesService.DB, rbacService)
	// ... rest of routes
}
```

## Permission Code Mapping

All permission codes are defined in `internal/constants/permissions.go`:

### Sales Module
- `SalesLeadCreate` = "leads.create"
- `SalesLeadRead` = "leads.read"
- `SalesLeadUpdate` = "leads.update"
- `SalesLeadDelete` = "leads.delete"

### HR Module
- `EmployeeCreate` = "employees.create"
- `EmployeeRead` = "employees.read"
- `EmployeeUpdate` = "employees.update"
- `EmployeeDelete` = "employees.delete"
- `PayrollGenerate` = "payroll.generate"
- `PayrollRead` = "payroll.read"

### GL Module
- `AccountCreate` = "accounts.create"
- `AccountRead` = "accounts.read"
- `AccountUpdate` = "accounts.update"
- `EntryPost` = "entries.post"
- `PeriodClose` = "period.close"

See `internal/constants/permissions.go` for complete list.

## Implementation Checklist

- [ ] Update `NewSalesHandler` constructor
- [ ] Add RBACService to SalesHandler struct
- [ ] Update `pkg/router/router.go` to pass rbacService to NewSalesHandler
- [ ] Add permission checks to:
  - [ ] CreateSalesLead
  - [ ] UpdateSalesLead
  - [ ] DeleteSalesLead
  - [ ] (10 total sales operations)
- [ ] Repeat for HR handlers (12 ops)
- [ ] Repeat for GL handlers (8 ops)
- [ ] Repeat for Real Estate handlers (8 ops)
- [ ] Repeat for Construction handlers (6 ops)
- [ ] Repeat for Civil handlers (5 ops)
- [ ] Repeat for Purchase handlers (7 ops)
- [ ] Repeat for remaining handlers
- [ ] Verify all operations now have permission checks
- [ ] Test with restricted user account

## Quick Reference Code

### Getting User and Tenant from Context
```go
userID, _ := r.Context().Value(middleware.UserIDKey).(int64)
tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
```

### Verifying Permission
```go
if err := h.RBACService.VerifyPermission(r.Context(), tenantID, userID, permissionCode); err != nil {
	h.respondError(w, http.StatusForbidden, "Permission denied")
	return
}
```

### Permission Code Format
- Module: `module_name`
- Action: `create`, `read`, `update`, `delete`, `post`, `generate`, `close`
- Format: `"module_name.action"` (e.g., "leads.create", "accounts.post")

## Estimated Effort
- ~10 minutes per handler (update constructor + 3-5 CRUD operations)
- ~20 handlers = ~3-4 hours total
- Recommended to do in batches: Sales → HR → GL → Others

## Next Steps
After Phase 2 completion:
1. Run full test suite to verify no regressions
2. Test with restricted user accounts
3. Move to Phase 3: Audit Logging
4. Move to Phase 4: Data Filtering by Role
