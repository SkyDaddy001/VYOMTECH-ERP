-- Gamification Module Schema Migration
-- Extends existing gamification tables with additional fields for enhanced tracking

-- Gamification Points History Table
CREATE TABLE IF NOT EXISTS gamification_points_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    points INT NOT NULL,
    reason VARCHAR(500),
    action_type VARCHAR(100), -- call_completed, lead_closed, challenge_completed, badge_earned, daily_bonus
    reference_id VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_action_type (action_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Gamification Badges Table
CREATE TABLE IF NOT EXISTS gamification_badges (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    badge_name VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    icon_url VARCHAR(500),
    criteria VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    UNIQUE KEY unique_badge (tenant_id, badge_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User Badges Table
CREATE TABLE IF NOT EXISTS user_badges (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    badge_id BIGINT NOT NULL,
    earned_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_badge_id (badge_id),
    UNIQUE KEY unique_user_badge (tenant_id, user_id, badge_id),
    FOREIGN KEY (badge_id) REFERENCES gamification_badges(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Gamification Challenges Table
CREATE TABLE IF NOT EXISTS gamification_challenges (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    challenge_name VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    target_value INT,
    reward_points INT,
    start_date DATETIME,
    end_date DATETIME,
    status VARCHAR(50), -- active, completed, expired, cancelled
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_start_date (start_date),
    INDEX idx_end_date (end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- User Challenge Progress Table
CREATE TABLE IF NOT EXISTS user_challenge_progress (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    challenge_id BIGINT NOT NULL,
    current_progress INT DEFAULT 0,
    target_value INT,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_challenge_id (challenge_id),
    UNIQUE KEY unique_user_challenge (tenant_id, user_id, challenge_id),
    FOREIGN KEY (challenge_id) REFERENCES gamification_challenges(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Leaderboards Table
CREATE TABLE IF NOT EXISTS leaderboards (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    total_points INT DEFAULT 0,
    user_rank INT,
    badges_count INT DEFAULT 0,
    challenges_completed INT DEFAULT 0,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_rank (user_rank),
    INDEX idx_total_points (total_points),
    UNIQUE KEY unique_leaderboard (tenant_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
