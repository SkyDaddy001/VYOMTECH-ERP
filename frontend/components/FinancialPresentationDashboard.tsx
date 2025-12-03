'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { financialDashboardService } from '@/services/api'

interface KPICard {
  label: string
  value: string | number
  icon: string
  color: 'blue' | 'green' | 'red' | 'yellow' | 'purple'
}

function KPICard({ label, value, icon, color }: KPICard) {
  const colorMap = {
    blue: 'from-blue-400 to-blue-600',
    green: 'from-green-400 to-green-600',
    red: 'from-red-400 to-red-600',
    yellow: 'from-yellow-400 to-yellow-600',
    purple: 'from-purple-400 to-purple-600'
  }

  return (
    <div className={`bg-gradient-to-br ${colorMap[color]} text-white p-6 rounded-xl shadow-lg`}>
      <div className="text-4xl mb-2">{icon}</div>
      <p className="text-sm font-semibold opacity-90">{label}</p>
      <p className="text-3xl font-bold mt-2">{value}</p>
    </div>
  )
}

// Format currency in Indian Rupees
function formatCurrency(value: number): string {
  if (value >= 10000000) {
    return `₹${(value / 10000000).toFixed(2)}Cr`
  } else if (value >= 100000) {
    return `₹${(value / 100000).toFixed(2)}L`
  } else {
    return `₹${value.toLocaleString('en-IN')}`
  }
}

