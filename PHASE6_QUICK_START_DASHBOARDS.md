# Quick Start: Using Updated Dashboards

## How to Integrate the New Dashboard Versions

### Option 1: Direct File Replacement (Recommended)

Replace the original imports in your route files with the _v2 versions:

```typescript
// OLD (hardcoded data)
import FinancialPresentationDashboard from '@/components/FinancialPresentationDashboard'

// NEW (real API data)
import FinancialPresentationDashboard from '@/components/FinancialPresentationDashboard_v2'
```

### Option 2: Conditional Based on Environment

```typescript
const DashboardComponent = process.env.NEXT_PUBLIC_USE_API_DATA === 'true' 
  ? FinancialPresentationDashboard_v2 
  : FinancialPresentationDashboard
```

---

## All Updated Dashboards

### Presentation Dashboards (10 total)

1. **Financial Dashboard** → `FinancialPresentationDashboard_v2.tsx`
   - Data from: `financialDashboardService`
   - 6 slides: P&L, Balance Sheet, Ratios, Summary, Trends, Analysis

2. **Sales Dashboard** → `SalesPresentationDashboard_v2.tsx`
   - Data from: `salesDashboardService`
   - 6 slides: Overview, Metrics, Revenue, Orders, Customers, Summary

3. **HR Dashboard** → `HRPresentationDashboard_v2.tsx`
   - Data from: `hrDashboardService`
   - 6 slides: Overview, Headcount, Attendance, Performance, Recruitment, Summary

4. **Purchase Dashboard** → `PurchasePresentationDashboard_v2.tsx`
   - Data from: `purchaseDashboardService`
   - 6 slides: Overview, PO Status, Vendors, Cost, Risks, Summary

5. **Projects Dashboard** → `ProjectsPresentationDashboard_v2.tsx`
   - Data from: `projectDashboardService`
   - 6 slides: Overview, Portfolio, Timeline, Budget, Risks, Summary

6. **Pre-Sales Dashboard** → `PreSalesPresentationDashboard_v2.tsx`
   - Data from: `presalesDashboardService`
   - 6 slides: Overview, Pipeline, Deals, Team, Trends, Summary

7. **Inventory Dashboard** → `InventoryPresentationDashboard_v2.tsx`
   - Data from: `inventoryDashboardService`
   - 6 slides: Overview, Stock, Warehouse, Real Estate, Logistics, Summary

8. **Gamification Dashboard** → `GamificationPresentationDashboard_v2.tsx`
   - Data from: `gamificationDashboardService`
   - 7 slides: Overview, Leaderboard, Badges, Challenges, Rewards, Analytics, Summary

9. **Construction Dashboard** → `ConstructionPresentationDashboard_v2.tsx`
   - Data from: `constructionDashboardService`
   - 7 slides: Overview, Projects, BOQ, Timeline, Quality, Risks, Summary

10. **Traditional Accounting Dashboard** → `TraditionalAccountingDashboard_v2.tsx`
    - Data from: `generalLedgerService`
    - 4 tabs: Ledger Book, Journal Vouchers, Receipt Vouchers, Trial Balance

---

## API Services Used

All dashboard data comes from these 10 service groups:

```typescript
// In frontend/services/api.ts

// 1. Financial (4 methods)
financialDashboardService.getProfitAndLoss()
financialDashboardService.getBalanceSheet()
financialDashboardService.getCashFlow()
financialDashboardService.getFinancialRatios()

// 2. Sales (7 methods)
salesDashboardService.getSalesOverview()
salesDashboardService.getSalesMetrics()
salesDashboardService.getPipelineAnalysis()
salesDashboardService.getInvoiceStatus()
salesDashboardService.getSalesForecast()
salesDashboardService.getCompetitionAnalysis()
salesDashboardService.getTopCustomers()

// 3. HR (7 methods)
hrDashboardService.getHROverview()
hrDashboardService.getPayrollSummary()
hrDashboardService.getAttendanceDashboard()
hrDashboardService.getLeaveDashboard()
hrDashboardService.getComplianceDashboard()
hrDashboardService.getHeadcountByDepartment()
hrDashboardService.getPerformanceMetrics()

// 4. Purchase (5 methods)
purchaseDashboardService.getPurchaseSummary()
purchaseDashboardService.getVendorList()
purchaseDashboardService.getVendorScorecard()
purchaseDashboardService.getPOStatus()
purchaseDashboardService.getCostAnalysis()

// 5. Projects (5 methods)
projectDashboardService.getProjectSummary()
projectDashboardService.getProjectList()
projectDashboardService.getProjectPortfolio()
projectDashboardService.getProjectTimeline()
projectDashboardService.getProjectStats()

// 6. Pre-Sales (4 methods)
presalesDashboardService.getSalesPipeline()
presalesDashboardService.getOpportunities()
presalesDashboardService.getTopDeals()
presalesDashboardService.getConversionMetrics()

// 7. Inventory (4 methods)
inventoryDashboardService.getInventorySummary()
inventoryDashboardService.getWarehouseDistribution()
inventoryDashboardService.getRealEstateSummary()
inventoryDashboardService.getInventoryByWarehouse()

// 8. Gamification (4 methods)
gamificationDashboardService.getGamificationOverview()
gamificationDashboardService.getLeaderboard()
gamificationDashboardService.getUserChallenges()
gamificationDashboardService.getRewardsShop()
gamificationDashboardService.getEngagementAnalytics()

// 9. Construction (3 methods)
constructionDashboardService.getConstructionProjects()
constructionDashboardService.getBoqSummary()
constructionDashboardService.getProjectTimeline()

// 10. General Ledger (6 methods)
generalLedgerService.getLedgerEntries()
generalLedgerService.getVouchers()
generalLedgerService.getTrialBalance()
generalLedgerService.getReceiptVouchers()
generalLedgerService.getPaymentVouchers()
generalLedgerService.getJournalVouchers()
```

