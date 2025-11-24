-- Gamification System - Migration
-- Adds complete points, badges, achievements, and leaderboards for multi-tenant call center

-- Gamification Configuration per Tenant
CREATE TABLE IF NOT EXISTS gamification_config (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL UNIQUE,
    enabled BOOLEAN DEFAULT TRUE,
    points_per_call INT DEFAULT 10,
    points_per_conversion INT DEFAULT 50,
    points_per_quality_review INT DEFAULT 5,
    points_per_feedback INT DEFAULT 3,
    points_decay_percent INT DEFAULT 5,
    decay_period_days INT DEFAULT 30,
    max_daily_points INT DEFAULT 500,
    leaderboard_reset_period ENUM('daily', 'weekly', 'monthly') DEFAULT 'weekly',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_tenant (tenant_id)
);

-- User Points (Current Period)
CREATE TABLE IF NOT EXISTS user_points (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    current_points INT DEFAULT 0,
    lifetime_points INT DEFAULT 0,
    period_start_date DATE NOT NULL,
    period_end_date DATE NOT NULL,
    daily_points INT DEFAULT 0,
    daily_date DATE,
    streak_days INT DEFAULT 0,
    last_action_date DATE,
    `rank` INT,
    rank_updated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_tenant_period (user_id, tenant_id, period_start_date),
    INDEX idx_tenant_points (tenant_id, current_points),
    INDEX idx_user_tenant (user_id, tenant_id),
    INDEX idx_daily_date (daily_date)
);

-- Point Transactions Log
CREATE TABLE IF NOT EXISTS point_transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    points INT NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    action_id VARCHAR(255),
    description TEXT,
    multiplier DECIMAL(3,2) DEFAULT 1.0,
    bonus_reason VARCHAR(100),
    status ENUM('pending', 'awarded', 'revoked') DEFAULT 'awarded',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_user_tenant_action (user_id, tenant_id, action_type),
    INDEX idx_created_at (created_at),
    INDEX idx_status (status)
);

-- Badges and Achievements
CREATE TABLE IF NOT EXISTS badge (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36),
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    category ENUM('milestone', 'skill', 'behavior', 'performance', 'social') NOT NULL,
    requirement_type VARCHAR(50) NOT NULL,
    requirement_value INT,
    points_reward INT DEFAULT 0,
    rarity ENUM('common', 'uncommon', 'rare', 'epic', 'legendary') DEFAULT 'uncommon',
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    UNIQUE KEY unique_tenant_code (tenant_id, code),
    INDEX idx_category (category),
    INDEX idx_rarity (rarity)
);

-- User Badges Earned
CREATE TABLE IF NOT EXISTS user_badge (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    badge_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    earned_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    progress_percent INT DEFAULT 100,
    notified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_id) REFERENCES badge(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_badge (user_id, badge_id),
    INDEX idx_user_tenant (user_id, tenant_id),
    INDEX idx_earned_date (earned_date)
);

-- Challenges and Quests
CREATE TABLE IF NOT EXISTS challenge (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    challenge_type ENUM('daily', 'weekly', 'monthly', 'seasonal', 'special') NOT NULL,
    status ENUM('active', 'inactive', 'completed') DEFAULT 'active',
    objective_type VARCHAR(50) NOT NULL,
    objective_target INT,
    objective_current INT DEFAULT 0,
    points_reward INT,
    badge_reward_id BIGINT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    difficulty ENUM('easy', 'medium', 'hard', 'extreme') DEFAULT 'medium',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_reward_id) REFERENCES badge(id) ON DELETE SET NULL,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_type (challenge_type),
    INDEX idx_end_date (end_date)
);

