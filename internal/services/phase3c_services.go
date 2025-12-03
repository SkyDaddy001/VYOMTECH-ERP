package services

import (
	"database/sql"

	"vyomtech-backend/pkg/logger"
)

// Phase3CServices contains all Phase 3C (Modular Monetization) services
type Phase3CServices struct {
	ModuleService  *ModuleService
	CompanyService *CompanyService
	BillingService *BillingService
}

// NewPhase3CServices initializes all Phase 3C services
func NewPhase3CServices(db *sql.DB, log *logger.Logger) *Phase3CServices {
	moduleService := NewModuleService(db, log)
	companyService := NewCompanyService(db, log)
	billingService := NewBillingService(db, log, moduleService)

	return &Phase3CServices{
		ModuleService:  moduleService,
		CompanyService: companyService,
		BillingService: billingService,
	}
}
