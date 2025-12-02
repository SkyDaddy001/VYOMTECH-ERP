import React from 'react'

interface StatCardProps {
  label: string
  value: string | number
  icon?: React.ReactNode
  trend?: 'up' | 'down'
  trendValue?: string
  className?: string
}

export const StatCard = React.forwardRef<HTMLDivElement, StatCardProps>(
  ({ label, value, icon, trend, trendValue, className = '' }, ref) => {
    return (
      <div
        ref={ref}
        className={`bg-white rounded-lg p-4 md:p-6 shadow-sm border border-gray-100 hover:shadow-md transition ${className}`}
      >
        <div className="flex items-start justify-between">
          <div>
            <p className="text-xs md:text-sm text-gray-600 font-medium">{label}</p>
            <p className="text-2xl md:text-3xl font-bold text-gray-900 mt-2">{value}</p>
            {trend && trendValue && (
              <p
                className={`text-xs mt-2 ${
                  trend === 'up' ? 'text-green-600' : 'text-red-600'
                }`}
              >
                {trend === 'up' ? '↑' : '↓'} {trendValue}
              </p>
            )}
          </div>
          {icon && <div className="text-3xl opacity-80">{icon}</div>}
        </div>
      </div>
    )
  }
)

StatCard.displayName = 'StatCard'
