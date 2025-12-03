/**
 * Dashboard Sidebar Navigation
 * Spreadsheet-style navigation for easy access
 */

'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Menu, X, BarChart3, Settings } from 'lucide-react';

const dashboardMenuItems = [
  { label: 'Projects', href: '/dashboard/projects', icon: '📊' },
  { label: 'Sites', href: '/dashboard/sites', icon: '🏗️' },
  { label: 'BOQ', href: '/dashboard/boq', icon: '📋' },
  { label: 'Progress Tracking', href: '/dashboard/progress', icon: '📈' },
  { label: 'Safety Incidents', href: '/dashboard/safety', icon: '⚠️' },
  { label: 'Compliance', href: '/dashboard/compliance', icon: '✓' },
  { label: 'Equipment', href: '/dashboard/equipment', icon: '🔧' },
  { label: 'Permits', href: '/dashboard/permits', icon: '📄' },
];

interface DashboardLayoutProps {
  children: React.ReactNode;
}

const DashboardLayout: React.FC<DashboardLayoutProps> = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <div
        className={`${
          sidebarOpen ? 'w-64' : 'w-0'
        } bg-gray-900 text-white transition-all duration-300 overflow-hidden flex flex-col`}
      >
        {/* Logo */}
        <div className="px-6 py-4 border-b border-gray-800">
          <h1 className="text-xl font-bold">VYOMTECH ERP</h1>
          <p className="text-xs text-gray-400">Spreadsheet-Style Interface</p>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto px-3 py-4">
          {dashboardMenuItems.map(item => (
            <Link
              key={item.href}
              href={item.href}
              className="flex items-center gap-3 px-4 py-3 mb-2 rounded-lg hover:bg-gray-800 transition-colors"
            >
              <span className="text-xl">{item.icon}</span>
              <span className="text-sm font-medium">{item.label}</span>
            </Link>
          ))}
        </nav>

        {/* Footer */}
        <div className="px-3 py-4 border-t border-gray-800">
          <button className="flex items-center gap-2 px-4 py-2 w-full text-sm rounded-lg hover:bg-gray-800 transition-colors">
            <Settings size={16} />
            Settings
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Top Bar */}
        <div className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            {sidebarOpen ? <X size={20} /> : <Menu size={20} />}
          </button>

          <div className="flex items-center gap-4">
            <div className="text-sm text-gray-600">
              Logged in as: <strong>Admin User</strong>
            </div>
            <button className="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 text-sm">
              Logout
            </button>
          </div>
        </div>

        {/* Page Content */}
        <div className="flex-1 overflow-auto">
          {children}
        </div>
      </div>
    </div>
  );
};

export default DashboardLayout;
