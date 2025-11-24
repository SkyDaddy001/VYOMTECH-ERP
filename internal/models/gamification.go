package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ==================== GAMIFICATION CONFIG ====================

// GamificationConfig holds tenant-wide gamification settings
type GamificationConfig struct {
	ID                     int64     `json:"id" db:"id"`
	TenantID               string    `json:"tenant_id" db:"tenant_id"`
	Enabled                bool      `json:"enabled" db:"enabled"`
	PointsPerCall          int       `json:"points_per_call" db:"points_per_call"`
	PointsPerConversion    int       `json:"points_per_conversion" db:"points_per_conversion"`
	PointsPerQualityReview int       `json:"points_per_quality_review" db:"points_per_quality_review"`
	PointsPerFeedback      int       `json:"points_per_feedback" db:"points_per_feedback"`
	PointsDecayPercent     int       `json:"points_decay_percent" db:"points_decay_percent"`
	DecayPeriodDays        int       `json:"decay_period_days" db:"decay_period_days"`
	MaxDailyPoints         int       `json:"max_daily_points" db:"max_daily_points"`
	LeaderboardResetPeriod string    `json:"leaderboard_reset_period" db:"leaderboard_reset_period"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// ==================== POINTS SYSTEM ====================

// UserPoints tracks user's current and lifetime points
type UserPoints struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	CurrentPoints   int       `json:"current_points" db:"current_points"`
	LifetimePoints  int       `json:"lifetime_points" db:"lifetime_points"`
	PeriodStartDate time.Time `json:"period_start_date" db:"period_start_date"`
	PeriodEndDate   time.Time `json:"period_end_date" db:"period_end_date"`
	DailyPoints     int       `json:"daily_points" db:"daily_points"`
	DailyDate       time.Time `json:"daily_date" db:"daily_date"`
	StreakDays      int       `json:"streak_days" db:"streak_days"`
	LastActionDate  time.Time `json:"last_action_date" db:"last_action_date"`
	Rank            int       `json:"rank" db:"rank"`
	RankUpdatedAt   time.Time `json:"rank_updated_at" db:"rank_updated_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// PointTransaction logs individual point award/revoke events
type PointTransaction struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	TenantID    string    `json:"tenant_id" db:"tenant_id"`
	Points      int       `json:"points" db:"points"`
	ActionType  string    `json:"action_type" db:"action_type"`
	ActionID    string    `json:"action_id" db:"action_id"`
	Description string    `json:"description" db:"description"`
	Multiplier  float64   `json:"multiplier" db:"multiplier"`
	BonusReason string    `json:"bonus_reason" db:"bonus_reason"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ==================== BADGES & ACHIEVEMENTS ====================

// Badge represents an achievement badge
type Badge struct {
	ID               int64     `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	Code             string    `json:"code" db:"code"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	IconURL          string    `json:"icon_url" db:"icon_url"`
	Category         string    `json:"category" db:"category"` // milestone, skill, behavior, performance, social
	RequirementType  string    `json:"requirement_type" db:"requirement_type"`
	RequirementValue int       `json:"requirement_value" db:"requirement_value"`
	PointsReward     int       `json:"points_reward" db:"points_reward"`
	Rarity           string    `json:"rarity" db:"rarity"` // common, uncommon, rare, epic, legendary
	Active           bool      `json:"active" db:"active"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserBadge represents a badge earned by a user
type UserBadge struct {
	ID              int64     `json:"id" db:"id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	BadgeID         int64     `json:"badge_id" db:"badge_id"`
	TenantID        string    `json:"tenant_id" db:"tenant_id"`
	EarnedDate      time.Time `json:"earned_date" db:"earned_date"`
	ProgressPercent int       `json:"progress_percent" db:"progress_percent"`
	Notified        bool      `json:"notified" db:"notified"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	Badge           *Badge    `json:"badge,omitempty" db:"-"`
}

// ==================== CHALLENGES ====================

// Challenge represents a challenge/quest
type Challenge struct {
	ID               int64     `json:"id" db:"id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	ChallengeType    string    `json:"challenge_type" db:"challenge_type"` // daily, weekly, monthly, seasonal, special
	Status           string    `json:"status" db:"status"`                 // active, inactive, completed
	ObjectiveType    string    `json:"objective_type" db:"objective_type"`
	ObjectiveTarget  int       `json:"objective_target" db:"objective_target"`
	ObjectiveCurrent int       `json:"objective_current" db:"objective_current"`
	PointsReward     int       `json:"points_reward" db:"points_reward"`
	BadgeRewardID    *int64    `json:"badge_reward_id" db:"badge_reward_id"`
	StartDate        time.Time `json:"start_date" db:"start_date"`
	EndDate          time.Time `json:"end_date" db:"end_date"`
	Difficulty       string    `json:"difficulty" db:"difficulty"` // easy, medium, hard, extreme
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserChallenge tracks user progress on a challenge
type UserChallenge struct {
	ID            int64      `json:"id" db:"id"`
	UserID        int64      `json:"user_id" db:"user_id"`
	ChallengeID   int64      `json:"challenge_id" db:"challenge_id"`
	TenantID      string     `json:"tenant_id" db:"tenant_id"`
	Progress      int        `json:"progress" db:"progress"`
	Completed     bool       `json:"completed" db:"completed"`
	CompletedDate *time.Time `json:"completed_date" db:"completed_date"`
	PointsEarned  int        `json:"points_earned" db:"points_earned"`
	BadgeEarnedID *int64     `json:"badge_earned_id" db:"badge_earned_id"`
	Status        string     `json:"status" db:"status"` // in_progress, completed, failed, abandoned
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	Challenge     *Challenge `json:"challenge,omitempty" db:"-"`
}

// ==================== LEADERBOARDS ====================

// LeaderboardEntry represents a user's position on leaderboard
type LeaderboardEntry struct {
	ID           int64     `json:"id" db:"id"`
	TenantID     string    `json:"tenant_id" db:"tenant_id"`
	PeriodType   string    `json:"period_type" db:"period_type"` // daily, weekly, monthly, all_time
	PeriodDate   time.Time `json:"period_date" db:"period_date"`
	UserID       int64     `json:"user_id" db:"user_id"`
	Rank         int       `json:"rank" db:"rank"`
	Points       int       `json:"points" db:"points"`
	PointsChange int       `json:"points_change" db:"points_change"`
	PreviousRank *int      `json:"previous_rank" db:"previous_rank"`
	BadgesCount  int       `json:"badges_count" db:"badges_count"`
	StreakDays   int       `json:"streak_days" db:"streak_days"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	User         *User     `json:"user,omitempty" db:"-"`
}

// ==================== REWARDS ====================

// RewardItem represents an item users can redeem points for
type RewardItem struct {
	ID                int64             `json:"id" db:"id"`
	TenantID          string            `json:"tenant_id" db:"tenant_id"`
	Name              string            `json:"name" db:"name"`
	Description       string            `json:"description" db:"description"`
	Category          string            `json:"category" db:"category"`
	PointsCost        int               `json:"points_cost" db:"points_cost"`
	Stock             int               `json:"stock" db:"stock"` // -1 = unlimited
	StockRemaining    int               `json:"stock_remaining" db:"stock_remaining"`
	ImageURL          string            `json:"image_url" db:"image_url"`
	RedemptionType    string            `json:"redemption_type" db:"redemption_type"` // digital, physical, experience, discount
	RedemptionDetails RedemptionDetails `json:"redemption_details" db:"redemption_details"`
	Active            bool              `json:"active" db:"active"`
	Featured          bool              `json:"featured" db:"featured"`
	ExpiryDate        *time.Time        `json:"expiry_date" db:"expiry_date"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
}

// RedemptionDetails holds redemption-specific information
type RedemptionDetails struct {
	Code         string `json:"code,omitempty"`
	Link         string `json:"link,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	ValidDays    int    `json:"valid_days,omitempty"`
}

// Scan implements sql.Scanner interface
func (rd *RedemptionDetails) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), rd)
}

