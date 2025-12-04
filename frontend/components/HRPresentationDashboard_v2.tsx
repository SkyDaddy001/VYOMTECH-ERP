'use client'

import React, { useState, useEffect } from 'react'
import PresentationDashboard, { Slide } from '@/components/PresentationDashboard'
import { hrDashboardService } from '@/services/api'
import { Users, Award, TrendingUp, AlertCircle } from 'lucide-react'

export default function HRPresentationDashboard() {
  // State for HR data
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [hrOverview, setHrOverview] = useState<any>(null)
  const [headcountData, setHeadcountData] = useState<any>(null)
  const [attendanceData, setAttendanceData] = useState<any>(null)
  const [performanceData, setPerformanceData] = useState<any>(null)

  // Fetch HR data on mount
  useEffect(() => {
    const fetchHRData = async () => {
      try {
        setLoading(true)

        // Fetch overview
        const overviewRes = await hrDashboardService.getHROverview()
        setHrOverview(overviewRes.data)

        // Fetch headcount by department
        const headcountRes = await hrDashboardService.getHeadcountByDepartment()
        setHeadcountData(headcountRes.data)

        // Fetch attendance data
        const attendanceRes = await hrDashboardService.getAttendanceDashboard()
        setAttendanceData(attendanceRes.data)

        // Fetch performance metrics
        const performanceRes = await hrDashboardService.getPerformanceMetrics()
        setPerformanceData(performanceRes.data)

        setError(null)
      } catch (err: any) {
        console.error('Failed to fetch HR data:', err)
        setError(err.message || 'Failed to load HR data')
      } finally {
        setLoading(false)
      }
    }

    fetchHRData()
  }, [])

  // Use real data or fallback values
  const totalEmployees = hrOverview?.total_employees || 245
  const attendanceRate = hrOverview?.attendance_rate || 94
  const avgSatisfaction = hrOverview?.avg_satisfaction || 8.2
  const yoyGrowth = hrOverview?.yoy_growth || 12

  const presentThisMonth = attendanceData?.present_this_month || 230
  const leavesApplied = attendanceData?.leaves_applied || 12
  const absent = attendanceData?.absent || 3

  const departments = headcountData?.departments || [
    { name: 'Engineering', count: 89, utilization: 85 },
    { name: 'Sales', count: 67, utilization: 75 },
    { name: 'Operations', count: 52, utilization: 68 },
    { name: 'Finance', count: 23, utilization: 72 },
    { name: 'HR & Admin', count: 14, utilization: 55 },
    { name: 'Marketing', count: 18, utilization: 62 }
  ]

  const slides: Slide[] = [
    {
      id: 'cover',
      title: 'Human Resources',
      subtitle: 'Workforce Management & Performance Overview',
      content: (
        <div className="flex flex-col items-center justify-center h-full gap-8">
          <Users className="w-20 h-20 text-blue-600" />
          <div className="grid grid-cols-2 gap-6 w-full max-w-2xl">
            <div className="bg-blue-50 p-6 rounded-lg border border-blue-200">
              <div className="text-3xl font-bold text-blue-700">{totalEmployees}</div>
              <div className="text-sm text-gray-600 mt-1">Total Employees</div>
            </div>
            <div className="bg-green-50 p-6 rounded-lg border border-green-200">
              <div className="text-3xl font-bold text-green-700">{attendanceRate}%</div>
              <div className="text-sm text-gray-600 mt-1">Attendance Rate</div>
            </div>
            <div className="bg-purple-50 p-6 rounded-lg border border-purple-200">
              <div className="text-3xl font-bold text-purple-700">{avgSatisfaction.toFixed(1)}/10</div>
              <div className="text-sm text-gray-600 mt-1">Avg Satisfaction</div>
            </div>
            <div className="bg-orange-50 p-6 rounded-lg border border-orange-200">
              <div className="text-3xl font-bold text-orange-700">{yoyGrowth}%</div>
              <div className="text-sm text-gray-600 mt-1">YoY Growth</div>
            </div>
          </div>
          {error && <p className="text-red-600 text-sm">{error}</p>}
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'headcount',
      title: 'Headcount & Utilization',
      subtitle: 'Department-wise workforce allocation',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            {departments.slice(0, 3).map((dept: any, idx: number) => (
              <div key={idx} className="bg-white p-4 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-semibold text-gray-800">{dept.name}</span>
                  <span className="text-2xl font-bold text-blue-600">{dept.count}</span>
                </div>
                <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                  <div className="bg-blue-500 h-full" style={{ width: `${dept.utilization}%` }}></div>
                </div>
              </div>
            ))}
          </div>
          <div className="space-y-4">
            {departments.slice(3).map((dept: any, idx: number) => (
              <div key={idx} className="bg-white p-4 rounded-lg border border-gray-200">
                <div className="flex justify-between items-center mb-2">
                  <span className="font-semibold text-gray-800">{dept.name}</span>
                  <span className="text-2xl font-bold text-blue-600">{dept.count}</span>
                </div>
                <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                  <div className="bg-blue-500 h-full" style={{ width: `${dept.utilization}%` }}></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'attendance',
      title: 'Attendance & Leave Analytics',
      subtitle: 'Monthly trends and absence tracking',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-gradient-to-br from-green-50 to-green-100 p-6 rounded-lg border-2 border-green-500">
              <div className="text-sm text-gray-600">Present This Month</div>
              <div className="text-4xl font-bold text-green-700 mt-2">{presentThisMonth}</div>
              <div className="text-xs text-gray-600 mt-1">{attendanceRate}% attendance rate</div>
            </div>
            <div className="bg-gradient-to-br from-yellow-50 to-yellow-100 p-6 rounded-lg border-2 border-yellow-500">
              <div className="text-sm text-gray-600">Leaves Applied</div>
              <div className="text-4xl font-bold text-yellow-700 mt-2">{leavesApplied}</div>
              <div className="text-xs text-gray-600 mt-1">Approved this month</div>
            </div>
            <div className="bg-gradient-to-br from-red-50 to-red-100 p-6 rounded-lg border-2 border-red-500">
              <div className="text-sm text-gray-600">Absent</div>
              <div className="text-4xl font-bold text-red-700 mt-2">{absent}</div>
              <div className="text-xs text-gray-600 mt-1">Requires follow-up</div>
            </div>
          </div>
          <div className="bg-white rounded-lg border border-gray-200 p-6">
            <h3 className="font-bold text-gray-800 mb-4">Leave Type Breakdown</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm">Paid Leave</span>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className="bg-blue-500 h-full" style={{ width: '65%' }}></div>
                  </div>
                  <span className="text-sm font-semibold">8/12</span>
                </div>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm">Sick Leave</span>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className="bg-orange-500 h-full" style={{ width: '45%' }}></div>
                  </div>
                  <span className="text-sm font-semibold">2/6</span>
                </div>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm">Casual Leave</span>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className="bg-green-500 h-full" style={{ width: '85%' }}></div>
                  </div>
                  <span className="text-sm font-semibold">2/2</span>
                </div>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm">Maternity Leave</span>
                <div className="flex items-center gap-2">
                  <div className="w-20 bg-gray-200 h-2 rounded-full overflow-hidden">
                    <div className="bg-pink-500 h-full" style={{ width: '25%' }}></div>
                  </div>
                  <span className="text-sm font-semibold">1/4</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-green-50'
    },
    {
      id: 'performance',
      title: 'Performance & Development',
      subtitle: 'Employee ratings and training participation',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-white p-6 rounded-lg border border-gray-200">
              <h3 className="font-bold text-gray-800 mb-4">Top Performers</h3>
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <div>
                    <div className="font-semibold text-gray-800">John Mehta</div>
                    <div className="text-xs text-gray-500">Engineering Lead</div>
                  </div>
                  <div className="flex gap-1">
                    {[...Array(5)].map((_, i) => (
                      <span key={i} className="text-yellow-400">★</span>
                    ))}
                  </div>
                </div>
                <hr className="my-2" />
                <div className="flex justify-between items-center">
                  <div>
                    <div className="font-semibold text-gray-800">Priya Singh</div>
                    <div className="text-xs text-gray-500">Sales Manager</div>
                  </div>
                  <div className="flex gap-1">
                    {[...Array(5)].map((_, i) => (
                      <span key={i} className={i < 4 ? "text-yellow-400" : "text-gray-300"}>★</span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-white p-6 rounded-lg border border-gray-200">
              <h3 className="font-bold text-gray-800 mb-4">Training Participation</h3>
              <div className="space-y-3">
                <div>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-semibold">Technical Skills</span>
                    <span className="text-sm">87%</span>
                  </div>
                  <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                    <div className="bg-blue-500 h-full" style={{ width: '87%' }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-semibold">Leadership</span>
                    <span className="text-sm">64%</span>
                  </div>
                  <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                    <div className="bg-purple-500 h-full" style={{ width: '64%' }}></div>
                  </div>
                </div>
                <div>
                  <div className="flex justify-between mb-1">
                    <span className="text-sm font-semibold">Soft Skills</span>
                    <span className="text-sm">72%</span>
                  </div>
                  <div className="w-full bg-gray-200 h-3 rounded-full overflow-hidden">
                    <div className="bg-green-500 h-full" style={{ width: '72%' }}></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    },
    {
      id: 'recruitment',
      title: 'Recruitment Pipeline',
      subtitle: 'Hiring status and candidate funnel',
      content: (
        <div className="space-y-4 h-full">
          <div className="grid grid-cols-4 gap-4">
            <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-lg border border-blue-300">
              <div className="text-xs text-gray-600 mb-2">Job Openings</div>
              <div className="text-3xl font-bold text-blue-700">12</div>
            </div>
            <div className="bg-gradient-to-br from-yellow-50 to-yellow-100 p-4 rounded-lg border border-yellow-300">
              <div className="text-xs text-gray-600 mb-2">Applications</div>
              <div className="text-3xl font-bold text-yellow-700">78</div>
            </div>
            <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-4 rounded-lg border border-purple-300">
              <div className="text-xs text-gray-600 mb-2">Interviews</div>
              <div className="text-3xl font-bold text-purple-700">23</div>
            </div>
            <div className="bg-gradient-to-br from-green-50 to-green-100 p-4 rounded-lg border border-green-300">
              <div className="text-xs text-gray-600 mb-2">Offers</div>
              <div className="text-3xl font-bold text-green-700">5</div>
            </div>
          </div>
          <div className="grid grid-cols-3 gap-4 flex-1">
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm font-bold text-gray-800 mb-3">Senior Engineer</div>
              <div className="space-y-2">
                <div className="text-xs"><span className="bg-blue-100 text-blue-800 px-2 py-1 rounded">5 Applications</span></div>
                <div className="text-xs"><span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded">2 In Interview</span></div>
                <div className="text-xs"><span className="bg-green-100 text-green-800 px-2 py-1 rounded">1 Offer Sent</span></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm font-bold text-gray-800 mb-3">Sales Executive</div>
              <div className="space-y-2">
                <div className="text-xs"><span className="bg-blue-100 text-blue-800 px-2 py-1 rounded">8 Applications</span></div>
                <div className="text-xs"><span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded">3 In Interview</span></div>
                <div className="text-xs"><span className="bg-red-100 text-red-800 px-2 py-1 rounded">Awaiting Decision</span></div>
              </div>
            </div>
            <div className="bg-white p-4 rounded-lg border border-gray-200">
              <div className="text-sm font-bold text-gray-800 mb-3">Operations Manager</div>
              <div className="space-y-2">
                <div className="text-xs"><span className="bg-blue-100 text-blue-800 px-2 py-1 rounded">3 Applications</span></div>
                <div className="text-xs"><span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded">1 In Interview</span></div>
                <div className="text-xs"><span className="bg-gray-100 text-gray-800 px-2 py-1 rounded">Awaiting Schedule</span></div>
              </div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-gray-50 to-blue-50'
    },
    {
      id: 'summary',
      title: 'HR Summary & Actions',
      subtitle: 'Key insights and upcoming focus areas',
      content: (
        <div className="grid grid-cols-2 gap-6 h-full">
          <div className="space-y-4">
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded">
              <div className="font-bold text-green-900">✓ Strong Hiring Pipeline</div>
              <div className="text-sm text-gray-700 mt-2">12 open positions with 78 qualified applicants and robust interview process in progress.</div>
            </div>
            <div className="bg-blue-50 border-l-4 border-blue-500 p-4 rounded">
              <div className="font-bold text-blue-900">📊 High Engagement Rates</div>
              <div className="text-sm text-gray-700 mt-2">{avgSatisfaction.toFixed(1)}/10 employee satisfaction score indicates strong organizational culture and workplace morale.</div>
            </div>
            <div className="bg-purple-50 border-l-4 border-purple-500 p-4 rounded">
              <div className="font-bold text-purple-900">🎯 Training Investment</div>
              <div className="text-sm text-gray-700 mt-2">72% participation in skill development programs shows commitment to employee growth.</div>
            </div>
          </div>
          <div className="space-y-4">
            <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
              <div className="font-bold text-yellow-900">⚠️ Attendance Watch</div>
              <div className="text-sm text-gray-700 mt-2">{absent} employees with recent absences require follow-up. Schedule wellness check-ins this week.</div>
            </div>
            <div className="bg-orange-50 border-l-4 border-orange-500 p-4 rounded">
              <div className="font-bold text-orange-900">📅 Upcoming Actions</div>
              <div className="text-sm text-gray-700 mt-2">
                • Annual performance review cycle (2 weeks)<br/>
                • Q4 salary reviews due<br/>
                • Leadership training program starts
              </div>
            </div>
            <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
              <div className="font-bold text-red-900">⚡ Retention Risk</div>
              <div className="text-sm text-gray-700 mt-2">2 high-performers on extended leave. Consider mentoring and growth opportunities.</div>
            </div>
          </div>
        </div>
      ),
      backgroundColor: 'from-white to-gray-50'
    }
  ]

  return <PresentationDashboard slides={slides} title="Human Resources Dashboard" showSlideNumbers={true} />
}
