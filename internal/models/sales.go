package models

import "time"

// ============================================================================
// SALES LEADS
// ============================================================================

type SalesLead struct {
	ID                  string     `json:"id" db:"id"`
	TenantID            string     `json:"tenant_id" db:"tenant_id"`
	LeadCode            string     `json:"lead_code" db:"lead_code"`
	FirstName           string     `json:"first_name" db:"first_name"`
	LastName            string     `json:"last_name" db:"last_name"`
	Email               string     `json:"email" db:"email"`
	Phone               string     `json:"phone" db:"phone"`
	CompanyName         string     `json:"company_name" db:"company_name"`
	Industry            string     `json:"industry" db:"industry"`
	Status              string     `json:"status" db:"status"` // new, contacted, qualified, negotiation, converted, lost
	DetailedStatus      *string    `json:"detailed_status" db:"detailed_status"`
	PipelineStage       *string    `json:"pipeline_stage" db:"pipeline_stage"`
	Probability         float64    `json:"probability" db:"probability"`
	Source              string     `json:"source" db:"source"` // website, email, phone, referral, event, social
	CampaignID          *string    `json:"campaign_id" db:"campaign_id"`
	AssignedTo          *string    `json:"assigned_to" db:"assigned_to"`
	AssignedDate        *time.Time `json:"assigned_date" db:"assigned_date"`
	ConvertedToCustomer bool       `json:"converted_to_customer" db:"converted_to_customer"`
	CustomerID          *string    `json:"customer_id" db:"customer_id"`
	NextActionDate      *time.Time `json:"next_action_date" db:"next_action_date"`
	NextActionNotes     *string    `json:"next_action_notes" db:"next_action_notes"`
	CaptureDateA        *time.Time `json:"capture_date_a" db:"capture_date_a"`
	CaptureDateB        *time.Time `json:"capture_date_b" db:"capture_date_b"`
	CaptureDateC        *time.Time `json:"capture_date_c" db:"capture_date_c"`
	CaptureDateD        *time.Time `json:"capture_date_d" db:"capture_date_d"`
	LastStatusChange    *time.Time `json:"last_status_change" db:"last_status_change"`
	CreatedBy           *string    `json:"created_by" db:"created_by"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// SALES CUSTOMERS
// ============================================================================

type SalesCustomer struct {
	ID                 string     `json:"id" db:"id"`
	TenantID           string     `json:"tenant_id" db:"tenant_id"`
	CustomerCode       string     `json:"customer_code" db:"customer_code"`
	CustomerName       string     `json:"customer_name" db:"customer_name"`
	BusinessName       string     `json:"business_name" db:"business_name"`
	BusinessType       string     `json:"business_type" db:"business_type"` // individual, proprietorship, partnership, pvt_ltd, public_ltd
	Industry           string     `json:"industry" db:"industry"`
	PrimaryContactName string     `json:"primary_contact_name" db:"primary_contact_name"`
	PrimaryEmail       string     `json:"primary_email" db:"primary_email"`
	PrimaryPhone       string     `json:"primary_phone" db:"primary_phone"`
	BillingAddress     string     `json:"billing_address" db:"billing_address"`
	BillingCity        string     `json:"billing_city" db:"billing_city"`
	BillingState       string     `json:"billing_state" db:"billing_state"`
	BillingCountry     string     `json:"billing_country" db:"billing_country"`
	BillingZip         string     `json:"billing_zip" db:"billing_zip"`
	ShippingAddress    string     `json:"shipping_address" db:"shipping_address"`
	ShippingCity       string     `json:"shipping_city" db:"shipping_city"`
	ShippingState      string     `json:"shipping_state" db:"shipping_state"`
	ShippingCountry    string     `json:"shipping_country" db:"shipping_country"`
	ShippingZip        string     `json:"shipping_zip" db:"shipping_zip"`
	PANNumber          string     `json:"pan_number" db:"pan_number"`
	GSTNumber          string     `json:"gst_number" db:"gst_number"`
	CreditLimit        float64    `json:"credit_limit" db:"credit_limit"`
	CreditDays         int        `json:"credit_days" db:"credit_days"`
	PaymentTerms       string     `json:"payment_terms" db:"payment_terms"`
	CustomerCategory   string     `json:"customer_category" db:"customer_category"` // gold, silver, bronze, regular
	Status             string     `json:"status" db:"status"`                       // active, inactive, blocked
	CurrentBalance     float64    `json:"current_balance" db:"current_balance"`
	CreatedBy          *string    `json:"created_by" db:"created_by"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// SALES QUOTATIONS
// ============================================================================

type SalesQuotation struct {
	ID                 string               `json:"id" db:"id"`
	TenantID           string               `json:"tenant_id" db:"tenant_id"`
	QuotationNumber    string               `json:"quotation_number" db:"quotation_number"`
	CustomerID         string               `json:"customer_id" db:"customer_id"`
	QuotationDate      time.Time            `json:"quotation_date" db:"quotation_date"`
	ValidUntil         *time.Time           `json:"valid_until" db:"valid_until"`
	SubtotalAmount     float64              `json:"subtotal_amount" db:"subtotal_amount"`
	DiscountAmount     float64              `json:"discount_amount" db:"discount_amount"`
	TaxAmount          float64              `json:"tax_amount" db:"tax_amount"`
	TotalAmount        float64              `json:"total_amount" db:"total_amount"`
	Status             string               `json:"status" db:"status"` // draft, sent, accepted, rejected, expired, converted_to_order
	ConvertedToOrder   bool                 `json:"converted_to_order" db:"converted_to_order"`
	SalesOrderID       *string              `json:"sales_order_id" db:"sales_order_id"`
	Notes              *string              `json:"notes" db:"notes"`
	TermsAndConditions *string              `json:"terms_and_conditions" db:"terms_and_conditions"`
	CreatedBy          *string              `json:"created_by" db:"created_by"`
	CreatedAt          time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time           `json:"deleted_at" db:"deleted_at"`
	Items              []SalesQuotationItem `json:"items,omitempty"`
}

type SalesQuotationItem struct {
	ID                 string    `json:"id" db:"id"`
	TenantID           string    `json:"tenant_id" db:"tenant_id"`
	QuotationID        string    `json:"quotation_id" db:"quotation_id"`
	LineNumber         int       `json:"line_number" db:"line_number"`
	Description        string    `json:"description" db:"description"`
	ProductServiceCode string    `json:"product_service_code" db:"product_service_code"`
	Quantity           float64   `json:"quantity" db:"quantity"`
	UnitPrice          float64   `json:"unit_price" db:"unit_price"`
	LineTotal          float64   `json:"line_total" db:"line_total"`
	DiscountPercent    float64   `json:"discount_percent" db:"discount_percent"`
	DiscountAmount     float64   `json:"discount_amount" db:"discount_amount"`
	HSNCode            string    `json:"hsn_code" db:"hsn_code"`
	TaxRate            float64   `json:"tax_rate" db:"tax_rate"`
	TaxAmount          float64   `json:"tax_amount" db:"tax_amount"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// SALES ORDERS
// ============================================================================

type SalesOrder struct {
	ID                   string           `json:"id" db:"id"`
	TenantID             string           `json:"tenant_id" db:"tenant_id"`
	OrderNumber          string           `json:"order_number" db:"order_number"`
	CustomerID           string           `json:"customer_id" db:"customer_id"`
	QuotationID          *string          `json:"quotation_id" db:"quotation_id"`
	OrderDate            time.Time        `json:"order_date" db:"order_date"`
	RequiredByDate       *time.Time       `json:"required_by_date" db:"required_by_date"`
	DeliveryLocation     string           `json:"delivery_location" db:"delivery_location"`
	DeliveryInstructions *string          `json:"delivery_instructions" db:"delivery_instructions"`
	SubtotalAmount       float64          `json:"subtotal_amount" db:"subtotal_amount"`
	DiscountAmount       float64          `json:"discount_amount" db:"discount_amount"`
	TaxAmount            float64          `json:"tax_amount" db:"tax_amount"`
	TotalAmount          float64          `json:"total_amount" db:"total_amount"`
	InvoicedAmount       float64          `json:"invoiced_amount" db:"invoiced_amount"`
	PendingAmount        float64          `json:"pending_amount" db:"pending_amount"`
	Status               string           `json:"status" db:"status"` // draft, confirmed, partially_invoiced, invoiced, partially_delivered, delivered, cancelled
	Notes                *string          `json:"notes" db:"notes"`
	CreatedBy            *string          `json:"created_by" db:"created_by"`
	CreatedAt            time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time       `json:"deleted_at" db:"deleted_at"`
	Items                []SalesOrderItem `json:"items,omitempty"`
}

type SalesOrderItem struct {
	ID                 string    `json:"id" db:"id"`
	TenantID           string    `json:"tenant_id" db:"tenant_id"`
	OrderID            string    `json:"order_id" db:"order_id"`
	LineNumber         int       `json:"line_number" db:"line_number"`
	Description        string    `json:"description" db:"description"`
	ProductServiceCode string    `json:"product_service_code" db:"product_service_code"`
	OrderedQuantity    float64   `json:"ordered_quantity" db:"ordered_quantity"`
	InvoicedQuantity   float64   `json:"invoiced_quantity" db:"invoiced_quantity"`
	UnitPrice          float64   `json:"unit_price" db:"unit_price"`
	LineTotal          float64   `json:"line_total" db:"line_total"`
	DiscountPercent    float64   `json:"discount_percent" db:"discount_percent"`
	DiscountAmount     float64   `json:"discount_amount" db:"discount_amount"`
	HSNCode            string    `json:"hsn_code" db:"hsn_code"`
	TaxRate            float64   `json:"tax_rate" db:"tax_rate"`
	TaxAmount          float64   `json:"tax_amount" db:"tax_amount"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// SALES INVOICES
// ============================================================================

type SalesInvoice struct {
	ID                string             `json:"id" db:"id"`
	TenantID          string             `json:"tenant_id" db:"tenant_id"`
	InvoiceNumber     string             `json:"invoice_number" db:"invoice_number"`
	CustomerID        string             `json:"customer_id" db:"customer_id"`
	SalesOrderID      *string            `json:"sales_order_id" db:"sales_order_id"`
	InvoiceDate       time.Time          `json:"invoice_date" db:"invoice_date"`
	DueDate           *time.Time         `json:"due_date" db:"due_date"`
	SubtotalAmount    float64            `json:"subtotal_amount" db:"subtotal_amount"`
	DiscountAmount    float64            `json:"discount_amount" db:"discount_amount"`
	CGSTAmount        float64            `json:"cgst_amount" db:"cgst_amount"`
	SGSTAmount        float64            `json:"sgst_amount" db:"sgst_amount"`
	IGSTAmount        float64            `json:"igst_amount" db:"igst_amount"`
	TotalTax          float64            `json:"total_tax" db:"total_tax"`
	TotalAmount       float64            `json:"total_amount" db:"total_amount"`
	PaymentStatus     string             `json:"payment_status" db:"payment_status"` // unpaid, partially_paid, paid, overdue, cancelled
	PaidAmount        float64            `json:"paid_amount" db:"paid_amount"`
	PendingAmount     float64            `json:"pending_amount" db:"pending_amount"`
	ARPostingStatus   string             `json:"ar_posting_status" db:"ar_posting_status"` // not_posted, posted, reversed
	GLReferenceNumber *string            `json:"gl_reference_number" db:"gl_reference_number"`
	DocumentStatus    string             `json:"document_status" db:"document_status"` // draft, issued, cancelled
	Notes             *string            `json:"notes" db:"notes"`
	CreatedBy         *string            `json:"created_by" db:"created_by"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time         `json:"deleted_at" db:"deleted_at"`
	Items             []SalesInvoiceItem `json:"items,omitempty"`
}

type SalesInvoiceItem struct {
	ID              string    `json:"id" db:"id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	InvoiceID       string    `json:"invoice_id" db:"invoice_id"`
	LineNumber      int       `json:"line_number" db:"line_number"`
	Description     string    `json:"description" db:"description"`
	HSNCode         string    `json:"hsn_code" db:"hsn_code"`
	Quantity        float64   `json:"quantity" db:"quantity"`
	UnitPrice       float64   `json:"unit_price" db:"unit_price"`
	LineTotal       float64   `json:"line_total" db:"line_total"`
	DiscountPercent float64   `json:"discount_percent" db:"discount_percent"`
	DiscountAmount  float64   `json:"discount_amount" db:"discount_amount"`
	CGSTRate        float64   `json:"cgst_rate" db:"cgst_rate"`
	CGSTAmount      float64   `json:"cgst_amount" db:"cgst_amount"`
	SGSTRate        float64   `json:"sgst_rate" db:"sgst_rate"`
	SGSTAmount      float64   `json:"sgst_amount" db:"sgst_amount"`
	IGSTRate        float64   `json:"igst_rate" db:"igst_rate"`
	IGSTAmount      float64   `json:"igst_amount" db:"igst_amount"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// PAYMENTS
// ============================================================================

type SalesPayment struct {
	ID              string    `json:"id" db:"id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	InvoiceID       string    `json:"invoice_id" db:"invoice_id"`
	PaymentDate     time.Time `json:"payment_date" db:"payment_date"`
	PaymentAmount   float64   `json:"payment_amount" db:"payment_amount"`
	PaymentMethod   string    `json:"payment_method" db:"payment_method"` // cheque, bank_transfer, cash, credit_card, digital_payment
	ReferenceNumber string    `json:"reference_number" db:"reference_number"`
	PaymentStatus   string    `json:"payment_status" db:"payment_status"` // initiated, processed, confirmed, failed, cancelled
	Notes           *string   `json:"notes" db:"notes"`
	CreatedBy       *string   `json:"created_by" db:"created_by"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// CUSTOMER CONTACTS
// ============================================================================

type SalesCustomerContact struct {
	ID               string    `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	CustomerID       string    `json:"customer_id" db:"customer_id"`
	ContactName      string    `json:"contact_name" db:"contact_name"`
	ContactTitle     *string   `json:"contact_title" db:"contact_title"`
	Email            string    `json:"email" db:"email"`
	Phone            string    `json:"phone" db:"phone"`
	Mobile           string    `json:"mobile" db:"mobile"`
	ContactRole      string    `json:"contact_role" db:"contact_role"` // decision_maker, finance, technical, other
	IsPrimaryContact bool      `json:"is_primary_contact" db:"is_primary_contact"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// DELIVERY TRACKING
// ============================================================================

type SalesDeliveryNote struct {
	ID                   string     `json:"id" db:"id"`
	TenantID             string     `json:"tenant_id" db:"tenant_id"`
	OrderID              string     `json:"order_id" db:"order_id"`
	DeliveryDate         time.Time  `json:"delivery_date" db:"delivery_date"`
	DeliveryLocation     string     `json:"delivery_location" db:"delivery_location"`
	DeliveredBy          string     `json:"delivered_by" db:"delivered_by"`
	VehicleNumber        string     `json:"vehicle_number" db:"vehicle_number"`
	DeliveredQuantity    float64    `json:"delivered_quantity" db:"delivered_quantity"`
	DeliveryStatus       string     `json:"delivery_status" db:"delivery_status"` // in_transit, delivered, pending_pod
	PODReceived          bool       `json:"pod_received" db:"pod_received"`
	PODDate              *time.Time `json:"pod_date" db:"pod_date"`
	ReceiverName         *string    `json:"receiver_name" db:"receiver_name"`
	ReceiverSignatureURL *string    `json:"receiver_signature_url" db:"receiver_signature_url"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// CREDIT NOTES
// ============================================================================

type SalesCreditNote struct {
	ID               string                `json:"id" db:"id"`
	TenantID         string                `json:"tenant_id" db:"tenant_id"`
	InvoiceID        string                `json:"invoice_id" db:"invoice_id"`
	CreditNoteNumber string                `json:"credit_note_number" db:"credit_note_number"`
	CreditNoteDate   time.Time             `json:"credit_note_date" db:"credit_note_date"`
	Reason           string                `json:"reason" db:"reason"`
	TotalAmount      float64               `json:"total_amount" db:"total_amount"`
	Status           string                `json:"status" db:"status"`                       // draft, issued, cancelled
	ARPostingStatus  string                `json:"ar_posting_status" db:"ar_posting_status"` // not_posted, posted, reversed
	CreatedBy        *string               `json:"created_by" db:"created_by"`
	CreatedAt        time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at" db:"updated_at"`
	Items            []SalesCreditNoteItem `json:"items,omitempty"`
}

type SalesCreditNoteItem struct {
	ID           string    `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	CreditNoteID string    `json:"credit_note_id" db:"credit_note_id"`
	LineNumber   int       `json:"line_number" db:"line_number"`
	Description  string    `json:"description" db:"description"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	UnitPrice    float64   `json:"unit_price" db:"unit_price"`
	LineTotal    float64   `json:"line_total" db:"line_total"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ============================================================================
// PERFORMANCE METRICS
// ============================================================================

type SalesPerformanceMetrics struct {
	ID                 string    `json:"id" db:"id"`
	TenantID           string    `json:"tenant_id" db:"tenant_id"`
	SalespersonID      string    `json:"salesperson_id" db:"salesperson_id"`
	MetricMonth        string    `json:"metric_month" db:"metric_month"`
	LeadsGenerated     int       `json:"leads_generated" db:"leads_generated"`
	LeadsQualified     int       `json:"leads_qualified" db:"leads_qualified"`
	ConversionRate     float64   `json:"conversion_rate" db:"conversion_rate"`
	TotalOrders        int       `json:"total_orders" db:"total_orders"`
	TotalOrderValue    float64   `json:"total_order_value" db:"total_order_value"`
	AverageOrderValue  float64   `json:"average_order_value" db:"average_order_value"`
	TargetAmount       float64   `json:"target_amount" db:"target_amount"`
	ActualAmount       float64   `json:"actual_amount" db:"actual_amount"`
	AchievementPercent float64   `json:"achievement_percent" db:"achievement_percent"`
	CommissionRate     float64   `json:"commission_rate" db:"commission_rate"`
	CommissionAmount   float64   `json:"commission_amount" db:"commission_amount"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// CAMPAIGNS
// ============================================================================

type SalesCampaign struct {
	ID           string     `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	CampaignCode string     `json:"campaign_code" db:"campaign_code"`
	CampaignName string     `json:"campaign_name" db:"campaign_name"`
	CampaignType string     `json:"campaign_type" db:"campaign_type"` // email, social, referral, event, digital, traditional, direct, outbound
	Description  *string    `json:"description" db:"description"`
	StartDate    time.Time  `json:"start_date" db:"start_date"`
	EndDate      *time.Time `json:"end_date" db:"end_date"`
	Budget       *float64   `json:"budget" db:"budget"`
	ExpectedROI  *float64   `json:"expected_roi" db:"expected_roi"`
	AssignedTo   *string    `json:"assigned_to" db:"assigned_to"`
	Status       string     `json:"status" db:"status"` // active, inactive, completed, paused
	CreatedBy    *string    `json:"created_by" db:"created_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// LEAD SOURCES AND SUBSOURCES
// ============================================================================

type SalesLeadSource struct {
	ID            string     `json:"id" db:"id"`
	TenantID      string     `json:"tenant_id" db:"tenant_id"`
	SourceCode    string     `json:"source_code" db:"source_code"`
	SourceName    string     `json:"source_name" db:"source_name"`
	SourceType    string     `json:"source_type" db:"source_type"` // website, email, phone, referral, event, social, direct, digital, traditional, partner
	SubsourceName *string    `json:"subsource_name" db:"subsource_name"`
	Channel       *string    `json:"channel" db:"channel"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	Description   *string    `json:"description" db:"description"`
	CreatedBy     *string    `json:"created_by" db:"created_by"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// LEAD MILESTONES
// ============================================================================

type SalesLeadMilestone struct {
	ID                string                 `json:"id" db:"id"`
	TenantID          string                 `json:"tenant_id" db:"tenant_id"`
	LeadID            string                 `json:"lead_id" db:"lead_id"`
	MilestoneType     string                 `json:"milestone_type" db:"milestone_type"` // lead_generated, contacted, site_visit, revisit, demo, proposal, negotiation, booking, cancellation, reengaged
	MilestoneDate     time.Time              `json:"milestone_date" db:"milestone_date"`
	MilestoneTime     *string                `json:"milestone_time" db:"milestone_time"`
	Notes             *string                `json:"notes" db:"notes"`
	LocationLatitude  *float64               `json:"location_latitude" db:"location_latitude"`
	LocationLongitude *float64               `json:"location_longitude" db:"location_longitude"`
	LocationName      *string                `json:"location_name" db:"location_name"`
	StatusBefore      *string                `json:"status_before" db:"status_before"`
	StatusAfter       *string                `json:"status_after" db:"status_after"`
	VisitedBy         *string                `json:"visited_by" db:"visited_by"`
	DurationMinutes   *int                   `json:"duration_minutes" db:"duration_minutes"`
	Outcome           *string                `json:"outcome" db:"outcome"` // positive, neutral, negative
	FollowUpDate      *time.Time             `json:"follow_up_date" db:"follow_up_date"`
	FollowUpRequired  bool                   `json:"follow_up_required" db:"follow_up_required"`
	Metadata          map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedBy         *string                `json:"created_by" db:"created_by"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
}

// ============================================================================
// LEAD ENGAGEMENT
// ============================================================================

type SalesLeadEngagement struct {
	ID                string     `json:"id" db:"id"`
	TenantID          string     `json:"tenant_id" db:"tenant_id"`
	LeadID            string     `json:"lead_id" db:"lead_id"`
	EngagementType    string     `json:"engagement_type" db:"engagement_type"` // email_sent, call_made, message_sent, meeting_scheduled, proposal_sent, quote_sent
	EngagementDate    time.Time  `json:"engagement_date" db:"engagement_date"`
	EngagementChannel *string    `json:"engagement_channel" db:"engagement_channel"` // email, phone, sms, whatsapp, in_person, video
	Subject           *string    `json:"subject" db:"subject"`
	Notes             *string    `json:"notes" db:"notes"`
	Status            string     `json:"status" db:"status"` // completed, pending, failed
	ResponseReceived  bool       `json:"response_received" db:"response_received"`
	ResponseDate      *time.Time `json:"response_date" db:"response_date"`
	ResponseNotes     *string    `json:"response_notes" db:"response_notes"`
	AssignedTo        *string    `json:"assigned_to" db:"assigned_to"`
	DurationSeconds   *int       `json:"duration_seconds" db:"duration_seconds"`
	CreatedBy         *string    `json:"created_by" db:"created_by"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}

// ============================================================================
// BOOKINGS
// ============================================================================

type SalesBooking struct {
	ID                       string     `json:"id" db:"id"`
	TenantID                 string     `json:"tenant_id" db:"tenant_id"`
	CustomerID               string     `json:"customer_id" db:"customer_id"`
	LeadID                   *string    `json:"lead_id" db:"lead_id"`
	BookingCode              string     `json:"booking_code" db:"booking_code"`
	BookingDate              time.Time  `json:"booking_date" db:"booking_date"`
	BookingAmount            float64    `json:"booking_amount" db:"booking_amount"`
	UnitType                 *string    `json:"unit_type" db:"unit_type"`
	UnitCount                int        `json:"unit_count" db:"unit_count"`
	UnitsBooked              int        `json:"units_booked" db:"units_booked"`
	UnitsAvailable           *int       `json:"units_available" db:"units_available"`
	DeliveryDate             *time.Time `json:"delivery_date" db:"delivery_date"`
	Status                   string     `json:"status" db:"status"` // confirmed, pending, cancelled, completed, on_hold
	CancellationDate         *time.Time `json:"cancellation_date" db:"cancellation_date"`
	CancellationReason       *string    `json:"cancellation_reason" db:"cancellation_reason"`
	CancellationRefundAmount *float64   `json:"cancellation_refund_amount" db:"cancellation_refund_amount"`
	Notes                    *string    `json:"notes" db:"notes"`
	CreatedBy                *string    `json:"created_by" db:"created_by"`
	CreatedAt                time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at" db:"deleted_at"`
}

// ============================================================================
// ACCOUNT LEDGER
// ============================================================================

type SalesAccountLedger struct {
	ID                    string    `json:"id" db:"id"`
	TenantID              string    `json:"tenant_id" db:"tenant_id"`
	CustomerID            string    `json:"customer_id" db:"customer_id"`
	LedgerCode            string    `json:"ledger_code" db:"ledger_code"`
	LedgerDate            time.Time `json:"ledger_date" db:"ledger_date"`
	TransactionType       string    `json:"transaction_type" db:"transaction_type"` // invoice, payment, credit_note, debit_note, adjustment
	ReferenceDocumentType *string   `json:"reference_document_type" db:"reference_document_type"`
	ReferenceDocumentID   *string   `json:"reference_document_id" db:"reference_document_id"`
	DebitAmount           float64   `json:"debit_amount" db:"debit_amount"`
	CreditAmount          float64   `json:"credit_amount" db:"credit_amount"`
	BalanceAfter          float64   `json:"balance_after" db:"balance_after"`
	Description           *string   `json:"description" db:"description"`
}

// ============================================================================
// LEAD STATUS LOG - Track all status changes
// ============================================================================

type LeadStatusLog struct {
	ID               string    `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	LeadID           string    `json:"lead_id" db:"lead_id"`
	OldStatus        *string   `json:"old_status" db:"old_status"`
	NewStatus        string    `json:"new_status" db:"new_status"`
	OldPipelineStage *string   `json:"old_pipeline_stage" db:"old_pipeline_stage"`
	NewPipelineStage *string   `json:"new_pipeline_stage" db:"new_pipeline_stage"`
	ChangedBy        *string   `json:"changed_by" db:"changed_by"`
	ChangeReason     *string   `json:"change_reason" db:"change_reason"`
	CaptureDateType  *string   `json:"capture_date_type" db:"capture_date_type"` // capture_date_a, capture_date_b, etc
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// ============================================================================
// LEAD PIPELINE CONFIG - Store pipeline configurations
// ============================================================================

type LeadPipelineConfig struct {
	ID            string    `json:"id" db:"id"`
	TenantID      string    `json:"tenant_id" db:"tenant_id"`
	Status        string    `json:"status" db:"status"`
	PipelineStage string    `json:"pipeline_stage" db:"pipeline_stage"`
	Phase         *string   `json:"phase" db:"phase"`
	ColorCode     *string   `json:"color_code" db:"color_code"`
	Icon          *string   `json:"icon" db:"icon"`
	Description   *string   `json:"description" db:"description"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	SortOrder     int       `json:"sort_order" db:"sort_order"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
