-- ============================================================
-- MIGRATION 041: AI-POWERED RECOMMENDATIONS
-- Date: December 23, 2025
-- Purpose: AI/ML infrastructure for intelligent recommendations,
--          prediction, anomaly detection, and business insights
-- ============================================================

SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- AI MODELS - ML MODEL CONFIGURATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `ai_models` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_name` VARCHAR(255) NOT NULL,
    `model_type` VARCHAR(100) NOT NULL, -- recommendation, prediction, clustering, classification, anomaly_detection
    `model_version` VARCHAR(50) NOT NULL,
    `description` TEXT,
    `algorithm_name` VARCHAR(255), -- e.g., collaborative_filtering, gradient_boosting, neural_network
    `input_features` JSON NOT NULL, -- Array of feature names used by model
    `output_format` VARCHAR(100), -- json, float, int, string, array
    `accuracy_score` DECIMAL(5,2), -- Percentage 0-100
    `precision_score` DECIMAL(5,2),
    `recall_score` DECIMAL(5,2),
    `f1_score` DECIMAL(5,2),
    `training_dataset_size` INT,
    `training_duration_seconds` INT,
    `training_completed_at` TIMESTAMP NULL,
    `last_updated_by` INT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, training, active, inactive, deprecated, archived
    `is_production` BOOLEAN DEFAULT 0,
    `hyperparameters` JSON, -- Model tuning parameters
    `metadata` JSON,
    `training_notes` TEXT,
    `model_file_path` VARCHAR(500), -- Path to serialized model file
    `model_file_size_bytes` BIGINT,
    `next_retraining_date` TIMESTAMP NULL,
    `retraining_frequency_days` INT DEFAULT 30,
    `min_data_points_for_training` INT DEFAULT 100,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_type` (`model_type`),
    KEY `idx_status` (`status`),
    KEY `idx_is_production` (`is_production`),
    KEY `idx_created_at` (`created_at`),
    UNIQUE KEY `uk_tenant_model_version` (`tenant_id`, `model_name`, `model_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- RECOMMENDATION ENGINE - STRATEGY CONFIGURATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `recommendation_engine` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `engine_name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `recommendation_type` VARCHAR(100) NOT NULL, -- product, service, lead, property, course, job
    `model_id` VARCHAR(36) NOT NULL,
    `strategy_type` VARCHAR(100), -- collaborative_filtering, content_based, hybrid, rule_based
    `scoring_algorithm` VARCHAR(100), -- weighted_score, percentile_rank, z_score
    `min_confidence_threshold` DECIMAL(5,3), -- 0-1, minimum confidence to show recommendation
    `max_recommendations_per_user` INT DEFAULT 10,
    `max_recommendations_per_session` INT DEFAULT 5,
    `recommendation_ttl_hours` INT DEFAULT 24, -- Time to live for cached recommendations
    `enable_personalization` BOOLEAN DEFAULT 1,
    `enable_real_time_updates` BOOLEAN DEFAULT 1,
    `enable_a_b_testing` BOOLEAN DEFAULT 0,
    `a_b_test_variant` VARCHAR(50), -- control, variant_a, variant_b
    `filters` JSON, -- Filtering criteria: exclude_purchased, exclude_viewed, category_whitelist, etc.
    `boost_rules` JSON, -- Boosting factors for certain categories/items
    `penalty_rules` JSON, -- Penalty factors for certain categories/items
    `feedback_enabled` BOOLEAN DEFAULT 1,
    `feedback_collection_rate` DECIMAL(3,2) DEFAULT 1.0, -- 0-1, sample rate for feedback
    `ranking_factors` JSON, -- Weights for different ranking factors
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, testing, active, inactive, archived
    `performance_metrics` JSON, -- CTR, conversion_rate, user_satisfaction, etc.
    `created_by` INT,
    `last_modified_by` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`model_id`) REFERENCES `ai_models`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_recommendation_type` (`recommendation_type`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    UNIQUE KEY `uk_tenant_engine_name` (`tenant_id`, `engine_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- USER RECOMMENDATIONS - PERSONALIZED RECOMMENDATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `user_recommendations` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `engine_id` VARCHAR(36) NOT NULL,
    `model_id` VARCHAR(36) NOT NULL,
    `recommended_item_id` VARCHAR(255) NOT NULL, -- ID of recommended product/service/property/lead
    `recommended_item_type` VARCHAR(100), -- product, service, lead, property, course, job
    `recommended_item_title` VARCHAR(500),
    `recommendation_score` DECIMAL(10,4), -- 0-100 or 0-1 depending on model
    `confidence_level` DECIMAL(5,3), -- 0-1
    `rank_position` INT, -- Position in the recommendation list (1st, 2nd, etc.)
    `recommendation_reason` TEXT, -- Why this was recommended (e.g., "Users similar to you liked this")
    `recommendation_factors` JSON, -- Contributing factors: {factor: weight}
    `click_probability` DECIMAL(5,3), -- Predicted probability of click
    `conversion_probability` DECIMAL(5,3), -- Predicted probability of conversion
    `personalization_level` VARCHAR(50), -- high, medium, low
    `variant_served` VARCHAR(50), -- For A/B testing
    `expiry_date` TIMESTAMP NULL,
    `displayed_at` TIMESTAMP NULL,
    `clicked_at` TIMESTAMP NULL,
    `converted_at` TIMESTAMP NULL,
    `user_feedback` VARCHAR(50), -- helpful, not_helpful, irrelevant, seen_before, not_interested
    `feedback_reason` VARCHAR(255),
    `feedback_timestamp` TIMESTAMP NULL,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`engine_id`) REFERENCES `recommendation_engine`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`model_id`) REFERENCES `ai_models`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_engine_id` (`engine_id`),
    KEY `idx_model_id` (`model_id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_expiry_date` (`expiry_date`),
    KEY `idx_user_feedback` (`user_feedback`),
    KEY `idx_display_clicked` (`displayed_at`, `clicked_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- RECOMMENDATION FEEDBACK - USER FEEDBACK
