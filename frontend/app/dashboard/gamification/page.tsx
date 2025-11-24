import { Metadata } from 'next'
import GamificationDashboard from '@/components/dashboard/GamificationDashboard'
import { RewardsShop } from '@/components/dashboard/RewardsShop'

export const metadata: Metadata = {
  title: 'Gamification | Call Center',
  description: 'Points, badges, challenges, and rewards'
}

export default function GamificationPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold mb-2">🎮 Gamification Center</h1>
        <p className="text-gray-600">Earn points, unlock badges, and compete on leaderboards</p>
      </div>

      <div className="tabs">
        <input type="radio" name="tabs" id="tab1" checked hidden defaultChecked />
        <label htmlFor="tab1" className="tab-label">Dashboard</label>
        <div className="tab-content">
          <GamificationDashboard />
        </div>

        <input type="radio" name="tabs" id="tab2" hidden />
        <label htmlFor="tab2" className="tab-label">Rewards Shop</label>
        <div className="tab-content">
          <RewardsShop />
        </div>
      </div>
    </div>
  )
}
