/**
 * CRUD, RBAC, and Isolation Tests
 * Comprehensive testing for Create, Read, Update, Delete operations
 * Role-Based Access Control and Data Isolation verification
 */

import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

/**
 * CRUD OPERATIONS TESTS
 */
describe('CRUD Operations', () => {
  /**
   * CREATE Operations
   */
  describe('Create Operations', () => {
    it('should create invoice with all required fields', async () => {
      const invoiceData = {
        invoiceNumber: 'INV-001',
        date: '2024-01-15',
        customerName: 'ACME Corp',
        customerEmail: 'contact@acme.com',
        taxId: '27ABCDE1234F2Z5',
        items: [
          { description: 'Service', quantity: 1, unitPrice: 10000, taxRate: 18 },
        ],
        paymentTerms: 'NET30',
        notes: 'Test invoice',
      };

      const response = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify(invoiceData),
      });

      expect(response.status).toBe(201);
      const created = await response.json();
      expect(created.id).toBeDefined();
      expect(created.invoiceNumber).toBe('INV-001');
      expect(created.customerId).toBeDefined();
      expect(created.tenantId).toBe('tenant-123');
    });

    it('should create sales order with line items', async () => {
      const orderData = {
        orderNumber: 'SO-001',
        date: '2024-01-15',
        dueDate: '2024-02-15',
        customerName: 'TechStart Inc',
        customerPhone: '+91-9876543210',
        deliveryAddress: 'Mumbai, India',
        items: [
          {
            productName: 'Widget A',
            sku: 'WID-001',
            quantity: 10,
            unitPrice: 5000,
            discount: 0,
          },
        ],
        orderStatus: 'DRAFT',
      };

      const response = await fetch('/api/v1/sales-orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify(orderData),
      });

      expect(response.status).toBe(201);
      const created = await response.json();
      expect(created.id).toBeDefined();
      expect(created.orderNumber).toBe('SO-001');
      expect(created.items.length).toBe(1);
    });

    it('should create BOQ with precision calculations', async () => {
      const boqData = {
        projectCode: 'PRJ-001',
        projectName: 'Commercial Complex',
        contractorName: 'BuildCorp Ltd',
        contractorContact: '+91-9876543210',
        boqDate: '2024-01-15',
        projectLocation: 'Mumbai',
        items: [
          {
            description: 'Excavation',
            specification: '1m depth',
            quantity: 500.5,
            unit: 'cum',
            ratePerUnit: 2500.5,
            progress: 100,
          },
        ],
        contingency: 5,
      };

      const response = await fetch('/api/v1/boq', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify(boqData),
      });

      expect(response.status).toBe(201);
      const created = await response.json();
      expect(created.id).toBeDefined();
      // Verify 0.01 rupee precision: 500.5 × 2500.5 = 1,251,250.25
      expect(created.items[0].amount).toBe(1251250.25);
    });

    it('should reject create without required fields', async () => {
      const invalidData = {
        date: '2024-01-15',
        // Missing invoiceNumber, customerName
      };

      const response = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify(invalidData),
      });

      expect(response.status).toBe(400);
      const error = await response.json();
      expect(error.error).toContain('required');
    });

    it('should auto-generate ID as UUID string', async () => {
      const invoiceData = {
        invoiceNumber: 'INV-UUID-TEST',
        date: '2024-01-15',
        customerName: 'Test Corp',
        customerEmail: 'test@corp.com',
        taxId: '',
        items: [],
        paymentTerms: 'NET30',
        notes: '',
      };

      const response = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify(invoiceData),
      });

      const created = await response.json();
      // Verify UUID format (36 chars with hyphens)
      expect(created.id).toMatch(
        /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i
      );
    });

    it('should create with timestamps (created_at, updated_at)', async () => {
      const invoiceData = {
        invoiceNumber: 'INV-TIMESTAMP',
        date: '2024-01-15',
        customerName: 'Timestamp Test',
        customerEmail: 'test@timestamp.com',
        taxId: '',
        items: [],
        paymentTerms: 'NET30',
        notes: '',
      };

      const response = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify(invoiceData),
      });

      const created = await response.json();
      expect(created.createdAt).toBeDefined();
      expect(created.updatedAt).toBeDefined();
      expect(new Date(created.createdAt).getTime()).toBeGreaterThan(0);
    });
  });

  /**
   * READ Operations
   */
  describe('Read Operations', () => {
    it('should read invoice by ID', async () => {
      // Create first
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-READ-TEST',
          date: '2024-01-15',
          customerName: 'Read Test',
          customerEmail: 'read@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;

      // Read
      const response = await fetch(`/api/v1/invoices/${invoiceId}`, {
        headers: {
          'X-Tenant-ID': 'tenant-123',
        },
      });

      expect(response.status).toBe(200);
      const invoice = await response.json();
      expect(invoice.id).toBe(invoiceId);
      expect(invoice.invoiceNumber).toBe('INV-READ-TEST');
    });

    it('should list invoices with pagination', async () => {
      const response = await fetch(
        '/api/v1/invoices?limit=10&offset=0&sort=date&order=desc',
        {
          headers: {
            'X-Tenant-ID': 'tenant-123',
          },
        }
      );

      expect(response.status).toBe(200);
      const result = await response.json();
      expect(result.data).toBeDefined();
      expect(Array.isArray(result.data)).toBe(true);
      expect(result.total).toBeDefined();
      expect(result.limit).toBe(10);
      expect(result.offset).toBe(0);
    });

    it('should filter invoices by status', async () => {
      const response = await fetch(
        '/api/v1/invoices?status=PAID&limit=50',
        {
          headers: {
            'X-Tenant-ID': 'tenant-123',
          },
        }
      );

      expect(response.status).toBe(200);
      const result = await response.json();
      // All returned invoices should have status=PAID
      result.data.forEach((invoice: any) => {
        expect(invoice.status).toBe('PAID');
      });
    });

    it('should search invoices by customer name', async () => {
      const response = await fetch(
        '/api/v1/invoices?search=ACME&limit=50',
        {
          headers: {
            'X-Tenant-ID': 'tenant-123',
          },
        }
      );

      expect(response.status).toBe(200);
      const result = await response.json();
      // Results should contain customer name matching search
      result.data.forEach((invoice: any) => {
        expect(invoice.customerName.toLowerCase()).toContain('acme');
      });
    });

    it('should return 404 for non-existent invoice', async () => {
      const response = await fetch(
        '/api/v1/invoices/00000000-0000-0000-0000-000000000000',
        {
          headers: {
            'X-Tenant-ID': 'tenant-123',
          },
        }
      );

      expect(response.status).toBe(404);
    });

    it('should read BOQ items with progress tracking', async () => {
      // Create BOQ
      const createRes = await fetch('/api/v1/boq', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          projectCode: 'PRJ-BOQ-READ',
          projectName: 'BOQ Read Test',
          contractorName: 'Test Contractor',
          contractorContact: '+91-9876543210',
          boqDate: '2024-01-15',
          projectLocation: 'Mumbai',
          items: [
            {
              description: 'Item 1',
              specification: 'Spec 1',
              quantity: 100,
              unit: 'nos',
              ratePerUnit: 1000,
              progress: 50,
            },
          ],
          contingency: 5,
        }),
      });

      const created = await createRes.json();
      const boqId = created.id;

      // Read
      const response = await fetch(`/api/v1/boq/${boqId}`, {
        headers: {
          'X-Tenant-ID': 'tenant-123',
        },
      });

      expect(response.status).toBe(200);
      const boq = await response.json();
      expect(boq.items[0].progress).toBe(50);
      expect(boq.items[0].amount).toBe(100000); // 100 × 1000
    });
  });

  /**
   * UPDATE Operations
   */
  describe('Update Operations', () => {
    it('should update invoice customer name', async () => {
      // Create
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-UPDATE-TEST',
          date: '2024-01-15',
          customerName: 'Old Name',
          customerEmail: 'old@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;

      // Update
      const updateRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          customerName: 'New Name',
          customerEmail: 'new@test.com',
        }),
      });

      expect(updateRes.status).toBe(200);
      const updated = await updateRes.json();
      expect(updated.customerName).toBe('New Name');
      expect(updated.customerEmail).toBe('new@test.com');
      expect(updated.id).toBe(invoiceId); // ID should not change
    });

    it('should update invoice status through state machine', async () => {
      // Create with DRAFT status
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-STATUS-FLOW',
          date: '2024-01-15',
          customerName: 'Status Flow Test',
          customerEmail: 'status@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;
      expect(created.status).toBe('DRAFT');

      // Transition DRAFT → SENT
      let updateRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ status: 'SENT' }),
      });

      expect(updateRes.status).toBe(200);
      let updated = await updateRes.json();
      expect(updated.status).toBe('SENT');

      // Transition SENT → PAID
      updateRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ status: 'PAID' }),
      });

      expect(updateRes.status).toBe(200);
      updated = await updateRes.json();
      expect(updated.status).toBe('PAID');

      // Reject invalid transition PAID → DRAFT (reverse not allowed)
      updateRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ status: 'DRAFT' }),
      });

      expect(updateRes.status).toBe(400);
      const error = await updateRes.json();
      expect(error.error).toContain('invalid transition');
    });

    it('should update BOQ item progress', async () => {
      // Create BOQ
      const createRes = await fetch('/api/v1/boq', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          projectCode: 'PRJ-BOQ-UPDATE',
          projectName: 'BOQ Update Test',
          contractorName: 'Test Contractor',
          contractorContact: '+91-9876543210',
          boqDate: '2024-01-15',
          projectLocation: 'Mumbai',
          items: [
            {
              description: 'Item 1',
              specification: 'Spec 1',
              quantity: 100,
              unit: 'nos',
              ratePerUnit: 1000,
              progress: 0,
            },
          ],
          contingency: 5,
        }),
      });

      const created = await createRes.json();
      const boqId = created.id;

      // Update progress to 50%
      const updateRes = await fetch(`/api/v1/boq/${boqId}/items/0`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ progress: 50 }),
      });

      expect(updateRes.status).toBe(200);
      const updated = await updateRes.json();
      expect(updated.items[0].progress).toBe(50);

      // Reject progress > 100%
      const invalidRes = await fetch(`/api/v1/boq/${boqId}/items/0`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ progress: 150 }),
      });

      expect(invalidRes.status).toBe(400);
    });

    it('should update updatedAt timestamp on modification', async () => {
      // Create
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-TIMESTAMP-UPDATE',
          date: '2024-01-15',
          customerName: 'Timestamp Test',
          customerEmail: 'timestamp@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;
      const originalUpdatedAt = new Date(created.updatedAt).getTime();

      // Wait 100ms
      await new Promise((resolve) => setTimeout(resolve, 100));

      // Update
      const updateRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ customerName: 'Updated Name' }),
      });

      const updated = await updateRes.json();
      const newUpdatedAt = new Date(updated.updatedAt).getTime();
      expect(newUpdatedAt).toBeGreaterThan(originalUpdatedAt);
    });
  });

  /**
   * DELETE Operations
   */
  describe('Delete Operations', () => {
    it('should soft-delete invoice (mark as deleted)', async () => {
      // Create
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-SOFT-DELETE',
          date: '2024-01-15',
          customerName: 'Soft Delete Test',
          customerEmail: 'delete@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;

      // Delete (soft)
      const deleteRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': 'tenant-123',
        },
      });

      expect(deleteRes.status).toBe(200);

      // Verify soft-delete (deleted_at set, but record still exists)
      const getRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        headers: {
          'X-Tenant-ID': 'tenant-123',
        },
      });

      // Should not appear in normal queries
      expect(getRes.status).toBe(404);
    });

    it('should prevent deletion of PAID invoice', async () => {
      // Create and mark as PAID
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-PAID-DELETE',
          date: '2024-01-15',
          customerName: 'Paid Delete Test',
          customerEmail: 'paid@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;

      // Mark as PAID
      await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({ status: 'PAID' }),
      });

      // Try to delete
      const deleteRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': 'tenant-123',
        },
      });

      expect(deleteRes.status).toBe(403);
      const error = await deleteRes.json();
      expect(error.error).toContain('cannot delete');
    });

    it('should allow deletion of DRAFT invoice', async () => {
      // Create in DRAFT status
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-DRAFT-DELETE',
          date: '2024-01-15',
          customerName: 'Draft Delete Test',
          customerEmail: 'draft@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;

      // Delete DRAFT invoice
      const deleteRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': 'tenant-123',
        },
      });

      expect(deleteRes.status).toBe(200);
    });
  });
});

