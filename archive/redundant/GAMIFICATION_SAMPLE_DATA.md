# 🎮 Gamification System - Sample Data & Examples

Complete sample data, API responses, and real-world examples.

---

## 📊 Sample Database Data

### Initialize Default Configuration
```sql
-- Gamification config
INSERT INTO gamification_config (tenant_id, enabled, points_per_call, points_per_conversion, 
                                 points_per_quality_review, points_per_feedback, 
                                 max_daily_points, leaderboard_reset_period)
VALUES ('tenant-001', true, 10, 50, 5, 3, 500, 'weekly');

-- Create badges
INSERT INTO badge (tenant_id, code, name, description, icon_url, category, 
                   requirement_type, requirement_value, points_reward, rarity, active)
VALUES
('tenant-001', 'welcome', 'Welcome Aboard', 'Complete first call', '🎉', 'milestone', 'total_calls', 1, 10, 'common', true),
('tenant-001', 'speed_demon', 'Speed Demon', 'Complete 5 calls in under 10 minutes', '⚡', 'skill', 'fast_calls', 5, 50, 'rare', true),
('tenant-001', 'perfect_week', 'Perfect Week', 'Maintain 100% quality for a week', '⭐', 'performance', 'perfect_week', 1, 200, 'epic', true),
('tenant-001', 'team_player', 'Team Player', 'Help 10 teammates', '🤝', 'social', 'helped_teammates', 10, 100, 'uncommon', true),
('tenant-001', 'century_club', 'Century Club', 'Complete 100 calls', '💯', 'milestone', 'total_calls', 100, 250, 'epic', true),
('tenant-001', 'streak_7', 'Week Warrior', '7-day streak', '🔥', 'behavior', 'streak_days', 7, 75, 'uncommon', true),
('tenant-001', 'streak_30', 'Month Master', '30-day streak', '🔥🔥', 'behavior', 'streak_days', 30, 500, 'legendary', true),
('tenant-001', 'conversions_10', 'Closer', '10 lead conversions', '💰', 'performance', 'conversions', 10, 150, 'rare', true),
('tenant-001', 'feedback_hero', 'Feedback Hero', '50 positive feedbacks', '💬', 'social', 'positive_feedback', 50, 100, 'rare', true),
('tenant-001', 'level_10', 'Legend', 'Reach Level 10', '👑', 'milestone', 'level', 10, 1000, 'legendary', true);

-- Create challenges
INSERT INTO challenge (tenant_id, name, description, challenge_type, status, objective_type, 
                      objective_target, points_reward, difficulty, start_date, end_date)
VALUES
-- Daily Challenges
('tenant-001', 'Daily Grind', 'Complete 10 calls', 'daily', 'active', 'total_calls', 10, 25, 'easy', NOW(), DATE_ADD(NOW(), INTERVAL 1 DAY)),
('tenant-001', 'Conversion Blitz', '3 lead conversions today', 'daily', 'active', 'conversions', 3, 50, 'medium', NOW(), DATE_ADD(NOW(), INTERVAL 1 DAY)),

-- Weekly Challenges
('tenant-001', 'Weekly Champion', 'Complete 50 calls this week', 'weekly', 'active', 'total_calls', 50, 100, 'medium', DATE_SUB(NOW(), INTERVAL DAYOFWEEK(NOW())-1 DAY), DATE_ADD(DATE_SUB(NOW(), INTERVAL DAYOFWEEK(NOW())-1 DAY), INTERVAL 7 DAY)),
('tenant-001', 'Conversion King', '15 conversions this week', 'weekly', 'active', 'conversions', 15, 150, 'hard', DATE_SUB(NOW(), INTERVAL DAYOFWEEK(NOW())-1 DAY), DATE_ADD(DATE_SUB(NOW(), INTERVAL DAYOFWEEK(NOW())-1 DAY), INTERVAL 7 DAY)),

-- Monthly Challenges
('tenant-001', 'Monthly Marathon', 'Complete 200 calls this month', 'monthly', 'active', 'total_calls', 200, 300, 'hard', DATE(CONCAT(YEAR(NOW()), '-', MONTH(NOW()), '-01')), LAST_DAY(NOW())),
('tenant-001', 'Perfect Month', '95%+ quality rating for the month', 'monthly', 'active', 'quality_rating', 95, 500, 'hard', DATE(CONCAT(YEAR(NOW()), '-', MONTH(NOW()), '-01')), LAST_DAY(NOW()));

-- Create player levels
INSERT INTO player_level (tenant_id, level_number, name, points_required, points_total_required, benefits)
VALUES
('tenant-001', 1, 'Rookie', 0, 0, 'Welcome to the game'),
('tenant-001', 2, 'Bronze', 1000, 1000, '+50 daily point limit'),
('tenant-001', 3, 'Bronze', 1000, 2000, '+50 daily point limit'),
('tenant-001', 4, 'Silver', 1000, 3000, 'Can create team challenges'),
('tenant-001', 5, 'Silver', 1000, 4000, 'Invite teammates to challenges'),
('tenant-001', 6, 'Gold', 1000, 5000, 'Premium rewards unlock'),
('tenant-001', 7, 'Gold', 1000, 6000, 'Leaderboard freeze protection'),
('tenant-001', 8, 'Platinum', 1000, 7000, 'Exclusive badges'),
('tenant-001', 9, 'Diamond', 1000, 8000, 'VIP customer support'),
('tenant-001', 10, 'Legend', 1000, 9000, 'Lifetime perks');

-- Create rewards
INSERT INTO reward_item (tenant_id, name, description, category, points_cost, stock, image_url, 
                         redemption_type, active, featured, expiry_date)
VALUES
('tenant-001', '10% Off Merch', '10% discount on company merchandise', 'discount', 500, -1, '/rewards/merch-10.png', 'discount', true, true, DATE_ADD(NOW(), INTERVAL 90 DAY)),
('tenant-001', '25% Off Coffee', '25% off premium coffee subscription', 'discount', 250, -1, '/rewards/coffee-25.png', 'discount', true, true, DATE_ADD(NOW(), INTERVAL 90 DAY)),
('tenant-001', 'Free Lunch', 'Free lunch at company cafeteria', 'experience', 300, 50, '/rewards/lunch.png', 'experience', true, true, DATE_ADD(NOW(), INTERVAL 30 DAY)),
('tenant-001', 'Premium Headphones', 'Premium noise-canceling headphones', 'physical', 1000, 20, '/rewards/headphones.png', 'physical', true, false, DATE_ADD(NOW(), INTERVAL 180 DAY)),
('tenant-001', 'Extra PTO Day', 'One additional paid time off day', 'experience', 800, -1, '/rewards/pto.png', 'experience', true, false, DATE_ADD(NOW(), INTERVAL 365 DAY)),
('tenant-001', '$50 Gift Card', 'Amazon gift card $50 value', 'digital', 400, 100, '/rewards/amazon-50.png', 'digital', true, true, DATE_ADD(NOW(), INTERVAL 60 DAY)),
('tenant-001', 'Massage Voucher', 'Corporate massage session ($100 value)', 'experience', 750, 30, '/rewards/massage.png', 'experience', true, false, DATE_ADD(NOW(), INTERVAL 120 DAY)),
('tenant-001', 'Premium Parking Spot', 'Reserved parking for 1 month', 'experience', 200, -1, '/rewards/parking.png', 'experience', true, true, DATE_ADD(NOW(), INTERVAL 45 DAY));

-- Create user data (example agents)
INSERT INTO user (email, password_hash, role, status, tenant_id, first_name, last_name)
VALUES
('sarah.chen@company.com', 'hashed_pwd_1', 'agent', 'active', 'tenant-001', 'Sarah', 'Chen'),
('mike.johnson@company.com', 'hashed_pwd_2', 'agent', 'active', 'tenant-001', 'Mike', 'Johnson'),
('emma.davis@company.com', 'hashed_pwd_3', 'agent', 'active', 'tenant-001', 'Emma', 'Davis'),
('tom.wilson@company.com', 'hashed_pwd_4', 'agent', 'active', 'tenant-001', 'Tom', 'Wilson'),
('alice.zhang@company.com', 'hashed_pwd_5', 'agent', 'active', 'tenant-001', 'Alice', 'Zhang');

-- Create user points (sample data)
INSERT INTO user_points (user_id, tenant_id, current_points, lifetime_points, period_start_date, 
                        period_end_date, daily_points, daily_date, streak_days, rank)
VALUES
(1, 'tenant-001', 450, 2450, CURDATE(), DATE_ADD(CURDATE(), INTERVAL 7 DAY), 145, CURDATE(), 7, 1),
(2, 'tenant-001', 380, 2380, CURDATE(), DATE_ADD(CURDATE(), INTERVAL 7 DAY), 120, CURDATE(), 3, 2),
(3, 'tenant-001', 310, 2310, CURDATE(), DATE_ADD(CURDATE(), INTERVAL 7 DAY), 95, CURDATE(), 5, 3),
(4, 'tenant-001', 245, 2245, CURDATE(), DATE_ADD(CURDATE(), INTERVAL 7 DAY), 70, CURDATE(), 2, 4),
(5, 'tenant-001', 180, 1980, CURDATE(), DATE_ADD(CURDATE(), INTERVAL 7 DAY), 55, CURDATE(), 1, 5);

-- Award some badges
INSERT INTO user_badge (user_id, badge_id, tenant_id, earned_date)
VALUES
(1, 1, 'tenant-001', DATE_SUB(NOW(), INTERVAL 30 DAY)),  -- Sarah: Welcome
(1, 6, 'tenant-001', DATE_SUB(NOW(), INTERVAL 7 DAY)),   -- Sarah: Week Warrior
(2, 1, 'tenant-001', DATE_SUB(NOW(), INTERVAL 60 DAY)),  -- Mike: Welcome
(3, 1, 'tenant-001', DATE_SUB(NOW(), INTERVAL 45 DAY)),  -- Emma: Welcome
(3, 8, 'tenant-001', DATE_SUB(NOW(), INTERVAL 10 DAY));  -- Emma: Closer

-- Create user levels
INSERT INTO user_level (user_id, tenant_id, current_level, lifetime_level, progress_to_next)
VALUES
(1, 'tenant-001', 3, 3, 450),   -- Sarah: Level 3, 450/1000 to Level 4
(2, 'tenant-001', 2, 2, 380),   -- Mike: Level 2, 380/1000 to Level 3
(3, 'tenant-001', 2, 2, 310),   -- Emma: Level 2, 310/1000 to Level 3
(4, 'tenant-001', 2, 2, 245),   -- Tom: Level 2, 245/1000 to Level 3
(5, 'tenant-001', 1, 1, 180);   -- Alice: Level 1, 180/1000 to Level 2

-- Create user challenges
INSERT INTO user_challenge (user_id, challenge_id, tenant_id, progress, status)
VALUES
(1, 1, 'tenant-001', 8, 'in_progress'),   -- Sarah: 8/10 daily calls
(1, 3, 'tenant-001', 42, 'in_progress'),  -- Sarah: 42/50 weekly calls
(2, 1, 'tenant-001', 6, 'in_progress'),   -- Mike: 6/10 daily calls
(3, 2, 'tenant-001', 2, 'in_progress'),   -- Emma: 2/3 conversions
(4, 1, 'tenant-001', 4, 'in_progress');   -- Tom: 4/10 daily calls

-- Create leaderboard entries
INSERT INTO leaderboard (tenant_id, period_type, period_date, user_id, rank, points, badges_count, streak_days)
VALUES
('tenant-001', 'weekly', CURDATE(), 1, 1, 450, 2, 7),
('tenant-001', 'weekly', CURDATE(), 2, 2, 380, 1, 3),
('tenant-001', 'weekly', CURDATE(), 3, 3, 310, 1, 5),
('tenant-001', 'weekly', CURDATE(), 4, 4, 245, 0, 2),
('tenant-001', 'weekly', CURDATE(), 5, 5, 180, 0, 1);
```

