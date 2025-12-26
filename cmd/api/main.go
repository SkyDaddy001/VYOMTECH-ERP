package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Create Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheck)

	// Register mock OAuth routes for development/testing
	RegisterMockOAuthRoutes(router)

	// Register test data routes for development/testing
	RegisterTestDataRoutes(router)

	// Register Phase 8: API Endpoints
	registerGoogleAdsRoutes(router)
	registerMetaAdsRoutes(router)
	registerROIRoutes(router)

	// Register Phase 8: Test & Documentation Routes
	registerPhase8TestRoutes(router)

	// ============================================================
	// PHASE 10: Initialize Sync Job Scheduler
	// ============================================================
	ctx := context.Background()
	scheduler = initSyncJobScheduler()

	// Register Phase 10 sync jobs
	if err := initSyncJobs(scheduler); err != nil {
		log.Printf("Failed to initialize sync jobs: %v", err)
	}

	// Start scheduler
	scheduler.Start(ctx)

	// Register Phase 10: Sync Status Endpoints
	registerPhase10SyncRoutes(router)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		v1.POST("/auth/login", login)
		v1.POST("/auth/register", register)
		v1.POST("/auth/refresh", refresh)

		// Protected routes
		protected := v1.Group("")
		protected.Use(authMiddleware())
		{
			// Tenant routes
			protected.GET("/tenants", listTenants)
			protected.GET("/tenants/:id", getTenant)
			protected.PUT("/tenants/:id", updateTenant)
			protected.POST("/tenants", createTenant)

			// User routes
			protected.GET("/users", listUsers)
			protected.GET("/users/:id", getUser)
			protected.PUT("/users/:id", updateUser)
			protected.DELETE("/users/:id", deleteUser)

			// User Count & Seat Management routes
			protected.GET("/tenant/users/count", getTenantUserCount)
			protected.GET("/tenant/users/breakdown", getTenantUserBreakdown)
			protected.GET("/tenant/users/history", getUserCountHistory)
			protected.GET("/tenant/users/active", getActiveUsersRealtime)
			protected.GET("/tenant/users/seat-available", checkSeatAvailability)
			protected.GET("/tenant/users/list", listTenantUsers)
			protected.GET("/users/:id/activity", getUserActivity)
			protected.GET("/tenant/users/stats", getUserCountStats)
			protected.GET("/tenant/users/report", exportUserCountReport)
			protected.PATCH("/tenant/users/limit", updateTenantUserLimit)
			protected.POST("/users/activity", recordUserActivity)
			protected.GET("/tenant/billing/overage-check", checkBillingOverage)

			// Call Center routes
			protected.GET("/call-centers", listCallCenters)
			protected.GET("/call-centers/:id", getCallCenter)
			protected.PUT("/call-centers/:id", updateCallCenter)
			protected.POST("/call-centers", createCallCenter)

			// Agent routes
			protected.GET("/agents", listAgents)
			protected.GET("/agents/:id", getAgent)
			protected.PUT("/agents/:id", updateAgent)
			protected.POST("/agents", createAgent)

			// Call routes
			protected.GET("/calls", listCalls)
			protected.GET("/calls/:id", getCall)
			protected.PUT("/calls/:id", updateCall)
			protected.POST("/calls", createCall)

			// Campaign routes
			protected.GET("/campaigns", listCampaigns)
			protected.GET("/campaigns/:id", getCampaign)
			protected.PUT("/campaigns/:id", updateCampaign)
			protected.POST("/campaigns", createCampaign)

			// Sales Lead routes
			protected.GET("/sales-leads", listSalesLeads)
			protected.GET("/sales-leads/:id", getSalesLead)
			protected.PUT("/sales-leads/:id", updateSalesLead)
			protected.POST("/sales-leads", createSalesLead)

			// Inventory routes
			protected.GET("/inventory-items", listInventoryItems)
			protected.GET("/inventory-items/:id", getInventoryItem)
			protected.PUT("/inventory-items/:id", updateInventoryItem)
			protected.POST("/inventory-items", createInventoryItem)
		}
	}

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 VYOM ERP Backend starting on http://localhost:%s\n", port)
	log.Printf("📊 Database: vyom_lms (230 models ready)\n")
	log.Printf("🔗 Health Check: http://localhost:%s/health\n", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// ============================================================================
// MIDDLEWARE
// ============================================================================

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation with Prisma
		// For now, pass all requests
		c.Next()
	}
}

