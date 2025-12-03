'use client'

import React from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { Briefcase, TrendingUp, Clock, Users } from 'lucide-react'

export default function ProjectsPresentationDashboard() {
  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Project Management',
      subtitle: 'Portfolio Overview & Delivery Status',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <Briefcase className="w-20 h-20 text-blue-600" />
          <div className="grid grid-cols-2 gap-6 w-full max-w-2xl">
            <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
              <div className="text-3xl font-bold text-blue-700">18</div>
              <div className="text-sm text-gray-600 mt-1">Active Projects</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">₹85 Cr</div>
              <div className="text-sm text-gray-600 mt-1">Total Value</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">62%</div>
              <div className="text-sm text-gray-600 mt-1">Avg Completion</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">245</div>
              <div className="text-sm text-gray-600 mt-1">Team Members</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'projects',
      title: 'Project Portfolio Status',
      subtitle: 'Current projects and their completion status',
      content: (
        <div className="space-y-3 h-full overflow-y-auto">
          <div className="grid grid-cols-2 gap-4">
            {[
              { name: 'Enterprise Portal', progress: 85, status: 'On Track', team: 12 },
              { name: 'Mobile App v2', progress: 72, status: 'On Track', team: 8 },
              { name: 'Analytics Platform', progress: 58, status: 'On Track', team: 6 },
              { name: 'Legacy Migration', progress: 45, status: 'Behind', team: 10 },
              { name: 'API Modernization', progress: 68, status: 'On Track', team: 7 },
              { name: 'Security Audit', progress: 91, status: 'Completing', team: 5 }
            ].map((proj, i) => (
              <div key={i} className="bg-white p-4 rounded-lg border border-gray-200">
                <div className="flex justify-between items-start mb-2">
                  <div>
                    <h3 className="font-bold text-gray-800 text-sm">{proj.name}</h3>
                    <div className="text-xs text-gray-500">{proj.team} team members</div>
                  </div>
                  <span className={`text-xs px-2 py-1 rounded font-bold ${
                    proj.status === 'On Track' ? 'bg-green-100 text-green-800' :
                    proj.status === 'Behind' ? 'bg-red-100 text-red-800' :
                    'bg-blue-100 text-blue-800'
                  }`}>{proj.status}</span>
                </div>
                <div className="relative w-full bg-gray-200 h-2.5 rounded-full overflow-hidden">
                  <div className={`h-full ${
                    proj.progress >= 80 ? 'bg-green-500' :
                    proj.progress >= 50 ? 'bg-blue-500' :
                    'bg-yellow-500'
                  }`} style={{ width: `${proj.progress}%` }}></div>
                </div>
                <div className="text-xs text-gray-600 mt-1 text-right">{proj.progress}%</div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'timeline',
      title: 'Delivery Timeline',
      subtitle: 'Upcoming milestones and deliverables',
      content: (
        <div className="space-y-4 h-full">
          <div className="grid grid-cols-4 gap-3">
            {[
              { month: 'Dec 2024', status: 'Complete', count: 8 },
              { month: 'Jan 2025', status: 'In Progress', count: 12 },
              { month: 'Feb 2025', status: 'Planned', count: 15 },
              { month: 'Mar 2025', status: 'Planned', count: 10 }
            ].map((m, i) => (
              <div key={i} className={`p-3 rounded-lg border ${
                m.status === 'Complete' ? 'bg-green-50 border-green-300' :
                m.status === 'In Progress' ? 'bg-blue-50 border-blue-300' :
                'bg-gray-50 border-gray-300'
              }`}>
                <div className="text-xs font-bold text-gray-700">{m.month}</div>
                <div className="text-2xl font-bold mt-2">{m.count}</div>
                <div className="text-xs text-gray-600 mt-1">{m.status}</div>
              </div>
            ))}
          </div>
          <div className="space-y-2">
            <h3 className="font-bold text-gray-800">Critical Milestones (Next 90 Days)</h3>
            {[
              { name: 'Enterprise Portal - UAT Complete', date: 'Dec 15', status: '85%' },
              { name: 'Mobile App v2 - Beta Release', date: 'Dec 22', status: '72%' },
              { name: 'Analytics Dashboard - Phase 1', date: 'Jan 10', status: '58%' },
              { name: 'Legacy System Cutover', date: 'Jan 28', status: '45%' }
            ].map((m, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200 flex justify-between items-center">
                <div>
                  <div className="font-semibold text-gray-800 text-sm">{m.name}</div>
                  <div className="text-xs text-gray-500">{m.date}</div>
                </div>
                <div className="text-sm font-bold text-blue-600">{m.status}</div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'budget',
      title: 'Budget & Resource Allocation',
      subtitle: 'Project budgets and resource utilization',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Budget Status by Project</h3>
            {[
              { name: 'Enterprise Portal', budget: '₹12 Cr', spent: '₹10.2 Cr', variance: '-15%', status: 'green' },
              { name: 'Mobile App v2', budget: '₹8 Cr', spent: '₹5.8 Cr', variance: '-27%', status: 'green' },
              { name: 'Analytics Platform', budget: '₹6.5 Cr', spent: '₹3.8 Cr', variance: '-42%', status: 'green' },
              { name: 'Legacy Migration', budget: '₹9 Cr', spent: '₹5.2 Cr', variance: '+8%', status: 'red' }
            ].map((p, i) => (
              <div key={i} className="bg-white p-3 rounded-lg border border-gray-200">
                <div className="flex justify-between items-start mb-1">
                  <span className="font-semibold text-sm text-gray-800">{p.name}</span>
                  <span className={`text-xs px-2 py-0.5 rounded font-bold ${p.status === 'green' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>{p.variance}</span>
                </div>
                <div className="text-xs text-gray-600">{p.spent} / {p.budget}</div>
              </div>
            ))}
          </div>
          <div className="space-y-3">
            <h3 className="font-bold text-gray-800 mb-2">Resource Utilization</h3>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm mb-3">
                <span className="font-semibold text-gray-800">Total Team Capacity</span>
                <span className="float-right text-blue-600 font-bold">245 / 300</span>
              </div>
              <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                <div className="bg-blue-500 h-full" style={{ width: '82%' }}></div>
              </div>
              <div className="text-xs text-gray-600 mt-1">Utilization: 82%</div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <h4 className="font-semibold text-gray-800 text-sm mb-3">Skills Allocation</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Backend Engineers</span>
                  <span className="font-bold">78 (32%)</span>
                </div>
                <div className="flex justify-between">
                  <span>Frontend Engineers</span>
                  <span className="font-bold">65 (27%)</span>
                </div>
                <div className="flex justify-between">
                  <span>QA Engineers</span>
                  <span className="font-bold">52 (21%)</span>
                </div>
                <div className="flex justify-between">
                  <span>DevOps/Infrastructure</span>
                  <span className="font-bold">28 (11%)</span>
                </div>
                <div className="flex justify-between">
                  <span>Product Managers</span>
                  <span className="font-bold">22 (9%)</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-green-50'
    },
    {
      id: 'risks',
      title: 'Risks & Issues Tracker',
      subtitle: 'Critical blockers and mitigation plans',
      content: (
        <div className="space-y-3 h-full overflow-y-auto">
          <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
            <div className="font-bold text-red-900 flex items-center gap-2">
              <span className="text-lg">🔴 CRITICAL</span>
            </div>
            <div className="text-sm text-gray-700 mt-1">Legacy Migration: Database migration tool failed. 3-day delay expected.</div>
            <div className="mt-2 text-xs"><span className="bg-red-100 px-2 py-1 rounded">Owner: Ramesh Kumar</span></div>
          </div>
          <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
            <div className="font-bold text-orange-900">🟠 HIGH</div>
            <div className="text-sm text-gray-700 mt-1">Enterprise Portal: API performance degradation. SLA impact ₹8L if not resolved in 48 hours.</div>
            <div className="mt-2 text-xs"><span className="bg-orange-100 px-2 py-1 rounded">Owner: Priya Sharma</span></div>
          </div>
          <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
            <div className="font-bold text-yellow-900">🟡 MEDIUM</div>
            <div className="text-sm text-gray-700 mt-1">Mobile App: Third-party dependency update required. May conflict with existing integrations.</div>
            <div className="mt-2 text-xs"><span className="bg-yellow-100 px-2 py-1 rounded">Owner: Arjun Singh</span></div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-red-50'
    },
    {
      id: 'summary',
      title: 'Summary & Strategic Focus',
      subtitle: 'Key achievements and Q1 priorities',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ On-Time Delivery</div>
              <div className="text-sm text-gray-700 mt-2">8 of 9 projects tracking on schedule. Strong delivery track record maintained this quarter.</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 Resource Efficiency</div>
              <div className="text-sm text-gray-700 mt-2">82% utilization rate with flexibility for new initiatives. Cost per completed project down 12%.</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="font-bold text-purple-900">⭐ Team Satisfaction</div>
              <div className="text-sm text-gray-700 mt-2">NPS Score: 8.4/10. Strong team morale with 94% retention rate across projects.</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">⚡ Q1 2025 Priorities</div>
              <div className="text-sm text-gray-700 mt-2">
                • Enterprise Portal - Launch to production<br/>
                • Mobile App - Beta to 10K users<br/>
                • Legacy Migration - Complete by Jan 28<br/>
                • Infrastructure - 99.95% uptime target
              </div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">🎯 Risk Mitigation</div>
              <div className="text-sm text-gray-700 mt-2">
                • Database migration - Backup plan ready<br/>
                • API performance - Caching layer added<br/>
                • Team bandwidth - Contractor onboarded<br/>
                • External dependencies - Vendor SLAs updated
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    }
  ]

  return <PresentationDashboard slides={slides} title="Project Management Dashboard" showSlideNumbers={true} />
}