---

## 🔌 API Response Examples

### Get Gamification Profile
```json
GET /api/v1/gamification/profile
Authorization: Bearer eyJhbGc...

{
  "user_points": {
    "id": 1,
    "user_id": 1,
    "tenant_id": "tenant-001",
    "current_points": 450,
    "lifetime_points": 2450,
    "period_start_date": "2025-11-17",
    "period_end_date": "2025-11-24",
    "daily_points": 145,
    "daily_date": "2025-11-22",
    "streak_days": 7,
    "last_action_date": "2025-11-22",
    "rank": 1,
    "rank_updated_at": "2025-11-22T14:32:10Z"
  },
  "current_level": {
    "id": 1,
    "user_id": 1,
    "tenant_id": "tenant-001",
    "current_level": 3,
    "lifetime_level": 3,
    "level_up_date": "2025-11-10T00:00:00Z",
    "progress_to_next": 450
  },
  "badges": [
    {
      "id": 1,
      "user_id": 1,
      "badge_id": 1,
      "tenant_id": "tenant-001",
      "earned_date": "2025-10-23T00:00:00Z",
      "badge": {
        "id": 1,
        "code": "welcome",
        "name": "Welcome Aboard",
        "description": "Complete first call",
        "icon_url": "🎉",
        "category": "milestone",
        "rarity": "common",
        "points_reward": 10
      }
    },
    {
      "id": 2,
      "user_id": 1,
      "badge_id": 6,
      "tenant_id": "tenant-001",
      "earned_date": "2025-11-15T00:00:00Z",
      "badge": {
        "id": 6,
        "code": "streak_7",
        "name": "Week Warrior",
        "description": "7-day streak",
        "icon_url": "🔥",
        "category": "behavior",
        "rarity": "uncommon",
        "points_reward": 75
      }
    }
  ],
  "active_challenges": [
    {
      "id": 1,
      "user_id": 1,
      "challenge_id": 1,
      "tenant_id": "tenant-001",
      "progress": 8,
      "completed": false,
      "status": "in_progress",
      "challenge": {
        "id": 1,
        "name": "Daily Grind",
        "description": "Complete 10 calls",
        "challenge_type": "daily",
        "objective_type": "total_calls",
        "objective_target": 10,
        "points_reward": 25,
        "difficulty": "easy"
      }
    },
    {
      "id": 2,
      "user_id": 1,
      "challenge_id": 3,
      "tenant_id": "tenant-001",
      "progress": 42,
      "completed": false,
      "status": "in_progress",
      "challenge": {
        "id": 3,
        "name": "Weekly Champion",
        "description": "Complete 50 calls this week",
        "challenge_type": "weekly",
        "objective_type": "total_calls",
        "objective_target": 50,
        "points_reward": 100,
        "difficulty": "medium"
      }
    }
  ],
  "leaderboard_rank": {
    "id": 1,
    "tenant_id": "tenant-001",
    "period_type": "weekly",
    "period_date": "2025-11-22",
    "user_id": 1,
    "rank": 1,
    "points": 450,
    "points_change": 15,
    "previous_rank": 1,
    "badges_count": 2,
    "streak_days": 7
  },
  "total_rewards": 4,
  "streak_days": 7
}
```

