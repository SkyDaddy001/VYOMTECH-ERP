package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"multi-tenant-ai-callcenter/internal/models"
	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// GamificationHandler handles all gamification-related endpoints
type GamificationHandler struct {
	gamificationService *services.GamificationService
	logger              *logger.Logger
}

// NewGamificationHandler creates a new GamificationHandler
func NewGamificationHandler(gamificationService *services.GamificationService, logger *logger.Logger) *GamificationHandler {
	return &GamificationHandler{
		gamificationService: gamificationService,
		logger:              logger,
	}
}

// AwardPointsRequest is the request body for awarding points
type AwardPointsRequest struct {
	ActionType  string `json:"actionType"`
	Points      int    `json:"points"`
	Description string `json:"description"`
	BonusReason string `json:"bonusReason,omitempty"`
}

// RevokePointsRequest is the request body for revoking points
type RevokePointsRequest struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

// AwardBadgeRequest is the request body for awarding a badge
type AwardBadgeRequest struct {
	BadgeID int64 `json:"badgeId"`
}

// CreateBadgeRequest is the request body for creating a badge
type CreateBadgeRequest struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	IconURL          string `json:"iconUrl"`
	Category         string `json:"category"`
	RequirementType  string `json:"requirementType"`
	RequirementValue int    `json:"requirementValue"`
	PointsReward     int    `json:"pointsReward"`
	Rarity           string `json:"rarity"`
}

// CreateChallengeRequest is the request body for creating a challenge
type CreateChallengeRequest struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Type             string  `json:"type"`
	ObjectiveType    string  `json:"objectiveType"`
	ObjectiveValue   int     `json:"objectiveValue"`
	PointsReward     int     `json:"pointsReward"`
	DurationDays     int     `json:"durationDays"`
	MaxParticipants  int     `json:"maxParticipants"`
	StartDate        string  `json:"startDate"`
	EndDate          string  `json:"endDate"`
	ResetFrequency   string  `json:"resetFrequency"`
	CompletionBonus  int     `json:"completionBonus"`
	PointsMultiplier float64 `json:"pointsMultiplier"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse is a standard success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetUserPoints retrieves the current user's points
// GET /api/v1/gamification/points
func (h *GamificationHandler) GetUserPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userPoints, err := h.gamificationService.GetUserPoints(ctx, userID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get user points", "userID", userID, "error", err)
		http.Error(w, "failed to get user points", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userPoints)
}

// AwardPoints awards points to the user
// POST /api/v1/gamification/points/award
func (h *GamificationHandler) AwardPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req AwardPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.gamificationService.AwardPoints(ctx, userID, tenantID, req.ActionType, req.Points, req.Description, req.BonusReason); err != nil {
		h.logger.Error("Failed to award points", "userID", userID, "error", err)
		http.Error(w, "failed to award points", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "points awarded successfully",
	})
}

// RevokePoints revokes points from the user
// POST /api/v1/gamification/points/revoke
func (h *GamificationHandler) RevokePoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req RevokePointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.gamificationService.RevokePoints(ctx, userID, tenantID, req.Points, req.Reason); err != nil {
		h.logger.Error("Failed to revoke points", "userID", userID, "error", err)
		http.Error(w, "failed to revoke points", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "points revoked successfully",
	})
}

// GetUserBadges retrieves all badges for the user
// GET /api/v1/gamification/badges
func (h *GamificationHandler) GetUserBadges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	badges, err := h.gamificationService.GetUserBadges(ctx, userID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get user badges", "userID", userID, "error", err)
		http.Error(w, "failed to get user badges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(badges)
}

// AwardBadge awards a badge to the user
// POST /api/v1/gamification/badges/award
func (h *GamificationHandler) AwardBadge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req AwardBadgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.gamificationService.AwardBadge(ctx, userID, req.BadgeID, tenantID); err != nil {
		h.logger.Error("Failed to award badge", "userID", userID, "badgeID", req.BadgeID, "error", err)
		http.Error(w, "failed to award badge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "badge awarded successfully",
	})
}

// CreateBadge creates a new badge
// POST /api/v1/gamification/badges
func (h *GamificationHandler) CreateBadge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateBadgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	badge := &models.Badge{
		TenantID:         tenantID,
		Code:             req.Code,
		Name:             req.Name,
		Description:      req.Description,
		IconURL:          req.IconURL,
		Category:         req.Category,
		RequirementType:  req.RequirementType,
		RequirementValue: req.RequirementValue,
		PointsReward:     req.PointsReward,
		Rarity:           req.Rarity,
		Active:           true,
	}

	if err := h.gamificationService.CreateBadge(ctx, badge); err != nil {
		h.logger.Error("Failed to create badge", "error", err)
		http.Error(w, "failed to create badge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "badge created successfully",
		Data:    badge,
	})
}

// GetUserChallenges retrieves all challenges for the user
// GET /api/v1/gamification/challenges
func (h *GamificationHandler) GetUserChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	challenges, err := h.gamificationService.GetUserChallenges(ctx, userID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get user challenges", "userID", userID, "error", err)
		http.Error(w, "failed to get user challenges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

// GetActiveChallenges retrieves all active challenges
// GET /api/v1/gamification/challenges/active
func (h *GamificationHandler) GetActiveChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	challenges, err := h.gamificationService.GetActiveChallenges(ctx, tenantID)
	if err != nil {
		h.logger.Error("Failed to get active challenges", "error", err)
		http.Error(w, "failed to get active challenges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

// CreateChallenge creates a new challenge
// POST /api/v1/gamification/challenges
func (h *GamificationHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		http.Error(w, "invalid start date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		http.Error(w, "invalid end date format", http.StatusBadRequest)
		return
	}

	challenge := &models.Challenge{
		TenantID:        tenantID,
		Name:            req.Name,
		Description:     req.Description,
		ChallengeType:   req.Type,
		Status:          "active",
		ObjectiveType:   req.ObjectiveType,
		ObjectiveTarget: req.ObjectiveValue,
		PointsReward:    req.PointsReward,
		StartDate:       startDate,
		EndDate:         endDate,
	}

	if err := h.gamificationService.CreateChallenge(ctx, challenge); err != nil {
		h.logger.Error("Failed to create challenge", "error", err)
		http.Error(w, "failed to create challenge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "challenge created successfully",
		Data:    challenge,
	})
}

// GetLeaderboard retrieves the leaderboard
// GET /api/v1/gamification/leaderboard?period=weekly&limit=100
func (h *GamificationHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	// Get query parameters
	periodType := r.URL.Query().Get("period")
	if periodType == "" {
		periodType = "weekly"
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	leaderboard, err := h.gamificationService.GetLeaderboard(ctx, tenantID, periodType, limit)
	if err != nil {
		h.logger.Error("Failed to get leaderboard", "error", err)
		http.Error(w, "failed to get leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

// GetGamificationProfile retrieves the user's gamification profile
// GET /api/v1/gamification/profile
func (h *GamificationHandler) GetGamificationProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		h.logger.Error("Failed to extract user ID from context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		h.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	profile, err := h.gamificationService.GetGamificationProfile(ctx, userID, tenantID)
	if err != nil {
		h.logger.Error("Failed to get gamification profile", "userID", userID, "error", err)
		http.Error(w, "failed to get gamification profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