-- User Challenge Progress
CREATE TABLE IF NOT EXISTS user_challenge (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    challenge_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    progress INT DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    completed_date TIMESTAMP,
    points_earned INT DEFAULT 0,
    badge_earned_id BIGINT,
    status ENUM('in_progress', 'completed', 'failed', 'abandoned') DEFAULT 'in_progress',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (challenge_id) REFERENCES challenge(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (badge_earned_id) REFERENCES badge(id) ON DELETE SET NULL,
    UNIQUE KEY unique_user_challenge (user_id, challenge_id),
    INDEX idx_user_tenant (user_id, tenant_id),
    INDEX idx_status (status)
);

-- Leaderboards (Pre-calculated for performance)
CREATE TABLE IF NOT EXISTS leaderboard (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    period_type ENUM('daily', 'weekly', 'monthly', 'all_time') NOT NULL,
    period_date DATE NOT NULL,
    user_id BIGINT NOT NULL,
    `rank` INT NOT NULL,
    points INT NOT NULL,
    points_change INT DEFAULT 0,
    previous_rank INT,
    badges_count INT DEFAULT 0,
    streak_days INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    UNIQUE KEY unique_leaderboard (tenant_id, period_type, period_date, user_id),
    INDEX idx_tenant_period_rank (tenant_id, period_type, `rank`),
    INDEX idx_user_tenant (user_id, tenant_id)
);

-- Rewards and Redemption
CREATE TABLE IF NOT EXISTS reward_item (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    points_cost INT NOT NULL,
    stock INT DEFAULT -1,
    stock_remaining INT DEFAULT -1,
    image_url VARCHAR(500),
    redemption_type ENUM('digital', 'physical', 'experience', 'discount') NOT NULL,
    redemption_details JSON,
    active BOOLEAN DEFAULT TRUE,
    featured BOOLEAN DEFAULT FALSE,
    expiry_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_tenant_active (tenant_id, active),
    INDEX idx_category (category)
);

-- User Redemptions
CREATE TABLE IF NOT EXISTS user_redemption (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    reward_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    points_spent INT NOT NULL,
    quantity INT DEFAULT 1,
    status ENUM('pending', 'approved', 'completed', 'cancelled') DEFAULT 'pending',
    redemption_code VARCHAR(100),
    redeemed_date TIMESTAMP,
    completed_date TIMESTAMP,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (reward_id) REFERENCES reward_item(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_user_tenant_status (user_id, tenant_id, status),
    INDEX idx_created_at (created_at)
);

-- Team/Department Challenges
CREATE TABLE IF NOT EXISTS team_challenge (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    team_id BIGINT,
    objective_type VARCHAR(50) NOT NULL,
    objective_target INT,
    team_reward_points INT,
    individual_bonus_multiplier DECIMAL(2,1) DEFAULT 1.5,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    status ENUM('active', 'completed', 'cancelled') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_end_date (end_date)
);

-- Milestones/Tiers
CREATE TABLE IF NOT EXISTS player_level (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    level_number INT NOT NULL,
    name VARCHAR(100),
    description TEXT,
    points_required INT NOT NULL,
    points_total_required INT NOT NULL,
    icon_url VARCHAR(500),
    benefits TEXT,
    unlock_privileges JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    UNIQUE KEY unique_tenant_level (tenant_id, level_number),
    INDEX idx_tenant (tenant_id)
);

-- User Levels/Tiers
CREATE TABLE IF NOT EXISTS user_level (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    current_level INT DEFAULT 1,
    lifetime_level INT DEFAULT 1,
    previous_level INT,
    level_up_date TIMESTAMP,
    progress_to_next INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_tenant_level (user_id, tenant_id),
    INDEX idx_user_tenant (user_id, tenant_id)
);

-- Achievements History for Analytics
CREATE TABLE IF NOT EXISTS achievement_event (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    event_name VARCHAR(255),
    details JSON,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES `user`(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE,
    INDEX idx_user_tenant_event (user_id, tenant_id, event_type),
    INDEX idx_timestamp (timestamp)
);

-- Indexes for performance
DROP INDEX idx_point_transactions_user_date ON point_transactions;
CREATE INDEX idx_point_transactions_user_date ON point_transactions(user_id, created_at);

DROP INDEX idx_user_points_tenant_date ON user_points;
CREATE INDEX idx_user_points_tenant_date ON user_points(tenant_id, daily_date);

DROP INDEX idx_badge_tenant_active ON badge;
CREATE INDEX idx_badge_tenant_active ON badge(tenant_id, active);

DROP INDEX idx_challenge_tenant_active ON challenge;
CREATE INDEX idx_challenge_tenant_active ON challenge(tenant_id, status, end_date);
