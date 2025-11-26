# Phase 3E - Unified Codebase Implementation Guide

## Overview
This document provides unified implementation guidelines for all Phase 3E modules following the established code patterns from Phase 1-3D.

**Date**: November 25, 2025  
**Status**: Implementation Started  
**Completed Modules**: Purchase (GRN/MRN/Contracts), Module Routes Created

---

## 1. Architecture Overview

### Frontend Structure (Next.js 16)
```
frontend/
├── app/dashboard/
│   ├── purchase/
│   │   └── page.tsx              # Module main page
│   ├── sales/
│   │   └── page.tsx
│   ├── hr/
│   │   └── page.tsx
│   ├── accounts/
│   │   └── page.tsx
│   ├── construction/
│   │   └── page.tsx
│   ├── civil/
│   │   └── page.tsx
│   └── presales/
│       └── page.tsx
└── components/modules/
    ├── Purchase/
    │   ├── PurchaseDashboard.tsx
    │   ├── VendorManagement.tsx
    │   ├── PurchaseOrderManagement.tsx
    │   ├── GRNManagement.tsx
    │   └── ContractManagement.tsx
    ├── Sales/
    │   └── [Components...]
    └── [Other Modules...]
```

### Backend Structure (Go + GORM)
```
internal/
├── models/
│   ├── purchase.go               # All purchase-related models
│   ├── sales.go
│   ├── hr.go
│   └── [Others...]
├── handlers/
│   ├── purchase_handler.go       # All purchase handlers
│   ├── sales_handler.go
│   └── [Others...]
└── services/
    ├── purchase_service.go       # Business logic
    └── [Others...]

migrations/
├── 008_purchase_module_schema.sql
├── 009_sales_module_schema.sql
└── [Others...]
```

---

## 2. Frontend Standards (Next.js + TypeScript + Tailwind)

### Module Page Pattern
```typescript
'use client'

import { useState } from 'react'
import DashboardLayout from '@/components/layouts/DashboardLayout'
import ComponentName from '@/components/modules/ModuleName/ComponentName'

type TabType = 'dashboard' | 'tab1' | 'tab2' | 'tab3'

export default function ModulePage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'tab1', label: 'Tab 1' },
    // ... more tabs
  ]

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Header with gradient */}
        <div className="bg-gradient-to-r from-[color]-600 to-[color]-800 rounded-lg p-6 text-white">
          <h1 className="text-3xl font-bold">Module Name</h1>
          <p className="text-[color]-100 mt-2">Description</p>
        </div>

        {/* Tab Navigation */}
        <div className="flex gap-2 border-b border-gray-200">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-3 font-medium border-b-2 transition ${
                activeTab === tab.id
                  ? 'border-[color]-600 text-[color]-600'
                  : 'border-transparent text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Content */}
        <div className="bg-white rounded-lg shadow">
          {/* Tab content rendering */}
        </div>
      </div>
    </DashboardLayout>
  )
}
```

### Component Pattern
```typescript
'use client'

import { useState, useEffect } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'

interface DataType {
  id: string
  [fields...]
}

export default function ComponentName() {
  const [data, setData] = useState<DataType[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const response = await axios.get('/api/v1/module/endpoint')
      setData(response.data || [])
    } catch (error) {
      toast.error('Failed to fetch data')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  // ... component logic

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-bold">Component Title</h2>
        <button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
          + Add Item
        </button>
      </div>

      {/* Form Modal */}
      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          {/* Form JSX */}
        </div>
      )}

      {/* Data Table/List */}
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          {/* Table structure */}
        </table>
      </div>
    </div>
  )
}
```

### Color Scheme by Module
```
HR & Payroll       : orange-600 (👨‍💼)
Accounts (GL)      : indigo-600 (💰)
Sales              : green-600 (🛒)
Purchase           : blue-600 (📦)
Construction       : red-600 (🏗️)
Civil              : teal-600 (🌉)
Post Sales         : pink-600 (⭐)
```

---

## 3. Backend Standards (Go + GORM)

