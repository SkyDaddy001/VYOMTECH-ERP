'use client'

import { useState, useEffect } from 'react'
import axios from 'axios'

interface DashboardStats {
  totalVendors: number
  activeOrders: number
  pendingGRN: number
  totalContracts: number
  monthlySpend: number
  orderValue: number
}

export default function PurchaseDashboard() {
  const [stats, setStats] = useState<DashboardStats>({
    totalVendors: 0,
    activeOrders: 0,
    pendingGRN: 0,
    totalContracts: 0,
    monthlySpend: 0,
    orderValue: 0,
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      const response = await axios.get('/api/v1/purchase/dashboard/stats')
      setStats(response.data)
    } catch (error) {
      console.error('Failed to fetch stats:', error)
    } finally {
      setLoading(false)
    }
  }

  const StatCard = ({ title, value, icon, color }: { title: string; value: string | number; icon: string; color: string }) => (
    <div className={`${color} rounded-lg p-6 text-white shadow-lg`}>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm opacity-90">{title}</p>
          <p className="text-3xl font-bold mt-2">{value}</p>
        </div>
        <span className="text-5xl opacity-20">{icon}</span>
      </div>
    </div>
  )

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-gray-500">Loading dashboard...</div>
      </div>
    )
  }

  return (
    <div className="p-6 space-y-6">
      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <StatCard
          title="Total Vendors"
          value={stats.totalVendors}
          icon="🏢"
          color="bg-gradient-to-br from-blue-500 to-blue-600"
        />
        <StatCard
          title="Active Orders"
          value={stats.activeOrders}
          icon="📋"
          color="bg-gradient-to-br from-purple-500 to-purple-600"
        />
        <StatCard
          title="Pending GRN"
          value={stats.pendingGRN}
          icon="📦"
          color="bg-gradient-to-br from-orange-500 to-orange-600"
        />
        <StatCard
          title="Total Contracts"
          value={stats.totalContracts}
          icon="📄"
          color="bg-gradient-to-br from-green-500 to-green-600"
        />
        <StatCard
          title="Monthly Spend"
          value={`₹${stats.monthlySpend.toLocaleString('en-IN')}`}
          icon="💰"
          color="bg-gradient-to-br from-pink-500 to-pink-600"
        />
        <StatCard
          title="Avg Order Value"
          value={`₹${stats.orderValue.toLocaleString('en-IN')}`}
          icon="💳"
          color="bg-gradient-to-br from-indigo-500 to-indigo-600"
        />
      </div>

      {/* Recent Activity */}
      <div className="bg-gray-50 rounded-lg p-6">
        <h3 className="text-lg font-semibold mb-4">Recent Activity</h3>
        <div className="space-y-3">
          <p className="text-gray-600">Recent purchase orders and activities will appear here</p>
        </div>
      </div>
    </div>
  )
}
