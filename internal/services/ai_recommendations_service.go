package services

import (
	"context"
	"database/sql"
	"fmt"

	"vyomtech-backend/internal/models"
)

// AIService handles all AI/ML operations including recommendations, predictions, and insights
type AIService struct {
	db *sql.DB
}

// NewAIService creates a new AIService
func NewAIService(db *sql.DB) *AIService {
	return &AIService{
		db: db,
	}
}

// ============================================================
// AI MODEL METHODS
// ============================================================

// CreateAIModel creates a new AI model
func (as *AIService) CreateAIModel(ctx context.Context, model *models.AIModel) error {
	query := `
		INSERT INTO ai_models (
			id, tenant_id, model_name, model_type, model_version, description, algorithm_name,
			input_features, output_format, status, is_production, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		model.ID, model.TenantID, model.ModelName, model.ModelType, model.ModelVersion,
		model.Description, model.AlgorithmName, model.InputFeatures, model.OutputFormat,
		model.Status, model.IsProduction,
	)
	return err
}

// GetAIModel retrieves an AI model
func (as *AIService) GetAIModel(ctx context.Context, id string, tenantID string) (*models.AIModel, error) {
	query := `
		SELECT id, tenant_id, model_name, model_type, model_version, description, algorithm_name,
		       input_features, output_format, accuracy_score, precision_score, recall_score, f1_score,
		       training_dataset_size, training_duration_seconds, training_completed_at, last_updated_by,
		       status, is_production, hyperparameters, metadata, training_notes, model_file_path,
		       model_file_size_bytes, next_retraining_date, retraining_frequency_days,
		       min_data_points_for_training, created_at, updated_at
		FROM ai_models WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	model := &models.AIModel{}
	err := as.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&model.ID, &model.TenantID, &model.ModelName, &model.ModelType, &model.ModelVersion,
		&model.Description, &model.AlgorithmName, &model.InputFeatures, &model.OutputFormat,
		&model.AccuracyScore, &model.PrecisionScore, &model.RecallScore, &model.F1Score,
		&model.TrainingDatasetSize, &model.TrainingDurationSeconds, &model.TrainingCompletedAt,
		&model.LastUpdatedBy, &model.Status, &model.IsProduction, &model.Hyperparameters,
		&model.Metadata, &model.TrainingNotes, &model.ModelFilePath, &model.ModelFileSizeBytes,
		&model.NextRetrainingDate, &model.RetrainingFrequencyDays, &model.MinDataPointsForTraining,
		&model.CreatedAt, &model.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ai model not found")
	}
	return model, err
}

