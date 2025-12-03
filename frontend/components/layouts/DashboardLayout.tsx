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
    { href: '/dashboard', label: 'Overview', icon: '📊' },
    { href: '/dashboard/sales', label: 'Sales', icon: '📈' },
    { href: '/dashboard/presales', label: 'Pre-Sales', icon: '🎯' },
    { href: '/dashboard/leads', label: 'Leads', icon: '🔍' },
    { href: '/dashboard/accounts', label: 'Finance', icon: '💰' },
    { href: '/dashboard/ledgers', label: 'Ledgers', icon: '📑' },
    { href: '/dashboard/purchase', label: 'Purchase', icon: '📦' },
    { href: '/dashboard/hr', label: 'HR', icon: '👨‍💼' },
    { href: '/dashboard/projects', label: 'Projects', icon: '📌' },
    { href: '/dashboard/construction', label: 'Construction', icon: '🏗️' },
    { href: '/dashboard/real-estate', label: 'Real Estate', icon: '🏢' },
    { href: '/dashboard/civil', label: 'Civil', icon: '🏛️' },
    { href: '/dashboard/units', label: 'Units', icon: '🏠' },
    { href: '/dashboard/bookings', label: 'Bookings', icon: '📅' },
    { href: '/dashboard/marketing', label: 'Marketing', icon: '📣' },
    { href: '/dashboard/campaigns', label: 'Campaigns', icon: '🎪' },
    { href: '/dashboard/workflows', label: 'Workflows', icon: '⚙️' },
    { href: '/dashboard/reports', label: 'Reports', icon: '📋' },
    { href: '/dashboard/users', label: 'Users', icon: '👥' },
    { href: '/dashboard/tenants', label: 'Tenants', icon: '🏢' },
    { href: '/dashboard/company', label: 'Company', icon: '🏛️' },
    { href: '/dashboard/agents', label: 'Agents', icon: '📞' },
    { href: '/dashboard/calls', label: 'Calls', icon: '☎️' },
  ]

  return (
    <div className="flex flex-col md:flex-row h-screen bg-white">
      {/* Mobile Header */}
      <div className="md:hidden bg-black text-white px-4 py-3 flex items-center justify-between">
        <h1 className="text-lg font-bold">VYOM ERP</h1>
        <button
          onClick={() => setSidebarOpen(!sidebarOpen)}
          className="text-2xl"
        >
          ☰
        </button>
      </div>

      {/* Sidebar - Minimalistic */}
      <aside
        className={`${
          sidebarOpen ? 'block' : 'hidden'
        } md:block w-full md:w-16 bg-black text-white flex flex-col`}
      >
        {/* Logo */}
        <div className="hidden md:flex items-center justify-center h-16 border-b border-gray-800">
          <span className="text-2xl font-bold">F.</span>
        </div>

        {/* Navigation */}
        <nav className="flex md:flex-col flex-1 overflow-x-auto md:overflow-y-auto py-2 md:py-4">
          {menuItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="flex-shrink-0 flex md:flex-col items-center justify-center px-4 md:px-3 py-4 hover:bg-gray-900 transition text-center md:text-left"
              title={item.label}
            >
              <span className="text-2xl">{item.icon}</span>
              <span className="hidden md:hidden text-xs mt-1 md:block">{item.label}</span>
            </Link>
          ))}
        </nav>

        {/* Logout Button */}
        <div className="border-t border-gray-800 p-3">
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center md:justify-start px-3 py-2 hover:bg-gray-900 rounded transition text-sm font-medium"
            title="Logout"
          >
            <span className="text-xl">🚪</span>
          </button>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header - Clean and Minimal */}
        <header className="hidden md:flex bg-white border-b border-gray-200 px-6 py-4 items-center justify-between">
          <div className="flex items-center gap-4">
            <h2 className="text-xl font-semibold text-gray-900">Hello {user?.email?.split('@')[0]}!</h2>
            <span className="text-gray-400">It's good to see you again.</span>
          </div>
          <div className="flex items-center gap-4">
            <button className="text-gray-600 hover:text-gray-900 transition">🔔</button>
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full text-white flex items-center justify-center font-bold text-sm">
              {user?.email?.charAt(0).toUpperCase()}
            </div>
          </div>
        </header>

        {/* Quick Shortcuts Bar */}
        <div className="bg-white border-b border-gray-200 px-6 py-3 hidden md:block">
          <div className="flex items-center gap-2 overflow-x-auto">
            <span className="text-xs font-semibold text-gray-600 whitespace-nowrap">Quick Links:</span>
            <Link href="/dashboard/sales" className="px-3 py-1 bg-blue-50 hover:bg-blue-100 text-blue-700 text-xs rounded transition whitespace-nowrap">📈 Sales</Link>
            <Link href="/dashboard/leads" className="px-3 py-1 bg-green-50 hover:bg-green-100 text-green-700 text-xs rounded transition whitespace-nowrap">🔍 Leads</Link>
            <Link href="/dashboard/calls" className="px-3 py-1 bg-purple-50 hover:bg-purple-100 text-purple-700 text-xs rounded transition whitespace-nowrap">☎️ Calls</Link>
            <Link href="/dashboard/projects" className="px-3 py-1 bg-orange-50 hover:bg-orange-100 text-orange-700 text-xs rounded transition whitespace-nowrap">📌 Projects</Link>
            <Link href="/dashboard/users" className="px-3 py-1 bg-red-50 hover:bg-red-100 text-red-700 text-xs rounded transition whitespace-nowrap">👥 Users</Link>
            <Link href="/dashboard/reports" className="px-3 py-1 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 text-xs rounded transition whitespace-nowrap">📋 Reports</Link>
            <Link href="/dashboard/workflows" className="px-3 py-1 bg-cyan-50 hover:bg-cyan-100 text-cyan-700 text-xs rounded transition whitespace-nowrap">⚙️ Workflows</Link>
          </div>
        </div>

        {/* Content Area - Mobile First */}
        <main className="flex-1 overflow-auto p-4 md:p-6 bg-gray-50">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}
