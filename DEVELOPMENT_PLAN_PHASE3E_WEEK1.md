# Phase 3E Development Plan - Week 1 (Nov 25 - Dec 1, 2025)

## Overview
Complete Sales module implementation (0% → 100%)

---

## Module: Sales (CRM + Order Management)

### What is the Sales Module?
The Sales module manages the entire customer sales lifecycle:
1. **Lead Management** - Track and qualify leads
2. **Customer Quotations** - Create and track quotes
3. **Sales Orders** - Convert quotes to orders
4. **Invoicing** - Generate sales invoices (posts to GL)
5. **Revenue Recognition** - Track revenue realization

### Business Flow
```
Lead → Quotation → Sales Order → Invoice → Payment → GL Entry
```

---

## Database Schema (18 tables, 230+ fields)

### Core Tables

1. **sales_leads**
   - Track prospective customers
   - Lead source, status, probability
   - Contact person details

2. **sales_customers**
   - Customer master data
   - Credit limit, payment terms
   - Multi-location support

3. **sales_quotations**
   - Quote management
   - Line items with pricing
   - Expiry tracking

4. **sales_quotation_items**
   - Quote line items
   - Product/service with quantity
   - Discount, tax calculation

5. **sales_orders**
   - Sales order creation
   - Reference to quotation
   - Order status tracking

6. **sales_order_items**
   - Order line items
   - Stock allocation
   - Delivery scheduling

7. **sales_invoices**
   - Tax invoice generation
   - Line-item level details
   - Payment tracking

8. **sales_invoice_items**
   - Invoice line items
   - Quantity invoiced
   - Tax breakdown (CGST, SGST, IGST)

9. **sales_payments**
   - Payment receipts
   - Payment mode (cheque, transfer, cash)
   - Reconciliation status

10. **sales_debit_notes**
    - Customer debit notes
    - Reason tracking

11. **sales_credit_notes**
    - Customer credit notes
    - Reason tracking

12. **sales_returns**
    - Returned items
    - Return reason
    - Credit note link

13. **sales_delivery_notes**
    - Delivery tracking
    - Dispatch information
    - Proof of delivery

14. **sales_performance_metrics**
    - Salesperson metrics
    - Commission calculation
    - Targets vs actuals

15. **sales_customer_contacts**
    - Multiple contacts per customer
    - Role tracking (decision maker, finance, tech)

16. **sales_pricing_rules**
    - Volume-based discounts
    - Customer-specific pricing
    - Promotional offers

17. **sales_revenue_schedule**
    - Revenue recognition
    - Performance obligation tracking
    - GL posting schedule

18. **sales_customer_balances**
    - Customer account balance
    - AR aging
    - Receivables tracking

---

## Implementation Tasks

### 1. Backend (Go) - Day 1-2
- [ ] Create 18 database models
- [ ] Write database migrations
- [ ] Create sales_handler.go with CRUD operations
- [ ] Implement 35+ REST endpoints

### 2. Frontend (React/Next.js) - Day 2-3
- [ ] Create Sales module layout
- [ ] Implement Lead management UI
- [ ] Implement Quotation creation UI
- [ ] Implement Sales order UI
- [ ] Implement Invoice UI

### 3. API Integration - Day 3-4
- [ ] Connect frontend to backend APIs
- [ ] Implement data validation
- [ ] Error handling

### 4. Testing - Day 4-5
- [ ] Unit testing
- [ ] Integration testing
- [ ] Manual testing

### 5. GL Integration - Day 5
- [ ] AR account posting
- [ ] Revenue posting
- [ ] Tax posting

---

## Endpoints to Implement (35+)

### Lead Management (8 endpoints)
- POST /api/v1/sales/leads - Create lead
- GET /api/v1/sales/leads - List leads
- GET /api/v1/sales/leads/{id} - Get lead
- PUT /api/v1/sales/leads/{id} - Update lead
- DELETE /api/v1/sales/leads/{id} - Delete lead
- POST /api/v1/sales/leads/{id}/qualify - Qualify lead
- POST /api/v1/sales/leads/{id}/convert - Convert to customer
- GET /api/v1/sales/leads/stats - Lead statistics

### Customer Management (7 endpoints)
- POST /api/v1/sales/customers - Create customer
- GET /api/v1/sales/customers - List customers
- GET /api/v1/sales/customers/{id} - Get customer
- PUT /api/v1/sales/customers/{id} - Update customer
- DELETE /api/v1/sales/customers/{id} - Delete customer
- GET /api/v1/sales/customers/{id}/balance - AR balance
- POST /api/v1/sales/customers/{id}/contacts - Add contact

### Quotation Management (8 endpoints)
- POST /api/v1/sales/quotations - Create quotation
- GET /api/v1/sales/quotations - List quotations
- GET /api/v1/sales/quotations/{id} - Get quotation
- PUT /api/v1/sales/quotations/{id} - Update quotation
- DELETE /api/v1/sales/quotations/{id} - Delete quotation
- POST /api/v1/sales/quotations/{id}/convert - Convert to order
- POST /api/v1/sales/quotations/{id}/send - Send to customer
- POST /api/v1/sales/quotations/{id}/items - Add line item

### Sales Order Management (8 endpoints)
- POST /api/v1/sales/orders - Create order
- GET /api/v1/sales/orders - List orders
- GET /api/v1/sales/orders/{id} - Get order
- PUT /api/v1/sales/orders/{id} - Update order
- DELETE /api/v1/sales/orders/{id} - Delete order
- POST /api/v1/sales/orders/{id}/confirm - Confirm order
- POST /api/v1/sales/orders/{id}/invoice - Create invoice
- POST /api/v1/sales/orders/{id}/items - Add line item

### Invoice Management (4 endpoints)
- POST /api/v1/sales/invoices - Create invoice
- GET /api/v1/sales/invoices - List invoices
- GET /api/v1/sales/invoices/{id} - Get invoice
- POST /api/v1/sales/invoices/{id}/post - Post to GL

---

## Frontend Components (8)

### 1. SalesModule.tsx (Main container)
### 2. LeadManagement.tsx
### 3. CustomerManagement.tsx
### 4. QuotationManagement.tsx
### 5. SalesOrderManagement.tsx
### 6. InvoiceManagement.tsx
### 7. PaymentReceipt.tsx
### 8. PerformanceMetrics.tsx

---

## Success Criteria

✅ All 35+ endpoints working  
✅ CRUD operations for all tables  
✅ GL integration functioning  
✅ Frontend pages responsive  
✅ API documentation complete  
✅ Unit tests passing (80%+ coverage)  
✅ Integration tests passing  

---

## Start Date: November 25, 2025
## Target Completion: December 1, 2025
## Est. Lines of Code: 2,500+ (Backend + Frontend)