// ListAIModels lists all AI models for a tenant
func (as *AIService) ListAIModels(ctx context.Context, tenantID string, modelType *string, limit int, offset int) ([]*models.AIModel, int, error) {
	query := `SELECT COUNT(*) FROM ai_models WHERE tenant_id = ? AND deleted_at IS NULL`
	args := []interface{}{tenantID}

	if modelType != nil {
		query = query + ` AND model_type = ?`
		args = append(args, *modelType)
	}

	var total int
	err := as.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query = `
		SELECT id, tenant_id, model_name, model_type, model_version, description, algorithm_name,
		       input_features, output_format, accuracy_score, precision_score, recall_score, f1_score,
		       training_dataset_size, training_duration_seconds, training_completed_at, last_updated_by,
		       status, is_production, hyperparameters, metadata, training_notes, model_file_path,
		       model_file_size_bytes, next_retraining_date, retraining_frequency_days,
		       min_data_points_for_training, created_at, updated_at
		FROM ai_models WHERE tenant_id = ? AND deleted_at IS NULL
	`
	args = []interface{}{tenantID}

	if modelType != nil {
		query = query + ` AND model_type = ?`
		args = append(args, *modelType)
	}

	query = query + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := as.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	modelsList := make([]*models.AIModel, 0)
	for rows.Next() {
		model := &models.AIModel{}
		err := rows.Scan(
			&model.ID, &model.TenantID, &model.ModelName, &model.ModelType, &model.ModelVersion,
			&model.Description, &model.AlgorithmName, &model.InputFeatures, &model.OutputFormat,
			&model.AccuracyScore, &model.PrecisionScore, &model.RecallScore, &model.F1Score,
			&model.TrainingDatasetSize, &model.TrainingDurationSeconds, &model.TrainingCompletedAt,
			&model.LastUpdatedBy, &model.Status, &model.IsProduction, &model.Hyperparameters,
			&model.Metadata, &model.TrainingNotes, &model.ModelFilePath, &model.ModelFileSizeBytes,
			&model.NextRetrainingDate, &model.RetrainingFrequencyDays, &model.MinDataPointsForTraining,
			&model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		modelsList = append(modelsList, model)
	}

	return modelsList, total, rows.Err()
}

// ============================================================
// RECOMMENDATION ENGINE METHODS
// ============================================================

// CreateRecommendationEngine creates a new recommendation engine
func (as *AIService) CreateRecommendationEngine(ctx context.Context, engine *models.RecommendationEngine) error {
	query := `
		INSERT INTO recommendation_engine (
			id, tenant_id, engine_name, description, recommendation_type, model_id,
			strategy_type, scoring_algorithm, min_confidence_threshold, max_recommendations_per_user,
			max_recommendations_per_session, recommendation_ttl_hours, enable_personalization,
			enable_real_time_updates, enable_a_b_testing, feedback_enabled, status,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		engine.ID, engine.TenantID, engine.EngineName, engine.Description, engine.RecommendationType,
		engine.ModelID, engine.StrategyType, engine.ScoringAlgorithm, engine.MinConfidenceThreshold,
		engine.MaxRecommendationsPerUser, engine.MaxRecommendationsPerSession, engine.RecommendationTTLHours,
		engine.EnablePersonalization, engine.EnableRealTimeUpdates, engine.EnableABTesting,
		engine.FeedbackEnabled, engine.Status,
	)
	return err
}

// GetRecommendationEngine retrieves a recommendation engine
func (as *AIService) GetRecommendationEngine(ctx context.Context, id string, tenantID string) (*models.RecommendationEngine, error) {
	query := `
		SELECT id, tenant_id, engine_name, description, recommendation_type, model_id,
		       strategy_type, scoring_algorithm, min_confidence_threshold, max_recommendations_per_user,
		       max_recommendations_per_session, recommendation_ttl_hours, enable_personalization,
		       enable_real_time_updates, enable_a_b_testing, a_b_test_variant, filters, boost_rules,
		       penalty_rules, feedback_enabled, feedback_collection_rate, ranking_factors, status,
		       performance_metrics, created_by, last_modified_by, created_at, updated_at
		FROM recommendation_engine WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	engine := &models.RecommendationEngine{}
	err := as.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&engine.ID, &engine.TenantID, &engine.EngineName, &engine.Description, &engine.RecommendationType,
		&engine.ModelID, &engine.StrategyType, &engine.ScoringAlgorithm, &engine.MinConfidenceThreshold,
		&engine.MaxRecommendationsPerUser, &engine.MaxRecommendationsPerSession, &engine.RecommendationTTLHours,
		&engine.EnablePersonalization, &engine.EnableRealTimeUpdates, &engine.EnableABTesting,
		&engine.ABTestVariant, &engine.Filters, &engine.BoostRules, &engine.PenaltyRules,
		&engine.FeedbackEnabled, &engine.FeedbackCollectionRate, &engine.RankingFactors, &engine.Status,
		&engine.PerformanceMetrics, &engine.CreatedBy, &engine.LastModifiedBy, &engine.CreatedAt, &engine.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recommendation engine not found")
	}
	return engine, err
}

// ============================================================
// USER RECOMMENDATIONS METHODS
// ============================================================

