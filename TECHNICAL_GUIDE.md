# VYOMTECH ERP - Technical Implementation Guide

**Purpose**: Developer reference for implementation, architecture, and best practices  
**Audience**: Developers, Architects, DevOps  
**Last Updated**: November 25, 2025  

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Project Structure](#project-structure)
3. [Frontend Development](#frontend-development)
4. [Backend Development](#backend-development)
5. [Database Design](#database-design)
6. [API Standards](#api-standards)
7. [Testing](#testing)
8. [Deployment](#deployment)
9. [Troubleshooting](#troubleshooting)

---

## Quick Start

### Clone & Setup

```bash
# Clone repository
git clone https://github.com/SkyDaddy001/VYOMTECH-ERP.git
cd VYOMTECH-ERP

# Frontend setup
cd frontend
npm install
npm run dev

# Backend setup (in new terminal)
cd backend
go mod download
go run cmd/main.go

# Database setup
mysql -u root -p < migrations/init.sql
```

### Access Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Docs**: http://localhost:8080/swagger

---

## Project Structure

### Frontend

```
frontend/
├── app/                          # Next.js App Router
│   ├── auth/                     # Authentication pages
│   │   ├── login/
│   │   ├── register/
│   │   └── reset-password/
│   ├── dashboard/                # Protected routes
│   │   ├── page.tsx             # Main dashboard
│   │   ├── purchase/             # Purchase module routes
│   │   │   └── page.tsx
│   │   ├── sales/
│   │   ├── hr/
│   │   ├── accounts/
│   │   ├── construction/
│   │   ├── civil/
│   │   └── presales/
│   ├── layout.tsx               # Root layout
│   ├── page.tsx                 # Home page
│   └── globals.css              # Global styles
├── components/                   # React components
│   ├── layouts/
│   │   └── DashboardLayout.tsx  # Main layout component
│   ├── dashboard/               # Dashboard components
│   ├── modules/                 # Module-specific components
│   │   ├── Purchase/
│   │   │   ├── PurchaseDashboard.tsx
│   │   │   ├── VendorManagement.tsx
│   │   │   ├── PurchaseOrderManagement.tsx
│   │   │   ├── GRNManagement.tsx
│   │   │   └── ContractManagement.tsx
│   │   ├── Sales/
│   │   ├── HR/
│   │   └── [Other modules...]
│   ├── auth/                    # Auth components
│   ├── providers/               # Context providers
│   └── common/                  # Reusable components
├── hooks/                        # Custom React hooks
│   ├── useAuth.ts
│   ├── useAPI.ts
│   └── [Other hooks...]
├── services/                     # API services
│   ├── authService.ts
│   ├── purchaseService.ts
│   ├── salesService.ts
│   └── [Other services...]
├── types/                        # TypeScript types
│   ├── index.ts
│   ├── auth.ts
│   ├── purchase.ts
│   └── [Other types...]
├── contexts/                     # React contexts
│   ├── AuthContext.tsx
│   └── TenantContext.tsx
├── utils/                        # Utilities
│   ├── api.ts
│   ├── constants.ts
│   └── helpers.ts
├── package.json
├── tsconfig.json
├── next.config.js
├── tailwind.config.js
└── postcss.config.js
```

### Backend

```
backend/ (cmd/main.go)
internal/
├── config/                       # Configuration
│   └── config.go
├── db/                          # Database setup
│   ├── connection.go
│   └── migrations.go
├── handlers/                     # HTTP handlers
│   ├── auth_handler.go
│   ├── purchase_handler.go      # ✅ 14+ endpoints
│   ├── sales_handler.go
│   ├── hr_handler.go
│   ├── accounts_handler.go
│   └── [Other handlers...]
├── middleware/                   # HTTP middleware
│   ├── auth.go
│   ├── tenant.go
│   └── logger.go
├── models/                       # Data models
│   ├── user.go
│   ├── purchase.go              # ✅ 15+ tables
│   ├── sales.go
│   ├── hr.go
│   ├── accounts.go
│   └── [Other models...]
├── services/                     # Business logic
│   ├── auth_service.go
│   ├── purchase_service.go
│   ├── sales_service.go
│   └── [Other services...]
└── pkg/
    ├── logger/
    ├── utils/
    └── validators/

migrations/                       # Database schemas
├── 001_initial_schema.sql
├── 002_multi_tenant_users.sql
├── 008_purchase_module_schema.sql ✅
├── 009_sales_module_schema.sql
└── [Other migrations...]

cmd/
└── main.go                       # Application entry point

go.mod                           # Go module file
go.sum                           # Dependency lock file
```

---

## Frontend Development

### Component Pattern

```typescript
'use client'  // Client component for interactivity

import { useState, useEffect } from 'react'
import axios from 'axios'
import toast from 'react-hot-toast'

interface DataType {
  id: string
  name: string
  status: string
  created_at: string
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
      const response = await axios.get('/api/v1/module/resource')
      setData(response.data || [])
    } catch (error) {
      toast.error('Failed to fetch data')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (formData: Partial<DataType>) => {
    try {
      await axios.post('/api/v1/module/resource', formData)
      toast.success('Created successfully')
      fetchData()
    } catch (error) {
      toast.error('Failed to create')
    }
  }

  if (loading) return <div>Loading...</div>

  return (
    <div className="p-6 space-y-6">
      <h2 className="text-2xl font-bold">Component Title</h2>
      {/* Content */}
    </div>
  )
}
```

### Styling Pattern (Tailwind)

```typescript
// Color scheme by module
export const colors = {
  purchase: 'blue',      // 📦
  sales: 'green',        // 🛒
  hr: 'orange',          // 👨‍💼
  accounts: 'indigo',    // 💰
  construction: 'red',   // 🏗️
  civil: 'teal',         // 🌉
  postsales: 'pink',     // ⭐
}

// Button styles
className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition"

// Input styles
className="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"

// Card styles
className="bg-white rounded-lg shadow p-6"

// Table styles
className="w-full border-collapse"
```

### API Service Pattern

```typescript
// services/purchaseService.ts
import axios, { AxiosInstance } from 'axios'

class PurchaseService {
  private api: AxiosInstance

  constructor() {
    this.api = axios.create({
      baseURL: '/api/v1/purchase',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }

  // Vendors
  async createVendor(data: VendorData) {
    return this.api.post('/vendors', data)
  }

  async getVendors(page = 1, limit = 20) {
    return this.api.get('/vendors', { params: { page, limit } })
  }

  // Purchase Orders
  async createPO(data: POData) {
    return this.api.post('/orders', data)
  }

  async getPOs() {
    return this.api.get('/orders')
  }

  // GRN
  async createGRN(data: GRNData) {
    return this.api.post('/grn', data)
  }

  async qualityCheckGRN(id: string, data: QCData) {
    return this.api.post(`/grn/${id}/quality-check`, data)
  }
}

export default new PurchaseService()
```

---

## Backend Development

### Handler Pattern

```go
package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "gorm.io/gorm"
    "github.com/SkyDaddy001/VYOMTECH-ERP/internal/models"
)

type PurchaseHandler struct {
    DB *gorm.DB
}

// CreateVendor - POST /api/v1/purchase/vendors
func (h *PurchaseHandler) CreateVendor(w http.ResponseWriter, r *http.Request) {
    // Extract tenant
    tenantID := r.Header.Get("X-Tenant-ID")
    if tenantID == "" {
        http.Error(w, `{"error": "Missing tenant"}`, http.StatusBadRequest)
        return
    }

    // Parse request
    var vendor models.Vendor
    if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
        http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
        return
    }

    // Set metadata
    vendor.TenantID = tenantID
    vendor.Status = "active"
    vendor.CreatedAt = time.Now()

    // Save to DB
    if err := h.DB.Create(&vendor).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    // Return response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(vendor)
}

// ListVendors - GET /api/v1/purchase/vendors
func (h *PurchaseHandler) ListVendors(w http.ResponseWriter, r *http.Request) {
    tenantID := r.Header.Get("X-Tenant-ID")

    var vendors []models.Vendor
    if err := h.DB.Where("tenant_id = ?", tenantID).Find(&vendors).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(vendors)
}

// UpdateVendor - PUT /api/v1/purchase/vendors/{id}
func (h *PurchaseHandler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var vendor models.Vendor
    if err := json.NewDecoder(r.Body).Decode(&vendor); err != nil {
        http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
        return
    }

    vendor.UpdatedAt = time.Now()
    if err := h.DB.Model(&vendor).Where("id = ?", id).Updates(&vendor).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(vendor)
}

// DeleteVendor - DELETE /api/v1/purchase/vendors/{id}
func (h *PurchaseHandler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if err := h.DB.Where("id = ?", id).Delete(&models.Vendor{}).Error; err != nil {
        http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Deleted successfully"})
}
```

### Model Pattern

```go
package models

import "time"

// Vendor represents a supplier/vendor
type Vendor struct {
    ID           string    `json:"id" gorm:"primaryKey"`
    TenantID     string    `json:"tenant_id" gorm:"index"`
    VendorCode   string    `json:"vendor_code" gorm:"unique"`
    Name         string    `json:"name" gorm:"index"`
    Email        string    `json:"email"`
    Phone        string    `json:"phone"`
    Status       string    `json:"status" gorm:"index"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    DeletedAt    *time.Time `json:"deleted_at"`
    CreatedBy    string    `json:"created_by"`
    
    // Relations
    Contacts []VendorContact `json:"contacts,omitempty" gorm:"foreignKey:VendorID"`
    Orders   []PurchaseOrder  `json:"orders,omitempty" gorm:"foreignKey:VendorID"`
}

func (Vendor) TableName() string {
    return "vendors"
}
```

### Route Registration

```go
// In cmd/main.go

func setupRoutes(router *mux.Router, db *gorm.DB) {
    // Purchase routes
    purchaseHandler := &handlers.PurchaseHandler{DB: db}
    
    // Vendor endpoints
    router.HandleFunc("/api/v1/purchase/vendors", purchaseHandler.CreateVendor).Methods("POST")
    router.HandleFunc("/api/v1/purchase/vendors", purchaseHandler.ListVendors).Methods("GET")
    router.HandleFunc("/api/v1/purchase/vendors/{id}", purchaseHandler.GetVendor).Methods("GET")
    router.HandleFunc("/api/v1/purchase/vendors/{id}", purchaseHandler.UpdateVendor).Methods("PUT")
    router.HandleFunc("/api/v1/purchase/vendors/{id}", purchaseHandler.DeleteVendor).Methods("DELETE")
    
    // Purchase Order endpoints
    router.HandleFunc("/api/v1/purchase/orders", purchaseHandler.CreatePurchaseOrder).Methods("POST")
    router.HandleFunc("/api/v1/purchase/orders", purchaseHandler.ListPurchaseOrders).Methods("GET")
    // ... more routes
}
```

---

## Database Design

### Multi-Tenant Schema Pattern

```sql
-- All tables follow this pattern
CREATE TABLE resource_name (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    -- Core fields
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    -- Audit trail
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(36),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_by VARCHAR(36),
    deleted_at TIMESTAMP NULL,
    -- Indexes
    UNIQUE KEY unique_code (tenant_id, code),
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
```

### Example: Purchase Module Tables

**Vendors**
```sql
CREATE TABLE vendors (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    vendor_code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    tax_id VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_status (tenant_id, status)
);
```

**Purchase Orders**
```sql
CREATE TABLE purchase_orders (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    po_number VARCHAR(50) NOT NULL UNIQUE,
    vendor_id VARCHAR(36) NOT NULL,
    po_date DATE,
    total_amount DECIMAL(15, 2),
    status VARCHAR(20) DEFAULT 'draft',
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_vendor (vendor_id),
    INDEX idx_po_number (po_number),
    INDEX idx_status (status),
    FOREIGN KEY (vendor_id) REFERENCES vendors(id)
);
```

**GRN (Goods Receipt Notes)**
```sql
CREATE TABLE goods_receipts (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    grn_number VARCHAR(50) NOT NULL UNIQUE,
    po_id VARCHAR(36) NOT NULL,
    received_date DATE,
    total_received DECIMAL(10, 2),
    inspection_status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_po_id (po_id),
    INDEX idx_inspection_status (inspection_status),
    FOREIGN KEY (po_id) REFERENCES purchase_orders(id)
);
```

---

## API Standards

### Request Format

```json
// Headers
{
    "Content-Type": "application/json",
    "Authorization": "Bearer {jwt_token}",
    "X-Tenant-ID": "tenant_123"
}

// Body (example: Create Vendor)
{
    "vendor_code": "VENDOR001",
    "name": "ABC Suppliers",
    "email": "contact@abc.com",
    "phone": "+1234567890",
    "payment_terms": "30 days"
}
```

### Response Format

```json
// Success (200, 201)
{
    "id": "vendor_123",
    "vendor_code": "VENDOR001",
    "name": "ABC Suppliers",
    "status": "active",
    "created_at": "2025-11-25T10:00:00Z"
}

// Error (400, 500)
{
    "error": "Vendor code already exists",
    "code": "VENDOR_DUPLICATE",
    "timestamp": "2025-11-25T10:00:00Z"
}
```

### Pagination

```json
// Query
GET /api/v1/purchase/vendors?page=1&limit=20

// Response
{
    "data": [...],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 150,
        "pages": 8
    }
}
```

---

## Testing

### Frontend Testing

```typescript
// __tests__/components/VendorManagement.test.tsx
import { render, screen, fireEvent } from '@testing-library/react'
import VendorManagement from '@/components/modules/Purchase/VendorManagement'

jest.mock('axios')

describe('VendorManagement', () => {
  test('renders vendor list', async () => {
    render(<VendorManagement />)
    expect(screen.getByText('Vendor Management')).toBeInTheDocument()
  })

  test('adds vendor on form submit', async () => {
    render(<VendorManagement />)
    // Test implementation
  })
})
```

### Backend Testing

```go
// internal/handlers/purchase_handler_test.go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateVendor(t *testing.T) {
    // Test implementation
}

func TestListVendors(t *testing.T) {
    // Test implementation
}
```

---

## Deployment

### Docker Build

```bash
# Frontend
cd frontend
docker build -t vyomtech-frontend:latest .

# Backend
cd backend
docker build -t vyomtech-backend:latest .
```

### Docker Compose

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: vyomtech
    ports:
      - "3306:3306"

  backend:
    build: ./backend
    depends_on:
      - mysql
    environment:
      DB_HOST: mysql
      DB_USER: root
      DB_PASSWORD: password
      DB_NAME: vyomtech
    ports:
      - "8080:8080"

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
```

---

## Troubleshooting

### Common Issues

**Issue**: Database connection refused  
**Solution**: Check MySQL is running, credentials are correct, database exists

**Issue**: API returning 401 Unauthorized  
**Solution**: Verify JWT token is valid, not expired, and X-Tenant-ID header is set

**Issue**: Frontend not connecting to API  
**Solution**: Check API URL in frontend env vars, CORS is enabled, backend is running

**Issue**: Performance degradation  
**Solution**: Check database indexes, run EXPLAIN on slow queries, consider caching

---

**Document**: VYOMTECH ERP Technical Implementation Guide  
**Version**: 1.0  
**Last Updated**: November 25, 2025
