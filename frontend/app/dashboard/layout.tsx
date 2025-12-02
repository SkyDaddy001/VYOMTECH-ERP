'use client'

import { ReactNode } from 'react'
import { useAuth } from '@/hooks/useAuth'
import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import DashboardLayout from '@/components/layouts/DashboardLayout'

interface DashboardLayoutWrapperProps {
  children: ReactNode
}

export default function DashboardLayoutWrapper({ children }: DashboardLayoutWrapperProps) {
  const { user, loading } = useAuth()
  const router = useRouter()
  const [hasRedirected, setHasRedirected] = useState(false)

  useEffect(() => {
    // Only redirect if auth check is done, user is not logged in, and we haven't already redirected
    if (!loading && !user && !hasRedirected) {
      console.warn('No user on dashboard, redirecting to login')
      setHasRedirected(true)
      router.push('/auth/login')
    }
  }, [user, loading, router, hasRedirected])

  if (loading) {
    return <div className="flex items-center justify-center h-screen">Loading...</div>
  }

  if (!user) {
    // Show nothing while redirecting
    return null
  }

  return <DashboardLayout>{children}</DashboardLayout>
}
