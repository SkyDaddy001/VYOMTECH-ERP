# Database Setup & Sample Partner Logins Guide
**Date**: December 3, 2025  
**Status**: ✅ Complete  

---

## 📋 Overview

This guide covers:
1. Database migrations verification
2. Sample partner setup
3. Login credentials for testing
4. Testing procedures

---

## 🗄️ Database Migrations

### Total Migrations: 24 files

| Migration | Component | Status |
|-----------|-----------|--------|
| 001-005 | Core schema & modules | ✅ |
| 006-021 | Customization & modules | ✅ |
| 022 | External Partner System | ✅ |
| 023 | Partner Sources & Credit Policies | ✅ |
| 024 | Sample Partner Logins | ✅ |

### Required Tables

All the following tables must exist before using the sample logins:

```sql
✅ partners                           (36 columns)
✅ partner_users                      (13 columns)
✅ partner_leads                      (50+ columns)
✅ partner_lead_credits              (12 columns)
✅ partner_payouts                    (30+ columns)
✅ partner_payout_details            (10 columns)
✅ partner_activities                (15 columns)
✅ partner_sources                    (14 columns)
✅ partner_credit_policies            (36 columns)
✅ partner_credit_policy_mappings     (7 columns)
```

---

## 👥 Sample Partners & Users Created

### 1. Property Portal (🌐)
**Partner Type**: `portal`  
**Organization**: PropTech Portal  
**Contact**: Rajesh Kumar (rajesh@proptech.com)

| Email | Password | Role | Purpose |
|-------|----------|------|---------|
| portal_admin@proptech.com | password123 | admin | Full portal control |
| manager@proptech.com | password123 | lead_manager | Lead submission & management |

**Use Case**: Real estate portal submitting property leads

---

### 2. Channel Partner (🔗)
**Partner Type**: `channel_partner`  
**Organization**: BuildTech Solutions  
**Contact**: Priya Singh (priya@buildtech.in)

| Email | Password | Role | Purpose |
|-------|----------|------|---------|
| channel_admin@buildtech.in | password123 | admin | Full partner control |
| agent@buildtech.in | password123 | lead_manager | Agent lead submission |

**Use Case**: Channel partner/reseller submitting leads

---

### 3. Vendor Reference (🏭)
**Partner Type**: `vendor`  
**Organization**: Premium Vendors Inc  
**Contact**: Arun Patel (arun@vendors.in)

| Email | Password | Role | Purpose |
|-------|----------|------|---------|
| vendor_admin@vendors.in | password123 | admin | Full vendor control |
| user@vendors.in | password123 | viewer | View-only access |

**Use Case**: Vendor/supplier referring leads

---

### 4. Customer Reference (👤)
**Partner Type**: `customer`  
**Organization**: Happy Customers Ltd  
**Contact**: Neha Sharma (neha@customers.in)

| Email | Password | Role | Purpose |
|-------|----------|------|---------|
| customer_admin@customers.in | password123 | admin | Full customer control |
| referrer@customers.in | password123 | lead_manager | Submit referral leads |

**Use Case**: Customer/client referring leads

---

## 🔐 Role-Based Access Control

### Admin Role
- ✅ Full access to partner dashboard
- ✅ Manage users
- ✅ Submit leads
- ✅ View payouts
- ✅ Access analytics

### Lead Manager Role
- ✅ Submit leads
- ✅ View submitted leads
- ✅ Track lead status
- ✅ View earnings
- ❌ Manage users
- ❌ Access admin features

### Viewer Role
- ✅ View dashboard
- ✅ View analytics
- ✅ View earnings
- ❌ Submit leads
- ❌ Manage anything

---

## 🚀 Setup Instructions

### Option 1: Linux/macOS

```bash
# Make script executable
chmod +x scripts/setup-database.sh

# Run setup
./scripts/setup-database.sh

# Or with custom database credentials
DB_HOST=localhost \
DB_PORT=3306 \
DB_NAME=callcenter \
DB_USER=callcenter_user \
DB_PASSWORD=secure_app_pass \
./scripts/setup-database.sh
```

