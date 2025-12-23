'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { FiMenu, FiX } from 'react-icons/fi'
import { useState } from 'react'

const menuItems = [
  { label: 'Dashboard', href: '/dashboard' },
  { label: 'Leads', href: '/dashboard/leads' },
  { label: 'Campaigns', href: '/dashboard/campaigns' },
  { label: 'Calls', href: '/dashboard/calls' },
  { label: 'Agents', href: '/dashboard/agents' },
  { label: 'Reports', href: '/dashboard/reports' },
  { label: 'Settings', href: '/dashboard/settings' },
]

export const Sidebar = () => {
  const pathname = usePathname()
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      {/* Mobile menu button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed top-4 left-4 z-50 lg:hidden p-2 hover:bg-gray-100 rounded-sm"
      >
        {isOpen ? <FiX size={20} /> : <FiMenu size={20} />}
      </button>

      {/* Sidebar */}
      <aside
        className={`fixed left-0 top-0 z-40 w-64 h-screen bg-white border-r border-gray-200 transition-transform duration-300 ease-in-out ${
          isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        }`}
      >
        <div className="p-6 border-b border-gray-200">
          <h1 className="text-xl font-semibold text-gray-900">VYOMTECH</h1>
          <p className="text-xs text-gray-600 mt-1">AI Call Center ERP</p>
        </div>

        <nav className="flex-1 overflow-y-auto py-6">
          {menuItems.map((item) => {
            const isActive = pathname === item.href
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={() => setIsOpen(false)}
                className={`block px-6 py-3 text-sm font-medium transition ${
                  isActive
                    ? 'bg-gray-100 text-gray-900 border-l-2 border-gray-900'
                    : 'text-gray-700 hover:bg-gray-50'
                }`}
              >
                {item.label}
              </Link>
            )
          })}
        </nav>

        <div className="p-6 border-t border-gray-200">
          <p className="text-xs text-gray-600 mb-2 uppercase tracking-wide">Logged in as</p>
          <p className="text-sm font-medium text-gray-900">Demo User</p>
        </div>
      </aside>

      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-30 lg:hidden"
          onClick={() => setIsOpen(false)}
        ></div>
      )}
    </>
  )
}
