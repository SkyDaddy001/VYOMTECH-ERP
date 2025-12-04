'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'

export default function Home() {
  const router = useRouter()

  useEffect(() => {
    // Check if user is authenticated
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null
    
    if (token) {
      // User is authenticated, go to dashboard
      router.push('/dashboard')
    } else {
      // No token, go to login
      router.push('/login')
    }
  }, [router])

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">VYOMTECH ERP</h1>
        <p className="text-gray-600">Redirecting...</p>
      </div>
    </div>
  )
}