### Option 2: Windows

```batch
# Run the batch script
scripts\setup-database.bat

# Or set environment variables first
set DB_HOST=localhost
set DB_PORT=3306
set DB_NAME=callcenter
set DB_USER=callcenter_user
set DB_PASSWORD=secure_app_pass
scripts\setup-database.bat
```

### Option 3: Manual SQL Execution

```bash
# Run the migration directly
mysql -h localhost -u callcenter_user -p callcenter < migrations/024_sample_partner_logins.sql
```

---

## 🧪 Testing Sample Logins

### Test 1: Portal Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/partner/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "portal_admin@proptech.com",
    "password": "password123"
  }'
```

**Expected Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "partner_id": 1,
  "partner_type": "portal",
  "organization_name": "PropTech Portal",
  "role": "admin"
}
```

### Test 2: Channel Partner Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/partner/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "channel_admin@buildtech.in",
    "password": "password123"
  }'
```

### Test 3: Vendor Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/partner/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "vendor_admin@vendors.in",
    "password": "password123"
  }'
```

### Test 4: Customer Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/partner/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer_admin@customers.in",
    "password": "password123"
  }'
```

---

## 📊 Verify Setup

### Check Partner Count
```sql
SELECT partner_type, COUNT(*) as count FROM partners GROUP BY partner_type;
```

Expected Output:
```
| portal          | 1 |
| channel_partner | 1 |
| vendor          | 1 |
| customer        | 1 |
```

### Check User Count per Partner
```sql
SELECT p.organization_name, COUNT(pu.id) as users
FROM partners p
LEFT JOIN partner_users pu ON p.id = pu.partner_id
GROUP BY p.organization_name;
```

Expected Output:
```
| PropTech Portal         | 2 |
| BuildTech Solutions     | 2 |
| Premium Vendors Inc     | 2 |
| Happy Customers Ltd     | 2 |
```

### List All Sample Users
```sql
SELECT 
  p.partner_type,
  p.organization_name,
  pu.email,
  pu.role,
  pu.is_active
FROM partner_users pu
JOIN partners p ON pu.partner_id = p.id
ORDER BY p.partner_type, pu.email;
```

---

## 🔄 Lead Submission Flow

### Step 1: Partner Login
```bash
# Portal admin logs in
curl -X POST http://localhost:8080/api/v1/auth/partner/login \
  -d '{"email": "portal_admin@proptech.com", "password": "password123"}' \
  | jq '.token'
```

### Step 2: Submit Lead
```bash
TOKEN="<token_from_step1>"

curl -X POST http://localhost:8080/api/v1/partners/leads \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lead_source": "property_portal",
    "lead_title": "Residential Property Enquiry",
    "lead_description": "Customer interested in 2BHK apartment",
    "customer_name": "Rajesh Patel",
    "customer_email": "rajesh.patel@email.com",
    "customer_phone": "+91-9123456789",
    "property_type": "residential",
    "budget_min": 5000000,
    "budget_max": 7000000,
    "location": "Bangalore"
  }'
```

### Step 3: Track Lead Status
```bash
curl -X GET http://localhost:8080/api/v1/partners/leads \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.data[] | {id, status, lead_title, quality_score}'
```

---

## 💰 Credit Policy Testing

### View Available Policies
```bash
curl -X GET http://localhost:8080/api/v1/partners/1/credit-policies \
  -H "Authorization: Bearer $TOKEN"
```

### Calculate Lead Credit
```bash
# After lead approval
curl -X POST http://localhost:8080/api/v1/partners/leads/1/calculate-credit \
  -H "Authorization: Bearer $TOKEN"