-- ============================================================
CREATE TABLE IF NOT EXISTS `recommendation_feedback` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `recommendation_id` VARCHAR(36) NOT NULL,
    `engine_id` VARCHAR(36) NOT NULL,
    `feedback_type` VARCHAR(50) NOT NULL, -- explicit (user clicked/rated), implicit (view duration, click), conversion (actually purchased)
    `feedback_value` VARCHAR(100), -- helpful, not_helpful, irrelevant, interesting, purchased, viewed, skipped
    `rating_score` TINYINT, -- 1-5 star rating, if applicable
    `detailed_feedback` TEXT, -- User's written feedback
    `device_type` VARCHAR(50), -- mobile, tablet, desktop
    `user_context` JSON, -- Context: time_of_day, day_of_week, user_segment, etc.
    `recommendation_context` JSON, -- Context: position_in_list, variant, personalization_level
    `external_event_triggered` VARCHAR(255), -- External event that led to feedback (e.g., "purchase_completed")
    `session_id` VARCHAR(255),
    `response_time_seconds` INT, -- How long user took to provide feedback
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`recommendation_id`) REFERENCES `user_recommendations`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`engine_id`) REFERENCES `recommendation_engine`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_recommendation_id` (`recommendation_id`),
    KEY `idx_engine_id` (`engine_id`),
    KEY `idx_feedback_type` (`feedback_type`),
    KEY `idx_feedback_value` (`feedback_value`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ML TRAINING DATA - DATASETS FOR MODEL TRAINING
-- ============================================================
CREATE TABLE IF NOT EXISTS `ml_training_data` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_id` VARCHAR(36) NOT NULL,
    `dataset_name` VARCHAR(255) NOT NULL,
    `dataset_version` VARCHAR(50),
    `data_type` VARCHAR(50), -- historical, synthetic, augmented, external
    `source_table` VARCHAR(255), -- Original table where data comes from
    `total_records` INT,
    `valid_records` INT,
    `invalid_records` INT,
    `duplicate_records` INT,
    `missing_value_percentage` DECIMAL(5,2),
    `outliers_detected` INT,
    `outliers_removed` INT,
    `feature_scaling_method` VARCHAR(100), -- standardization, normalization, robust_scaling
    `feature_encoding_method` VARCHAR(100), -- one_hot, label, ordinal, target
    `train_test_split_ratio` VARCHAR(20) DEFAULT '0.8:0.2', -- e.g., "0.8:0.2"
    `validation_set_ratio` DECIMAL(3,2), -- e.g., 0.1 for 10%
    `class_imbalance_handling` VARCHAR(100), -- oversampling, undersampling, smote, weighted_loss
    `temporal_split_date` DATE, -- For time-series: data before this for training
    `data_collection_start_date` DATE,
    `data_collection_end_date` DATE,
    `preprocessing_steps` JSON, -- Array of preprocessing steps applied
    `quality_score` DECIMAL(5,2), -- 0-100
    `file_path` VARCHAR(500),
    `file_size_mb` DECIMAL(10,2),
    `hash_value` VARCHAR(255), -- MD5/SHA for data integrity
    `created_by` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`model_id`) REFERENCES `ai_models`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`),
    KEY `idx_created_at` (`created_at`),
    UNIQUE KEY `uk_tenant_dataset_version` (`tenant_id`, `dataset_name`, `dataset_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- PREDICTION RESULTS - ML PREDICTION OUTPUTS
