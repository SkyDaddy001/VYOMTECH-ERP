package routes

import (
	"net/http"
	"strconv"

	"vyom-erp/cmd/api/handlers"
	"vyom-erp/cmd/api/middleware"
	"vyom-erp/cmd/api/services"

	"github.com/gin-gonic/gin"
)

// RegisterSalesRoutes registers all sales and CRM routes
func RegisterSalesRoutes(r *gin.Engine, services map[string]interface{}) {
	leadService := services["leadService"].(*services.LeadService)
	oppService := services["opportunityService"].(*services.OpportunityService)
	activityService := services["activityService"].(*services.ActivityService)
	accountService := services["accountService"].(*services.AccountService)
	contactService := services["contactService"].(*services.ContactService)
	interactionService := services["interactionService"].(*services.InteractionService)

	// Leads Group
	leads := r.Group("/api/leads")
	leads.Use(middleware.AuthMiddleware())
	{
		leads.POST("", handlers.CreateLeadHandler(leadService))
		leads.GET("/:id", handlers.GetLeadHandler(leadService))
		leads.GET("", handlers.ListLeadsHandler(leadService))
		leads.PATCH("/:id/status", handlers.UpdateLeadStatusHandler(leadService))
	}

	// Opportunities Group
	opportunities := r.Group("/api/opportunities")
	opportunities.Use(middleware.AuthMiddleware())
	{
		opportunities.POST("", handlers.CreateOpportunityHandler(oppService))
		opportunities.GET("/:id", handlers.GetOpportunityHandler(oppService))
		opportunities.GET("", handlers.ListOpportunitiesHandler(oppService))
		opportunities.PATCH("/:id/stage", handlers.UpdateOpportunityStageHandler(oppService))
		opportunities.GET("/dashboard/pipeline", handlers.GetPipelineHandler(oppService))
	}

	// Activities Group
	activities := r.Group("/api/activities")
	activities.Use(middleware.AuthMiddleware())
	{
		activities.POST("", handlers.CreateActivityHandler(activityService))
		activities.GET("/:id", handlers.GetActivityHandler(activityService))
		activities.GET("/lead/:leadId", handlers.GetLeadActivitiesHandler(activityService))
	}

	// CRM Accounts Group
	accounts := r.Group("/api/crm/accounts")
	accounts.Use(middleware.AuthMiddleware())
	{
		accounts.POST("", handlers.CreateAccountHandler(accountService))
		accounts.GET("/:id", handlers.GetAccountHandler(accountService))
		accounts.GET("", handlers.ListAccountsHandler(accountService))
	}

	// Contacts Group
	contacts := r.Group("/api/crm/contacts")
	contacts.Use(middleware.AuthMiddleware())
	{
		contacts.POST("", handlers.CreateContactHandler(contactService))
		contacts.GET("/account/:accountId", handlers.GetAccountContactsHandler(contactService))
	}

	// Interactions Group
	interactions := r.Group("/api/crm/interactions")
	interactions.Use(middleware.AuthMiddleware())
	{
		interactions.POST("", handlers.CreateInteractionHandler(interactionService))
		interactions.GET("/account/:accountId", handlers.GetAccountInteractionsHandler(interactionService))
	}

	// Sales Dashboard
	dashboard := r.Group("/api/dashboard/sales")
	dashboard.Use(middleware.AuthMiddleware())
	{
		dashboard.GET("", handlers.GetSalesDashboardHandler(services["salesAnalyticsService"].(*services.SalesAnalyticsService)))
	}
}

// RegisterFinanceRoutes registers all finance and accounting routes
func RegisterFinanceRoutes(r *gin.Engine, services map[string]interface{}) {
	coaService := services["chartOfAccountsService"].(*services.ChartOfAccountsService)
	jeService := services["journalEntryService"].(*services.JournalEntryService)
	glService := services["generalLedgerService"].(*services.GeneralLedgerService)
	tbService := services["trialBalanceService"].(*services.TrialBalanceService)
	arService := services["accountsReceivableService"].(*services.AccountsReceivableService)
	apService := services["accountsPayableService"].(*services.AccountsPayableService)
	fsService := services["financialStatementsService"].(*services.FinancialStatementsService)

	// Chart of Accounts
	coa := r.Group("/api/accounting/chart-of-accounts")
	coa.Use(middleware.AuthMiddleware())
	{
		coa.POST("", handlers.CreateChartOfAccountsHandler(coaService))
		coa.GET("/:id", handlers.GetChartOfAccountsHandler(coaService))
		coa.GET("", handlers.ListChartOfAccountsHandler(coaService))
	}

	// Journal Entries
	je := r.Group("/api/accounting/journal-entries")
	je.Use(middleware.AuthMiddleware())
	{
		je.POST("", handlers.CreateJournalEntryHandler(jeService))
		je.GET("/:id", handlers.GetJournalEntryHandler(jeService))
		je.GET("", handlers.ListJournalEntriesHandler(jeService))
		je.POST("/:id/post", handlers.PostJournalEntryHandler(jeService))
	}

	// General Ledger
	gl := r.Group("/api/accounting/general-ledger")
	gl.Use(middleware.AuthMiddleware())
	{
		gl.GET("/account/:accountId", handlers.GetAccountLedgerHandler(glService))
	}

	// Trial Balance
	tb := r.Group("/api/accounting/trial-balance")
	tb.Use(middleware.AuthMiddleware())
	{
		tb.GET("", handlers.GetTrialBalanceHandler(tbService))
	}

	// Accounts Receivable
	ar := r.Group("/api/accounting/receivable")
	ar.Use(middleware.AuthMiddleware())
	{
		ar.POST("", handlers.CreateARHandler(arService))
		ar.GET("", handlers.GetOutstandingARHandler(arService))
	}

	// Accounts Payable
	ap := r.Group("/api/accounting/payable")
	ap.Use(middleware.AuthMiddleware())
	{
		ap.POST("", handlers.CreateAPHandler(apService))
		ap.GET("", handlers.GetOutstandingAPHandler(apService))
	}

	// Financial Statements
	fs := r.Group("/api/accounting/statements")
	fs.Use(middleware.AuthMiddleware())
	{
		fs.GET("/income", handlers.GetIncomeStatementHandler(fsService))
		fs.GET("/balance-sheet", handlers.GetBalanceSheetHandler(fsService))
	}

	// Financial Dashboard
	fd := r.Group("/api/dashboard/finance")
	fd.Use(middleware.AuthMiddleware())
	{
		fd.GET("", handlers.GetFinancialDashboardHandler(arService, apService, coaService))
	}
}

