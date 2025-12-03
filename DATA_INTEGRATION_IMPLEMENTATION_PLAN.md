# Data Integration Implementation Plan - Phase 6
**Status**: Ready for Implementation  
**Last Updated**: Current Session  
**Objective**: Replace hardcoded data with real database values for all 10 presentation dashboards and 4 traditional accounting components

---

## 📋 Executive Summary

### Current State
- ✅ 10 Presentation Dashboard components created with hardcoded sample data
- ✅ 4 Traditional accounting components with hardcoded sample data
- ✅ Backend: 30+ services, 54 handlers, 32 data models
- ✅ Frontend API client configured with auth headers and multi-tenancy
- ✅ Database: MySQL with GORM models

### Tasks
1. **Replace hardcoded data** with real API calls (10 dashboards × ~500-750 lines each)
2. **Identify missing endpoints** and create them in backend
3. **Implement missing features**: pagination, filtering, date ranges, export
4. **Connect existing backend handlers** to frontend dashboards

### Scope: ~7,000 lines across frontend + backend

---

## 🗂️ Dashboard-to-Service Mapping

### Dashboard 1: FinancialPresentationDashboard.tsx (299 lines)
**Status**: ✅ HIGH PRIORITY - Backend endpoints exist  
**Service**: GLService  
**Existing Endpoints**:
- `POST /api/v1/dashboard/profit-and-loss` → GetProfitAndLoss()
- `POST /api/v1/dashboard/balance-sheet` → GetBalanceSheet()
- `POST /api/v1/dashboard/cash-flow` → GetCashFlow()
- `GET /api/v1/dashboard/ratios` → GetFinancialRatios()

**Current Hardcoded Data**:
```
- Total Assets: ₹1.2Cr (hardcoded)
- Total Liabilities: ₹42L (hardcoded)
- Revenue: ₹85L (hardcoded)
- Net Profit: ₹33L (hardcoded)
- 6-month chart data (Jan-Jun): All hardcoded
```

**Integration Plan**:
1. Create API service methods in `frontend/services/api.ts`:
   - `getFinancialOverview(startDate, endDate)`
   - `getBalanceSheetData(asOfDate)`
   - `getProfitAndLossData(startDate, endDate)`
   - `getCashFlowData(startDate, endDate)`
   - `getFinancialRatios()`

2. Replace hardcoded data with API calls in component lifecycle
3. Add state management (useState/useQuery) for data fetching
4. Add loading states and error handling

---

### Dashboard 2: SalesPresentationDashboard.tsx (450 lines)
**Status**: 🔄 MEDIUM PRIORITY - Backend partially exists  
**Service**: SalesService + SalesDashboardHandler  
**Existing Endpoints**:
- `GET /api/v1/dashboard/sales/overview` → GetSalesOverview()
- `GET /api/v1/dashboard/sales/pipeline` → GetPipelineAnalysis()
- `GET /api/v1/dashboard/sales/metrics` → GetSalesMetrics()
- `GET /api/v1/dashboard/sales/invoice-status` → GetInvoiceStatus()
- `GET /api/v1/dashboard/sales/forecast` → GetForecast()
- `GET /api/v1/dashboard/sales/competition` → GetCompetitionAnalysis()

**Current Hardcoded Data**:
```
- 142 invoices (hardcoded)
- 89 orders (hardcoded)
- ₹24.5L revenue (hardcoded)
- 4 top customers (hardcoded)
- 6-month trend data (hardcoded)
```

**Integration Plan**:
1. Create API service methods:
   - `getSalesOverview(dateRange)`
   - `getPipelineAnalysis()`
   - `getSalesMetrics()`
   - `getInvoiceStatus()`
   - `getTopCustomers(limit)`

2. Replace all hardcoded arrays with API calls
3. Add pagination for customer/invoice lists

---

### Dashboard 3: HRPresentationDashboard.tsx (550 lines)
**Status**: 🔄 MEDIUM PRIORITY - Backend exists (HRService, HRDashboardHandler)  
**Service**: HRService  
**Existing Endpoints**:
- `GET /api/v1/dashboard/hr/overview` → GetHROverview() (in HRDashboardHandler)
- `GET /api/v1/dashboard/hr/headcount` → GetHeadcount()
- `GET /api/v1/dashboard/hr/performance` → GetPerformanceMetrics()
- `GET /api/v1/dashboard/hr/attendance` → GetAttendanceData()