// Value implements driver.Valuer interface
func (rd RedemptionDetails) Value() (driver.Value, error) {
	return json.Marshal(rd)
}

// UserRedemption tracks a user's reward redemption
type UserRedemption struct {
	ID             int64       `json:"id" db:"id"`
	UserID         int64       `json:"user_id" db:"user_id"`
	RewardID       int64       `json:"reward_id" db:"reward_id"`
	TenantID       string      `json:"tenant_id" db:"tenant_id"`
	PointsSpent    int         `json:"points_spent" db:"points_spent"`
	Quantity       int         `json:"quantity" db:"quantity"`
	Status         string      `json:"status" db:"status"` // pending, approved, completed, cancelled
	RedemptionCode string      `json:"redemption_code" db:"redemption_code"`
	RedeemedDate   *time.Time  `json:"redeemed_date" db:"redeemed_date"`
	CompletedDate  *time.Time  `json:"completed_date" db:"completed_date"`
	Notes          string      `json:"notes" db:"notes"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
	Reward         *RewardItem `json:"reward,omitempty" db:"-"`
}

// ==================== LEVELS/TIERS ====================

// PlayerLevel represents a level/tier in the gamification system
type PlayerLevel struct {
	ID                  int64           `json:"id" db:"id"`
	TenantID            string          `json:"tenant_id" db:"tenant_id"`
	LevelNumber         int             `json:"level_number" db:"level_number"`
	Name                string          `json:"name" db:"name"`
	Description         string          `json:"description" db:"description"`
	PointsRequired      int             `json:"points_required" db:"points_required"`
	PointsTotalRequired int             `json:"points_total_required" db:"points_total_required"`
	IconURL             string          `json:"icon_url" db:"icon_url"`
	Benefits            string          `json:"benefits" db:"benefits"`
	UnlockPrivileges    json.RawMessage `json:"unlock_privileges" db:"unlock_privileges"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
}

