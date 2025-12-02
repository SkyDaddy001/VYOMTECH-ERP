# Dashboard API Quick Reference

## Financial Dashboard
**Base URL:** `/api/v1/dashboard/financial`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/profit-and-loss` | POST | P&L statement (requires: start_date, end_date) |
| `/balance-sheet` | POST | Balance sheet as of date (requires: as_of_date) |
| `/cash-flow` | POST | Cash flow statement (requires: start_date, end_date) |
| `/ratios` | GET | Financial ratios analysis |

---

## HR Dashboard
**Base URL:** `/api/v1/dashboard/hr`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/overview` | GET | HR department overview with headcount, attrition |
| `/payroll` | POST | Payroll summary (requires: payroll_month) |
| `/attendance` | POST | Attendance metrics (requires: start_date, end_date) |
| `/leaves` | GET | Leave requests and analytics |
| `/compliance` | GET | HR compliance status (ESI, EPF, PT, Gratuity) |

---

## Compliance Dashboard
**Base URL:** `/api/v1/dashboard/compliance`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/rera-status` | GET | RERA collections and fund utilization |
| `/hr-status` | GET | Labour law compliance status |
| `/tax-status` | GET | Income Tax and GST compliance |
| `/health-score` | GET | Overall compliance health score |
| `/documentation` | GET | Documentation upload tracking |

---

## Sales Dashboard
**Base URL:** `/api/v1/dashboard/sales`

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/overview` | GET | Sales overview with YTD revenue and pipeline |
| `/pipeline` | GET | Pipeline analysis by stage and region |
| `/metrics` | POST | Sales metrics (requires: period - monthly/quarterly/annual) |
| `/forecast` | GET | Sales forecast by rep and product |
| `/invoices` | GET | Invoice tracking with aging analysis |
| `/competition` | GET | Competitive intelligence and win/loss analysis |

---

## Response Format Example

### Financial P&L Response
```json
{
  "period": {
    "start": "2024-01-01",
    "end": "2024-12-31"
  },
  "income": {
    "sales_revenue": 1000000,
    "service_revenue": 200000,
    "other_income": 50000
  },
  "expenses": {
    "cogs": 400000,
    "operating_expenses": 150000,
    "administrative": 100000,
    "depreciation": 50000,
    "finance_costs": 25000
  },
  "profit_summary": {
    "gross_profit": 650000,
    "operating_profit": 400000,
    "net_profit": 375000,
    "net_margin": 25.5
  }
}
```

### HR Overview Response
```json
{
  "workforce": {
    "total_employees": 150,
    "active_employees": 145,
    "on_leave": 3,
    "inactive": 2,
    "contractors": 0
  },
  "departments": { },
  "positions": { },
  "headcount_trend": [],
  "attrition": {
    "current_month": 0.67,
    "ytd_attrition": 5.33,
    "voluntary_exit": 4,
    "involuntary_exit": 2
  }
}
```

### Compliance Status Response
```json
{
  "overall_status": "Compliant",
  "compliance_modules": {
    "esi": {
      "status": "Compliant",
      "employees": 150,
      "violations": 0
    },
    "epf": {
      "status": "Compliant",
      "employees": 145,
      "violations": 0
    },
    "professional_tax": {
      "status": "Compliant",
      "employees": 95,
      "violations": 0
    },
    "gratuity": {
      "status": "Compliant",
      "eligible_employees": 42,
      "accrued_liability": 2500000
    }
  },
  "violation_summary": {
    "critical": 0,
    "high": 0,
    "medium": 0,
    "low": 0
  },
  "upcoming_deadlines": [
    {
      "compliance_item": "ESI Return Filing",
      "due_date": "2024-12-15",
      "status": "Pending"
    }
  ]
}
```

### Sales Pipeline Response
```json
{
  "pipeline_summary": {
    "total_opportunities": 45,
    "total_value": 5000000,
    "weighted_value": 3250000,
    "average_deal_size": 111111
  },
  "by_stage": [
    {
      "stage": "Prospecting",
      "opportunities": 15,
      "value": 1500000,
      "count_change": 2,
      "value_change": 100000
    },
    {
      "stage": "Closed Won",
      "opportunities": 8,
      "value": 1200000,
      "count_change": 1,
      "value_change": 200000
    }
  ],
  "aged_pipeline": {
    "0_to_30_days": 1500000,
    "31_to_60_days": 1200000,
    "61_to_90_days": 800000,
    "over_90_days": 500000
  }
}
```

---

## Error Handling

All endpoints return appropriate HTTP status codes:
- **200 OK** - Request successful
- **400 Bad Request** - Invalid request body/parameters
- **401 Unauthorized** - Missing or invalid authentication
- **403 Forbidden** - Insufficient permissions
- **500 Internal Server Error** - Server error

Example error response:
```json
{
  "error": "Invalid request body",
  "status": 400
}
```

---

## Authentication

All dashboard endpoints require:
- **Header:** `Authorization: Bearer <JWT_TOKEN>`
- **Header:** `X-Tenant-ID: <TENANT_ID>` (for multi-tenant isolation)

---

## Data Refresh Patterns

- **Financial Dashboard:** Daily (end of business)
- **HR Dashboard:** Real-time (updated with attendance/payroll transactions)
- **Compliance Dashboard:** Daily (overnight batch reconciliation)
- **Sales Dashboard:** Real-time (updated with invoice/opportunity changes)

---

## Future Enhancements

- [ ] Custom date range filtering on all endpoints
- [ ] Department/segment filtering
- [ ] Export to PDF/Excel formats
- [ ] Role-based dashboard customization
- [ ] WebSocket real-time updates
- [ ] Mobile-optimized response formats
- [ ] Caching layer for improved performance
- [ ] Advanced forecasting with ML models
