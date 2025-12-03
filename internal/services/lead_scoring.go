package services

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"vyomtech-backend/internal/models"
	"vyomtech-backend/pkg/logger"
)

type LeadScoringService struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewLeadScoringService(db *sql.DB, logger *logger.Logger) *LeadScoringService {
	return &LeadScoringService{
		db:     db,
		logger: logger,
	}
}

// CalculateLeadScore calculates a composite lead score (0-100)
// Components:
// - Source Quality (0-25): Different sources have different quality
// - Engagement (0-25): Contact info completeness
// - Conversion Probability (0-30): Based on lead status and actions
// - Urgency (0-20): How soon action is needed
func (s *LeadScoringService) CalculateLeadScore(ctx context.Context, leadID int64, tenantID string) (*models.LeadScore, error) {
	// Get lead details
	var source, status string
	var email, phone, company sql.NullString
	var priority sql.NullString
	var createdAt time.Time

	query := "SELECT source, email, phone, company, status, priority, created_at " +
		"FROM `lead` WHERE id = ? AND tenant_id = ?"
	err := s.db.QueryRowContext(ctx, query, leadID, tenantID).Scan(
		&source, &email, &phone, &company, &status, &priority, &createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lead not found")
		}
		s.logger.Error("error fetching lead for scoring", map[string]interface{}{
			"lead_id": leadID, "error": err,
		})
		return nil, err
	}

	// Get existing score for comparison
	oldScore, _ := s.GetLeadScore(ctx, leadID, tenantID)
	var oldScoreValue *float64
	if oldScore != nil {
		oldScoreValue = &oldScore.OverallScore
	}

	// Calculate individual scores
	sourceScore := s.calculateSourceScore(source)                               // 0-25
	engagementScore := s.calculateEngagementScore(email.String, phone.String)   // 0-25
	conversionProb := s.calculateConversionProbability(status, priority.String) // 0-30
	urgencyScore := s.calculateUrgencyScore(createdAt)                          // 0-20

	// Weighted overall score (0-100)
	overall := (sourceScore * 0.25) + (engagementScore * 0.25) +
		(conversionProb * 0.30) + (urgencyScore * 0.20)

	// Round to 2 decimal places
	overall = math.Round(overall*100) / 100

	// Determine category
	category := s.categorizeScore(overall)

	// Calculate score change
	var scoreChange *float64
	if oldScoreValue != nil {
		change := overall - *oldScoreValue
		scoreChange = &change
	}

	leadScore := &models.LeadScore{
		LeadID:                leadID,
		TenantID:              tenantID,
		SourceQualityScore:    sourceScore,
		EngagementScore:       engagementScore,
		ConversionProbability: conversionProb,
		UrgencyScore:          urgencyScore,
		OverallScore:          overall,
		ScoreCategory:         category,
		PreviousScore:         oldScoreValue,
		ScoreChange:           scoreChange,
		CalculationMethod:     "weighted",
		LastCalculated:        time.Now(),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	return leadScore, nil
}