// UserLevel represents a user's current level
type UserLevel struct {
	ID             int64      `json:"id" db:"id"`
	UserID         int64      `json:"user_id" db:"user_id"`
	TenantID       string     `json:"tenant_id" db:"tenant_id"`
	CurrentLevel   int        `json:"current_level" db:"current_level"`
	LifetimeLevel  int        `json:"lifetime_level" db:"lifetime_level"`
	PreviousLevel  *int       `json:"previous_level" db:"previous_level"`
	LevelUpDate    *time.Time `json:"level_up_date" db:"level_up_date"`
	ProgressToNext int        `json:"progress_to_next" db:"progress_to_next"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// ==================== TEAM CHALLENGES ====================

// TeamChallenge represents a team-based challenge
type TeamChallenge struct {
	ID                        int64     `json:"id" db:"id"`
	TenantID                  string    `json:"tenant_id" db:"tenant_id"`
	Name                      string    `json:"name" db:"name"`
	Description               string    `json:"description" db:"description"`
	TeamID                    *int64    `json:"team_id" db:"team_id"`
	ObjectiveType             string    `json:"objective_type" db:"objective_type"`
	ObjectiveTarget           int       `json:"objective_target" db:"objective_target"`
	TeamRewardPoints          int       `json:"team_reward_points" db:"team_reward_points"`
	IndividualBonusMultiplier float64   `json:"individual_bonus_multiplier" db:"individual_bonus_multiplier"`
	StartDate                 time.Time `json:"start_date" db:"start_date"`
	EndDate                   time.Time `json:"end_date" db:"end_date"`
	Status                    string    `json:"status" db:"status"` // active, completed, cancelled
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}

// ==================== ACHIEVEMENT ANALYTICS ====================

// AchievementEvent logs user achievement activities for analytics
type AchievementEvent struct {
	ID        int64           `json:"id" db:"id"`
	UserID    int64           `json:"user_id" db:"user_id"`
	TenantID  string          `json:"tenant_id" db:"tenant_id"`
	EventType string          `json:"event_type" db:"event_type"`
	EventName string          `json:"event_name" db:"event_name"`
	Details   json.RawMessage `json:"details" db:"details"`
	Timestamp time.Time       `json:"timestamp" db:"timestamp"`
}

// ==================== GAMIFICATION PROFILE ====================

// GamificationProfile holds complete gamification data for a user
type GamificationProfile struct {
	UserPoints       *UserPoints       `json:"user_points"`
	CurrentLevel     *UserLevel        `json:"current_level"`
	Badges           []*UserBadge      `json:"badges"`
	ActiveChallenges []*UserChallenge  `json:"active_challenges"`
	LeaderboardRank  *LeaderboardEntry `json:"leaderboard_rank"`
	TotalRewards     int               `json:"total_rewards"`
	StreakDays       int               `json:"streak_days"`
}

// LeaderboardData holds leaderboard information
type LeaderboardData struct {
	Entries    []*LeaderboardEntry `json:"entries"`
	UserRank   *LeaderboardEntry   `json:"user_rank"`
	PeriodType string              `json:"period_type"`
	PeriodDate time.Time           `json:"period_date"`
}
