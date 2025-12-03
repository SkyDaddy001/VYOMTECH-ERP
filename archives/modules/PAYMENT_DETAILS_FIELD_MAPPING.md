# Payment Details Field Mapping

## Overview
This document maps real-world payment transaction details to the PropertyPaymentReceipt data model and database schema. It supports comprehensive payment tracking with receipt generation, payment categorization, and financial reconciliation.

---

## Real-World Payment Example
```
BLOCK B - UNIT 406 (LML - THE LEAGUE ONE)
Client Name: Dr. Nagaraju & Ms. Sakthi Abirami .N
Payment Date: 15-Apr-24
Payment Mode: Cheque/558471/PNB/15-04-2024
Paid By: CUSTOMER PAYMENT
Receipt No: TLO/001
Towards: APARTMENT COST
Amount: 95,238
Received In: Acc No: 7729200809
```

---

## Field Mapping Reference

| Real-World Field | Model Field | Database Column | Type | Comment |
|------------------|-------------|-----------------|------|---------|
| BLOCK B - UNIT 406 | UnitID | unit_id | VARCHAR(36) | Foreign key to property_unit |
| Client Name | CustomerName | customer_name | VARCHAR(200) | Denormalized customer name for reporting |
| Payment Date | PaymentDate | payment_date | DATE | Date payment was received |
| Receipt No | ReceiptNumber | receipt_number | VARCHAR(50) | Unique receipt identifier per tenant |
| Cheque No | ChequeNumber | cheque_number | VARCHAR(50) | Check number (for cheque mode) |
| Bank | BankName | bank_name | VARCHAR(200) | Issuing bank name (for cheque mode) |
| Cheque Date | ChequeDate | cheque_date | DATE | Date on cheque (for cheque mode) |
| Payment Mode | PaymentMode | payment_mode | VARCHAR(50) | CASH/CHEQUE/NEFT/RTGS/ONLINE/DD |
| Amount | PaymentAmount | payment_amount | DECIMAL(18,2) | Amount received |
| Towards | TowardsDescription | towards_description | VARCHAR(200) | Payment category/purpose |
| Paid By | PaidBy | paid_by | VARCHAR(100) | Customer/Agent/Representative |
| Acc No | ReceivedInBankAccount | received_in_bank_account | VARCHAR(50) | Account where payment received |
| NEFT/RTGS Ref | TransactionID | transaction_id | VARCHAR(100) | Bank transaction reference |
| Online Ref | AccountNumber | account_number | VARCHAR(50) | Account number (for online) |

---

## Payment Categories (towards_description)

Standard payment categorization enums:

```
APARTMENT_COST          - Base apartment/flat cost
PARKING                 - Car parking charges
MAINTENANCE             - Annual maintenance charges
CORPUS                  - Sinking fund/corpus charges
REGISTRATION            - Registration and stamp duty
SOCIETY_REGISTRATION    - Society registration charges
TRANSFER_CHARGES        - Transfer/mutation charges
BROKERAGE               - Agent brokerage fees
OTHER_WORKS             - Other development charges
INSTALLMENT_PENALTY     - Late payment penalty
INTEREST_CHARGES        - Interest on delayed payment
ADJUSTMENT              - Payment adjustment/reverse
```

---

## Payment Modes (payment_mode)

Supported payment methods with associated fields:

### CASH
- **Fields Used**: None
- **Validation**: payment_amount must be positive
- **GL Impact**: Direct cash receipt to bank account

### CHEQUE
- **Fields Used**: bank_name, cheque_number, cheque_date
- **Validation**: All three fields required
- **GL Impact**: Cheque clearing → bank account
- **Flow**: Cheque Received → Presented → Cleared/Bounced

### NEFT (National Electronic Funds Transfer)
- **Fields Used**: transaction_id, account_number
- **Validation**: Both fields required
- **GL Impact**: Direct NEFT receipt to bank account
- **Status Flow**: Initiated → In-Progress → Completed/Failed

### RTGS (Real Time Gross Settlement)
- **Fields Used**: transaction_id, account_number
- **Validation**: Both fields required
- **GL Impact**: Direct RTGS receipt to bank account
- **Status Flow**: Initiated → In-Progress → Completed/Failed