// ============================================================================
// HEALTH & AUTH HANDLERS
// ============================================================================

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "VYOM ERP Backend is running",
		"version": "1.0.0",
		"time":    time.Now(),
	})
}

func login(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate credentials using Prisma client
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  "generated_token",
		"refreshToken": "refresh_token",
		"expiresIn":    3600,
	})
}

func register(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Create user using Prisma client
	c.JSON(http.StatusCreated, req)
}

func refresh(c *gin.Context) {
	// TODO: Refresh token using Prisma client
	c.JSON(http.StatusOK, gin.H{"accessToken": "new_token", "expiresIn": 3600})
}

// ============================================================================
// TENANT HANDLERS
// ============================================================================

func createTenant(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.Tenant.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listTenants(c *gin.Context) {
	// TODO: Query using prisma.Tenant.FindMany()
	c.JSON(http.StatusOK, gin.H{"tenants": []interface{}{}})
}

func getTenant(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.Tenant.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateTenant(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.Tenant.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ============================================================================
// USER HANDLERS
// ============================================================================

func createUser(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.User.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listUsers(c *gin.Context) {
	// TODO: Query using prisma.User.FindMany()
	c.JSON(http.StatusOK, gin.H{"users": []interface{}{}})
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.User.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.User.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func deleteUser(c *gin.Context) {
	_ = c.Param("id") // TODO: Use id with prisma.User.Delete(id)
	c.JSON(http.StatusNoContent, nil)
}

// ============================================================================
// CALL CENTER HANDLERS
// ============================================================================

func createCallCenter(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.CallCenter.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listCallCenters(c *gin.Context) {
	// TODO: Query using prisma.CallCenter.FindMany()
	c.JSON(http.StatusOK, gin.H{"callCenters": []interface{}{}})
}

func getCallCenter(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.CallCenter.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateCallCenter(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.CallCenter.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ============================================================================
// AGENT HANDLERS
// ============================================================================

func createAgent(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.Agent.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listAgents(c *gin.Context) {
	// TODO: Query using prisma.Agent.FindMany()
	c.JSON(http.StatusOK, gin.H{"agents": []interface{}{}})
}

func getAgent(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.Agent.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateAgent(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.Agent.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ============================================================================
// CALL HANDLERS
// ============================================================================

func createCall(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.Call.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listCalls(c *gin.Context) {
	// TODO: Query using prisma.Call.FindMany()
	c.JSON(http.StatusOK, gin.H{"calls": []interface{}{}})
}

func getCall(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.Call.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateCall(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.Call.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ============================================================================
// CAMPAIGN HANDLERS
// ============================================================================

func createCampaign(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.Campaign.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listCampaigns(c *gin.Context) {
	// TODO: Query using prisma.Campaign.FindMany()
	c.JSON(http.StatusOK, gin.H{"campaigns": []interface{}{}})
}

func getCampaign(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.Campaign.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateCampaign(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.Campaign.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ============================================================================
// SALES LEAD HANDLERS
// ============================================================================

func createSalesLead(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.SalesLead.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listSalesLeads(c *gin.Context) {
	// TODO: Query using prisma.SalesLead.FindMany()
	c.JSON(http.StatusOK, gin.H{"salesLeads": []interface{}{}})
}

func getSalesLead(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.SalesLead.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateSalesLead(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.SalesLead.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// ============================================================================
// INVENTORY HANDLERS
// ============================================================================

func createInventoryItem(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Create using prisma.InventoryItem.CreateOne()
	c.JSON(http.StatusCreated, req)
}

func listInventoryItems(c *gin.Context) {
	// TODO: Query using prisma.InventoryItem.FindMany()
	c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
}

func getInventoryItem(c *gin.Context) {
	id := c.Param("id")
	// TODO: Query using prisma.InventoryItem.FindUnique(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func updateInventoryItem(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Update using prisma.InventoryItem.Update(id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}
