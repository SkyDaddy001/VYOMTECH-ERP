'use client'

import { FiBell, FiUser, FiLogOut } from 'react-icons/fi'
import { useRouter } from 'next/navigation'

export const Header = () => {
  const router = useRouter()

  const handleLogout = () => {
    localStorage.removeItem('auth_token')
    router.push('/login')
  }

  return (
    <header className="fixed top-0 right-0 left-0 lg:left-64 z-40 bg-white border-b border-gray-200 h-16">
      <div className="flex items-center justify-between px-6 h-full">
        <div className="hidden lg:block">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-gray-900">Dashboard</h2>
        </div>

        <div className="flex items-center gap-2">
          <button className="p-2 hover:bg-gray-100 rounded-sm relative">
            <FiBell size={18} className="text-gray-700" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-gray-900 rounded-full"></span>
          </button>

          <button className="p-2 hover:bg-gray-100 rounded-sm">
            <FiUser size={18} className="text-gray-700" />
          </button>

          <button
            onClick={handleLogout}
            className="p-2 hover:bg-gray-100 rounded-sm text-gray-700 hover:text-gray-900"
          >
            <FiLogOut size={18} />
          </button>
        </div>
      </div>
    </header>
  )
}
