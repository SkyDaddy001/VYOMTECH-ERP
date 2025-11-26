package models

import (
	"encoding/json"
	"time"
)

// ============================================================================
// VENDOR MODELS
// ============================================================================

type Vendor struct {
	ID           string          `json:"id" gorm:"primaryKey"`
	TenantID     string          `json:"tenant_id" gorm:"index"`
	VendorCode   string          `json:"vendor_code" gorm:"unique"`
	Name         string          `json:"name"`
	Email        string          `json:"email"`
	Phone        string          `json:"phone"`
	Address      string          `json:"address"`
	City         string          `json:"city"`
	State        string          `json:"state"`
	Country      string          `json:"country"`
	PostalCode   string          `json:"postal_code"`
	TaxID        string          `json:"tax_id"`
	PaymentTerms string          `json:"payment_terms"`
	VendorType   string          `json:"vendor_type"`
	Rating       float64         `json:"rating"`
	IsActive     bool            `json:"is_active"`
	IsBlocked    bool            `json:"is_blocked"`
	CreatedBy    string          `json:"created_by"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at"`
	Status       string          `json:"status"`
	Contacts     []VendorContact `json:"contacts,omitempty" gorm:"foreignKey:VendorID"`
	Addresses    []VendorAddress `json:"addresses,omitempty" gorm:"foreignKey:VendorID"`
}

type VendorContact struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	TenantID  string    `json:"tenant_id"`
	VendorID  string    `json:"vendor_id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

type VendorAddress struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	TenantID  string    `json:"tenant_id"`
	VendorID  string    `json:"vendor_id"`
	Type      string    `json:"type"`
	Line1     string    `json:"line1"`
	Line2     string    `json:"line2"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Country   string    `json:"country"`
	PostCode  string    `json:"post_code"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

// ============================================================================
// PURCHASE ORDER MODELS
// ============================================================================

type PurchaseRequisition struct {
	ID                string     `json:"id" gorm:"primaryKey"`
	TenantID          string     `json:"tenant_id" gorm:"index"`
	RequisitionNumber string     `json:"requisition_number" gorm:"unique"`
	RequesterID       string     `json:"requester_id"`
	Department        string     `json:"department"`
	RequestDate       time.Time  `json:"request_date"`
	RequiredByDate    *time.Time `json:"required_by_date"`
	Purpose           string     `json:"purpose"`
	Status            string     `json:"status"` // Draft, Submitted, Approved, Rejected, Converted_to_PO
	ApprovedBy        *string    `json:"approved_by"`
	ApprovedAt        *time.Time `json:"approved_at"`
	RejectionReason   string     `json:"rejection_reason"`
	CreatedBy         string     `json:"created_by"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
}

