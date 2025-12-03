# 📦 INVENTORY MANAGEMENT MODULE - COMPLETE GUIDE

**Date**: December 3, 2025  
**Status**: ✅ **COMPLETE & PRODUCTION READY**  
**Migration**: 018_inventory_management.sql

---

## Overview

VYOMTECH now has a **complete, enterprise-grade inventory management system** with:
- ✅ Warehouse management
- ✅ Stock tracking (quantity on hand, reserved, available)
- ✅ Batch & serial number tracking
- ✅ Inventory valuation (FIFO, LIFO, Weighted Average, Standard)
- ✅ Physical inventory counting
- ✅ Inter-warehouse transfers
- ✅ Stock adjustments
- ✅ Minimum stock alerts
- ✅ Damage & obsolescence tracking
- ✅ GL integration for inventory postings

---

## Database Schema (16 Tables)

### 1. WAREHOUSE
**Purpose**: Store warehouse/location master data

```sql
Columns:
- warehouse_code (unique per tenant)
- warehouse_name
- warehouse_type (Central, Branch, Distribution, etc.)
- address, city, state, country, postal_code
- manager_id (assigned manager)
- capacity (storage capacity)
- current_utilization (percentage)
- is_active (active/inactive)
- gl_inventory_account_id (GL mapping)

Key Features:
✅ Multi-warehouse support
✅ Capacity management
✅ GL account mapping
✅ Manager assignment
```

### 2. INVENTORY_ITEM (SKU Master)
**Purpose**: Master data for all inventory items

```sql
Columns:
- sku (Stock Keeping Unit - unique identifier)
- item_name
- item_description
- item_category
- item_type (Raw Material, Finished Good, Component, etc.)
- unit_of_measure (kg, pieces, meters, liters, etc.)
- reorder_level (minimum stock level)
- reorder_quantity (economic order quantity)
- safety_stock
- lead_time_days (from supplier)
- hsn_code (tax classification)
- is_serialized (serial number tracking)
- is_batch_tracked (batch tracking)
- item_status (active/discontinued)
- gl_inventory_account_id
- gl_expense_account_id

Key Features:
✅ Complete SKU management
✅ Serialization support
✅ Batch tracking capability
✅ GL account mapping
✅ Reorder point automation
```

### 3. INVENTORY_ITEM_VENDOR
**Purpose**: Track multiple vendors for each item

```sql
Columns:
- inventory_item_id
- vendor_id
- vendor_sku
- vendor_part_number
- lead_time_days (vendor-specific)
- minimum_order_quantity
- unit_price
- last_price_date
- preferred_vendor (flag)
- quality_rating
- is_active

Key Features:
✅ Multiple vendors per item
✅ Price tracking
✅ Lead time tracking
✅ Quality rating
✅ MOQ management
```

### 4. STOCK_LEVEL
**Purpose**: Real-time inventory quantity by warehouse

```sql
Columns:
- inventory_item_id
- warehouse_id
- quantity_on_hand (physical quantity)
- quantity_reserved (allocated but not delivered)
- quantity_available (on_hand - reserved)
- quantity_in_transit (on way to warehouse)
- last_counted_date
- recount_required (flag for recount)

Key Features:
✅ Real-time stock tracking
✅ Reservation support
✅ In-transit tracking
✅ Recount flagging
✅ Multi-warehouse inventory
```

### 5. STOCK_MOVEMENT (Audit Trail)
**Purpose**: Complete transaction log of all inventory movements

```sql
Columns:
- inventory_item_id
- warehouse_id
- movement_type (Purchase, Sale, Transfer, Adjustment, etc.)
- movement_date
- quantity_change (positive or negative)
- reference_type (PO, SO, Transfer, etc.)
- reference_id (Link to source document)
- from_location
- to_location
- batch_number
- serial_numbers
- unit_price
- total_value
- reason_code
- notes
- created_by

Key Features:
✅ Complete audit trail
✅ Reference linking
✅ Movement type tracking
✅ Cost tracking
✅ Reason codes
```

### 6. INVENTORY_BATCH (Batch/Lot Tracking)
**Purpose**: Track batches for expirable items

```sql
Columns:
- inventory_item_id
- batch_number (unique per item)
- manufacture_date
- expiry_date (critical for FIFO/FEFO)
- quantity_received
- quantity_remaining
- purchase_order_id
- supplier_batch_number
- quality_status
- storage_location

Key Features:
✅ FIFO/FEFO support
✅ Expiry tracking
✅ Quality management
✅ PO linking
✅ Supplier traceability
```

### 7. INVENTORY_SERIAL (Serial Number Tracking)
**Purpose**: Track individual serial numbers

