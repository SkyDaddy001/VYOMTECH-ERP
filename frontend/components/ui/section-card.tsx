import React from 'react'

interface SectionCardProps {
  title: string
  children: React.ReactNode
  className?: string
  action?: React.ReactNode
}

export const SectionCard = React.forwardRef<HTMLDivElement, SectionCardProps>(
  ({ title, children, className = '', action }, ref) => {
    return (
      <div
        ref={ref}
        className={`bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden ${className}`}
      >
        <div className="px-4 md:px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900">{title}</h2>
          {action && <div>{action}</div>}
        </div>
        <div className="px-4 md:px-6 py-4">
          {children}
        </div>
      </div>
    )
  }
)

SectionCard.displayName = 'SectionCard'