type PurchaseOrder struct {
	ID                  string     `json:"id" gorm:"primaryKey"`
	TenantID            string     `json:"tenant_id" gorm:"index"`
	PONumber            string     `json:"po_number" gorm:"unique"`
	VendorID            string     `json:"vendor_id"`
	RequisitionID       *string    `json:"requisition_id"`
	PODate              time.Time  `json:"po_date"`
	DeliveryDate        *time.Time `json:"delivery_date"`
	TotalAmount         float64    `json:"total_amount"`
	TaxAmount           float64    `json:"tax_amount"`
	ShippingAmount      float64    `json:"shipping_amount"`
	DiscountAmount      float64    `json:"discount_amount"`
	NetAmount           float64    `json:"net_amount"`
	PaymentTerms        string     `json:"payment_terms"`
	DeliveryLocation    string     `json:"delivery_location"`
	SpecialInstructions string     `json:"special_instructions"`
	Status              string     `json:"status"` // Draft, Sent, Acknowledged, Partial_Received, Fully_Received, Cancelled, Closed
	SentToVendorAt      *time.Time `json:"sent_to_vendor_at"`
	AcknowledgedAt      *time.Time `json:"acknowledged_at"`
	CreatedBy           string     `json:"created_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`

	// Relations
	Vendor      *Vendor              `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	LineItems   []POLineItem         `json:"line_items,omitempty" gorm:"foreignKey:POID"`
	Requisition *PurchaseRequisition `json:"requisition,omitempty" gorm:"foreignKey:RequisitionID"`
}

type POLineItem struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	TenantID      string    `json:"tenant_id"`
	POID          string    `json:"po_id"`
	LineNumber    int       `json:"line_number"`
	ProductCode   string    `json:"product_code"`
	Description   string    `json:"description"`
	Quantity      float64   `json:"quantity"`
	Unit          string    `json:"unit"`
	UnitPrice     float64   `json:"unit_price"`
	LineTotal     float64   `json:"line_total"`
	HSNCode       string    `json:"hsn_code"`
	TaxRate       float64   `json:"tax_rate"`
	TaxAmount     float64   `json:"tax_amount"`
	Specification string    `json:"specification"`
	CreatedAt     time.Time `json:"created_at"`
}

// ============================================================================
// GOODS RECEIPT & QUALITY MODELS
// ============================================================================

type GoodsReceipt struct {
	ID                    string     `json:"id" gorm:"primaryKey"`
	TenantID              string     `json:"tenant_id" gorm:"index"`
	GRNNumber             string     `json:"grn_number" gorm:"unique"`
	POID                  string     `json:"po_id"`
	ReceiptDate           time.Time  `json:"receipt_date"`
	ReceivedBy            string     `json:"received_by"`
	TotalQuantityReceived float64    `json:"total_quantity_received"`
	TotalQuantityAccepted float64    `json:"total_quantity_accepted"`
	TotalQuantityRejected float64    `json:"total_quantity_rejected"`
	DeliveryNoteNumber    string     `json:"delivery_note_number"`
	VehicleNumber         string     `json:"vehicle_number"`
	DriverName            string     `json:"driver_name"`
	DriverPhone           string     `json:"driver_phone"`
	Remarks               string     `json:"remarks"`
	Status                string     `json:"status"`    // Received, QC_In_Progress, QC_Passed, QC_Failed, Partial_Accepted, Rejected
	QCStatus              string     `json:"qc_status"` // Pending, In_Progress, Passed, Failed, Partial
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"deleted_at"`

	// Relations
	PurchaseOrder *PurchaseOrder        `json:"purchase_order,omitempty" gorm:"foreignKey:POID"`
	LineItems     []GRNLineItem         `json:"line_items,omitempty" gorm:"foreignKey:GRNID"`
	Inspections   []QualityInspection   `json:"inspections,omitempty" gorm:"foreignKey:GRNID"`
	MaterialNotes []MaterialReceiptNote `json:"material_notes,omitempty" gorm:"foreignKey:GRNID"`
}

type GRNLineItem struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	TenantID         string     `json:"tenant_id"`
	GRNID            string     `json:"grn_id"`
	POLineItemID     *string    `json:"po_line_item_id"`
	LineNumber       int        `json:"line_number"`
	ProductCode      string     `json:"product_code"`
	Description      string     `json:"description"`
	POQuantity       float64    `json:"po_quantity"`
	ReceivedQuantity float64    `json:"received_quantity"`
	AcceptedQuantity float64    `json:"accepted_quantity"`
	RejectedQuantity float64    `json:"rejected_quantity"`
	Unit             string     `json:"unit"`
	RejectionReason  string     `json:"rejection_reason"`
	BatchNumber      string     `json:"batch_number"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	CreatedAt        time.Time  `json:"created_at"`
}

type QualityInspection struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	TenantID          string    `json:"tenant_id"`
	GRNID             string    `json:"grn_id"`
	GRNLineItemID     *string   `json:"grn_line_item_id"`
	InspectionDate    time.Time `json:"inspection_date"`
	InspectedBy       string    `json:"inspected_by"`
	InspectionType    string    `json:"inspection_type"` // Visual, Functional, Dimensional, Lab_Test, Batch_Test
	Status            string    `json:"status"`          // Passed, Failed, Partial_Pass, Conditional_Pass
	QuantityInspected float64   `json:"quantity_inspected"`
	QuantityPassed    float64   `json:"quantity_passed"`
	QuantityFailed    float64   `json:"quantity_failed"`
	DefectsFound      string    `json:"defects_found"`
	QualityScore      float64   `json:"quality_score"`
	Notes             string    `json:"notes"`
	CertificateNumber string    `json:"certificate_number"`
	CreatedAt         time.Time `json:"created_at"`
}

// ============================================================================
// MATERIAL RECEIPT NOTE (MRN) MODELS
// ============================================================================

type MaterialReceiptNote struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	TenantID              string    `json:"tenant_id" gorm:"index"`
	MRNNumber             string    `json:"mrn_number" gorm:"unique"`
	GRNID                 string    `json:"grn_id"`
	WarehouseID           *string   `json:"warehouse_id"`
	ReceiptDate           time.Time `json:"receipt_date"`
	AcceptedBy            string    `json:"accepted_by"`
	StorageLocation       string    `json:"storage_location"`
	TotalQuantityAccepted float64   `json:"total_quantity_accepted"`
	TotalQuantityInStock  float64   `json:"total_quantity_in_stock"`
	Status                string    `json:"status"` // Received, Stored, Available
	Remarks               string    `json:"remarks"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Relations
	GoodsReceipt *GoodsReceipt `json:"goods_receipt,omitempty" gorm:"foreignKey:GRNID"`
	LineItems    []MRNLineItem `json:"line_items,omitempty" gorm:"foreignKey:MRNID"`
}