**Current Hardcoded Data**:
```
- 245 total employees (hardcoded)
- 94% attendance rate (hardcoded)
- 6 departments with counts: Engineering 89, Sales 67, Operations 52, Finance 23, HR 14, Marketing 18
- Performance data (salary info, satisfaction)
- Leave breakdown
```

**Integration Plan**:
1. Create API service methods:
   - `getHRMetrics()` → overall KPIs
   - `getHeadcountByDepartment()`
   - `getAttendanceData(month)`
   - `getLeaveData()`
   - `getPerformanceMetrics()`

2. Replace all hardcoded employee counts and metrics
3. Query actual employee data from HRService

---

### Dashboard 4: PurchasePresentationDashboard.tsx (600 lines)
**Status**: 🔄 MEDIUM-LOW PRIORITY - Backend may need endpoints  
**Service**: (Needs analysis - likely PurchaseService or needs creation)  
**Potential Endpoints** (to be verified/created):
- `/api/v1/purchase/summary` → PO counts, values
- `/api/v1/purchase/vendors` → Vendor list with scores
- `/api/v1/purchase/po-status` → PO status breakdown
- `/api/v1/purchase/cost-analysis` → Cost savings

**Current Hardcoded Data**:
```
- 47 vendors (hardcoded)
- ₹12.5 Cr purchase value (hardcoded)
- ₹1.15 Cr savings (hardcoded)
- PO status: Draft 23, Confirmed 78, In Transit 34, Received 52, Delayed 8, Completed 89
- Vendor scorecard data
```

**Integration Plan**:
1. Check if PurchaseService exists in backend
2. If not, create PurchaseHandler with:
   - `GetPurchaseSummary()` - KPIs
   - `GetVendorList()` - Vendor list with scores
   - `GetPOStatus()` - PO breakdown by status
   - `GetCostAnalysis()` - Savings metrics

3. Create API service methods in frontend
4. Replace all hardcoded vendor and PO data

---

### Dashboard 5: ProjectsPresentationDashboard.tsx (650 lines)
**Status**: 🔄 MEDIUM-LOW PRIORITY - ProjectManagementService exists  
**Service**: ProjectManagementService  
**Potential Endpoints** (likely need creation):
- `/api/v1/projects/summary` → Overall project KPIs
- `/api/v1/projects/list` → Project list with progress
- `/api/v1/projects/portfolio` → Portfolio view
- `/api/v1/projects/timeline` → Timeline data

**Current Hardcoded Data**:
```
- 18 projects (hardcoded)
- ₹85 Cr total value (hardcoded)
- 62% average completion (hardcoded)
- 6 sample projects with details
- 245 team members (hardcoded)
- Budget and resource data
```

**Integration Plan**:
1. Create ProjectDashboardHandler with endpoints
2. Create API service methods:
   - `getProjectPortfolio()`
   - `getProjectList()`
   - `getProjectTimeline()`
   - `getProjectStats()`

3. Replace all hardcoded project data

---

### Dashboard 6: PreSalesPresentationDashboard.tsx (700 lines)
**Status**: 🔄 MEDIUM PRIORITY - Likely maps to SalesService (LeadService exists)  
**Service**: SalesService + LeadService  
**Potential Endpoints**:
- `/api/v1/sales/pipeline` → Pipeline summary
- `/api/v1/sales/opportunities` → Opportunity list
- `/api/v1/sales/deals` → Deal progression

**Current Hardcoded Data**:
```
- ₹42 Cr pipeline value (hardcoded)
- 34% conversion rate (hardcoded)
- 127 opportunities (hardcoded)
- ₹18 Cr expected value (hardcoded)
- 5 top deals (hardcoded)
```

**Integration Plan**:
1. Use existing SalesService endpoints
2. Create API service methods:
   - `getSalesPipeline()`
   - `getOpportunities()`
   - `getTopDeals(limit)`
   - `getConversionMetrics()`

3. Replace hardcoded pipeline and opportunity data

---

### Dashboard 7: InventoryPresentationDashboard.tsx (700 lines)
**Status**: 🔄 LOW PRIORITY - May need InventoryService  
**Service**: (Needs analysis - likely Phase 3C related)  
**Potential Endpoints** (to be created):
- `/api/v1/inventory/summary` → Inventory KPIs
- `/api/v1/inventory/warehouses` → Warehouse distribution
- `/api/v1/inventory/real-estate` → Property portfolio