// GenerateRecommendations generates recommendations for a user
func (as *AIService) GenerateRecommendations(ctx context.Context, recommendation *models.UserRecommendation) error {
	query := `
		INSERT INTO user_recommendations (
			id, tenant_id, user_id, engine_id, model_id, recommended_item_id, recommended_item_type,
			recommended_item_title, recommendation_score, confidence_level, rank_position,
			recommendation_reason, click_probability, conversion_probability, personalization_level,
			variant_served, expiry_date, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		recommendation.ID, recommendation.TenantID, recommendation.UserID, recommendation.EngineID,
		recommendation.ModelID, recommendation.RecommendedItemID, recommendation.RecommendedItemType,
		recommendation.RecommendedItemTitle, recommendation.RecommendationScore, recommendation.ConfidenceLevel,
		recommendation.RankPosition, recommendation.RecommendationReason, recommendation.ClickProbability,
		recommendation.ConversionProbability, recommendation.PersonalizationLevel, recommendation.VariantServed,
		recommendation.ExpiryDate,
	)
	return err
}

// GetUserRecommendations retrieves recommendations for a user
func (as *AIService) GetUserRecommendations(ctx context.Context, userID int, engineID string, tenantID string, limit int, offset int) ([]*models.UserRecommendation, error) {
	query := `
		SELECT id, tenant_id, user_id, engine_id, model_id, recommended_item_id, recommended_item_type,
		       recommended_item_title, recommendation_score, confidence_level, rank_position,
		       recommendation_reason, recommendation_factors, click_probability, conversion_probability,
		       personalization_level, variant_served, expiry_date, displayed_at, clicked_at,
		       converted_at, user_feedback, feedback_reason, feedback_timestamp, metadata,
		       created_at, updated_at
		FROM user_recommendations 
		WHERE user_id = ? AND engine_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY rank_position ASC LIMIT ? OFFSET ?
	`

	rows, err := as.db.QueryContext(ctx, query, userID, engineID, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recommendations := make([]*models.UserRecommendation, 0)
	for rows.Next() {
		rec := &models.UserRecommendation{}
		err := rows.Scan(
			&rec.ID, &rec.TenantID, &rec.UserID, &rec.EngineID, &rec.ModelID, &rec.RecommendedItemID,
			&rec.RecommendedItemType, &rec.RecommendedItemTitle, &rec.RecommendationScore, &rec.ConfidenceLevel,
			&rec.RankPosition, &rec.RecommendationReason, &rec.RecommendationFactors, &rec.ClickProbability,
			&rec.ConversionProbability, &rec.PersonalizationLevel, &rec.VariantServed, &rec.ExpiryDate,
			&rec.DisplayedAt, &rec.ClickedAt, &rec.ConvertedAt, &rec.UserFeedback, &rec.FeedbackReason,
			&rec.FeedbackTimestamp, &rec.Metadata, &rec.CreatedAt, &rec.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, rows.Err()
}

// ============================================================
// RECOMMENDATION FEEDBACK METHODS
// ============================================================

// SubmitRecommendationFeedback submits feedback on a recommendation
func (as *AIService) SubmitRecommendationFeedback(ctx context.Context, feedback *models.RecommendationFeedback) error {
	query := `
		INSERT INTO recommendation_feedback (
			id, tenant_id, user_id, recommendation_id, engine_id, feedback_type, feedback_value,
			rating_score, detailed_feedback, device_type, session_id, response_time_seconds,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		feedback.ID, feedback.TenantID, feedback.UserID, feedback.RecommendationID, feedback.EngineID,
		feedback.FeedbackType, feedback.FeedbackValue, feedback.RatingScore, feedback.DetailedFeedback,
		feedback.DeviceType, feedback.SessionID, feedback.ResponseTimeSeconds,
	)
	if err == nil {
		// Also update the user_recommendations table
		updateQuery := `
			UPDATE user_recommendations 
			SET user_feedback = ?, feedback_reason = ?, feedback_timestamp = NOW(), updated_at = NOW()
			WHERE id = ? AND tenant_id = ?
		`
		as.db.ExecContext(ctx, updateQuery, feedback.FeedbackValue, feedback.DetailedFeedback,
			feedback.RecommendationID, feedback.TenantID)
	}
	return err
}