type MRNLineItem struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	TenantID         string     `json:"tenant_id"`
	MRNID            string     `json:"mrn_id"`
	GRNLineItemID    *string    `json:"grn_line_item_id"`
	LineNumber       int        `json:"line_number"`
	ProductCode      string     `json:"product_code"`
	Description      string     `json:"description"`
	QuantityAccepted float64    `json:"quantity_accepted"`
	Unit             string     `json:"unit"`
	StorageLocation  string     `json:"storage_location"`
	BatchNumber      string     `json:"batch_number"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	WarehouseNotes   string     `json:"warehouse_notes"`
	CreatedAt        time.Time  `json:"created_at"`
}

// ============================================================================
// CONTRACT MODELS
// ============================================================================

type Contract struct {
	ID                 string     `json:"id" gorm:"primaryKey"`
	TenantID           string     `json:"tenant_id" gorm:"index"`
	ContractNumber     string     `json:"contract_number" gorm:"unique"`
	VendorID           string     `json:"vendor_id"`
	ContractDate       time.Time  `json:"contract_date"`
	StartDate          *time.Time `json:"start_date"`
	EndDate            *time.Time `json:"end_date"`
	ContractType       string     `json:"contract_type"` // Material, Labour, Service, Hybrid
	BOQID              *string    `json:"boq_id"`
	TotalContractValue float64    `json:"total_contract_value"`
	Currency           string     `json:"currency"`
	PaymentTerms       string     `json:"payment_terms"`
	DeliverySchedule   string     `json:"delivery_schedule"`
	ContractStatus     string     `json:"contract_status"` // Draft, Sent, Accepted, Active, Completed, Cancelled
	ContractFilePath   string     `json:"contract_file_path"`
	SignedDate         *time.Time `json:"signed_date"`
	SignedByVendor     bool       `json:"signed_by_vendor"`
	SignedByCompany    bool       `json:"signed_by_company"`
	CompanySignatory   *string    `json:"company_signatory"`
	VendorSignatory    string     `json:"vendor_signatory"`
	CreatedBy          string     `json:"created_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Relations
	Vendor    *Vendor            `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	LineItems []ContractLineItem `json:"line_items,omitempty" gorm:"foreignKey:ContractID"`
	Materials []ContractMaterial `json:"materials,omitempty" gorm:"foreignKey:ContractID"`
	Labour    []ContractLabour   `json:"labour,omitempty" gorm:"foreignKey:ContractID"`
	Services  []ContractService  `json:"services,omitempty" gorm:"foreignKey:ContractID"`
}

type ContractLineItem struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	TenantID         string     `json:"tenant_id"`
	ContractID       string     `json:"contract_id"`
	LineNumber       int        `json:"line_number"`
	ItemCode         string     `json:"item_code"`
	ItemDescription  string     `json:"item_description"`
	ItemType         string     `json:"item_type"` // Material, Labour, Service, Hybrid
	Quantity         float64    `json:"quantity"`
	Unit             string     `json:"unit"`
	UnitPrice        float64    `json:"unit_price"`
	LineTotal        float64    `json:"line_total"`
	Specification    string     `json:"specification"`
	DeliveryLocation string     `json:"delivery_location"`
	DeliveryDate     *time.Time `json:"delivery_date"`
	CreatedAt        time.Time  `json:"created_at"`
}

type ContractMaterial struct {
	ID                  string    `json:"id" gorm:"primaryKey"`
	TenantID            string    `json:"tenant_id"`
	ContractID          string    `json:"contract_id"`
	MaterialCode        string    `json:"material_code"`
	MaterialDescription string    `json:"material_description"`
	Quantity            float64   `json:"quantity"`
	Unit                string    `json:"unit"`
	UnitPrice           float64   `json:"unit_price"`
	TotalPrice          float64   `json:"total_price"`
	HSNCode             string    `json:"hsn_code"`
	SupplierNotes       string    `json:"supplier_notes"`
	CreatedAt           time.Time `json:"created_at"`
}

type ContractLabour struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	TenantID        string    `json:"tenant_id"`
	ContractID      string    `json:"contract_id"`
	SkillType       string    `json:"skill_type"`
	LabourCategory  string    `json:"labour_category"` // Skilled, Semi-Skilled, Unskilled
	NumberOfWorkers int       `json:"number_of_workers"`
	DurationDays    int       `json:"duration_days"`
	DailyRate       float64   `json:"daily_rate"`
	TotalLabourCost float64   `json:"total_labour_cost"`
	WorkDescription string    `json:"work_description"`
	LabourNotes     string    `json:"labour_notes"`
	CreatedAt       time.Time `json:"created_at"`
}

type ContractService struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	TenantID              string    `json:"tenant_id"`
	ContractID            string    `json:"contract_id"`
	ServiceCode           string    `json:"service_code"`
	ServiceDescription    string    `json:"service_description"`
	ServiceType           string    `json:"service_type"`
	UnitOfService         string    `json:"unit_of_service"` // hour, day, project, etc
	Quantity              float64   `json:"quantity"`
	UnitPrice             float64   `json:"unit_price"`
	TotalServiceCost      float64   `json:"total_service_cost"`
	ServiceLevelAgreement string    `json:"service_level_agreement"`
	ServiceNotes          string    `json:"service_notes"`
	CreatedAt             time.Time `json:"created_at"`
}

// ============================================================================
// INVOICE & PAYMENT MODELS
// ============================================================================

type VendorInvoice struct {
	ID              string     `json:"id" gorm:"primaryKey"`
	TenantID        string     `json:"tenant_id" gorm:"index"`
	InvoiceNumber   string     `json:"invoice_number" gorm:"unique"`
	VendorID        string     `json:"vendor_id"`
	POID            *string    `json:"po_id"`
	GRNID           *string    `json:"grn_id"`
	InvoiceDate     time.Time  `json:"invoice_date"`
	DueDate         *time.Time `json:"due_date"`
	InvoiceAmount   float64    `json:"invoice_amount"`
	TaxAmount       float64    `json:"tax_amount"`
	DiscountAmount  float64    `json:"discount_amount"`
	TotalPayable    float64    `json:"total_payable"`
	Status          string     `json:"status"`         // Received, Approved, Rejected, Paid, Partially_Paid
	MatchedStatus   string     `json:"matched_status"` // Not_Matched, PO_Matched, GRN_Matched, Three_Way_Match
	ThreeWayMatch   bool       `json:"three_way_match"`
	ReceivedAt      *time.Time `json:"received_at"`
	ApprovedAt      *time.Time `json:"approved_at"`
	ApprovedBy      *string    `json:"approved_by"`
	RejectionReason string     `json:"rejection_reason"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relations
	Vendor    *Vendor           `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	LineItems []InvoiceLineItem `json:"line_items,omitempty" gorm:"foreignKey:InvoiceID"`
	Payments  []VendorPayment   `json:"payments,omitempty" gorm:"foreignKey:InvoiceID"`
}

type PurchaseInvoiceLineItem struct {
	ID          string    `json:"id" db:"id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	InvoiceID   string    `json:"invoice_id" db:"invoice_id"`
	LineNumber  int       `json:"line_number" db:"line_number"`
	Description string    `json:"description" db:"description"`
	Quantity    float64   `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	LineTotal   float64   `json:"line_total" db:"line_total"`
	HSNCode     string    `json:"hsn_code" db:"hsn_code"`
	TaxRate     float64   `json:"tax_rate" db:"tax_rate"`
	TaxAmount   float64   `json:"tax_amount" db:"tax_amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type VendorPayment struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	TenantID        string    `json:"tenant_id"`
	InvoiceID       string    `json:"invoice_id"`
	PaymentDate     time.Time `json:"payment_date"`
	PaymentAmount   float64   `json:"payment_amount"`
	PaymentMethod   string    `json:"payment_method"` // Cheque, Bank_Transfer, Cash, Credit_Card, Digital_Payment
	ReferenceNumber string    `json:"reference_number"`
	PaymentStatus   string    `json:"payment_status"` // Initiated, Processed, Confirmed, Cancelled
	PaidBy          *string   `json:"paid_by"`
	CreatedAt       time.Time `json:"created_at"`
}