**Current Hardcoded Data**:
```
- ₹8.5 Cr inventory value (hardcoded)
- 12,450 units (hardcoded)
- 24 warehouses (hardcoded)
- 92% utilization (hardcoded)
- Real estate: ₹90 Cr, 4 properties (2 owned, 2 leased)
- Warehouse distribution data
```

**Integration Plan**:
1. Create or verify InventoryHandler
2. Create API service methods for warehouse and inventory data
3. Create real estate endpoints (possibly via RealEstateService)
4. Replace all hardcoded inventory and real estate data

---

### Dashboard 8: GamificationPresentationDashboard.tsx (750 lines)
**Status**: 🔄 MEDIUM PRIORITY - GamificationService exists  
**Service**: GamificationService  
**Existing Endpoints** (need verification):
- Gamification handlers should exist from Phase 1
- Need dashboard aggregations

**Current Hardcoded Data**:
```
- 3.2M points (hardcoded)
- 245 users engaged (hardcoded)
- 1,250+ badges earned (hardcoded)
- 87% engagement rate (hardcoded)
- Leaderboard (hardcoded)
- Challenge data (hardcoded)
- Rewards shop inventory (hardcoded)
```

**Integration Plan**:
1. Create GamificationDashboardHandler if not exists
2. Create API service methods:
   - `getGamificationOverview()`
   - `getLeaderboard(limit, page)`
   - `getUserChallenges()`
   - `getRewardsShop()`
   - `getEngagementAnalytics()`

3. Replace all hardcoded gamification data
4. Implement real-time leaderboard updates

---

### Dashboard 9: ConstructionPresentationDashboard.tsx (550+ lines)
**Status**: 🔄 MEDIUM PRIORITY - ConstructionService + BOQService exist  
**Service**: ConstructionService + BOQService  
**Potential Endpoints**:
- `/api/v1/construction/projects` → Project list
- `/api/v1/construction/boq` → BOQ data
- `/api/v1/construction/timeline` → Timeline tracking
- `/api/v1/construction/safety` → Safety metrics

**Current Hardcoded Data**:
```
- 12 projects (hardcoded)
- ₹45.8 Cr BOQ total (hardcoded)
- 68% average completion (hardcoded)
- 245 workers (hardcoded)
- 5 sample projects with progress
- Safety metrics and risks
```

**Integration Plan**:
1. Create ConstructionDashboardHandler
2. Create API service methods:
   - `getConstructionProjects()`
   - `getBoqSummary()`
   - `getProjectTimeline(projectId)`
   - `getSafetyMetrics()`
   - `getWorkerAllocation()`

3. Replace all hardcoded construction data

---

### Dashboard 10: TraditionalAccountingDashboard.tsx (300 lines)
**Status**: ✅ HIGH PRIORITY - GLService endpoints exist  
**Service**: GLService  
**Uses Sub-components**:
- LedgerBook.tsx - General ledger entries
- TraditionalVoucher.tsx - Journal vouchers
- ReceiptVoucher.tsx - Receipt/payment vouchers
- TrialBalance.tsx - Trial balance view

**Current Hardcoded Data**:
```
- All ledger entries (hardcoded)
- All voucher data (hardcoded)
- Receipt/payment details (hardcoded)
- Trial balance entries (hardcoded)
```

**Integration Plan**:
1. Create API service methods:
   - `getLedgerEntries(accountCode, dateRange)`
   - `getVouchers(type, dateRange)`
   - `getTrialBalance(asOfDate)`
   - `getReceiptVouchers()`
   - `getPaymentVouchers()`

2. Update each sub-component to use real GL data
3. Replace all hardcoded transaction data

---

## 🔧 Implementation Sequence & Priority

### Phase 6A: High Priority (Days 1-3)
**Estimated**: 20-30 hours  
**Impact**: Core financial reporting

1. **FinancialPresentationDashboard** (3-4 hours)
   - Endpoints already exist
   - Simple API integration
   - No backend changes needed

2. **TraditionalAccountingDashboard + sub-components** (4-5 hours)
   - Endpoints already exist
   - Multiple sub-components to update
   - More complex state management

