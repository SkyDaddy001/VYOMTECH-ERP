'use client'

import { ReactNode, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/hooks/useAuth'
import { TenantSwitcher } from '@/components/dashboard/TenantSwitcher'

interface DashboardLayoutProps {
  children: ReactNode
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const router = useRouter()
  const { user, logout } = useAuth()
  const [sidebarOpen, setSidebarOpen] = useState(true)

  const handleLogout = () => {
    logout()
    router.push('/auth/login')
  }

  const menuItems = [
    { href: '/dashboard', label: 'Dashboard', icon: '📊' },
    { href: '/dashboard/agents', label: 'Agents', icon: '👥' },
    { href: '/dashboard/calls', label: 'Calls', icon: '📞' },
    { href: '/dashboard/leads', label: 'Leads', icon: '📋' },
    { href: '/dashboard/campaigns', label: 'Campaigns', icon: '🎯' },
    { href: '/dashboard/workflows', label: 'Workflows', icon: '⚙️' },
    { href: '/dashboard/scheduled-tasks', label: 'Scheduled Tasks', icon: '⏱️' },
    { href: '/dashboard/reports', label: 'Reports', icon: '📈' },
  ]

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <aside
        className={`${
          sidebarOpen ? 'w-64' : 'w-20'
        } bg-gradient-to-b from-blue-900 to-blue-800 text-white transition-all duration-300 flex flex-col overflow-y-auto`}
      >
        <div className="p-4 border-b border-blue-700">
          <div className="flex items-center justify-between">
            {sidebarOpen && <h1 className="text-xl font-bold">CallCenter</h1>}
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="p-2 hover:bg-blue-700 rounded-lg transition"
            >
              {sidebarOpen ? '←' : '→'}
            </button>
          </div>
        </div>

        {/* Tenant Switcher */}
        {sidebarOpen && <TenantSwitcher />}        <nav className="flex-1 overflow-y-auto py-4">
          {menuItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="flex items-center px-4 py-3 hover:bg-blue-700 transition"
            >
              <span className="text-xl">{item.icon}</span>
              {sidebarOpen && <span className="ml-3">{item.label}</span>}
            </Link>
          ))}
        </nav>

        <div className="p-4 border-t border-blue-700">
          <button
            onClick={handleLogout}
            className="w-full flex items-center px-3 py-2 bg-red-600 hover:bg-red-700 rounded-lg transition text-sm font-medium"
          >
            <span>🚪</span>
            {sidebarOpen && <span className="ml-2">Logout</span>}
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="bg-white shadow-sm px-6 py-4 flex items-center justify-between">
          <h2 className="text-2xl font-bold text-gray-800">Multi-Tenant AI Call Center</h2>
          <div className="flex items-center gap-4">
            <span className="text-gray-600">{user?.email}</span>
            <span className="inline-block w-10 h-10 bg-blue-500 rounded-full text-white flex items-center justify-center font-bold">
              {user?.email?.charAt(0).toUpperCase()}
            </span>
          </div>
        </header>

        {/* Content Area */}
        <main className="flex-1 overflow-auto p-6">
          {children}
        </main>
      </div>
    </div>
  )
}
