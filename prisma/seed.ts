// Prisma seed script for comprehensive data setup
// Run with: npx prisma db seed

import prisma, { generateId } from './client';
// @ts-ignore - bcryptjs types will be available after npm install @types/bcryptjs
import { hash } from 'bcryptjs';

async function hashPassword(password: string): Promise<string> {
  return hash(password, 10);
}

async function main() {
  console.log('🚀 Starting comprehensive database seed...\n');

  try {
    // ============================================================
    // SYSTEM MASTER ADMIN (Can manage all tenants)
    // ============================================================
    console.log('📌 Creating Master Admin User...');
    
    const masterAdminUser = await prisma.user.create({
      data: {
        id: generateId(),
        email: 'master.admin@vyomtech.com',
        passwordHash: await hashPassword('Master@123'),
        role: 'master_admin',
        tenantId: 'system', // System tenant
        currentTenantId: null,
      },
    });
    console.log('  ✓ Master Admin created: master.admin@vyomtech.com / Master@123\n');

    // ============================================================
    // CREATE ABC GROUP TENANT
    // ============================================================
    console.log('🏢 Creating ABC GROUP Tenant...');
    
    const abcGroupTenant = await prisma.tenant.create({
      data: {
        id: generateId(),
        name: 'ABC GROUP',
        domain: 'abc-group.vyomtech.local',
        status: 'active',
        maxUsers: 500,
        maxConcurrentCalls: 100,
        aiBudgetMonthly: 10000,
      },
    });
    console.log(`  ✓ Tenant created: ${abcGroupTenant.name}\n`);

    // ============================================================
    // CREATE COMPANIES UNDER ABC GROUP
    // ============================================================
    console.log('🏛️  Creating Companies under ABC GROUP...');
    
    const abcLlp = await prisma.tenant.create({
      data: {
        id: generateId(),
        name: 'ABC COMPANY LLP',
        domain: 'abc-llp.vyomtech.local',
        status: 'active',
        maxUsers: 250,
        maxConcurrentCalls: 50,
        aiBudgetMonthly: 5000,
      },
    });
    console.log(`  ✓ Company 1: ${abcLlp.name}`);

    const abcPvtLtd = await prisma.tenant.create({
      data: {
        id: generateId(),
        name: 'ABC COMPANY PVT LTD',
        domain: 'abc-pvtltd.vyomtech.local',
        status: 'active',
        maxUsers: 250,
        maxConcurrentCalls: 50,
        aiBudgetMonthly: 5000,
      },
    });
    console.log(`  ✓ Company 2: ${abcPvtLtd.name}\n`);

    // ============================================================
    // CREATE RBAC ROLES & PERMISSIONS
    // ============================================================
    console.log('🔐 Setting up RBAC Roles & Permissions...');
    
    const roles: { [key: string]: any } = {};
    
    const roleNames = ['Tenant Admin', 'Finance Manager', 'HR Manager', 'Sales Manager', 'Inventory Manager', 'User'];
    for (const roleName of roleNames) {
      const role = await prisma.role.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          roleName: roleName,
          roleCode: roleName.toUpperCase().replace(/\s+/g, '_'),
          description: `${roleName} role with appropriate permissions`,
          isActive: true,
        },
      });
      roles[roleName] = role;
      console.log(`  ✓ Role created: ${roleName}`);
    }

    // Create permissions
    const permissions: { [key: string]: any } = {};
    const permissionsList = [
      'create_user', 'read_user', 'update_user', 'delete_user',
      'create_company', 'read_company', 'update_company',
      'create_invoice', 'read_invoice', 'update_invoice', 'approve_invoice',
      'create_payroll', 'read_payroll', 'approve_payroll',
      'create_inventory', 'read_inventory', 'update_inventory',
      'view_dashboard', 'export_data', 'view_reports',
    ];

    for (const permCode of permissionsList) {
      const perm = await prisma.permission.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          permissionCode: permCode,
          permissionName: permCode.replace(/_/g, ' ').toUpperCase(),
          description: `Permission to ${permCode}`,
          module: permCode.split('_')[1] || 'system',
          isActive: true,
        },
      });
      permissions[permCode] = perm;
    }
    console.log(`  ✓ Created ${permissionsList.length} permissions\n`);

    // Assign permissions to Tenant Admin role
    const tenantAdminRole = roles['Tenant Admin'];
    const adminPerms = ['create_user', 'read_user', 'update_user', 'delete_user', 'view_dashboard', 'export_data'];
    for (const permCode of adminPerms) {
      await prisma.rolePermission.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          roleId: tenantAdminRole.id,
          permissionId: permissions[permCode].id,
        },
      });
    }
    console.log(`  ✓ Assigned ${adminPerms.length} permissions to Tenant Admin\n`);

    // ============================================================
    // CREATE USERS FOR ABC LLP
    // ============================================================
    console.log('👥 Creating Users for ABC LLP...');
    
    const users: { [key: string]: any } = {};
    
    const userData = [
      { email: 'admin.llp@abc-group.com', name: 'Rajesh Kumar', role: 'Tenant Admin', password: 'Admin@123' },
      { email: 'finance.llp@abc-group.com', name: 'Priya Singh', role: 'Finance Manager', password: 'Finance@123' },
      { email: 'hr.llp@abc-group.com', name: 'Amit Patel', role: 'HR Manager', password: 'HR@123' },
      { email: 'sales.llp@abc-group.com', name: 'Neha Sharma', role: 'Sales Manager', password: 'Sales@123' },
      { email: 'inventory.llp@abc-group.com', name: 'Rohan Desai', role: 'Inventory Manager', password: 'Inv@123' },
    ];

    for (const userData_item of userData) {
      const user = await prisma.user.create({
        data: {
          id: generateId(),
          email: userData_item.email,
          passwordHash: await hashPassword(userData_item.password),
          role: 'user',
          tenantId: abcLlp.id,
          currentTenantId: abcLlp.id,
        },
      });
      users[userData_item.email] = user;
      
      // Assign role
      const userRole = roles[userData_item.role];
      await prisma.userRole.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          userId: user.id,
          roleId: userRole.id,
        },
      });
      
      console.log(`  ✓ ${userData_item.name} (${userData_item.email})`);
    }
    console.log();

    // ============================================================
    // CREATE HR RECORDS FOR USERS
    // ============================================================
    console.log('📋 Creating HR Records for Users...');
    
    const employees: { [key: string]: any } = {};
    
    const hrData = [
      { userEmail: 'admin.llp@abc-group.com', empName: 'Rajesh Kumar', desig: 'Director', dept: 'Management', salary: 150000 },
      { userEmail: 'finance.llp@abc-group.com', empName: 'Priya Singh', desig: 'Finance Manager', dept: 'Finance', salary: 80000 },
      { userEmail: 'hr.llp@abc-group.com', empName: 'Amit Patel', desig: 'HR Manager', dept: 'Human Resources', salary: 75000 },
      { userEmail: 'sales.llp@abc-group.com', empName: 'Neha Sharma', desig: 'Sales Manager', dept: 'Sales', salary: 85000 },
      { userEmail: 'inventory.llp@abc-group.com', empName: 'Rohan Desai', desig: 'Inventory Manager', dept: 'Operations', salary: 70000 },
    ];

    for (const hrItem of hrData) {
      const user = users[hrItem.userEmail];
      
      const employee = await prisma.employee.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          userId: user.id,
          empName: hrItem.empName,
          empCode: `EMP${Math.random().toString(36).substr(2, 4).toUpperCase()}`,
          designation: hrItem.desig,
          department: hrItem.dept,
          dob: new Date('1990-01-01'),
          joinDate: new Date('2020-01-15'),
          ctcAmount: hrItem.salary,
          basicSalary: hrItem.salary * 0.6,
          pfPercentage: 12,
          esiPercentage: 0.75,
          status: 'active',
        },
      });
      employees[hrItem.userEmail] = employee;
      console.log(`  ✓ ${hrItem.empName} - ${hrItem.desig}`);
    }
    console.log();

    // Create payroll records
    console.log('💰 Creating Payroll Records...');
    
    for (const [email, employee] of Object.entries(employees)) {
      const payroll = await prisma.payroll.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          empId: employee.id,
          month: new Date().getMonth() + 1,
          year: new Date().getFullYear(),
          basicSalary: employee.basicSalary,
          dearness: employee.basicSalary * 0.1,
          pfAmount: employee.basicSalary * (employee.pfPercentage! / 100),
          esiAmount: employee.basicSalary * (employee.esiPercentage! / 100),
          netSalary: employee.basicSalary + (employee.basicSalary * 0.1) - (employee.basicSalary * (employee.pfPercentage! / 100)) - (employee.basicSalary * (employee.esiPercentage! / 100)),
          status: 'pending',
        },
      });
    }
    console.log('  ✓ Payroll records created\n');

    // ============================================================
    // CREATE CHART OF ACCOUNTS (Finance)
    // ============================================================
    console.log('📊 Creating Chart of Accounts...');
    
    const coas: { [key: string]: any } = {};
    
    const coaData = [
      { code: '1001', name: 'Cash in Hand', type: 'asset' },
      { code: '1002', name: 'Bank Accounts', type: 'asset' },
      { code: '2001', name: 'Accounts Payable', type: 'liability' },
      { code: '3001', name: 'Capital Account', type: 'equity' },
      { code: '4001', name: 'Sales Revenue', type: 'income' },
      { code: '5001', name: 'Cost of Goods Sold', type: 'expense' },
      { code: '5002', name: 'Salaries & Wages', type: 'expense' },
      { code: '5003', name: 'Rent Expense', type: 'expense' },
      { code: '5004', name: 'Utilities', type: 'expense' },
    ];

    for (const coaItem of coaData) {
      const coa = await prisma.chartOfAccount.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          accountCode: coaItem.code,
          accountName: coaItem.name,
          accountType: coaItem.type,
          isActive: true,
          openingBalance: 100000,
        },
      });
      coas[coaItem.code] = coa;
      console.log(`  ✓ ${coaItem.code} - ${coaItem.name}`);
    }
    console.log();

    // ============================================================
    // CREATE JOURNAL ENTRIES (Finance)
    // ============================================================
    console.log('💳 Creating Journal Entries...');
    
    const je = await prisma.journalEntry.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        jvNo: `JV-${Date.now()}`,
        jvDate: new Date(),
        narration: 'Opening Balance',
        status: 'approved',
        postedDate: new Date(),
        journalEntryDetails: {
          create: [
            {
              id: generateId(),
              tenantId: abcLlp.id,
              coaId: coas['1002'].id,
              debit: 500000,
              credit: 0,
            },
            {
              id: generateId(),
              tenantId: abcLlp.id,
              coaId: coas['3001'].id,
              debit: 0,
              credit: 500000,
            },
          ],
        },
      },
    });
    console.log('  ✓ Journal entry created with 2 line items\n');

    // ============================================================
    // CREATE VENDORS & PURCHASE ORDERS (Procurement)
    // ============================================================
    console.log('🏭 Creating Vendors & Purchase Orders...');
    
    const vendor = await prisma.vendor.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        vendorName: 'Tech Supplies India',
        vendorCode: 'VENDOR-001',
        contactEmail: 'info@techsupplies.com',
        contactPhone: '+91-9876543210',
        address: '123 Industrial Area, Mumbai',
        city: 'Mumbai',
        state: 'Maharashtra',
        country: 'India',
        pinCode: '400001',
        gstIn: '27AABCT1234A1Z1',
        status: 'active',
      },
    });
    console.log(`  ✓ Vendor created: ${vendor.vendorName}`);

    const po = await prisma.purchaseOrder.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        vendorId: vendor.id,
        poNumber: `PO-${Date.now().toString().slice(-5)}`,
        poDate: new Date(),
        deliveryDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000), // 30 days later
        totalAmount: 150000,
        taxAmount: 27000,
        grandTotal: 177000,
        status: 'approved',
      },
    });
    console.log(`  ✓ Purchase Order created: ${po.poNumber}\n`);

    // ============================================================
    // CREATE WAREHOUSES & INVENTORY (Inventory)
    // ============================================================
    console.log('📦 Creating Warehouses & Inventory...');
    
    const warehouse = await prisma.warehouse.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        warehouseName: 'Main Warehouse',
        warehouseCode: 'WH-001',
        location: 'Mumbai',
        capacity: 10000,
        currentUtilization: 4500,
        isActive: true,
      },
    });
    console.log(`  ✓ Warehouse created: ${warehouse.warehouseName}`);

    const inventoryItems = [];
    const itemData = [
      { name: 'Computer Monitors', sku: 'MON-001', price: 15000, qty: 50 },
      { name: 'Keyboards', sku: 'KEY-001', price: 2000, qty: 200 },
      { name: 'Mouse', sku: 'MOU-001', price: 800, qty: 300 },
    ];

    for (const item of itemData) {
      const invItem = await prisma.inventoryItem.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          itemName: item.name,
          itemSku: item.sku,
          unitPrice: item.price,
          reorderLevel: 10,
          isActive: true,
        },
      });
      inventoryItems.push(invItem);

      // Create stock level
      await prisma.stockLevel.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          warehouseId: warehouse.id,
          itemId: invItem.id,
          quantityOnHand: item.qty,
          quantityReserved: 0,
          quantityAvailable: item.qty,
          lastStockCount: new Date(),
        },
      });
      console.log(`  ✓ Inventory Item: ${item.name} (${item.qty} units)`);
    }
    console.log();

    // ============================================================
    // CREATE CUSTOMERS & SALES ORDERS (CRM/Sales)
    // ============================================================
    console.log('👨‍💼 Creating Customers & Sales Orders...');
    
    const customer = await prisma.salesCustomer.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        customerName: 'Global Corp Ltd',
        customerCode: 'CUST-001',
        contactEmail: 'contact@globalcorp.com',
        contactPhone: '+91-8765432109',
        billingAddress: '456 Commercial Street, Bangalore',
        shippingAddress: '456 Commercial Street, Bangalore',
        city: 'Bangalore',
        state: 'Karnataka',
        country: 'India',
        pinCode: '560001',
        gstIn: '29AABCT5678B1Z2',
        creditLimit: 500000,
        creditPeriod: 30,
        status: 'active',
      },
    });
    console.log(`  ✓ Customer created: ${customer.customerName}`);

    const so = await prisma.salesOrder.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        customerId: customer.id,
        soNumber: `SO-${Date.now().toString().slice(-5)}`,
        soDate: new Date(),
        deliveryDate: new Date(Date.now() + 15 * 24 * 60 * 60 * 1000),
        totalAmount: 300000,
        taxAmount: 54000,
        grandTotal: 354000,
        status: 'confirmed',
      },
    });
    console.log(`  ✓ Sales Order created: ${so.soNumber}\n`);

    // ============================================================
    // CREATE REAL ESTATE PROJECTS (Real Estate)
    // ============================================================
    console.log('🏗️  Creating Real Estate Projects...');
    
    const project = await prisma.constructionProject.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        projectName: 'Grand Heights Complex',
        projectCode: 'PROJ-001',
        location: 'Pune',
        totalPlotArea: 50000,
        totalBuildupArea: 100000,
        totalUnits: 6,
        estimatedValue: 50000000,
        status: 'active',
        startDate: new Date('2023-01-01'),
        endDate: new Date('2025-12-31'),
      },
    });
    console.log(`  ✓ Project created: ${project.projectName}`);

    // Create blocks within project
    for (let blockNum = 1; blockNum <= 2; blockNum++) {
      const block = await prisma.constructionBlock.create({
        data: {
          id: generateId(),
          tenantId: abcLlp.id,
          projectId: project.id,
          blockName: `Block ${String.fromCharCode(64 + blockNum)}`,
          blockCode: `BLK-${blockNum.toString().padStart(3, '0')}`,
          totalUnits: 3,
          totalArea: 50000,
          constructionStatus: 'ongoing',
          estimatedCompletionDate: new Date('2025-06-30'),
        },
      });

      // Create units within block
      for (let unitNum = 1; unitNum <= 3; unitNum++) {
        const unit = await prisma.constructionUnit.create({
          data: {
            id: generateId(),
            tenantId: abcLlp.id,
            projectId: project.id,
            blockId: block.id,
            unitName: `Unit ${blockNum}-${unitNum}`,
            unitCode: `UNIT-${blockNum}-${unitNum}`,
            superBuiltupArea: 1200,
            carpetArea: 900,
            unitPrice: 4500000,
            status: 'available',
          },
        });
        
        // Create booking for first unit
        if (unitNum === 1) {
          const booking = await prisma.booking.create({
            data: {
              id: generateId(),
              tenantId: abcLlp.id,
              customerId: customer.id,
              unitId: unit.id,
              bookingDate: new Date(),
              bookingAmount: 450000,
              expectedPossessionDate: new Date('2025-12-31'),
              status: 'confirmed',
            },
          });

          // Create payment plan
          for (let installment = 1; installment <= 3; installment++) {
            await prisma.paymentPlan.create({
              data: {
                id: generateId(),
                tenantId: abcLlp.id,
                bookingId: booking.id,
                installmentNo: installment,
                dueDate: new Date(Date.now() + installment * 90 * 24 * 60 * 60 * 1000),
                amount: 1500000,
                status: 'pending',
              },
            });
          }
          console.log(`    ✓ Unit ${unitNum}-${blockNum} created with booking`);
        } else {
          console.log(`    ✓ Unit ${unitNum}-${blockNum} created`);
        }
      }
    }
    console.log();

    // Create system config
    await prisma.systemConfig.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        configKey: 'app_name',
        configValue: 'VyomTech ERP System',
        dataType: 'string',
        description: 'Application name',
        isGlobal: false,
      },
    });
    console.log('✓ Created system config');

    // Create sample team
    const team = await prisma.team.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        name: 'Sales Team',
        description: 'Default sales team',
      },
    });
    console.log('✓ Created team:', team.id);

    // Create sample document categories
    const docCategory = await prisma.documentCategory.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        categoryName: 'Identity Documents',
        categoryCode: 'ID_DOCS',
        description: 'Government-issued identity documents',
        displayOrder: 1,
        isActive: true,
      },
    });
    console.log('✓ Created document category:', docCategory.id);

    // Create document types
    await prisma.documentType.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        categoryId: docCategory.id,
        typeName: 'Aadhar Card',
        typeCode: 'AADHAR',
        description: 'Indian Aadhar identity card',
        isMandatory: true,
        isIdentityProof: true,
        fileFormats: 'pdf,jpg,png',
        maxFileSize: 5242880, // 5MB
      },
    });
    console.log('✓ Created document type');

    // Create sample loan
    const loan = await prisma.loan.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        loanType: 'Home Loan',
        loanAmount: 5000000,
        interestRate: 7.5,
        tenure: 180,
        status: 'pending',
        loanOfficerId: adminUser.id,
      },
    });
    console.log('✓ Created loan:', loan.id);

    // Create bank financing
    await prisma.bankFinancing.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        loanId: loan.id,
        bankName: 'HDFC Bank',
        bankCode: 'HDFC0000001',
        sanctionAmount: 5000000,
        disbursedAmount: 0,
        status: 'pending',
      },
    });
    console.log('✓ Created bank financing');

    // Create sample employee
    const employee = await prisma.employee.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        userId: adminUser.id,
        firstName: 'Admin',
        lastName: 'User',
        email: 'admin@vyomtech.local',
        designation: 'System Administrator',
        department: 'IT',
        status: 'active',
      },
    });
    console.log('✓ Created employee:', employee.id);

    // Create accounting account
    await prisma.account.create({
      data: {
        id: generateId(),
        tenantId: tenant.id,
        accountCode: '1001',
        accountName: 'Cash',
        accountType: 'Asset',
        subType: 'Current Asset',
        balance: 100000,
      },
    });

    // ============================================================
    // CREATE SYSTEM CONFIGURATION
    // ============================================================
    console.log('⚙️  Creating System Configuration...');
    
    await prisma.systemConfig.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        configKey: 'company_name',
        configValue: 'ABC COMPANY LLP',
      },
    });

    await prisma.systemConfig.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        configKey: 'financial_year_start',
        configValue: '01-04',
      },
    });

    await prisma.systemConfig.create({
      data: {
        id: generateId(),
        tenantId: abcLlp.id,
        configKey: 'gst_applicable',
        configValue: 'true',
      },
    });

    console.log('  ✓ System configurations created\n');

    // ============================================================
    // SUMMARY
    // ============================================================
    console.log('═'.repeat(60));
    console.log('✅ DATABASE SEED COMPLETED SUCCESSFULLY!');
    console.log('═'.repeat(60));
    console.log('\n📊 SUMMARY OF CREATED DATA:\n');
    console.log('👤 SYSTEM ADMIN:');
    console.log('   Email: master.admin@vyomtech.com');
    console.log('   Password: Master@123');
    console.log('   Role: Master Admin (System-wide)\n');

    console.log('🏢 TENANTS:');
    console.log('   1. ABC GROUP (Main Tenant)');
    console.log('   2. ABC COMPANY LLP');
    console.log('   3. ABC COMPANY PVT LTD\n');

    console.log('👥 ABC LLP USERS & CREDENTIALS:');
    userData.forEach((user) => {
      console.log(`   • ${user.name}`);
      console.log(`     Email: ${user.email}`);
      console.log(`     Password: ${user.password}`);
      console.log(`     Role: ${user.role}\n`);
    });

    console.log('📊 DATA CREATED:');
    console.log('   ✓ 5 User Accounts with RBAC Roles');
    console.log('   ✓ 5 HR Employee Records with Payroll');
    console.log('   ✓ 9 Chart of Accounts (GL Structure)');
    console.log('   ✓ 1 Journal Entry with 2 Line Items');
    console.log('   ✓ 1 Vendor & 1 Purchase Order');
    console.log('   ✓ 1 Warehouse with 3 Inventory Items');
    console.log('   ✓ 1 Customer & 1 Sales Order');
    console.log('   ✓ 1 Real Estate Project (Grand Heights Complex)');
    console.log('   ✓ 2 Project Blocks with 3 Units Each');
    console.log('   ✓ 1 Unit Booking with 3-Installment Payment Plan');
    console.log('   ✓ 3 System Configuration Values\n');

    console.log('═'.repeat(60));
    console.log('🚀 System is ready for testing!');
    console.log('═'.repeat(60));

  } catch (error) {
    console.error('❌ Error seeding database:', error);
    throw error;
  }
}

main()
  .catch((error) => {
    console.error(error);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
