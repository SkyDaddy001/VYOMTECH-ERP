'use client'

import { useState } from 'react'
import Link from 'next/link'

type TenantMode = 'create' | 'join'

interface RegisterFormProps {
  onSubmit: (data: {
    email: string
    password: string
    name: string
    tenantMode: TenantMode
    tenantName?: string
    tenantDomain?: string
    tenantCode?: string
  }) => Promise<void>
  loading?: boolean
}

export default function RegisterForm({ onSubmit, loading = false }: RegisterFormProps) {
  const [tenantMode, setTenantMode] = useState<TenantMode>('create')
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
    tenantName: '',
    tenantDomain: '',
    tenantCode: '',
  })
  const [errors, setErrors] = useState<Record<string, string>>({})

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  const handleModeChange = (mode: TenantMode) => {
    setTenantMode(mode)
    setErrors({})
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})

    const newErrors: Record<string, string> = {}
    if (!formData.email) newErrors.email = 'Email is required'
    if (!formData.password) newErrors.password = 'Password is required'
    if (!formData.confirmPassword) newErrors.confirmPassword = 'Please confirm password'
    if (!formData.name) newErrors.name = 'Name is required'

    if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match'
    }

    if (formData.password && formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters'
    }

    if (tenantMode === 'create' && !formData.tenantName) {
      newErrors.tenantName = 'Tenant name is required'
    }

    if (tenantMode === 'join' && !formData.tenantCode) {
      newErrors.tenantCode = 'Tenant code is required'
    }

    if (Object.keys(newErrors).length === 0) {
      await onSubmit({
        email: formData.email,
        password: formData.password,
        name: formData.name,
        tenantMode,
        tenantName: tenantMode === 'create' ? formData.tenantName : undefined,
        tenantDomain: tenantMode === 'create' ? formData.tenantDomain : undefined,
        tenantCode: tenantMode === 'join' ? formData.tenantCode : undefined,
      })
    } else {
      setErrors(newErrors)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-8">
      <h1 className="text-3xl font-bold text-center mb-8 text-gray-800">
        Create Account
      </h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* User Information Section */}
        <div className="pb-6 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">User Information</h2>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Full Name
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="John Doe"
              disabled={loading}
            />
            {errors.name && <p className="text-red-500 text-sm mt-1">{errors.name}</p>}
          </div>

          <div className="mt-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="your@email.com"
              disabled={loading}
            />
            {errors.email && <p className="text-red-500 text-sm mt-1">{errors.email}</p>}
          </div>

          <div className="mt-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="••••••••"
              disabled={loading}
            />
            {errors.password && <p className="text-red-500 text-sm mt-1">{errors.password}</p>}
          </div>

          <div className="mt-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Confirm Password
            </label>
            <input
              type="password"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="••••••••"
              disabled={loading}
            />
            {errors.confirmPassword && (
              <p className="text-red-500 text-sm mt-1">{errors.confirmPassword}</p>
            )}
          </div>
        </div>

        {/* Tenant Selection Section */}
        <div className="pt-6">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">Tenant Selection</h2>
          
          <div className="space-y-4">
            {/* Create New Tenant Option */}
            <div className="flex items-start gap-4">
              <div className="mt-1">
                <input
                  type="radio"
                  name="tenantMode"
                  value="create"
                  checked={tenantMode === 'create'}
                  onChange={() => handleModeChange('create')}
                  disabled={loading}
                  className="w-4 h-4"
                />
              </div>
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-2 cursor-pointer">
                  Create New Tenant
                </label>
                {tenantMode === 'create' && (
                  <div className="space-y-3 pl-4 bg-blue-50 p-3 rounded-lg">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Tenant Name *
                      </label>
                      <input
                        type="text"
                        name="tenantName"
                        value={formData.tenantName}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="e.g., Acme Corp"
                        disabled={loading}
                      />
                      {errors.tenantName && (
                        <p className="text-red-500 text-sm mt-1">{errors.tenantName}</p>
                      )}
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Tenant Domain (Optional)
                      </label>
                      <input
                        type="text"
                        name="tenantDomain"
                        value={formData.tenantDomain}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="e.g., acme.callcenter.com"
                        disabled={loading}
                      />
                      <p className="text-xs text-gray-500 mt-1">
                        Used for custom domain setup
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Join Existing Tenant Option */}
            <div className="flex items-start gap-4">
              <div className="mt-1">
                <input
                  type="radio"
                  name="tenantMode"
                  value="join"
                  checked={tenantMode === 'join'}
                  onChange={() => handleModeChange('join')}
                  disabled={loading}
                  className="w-4 h-4"
                />
              </div>
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-2 cursor-pointer">
                  Join Existing Tenant
                </label>
                {tenantMode === 'join' && (
                  <div className="pl-4 bg-green-50 p-3 rounded-lg">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Tenant Invitation Code *
                    </label>
                    <input
                      type="text"
                      name="tenantCode"
                      value={formData.tenantCode}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="Enter your tenant code"
                      disabled={loading}
                    />
                    {errors.tenantCode && (
                      <p className="text-red-500 text-sm mt-1">{errors.tenantCode}</p>
                    )}
                    <p className="text-xs text-gray-500 mt-2">
                      Ask your tenant administrator for the invitation code
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
        >
          {loading ? 'Creating Account...' : 'Sign Up'}
        </button>
      </form>

      <p className="text-center mt-6 text-gray-600">
        Already have an account?{' '}
        <Link href="/auth/login" className="text-blue-600 hover:text-blue-700 font-medium">
          Login here
        </Link>
      </p>
    </div>
  )
}