### Handler Structure Pattern
```go
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/SkyDaddy001/VYOMTECH-ERP/internal/models"
    "github.com/gorilla/mux"
    "gorm.io/gorm"
)

type ModuleHandler struct {
    DB *gorm.DB
}

// Create - POST /api/v1/module/resource
func (h *ModuleHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req models.ResourceModel
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
        return
    }

    // Extract tenant from header
    tenantID := r.Header.Get("X-Tenant-ID")
    req.TenantID = tenantID
    req.CreatedAt = time.Now()

    // Save to DB
    if err := h.DB.Create(&req).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(req)
}

// Read - GET /api/v1/module/resource
func (h *ModuleHandler) Read(w http.ResponseWriter, r *http.Request) {
    tenantID := r.Header.Get("X-Tenant-ID")
    var resources []models.ResourceModel

    if err := h.DB.Where("tenant_id = ?", tenantID).Find(&resources).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resources)
}

// Update - PUT /api/v1/module/resource/{id}
func (h *ModuleHandler) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var req models.ResourceModel
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
        return
    }

    req.UpdatedAt = time.Now()
    if err := h.DB.Model(&req).Where("id = ?", id).Updates(&req).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(req)
}

// Delete - DELETE /api/v1/module/resource/{id}
func (h *ModuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if err := h.DB.Where("id = ?", id).Delete(&models.ResourceModel{}).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}
```

### Model Structure Pattern
```go
package models

import "time"

type ResourceModel struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    TenantID  string    `json:"tenant_id" gorm:"index"`
    // Core fields
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    // Timestamps
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at"`
    // Relations
    RelatedItems []RelatedModel `json:"related_items,omitempty" gorm:"foreignKey:ResourceID"`
}

