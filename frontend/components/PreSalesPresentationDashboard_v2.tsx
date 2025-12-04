'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { presalesDashboardService } from '@/services/api'
import { HandshakeIcon, Target, Users, TrendingUp } from 'lucide-react'

export default function PreSalesPresentationDashboard() {
  // State for pre-sales data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [salesPipeline, setSalesPipeline] = useState<any>(null)
  const [opportunities, setOpportunities] = useState<any>(null)
  const [topDeals, setTopDeals] = useState<any>(null)
  const [conversionMetrics, setConversionMetrics] = useState<any>(null)

  // Fetch pre-sales data on mount
  useEffect(() => {
    const fetchPreSalesData = async () => {
      try {
        setLoading(true)

        // Fetch sales pipeline
        const pipelineRes = await presalesDashboardService.getSalesPipeline()
        setSalesPipeline(pipelineRes.data)

        // Fetch opportunities
        const opportunitiesRes = await presalesDashboardService.getOpportunities()
        setOpportunities(opportunitiesRes.data)

        // Fetch top deals
        const dealsRes = await presalesDashboardService.getTopDeals()
        setTopDeals(dealsRes.data)

        // Fetch conversion metrics
        const conversionRes = await presalesDashboardService.getConversionMetrics()
        setConversionMetrics(conversionRes.data)

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch pre-sales data:', err)
        setError(err.message || 'Failed to load pre-sales data')
      } finally {
        setLoading(false)
      }
    }

    fetchPreSalesData()
  }, [])

  // Use real data or fallback values
  const pipelineValue = salesPipeline?.pipeline_value || 4200000000
  const conversionRate = salesPipeline?.conversion_rate || 34
  const activeOpportunities = opportunities?.total || 127
  const expectedRevenue = salesPipeline?.expected_revenue || 1800000000

  const stages = salesPipeline?.stages || [
    { stage: 'Lead', count: 245, value: 8200000000, color: 'blue' },
    { stage: 'Qualified', count: 156, value: 5800000000, color: 'blue' },
    { stage: 'Proposal', count: 87, value: 4200000000, color: 'purple' },
    { stage: 'Negotiation', count: 34, value: 2200000000, color: 'orange' },
    { stage: 'Closing', count: 12, value: 800000000, color: 'green' }
  ]

  const deals = topDeals?.deals || [
    { name: 'Global Tech Corp - Enterprise Suite', value: 320000000, stage: 'Proposal', probability: 75, owner: 'Rahul Sharma', nextStep: 'Technical demo - Dec 10' },
    { name: 'Finance Plus - Integration Module', value: 280000000, stage: 'Negotiation', probability: 85, owner: 'Priya Nair', nextStep: 'Contract review - Dec 8' },
    { name: 'Manufacturing Ltd - MES System', value: 250000000, stage: 'Proposal', probability: 62, owner: 'Arjun Singh', nextStep: 'Commercial discussion - Dec 15' },
    { name: 'Healthcare Network - EMR Suite', value: 210000000, stage: 'Qualified', probability: 48, owner: 'Meera Kapoor', nextStep: 'Requirements workshop - Dec 12' }
  ]

  const teamMembers = [
    { name: 'Rahul Sharma', deals: 18, value: 850000000, achievement: 125 },
    { name: 'Priya Nair', deals: 15, value: 720000000, achievement: 118 },
    { name: 'Arjun Singh', deals: 12, value: 580000000, achievement: 94 },
    { name: 'Meera Kapoor', deals: 10, value: 450000000, achievement: 82 }
  ]

  const metrics = conversionMetrics || {
    lead_to_qualified: 64,
    qualified_to_proposal: 56,
    proposal_to_negotiation: 39,
    negotiation_to_close: 35
  }

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
      title: 'Pre-Sales & Opportunities',
      subtitle: 'Pipeline Health & Conversion Analysis',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <Target className="w-20 h-20 text-blue-600" />
          <div className="grid grid-cols-2 gap-6 w-full max-w-2xl">
            <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
              <div className="text-3xl font-bold text-blue-700">{formatCurrency(pipelineValue)}</div>
              <div className="text-sm text-gray-600 mt-1">Pipeline Value</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">{conversionRate}%</div>
              <div className="text-sm text-gray-600 mt-1">Conversion Rate</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">{activeOpportunities}</div>
              <div className="text-sm text-gray-600 mt-1">Active Opportunities</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">{formatCurrency(expectedRevenue)}</div>
              <div className="text-sm text-gray-600 mt-1">Expected Revenue</div>
            </div>
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'pipeline',
      title: 'Sales Pipeline Funnel',
      subtitle: 'Opportunity distribution across stages',
      content: (
        <div className="space-y-3 h-full">
          <div className="grid grid-cols-5 gap-2">
            {stages.map((stage: any, i: number) => {
              const maxCount = stages[0].count
              const width = (stage.count / maxCount) * 100
              return (
                <div key={i} className="text-center">
                  <div className={`bg-${stage.color}-100 border-2 border-${stage.color}-300 rounded-lg p-4 mb-2`}>
                    <div className="text-2xl font-bold text-gray-800">{stage.count}</div>
                    <div className="text-xs text-gray-600 mt-1">{formatCurrency(stage.value)}</div>
                  </div>
                  <div className="text-xs font-bold text-gray-700">{stage.stage}</div>
                  <div className="w-full h-1 bg-gray-200 rounded-full mt-2">
                    <div className={`h-full bg-${stage.color}-500`} style={{ width: `${width}%` }}></div>
                  </div>
                </div>
              )
            })}
          </div>
          <div className="bg-white p-4 rounded-lg border border-gray-200 mt-4">
            <h3 className="font-bold text-gray-800 mb-3">Conversion Metrics</h3>
            <div className="grid grid-cols-4 gap-4 text-center">
              <div>
                <div className="text-2xl font-bold text-blue-600">{metrics.lead_to_qualified}%</div>
                <div className="text-xs text-gray-600">Lead to Qualified</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-purple-600">{metrics.qualified_to_proposal}%</div>
                <div className="text-xs text-gray-600">Qualified to Proposal</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-orange-600">{metrics.proposal_to_negotiation}%</div>
                <div className="text-xs text-gray-600">Proposal to Negotiation</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-green-600">{metrics.negotiation_to_close}%</div>
                <div className="text-xs text-gray-600">Negotiation to Close</div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'top-deals',
      title: 'Top Opportunities',
      subtitle: 'High-value deals requiring immediate attention',
      content: (
        <div className="space-y-3 h-full overflow-y-auto">
          {deals.map((deal: any, i: number) => (
            <div key={i} className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="flex justify-between items-start mb-2">
                <div>
                  <h4 className="font-bold text-gray-800 text-sm">{deal.name}</h4>
                  <div className="text-xs text-gray-500 mt-1">Owner: {deal.owner}</div>
                </div>
                <div className="text-right">
                  <div className="text-lg font-bold text-blue-600">{formatCurrency(deal.value)}</div>
                  <div className={`text-xs px-2 py-1 rounded mt-1 font-bold ${
                    deal.probability >= 75 ? 'bg-green-100 text-green-800' :
                    deal.probability >= 50 ? 'bg-yellow-100 text-yellow-800' :
                    'bg-orange-100 text-orange-800'
                  }`}>{deal.probability}%</div>
                </div>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-0.5 rounded">{deal.stage}</span>
                <span className="text-xs text-gray-600">{deal.nextStep}</span>
              </div>
            </div>
          ))}
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'team-performance',
      title: 'Team Performance & Metrics',
      subtitle: 'Sales team productivity and targets',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Pre-Sales Executive Performance</h3>
            {teamMembers.map((person: any, i: number) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center mb-1">
                  <span className="font-semibold text-sm text-gray-800">{person.name}</span>
                  <span className={`text-xs px-2 py-0.5 rounded font-bold ${
                    person.achievement >= 100 ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                  }`}>{person.achievement}%</span>
                </div>
                <div className="text-xs text-gray-600">{person.deals} deals | {formatCurrency(person.value)}</div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Monthly Targets & Progress</h3>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm mb-3">
                <span className="font-semibold text-gray-800">Revenue Target</span>
                <span className="float-right font-bold text-blue-600">{formatCurrency(expectedRevenue)} / {formatCurrency(1600000000)}</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-green-500 h-full" style={{ width: '112%' }}></div>
              </div>
              <div className="text-xs text-gray-600 mt-1">On track: 112% of target</div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm mb-3">
                <span className="font-semibold text-gray-800">Deal Count</span>
                <span className="float-right font-bold text-green-600">55 / 45</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-blue-500 h-full" style={{ width: '122%' }}></div>
              </div>
              <div className="text-xs text-gray-600 mt-1">Excellent: 122% of target</div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm mb-3">
                <span className="font-semibold text-gray-800">Conversion Rate</span>
                <span className="float-right font-bold text-purple-600">{conversionRate}%</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-purple-500 h-full" style={{ width: `${conversionRate * 2}%` }}></div>
              </div>
              <div className="text-xs text-gray-600 mt-1">Target: 50% (industry avg: 28%)</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-green-50'
    },
    {
      id: 'industry-trends',
      title: 'Market Trends & Opportunities',
      subtitle: 'Market analysis and emerging segments',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Industry Segments - Demand</h3>
            {[
              { segment: 'Manufacturing', demand: 'Very High', growth: '+18%', pipeline: formatCurrency(1200000000) },
              { segment: 'Financial Services', demand: 'High', growth: '+14%', pipeline: formatCurrency(1000000000) },
              { segment: 'Healthcare', demand: 'High', growth: '+12%', pipeline: formatCurrency(800000000) },
              { segment: 'Retail/E-commerce', demand: 'Medium', growth: '+8%', pipeline: formatCurrency(600000000) }
            ].map((seg, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center">
                  <div>
                    <div className="font-semibold text-gray-800 text-sm">{seg.segment}</div>
                    <div className="text-xs text-gray-600">{seg.pipeline} in pipeline</div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-bold text-green-600">{seg.growth}</div>
                    <div className="text-xs text-gray-600">{seg.demand}</div>
                  </div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Competitive Landscape</h3>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">💡 Key Strength</div>
              <div className="text-sm text-gray-700 mt-1">Best-in-class implementation time. 40% faster than competitors. Huge market differentiator.</div>
            </div>
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ Market Opportunity</div>
              <div className="text-sm text-gray-700 mt-1">Competitors weak in customization. Position as "flexible ERP" targeting mid-market. TAM: ₹500 Cr.</div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">⚠️ Threat</div>
              <div className="text-sm text-gray-700 mt-1">Low-cost competitors entering market. Response: Bundled offerings, extended support packages.</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'summary',
      title: 'Pre-Sales Summary & Outlook',
      subtitle: 'Achievements and strategic initiatives',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ Strong Pipeline Health</div>
              <div className="text-sm text-gray-700 mt-2">{formatCurrency(pipelineValue)} pipeline with {conversionRate}% conversion rate. 3x coverage of monthly targets. Healthy pipeline distribution.</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 Team Excellence</div>
              <div className="text-sm text-gray-700 mt-2">125% of quarterly targets achieved. Avg deal value ₹3.3 Cr. Top performers incentivized for Q1.</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="font-bold text-purple-900">🎯 Market Position</div>
              <div className="text-sm text-gray-700 mt-2">Gaining momentum in manufacturing. Healthcare segment emerging. 3-5 year market opportunity: ₹1000+ Cr.</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">⚡ Next 30 Days Focus</div>
              <div className="text-sm text-gray-700 mt-2">
                • Close {formatCurrency(800000000)} deals in negotiation<br/>
                • 3 product demos scheduled<br/>
                • Healthcare vertical launch<br/>
                • Competitive positioning updates
              </div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">🚀 Strategic Initiatives</div>
              <div className="text-sm text-gray-700 mt-2">
                • Industry-specific solutions<br/>
                • Partner ecosystem development<br/>
                • Solution accelerators for segments<br/>
                • Thought leadership content
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    }
  ]

  return <PresentationDashboard slides={slides} title="Pre-Sales Dashboard" showSlideNumbers={true} />
}
