'use client'

import { useState } from 'react'
import Link from 'next/link'

interface LoginFormProps {
  onSubmit: (email: string, password: string) => Promise<void>
  loading?: boolean
  apiError?: string | null
}

interface TestCredential {
  email: string
  password: string
  role: string
  description: string
}

const TEST_CREDENTIALS: TestCredential[] = [
  {
    email: 'demo@vyomtech.com',
    password: 'DemoPass@123',
    role: 'Admin',
    description: 'Full system access',
  },
  {
    email: 'agent@vyomtech.com',
    password: 'AgentPass@123',
    role: 'Agent',
    description: 'Agent dashboard & call management',
  },
  {
    email: 'manager@vyomtech.com',
    password: 'ManagerPass@123',
    role: 'Manager',
    description: 'Team management & reporting',
  },
  {
    email: 'sales@vyomtech.com',
    password: 'SalesPass@123',
    role: 'Sales',
    description: 'Sales pipeline & leads',
  },
  {
    email: 'hr@vyomtech.com',
    password: 'HRPass@123',
    role: 'HR Staff',
    description: 'HR & employee management',
  },
]

export default function LoginForm({ onSubmit, loading = false, apiError = null }: LoginFormProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [showTestCredentials, setShowTestCredentials] = useState(true)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})

    if (!email) setErrors((prev) => ({ ...prev, email: 'Email is required' }))
    if (!password) setErrors((prev) => ({ ...prev, password: 'Password is required' }))

    if (email && password) {
      await onSubmit(email, password)
    }
  }

  const handleQuickLogin = (testEmail: string, testPassword: string) => {
    setEmail(testEmail)
    setPassword(testPassword)
    // Trigger login after state update
    setTimeout(() => {
      onSubmit(testEmail, testPassword)
    }, 0)
  }

  return (
    <div className="space-y-6">
      {/* Main Login Form */}
      <div className="bg-white rounded-lg shadow-lg p-8">
        <h1 className="text-3xl font-bold text-center mb-2 text-gray-800">
          VYOMTECH ERP
        </h1>
        <p className="text-center text-gray-500 text-sm mb-8">
          Multi-Tenant Business Management System
        </p>

        <form onSubmit={handleSubmit} className="space-y-4">
          {apiError && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-800 text-sm font-medium">Error: {apiError}</p>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="your@email.com"
              disabled={loading}
            />
            {errors.email && <p className="text-red-500 text-sm mt-1">{errors.email}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="••••••••"
              disabled={loading}
            />
            {errors.password && <p className="text-red-500 text-sm mt-1">{errors.password}</p>}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded-lg transition disabled:opacity-50"
          >
            {loading ? 'Logging in...' : 'Sign In'}
          </button>
        </form>

        <p className="text-center mt-6 text-gray-600">
          Don't have an account?{' '}
          <Link href="/auth/register" className="text-blue-600 hover:text-blue-700 font-medium">
            Register here
          </Link>
        </p>
      </div>

      {/* Test Credentials Section */}
      <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg shadow-md border border-green-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <div className="bg-green-500 rounded-full w-8 h-8 flex items-center justify-center">
              <span className="text-white text-lg font-bold">✓</span>
            </div>
            <h2 className="text-lg font-bold text-gray-800">
              Demo Test Credentials
            </h2>
          </div>
          <button
            type="button"
            onClick={() => setShowTestCredentials(!showTestCredentials)}
            className="text-sm text-green-600 hover:text-green-700 font-medium"
          >
            {showTestCredentials ? 'Hide' : 'Show'}
          </button>
        </div>

        {showTestCredentials && (
          <div className="space-y-3">
            <p className="text-sm text-gray-600 mb-4">
              Click any credential below to auto-fill and login instantly:
            </p>
            {TEST_CREDENTIALS.map((cred, index) => (
              <div
                key={index}
                onClick={() => handleQuickLogin(cred.email, cred.password)}
                className="bg-white rounded-lg p-4 cursor-pointer hover:shadow-md transition hover:bg-blue-50 border border-green-200"
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="inline-block px-2 py-1 bg-blue-100 text-blue-700 text-xs font-semibold rounded">
                        {cred.role}
                      </span>
                      <span className="text-xs text-gray-500">{cred.description}</span>
                    </div>
                    <p className="font-mono text-sm text-gray-700">
                      <span className="text-gray-500">Email: </span>
                      <span className="font-medium">{cred.email}</span>
                    </p>
                    <p className="font-mono text-sm text-gray-700">
                      <span className="text-gray-500">Password: </span>
                      <span className="font-medium">{cred.password}</span>
                    </p>
                  </div>
                  <div className="ml-4 text-blue-500 hover:text-blue-700">
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </div>
            ))}

            <div className="mt-4 p-4 bg-blue-50 rounded-lg border border-blue-200">
              <p className="text-xs text-blue-700">
                <strong>💡 Tip:</strong> These are demo accounts for testing. All data is sample data and can be reset anytime.
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Environment Info */}
      <div className="text-center text-xs text-gray-500 space-y-1">
        <p>Environment: <span className="font-mono bg-gray-100 px-2 py-1 rounded">Development</span></p>
        <p>API: <span className="font-mono bg-gray-100 px-2 py-1 rounded">http://localhost:8080</span></p>
      </div>
    </div>
  )
}
