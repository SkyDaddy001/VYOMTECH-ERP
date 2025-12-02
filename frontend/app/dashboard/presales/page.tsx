'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import InteractionList from '@/components/modules/PostSales/InteractionList'
import InteractionForm from '@/components/modules/PostSales/InteractionForm'
import DocumentList from '@/components/modules/PostSales/DocumentList'
import SnagList from '@/components/modules/PostSales/SnagList'
import ChangeRequestList from '@/components/modules/PostSales/ChangeRequestList'
import { postSalesService } from '@/services/postsales.service'
import { CustomerInteraction, DocumentTracker, SnagList as SnagListType, ChangeRequest, PostSalesMetrics } from '@/types/postsales'

type TabType = 'interactions' | 'documents' | 'snags' | 'crms' | 'kpi'
type FormType = 'interaction' | null

export default function PostSalesPage() {
  const [activeTab, setActiveTab] = useState<TabType>('interactions')
  const [showForm, setShowForm] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [formType, setFormType] = useState<FormType>(null)

  // Data states
  const [interactions, setInteractions] = useState<CustomerInteraction[]>([])
  const [documents, setDocuments] = useState<DocumentTracker[]>([])
  const [snags, setSnags] = useState<SnagListType[]>([])
  const [changeRequests, setChangeRequests] = useState<ChangeRequest[]>([])
  const [metrics, setMetrics] = useState<PostSalesMetrics | null>(null)

  // Loading states
  const [interactionsLoading, setInteractionsLoading] = useState(false)
  const [documentsLoading, setDocumentsLoading] = useState(false)
  const [snagsLoading, setSnagsLoading] = useState(false)
  const [crmsLoading, setCrmsLoading] = useState(false)

  // Load interactions
  const loadInteractions = async () => {
    setInteractionsLoading(true)
    try {
      const data = await postSalesService.getInteractions()
      setInteractions(data)
    } catch (error) {
      toast.error('Failed to load interactions')
    } finally {
      setInteractionsLoading(false)
    }
  }

  // Load documents
  const loadDocuments = async () => {
    setDocumentsLoading(true)
    try {
      const data = await postSalesService.getDocuments()
      setDocuments(data)
    } catch (error) {
      toast.error('Failed to load documents')
    } finally {
      setDocumentsLoading(false)
    }
  }

  // Load snags
  const loadSnags = async () => {
    setSnagsLoading(true)
    try {
      const data = await postSalesService.getSnags()
      setSnags(data)
    } catch (error) {
      toast.error('Failed to load snags')
    } finally {
      setSnagsLoading(false)
    }
  }

  // Load change requests
  const loadChangeRequests = async () => {
    setCrmsLoading(true)
    try {
      const data = await postSalesService.getChangeRequests()
      setChangeRequests(data)
    } catch (error) {
      toast.error('Failed to load change requests')
    } finally {
      setCrmsLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    try {
      const data = await postSalesService.getMetrics()
      setMetrics(data)
    } catch (error) {
      toast.error('Failed to load metrics')
    }
  }

  // Load data on tab change
  useEffect(() => {
    loadMetrics()
    switch (activeTab) {
      case 'interactions':
        loadInteractions()
        break
      case 'documents':
        loadDocuments()
        break
      case 'snags':
        loadSnags()
        break
      case 'crms':
        loadChangeRequests()
        break
    }
  }, [activeTab])

  // Interaction CRUD
  const handleCreateInteraction = () => {
    setEditingItem(null)
    setFormType('interaction')
    setShowForm(true)
  }

  const handleEditInteraction = (interaction: CustomerInteraction) => {
    setEditingItem(interaction)
    setFormType('interaction')
    setShowForm(true)
  }

  const handleDeleteInteraction = async (interaction: CustomerInteraction) => {
    if (!confirm('Are you sure?')) return
    try {
      await postSalesService.deleteInteraction(interaction.id || '')
      toast.success('Interaction deleted!')
      loadInteractions()
    } catch (error) {
      toast.error('Failed to delete interaction')
    }
  }

  const handleUpdateInteractionStatus = async (interaction: CustomerInteraction, status: string) => {
    try {
      await postSalesService.updateInteractionStatus(interaction.id || '', status)
      toast.success('Interaction status updated!')
      loadInteractions()
    } catch (error) {
      toast.error('Failed to update interaction status')
    }
  }

  const handleSubmitInteraction = async (data: Partial<CustomerInteraction>) => {
    try {
      if (editingItem) {
        await postSalesService.updateInteraction(editingItem.id, data)
      } else {
        await postSalesService.createInteraction(data)
      }
      setShowForm(false)
      loadInteractions()
    } catch (error) {
      toast.error('Failed to save interaction')
    }
  }

  const handleUpdateDocumentStatus = async (doc: DocumentTracker, status: string) => {
    try {
      await postSalesService.updateDocumentStatus(doc.id || '', status)
      toast.success('Document status updated!')
      loadDocuments()
    } catch (error) {
      toast.error('Failed to update document status')
    }
  }

  const handleUpdateSnagStatus = async (snag: SnagListType, status: string) => {
    try {
      await postSalesService.updateSnagStatus(snag.id || '', status)
      toast.success('Snag status updated!')
      loadSnags()
    } catch (error) {
      toast.error('Failed to update snag status')
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
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Post-Sales & CRM Management</h1>
          <p className="text-gray-600">Manage customer interactions, documents, snags, and change requests</p>
        </div>

        {/* KPI Metrics */}
        {metrics && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Payment Collection</p>
              <p className="text-2xl font-bold text-green-600 mt-1">{metrics.payment_collection_percentage}%</p>
              <p className="text-xs text-gray-500 mt-1">Target: 98%+</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Agreement Signing TAT</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.agreement_signing_tat} days</p>
              <p className="text-xs text-gray-500 mt-1">Target: &lt;14 days</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Snag Resolution Time</p>
              <p className="text-2xl font-bold text-purple-600 mt-1">{metrics.resolved_snags} resolved</p>
              <p className="text-xs text-gray-500 mt-1">Target: &lt;10 days</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">NPS Score</p>
              <p className="text-2xl font-bold text-orange-600 mt-1">{metrics.nps_score.toFixed(1)}</p>
              <p className="text-xs text-gray-500 mt-1">Target: 8+</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6 border border-gray-200">
          <div className="flex flex-wrap border-b border-gray-200">
            <button
              onClick={() => {
                setActiveTab('interactions')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'interactions'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Interactions
            </button>
            <button
              onClick={() => {
                setActiveTab('documents')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'documents'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Documents
            </button>
            <button
              onClick={() => {
                setActiveTab('snags')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'snags'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Snags
            </button>
            <button
              onClick={() => {
                setActiveTab('crms')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'crms'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              CRM - Change Requests
            </button>
            <button
              onClick={() => {
                setActiveTab('kpi')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'kpi'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              KPI Dashboard
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {/* Interactions Tab */}
            {activeTab === 'interactions' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateInteraction}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Log Interaction
                    </button>
                    <InteractionList
                      interactions={interactions}
                      loading={interactionsLoading}
                      onEdit={handleEditInteraction}
                      onDelete={handleDeleteInteraction}
                      onUpdateStatus={handleUpdateInteractionStatus}
                    />
                  </div>
                ) : (
                  <InteractionForm
                    interaction={editingItem}
                    onSubmit={handleSubmitInteraction}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Documents Tab */}
            {activeTab === 'documents' && (
              <DocumentList
                documents={documents}
                loading={documentsLoading}
                onStatusChange={handleUpdateDocumentStatus}
                onEdit={() => toast('Edit document coming soon', { icon: '✏️' })}
              />
            )}

            {/* Snags Tab */}
            {activeTab === 'snags' && (
              <SnagList
                snags={snags}
                loading={snagsLoading}
                onStatusChange={handleUpdateSnagStatus}
                onEdit={() => toast('Edit snag coming soon', { icon: '✏️' })}
                onDelete={() => toast('Delete functionality coming soon', { icon: '🗑️' })}
              />
            )}

            {/* Change Requests (CRM) Tab */}
            {activeTab === 'crms' && (
              <ChangeRequestList
                changeRequests={changeRequests}
                loading={crmsLoading}
                onEdit={() => toast('Edit CRM coming soon', { icon: '✏️' })}
                onDelete={() => toast('Delete functionality coming soon', { icon: '🗑️' })}
              />
            )}

            {/* KPI Dashboard Tab */}
            {activeTab === 'kpi' && metrics && (
              <div className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {/* Payment Collection */}
                  <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-6 border border-green-200">
                    <h3 className="text-lg font-semibold text-green-900 mb-4">Payment Collection</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-green-700">Current</span>
                        <span className="text-2xl font-bold text-green-600">{metrics.payment_collection_percentage}%</span>
                      </div>
                      <div className="w-full bg-green-200 rounded-full h-3">
                        <div
                          className="bg-green-600 h-3 rounded-full transition-all"
                          style={{ width: `${metrics.payment_collection_percentage}%` }}
                        />
                      </div>
                      <p className="text-xs text-green-700">Target: 98%+ | Status: {metrics.payment_collection_percentage >= 98 ? '✓ On Track' : '⚠ Below Target'}</p>
                    </div>
                  </div>

                  {/* Document TAT */}
                  <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-6 border border-blue-200">
                    <h3 className="text-lg font-semibold text-blue-900 mb-4">Document TAT</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-blue-700">Average Hours</span>
                        <span className="text-2xl font-bold text-blue-600">{metrics.avg_document_tat_hours}h</span>
                      </div>
                      <p className="text-xs text-blue-700">Target: &lt;48 hours</p>
                      <p className="text-xs text-blue-700">Status: {metrics.avg_document_tat_hours <= 48 ? '✓ On Track' : '⚠ Exceeding'}</p>
                    </div>
                  </div>

                  {/* Snag Resolution */}
                  <div className="bg-gradient-to-br from-purple-50 to-purple-100 rounded-lg p-6 border border-purple-200">
                    <h3 className="text-lg font-semibold text-purple-900 mb-4">Snag Resolution SLA</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-purple-700">Resolved Snags</span>
                        <span className="text-2xl font-bold text-purple-600">{metrics.resolved_snags}</span>
                      </div>
                      <p className="text-xs text-purple-700">Pending: {metrics.pending_snags}</p>
                      <p className="text-xs text-purple-700">Avg Target: &lt;10 days</p>
                    </div>
                  </div>

                  {/* NPS Score */}
                  <div className="bg-gradient-to-br from-orange-50 to-orange-100 rounded-lg p-6 border border-orange-200">
                    <h3 className="text-lg font-semibold text-orange-900 mb-4">NPS Score</h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm text-orange-700">Current NPS</span>
                        <span className="text-2xl font-bold text-orange-600">{metrics.nps_score.toFixed(1)}</span>
                      </div>
                      <p className="text-xs text-orange-700">Target: 8+</p>
                      <p className="text-xs text-orange-700">Status: {metrics.nps_score >= 8 ? '✓ Excellent' : '⚠ Needs Improvement'}</p>
                    </div>
                  </div>
                </div>

                {/* Interaction Status Summary */}
                <div className="bg-white rounded-lg p-6 border border-gray-200">
                  <h3 className="text-lg font-semibold text-gray-900 mb-4">Interaction Status</h3>
                  <div className="grid grid-cols-3 gap-4">
                    <div>
                      <p className="text-sm text-gray-600">Pending Follow-ups</p>
                      <p className="text-2xl font-bold text-yellow-600 mt-1">{metrics.pending_interactions || 0}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-600">Resolved</p>
                      <p className="text-2xl font-bold text-green-600 mt-1">{metrics.resolved_interactions || 0}</p>
                    </div>
                    <div>
                      <p className="text-sm text-gray-600">Escalated</p>
                      <p className="text-2xl font-bold text-red-600 mt-1">{metrics.escalated_issues || 0}</p>
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
