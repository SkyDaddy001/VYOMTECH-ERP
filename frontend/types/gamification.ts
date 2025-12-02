export interface GamificationRule {
  id: string
  rule_name: string
  rule_type: 'achievement' | 'milestone' | 'competition' | 'challenge'
  description: string
  points_awarded: number
  badge_name?: string
  badge_icon?: string
  trigger_condition: string
  is_active: boolean
  created_at?: string
}

export interface UserAchievement {
  id: string
  user_id: string
  user_name: string
  achievement_id: string
  achievement_name: string
  achievement_date: string
  points_earned: number
  badge_unlocked: boolean
  badge_name?: string
}

export interface Leaderboard {
  id: string
  user_id: string
  rank: number
  user_name: string
  total_points: number
  achievements_count: number
  badges_count: number
  current_streak: number
  last_activity: string
}

export interface Challenge {
  id: string
  challenge_name: string
  description: string
  start_date: string
  end_date: string
  challenge_type: 'individual' | 'team' | 'department'
  objective: string
  reward_points: number
  participation_count: number
  winner?: string
  status: 'active' | 'completed' | 'upcoming' | 'cancelled'
}

export interface Badge {
  id: string
  badge_name: string
  description: string
  icon_url: string
  requirement: string
  points_value: number
  rarity: 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary'
  total_earned: number
}

export interface GamificationDashboard {
  total_users_participated: number
  total_points_distributed: number
  active_challenges: number
  completed_challenges: number
  total_badges: number
  top_performer_name: string
  top_performer_points: number
  average_user_points: number
  engagement_rate_percentage: number
}
