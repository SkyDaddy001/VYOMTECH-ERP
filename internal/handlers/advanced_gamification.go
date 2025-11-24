package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"multi-tenant-ai-callcenter/internal/services"
	"multi-tenant-ai-callcenter/pkg/logger"
)

// AdvancedGamificationHandler handles advanced gamification requests
type AdvancedGamificationHandler struct {
	advGamService *services.AdvancedGamificationService
	logger        *logger.Logger
}

// NewAdvancedGamificationHandler creates a new AdvancedGamificationHandler
func NewAdvancedGamificationHandler(advGamService *services.AdvancedGamificationService, logger *logger.Logger) *AdvancedGamificationHandler {
	return &AdvancedGamificationHandler{
		advGamService: advGamService,
		logger:        logger,
	}
}

// CreateTeamCompetitionRequest is the request body for creating a team competition
type CreateTeamCompetitionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"` // ISO8601
	EndDate     string `json:"end_date"`   // ISO8601
	Objective   string `json:"objective"`  // most_calls, highest_conversion, best_engagement
	Prize       string `json:"prize"`
}

// CreateAdvancedChallengeRequest is the request body for creating an advanced challenge
type CreateAdvancedChallengeRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Duration      string `json:"duration"`   // weekly, monthly, daily
	StartDate     string `json:"start_date"` // ISO8601
	EndDate       string `json:"end_date"`
	Objective     string `json:"objective"` // calls_made, leads_converted
	TargetValue   int    `json:"target_value"`
	RewardPoints  int    `json:"reward_points"`
	RewardBadgeID *int64 `json:"reward_badge_id,omitempty"`
	Difficulty    string `json:"difficulty"` // easy, medium, hard
}

// CreateRewardRequest is the request body for creating a reward
type CreateRewardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PointsCost  int    `json:"points_cost"`
	Category    string `json:"category"` // discount, gift, recognition
	Stock       int    `json:"stock"`
}

// CreateTeamCompetition creates a new team competition
// POST /api/v1/gamification/competitions
func (agh *AdvancedGamificationHandler) CreateTeamCompetition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateTeamCompetitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	competition := &services.TeamCompetition{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		Objective:   req.Objective,
		Prize:       req.Prize,
	}

	err = agh.advGamService.CreateTeamCompetition(ctx, competition)
	if err != nil {
		agh.logger.Error("Failed to create competition", "error", err)
		http.Error(w, "failed to create competition", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(competition)
}

// GetTeamLeaderboard retrieves leaderboard for a competition
// GET /api/v1/gamification/competitions/{id}/leaderboard
func (agh *AdvancedGamificationHandler) GetTeamLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	compIDStr := r.URL.Query().Get("id")
	if compIDStr == "" {
		http.Error(w, "competition id required", http.StatusBadRequest)
		return
	}

	compID, err := strconv.ParseInt(compIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid competition id", http.StatusBadRequest)
		return
	}

	leaderboard, err := agh.advGamService.GetTeamLeaderboard(ctx, tenantID, compID)
	if err != nil {
		agh.logger.Error("Failed to get leaderboard", "error", err)
		http.Error(w, "failed to get leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"competition_id": compID,
		"leaderboard":    leaderboard,
	})
}

// CreateChallenge creates a new challenge
// POST /api/v1/gamification/challenges
func (agh *AdvancedGamificationHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateAdvancedChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format", http.StatusBadRequest)
		return
	}

	challenge := &services.Challenge{
		TenantID:      tenantID,
		Title:         req.Title,
		Description:   req.Description,
		Duration:      req.Duration,
		StartDate:     startDate,
		EndDate:       endDate,
		Objective:     req.Objective,
		TargetValue:   req.TargetValue,
		RewardPoints:  req.RewardPoints,
		RewardBadgeID: req.RewardBadgeID,
		Difficulty:    req.Difficulty,
		IsActive:      true,
	}

	err = agh.advGamService.CreateChallenge(ctx, challenge)
	if err != nil {
		agh.logger.Error("Failed to create challenge", "error", err)
		http.Error(w, "failed to create challenge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(challenge)
}

// GetActiveChallenges retrieves all active challenges
// GET /api/v1/gamification/challenges/active
func (agh *AdvancedGamificationHandler) GetActiveChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	challenges, err := agh.advGamService.GetActiveChallenges(ctx, tenantID)
	if err != nil {
		agh.logger.Error("Failed to get challenges", "error", err)
		http.Error(w, "failed to get challenges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":      len(challenges),
		"challenges": challenges,
	})
}

// CreateReward creates a new redeemable reward
// POST /api/v1/gamification/rewards
func (agh *AdvancedGamificationHandler) CreateReward(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	var req CreateRewardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	reward := &services.Reward{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		PointsCost:  req.PointsCost,
		Category:    req.Category,
		Stock:       req.Stock,
		IsAvailable: true,
	}

	err := agh.advGamService.CreateReward(ctx, reward)
	if err != nil {
		agh.logger.Error("Failed to create reward", "error", err)
		http.Error(w, "failed to create reward", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reward)
}

// GetAvailableRewards retrieves available rewards
// GET /api/v1/gamification/rewards
func (agh *AdvancedGamificationHandler) GetAvailableRewards(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	rewards, err := agh.advGamService.GetAvailableRewards(ctx, tenantID)
	if err != nil {
		agh.logger.Error("Failed to get rewards", "error", err)
		http.Error(w, "failed to get rewards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":   len(rewards),
		"rewards": rewards,
	})
}

// RedeemReward redeems a reward for the user
// POST /api/v1/gamification/redeem
func (agh *AdvancedGamificationHandler) RedeemReward(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		agh.logger.Error("Failed to extract user ID from context")
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	rewardIDStr := r.URL.Query().Get("reward_id")
	if rewardIDStr == "" {
		http.Error(w, "reward_id required", http.StatusBadRequest)
		return
	}

	rewardID, err := strconv.ParseInt(rewardIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid reward_id", http.StatusBadRequest)
		return
	}

	err = agh.advGamService.RedeemReward(ctx, tenantID, userID, rewardID)
	if err != nil {
		agh.logger.Error("Failed to redeem reward", "error", err)
		http.Error(w, "failed to redeem reward", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Reward redeemed successfully",
	})
}

// GetGamificationLeaderboard retrieves the gamification leaderboard
// GET /api/v1/gamification/leaderboard?limit=50
func (agh *AdvancedGamificationHandler) GetGamificationLeaderboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	leaderboard, err := agh.advGamService.GetGamificationLeaderboard(ctx, tenantID, limit)
	if err != nil {
		agh.logger.Error("Failed to get leaderboard", "error", err)
		http.Error(w, "failed to get leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"limit":       limit,
		"total":       len(leaderboard),
		"leaderboard": leaderboard,
	})
}

// GetGamificationStats retrieves gamification statistics
// GET /api/v1/gamification/stats
func (agh *AdvancedGamificationHandler) GetGamificationStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, ok := ctx.Value("tenantID").(string)
	if !ok {
		agh.logger.Error("Failed to extract tenant ID from context")
		http.Error(w, "tenant id not found", http.StatusBadRequest)
		return
	}

	stats, err := agh.advGamService.GetGamificationStats(ctx, tenantID)
	if err != nil {
		agh.logger.Error("Failed to get stats", "error", err)
		http.Error(w, "failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