// TableName specifies the table name
func (ResourceModel) TableName() string {
    return "resource_models"
}
```

### Route Registration Pattern
```go
// In main.go or routes setup
func setupRoutes(r *mux.Router, db *gorm.DB) {
    // Purchase Routes
    purchaseHandler := &handlers.PurchaseHandler{DB: db}
    
    // Vendor Routes
    r.HandleFunc("/api/v1/purchase/vendors", purchaseHandler.CreateVendor).Methods("POST")
    r.HandleFunc("/api/v1/purchase/vendors", purchaseHandler.ListVendors).Methods("GET")
    r.HandleFunc("/api/v1/purchase/vendors/{id}", purchaseHandler.GetVendor).Methods("GET")
    r.HandleFunc("/api/v1/purchase/vendors/{id}", purchaseHandler.UpdateVendor).Methods("PUT")
    r.HandleFunc("/api/v1/purchase/vendors/{id}", purchaseHandler.DeleteVendor).Methods("DELETE")
    
    // Purchase Order Routes
    r.HandleFunc("/api/v1/purchase/orders", purchaseHandler.CreatePurchaseOrder).Methods("POST")
    r.HandleFunc("/api/v1/purchase/orders", purchaseHandler.ListPurchaseOrders).Methods("GET")
    // ... more routes
}
```

---

## 4. Purchase Module Implementation Details

### GRN/MRN Features
- Log material receipts against Purchase Orders
- Quality inspection workflow
- Accept/Reject functionality
- Quantity reconciliation

### Contract Management
- **Material Contracts**: Supply agreements for goods
- **Labour Contracts**: Contractor service agreements
- **Service Contracts**: Professional services
- **Hybrid Contracts**: Combined material + labour, or material + service
- **BOQ Integration**: Link contracts to Bill of Quantities

### API Endpoints (Purchase)

#### Vendors
```
POST   /api/v1/purchase/vendors                 # Create vendor
GET    /api/v1/purchase/vendors                 # List vendors
GET    /api/v1/purchase/vendors/{id}           # Get vendor details
PUT    /api/v1/purchase/vendors/{id}           # Update vendor
DELETE /api/v1/purchase/vendors/{id}           # Delete vendor
```

#### Purchase Orders
```
POST   /api/v1/purchase/orders                 # Create PO
GET    /api/v1/purchase/orders                 # List POs
GET    /api/v1/purchase/orders/{id}           # Get PO details
PUT    /api/v1/purchase/orders/{id}           # Update PO
POST   /api/v1/purchase/orders/{id}/approve   # Approve PO
POST   /api/v1/purchase/orders/{id}/cancel    # Cancel PO
```

#### GRN/MRN
```
POST   /api/v1/purchase/grn                    # Create GRN
GET    /api/v1/purchase/grn                    # List GRNs
POST   /api/v1/purchase/grn/{id}/quality-check # QC Inspection
POST   /api/v1/purchase/grn/{id}/accept        # Accept GRN
POST   /api/v1/purchase/grn/{id}/reject        # Reject GRN
```

#### Contracts
```
POST   /api/v1/purchase/contracts              # Create contract
GET    /api/v1/purchase/contracts              # List contracts
GET    /api/v1/purchase/contracts/{id}        # Get contract details
PUT    /api/v1/purchase/contracts/{id}        # Update contract
POST   /api/v1/purchase/contracts/{id}/link-boq # Link to BOQ
```

---

## 5. Module Navigation Structure

All modules are accessible via:
- `/dashboard/purchase` - Purchase Module
- `/dashboard/sales` - Sales Module
- `/dashboard/hr` - HR & Payroll Module
- `/dashboard/accounts` - Accounts (GL) Module
- `/dashboard/construction` - Construction Module
- `/dashboard/civil` - Civil Module
- `/dashboard/presales` - Post Sales Module

Each module route automatically includes navigation to the main dashboard via the DashboardLayout sidebar.

---

## 6. Database Schema Standards

### Multi-Tenant Pattern
```sql
-- All tables include tenant_id for isolation
CREATE TABLE vendors (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    vendor_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_created_at (created_at)
);
```

### Audit Trail Pattern
```sql
-- All transactional tables include audit columns
ALTER TABLE vendors ADD COLUMN (
    created_by VARCHAR(36),
    updated_by VARCHAR(36),
    deleted_by VARCHAR(36)
);
```

---

## 7. API Response Standards

### Success Response
```json
{
    "id": "resource_id",
    "name": "Resource Name",
    "status": "active",
    "created_at": "2025-11-25T10:00:00Z",
    ...
}
```

### Error Response
```json
{
    "error": "Error message description",
    "code": "ERROR_CODE",
    "timestamp": "2025-11-25T10:00:00Z"
}
```

### List Response
```json
{
    "data": [...],
    "total": 100,
    "page": 1,
    "limit": 20
}
```

---

## 8. Testing Standards

### Frontend Components
- Use `jest` and `@testing-library/react`
- Test user interactions (clicks, form submissions)
- Mock API calls with `axios`
- Verify UI elements render correctly

### Backend Handlers
- Use Go's `testing` package
- Mock database with `gorm` test helpers
- Test all CRUD operations
- Verify error handling

---

## 9. Development Workflow

### 1. Create Module Routes
- Add page.tsx in `frontend/app/dashboard/[module]/`
- Include tab navigation
- Import component modules

### 2. Create Components
- Create components in `frontend/components/modules/[Module]/`
- Follow component pattern (state, effects, handlers)
- Use Tailwind CSS for styling

### 3. Create Backend Models
- Define models in `internal/models/[module].go`
- Include all relationships
- Add JSON and GORM tags

### 4. Create Backend Handlers
- Create `internal/handlers/[module]_handler.go`
- Implement CRUD handlers
- Handle tenant isolation

### 5. Create Database Migration
- Create `migrations/NNN_[module]_schema.sql`
- Define all tables with proper indexes
- Include foreign key relationships

### 6. Register Routes
- Add routes in `main.go` or routing setup
- Follow RESTful patterns
- Add middleware for auth/tenant validation

---

## 10. Completed Phase 3E Work

### ✅ Completed
1. **Purchase Module**
   - Database schema migration created
   - Backend models (Vendor, PO, GRN, Contract)
   - Purchase handlers (CRUD operations)
   - Frontend pages and components
   - GRN/MRN logging
   - Contract management with BOQ integration

2. **Module Navigation**
   - All 7 module routes created
   - Dashboard sidebar updated
   - Consistent UI/UX across modules
   - Color-coded by module

3. **Unified Codebase**
   - Next.js 16 frontend (no Ant Design)
   - Go backend with GORM
   - TypeScript for type safety
   - Tailwind CSS for styling
   - RESTful API structure

### 📋 Pending
1. **Sales Module** - Full implementation
2. **HR & Payroll Module** - Full implementation
3. **Accounts (GL) Module** - Full implementation
4. **Construction Module** - Full implementation
5. **Civil Module** - Full implementation
6. **Post Sales Module** - Full implementation

### 🔄 Integration Points
- All modules → GL (Accounts) for posting
- Purchase → Inventory (stock tracking)
- Sales → AR (Accounts Receivable)
- HR → GL (Salary Expense)
- Construction → GL & Purchase (material tracking)

---

## 11. Quick Start for New Modules

To implement a new module:

1. Create page: `frontend/app/dashboard/[module]/page.tsx`
2. Create components in `frontend/components/modules/[Module]/`
3. Create models in `internal/models/[module].go`
4. Create handlers in `internal/handlers/[module]_handler.go`
5. Create migration in `migrations/NNN_[module]_schema.sql`
6. Register routes in route setup
7. Update DashboardLayout with module link

---

## 12. Deployment Checklist

- [ ] All database migrations run successfully
- [ ] Backend handlers deployed
- [ ] Frontend components compiled
- [ ] API routes registered
- [ ] Environment variables configured
- [ ] Database backups created
- [ ] Multi-tenant isolation verified
- [ ] SSL/TLS configured
- [ ] Rate limiting configured
- [ ] Monitoring set up

---

**Document Version**: 1.0  
**Last Updated**: November 25, 2025  
**Maintained By**: Development Team