// ============================================================================
// VENDOR PERFORMANCE & AUDIT MODELS
// ============================================================================

type VendorPerformanceMetrics struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	TenantID              string    `json:"tenant_id"`
	VendorID              string    `json:"vendor_id"`
	MetricMonth           time.Time `json:"metric_month"`
	TotalOrders           int       `json:"total_orders"`
	OrdersDeliveredOnTime int       `json:"orders_delivered_on_time"`
	OnTimeDeliveryRate    float64   `json:"on_time_delivery_rate"`
	TotalGRN              int       `json:"total_grn"`
	GRNAcceptedFirstTime  int       `json:"grn_accepted_first_time"`
	QualityAcceptanceRate float64   `json:"quality_acceptance_rate"`
	TotalInvoices         int       `json:"total_invoices"`
	InvoiceDiscrepancies  int       `json:"invoice_discrepancies"`
	AvgResponseTimeHours  float64   `json:"avg_response_time_hours"`
	OverallRating         float64   `json:"overall_rating"`
	Notes                 string    `json:"notes"`
	CalculatedAt          time.Time `json:"calculated_at"`
}

type PurchaseApproval struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	TenantID       string     `json:"tenant_id"`
	POID           *string    `json:"po_id"`
	ApprovalLevel  int        `json:"approval_level"`
	ApproverID     string     `json:"approver_id"`
	ApprovalStatus string     `json:"approval_status"` // Pending, Approved, Rejected, On_Hold
	ApprovalDate   *time.Time `json:"approval_date"`
	Comments       string     `json:"comments"`
	CreatedAt      time.Time  `json:"created_at"`
}