### Get Leaderboard
```json
GET /api/v1/gamification/leaderboard?period=weekly&limit=10
Authorization: Bearer eyJhbGc...

{
  "entries": [
    {
      "id": 1,
      "tenant_id": "tenant-001",
      "period_type": "weekly",
      "period_date": "2025-11-22",
      "user_id": 1,
      "rank": 1,
      "points": 450,
      "points_change": 15,
      "previous_rank": 1,
      "badges_count": 2,
      "streak_days": 7,
      "user": {
        "id": 1,
        "email": "sarah.chen@company.com",
        "first_name": "Sarah",
        "last_name": "Chen"
      }
    },
    {
      "id": 2,
      "tenant_id": "tenant-001",
      "period_type": "weekly",
      "period_date": "2025-11-22",
      "user_id": 2,
      "rank": 2,
      "points": 380,
      "points_change": 5,
      "previous_rank": 2,
      "badges_count": 1,
      "streak_days": 3,
      "user": {
        "id": 2,
        "email": "mike.johnson@company.com",
        "first_name": "Mike",
        "last_name": "Johnson"
      }
    },
    {
      "id": 3,
      "tenant_id": "tenant-001",
      "period_type": "weekly",
      "period_date": "2025-11-22",
      "user_id": 3,
      "rank": 3,
      "points": 310,
      "points_change": -10,
      "previous_rank": 2,
      "badges_count": 1,
      "streak_days": 5,
      "user": {
        "id": 3,
        "email": "emma.davis@company.com",
        "first_name": "Emma",
        "last_name": "Davis"
      }
    }
  ],
  "user_rank": {
    "id": 1,
    "tenant_id": "tenant-001",
    "period_type": "weekly",
    "rank": 1,
    "points": 450
  },
  "period_type": "weekly"
}
```

