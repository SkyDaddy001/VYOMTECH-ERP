'use client'

import { useState, useEffect } from 'react'
import { TenantInfo } from './TenantInfo'

interface StatCard {
  title: string
  value: string | number
  icon: string
  color: string
}

export default function DashboardContent() {
  const [stats, setStats] = useState<StatCard[]>([
    {
      title: 'Online Agents',
      value: 12,
      icon: '👥',
      color: 'bg-green-500',
    },
    {
      title: 'Calls Today',
      value: 342,
      icon: '📞',
      color: 'bg-blue-500',
    },
    {
      title: 'Avg Handle Time',
      value: '5m 32s',
      icon: '⏱️',
      color: 'bg-purple-500',
    },
    {
      title: 'Customer Satisfaction',
      value: '94%',
      icon: '😊',
      color: 'bg-yellow-500',
    },
    {
      title: 'Revenue Today',
      value: '$12,450',
      icon: '💰',
      color: 'bg-pink-500',
    },
    {
      title: 'Queue Length',
      value: 8,
      icon: '📋',
      color: 'bg-orange-500',
    },
  ])

  return (
    <div className="space-y-8">
      {/* Tenant Info */}
      <TenantInfo />

      {/* Welcome Section */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 rounded-lg shadow-lg p-8 text-white">
        <h1 className="text-4xl font-bold mb-2">Welcome to Call Center Dashboard</h1>
        <p className="text-blue-100">
          Monitor calls, manage agents, and track performance in real-time
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {stats.map((stat, index) => (
          <div key={index} className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm font-medium">{stat.title}</p>
                <p className="text-3xl font-bold text-gray-800 mt-2">{stat.value}</p>
              </div>
              <div className={`${stat.color} w-16 h-16 rounded-full flex items-center justify-center text-3xl`}>
                {stat.icon}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Calls */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-bold text-gray-800 mb-4">Recent Calls</h3>
          <div className="space-y-3">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-gray-800">Customer #{1000 + i}</p>
                  <p className="text-sm text-gray-500">5 min ago</p>
                </div>
                <span className="text-green-600 font-medium">Completed</span>
              </div>
            ))}
          </div>
        </div>

        {/* Active Agents */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-bold text-gray-800 mb-4">Active Agents</h3>
          <div className="space-y-3">
            {['John Doe', 'Jane Smith', 'Mike Johnson', 'Sarah Williams', 'Tom Brown'].map(
              (agent, i) => (
                <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center">
                    <div className="w-3 h-3 bg-green-500 rounded-full mr-3"></div>
                    <p className="font-medium text-gray-800">{agent}</p>
                  </div>
                  <span className="text-sm text-gray-500">Call #{Math.floor(Math.random() * 100)}</span>
                </div>
              )
            )}
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-bold text-gray-800 mb-4">Quick Actions</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button className="p-4 bg-blue-50 hover:bg-blue-100 rounded-lg text-center transition">
            <p className="text-2xl mb-2">➕</p>
            <p className="font-medium text-blue-900">New Agent</p>
          </button>
          <button className="p-4 bg-green-50 hover:bg-green-100 rounded-lg text-center transition">
            <p className="text-2xl mb-2">📞</p>
            <p className="font-medium text-green-900">Start Call</p>
          </button>
          <button className="p-4 bg-purple-50 hover:bg-purple-100 rounded-lg text-center transition">
            <p className="text-2xl mb-2">📋</p>
            <p className="font-medium text-purple-900">New Campaign</p>
          </button>
          <button className="p-4 bg-orange-50 hover:bg-orange-100 rounded-lg text-center transition">
            <p className="text-2xl mb-2">📊</p>
            <p className="font-medium text-orange-900">View Reports</p>
          </button>
        </div>
      </div>
    </div>
  )
}
