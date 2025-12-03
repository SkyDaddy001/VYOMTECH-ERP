# Project Management API Documentation

## Overview

The Project Management API provides comprehensive endpoints for managing real estate projects including customer profiles, area statements, cost sheets, bank financing, disbursements, and payment stages.

**Base URL**: `/api/v1/project-management`
**Authentication**: Required (Bearer Token)
**Tenant Isolation**: Automatic (based on tenant context)

---

## Table of Contents

1. [Customer Profile Endpoints](#customer-profile-endpoints)
2. [Area Statement Endpoints](#area-statement-endpoints)
3. [Cost Sheet Endpoints](#cost-sheet-endpoints)
4. [Cost Configuration Endpoints](#cost-configuration-endpoints)
5. [Bank Financing Endpoints](#bank-financing-endpoints)
6. [Disbursement Schedule Endpoints](#disbursement-schedule-endpoints)
7. [Payment Stage Endpoints](#payment-stage-endpoints)
8. [Reporting Endpoints](#reporting-endpoints)

---

## Customer Profile Endpoints

### Create Customer Profile
**Endpoint**: `POST /api/v1/project-management/customers`

Creates a new customer profile with comprehensive KYC information including primary applicant, up to 3 co-applicants, addresses, employment details, and financing information.

**Request Body**:
```json
{
  "customer_code": "CUST-001",
  "unit_id": "unit-uuid",
  "first_name": "John",
  "middle_name": "Michael",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_primary": "+91-9999999999",
  "phone_secondary": "+91-8888888888",
  "company_name": "Tech Corp",
  "designation": "Manager",
  "pan_number": "ABCDE1234F",
  "aadhar_number": "123456789012",
  "communication_address_line1": "123 Main Street",
  "communication_city": "Bangalore",
  "communication_state": "Karnataka",
  "communication_country": "India",
  "communication_zip": "560001",
  "permanent_address_line1": "456 Oak Avenue",
  "permanent_city": "Bangalore",
  "permanent_state": "Karnataka",
  "permanent_country": "India",
  "permanent_zip": "560001",
  "profession": "Software Engineer",
  "employer_name": "TechCorp Ltd",
  "employment_type": "Full-Time",
  "monthly_income": 150000,
  "customer_type": "INDIVIDUAL"
}
```

**Response** (201 Created):
```json
{
  "id": "customer-uuid",
  "tenant_id": "tenant-uuid",
  "customer_code": "CUST-001",
  "first_name": "John",
  "email": "john.doe@example.com",
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

### Get Customer Profile
**Endpoint**: `GET /api/v1/project-management/customers/{id}`

Retrieves a specific customer profile by ID.

**Response** (200 OK):
```json
{
  "id": "customer-uuid",
  "tenant_id": "tenant-uuid",
  "customer_code": "CUST-001",
  "first_name": "John",
  "email": "john.doe@example.com",
  "phone_primary": "+91-9999999999",
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

---

## Area Statement Endpoints

### Create Area Statement
**Endpoint**: `POST /api/v1/project-management/area-statements`

Creates a complete area statement for a unit with detailed area measurements including carpet area, plinth area, SBUA, and additional areas.

**Request Body**:
```json
{
  "project_id": "project-uuid",
  "block_id": "block-uuid",
  "unit_id": "unit-uuid",
  "apt_no": "406",
  "floor": "4",
  "unit_type": "2BHK",
  "facing": "NORTH",
  "rera_carpet_area_sqft": 1075,
  "rera_carpet_area_sqm": 100,
  "carpet_area_with_balcony_sqft": 1175,
  "carpet_area_with_balcony_sqm": 109,
  "plinth_area_sqft": 1250,
  "plinth_area_sqm": 116.13,
  "sbua_sqft": 1500,
  "sbua_sqm": 139.35,
  "uds_per_sqft": 0.5,
  "uds_total_sqft": 750,
  "balcony_area_sqft": 100,
  "balcony_area_sqm": 9.29,
  "utility_area_sqft": 0,
  "utility_area_sqm": 0,
  "garden_area_sqft": 0,
  "garden_area_sqm": 0,
  "parking_area_sqft": 40,
  "parking_area_sqm": 3.72,
  "common_area_sqft": 0,
  "common_area_sqm": 0,
  "alloted_to": "John Doe",
  "key_holder": "John Doe",
  "percentage_allocation": 0.85,
  "noc_taken": "YES",
  "noc_date": "2025-10-15",
  "noc_document_url": "https://example.com/noc.pdf",
  "area_type": "CARPET_AREA",
  "description": "Area statement for Unit 406, Block A"
}
```

**Response** (201 Created):
```json
{
  "id": "area-statement-uuid",
  "tenant_id": "tenant-uuid",
  "project_id": "project-uuid",
  "unit_id": "unit-uuid",
  "apt_no": "406",
  "floor": "4",
  "unit_type": "2BHK",
  "facing": "NORTH",
  "rera_carpet_area_sqft": 1075,
  "sbua_sqft": 1500,
  "noc_taken": "YES",
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

---

## Cost Sheet Endpoints

### Update Cost Sheet
**Endpoint**: `POST /api/v1/project-management/cost-sheets`

Updates or creates a cost sheet for a unit with comprehensive cost breakdown including car parking, statutory charges, other charges, and GST.

**Request Body**:
```json
{
  "unit_id": "unit-uuid",
  "block_name": "Block A",
  "sbua": 1500,
  "rate_per_sqft": 5000,
  "car_parking_cost": 750000,
  "plc": 100000,
  "statutory_approval_charge": 200000,
  "legal_documentation_charge": 50000,
  "amenities_equipment_charge": 300000,
  "other_charges_1": 30000,
  "other_charges_1_name": "CMWSSB",
  "other_charges_1_type": "PER_SQFT",
  "other_charges_2": 20000,
  "other_charges_2_name": "Water Tax",
  "other_charges_2_type": "LUMPSUM",
  "other_charges_3": 15000,
  "other_charges_3_name": "Electricity Deposit",
  "other_charges_3_type": "LUMPSUM",
  "other_charges_4": 10000,
  "other_charges_4_name": "Maintenance",
  "other_charges_4_type": "PER_SQFT",
  "other_charges_5": 5000,
  "other_charges_5_name": "Club Membership",
  "other_charges_5_type": "LUMPSUM",
  "apartment_cost_excluding_govt": 7500000,
  "actual_sold_price_excluding_govt": 7500000,
  "gst_applicable": true,
  "gst_percentage": 5,
  "gst_amount": 375000,
  "grand_total": 7875000,
  "club_membership": 100000,
  "registration_charge": 50000
}
```

**Response** (200 OK):
```json
{
  "message": "Cost sheet updated successfully",
  "unit_id": "unit-uuid"
}
```

---

## Cost Configuration Endpoints

### Create Cost Configuration
**Endpoint**: `POST /api/v1/project-management/cost-configurations`

Creates a project-wide cost configuration for charges like CMWSSB, Water Tax, etc. with support for per-sqft or lump-sum calculation.

**Request Body**:
```json
{
  "project_id": "project-uuid",
  "config_name": "CMWSSB",
  "config_type": "OTHER_CHARGE_1",
  "charge_type": "PER_SQFT",
  "charge_amount": 5,
  "display_order": 1,
  "is_mandatory": true,
  "applicable_for_unit_type": "1BHK,2BHK,3BHK",
  "description": "CMWSSB charges calculated per square foot"
}
```

**Response** (201 Created):
```json
{
  "id": "config-uuid",
  "tenant_id": "tenant-uuid",
  "project_id": "project-uuid",
  "config_name": "CMWSSB",
  "charge_type": "PER_SQFT",
  "charge_amount": 5,
  "active": true,
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

---

## Bank Financing Endpoints

### Create Bank Financing
**Endpoint**: `POST /api/v1/project-management/bank-financing`

Initiates bank financing tracking with sanction amount, bank details, and commitment tracking.

**Request Body**:
```json
{
  "project_id": "project-uuid",
  "block_id": "block-uuid",
  "unit_id": "unit-uuid",
  "customer_id": "customer-uuid",
  "apt_no": "406",
  "block_name": "Block A",
  "apartment_cost": 7500000,
  "bank_name": "HDFC Bank",
  "banker_reference_no": "HDFC/2025/001234",
  "sanctioned_amount": 6000000,
  "sanctioned_date": "2025-11-15",
  "total_commitment": 7500000
}
```

**Response** (201 Created):
```json
{
  "id": "financing-uuid",
  "tenant_id": "tenant-uuid",
  "project_id": "project-uuid",
  "unit_id": "unit-uuid",
  "apt_no": "406",
  "apartment_cost": 7500000,
  "sanctioned_amount": 6000000,
  "total_disbursed_amount": 0,
  "remaining_disbursement": 6000000,
  "disbursement_status": "PENDING",
  "collection_status": "PENDING",
  "outstanding_amount": 7500000,
  "active": true,
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

---

## Disbursement Schedule Endpoints

### Create Disbursement Schedule
**Endpoint**: `POST /api/v1/project-management/disbursement-schedule`

Creates a disbursement schedule linked to project milestones.

**Request Body**:
```json
{
  "financing_id": "financing-uuid",
  "unit_id": "unit-uuid",
  "customer_id": "customer-uuid",
  "disbursement_no": 1,
  "expected_disbursement_date": "2025-12-15",
  "expected_disbursement_amount": 3000000,
  "disbursement_percentage": 50,
  "linked_milestone_id": "milestone-uuid",
  "milestone_stage": "FOUNDATION"
}
```

**Response** (201 Created):
```json
{
  "id": "disbursement-uuid",
  "tenant_id": "tenant-uuid",
  "financing_id": "financing-uuid",
  "unit_id": "unit-uuid",
  "disbursement_no": 1,
  "expected_disbursement_date": "2025-12-15T00:00:00Z",
  "expected_disbursement_amount": 3000000,
  "disbursement_percentage": 50,
  "milestone_stage": "FOUNDATION",
  "disbursement_status": "PENDING",
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

### Update Disbursement
**Endpoint**: `PUT /api/v1/project-management/disbursement/{id}`

Records actual disbursement with bank details and reference numbers.

**Request Body**:
```json
{
  "actual_disbursement_date": "2025-12-15",
  "actual_disbursement_amount": 3000000,
  "disbursement_status": "CLEARED",
  "cheque_no": "CHQ-123456",
  "bank_reference_no": "HDFC/TXN/987654",
  "neft_ref_id": "HDFC0001234567"
}
```

**Response** (200 OK):
```json
{
  "message": "Disbursement updated successfully",
  "disbursement_id": "disbursement-uuid"
}
```

---

## Payment Stage Endpoints

### Create Payment Stage
**Endpoint**: `POST /api/v1/project-management/payment-stages`

Sets up payment stages with percentage-based cost allocation.

**Request Body**:
```json
{
  "project_id": "project-uuid",
  "unit_id": "unit-uuid",
  "customer_id": "customer-uuid",
  "stage_name": "BOOKING",
  "stage_number": 1,
  "stage_description": "Booking stage - 20% of cost",
  "stage_percentage": 20,
  "apartment_cost": 7500000,
  "due_date": "2025-12-20"
}
```

**Response** (201 Created):
```json
{
  "id": "stage-uuid",
  "tenant_id": "tenant-uuid",
  "project_id": "project-uuid",
  "unit_id": "unit-uuid",
  "stage_name": "BOOKING",
  "stage_number": 1,
  "stage_percentage": 20,
  "stage_due_amount": 1500000,
  "apartment_cost": 7500000,
  "amount_due": 1500000,
  "amount_received": 0,
  "amount_pending": 1500000,
  "collection_status": "PENDING",
  "due_date": "2025-12-20T00:00:00Z",
  "created_at": "2025-12-03T10:30:00Z",
  "updated_at": "2025-12-03T10:30:00Z"
}
```

### Record Payment Collection
**Endpoint**: `PUT /api/v1/project-management/payment-stages/{id}/collection`

Records collection/payment for a payment stage.

**Request Body**:
```json
{
  "amount_received": 1500000,
  "payment_received_date": "2025-12-18",
  "payment_mode": "NEFT",
  "reference_no": "NEFT/TXN/123456",
  "collection_status": "COMPLETED"
}
```

**Response** (200 OK):
```json
{
  "message": "Payment collected successfully",
  "stage_id": "stage-uuid"
}
```

---

## Reporting Endpoints

### Bank Financing Report
**Endpoint**: `GET /api/v1/project-management/reports/bank-financing?project_id={project_id}`

Generates bank financing report with sanctioned amounts, disbursements, and collections.

**Query Parameters**:
- `project_id` (optional): Filter by project

**Response** (200 OK):
```json
{
  "count": 5,
  "data": [
    {
      "id": "financing-uuid",
      "unit_id": "unit-uuid",
      "apt_no": "406",
      "apartment_cost": 7500000,
      "sanctioned_amount": 6000000,
      "total_disbursed_amount": 3000000,
      "remaining_disbursement": 3000000,
      "total_collection_from_unit": 1500000,
      "disbursement_status": "PARTIAL",
      "collection_status": "PARTIAL",
      "noc_received": true
    }
  ]
}
```

### Payment Stage Report
**Endpoint**: `GET /api/v1/project-management/reports/payment-stages?project_id={project_id}&unit_id={unit_id}`

Generates payment stage report with collection status and variance analysis.

**Query Parameters**:
- `project_id` (optional): Filter by project
- `unit_id` (optional): Filter by unit

**Response** (200 OK):
```json
{
  "count": 5,
  "data": [
    {
      "id": "stage-uuid",
      "stage_name": "BOOKING",
      "stage_number": 1,
      "stage_percentage": 20,
      "stage_due_amount": 1500000,
      "amount_due": 1500000,
      "amount_received": 1500000,
      "amount_pending": 0,
      "collection_status": "COMPLETED",
      "due_date": "2025-12-20T00:00:00Z",
      "payment_received_date": "2025-12-18T00:00:00Z",
      "payment_mode": "NEFT"
    }
  ]
}
```

---

## Error Responses

All endpoints follow standard error response format:

**400 Bad Request**:
```json
{
  "error": "Invalid request body"
}
```

**401 Unauthorized**:
```json
{
  "error": "Unauthorized - Missing or invalid token"
}
```

**404 Not Found**:
```json
{
  "error": "Resource not found"
}
```

**500 Internal Server Error**:
```json
{
  "error": "Failed to process request"
}
```

---

## Data Types & Field Explanations

### Charge Types
- `PER_SQFT`: Charge calculated as amount × SBUA
- `LUMPSUM`: Fixed charge amount regardless of unit size

### Collection Status
- `PENDING`: Not yet collected
- `PARTIAL`: Partially collected
- `COMPLETED`: Fully collected
- `OVERDUE`: Due date passed without collection

### Disbursement Status
- `PENDING`: Awaiting disbursement
- `PARTIAL`: Partially disbursed
- `COMPLETED`: Fully disbursed
- `PENDING_DOCUMENTS`: Waiting for required documents

### Payment Mode
- `CASH`: Direct cash payment
- `CHEQUE`: Cheque payment
- `NEFT`: National Electronic Funds Transfer
- `RTGS`: Real Time Gross Settlement
- `ONLINE`: Online transfer
- `DD`: Demand Draft

### Unit Type
- `1BHK`: One bedroom, hall, kitchen
- `2BHK`: Two bedroom, hall, kitchen
- `3BHK`: Three bedroom, hall, kitchen
- `STUDIO`: Studio apartment

### NOC Status
- `YES`: NOC received
- `NO`: NOC not received
- `PENDING`: NOC pending
- `NA`: NOC not applicable

---

## Example Workflows

### Complete Unit Setup Workflow

1. **Create Customer Profile**
   ```
   POST /api/v1/project-management/customers
   ```

2. **Create Area Statement**
   ```
   POST /api/v1/project-management/area-statements
   ```

3. **Update Cost Sheet**
   ```
   POST /api/v1/project-management/cost-sheets
   ```

4. **Create Bank Financing**
   ```
   POST /api/v1/project-management/bank-financing
   ```

5. **Create Payment Stages**
   ```
   POST /api/v1/project-management/payment-stages (5 times for 5 stages)
   ```

6. **Create Disbursement Schedule**
   ```
   POST /api/v1/project-management/disbursement-schedule
   ```

### Payment Collection Workflow

1. **Record Disbursement**
   ```
   PUT /api/v1/project-management/disbursement/{id}
   ```

2. **Record Stage Collection**
   ```
   PUT /api/v1/project-management/payment-stages/{id}/collection
   ```

3. **Generate Reports**
   ```
   GET /api/v1/project-management/reports/bank-financing
   GET /api/v1/project-management/reports/payment-stages
   ```

---

## Implementation Notes

- All amounts are in INR (₹)
- All dates are in `YYYY-MM-DD` format
- All IDs are UUID v4 strings
- Timestamps are in ISO 8601 format (UTC)
- Multi-tenancy is automatic based on authentication token
- All API calls require valid JWT token in Authorization header
- Request/response payload limit: 10 MB
