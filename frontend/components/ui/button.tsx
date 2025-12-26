import React from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  variant?: 'default' | 'primary' | 'secondary' | 'danger' | 'success' | 'warning' | 'outline';
  size?: 'sm' | 'md' | 'lg';
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className = '',
      children,
      variant = 'default',
      size = 'md',
      disabled = false,
      ...props
    },
    ref
  ) => {
    const variantClasses = {
      default: 'bg-gray-200 hover:bg-gray-300 text-gray-900',
      primary: 'bg-blue-600 hover:bg-blue-700 text-white',
      secondary: 'bg-gray-100 hover:bg-gray-200 text-gray-900',
      danger: 'bg-red-600 hover:bg-red-700 text-white',
      success: 'bg-green-600 hover:bg-green-700 text-white',
      warning: 'bg-yellow-600 hover:bg-yellow-700 text-white',
      outline: 'border border-gray-300 hover:bg-gray-50 text-gray-900',
    };

    const sizeClasses = {
      sm: 'px-3 py-1 text-sm',
      md: 'px-4 py-2 text-base',
      lg: 'px-6 py-3 text-lg',
    };

    const disabledClass = disabled ? 'opacity-50 cursor-not-allowed' : '';

    return (
      <button
        ref={ref}
        className={`rounded font-medium transition-colors ${variantClasses[variant]} ${sizeClasses[size]} ${disabledClass} ${className}`}
        disabled={disabled}
        {...props}
      >
        {children}
      </button>
    );
  }
);
Button.displayName = 'Button';