// ===== HANDLER STUBS (implementations in separate handlers file) =====

// Lead Handlers
func (h *Handler) CreateLeadHandler(service *services.LeadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create lead handler"})
	}
}

func (h *Handler) GetLeadHandler(service *services.LeadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get lead handler"})
	}
}

func (h *Handler) ListLeadsHandler(service *services.LeadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		c.JSON(http.StatusOK, gin.H{"limit": limit, "offset": offset})
	}
}

func (h *Handler) UpdateLeadStatusHandler(service *services.LeadService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Update lead status handler"})
	}
}

// Opportunity Handlers
func (h *Handler) CreateOpportunityHandler(service *services.OpportunityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create opportunity handler"})
	}
}

func (h *Handler) GetOpportunityHandler(service *services.OpportunityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get opportunity handler"})
	}
}

func (h *Handler) ListOpportunitiesHandler(service *services.OpportunityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "List opportunities handler"})
	}
}

func (h *Handler) UpdateOpportunityStageHandler(service *services.OpportunityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Update opportunity stage handler"})
	}
}

func (h *Handler) GetPipelineHandler(service *services.OpportunityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get pipeline handler"})
	}
}

// Activity Handlers
func (h *Handler) CreateActivityHandler(service *services.ActivityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create activity handler"})
	}
}

func (h *Handler) GetActivityHandler(service *services.ActivityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get activity handler"})
	}
}

func (h *Handler) GetLeadActivitiesHandler(service *services.ActivityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get lead activities handler"})
	}
}

// Account Handlers
func (h *Handler) CreateAccountHandler(service *services.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create account handler"})
	}
}

func (h *Handler) GetAccountHandler(service *services.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get account handler"})
	}
}

func (h *Handler) ListAccountsHandler(service *services.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "List accounts handler"})
	}
}

// Contact Handlers
func (h *Handler) CreateContactHandler(service *services.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create contact handler"})
	}
}

func (h *Handler) GetAccountContactsHandler(service *services.ContactService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get account contacts handler"})
	}
}

// Interaction Handlers
func (h *Handler) CreateInteractionHandler(service *services.InteractionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create interaction handler"})
	}
}

func (h *Handler) GetAccountInteractionsHandler(service *services.InteractionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get account interactions handler"})
	}
}

// Dashboard Handlers
func (h *Handler) GetSalesDashboardHandler(service *services.SalesAnalyticsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get sales dashboard handler"})
	}
}

// Finance Handlers
func (h *Handler) CreateChartOfAccountsHandler(service *services.ChartOfAccountsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create chart of accounts handler"})
	}
}

func (h *Handler) GetChartOfAccountsHandler(service *services.ChartOfAccountsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get chart of accounts handler"})
	}
}

func (h *Handler) ListChartOfAccountsHandler(service *services.ChartOfAccountsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "List chart of accounts handler"})
	}
}

func (h *Handler) CreateJournalEntryHandler(service *services.JournalEntryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create journal entry handler"})
	}
}

func (h *Handler) GetJournalEntryHandler(service *services.JournalEntryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get journal entry handler"})
	}
}

func (h *Handler) ListJournalEntriesHandler(service *services.JournalEntryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "List journal entries handler"})
	}
}

func (h *Handler) PostJournalEntryHandler(service *services.JournalEntryService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Post journal entry handler"})
	}
}

func (h *Handler) GetAccountLedgerHandler(service *services.GeneralLedgerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get account ledger handler"})
	}
}

func (h *Handler) GetTrialBalanceHandler(service *services.TrialBalanceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get trial balance handler"})
	}
}

func (h *Handler) CreateARHandler(service *services.AccountsReceivableService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create AR handler"})
	}
}

func (h *Handler) GetOutstandingARHandler(service *services.AccountsReceivableService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get outstanding AR handler"})
	}
}

func (h *Handler) CreateAPHandler(service *services.AccountsPayableService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Create AP handler"})
	}
}

func (h *Handler) GetOutstandingAPHandler(service *services.AccountsPayableService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get outstanding AP handler"})
	}
}

func (h *Handler) GetIncomeStatementHandler(service *services.FinancialStatementsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get income statement handler"})
	}
}

func (h *Handler) GetBalanceSheetHandler(service *services.FinancialStatementsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get balance sheet handler"})
	}
}

func (h *Handler) GetFinancialDashboardHandler(arService *services.AccountsReceivableService, apService *services.AccountsPayableService, coaService *services.ChartOfAccountsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Get financial dashboard handler"})
	}
}

type Handler struct{}
