'use client'

import { useState } from 'react'
import { SectionCard } from '@/components/ui/section-card'
import { StatCard } from '@/components/ui/stat-card'

interface Campaign {
  id: string
  name: string
  status: 'planned' | 'active' | 'completed' | 'paused'
  startDate: string
  endDate: string
  budget: number
  spent: number
  leads: number
  conversions: number
}

export default function CampaignsPage() {
  const [campaigns] = useState<Campaign[]>([
    { id: '1', name: 'Q4 Email Campaign', status: 'active', startDate: '2025-10-01', endDate: '2025-12-31', budget: 50000, spent: 35000, leads: 450, conversions: 35 },
    { id: '2', name: 'Product Launch - Feb 2026', status: 'planned', startDate: '2026-02-01', endDate: '2026-02-28', budget: 75000, spent: 0, leads: 0, conversions: 0 },
    { id: '3', name: 'Referral Program Q3', status: 'completed', startDate: '2025-07-01', endDate: '2025-09-30', budget: 30000, spent: 28500, leads: 320, conversions: 48 },
    { id: '4', name: 'Social Media Campaign', status: 'active', startDate: '2025-11-01', endDate: '2025-11-30', budget: 40000, spent: 12000, leads: 280, conversions: 18 },
  ])

  const stats = [
    { label: 'Total Campaigns', value: campaigns.length, icon: '🎯' },
    { label: 'Active Now', value: campaigns.filter(c => c.status === 'active').length, icon: '▶️' },
    { label: 'Total Leads', value: campaigns.reduce((sum, c) => sum + c.leads, 0), icon: '📋' },
    { label: 'Total Conversions', value: campaigns.reduce((sum, c) => sum + c.conversions, 0), icon: '✅' },
  ]

  const statusColors = {
    'planned': 'bg-gray-100 text-gray-800',
    'active': 'bg-green-100 text-green-800',
    'completed': 'bg-blue-100 text-blue-800',
    'paused': 'bg-yellow-100 text-yellow-800',
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-gradient-to-r from-purple-600 to-purple-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Campaigns</h1>
        <p className="text-purple-100 mt-2">Create and manage marketing campaigns</p>
      </div>

      {/* KPIs */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat, i) => (
          <StatCard key={i} label={stat.label} value={stat.value} icon={stat.icon} />
        ))}
      </div>

      {/* Campaigns List */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {campaigns.map((campaign) => (
          <SectionCard key={campaign.id} title={campaign.name} action={
            <span className={`text-xs font-medium px-2 py-1 rounded ${statusColors[campaign.status]}`}>
              {campaign.status.charAt(0).toUpperCase() + campaign.status.slice(1)}
            </span>
          }>
            <div className="space-y-3">
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Period:</span>
                <span className="font-medium">{campaign.startDate} to {campaign.endDate}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-600">Budget Spent:</span>
                <span className="font-medium">${campaign.spent.toLocaleString()} / ${campaign.budget.toLocaleString()}</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-blue-600 h-2 rounded-full" style={{ width: `${(campaign.spent / campaign.budget) * 100}%` }} />
              </div>
              <div className="grid grid-cols-3 gap-2 pt-2">
                <div className="text-center">
                  <p className="text-2xl font-bold text-gray-900">{campaign.leads}</p>
                  <p className="text-xs text-gray-600">Leads</p>
                </div>
                <div className="text-center">
                  <p className="text-2xl font-bold text-gray-900">{campaign.conversions}</p>
                  <p className="text-xs text-gray-600">Conversions</p>
                </div>
                <div className="text-center">
                  <p className="text-2xl font-bold text-gray-900">{campaign.leads > 0 ? ((campaign.conversions / campaign.leads) * 100).toFixed(1) : 0}%</p>
                  <p className="text-xs text-gray-600">Conversion</p>
                </div>
              </div>
            </div>
          </SectionCard>
        ))}
      </div>
    </div>
  )
}
