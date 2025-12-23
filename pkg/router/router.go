package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"vyomtech-backend/internal/handlers"
	"vyomtech-backend/internal/middleware"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/logger"
)

// SetupRoutesWithServices configures all API routes with basic services
func SetupRoutesWithServices(
	authService *services.AuthService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	log *logger.Logger,
) *mux.Router {
	return setupRoutes(authService, nil, passwordResetHandler, agentService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, log)
}

// SetupRoutesWithTenant configures all API routes with tenant service
func SetupRoutesWithTenant(
	authService *services.AuthService,
	tenantService *services.TenantService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	log *logger.Logger,
) *mux.Router {
	return setupRoutes(authService, tenantService, passwordResetHandler, agentService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, log)
}

// SetupRoutesWithGamification configures all API routes with gamification features
func SetupRoutesWithGamification(
	authService *services.AuthService,
	tenantService *services.TenantService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	gamificationService *services.GamificationService,
	log *logger.Logger,
) *mux.Router {
	return setupRoutes(authService, tenantService, passwordResetHandler, agentService, gamificationService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, log)
}

// SetupRoutesWithCoreFeatures configures all API routes with core features and real-time support
func SetupRoutesWithCoreFeatures(
	authService *services.AuthService,
	tenantService *services.TenantService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	gamificationService *services.GamificationService,
	leadService *services.LeadService,
	callService *services.CallService,
	campaignService *services.CampaignService,
	aiOrchestrator *services.AIOrchestrator,
	log *logger.Logger,
) *mux.Router {
	return setupRoutes(authService, tenantService, passwordResetHandler, agentService, gamificationService, leadService, callService, campaignService, aiOrchestrator, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, log)
}

// SetupRoutesWithRealtime configures all API routes with realtime WebSocket support
func SetupRoutesWithRealtime(
	authService *services.AuthService,
	tenantService *services.TenantService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	gamificationService *services.GamificationService,
	leadService *services.LeadService,
	callService *services.CallService,
	campaignService *services.CampaignService,
	aiOrchestrator *services.AIOrchestrator,
	webSocketHub *services.WebSocketHub,
	leadScoringService *services.LeadScoringService,
	dashboardService *services.DashboardService,
	taskService services.TaskService,
	notificationService services.NotificationService,
	customizationService services.TenantCustomizationService,
	log *logger.Logger,
) *mux.Router {
	return setupRoutes(authService, tenantService, passwordResetHandler, agentService, gamificationService, leadService, callService, campaignService, aiOrchestrator, webSocketHub, leadScoringService, dashboardService, taskService, notificationService, customizationService, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, log)
}

// SetupRoutesWithPhase3C configures all API routes including Phase 3C (Modular Monetization) and Phase 3.2 Mobile
func SetupRoutesWithPhase3C(
	authService *services.AuthService,
	tenantService *services.TenantService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	gamificationService *services.GamificationService,
	leadService *services.LeadService,
	callService *services.CallService,
	campaignService *services.CampaignService,
	aiOrchestrator *services.AIOrchestrator,
	webSocketHub *services.WebSocketHub,
	leadScoringService *services.LeadScoringService,
	dashboardService *services.DashboardService,
	taskService services.TaskService,
	notificationService services.NotificationService,
	customizationService services.TenantCustomizationService,
	phase3cServices *services.Phase3CServices,
	salesService *services.SalesService,
	realEstateService *services.RealEstateService,
	civilService *services.CivilService,
	constructionService *services.ConstructionService,
	boqService *services.BOQService,
	hrService *services.HRService,
	glService *services.GLService,
	rbacService *services.RBACService,
	reraComplianceHandler *handlers.RERAComplianceHandler,
	hrComplianceHandler *handlers.HRComplianceHandler,
	taxComplianceHandler *handlers.TaxComplianceHandler,
	financialDashboardHandler *handlers.FinancialDashboardHandler,
	hrDashboardHandler *handlers.HRDashboardHandler,
	complianceDashboardHandler *handlers.ComplianceDashboardHandler,
	salesDashboardHandler *handlers.SalesDashboardHandler,
	brokerHandler *handlers.BrokerHandler,
	jointApplicantHandler *handlers.JointApplicantHandler,
	documentHandler *handlers.DocumentHandler,
	possessionHandler *handlers.PossessionHandler,
	titleHandler *handlers.TitleHandler,
	customerPortalHandler *handlers.CustomerPortalHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	userAdminHandler *handlers.UserAdminHandler,
	tenantAdminHandler *handlers.TenantAdminHandler,
	mobileHandler *handlers.MobileHandler,
	aiRecommendationsHandler *handlers.AIRecommendationsHandler,
	log *logger.Logger,
) *mux.Router {
	return setupRoutes(authService, tenantService, passwordResetHandler, agentService, gamificationService, leadService, callService, campaignService, aiOrchestrator, webSocketHub, leadScoringService, dashboardService, taskService, notificationService, customizationService, phase3cServices, salesService, realEstateService, civilService, constructionService, boqService, hrService, glService, rbacService, reraComplianceHandler, hrComplianceHandler, taxComplianceHandler, financialDashboardHandler, hrDashboardHandler, complianceDashboardHandler, salesDashboardHandler, brokerHandler, jointApplicantHandler, documentHandler, possessionHandler, titleHandler, customerPortalHandler, analyticsHandler, userAdminHandler, tenantAdminHandler, mobileHandler, aiRecommendationsHandler, log)
}

