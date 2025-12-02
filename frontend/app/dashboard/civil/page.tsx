'use client'

import { useState } from 'react'
import { AlertTriangle, CheckCircle, Clock, MapPin, Users } from 'lucide-react'
import type { Site, SafetyIncident, Compliance, Permit, CivilDashboard } from '@/types/civil'

type TabType = 'dashboard' | 'sites' | 'safety' | 'compliance' | 'permits'

export default function CivilPage() {
  const [activeTab, setActiveTab] = useState<TabType>('dashboard')

  // Mock data
  const dashboard: CivilDashboard = {
    total_sites: 12,
    active_sites: 8,
    total_incidents: 23,
    critical_incidents: 2,
    compliance_status: { compliant: 18, non_compliant: 3, in_progress: 4 },
    pending_permits: 5,
    workforce_total: 247,
    safety_score: 8.4,
  }

  const sites: Site[] = [
    {
      id: '1',
      site_name: 'Downtown Tower A',
      location: 'Downtown District',
      project_id: 'P001',
      site_manager: 'John Smith',
      start_date: '2024-01-15',
      expected_end_date: '2025-12-31',
      current_status: 'active',
      site_area_sqm: 15000,
      workforce_count: 85,
    },
    {
      id: '2',
      site_name: 'Industrial Complex Phase 2',
      location: 'Industrial Zone',
      project_id: 'P002',
      site_manager: 'Sarah Johnson',
      start_date: '2024-03-01',
      expected_end_date: '2025-06-30',
      current_status: 'active',
      site_area_sqm: 25000,
      workforce_count: 120,
    },
  ]

  const incidents: SafetyIncident[] = [
    {
      id: '1',
      site_id: '1',
      incident_type: 'accident',
      severity: 'high',
      incident_date: '2024-11-28',
      description: 'Worker fall from scaffolding',
      reported_by: 'John Smith',
      status: 'resolved',
      incident_number: 'INC-2024-001',
    },
    {
      id: '2',
      site_id: '2',
      incident_type: 'near_miss',
      severity: 'medium',
      incident_date: '2024-11-29',
      description: 'Equipment malfunction warning',
      reported_by: 'Sarah Johnson',
      status: 'investigating',
      incident_number: 'INC-2024-002',
    },
  ]

  const compliances: Compliance[] = [
    {
      id: '1',
      site_id: '1',
      compliance_type: 'safety',
      requirement: 'Monthly Safety Audit',
      due_date: '2024-12-15',
      status: 'compliant',
      last_audit_date: '2024-11-15',
      audit_result: 'pass',
      notes: 'All safety protocols met',
    },
    {
      id: '2',
      site_id: '2',
      compliance_type: 'environmental',
      requirement: 'Environmental Impact Assessment',
      due_date: '2024-12-20',
      status: 'in_progress',
      last_audit_date: '2024-10-01',
      audit_result: 'pending',
      notes: 'Assessment underway',
    },
  ]

  const permits: Permit[] = [
    {
      id: '1',
      site_id: '1',
      permit_type: 'Construction',
      permit_number: 'PERMIT-2024-001',
      issued_date: '2024-01-15',
      expiry_date: '2025-12-31',
      issuing_authority: 'City Planning Department',
      status: 'active',
    },
    {
      id: '2',
      site_id: '2',
      permit_type: 'Environmental',
      permit_number: 'PERMIT-2024-002',
      issued_date: '2024-03-01',
      expiry_date: '2025-02-28',
      issuing_authority: 'Environmental Agency',
      status: 'active',
    },
  ]

  const tabs: Array<{ id: TabType; label: string }> = [
    { id: 'dashboard', label: 'Dashboard' },
    { id: 'sites', label: 'Site Management' },
    { id: 'safety', label: 'Safety & Incidents' },
    { id: 'compliance', label: 'Compliance' },
    { id: 'permits', label: 'Permits' },
  ]

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-teal-600 to-teal-800 rounded-lg p-6 text-white">
        <h1 className="text-3xl font-bold">Civil Engineering Module</h1>
        <p className="text-teal-100 mt-2">Manage site operations, safety, compliance, and permits</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200 overflow-x-auto">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-3 font-medium border-b-2 transition whitespace-nowrap ${
              activeTab === tab.id
                ? 'border-teal-600 text-teal-600'
                : 'border-transparent text-gray-600 hover:text-gray-900'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'dashboard' && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm">Total Sites</p>
                <p className="text-3xl font-bold text-gray-900">{dashboard.total_sites}</p>
              </div>
              <MapPin className="text-teal-600" size={32} />
            </div>
            <p className="text-green-600 text-sm mt-2">{dashboard.active_sites} active</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm">Workforce</p>
                <p className="text-3xl font-bold text-gray-900">{dashboard.workforce_total}</p>
              </div>
              <Users className="text-blue-600" size={32} />
            </div>
            <p className="text-gray-600 text-sm mt-2">Personnel deployed</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm">Safety Score</p>
                <p className="text-3xl font-bold text-gray-900">{dashboard.safety_score}/10</p>
              </div>
              <CheckCircle className="text-green-600" size={32} />
            </div>
            <p className="text-green-600 text-sm mt-2">Excellent</p>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm">Critical Incidents</p>
                <p className="text-3xl font-bold text-gray-900">{dashboard.critical_incidents}</p>
              </div>
              <AlertTriangle className="text-red-600" size={32} />
            </div>
            <p className="text-gray-600 text-sm mt-2">Requires attention</p>
          </div>
        </div>
      )}

      {activeTab === 'sites' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Site Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Location</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Manager</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Workforce</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {sites.map((site) => (
                <tr key={site.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{site.site_name}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{site.location}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{site.site_manager}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{site.workforce_count}</td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      {site.current_status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'safety' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Incident #</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Type</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Severity</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {incidents.map((incident) => (
                <tr key={incident.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{incident.incident_number}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{incident.incident_type}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        incident.severity === 'critical'
                          ? 'bg-red-100 text-red-800'
                          : incident.severity === 'high'
                            ? 'bg-orange-100 text-orange-800'
                            : 'bg-yellow-100 text-yellow-800'
                      }`}
                    >
                      {incident.severity}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{incident.incident_date}</td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      {incident.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'compliance' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Requirement</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Type</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Due Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Audit Result</th>
              </tr>
            </thead>
            <tbody>
              {compliances.map((compliance) => (
                <tr key={compliance.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm text-gray-900">{compliance.requirement}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{compliance.compliance_type}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{compliance.due_date}</td>
                  <td className="px-6 py-4">
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-medium ${
                        compliance.status === 'compliant'
                          ? 'bg-green-100 text-green-800'
                          : compliance.status === 'non_compliant'
                            ? 'bg-red-100 text-red-800'
                            : 'bg-yellow-100 text-yellow-800'
                      }`}
                    >
                      {compliance.status}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      {compliance.audit_result}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {activeTab === 'permits' && (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Permit #</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Type</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Authority</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Issued Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Expiry Date</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-900">Status</th>
              </tr>
            </thead>
            <tbody>
              {permits.map((permit) => (
                <tr key={permit.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{permit.permit_number}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{permit.permit_type}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{permit.issuing_authority}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{permit.issued_date}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{permit.expiry_date}</td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      {permit.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
