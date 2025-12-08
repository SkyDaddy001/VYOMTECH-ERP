package handlers

import "time"

// ============================================================================
// API Response Structures - Consistent across all modules
// ============================================================================

// APIResponse is the standard API response wrapper
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// ============================================================================
// Dashboard Response Structures
// ============================================================================

// FinancialDashboardResponse contains financial metrics
type FinancialDashboardResponse struct {
	TotalRevenue       float64       `json:"total_revenue"`
	TotalExpenses      float64       `json:"total_expenses"`
	NetProfit          float64       `json:"net_profit"`
	ProfitMargin       float64       `json:"profit_margin"`
	CashPosition       float64       `json:"cash_position"`
	AccountsReceivable float64       `json:"accounts_receivable"`
	AccountsPayable    float64       `json:"accounts_payable"`
	TrialBalance       TrialBalance  `json:"trial_balance"`
	RevenueByMonth     []MonthlyData `json:"revenue_by_month"`
	TopExpenses        []ExpenseItem `json:"top_expenses"`
}

// TrialBalance represents GL trial balance
type TrialBalance struct {
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	Balance     float64 `json:"balance"`
	IsBalanced  bool    `json:"is_balanced"`
}

// MonthlyData represents monthly financial data
type MonthlyData struct {
	Month   string  `json:"month"`
	Amount  float64 `json:"amount"`
	Percent float64 `json:"percent"`
}

// ExpenseItem represents an expense line item
type ExpenseItem struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Percent  float64 `json:"percent"`
}

// SalesDashboardResponse contains sales metrics
type SalesDashboardResponse struct {
	TotalQuotations   int64                `json:"total_quotations"`
	TotalOrders       int64                `json:"total_orders"`
	TotalSales        float64              `json:"total_sales"`
	ConversionRate    float64              `json:"conversion_rate"`
	AverageOrderValue float64              `json:"average_order_value"`
	TopProducts       []ProductSale        `json:"top_products"`
	SalesByMonth      []MonthlyData        `json:"sales_by_month"`
	OrderStatus       OrderStatusBreakdown `json:"order_status"`
	RecentOrders      []OrderSummary       `json:"recent_orders"`
}

// ProductSale represents a product sale summary
type ProductSale struct {
	ProductName string  `json:"product_name"`
	Quantity    int64   `json:"quantity"`
	Amount      float64 `json:"amount"`
	Percent     float64 `json:"percent"`
}

// OrderStatusBreakdown represents order status distribution
type OrderStatusBreakdown struct {
	Pending   int64 `json:"pending"`
	Confirmed int64 `json:"confirmed"`
	Shipped   int64 `json:"shipped"`
	Delivered int64 `json:"delivered"`
}

// OrderSummary represents a summary of an order
type OrderSummary struct {
	OrderNo      string    `json:"order_no"`
	CustomerName string    `json:"customer_name"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	OrderDate    time.Time `json:"order_date"`
}

// HRDashboardResponse contains HR metrics
type HRDashboardResponse struct {
	TotalEmployees  int64                 `json:"total_employees"`
	ActiveEmployees int64                 `json:"active_employees"`
	AverageSalary   float64               `json:"average_salary"`
	TotalPayroll    float64               `json:"total_payroll"`
	EmployeesByDept []DepartmentData      `json:"employees_by_dept"`
	LeaveBalance    map[string]LeaveStats `json:"leave_balance"`
	Attendance      AttendanceStats       `json:"attendance"`
	UpcomingLeaves  []LeaveEvent          `json:"upcoming_leaves"`
	PayrollStatus   PayrollSummary        `json:"payroll_status"`
}

// DepartmentData represents department statistics
type DepartmentData struct {
	Department string  `json:"department"`
	HeadCount  int64   `json:"head_count"`
	AverageSal float64 `json:"average_salary"`
	Percent    float64 `json:"percent"`
}

// LeaveStats represents leave statistics
type LeaveStats struct {
	LeaveType    string `json:"leave_type"`
	TotalAllowed int64  `json:"total_allowed"`
	Used         int64  `json:"used"`
	Balance      int64  `json:"balance"`
	Pending      int64  `json:"pending"`
}

// AttendanceStats represents attendance statistics
type AttendanceStats struct {
	Present    int64   `json:"present"`
	Absent     int64   `json:"absent"`
	OnLeave    int64   `json:"on_leave"`
	Percentage float64 `json:"percentage"`
}

// LeaveEvent represents a leave event
type LeaveEvent struct {
	EmployeeName string    `json:"employee_name"`
	LeaveType    string    `json:"leave_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Days         int64     `json:"days"`
	Status       string    `json:"status"`
}