// GetLeadScore retrieves existing lead score
func (s *LeadScoringService) GetLeadScore(ctx context.Context, leadID int64, tenantID string) (*models.LeadScore, error) {
	score := &models.LeadScore{}
	query := "SELECT id, lead_id, tenant_id, source_quality_score, engagement_score, " +
		"conversion_probability, urgency_score, overall_score, " +
		"score_category, previous_score, score_change, reason_text, " +
		"calculation_method, last_calculated, created_at, updated_at " +
		"FROM `lead_scores` WHERE lead_id = ? AND tenant_id = ?"

	err := s.db.QueryRowContext(ctx, query, leadID, tenantID).Scan(
		&score.ID, &score.LeadID, &score.TenantID, &score.SourceQualityScore,
		&score.EngagementScore, &score.ConversionProbability, &score.UrgencyScore,
		&score.OverallScore, &score.ScoreCategory, &score.PreviousScore, &score.ScoreChange,
		&score.ReasonText, &score.CalculationMethod, &score.LastCalculated,
		&score.CreatedAt, &score.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return score, nil
}

// SaveLeadScore saves or updates lead score in database
func (s *LeadScoringService) SaveLeadScore(ctx context.Context, score *models.LeadScore) error {
	// Check if score exists
	existing, err := s.GetLeadScore(ctx, score.LeadID, score.TenantID)
	if err != nil {
		return err
	}

	if existing == nil {
		// Insert new score
		query := "INSERT INTO `lead_scores` (" +
			"lead_id, tenant_id, source_quality_score, engagement_score, " +
			"conversion_probability, urgency_score, overall_score, " +
			"score_category, calculation_method, last_calculated" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		_, err := s.db.ExecContext(ctx, query,
			score.LeadID, score.TenantID, score.SourceQualityScore, score.EngagementScore,
			score.ConversionProbability, score.UrgencyScore, score.OverallScore,
			score.ScoreCategory, score.CalculationMethod, score.LastCalculated,
		)
		return err
	}

	// Update existing score
	query := "UPDATE `lead_scores` " +
		"SET source_quality_score = ?, engagement_score = ?," +
		"    conversion_probability = ?, urgency_score = ?," +
		"    overall_score = ?, score_category = ?," +
		"    previous_score = ?, score_change = ?," +
		"    last_calculated = NOW() " +
		"WHERE lead_id = ? AND tenant_id = ?"

	_, err = s.db.ExecContext(ctx, query,
		score.SourceQualityScore, score.EngagementScore,
		score.ConversionProbability, score.UrgencyScore,
		score.OverallScore, score.ScoreCategory,
		existing.OverallScore, score.ScoreChange,
		score.LeadID, score.TenantID,
	)
	return err
}

// UpdateLeadScore recalculates and updates lead score
func (s *LeadScoringService) UpdateLeadScore(ctx context.Context, leadID int64, tenantID string) (*models.LeadScore, error) {
	score, err := s.CalculateLeadScore(ctx, leadID, tenantID)
	if err != nil {
		return nil, err
	}

	err = s.SaveLeadScore(ctx, score)
	if err != nil {
		return nil, err
	}

	return score, nil
}

// GetLeadsByCategory retrieves leads in a specific score category
func (s *LeadScoringService) GetLeadsByCategory(ctx context.Context, tenantID, category string, limit int) ([]models.LeadScore, error) {
	scores := []models.LeadScore{}
	query := "SELECT id, lead_id, tenant_id, source_quality_score, engagement_score, " +
		"conversion_probability, urgency_score, overall_score, " +
		"score_category, previous_score, score_change, reason_text, " +
		"calculation_method, last_calculated, created_at, updated_at " +
		"FROM `lead_scores` " +
		"WHERE tenant_id = ? AND score_category = ? " +
		"ORDER BY overall_score DESC " +
		"LIMIT ?"

	rows, err := s.db.QueryContext(ctx, query, tenantID, category, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		score := models.LeadScore{}
		err := rows.Scan(
			&score.ID, &score.LeadID, &score.TenantID, &score.SourceQualityScore,
			&score.EngagementScore, &score.ConversionProbability, &score.UrgencyScore,
			&score.OverallScore, &score.ScoreCategory, &score.PreviousScore, &score.ScoreChange,
			&score.ReasonText, &score.CalculationMethod, &score.LastCalculated,
			&score.CreatedAt, &score.UpdatedAt,
		)
		if err != nil {
			s.logger.Error("error scanning lead score", map[string]interface{}{
				"error": err,
			})
			continue
		}
		scores = append(scores, score)
	}

	return scores, rows.Err()
}

// BatchCalculateScores recalculates all lead scores for recent leads
func (s *LeadScoringService) BatchCalculateScores(ctx context.Context, tenantID string) error {
	query := "SELECT id FROM `lead` " +
		"WHERE tenant_id = ? AND created_at > DATE_SUB(NOW(), INTERVAL 30 DAY)"

	rows, err := s.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var leadID int64
		if err := rows.Scan(&leadID); err != nil {
			continue
		}
		_, _ = s.UpdateLeadScore(ctx, leadID, tenantID)
		count++
	}

	s.logger.Info("batch score calculation completed", map[string]interface{}{
		"tenant_id": tenantID,
		"count":     count,
	})

	return rows.Err()
}

// ========== Scoring Helper Methods ==========

// calculateSourceScore: Different sources have different quality (0-25)
func (s *LeadScoringService) calculateSourceScore(source string) float64 {
	scores := map[string]float64{
		"direct_website": 25.0,
		"referral":       24.0,
		"google_ads":     22.0,
		"facebook_ads":   20.0,
		"instagram_ads":  18.0,
		"linkedin_ads":   21.0,
		"cold_call":      15.0,
		"event":          19.0,
		"email_campaign": 17.0,
		"other":          10.0,
	}
	if score, exists := scores[source]; exists {
		return score
	}
	return 10.0
}

// calculateEngagementScore: Contact info completeness (0-25)
func (s *LeadScoringService) calculateEngagementScore(email, phone string) float64 {
	score := 0.0
	if email != "" && email != "null" {
		score += 12.5
	}
	if phone != "" && phone != "null" {
		score += 12.5
	}
	return score
}

// calculateConversionProbability: Based on lead status and priority (0-30)
func (s *LeadScoringService) calculateConversionProbability(status, priority string) float64 {
	score := 5.0

	// Status bonus
	switch status {
	case "converted":
		score = 30.0
	case "qualified":
		score = 25.0
	case "contacted":
		score = 15.0
	case "new":
		score = 5.0
	default:
		score = 5.0
	}

	// Priority bonus
	switch priority {
	case "high":
		score += 5.0
	case "medium":
		score += 2.0
	}

	if score > 30.0 {
		score = 30.0
	}

	return score
}

// calculateUrgencyScore: How soon action is needed based on lead age (0-20)
func (s *LeadScoringService) calculateUrgencyScore(createdAt time.Time) float64 {
	now := time.Now()
	ageHours := now.Sub(createdAt).Hours()
	ageDays := ageHours / 24

	// Newer leads are more urgent
	if ageDays <= 1 {
		return 20.0 // Created today
	} else if ageDays <= 3 {
		return 16.0 // Last 3 days
	} else if ageDays <= 7 {
		return 12.0 // Within a week
	} else if ageDays <= 14 {
		return 8.0 // Within 2 weeks
	} else if ageDays <= 30 {
		return 5.0 // Within a month
	}

	return 2.0 // Older than a month
}

// categorizeScore: Convert numeric score to category
func (s *LeadScoringService) categorizeScore(score float64) string {
	if score >= 75.0 {
		return models.ScoreCategoryHot
	} else if score >= 50.0 {
		return models.ScoreCategoryWarm
	} else if score >= 25.0 {
		return models.ScoreCategoryCold
	}
	return models.ScoreCategoryNurture
}
