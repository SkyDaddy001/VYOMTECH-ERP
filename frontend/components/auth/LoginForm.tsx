'use client'

import { useState } from 'react'
import Link from 'next/link'

interface LoginFormProps {
  onSubmit: (email: string, password: string) => Promise<void>
  loading?: boolean
  apiError?: string | null
}

export default function LoginForm({ onSubmit, loading = false, apiError = null }: LoginFormProps) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [errors, setErrors] = useState<Record<string, string>>({})

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrors({})

    if (!email) setErrors((prev) => ({ ...prev, email: 'Email is required' }))
    if (!password) setErrors((prev) => ({ ...prev, password: 'Password is required' }))

    if (email && password) {
      await onSubmit(email, password)
    }
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-8">
      <h1 className="text-3xl font-bold text-center mb-8 text-gray-800">
        Call Center Login
      </h1>

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
  )
}