-- ============================================================
CREATE TABLE IF NOT EXISTS `prediction_results` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_id` VARCHAR(36) NOT NULL,
    `input_data_id` VARCHAR(255), -- ID of the input record (lead, customer, property, etc.)
    `input_data_type` VARCHAR(100), -- lead, customer, property, transaction
    `prediction_type` VARCHAR(100), -- conversion_probability, churn_risk, property_price, lead_score
    `predicted_value` VARCHAR(500), -- The prediction result
    `predicted_class` VARCHAR(100), -- For classification: the class label
    `confidence_score` DECIMAL(5,3), -- 0-1
    `probability_distribution` JSON, -- For multi-class: {class: probability}
    `feature_importance` JSON, -- Most important features: {feature: importance_score}
    `processing_time_ms` INT, -- How long prediction took
    `prediction_method` VARCHAR(100), -- batch, real_time, scheduled
    `batch_id` VARCHAR(255), -- If part of batch prediction
    `execution_environment` VARCHAR(50), -- cpu, gpu, tpu
    `model_inference_version` VARCHAR(50), -- Version of model used
    `input_features_hash` VARCHAR(255), -- Hash of input features for validation
    `drift_detected` BOOLEAN DEFAULT 0, -- Whether data drift detected
    `drift_magnitude` DECIMAL(5,3), -- 0-1, measure of drift
    `actual_outcome` VARCHAR(500), -- Actual value if later observed
    `actual_outcome_date` TIMESTAMP NULL,
    `prediction_accuracy` BOOLEAN, -- True if prediction matched actual
    `feedback_loop_triggered` BOOLEAN DEFAULT 0,
    `metadata` JSON,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`model_id`) REFERENCES `ai_models`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`),
    KEY `idx_input_data_id` (`input_data_id`),
    KEY `idx_prediction_type` (`prediction_type`),
    KEY `idx_confidence_score` (`confidence_score`),
    KEY `idx_drift_detected` (`drift_detected`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- RECOMMENDATION HISTORY - HISTORICAL RECOMMENDATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS `recommendation_history` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` INT NOT NULL,
    `engine_id` VARCHAR(36) NOT NULL,
    `batch_timestamp` TIMESTAMP,
    `recommendation_batch_id` VARCHAR(255), -- Group recommendations generated in same batch
    `total_recommendations_shown` INT,
    `total_recommendations_clicked` INT,
    `total_recommendations_converted` INT,
    `ctr` DECIMAL(5,3), -- Click-through rate
    `conversion_rate` DECIMAL(5,3),
    `avg_recommendation_score` DECIMAL(10,4),
    `session_duration_seconds` INT,
    `device_type` VARCHAR(50),
    `user_segment` VARCHAR(100),
    `personalization_used` BOOLEAN,
    `variant_served` VARCHAR(50),
    `ab_test_group` VARCHAR(50),
    `top_clicked_item_id` VARCHAR(255),
    `top_clicked_item_rank` INT,
    `top_converted_item_id` VARCHAR(255),
    `top_converted_item_rank` INT,
    `unique_items_recommended` INT,
    `repeat_recommendations` INT,
    `user_segment_at_time` VARCHAR(100),
    `contextual_factors` JSON, -- Time of day, day of week, etc.
    `system_performance` JSON, -- Response time, recommendations generated, etc.
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`engine_id`) REFERENCES `recommendation_engine`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_engine_id` (`engine_id`),
    KEY `idx_batch_timestamp` (`batch_timestamp`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- ANOMALY DETECTION - DETECTED ANOMALIES IN DATA
-- ============================================================
CREATE TABLE IF NOT EXISTS `anomaly_detection` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `anomaly_type` VARCHAR(100) NOT NULL, -- data_drift, performance_degradation, outlier, fraud, system_anomaly
    `severity_level` VARCHAR(50) NOT NULL, -- critical, high, medium, low, info
    `affected_model_id` VARCHAR(36),
    `affected_entity_id` VARCHAR(255), -- ID of the entity with anomaly
    `affected_entity_type` VARCHAR(100), -- lead, customer, property, transaction
    `anomaly_description` TEXT,
    `anomaly_metric` VARCHAR(255), -- e.g., "prediction_confidence", "conversion_rate", "user_activity"
    `expected_value` DECIMAL(15,4),
    `actual_value` DECIMAL(15,4),
    `deviation_percentage` DECIMAL(10,2), -- How much deviation from expected
    `z_score` DECIMAL(10,4), -- Statistical z-score
    `detection_method` VARCHAR(100), -- isolation_forest, statistical, rule_based, ml_model
    `detection_confidence` DECIMAL(5,3),
    `root_cause_analysis` TEXT,
    `recommended_action` TEXT,
    `is_acknowledged` BOOLEAN DEFAULT 0,
    `acknowledged_by` INT,
    `acknowledged_at` TIMESTAMP NULL,
    `acknowledged_notes` TEXT,
    `is_resolved` BOOLEAN DEFAULT 0,
    `resolved_by` INT,
    `resolved_at` TIMESTAMP NULL,
    `resolution_action_taken` TEXT,
    `affected_metrics` JSON, -- Metrics impacted by this anomaly
    `related_anomalies` JSON, -- IDs of related anomalies
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`affected_model_id`) REFERENCES `ai_models`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`acknowledged_by`) REFERENCES `user`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`resolved_by`) REFERENCES `user`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_anomaly_type` (`anomaly_type`),
    KEY `idx_severity_level` (`severity_level`),
    KEY `idx_affected_model_id` (`affected_model_id`),
    KEY `idx_is_acknowledged` (`is_acknowledged`),
    KEY `idx_is_resolved` (`is_resolved`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- AI INSIGHTS - AI-GENERATED BUSINESS INSIGHTS
-- ============================================================
CREATE TABLE IF NOT EXISTS `ai_insights` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `insight_type` VARCHAR(100) NOT NULL, -- trend, pattern, opportunity, risk, anomaly, forecast
    `insight_category` VARCHAR(100), -- sales, marketing, operations, finance, hr, customer
    `title` VARCHAR(255) NOT NULL,
    `description` TEXT NOT NULL,
    `generated_by_model_id` VARCHAR(36),
    `confidence_score` DECIMAL(5,3), -- 0-1, how confident in the insight
    `impact_score` DECIMAL(5,3), -- 0-1, potential business impact
    `relevance_score` DECIMAL(5,3), -- 0-1, relevance to the tenant
    `supporting_data` JSON, -- Key numbers, metrics, evidence
    `visualizations` JSON, -- Chart configurations for visualization
    `recommendation` TEXT, -- What action to take based on insight
    `affected_kpis` JSON, -- Array of KPIs affected by this insight
    `time_period` VARCHAR(100), -- Period the insight covers: "last_7_days", "last_30_days", "year_to_date"
    `trend_direction` VARCHAR(50), -- upward, downward, stable, volatile
    `forecast_period_days` INT, -- Days into future for forecasts
    `historical_accuracy_rate` DECIMAL(5,2), -- % of similar insights that proved accurate
    `audience` VARCHAR(100), -- who should see this: all_users, managers, executives, data_team
    `priority_level` VARCHAR(50), -- critical, high, medium, low
    `action_items` JSON, -- Suggested follow-up actions
    `created_by` VARCHAR(36), -- User or system that created insight
    `is_shared` BOOLEAN DEFAULT 0,
    `shared_with_user_ids` JSON, -- Array of user IDs shared with
    `is_archived` BOOLEAN DEFAULT 0,
    `archived_reason` VARCHAR(255),
    `feedback_score` DECIMAL(3,2), -- User feedback: 1-5
    `feedback_count` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`generated_by_model_id`) REFERENCES `ai_models`(`id`) ON DELETE SET NULL,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_insight_type` (`insight_type`),
    KEY `idx_insight_category` (`insight_category`),
    KEY `idx_confidence_score` (`confidence_score`),
    KEY `idx_impact_score` (`impact_score`),
    KEY `idx_priority_level` (`priority_level`),
    KEY `idx_is_archived` (`is_archived`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- MODEL PERFORMANCE - ML MODEL PERFORMANCE METRICS
-- ============================================================
CREATE TABLE IF NOT EXISTS `model_performance` (
    `id` CHAR(36) PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `model_id` VARCHAR(36) NOT NULL,
    `performance_date` DATE NOT NULL,
    `total_predictions` INT,
    `successful_predictions` INT,
    `failed_predictions` INT,
    `accuracy` DECIMAL(5,3), -- 0-1
    `precision` DECIMAL(5,3), -- 0-1
    `recall` DECIMAL(5,3), -- 0-1
    `f1_score` DECIMAL(5,3), -- 0-1
    `auc_roc` DECIMAL(5,3), -- 0-1, Area under ROC curve
    `log_loss` DECIMAL(10,4),
    `mean_absolute_error` DECIMAL(10,4),
    `mean_squared_error` DECIMAL(10,4),
    `root_mean_squared_error` DECIMAL(10,4),
    `avg_inference_time_ms` DECIMAL(10,2),
    `p95_inference_time_ms` DECIMAL(10,2),
    `p99_inference_time_ms` DECIMAL(10,2),
    `throughput_predictions_per_second` DECIMAL(10,2),
    `error_rate` DECIMAL(5,3), -- 0-1
    `prediction_distribution` JSON, -- Distribution of predictions
    `feature_importance_snapshot` JSON, -- Top features at this time
    `data_drift_score` DECIMAL(5,3), -- 0-1, measure of data drift
    `model_degradation_detected` BOOLEAN DEFAULT 0,
    `degradation_reason` VARCHAR(255),
    `retraining_recommended` BOOLEAN DEFAULT 0,
    `confidence_interval_lower` DECIMAL(5,3),
    `confidence_interval_upper` DECIMAL(5,3),
    `segment_performance` JSON, -- Performance by user segment, category, etc.
    `error_analysis` JSON, -- Types of errors and their frequencies
    `false_positive_rate` DECIMAL(5,3),
    `false_negative_rate` DECIMAL(5,3),
    `baseline_model_comparison` JSON, -- Comparison with baseline model
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`tenant_id`) REFERENCES `tenant`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`model_id`) REFERENCES `ai_models`(`id`) ON DELETE CASCADE,
    KEY `idx_tenant_id` (`tenant_id`),
    KEY `idx_model_id` (`model_id`),
    KEY `idx_performance_date` (`performance_date`),
    KEY `idx_accuracy` (`accuracy`),
    KEY `idx_model_degradation` (`model_degradation_detected`),
    UNIQUE KEY `uk_tenant_model_date` (`tenant_id`, `model_id`, `performance_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================
-- INDEXES FOR COMMON QUERIES
-- ============================================================
CREATE INDEX idx_ai_models_tenant_type_status ON ai_models(tenant_id, model_type, status);
CREATE INDEX idx_recommendation_engine_tenant_type ON recommendation_engine(tenant_id, recommendation_type);
CREATE INDEX idx_user_recommendations_tenant_engine ON user_recommendations(tenant_id, engine_id, created_at DESC);
CREATE INDEX idx_prediction_results_tenant_model_type ON prediction_results(tenant_id, model_id, prediction_type);
CREATE INDEX idx_anomaly_detection_tenant_severity ON anomaly_detection(tenant_id, severity_level, is_resolved);
CREATE INDEX idx_ai_insights_tenant_priority ON ai_insights(tenant_id, priority_level, is_archived);
CREATE INDEX idx_model_performance_tenant_model_date ON model_performance(tenant_id, model_id, performance_date DESC);

SET FOREIGN_KEY_CHECKS = 1;
