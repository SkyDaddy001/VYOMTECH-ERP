package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"multi-tenant-ai-callcenter/internal/config"
	"multi-tenant-ai-callcenter/internal/db"
	"multi-tenant-ai-callcenter/internal/handlers"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/auth"
	"multi-tenant-ai-callcenter/pkg/logger"
	"multi-tenant-ai-callcenter/pkg/router"
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

	// Compliance Services (RERA, HR, Tax)
	reraComplianceService := services.NewRERAComplianceService(dbConn)
	hrComplianceService := services.NewHRComplianceService(dbConn)
	taxComplianceService := services.NewTaxComplianceService(dbConn)

	// Initialize handlers
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)

	// Compliance Handlers
	reraComplianceHandler := handlers.NewRERAComplianceHandler(reraComplianceService)
	hrComplianceHandler := handlers.NewHRComplianceHandler(hrComplianceService)
	taxComplianceHandler := handlers.NewTaxComplianceHandler(taxComplianceService)

	// Dashboard Handlers
	financialDashboardHandler := handlers.NewFinancialDashboardHandler(glService)
	hrDashboardHandler := handlers.NewHRDashboardHandler(hrService, hrComplianceService)
	complianceDashboardHandler := handlers.NewComplianceDashboardHandler(reraComplianceService, hrComplianceService, taxComplianceService)
	salesDashboardHandler := handlers.NewSalesDashboardHandler(salesService)

	// Setup router with all services
	r := router.SetupRoutesWithPhase3C(authService, tenantService, passwordResetHandler, agentService, gamificationService, leadService, callService, campaignService, aiOrchestrator, webSocketHub, leadScoringService, dashboardService, taskService, notificationService, tenantCustomizationService, phase3cServices, salesService, realEstateService, civilService, constructionService, boqService, hrService, glService, reraComplianceHandler, hrComplianceHandler, taxComplianceHandler, financialDashboardHandler, hrDashboardHandler, complianceDashboardHandler, salesDashboardHandler, log)

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