### ONLINE
- **Fields Used**: account_number, transaction_id
- **Validation**: At least transaction_id required
- **GL Impact**: Online gateway → bank account
- **Status Flow**: Initiated → Confirmed → Settled

### DD (Demand Draft)
- **Fields Used**: bank_name, cheque_number (use as DD ref), cheque_date
- **Validation**: bank_name and reference required
- **GL Impact**: DD clearing → bank account

---

## Payment Status Lifecycle

```
PENDING
  ↓
RECEIVED (payment received, not yet confirmed)
  ↓
PROCESSED (validated, allocated to unit/installment)
  ↓
CLEARED (for cheque: cleared; for electronic: settled)
  ↓
POSTED_TO_GL (posted to general ledger for accounting)

OR (Error path):

PENDING → BOUNCED (cheque bounced, NEFT/RTGS failed)
       → CANCELLED (manual cancellation)
```

---

## Database Schema (Enhanced)

### Table: property_payment_receipt

```sql
CREATE TABLE `property_payment_receipt` (
    -- Primary & Foreign Keys
    `id` VARCHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `customer_id` VARCHAR(36) NOT NULL,
    `unit_id` VARCHAR(36) NOT NULL,
    `installment_id` VARCHAR(36),
    
    -- Core Receipt Information
    `receipt_number` VARCHAR(50) NOT NULL,
    `receipt_date` DATE,
    `payment_date` DATE NOT NULL,
    
    -- Payment Information
    `payment_mode` VARCHAR(50) NOT NULL,
    `payment_amount` DECIMAL(18,2),
    `payment_status` VARCHAR(50) DEFAULT 'PENDING',
    
    -- Amount Reconciliation
    `installment_amount_due` DECIMAL(18,2),
    `shortfall_amount` DECIMAL(18,2),
    `excess_amount` DECIMAL(18,2),
    
    -- Cheque Details (for CHEQUE mode)
    `bank_name` VARCHAR(200),
    `cheque_number` VARCHAR(50),
    `cheque_date` DATE,
    
    -- Electronic Transfer Details (NEFT/RTGS/Online)
    `transaction_id` VARCHAR(100),
    `account_number` VARCHAR(50),
    
    -- Receipt Categorization
    `towards_description` VARCHAR(200),
    `received_in_bank_account` VARCHAR(50),
    `paid_by` VARCHAR(100),
    
    -- Customer Information (Denormalized)
    `customer_name` VARCHAR(200),
    
    -- Accounting
    `gl_account_id` VARCHAR(36),
    
    -- Metadata
    `remarks` TEXT,
    `created_by` VARCHAR(36),
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign Keys
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`),
    FOREIGN KEY (`customer_id`) REFERENCES `property_customer_profile`(`id`),
    FOREIGN KEY (`unit_id`) REFERENCES `property_unit`(`id`),
    FOREIGN KEY (`installment_id`) REFERENCES `installment`(`id`),
    
    -- Indexes
    UNIQUE KEY `unique_receipt` (`tenant_id`, `receipt_number`),
    KEY `idx_tenant` (`tenant_id`),
    KEY `idx_customer` (`customer_id`),
    KEY `idx_unit` (`unit_id`),
    KEY `idx_date` (`payment_date`),
    KEY `idx_status` (`payment_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
COMMENT='Payment receipts with detailed transaction tracking';
```

---

## Go Model Structure

### PropertyPaymentReceipt Data Model

```go
type PropertyPaymentReceipt struct {
    ID                      string       // UUID
    TenantID                string       // Multi-tenant isolation
    CustomerID              string       // FK to property_customer_profile
    CustomerName            string       // Denormalized for reporting
    UnitID                  string       // FK to property_unit
    InstallmentID           string       // FK to installment (nullable)
    
    // Receipt Details
    ReceiptNumber           string       // Unique per tenant
    ReceiptDate             *time.Time   // When receipt was generated
    PaymentDate             *time.Time   // When payment was received
    
    // Payment Information
    PaymentMode             string       // CASH/CHEQUE/NEFT/RTGS/ONLINE/DD
    PaymentAmount           float64      // Decimal amount
    PaymentStatus           string       // PENDING/RECEIVED/PROCESSED/CLEARED/BOUNCED/CANCELLED
    InstallmentAmountDue    float64      // Expected amount for installment
    ShortfallAmount         float64      // If payment < due amount
    ExcessAmount            float64      // If payment > due amount
    
    // Payment Mode Details (Cheque)
    BankName                string       // Issuing bank
    ChequeNumber            string       // Check/DD reference
    ChequeDate              *time.Time   // Date on check
    
    // Electronic Payment Details
    TransactionID           string       // NEFT/RTGS/Online reference
    AccountNumber           string       // Source/destination account
    
    // Receipt Categorization
    TowardsDescription      string       // APARTMENT_COST, MAINTENANCE, CORPUS, etc.
    ReceivedInBankAccount   string       // Account receiving payment
    PaidBy                  string       // Customer/Agent/Representative
    
    // Accounting Integration
    GLAccountID             string       // GL account for posting
    
    // Metadata
    Remarks                 string       // Additional notes
    CreatedBy               string       // User who created receipt
    CreatedAt               time.Time    // Creation timestamp
    UpdatedAt               time.Time    // Last update timestamp
    DeletedAt               sql.NullTime // Soft delete support
}
```

### CreatePaymentReceiptRequest API Input

```go
type CreatePaymentReceiptRequest struct {
    // Core References
    CustomerID              string       // Required
    CustomerName            string       // Optional, will denormalize
    UnitID                  string       // Required
    InstallmentID           string       // Optional
    
    // Receipt Details
    PaymentDate             *time.Time   // Required
    PaymentMode             string       // Required: CASH/CHEQUE/NEFT/RTGS/ONLINE/DD
    PaymentAmount           float64      // Required
    
    // Payment Mode Details
    BankName                string       // For CHEQUE/DD
    ChequeNumber            string       // For CHEQUE/DD
    ChequeDate              *time.Time   // For CHEQUE/DD
    TransactionID           string       // For NEFT/RTGS/ONLINE
    AccountNumber           string       // For electronic transfers
    
    // Receipt Details
    TowardsDescription      string       // Payment category
    ReceivedInBankAccount   string       // Receiving account
    PaidBy                  string       // Payer identifier
    
    // Additional
    Remarks                 string       // Optional notes
}
```

---

## Implementation Flow

### 1. Create Payment Receipt
```
POST /api/v1/projects/{projectId}/payment-receipts
```

**Request Body Example:**
```json
{
  "customer_id": "cust_001",
  "customer_name": "Dr. Nagaraju & Ms. Sakthi Abirami .N",
  "unit_id": "unit_406",
  "installment_id": "inst_001",
  "payment_date": "2024-04-15",
  "payment_mode": "CHEQUE",
  "payment_amount": 95238.00,
  "bank_name": "PNB",
  "cheque_number": "558471",
  "cheque_date": "2024-04-15",
  "towards_description": "APARTMENT_COST",
  "received_in_bank_account": "7729200809",
  "paid_by": "CUSTOMER",
  "remarks": "Payment received for unit 406"
}
```

**Response:**
```json
{
  "id": "receipt_001",
  "receipt_number": "TLO/001",
  "payment_status": "RECEIVED",
  "created_at": "2024-04-15T10:30:00Z"
}
```

### 2. Payment Reconciliation
- Compare payment_amount with installment_amount_due
- Calculate shortfall_amount (if payment < due)
- Calculate excess_amount (if payment > due)
- Generate GL posting entry

### 3. GL Integration Points
- **GL Account ID**: Links payment to specific accounting code
- **Payment Status**: CLEARED triggers GL posting
- **Amount**: Posted to appropriate revenue/bank account

---

## Multi-Tenant Isolation

All payment records include `tenant_id` for secure multi-tenant isolation:

```go
// Database queries must include tenant_id
SELECT * FROM property_payment_receipt 
WHERE tenant_id = 'tenant_xyz' 
  AND customer_id = 'cust_001'
```

---

## Payment Validation Rules

| Field | Validation | Error |
|-------|-----------|-------|
| payment_amount | > 0 | "Amount must be positive" |
| payment_date | ≤ today | "Payment date cannot be future" |
| cheque_number | Required if mode=CHEQUE | "Cheque number required" |
| transaction_id | Required if mode=NEFT/RTGS | "Transaction ID required" |
| towards_description | Standard enum or custom | Validate against allowed values |
| customer_id | Must exist | Foreign key constraint |
| unit_id | Must exist | Foreign key constraint |
| receipt_number | Unique per tenant | Unique constraint |

---

## Reporting & Analytics

### Common Queries

**Payment Summary by Towards Description:**
```sql
SELECT towards_description, 
       COUNT(*) as count,
       SUM(payment_amount) as total_amount
FROM property_payment_receipt
WHERE tenant_id = 'tenant_xyz'
  AND payment_date BETWEEN ? AND ?
GROUP BY towards_description
ORDER BY total_amount DESC;
```

**Payment Status Breakdown:**
```sql
SELECT payment_status, 
       COUNT(*) as count,
       SUM(payment_amount) as amount
FROM property_payment_receipt
WHERE tenant_id = 'tenant_xyz'
GROUP BY payment_status;
```

**Overdue Payments (Shortfall):**
```sql
SELECT customer_id, unit_id, SUM(shortfall_amount) as total_shortfall
FROM property_payment_receipt
WHERE tenant_id = 'tenant_xyz'
  AND shortfall_amount > 0
  AND payment_date < DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY customer_id, unit_id;
```

---

## File Changes Summary

### 1. Migration Updated: `migrations/022_project_management_system.sql`
- **Added columns to property_payment_receipt table:**
  - `customer_name` (VARCHAR 200) - Denormalized
  - `account_number` (VARCHAR 50) - For electronic transfers
  - `towards_description` (VARCHAR 200) - Payment category
  - `received_in_bank_account` (VARCHAR 50) - Receiving account
  - `paid_by` (VARCHAR 100) - Payer type

### 2. Models Updated: `internal/models/project_management.go`
- **Enhanced PropertyPaymentReceipt struct:**
  - Added `CustomerName` field
  - Added `AccountNumber` field
  - Added `TowardsDescription` field
  - Added `ReceivedInBankAccount` field
  - Added `PaidBy` field
  - Organized fields into logical sections with comments

- **Enhanced CreatePaymentReceiptRequest struct:**
  - Added `CustomerName` field
  - Added `AccountNumber` field
  - Added `TowardsDescription` field
  - Added `ReceivedInBankAccount` field
  - Added `PaidBy` field
  - Organized with section comments

### 3. New Document: `PAYMENT_DETAILS_FIELD_MAPPING.md`
- Complete payment field mapping reference
- Real-world example with field mappings
- Payment categories and modes
- Database schema documentation
- Go model structures
- Implementation flow examples
- Multi-tenant isolation details
- Validation rules
- Reporting queries

---

## Next Steps

1. **Service Layer**: Implement payment receipt CRUD operations
   - CreatePaymentReceipt service
   - UpdatePaymentStatus service
   - GetPaymentReceipts query
   - GL posting integration

2. **API Handlers**: REST endpoints for payment management
   - POST /api/v1/projects/{projectId}/payment-receipts
   - GET /api/v1/projects/{projectId}/payment-receipts
   - GET /api/v1/projects/{projectId}/payment-receipts/{receiptId}
   - PUT /api/v1/projects/{projectId}/payment-receipts/{receiptId}
   - DELETE /api/v1/projects/{projectId}/payment-receipts/{receiptId}

3. **GL Integration**: Payment posting to GL
   - Implement GL posting logic on payment status change to CLEARED
   - Create GL transaction entry with debit/credit mapping
   - Link payment receipt to GL transaction

4. **Frontend**: Payment receipt UI
   - Payment receipt entry form
   - Receipt list view with filters (status, date range, towards)
   - Receipt detail view
   - Payment reconciliation dashboard
   - Amount vs. due analysis

---

## Status

✅ **COMPLETE**: Payment details model enhancement
✅ **COMPLETE**: Database schema update
✅ **COMPLETE**: Field mapping documentation
🔄 **PENDING**: Service layer implementation
🔄 **PENDING**: API handlers
🔄 **PENDING**: GL integration
🔄 **PENDING**: Frontend UI

---

## Related Documents

- `PROJECT_MANAGEMENT_INDEX.md` - System overview
- `PROJECT_MANAGEMENT_INTEGRATION.md` - Architecture details
- `PROJECT_MANAGEMENT_QUICK_REFERENCE.md` - Developer guide
- `CUSTOMER_FIELD_MAPPING.md` - Customer details mapping