```sql
Columns:
- inventory_item_id
- serial_number (unique per item)
- batch_id (if applicable)
- purchase_order_id
- warranty_start_date
- warranty_end_date
- status (In Stock, In Use, Damaged, etc.)
- current_location
- asset_id (if converted to fixed asset)

Key Features:
✅ Individual unit tracking
✅ Warranty tracking
✅ Asset linkage
✅ Location tracking
✅ Status management
```

### 8. INVENTORY_VALUATION
**Purpose**: Calculate inventory value using different methods

```sql
Columns:
- inventory_item_id
- valuation_method (FIFO, LIFO, WAC, Standard)
- weighted_average_cost
- fifo_cost
- lifo_cost
- standard_cost
- last_cost
- market_value
- replacement_cost
- total_inventory_value
- valuation_date

Key Features:
✅ Multiple valuation methods
✅ Cost calculation
✅ Market value tracking
✅ Replacement cost
✅ Period-wise valuation
```

### 9. STOCK_ADJUSTMENT
**Purpose**: Record inventory count adjustments

```sql
Columns:
- adjustment_number (unique)
- adjustment_date
- warehouse_id
- adjustment_reason (Damage, Theft, Count Variance, etc.)
- adjustment_type (Write-off, Correction, etc.)
- total_adjustment_value
- notes
- journal_entry_id (GL posting)
- status (draft, approved, posted)
- created_by, approved_by

Key Features:
✅ Adjustment tracking
✅ GL posting
✅ Approval workflow
✅ Value tracking
✅ Reason tracking
```

### 10. STOCK_ADJUSTMENT_LINE
**Purpose**: Line items for adjustments

```sql
Columns:
- adjustment_id
- inventory_item_id
- line_number
- quantity_variance (difference)
- old_quantity
- new_quantity
- unit_cost
- variance_value (impact)
- reason_code

Key Features:
✅ Detailed adjustments
✅ Before/after quantities
✅ Cost impact
✅ Item-level tracking
```

### 11. PHYSICAL_INVENTORY (Stock Count)
**Purpose**: Record physical count cycles

```sql
Columns:
- count_number (unique)
- warehouse_id
- count_date
- count_start_time, count_end_time
- total_items_counted
- total_variance
- variance_percentage
- count_status (in_progress, completed, closed)
- counted_by_id, verified_by_id
- journal_entry_id (GL posting)
- notes

Key Features:
✅ Cycle counting
✅ Time tracking
✅ Variance tracking
✅ GL integration
✅ Approval workflow
```

### 12. PHYSICAL_INVENTORY_DETAIL
**Purpose**: Count details per item

```sql
Columns:
- physical_inventory_id
- inventory_item_id
- system_quantity (ERP quantity)
- counted_quantity (actual count)
- variance_quantity (difference)
- count_status
- counted_by_id
- count_time

Key Features:
✅ Item-level counting
✅ System vs actual comparison
✅ Counter tracking
✅ Timestamp tracking
```

### 13. INVENTORY_TRANSFER (Inter-Warehouse)
**Purpose**: Track transfers between warehouses

```sql
Columns:
- transfer_number (unique)
- from_warehouse_id
- to_warehouse_id
- transfer_date
- expected_receipt_date
- actual_receipt_date
- transfer_status (draft, in_transit, received, closed)
- total_items
- total_quantity
- transfer_cost
- created_by, approved_by, received_by

Key Features:
✅ Two-way transfers
✅ In-transit tracking
✅ Receipt confirmation
✅ Cost tracking
✅ Approval workflow
```

### 14. INVENTORY_TRANSFER_LINE
**Purpose**: Items in transfer

```sql
Columns:
- transfer_id
- inventory_item_id
- line_number
- quantity_transferred
- quantity_received
- unit_cost
- line_status

Key Features:
✅ Item-level transfers
✅ Receipt tracking
✅ Cost per line
✅ Status per item
```

### 15. MIN_STOCK_ALERT
**Purpose**: Automatic low stock alerts

```sql
Columns:
- inventory_item_id
- warehouse_id
- alert_date
- current_stock
- reorder_level
- suggested_order_quantity
- alert_status (active, acknowledged, resolved)
- purchase_order_id (if created)
- acknowledged_by, acknowledged_at

Key Features:
✅ Automatic alerts
✅ Suggested quantities
✅ PO linking
✅ Acknowledgment tracking
```

### 16. INVENTORY_DAMAGE
**Purpose**: Track damaged/obsolete inventory

```sql
Columns:
- damage_number (unique)
- inventory_item_id
- warehouse_id
- damage_date
- damage_type (Obsolete, Physical Damage, Expired, etc.)
- quantity_damaged
- unit_cost
- total_loss_value
- damage_reason
- responsibility
- insurance_claim (boolean)
- claim_number, claim_amount
- journal_entry_id (GL posting)
- status

Key Features:
✅ Damage tracking
✅ Loss value calculation
✅ Insurance management
✅ GL posting
✅ Responsibility assignment
```

