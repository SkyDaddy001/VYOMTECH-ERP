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
        className="fixed top-4 left-4 z-50 lg:hidden p-2 hover:bg-gray-100 rounded-lg"
      >
        {isOpen ? <FiX size={24} /> : <FiMenu size={24} />}
      </button>

      {/* Sidebar */}
      <aside
        className={`fixed left-0 top-0 z-40 w-64 h-screen bg-gradient-to-b from-gray-900 to-gray-800 text-white transition-transform duration-300 ease-in-out ${
          isOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
        }`}
      >
        <div className="p-6 border-b border-gray-700">
          <h1 className="text-2xl font-bold">VYOMTECH</h1>
          <p className="text-xs text-gray-400 mt-1">AI Call Center ERP</p>
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
                    ? 'bg-blue-600 border-l-4 border-blue-400'
                    : 'text-gray-300 hover:bg-gray-700'
                }`}
              >
                {item.label}
              </Link>
            )
          })}
        </nav>

        <div className="p-6 border-t border-gray-700">
          <p className="text-xs text-gray-400 mb-4">Logged in as</p>
          <p className="text-sm font-medium text-gray-100">Demo User</p>
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
