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
  <div className="bg-white rounded-sm p-6 border-b border-gray-200 hover:border-gray-300 transition-colors">
    <div className="flex items-center justify-between">
      <div>
        <p className="text-gray-600 text-xs font-semibold uppercase tracking-wide">{title}</p>
        <p className="text-2xl font-normal text-gray-900 mt-3">{value}</p>
        {trend !== undefined && (
          <p className={`text-xs mt-3 font-medium ${trend > 0 ? 'text-gray-700' : 'text-gray-500'}`}>
            {trend > 0 ? '↑' : '↓'} {Math.abs(trend)}%
          </p>
        )}
      </div>
      <div className="text-2xl text-gray-400">{Icon}</div>
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