### Get Rewards
```json
GET /api/v1/gamification/rewards
Authorization: Bearer eyJhbGc...

{
  "rewards": [
    {
      "id": 1,
      "tenant_id": "tenant-001",
      "name": "10% Off Merch",
      "description": "10% discount on company merchandise",
      "category": "discount",
      "points_cost": 500,
      "stock": -1,
      "image_url": "/rewards/merch-10.png",
      "redemption_type": "discount",
      "active": true,
      "featured": true,
      "expiry_date": "2026-02-20"
    },
    {
      "id": 5,
      "tenant_id": "tenant-001",
      "name": "Extra PTO Day",
      "description": "One additional paid time off day",
      "category": "experience",
      "points_cost": 800,
      "stock": -1,
      "image_url": "/rewards/pto.png",
      "redemption_type": "experience",
      "active": true,
      "featured": false,
      "expiry_date": "2026-11-22"
    }
  ]
}
```

### Redeem Reward
```json
POST /api/v1/gamification/redeem
Authorization: Bearer eyJhbGc...
Content-Type: application/json

{
  "reward_id": 1,
  "quantity": 1
}

Response (200):
{
  "success": true,
  "redemption": {
    "id": 101,
    "user_id": 1,
    "reward_id": 1,
    "status": "pending",
    "redemption_code": "REWARD-12345",
    "points_spent": 500,
    "created_at": "2025-11-22T15:30:00Z",
    "expiry_date": "2026-02-20"
  },
  "message": "Reward redeemed successfully. Code: REWARD-12345"
}
```

