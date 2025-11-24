package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// TeamCompetition represents a team-based competition
type TeamCompetition struct {
	ID          int64     `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Objective   string    `json:"objective"` // most_calls, highest_conversion, best_engagement
	Prize       string    `json:"prize"`
	Status      string    `json:"status"` // pending, active, completed
	CreatedAt   time.Time `json:"created_at"`
}

// Challenge represents a challenge users can complete
type Challenge struct {
	ID            int64     `json:"id"`
	TenantID      string    `json:"tenant_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Duration      string    `json:"duration"` // weekly, monthly, daily
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Objective     string    `json:"objective"` // calls_made, leads_converted, etc
	TargetValue   int       `json:"target_value"`
	RewardPoints  int       `json:"reward_points"`
	RewardBadgeID *int64    `json:"reward_badge_id,omitempty"`
	Difficulty    string    `json:"difficulty"` // easy, medium, hard
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

// AchievementTier represents achievement tiers
type AchievementTier struct {
	ID        int64                  `json:"id"`
	TenantID  string                 `json:"tenant_id"`
	Name      string                 `json:"name"`
	Level     int                    `json:"level"`
	MinPoints int                    `json:"min_points"`
	MaxPoints int                    `json:"max_points"`
	Rewards   map[string]interface{} `json:"rewards"`
	Badge     string                 `json:"badge"`
	CreatedAt time.Time              `json:"created_at"`
}

// Reward represents a redeemable reward
type Reward struct {
	ID          int64     `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PointsCost  int       `json:"points_cost"`
	Category    string    `json:"category"` // discount, gift, recognition
	Stock       int       `json:"stock"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserRedemption tracks user reward redemptions
type UserRedemption struct {
	ID         int64      `json:"id"`
	TenantID   string     `json:"tenant_id"`
	UserID     int64      `json:"user_id"`
	RewardID   int64      `json:"reward_id"`
	Status     string     `json:"status"` // pending, approved, redeemed
	RedeemedAt *time.Time `json:"redeemed_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// AdvancedGamificationService handles advanced gamification features
type AdvancedGamificationService struct {
	db *sql.DB
}

// NewAdvancedGamificationService creates a new AdvancedGamificationService
func NewAdvancedGamificationService(db *sql.DB) *AdvancedGamificationService {
	return &AdvancedGamificationService{
		db: db,
	}
}

// CreateTeamCompetition creates a new team competition
func (ags *AdvancedGamificationService) CreateTeamCompetition(ctx context.Context, competition *TeamCompetition) error {
	query := `
		INSERT INTO team_competition (tenant_id, name, description, start_date, end_date, objective, prize, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := ags.db.ExecContext(ctx, query,
		competition.TenantID, competition.Name, competition.Description,
		competition.StartDate, competition.EndDate, competition.Objective,
		competition.Prize, "pending")
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	competition.ID = id
	return err
}

// GetTeamLeaderboard retrieves team leaderboard for a competition
func (ags *AdvancedGamificationService) GetTeamLeaderboard(ctx context.Context, tenantID string, competitionID int64) ([]map[string]interface{}, error) {
	leaderboard := make([]map[string]interface{}, 0)

	// This would aggregate points/metrics by team based on the competition objective
	// For now, return structure
	return leaderboard, nil
}

// CreateChallenge creates a new challenge
func (ags *AdvancedGamificationService) CreateChallenge(ctx context.Context, challenge *Challenge) error {
	query := `
		INSERT INTO challenge (tenant_id, title, description, duration, start_date, end_date, objective, target_value, reward_points, reward_badge_id, difficulty, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := ags.db.ExecContext(ctx, query,
		challenge.TenantID, challenge.Title, challenge.Description, challenge.Duration,
		challenge.StartDate, challenge.EndDate, challenge.Objective, challenge.TargetValue,
		challenge.RewardPoints, challenge.RewardBadgeID, challenge.Difficulty, challenge.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	challenge.ID = id
	return err
}

// GetActiveChallenges retrieves active challenges for a tenant
func (ags *AdvancedGamificationService) GetActiveChallenges(ctx context.Context, tenantID string) ([]Challenge, error) {
	challenges := make([]Challenge, 0)

	query := `
		SELECT id, tenant_id, title, description, duration, start_date, end_date, objective, target_value, reward_points, reward_badge_id, difficulty, is_active, created_at
		FROM challenge
		WHERE tenant_id = ? AND is_active = true AND end_date > NOW()
		ORDER BY start_date DESC
	`

	rows, err := ags.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ch Challenge
		var badgeID sql.NullInt64

		if err := rows.Scan(&ch.ID, &ch.TenantID, &ch.Title, &ch.Description, &ch.Duration,
			&ch.StartDate, &ch.EndDate, &ch.Objective, &ch.TargetValue, &ch.RewardPoints,
			&badgeID, &ch.Difficulty, &ch.IsActive, &ch.CreatedAt); err == nil {
			if badgeID.Valid {
				ch.RewardBadgeID = &badgeID.Int64
			}
			challenges = append(challenges, ch)
		}
	}

	return challenges, nil
}

// CheckChallengeCompletion checks if a user has completed a challenge
func (ags *AdvancedGamificationService) CheckChallengeCompletion(ctx context.Context, userID int64, challengeID int64) (bool, error) {
	// Implementation would check user metrics against challenge objectives
	// This is a placeholder that returns false
	return false, nil
}

// CreateAchievementTier creates a new achievement tier
func (ags *AdvancedGamificationService) CreateAchievementTier(ctx context.Context, tier *AchievementTier) error {
	query := `
		INSERT INTO achievement_tier (tenant_id, name, level, min_points, max_points, rewards, badge, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := ags.db.ExecContext(ctx, query,
		tier.TenantID, tier.Name, tier.Level, tier.MinPoints, tier.MaxPoints,
		fmt.Sprintf("%v", tier.Rewards), tier.Badge)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	tier.ID = id
	return err
}

// GetUserAchievementTier gets the current tier for a user
func (ags *AdvancedGamificationService) GetUserAchievementTier(ctx context.Context, tenantID string, userID int64) (*AchievementTier, error) {
	// Get user's current points
	query := `
		SELECT total_points FROM user_gamification WHERE user_id = ? AND tenant_id = ?
	`
	var totalPoints int
	err := ags.db.QueryRowContext(ctx, query, userID, tenantID).Scan(&totalPoints)
	if err != nil {
		return nil, err
	}

	// Find tier based on points
	tierQuery := `
		SELECT id, tenant_id, name, level, min_points, max_points, rewards, badge, created_at
		FROM achievement_tier
		WHERE tenant_id = ? AND min_points <= ? AND max_points >= ?
		LIMIT 1
	`

	tier := &AchievementTier{}
	var rewardsStr string

	err = ags.db.QueryRowContext(ctx, tierQuery, tenantID, totalPoints, totalPoints).
		Scan(&tier.ID, &tier.TenantID, &tier.Name, &tier.Level, &tier.MinPoints, &tier.MaxPoints,
			&rewardsStr, &tier.Badge, &tier.CreatedAt)

	if err != nil {
		return nil, err
	}

	tier.Rewards = make(map[string]interface{})
	return tier, nil
}

// CreateReward creates a new redeemable reward
func (ags *AdvancedGamificationService) CreateReward(ctx context.Context, reward *Reward) error {
	query := `
		INSERT INTO reward (tenant_id, name, description, points_cost, category, stock, is_available, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := ags.db.ExecContext(ctx, query,
		reward.TenantID, reward.Name, reward.Description, reward.PointsCost,
		reward.Category, reward.Stock, reward.IsAvailable)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	reward.ID = id
	return err
}

// GetAvailableRewards retrieves available rewards for a tenant
func (ags *AdvancedGamificationService) GetAvailableRewards(ctx context.Context, tenantID string) ([]Reward, error) {
	rewards := make([]Reward, 0)

	query := `
		SELECT id, tenant_id, name, description, points_cost, category, stock, is_available, created_at
		FROM reward
		WHERE tenant_id = ? AND is_available = true AND stock > 0
		ORDER BY points_cost ASC
	`

	rows, err := ags.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Reward
		if err := rows.Scan(&r.ID, &r.TenantID, &r.Name, &r.Description, &r.PointsCost,
			&r.Category, &r.Stock, &r.IsAvailable, &r.CreatedAt); err == nil {
			rewards = append(rewards, r)
		}
	}

	return rewards, nil
}

// RedeemReward redeems a reward for a user
func (ags *AdvancedGamificationService) RedeemReward(ctx context.Context, tenantID string, userID int64, rewardID int64) error {
	// Check user has enough points
	userQuery := `SELECT total_points FROM user_gamification WHERE user_id = ? AND tenant_id = ?`
	var userPoints int
	err := ags.db.QueryRowContext(ctx, userQuery, userID, tenantID).Scan(&userPoints)
	if err != nil {
		return err
	}

	// Get reward cost
	rewardQuery := `SELECT points_cost FROM reward WHERE id = ? AND tenant_id = ?`
	var cost int
	err = ags.db.QueryRowContext(ctx, rewardQuery, rewardID, tenantID).Scan(&cost)
	if err != nil {
		return err
	}

	if userPoints < cost {
		return fmt.Errorf("insufficient points")
	}

	// Insert redemption record
	query := `
		INSERT INTO user_redemption (tenant_id, user_id, reward_id, status, created_at)
		VALUES (?, ?, ?, 'pending', NOW())
	`

	_, err = ags.db.ExecContext(ctx, query, tenantID, userID, rewardID)
	if err != nil {
		return err
	}

	// Deduct points from user
	updateQuery := `
		UPDATE user_gamification
		SET total_points = total_points - ?
		WHERE user_id = ? AND tenant_id = ?
	`
	_, err = ags.db.ExecContext(ctx, updateQuery, cost, userID, tenantID)
	return err
}

// GetUserRedemptions retrieves user's reward redemptions
func (ags *AdvancedGamificationService) GetUserRedemptions(ctx context.Context, tenantID string, userID int64) ([]UserRedemption, error) {
	redemptions := make([]UserRedemption, 0)

	query := `
		SELECT id, tenant_id, user_id, reward_id, status, redeemed_at, created_at
		FROM user_redemption
		WHERE tenant_id = ? AND user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := ags.db.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ur UserRedemption
		var redeemedAt sql.NullTime

		if err := rows.Scan(&ur.ID, &ur.TenantID, &ur.UserID, &ur.RewardID, &ur.Status, &redeemedAt, &ur.CreatedAt); err == nil {
			if redeemedAt.Valid {
				ur.RedeemedAt = &redeemedAt.Time
			}
			redemptions = append(redemptions, ur)
		}
	}

	return redemptions, nil
}

// GetGamificationLeaderboard retrieves the full leaderboard for a tenant
func (ags *AdvancedGamificationService) GetGamificationLeaderboard(ctx context.Context, tenantID string, limit int) ([]map[string]interface{}, error) {
	leaderboard := make([]map[string]interface{}, 0)

	query := `
		SELECT ug.user_id, ug.total_points, COUNT(DISTINCT ub.badge_id) as badge_count
		FROM user_gamification ug
		LEFT JOIN user_badge ub ON ug.user_id = ub.user_id
		WHERE ug.tenant_id = ?
		GROUP BY ug.user_id
		ORDER BY ug.total_points DESC
		LIMIT ?
	`

	rows, err := ags.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rank := 1
	for rows.Next() {
		var userID int64
		var points, badgeCount int

		if err := rows.Scan(&userID, &points, &badgeCount); err == nil {
			entry := map[string]interface{}{
				"rank":        rank,
				"user_id":     userID,
				"points":      points,
				"badge_count": badgeCount,
			}
			leaderboard = append(leaderboard, entry)
			rank++
		}
	}

	return leaderboard, nil
}

// GetGamificationStats returns gamification statistics for a tenant
func (ags *AdvancedGamificationService) GetGamificationStats(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total users engaged
	usersQuery := `SELECT COUNT(*) FROM user_gamification WHERE tenant_id = ? AND total_points > 0`
	var engagedUsers int
	ags.db.QueryRowContext(ctx, usersQuery, tenantID).Scan(&engagedUsers)
	stats["engaged_users"] = engagedUsers

	// Active challenges
	challengesQuery := `SELECT COUNT(*) FROM challenge WHERE tenant_id = ? AND is_active = true`
	var activeChallenges int
	ags.db.QueryRowContext(ctx, challengesQuery, tenantID).Scan(&activeChallenges)
	stats["active_challenges"] = activeChallenges

	// Available rewards
	rewardsQuery := `SELECT COUNT(*) FROM reward WHERE tenant_id = ? AND is_available = true AND stock > 0`
	var availableRewards int
	ags.db.QueryRowContext(ctx, rewardsQuery, tenantID).Scan(&availableRewards)
	stats["available_rewards"] = availableRewards

	// Avg points per user
	avgQuery := `SELECT AVG(total_points) FROM user_gamification WHERE tenant_id = ?`
	var avgPoints sql.NullFloat64
	ags.db.QueryRowContext(ctx, avgQuery, tenantID).Scan(&avgPoints)
	stats["avg_points_per_user"] = avgPoints.Float64

	return stats, nil
}
