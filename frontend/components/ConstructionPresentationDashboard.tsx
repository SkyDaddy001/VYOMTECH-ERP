'use client'

import React from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'

function StatBox({ title, value, subtitle, icon }: { title: string; value: string; subtitle?: string; icon: string }) {
  return (
    <div className="bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-6 border border-gray-200">
      <div className="flex items-start justify-between mb-4">
        <div>
          <p className="text-sm font-semibold text-gray-600">{title}</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">{value}</p>
        </div>
        <span className="text-4xl">{icon}</span>
      </div>
      {subtitle && <p className="text-xs text-gray-600">{subtitle}</p>}
    </div>
  )
}

export default function ConstructionPresentationDashboard() {
  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Construction Dashboard',
      subtitle: 'Project Progress & BOQ Tracking',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <div className="text-center">
            <h3 className="text-5xl font-bold text-orange-600 mb-4">🏗️ CONSTRUCTION REPORT</h3>
            <p className="text-2xl text-gray-700 mb-2">December 2025</p>
            <p className="text-xl text-gray-600">Multi-Project Overview</p>
          </div>
          <div className="grid grid-cols-2 gap-6 mt-8 w-full max-w-2xl">
            <StatBox title="Active Projects" value="12" icon="🏢" />
            <StatBox title="Avg Completion" value="68%" icon="📊" />
          </div>
          <p className="text-gray-600 text-lg mt-8">Explore project details →</p>
        </div>
      )
    },
    {
      id: 'projects-overview',
      title: 'Projects Overview',
      subtitle: 'Current Status Snapshot',
      content: (
        <div className="space-y-3">
          {[
            { name: 'Metro Station Expansion', progress: 92, status: 'On Track' },
            { name: 'Commercial Complex - Phase 1', progress: 68, status: 'On Track' },
            { name: 'Residential Tower A', progress: 55, status: 'Delayed' },
            { name: 'Highway Overpass - Section 2', progress: 78, status: 'On Track' },
            { name: 'Shopping Mall - Foundation', progress: 32, status: 'On Track' }
          ].map((project, idx) => (
            <div key={idx} className="bg-white border border-gray-200 rounded-lg p-4">
              <div className="flex justify-between items-start mb-3">
                <div>
                  <p className="font-semibold text-gray-900">{project.name}</p>
                  <p className={`text-xs font-semibold ${project.status === 'On Track' ? 'text-green-600' : 'text-red-600'}`}>
                    {project.status}
                  </p>
                </div>
                <span className="text-2xl font-bold text-blue-600">{project.progress}%</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div
                  className={`h-full rounded-full transition-all ${project.status === 'On Track' ? 'bg-green-500' : 'bg-red-500'}`}
                  style={{ width: `${project.progress}%` }}
                ></div>
              </div>
            </div>
          ))}
        </div>
      )
    },
    {
      id: 'boq-summary',
      title: 'Bill of Quantities (BOQ) Summary',
      subtitle: 'Cost & Progress Analysis',
      content: (
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <StatBox title="Total BOQ Value" value="₹45.8 Cr" icon="💰" subtitle="All projects combined" />
            <StatBox title="Executed Value" value="₹31.2 Cr" icon="✓" subtitle="68% of total BOQ" />
            <StatBox title="Remaining Budget" value="₹14.6 Cr" icon="📈" subtitle="32% balance available" />
          </div>
          <div className="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-6 space-y-4">
            <h3 className="font-bold text-gray-900 text-lg">Budget vs Actual</h3>
            <div className="space-y-3">
              <div>
                <div className="flex justify-between mb-1">
                  <span className="text-sm font-semibold text-gray-700">Civil Works</span>
                  <span className="text-sm font-bold text-gray-900">₹18.5 Cr / ₹20 Cr</span>
                </div>
                <div className="w-full bg-gray-300 h-2 rounded-full overflow-hidden">
                  <div className="bg-green-500 h-full" style={{ width: '92.5%' }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between mb-1">
                  <span className="text-sm font-semibold text-gray-700">Material Cost</span>
                  <span className="text-sm font-bold text-gray-900">₹9.2 Cr / ₹15 Cr</span>
                </div>
                <div className="w-full bg-gray-300 h-2 rounded-full overflow-hidden">
                  <div className="bg-orange-500 h-full" style={{ width: '61.3%' }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between mb-1">
                  <span className="text-sm font-semibold text-gray-700">Labor Cost</span>
                  <span className="text-sm font-bold text-gray-900">₹3.5 Cr / ₹5 Cr</span>
                </div>
                <div className="w-full bg-gray-300 h-2 rounded-full overflow-hidden">
                  <div className="bg-yellow-500 h-full" style={{ width: '70%' }}></div>
                </div>
              </div>
              <div>
                <div className="flex justify-between mb-1">
                  <span className="text-sm font-semibold text-gray-700">Equipment</span>
                  <span className="text-sm font-bold text-gray-900">₹2.4 Cr / ₹5.8 Cr</span>
                </div>
                <div className="w-full bg-gray-300 h-2 rounded-full overflow-hidden">
                  <div className="bg-red-500 h-full" style={{ width: '41.4%' }}></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'timeline',
      title: 'Project Timeline',
      subtitle: 'Schedule & Milestones',
      content: (
        <div className="space-y-6">
          {[
            { title: 'Site Mobilization & Leveling', months: 2, completed: 2, color: 'bg-green-500' },
            { title: 'Foundation & Pile Works', months: 4, completed: 3.5, color: 'bg-green-500' },
            { title: 'Structural Steel/RCC Work', months: 6, completed: 3, color: 'bg-yellow-500' },
            { title: 'MEP Installation', months: 5, completed: 1, color: 'bg-orange-500' },
            { title: 'Interior Finishing', months: 4, completed: 0, color: 'bg-gray-300' },
            { title: 'Testing & Commissioning', months: 2, completed: 0, color: 'bg-gray-300' }
          ].map((milestone, idx) => (
            <div key={idx}>
              <p className="text-sm font-semibold text-gray-700 mb-2">{milestone.title}</p>
              <div className="flex gap-2 mb-1">
                {Array.from({ length: milestone.months }).map((_, i) => (
                  <div
                    key={i}
                    className={`flex-1 h-4 rounded ${
                      i < milestone.completed ? milestone.color : 'bg-gray-200'
                    }`}
                  />
                ))}
              </div>
              <p className="text-xs text-gray-600">
                {milestone.completed.toFixed(1)} of {milestone.months} months ({(milestone.completed / milestone.months * 100).toFixed(0)}%)
              </p>
            </div>
          ))}
        </div>
      )
    },
    {
      id: 'quality-safety',
      title: 'Quality & Safety',
      subtitle: 'Compliance & Performance Metrics',
      content: (
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-4">
            <div className="bg-gradient-to-br from-green-100 to-green-50 rounded-lg p-6 border-l-4 border-green-500">
              <p className="text-sm font-semibold text-gray-600 mb-1">Safety Record</p>
              <p className="text-4xl font-bold text-green-600">0</p>
              <p className="text-xs text-gray-600 mt-2">Days without accident</p>
            </div>
            <div className="bg-gradient-to-br from-blue-100 to-blue-50 rounded-lg p-6 border-l-4 border-blue-500">
              <p className="text-sm font-semibold text-gray-600 mb-1">Quality Score</p>
              <p className="text-4xl font-bold text-blue-600">94%</p>
              <p className="text-xs text-gray-600 mt-2">Inspection compliance</p>
            </div>
            <div className="bg-gradient-to-br from-purple-100 to-purple-50 rounded-lg p-6 border-l-4 border-purple-500">
              <p className="text-sm font-semibold text-gray-600 mb-1">Workforce</p>
              <p className="text-4xl font-bold text-purple-600">245</p>
              <p className="text-xs text-gray-600 mt-2">Active workers on-site</p>
            </div>
          </div>

          <div className="bg-gray-50 rounded-lg p-6 border border-gray-200">
            <h3 className="font-bold text-gray-900 mb-4">Quality Checklist</h3>
            <div className="space-y-3">
              {[
                { item: 'Concrete Strength Testing', status: true },
                { item: 'Structural Alignment', status: true },
                { item: 'Waterproofing Work', status: true },
                { item: 'Electrical Installation', status: false },
                { item: 'HVAC System', status: false }
              ].map((check, idx) => (
                <div key={idx} className="flex items-center gap-3">
                  <span className={`text-xl ${check.status ? '✓' : '○'}`}>
                    {check.status ? '✓' : '○'}
                  </span>
                  <span className={`text-sm ${check.status ? 'text-gray-700 line-through' : 'text-gray-700'}`}>
                    {check.item}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )
    },
    {
      id: 'risks-issues',
      title: 'Risks & Issues',
      subtitle: 'Current Challenges & Mitigation',
      content: (
        <div className="space-y-4">
          <div className="bg-red-50 border-l-4 border-red-500 p-6 rounded-lg">
            <p className="font-bold text-red-900 mb-2">🔴 High Priority - Supply Chain Delay</p>
            <p className="text-gray-700 text-sm mb-3">Steel reinforcement delivery delayed by 2 weeks, impacting structural work.</p>
            <p className="text-xs font-semibold text-red-700">Mitigation: Alternative supplier identified, expedited delivery arranged.</p>
          </div>

          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 rounded-lg">
            <p className="font-bold text-yellow-900 mb-2">🟡 Medium Priority - Weather Impact</p>
            <p className="text-gray-700 text-sm mb-3">Monsoon season approaching, may affect concrete curing and excavation work.</p>
            <p className="text-xs font-semibold text-yellow-700">Mitigation: Covered areas being prepared, contingency schedule ready.</p>
          </div>

          <div className="bg-blue-50 border-l-4 border-blue-500 p-6 rounded-lg">
            <p className="font-bold text-blue-900 mb-2">🔵 Low Priority - Permit Update</p>
            <p className="text-gray-700 text-sm mb-3">Municipal approval for MEP work stage in progress.</p>
            <p className="text-xs font-semibold text-blue-700">Status: Expected within 1 week, no impact on current schedule.</p>
          </div>
        </div>
      )
    },
    {
      id: 'summary',
      title: 'Executive Summary',
      subtitle: 'Key Takeaways & Next Steps',
      content: (
        <div className="space-y-4">
          <div className="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg">
            <p className="font-bold text-green-900 mb-2">✓ Overall Progress: 68%</p>
            <p className="text-gray-700 text-sm">Projects are progressing well with overall 68% completion. Most projects tracking on schedule.</p>
          </div>

          <div className="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg">
            <p className="font-bold text-green-900 mb-2">✓ Budget Status: On Track</p>
            <p className="text-gray-700 text-sm">Executed ₹31.2 Cr out of ₹45.8 Cr BOQ. Cost overrun controlled at -2.1%.</p>
          </div>

          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 rounded-lg">
            <p className="font-bold text-yellow-900 mb-2">⚠️ Residential Tower A - Watch List</p>
            <p className="text-gray-700 text-sm">Currently delayed by 2 weeks. Accelerated schedule being implemented.</p>
          </div>

          <div className="bg-blue-50 border-l-4 border-blue-500 p-6 rounded-lg">
            <p className="font-bold text-blue-900 mb-2">📅 Next Milestones (Next 30 Days)</p>
            <ul className="text-gray-700 text-sm space-y-1 ml-4 list-disc">
              <li>Complete structural work on 3 projects</li>
              <li>Commence MEP installation phase</li>
              <li>Resolve supply chain delays</li>
            </ul>
          </div>
        </div>
      )
    }
  ]

  return <PresentationDashboard slides={slides} title="Construction Project Dashboard" />
}
