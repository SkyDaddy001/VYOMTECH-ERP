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
          <h2 className="text-lg font-semibold text-gray-900">Dashboard</h2>
        </div>

        <div className="flex items-center gap-4">
          <button className="p-2 hover:bg-gray-100 rounded-lg relative">
            <FiBell size={20} className="text-gray-600" />
            <span className="absolute top-0 right-0 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>

          <button className="p-2 hover:bg-gray-100 rounded-lg">
            <FiUser size={20} className="text-gray-600" />
          </button>

          <button
            onClick={handleLogout}
            className="p-2 hover:bg-gray-100 rounded-lg text-gray-600 hover:text-red-600"
          >
            <FiLogOut size={20} />
          </button>
        </div>
      </div>
    </header>
  )
}