```

Expected Response:
```json
{
  "partner_lead_id": 1,
  "calculated_credit": 12.60,
  "policy_id": 2,
  "policy_name": "Tiered Pricing - Volume Based",
  "tier_applied": "Tier 2",
  "base_credit": 12.00,
  "quality_bonus": 0.60,
  "reason": "Applied tiered policy with 5% quality bonus"
}
```

---

## 🔑 Password Information

### Default Password
```
password123
```

### Bcrypt Hash
```
$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm
```

### To Change Password
```bash
# Using API (if implemented)
curl -X POST http://localhost:8080/api/v1/partners/users/me/change-password \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "current_password": "password123",
    "new_password": "NewSecurePassword123!"
  }'
```

---

## 🐛 Troubleshooting

### Issue: "Access Denied" Error
```
Solution: Check database connection
  - Verify MySQL is running
  - Check DB credentials
  - Ensure database 'callcenter' exists
```

### Issue: "No partners found"
```
Solution: Run migration 024
  - Verify migration file exists
  - Check migration executed: 
    SELECT COUNT(*) FROM partners;
```

### Issue: "Unknown column 'X'"
```
Solution: Ensure all migrations are applied
  - Check all migration files are present
  - Verify migration 022 & 023 executed
  - Check table schema: DESCRIBE partners;
```

### Issue: "Login fails with correct credentials"
```
Solution: Verify partner_users table
  - Check user exists: 
    SELECT * FROM partner_users WHERE email='X';
  - Check user is active: is_active = TRUE
  - Check password_hash is set
```

---

## 📋 Migration Checklist

- [x] Migration 022: External Partner System
  - [x] partners table (36 columns)
  - [x] partner_users table (13 columns)
  - [x] partner_leads table
  - [x] partner_lead_credits table
  - [x] partner_payouts table
  - [x] partner_payout_details table
  - [x] partner_activities table

- [x] Migration 023: Partner Sources & Credit Policies
  - [x] partner_sources table (14 columns)
  - [x] partner_credit_policies table (36 columns)
  - [x] partner_credit_policy_mappings table (7 columns)

- [x] Migration 024: Sample Partner Logins
  - [x] 4 sample partners (portal, channel, vendor, customer)
  - [x] 8 sample partner users (2 per partner)
  - [x] All users with password hashes

---

## 🎓 Sample Data Summary

| Component | Count |
|-----------|-------|
| Sample Partners | 4 |
| Sample Users | 8 |
| Partners with 2 users each | 4 |
| Default Password | password123 |
| Password Hash Algorithm | bcrypt |
| Admin Users | 4 |
| Lead Manager Users | 3 |
| Viewer Users | 1 |

---

## 🔗 Related Documentation

- PHASE5_PARTNER_SOURCES_CREDIT_POLICIES.md - Credit policy details
- PHASE5_QUICK_REFERENCE.md - API quick reference
- migrations/022_external_partner_system.sql - Partner system schema
- migrations/023_partner_sources_and_credit_policies.sql - Policies schema
- migrations/024_sample_partner_logins.sql - Sample data

---

## ✅ Verification Commands

```bash
# Check all partner tables exist
mysql -u callcenter_user -p callcenter -e "
  SHOW TABLES LIKE 'partner%';
"

# Verify sample data loaded
mysql -u callcenter_user -p callcenter -e "
  SELECT 
    p.partner_type,
    p.organization_name,
    COUNT(pu.id) as users
  FROM partners p
  LEFT JOIN partner_users pu ON p.id = pu.partner_id
  GROUP BY p.partner_type, p.organization_name;
"

# List all sample users
mysql -u callcenter_user -p callcenter -e "
  SELECT 
    p.organization_name,
    pu.email,
    pu.role
  FROM partner_users pu
  JOIN partners p ON pu.partner_id = p.id
  ORDER BY p.organization_name, pu.email;
"
```

---

**Status**: ✅ Database setup complete with 4 sample partners and 8 test users

**Next Steps**: 
1. Start the application
2. Test login with sample credentials
3. Submit test leads
4. Verify credit calculation

