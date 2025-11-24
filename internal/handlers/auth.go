package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"multi-tenant-ai-callcenter/internal/middleware"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

type AuthHandler struct {
	authService *services.AuthService
	logger      *logger.Logger
}

func NewAuthHandler(authService *services.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterRequest defines the registration request structure
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"omitempty"`
	Role     string `json:"role" binding:"omitempty"`
	TenantID string `json:"tenant_id" binding:"omitempty"`
}

// LoginRequest defines the login request structure
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ChangePasswordRequest defines the change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// AuthResponse defines the authentication response
type AuthResponse struct {
	Token   string    `json:"token"`
	User    *UserInfo `json:"user"`
	Message string    `json:"message,omitempty"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.authService.Register(ctx, req.Email, req.Password, req.Role, req.TenantID)
	if err != nil {
		h.logger.Warn("Registration failed", "error", err)
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email, user.Role, user.TenantID)
	if err != nil {
		h.logger.Error("Failed to generate token", "error", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
		Message: "User registered successfully",
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Login failed", "error", err, "email", req.Email)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get user details for response
	user, err := h.authService.ValidateToken(token)
	if err != nil {
		h.logger.Error("Failed to validate token", "error", err)
		http.Error(w, "Failed to validate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User: &UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
		Message: "Login successful",
	})
}

// ChangePassword handles password changes for authenticated users
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.authService.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword); err != nil {
		h.logger.Warn("Failed to change password", "error", err, "user_id", userID)
		http.Error(w, "Failed to change password", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserIDKey)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid": true,
		"user":  user,
	})
}

// AgentHandler handles agent-related requests
type AgentHandler struct {
	agentService *services.AgentService
	logger       *logger.Logger
}

func NewAgentHandler(agentService *services.AgentService, logger *logger.Logger) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
		logger:       logger,
	}
}

// GetAgent retrieves a specific agent
func (h *AgentHandler) GetAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid agent ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	agent, err := h.agentService.GetAgent(ctx, agentID)
	if err != nil {
		h.logger.Warn("Failed to get agent", "error", err, "agent_id", agentID)
		http.Error(w, "Agent not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(agent)
}

// GetAgentsByTenant retrieves all agents for a tenant
func (h *AgentHandler) GetAgentsByTenant(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Missing tenant context", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	agents, err := h.agentService.GetAgentsByTenant(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to get agents", "error", err, "tenant_id", tenantID)
		http.Error(w, "Failed to retrieve agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"agents": agents,
		"count":  len(agents),
	})
}

// UpdateAgentAvailability updates agent availability status
func (h *AgentHandler) UpdateAgentAvailability(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	availability, ok := req["availability"]
	if !ok || availability == "" {
		http.Error(w, "Missing availability field", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.agentService.UpdateAgentAvailability(ctx, userID, availability); err != nil {
		h.logger.Error("Failed to update availability", "error", err)
		http.Error(w, "Failed to update availability", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Availability updated successfully",
	})
}

// GetAvailableAgents retrieves available agents for a tenant
func (h *AgentHandler) GetAvailableAgents(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Missing tenant context", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	agents, err := h.agentService.GetAvailableAgents(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to get available agents", "error", err)
		http.Error(w, "Failed to retrieve available agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"available_agents": agents,
		"count":            len(agents),
	})
}

// GetAgentStats retrieves agent statistics for a tenant
func (h *AgentHandler) GetAgentStats(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := r.Context().Value(middleware.TenantIDKey).(string)
	if !ok {
		http.Error(w, "Missing tenant context", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	stats, err := h.agentService.GetAgentStats(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to get agent stats", "error", err)
		http.Error(w, "Failed to retrieve agent stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
