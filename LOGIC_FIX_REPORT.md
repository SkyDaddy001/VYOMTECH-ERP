# System Logic Fixes - December 3, 2025

## Critical Issues Identified

### 1. CALLS LOGIC (FIXED ✅)

**Issues Found:**
- Call model using `int64` for IDs while Agent/Lead models use UUID strings
- Service using `int64` instead of string for agent_id and lead_id
- Handler converting UUID string IDs to int64 with `strconv.ParseInt()` - losing data integrity
- Type mismatch between User Context (string) and handler usage (int64)

**Fixes Applied:**
- ✅ Changed Call.ID, Call.AgentID, Call.LeadID from `int64` to `string` in model
- ✅ Updated CallService CreateCall to generate UUID if not provided
- ✅ Updated CallService GetCall signature to accept `id string` instead of `int64`
- ✅ Updated CallService EndCall signature to accept `id string` instead of `int64`
- ✅ Updated CallFilter to use `AgentID string` and `LeadID string` instead of `int64`
- ✅ Removed `strconv.ParseInt()` from CallHandler - now passes UUID directly
- ✅ Fixed userID context extraction to expect string instead of int64
- ✅ Added `github.com/google/uuid` import to call service

**Database Impact:**
- Ensure `call` table has VARCHAR/TEXT columns for: id, agent_id, lead_id (not INT)

---

### 2. PROJECT MANAGEMENT LOGIC (CRITICAL - NEEDS FIX)

**Issues Found:**
- **CRITICAL:** Using PostgreSQL syntax `$1, $2, $3...` instead of MySQL `?, ?, ?...`
- Affects ALL INSERT, UPDATE, DELETE operations in project_management_service.go
- Service will FAIL at runtime when executing queries
- All property customer profiles, units, BOQs will fail to create

**Fixes Required:**
- Replace ALL `$1, $2, $3...` with `?, ?, ?...` throughout the file
- Remove `RETURNING` clause (MySQL doesn't support it) - use `LAST_INSERT_ID()` if needed
- Verify all table names exist in migrations
- Fix any field references that don't match actual schema

**Methods Affected:**
1. CreateCustomerProfile() - Line 35
2. CreatePropertyUnit() - Line 161  
3. CreateBOQ() - Multiple locations
4. CreateUnitAreaStatement() - Line 161
5. All SELECT/UPDATE/DELETE operations

**Status:** ❌ NOT YET FIXED - Requires comprehensive string replacement

---

### 3. SALES LOGIC (NEEDS AUDIT)

**Potential Issues:**
- Sales leads using UUID strings (correct)
- Sales customers using UUID strings (correct)
- Need to verify:
  - All sales handlers properly extracting tenant from context
  - Sales quotations/orders logic validating customer existence before insertion
  - Discount calculations and tax computations are correct
  - Quote-to-order conversion logic preserves all required fields
  - Payment term calculations don't have off-by-one errors

**Status:** 🟡 PENDING DETAILED AUDIT

---

### 4. POST-SALES / SERVICE LOGIC (NEEDS AUDIT)

**Potential Issues:**
- Service request assignment validation
- Service resource allocation logic
- Warranty/support scope calculations
- SLA tracking and escalation triggers
- Contract term validations

**Status:** 🟡 PENDING DETAILED AUDIT

---

### 5. ACCOUNTS / GENERAL LEDGER (NEEDS AUDIT)

**Potential Issues:**
- Journal entry balance validation (debits must equal credits)
- Account type restrictions (cash can't have certain operations)
- Ledger posting sequence (must be chronological)
- Trial balance calculations
- Account reconciliation logic
- Tax calculations and compliance checks

**Status:** 🟡 PENDING DETAILED AUDIT

---

## Fix Priority

1. **IMMEDIATE (P0):** Fix project_management_service.go SQL syntax
   - Impact: Construction projects, property units won't work
   - Effort: High (large file with many queries)
   - Timeline: 1-2 hours

2. **HIGH (P1):** Complete Call fixes verification
   - Impact: All call operations
   - Effort: Already done, just need verification
   - Timeline: 15 minutes

3. **MEDIUM (P2):** Complete Sales logic audit
   - Impact: Sales operations
   - Effort: Medium (multiple files to check)
   - Timeline: 1 hour

4. **MEDIUM (P2):** Complete Post-Sales audit
   - Impact: Service and warranty operations
   - Effort: High (business logic verification)
   - Timeline: 2 hours

5. **HIGH (P1):** Complete Accounts logic audit
   - Impact: Financial reporting accuracy
   - Effort: High (complex business rules)
   - Timeline: 2+ hours

---

## Testing Plan

After fixes:

1. **Calls Testing**
   ```bash
   curl -X POST http://localhost:8080/api/v1/calls \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"lead_id":"<uuid>","agent_id":"<uuid>"}'
   ```

2. **Project Testing**
   - Create property unit
   - Create BOQ
   - Verify records in database

3. **Sales Testing**
   - Create quotation
   - Convert to order
   - Create invoice

4. **Accounts Testing**
   - Create journal entry
   - Verify GL posting
   - Run trial balance

---

## Files to Fix

### Immediate:
- `internal/services/project_management_service.go` - **PRIMARY FOCUS**

### Medium Term:
- `internal/services/sales_service.go`
- `internal/handlers/sales_handler.go`
- `internal/handlers/sales_quotations_orders.go`
- `internal/handlers/sales_invoices_payments.go`

### Post-Sales:
- Any service management handlers/services

### Accounts:
- `internal/services/gl_service.go`
- `internal/handlers/gl_handler.go`

---

## Next Steps

1. Run `go build` to identify compilation errors
2. Execute fixes in priority order
3. Run unit tests
4. Perform integration testing with Docker environment
5. Validate with comprehensive test script

---

**Last Updated:** December 3, 2025 15:58 UTC
**Status:** In Progress - Call fixes complete, Project fixes pending
