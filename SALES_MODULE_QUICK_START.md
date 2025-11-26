# SALES MODULE - QUICK START GUIDE

## 🚀 Getting Started

### 1. Backend Setup
```bash
# Build the project
go build ./cmd

# Run the server
./cmd

# Server listens on http://localhost:8080
```

### 2. Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Access at http://localhost:3000
```

### 3. Database Setup
```bash
# Apply migrations
psql -f migrations/009_sales_module_schema.sql

# Verify tables
psql -c "SELECT tablename FROM pg_tables WHERE schemaname='public' AND tablename LIKE 'sales_%'"
```

---

## 📡 API Quick Reference

### Base URL
```
http://localhost:8080/api/v1/sales
```

### Required Headers
```
X-Tenant-ID: your-tenant-id
X-User-ID: your-user-id
Authorization: Bearer your-jwt-token
Content-Type: application/json
```

### Common Operations

#### Create Lead
```bash
curl -X POST http://localhost:8080/api/v1/sales/leads \
  -H "X-Tenant-ID: tenant-1" \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "9876543210",
    "company_name": "ACME Corp",
    "industry": "Technology",
    "source": "website",
    "status": "new"
  }'
```

#### Create Customer
```bash
curl -X POST http://localhost:8080/api/v1/sales/customers \
  -H "X-Tenant-ID: tenant-1" \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "ACME Industries",
    "business_name": "ACME Industries Ltd",
    "business_type": "pvt_ltd",
    "primary_contact_name": "John Doe",
    "primary_email": "john@acme.com",
    "primary_phone": "9876543210",
    "gst_number": "27AAACT1234H1Z0",
    "credit_limit": 500000,
    "billing_city": "Bangalore"
  }'
```

#### Create Quotation
```bash
curl -X POST http://localhost:8080/api/v1/sales/quotations \
  -H "X-Tenant-ID: tenant-1" \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "customer-uuid",
    "quotation_date": "2025-11-25",
    "validity_date": "2025-12-25",
    "discount_percent": 5,
    "items": [
      {
        "item_name": "Product A",
        "quantity": 10,
        "unit_price": 1000,
        "tax_rate": 18
      }
    ],
    "notes": "Special pricing for bulk order"
  }'
```

#### Record Payment
```bash
curl -X POST http://localhost:8080/api/v1/sales/payments \
  -H "X-Tenant-ID: tenant-1" \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "invoice_id": "invoice-uuid",
    "payment_date": "2025-11-25",
    "payment_method": "bank_transfer",
    "amount_paid": 50000,
    "reference_number": "TXN123456",
    "remarks": "Payment received"
  }'
```

---

## 💾 Frontend Usage

### Access Sales Module
1. Navigate to `http://localhost:3000`
2. Login with credentials
3. Go to Dashboard → Sales
4. Choose tab: Leads | Customers | Quotations | Orders | Invoices | Payments

### Creating a Lead
1. Click "New Lead" button
2. Fill in form:
   - First Name, Last Name
   - Email, Phone
   - Company Name, Industry
   - Select Source and Status
3. Click "Create Lead"

### Creating a Customer
1. Click "New Customer" button
2. Fill in form:
   - Customer Name, Business Name
   - Business Type, Industry
   - Contact Info (Name, Email, Phone)
   - GST Number, Credit Limit
3. Click "Create Customer"

### Creating a Quotation
1. Click "New Quotation" button
2. Select Customer
3. Set dates
4. Add line items (click "+ Add Item")
5. Set tax rate and quantity
6. Review totals
7. Click "Create Quotation"

### Converting Quote to Order
1. Go to Quotations tab
2. Click "→" button on accepted quotation
3. Confirm conversion
4. Order created automatically

### Recording Payment
1. Go to Payments tab
2. Click "Record Payment"
3. Select Invoice
4. Enter Amount Paid
5. Select Payment Method
6. Enter Reference Number
7. Click "Record Payment"
8. Download receipt

---

## 🔍 Testing Workflows

### Complete Lead-to-Cash Flow
```
1. Create Lead (Leads tab)
2. Create Customer (Customers tab)
3. Create Quotation (Quotations tab)
4. Convert to Order (Orders tab)
5. Create Invoice (Invoices tab)
6. Record Payment (Payments tab)
```

### Status Transitions
```
Lead:
  new → contacted → qualified → negotiation → converted/lost

Order:
  draft → confirmed → invoiced → delivered

Invoice:
  unpaid → partially_paid → paid
```