---

## 📱 Frontend Component Examples

### Display Points Widget
```tsx
<GamificationIndicator 
  points={450}
  level={3}
  rank={1}
  streak={7}
/>

// Renders:
// ⭐ 450 | Lv 3 | #1 | 🔥 7
```

### Display Challenge Progress
```tsx
<ChallengeCard
  name="Daily Grind"
  progress={8}
  target={10}
  reward={25}
/>

// Renders:
// Daily Grind
// ████████░░ 8/10
// +25 pts
```

### Display Badge
```tsx
<BadgeIcon
  name="Week Warrior"
  icon="🔥"
  rarity="uncommon"
  earnedDate="2025-11-15"
/>

// Renders badge with:
// - Fire emoji (🔥)
// - Green border (uncommon)
// - Earned date
```

### Leaderboard Row
```tsx
<LeaderboardRow
  rank={1}
  name="Sarah Chen"
  points={450}
  badges={2}
  streak={7}
  trend={15}  // +15 since yesterday
/>

// Renders:
// 🥇 Sarah Chen | 450 pts ⭐2 🔥7 | ↑15
```

---

## 🎯 Real-World Scenarios

### Scenario 1: Day in the Life of a Top Performer

**Morning (8:00 AM)**
- Sarah logs in
- Sees: 450 points, Rank #1, 7-day streak
- Dashboard shows: Daily Grind (8/10 done), Weekly Champion (42/50)

