'use client'

import { useDashboardStats } from '@/hooks/use-dashboard'
import { formatCurrency } from '@/lib/utils'
import { FiTrendingUp, FiUsers, FiPhoneCall, FiTarget } from 'react-icons/fi'

const StatCard = ({
  title,
  value,
  icon: Icon,
  trend,
  color,
}: {
  title: string
  value: string | number
  icon: React.ReactNode
  trend?: number
  color: string
}) => (
  <div className={`bg-white rounded-lg shadow p-6 border-l-4 ${color}`}>
    <div className="flex items-center justify-between">
      <div>
        <p className="text-gray-600 text-sm font-medium">{title}</p>
        <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
        {trend !== undefined && (
          <p className={`text-xs mt-2 ${trend > 0 ? 'text-green-600' : 'text-red-600'}`}>
            {trend > 0 ? '↑' : '↓'} {Math.abs(trend)}% from last month
          </p>
        )}
      </div>
      <div className="text-3xl text-gray-400">{Icon}</div>
    </div>
  </div>
)

export const StatsOverview = () => {
  const { stats, loading } = useDashboardStats()

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="bg-gray-200 rounded-lg h-32 animate-pulse"></div>
        ))}
      </div>
    )
  }

  if (!stats) return null

  const statsArray = [
    {
      title: 'Total Leads',
      value: stats.totalLeads,
      icon: <FiUsers />,
      color: 'border-blue-500',
      trend: 12,
    },
    {
      title: 'Active Calls',
      value: stats.totalCalls,
      icon: <FiPhoneCall />,
      color: 'border-green-500',
      trend: 8,
    },
    {
      title: 'Conversion Rate',
      value: `${stats.conversionRate.toFixed(1)}%`,
      icon: <FiTrendingUp />,
      color: 'border-purple-500',
      trend: 3,
    },
    {
      title: 'Revenue',
      value: formatCurrency(stats.totalRevenue),
      icon: <FiTarget />,
      color: 'border-orange-500',
      trend: 15,
    },
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {statsArray.map((stat) => (
        <StatCard key={stat.title} {...stat} />
      ))}
    </div>
  )
}