/**
 * RBAC (Role-Based Access Control) TESTS
 */
describe('RBAC (Role-Based Access Control)', () => {
  /**
   * Read Permissions
   */
  describe('Read Permissions', () => {
    it('Admin role can read all invoices', async () => {
      const response = await fetch('/api/v1/invoices?limit=100', {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
      });

      expect(response.status).toBe(200);
      const result = await response.json();
      expect(result.data).toBeDefined();
    });

    it('Sales user can read sales data', async () => {
      const response = await fetch('/api/v1/sales-orders?limit=50', {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': 'user-sales-1',
        },
      });

      expect(response.status).toBe(200);
      const result = await response.json();
      expect(result.data).toBeDefined();
    });

    it('Accountant can read GL transactions', async () => {
      const response = await fetch('/api/v1/journal-entries?limit=50', {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'accountant',
          'X-User-ID': 'user-accountant-1',
        },
      });

      expect(response.status).toBe(200);
      const result = await response.json();
      expect(result.data).toBeDefined();
    });

    it('Sales user cannot read accounting data', async () => {
      const response = await fetch('/api/v1/chart-of-accounts?limit=50', {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': 'user-sales-1',
        },
      });

      expect(response.status).toBe(403);
      const error = await response.json();
      expect(error.error).toContain('unauthorized');
    });

    it('Guest role has read-only access', async () => {
      // Can read
      const readRes = await fetch('/api/v1/invoices?limit=10', {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'guest',
          'X-User-ID': 'user-guest-1',
        },
      });

      expect(readRes.status).toBe(200);

      // Cannot write
      const writeRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'guest',
          'X-User-ID': 'user-guest-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-GUEST',
          date: '2024-01-15',
          customerName: 'Guest Test',
          customerEmail: 'guest@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      expect(writeRes.status).toBe(403);
    });
  });

  /**
   * Create Permissions
   */
  describe('Create Permissions', () => {
    it('Admin can create invoices', async () => {
      const response = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-ADMIN-CREATE',
          date: '2024-01-15',
          customerName: 'Admin Create',
          customerEmail: 'admin@create.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      expect(response.status).toBe(201);
    });

    it('Sales user can create sales orders', async () => {
      const response = await fetch('/api/v1/sales-orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': 'user-sales-1',
        },
        body: JSON.stringify({
          orderNumber: 'SO-SALES-CREATE',
          date: '2024-01-15',
          dueDate: '2024-02-15',
          customerName: 'Sales Create',
          customerPhone: '+91-9876543210',
          deliveryAddress: 'Mumbai',
          items: [],
          orderStatus: 'DRAFT',
        }),
      });

      expect(response.status).toBe(201);
    });

    it('Accountant cannot create invoices', async () => {
      const response = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'accountant',
          'X-User-ID': 'user-accountant-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-ACCOUNTANT-CREATE',
          date: '2024-01-15',
          customerName: 'Accountant Create',
          customerEmail: 'acc@create.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      expect(response.status).toBe(403);
    });

    it('Accountant can create GL entries', async () => {
      const response = await fetch('/api/v1/journal-entries', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'accountant',
          'X-User-ID': 'user-accountant-1',
        },
        body: JSON.stringify({
          date: '2024-01-15',
          description: 'Journal entry',
          items: [
            { account: '1101', debit: 10000 },
            { account: '4001', credit: 10000 },
          ],
        }),
      });

      expect(response.status).toBe(201);
    });
  });

  /**
   * Update Permissions
   */
  describe('Update Permissions', () => {
    it('Admin can update any invoice', async () => {
      // Create first
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-ADMIN-UPDATE',
          date: '2024-01-15',
          customerName: 'Admin Update',
          customerEmail: 'admin@update.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // Update by another admin
      const updateRes = await fetch(`/api/v1/invoices/${created.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-2',
        },
        body: JSON.stringify({ customerName: 'Updated by Admin 2' }),
      });

      expect(updateRes.status).toBe(200);
    });

    it('Sales user can only update own sales orders', async () => {
      const userId = 'user-sales-1';

      // Create by sales user 1
      const createRes = await fetch('/api/v1/sales-orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': userId,
        },
        body: JSON.stringify({
          orderNumber: 'SO-SALES-OWN',
          date: '2024-01-15',
          dueDate: '2024-02-15',
          customerName: 'Own Order',
          customerPhone: '+91-9876543210',
          deliveryAddress: 'Mumbai',
          items: [],
          orderStatus: 'DRAFT',
        }),
      });

      const created = await createRes.json();

      // Update by same user (allowed)
      let updateRes = await fetch(`/api/v1/sales-orders/${created.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': userId,
        },
        body: JSON.stringify({ orderStatus: 'CONFIRMED' }),
      });

      expect(updateRes.status).toBe(200);

      // Update by different user (denied)
      updateRes = await fetch(`/api/v1/sales-orders/${created.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': 'user-sales-2',
        },
        body: JSON.stringify({ orderStatus: 'READY' }),
      });

      expect(updateRes.status).toBe(403);
    });

    it('Cannot update createdAt or id fields', async () => {
      // Create
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-READONLY-FIELDS',
          date: '2024-01-15',
          customerName: 'Readonly Fields',
          customerEmail: 'readonly@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const originalId = created.id;

      // Try to update id
      const updateRes = await fetch(`/api/v1/invoices/${originalId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
        body: JSON.stringify({
          id: 'different-id',
          customerName: 'Updated Name',
        }),
      });

      const updated = await updateRes.json();
      expect(updated.id).toBe(originalId); // ID unchanged
    });
  });

  /**
   * Delete Permissions
   */
  describe('Delete Permissions', () => {
    it('Admin can delete any invoice', async () => {
      // Create
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-ADMIN-DELETE',
          date: '2024-01-15',
          customerName: 'Admin Delete',
          customerEmail: 'admin@delete.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // Delete by admin
      const deleteRes = await fetch(`/api/v1/invoices/${created.id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
      });

      expect(deleteRes.status).toBe(200);
    });

    it('Sales user cannot delete invoices', async () => {
      // Create
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'admin',
          'X-User-ID': 'user-admin-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-SALES-DELETE',
          date: '2024-01-15',
          customerName: 'Sales Delete',
          customerEmail: 'sales@delete.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // Try delete by sales user
      const deleteRes = await fetch(`/api/v1/invoices/${created.id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-Role': 'sales',
          'X-User-ID': 'user-sales-1',
        },
      });

      expect(deleteRes.status).toBe(403);
    });
  });
});