// GetRecommendationFeedback retrieves feedback for a recommendation
func (as *AIService) GetRecommendationFeedback(ctx context.Context, recommendationID string, tenantID string) (*models.RecommendationFeedback, error) {
	query := `
		SELECT id, tenant_id, user_id, recommendation_id, engine_id, feedback_type, feedback_value,
		       rating_score, detailed_feedback, device_type, user_context, recommendation_context,
		       external_event_triggered, session_id, response_time_seconds, created_at, updated_at
		FROM recommendation_feedback 
		WHERE recommendation_id = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	feedback := &models.RecommendationFeedback{}
	err := as.db.QueryRowContext(ctx, query, recommendationID, tenantID).Scan(
		&feedback.ID, &feedback.TenantID, &feedback.UserID, &feedback.RecommendationID, &feedback.EngineID,
		&feedback.FeedbackType, &feedback.FeedbackValue, &feedback.RatingScore, &feedback.DetailedFeedback,
		&feedback.DeviceType, &feedback.UserContext, &feedback.RecommendationContext,
		&feedback.ExternalEventTriggered, &feedback.SessionID, &feedback.ResponseTimeSeconds,
		&feedback.CreatedAt, &feedback.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("feedback not found")
	}
	return feedback, err
}

// ============================================================
// PREDICTION METHODS
// ============================================================

// MakePrediction creates a prediction result
func (as *AIService) MakePrediction(ctx context.Context, prediction *models.PredictionResult) error {
	query := `
		INSERT INTO prediction_results (
			id, tenant_id, model_id, input_data_id, input_data_type, prediction_type,
			predicted_value, predicted_class, confidence_score, probability_distribution,
			feature_importance, processing_time_ms, prediction_method, batch_id,
			execution_environment, model_inference_version, input_features_hash, drift_detected,
			drift_magnitude, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		prediction.ID, prediction.TenantID, prediction.ModelID, prediction.InputDataID, prediction.InputDataType,
		prediction.PredictionType, prediction.PredictedValue, prediction.PredictedClass, prediction.ConfidenceScore,
		prediction.ProbabilityDistribution, prediction.FeatureImportance, prediction.ProcessingTimeMs,
		prediction.PredictionMethod, prediction.BatchID, prediction.ExecutionEnvironment, prediction.ModelInferenceVersion,
		prediction.InputFeaturesHash, prediction.DriftDetected, prediction.DriftMagnitude,
	)
	return err
}

// GetPredictionResult retrieves a prediction result
func (as *AIService) GetPredictionResult(ctx context.Context, id string, tenantID string) (*models.PredictionResult, error) {
	query := `
		SELECT id, tenant_id, model_id, input_data_id, input_data_type, prediction_type,
		       predicted_value, predicted_class, confidence_score, probability_distribution,
		       feature_importance, processing_time_ms, prediction_method, batch_id,
		       execution_environment, model_inference_version, input_features_hash, drift_detected,
		       drift_magnitude, actual_outcome, actual_outcome_date, prediction_accuracy,
		       feedback_loop_triggered, metadata, created_at, updated_at
		FROM prediction_results WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	prediction := &models.PredictionResult{}
	err := as.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&prediction.ID, &prediction.TenantID, &prediction.ModelID, &prediction.InputDataID,
		&prediction.InputDataType, &prediction.PredictionType, &prediction.PredictedValue,
		&prediction.PredictedClass, &prediction.ConfidenceScore, &prediction.ProbabilityDistribution,
		&prediction.FeatureImportance, &prediction.ProcessingTimeMs, &prediction.PredictionMethod,
		&prediction.BatchID, &prediction.ExecutionEnvironment, &prediction.ModelInferenceVersion,
		&prediction.InputFeaturesHash, &prediction.DriftDetected, &prediction.DriftMagnitude,
		&prediction.ActualOutcome, &prediction.ActualOutcomeDate, &prediction.PredictionAccuracy,
		&prediction.FeedbackLoopTriggered, &prediction.Metadata, &prediction.CreatedAt, &prediction.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("prediction not found")
	}
	return prediction, err
}

// ============================================================
// ANOMALY DETECTION METHODS
// ============================================================

// DetectAnomaly records a detected anomaly
func (as *AIService) DetectAnomaly(ctx context.Context, anomaly *models.AnomalyDetection) error {
	query := `
		INSERT INTO anomaly_detection (
			id, tenant_id, anomaly_type, severity_level, affected_model_id, affected_entity_id,
			affected_entity_type, anomaly_description, anomaly_metric, expected_value, actual_value,
			deviation_percentage, z_score, detection_method, detection_confidence,
			root_cause_analysis, recommended_action, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		anomaly.ID, anomaly.TenantID, anomaly.AnomalyType, anomaly.SeverityLevel, anomaly.AffectedModelID,
		anomaly.AffectedEntityID, anomaly.AffectedEntityType, anomaly.AnomalyDescription, anomaly.AnomalyMetric,
		anomaly.ExpectedValue, anomaly.ActualValue, anomaly.DeviationPercentage, anomaly.ZScore,
		anomaly.DetectionMethod, anomaly.DetectionConfidence, anomaly.RootCauseAnalysis, anomaly.RecommendedAction,
	)
	return err
}

