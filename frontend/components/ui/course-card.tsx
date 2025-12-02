import React from 'react'

interface CourseCardProps {
  title: string
  description?: string
  progress?: number
  icon?: React.ReactNode
  onClick?: () => void
  className?: string
  status?: 'completed' | 'in-progress' | 'upcoming'
}

export const CourseCard = React.forwardRef<HTMLDivElement, CourseCardProps>(
  ({ title, description, progress, icon, onClick, className = '', status = 'in-progress' }, ref) => {
    const statusColor = {
      completed: 'bg-green-50 border-green-200',
      'in-progress': 'bg-blue-50 border-blue-200',
      upcoming: 'bg-gray-50 border-gray-200',
    }

    return (
      <div
        ref={ref}
        onClick={onClick}
        className={`bg-white rounded-lg p-4 md:p-5 shadow-sm border hover:shadow-md transition cursor-pointer ${statusColor[status]} ${className}`}
      >
        <div className="flex items-start gap-3">
          {icon && <div className="text-2xl flex-shrink-0">{icon}</div>}
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-gray-900 truncate">{title}</h3>
            {description && (
              <p className="text-sm text-gray-600 mt-1 line-clamp-2">{description}</p>
            )}
            {progress !== undefined && (
              <div className="mt-3">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-xs text-gray-600">Progress</span>
                  <span className="text-xs font-medium text-gray-900">{progress}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full transition"
                    style={{ width: `${progress}%` }}
                  />
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    )
  }
)

CourseCard.displayName = 'CourseCard'