// PayrollSummary represents payroll summary
type PayrollSummary struct {
	LastPayroll        time.Time `json:"last_payroll"`
	NextPayroll        time.Time `json:"next_payroll"`
	TotalEmployees     int64     `json:"total_employees"`
	ProcessedCount     int64     `json:"processed_count"`
	PendingCount       int64     `json:"pending_count"`
	TotalPayrollAmount float64   `json:"total_payroll_amount"`
	Status             string    `json:"status"`
}

// RealEstateDashboardResponse contains real estate metrics
type RealEstateDashboardResponse struct {
	TotalProjects     int64             `json:"total_projects"`
	OngoingProjects   int64             `json:"ongoing_projects"`
	CompletedProjects int64             `json:"completed_projects"`
	TotalUnits        int64             `json:"total_units"`
	SoldUnits         int64             `json:"sold_units"`
	AvailableUnits    int64             `json:"available_units"`
	TotalRevenue      float64           `json:"total_revenue"`
	AvgUnitPrice      float64           `json:"avg_unit_price"`
	ProjectProgress   []ProjectProgress `json:"project_progress"`
	SalesPerformance  SalesPerformance  `json:"sales_performance"`
	TopProjects       []TopProject      `json:"top_projects"`
}

// ProjectProgress represents project progress data
type ProjectProgress struct {
	ProjectName    string `json:"project_name"`
	Location       string `json:"location"`
	Progress       int64  `json:"progress"`
	CompletedUnits int64  `json:"completed_units"`
	SoldUnits      int64  `json:"sold_units"`
	Status         string `json:"status"`
}

// SalesPerformance represents sales performance data
type SalesPerformance struct {
	TotalListings  int64   `json:"total_listings"`
	ActiveListings int64   `json:"active_listings"`
	ConversionRate float64 `json:"conversion_rate"`
	AveragePrice   float64 `json:"average_price"`
	MonthlyGrowth  float64 `json:"monthly_growth"`
}

// TopProject represents a top performing project
type TopProject struct {
	ProjectName string  `json:"project_name"`
	Location    string  `json:"location"`
	Revenue     float64 `json:"revenue"`
	SoldUnits   int64   `json:"sold_units"`
	Percent     float64 `json:"percent"`
}

// ConstructionDashboardResponse contains construction metrics
type ConstructionDashboardResponse struct {
	TotalProjects      int64                `json:"total_projects"`
	OngoingProjects    int64                `json:"ongoing_projects"`
	CompletedProjects  int64                `json:"completed_projects"`
	TotalValue         float64              `json:"total_value"`
	AvgCompletion      float64              `json:"avg_completion"`
	ProjectMetrics     []ConstructionMetric `json:"project_metrics"`
	BOQSummary         BOQSummaryData       `json:"boq_summary"`
	SafetyMetrics      SafetyData           `json:"safety_metrics"`
	ResourceAllocation ResourceData         `json:"resource_allocation"`
}

// ConstructionMetric represents construction project metrics
type ConstructionMetric struct {
	ProjectName      string  `json:"project_name"`
	Location         string  `json:"location"`
	Progress         int64   `json:"progress"`
	Value            float64 `json:"value"`
	WorkersOnSite    int64   `json:"workers_on_site"`
	ScheduleVariance float64 `json:"schedule_variance"`
	CostVariance     float64 `json:"cost_variance"`
}

// BOQSummaryData represents BOQ summary
type BOQSummaryData struct {
	TotalItems     int64   `json:"total_items"`
	CompletedItems int64   `json:"completed_items"`
	PendingItems   int64   `json:"pending_items"`
	CompletionRate float64 `json:"completion_rate"`
	TotalValue     float64 `json:"total_value"`
}

// SafetyData represents safety metrics
type SafetyData struct {
	TotalIncidents   int64   `json:"total_incidents"`
	LostTimeInjuries int64   `json:"lost_time_injuries"`
	SafetyScore      float64 `json:"safety_score"`
	Status           string  `json:"status"`
}

// ResourceData represents resource allocation
type ResourceData struct {
	TotalWorkers   int64   `json:"total_workers"`
	EquipmentCount int64   `json:"equipment_count"`
	VehicleCount   int64   `json:"vehicle_count"`
	Utilization    float64 `json:"utilization"`
}

// ComplianceDashboardResponse contains compliance metrics
type ComplianceDashboardResponse struct {
	OverallStatus     string             `json:"overall_status"`
	ComplianceScore   float64            `json:"compliance_score"`
	RERAStatus        ComplianceStatus   `json:"rera_status"`
	TaxStatus         ComplianceStatus   `json:"tax_status"`
	LaborStatus       ComplianceStatus   `json:"labor_status"`
	PendingActions    []ComplianceAction `json:"pending_actions"`
	UpcomingDeadlines []Deadline         `json:"upcoming_deadlines"`
	AuditHistory      []AuditRecord      `json:"audit_history"`
}