type PurchaseAuditLog struct {
	ID         string          `json:"id" gorm:"primaryKey"`
	TenantID   string          `json:"tenant_id"`
	EntityType string          `json:"entity_type"`
	EntityID   string          `json:"entity_id"`
	ActionType string          `json:"action_type"`
	OldValues  json.RawMessage `json:"old_values" gorm:"type:json"`
	NewValues  json.RawMessage `json:"new_values" gorm:"type:json"`
	ChangedBy  *string         `json:"changed_by"`
	CreatedAt  time.Time       `json:"created_at"`
}

// ============================================================================
// TABLE NAME MAPPINGS
// ============================================================================

func (Vendor) TableName() string                   { return "vendors" }
func (VendorContact) TableName() string            { return "vendor_contacts" }
func (VendorAddress) TableName() string            { return "vendor_addresses" }
func (PurchaseRequisition) TableName() string      { return "purchase_requisitions" }
func (PurchaseOrder) TableName() string            { return "purchase_orders" }
func (POLineItem) TableName() string               { return "po_line_items" }
func (GoodsReceipt) TableName() string             { return "goods_receipts" }
func (GRNLineItem) TableName() string              { return "grn_line_items" }
func (QualityInspection) TableName() string        { return "quality_inspections" }
func (MaterialReceiptNote) TableName() string      { return "material_receipt_notes" }
func (MRNLineItem) TableName() string              { return "mrn_line_items" }
func (Contract) TableName() string                 { return "contracts" }
func (ContractLineItem) TableName() string         { return "contract_line_items" }
func (ContractMaterial) TableName() string         { return "contract_materials" }
func (ContractLabour) TableName() string           { return "contract_labour" }
func (ContractService) TableName() string          { return "contract_services" }
func (VendorInvoice) TableName() string            { return "vendor_invoices" }
func (InvoiceLineItem) TableName() string          { return "invoice_line_items" }
func (VendorPayment) TableName() string            { return "vendor_payments" }
func (VendorPerformanceMetrics) TableName() string { return "vendor_performance_metrics" }
func (PurchaseApproval) TableName() string         { return "purchase_approvals" }
func (PurchaseAuditLog) TableName() string         { return "purchase_audit_log" }
