'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import CampaignList from '@/components/modules/Marketing/CampaignList'
import CampaignForm from '@/components/modules/Marketing/CampaignForm'
import LeadList from '@/components/modules/Marketing/LeadList'
import LeadForm from '@/components/modules/Marketing/LeadForm'
import { marketingService } from '@/services/marketing.service'
import { Campaign, Lead, MarketingMetrics } from '@/types/marketing'

type TabType = 'campaigns' | 'leads' | 'analytics'
type FormType = 'campaign' | 'lead' | null

export default function MarketingPage() {
  const [activeTab, setActiveTab] = useState<TabType>('campaigns')
  const [showForm, setShowForm] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [formType, setFormType] = useState<FormType>(null)

  // Data states
  const [campaigns, setCampaigns] = useState<Campaign[]>([])
  const [leads, setLeads] = useState<Lead[]>([])
  const [metrics, setMetrics] = useState<MarketingMetrics | null>(null)

  // Loading states
  const [campaignsLoading, setCampaignsLoading] = useState(false)
  const [leadsLoading, setLeadsLoading] = useState(false)

  // Load campaigns
  const loadCampaigns = async () => {
    setCampaignsLoading(true)
    try {
      const data = await marketingService.getCampaigns()
      setCampaigns(data)
    } catch (error) {
      toast.error('Failed to load campaigns')
    } finally {
      setCampaignsLoading(false)
    }
  }

  // Load leads
  const loadLeads = async () => {
    setLeadsLoading(true)
    try {
      const data = await marketingService.getLeads()
      setLeads(data)
    } catch (error) {
      toast.error('Failed to load leads')
    } finally {
      setLeadsLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    try {
      const data = await marketingService.getMetrics()
      setMetrics(data)
    } catch (error) {
      toast.error('Failed to load metrics')
    }
  }

  // Load data on tab change
  useEffect(() => {
    loadMetrics()
    switch (activeTab) {
      case 'campaigns':
        loadCampaigns()
        break
      case 'leads':
        loadLeads()
        break
    }
  }, [activeTab])

  // Campaign CRUD
  const handleCreateCampaign = () => {
    setEditingItem(null)
    setFormType('campaign')
    setShowForm(true)
  }

  const handleEditCampaign = (campaign: Campaign) => {
    setEditingItem(campaign)
    setFormType('campaign')
    setShowForm(true)
  }

  const handleDeleteCampaign = async (campaign: Campaign) => {
    if (!confirm('Are you sure?')) return
    try {
      await marketingService.deleteCampaign(campaign.id || '')
      toast.success('Campaign deleted!')
      loadCampaigns()
    } catch (error) {
      toast.error('Failed to delete campaign')
    }
  }

  const handleSubmitCampaign = async (data: Partial<Campaign>) => {
    try {
      if (editingItem) {
        await marketingService.updateCampaign(editingItem.id, data)
      } else {
        await marketingService.createCampaign(data)
      }
      setShowForm(false)
      loadCampaigns()
    } catch (error) {
      throw error
    }
  }

  // Lead CRUD
  const handleCreateLead = () => {
    setEditingItem(null)
    setFormType('lead')
    setShowForm(true)
  }

  const handleEditLead = (lead: Lead) => {
    setEditingItem(lead)
    setFormType('lead')
    setShowForm(true)
  }

  const handleDeleteLead = async (lead: Lead) => {
    if (!confirm('Are you sure?')) return
    try {
      await marketingService.deleteLead(lead.id || '')
      toast.success('Lead deleted!')
      loadLeads()
    } catch (error) {
      toast.error('Failed to delete lead')
    }
  }

  const handleQualifyLead = async (lead: Lead) => {
    try {
      await marketingService.qualifyLead(lead.id || '')
      toast.success('Lead qualified!')
      loadLeads()
    } catch (error) {
      toast.error('Failed to qualify lead')
    }
  }

  const handleConvertLead = async (lead: Lead) => {
    try {
      await marketingService.convertLead(lead.id || '')
      toast.success('Lead converted to booking!')
      loadLeads()
    } catch (error) {
      toast.error('Failed to convert lead')
    }
  }

  const handleSubmitLead = async (data: Partial<Lead>) => {
    try {
      if (editingItem) {
        await marketingService.updateLead(editingItem.id, data)
      } else {
        await marketingService.createLead(data)
      }
      setShowForm(false)
      loadLeads()
    } catch (error) {
      throw error
    }
  }

  const closeForm = () => {
    setShowForm(false)
    setEditingItem(null)
    setFormType(null)
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Real Estate Marketing</h1>
          <p className="text-gray-600">Manage campaigns, leads, and marketing analytics</p>
        </div>

        {/* Real Estate Marketing KPIs */}
        {metrics && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Active Campaigns</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.active_campaigns}</p>
              <p className="text-xs text-gray-500 mt-1">Projects launched</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Leads</p>
              <p className="text-2xl font-bold text-green-600 mt-1">{metrics.total_leads}</p>
              <p className="text-xs text-gray-500 mt-1">Property inquiries</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Qualified Leads</p>
              <p className="text-2xl font-bold text-purple-600 mt-1">{metrics.qualified_leads}</p>
              <p className="text-xs text-gray-500 mt-1">Ready for sales</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Lead to Booking</p>
              <p className="text-2xl font-bold text-orange-600 mt-1">{metrics.conversion_rate_inquiry_to_booking.toFixed(1)}%</p>
              <p className="text-xs text-gray-500 mt-1">Conversion rate</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6 border border-gray-200">
          <div className="flex flex-wrap border-b border-gray-200">
            <button
              onClick={() => {
                setActiveTab('campaigns')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'campaigns'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Project Campaigns
            </button>
            <button
              onClick={() => {
                setActiveTab('leads')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'leads'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Property Inquiries
            </button>
            <button
              onClick={() => {
                setActiveTab('analytics')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'analytics'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Analytics
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {/* Campaigns Tab */}
            {activeTab === 'campaigns' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateCampaign}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Create Campaign
                    </button>
                    <CampaignList
                      campaigns={campaigns}
                      loading={campaignsLoading}
                      onEdit={handleEditCampaign}
                      onDelete={handleDeleteCampaign}
                    />
                  </div>
                ) : (
                  <CampaignForm
                    campaign={editingItem}
                    onSubmit={handleSubmitCampaign}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Leads Tab */}
            {activeTab === 'leads' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateLead}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Add Property Inquiry
                    </button>
                    <LeadList
                      leads={leads}
                      loading={leadsLoading}
                      onEdit={handleEditLead}
                      onDelete={handleDeleteLead}
                      onQualify={handleQualifyLead}
                      onConvert={handleConvertLead}
                    />
                  </div>
                ) : (
                  <LeadForm
                    lead={editingItem}
                    onSubmit={handleSubmitLead}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Analytics Tab */}
            {activeTab === 'analytics' && metrics && (
              <div className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  {/* Lead Generation */}
                  <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-6 border border-blue-200">
                    <h3 className="text-lg font-semibold text-blue-900 mb-4">Lead Generation</h3>
                    <div className="space-y-3">
                      <div>
                        <p className="text-sm text-blue-700">Total Leads Generated</p>
                        <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.total_leads}</p>
                      </div>
                      <div>
                        <p className="text-sm text-blue-700">Site Visits Completed</p>
                        <p className="text-xl font-bold text-blue-600">{metrics.site_visits_completed}</p>
                      </div>
                      <div>
                        <p className="text-sm text-blue-700">Cost per Lead</p>
                        <p className="text-lg font-bold text-blue-600">₹{metrics.avg_cost_per_lead.toLocaleString()}</p>
                      </div>
                    </div>
                  </div>

                  {/* Conversion Performance */}
                  <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-6 border border-green-200">
                    <h3 className="text-lg font-semibold text-green-900 mb-4">Conversion Performance</h3>
                    <div className="space-y-3">
                      <div>
                        <p className="text-sm text-green-700">Lead to Booking Rate</p>
                        <p className="text-2xl font-bold text-green-600 mt-1">{metrics.conversion_rate_inquiry_to_booking.toFixed(1)}%</p>
                      </div>
                      <div>
                        <p className="text-sm text-green-700">Bookings from Leads</p>
                        <p className="text-xl font-bold text-green-600">{metrics.bookings_from_leads}</p>
                      </div>
                      <div>
                        <p className="text-sm text-green-700">Cost per Booking</p>
                        <p className="text-lg font-bold text-green-600">₹{metrics.avg_cost_per_booking.toLocaleString()}</p>
                      </div>
                    </div>
                  </div>

                  {/* Budget Efficiency */}
                  <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-6 border border-purple-200">
                    <h3 className="text-lg font-semibold text-purple-900 mb-4">Budget Efficiency</h3>
                    <div className="space-y-3">
                      <div>
                        <p className="text-sm text-purple-700">Total Budget</p>
                        <p className="text-2xl font-bold text-purple-600 mt-1">₹{metrics.total_budget.toLocaleString()}</p>
                      </div>
                      <div>
                        <p className="text-sm text-purple-700">Amount Spent</p>
                        <p className="text-xl font-bold text-purple-600">₹{metrics.total_spent.toLocaleString()}</p>
                      </div>
                      <div>
                        <p className="text-sm text-purple-700">Marketing ROI</p>
                        <p className="text-lg font-bold text-purple-600">{metrics.avg_roi.toFixed(1)}%</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
