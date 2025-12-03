package services

import (
	"context"
	"database/sql"
	"fmt"
	"vyomtech-backend/internal/models"
	"time"
)

// GamificationService handles all gamification logic
type GamificationService struct {
	db *sql.DB
}

// NewGamificationService creates a new gamification service
func NewGamificationService(db *sql.DB) *GamificationService {
	return &GamificationService{db: db}
}

// ==================== POINTS MANAGEMENT ====================

// AwardPoints awards points to a user
func (gs *GamificationService) AwardPoints(ctx context.Context, userID int64, tenantID string, actionType string, points int, description string, bonusReason string) error {
	multiplier := 1.0

	// Apply bonus multipliers
	if bonusReason != "" {
		switch bonusReason {
		case "perfect_call":
			multiplier = 2.0
		case "streak_bonus":
			multiplier = 1.5
		case "level_unlock":
			multiplier = 1.25
		}
	}

	finalPoints := int(float64(points) * multiplier)

	// Start transaction
	tx, err := gs.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Log transaction
	_, err = tx.ExecContext(ctx,
		`INSERT INTO point_transactions (user_id, tenant_id, points, action_type, description, multiplier, bonus_reason, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, tenantID, finalPoints, actionType, description, multiplier, bonusReason, "awarded")
	if err != nil {
		return err
	}

	// Update user points
	_, err = tx.ExecContext(ctx,
		`UPDATE user_points 
		 SET current_points = current_points + ?, lifetime_points = lifetime_points + ?, daily_points = daily_points + ?, updated_at = NOW()
		 WHERE user_id = ? AND tenant_id = ?`,
		finalPoints, finalPoints, finalPoints, userID, tenantID)
	if err != nil {
		return err
	}

	// Update streak
	today := time.Now().Format("2006-01-02")
	_, err = tx.ExecContext(ctx,
		`UPDATE user_points 
		 SET streak_days = IF(daily_date = ?, streak_days + 1, 1), 
		     daily_date = CURDATE(),
		     last_action_date = CURDATE()
		 WHERE user_id = ? AND tenant_id = ?`,
		today, userID, tenantID)
	if err != nil {
		return err
	}

	// Check for level up
	err = gs.checkAndUpdateLevel(tx, ctx, userID, tenantID)
	if err != nil {
		return err
	}

	// Check for challenge completion
	err = gs.updateChallengeProgress(tx, ctx, userID, tenantID, actionType, 1)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RevokePoints removes points from a user
func (gs *GamificationService) RevokePoints(ctx context.Context, userID int64, tenantID string, points int, reason string) error {
	tx, err := gs.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Log transaction
	_, err = tx.ExecContext(ctx,
		`INSERT INTO point_transactions (user_id, tenant_id, points, action_type, description, status)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userID, tenantID, -points, "revoke", reason, "revoked")
	if err != nil {
		return err
	}

	// Update points
	_, err = tx.ExecContext(ctx,
		`UPDATE user_points SET current_points = MAX(0, current_points - ?), updated_at = NOW()
		 WHERE user_id = ? AND tenant_id = ?`,
		points, userID, tenantID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetUserPoints retrieves user's current points
func (gs *GamificationService) GetUserPoints(ctx context.Context, userID int64, tenantID string) (*models.UserPoints, error) {
	up := &models.UserPoints{}
	err := gs.db.QueryRowContext(ctx,
		`SELECT id, user_id, tenant_id, current_points, lifetime_points, period_start_date, 
		        period_end_date, daily_points, daily_date, streak_days, last_action_date, rank, rank_updated_at
		 FROM user_points WHERE user_id = ? AND tenant_id = ?`,
		userID, tenantID).Scan(
		&up.ID, &up.UserID, &up.TenantID, &up.CurrentPoints, &up.LifetimePoints, &up.PeriodStartDate,
		&up.PeriodEndDate, &up.DailyPoints, &up.DailyDate, &up.StreakDays, &up.LastActionDate, &up.Rank, &up.RankUpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return up, err
}

// ==================== BADGES ====================

// CreateBadge creates a new badge
func (gs *GamificationService) CreateBadge(ctx context.Context, badge *models.Badge) error {
	return gs.db.QueryRowContext(ctx,
		`INSERT INTO badge (tenant_id, code, name, description, icon_url, category, requirement_type, requirement_value, points_reward, rarity, active)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 RETURNING id`,
		badge.TenantID, badge.Code, badge.Name, badge.Description, badge.IconURL, badge.Category,
		badge.RequirementType, badge.RequirementValue, badge.PointsReward, badge.Rarity, badge.Active).Scan(&badge.ID)
}

// AwardBadge awards a badge to a user
func (gs *GamificationService) AwardBadge(ctx context.Context, userID int64, badgeID int64, tenantID string) error {
	badge := &models.Badge{}
	err := gs.db.QueryRowContext(ctx,
		`SELECT id, points_reward FROM badge WHERE id = ? AND tenant_id = ?`, badgeID, tenantID).Scan(&badge.ID, &badge.PointsReward)
	if err != nil {
		return err
	}

	tx, err := gs.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if already has badge
	count := 0
	err = tx.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_badge WHERE user_id = ? AND badge_id = ?`, userID, badgeID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("user already has this badge")
	}

	// Award badge
	_, err = tx.ExecContext(ctx,
		`INSERT INTO user_badge (user_id, badge_id, tenant_id, earned_date) VALUES (?, ?, ?, NOW())`,
		userID, badgeID, tenantID)
	if err != nil {
		return err
	}

	// Award bonus points
	if badge.PointsReward > 0 {
		_, err = tx.ExecContext(ctx,
			`UPDATE user_points SET current_points = current_points + ?, lifetime_points = lifetime_points + ?
			 WHERE user_id = ? AND tenant_id = ?`,
			badge.PointsReward, badge.PointsReward, userID, tenantID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetUserBadges gets all badges earned by a user
func (gs *GamificationService) GetUserBadges(ctx context.Context, userID int64, tenantID string) ([]*models.UserBadge, error) {
	rows, err := gs.db.QueryContext(ctx,
		`SELECT ub.id, ub.user_id, ub.badge_id, ub.tenant_id, ub.earned_date, ub.progress_percent, ub.notified,
		        b.id, b.code, b.name, b.description, b.icon_url, b.category, b.requirement_type, b.requirement_value, b.points_reward, b.rarity
		 FROM user_badge ub
		 JOIN badge b ON ub.badge_id = b.id
		 WHERE ub.user_id = ? AND ub.tenant_id = ?
		 ORDER BY ub.earned_date DESC`,
		userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []*models.UserBadge
	for rows.Next() {
		ub := &models.UserBadge{Badge: &models.Badge{}}
		err = rows.Scan(&ub.ID, &ub.UserID, &ub.BadgeID, &ub.TenantID, &ub.EarnedDate, &ub.ProgressPercent, &ub.Notified,
			&ub.Badge.ID, &ub.Badge.Code, &ub.Badge.Name, &ub.Badge.Description, &ub.Badge.IconURL, &ub.Badge.Category,
			&ub.Badge.RequirementType, &ub.Badge.RequirementValue, &ub.Badge.PointsReward, &ub.Badge.Rarity)
		if err != nil {
			return nil, err
		}
		badges = append(badges, ub)
	}

	return badges, rows.Err()
}

// ==================== CHALLENGES ====================

// CreateChallenge creates a new challenge
func (gs *GamificationService) CreateChallenge(ctx context.Context, challenge *models.Challenge) error {
	return gs.db.QueryRowContext(ctx,
		`INSERT INTO challenge (tenant_id, name, description, challenge_type, status, objective_type, objective_target, points_reward, badge_reward_id, start_date, end_date, difficulty)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 RETURNING id`,
		challenge.TenantID, challenge.Name, challenge.Description, challenge.ChallengeType, challenge.Status,
		challenge.ObjectiveType, challenge.ObjectiveTarget, challenge.PointsReward, challenge.BadgeRewardID,
		challenge.StartDate, challenge.EndDate, challenge.Difficulty).Scan(&challenge.ID)
}

// GetActiveChallenges gets all active challenges for a tenant
func (gs *GamificationService) GetActiveChallenges(ctx context.Context, tenantID string) ([]*models.Challenge, error) {
	rows, err := gs.db.QueryContext(ctx,
		`SELECT id, tenant_id, name, description, challenge_type, status, objective_type, objective_target, 
		        objective_current, points_reward, badge_reward_id, start_date, end_date, difficulty
		 FROM challenge
		 WHERE tenant_id = ? AND status = 'active' AND end_date > NOW()
		 ORDER BY end_date ASC`,
		tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*models.Challenge
	for rows.Next() {
		c := &models.Challenge{}
		err = rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Description, &c.ChallengeType, &c.Status,
			&c.ObjectiveType, &c.ObjectiveTarget, &c.ObjectiveCurrent, &c.PointsReward, &c.BadgeRewardID,
			&c.StartDate, &c.EndDate, &c.Difficulty)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	return challenges, rows.Err()
}

// GetUserChallenges gets challenges for a specific user
func (gs *GamificationService) GetUserChallenges(ctx context.Context, userID int64, tenantID string) ([]*models.UserChallenge, error) {
	rows, err := gs.db.QueryContext(ctx,
		`SELECT uc.id, uc.user_id, uc.challenge_id, uc.tenant_id, uc.progress, uc.completed, uc.completed_date, uc.points_earned, uc.badge_earned_id, uc.status,
		        c.id, c.name, c.description, c.challenge_type, c.status, c.objective_type, c.objective_target, c.points_reward, c.difficulty
		 FROM user_challenge uc
		 JOIN challenge c ON uc.challenge_id = c.id
		 WHERE uc.user_id = ? AND uc.tenant_id = ? AND uc.status = 'in_progress'
		 ORDER BY c.end_date ASC`,
		userID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*models.UserChallenge
	for rows.Next() {
		uc := &models.UserChallenge{Challenge: &models.Challenge{}}
		err = rows.Scan(&uc.ID, &uc.UserID, &uc.ChallengeID, &uc.TenantID, &uc.Progress, &uc.Completed, &uc.CompletedDate, &uc.PointsEarned, &uc.BadgeEarnedID, &uc.Status,
			&uc.Challenge.ID, &uc.Challenge.Name, &uc.Challenge.Description, &uc.Challenge.ChallengeType, &uc.Challenge.Status,
			&uc.Challenge.ObjectiveType, &uc.Challenge.ObjectiveTarget, &uc.Challenge.PointsReward, &uc.Challenge.Difficulty)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, uc)
	}

	return challenges, rows.Err()
}

// ==================== LEADERBOARDS ====================

// UpdateLeaderboard recalculates leaderboard rankings
func (gs *GamificationService) UpdateLeaderboard(ctx context.Context, tenantID string, periodType string) error {
	periodDate := time.Now()

	tx, err := gs.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all users and their points
	rows, err := tx.QueryContext(ctx,
		`SELECT user_id, current_points, streak_days 
		 FROM user_points 
		 WHERE tenant_id = ? 
		 ORDER BY current_points DESC`,
		tenantID)
	if err != nil {
		return err
	}
	defer rows.Close()

	rank := 1
	for rows.Next() {
		var userID int64
		var points int
		var streakDays int

		err = rows.Scan(&userID, &points, &streakDays)
		if err != nil {
			return err
		}

		// Count badges
		badgeCount := 0
		err = tx.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM user_badge WHERE user_id = ? AND tenant_id = ?`, userID, tenantID).Scan(&badgeCount)
		if err != nil {
			return err
		}

		// Upsert leaderboard entry
		_, err = tx.ExecContext(ctx,
			`INSERT INTO leaderboard (tenant_id, period_type, period_date, user_id, rank, points, badges_count, streak_days, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
			 ON DUPLICATE KEY UPDATE rank = ?, points = ?, badges_count = ?, streak_days = ?, updated_at = NOW()`,
			tenantID, periodType, periodDate, userID, rank, points, badgeCount, streakDays,
			rank, points, badgeCount, streakDays)
		if err != nil {
			return err
		}

		rank++
	}

	return tx.Commit()
}

// GetLeaderboard gets top users on leaderboard
func (gs *GamificationService) GetLeaderboard(ctx context.Context, tenantID string, periodType string, limit int) ([]*models.LeaderboardEntry, error) {
	rows, err := gs.db.QueryContext(ctx,
		`SELECT id, tenant_id, period_type, period_date, user_id, rank, points, points_change, previous_rank, badges_count, streak_days
		 FROM leaderboard
		 WHERE tenant_id = ? AND period_type = ?
		 ORDER BY rank ASC
		 LIMIT ?`,
		tenantID, periodType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.LeaderboardEntry
	for rows.Next() {
		e := &models.LeaderboardEntry{}
		err = rows.Scan(&e.ID, &e.TenantID, &e.PeriodType, &e.PeriodDate, &e.UserID, &e.Rank, &e.Points, &e.PointsChange, &e.PreviousRank, &e.BadgesCount, &e.StreakDays)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// ==================== LEVELS ====================

// checkAndUpdateLevel checks if user should level up
func (gs *GamificationService) checkAndUpdateLevel(tx *sql.Tx, ctx context.Context, userID int64, tenantID string) error {
	// Get user's current points
	var currentPoints int
	err := tx.QueryRowContext(ctx,
		`SELECT lifetime_points FROM user_points WHERE user_id = ? AND tenant_id = ?`,
		userID, tenantID).Scan(&currentPoints)
	if err != nil {
		return err
	}

	// Get user's current level
	var currentLevel int
	err = tx.QueryRowContext(ctx,
		`SELECT current_level FROM user_level WHERE user_id = ? AND tenant_id = ?`,
		userID, tenantID).Scan(&currentLevel)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Check if user qualifies for next level
	nextLevel := &models.PlayerLevel{}
	err = tx.QueryRowContext(ctx,
		`SELECT level_number, points_total_required FROM player_level 
		 WHERE tenant_id = ? AND level_number = ?`,
		tenantID, currentLevel+1).Scan(&nextLevel.LevelNumber, &nextLevel.PointsTotalRequired)
	if err == sql.ErrNoRows {
		return nil // No next level
	}
	if err != nil {
		return err
	}

	if currentPoints >= nextLevel.PointsTotalRequired {
		// Level up!
		_, err = tx.ExecContext(ctx,
			`UPDATE user_level SET current_level = ?, lifetime_level = ?, level_up_date = NOW(), previous_level = ? 
			 WHERE user_id = ? AND tenant_id = ?`,
			currentLevel+1, currentLevel+1, currentLevel, userID, tenantID)
		if err != nil {
			return err
		}

		// Award bonus for leveling up
		_, err = tx.ExecContext(ctx,
			`INSERT INTO point_transactions (user_id, tenant_id, points, action_type, description, bonus_reason, status)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			userID, tenantID, 100, "level_up", fmt.Sprintf("Reached Level %d", currentLevel+1), "level_unlock", "awarded")
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE user_points SET current_points = current_points + 100 WHERE user_id = ? AND tenant_id = ?`,
			userID, tenantID)
		if err != nil {
			return err
		}
	}

	return nil
}

// updateChallengeProgress updates progress on challenges
func (gs *GamificationService) updateChallengeProgress(tx *sql.Tx, ctx context.Context, userID int64, tenantID string, actionType string, increment int) error {
	// Find relevant challenges
	rows, err := tx.QueryContext(ctx,
		`SELECT uc.id, uc.challenge_id, c.objective_target 
		 FROM user_challenge uc
		 JOIN challenge c ON uc.challenge_id = c.id
		 WHERE uc.user_id = ? AND uc.tenant_id = ? AND uc.status = 'in_progress' AND c.objective_type = ?`,
		userID, tenantID, actionType)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ucID int64
		var challengeID int64
		var target int

		err = rows.Scan(&ucID, &challengeID, &target)
		if err != nil {
			return err
		}

		// Update progress
		var newProgress int
		err = tx.QueryRowContext(ctx,
			`UPDATE user_challenge SET progress = progress + ? WHERE id = ? RETURNING progress`,
			increment, ucID).Scan(&newProgress)
		if err != nil {
			return err
		}

		// Check if completed
		if newProgress >= target {
			// Get points and badge reward
			var pointsReward int64
			var badgeRewardID *int64
			err = tx.QueryRowContext(ctx,
				`SELECT points_reward, badge_reward_id FROM challenge WHERE id = ?`, challengeID).Scan(&pointsReward, &badgeRewardID)
			if err != nil {
				return err
			}

			// Mark as completed
			_, err = tx.ExecContext(ctx,
				`UPDATE user_challenge SET completed = TRUE, completed_date = NOW(), points_earned = ?, status = 'completed' 
				 WHERE id = ?`,
				pointsReward, ucID)
			if err != nil {
				return err
			}

			// Award points
			if pointsReward > 0 {
				_, err = tx.ExecContext(ctx,
					`INSERT INTO point_transactions (user_id, tenant_id, points, action_type, description, status)
					 VALUES (?, ?, ?, ?, ?, ?)`,
					userID, tenantID, pointsReward, "challenge_complete", "Challenge Completed", "awarded")
				if err != nil {
					return err
				}
			}
		}
	}

	return rows.Err()
}

// ==================== GAMIFICATION PROFILE ====================

// GetGamificationProfile gets complete gamification profile for user
func (gs *GamificationService) GetGamificationProfile(ctx context.Context, userID int64, tenantID string) (*models.GamificationProfile, error) {
	profile := &models.GamificationProfile{}

	// Get points
	points, err := gs.GetUserPoints(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	profile.UserPoints = points

	// Get level
	level := &models.UserLevel{}
	err = gs.db.QueryRowContext(ctx,
		`SELECT id, user_id, tenant_id, current_level, lifetime_level, progress_to_next
		 FROM user_level WHERE user_id = ? AND tenant_id = ?`, userID, tenantID).Scan(
		&level.ID, &level.UserID, &level.TenantID, &level.CurrentLevel, &level.LifetimeLevel, &level.ProgressToNext)
	if err == nil {
		profile.CurrentLevel = level
	}

	// Get badges
	badges, err := gs.GetUserBadges(ctx, userID, tenantID)
	if err == nil {
		profile.Badges = badges
	}

	// Get active challenges
	challenges, err := gs.GetUserChallenges(ctx, userID, tenantID)
	if err == nil {
		profile.ActiveChallenges = challenges
	}

	return profile, nil
}