// GetAnomalies retrieves anomalies for a tenant
func (as *AIService) GetAnomalies(ctx context.Context, tenantID string, severityLevel *string, limit int, offset int) ([]*models.AnomalyDetection, int, error) {
	countQuery := `SELECT COUNT(*) FROM anomaly_detection WHERE tenant_id = ? AND deleted_at IS NULL`
	args := []interface{}{tenantID}

	if severityLevel != nil {
		countQuery = countQuery + ` AND severity_level = ?`
		args = append(args, *severityLevel)
	}

	var total int
	err := as.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, tenant_id, anomaly_type, severity_level, affected_model_id, affected_entity_id,
		       affected_entity_type, anomaly_description, anomaly_metric, expected_value, actual_value,
		       deviation_percentage, z_score, detection_method, detection_confidence,
		       root_cause_analysis, recommended_action, is_acknowledged, is_resolved,
		       created_at, updated_at
		FROM anomaly_detection WHERE tenant_id = ? AND deleted_at IS NULL
	`
	args = []interface{}{tenantID}

	if severityLevel != nil {
		query = query + ` AND severity_level = ?`
		args = append(args, *severityLevel)
	}

	query = query + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := as.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	anomalies := make([]*models.AnomalyDetection, 0)
	for rows.Next() {
		anomaly := &models.AnomalyDetection{}
		err := rows.Scan(
			&anomaly.ID, &anomaly.TenantID, &anomaly.AnomalyType, &anomaly.SeverityLevel, &anomaly.AffectedModelID,
			&anomaly.AffectedEntityID, &anomaly.AffectedEntityType, &anomaly.AnomalyDescription, &anomaly.AnomalyMetric,
			&anomaly.ExpectedValue, &anomaly.ActualValue, &anomaly.DeviationPercentage, &anomaly.ZScore,
			&anomaly.DetectionMethod, &anomaly.DetectionConfidence, &anomaly.RootCauseAnalysis, &anomaly.RecommendedAction,
			&anomaly.IsAcknowledged, &anomaly.IsResolved, &anomaly.CreatedAt, &anomaly.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		anomalies = append(anomalies, anomaly)
	}

	return anomalies, total, rows.Err()
}

// ============================================================
// AI INSIGHTS METHODS
// ============================================================

// GenerateInsight creates an AI insight
func (as *AIService) GenerateInsight(ctx context.Context, insight *models.AIInsight) error {
	query := `
		INSERT INTO ai_insights (
			id, tenant_id, insight_type, insight_category, title, description,
			generated_by_model_id, confidence_score, impact_score, relevance_score,
			recommendation, time_period, trend_direction, audience, priority_level,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		insight.ID, insight.TenantID, insight.InsightType, insight.InsightCategory, insight.Title,
		insight.Description, insight.GeneratedByModelID, insight.ConfidenceScore, insight.ImpactScore,
		insight.RelevanceScore, insight.Recommendation, insight.TimePeriod, insight.TrendDirection,
		insight.Audience, insight.PriorityLevel,
	)
	return err
}

