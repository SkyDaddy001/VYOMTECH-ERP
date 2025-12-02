'use client'

import { useState } from 'react'
import CallRecordForm from './CallRecordForm'
import CallRecordList from './CallRecordList'
import LeadList from './LeadList'
import { CallRecord, Lead } from '@/types/presales'

export default function PresalesDashboard() {
  const [activeTab, setActiveTab] = useState<'overview' | 'calls' | 'leads'>('overview')
  const [callRecords, setCallRecords] = useState<CallRecord[]>([])
  const [leads, setLeads] = useState<Lead[]>([])
  const [showNewCallForm, setShowNewCallForm] = useState(false)

  const handleSaveCallRecord = async (data: Partial<CallRecord>) => {
    try {
      // TODO: Make API call to save call record
      const newRecord: CallRecord = {
        id: Date.now().toString(),
        ...data as CallRecord
      }
      setCallRecords([newRecord, ...callRecords])
      setShowNewCallForm(false)
      
      // Auto-create or update lead
      const existingLead = leads.find(l => l.customer_phone === data.customer_phone)
      if (!existingLead) {
        const newLead: Lead = {
          id: Date.now().toString(),
          customer_name: data.customer_name || '',
          customer_phone: data.customer_phone || '',
          email: data.email,
          source: 'call',
          status: data.call_outcome === 'converted' ? 'converted' : 'contacted',
          priority: determinePriority(data),
          assigned_to: data.presales_agent_id || '',
          first_contact_date: new Date().toISOString(),
          last_contact_date: new Date().toISOString(),
          call_records: [newRecord as CallRecord]
        }
        setLeads([newLead, ...leads])
      } else {
        // Update existing lead
        setLeads(leads.map(l => 
          l.id === existingLead.id 
            ? { 
                ...l, 
                status: data.call_outcome === 'converted' ? 'converted' : 'contacted',
                last_contact_date: new Date().toISOString(),
                next_follow_up_date: data.follow_up_date,
                call_records: [...(l.call_records || []), newRecord as CallRecord]
              }
            : l
        ))
      }
    } catch (error) {
      console.error('Error saving call record:', error)
      alert('Failed to save call record')
    }
  }

  const determinePriority = (data: Partial<CallRecord>): 'low' | 'medium' | 'high' | 'urgent' => {
    let priority: 'low' | 'medium' | 'high' | 'urgent' = 'low'
    
    if (data.call_outcome === 'interested') {
      priority = 'high'
    } else if (data.call_outcome === 'follow_up_needed') {
      priority = 'medium'
    } else if (data.call_outcome === 'converted') {
      priority = 'urgent'
    }
    
    // Boost priority if budget is high
    if (data.budget_range === 'above_5_crores' || data.budget_range === '2_to_5_crores') {
      priority = 'urgent'
    }
    
    return priority
  }

  const getMetrics = () => {
    return {
      totalCalls: callRecords.length,
      callsThisMonth: callRecords.filter(c => {
        const callDate = new Date(c.call_date)
        const now = new Date()
        return callDate.getMonth() === now.getMonth() && callDate.getFullYear() === now.getFullYear()
      }).length,
      interestedCount: callRecords.filter(c => c.call_outcome === 'interested').length,
      convertedCount: callRecords.filter(c => c.call_outcome === 'converted').length,
      followUpPending: callRecords.filter(c => c.follow_up_required && c.call_outcome !== 'converted').length,
      totalLeads: leads.length,
      qualifiedLeads: leads.filter(l => ['qualified', 'converted'].includes(l.status)).length,
      conversionRate: callRecords.length > 0 
        ? ((callRecords.filter(c => c.call_outcome === 'converted').length / callRecords.length) * 100).toFixed(1)
        : '0'
    }
  }

  const metrics = getMetrics()

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
            <span className="text-4xl">📞</span> Presales Management
          </h1>
          <p className="text-gray-600 mt-2">Track call records and manage leads efficiently</p>
        </div>

        {/* Tabs */}
        <div className="flex gap-4 mb-6 border-b border-gray-200">
          {[
            { id: 'overview', label: '📊 Overview', icon: 'overview' },
            { id: 'calls', label: '📞 Call Records', icon: 'calls' },
            { id: 'leads', label: '👥 Leads', icon: 'leads' }
          ].map(tab => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as any)}
              className={`px-4 py-3 font-medium transition-colors border-b-2 -mb-px ${
                activeTab === tab.id
                  ? 'text-blue-600 border-blue-600'
                  : 'text-gray-600 border-transparent hover:text-gray-800'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Overview Tab */}
        {activeTab === 'overview' && (
          <div className="space-y-6">
            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Total Calls</p>
                    <p className="text-3xl font-bold text-gray-900">{metrics.totalCalls}</p>
                  </div>
                  <span className="text-4xl">📞</span>
                </div>
              </div>

              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">This Month</p>
                    <p className="text-3xl font-bold text-gray-900">{metrics.callsThisMonth}</p>
                  </div>
                  <span className="text-4xl">📅</span>
                </div>
              </div>

              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Conversion Rate</p>
                    <p className="text-3xl font-bold text-gray-900">{metrics.conversionRate}%</p>
                  </div>
                  <span className="text-4xl">📈</span>
                </div>
              </div>

              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Total Leads</p>
                    <p className="text-3xl font-bold text-gray-900">{metrics.totalLeads}</p>
                  </div>
                  <span className="text-4xl">👥</span>
                </div>
              </div>
            </div>

            {/* Detailed Metrics */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Call Metrics</h3>
                <div className="space-y-3">
                  <div className="flex justify-between items-center py-2 border-b border-gray-200">
                    <span className="text-gray-700">Interested Leads</span>
                    <span className="font-semibold text-green-600">{metrics.interestedCount}</span>
                  </div>
                  <div className="flex justify-between items-center py-2 border-b border-gray-200">
                    <span className="text-gray-700">Converted to Booking</span>
                    <span className="font-semibold text-purple-600">{metrics.convertedCount}</span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                    <span className="text-gray-700">Follow-ups Pending</span>
                    <span className="font-semibold text-orange-600">{metrics.followUpPending}</span>
                  </div>
                </div>
              </div>

              <div className="bg-white rounded-lg border border-gray-200 p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Lead Metrics</h3>
                <div className="space-y-3">
                  <div className="flex justify-between items-center py-2 border-b border-gray-200">
                    <span className="text-gray-700">Total Active Leads</span>
                    <span className="font-semibold text-blue-600">{metrics.totalLeads}</span>
                  </div>
                  <div className="flex justify-between items-center py-2 border-b border-gray-200">
                    <span className="text-gray-700">Qualified Leads</span>
                    <span className="font-semibold text-green-600">{metrics.qualifiedLeads}</span>
                  </div>
                  <div className="flex justify-between items-center py-2">
                    <span className="text-gray-700">Conversion Rate</span>
                    <span className="font-semibold text-blue-600">{metrics.conversionRate}%</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Quick Action Button */}
            <div className="bg-blue-50 rounded-lg border border-blue-200 p-6 text-center">
              <p className="text-blue-900 mb-4 font-medium">Ready to log a new call?</p>
              <button
                onClick={() => {
                  setShowNewCallForm(true)
                  setActiveTab('calls')
                }}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors"
              >
                + Create New Call Record
              </button>
            </div>
          </div>
        )}

        {/* Call Records Tab */}
        {activeTab === 'calls' && (
          <div>
            {showNewCallForm ? (
              <CallRecordForm
                onSubmit={handleSaveCallRecord}
                onCancel={() => setShowNewCallForm(false)}
              />
            ) : (
              <CallRecordList
                records={callRecords}
                onCreateNew={() => setShowNewCallForm(true)}
              />
            )}
          </div>
        )}

        {/* Leads Tab */}
        {activeTab === 'leads' && (
          <LeadList leads={leads} />
        )}
      </div>
    </div>
  )
}