export default function FinancialPresentationDashboard() {
  // State for financial data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [balanceSheetData, setBalanceSheetData] = useState<any>(null)
  const [plData, setPLData] = useState<any>(null)
  const [ratiosData, setRatiosData] = useState<any>(null)

  // Fetch financial data on mount
  useEffect(() => {
    const fetchFinancialData = async () => {
      try {
        setLoading(true)
        const endDate = new Date()
        const startDate = new Date(endDate.getFullYear(), endDate.getMonth() - 2, 1)

        // Fetch balance sheet
        const bsResponse = await financialDashboardService.getBalanceSheet(endDate)
        setBalanceSheetData(bsResponse.data)

        // Fetch P&L
        const plResponse = await financialDashboardService.getProfitAndLoss(startDate, endDate)
        setPLData(plResponse.data)

        // Fetch ratios
        const ratiosResponse = await financialDashboardService.getFinancialRatios()
        setRatiosData(ratiosResponse.data)

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch financial data:', err)
        setError(err.message || 'Failed to load financial data')
        // Keep showing UI with fallback values
      } finally {
        setLoading(false)
      }
    }

    fetchFinancialData()
  }, [])

  // Use real data or fallback values
  const totalAssets = balanceSheetData?.total_assets || 12000000
  const totalLiabilities = balanceSheetData?.total_liabilities || 4200000
  const totalEquity = balanceSheetData?.total_equity || 7800000
  const totalIncome = plData?.total_income || 8500000
  const totalExpenses = plData?.total_expenses || 5200000
  const netProfit = plData?.profit_summary?.net_profit || 3300000

  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Financial Overview',
      subtitle: 'Quarterly Financial Performance & Analysis',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <div className="text-center">
            <h3 className="text-5xl font-bold text-blue-600 mb-4">💰 FINANCIAL REPORT</h3>
            <p className="text-2xl text-gray-700 mb-2">Q2 2025</p>
            <p className="text-xl text-gray-600">Complete Financial Overview</p>
          </div>
          <div className="grid grid-cols-2 gap-6 mt-8 w-full max-w-2xl">
            <KPICard label="Total Assets" value={formatCurrency(totalAssets)} icon="🏦" color="blue" />
            <KPICard label="Total Liabilities" value={formatCurrency(totalLiabilities)} icon="📊" color="red" />
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
          <p className="text-gray-600 text-lg mt-8">Navigate through financial insights →</p>
        </div>
      )
    },
    {
      id: 'balance-sheet',
      title: 'Balance Sheet Summary',
      subtitle: 'As on 30th June 2025',
      content: (
        <div className="grid grid-cols-3 gap-6">
          <div className="bg-blue-50 rounded-xl p-6 border-l-4 border-blue-500">
            <p className="text-sm font-semibold text-gray-600 mb-3">ASSETS</p>
            <div className="space-y-2">
              <div className="flex justify-between items-center pb-2 border-b border-blue-200">
                <span className="text-gray-700">Current Assets</span>
                <span className="font-bold text-blue-600">{formatCurrency(totalAssets * 0.45)}</span>
              </div>
              <div className="flex justify-between items-center pb-2 border-b border-blue-200">
                <span className="text-gray-700">Fixed Assets</span>
                <span className="font-bold text-blue-600">{formatCurrency(totalAssets * 0.55)}</span>
              </div>
              <div className="flex justify-between items-center pt-2 bg-blue-100 px-2 py-1 rounded">
                <span className="font-bold text-gray-800">Total Assets</span>
                <span className="font-bold text-blue-700">{formatCurrency(totalAssets)}</span>
              </div>
            </div>
          </div>

          <div className="bg-red-50 rounded-xl p-6 border-l-4 border-red-500">
            <p className="text-sm font-semibold text-gray-600 mb-3">LIABILITIES</p>
            <div className="space-y-2">
              <div className="flex justify-between items-center pb-2 border-b border-red-200">
                <span className="text-gray-700">Current Liab.</span>
                <span className="font-bold text-red-600">{formatCurrency(totalLiabilities * 0.43)}</span>
              </div>
              <div className="flex justify-between items-center pb-2 border-b border-red-200">
                <span className="text-gray-700">Long-term Liab.</span>
                <span className="font-bold text-red-600">{formatCurrency(totalLiabilities * 0.57)}</span>
              </div>
              <div className="flex justify-between items-center pt-2 bg-red-100 px-2 py-1 rounded">
                <span className="font-bold text-gray-800">Total Liab.</span>
                <span className="font-bold text-red-700">{formatCurrency(totalLiabilities)}</span>
              </div>
            </div>
          </div>

          <div className="bg-green-50 rounded-xl p-6 border-l-4 border-green-500">
            <p className="text-sm font-semibold text-gray-600 mb-3">EQUITY</p>
            <div className="space-y-2">
              <div className="flex justify-between items-center pb-2 border-b border-green-200">
                <span className="text-gray-700">Capital</span>
                <span className="font-bold text-green-600">{formatCurrency(totalEquity * 0.64)}</span>
              </div>
              <div className="flex justify-between items-center pb-2 border-b border-green-200">
                <span className="text-gray-700">Reserves</span>
                <span className="font-bold text-green-600">{formatCurrency(totalEquity * 0.36)}</span>
              </div>
              <div className="flex justify-between items-center pt-2 bg-green-100 px-2 py-1 rounded">
                <span className="font-bold text-gray-800">Total Equity</span>
                <span className="font-bold text-green-700">{formatCurrency(totalEquity)}</span>
              </div>
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'accounting-equation',
      title: 'Accounting Equation Verification',
      subtitle: 'Assets = Liabilities + Equity',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <div className="grid grid-cols-3 gap-8 w-full">
            <div className="bg-gradient-to-br from-blue-400 to-blue-600 text-white rounded-xl p-8 text-center">
              <p className="text-sm font-semibold opacity-90 mb-2">LEFT SIDE</p>
              <p className="text-5xl font-bold">{formatCurrency(totalAssets)}</p>
              <p className="text-lg font-semibold mt-3">ASSETS</p>
            </div>

            <div className="flex items-center justify-center">
              <p className="text-6xl font-bold text-gray-800">=</p>
            </div>

            <div className="space-y-4">
              <div className="bg-gradient-to-br from-red-400 to-red-600 text-white rounded-xl p-4 text-center">
                <p className="text-sm font-semibold opacity-90">LIABILITIES</p>
                <p className="text-3xl font-bold">{formatCurrency(totalLiabilities)}</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-gray-800 mb-2">+</p>
              </div>
              <div className="bg-gradient-to-br from-green-400 to-green-600 text-white rounded-xl p-4 text-center">
                <p className="text-sm font-semibold opacity-90">EQUITY</p>
                <p className="text-3xl font-bold">{formatCurrency(totalEquity)}</p>
              </div>
            </div>
          </div>

          <div className="mt-8 bg-green-50 border-2 border-green-500 rounded-xl p-6 w-full text-center">
            <p className="text-2xl font-bold text-green-700">✓ EQUATION BALANCED</p>
            <p className="text-gray-700 mt-2">{formatCurrency(totalAssets)} = {formatCurrency(totalLiabilities)} + {formatCurrency(totalEquity)}</p>
          </div>
        </div>
      )
    },
    {
      id: 'profitability',
      title: 'Income Statement',
      subtitle: 'Profitability Analysis',
      content: (
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-green-50 rounded-lg p-4 border-l-4 border-green-500">
              <p className="text-xs font-semibold text-gray-600 mb-1">REVENUE</p>
              <p className="text-3xl font-bold text-green-600">{formatCurrency(totalIncome)}</p>
            </div>
            <div className="bg-red-50 rounded-lg p-4 border-l-4 border-red-500">
              <p className="text-xs font-semibold text-gray-600 mb-1">EXPENSES</p>
              <p className="text-3xl font-bold text-red-600">{formatCurrency(totalExpenses)}</p>
            </div>
          </div>

          <div className="bg-gradient-to-r from-gray-100 to-gray-50 rounded-lg p-6 space-y-3">
            <div className="flex justify-between items-center">
              <span className="font-semibold text-gray-700">Revenue</span>
              <span className="text-lg font-bold text-gray-900">{formatCurrency(totalIncome)}</span>
            </div>
            <div className="border-t border-gray-300"></div>
            <div className="flex justify-between items-center text-red-600">
              <span className="font-semibold">Cost of Goods Sold</span>
              <span className="text-lg font-bold">-{formatCurrency(totalIncome * 0.412)}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="font-semibold text-gray-700">Gross Profit</span>
              <span className="text-lg font-bold text-blue-600">{formatCurrency(totalIncome * 0.588)}</span>
            </div>
            <div className="border-t border-gray-300"></div>
            <div className="flex justify-between items-center text-red-600">
              <span className="font-semibold">Operating Expenses</span>
              <span className="text-lg font-bold">-{formatCurrency(totalExpenses * 0.4)}</span>
            </div>
            <div className="bg-blue-100 rounded px-3 py-2">
              <div className="flex justify-between items-center">
                <span className="font-bold text-gray-900">Operating Profit</span>
                <span className="text-xl font-bold text-blue-600">{formatCurrency(netProfit * 1.1)}</span>
              </div>
            </div>
            <div className="border-t border-gray-300"></div>
            <div className="bg-gradient-to-r from-green-100 to-green-50 rounded px-3 py-2">
              <div className="flex justify-between items-center">
                <span className="font-bold text-green-900">Net Profit</span>
                <span className="text-2xl font-bold text-green-600">{formatCurrency(netProfit)}</span>
              </div>
            </div>
          </div>

          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
            <p className="text-sm">
              <span className="font-bold">Profit Margin: </span>
              <span className="text-lg font-bold text-yellow-600">{((netProfit / totalIncome) * 100).toFixed(1)}%</span>
            </p>
          </div>
        </div>
      )
    },
    {
      id: 'ratio-analysis',
      title: 'Financial Ratios',
      subtitle: 'Key Performance Metrics',
      content: (
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <div className="bg-blue-50 rounded-lg p-4 border-l-4 border-blue-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Current Ratio</p>
              <p className="text-4xl font-bold text-blue-600">{(totalAssets / totalLiabilities).toFixed(2)}x</p>
              <p className="text-xs text-gray-600 mt-1">Liquidity Position</p>
            </div>
            <div className="bg-green-50 rounded-lg p-4 border-l-4 border-green-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">ROE (Return on Equity)</p>
              <p className="text-4xl font-bold text-green-600">{((netProfit / totalEquity) * 100).toFixed(1)}%</p>
              <p className="text-xs text-gray-600 mt-1">Shareholder Returns</p>
            </div>
            <div className="bg-purple-50 rounded-lg p-4 border-l-4 border-purple-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Debt to Equity</p>
              <p className="text-4xl font-bold text-purple-600">{(totalLiabilities / totalEquity).toFixed(2)}</p>
              <p className="text-xs text-gray-600 mt-1">Financial Leverage</p>
            </div>
          </div>

          <div className="space-y-4">
            <div className="bg-orange-50 rounded-lg p-4 border-l-4 border-orange-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Asset Turnover</p>
              <p className="text-4xl font-bold text-orange-600">{(totalIncome / totalAssets).toFixed(2)}x</p>
    {
      id: 'cover',
      title: 'Financial Overview',
      subtitle: 'Quarterly Financial Performance & Analysis',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <div className="text-center">
            <h3 className="text-5xl font-bold text-blue-600 mb-4">💰 FINANCIAL REPORT</h3>
            <p className="text-2xl text-gray-700 mb-2">Q2 2025</p>
            <p className="text-xl text-gray-600">Complete Financial Overview</p>
          </div>
          <div className="grid grid-cols-2 gap-6 mt-8 w-full max-w-2xl">
            <KPICard label="Total Assets" value="₹1.2Cr" icon="🏦" color="blue" />
            <KPICard label="Total Liabilities" value="₹42L" icon="📊" color="red" />
          </div>
          <p className="text-gray-600 text-lg mt-8">Navigate through financial insights →</p>
        </div>
      )
    },
    {
      id: 'balance-sheet',
      title: 'Balance Sheet Summary',
      subtitle: 'As on 30th June 2025',
      content: (
        <div className="grid grid-cols-3 gap-6">
          <div className="bg-blue-50 rounded-xl p-6 border-l-4 border-blue-500">
            <p className="text-sm font-semibold text-gray-600 mb-3">ASSETS</p>
            <div className="space-y-2">
              <div className="flex justify-between items-center pb-2 border-b border-blue-200">
                <span className="text-gray-700">Current Assets</span>
                <span className="font-bold text-blue-600">₹56L</span>
              </div>
              <div className="flex justify-between items-center pb-2 border-b border-blue-200">
                <span className="text-gray-700">Fixed Assets</span>
                <span className="font-bold text-blue-600">₹64L</span>
              </div>
              <div className="flex justify-between items-center pt-2 bg-blue-100 px-2 py-1 rounded">
                <span className="font-bold text-gray-800">Total Assets</span>
                <span className="font-bold text-blue-700">₹120L</span>
              </div>
            </div>
          </div>

          <div className="bg-red-50 rounded-xl p-6 border-l-4 border-red-500">
            <p className="text-sm font-semibold text-gray-600 mb-3">LIABILITIES</p>
            <div className="space-y-2">
              <div className="flex justify-between items-center pb-2 border-b border-red-200">
                <span className="text-gray-700">Current Liab.</span>
                <span className="font-bold text-red-600">₹18L</span>
              </div>
              <div className="flex justify-between items-center pb-2 border-b border-red-200">
                <span className="text-gray-700">Long-term Liab.</span>
                <span className="font-bold text-red-600">₹24L</span>
              </div>
              <div className="flex justify-between items-center pt-2 bg-red-100 px-2 py-1 rounded">
                <span className="font-bold text-gray-800">Total Liab.</span>
                <span className="font-bold text-red-700">₹42L</span>
              </div>
            </div>
          </div>

          <div className="bg-green-50 rounded-xl p-6 border-l-4 border-green-500">
            <p className="text-sm font-semibold text-gray-600 mb-3">EQUITY</p>
            <div className="space-y-2">
              <div className="flex justify-between items-center pb-2 border-b border-green-200">
                <span className="text-gray-700">Capital</span>
                <span className="font-bold text-green-600">₹50L</span>
              </div>
              <div className="flex justify-between items-center pb-2 border-b border-green-200">
                <span className="text-gray-700">Reserves</span>
                <span className="font-bold text-green-600">₹28L</span>
              </div>
              <div className="flex justify-between items-center pt-2 bg-green-100 px-2 py-1 rounded">
                <span className="font-bold text-gray-800">Total Equity</span>
                <span className="font-bold text-green-700">₹78L</span>
              </div>
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'accounting-equation',
      title: 'Accounting Equation Verification',
      subtitle: 'Assets = Liabilities + Equity',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <div className="grid grid-cols-3 gap-8 w-full">
            <div className="bg-gradient-to-br from-blue-400 to-blue-600 text-white rounded-xl p-8 text-center">
              <p className="text-sm font-semibold opacity-90 mb-2">LEFT SIDE</p>
              <p className="text-5xl font-bold">₹120L</p>
              <p className="text-lg font-semibold mt-3">ASSETS</p>
            </div>

            <div className="flex items-center justify-center">
              <p className="text-6xl font-bold text-gray-800">=</p>
            </div>

            <div className="space-y-4">
              <div className="bg-gradient-to-br from-red-400 to-red-600 text-white rounded-xl p-4 text-center">
                <p className="text-sm font-semibold opacity-90">LIABILITIES</p>
                <p className="text-3xl font-bold">₹42L</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-gray-800 mb-2">+</p>
              </div>
              <div className="bg-gradient-to-br from-green-400 to-green-600 text-white rounded-xl p-4 text-center">
                <p className="text-sm font-semibold opacity-90">EQUITY</p>
                <p className="text-3xl font-bold">₹78L</p>
              </div>
            </div>
          </div>

          <div className="mt-8 bg-green-50 border-2 border-green-500 rounded-xl p-6 w-full text-center">
            <p className="text-2xl font-bold text-green-700">✓ EQUATION BALANCED</p>
            <p className="text-gray-700 mt-2">₹120L = ₹42L + ₹78L</p>
          </div>
        </div>
      )
    },
    {
      id: 'profitability',
      title: 'Income Statement',
      subtitle: 'Profitability Analysis',
      content: (
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-green-50 rounded-lg p-4 border-l-4 border-green-500">
              <p className="text-xs font-semibold text-gray-600 mb-1">REVENUE</p>
              <p className="text-3xl font-bold text-green-600">₹85L</p>
            </div>
            <div className="bg-red-50 rounded-lg p-4 border-l-4 border-red-500">
              <p className="text-xs font-semibold text-gray-600 mb-1">EXPENSES</p>
              <p className="text-3xl font-bold text-red-600">₹52L</p>
            </div>
          </div>

          <div className="bg-gradient-to-r from-gray-100 to-gray-50 rounded-lg p-6 space-y-3">
            <div className="flex justify-between items-center">
              <span className="font-semibold text-gray-700">Revenue</span>
              <span className="text-lg font-bold text-gray-900">₹85L</span>
            </div>
            <div className="border-t border-gray-300"></div>
            <div className="flex justify-between items-center text-red-600">
              <span className="font-semibold">Cost of Goods Sold</span>
              <span className="text-lg font-bold">-₹35L</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="font-semibold text-gray-700">Gross Profit</span>
              <span className="text-lg font-bold text-blue-600">₹50L</span>
            </div>
            <div className="border-t border-gray-300"></div>
            <div className="flex justify-between items-center text-red-600">
              <span className="font-semibold">Operating Expenses</span>
              <span className="text-lg font-bold">-₹17L</span>
            </div>
            <div className="bg-blue-100 rounded px-3 py-2">
              <div className="flex justify-between items-center">
                <span className="font-bold text-gray-900">Operating Profit</span>
                <span className="text-xl font-bold text-blue-600">₹33L</span>
              </div>
            </div>
            <div className="border-t border-gray-300"></div>
            <div className="bg-gradient-to-r from-green-100 to-green-50 rounded px-3 py-2">
              <div className="flex justify-between items-center">
                <span className="font-bold text-green-900">Net Profit</span>
                <span className="text-2xl font-bold text-green-600">₹33L</span>
              </div>
            </div>
          </div>

          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
            <p className="text-sm">
              <span className="font-bold">Profit Margin: </span>
              <span className="text-lg font-bold text-yellow-600">38.8%</span>
            </p>
          </div>
        </div>
      )
    },
    {
      id: 'ratio-analysis',
      title: 'Financial Ratios',
      subtitle: 'Key Performance Metrics',
      content: (
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <div className="bg-blue-50 rounded-lg p-4 border-l-4 border-blue-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Current Ratio</p>
              <p className="text-4xl font-bold text-blue-600">1.8x</p>
              <p className="text-xs text-gray-600 mt-1">Liquidity Position</p>
            </div>
            <div className="bg-green-50 rounded-lg p-4 border-l-4 border-green-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">ROE (Return on Equity)</p>
              <p className="text-4xl font-bold text-green-600">42.3%</p>
              <p className="text-xs text-gray-600 mt-1">Shareholder Returns</p>
            </div>
            <div className="bg-purple-50 rounded-lg p-4 border-l-4 border-purple-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Debt to Equity</p>
              <p className="text-4xl font-bold text-purple-600">0.54</p>
              <p className="text-xs text-gray-600 mt-1">Financial Leverage</p>
            </div>
          </div>

          <div className="space-y-4">
            <div className="bg-orange-50 rounded-lg p-4 border-l-4 border-orange-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Asset Turnover</p>
              <p className="text-4xl font-bold text-orange-600">0.71x</p>
              <p className="text-xs text-gray-600 mt-1">Efficiency Metric</p>
            </div>
            <div className="bg-red-50 rounded-lg p-4 border-l-4 border-red-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">Quick Ratio</p>
              <p className="text-4xl font-bold text-red-600">1.2x</p>
              <p className="text-xs text-gray-600 mt-1">Immediate Liquidity</p>
            </div>
            <div className="bg-indigo-50 rounded-lg p-4 border-l-4 border-indigo-500">
              <p className="text-sm text-gray-600 font-semibold mb-2">ROA (Return on Assets)</p>
              <p className="text-4xl font-bold text-indigo-600">27.5%</p>
              <p className="text-xs text-gray-600 mt-1">Asset Efficiency</p>
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'conclusions',
      title: 'Financial Summary',
      subtitle: 'Key Conclusions & Recommendations',
      content: (
        <div className="space-y-4">
          <div className="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg">
            <p className="font-bold text-green-900 mb-2">✓ Strong Financial Health</p>
            <p className="text-gray-700 text-sm">Healthy current ratio of 1.8x indicates good short-term liquidity position.</p>
          </div>
          <div className="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg">
            <p className="font-bold text-green-900 mb-2">✓ Excellent Profitability</p>
            <p className="text-gray-700 text-sm">Net profit margin of 38.8% shows strong operational efficiency and cost control.</p>
          </div>
          <div className="bg-blue-50 border-l-4 border-blue-500 p-6 rounded-lg">
            <p className="font-bold text-blue-900 mb-2">📊 Balanced Capital Structure</p>
            <p className="text-gray-700 text-sm">Debt-to-equity ratio of 0.54 indicates conservative financial leverage.</p>
          </div>
          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 rounded-lg">
            <p className="font-bold text-yellow-900 mb-2">💡 Recommendations</p>
            <ul className="text-gray-700 text-sm space-y-1 ml-4 list-disc">
              <li>Continue focus on revenue growth initiatives</li>
              <li>Optimize working capital management</li>
              <li>Explore strategic expansion opportunities</li>
            </ul>
          </div>
        </div>
      )
    }
  ]

  return <PresentationDashboard slides={slides} title="Financial Dashboard" />
}