// GetInsights retrieves insights for a tenant
func (as *AIService) GetInsights(ctx context.Context, tenantID string, insightType *string, limit int, offset int) ([]*models.AIInsight, int, error) {
	countQuery := `SELECT COUNT(*) FROM ai_insights WHERE tenant_id = ? AND is_archived = 0 AND deleted_at IS NULL`
	args := []interface{}{tenantID}

	if insightType != nil {
		countQuery = countQuery + ` AND insight_type = ?`
		args = append(args, *insightType)
	}

	var total int
	err := as.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, tenant_id, insight_type, insight_category, title, description,
		       generated_by_model_id, confidence_score, impact_score, relevance_score,
		       supporting_data, visualizations, recommendation, affected_kpis, time_period,
		       trend_direction, forecast_period_days, historical_accuracy_rate, audience,
		       priority_level, is_shared, is_archived, feedback_score, feedback_count,
		       created_at, updated_at
		FROM ai_insights WHERE tenant_id = ? AND is_archived = 0 AND deleted_at IS NULL
	`
	args = []interface{}{tenantID}

	if insightType != nil {
		query = query + ` AND insight_type = ?`
		args = append(args, *insightType)
	}

	query = query + ` ORDER BY priority_level DESC, confidence_score DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := as.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	insights := make([]*models.AIInsight, 0)
	for rows.Next() {
		insight := &models.AIInsight{}
		err := rows.Scan(
			&insight.ID, &insight.TenantID, &insight.InsightType, &insight.InsightCategory, &insight.Title,
			&insight.Description, &insight.GeneratedByModelID, &insight.ConfidenceScore, &insight.ImpactScore,
			&insight.RelevanceScore, &insight.SupportingData, &insight.Visualizations, &insight.Recommendation,
			&insight.AffectedKPIs, &insight.TimePeriod, &insight.TrendDirection, &insight.ForecastPeriodDays,
			&insight.HistoricalAccuracyRate, &insight.Audience, &insight.PriorityLevel, &insight.IsShared,
			&insight.IsArchived, &insight.FeedbackScore, &insight.FeedbackCount, &insight.CreatedAt, &insight.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		insights = append(insights, insight)
	}

	return insights, total, rows.Err()
}

// ============================================================
// MODEL PERFORMANCE METHODS
// ============================================================

// RecordModelPerformance records model performance metrics
func (as *AIService) RecordModelPerformance(ctx context.Context, perf *models.ModelPerformance) error {
	query := `
		INSERT INTO model_performance (
			id, tenant_id, model_id, performance_date, total_predictions, successful_predictions,
			failed_predictions, accuracy, precision, recall, f1_score, auc_roc, log_loss,
			mean_absolute_error, mean_squared_error, root_mean_squared_error, avg_inference_time_ms,
			p95_inference_time_ms, p99_inference_time_ms, throughput_predictions_per_second,
			error_rate, data_drift_score, model_degradation_detected, retraining_recommended,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		perf.ID, perf.TenantID, perf.ModelID, perf.PerformanceDate, perf.TotalPredictions,
		perf.SuccessfulPredictions, perf.FailedPredictions, perf.Accuracy, perf.Precision, perf.Recall,
		perf.F1Score, perf.AUCROC, perf.LogLoss, perf.MeanAbsoluteError, perf.MeanSquaredError,
		perf.RootMeanSquaredError, perf.AvgInferenceTimeMs, perf.P95InferenceTimeMs, perf.P99InferenceTimeMs,
		perf.ThroughputPredictionsPerSec, perf.ErrorRate, perf.DataDriftScore, perf.ModelDegradationDetected,
		perf.RetrainingRecommended,
	)
	return err
}

// GetModelPerformance retrieves model performance metrics
func (as *AIService) GetModelPerformance(ctx context.Context, modelID string, tenantID string) (*models.ModelPerformance, error) {
	query := `
		SELECT id, tenant_id, model_id, performance_date, total_predictions, successful_predictions,
		       failed_predictions, accuracy, precision, recall, f1_score, auc_roc, log_loss,
		       mean_absolute_error, mean_squared_error, root_mean_squared_error, avg_inference_time_ms,
		       p95_inference_time_ms, p99_inference_time_ms, throughput_predictions_per_second,
		       error_rate, prediction_distribution, feature_importance_snapshot, data_drift_score,
		       model_degradation_detected, degradation_reason, retraining_recommended,
		       confidence_interval_lower, confidence_interval_upper, segment_performance,
		       error_analysis, false_positive_rate, false_negative_rate, baseline_model_comparison,
		       notes, created_at, updated_at
		FROM model_performance WHERE model_id = ? AND tenant_id = ? AND deleted_at IS NULL
		ORDER BY performance_date DESC LIMIT 1
	`

	perf := &models.ModelPerformance{}
	err := as.db.QueryRowContext(ctx, query, modelID, tenantID).Scan(
		&perf.ID, &perf.TenantID, &perf.ModelID, &perf.PerformanceDate, &perf.TotalPredictions,
		&perf.SuccessfulPredictions, &perf.FailedPredictions, &perf.Accuracy, &perf.Precision, &perf.Recall,
		&perf.F1Score, &perf.AUCROC, &perf.LogLoss, &perf.MeanAbsoluteError, &perf.MeanSquaredError,
		&perf.RootMeanSquaredError, &perf.AvgInferenceTimeMs, &perf.P95InferenceTimeMs, &perf.P99InferenceTimeMs,
		&perf.ThroughputPredictionsPerSec, &perf.ErrorRate, &perf.PredictionDistribution,
		&perf.FeatureImportanceSnapshot, &perf.DataDriftScore, &perf.ModelDegradationDetected,
		&perf.DegradationReason, &perf.RetrainingRecommended, &perf.ConfidenceIntervalLower,
		&perf.ConfidenceIntervalUpper, &perf.SegmentPerformance, &perf.ErrorAnalysis,
		&perf.FalsePositiveRate, &perf.FalseNegativeRate, &perf.BaselineModelComparison,
		&perf.Notes, &perf.CreatedAt, &perf.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("performance metrics not found")
	}
	return perf, err
}

// ============================================================
// TRAINING DATA METHODS
// ============================================================

// CreateTrainingDataset creates a training dataset record
func (as *AIService) CreateTrainingDataset(ctx context.Context, dataset *models.MLTrainingData) error {
	query := `
		INSERT INTO ml_training_data (
			id, tenant_id, model_id, dataset_name, dataset_version, data_type, source_table,
			total_records, valid_records, invalid_records, duplicate_records, missing_value_percentage,
			outliers_detected, outliers_removed, feature_scaling_method, feature_encoding_method,
			train_test_split_ratio, validation_set_ratio, class_imbalance_handling,
			temporal_split_date, data_collection_start_date, data_collection_end_date,
			preprocessing_steps, quality_score, file_path, file_size_mb, hash_value,
			created_by, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		dataset.ID, dataset.TenantID, dataset.ModelID, dataset.DatasetName, dataset.DatasetVersion,
		dataset.DataType, dataset.SourceTable, dataset.TotalRecords, dataset.ValidRecords, dataset.InvalidRecords,
		dataset.DuplicateRecords, dataset.MissingValuePercentage, dataset.OutliersDetected, dataset.OutliersRemoved,
		dataset.FeatureScalingMethod, dataset.FeatureEncodingMethod, dataset.TrainTestSplitRatio,
		dataset.ValidationSetRatio, dataset.ClassImbalanceHandling, dataset.TemporalSplitDate,
		dataset.DataCollectionStartDate, dataset.DataCollectionEndDate, dataset.PreprocessingSteps,
		dataset.QualityScore, dataset.FilePath, dataset.FileSizeMB, dataset.HashValue, dataset.CreatedBy,
	)
	return err
}

// GetTrainingDataset retrieves a training dataset
func (as *AIService) GetTrainingDataset(ctx context.Context, id string, tenantID string) (*models.MLTrainingData, error) {
	query := `
		SELECT id, tenant_id, model_id, dataset_name, dataset_version, data_type, source_table,
		       total_records, valid_records, invalid_records, duplicate_records, missing_value_percentage,
		       outliers_detected, outliers_removed, feature_scaling_method, feature_encoding_method,
		       train_test_split_ratio, validation_set_ratio, class_imbalance_handling,
		       temporal_split_date, data_collection_start_date, data_collection_end_date,
		       preprocessing_steps, quality_score, file_path, file_size_mb, hash_value,
		       created_by, created_at, updated_at
		FROM ml_training_data WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	dataset := &models.MLTrainingData{}
	err := as.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&dataset.ID, &dataset.TenantID, &dataset.ModelID, &dataset.DatasetName, &dataset.DatasetVersion,
		&dataset.DataType, &dataset.SourceTable, &dataset.TotalRecords, &dataset.ValidRecords,
		&dataset.InvalidRecords, &dataset.DuplicateRecords, &dataset.MissingValuePercentage,
		&dataset.OutliersDetected, &dataset.OutliersRemoved, &dataset.FeatureScalingMethod,
		&dataset.FeatureEncodingMethod, &dataset.TrainTestSplitRatio, &dataset.ValidationSetRatio,
		&dataset.ClassImbalanceHandling, &dataset.TemporalSplitDate, &dataset.DataCollectionStartDate,
		&dataset.DataCollectionEndDate, &dataset.PreprocessingSteps, &dataset.QualityScore,
		&dataset.FilePath, &dataset.FileSizeMB, &dataset.HashValue, &dataset.CreatedBy,
		&dataset.CreatedAt, &dataset.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("training dataset not found")
	}
	return dataset, err
}

// ============================================================
// RECOMMENDATION HISTORY METHODS
// ============================================================

// RecordRecommendationHistory records historical recommendation data
func (as *AIService) RecordRecommendationHistory(ctx context.Context, history *models.RecommendationHistory) error {
	query := `
		INSERT INTO recommendation_history (
			id, tenant_id, user_id, engine_id, batch_timestamp, recommendation_batch_id,
			total_recommendations_shown, total_recommendations_clicked, total_recommendations_converted,
			ctr, conversion_rate, avg_recommendation_score, session_duration_seconds,
			device_type, user_segment, personalization_used, variant_served, ab_test_group,
			unique_items_recommended, repeat_recommendations, user_segment_at_time,
			contextual_factors, system_performance, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	_, err := as.db.ExecContext(ctx, query,
		history.ID, history.TenantID, history.UserID, history.EngineID, history.BatchTimestamp,
		history.RecommendationBatchID, history.TotalRecommendationsShown, history.TotalRecommendationsClicked,
		history.TotalRecommendationsConverted, history.CTR, history.ConversionRate, history.AvgRecommendationScore,
		history.SessionDurationSeconds, history.DeviceType, history.UserSegment, history.PersonalizationUsed,
		history.VariantServed, history.ABTestGroup, history.UniqueItemsRecommended, history.RepeatRecommendations,
		history.UserSegmentAtTime, history.ContextualFactors, history.SystemPerformance,
	)
	return err
}

// GetRecommendationStats retrieves recommendation statistics for an engine
func (as *AIService) GetRecommendationStats(ctx context.Context, tenantID string, engineID string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_recommendations_shown), 0) as total_shown,
			COALESCE(SUM(total_recommendations_clicked), 0) as total_clicked,
			COALESCE(SUM(total_recommendations_converted), 0) as total_converted,
			COALESCE(AVG(ctr), 0) as avg_ctr,
			COALESCE(AVG(conversion_rate), 0) as avg_conversion_rate,
			COUNT(*) as total_sessions
		FROM recommendation_history 
		WHERE engine_id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	var totalShown, totalClicked, totalConverted, totalSessions int
	var avgCTR, avgConversionRate float64

	err := as.db.QueryRowContext(ctx, query, engineID, tenantID).Scan(
		&totalShown, &totalClicked, &totalConverted, &avgCTR, &avgConversionRate, &totalSessions,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_recommendations_shown":     totalShown,
		"total_recommendations_clicked":   totalClicked,
		"total_recommendations_converted": totalConverted,
		"average_ctr":                     avgCTR,
		"average_conversion_rate":         avgConversionRate,
		"total_sessions":                  totalSessions,
	}

	return stats, nil
}

