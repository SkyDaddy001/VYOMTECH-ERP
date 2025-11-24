'use client'

import { useAuth } from '@/hooks/useAuth'
import { redirect } from 'next/navigation'
import DashboardLayout from '@/components/layouts/DashboardLayout'
import DashboardContent from '@/components/dashboard/DashboardContent'
import GamificationDashboard from '@/components/dashboard/GamificationDashboard'

export default function DashboardPage() {
  const { user, loading } = useAuth()

  if (loading) {
    return <div className="flex items-center justify-center h-screen">Loading...</div>
  }

  if (!user) {
    redirect('/auth/login')
  }

  return (
    <DashboardLayout>
      <DashboardContent />
      <div className="mt-8">
        <GamificationDashboard />
      </div>
    </DashboardLayout>
  )
}
