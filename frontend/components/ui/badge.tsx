import React from 'react';

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  children: React.ReactNode;
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'danger';
}

export const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
  ({ className = '', children, variant = 'default', ...props }, ref) => {
    const variantClasses = {
      default: 'bg-gray-100 text-gray-900',
      primary: 'bg-blue-100 text-blue-900',
      secondary: 'bg-gray-200 text-gray-800',
      success: 'bg-green-100 text-green-900',
      warning: 'bg-yellow-100 text-yellow-900',
      danger: 'bg-red-100 text-red-900',
    };

    return (
      <span
        ref={ref}
        className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium ${variantClasses[variant]} ${className}`}
        {...props}
      >
        {children}
      </span>
    );
  }
);
Badge.displayName = 'Badge';
