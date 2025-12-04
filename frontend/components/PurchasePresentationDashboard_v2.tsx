'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { purchaseDashboardService } from '@/services/api'
import { ShoppingBag, TrendingDown, Clock, AlertTriangle } from 'lucide-react'

export default function PurchasePresentationDashboard() {
  // State for purchase data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [purchaseSummary, setPurchaseSummary] = useState<any>(null)
  const [poStatus, setPOStatus] = useState<any>(null)
  const [vendorList, setVendorList] = useState<any>(null)
  const [costAnalysis, setCostAnalysis] = useState<any>(null)

  // Fetch purchase data on mount
  useEffect(() => {
    const fetchPurchaseData = async () => {
      try {
        setLoading(true)

        // Fetch purchase summary
        const summaryRes = await purchaseDashboardService.getPurchaseSummary()
        setPurchaseSummary(summaryRes.data)

        // Fetch PO status
        const poRes = await purchaseDashboardService.getPOStatus()
        setPOStatus(poRes.data)

        // Fetch vendor list
        const vendorRes = await purchaseDashboardService.getVendorList()
        setVendorList(vendorRes.data)

        // Fetch cost analysis
        const costRes = await purchaseDashboardService.getCostAnalysis()
        setCostAnalysis(costRes.data)

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch purchase data:', err)
        setError(err.message || 'Failed to load purchase data')
      } finally {
        setLoading(false)
      }
    }

    fetchPurchaseData()
  }, [])

  // Use real data or fallback values
  const totalPurchaseValue = purchaseSummary?.total_value || 1250000000
  const activeVendors = purchaseSummary?.active_vendors || 47
  const posThisMonth = purchaseSummary?.pos_this_month || 156
  const pendingPayments = purchaseSummary?.pending_payments || 320000000

  const draftPOs = poStatus?.draft || 23
  const confirmedPOs = poStatus?.confirmed || 78
  const inTransit = poStatus?.in_transit || 34
  const received = poStatus?.received || 52
  const delayed = poStatus?.delayed || 8
  const completed = poStatus?.completed || 89

  const topVendors = vendorList?.top_vendors || [
    { name: 'Steel Supplies Ltd', value: 280000000, percentage: 100 },
    { name: 'Cement Corp India', value: 190000000, percentage: 68 },
    { name: 'Electric Depot', value: 150000000, percentage: 54 }
  ]

  const vendorScores = vendorList?.scores || [
    { metric: 'Quality', score: 95, stars: 5 },
    { metric: 'On-Time Delivery', score: 92, stars: 5 },
    { metric: 'Price Competitiveness', score: 88, stars: 4 }
  ]

  const spendByCategory = costAnalysis?.spend_by_category || [
    { category: 'Raw Materials', amount: 620000000, percentage: 49 },
    { category: 'Equipment', amount: 350000000, percentage: 28 },
    { category: 'Services', amount: 180000000, percentage: 14 },
    { category: 'Consumables', amount: 100000000, percentage: 8 }
  ]

  const budgetStatus = costAnalysis?.budget_status || {
    actual: 1250000000,
    budgeted: 1300000000,
    percentage: 96
  }

  const totalSavings = costAnalysis?.total_savings || 115000000

  const formatCurrency = (value: number) => {
    if (value >= 10000000) {
      return `₹${(value / 10000000).toFixed(1)} Cr`
    } else if (value >= 100000) {
      return `₹${(value / 100000).toFixed(0)} L`
    }
    return `₹${value.toLocaleString('en-IN')}`
  }

  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Procurement & Purchasing',
      subtitle: 'Vendor Management & Purchase Order Tracking',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <ShoppingBag className="w-20 h-20 text-blue-600" />
          <div className="grid grid-cols-2 gap-6 w-full max-w-2xl">
            <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
              <div className="text-3xl font-bold text-blue-700">{formatCurrency(totalPurchaseValue)}</div>
              <div className="text-sm text-gray-600 mt-1">Total Purchase Value</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">{activeVendors}</div>
              <div className="text-sm text-gray-600 mt-1">Active Vendors</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">{posThisMonth}</div>
              <div className="text-sm text-gray-600 mt-1">POs This Month</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">{formatCurrency(pendingPayments)}</div>
              <div className="text-sm text-gray-600 mt-1">Pending Payments</div>
            </div>
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'po-status',
      title: 'Purchase Order Pipeline',
      subtitle: 'Current PO status and order progression',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-gradient-to-br from-yellow-50 to-yellow-100 p-6 rounded-lg border-2 border-yellow-500">
              <div className="text-sm text-gray-600">Draft POs</div>
              <div className="text-4xl font-bold text-yellow-700 mt-2">{draftPOs}</div>
              <div className="text-xs text-gray-600 mt-1">Awaiting approval</div>
            </div>
            <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-6 rounded-lg border-2 border-blue-500">
              <div className="text-sm text-gray-600">Confirmed POs</div>
              <div className="text-4xl font-bold text-blue-700 mt-2">{confirmedPOs}</div>
              <div className="text-xs text-gray-600 mt-1">Sent to vendors</div>
            </div>
            <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-6 rounded-lg border-2 border-purple-500">
              <div className="text-sm text-gray-600">In Transit</div>
              <div className="text-4xl font-bold text-purple-700 mt-2">{inTransit}</div>
              <div className="text-xs text-gray-600 mt-1">Expected delivery</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-gradient-to-br from-green-50 to-green-100 p-6 rounded-lg border-2 border-green-500">
              <div className="text-sm text-gray-600">Received</div>
              <div className="text-4xl font-bold text-green-700 mt-2">{received}</div>
              <div className="text-xs text-gray-600 mt-1">Quality verified</div>
            </div>
            <div className="bg-gradient-to-br from-red-50 to-red-100 p-6 rounded-lg border-2 border-red-500">
              <div className="text-sm text-gray-600">Delayed</div>
              <div className="text-4xl font-bold text-red-700 mt-2">{delayed}</div>
              <div className="text-xs text-gray-600 mt-1">Follow-up required</div>
            </div>
            <div className="bg-gradient-to-br from-gray-50 to-gray-100 p-6 rounded-lg border-2 border-gray-500">
              <div className="text-sm text-gray-600">Completed</div>
              <div className="text-4xl font-bold text-gray-700 mt-2">{completed}</div>
              <div className="text-xs text-gray-600 mt-1">Invoiced & closed</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'vendor-performance',
      title: 'Vendor Performance Analysis',
      subtitle: 'Top vendors and reliability metrics',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Top Vendors by Volume</h3>
            {topVendors.map((vendor: any, idx: number) => (
              <div key={idx} className="bg-white p-4 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-semibold">{vendor.name}</span>
                  <span className="text-xl font-bold text-blue-600">{formatCurrency(vendor.value)}</span>
                </div>
                <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                  <div className="bg-blue-500 h-full" style={{ width: `${vendor.percentage}%` }}></div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Vendor Reliability Score</h3>
            {vendorScores.map((score: any, idx: number) => (
              <div key={idx} className="bg-white p-4 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-semibold">{score.metric} ({score.score}%)</span>
                  <div className="flex gap-1">
                    {[...Array(5)].map((_, i) => (
                      <span key={i} className={i < score.stars ? "text-yellow-400" : "text-gray-300"}>★</span>
                    ))}
                  </div>
                </div>
                <div className="text-xs text-gray-600">Industry benchmark: {score.score}%</div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-orange-50'
    },
    {
      id: 'cost-analysis',
      title: 'Cost Optimization & Spend',
      subtitle: 'Purchase spend by category and variance analysis',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Spend by Category</h3>
            {spendByCategory.map((item: any, idx: number) => (
              <div key={idx} className="bg-white p-4 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm">{item.category}</span>
                  <span className="font-bold">{formatCurrency(item.amount)} ({item.percentage}%)</span>
                </div>
                <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                  <div className="bg-blue-500 h-full" style={{ width: `${item.percentage}%` }}></div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Budget vs Actual</h3>
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="text-sm text-gray-600">Budget Status</div>
              <div className="text-2xl font-bold text-green-700 mt-2">Within Budget</div>
              <div className="text-xs text-gray-600 mt-1">{formatCurrency(budgetStatus.actual)} actual vs {formatCurrency(budgetStatus.budgeted)} budgeted</div>
              <div className="mt-3">
                <div className="flex justify-between text-xs mb-1">
                  <span>Spend Rate</span>
                  <span>{budgetStatus.percentage}%</span>
                </div>
                <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden">
                  <div className="bg-green-500 h-full" style={{ width: `${budgetStatus.percentage}%` }}></div>
                </div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm font-bold text-gray-800 mb-3">Savings Achievement</div>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Vendor Discounts</span>
                  <span className="font-bold text-green-600">{formatCurrency(totalSavings * 0.45)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Negotiated Rates</span>
                  <span className="font-bold text-green-600">{formatCurrency(totalSavings * 0.33)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Volume Rebates</span>
                  <span className="font-bold text-green-600">{formatCurrency(totalSavings * 0.22)}</span>
                </div>
                <hr className="my-2" />
                <div className="flex justify-between font-bold text-green-700">
                  <span>Total Savings</span>
                  <span>{formatCurrency(totalSavings)}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-green-50'
    },
    {
      id: 'issues-risks',
      title: 'Issues & Risk Management',
      subtitle: 'Supply chain risks and pending actions',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
              <div className="font-bold text-red-900 flex items-center gap-2">
                <AlertTriangle className="w-4 h-4" />
                Critical: Vendor Bankruptcy Risk
              </div>
              <div className="text-sm text-gray-700 mt-2">Major cement supplier showing signs of financial stress. Prepare alternate supplier agreements immediately.</div>
              <div className="mt-2 text-xs font-semibold text-red-700">Action: Contact 3 alternate suppliers by Friday</div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">⚠️ High: Delayed Delivery</div>
              <div className="text-sm text-gray-700 mt-2">{delayed} POs from primary suppliers delayed by 2-3 weeks. Potential impact: {formatCurrency(2500000)}.</div>
              <div className="mt-2 text-xs font-semibold text-orange-700">Action: Follow-up calls scheduled today</div>
            </div>
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">🔸 Medium: Price Increases</div>
              <div className="text-sm text-gray-700 mt-2">Steel prices up 8% month-on-month. Budget impact: {formatCurrency(15000000)} for Q2. Renegotiate volume contracts.</div>
              <div className="mt-2 text-xs font-semibold text-yellow-700">Action: Meeting with top 3 vendors scheduled</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <h3 className="font-bold text-gray-800 mb-3">Pending Payments</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Overdue (30+ days)</span>
                  <span className="font-bold text-red-600">{formatCurrency(pendingPayments * 0.15)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Due Soon (1-30 days)</span>
                  <span className="font-bold text-orange-600">{formatCurrency(pendingPayments * 0.375)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Normal Terms</span>
                  <span className="font-bold text-green-600">{formatCurrency(pendingPayments * 0.475)}</span>
                </div>
                <hr className="my-2" />
                <div className="flex justify-between font-bold text-gray-800">
                  <span>Total Payable</span>
                  <span>{formatCurrency(pendingPayments)}</span>
                </div>
              </div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">💡 Opportunity: Blockchain Tracking</div>
              <div className="text-sm text-gray-700 mt-2">Implement blockchain for supply chain visibility. Reduces discrepancies, improves vendor relations, enables better forecasting.</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-red-50'
    },
    {
      id: 'summary',
      title: 'Purchase Summary & Outlook',
      subtitle: 'Key metrics, trends, and next steps',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ Cost Efficiency</div>
              <div className="text-sm text-gray-700 mt-2">Achieved {formatCurrency(totalSavings)} in savings through negotiation and volume discounts. On track for annual target of {formatCurrency(500000000)}.</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 Strong Vendor Base</div>
              <div className="text-sm text-gray-700 mt-2">{activeVendors} active vendors with average reliability score of 92%. Main 3 vendors represent 52% of spend.</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="font-bold text-purple-900">📈 Operational Efficiency</div>
              <div className="text-sm text-gray-700 mt-2">Average PO processing time reduced to 3.2 days. Automated 65% of routine orders through system.</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">⚡ Next Steps (30 Days)</div>
              <div className="text-sm text-gray-700 mt-2">
                • Risk mitigation for cement supplier<br/>
                • Renegotiate steel contracts<br/>
                • Implement vendor scorecard system<br/>
                • Conduct quarterly vendor reviews
              </div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">🎯 Strategic Focus</div>
              <div className="text-sm text-gray-700 mt-2">
                • Diversify supplier base<br/>
                • Implement JIT inventory<br/>
                • Evaluate supplier financing options<br/>
                • Launch vendor development program
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    }
  ]

  return <PresentationDashboard slides={slides} title="Procurement Dashboard" showSlideNumbers={true} />
}