### Tax Calculations
```
Subtotal = Σ(Qty × Price)
Tax = Σ(Qty × Price × Tax%)
Discount = Subtotal × Discount%
Total = Subtotal + Tax - Discount
```

---

## 🐛 Troubleshooting

### Backend Won't Start
```bash
# Check if port 8080 is available
lsof -i :8080

# Check database connection
psql -h localhost -U user -d vyomtech -c "SELECT 1"

# View logs
tail -f logs/app.log
```

### API Returns 401
- Verify JWT token is valid
- Check X-Tenant-ID header is set
- Verify Authorization header format: `Bearer {token}`

### Frontend Components Not Loading
```bash
# Clear Next.js cache
rm -rf .next

# Reinstall dependencies
npm install

# Rebuild
npm run build
```

### Database Issues
```sql
-- Check if tables exist
SELECT * FROM pg_tables WHERE schemaname='public' AND tablename LIKE 'sales_%';

-- Check row counts
SELECT tablename, COUNT(*) FROM pg_tables WHERE schemaname='public' GROUP BY tablename;

-- Reset for testing
DELETE FROM sales_leads WHERE tenant_id = 'test-tenant';
```

---

## 📊 Useful Queries

### Get All Leads
```bash
curl -X GET "http://localhost:8080/api/v1/sales/leads?limit=100" \
  -H "X-Tenant-ID: tenant-1"
```

### Get Lead Statistics
```bash
curl -X GET "http://localhost:8080/api/v1/sales/leads/stats" \
  -H "X-Tenant-ID: tenant-1"
```

### List Orders by Status
```bash
curl -X GET "http://localhost:8080/api/v1/sales/orders?status=confirmed" \
  -H "X-Tenant-ID: tenant-1"
```

### Get Invoice Details
```bash
curl -X GET "http://localhost:8080/api/v1/sales/invoices/{id}" \
  -H "X-Tenant-ID: tenant-1"
```

---

## 📈 Performance Tips

### Frontend
- Use React DevTools to monitor component renders
- Enable production build for better performance
- Use browser DevTools Network tab to monitor API calls

### Backend
- Monitor database query performance with `EXPLAIN ANALYZE`
- Use connection pooling for database
- Enable caching for frequently accessed data

### Database
- Regularly analyze tables: `ANALYZE sales_leads;`
- Monitor index usage: `SELECT * FROM pg_stat_user_indexes;`
- Check table sizes: `SELECT tablename, pg_size_pretty(pg_total_relation_size(tablename)) FROM pg_tables WHERE schemaname='public';`

---

## 🔐 Security Checklist

- [x] All endpoints require JWT authentication
- [x] Multi-tenant isolation at database level
- [x] Parameterized queries prevent SQL injection
- [x] Input validation on frontend and backend
- [x] CORS configured for allowed origins
- [x] X-Tenant-ID header validation
- [x] Soft deletes preserve audit trail
- [x] HTTPS recommended for production

---

## 📝 Environment Variables

```env
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=vyomtech

# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

---

## 🎓 Learning Resources

### Backend
- **Handlers**: `internal/handlers/sales_*.go`
- **Models**: `internal/models/sales.go`
- **Router**: `pkg/router/router.go`
- **Service**: `internal/services/sales_service.go`

### Frontend
- **Main Page**: `frontend/app/dashboard/sales/page.tsx`
- **Lead Component**: `frontend/components/modules/Sales/LeadManagement.tsx`
- **Customer Component**: `frontend/components/modules/Sales/CustomerManagement.tsx`
- **Quotation Component**: `frontend/components/modules/Sales/QuotationManagement.tsx`
- **Order Component**: `frontend/components/modules/Sales/SalesOrderManagement.tsx`
- **Invoice Component**: `frontend/components/modules/Sales/InvoiceManagement.tsx`
- **Payment Component**: `frontend/components/modules/Sales/PaymentReceipt.tsx`

---

## 📞 Support

For issues or questions:
1. Check the comprehensive documentation: `SALES_MODULE_COMPLETE.md`
2. Review API reference in documentation
3. Check database schema: `migrations/009_sales_module_schema.sql`
4. Review component source code for implementation details

---

**Last Updated**: November 25, 2025  
**Version**: 1.0.0  
**Status**: Production Ready ✅