/**
 * DATA ISOLATION TESTS
 */
describe('Data Isolation', () => {
  /**
   * Tenant Isolation
   */
  describe('Tenant Isolation', () => {
    it('Tenant-1 cannot read Tenant-2 invoices', async () => {
      // Create invoice in Tenant-2
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-2',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-TENANT-2',
          date: '2024-01-15',
          customerName: 'Tenant 2 Customer',
          customerEmail: 'tenant2@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();
      const invoiceId = created.id;

      // Try to read from Tenant-1
      const readRes = await fetch(`/api/v1/invoices/${invoiceId}`, {
        headers: {
          'X-Tenant-ID': 'tenant-1',
        },
      });

      expect(readRes.status).toBe(404);
    });

    it('Tenant-1 cannot update Tenant-2 invoices', async () => {
      // Create in Tenant-2
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-2',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-TENANT-2-UPDATE',
          date: '2024-01-15',
          customerName: 'Tenant 2',
          customerEmail: 'tenant2@update.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // Try update from Tenant-1
      const updateRes = await fetch(`/api/v1/invoices/${created.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-1',
        },
        body: JSON.stringify({ customerName: 'Hacked' }),
      });

      expect(updateRes.status).toBe(404);
    });

    it('Tenant-1 cannot delete Tenant-2 invoices', async () => {
      // Create in Tenant-2
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-2',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-TENANT-2-DELETE',
          date: '2024-01-15',
          customerName: 'Tenant 2 Delete',
          customerEmail: 'tenant2@delete.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // Try delete from Tenant-1
      const deleteRes = await fetch(`/api/v1/invoices/${created.id}`, {
        method: 'DELETE',
        headers: {
          'X-Tenant-ID': 'tenant-1',
        },
      });

      expect(deleteRes.status).toBe(404);
    });

    it('List queries only return current tenant data', async () => {
      // Create in Tenant-1
      await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-LIST-T1',
          date: '2024-01-15',
          customerName: 'Tenant 1',
          customerEmail: 'tenant1@list.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      // Create in Tenant-2
      await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-2',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-LIST-T2',
          date: '2024-01-15',
          customerName: 'Tenant 2',
          customerEmail: 'tenant2@list.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      // List from Tenant-1
      const t1Res = await fetch('/api/v1/invoices?limit=100', {
        headers: {
          'X-Tenant-ID': 'tenant-1',
        },
      });

      const t1Data = await t1Res.json();

      // List from Tenant-2
      const t2Res = await fetch('/api/v1/invoices?limit=100', {
        headers: {
          'X-Tenant-ID': 'tenant-2',
        },
      });

      const t2Data = await t2Res.json();

      // Verify isolation
      const t1Invoices = t1Data.data;
      const t2Invoices = t2Data.data;

      t1Invoices.forEach((inv: any) => {
        expect(inv.invoiceNumber).toContain('T1');
      });

      t2Invoices.forEach((inv: any) => {
        expect(inv.invoiceNumber).toContain('T2');
      });
    });
  });

  /**
   * User Isolation
   */
  describe('User Isolation', () => {
    it('User cannot see draft invoices of other users', async () => {
      // User-1 creates draft
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-ID': 'user-1',
          'X-User-Role': 'sales',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-USER1-DRAFT',
          date: '2024-01-15',
          customerName: 'User 1 Draft',
          customerEmail: 'user1@draft.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // User-2 tries to view
      const readRes = await fetch(`/api/v1/invoices/${created.id}`, {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-ID': 'user-2',
          'X-User-Role': 'sales',
        },
      });

      // Should not be visible (privacy for draft)
      expect(readRes.status).toBe(404);
    });

    it('User can see published invoices of other users', async () => {
      // User-1 creates and publishes
      const createRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-ID': 'user-1',
          'X-User-Role': 'sales',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-USER1-SENT',
          date: '2024-01-15',
          customerName: 'User 1 Sent',
          customerEmail: 'user1@sent.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const created = await createRes.json();

      // Publish (SENT status)
      await fetch(`/api/v1/invoices/${created.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-123',
          'X-User-ID': 'user-1',
          'X-User-Role': 'sales',
        },
        body: JSON.stringify({ status: 'SENT' }),
      });

      // User-2 can view
      const readRes = await fetch(`/api/v1/invoices/${created.id}`, {
        headers: {
          'X-Tenant-ID': 'tenant-123',
          'X-User-ID': 'user-2',
          'X-User-Role': 'sales',
        },
      });

      expect(readRes.status).toBe(200);
    });
  });

  /**
   * Related Data Isolation
   */
  describe('Related Data Isolation', () => {
    it('Cannot associate invoice with customer from different tenant', async () => {
      // Create customer in Tenant-2
      const customerRes = await fetch('/api/v1/customers', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-2',
        },
        body: JSON.stringify({
          name: 'Tenant 2 Customer',
          email: 'tenant2@customer.com',
        }),
      });

      const customer = await customerRes.json();

      // Try to create invoice in Tenant-1 with Tenant-2 customer
      const invoiceRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-CROSS-TENANT',
          date: '2024-01-15',
          customerId: customer.id,
          customerName: 'Hacked',
          customerEmail: 'hacked@test.com',
          taxId: '',
          items: [],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      expect(invoiceRes.status).toBe(404); // Customer not found in Tenant-1
    });

    it('Invoice and its GL entries must be in same tenant', async () => {
      // Create invoice in Tenant-1
      const invoiceRes = await fetch('/api/v1/invoices', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': 'tenant-1',
        },
        body: JSON.stringify({
          invoiceNumber: 'INV-GL-ISOLATION',
          date: '2024-01-15',
          customerName: 'GL Test',
          customerEmail: 'gl@test.com',
          taxId: '',
          items: [{ description: 'Item', quantity: 1, unitPrice: 1000, taxRate: 18 }],
          paymentTerms: 'NET30',
          notes: '',
        }),
      });

      const invoice = await invoiceRes.json();

      // GL entry should be in same tenant
      const glRes = await fetch('/api/v1/journal-entries?invoiceId=' + invoice.id, {
        headers: {
          'X-Tenant-ID': 'tenant-1',
        },
      });

      const glData = await glRes.json();
      expect(glData.data.length).toBeGreaterThan(0);

      // Tenant-2 should not see GL entry
      const wrongTenantRes = await fetch(
        '/api/v1/journal-entries?invoiceId=' + invoice.id,
        {
          headers: {
            'X-Tenant-ID': 'tenant-2',
          },
        }
      );

      const wrongData = await wrongTenantRes.json();
      expect(wrongData.data.length).toBe(0);
    });
  });
});
