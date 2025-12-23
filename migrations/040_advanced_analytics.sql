-- Phase 3.1: Advanced Analytics
-- Comprehensive analytics, dashboard, and reporting infrastructure

-- ============================================
-- Analytics Dashboard Tables
-- ============================================

-- Real-time dashboard metrics and KPIs
CREATE TABLE IF NOT EXISTS analytics_dashboards (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  dashboard_name VARCHAR(255) NOT NULL,
  dashboard_description TEXT,
  dashboard_type ENUM('executive', 'manager', 'agent', 'custom') DEFAULT 'custom',
  layout_config JSON,
  is_public BOOLEAN DEFAULT FALSE,
  is_default BOOLEAN DEFAULT FALSE,
  widget_ids JSON,
  refresh_interval INT DEFAULT 300,
  data_date_range VARCHAR(50),
  last_refreshed_at TIMESTAMP,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_user (tenant_id, user_id),
  INDEX idx_dashboard_type (tenant_id, dashboard_type),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Analytics widgets (building blocks for dashboards)
CREATE TABLE IF NOT EXISTS analytics_widgets (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  widget_name VARCHAR(255) NOT NULL,
  widget_type ENUM('metric', 'chart', 'table', 'gauge', 'timeline', 'heatmap', 'funnel') NOT NULL,
  metric_source VARCHAR(100),
  metric_key VARCHAR(100),
  aggregation_type ENUM('sum', 'avg', 'count', 'min', 'max', 'distinct') DEFAULT 'count',
  time_period ENUM('hour', 'day', 'week', 'month', 'quarter', 'year', 'custom') DEFAULT 'day',
  comparison_enabled BOOLEAN DEFAULT FALSE,
  comparison_period ENUM('previous', 'year_ago', 'custom'),
  drill_down_enabled BOOLEAN DEFAULT FALSE,
  drill_down_target VARCHAR(255),
  widget_config JSON,
  refresh_interval INT DEFAULT 300,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant (tenant_id),
  INDEX idx_widget_type (widget_type),
  INDEX idx_metric_source (metric_source)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Real-time KPI tracking
CREATE TABLE IF NOT EXISTS analytics_kpis (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  kpi_name VARCHAR(255) NOT NULL,
  kpi_description TEXT,
  kpi_category ENUM('sales', 'customer', 'finance', 'operations', 'performance', 'quality') NOT NULL,
  metric_source VARCHAR(100),
  calculation_method VARCHAR(255),
  target_value DECIMAL(15, 2),
  current_value DECIMAL(15, 2),
  previous_value DECIMAL(15, 2),
  threshold_warning DECIMAL(15, 2),
  threshold_critical DECIMAL(15, 2),
  status ENUM('on_track', 'at_risk', 'critical', 'exceeded') DEFAULT 'on_track',
  percentage_change DECIMAL(5, 2),
  trend ENUM('up', 'down', 'stable') DEFAULT 'stable',
  time_period ENUM('hour', 'day', 'week', 'month', 'quarter', 'year') DEFAULT 'month',
  measurement_date DATE,
  last_updated_at TIMESTAMP,
  owner_id BIGINT,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_category (tenant_id, kpi_category),
  INDEX idx_kpi_status (tenant_id, status),
  INDEX idx_measurement_date (measurement_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Report Management Tables
-- ============================================

-- Report definitions and templates
CREATE TABLE IF NOT EXISTS analytics_reports (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  report_name VARCHAR(255) NOT NULL,
  report_description TEXT,
  report_type ENUM('operational', 'financial', 'sales', 'customer', 'performance', 'compliance', 'custom') NOT NULL,
  report_category VARCHAR(100),
  report_template VARCHAR(100),
  data_source VARCHAR(100),
  filters JSON,
  grouping JSON,
  sorting JSON,
  metrics JSON,
  dimensions JSON,
  format ENUM('pdf', 'excel', 'csv', 'json', 'html') DEFAULT 'pdf',
  time_period ENUM('daily', 'weekly', 'monthly', 'quarterly', 'yearly', 'custom') DEFAULT 'monthly',
  schedule_frequency ENUM('once', 'daily', 'weekly', 'monthly', 'quarterly', 'yearly') DEFAULT 'once',
  schedule_day INT,
  schedule_time TIME,
  is_public BOOLEAN DEFAULT FALSE,
  is_scheduled BOOLEAN DEFAULT FALSE,
  last_generated_at TIMESTAMP NULL,
  next_scheduled_run TIMESTAMP NULL,
  recipient_emails JSON,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_type (tenant_id, report_type),
  INDEX idx_next_scheduled (next_scheduled_run),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Generated report instances
CREATE TABLE IF NOT EXISTS analytics_report_instances (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  report_id BIGINT NOT NULL,
  report_name VARCHAR(255) NOT NULL,
  report_type VARCHAR(100),
  generation_start_time TIMESTAMP,
  generation_end_time TIMESTAMP,
  generation_duration_seconds INT,
  status ENUM('pending', 'generating', 'completed', 'failed', 'archived') DEFAULT 'pending',
  file_path VARCHAR(500),
  file_size BIGINT,
  file_format VARCHAR(20),
  s3_bucket VARCHAR(255),
  s3_key VARCHAR(500),
  record_count BIGINT,
  data_from_date DATE,
  data_to_date DATE,
  generated_by BIGINT,
  download_count INT DEFAULT 0,
  last_downloaded_at TIMESTAMP NULL,
  error_message TEXT,
  metadata JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_report (tenant_id, report_id),
  INDEX idx_status (status),
  INDEX idx_created_at (created_at),
  FOREIGN KEY (report_id) REFERENCES analytics_reports(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Time Series and Trend Analysis Tables
-- ============================================

-- Time series data points for trend analysis
CREATE TABLE IF NOT EXISTS analytics_timeseries (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  metric_name VARCHAR(255) NOT NULL,
  metric_category VARCHAR(100),
  metric_value DECIMAL(15, 4),
  metric_unit VARCHAR(50),
  time_bucket TIMESTAMP NOT NULL,
  time_granularity ENUM('minute', 'hour', 'day', 'week', 'month') DEFAULT 'day',
  dimension_values JSON,
  cumulative_value DECIMAL(15, 4),
  trend_indicator ENUM('increasing', 'decreasing', 'stable', 'volatile') DEFAULT 'stable',
  is_anomaly BOOLEAN DEFAULT FALSE,
  anomaly_score DECIMAL(5, 2),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_tenant_metric (tenant_id, metric_name),
  INDEX idx_time_bucket (time_bucket),
  INDEX idx_metric_category (metric_category),
  INDEX idx_anomaly (is_anomaly)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Trend analysis and forecasting
CREATE TABLE IF NOT EXISTS analytics_trends (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  metric_name VARCHAR(255) NOT NULL,
  analysis_type ENUM('linear', 'exponential', 'moving_average', 'seasonal') NOT NULL,
  time_period_days INT,
  baseline_value DECIMAL(15, 4),
  current_value DECIMAL(15, 4),
  growth_rate DECIMAL(5, 2),
  volatility DECIMAL(5, 2),
  forecast_value DECIMAL(15, 4),
  forecast_confidence DECIMAL(5, 2),
  forecast_horizon_days INT,
  correlation_with_other_metrics JSON,
  seasonal_pattern VARCHAR(100),
  analysis_date DATE,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_metric (tenant_id, metric_name),
  INDEX idx_analysis_date (analysis_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Data Aggregation and Caching Tables
-- ============================================

-- Pre-aggregated metrics for fast querying
CREATE TABLE IF NOT EXISTS analytics_aggregates (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  aggregate_name VARCHAR(255) NOT NULL,
  aggregate_type ENUM('hourly', 'daily', 'weekly', 'monthly') NOT NULL,
  aggregate_date DATE,
  aggregate_time TIME,
  dimension_keys JSON,
  dimension_values JSON,
  metric_data JSON,
  record_count BIGINT,
  last_computed_at TIMESTAMP,
  computation_duration_ms INT,
  data_freshness_minutes INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_tenant_date (tenant_id, aggregate_date),
  INDEX idx_aggregate_type (aggregate_type),
  UNIQUE KEY unique_aggregate (tenant_id, aggregate_name, aggregate_type, aggregate_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Drill-down and Segmentation Tables
-- ============================================

-- Detailed breakdown for drill-down analytics
CREATE TABLE IF NOT EXISTS analytics_drilldown (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  parent_metric_id BIGINT,
  parent_metric_name VARCHAR(255),
  drill_level INT,
  drill_dimension VARCHAR(100),
  segment_name VARCHAR(255),
  segment_value VARCHAR(255),
  metric_value DECIMAL(15, 4),
  metric_percentage DECIMAL(5, 2),
  record_count BIGINT,
  contribution_rank INT,
  analysis_date DATE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_tenant_parent (tenant_id, parent_metric_id),
  INDEX idx_drill_dimension (drill_dimension),
  INDEX idx_analysis_date (analysis_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- User Preferences and Configurations
-- ============================================

-- User dashboard and report preferences
CREATE TABLE IF NOT EXISTS analytics_user_preferences (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  default_dashboard_id BIGINT,
  default_report_id BIGINT,
  preferred_time_zone VARCHAR(50),
  preferred_currency VARCHAR(10),
  date_format VARCHAR(20),
  number_format VARCHAR(20),
  auto_refresh_enabled BOOLEAN DEFAULT TRUE,
  auto_refresh_interval INT DEFAULT 300,
  email_notifications_enabled BOOLEAN DEFAULT TRUE,
  report_frequency VARCHAR(50),
  theme_preference ENUM('light', 'dark', 'auto') DEFAULT 'auto',
  chart_preferences JSON,
  saved_filters JSON,
  export_format_default VARCHAR(20),
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_user (tenant_id, user_id),
  UNIQUE KEY unique_user_prefs (tenant_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Analytics Activity Audit
-- ============================================

-- Analytics access and usage audit trail
CREATE TABLE IF NOT EXISTS analytics_audit_log (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  user_id BIGINT,
  action_type ENUM('view_dashboard', 'create_dashboard', 'edit_dashboard', 'delete_dashboard', 'generate_report', 'download_report', 'share_dashboard', 'export_data', 'view_widget', 'drill_down') NOT NULL,
  resource_type VARCHAR(100),
  resource_id BIGINT,
  resource_name VARCHAR(255),
  action_details JSON,
  ip_address VARCHAR(45),
  user_agent TEXT,
  device_type VARCHAR(50),
  duration_seconds INT,
  status ENUM('success', 'failure', 'error') DEFAULT 'success',
  error_message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_tenant_date (tenant_id, created_at),
  INDEX idx_user_date (user_id, created_at),
  INDEX idx_action_type (action_type),
  INDEX idx_resource (resource_type, resource_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Custom Metrics and KPI Configuration
-- ============================================

-- Custom metric definitions
CREATE TABLE IF NOT EXISTS analytics_custom_metrics (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  metric_name VARCHAR(255) NOT NULL,
  metric_description TEXT,
  metric_formula VARCHAR(500),
  formula_type ENUM('simple', 'complex', 'custom_sql') NOT NULL,
  data_source VARCHAR(100),
  source_table VARCHAR(100),
  source_columns JSON,
  filter_conditions JSON,
  grouping_dimensions JSON,
  output_unit VARCHAR(50),
  output_decimal_places INT DEFAULT 2,
  refresh_frequency ENUM('real_time', 'hourly', 'daily', 'weekly', 'manual') DEFAULT 'daily',
  last_refresh_at TIMESTAMP NULL,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant (tenant_id),
  INDEX idx_metric_name (metric_name),
  INDEX idx_data_source (data_source)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- Alerting and Threshold Management
-- ============================================

-- Analytics alerts and thresholds
CREATE TABLE IF NOT EXISTS analytics_alerts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  alert_name VARCHAR(255) NOT NULL,
  alert_type ENUM('threshold', 'anomaly', 'trend', 'comparison') NOT NULL,
  metric_name VARCHAR(255),
  condition ENUM('equals', 'not_equals', 'greater_than', 'less_than', 'between', 'contains') NOT NULL,
  threshold_value DECIMAL(15, 4),
  threshold_lower DECIMAL(15, 4),
  threshold_upper DECIMAL(15, 4),
  severity ENUM('low', 'medium', 'high', 'critical') DEFAULT 'medium',
  enabled BOOLEAN DEFAULT TRUE,
  alert_frequency ENUM('once', 'daily', 'hourly') DEFAULT 'once',
  notification_channels JSON,
  recipient_emails JSON,
  recipient_user_ids JSON,
  last_triggered_at TIMESTAMP NULL,
  trigger_count INT DEFAULT 0,
  created_by BIGINT,
  updated_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tenant_enabled (tenant_id, enabled),
  INDEX idx_severity (severity),
  INDEX idx_last_triggered (last_triggered_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Alert history and trigger logs
CREATE TABLE IF NOT EXISTS analytics_alert_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL,
  alert_id BIGINT NOT NULL,
  alert_name VARCHAR(255),
  metric_name VARCHAR(255),
  metric_value DECIMAL(15, 4),
  threshold_value DECIMAL(15, 4),
  condition_met BOOLEAN,
  triggered_at TIMESTAMP,
  acknowledged_at TIMESTAMP NULL,
  acknowledged_by BIGINT,
  acknowledged_notes TEXT,
  notification_sent BOOLEAN DEFAULT FALSE,
  notification_channels JSON,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  INDEX idx_tenant_date (tenant_id, triggered_at),
  INDEX idx_alert (alert_id),
  INDEX idx_acknowledged (acknowledged_at),
  FOREIGN KEY (alert_id) REFERENCES analytics_alerts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