3. **SalesPresentationDashboard** (3-4 hours)
   - Most endpoints exist
   - Some may need minor adjustments
   - High business value

### Phase 6B: Medium Priority (Days 3-5)
**Estimated**: 20-25 hours  
**Impact**: Operational dashboards

4. **HRPresentationDashboard** (3-4 hours)
   - Endpoints likely exist
   - High employee data volume

5. **PreSalesPresentationDashboard** (2-3 hours)
   - Uses existing SalesService endpoints
   - Reuses some patterns from Sales dashboard

6. **GamificationPresentationDashboard** (4-5 hours)
   - Real-time updates needed
   - Complex leaderboard logic

### Phase 6C: Lower Priority (Days 5-7)
**Estimated**: 15-20 hours  
**Impact**: Specialized functions

7. **ConstructionPresentationDashboard** (4-5 hours)
   - Need to verify BOQ endpoints
   - Integration with multiple services

8. **ProjectsPresentationDashboard** (3-4 hours)
   - ProjectManagementService endpoints may need creation

9. **PurchasePresentationDashboard** (2-3 hours)
   - PurchaseService may need verification

10. **InventoryPresentationDashboard** (3-4 hours)
    - Real estate + inventory integration
    - New endpoints likely needed

---

## 📝 Implementation Template

### Step 1: Analyze Backend Handler
```go
// Check: /internal/handlers/{service}_handler.go
// Identify: Available endpoints and methods
// Note: Required parameters and response format
```

### Step 2: Create Frontend API Service Methods
```typescript
// In: frontend/services/api.ts
// Add: Service-specific methods that wrap axios calls
// Pattern:
async getSalesOverview(dateRange?: { start: Date; end: Date }) {
  const response = await this.get('/api/v1/dashboard/sales/overview', {
    params: {
      start_date: dateRange?.start.toISOString(),
      end_date: dateRange?.end.toISOString(),
    }
  })
  return response.data
}
```

### Step 3: Create React Hooks (Optional)
```typescript
// In: frontend/hooks/use{Entity}.ts
// Pattern: useSalesData, useHRMetrics, etc.
// Benefits: Reusable, testable, separates data from UI
```

### Step 4: Update Dashboard Component
```typescript
// Replace hardcoded data with API calls
// Use useState for data, loading, error states
// Call API on component mount
// Add error boundaries and loading screens
```

### Step 5: Test & Validate
```
- Verify API responses match expected format
- Test with different date ranges
- Test multi-tenant isolation
- Test error handling
```

---

## 🛠️ Missing Endpoints to Create

### Backend Handlers to Create/Extend

