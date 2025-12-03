'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import LoginForm from '@/components/auth/LoginForm'
import { DemoCredentials } from '@/app/demo-credentials'
import toast from 'react-hot-toast'

export default function LoginPage() {
  const router = useRouter()
  const { login, error } = useAuth()
  const [loading, setLoading] = useState(false)

  const handleLogin = async (email: string, password: string) => {
    try {
      setLoading(true)
      await login(email, password)
      toast.success('Login successful!')
      router.push('/dashboard')
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  const handleSelectDemoCredential = (email: string, password: string) => {
    handleLogin(email, password)
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4">
      <div className="w-full max-w-md">
        <LoginForm onSubmit={handleLogin} loading={loading} apiError={error} />
        <div className="mt-8">
          <DemoCredentials onSelectCredential={handleSelectDemoCredential} />
        </div>
      </div>
    </div>
  )
}