**During Day**
- Completes 2 more calls: +20 points (10 each)
- Gets positive feedback: +3 points
- Reaches 450 daily point limit at 3:00 PM

**End of Day (5:00 PM)**
- Points: 450 → Still 450 (daily limit)
- Lifetime: 2450 → 2473
- Checks rewards shop, considering "Extra PTO Day" (800 pts)

**Next Week**
- Completes "Weekly Champion" challenge: +100 points
- Completes "Perfect Week" badge: 200 bonus points
- Level up to Level 4!

---

### Scenario 2: Team Competition

**Weekly Challenge Active**: "Conversion King" - 15 conversions, +150 points

- Sarah: 8 conversions (on track)
- Mike: 5 conversions (behind)
- Emma: 12 conversions (leading)
- Tom: 3 conversions (struggling)

**Friday Afternoon**
- Emma completes 15th conversion
- Celebrates with badge unlock
- Gets team notification
- Sarah pushes hard, gets 3 more conversions

**Weekend**
- Final results: Emma wins (+150 points + badge)
- Sarah gets 2nd place bonus
- Team morale boost for next week

---

### Scenario 3: Reward Redemption

**Sarah accumulates 1200 points (across multiple months)**

**Decides to redeem:**
1. "Extra PTO Day" (800 pts) → Gets code REWARD-54321
2. Submits through HR portal
3. HR approves within 1 business day
4. Gets confirmation email with dates

**Later:**
1. "10% Off Merch" (500 pts) → Gets code REWARD-98765
2. Uses code at company store
3. Gets 10% discount on purchase

**Results:**
- Points remaining: 1200 - 800 - 500 = -100 (overclaimed, but allowed)
- Gets extra day off
- Gets discount on merch

---

## 📊 Analytics View

```
Weekly Leaderboard Summary:
┌─────┬──────────────┬────────┬──────────┬────────┐
│ # │ Name │ Points │ Badges │ Streak │
├─────┼──────────────┼────────┼──────────┼────────┤
│ 🥇 1 │ Sarah Chen │ 450 │ ⭐⭐ │ 🔥🔥🔥🔥🔥🔥🔥 │
│ 🥈 2 │ Mike Johnson │ 380 │ ⭐ │ 🔥🔥🔥 │
│ 🥉 3 │ Emma Davis │ 310 │ ⭐ │ 🔥🔥🔥🔥🔥 │
│ 4 │ Tom Wilson │ 245 │ │ 🔥🔥 │
│ 5 │ Alice Zhang │ 180 │ │ 🔥 │
└─────┴──────────────┴────────┴──────────┴────────┘

Total Points Distributed This Week: 1,565
Average Points per User: 313
Most Active Badge: Week Warrior (5 unlocks)
```

---

## 🚀 Configuration Tuning Examples

### High Engagement (Aggressive)
```sql
UPDATE gamification_config 
SET points_per_call = 20,
    points_per_conversion = 100,
    max_daily_points = 1000,
    points_decay_percent = 2
WHERE tenant_id = 'tenant-aggressive';
```

### Moderate Engagement (Balanced)
```sql
UPDATE gamification_config 
SET points_per_call = 10,
    points_per_conversion = 50,
    max_daily_points = 500,
    points_decay_percent = 5
WHERE tenant_id = 'tenant-balanced';
```

### Light Engagement (Casual)
```sql
UPDATE gamification_config 
SET points_per_call = 5,
    points_per_conversion = 25,
    max_daily_points = 250,
    points_decay_percent = 10
WHERE tenant_id = 'tenant-casual';
```

---

**Now you have everything to launch your gamification system!** 🎮✨