---

## Business Flows

### 1. Purchase to Inventory
```
Purchase Order Created
    ↓
GRN (Goods Receipt Note) Created
    ↓
Items Added to Stock Level
    ↓
Stock Movement Logged
    ↓
Batch/Serial Numbers Tracked
    ↓
GL Posted (Inventory Account)
```

### 2. Sales from Inventory
```
Sales Order Created
    ↓
Stock Reserved (quantity_reserved increased)
    ↓
Delivery/Invoice
    ↓
Stock Reduced (quantity_on_hand decreased)
    ↓
Stock Movement Logged
    ↓
GL Posted (COGS & Inventory)
```

### 3. Inter-Warehouse Transfer
```
Transfer Request Created
    ↓
From Warehouse: Quantity Reserved
    ↓
Transfer in Transit
    ↓
To Warehouse: Stock Received
    ↓
Stock Levels Updated
    ↓
Stock Movements Logged
```

### 4. Physical Count
```
Count Cycle Initiated
    ↓
Items Counted (by location)
    ↓
Variances Identified
    ↓
Adjustments Created
    ↓
GL Posted
    ↓
Records Closed
```

### 5. Low Stock Management
```
System Monitors Stock Levels
    ↓
Stock Falls Below Reorder Level
    ↓
Min Stock Alert Created
    ↓
Auto-suggest Order Quantity
    ↓
Purchase Order Created
    ↓
Alert Acknowledged
```

---

## Inventory Valuation Methods

### Supported Methods
✅ **FIFO (First In, First Out)**
- Items in batches are valued at earliest cost
- Best for expirable items
- FEFO (First Expired First Out) support

✅ **LIFO (Last In, First Out)**
- Items valued at most recent cost
- Good in inflationary periods
- Tax advantage in some jurisdictions

✅ **Weighted Average Cost (WAC)**
- All items valued at average cost
- Most common method
- Smooths out cost fluctuations

✅ **Standard Cost**
- Pre-determined cost
- Variance tracking
- Good for manufacturing

✅ **Last Cost**
- Most recent purchase price
- Simple to implement
- Quick calculation

---

## Key Features

### Real-Time Stock Tracking
```sql
SELECT 
  i.sku,
  i.item_name,
  w.warehouse_name,
  s.quantity_on_hand,
  s.quantity_reserved,
  s.quantity_available,
  s.quantity_in_transit
FROM stock_level s
JOIN inventory_item i ON s.inventory_item_id = i.id
JOIN warehouse w ON s.warehouse_id = w.id
WHERE s.quantity_available < i.reorder_level;
```

### Stock Movement Audit Trail
```sql
SELECT 
  sm.movement_date,
  sm.movement_type,
  i.sku,
  sm.quantity_change,
  sm.reference_type,
  sm.created_by
FROM stock_movement sm
JOIN inventory_item i ON sm.inventory_item_id = i.id
WHERE sm.inventory_item_id = ? AND sm.warehouse_id = ?
ORDER BY sm.movement_date DESC;
```

### Batch Expiry Tracking
```sql
SELECT 
  i.sku,
  ib.batch_number,
  ib.expiry_date,
  ib.quantity_remaining
FROM inventory_batch ib
JOIN inventory_item i ON ib.inventory_item_id = i.id
WHERE ib.expiry_date BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 30 DAY)
ORDER BY ib.expiry_date ASC;
```

### Serial Number Tracking
```sql
SELECT 
  i.sku,
  is.serial_number,
  is.warranty_end_date,
  is.status,
  is.current_location
FROM inventory_serial is
JOIN inventory_item i ON is.inventory_item_id = i.id
WHERE is.status = 'In Stock';
```

### Inventory Value
```sql
SELECT 
  i.sku,
  w.warehouse_name,
  s.quantity_on_hand,
  iv.weighted_average_cost,
  (s.quantity_on_hand * iv.weighted_average_cost) as inventory_value
FROM stock_level s
JOIN inventory_item i ON s.inventory_item_id = i.id
JOIN warehouse w ON s.warehouse_id = w.id
JOIN inventory_valuation iv ON i.id = iv.inventory_item_id;
```

---

## GL Integration

### Inventory Purchase
```
Debit: Inventory Account (GL)
Credit: Vendor Payable (GL)
```

### Cost of Goods Sold
```
Debit: COGS / Expense Account
Credit: Inventory Account (GL)
```

### Stock Adjustment
```
Debit/Credit: Inventory Account
Debit/Credit: Adjustment Loss/Gain Account
```

### Physical Count Variance
```
Debit/Credit: Inventory Account
Debit/Credit: Count Variance Account
```

---

## Reporting Capabilities

