package constants

// Permission codes format: module.action
// All permission codes should be defined here for consistency

// Sales Module Permissions
const (
	SalesLeadCreate = "leads.create"
	SalesLeadRead   = "leads.read"
	SalesLeadUpdate = "leads.update"
	SalesLeadDelete = "leads.delete"

	SalesCustomerCreate = "customers.create"
	SalesCustomerRead   = "customers.read"
	SalesCustomerUpdate = "customers.update"
	SalesCustomerDelete = "customers.delete"

	SalesInvoiceCreate = "invoices.create"
	SalesInvoiceRead   = "invoices.read"
	SalesInvoiceUpdate = "invoices.update"
	SalesInvoiceDelete = "invoices.delete"

	SalesPaymentCreate = "payments.create"
	SalesPaymentRead   = "payments.read"

	SalesReportExport = "sales.reports.export"
	SalesReportView   = "sales.reports.view"
)

// HR Module Permissions
const (
	EmployeeCreate = "employees.create"
	EmployeeRead   = "employees.read"
	EmployeeUpdate = "employees.update"
	EmployeeDelete = "employees.delete"

	PayrollExecute = "payroll.execute"
	PayrollView    = "payroll.view"

	SalaryView   = "salary.view"
	SalaryUpdate = "salary.update"

	HRReportExport = "hr.reports.export"
	HRReportView   = "hr.reports.view"
)

// GL Module Permissions
const (
	AccountCreate = "accounts.create"
	AccountRead   = "accounts.read"
	AccountUpdate = "accounts.update"
	AccountDelete = "accounts.delete"

	EntryPost    = "entries.post"
	EntryReverse = "entries.reverse"

	ReconcileExecute = "reconciliation.execute"

	PeriodClose = "period.close"
	PeriodOpen  = "period.open"

	GLReportExport = "gl.reports.export"
	GLReportView   = "gl.reports.view"
)

// Purchase Module Permissions
const (
	PurchaseOrderCreate  = "purchase_orders.create"
	PurchaseOrderRead    = "purchase_orders.read"
	PurchaseOrderUpdate  = "purchase_orders.update"
	PurchaseOrderDelete  = "purchase_orders.delete"
	PurchaseOrderApprove = "purchase_orders.approve"

	VendorCreate = "vendors.create"
	VendorRead   = "vendors.read"
	VendorUpdate = "vendors.update"
	VendorDelete = "vendors.delete"

	ReceiptCreate = "receipts.create"
	ReceiptRead   = "receipts.read"
)

// Real Estate Module Permissions
const (
	ProjectCreate = "projects.create"
	ProjectRead   = "projects.read"
	ProjectUpdate = "projects.update"
	ProjectDelete = "projects.delete"

	UnitCreate = "units.create"
	UnitRead   = "units.read"
	UnitUpdate = "units.update"

	BookingCreate = "bookings.create"
	BookingRead   = "bookings.read"
	BookingUpdate = "bookings.update"
)

// Construction Module Permissions
const (
	ConstructionProjectCreate = "construction.projects.create"
	ConstructionProjectRead   = "construction.projects.read"
	ConstructionProjectUpdate = "construction.projects.update"

	BOQCreate = "boq.create"
	BOQRead   = "boq.read"
	BOQUpdate = "boq.update"
)

// Admin Permissions
const (
	UsersCreate = "users.create"
	UsersRead   = "users.read"
	UsersUpdate = "users.update"
	UsersDelete = "users.delete"

	RolesCreate = "roles.create"
	RolesRead   = "roles.read"
	RolesUpdate = "roles.update"
	RolesDelete = "roles.delete"

	AuditView   = "audit.view"
	AuditExport = "audit.export"

	SettingsManage = "settings.manage"
)
