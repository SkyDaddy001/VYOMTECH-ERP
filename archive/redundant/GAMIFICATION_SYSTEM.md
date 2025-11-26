# 🎮 Complete Gamification & Points System Guide

Comprehensive implementation of a 100% gamified multi-tenant AI call center with points, badges, challenges, leaderboards, and rewards.

---

## 📋 Table of Contents
1. [Overview](#overview)
2. [System Components](#system-components)
3. [Points System](#points-system)
4. [Badges & Achievements](#badges--achievements)
5. [Challenges & Quests](#challenges--quests)
6. [Leaderboards](#leaderboards)
7. [Rewards & Redemption](#rewards--redemption)
8. [Levels & Progression](#levels--progression)
9. [Implementation Guide](#implementation-guide)
10. [API Reference](#api-reference)
11. [Configuration](#configuration)

---

## Overview

The gamification system transforms regular work activities into engaging, competitive experiences by:

- **Awarding Points** for every action (calls, conversions, feedback)
- **Unlocking Badges** for achievements and milestones
- **Running Challenges** to drive specific behaviors
- **Ranking Users** on leaderboards (daily, weekly, monthly, all-time)
- **Redeeming Rewards** with accumulated points
- **Progressing Through Levels** to unlock privileges
- **Maintaining Streaks** for consistency rewards

---

## System Components

### Database Tables

```
gamification_config          ← Tenant settings
├── user_points              ← User's current/lifetime points
├── point_transactions       ← Audit log of all points
├── badge                    ← Badge definitions
├── user_badge               ← Badges earned by users
├── challenge                ← Challenge definitions
├── user_challenge           ← Challenge progress
├── leaderboard              ← Pre-calculated rankings
├── player_level             ← Level definitions
├── user_level               ← User's current level
├── reward_item              ← Rewardable items
├── user_redemption          ← Redemption records
├── team_challenge           ← Team-based challenges
└── achievement_event        ← Analytics events
```

### Backend Services

**GamificationService** (`internal/services/gamification.go`)
- Points management (award, revoke, track)
- Badge operations (create, award, list)
- Challenge tracking and completion
- Leaderboard calculations
- Level progression
- Complete profile retrieval

### Frontend Components

**GamificationDashboard** (`frontend/components/dashboard/GamificationDashboard.tsx`)
- Points display and level progress
- Badge showcase
- Challenge tracker
- Leaderboard rankings
- Real-time stats

**RewardsShop** (`frontend/components/dashboard/RewardsShop.tsx`)
- Browse available rewards
- Filter by category
- Redeem for points
- Track redemption history

**PointsIndicator** (`frontend/components/dashboard/PointsIndicator.tsx`)
- Mini widget for sidebar
- Quick points check
- Daily progress
- Streak display

---

## Points System

### Point Sources

| Action | Points | Multiplier | Bonus |
|--------|--------|-----------|-------|
| Completed Call | 10 | 1x | - |
| Lead Conversion | 50 | 1x | - |
| Quality Review | 5 | 1x | - |
| Customer Feedback | 3 | 1x | - |
| Perfect Call (5★) | 10 | 2x | +10 |
| Streak Bonus | Varies | 1.5x | Active |
| Level Up | 100 | 1.25x | +25 |

### Configuration

```go
// Per-tenant settings
type GamificationConfig struct {
    PointsPerCall          int // 10
    PointsPerConversion    int // 50
    PointsPerQualityReview int // 5
    PointsPerFeedback      int // 3
    PointsDecayPercent     int // 5% per 30 days
    MaxDailyPoints         int // 500
    LeaderboardResetPeriod string // "weekly"
}
```

### Point Mechanics

**Daily Limits**
- Maximum 500 points per day
- Resets at midnight
- Prevents point farming

**Decay System**
- 5% point decay every 30 days for inactivity
- Encourages consistent engagement
- Can be customized per tenant

**Streak System**
- Increases by 1 day for daily activity
- Resets after 1 day of inactivity
- Multiplies points earned (1.5x bonus)

### Award Points Example

```go
// Backend: Award points for call completion
gs.AwardPoints(
    ctx,
    userID,
    tenantID,
    "call_completed",      // action type
    10,                    // base points
    "Completed 1 call",    // description
    "",                    // no special bonus
)

// With multiplier for perfect call
gs.AwardPoints(
    ctx,
    userID,
    tenantID,
    "perfect_call",
    10,
    "Perfect call (5★ rating)",
    "perfect_call",  // 2x multiplier
)
```

### Track Points (Frontend)

```tsx
// Get current points
const response = await fetch('/api/v1/gamification/points', {
  headers: { Authorization: `Bearer ${token}` }
})

const data = await response.json()
// {
//   current_points: 2450,
//   lifetime_points: 15000,
//   daily_points: 145,
//   streak_days: 7,
//   rank: 3
// }
```

---

## Badges & Achievements

### Badge Types

**Milestone Badges** 🏁
- 100 Calls Completed
- 1,000 Lifetime Points
- 30-Day Streak
- Level 10 Reached

**Skill Badges** 🎯
- Expert Closer
- Quality Champion
- Team Player
- Speed Demon

**Behavior Badges** ⭐
- Consistency King
- Early Bird
- Night Owl
- Marathon Runner

**Performance Badges** 🚀
- Top Performer
- Rising Star
- Undefeated
- Perfect Week

**Social Badges** 🤝
- Helpful Teammate
- Mentor
- Community Champion
- Social Butterfly

### Rarity Levels

```
Common      - Easy to earn, blue
Uncommon    - Moderate effort, green
Rare        - Challenging, blue
Epic        - Very difficult, purple
Legendary   - Extremely rare, gold
```

### Badge Rewards

```go
type Badge struct {
    ID               int64   // Unique ID
    Code             string  // "100_calls_completed"
    Name             string  // "Century Club"
    Description      string  // "Complete 100 calls"
    Category         string  // "milestone"
    RequirementType  string  // "total_calls"
    RequirementValue int     // 100
    PointsReward     int     // 250
    Rarity           string  // "epic"
    IconURL          string  // URL to badge icon
}
```

### Create Badge Example

```go
// Create badge
badge := &models.Badge{
    TenantID: tenantID,
    Code: "100_calls_completed",
    Name: "Century Club",
    Description: "Complete 100 calls",
    Category: "milestone",
    RequirementType: "total_calls",
    RequirementValue: 100,
    PointsReward: 250,
    Rarity: "epic",
}

err := gamificationService.CreateBadge(ctx, badge)
```

### Award Badge Example

```go
// Award badge when requirement met
if userCalls >= 100 {
    err := gamificationService.AwardBadge(ctx, userID, badgeID, tenantID)
    // Awards bonus points automatically
}
```

---

## Challenges & Quests

### Challenge Types

**Daily Challenges** 📅
- Reset daily
- Easy to moderate difficulty
- 10-25 points reward
- Examples: "10 calls", "3 conversions"

**Weekly Challenges** 📆
- Span Monday-Sunday
- Moderate difficulty
- 50-100 points reward
- Examples: "50 calls", "10 conversions"

**Monthly Challenges** 📊
- Full calendar month
- Hard difficulty
- 200-500 points reward
- Examples: "200 calls", "50 conversions", "90% satisfaction"

**Seasonal Challenges** 🎯
- Span 3 months
- Very difficult
- 1000+ points reward
- Examples: "Q1 Top Performer"

**Special Challenges** 🌟
- Limited time events
- Variable difficulty
- Special prizes
- Examples: "Black Friday Blitz"

### Challenge Structure

```go
type Challenge struct {
    ID              int64     // Unique ID
    Name            string    // "Call Blitz"
    Description     string    // "Complete 50 calls"
    ChallengeType   string    // "weekly"
    ObjectiveType   string    // "total_calls"
    ObjectiveTarget int       // 50
    PointsReward    int       // 100
    BadgeRewardID   *int64    // Optional badge
    Difficulty      string    // "medium"
    StartDate       time.Time
    EndDate         time.Time
}
```

### Track Progress Example

```tsx
// Frontend: Show challenge progress
<ChallengeCard>
  <Title>Call Blitz</Title>
  <Progress>{progress} / 50</Progress>
  <ProgressBar percent={(progress / 50) * 100} />
  <Reward>+100 pts</Reward>
</ChallengeCard>
```

---

## Leaderboards

### Leaderboard Periods

**Daily Leaderboard** 📊
- Resets at midnight
- Based on daily points
- Shows top 50 users
- Updated hourly

**Weekly Leaderboard** 📈
- Monday-Sunday
- Based on weekly points
- Shows top 100 users
- Updated every 6 hours

**Monthly Leaderboard** 📉
- Calendar month
- Based on monthly points
- Shows all users
- Updated daily

**All-Time Leaderboard** ⭐
- Lifetime points
- Lifetime rank
- Shows all users
- Updated daily

### Leaderboard Entry

```go
type LeaderboardEntry struct {
    Rank        int   // 1-N
    UserID      int64
    Points      int
    PointsChange int  // +50 vs yesterday
    PreviousRank *int // Was rank 5
    BadgesCount int
    StreakDays  int
}
```

### Display Example

```
🥇 #1  Sarah Chen        2,450 pts  ⭐15  🔥7
🥈 #2  Mike Johnson      2,380 pts  ⭐12  🔥3
🥉 #3  Emma Davis       2,310 pts  ⭐10  🔥5
   #4  Tom Wilson       2,245 pts  ⭐9   🔥2
```

### Fetch Leaderboard (Frontend)

```tsx
const response = await fetch(
  '/api/v1/gamification/leaderboard?period=weekly&limit=50',
  { headers: { Authorization: `Bearer ${token}` } }
)

const data = await response.json()
// {
//   entries: [...LeaderboardEntry],
//   user_rank: {...LeaderboardEntry},
//   period_type: "weekly"
// }
```

---

## Rewards & Redemption

### Reward Categories

**Discounts** 💰
- 10% off product
- Free month subscription
- Bulk discount

**Digital** 📱
- Premium features unlock
- Extra storage
- API credits

**Experiences** 🎉
- Training workshop
- Conference ticket
- Team lunch

**Physical** 📦
- Company swag
- Gift card
- Merchandise

### Reward Mechanics

```go
type RewardItem struct {
    ID                int64
    Name              string              // "10% Company Swag Discount"
    Description       string              // "Get 10% off merch"
    PointsCost        int                 // 500
    Category          string              // "discount"
    RedemptionType    string              // "discount"
    RedemptionDetails map[string]interface{} // {code: "XYZ123"}
    Stock             int                 // -1 = unlimited
    Featured          bool
}
```

### Redemption Flow

```
1. User views rewards shop
2. User selects reward
3. Check: points >= cost
4. Check: stock available
5. Deduct points
6. Mark as pending
7. Admin approves
8. Generate code/link
9. Send to user
10. User marks as completed
```

### Redeem Example

```tsx
// Frontend: Redeem reward
const handleRedeem = async (rewardId: number) => {
  const response = await fetch('/api/v1/gamification/redeem', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify({
      reward_id: rewardId,
      quantity: 1
    })
  })

  if (response.ok) {
    const data = await response.json()
    // {
    //   success: true,
    //   redemption_code: "REWARD-123456",
    //   expiry_date: "2025-12-22"
    // }
    showCode(data.redemption_code)
  }
}
```

---

## Levels & Progression

### Level System

Users progress through 10 levels based on lifetime points:

```
Level 1:   0 -     999 pts   🟤 Bronze
Level 2: 1,000 -   1,999 pts 🟤 Bronze
Level 3: 2,000 -   2,999 pts 🟤 Bronze
Level 4: 3,000 -   3,999 pts 🥈 Silver
Level 5: 4,000 -   4,999 pts 🥈 Silver
Level 6: 5,000 -   5,999 pts 🥈 Silver
Level 7: 6,000 -   6,999 pts 🥇 Gold
Level 8: 7,000 -   7,999 pts 🥇 Gold
Level 9: 8,000 -   8,999 pts 🥇 Gold
Level 10: 9,000+   pts      👑 Legend
```

### Level Up Rewards

When leveling up:
- Award 100 bonus points (1.25x multiplier)
- Unlock level-specific badge
- Increase daily point limit (by 50 per level)
- Display celebration notification

### Level Definition

```go
type PlayerLevel struct {
    LevelNumber         int    // 1-10
    Name                string // "Bronze"
    PointsRequired      int    // Points needed this level
    PointsTotalRequired int    // Total lifetime points needed
    IconURL             string // Badge icon
    UnlockPrivileges    JSON   // Special abilities
}
```

---

## Implementation Guide

### 1. Backend Setup

**Create database tables:**
```bash
mysql < migrations/003_gamification_system.sql
```

**Initialize gamification for tenant:**
```go
config := &models.GamificationConfig{
    TenantID: tenantID,
    Enabled: true,
    PointsPerCall: 10,
    PointsPerConversion: 50,
    MaxDailyPoints: 500,
}
// Insert into database
```

**Create initial badges:**
```go
badges := []models.Badge{
    {Code: "first_call", Name: "Welcome", PointsReward: 10},
    {Code: "100_calls", Name: "Century Club", PointsReward: 250},
    {Code: "streak_7", Name: "Week Warrior", PointsReward: 100},
}
for _, badge := range badges {
    gamificationService.CreateBadge(ctx, &badge)
}
```

**Create initial challenges:**
```go
challenges := []models.Challenge{
    {
        ChallengeType: "daily",
        ObjectiveType: "total_calls",
        ObjectiveTarget: 10,
        PointsReward: 25,
    },
    // ... more challenges
}
```

### 2. Hook Points in Existing Code

**When call completes:**
```go
// In call_handler.go
func (h *CallHandler) CompleteCall(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Award gamification points
    gamificationService.AwardPoints(
        ctx,
        call.AgentID,
        tenantID,
        "call_completed",
        10,
        fmt.Sprintf("Completed call: %s", callID),
        "",
    )
}
```

**When lead converts:**
```go
// In lead_handler.go
func (h *LeadHandler) ConvertLead(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Award conversion bonus
    gamificationService.AwardPoints(
        ctx,
        lead.AssignedAgentID,
        tenantID,
        "lead_converted",
        50,
        fmt.Sprintf("Lead converted: %s", leadID),
        "",
    )
}
```

### 3. Add Frontend Components

**Add to dashboard layout:**
```tsx
// frontend/app/layout.tsx
import { PointsIndicator } from '@/components/dashboard/PointsIndicator'

export default function RootLayout({ children }) {
  return (
    <>
      <header>
        {/* ... */}
        <PointsIndicator />
      </header>
      {children}
    </>
  )
}
```

**Add navigation link:**
```tsx
// frontend/components/layouts/Sidebar.tsx
<nav>
  <Link href="/dashboard/gamification">
    🎮 Gamification
  </Link>
</nav>
```

### 4. Deployment

```bash
# Build backend
go build -o bin/main cmd/main.go

# Run migrations
mysql < migrations/003_gamification_system.sql

# Restart services
podman-compose down && podman-compose up -d
```

---

## API Reference

### Points Endpoints

```bash
# Get user's current points
GET /api/v1/gamification/points
Response: { current_points: 2450, daily_points: 145, streak_days: 7 }

# Award points (admin)
POST /api/v1/gamification/points/award
{ user_id: 123, points: 50, action_type: "conversion" }

# Get point history
GET /api/v1/gamification/transactions
Response: [{ points: 10, action: "call", timestamp: "..." }]
```

### Badge Endpoints

```bash
# Get user's badges
GET /api/v1/gamification/badges
Response: [{ id: 1, name: "Century Club", earned_date: "..." }]

# Get badge list
GET /api/v1/gamification/badges/available
Response: [{ id: 1, name: "...", requirement: "..." }]

# Award badge (admin)
POST /api/v1/gamification/badges/{id}/award
{ user_id: 123 }
```

### Challenge Endpoints

```bash
# Get active challenges
GET /api/v1/gamification/challenges
Response: [{ id: 1, name: "Call Blitz", progress: 25, target: 50 }]

# Get user challenge progress
GET /api/v1/gamification/challenges/{id}/progress
Response: { progress: 25, target: 50, completed: false }
```

### Leaderboard Endpoints

```bash
# Get leaderboard
GET /api/v1/gamification/leaderboard?period=weekly&limit=50
Response: { 
  entries: [...], 
  user_rank: { rank: 3, points: 2450 }
}

# Update leaderboard
POST /api/v1/gamification/leaderboard/update
```

### Rewards Endpoints

```bash
# Get available rewards
GET /api/v1/gamification/rewards
Response: [{ id: 1, name: "Discount", cost: 500 }]

# Redeem reward
POST /api/v1/gamification/redeem
{ reward_id: 1, quantity: 1 }
Response: { code: "REWARD-123", expiry: "..." }

# Get redemption history
GET /api/v1/gamification/redemptions
```

### Gamification Profile

```bash
# Get complete profile
GET /api/v1/gamification/profile
Response: {
  user_points: {...},
  current_level: {...},
  badges: [...],
  active_challenges: [...],
  leaderboard_rank: {...}
}
```

---

## Configuration

### Per-Tenant Settings

```sql
INSERT INTO gamification_config (
    tenant_id,
    enabled,
    points_per_call,
    points_per_conversion,
    points_per_quality_review,
    points_per_feedback,
    points_decay_percent,
    decay_period_days,
    max_daily_points,
    leaderboard_reset_period
) VALUES (
    'tenant-123',
    true,
    10,      -- points per call
    50,      -- points per conversion
    5,       -- points per quality review
    3,       -- points per feedback
    5,       -- 5% decay
    30,      -- every 30 days
    500,     -- max daily
    'weekly' -- leaderboard period
);
```

### Customization Options

**Aggressive Gamification** (High engagement)
```
Points per call: 20
Points per conversion: 100
Max daily: 1000
Decay: 2% every 60 days
Challenges: Daily + Weekly
```

**Moderate Gamification** (Balanced)
```
Points per call: 10
Points per conversion: 50
Max daily: 500
Decay: 5% every 30 days
Challenges: Daily + Weekly + Monthly
```

**Light Gamification** (Low pressure)
```
Points per call: 5
Points per conversion: 25
Max daily: 250
Decay: 10% every 14 days
Challenges: Weekly only
```

---

## Analytics & Monitoring

### Key Metrics

```
- Average points per user per day
- Badge unlocking rate
- Challenge completion rate
- Reward redemption rate
- Leaderboard viewing frequency
- Level progression speed
- Point farming detection
```

### Dashboards

**Engagement Dashboard**
- Daily active participants
- Points awarded/day
- Challenges active
- Badge unlocks

**Behavior Dashboard**
- Top point earners
- Most earned badges
- Popular rewards
- Challenge success rates

---

## Best Practices

1. **Balance Difficulty**
   - Easy: 50% of users should complete
   - Medium: 25% of users should complete
   - Hard: 5% of users should complete

2. **Regular Updates**
   - New badges monthly
   - Seasonal challenges
   - Limited-time events
   - Trending rewards

3. **Fairness**
   - Transparent point calculations
   - No point buying (earn only)
   - Equal opportunity for all
   - Detect and prevent abuse

4. **Engagement**
   - Daily notifications for streaks
   - Level up celebrations
   - Badge unlocking notifications
   - Leaderboard position changes

5. **Feedback**
   - Show points earned immediately
   - Display level-up celebrations
   - Notify badge unlocks
   - Show rank changes

---

## Troubleshooting

### Points not updating
- Check auth token validity
- Verify user exists in database
- Check tenant ID matches
- Review point transaction logs

### Leaderboard not updating
- Run manual update: `POST /api/v1/gamification/leaderboard/update`
- Check scheduled job is running
- Verify user_points table has data

### Badges not awarding
- Check badge exists
- Verify user meets requirement
- Check badge not already awarded
- Review achievement_event logs

---

## Summary

Your system now has:
- ✅ Complete points system with decay and streaks
- ✅ 40+ badge types across 5 categories
- ✅ Daily/Weekly/Monthly challenges
- ✅ Multi-period leaderboards
- ✅ Reward redemption shop
- ✅ 10-level progression system
- ✅ Team challenges
- ✅ Full API implementation
- ✅ Beautiful frontend components
- ✅ Analytics tracking

**Start earning points immediately!** 🎮
