package models

import (
	"encoding/json"
	"time"
)

// ============================================================
// EXISTING AI MODELS (ORCHESTRATION)
// ============================================================

type AIRequest struct {
	TenantID    string                 `json:"tenant_id"`
	Query       string                 `json:"query"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Priority    string                 `json:"priority"` // low, medium, high
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
}

type AIResponse struct {
	Response       string        `json:"response"`
	Provider       string        `json:"provider"`
	TokensUsed     int           `json:"tokens_used"`
	ProcessingTime time.Duration `json:"processing_time"`
	Cost           float64       `json:"cost"`
	Cached         bool          `json:"cached"`
}

type AIProvider interface {
	Call(req *AIRequest) (*AIResponse, error)
	GetCostPerToken() float64
	IsAvailable() bool
}

type ProviderConfig struct {
	APIKey       string  `json:"api_key"`
	BaseURL      string  `json:"base_url,omitempty"`
	Model        string  `json:"model"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float64 `json:"temperature"`
	CostPerToken float64 `json:"cost_per_token"`
}

// ============================================================
// AI MODELS - ML MODEL CONFIGURATIONS
// ============================================================

// AIModel represents an ML model configuration
type AIModel struct {
	ID                       string           `json:"id" db:"id"`
	TenantID                 string           `json:"tenant_id" db:"tenant_id"`
	ModelName                string           `json:"model_name" db:"model_name"`
	ModelType                string           `json:"model_type" db:"model_type"` // recommendation, prediction, clustering, classification, anomaly_detection
	ModelVersion             string           `json:"model_version" db:"model_version"`
	Description              *string          `json:"description" db:"description"`
	AlgorithmName            *string          `json:"algorithm_name" db:"algorithm_name"`
	InputFeatures            json.RawMessage  `json:"input_features" db:"input_features"`
	OutputFormat             *string          `json:"output_format" db:"output_format"`
	AccuracyScore            *float64         `json:"accuracy_score" db:"accuracy_score"`
	PrecisionScore           *float64         `json:"precision_score" db:"precision_score"`
	RecallScore              *float64         `json:"recall_score" db:"recall_score"`
	F1Score                  *float64         `json:"f1_score" db:"f1_score"`
	TrainingDatasetSize      *int             `json:"training_dataset_size" db:"training_dataset_size"`
	TrainingDurationSeconds  *int             `json:"training_duration_seconds" db:"training_duration_seconds"`
	TrainingCompletedAt      *time.Time       `json:"training_completed_at" db:"training_completed_at"`
	LastUpdatedBy            *int             `json:"last_updated_by" db:"last_updated_by"`
	Status                   string           `json:"status" db:"status"` // draft, training, active, inactive, deprecated, archived
	IsProduction             bool             `json:"is_production" db:"is_production"`
	Hyperparameters          *json.RawMessage `json:"hyperparameters" db:"hyperparameters"`
	Metadata                 *json.RawMessage `json:"metadata" db:"metadata"`
	TrainingNotes            *string          `json:"training_notes" db:"training_notes"`
	ModelFilePath            *string          `json:"model_file_path" db:"model_file_path"`
	ModelFileSizeBytes       *int64           `json:"model_file_size_bytes" db:"model_file_size_bytes"`
	NextRetrainingDate       *time.Time       `json:"next_retraining_date" db:"next_retraining_date"`
	RetrainingFrequencyDays  *int             `json:"retraining_frequency_days" db:"retraining_frequency_days"`
	MinDataPointsForTraining *int             `json:"min_data_points_for_training" db:"min_data_points_for_training"`
	CreatedAt                time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt                *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// RECOMMENDATION ENGINE - STRATEGY CONFIGURATIONS
// ============================================================

// RecommendationEngine represents a recommendation strategy configuration
type RecommendationEngine struct {
	ID                           string           `json:"id" db:"id"`
	TenantID                     string           `json:"tenant_id" db:"tenant_id"`
	EngineName                   string           `json:"engine_name" db:"engine_name"`
	Description                  *string          `json:"description" db:"description"`
	RecommendationType           string           `json:"recommendation_type" db:"recommendation_type"` // product, service, lead, property, course, job
	ModelID                      string           `json:"model_id" db:"model_id"`
	StrategyType                 *string          `json:"strategy_type" db:"strategy_type"`
	ScoringAlgorithm             *string          `json:"scoring_algorithm" db:"scoring_algorithm"`
	MinConfidenceThreshold       *float64         `json:"min_confidence_threshold" db:"min_confidence_threshold"`
	MaxRecommendationsPerUser    int              `json:"max_recommendations_per_user" db:"max_recommendations_per_user"`
	MaxRecommendationsPerSession int              `json:"max_recommendations_per_session" db:"max_recommendations_per_session"`
	RecommendationTTLHours       *int             `json:"recommendation_ttl_hours" db:"recommendation_ttl_hours"`
	EnablePersonalization        bool             `json:"enable_personalization" db:"enable_personalization"`
	EnableRealTimeUpdates        bool             `json:"enable_real_time_updates" db:"enable_real_time_updates"`
	EnableABTesting              bool             `json:"enable_a_b_testing" db:"enable_a_b_testing"`
	ABTestVariant                *string          `json:"a_b_test_variant" db:"a_b_test_variant"`
	Filters                      *json.RawMessage `json:"filters" db:"filters"`
	BoostRules                   *json.RawMessage `json:"boost_rules" db:"boost_rules"`
	PenaltyRules                 *json.RawMessage `json:"penalty_rules" db:"penalty_rules"`
	FeedbackEnabled              bool             `json:"feedback_enabled" db:"feedback_enabled"`
	FeedbackCollectionRate       *float64         `json:"feedback_collection_rate" db:"feedback_collection_rate"`
	RankingFactors               *json.RawMessage `json:"ranking_factors" db:"ranking_factors"`
	Status                       string           `json:"status" db:"status"` // draft, testing, active, inactive, archived
	PerformanceMetrics           *json.RawMessage `json:"performance_metrics" db:"performance_metrics"`
	CreatedBy                    *int             `json:"created_by" db:"created_by"`
	LastModifiedBy               *int             `json:"last_modified_by" db:"last_modified_by"`
	CreatedAt                    time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt                    *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// USER RECOMMENDATIONS - PERSONALIZED RECOMMENDATIONS
// ============================================================

// UserRecommendation represents a personalized recommendation for a user
type UserRecommendation struct {
	ID                    string           `json:"id" db:"id"`
	TenantID              string           `json:"tenant_id" db:"tenant_id"`
	UserID                int              `json:"user_id" db:"user_id"`
	EngineID              string           `json:"engine_id" db:"engine_id"`
	ModelID               string           `json:"model_id" db:"model_id"`
	RecommendedItemID     string           `json:"recommended_item_id" db:"recommended_item_id"`
	RecommendedItemType   *string          `json:"recommended_item_type" db:"recommended_item_type"`
	RecommendedItemTitle  *string          `json:"recommended_item_title" db:"recommended_item_title"`
	RecommendationScore   float64          `json:"recommendation_score" db:"recommendation_score"`
	ConfidenceLevel       float64          `json:"confidence_level" db:"confidence_level"`
	RankPosition          int              `json:"rank_position" db:"rank_position"`
	RecommendationReason  *string          `json:"recommendation_reason" db:"recommendation_reason"`
	RecommendationFactors *json.RawMessage `json:"recommendation_factors" db:"recommendation_factors"`
	ClickProbability      *float64         `json:"click_probability" db:"click_probability"`
	ConversionProbability *float64         `json:"conversion_probability" db:"conversion_probability"`
	PersonalizationLevel  *string          `json:"personalization_level" db:"personalization_level"`
	VariantServed         *string          `json:"variant_served" db:"variant_served"`
	ExpiryDate            *time.Time       `json:"expiry_date" db:"expiry_date"`
	DisplayedAt           *time.Time       `json:"displayed_at" db:"displayed_at"`
	ClickedAt             *time.Time       `json:"clicked_at" db:"clicked_at"`
	ConvertedAt           *time.Time       `json:"converted_at" db:"converted_at"`
	UserFeedback          *string          `json:"user_feedback" db:"user_feedback"`
	FeedbackReason        *string          `json:"feedback_reason" db:"feedback_reason"`
	FeedbackTimestamp     *time.Time       `json:"feedback_timestamp" db:"feedback_timestamp"`
	Metadata              *json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt             time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt             *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// RECOMMENDATION FEEDBACK
// ============================================================

// RecommendationFeedback represents user feedback on a recommendation
type RecommendationFeedback struct {
	ID                     string           `json:"id" db:"id"`
	TenantID               string           `json:"tenant_id" db:"tenant_id"`
	UserID                 int              `json:"user_id" db:"user_id"`
	RecommendationID       string           `json:"recommendation_id" db:"recommendation_id"`
	EngineID               string           `json:"engine_id" db:"engine_id"`
	FeedbackType           string           `json:"feedback_type" db:"feedback_type"` // explicit, implicit, conversion
	FeedbackValue          string           `json:"feedback_value" db:"feedback_value"`
	RatingScore            *int             `json:"rating_score" db:"rating_score"`
	DetailedFeedback       *string          `json:"detailed_feedback" db:"detailed_feedback"`
	DeviceType             *string          `json:"device_type" db:"device_type"`
	UserContext            *json.RawMessage `json:"user_context" db:"user_context"`
	RecommendationContext  *json.RawMessage `json:"recommendation_context" db:"recommendation_context"`
	ExternalEventTriggered *string          `json:"external_event_triggered" db:"external_event_triggered"`
	SessionID              *string          `json:"session_id" db:"session_id"`
	ResponseTimeSeconds    *int             `json:"response_time_seconds" db:"response_time_seconds"`
	CreatedAt              time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt              *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// ML TRAINING DATA
// ============================================================

// MLTrainingData represents a dataset for model training
type MLTrainingData struct {
	ID                      string           `json:"id" db:"id"`
	TenantID                string           `json:"tenant_id" db:"tenant_id"`
	ModelID                 string           `json:"model_id" db:"model_id"`
	DatasetName             string           `json:"dataset_name" db:"dataset_name"`
	DatasetVersion          *string          `json:"dataset_version" db:"dataset_version"`
	DataType                *string          `json:"data_type" db:"data_type"`
	SourceTable             *string          `json:"source_table" db:"source_table"`
	TotalRecords            *int             `json:"total_records" db:"total_records"`
	ValidRecords            *int             `json:"valid_records" db:"valid_records"`
	InvalidRecords          *int             `json:"invalid_records" db:"invalid_records"`
	DuplicateRecords        *int             `json:"duplicate_records" db:"duplicate_records"`
	MissingValuePercentage  *float64         `json:"missing_value_percentage" db:"missing_value_percentage"`
	OutliersDetected        *int             `json:"outliers_detected" db:"outliers_detected"`
	OutliersRemoved         *int             `json:"outliers_removed" db:"outliers_removed"`
	FeatureScalingMethod    *string          `json:"feature_scaling_method" db:"feature_scaling_method"`
	FeatureEncodingMethod   *string          `json:"feature_encoding_method" db:"feature_encoding_method"`
	TrainTestSplitRatio     *string          `json:"train_test_split_ratio" db:"train_test_split_ratio"`
	ValidationSetRatio      *float64         `json:"validation_set_ratio" db:"validation_set_ratio"`
	ClassImbalanceHandling  *string          `json:"class_imbalance_handling" db:"class_imbalance_handling"`
	TemporalSplitDate       *time.Time       `json:"temporal_split_date" db:"temporal_split_date"`
	DataCollectionStartDate *time.Time       `json:"data_collection_start_date" db:"data_collection_start_date"`
	DataCollectionEndDate   *time.Time       `json:"data_collection_end_date" db:"data_collection_end_date"`
	PreprocessingSteps      *json.RawMessage `json:"preprocessing_steps" db:"preprocessing_steps"`
	QualityScore            *float64         `json:"quality_score" db:"quality_score"`
	FilePath                *string          `json:"file_path" db:"file_path"`
	FileSizeMB              *float64         `json:"file_size_mb" db:"file_size_mb"`
	HashValue               *string          `json:"hash_value" db:"hash_value"`
	CreatedBy               *int             `json:"created_by" db:"created_by"`
	CreatedAt               time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// PREDICTION RESULTS
// ============================================================

// PredictionResult represents a machine learning prediction output
type PredictionResult struct {
	ID                      string           `json:"id" db:"id"`
	TenantID                string           `json:"tenant_id" db:"tenant_id"`
	ModelID                 string           `json:"model_id" db:"model_id"`
	InputDataID             *string          `json:"input_data_id" db:"input_data_id"`
	InputDataType           *string          `json:"input_data_type" db:"input_data_type"`
	PredictionType          *string          `json:"prediction_type" db:"prediction_type"`
	PredictedValue          *string          `json:"predicted_value" db:"predicted_value"`
	PredictedClass          *string          `json:"predicted_class" db:"predicted_class"`
	ConfidenceScore         *float64         `json:"confidence_score" db:"confidence_score"`
	ProbabilityDistribution *json.RawMessage `json:"probability_distribution" db:"probability_distribution"`
	FeatureImportance       *json.RawMessage `json:"feature_importance" db:"feature_importance"`
	ProcessingTimeMs        *int             `json:"processing_time_ms" db:"processing_time_ms"`
	PredictionMethod        *string          `json:"prediction_method" db:"prediction_method"`
	BatchID                 *string          `json:"batch_id" db:"batch_id"`
	ExecutionEnvironment    *string          `json:"execution_environment" db:"execution_environment"`
	ModelInferenceVersion   *string          `json:"model_inference_version" db:"model_inference_version"`
	InputFeaturesHash       *string          `json:"input_features_hash" db:"input_features_hash"`
	DriftDetected           bool             `json:"drift_detected" db:"drift_detected"`
	DriftMagnitude          *float64         `json:"drift_magnitude" db:"drift_magnitude"`
	ActualOutcome           *string          `json:"actual_outcome" db:"actual_outcome"`
	ActualOutcomeDate       *time.Time       `json:"actual_outcome_date" db:"actual_outcome_date"`
	PredictionAccuracy      *bool            `json:"prediction_accuracy" db:"prediction_accuracy"`
	FeedbackLoopTriggered   bool             `json:"feedback_loop_triggered" db:"feedback_loop_triggered"`
	Metadata                *json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt               time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt               *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// RECOMMENDATION HISTORY
// ============================================================

// RecommendationHistory represents historical recommendations for a user
type RecommendationHistory struct {
	ID                            string           `json:"id" db:"id"`
	TenantID                      string           `json:"tenant_id" db:"tenant_id"`
	UserID                        int              `json:"user_id" db:"user_id"`
	EngineID                      string           `json:"engine_id" db:"engine_id"`
	BatchTimestamp                *time.Time       `json:"batch_timestamp" db:"batch_timestamp"`
	RecommendationBatchID         *string          `json:"recommendation_batch_id" db:"recommendation_batch_id"`
	TotalRecommendationsShown     *int             `json:"total_recommendations_shown" db:"total_recommendations_shown"`
	TotalRecommendationsClicked   *int             `json:"total_recommendations_clicked" db:"total_recommendations_clicked"`
	TotalRecommendationsConverted *int             `json:"total_recommendations_converted" db:"total_recommendations_converted"`
	CTR                           *float64         `json:"ctr" db:"ctr"`
	ConversionRate                *float64         `json:"conversion_rate" db:"conversion_rate"`
	AvgRecommendationScore        *float64         `json:"avg_recommendation_score" db:"avg_recommendation_score"`
	SessionDurationSeconds        *int             `json:"session_duration_seconds" db:"session_duration_seconds"`
	DeviceType                    *string          `json:"device_type" db:"device_type"`
	UserSegment                   *string          `json:"user_segment" db:"user_segment"`
	PersonalizationUsed           *bool            `json:"personalization_used" db:"personalization_used"`
	VariantServed                 *string          `json:"variant_served" db:"variant_served"`
	ABTestGroup                   *string          `json:"ab_test_group" db:"ab_test_group"`
	TopClickedItemID              *string          `json:"top_clicked_item_id" db:"top_clicked_item_id"`
	TopClickedItemRank            *int             `json:"top_clicked_item_rank" db:"top_clicked_item_rank"`
	TopConvertedItemID            *string          `json:"top_converted_item_id" db:"top_converted_item_id"`
	TopConvertedItemRank          *int             `json:"top_converted_item_rank" db:"top_converted_item_rank"`
	UniqueItemsRecommended        *int             `json:"unique_items_recommended" db:"unique_items_recommended"`
	RepeatRecommendations         *int             `json:"repeat_recommendations" db:"repeat_recommendations"`
	UserSegmentAtTime             *string          `json:"user_segment_at_time" db:"user_segment_at_time"`
	ContextualFactors             *json.RawMessage `json:"contextual_factors" db:"contextual_factors"`
	SystemPerformance             *json.RawMessage `json:"system_performance" db:"system_performance"`
	CreatedAt                     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt                     *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// ANOMALY DETECTION
// ============================================================

// AnomalyDetection represents a detected anomaly
type AnomalyDetection struct {
	ID                    string           `json:"id" db:"id"`
	TenantID              string           `json:"tenant_id" db:"tenant_id"`
	AnomalyType           string           `json:"anomaly_type" db:"anomaly_type"`
	SeverityLevel         string           `json:"severity_level" db:"severity_level"`
	AffectedModelID       *string          `json:"affected_model_id" db:"affected_model_id"`
	AffectedEntityID      *string          `json:"affected_entity_id" db:"affected_entity_id"`
	AffectedEntityType    *string          `json:"affected_entity_type" db:"affected_entity_type"`
	AnomalyDescription    *string          `json:"anomaly_description" db:"anomaly_description"`
	AnomalyMetric         *string          `json:"anomaly_metric" db:"anomaly_metric"`
	ExpectedValue         *float64         `json:"expected_value" db:"expected_value"`
	ActualValue           *float64         `json:"actual_value" db:"actual_value"`
	DeviationPercentage   *float64         `json:"deviation_percentage" db:"deviation_percentage"`
	ZScore                *float64         `json:"z_score" db:"z_score"`
	DetectionMethod       *string          `json:"detection_method" db:"detection_method"`
	DetectionConfidence   *float64         `json:"detection_confidence" db:"detection_confidence"`
	RootCauseAnalysis     *string          `json:"root_cause_analysis" db:"root_cause_analysis"`
	RecommendedAction     *string          `json:"recommended_action" db:"recommended_action"`
	IsAcknowledged        bool             `json:"is_acknowledged" db:"is_acknowledged"`
	AcknowledgedBy        *int             `json:"acknowledged_by" db:"acknowledged_by"`
	AcknowledgedAt        *time.Time       `json:"acknowledged_at" db:"acknowledged_at"`
	AcknowledgedNotes     *string          `json:"acknowledged_notes" db:"acknowledged_notes"`
	IsResolved            bool             `json:"is_resolved" db:"is_resolved"`
	ResolvedBy            *int             `json:"resolved_by" db:"resolved_by"`
	ResolvedAt            *time.Time       `json:"resolved_at" db:"resolved_at"`
	ResolutionActionTaken *string          `json:"resolution_action_taken" db:"resolution_action_taken"`
	AffectedMetrics       *json.RawMessage `json:"affected_metrics" db:"affected_metrics"`
	RelatedAnomalies      *json.RawMessage `json:"related_anomalies" db:"related_anomalies"`
	CreatedAt             time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt             *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// AI INSIGHTS
// ============================================================

// AIInsight represents an AI-generated business insight
type AIInsight struct {
	ID                     string           `json:"id" db:"id"`
	TenantID               string           `json:"tenant_id" db:"tenant_id"`
	InsightType            string           `json:"insight_type" db:"insight_type"`
	InsightCategory        *string          `json:"insight_category" db:"insight_category"`
	Title                  string           `json:"title" db:"title"`
	Description            string           `json:"description" db:"description"`
	GeneratedByModelID     *string          `json:"generated_by_model_id" db:"generated_by_model_id"`
	ConfidenceScore        *float64         `json:"confidence_score" db:"confidence_score"`
	ImpactScore            *float64         `json:"impact_score" db:"impact_score"`
	RelevanceScore         *float64         `json:"relevance_score" db:"relevance_score"`
	SupportingData         *json.RawMessage `json:"supporting_data" db:"supporting_data"`
	Visualizations         *json.RawMessage `json:"visualizations" db:"visualizations"`
	Recommendation         *string          `json:"recommendation" db:"recommendation"`
	AffectedKPIs           *json.RawMessage `json:"affected_kpis" db:"affected_kpis"`
	TimePeriod             *string          `json:"time_period" db:"time_period"`
	TrendDirection         *string          `json:"trend_direction" db:"trend_direction"`
	ForecastPeriodDays     *int             `json:"forecast_period_days" db:"forecast_period_days"`
	HistoricalAccuracyRate *float64         `json:"historical_accuracy_rate" db:"historical_accuracy_rate"`
	Audience               *string          `json:"audience" db:"audience"`
	PriorityLevel          *string          `json:"priority_level" db:"priority_level"`
	ActionItems            *json.RawMessage `json:"action_items" db:"action_items"`
	CreatedBy              *string          `json:"created_by" db:"created_by"`
	IsShared               bool             `json:"is_shared" db:"is_shared"`
	SharedWithUserIDs      *json.RawMessage `json:"shared_with_user_ids" db:"shared_with_user_ids"`
	IsArchived             bool             `json:"is_archived" db:"is_archived"`
	ArchivedReason         *string          `json:"archived_reason" db:"archived_reason"`
	FeedbackScore          *float64         `json:"feedback_score" db:"feedback_score"`
	FeedbackCount          *int             `json:"feedback_count" db:"feedback_count"`
	CreatedAt              time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt              *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// MODEL PERFORMANCE
// ============================================================

// ModelPerformance represents ML model performance metrics
type ModelPerformance struct {
	ID                          string           `json:"id" db:"id"`
	TenantID                    string           `json:"tenant_id" db:"tenant_id"`
	ModelID                     string           `json:"model_id" db:"model_id"`
	PerformanceDate             time.Time        `json:"performance_date" db:"performance_date"`
	TotalPredictions            *int             `json:"total_predictions" db:"total_predictions"`
	SuccessfulPredictions       *int             `json:"successful_predictions" db:"successful_predictions"`
	FailedPredictions           *int             `json:"failed_predictions" db:"failed_predictions"`
	Accuracy                    *float64         `json:"accuracy" db:"accuracy"`
	Precision                   *float64         `json:"precision" db:"precision"`
	Recall                      *float64         `json:"recall" db:"recall"`
	F1Score                     *float64         `json:"f1_score" db:"f1_score"`
	AUCROC                      *float64         `json:"auc_roc" db:"auc_roc"`
	LogLoss                     *float64         `json:"log_loss" db:"log_loss"`
	MeanAbsoluteError           *float64         `json:"mean_absolute_error" db:"mean_absolute_error"`
	MeanSquaredError            *float64         `json:"mean_squared_error" db:"mean_squared_error"`
	RootMeanSquaredError        *float64         `json:"root_mean_squared_error" db:"root_mean_squared_error"`
	AvgInferenceTimeMs          *float64         `json:"avg_inference_time_ms" db:"avg_inference_time_ms"`
	P95InferenceTimeMs          *float64         `json:"p95_inference_time_ms" db:"p95_inference_time_ms"`
	P99InferenceTimeMs          *float64         `json:"p99_inference_time_ms" db:"p99_inference_time_ms"`
	ThroughputPredictionsPerSec *float64         `json:"throughput_predictions_per_second" db:"throughput_predictions_per_second"`
	ErrorRate                   *float64         `json:"error_rate" db:"error_rate"`
	PredictionDistribution      *json.RawMessage `json:"prediction_distribution" db:"prediction_distribution"`
	FeatureImportanceSnapshot   *json.RawMessage `json:"feature_importance_snapshot" db:"feature_importance_snapshot"`
	DataDriftScore              *float64         `json:"data_drift_score" db:"data_drift_score"`
	ModelDegradationDetected    bool             `json:"model_degradation_detected" db:"model_degradation_detected"`
	DegradationReason           *string          `json:"degradation_reason" db:"degradation_reason"`
	RetrainingRecommended       bool             `json:"retraining_recommended" db:"retraining_recommended"`
	ConfidenceIntervalLower     *float64         `json:"confidence_interval_lower" db:"confidence_interval_lower"`
	ConfidenceIntervalUpper     *float64         `json:"confidence_interval_upper" db:"confidence_interval_upper"`
	SegmentPerformance          *json.RawMessage `json:"segment_performance" db:"segment_performance"`
	ErrorAnalysis               *json.RawMessage `json:"error_analysis" db:"error_analysis"`
	FalsePositiveRate           *float64         `json:"false_positive_rate" db:"false_positive_rate"`
	FalseNegativeRate           *float64         `json:"false_negative_rate" db:"false_negative_rate"`
	BaselineModelComparison     *json.RawMessage `json:"baseline_model_comparison" db:"baseline_model_comparison"`
	Notes                       *string          `json:"notes" db:"notes"`
	CreatedAt                   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt                   time.Time        `json:"updated_at" db:"updated_at"`
	DeletedAt                   *time.Time       `json:"deleted_at" db:"deleted_at"`
}

// ============================================================
// DTOs FOR API REQUESTS/RESPONSES
// ============================================================

// CreateAIModelRequest represents request to create an AI model
type CreateAIModelRequest struct {
	ModelName     string          `json:"model_name"`
	ModelType     string          `json:"model_type"`
	ModelVersion  string          `json:"model_version"`
	Description   *string         `json:"description"`
	AlgorithmName *string         `json:"algorithm_name"`
	InputFeatures json.RawMessage `json:"input_features"`
	OutputFormat  *string         `json:"output_format"`
}

// CreateRecommendationEngineRequest represents request to create a recommendation engine
type CreateRecommendationEngineRequest struct {
	EngineName                string   `json:"engine_name"`
	Description               *string  `json:"description"`
	RecommendationType        string   `json:"recommendation_type"`
	ModelID                   string   `json:"model_id"`
	StrategyType              *string  `json:"strategy_type"`
	ScoringAlgorithm          *string  `json:"scoring_algorithm"`
	MinConfidenceThreshold    *float64 `json:"min_confidence_threshold"`
	MaxRecommendationsPerUser int      `json:"max_recommendations_per_user"`
}

// GetUserRecommendationsRequest represents request to get user recommendations
type GetUserRecommendationsRequest struct {
	UserID   int      `json:"user_id"`
	EngineID string   `json:"engine_id"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	MinScore *float64 `json:"min_score"`
}

// SubmitFeedbackRequest represents request to submit recommendation feedback
type SubmitFeedbackRequest struct {
	RecommendationID string  `json:"recommendation_id"`
	FeedbackValue    string  `json:"feedback_value"`
	RatingScore      *int    `json:"rating_score"`
	DetailedFeedback *string `json:"detailed_feedback"`
}

// GenerateInsightsRequest represents request to generate AI insights
type GenerateInsightsRequest struct {
	InsightType     string  `json:"insight_type"`
	InsightCategory *string `json:"insight_category"`
	TimePeriod      *string `json:"time_period"`
}

// TrainModelRequest represents request to train a model
type TrainModelRequest struct {
	ModelID       string  `json:"model_id"`
	DatasetID     string  `json:"dataset_id"`
	TrainingNotes *string `json:"training_notes"`
}

// DetectAnomaliesRequest represents request to detect anomalies
type DetectAnomaliesRequest struct {
	AnomalyType string  `json:"anomaly_type"`
	ModelID     *string `json:"model_id"`
	DateRange   *string `json:"date_range"`
}

// APIResponse is a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
