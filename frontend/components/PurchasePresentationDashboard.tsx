'use client'

import React from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { ShoppingBag, TrendingDown, Clock, AlertTriangle } from 'lucide-react'

export default function PurchasePresentationDashboard() {
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
              <div className="text-3xl font-bold text-blue-700">₹12.5 Cr</div>
              <div className="text-sm text-gray-600 mt-1">Total Purchase Value</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">47</div>
              <div className="text-sm text-gray-600 mt-1">Active Vendors</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">156</div>
              <div className="text-sm text-gray-600 mt-1">POs This Month</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">₹3.2 Cr</div>
              <div className="text-sm text-gray-600 mt-1">Pending Payments</div>
            </div>
          </div>
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
              <div className="text-4xl font-bold text-yellow-700 mt-2">23</div>
              <div className="text-xs text-gray-600 mt-1">Awaiting approval</div>
            </div>
            <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-6 rounded-lg border-2 border-blue-500">
              <div className="text-sm text-gray-600">Confirmed POs</div>
              <div className="text-4xl font-bold text-blue-700 mt-2">78</div>
              <div className="text-xs text-gray-600 mt-1">Sent to vendors</div>
            </div>
            <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-6 rounded-lg border-2 border-purple-500">
              <div className="text-sm text-gray-600">In Transit</div>
              <div className="text-4xl font-bold text-purple-700 mt-2">34</div>
              <div className="text-xs text-gray-600 mt-1">Expected delivery</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-gradient-to-br from-green-50 to-green-100 p-6 rounded-lg border-2 border-green-500">
              <div className="text-sm text-gray-600">Received</div>
              <div className="text-4xl font-bold text-green-700 mt-2">52</div>
              <div className="text-xs text-gray-600 mt-1">Quality verified</div>
            </div>
            <div className="bg-gradient-to-br from-red-50 to-red-100 p-6 rounded-lg border-2 border-red-500">
              <div className="text-sm text-gray-600">Delayed</div>
              <div className="text-4xl font-bold text-red-700 mt-2">8</div>
              <div className="text-xs text-gray-600 mt-1">Follow-up required</div>
            </div>
            <div className="bg-gradient-to-br from-gray-50 to-gray-100 p-6 rounded-lg border-2 border-gray-500">
              <div className="text-sm text-gray-600">Completed</div>
              <div className="text-4xl font-bold text-gray-700 mt-2">89</div>
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
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold">Steel Supplies Ltd</span>
                <span className="text-xl font-bold text-blue-600">₹2.8 Cr</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-blue-500 h-full" style={{ width: '100%' }}></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold">Cement Corp India</span>
                <span className="text-xl font-bold text-green-600">₹1.9 Cr</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-green-500 h-full" style={{ width: '68%' }}></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold">Electric Depot</span>
                <span className="text-xl font-bold text-purple-600">₹1.5 Cr</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-purple-500 h-full" style={{ width: '54%' }}></div>
              </div>
            </div>
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Vendor Reliability Score</h3>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold">Quality (95%)</span>
                <div className="flex gap-1">
                  {[...Array(5)].map((_, i) => (
                    <span key={i} className={i < 5 ? "text-yellow-400" : "text-gray-300"}>★</span>
                  ))}
                </div>
              </div>
              <div className="text-xs text-gray-600">Defect rate: 0.8%</div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold">On-Time Delivery (92%)</span>
                <div className="flex gap-1">
                  {[...Array(5)].map((_, i) => (
                    <span key={i} className={i < 5 ? "text-yellow-400" : "text-gray-300"}>★</span>
                  ))}
                </div>
              </div>
              <div className="text-xs text-gray-600">Avg delay: 1.2 days</div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="font-semibold">Price Competitiveness (88%)</span>
                <div className="flex gap-1">
                  {[...Array(5)].map((_, i) => (
                    <span key={i} className={i < 4 ? "text-yellow-400" : "text-gray-300"}>★</span>
                  ))}
                </div>
              </div>
              <div className="text-xs text-gray-600">1-2% above market avg</div>
            </div>
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
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm">Raw Materials</span>
                <span className="font-bold">₹6.2 Cr (49%)</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-blue-500 h-full" style={{ width: '100%' }}></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm">Equipment</span>
                <span className="font-bold">₹3.5 Cr (28%)</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-green-500 h-full" style={{ width: '57%' }}></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm">Services</span>
                <span className="font-bold">₹1.8 Cr (14%)</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-purple-500 h-full" style={{ width: '29%' }}></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm">Consumables</span>
                <span className="font-bold">₹1.0 Cr (8%)</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-orange-500 h-full" style={{ width: '16%' }}></div>
              </div>
            </div>
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Budget vs Actual</h3>
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="text-sm text-gray-600">Budget Status</div>
              <div className="text-2xl font-bold text-green-700 mt-2">Within Budget</div>
              <div className="text-xs text-gray-600 mt-1">₹12.5 Cr actual vs ₹13.0 Cr budgeted</div>
              <div className="mt-3">
                <div className="flex justify-between text-xs mb-1">
                  <span>Spend Rate</span>
                  <span>96%</span>
                </div>
                <div className="w-full bg-gray-200 h-2 rounded-full overflow-hidden">
                  <div className="bg-green-500 h-full" style={{ width: '96%' }}></div>
                </div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm font-bold text-gray-800 mb-3">Savings Achievement</div>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Vendor Discounts</span>
                  <span className="font-bold text-green-600">₹52 L</span>
                </div>
                <div className="flex justify-between">
                  <span>Negotiated Rates</span>
                  <span className="font-bold text-green-600">₹38 L</span>
                </div>
                <div className="flex justify-between">
                  <span>Volume Rebates</span>
                  <span className="font-bold text-green-600">₹25 L</span>
                </div>
                <hr className="my-2" />
                <div className="flex justify-between font-bold text-green-700">
                  <span>Total Savings</span>
                  <span>₹1.15 Cr</span>
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
              <div className="text-sm text-gray-700 mt-2">8 POs from primary suppliers delayed by 2-3 weeks. Impact on project timeline: ₹25 L potential loss.</div>
              <div className="mt-2 text-xs font-semibold text-orange-700">Action: Follow-up calls scheduled today</div>
            </div>
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">🔸 Medium: Price Increases</div>
              <div className="text-sm text-gray-700 mt-2">Steel prices up 8% month-on-month. Budget impact ₹15 L for Q2. Renegotiate volume contracts.</div>
              <div className="mt-2 text-xs font-semibold text-yellow-700">Action: Meeting with top 3 vendors scheduled</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <h3 className="font-bold text-gray-800 mb-3">Pending Payments</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Overdue (30+ days)</span>
                  <span className="font-bold text-red-600">₹48 L</span>
                </div>
                <div className="flex justify-between">
                  <span>Due Soon (1-30 days)</span>
                  <span className="font-bold text-orange-600">₹1.2 Cr</span>
                </div>
                <div className="flex justify-between">
                  <span>Normal Terms</span>
                  <span className="font-bold text-green-600">₹1.52 Cr</span>
                </div>
                <hr className="my-2" />
                <div className="flex justify-between font-bold text-gray-800">
                  <span>Total Payable</span>
                  <span>₹3.2 Cr</span>
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
              <div className="text-sm text-gray-700 mt-2">Achieved ₹1.15 Cr in savings through negotiation and volume discounts. On track for annual target of ₹5 Cr.</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 Strong Vendor Base</div>
              <div className="text-sm text-gray-700 mt-2">47 active vendors with average reliability score of 92%. Main 3 vendors represent 52% of spend.</div>
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
