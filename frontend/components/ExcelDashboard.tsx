'use client';

/**
 * ExcelDashboard Component
 * Multi-tab, grid-based dashboard styled like Excel workbook
 * Displays spreadsheets with minimal chrome
 */

import React, { useState } from 'react';
import {
  Plus, FileText, ShoppingCart, Building2, TrendingUp, Download,
  Eye, EyeOff, Settings
} from 'lucide-react';

export interface DashboardTab {
  id: string;
  label: string;
  icon: React.ReactNode;
  content: React.ReactNode;
}

interface ExcelDashboardProps {
  tabs: DashboardTab[];
  title?: string;
  subtitle?: string;
  onAddSheet?: () => void;
}

const ExcelDashboard: React.FC<ExcelDashboardProps> = ({
  tabs,
  title = 'Dashboard',
  subtitle,
  onAddSheet,
}) => {
  const [activeTab, setActiveTab] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);

  const activeTabData = tabs[activeTab];

  return (
    <div className={`flex flex-col bg-white ${isFullscreen ? 'fixed inset-0 z-50' : 'rounded-lg shadow-sm border border-gray-200'}`}>
      {/* Workbook Header */}
      <div className="px-4 py-3 border-b border-gray-300 bg-gray-50">
        <div className="flex items-center justify-between mb-3">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{title}</h1>
            {subtitle && <p className="text-sm text-gray-600 mt-1">{subtitle}</p>}
          </div>
          <div className="flex gap-2">
            <button
              onClick={() => setIsFullscreen(!isFullscreen)}
              className="p-2 hover:bg-gray-200 rounded transition-colors"
              title={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
            >
              {isFullscreen ? '🔽' : '🔼'}
            </button>
          </div>
        </div>
      </div>

      {/* Workbook Tabs (like Excel sheets) */}
      <div className="flex border-b border-gray-300 bg-white px-4 overflow-x-auto">
        {tabs.map((tab, idx) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(idx)}
            className={`flex items-center gap-2 px-4 py-3 border-b-2 transition-colors whitespace-nowrap ${
              idx === activeTab
                ? 'border-blue-500 text-blue-600 font-semibold bg-blue-50'
                : 'border-transparent text-gray-600 hover:text-gray-900 hover:bg-gray-50'
            }`}
          >
            {tab.icon && <span className="text-lg">{tab.icon}</span>}
            {tab.label}
          </button>
        ))}

        {onAddSheet && (
          <button
            onClick={onAddSheet}
            className="flex items-center gap-2 px-4 py-3 border-b-2 border-transparent text-gray-600 hover:text-gray-900 hover:bg-gray-50 transition-colors"
            title="Add new sheet"
          >
            <Plus size={18} />
          </button>
        )}
      </div>

      {/* Tab Content */}
      <div className="flex-1 overflow-auto bg-white">
        {activeTabData && (
          <div className="p-4">
            {activeTabData.content}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="px-4 py-2 bg-gray-50 border-t border-gray-300 text-xs text-gray-600 flex items-center justify-between">
        <div>
          Tab: <span className="font-semibold">{activeTabData?.label}</span>
        </div>
        <div className="flex gap-4">
          <button className="p-1 hover:bg-gray-200 rounded" title="Export">
            <Download size={16} />
          </button>
          <button className="p-1 hover:bg-gray-200 rounded" title="Settings">
            <Settings size={16} />
          </button>
        </div>
      </div>
    </div>
  );
};

export default ExcelDashboard;
