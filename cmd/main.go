package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"vyomtech-backend/internal/config"
	"vyomtech-backend/internal/db"
	"vyomtech-backend/internal/handlers"
	"vyomtech-backend/internal/services"
	"vyomtech-backend/pkg/auth"
	"vyomtech-backend/pkg/logger"
	"vyomtech-backend/pkg/router"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize logger
	log := logger.New()
	log.Info("Starting Multi-Tenant AI Call Center application")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize database connection
	dbConn, err := db.NewDatabaseConnection(&cfg.Database, log)
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	// Seed demo users if in development mode
	if os.Getenv("APP_ENV") != "production" {
		seeder := services.NewDemoUserSeeder(dbConn, log)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := seeder.SeedDemoUsers(ctx); err != nil {
			log.Warn("Failed to seed demo users", "error", err)
		}
		cancel()
	}

	// Initialize and start demo reset scheduler
	demoResetService := services.NewDemoResetService(dbConn, log)
	if err := demoResetService.ResetDemoData(); err != nil {
		log.Warn("Initial demo data reset failed", "error", err)
	}
	demoResetService.StartScheduler()
	log.Info("Demo reset scheduler started - will reset every 30 days")

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// Initialize services
	authService := services.NewAuthService(dbConn, jwtManager, log)
	tenantService := services.NewTenantService(dbConn, log)
	emailService := services.NewEmailService(&cfg.Email, log)
	passwordResetService := services.NewPasswordResetService(dbConn, emailService, log)
	agentService := services.NewAgentService(dbConn, log)
	gamificationService := services.NewGamificationService(dbConn)
	leadService := services.NewLeadService(dbConn)
	callService := services.NewCallService(dbConn)
	campaignService := services.NewCampaignService(dbConn)
	aiOrchestrator := services.NewAIOrchestrator(dbConn, log)
	webSocketHub := services.NewWebSocketHub(log)

	// Phase 1 Services
	leadScoringService := services.NewLeadScoringService(dbConn, log)
	dashboardService := services.NewDashboardService(dbConn, log)

	// Phase 2 Services
	taskService := services.NewTaskService(dbConn)
	notificationService := services.NewNotificationService(dbConn)
	tenantCustomizationService := services.NewTenantCustomizationService(dbConn)

	// Phase 3C Services (Modular Monetization)
	phase3cServices := services.NewPhase3CServices(dbConn, log)

	// Sales Service
	salesService := services.NewSalesService(dbConn)

	// Real Estate Service
	realEstateService := services.NewRealEstateService(dbConn)

	// Broker Service (Phase 1.2 Real Estate)
	brokerService := services.NewBrokerService(dbConn)

	// Joint Applicant Service (Phase 1.3 Real Estate)
	jointApplicantService := services.NewJointApplicantService(dbConn)

	// Document Management Service (Phase 1.4 Real Estate)
	documentService := services.NewDocumentService(dbConn)

	// Possession Management Service (Phase 2.1 Real Estate)
	possessionService := services.NewPossessionService(dbConn)

	// Title Clearance Management Service (Phase 2.2 Real Estate)
	titleService := services.NewTitleService(dbConn)

	// Customer Portal Service (Phase 2.3 Real Estate)
	customerPortalService := services.NewCustomerPortalService(dbConn)

	// Advanced Analytics Service (Phase 3.1 Enhancement)
	analyticsService := services.NewAnalyticsService(dbConn)

	// AI Recommendations Service (Phase 3.3)
	aiService := services.NewAIService(dbConn)

	// Civil Engineering Service
	civilService := services.NewCivilService(dbConn)

	// Construction Service
	constructionService := services.NewConstructionService(dbConn)

	// BOQ Service
	boqService := services.NewBOQService(dbConn)

	// HR & Payroll Service
	hrService := services.NewHRService(dbConn)

	// GL (General Ledger) Service
	glService := services.NewGLService(dbConn)

	// RBAC Service for permission checking
	rbacService := services.NewRBACService(dbConn, log)

	// Compliance Services (RERA, HR, Tax)
	reraComplianceService := services.NewRERAComplianceService(dbConn)
	hrComplianceService := services.NewHRComplianceService(dbConn)
	taxComplianceService := services.NewTaxComplianceService(dbConn)

	// Initialize handlers
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)

	// Admin Handlers
	userAdminHandler := handlers.NewUserAdminHandler(dbConn, log)
	tenantAdminHandler := handlers.NewTenantAdminHandler(dbConn, log)

	// Compliance Handlers
	reraComplianceHandler := handlers.NewRERAComplianceHandler(reraComplianceService)
	hrComplianceHandler := handlers.NewHRComplianceHandler(hrComplianceService)
	taxComplianceHandler := handlers.NewTaxComplianceHandler(taxComplianceService)

	// Broker Handler (Phase 1.2 Real Estate)
	brokerHandler := handlers.NewBrokerHandler(brokerService)

	// Joint Applicant Handler (Phase 1.3 Real Estate)
	jointApplicantHandler := handlers.NewJointApplicantHandler(jointApplicantService)

	// Document Management Handler (Phase 1.4 Real Estate)
	documentHandler := handlers.NewDocumentHandler(documentService)

	// Possession Management Handler (Phase 2.1 Real Estate)
	possessionHandler := handlers.NewPossessionHandler(possessionService)

	// Title Clearance Management Handler (Phase 2.2 Real Estate)
	titleHandler := handlers.NewTitleHandler(titleService)

	// Customer Portal Handler (Phase 2.3 Real Estate)
	customerPortalHandler := handlers.NewCustomerPortalHandler(customerPortalService)

	// Advanced Analytics Handler (Phase 3.1 Enhancement)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, log)

	// AI Recommendations Handler (Phase 3.3)
	aiRecommendationsHandler := handlers.NewAIRecommendationsHandler(aiService, log)

	// Phase 3.2: Mobile App Features Service & Handler
	mobileService := services.NewMobileService(dbConn)
	mobileHandler := handlers.NewMobileHandler(mobileService, log)

	// Site Visit Service & Handler
	siteVisitService := services.NewSiteVisitService(dbConn)
	siteVisitHandler := handlers.NewSiteVisitHandler(siteVisitService, log.Logger)

	// Phase 3.4: Integration Hub Service & Handler
	integrationService := services.NewIntegrationService(dbConn)
	integrationHandler := handlers.NewIntegrationHandler(integrationService, log.Logger)

	// Dashboard Handlers
	financialDashboardHandler := handlers.NewFinancialDashboardHandler(glService)
	hrDashboardHandler := handlers.NewHRDashboardHandler(hrService, hrComplianceService)
	complianceDashboardHandler := handlers.NewComplianceDashboardHandler(reraComplianceService, hrComplianceService, taxComplianceService)
	salesDashboardHandler := handlers.NewSalesDashboardHandler(salesService)

	// Setup router with all services
	r := router.SetupRoutesWithPhase3C(authService, tenantService, passwordResetHandler, agentService, gamificationService, leadService, callService, campaignService, aiOrchestrator, webSocketHub, leadScoringService, dashboardService, taskService, notificationService, tenantCustomizationService, phase3cServices, salesService, realEstateService, civilService, constructionService, boqService, hrService, glService, rbacService, reraComplianceHandler, hrComplianceHandler, taxComplianceHandler, financialDashboardHandler, hrDashboardHandler, complianceDashboardHandler, salesDashboardHandler, brokerHandler, jointApplicantHandler, documentHandler, possessionHandler, titleHandler, customerPortalHandler, analyticsHandler, userAdminHandler, tenantAdminHandler, mobileHandler, aiRecommendationsHandler, siteVisitHandler, integrationHandler, log)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info("Server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
		}
	}()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Shutdown signal received, gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", "error", err)
	}

	log.Info("Server stopped")
}