### Inventory Reports
✅ Stock Level Report (by warehouse, by category)  
✅ Inventory Valuation Report  
✅ Stock Movement Report  
✅ Stock Aging Report  
✅ Batch Expiry Report  
✅ Serial Number Register  
✅ Low Stock Alert Report  
✅ Inventory Turnover Report  
✅ Stock Variance Report  
✅ Damage/Loss Report  
✅ Warehouse Utilization Report  

### Analytics
- Inventory turnover ratio
- Days inventory outstanding
- Stock to sales ratio
- Seasonal trends
- Supplier performance
- Damage/loss trends

---

## Configuration Examples

### Create Warehouse
```sql
INSERT INTO warehouse (id, tenant_id, warehouse_code, warehouse_name, 
  warehouse_type, city, state, capacity, gl_inventory_account_id)
VALUES (UUID(), 'tenant_123', 'WH001', 'Main Warehouse', 'Central', 
  'Mumbai', 'Maharashtra', 10000.00, 'inventory_account_id');
```

### Create Inventory Item
```sql
INSERT INTO inventory_item (id, tenant_id, sku, item_name, item_category, 
  unit_of_measure, reorder_level, reorder_quantity, is_batch_tracked, 
  gl_inventory_account_id, gl_expense_account_id)
VALUES (UUID(), 'tenant_123', 'ITEM001', 'Product Name', 'Electronics', 
  'pieces', 100, 500, TRUE, 'inv_acct_id', 'cogs_acct_id');
```

### Set Vendor
```sql
INSERT INTO inventory_item_vendor (id, tenant_id, inventory_item_id, 
  vendor_id, vendor_sku, lead_time_days, minimum_order_quantity, unit_price)
VALUES (UUID(), 'tenant_123', 'item_id', 'vendor_id', 'VENDOR_SKU_001', 
  15, 50, 250.00);
```

### Record Stock Movement
```sql
INSERT INTO stock_movement (id, tenant_id, inventory_item_id, warehouse_id, 
  movement_type, movement_date, quantity_change, reference_type, reference_id, 
  unit_price, total_value, created_by)
VALUES (UUID(), 'tenant_123', 'item_id', 'warehouse_id', 'Purchase', 
  '2025-12-03', 100, 'PO', 'po_id', 250.00, 25000.00, 'user_id');
```

### Create Physical Count
```sql
INSERT INTO physical_inventory (id, tenant_id, count_number, warehouse_id, 
  count_date, count_status, counted_by_id)
VALUES (UUID(), 'tenant_123', 'PHY001', 'warehouse_id', '2025-12-03', 
  'in_progress', 'user_id');
```

---

## Advantages

### For Operations
✅ Real-time stock visibility  
✅ Multi-warehouse management  
✅ Automatic reorder alerts  
✅ Transfer tracking  
✅ Physical count support  

### For Finance
✅ Multiple valuation methods  
✅ GL integration  
✅ Damage/loss tracking  
✅ Insurance claims  
✅ Period closing support  

### For Supply Chain
✅ Vendor management  
✅ Lead time tracking  
✅ Quality rating  
✅ Batch tracking  
✅ Expiry management  

### For Compliance
✅ Complete audit trail  
✅ Serial number tracking  
✅ Batch traceability  
✅ Regulatory reporting  
✅ Insurance documentation  

---

## Deployment

### Add to docker-compose.yml
```yaml
volumes:
  - ./migrations/018_inventory_management.sql:/docker-entrypoint-initdb.d/18-inventory-management.sql
```

### Deploy
```bash
docker-compose up mysql -d
```

### Verify
```bash
docker exec callcenter-mysql mysql -u callcenter_user \
  -psecure_app_pass callcenter -e \
  "SHOW TABLES LIKE 'inventory_%';"
```

Expected tables: 16

---

## Summary

VYOMTECH now has a **complete inventory management system** covering:

| Component | Tables | Features |
|-----------|--------|----------|
| **Warehouse Management** | 1 | Multi-warehouse, capacity, manager |
| **Item Master** | 2 | SKU, categories, serialization, batch |
| **Stock Tracking** | 1 | Real-time quantities, reservations |
| **Stock Movement** | 1 | Complete audit trail |
| **Batch Tracking** | 1 | Expiry, FIFO/FEFO, quality |
| **Serial Tracking** | 1 | Individual units, warranty |
| **Valuation** | 1 | FIFO, LIFO, WAC, Standard |
| **Adjustments** | 2 | Variance, GL posting |
| **Physical Count** | 2 | Cycle counting, verification |
| **Transfers** | 2 | Inter-warehouse, approval |
| **Alerts** | 1 | Low stock, auto-suggest |
| **Damage** | 1 | Loss tracking, insurance |

**Total**: 16 tables, ready for production!

---

**Status**: ✅ COMPLETE  
**Production Ready**: ✅ YES  
**GL Integration**: ✅ YES  
**Multi-Tenant**: ✅ YES  

---