**1. PurchaseHandler** (if doesn't exist)
```go
type PurchaseHandler struct {
  Service *services.PurchaseService
}

func (h *PurchaseHandler) GetPurchaseSummary(w http.ResponseWriter, r *http.Request)
func (h *PurchaseHandler) GetVendorList(w http.ResponseWriter, r *http.Request)
func (h *PurchaseHandler) GetPOStatus(w http.ResponseWriter, r *http.Request)
func (h *PurchaseHandler) GetCostAnalysis(w http.ResponseWriter, r *http.Request)
```

**2. ProjectDashboardHandler** (if doesn't exist)
```go
type ProjectDashboardHandler struct {
  Service *services.ProjectManagementService
}

func (h *ProjectDashboardHandler) GetProjectSummary(w http.ResponseWriter, r *http.Request)
func (h *ProjectDashboardHandler) GetProjectList(w http.ResponseWriter, r *http.Request)
func (h *ProjectDashboardHandler) GetProjectTimeline(w http.ResponseWriter, r *http.Request)
```

**3. GamificationDashboardHandler** (if doesn't exist)
```go
type GamificationDashboardHandler struct {
  Service *services.GamificationService
}

func (h *GamificationDashboardHandler) GetOverview(w http.ResponseWriter, r *http.Request)
func (h *GamificationDashboardHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request)
func (h *GamificationDashboardHandler) GetUserChallenges(w http.ResponseWriter, r *http.Request)
```

**4. InventoryDashboardHandler** (if doesn't exist)
```go
type InventoryDashboardHandler struct {
  Service *services.InventoryService // May need to create
}

func (h *InventoryDashboardHandler) GetInventorySummary(w http.ResponseWriter, r *http.Request)
func (h *InventoryDashboardHandler) GetWarehouseDistribution(w http.ResponseWriter, r *http.Request)
```

**5. Extend HRDashboardHandler** (if incomplete)
- Verify all required endpoints exist
- Add missing metrics endpoints

**6. Extend SalesDashboardHandler** (if incomplete)
- Verify all required endpoints exist
- Add top customers, invoice status details

---

## 🔌 API Service Methods to Create

### frontend/services/api.ts - New Methods

```typescript
// Financial
async getFinancialOverview(startDate?: Date, endDate?: Date)
async getBalanceSheetData(asOfDate?: Date)
async getProfitAndLossData(startDate?: Date, endDate?: Date)
async getCashFlowData(startDate?: Date, endDate?: Date)
async getFinancialRatios()

// Sales
async getSalesOverview(dateRange?: DateRange)
async getPipelineAnalysis()
async getSalesMetrics()
async getTopCustomers(limit: number)
async getInvoiceStatus()
async getSalesForcast()

// HR
async getHRMetrics()
async getHeadcountByDepartment()
async getAttendanceData(month?: Date)
async getLeaveData()
async getPerformanceMetrics()

// Purchase
async getPurchaseSummary()
async getVendorList()
async getPOStatus()
async getCostAnalysis()

// Projects
async getProjectSummary()
async getProjectList(page?: number, limit?: number)
async getProjectTimeline(projectId: string)
async getProjectStats()

// Inventory
async getInventorySummary()
async getWarehouseDistribution()
async getRealEstateSummary()

// Gamification
async getGamificationOverview()
async getLeaderboard(limit: number, page: number)
async getUserChallenges(userId: string)
async getRewardsShop()

// Accounting (GL)
async getLedgerEntries(accountCode: string, dateRange: DateRange)
async getVouchers(type: string, dateRange: DateRange)
async getTrialBalance(asOfDate: Date)
async getReceiptVouchers()
async getPaymentVouchers()
```

---

## ✅ Validation Checklist

### Before Starting Phase 6
- [ ] Backend services confirmed (30+ services running)
- [ ] All 261 tests passing
- [ ] Database connected and migrations applied
- [ ] Frontend API client working with auth headers
- [ ] Multi-tenant isolation verified
- [ ] Docker services running

### Per Dashboard
- [ ] Backend endpoints identified/created
- [ ] API service methods implemented
- [ ] Hardcoded data replaced with API calls
- [ ] Loading states added
- [ ] Error handling implemented
- [ ] Multi-tenant filtering applied
- [ ] Date range filtering working
- [ ] Pagination implemented (if applicable)

### Final Validation
- [ ] All 10 dashboards loading real data
- [ ] All 4 traditional accounting components working
- [ ] No hardcoded sample data remains
- [ ] All API calls multi-tenant aware
- [ ] Performance acceptable (< 3s initial load)
- [ ] Error cases handled gracefully

---

## 📊 Expected Outcomes

### Code Changes
- **Frontend**: ~5,000-6,000 lines modified across 10 dashboards
- **Backend**: ~1,500-2,000 lines of new handlers/endpoints
- **API Service**: ~100-150 new methods in api.ts

### Features Implemented
- Real-time data binding for all dashboards
- Multi-tenant data isolation
- Date range filtering
- Pagination and sorting
- Error handling and loading states
- Data refresh capabilities

### Business Value
- Complete visibility into all business operations
- Real-time KPI dashboards
- Traditional accounting integration
- PowerPoint-style professional presentation
- Full ERP functionality end-to-end

---

## 🚀 Next Steps (Immediate)

1. **Verify Backend Handlers**
   - Check which handlers are missing
   - List all available endpoints
   - Identify gaps

2. **Start with FinancialPresentationDashboard**
   - Already has endpoints
   - Quickest win
   - Sets template for others

3. **Create API Service Methods**
   - Batch update frontend/services/api.ts
   - Add TypeScript interfaces for responses

4. **Update Components**
   - Replace hardcoded data incrementally
   - Add loading/error states
   - Test multi-tenant isolation

5. **Create Missing Handlers**
   - PurchaseHandler if missing
   - ProjectDashboardHandler if missing
   - InventoryHandler if missing

---

**Ready to proceed? Start with backend verification or jump to Frontend phase 6A implementation.**