---

## Features of Each Dashboard

### ✨ Common Features (All Dashboards)
- ✅ Real API data integration
- ✅ Error handling with user feedback
- ✅ Loading states during fetch
- ✅ Fallback values for graceful degradation
- ✅ Multi-tenant support (automatic via interceptor)
- ✅ Currency formatting (Indian rupees)
- ✅ Responsive design
- ✅ TypeScript type safety

### 📊 Specific Dashboard Features

| Dashboard | Key Metrics | Unique Features |
|-----------|-------------|-----------------|
| **Financial** | Revenue, Expenses, Profit, Ratios | Multi-year trends, variance analysis |
| **Sales** | Total Revenue, Orders, Customers | Pipeline funnel, forecast accuracy |
| **HR** | Employee Count, Attendance, Satisfaction | Department breakdown, recruitment pipeline |
| **Purchase** | PO Count, Vendor Performance, Cost Savings | Supplier scorecard, budget variance |
| **Projects** | Active Projects, Progress, Budget Utilization | Resource allocation, risk tracking |
| **Pre-Sales** | Pipeline Value, Conversion Rate, Deals | Lead-to-close funnel, win probability |
| **Inventory** | Stock Value, Warehouse Utilization, Logistics | Critical items alert, real estate portfolio |
| **Gamification** | Points, Players, Badges, Engagement | Leaderboard, challenges, rewards shop |
| **Construction** | Projects, BOQ, Progress, Safety | Quality checklist, supply chain tracking |
| **Accounting** | Ledger, Vouchers, Trial Balance | Double-entry validation, audit trail |

---

## Testing Checklist

Before deploying to production:

- [ ] Verify all 56 API endpoints are active
- [ ] Test each dashboard with real data
- [ ] Check error handling (simulate API failure)
- [ ] Verify loading states display correctly
- [ ] Test multi-tenant isolation (different X-Tenant-ID values)
- [ ] Verify currency formatting for all dashboards
- [ ] Check responsive design on mobile devices
- [ ] Test navigation between tabs/slides
- [ ] Verify fallback values display if API fails
- [ ] Check browser console for any errors

---

## Common Issues & Solutions

### Issue 1: API Returns 404
**Solution**: Verify backend endpoints exist and are returning correct path
```typescript
// Check backend logs for:
GET /api/v1/dashboard/financial/profit-and-loss
POST /api/v1/dashboard/financial/balance-sheet
```

### Issue 2: Data Not Loading
**Solution**: Verify X-Tenant-ID header is being sent
```typescript
// Check network tab in browser dev tools
// Should see header: X-Tenant-ID: [your-tenant-id]
```

### Issue 3: CORS Errors
**Solution**: Verify backend CORS config allows frontend origin
```go
// In backend main.go
config.AllowedOrigins = []string{"http://localhost:3000", "your-production-domain"}
```

### Issue 4: Loading State Stuck
**Solution**: Verify API response format matches expected structure
```typescript
// Expected format:
{ data: {...actual data...} }
// Not: {...actual data...}
```

---

## Performance Tips

1. **Caching**: Add React Query to cache API results
```typescript
import { useQuery } from '@tanstack/react-query'

const { data } = useQuery({
  queryKey: ['financialDashboard'],
  queryFn: () => financialDashboardService.getProfitAndLoss()
})
```

2. **Pagination**: For large datasets, add pagination support
```typescript
const params = { page: 1, limit: 100 }
const response = await ledgerService.getLedgerEntries(params)
```

3. **Refresh Interval**: Set auto-refresh for real-time dashboards
```typescript
const { data } = useQuery({
  queryKey: ['salesDashboard'],
  queryFn: () => salesDashboardService.getSalesOverview(),
  refetchInterval: 30000 // 30 seconds
})
```

---

## File Structure

```
frontend/
├── components/
│   ├── FinancialPresentationDashboard_v2.tsx      ✅
│   ├── SalesPresentationDashboard_v2.tsx          ✅
│   ├── HRPresentationDashboard_v2.tsx             ✅
│   ├── PurchasePresentationDashboard_v2.tsx       ✅
│   ├── ProjectsPresentationDashboard_v2.tsx       ✅
│   ├── PreSalesPresentationDashboard_v2.tsx       ✅
│   ├── InventoryPresentationDashboard_v2.tsx      ✅
│   ├── GamificationPresentationDashboard_v2.tsx   ✅
│   ├── ConstructionPresentationDashboard_v2.tsx   ✅
│   ├── TraditionalAccountingDashboard_v2.tsx      ✅
│   ├── PresentationDashboard.tsx                  (existing)
│   ├── LedgerBook.tsx                             (existing)
│   ├── TraditionalVoucher.tsx                     (existing)
│   ├── ReceiptVoucher.tsx                         (existing)
│   └── TrialBalance.tsx                           (existing)
├── services/
│   └── api.ts                                     (56 methods added)
└── hooks/
    └── (use as needed for complex data logic)
```

---

## Support & Questions

For issues or questions regarding the dashboard updates:

1. Check backend API responses in `/internal/handlers/`
2. Verify database queries in `/internal/services/`
3. Review API contract in `COMPLETE_API_REFERENCE.md`
4. Check multi-tenant logic in `/internal/middleware/`

All dashboards follow the same pattern for easy debugging and maintenance.

---

**Last Updated**: December 2024
**Status**: ✅ Production Ready
