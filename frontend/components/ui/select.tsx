import React, { useState } from 'react'

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  onValueChange?: (value: string) => void
}

export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ onValueChange, onChange, ...props }, ref) => (
    <select
      ref={ref}
      className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${
        props.className || ''
      }`}
      onChange={(e) => {
        onValueChange?.(e.target.value)
        onChange?.(e)
      }}
      {...props}
    />
  )
)
Select.displayName = 'Select'

export const SelectContent = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ ...props }, ref) => <div ref={ref} {...props} />
)
SelectContent.displayName = 'SelectContent'

export const SelectItem = ({ value, children }: { value: string; children: React.ReactNode }) => (
  <option value={value}>{children}</option>
)

export const SelectTrigger = React.forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ ...props }, ref) => <div ref={ref} {...props} />
)
SelectTrigger.displayName = 'SelectTrigger'

export const SelectValue = ({ placeholder }: { placeholder?: string }) => <span>{placeholder}</span>