// ComplianceStatus represents status of a compliance area
type ComplianceStatus struct {
	Area              string    `json:"area"`
	Status            string    `json:"status"`
	LastAudit         time.Time `json:"last_audit"`
	NextAudit         time.Time `json:"next_audit"`
	IssueCount        int64     `json:"issue_count"`
	CompliancePercent float64   `json:"compliance_percent"`
}

// ComplianceAction represents a pending compliance action
type ComplianceAction struct {
	Action   string    `json:"action"`
	Type     string    `json:"type"`
	DueDate  time.Time `json:"due_date"`
	Priority string    `json:"priority"`
	Status   string    `json:"status"`
	Owner    string    `json:"owner"`
}

// Deadline represents an upcoming deadline
type Deadline struct {
	Description   string    `json:"description"`
	DueDate       time.Time `json:"due_date"`
	DaysRemaining int       `json:"days_remaining"`
	Priority      string    `json:"priority"`
	Responsible   string    `json:"responsible"`
}

// AuditRecord represents an audit record
type AuditRecord struct {
	Area      string    `json:"area"`
	AuditDate time.Time `json:"audit_date"`
	AuditedBy string    `json:"audited_by"`
	Findings  int64     `json:"findings"`
	Status    string    `json:"status"`
	FollowUp  time.Time `json:"follow_up"`
}

// CallCenterDashboardResponse contains call center metrics
type CallCenterDashboardResponse struct {
	TotalAgents      int64              `json:"total_agents"`
	ActiveAgents     int64              `json:"active_agents"`
	TotalCalls       int64              `json:"total_calls"`
	IncomingCalls    int64              `json:"incoming_calls"`
	OutgoingCalls    int64              `json:"outgoing_calls"`
	MissedCalls      int64              `json:"missed_calls"`
	AvgCallDuration  string             `json:"avg_call_duration"`
	CallConversion   float64            `json:"call_conversion"`
	AgentPerformance []AgentPerformance `json:"agent_performance"`
	CallsByHour      []HourlyData       `json:"calls_by_hour"`
	TopPerformers    []AgentRanking     `json:"top_performers"`
	QueueStatus      QueueStats         `json:"queue_status"`
}

// AgentPerformance represents individual agent performance
type AgentPerformance struct {
	AgentName      string  `json:"agent_name"`
	TotalCalls     int64   `json:"total_calls"`
	ConvertedLeads int64   `json:"converted_leads"`
	ConversionRate float64 `json:"conversion_rate"`
	AvgDuration    string  `json:"avg_duration"`
	Rating         float64 `json:"rating"`
	Status         string  `json:"status"`
}

// HourlyData represents hourly call data
type HourlyData struct {
	Hour  string `json:"hour"`
	Calls int64  `json:"calls"`
}

// AgentRanking represents agent ranking
type AgentRanking struct {
	Rank           int     `json:"rank"`
	AgentName      string  `json:"agent_name"`
	ConversionRate float64 `json:"conversion_rate"`
	Score          float64 `json:"score"`
}

// QueueStats represents queue statistics
type QueueStats struct {
	TotalInQueue   int64   `json:"total_in_queue"`
	WaitTime       string  `json:"wait_time"`
	AbandonedCalls int64   `json:"abandoned_calls"`
	ServiceLevel   float64 `json:"service_level"`
}

// ============================================================================
// Error Response
// ============================================================================

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ============================================================================
// List Response Structures
// ============================================================================

// LeadListResponse represents a list of leads
type LeadListResponse struct {
	Items []LeadItem `json:"items"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Total int64      `json:"total"`
}

// LeadItem represents a lead in the list
type LeadItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Status      string    `json:"status"`
	Source      string    `json:"source"`
	Score       int64     `json:"score"`
	LastContact time.Time `json:"last_contact"`
	CreatedAt   time.Time `json:"created_at"`
}

// CallListResponse represents a list of calls
type CallListResponse struct {
	Items []CallItem `json:"items"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Total int64      `json:"total"`
}

// CallItem represents a call in the list
type CallItem struct {
	ID        int64         `json:"id"`
	LeadName  string        `json:"lead_name"`
	AgentName string        `json:"agent_name"`
	Status    string        `json:"status"`
	Duration  time.Duration `json:"duration"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Outcome   string        `json:"outcome"`
	CreatedAt time.Time     `json:"created_at"`
}

// ============================================================================
// Statistics Response
// ============================================================================

// StatisticsResponse represents statistical data
type StatisticsResponse struct {
	Period     string                 `json:"period"`
	Metrics    map[string]interface{} `json:"metrics"`
	Comparison map[string]interface{} `json:"comparison"`
	Trends     map[string]interface{} `json:"trends"`
	Timestamp  time.Time              `json:"timestamp"`
}