// ============================================================
// BATCH OPERATIONS
// ============================================================

// BatchMarkRecommendationsDisplayed marks multiple recommendations as displayed
func (as *AIService) BatchMarkRecommendationsDisplayed(ctx context.Context, recommendationIDs []string, tenantID string) error {
	if len(recommendationIDs) == 0 {
		return nil
	}

	query := `
		UPDATE user_recommendations 
		SET displayed_at = NOW(), updated_at = NOW()
		WHERE tenant_id = ? AND displayed_at IS NULL AND id IN (`

	for i := 0; i < len(recommendationIDs); i++ {
		if i > 0 {
			query = query + `,`
		}
		query = query + `?`
	}
	query = query + `)`

	args := []interface{}{tenantID}
	for _, id := range recommendationIDs {
		args = append(args, id)
	}

	_, err := as.db.ExecContext(ctx, query, args...)
	return err
}

// GetAIStats returns overall AI system statistics
func (as *AIService) GetAIStats(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count models
	var modelCount int
	as.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ai_models WHERE tenant_id = ? AND deleted_at IS NULL`, tenantID).Scan(&modelCount)
	stats["model_count"] = modelCount

	// Count engines
	var engineCount int
	as.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM recommendation_engine WHERE tenant_id = ? AND deleted_at IS NULL`, tenantID).Scan(&engineCount)
	stats["engine_count"] = engineCount

	// Count predictions this month
	var predictionCount int
	as.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM prediction_results 
		WHERE tenant_id = ? AND deleted_at IS NULL AND DATE(created_at) >= DATE_SUB(NOW(), INTERVAL 30 DAY)
	`, tenantID).Scan(&predictionCount)
	stats["predictions_this_month"] = predictionCount

	// Count anomalies
	var anomalyCount int
	as.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM anomaly_detection 
		WHERE tenant_id = ? AND deleted_at IS NULL AND is_resolved = 0
	`, tenantID).Scan(&anomalyCount)
	stats["unresolved_anomalies"] = anomalyCount

	// Count insights
	var insightCount int
	as.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM ai_insights 
		WHERE tenant_id = ? AND deleted_at IS NULL AND is_archived = 0
	`, tenantID).Scan(&insightCount)
	stats["active_insights"] = insightCount

	return stats, nil
}