func setupRoutes(
	authService *services.AuthService,
	tenantService *services.TenantService,
	passwordResetHandler *handlers.PasswordResetHandler,
	agentService *services.AgentService,
	gamificationService *services.GamificationService,
	leadService *services.LeadService,
	callService *services.CallService,
	campaignService *services.CampaignService,
	aiOrchestrator *services.AIOrchestrator,
	webSocketHub *services.WebSocketHub,
	leadScoringService *services.LeadScoringService,
	dashboardService *services.DashboardService,
	taskService services.TaskService,
	notificationService services.NotificationService,
	customizationService services.TenantCustomizationService,
	phase3cServices *services.Phase3CServices,
	salesService *services.SalesService,
	realEstateService *services.RealEstateService,
	civilService *services.CivilService,
	constructionService *services.ConstructionService,
	boqService *services.BOQService,
	hrService *services.HRService,
	glService *services.GLService,
	rbacService *services.RBACService,
	reraComplianceHandler *handlers.RERAComplianceHandler,
	hrComplianceHandler *handlers.HRComplianceHandler,
	taxComplianceHandler *handlers.TaxComplianceHandler,
	financialDashboardHandler *handlers.FinancialDashboardHandler,
	hrDashboardHandler *handlers.HRDashboardHandler,
	complianceDashboardHandler *handlers.ComplianceDashboardHandler,
	salesDashboardHandler *handlers.SalesDashboardHandler,
	brokerHandler *handlers.BrokerHandler,
	jointApplicantHandler *handlers.JointApplicantHandler,
	documentHandler *handlers.DocumentHandler,
	possessionHandler *handlers.PossessionHandler,
	titleHandler *handlers.TitleHandler,
	customerPortalHandler *handlers.CustomerPortalHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	userAdminHandler *handlers.UserAdminHandler,
	tenantAdminHandler *handlers.TenantAdminHandler,
	mobileHandler *handlers.MobileHandler,
	aiRecommendationsHandler *handlers.AIRecommendationsHandler,
	log *logger.Logger,
) *mux.Router {
	r := mux.NewRouter()

	// Global middleware
	r.Use(middleware.RequestLoggingMiddleware(log))
	r.Use(middleware.ErrorRecoveryMiddleware(log))
	r.Use(middleware.CORSMiddleware())

	// Health check endpoints (no auth required)
	r.HandleFunc("/health", HealthCheck).Methods("GET")
	r.HandleFunc("/ready", ReadinessCheck).Methods("GET")

	// API v1 routes
	v1 := r.PathPrefix("/api/v1").Subrouter()

	// Authentication routes (no auth required)
	authHandler := handlers.NewAuthHandler(authService, log)
	authRoutes := v1.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Protected authentication routes
	protectedAuth := v1.PathPrefix("/auth").Subrouter()
	protectedAuth.Use(middleware.AuthMiddleware(authService, log))
	protectedAuth.HandleFunc("/validate", authHandler.ValidateToken).Methods("GET")
	protectedAuth.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")

	// Password reset routes
	resetRoutes := v1.PathPrefix("/password-reset").Subrouter()
	resetRoutes.HandleFunc("/request", passwordResetHandler.RequestReset).Methods("POST")
	resetRoutes.HandleFunc("/reset", passwordResetHandler.ResetPassword).Methods("POST")

	// Admin Tenant Routes (CRUD operations for tenant management - replaces old tenant routes)
	if tenantAdminHandler != nil {
		tenantAdminRoutes := v1.PathPrefix("/tenants").Subrouter()
		tenantAdminRoutes.HandleFunc("", tenantAdminHandler.ListTenants).Methods("GET")
		tenantAdminRoutes.HandleFunc("", tenantAdminHandler.CreateTenant).Methods("POST")
		tenantAdminRoutes.HandleFunc("/{id}", tenantAdminHandler.GetTenant).Methods("GET")
		tenantAdminRoutes.HandleFunc("/{id}", tenantAdminHandler.UpdateTenant).Methods("PUT")
		tenantAdminRoutes.HandleFunc("/{id}", tenantAdminHandler.DeleteTenant).Methods("DELETE")
		tenantAdminRoutes.HandleFunc("/{id}/users", tenantAdminHandler.GetTenantUsers).Methods("GET")
	}

	// Legacy tenant routes (if tenant service provided) - for backwards compatibility
	if tenantService != nil {
		tenantHandler := handlers.NewTenantHandler(tenantService, log)

		// Protected tenant routes (different path prefix to avoid conflicts)
		protectedTenantRoutes := v1.PathPrefix("/tenant").Subrouter()
		protectedTenantRoutes.Use(middleware.AuthMiddleware(authService, log))
		protectedTenantRoutes.HandleFunc("", tenantHandler.GetTenantInfo).Methods("GET")
		protectedTenantRoutes.HandleFunc("/users/count", tenantHandler.GetTenantUserCount).Methods("GET")

		// Multi-tenant routes (protected)
		multiTenantRoutes := v1.PathPrefix("/my-tenants").Subrouter()
		multiTenantRoutes.Use(middleware.AuthMiddleware(authService, log))
		multiTenantRoutes.HandleFunc("", tenantHandler.GetUserTenants).Methods("GET")
		multiTenantRoutes.HandleFunc("/{id}/switch", tenantHandler.SwitchTenant).Methods("POST")
		multiTenantRoutes.HandleFunc("/{id}/members", tenantHandler.AddTenantMember).Methods("POST")
		multiTenantRoutes.HandleFunc("/{id}/members/{email}", tenantHandler.RemoveTenantMember).Methods("DELETE")
	}

	// Admin User Routes (CRUD operations for user management)
	if userAdminHandler != nil {
		userAdminRoutes := v1.PathPrefix("/users").Subrouter()
		userAdminRoutes.Use(middleware.AuthMiddleware(authService, log))
		userAdminRoutes.Use(middleware.TenantIsolationMiddleware(log))
		userAdminRoutes.HandleFunc("", userAdminHandler.ListUsers).Methods("GET")
		userAdminRoutes.HandleFunc("", userAdminHandler.CreateUser).Methods("POST")
		userAdminRoutes.HandleFunc("/{id}", userAdminHandler.GetUser).Methods("GET")
		userAdminRoutes.HandleFunc("/{id}", userAdminHandler.UpdateUser).Methods("PUT")
		userAdminRoutes.HandleFunc("/{id}", userAdminHandler.DeleteUser).Methods("DELETE")
		userAdminRoutes.HandleFunc("/{id}/role", userAdminHandler.UpdateUserRole).Methods("PUT")
		userAdminRoutes.HandleFunc("/{id}/reset-password", userAdminHandler.ResetPassword).Methods("POST")
	}

	// Protected agent routes
	agentHandler := handlers.NewAgentHandler(agentService, log)
	agentRoutes := v1.PathPrefix("/agents").Subrouter()
	agentRoutes.Use(middleware.AuthMiddleware(authService, log))
	agentRoutes.Use(middleware.TenantIsolationMiddleware(log))

	agentRoutes.HandleFunc("/{id}", agentHandler.GetAgent).Methods("GET")
	agentRoutes.HandleFunc("", agentHandler.GetAgentsByTenant).Methods("GET")
	agentRoutes.HandleFunc("/status", agentHandler.UpdateAgentAvailability).Methods("PATCH")
	agentRoutes.HandleFunc("/available", agentHandler.GetAvailableAgents).Methods("GET")
	agentRoutes.HandleFunc("/stats", agentHandler.GetAgentStats).Methods("GET")

	// Gamification routes (protected)
	if gamificationService != nil {
		gamificationHandler := handlers.NewGamificationHandler(gamificationService, log)
		gamificationRoutes := v1.PathPrefix("/gamification").Subrouter()
		gamificationRoutes.Use(middleware.AuthMiddleware(authService, log))
		gamificationRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Points endpoints
		gamificationRoutes.HandleFunc("/points", gamificationHandler.GetUserPoints).Methods("GET")
		gamificationRoutes.HandleFunc("/points/award", gamificationHandler.AwardPoints).Methods("POST")
		gamificationRoutes.HandleFunc("/points/revoke", gamificationHandler.RevokePoints).Methods("POST")

		// Badges endpoints
		gamificationRoutes.HandleFunc("/badges", gamificationHandler.GetUserBadges).Methods("GET")
		gamificationRoutes.HandleFunc("/badges", gamificationHandler.CreateBadge).Methods("POST")
		gamificationRoutes.HandleFunc("/badges/award", gamificationHandler.AwardBadge).Methods("POST")

		// Challenges endpoints
		gamificationRoutes.HandleFunc("/challenges", gamificationHandler.GetUserChallenges).Methods("GET")
		gamificationRoutes.HandleFunc("/challenges/active", gamificationHandler.GetActiveChallenges).Methods("GET")
		gamificationRoutes.HandleFunc("/challenges", gamificationHandler.CreateChallenge).Methods("POST")

		// Leaderboard endpoints
		gamificationRoutes.HandleFunc("/leaderboard", gamificationHandler.GetLeaderboard).Methods("GET")

		// Profile endpoints
		gamificationRoutes.HandleFunc("/profile", gamificationHandler.GetGamificationProfile).Methods("GET")
	}

	// Lead routes (protected)
	if leadService != nil {
		leadHandler := handlers.NewLeadHandler(leadService, log)
		leadRoutes := v1.PathPrefix("/leads").Subrouter()
		leadRoutes.Use(middleware.AuthMiddleware(authService, log))
		leadRoutes.Use(middleware.TenantIsolationMiddleware(log))
		leadRoutes.HandleFunc("", leadHandler.GetLeads).Methods("GET")
		leadRoutes.HandleFunc("/stats", leadHandler.GetLeadStats).Methods("GET")
		leadRoutes.HandleFunc("", leadHandler.CreateLead).Methods("POST")
		// Individual lead operations
		leadRoutes.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "id required", http.StatusBadRequest)
				return
			}
			switch r.Method {
			case "GET":
				leadHandler.GetLead(w, r)
			case "PUT":
				leadHandler.UpdateLead(w, r)
			case "DELETE":
				leadHandler.DeleteLead(w, r)
			}
		}).Methods("GET", "PUT", "DELETE")
	}

	// Call routes (protected)
	if callService != nil {
		callHandler := handlers.NewCallHandler(callService, log)
		callRoutes := v1.PathPrefix("/calls").Subrouter()
		callRoutes.Use(middleware.AuthMiddleware(authService, log))
		callRoutes.Use(middleware.TenantIsolationMiddleware(log))
		callRoutes.HandleFunc("", callHandler.GetCalls).Methods("GET")
		callRoutes.HandleFunc("/stats", callHandler.GetCallStats).Methods("GET")
		callRoutes.HandleFunc("", callHandler.CreateCall).Methods("POST")
		// Individual call operations
		callRoutes.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "id required", http.StatusBadRequest)
				return
			}
			if r.Method == "GET" {
				callHandler.GetCall(w, r)
			} else if r.Method == "POST" && r.URL.Path == "/end" {
				callHandler.EndCall(w, r)
			}
		}).Methods("GET", "POST")
	}

	// AI routes (protected)
	if aiOrchestrator != nil {
		aiHandler := handlers.NewAIHandler(aiOrchestrator, log)
		aiRoutes := v1.PathPrefix("/ai").Subrouter()
		aiRoutes.Use(middleware.AuthMiddleware(authService, log))
		aiRoutes.HandleFunc("/query", aiHandler.ProcessAIQuery).Methods("POST")
		aiRoutes.HandleFunc("/providers", aiHandler.ListAIProviders).Methods("GET")
	}

	// Campaign routes (protected)
	if campaignService != nil {
		campaignHandler := handlers.NewCampaignHandler(campaignService, log)
		campaignRoutes := v1.PathPrefix("/campaigns").Subrouter()
		campaignRoutes.Use(middleware.AuthMiddleware(authService, log))
		campaignRoutes.Use(middleware.TenantIsolationMiddleware(log))
		campaignRoutes.HandleFunc("", campaignHandler.GetCampaigns).Methods("GET")
		campaignRoutes.HandleFunc("/stats", campaignHandler.GetCampaignStats).Methods("GET")
		campaignRoutes.HandleFunc("", campaignHandler.CreateCampaign).Methods("POST")
		// Individual campaign operations
		campaignRoutes.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Query().Get("id")
			if id == "" {
				http.Error(w, "id required", http.StatusBadRequest)
				return
			}
			switch r.Method {
			case "GET":
				campaignHandler.GetCampaign(w, r)
			case "PUT":
				campaignHandler.UpdateCampaign(w, r)
			case "DELETE":
				campaignHandler.DeleteCampaign(w, r)
			}
		}).Methods("GET", "PUT", "DELETE")
	}

	// WebSocket routes (if hub provided)
	if webSocketHub != nil {
		wsHandler := handlers.NewWebSocketHandler(webSocketHub, log)
		wsRoutes := v1.PathPrefix("/ws").Subrouter()
		wsRoutes.Use(middleware.AuthMiddleware(authService, log))
		wsRoutes.Use(middleware.TenantIsolationMiddleware(log))
		wsRoutes.HandleFunc("", wsHandler.HandleConnection).Methods("GET")
		wsRoutes.HandleFunc("/stats", wsHandler.GetConnectionStats).Methods("GET")
	}

	// Advanced Gamification routes (protected)
	// NOTE: These routes require the AdvancedGamificationService to be initialized and passed through router setup
	// Currently commented out - will be enabled when service is properly integrated
	// advGamService := services.NewAdvancedGamificationService(dbConn)
	// advGamHandler := handlers.NewAdvancedGamificationHandler(advGamService, log)
	// advGamRoutes := v1.PathPrefix("/gamification-advanced").Subrouter()
	// advGamRoutes.Use(middleware.AuthMiddleware(authService, log))
	// advGamRoutes.Use(middleware.TenantIsolationMiddleware(log))
	// advGamRoutes.HandleFunc("/competitions", advGamHandler.CreateTeamCompetition).Methods("POST")
	// advGamRoutes.HandleFunc("/competitions/leaderboard", advGamHandler.GetTeamLeaderboard).Methods("GET")
	// advGamRoutes.HandleFunc("/challenges", advGamHandler.CreateChallenge).Methods("POST")
	// advGamRoutes.HandleFunc("/challenges/active", advGamHandler.GetActiveChallenges).Methods("GET")
	// advGamRoutes.HandleFunc("/rewards", advGamHandler.GetAvailableRewards).Methods("GET")
	// advGamRoutes.HandleFunc("/rewards", advGamHandler.CreateReward).Methods("POST")
	// advGamRoutes.HandleFunc("/redeem", advGamHandler.RedeemReward).Methods("POST")
	// advGamRoutes.HandleFunc("/leaderboard", advGamHandler.GetGamificationLeaderboard).Methods("GET")
	// advGamRoutes.HandleFunc("/stats", advGamHandler.GetGamificationStats).Methods("GET")

	// Analytics routes (protected)
	// TODO: Pass database connection through router setup
	// analyticsService := services.NewAnalyticsService(dbConn)
	// analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, log)
	// analyticsRoutes := v1.PathPrefix("/analytics").Subrouter()
	// analyticsRoutes.Use(middleware.AuthMiddleware(authService, log))
	// analyticsRoutes.Use(middleware.TenantIsolationMiddleware(log))
	// analyticsRoutes.HandleFunc("/reports", analyticsHandler.GenerateReport).Methods("POST")
	// analyticsRoutes.HandleFunc("/export", analyticsHandler.ExportReport).Methods("POST")
	// analyticsRoutes.HandleFunc("/trends", analyticsHandler.GetTrends).Methods("GET")
	// analyticsRoutes.HandleFunc("/metrics", analyticsHandler.GetCustomMetrics).Methods("GET")

	// Automation routes (protected) - lead scoring, routing, workflows
	// TODO: Pass database connection through router setup
	// automationService := services.NewAutomationService(dbConn)
	// automationHandler := handlers.NewAutomationHandler(automationService, log)
	// automationRoutes := v1.PathPrefix("/automation").Subrouter()
	// automationRoutes.Use(middleware.AuthMiddleware(authService, log))
	// automationRoutes.Use(middleware.TenantIsolationMiddleware(log))
	// automationRoutes.HandleFunc("/leads/score", automationHandler.CalculateLeadScore).Methods("POST")
	// automationRoutes.HandleFunc("/leads/ranked", automationHandler.RankLeads).Methods("GET")
	// automationRoutes.HandleFunc("/leads/route", automationHandler.RouteLeadToAgent).Methods("POST")
	// automationRoutes.HandleFunc("/routing-rules", automationHandler.CreateRoutingRule).Methods("POST")
	// automationRoutes.HandleFunc("/schedule-campaign", automationHandler.ScheduleCampaign).Methods("POST")
	// automationRoutes.HandleFunc("/metrics", automationHandler.GetLeadScoringMetrics).Methods("GET")

	// Communication routes (protected) - email, SMS, WhatsApp, etc
	// TODO: Pass database connection through router setup
	// commService := services.NewCommunicationService(dbConn)
	// commHandler := handlers.NewCommunicationHandler(commService, log)
	// commRoutes := v1.PathPrefix("/communication").Subrouter()
	// commRoutes.Use(middleware.AuthMiddleware(authService, log))
	// commRoutes.Use(middleware.TenantIsolationMiddleware(log))
	// commRoutes.HandleFunc("/providers", commHandler.RegisterProvider).Methods("POST")
	// commRoutes.HandleFunc("/templates", commHandler.CreateTemplate).Methods("POST")
	// commRoutes.HandleFunc("/messages", commHandler.SendMessage).Methods("POST")
	// commRoutes.HandleFunc("/messages/status", commHandler.GetMessageStatus).Methods("GET")
	// commRoutes.HandleFunc("/stats", commHandler.GetMessageStats).Methods("GET")

	// Phase 1: Lead Scoring routes (protected)
	if leadScoringService != nil {
		_ = handlers.NewPhase1Handler(v1, leadScoringService, log)
		// Routes are registered in NewPhase1Handler
		// GET    /api/v1/leads/{id}/score
		// POST   /api/v1/leads/{id}/score/calculate
		// GET    /api/v1/leads/scores/category/{category}
		// POST   /api/v1/leads/scores/batch-calculate
	}

	// Dashboard routes (protected)
	if dashboardService != nil {
		_ = handlers.NewDashboardHandler(v1, dashboardService, log)
		// Routes are registered in NewDashboardHandler
		// GET    /api/v1/dashboard/analytics
		// GET    /api/v1/dashboard/activity-logs
		// GET    /api/v1/dashboard/users
		// GET    /api/v1/dashboard/usage
	}

	// Phase 2A: Task routes (protected)
	if taskService != nil {
		taskHandler := handlers.NewTaskHandler(taskService, log)
		taskRoutes := v1.PathPrefix("/tasks").Subrouter()
		taskRoutes.Use(middleware.AuthMiddleware(authService, log))
		taskRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Register all task routes
		taskHandler.RegisterRoutes(taskRoutes)
		// Endpoints include:
		// POST   /api/v1/tasks
		// GET    /api/v1/tasks
		// GET    /api/v1/tasks/{id}
		// PUT    /api/v1/tasks/{id}
		// DELETE /api/v1/tasks/{id}
		// POST   /api/v1/tasks/{id}/complete
		// GET    /api/v1/tasks/user/{userID}
		// GET    /api/v1/tasks/stats
	}

	// Phase 2A: Notification routes (protected)
	if notificationService != nil {
		notificationHandler := handlers.NewNotificationHandler(notificationService, log)
		notificationRoutes := v1.PathPrefix("/notifications").Subrouter()
		notificationRoutes.Use(middleware.AuthMiddleware(authService, log))
		notificationRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Register all notification routes
		notificationHandler.RegisterRoutes(notificationRoutes)
		// Endpoints include:
		// POST   /api/v1/notifications
		// GET    /api/v1/notifications
		// GET    /api/v1/notifications/{id}
		// DELETE /api/v1/notifications/{id}
		// POST   /api/v1/notifications/{id}/read
		// POST   /api/v1/notifications/{id}/archive
		// GET    /api/v1/notifications/user/{userID}/unread
		// GET    /api/v1/notifications/stats
		// GET    /api/v1/notifications/preferences
		// PUT    /api/v1/notifications/preferences
	}

	// Phase 2B: Tenant Customization routes (protected)
	if customizationService != nil {
		customizationHandler := handlers.NewCustomizationHandler(customizationService)
		customizationRoutes := v1.PathPrefix("/config").Subrouter()
		customizationRoutes.Use(middleware.AuthMiddleware(authService, log))
		customizationRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Register all customization routes
		customizationHandler.RegisterRoutes(customizationRoutes)
		// Endpoints include:
		// POST   /api/v1/config/task-statuses
		// GET    /api/v1/config/task-statuses
		// GET    /api/v1/config/task-statuses/{statusCode}
		// PUT    /api/v1/config/task-statuses/{statusCode}
		// DELETE /api/v1/config/task-statuses/{statusCode}
		// POST   /api/v1/config/task-stages
		// GET    /api/v1/config/task-stages
		// PUT    /api/v1/config/task-stages/{stageCode}
		// POST   /api/v1/config/status-transitions
		// GET    /api/v1/config/status-transitions
		// GET    /api/v1/config/status-transitions/check
		// POST   /api/v1/config/task-types
		// GET    /api/v1/config/task-types
		// PUT    /api/v1/config/task-types/{typeCode}
		// POST   /api/v1/config/priority-levels
		// GET    /api/v1/config/priority-levels
		// PUT    /api/v1/config/priority-levels/{priorityCode}
		// POST   /api/v1/config/notification-types
		// GET    /api/v1/config/notification-types
		// PUT    /api/v1/config/notification-types/{typeCode}
		// POST   /api/v1/config/custom-fields
		// GET    /api/v1/config/custom-fields
		// PUT    /api/v1/config/custom-fields/{fieldCode}
		// POST   /api/v1/config/automation-rules
		// GET    /api/v1/config/automation-rules
		// PUT    /api/v1/config/automation-rules/{ruleCode}
		// GET    /api/v1/config/all
	}

	// Phase 3C: Modular Monetization routes (protected)
	if phase3cServices != nil {
		moduleHandler := handlers.NewModuleHandler(phase3cServices.ModuleService, log)
		companyHandler := handlers.NewCompanyHandler(phase3cServices.CompanyService, log)
		billingHandler := handlers.NewBillingHandler(phase3cServices.BillingService, log, rbacService)

		// Module routes
		moduleRoutes := v1.PathPrefix("/modules").Subrouter()
		moduleRoutes.Use(middleware.AuthMiddleware(authService, log))
		moduleRoutes.Use(middleware.TenantIsolationMiddleware(log))
		moduleRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		moduleRoutes.HandleFunc("/register", moduleHandler.RegisterModule).Methods("POST")
		moduleRoutes.HandleFunc("", moduleHandler.ListModules).Methods("GET")
		moduleRoutes.HandleFunc("/subscribe", moduleHandler.SubscribeToModule).Methods("POST")
		moduleRoutes.HandleFunc("/toggle", moduleHandler.ToggleModule).Methods("PUT")
		moduleRoutes.HandleFunc("/usage", moduleHandler.GetModuleUsage).Methods("GET")
		moduleRoutes.HandleFunc("/subscriptions", moduleHandler.ListSubscriptions).Methods("GET")

		// Company routes
		companyRoutes := v1.PathPrefix("/companies").Subrouter()
		companyRoutes.Use(middleware.AuthMiddleware(authService, log))
		companyRoutes.Use(middleware.TenantIsolationMiddleware(log))
		companyRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		companyRoutes.HandleFunc("", companyHandler.CreateCompany).Methods("POST")
		companyRoutes.HandleFunc("", companyHandler.ListCompanies).Methods("GET")
		companyRoutes.HandleFunc("/{id}", companyHandler.GetCompany).Methods("GET")
		companyRoutes.HandleFunc("/{id}", companyHandler.UpdateCompany).Methods("PUT")
		companyRoutes.HandleFunc("/{companyId}/projects", companyHandler.CreateProject).Methods("POST")
		companyRoutes.HandleFunc("/{companyId}/projects", companyHandler.ListProjects).Methods("GET")
		companyRoutes.HandleFunc("/{companyId}/projects/{projectId}", companyHandler.GetProject).Methods("GET")
		companyRoutes.HandleFunc("/{companyId}/members", companyHandler.GetCompanyMembers).Methods("GET")
		companyRoutes.HandleFunc("/{companyId}/members", companyHandler.AddMemberToCompany).Methods("POST")
		companyRoutes.HandleFunc("/{companyId}/projects/{projectId}/members", companyHandler.AddMemberToProject).Methods("POST")
		companyRoutes.HandleFunc("/{companyId}/projects/{projectId}/members", companyHandler.GetProjectMembers).Methods("GET")
		companyRoutes.HandleFunc("/{companyId}/projects/{projectId}/members/{userId}", companyHandler.RemoveProjectMember).Methods("DELETE")

		// Billing routes
		billingRoutes := v1.PathPrefix("/billing").Subrouter()
		billingRoutes.Use(middleware.AuthMiddleware(authService, log))
		billingRoutes.Use(middleware.TenantIsolationMiddleware(log))
		billingRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		billingRoutes.HandleFunc("/plans", billingHandler.CreatePricingPlan).Methods("POST")
		billingRoutes.HandleFunc("/plans", billingHandler.ListPricingPlans).Methods("GET")
		billingRoutes.HandleFunc("/subscribe", billingHandler.SubscribeToPlan).Methods("POST")
		billingRoutes.HandleFunc("/usage", billingHandler.RecordUsageMetrics).Methods("POST")
		billingRoutes.HandleFunc("/usage", billingHandler.GetUsageMetrics).Methods("GET")
		billingRoutes.HandleFunc("/invoices", billingHandler.ListInvoices).Methods("GET")
		billingRoutes.HandleFunc("/invoices/{id}", billingHandler.GetInvoice).Methods("GET")
		billingRoutes.HandleFunc("/invoices/{id}/pay", billingHandler.MarkInvoiceAsPaid).Methods("PUT")
		billingRoutes.HandleFunc("/charges", billingHandler.CalculateMonthlyCharges).Methods("GET")

		// Endpoints include:
		// Module endpoints:
		// POST   /api/v1/modules/register
		// GET    /api/v1/modules
		// POST   /api/v1/modules/subscribe
		// PUT    /api/v1/modules/toggle
		// GET    /api/v1/modules/usage
		// GET    /api/v1/modules/subscriptions
		// Company endpoints:
		// POST   /api/v1/companies
		// GET    /api/v1/companies
		// GET    /api/v1/companies/{id}
		// PUT    /api/v1/companies/{id}
		// POST   /api/v1/companies/{companyId}/projects
		// GET    /api/v1/companies/{companyId}/projects
		// GET    /api/v1/companies/{companyId}/members
		// POST   /api/v1/companies/{companyId}/members
		// DELETE /api/v1/companies/{companyId}/members/{userId}
		// Billing endpoints:
		// POST   /api/v1/billing/plans
		// GET    /api/v1/billing/plans
		// POST   /api/v1/billing/subscribe
		// POST   /api/v1/billing/usage
		// GET    /api/v1/billing/invoices
		// GET    /api/v1/billing/charges
	}

	// Sales Module routes (protected)
	if salesService != nil {
		salesHandler := handlers.NewSalesHandler(salesService.DB, rbacService)
		salesRoutes := v1.PathPrefix("/sales").Subrouter()
		salesRoutes.Use(middleware.AuthMiddleware(authService, log))
		salesRoutes.Use(middleware.TenantIsolationMiddleware(log))
		salesRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "agent"},
			log,
		))

		// Lead endpoints
		salesRoutes.HandleFunc("/leads", salesHandler.ListSalesLeads).Methods("GET")
		salesRoutes.HandleFunc("/leads", salesHandler.CreateSalesLead).Methods("POST")
		salesRoutes.HandleFunc("/leads/{id}", salesHandler.GetSalesLead).Methods("GET")
		salesRoutes.HandleFunc("/leads/{id}", salesHandler.UpdateSalesLead).Methods("PUT")
		salesRoutes.HandleFunc("/leads/{id}", salesHandler.DeleteSalesLead).Methods("DELETE")

		// Customer endpoints
		salesRoutes.HandleFunc("/customers", salesHandler.ListSalesCustomers).Methods("GET")
		salesRoutes.HandleFunc("/customers", salesHandler.CreateSalesCustomer).Methods("POST")
		salesRoutes.HandleFunc("/customers/{id}", salesHandler.GetSalesCustomer).Methods("GET")
		salesRoutes.HandleFunc("/customers/{id}", salesHandler.UpdateSalesCustomer).Methods("PUT")

		// Quotation endpoints
		salesRoutes.HandleFunc("/quotations", salesHandler.ListSalesQuotations).Methods("GET")
		salesRoutes.HandleFunc("/quotations", salesHandler.CreateSalesQuotation).Methods("POST")
		salesRoutes.HandleFunc("/quotations/{id}", salesHandler.GetSalesQuotation).Methods("GET")

		// Sales Order endpoints
		salesRoutes.HandleFunc("/orders", salesHandler.ListSalesOrders).Methods("GET")
		salesRoutes.HandleFunc("/orders", salesHandler.CreateSalesOrder).Methods("POST")
		salesRoutes.HandleFunc("/orders/{id}", salesHandler.GetSalesOrder).Methods("GET")
		salesRoutes.HandleFunc("/orders/{id}/status", salesHandler.UpdateSalesOrderStatus).Methods("PUT")

		// Invoice endpoints
		salesRoutes.HandleFunc("/invoices", salesHandler.ListSalesInvoices).Methods("GET")
		salesRoutes.HandleFunc("/invoices", salesHandler.CreateSalesInvoice).Methods("POST")
		salesRoutes.HandleFunc("/invoices/{id}", salesHandler.GetSalesInvoice).Methods("GET")

		// Payment endpoints
		salesRoutes.HandleFunc("/payments", salesHandler.CreateSalesPayment).Methods("POST")

		// Delivery endpoints
		salesRoutes.HandleFunc("/delivery-notes", salesHandler.CreateDeliveryNote).Methods("POST")
		salesRoutes.HandleFunc("/delivery-notes/{id}/pod", salesHandler.UpdateDeliveryPOD).Methods("PUT")

		// Credit Note endpoints
		salesRoutes.HandleFunc("/credit-notes", salesHandler.CreateCreditNote).Methods("POST")

		// Metrics endpoints
		salesRoutes.HandleFunc("/metrics/{salesperson_id}", salesHandler.GetSalespersonMetrics).Methods("GET")
		salesRoutes.HandleFunc("/metrics/calculate", salesHandler.CalculateAndUpdateMetrics).Methods("POST")

		// Milestone & Tracking endpoints
		salesRoutes.HandleFunc("/milestones/lead", salesHandler.CreateLeadMilestone).Methods("POST")
		salesRoutes.HandleFunc("/milestones/lead/{lead_id}", salesHandler.GetLeadMilestones).Methods("GET")

		// Engagement endpoints
		salesRoutes.HandleFunc("/engagement", salesHandler.CreateLeadEngagement).Methods("POST")
		salesRoutes.HandleFunc("/engagement/{lead_id}", salesHandler.GetLeadEngagements).Methods("GET")

		// Booking endpoints
		salesRoutes.HandleFunc("/bookings", salesHandler.CreateBooking).Methods("POST")
		salesRoutes.HandleFunc("/bookings", salesHandler.GetBookings).Methods("GET")

		// Account Ledger endpoints
		salesRoutes.HandleFunc("/ledger", salesHandler.CreateLedgerEntry).Methods("POST")
		salesRoutes.HandleFunc("/ledger/{customer_id}", salesHandler.GetCustomerLedger).Methods("GET")

		// Campaign endpoints
		salesRoutes.HandleFunc("/campaigns", salesHandler.CreateCampaign).Methods("POST")
		salesRoutes.HandleFunc("/campaigns", salesHandler.GetCampaigns).Methods("GET")

		// Reporting & Analytics endpoints
		salesRoutes.HandleFunc("/reports/funnel", salesHandler.LeadFunnelAnalysis).Methods("GET")
		salesRoutes.HandleFunc("/reports/source-performance", salesHandler.LeadSourcePerformance).Methods("GET")
		salesRoutes.HandleFunc("/reports/bookings", salesHandler.BookingSummary).Methods("GET")
		salesRoutes.HandleFunc("/reports/customer-ledger/{customer_id}", salesHandler.CustomerLedgerSummary).Methods("GET")
		salesRoutes.HandleFunc("/reports/milestone-timeline/{lead_id}", salesHandler.MilestoneTimeline).Methods("GET")
		salesRoutes.HandleFunc("/reports/engagement-stats/{lead_id}", salesHandler.LeadEngagementStats).Methods("GET")
		salesRoutes.HandleFunc("/reports/dashboard", salesHandler.DashboardMetrics).Methods("GET")
	}

	// ============================================
	// RBAC (ROLE-BASED ACCESS CONTROL) ROUTES
	// ============================================
	if rbacService != nil {
		rbacHandler := handlers.NewRBACHandler(rbacService, rbacService.GetDB(), log)
		rbacRoutes := v1.PathPrefix("/rbac").Subrouter()
		rbacRoutes.Use(middleware.AuthMiddleware(authService, log))
		rbacRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// List roles and permissions (accessible to all authenticated users)
		rbacRoutes.HandleFunc("/roles", rbacHandler.ListRoles).Methods("GET")
		rbacRoutes.HandleFunc("/permissions", rbacHandler.ListPermissions).Methods("GET")
		rbacRoutes.HandleFunc("/roles/{id}", rbacHandler.GetRole).Methods("GET")

		// Admin-only operations
		adminRbacRoutes := rbacRoutes.PathPrefix("").Subrouter()
		adminRbacRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin"},
			log,
		))
		adminRbacRoutes.HandleFunc("/roles", rbacHandler.CreateRole).Methods("POST")
		adminRbacRoutes.HandleFunc("/roles/{id}/permissions", rbacHandler.AssignPermissions).Methods("PUT")
		adminRbacRoutes.HandleFunc("/roles/{id}", rbacHandler.DeleteRole).Methods("DELETE")

		// Phase 3.6: User Role Assignment & Membership
		adminRbacRoutes.HandleFunc("/users/{user_id}/roles", rbacHandler.AssignRoleToUser).Methods("POST")
		adminRbacRoutes.HandleFunc("/users/{user_id}/roles", rbacHandler.GetUserRoles).Methods("GET")
		adminRbacRoutes.HandleFunc("/users/{user_id}/roles/{role_id}", rbacHandler.RemoveRoleFromUser).Methods("DELETE")
		adminRbacRoutes.HandleFunc("/users/{user_id}/roles/{role_id}", rbacHandler.UpdateUserRole).Methods("PUT")
		adminRbacRoutes.HandleFunc("/roles/{role_id}/members", rbacHandler.GetRoleMembers).Methods("GET")

		// Phase 4: Advanced RBAC Features
		adminRbacRoutes.HandleFunc("/resource-access", rbacHandler.GrantResourceAccess).Methods("POST")
		adminRbacRoutes.HandleFunc("/time-based-permissions", rbacHandler.CreateTimeBasedPermission).Methods("POST")
		adminRbacRoutes.HandleFunc("/field-permissions", rbacHandler.SetFieldLevelPermission).Methods("POST")
		adminRbacRoutes.HandleFunc("/delegations", rbacHandler.DelegateRole).Methods("POST")
		adminRbacRoutes.HandleFunc("/bulk-assign", rbacHandler.BulkAssignPermissions).Methods("POST")
	}

	// ============================================
	// CIVIL ENGINEERING ROUTES
	// ============================================
	if civilService != nil {
		civilRoutes := v1.PathPrefix("/civil").Subrouter()
		civilRoutes.Use(middleware.AuthMiddleware(authService, log))
		civilRoutes.Use(middleware.TenantIsolationMiddleware(log))
		civilRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterCivilRoutes(civilRoutes, civilService, rbacService)
	}

	// ============================================
	// CONSTRUCTION MANAGEMENT ROUTES
	// ============================================

	if constructionService != nil {
		constructionRoutes := v1.PathPrefix("/construction").Subrouter()
		constructionRoutes.Use(middleware.AuthMiddleware(authService, log))
		constructionRoutes.Use(middleware.TenantIsolationMiddleware(log))
		constructionRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterConstructionRoutes(constructionRoutes, constructionService, rbacService)
	}

	// ============================================
	// BOQ MANAGEMENT ROUTES
	// ============================================
	if boqService != nil {
		boqRoutes := v1.PathPrefix("/boq").Subrouter()
		boqRoutes.Use(middleware.AuthMiddleware(authService, log))
		boqRoutes.Use(middleware.TenantIsolationMiddleware(log))
		boqRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterBOQRoutes(boqRoutes, boqService)
	}

	// ============================================
	// HR & PAYROLL MANAGEMENT ROUTES
	// ============================================
	if hrService != nil {
		hrRoutes := v1.PathPrefix("/hr").Subrouter()
		hrRoutes.Use(middleware.AuthMiddleware(authService, log))
		hrRoutes.Use(middleware.TenantIsolationMiddleware(log))
		hrRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "supervisor"},
			log,
		))
		handlers.RegisterHRRoutes(hrRoutes, hrService, rbacService)
	}

	// ============================================
	if realEstateService != nil {
		realEstateHandler := handlers.NewRealEstateHandler(realEstateService.DB, rbacService)
		realEstateRoutes := v1.PathPrefix("/real-estate").Subrouter()
		realEstateRoutes.Use(middleware.AuthMiddleware(authService, log))
		realEstateRoutes.Use(middleware.TenantIsolationMiddleware(log))
		realEstateRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))

		// Property Projects
		realEstateRoutes.HandleFunc("/projects", realEstateHandler.CreateProject).Methods("POST")
		realEstateRoutes.HandleFunc("/projects", realEstateHandler.GetProjects).Methods("GET")

		// Property Units
		realEstateRoutes.HandleFunc("/units", realEstateHandler.CreateUnit).Methods("POST")
		realEstateRoutes.HandleFunc("/projects/{project_id}/units", realEstateHandler.ListUnits).Methods("GET")

		// Customer Bookings
		realEstateRoutes.HandleFunc("/bookings", realEstateHandler.CreateBooking).Methods("POST")
		realEstateRoutes.HandleFunc("/bookings", realEstateHandler.GetBookings).Methods("GET")

		// Payments
		realEstateRoutes.HandleFunc("/payments", realEstateHandler.RecordPayment).Methods("POST")
		realEstateRoutes.HandleFunc("/bookings/{booking_id}/payments", realEstateHandler.GetPayments).Methods("GET")

		// Milestone Tracking
		realEstateRoutes.HandleFunc("/milestones", realEstateHandler.TrackMilestone).Methods("POST")
		realEstateRoutes.HandleFunc("/milestones/{booking_id}", realEstateHandler.GetMilestones).Methods("GET")

		// Account Ledger
		realEstateRoutes.HandleFunc("/ledger/{booking_id}", realEstateHandler.GetAccountLedger).Methods("GET")
	}

	// ============================================
	// PROJECT MANAGEMENT ROUTES
	// ============================================
	if realEstateService != nil {
		projectMgmtService := services.NewProjectManagementService(realEstateService.DB)
		projectMgmtHandler := handlers.NewProjectManagementHandler(projectMgmtService, realEstateService.DB)
		projectMgmtRoutes := v1.PathPrefix("/project-management").Subrouter()
		projectMgmtRoutes.Use(middleware.AuthMiddleware(authService, log))
		projectMgmtRoutes.Use(middleware.TenantIsolationMiddleware(log))
		projectMgmtRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))

		// Customer Profile endpoints
		projectMgmtRoutes.HandleFunc("/customers", projectMgmtHandler.CreateCustomerProfile).Methods("POST")
		projectMgmtRoutes.HandleFunc("/customers/{id}", projectMgmtHandler.GetCustomerProfile).Methods("GET")

		// Area Statement endpoints
		projectMgmtRoutes.HandleFunc("/area-statements", projectMgmtHandler.CreateAreaStatement).Methods("POST")

		// Cost Sheet endpoints
		projectMgmtRoutes.HandleFunc("/cost-sheets", projectMgmtHandler.UpdateCostSheet).Methods("POST")

		// Cost Configuration endpoints
		projectMgmtRoutes.HandleFunc("/cost-configurations", projectMgmtHandler.CreateProjectCostConfiguration).Methods("POST")

		// Bank Financing endpoints
		projectMgmtRoutes.HandleFunc("/bank-financing", projectMgmtHandler.CreateBankFinancing).Methods("POST")

		// Disbursement Schedule endpoints
		projectMgmtRoutes.HandleFunc("/disbursement-schedule", projectMgmtHandler.CreateDisbursementSchedule).Methods("POST")
		projectMgmtRoutes.HandleFunc("/disbursement/{id}", projectMgmtHandler.UpdateDisbursement).Methods("PUT")

		// Payment Stage endpoints
		projectMgmtRoutes.HandleFunc("/payment-stages", projectMgmtHandler.CreatePaymentStage).Methods("POST")
		projectMgmtRoutes.HandleFunc("/payment-stages/{id}/collection", projectMgmtHandler.RecordPaymentCollection).Methods("PUT")

		// Reporting endpoints
		projectMgmtRoutes.HandleFunc("/reports/bank-financing", projectMgmtHandler.GetBankFinancingReport).Methods("GET")
		projectMgmtRoutes.HandleFunc("/reports/payment-stages", projectMgmtHandler.GetPaymentStageReport).Methods("GET")
	}

	// ============================================
	// ============================================
	// BROKER MANAGEMENT (PHASE 1.2 - REAL ESTATE)
	// ============================================
	if brokerHandler != nil {
		brokerRoutes := v1.PathPrefix("/brokers").Subrouter()
		brokerRoutes.Use(middleware.AuthMiddleware(authService, log))
		brokerRoutes.Use(middleware.TenantIsolationMiddleware(log))
		brokerRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "sales"},
			log,
		))

		// Broker Profile endpoints
		brokerRoutes.HandleFunc("", brokerHandler.CreateBroker).Methods("POST")
		brokerRoutes.HandleFunc("", brokerHandler.ListBrokers).Methods("GET")
		brokerRoutes.HandleFunc("/{brokerId}", brokerHandler.GetBroker).Methods("GET")
		brokerRoutes.HandleFunc("/{brokerId}", brokerHandler.UpdateBroker).Methods("PUT")
		brokerRoutes.HandleFunc("/{brokerId}", brokerHandler.DeleteBroker).Methods("DELETE")

		// Commission Structure endpoints
		brokerRoutes.HandleFunc("/{brokerId}/commission-structure", brokerHandler.CreateCommissionStructure).Methods("POST")
		brokerRoutes.HandleFunc("/{brokerId}/commission-structure", brokerHandler.ListCommissionStructures).Methods("GET")

		// Booking Link endpoints
		brokerRoutes.HandleFunc("/booking-link", brokerHandler.CreateBookingLink).Methods("POST")
		brokerRoutes.HandleFunc("/{brokerId}/bookings", brokerHandler.ListBookingLinks).Methods("GET")
		brokerRoutes.HandleFunc("/{brokerId}/bookings/{linkId}/status", brokerHandler.UpdateBookingLinkStatus).Methods("PATCH")

		// Commission Payout endpoints
		brokerRoutes.HandleFunc("/payout", brokerHandler.CreatePayout).Methods("POST")
		brokerRoutes.HandleFunc("/payouts", brokerHandler.ListPayouts).Methods("GET")
		brokerRoutes.HandleFunc("/payouts/{payoutId}/status", brokerHandler.UpdatePayoutStatus).Methods("PATCH")

		// Reporting & Analytics endpoints
		brokerRoutes.HandleFunc("/{brokerId}/performance", brokerHandler.GetBrokerPerformance).Methods("GET")
		brokerRoutes.HandleFunc("/reports/top-performers", brokerHandler.GetTopPerformers).Methods("GET")
		brokerRoutes.HandleFunc("/reports/commission-due", brokerHandler.GetCommissionDueReport).Methods("GET")
	}

	// ============================================
	// ============================================
	// JOINT APPLICANTS (PHASE 1.3 - REAL ESTATE)
	// ============================================
	if jointApplicantHandler != nil {
		jaRoutes := v1.PathPrefix("/joint-applicants").Subrouter()
		jaRoutes.Use(middleware.AuthMiddleware(authService, log))
		jaRoutes.Use(middleware.TenantIsolationMiddleware(log))
		jaRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "sales"},
			log,
		))

		// Joint Applicant endpoints
		jaRoutes.HandleFunc("", jointApplicantHandler.CreateJointApplicant).Methods("POST")
		jaRoutes.HandleFunc("", jointApplicantHandler.ListJointApplicants).Methods("GET")
		jaRoutes.HandleFunc("/{id}", jointApplicantHandler.GetJointApplicant).Methods("GET")
		jaRoutes.HandleFunc("/{id}", jointApplicantHandler.UpdateJointApplicant).Methods("PUT")
		jaRoutes.HandleFunc("/{id}", jointApplicantHandler.DeleteJointApplicant).Methods("DELETE")

		// Document Management endpoints
		jaRoutes.HandleFunc("/{applicant_id}/documents", jointApplicantHandler.UploadDocument).Methods("POST")
		jaRoutes.HandleFunc("/{applicant_id}/documents/{document_id}/verify", jointApplicantHandler.VerifyDocument).Methods("PATCH")

		// Co-ownership Agreement endpoints
		jaRoutes.HandleFunc("/agreements", jointApplicantHandler.CreateCoOwnershipAgreement).Methods("POST")
		jaRoutes.HandleFunc("/agreements", jointApplicantHandler.ListCoOwnershipAgreements).Methods("GET")
		jaRoutes.HandleFunc("/agreements/{id}", jointApplicantHandler.GetCoOwnershipAgreement).Methods("GET")
		jaRoutes.HandleFunc("/agreements/{id}/status", jointApplicantHandler.UpdateCoOwnershipAgreementStatus).Methods("PATCH")

		// Income Verification endpoints
		jaRoutes.HandleFunc("/{applicant_id}/income-verification", jointApplicantHandler.CreateIncomeVerification).Methods("POST")
		jaRoutes.HandleFunc("/{applicant_id}/income-verification/{verification_id}/verify", jointApplicantHandler.VerifyIncome).Methods("PATCH")

		// Liability Tracking endpoints
		jaRoutes.HandleFunc("/{applicant_id}/liabilities", jointApplicantHandler.AddLiability).Methods("POST")

		// Summary & Reporting endpoints
		jaRoutes.HandleFunc("/summary/{booking_id}", jointApplicantHandler.GetJointApplicantSummary).Methods("GET")
	}

	// ============================================
	// DOCUMENT MANAGEMENT ROUTES
	// ============================================
	if documentHandler != nil {
		docRoutes := v1.PathPrefix("/documents").Subrouter()
		docRoutes.Use(middleware.AuthMiddleware(authService, log))
		docRoutes.Use(middleware.TenantIsolationMiddleware(log))
		docRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "sales", "user"},
			log,
		))

		// Document CRUD endpoints
		docRoutes.HandleFunc("", documentHandler.CreateDocument).Methods("POST")
		docRoutes.HandleFunc("", documentHandler.ListDocuments).Methods("GET")
		docRoutes.HandleFunc("/{id}", documentHandler.GetDocument).Methods("GET")
		docRoutes.HandleFunc("/{id}", documentHandler.UpdateDocument).Methods("PUT")
		docRoutes.HandleFunc("/{id}", documentHandler.DeleteDocument).Methods("DELETE")

		// Document verification endpoint
		docRoutes.HandleFunc("/{id}/verify", documentHandler.VerifyDocument).Methods("PATCH")

		// Document search endpoints
		docRoutes.HandleFunc("/by-type/{type_id}", documentHandler.GetDocumentsByType).Methods("GET")
		docRoutes.HandleFunc("/expiring", documentHandler.CheckDocumentExpiry).Methods("GET")

		// Document collection endpoints
		docRoutes.HandleFunc("/collections", documentHandler.CreateDocumentCollection).Methods("POST")
		docRoutes.HandleFunc("/collections", documentHandler.ListDocumentCollections).Methods("GET")
		docRoutes.HandleFunc("/collections/{id}", documentHandler.GetDocumentCollection).Methods("GET")
		docRoutes.HandleFunc("/collections/{id}/status", documentHandler.UpdateCollectionStatus).Methods("PATCH")
		docRoutes.HandleFunc("/collections/{id}/items", documentHandler.AddDocumentToCollection).Methods("POST")
		docRoutes.HandleFunc("/collections/{id}/items/{item_id}", documentHandler.RemoveDocumentFromCollection).Methods("DELETE")

		// Document sharing endpoints
		docRoutes.HandleFunc("/{id}/share", documentHandler.ShareDocument).Methods("POST")
		docRoutes.HandleFunc("/{id}/shares", documentHandler.GetDocumentShares).Methods("GET")
		docRoutes.HandleFunc("/{id}/shares/{share_id}", documentHandler.RevokeDocumentShare).Methods("DELETE")

		// Document template endpoints
		docRoutes.HandleFunc("/templates", documentHandler.CreateDocumentTemplate).Methods("POST")
		docRoutes.HandleFunc("/templates", documentHandler.ListDocumentTemplates).Methods("GET")

		// Document metadata endpoints
		docRoutes.HandleFunc("/categories", documentHandler.ListDocumentCategories).Methods("GET")
		docRoutes.HandleFunc("/types/{category_id}", documentHandler.ListDocumentTypes).Methods("GET")

		// Document reporting endpoint
		docRoutes.HandleFunc("/summary/{booking_id}", documentHandler.GetDocumentSummary).Methods("GET")
	}

	// ============================================
	// POSSESSION MANAGEMENT ROUTES (Phase 2.1)
	// ============================================
	if possessionHandler != nil {
		posRoutes := v1.PathPrefix("/possessions").Subrouter()
		posRoutes.Use(middleware.AuthMiddleware(authService, log))
		posRoutes.Use(middleware.TenantIsolationMiddleware(log))
		posRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "sales", "user"},
			log,
		))

		// Possession status endpoints
		posRoutes.HandleFunc("", possessionHandler.CreatePossession).Methods("POST")
		posRoutes.HandleFunc("", possessionHandler.ListPossessions).Methods("GET")
		posRoutes.HandleFunc("/{id}", possessionHandler.GetPossession).Methods("GET")
		posRoutes.HandleFunc("/{id}", possessionHandler.UpdatePossession).Methods("PUT")
		posRoutes.HandleFunc("/{id}", possessionHandler.DeletePossession).Methods("DELETE")

		// Possession document endpoints
		posRoutes.HandleFunc("/{possession_id}/documents", possessionHandler.CreatePossessionDocument).Methods("POST")
		posRoutes.HandleFunc("/{possession_id}/documents", possessionHandler.ListPossessionDocuments).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/documents/{doc_id}", possessionHandler.GetPossessionDocument).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/documents/{doc_id}", possessionHandler.UpdatePossessionDocument).Methods("PUT")
		posRoutes.HandleFunc("/{possession_id}/documents/{doc_id}/verify", possessionHandler.VerifyPossessionDocument).Methods("PATCH")

		// Possession registration endpoints
		posRoutes.HandleFunc("/{possession_id}/registrations", possessionHandler.CreateRegistration).Methods("POST")
		posRoutes.HandleFunc("/{possession_id}/registrations", possessionHandler.ListRegistrations).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/registrations/{reg_id}", possessionHandler.GetRegistration).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/registrations/{reg_id}", possessionHandler.UpdateRegistration).Methods("PUT")
		posRoutes.HandleFunc("/{possession_id}/registrations/{reg_id}/approve", possessionHandler.ApproveRegistration).Methods("PATCH")

		// Possession certificate endpoints
		posRoutes.HandleFunc("/{possession_id}/certificates", possessionHandler.CreateCertificate).Methods("POST")
		posRoutes.HandleFunc("/{possession_id}/certificates", possessionHandler.ListCertificates).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/certificates/{cert_id}", possessionHandler.GetCertificate).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/certificates/{cert_id}/verify", possessionHandler.VerifyCertificate).Methods("PATCH")

		// Possession approval endpoints
		posRoutes.HandleFunc("/{possession_id}/approvals", possessionHandler.CreateApproval).Methods("POST")
		posRoutes.HandleFunc("/{possession_id}/approvals", possessionHandler.ListApprovals).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/approvals/{approval_id}", possessionHandler.GetApproval).Methods("GET")
		posRoutes.HandleFunc("/{possession_id}/approvals/{approval_id}/approve", possessionHandler.ApprovePossession).Methods("PATCH")

		// Possession summary endpoint
		posRoutes.HandleFunc("/{possession_id}/summary", possessionHandler.GetPossessionSummary).Methods("GET")
	}

	// TITLE CLEARANCE MANAGEMENT ROUTES (Phase 2.2)
	// ============================================
	if titleHandler != nil {
		titleRoutes := v1.PathPrefix("/title-clearances").Subrouter()
		titleRoutes.Use(middleware.AuthMiddleware(authService, log))
		titleRoutes.Use(middleware.TenantIsolationMiddleware(log))
		titleRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "legal", "user"},
			log,
		))

		// Title clearance endpoints
		titleRoutes.HandleFunc("", titleHandler.CreateTitleClearance).Methods("POST")
		titleRoutes.HandleFunc("", titleHandler.ListTitleClearances).Methods("GET")
		titleRoutes.HandleFunc("/{id}", titleHandler.GetTitleClearance).Methods("GET")
		titleRoutes.HandleFunc("/{id}", titleHandler.UpdateTitleClearance).Methods("PUT")
		titleRoutes.HandleFunc("/{id}/summary", titleHandler.GetClearanceSummary).Methods("GET")

		// Title issue endpoints
		titleRoutes.HandleFunc("/{clearance_id}/issues", titleHandler.CreateTitleIssue).Methods("POST")
		titleRoutes.HandleFunc("/{clearance_id}/issues", titleHandler.ListTitleIssues).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/issues/{issue_id}", titleHandler.GetTitleIssue).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/issues/{issue_id}", titleHandler.UpdateTitleIssue).Methods("PUT")
		titleRoutes.HandleFunc("/{clearance_id}/issues/{issue_id}/resolve", titleHandler.ResolveTitleIssue).Methods("PATCH")

		// Title search report endpoints
		titleRoutes.HandleFunc("/{clearance_id}/search-reports", titleHandler.CreateSearchReport).Methods("POST")
		titleRoutes.HandleFunc("/{clearance_id}/search-reports", titleHandler.ListSearchReports).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/search-reports/{report_id}", titleHandler.GetSearchReport).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/search-reports/{report_id}/verify", titleHandler.VerifySearchReport).Methods("PATCH")

		// Title legal opinion endpoints
		titleRoutes.HandleFunc("/{clearance_id}/legal-opinions", titleHandler.CreateLegalOpinion).Methods("POST")
		titleRoutes.HandleFunc("/{clearance_id}/legal-opinions", titleHandler.ListLegalOpinions).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/legal-opinions/{opinion_id}", titleHandler.GetLegalOpinion).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/legal-opinions/{opinion_id}/review", titleHandler.ReviewLegalOpinion).Methods("PATCH")

		// Verification checklist endpoints
		titleRoutes.HandleFunc("/{clearance_id}/checklists", titleHandler.ListVerificationChecklists).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/checklists/{checklist_id}/verify", titleHandler.VerifyChecklistItem).Methods("PATCH")

		// Approval endpoints
		titleRoutes.HandleFunc("/{clearance_id}/approvals", titleHandler.ListClearanceApprovals).Methods("GET")
		titleRoutes.HandleFunc("/{clearance_id}/approvals/{approval_id}/approve", titleHandler.ApproveClearance).Methods("PATCH")
	}

	// CUSTOMER PORTAL ROUTES (Phase 2.3)
	// ============================================
	if customerPortalHandler != nil {
		customerRoutes := v1.PathPrefix("/customer").Subrouter()
		customerRoutes.Use(middleware.AuthMiddleware(authService, log))
		customerRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Customer profile endpoints
		profileRoutes := customerRoutes.PathPrefix("/profile").Subrouter()
		profileRoutes.HandleFunc("", customerPortalHandler.CreateProfile).Methods("POST")
		profileRoutes.HandleFunc("", customerPortalHandler.GetProfile).Methods("GET")
		profileRoutes.HandleFunc("", customerPortalHandler.UpdateProfile).Methods("PUT")

		// Customer notification endpoints
		notifRoutes := customerRoutes.PathPrefix("/notifications").Subrouter()
		notifRoutes.HandleFunc("", customerPortalHandler.GetNotifications).Methods("GET")
		notifRoutes.HandleFunc("/{id}/read", customerPortalHandler.MarkNotificationAsRead).Methods("PATCH")

		// Customer conversation endpoints
		convRoutes := customerRoutes.PathPrefix("/conversations").Subrouter()
		convRoutes.HandleFunc("", customerPortalHandler.CreateConversation).Methods("POST")
		convRoutes.HandleFunc("/{conversation_id}/messages", customerPortalHandler.SendMessage).Methods("POST")
		convRoutes.HandleFunc("/{conversation_id}/messages", customerPortalHandler.GetConversationMessages).Methods("GET")

		// Customer document upload endpoints
		docRoutes := customerRoutes.PathPrefix("/documents").Subrouter()
		docRoutes.HandleFunc("", customerPortalHandler.UploadDocument).Methods("POST")

		// Customer booking tracking endpoints
		bookingRoutes := customerRoutes.PathPrefix("/bookings").Subrouter()
		bookingRoutes.HandleFunc("/{booking_id}/tracking", customerPortalHandler.GetBookingTracking).Methods("GET")

		// Customer payment tracking endpoints
		paymentRoutes := customerRoutes.PathPrefix("/payments").Subrouter()
		paymentRoutes.HandleFunc("/{booking_id}/history", customerPortalHandler.GetPayments).Methods("GET")

		// Customer feedback endpoints
		feedbackRoutes := customerRoutes.PathPrefix("/feedback").Subrouter()
		feedbackRoutes.HandleFunc("", customerPortalHandler.CreateFeedback).Methods("POST")

		// Customer preferences endpoints
		prefRoutes := customerRoutes.PathPrefix("/preferences").Subrouter()
		prefRoutes.HandleFunc("", customerPortalHandler.UpdatePreferences).Methods("PUT")
	}

	// ANALYTICS ROUTES (Phase 3.1)
	// ============================================
	if analyticsHandler != nil {
		analyticsRoutes := v1.PathPrefix("/analytics").Subrouter()
		analyticsRoutes.Use(middleware.AuthMiddleware(authService, log))
		analyticsRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Report generation endpoints
		analyticsRoutes.HandleFunc("/reports", analyticsHandler.GenerateReport).Methods("POST")
		analyticsRoutes.HandleFunc("/export", analyticsHandler.ExportReport).Methods("POST")

		// Trends endpoints
		analyticsRoutes.HandleFunc("/trends", analyticsHandler.GetTrends).Methods("GET")

		// Custom metrics endpoints
		analyticsRoutes.HandleFunc("/metrics", analyticsHandler.GetCustomMetrics).Methods("GET")
	}

	// ============================================
	// ============================================
	// MOBILE APP FEATURES (Phase 3.2) ROUTES
	// ============================================
	if mobileHandler != nil {
		mobileRoutes := v1.PathPrefix("/mobile").Subrouter()
		mobileRoutes.Use(middleware.AuthMiddleware(authService, log))
		mobileRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// Mobile App Configuration endpoints
		mobileRoutes.HandleFunc("/apps", mobileHandler.CreateApp).Methods("POST")
		mobileRoutes.HandleFunc("/apps", mobileHandler.ListApps).Methods("GET")
		mobileRoutes.HandleFunc("/apps/{id}", mobileHandler.GetApp).Methods("GET")

		// Mobile Device Registration endpoints
		mobileRoutes.HandleFunc("/devices", mobileHandler.RegisterDevice).Methods("POST")
		mobileRoutes.HandleFunc("/devices/{id}", mobileHandler.GetDevice).Methods("GET")
		mobileRoutes.HandleFunc("/my-devices", mobileHandler.ListUserDevices).Methods("GET")

		// Mobile Session endpoints
		mobileRoutes.HandleFunc("/sessions", mobileHandler.CreateSession).Methods("POST")

		// Push Notification endpoints
		mobileRoutes.HandleFunc("/notifications", mobileHandler.SendNotification).Methods("POST")
		mobileRoutes.HandleFunc("/notifications", mobileHandler.GetUserNotifications).Methods("GET")
		mobileRoutes.HandleFunc("/notifications/{id}/read", mobileHandler.MarkNotificationAsRead).Methods("PUT")

		// Feature Flags endpoints
		mobileRoutes.HandleFunc("/features", mobileHandler.GetAppFeatures).Methods("GET")

		// Offline Data Sync endpoints
		mobileRoutes.HandleFunc("/sync", mobileHandler.SyncOfflineData).Methods("POST")

		// Crash Reporting endpoints
		mobileRoutes.HandleFunc("/crashes", mobileHandler.ReportCrash).Methods("POST")

		// Analytics endpoints
		mobileRoutes.HandleFunc("/events", mobileHandler.TrackEvent).Methods("POST")

		// App Update endpoints
		mobileRoutes.HandleFunc("/updates/latest", mobileHandler.GetLatestUpdate).Methods("GET")

		// User Settings endpoints
		mobileRoutes.HandleFunc("/settings", mobileHandler.GetUserSettings).Methods("GET")

		// Endpoints documentation:
		// POST   /api/v1/mobile/apps - Create new app configuration
		// GET    /api/v1/mobile/apps - List all apps
		// GET    /api/v1/mobile/apps/{id} - Get app by ID
		// POST   /api/v1/mobile/devices - Register mobile device
		// GET    /api/v1/mobile/devices/{id} - Get device info
		// GET    /api/v1/mobile/my-devices - Get user's devices
		// POST   /api/v1/mobile/sessions - Create mobile session
		// POST   /api/v1/mobile/notifications - Send notification
		// GET    /api/v1/mobile/notifications - Get user notifications
		// PUT    /api/v1/mobile/notifications/{id}/read - Mark as read
		// GET    /api/v1/mobile/features - Get feature flags
		// POST   /api/v1/mobile/sync - Sync offline data
		// POST   /api/v1/mobile/crashes - Report crash
		// POST   /api/v1/mobile/events - Track analytics event
		// GET    /api/v1/mobile/updates/latest - Get latest update
		// GET    /api/v1/mobile/settings - Get user settings
	}

	// ============================================
	// AI POWERED RECOMMENDATIONS (PHASE 3.3)
	// ============================================
	if aiRecommendationsHandler != nil {
		aiRoutes := v1.PathPrefix("/ai").Subrouter()
		aiRoutes.Use(middleware.AuthMiddleware(authService, log))
		aiRoutes.Use(middleware.TenantIsolationMiddleware(log))

		// AI Models endpoints
		aiRoutes.HandleFunc("/models", aiRecommendationsHandler.CreateAIModel).Methods("POST")
		aiRoutes.HandleFunc("/models", aiRecommendationsHandler.ListAIModels).Methods("GET")
		aiRoutes.HandleFunc("/models/{id}", aiRecommendationsHandler.GetAIModel).Methods("GET")

		// Recommendation Engine endpoints
		aiRoutes.HandleFunc("/engines", aiRecommendationsHandler.CreateRecommendationEngine).Methods("POST")
		aiRoutes.HandleFunc("/engines/{id}", aiRecommendationsHandler.GetRecommendationEngine).Methods("GET")

		// User Recommendations endpoints
		aiRoutes.HandleFunc("/recommendations/generate", aiRecommendationsHandler.GenerateRecommendations).Methods("POST")
		aiRoutes.HandleFunc("/recommendations", aiRecommendationsHandler.GetUserRecommendations).Methods("GET")

		// Feedback endpoints
		aiRoutes.HandleFunc("/feedback", aiRecommendationsHandler.SubmitRecommendationFeedback).Methods("POST")
		aiRoutes.HandleFunc("/feedback", aiRecommendationsHandler.GetRecommendationFeedback).Methods("GET")

		// Prediction endpoints
		aiRoutes.HandleFunc("/predictions", aiRecommendationsHandler.MakePrediction).Methods("POST")
		aiRoutes.HandleFunc("/predictions/{id}", aiRecommendationsHandler.GetPredictionResult).Methods("GET")

		// Anomaly Detection endpoints
		aiRoutes.HandleFunc("/anomalies/detect", aiRecommendationsHandler.DetectAnomalies).Methods("POST")
		aiRoutes.HandleFunc("/anomalies", aiRecommendationsHandler.GetAnomalies).Methods("GET")

		// AI Insights endpoints
		aiRoutes.HandleFunc("/insights/generate", aiRecommendationsHandler.GenerateInsights).Methods("POST")
		aiRoutes.HandleFunc("/insights", aiRecommendationsHandler.GetInsights).Methods("GET")

		// Model Performance endpoints
		aiRoutes.HandleFunc("/performance", aiRecommendationsHandler.GetModelPerformance).Methods("GET")

		// Statistics endpoints
		aiRoutes.HandleFunc("/stats", aiRecommendationsHandler.GetAIStats).Methods("GET")
		aiRoutes.HandleFunc("/recommendations/stats", aiRecommendationsHandler.GetRecommendationStats).Methods("GET")

		// Endpoints documentation:
		// POST   /api/v1/ai/models - Create AI model
		// GET    /api/v1/ai/models - List AI models
		// GET    /api/v1/ai/models/{id} - Get AI model
		// POST   /api/v1/ai/engines - Create recommendation engine
		// GET    /api/v1/ai/engines/{id} - Get recommendation engine
		// POST   /api/v1/ai/recommendations/generate - Generate recommendations
		// GET    /api/v1/ai/recommendations - Get user recommendations
		// POST   /api/v1/ai/feedback - Submit recommendation feedback
		// GET    /api/v1/ai/feedback - Get recommendation feedback
		// POST   /api/v1/ai/predictions - Make prediction
		// GET    /api/v1/ai/predictions/{id} - Get prediction result
		// POST   /api/v1/ai/anomalies/detect - Detect anomalies
		// GET    /api/v1/ai/anomalies - Get anomalies
		// POST   /api/v1/ai/insights/generate - Generate insights
		// GET    /api/v1/ai/insights - Get insights
		// GET    /api/v1/ai/performance - Get model performance
		// GET    /api/v1/ai/stats - Get AI system statistics
		// GET    /api/v1/ai/recommendations/stats - Get recommendation statistics
	}

	// ============================================
	// ============================================
	// GENERAL LEDGER (GL) ROUTES
	// ============================================
	if glService != nil {
		glRoutes := v1.PathPrefix("/gl").Subrouter()
		glRoutes.Use(middleware.AuthMiddleware(authService, log))
		glRoutes.Use(middleware.TenantIsolationMiddleware(log))
		glRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "accountant"},
			log,
		))
		handlers.RegisterGLRoutes(glRoutes, glService, rbacService)
	}

	// Compliance Routes (RERA, HR, Tax)
	if reraComplianceHandler != nil {
		reraRoutes := v1.PathPrefix("/rera-compliance").Subrouter()
		reraRoutes.Use(middleware.AuthMiddleware(authService, log))
		reraRoutes.Use(middleware.TenantIsolationMiddleware(log))
		reraRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterRERARoutes(reraRoutes, reraComplianceHandler)
	}

	if hrComplianceHandler != nil {
		hrComplianceRoutes := v1.PathPrefix("/hr-compliance").Subrouter()
		hrComplianceRoutes.Use(middleware.AuthMiddleware(authService, log))
		hrComplianceRoutes.Use(middleware.TenantIsolationMiddleware(log))
		hrComplianceRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "supervisor"},
			log,
		))
		handlers.RegisterHRComplianceRoutes(hrComplianceRoutes, hrComplianceHandler)
	}

	if taxComplianceHandler != nil {
		taxRoutes := v1.PathPrefix("/tax-compliance").Subrouter()
		taxRoutes.Use(middleware.AuthMiddleware(authService, log))
		taxRoutes.Use(middleware.TenantIsolationMiddleware(log))
		taxRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterTaxComplianceRoutes(taxRoutes, taxComplianceHandler)
	}

	// Dashboard Routes
	if financialDashboardHandler != nil {
		finDashRoutes := v1.PathPrefix("/financial-dashboard").Subrouter()
		finDashRoutes.Use(middleware.AuthMiddleware(authService, log))
		finDashRoutes.Use(middleware.TenantIsolationMiddleware(log))
		finDashRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterFinancialDashboardRoutes(finDashRoutes, financialDashboardHandler)
	}

	if hrDashboardHandler != nil {
		hrDashRoutes := v1.PathPrefix("/hr-dashboard").Subrouter()
		hrDashRoutes.Use(middleware.AuthMiddleware(authService, log))
		hrDashRoutes.Use(middleware.TenantIsolationMiddleware(log))
		hrDashRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "supervisor"},
			log,
		))
		handlers.RegisterHRDashboardRoutes(hrDashRoutes, hrDashboardHandler)
	}

	if complianceDashboardHandler != nil {
		compDashRoutes := v1.PathPrefix("/compliance-dashboard").Subrouter()
		compDashRoutes.Use(middleware.AuthMiddleware(authService, log))
		compDashRoutes.Use(middleware.TenantIsolationMiddleware(log))
		compDashRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager"},
			log,
		))
		handlers.RegisterComplianceDashboardRoutes(compDashRoutes, complianceDashboardHandler)
	}

	if salesDashboardHandler != nil {
		salesDashRoutes := v1.PathPrefix("/sales-dashboard").Subrouter()
		salesDashRoutes.Use(middleware.AuthMiddleware(authService, log))
		salesDashRoutes.Use(middleware.TenantIsolationMiddleware(log))
		salesDashRoutes.Use(middleware.PermissionBasedAccessMiddleware(
			rbacService,
			[]string{"admin", "manager", "agent"},
			log,
		))
		handlers.RegisterSalesDashboardRoutes(salesDashRoutes, salesDashboardHandler)
	}

	// OPTIONS handler for CORS preflight requests
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusOK)
	})

	// 404 handler
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	log.Info("Routes configured successfully")
	return r
}

// HealthCheck returns service health status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// ReadinessCheck returns service readiness status
func ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}

// NotFoundHandler handles 404 responses
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"endpoint not found"}`))
}

// Placeholder handlers for leads, calls, campaigns, etc.

func GetLeads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"leads":[]}`))
}

func GetLead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"lead":{}}`))
}

func CreateLead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"lead created"}`))
}

func UpdateLead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"lead updated"}`))
}

func GetCalls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"calls":[]}`))
}

func GetCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"call":{}}`))
}

func InitiateCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"call initiated"}`))
}

func EndCall(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"call ended"}`))
}

func ProcessAIQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response":""}`))
}

func ListAIProviders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"providers":[]}`))
}

func GetCampaigns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"campaigns":[]}`))
}

func GetCampaign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"campaign":{}}`))
}

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"campaign created"}`))
}

func UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"campaign updated"}`))
}
