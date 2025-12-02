'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import SalesOrderList from '@/components/modules/Sales/SalesOrderList'
import SalesOrderForm from '@/components/modules/Sales/SalesOrderForm'
import SalesTargetList from '@/components/modules/Sales/SalesTargetList'
import SalesTargetForm from '@/components/modules/Sales/SalesTargetForm'
import SalesQuotaList from '@/components/modules/Sales/SalesQuotaList'
import SalesQuotaForm from '@/components/modules/Sales/SalesQuotaForm'
import { salesService } from '@/services/sales.service'
import { SalesOrder, SalesTarget, SalesQuota, SalesMetrics, SalesForecast } from '@/types/sales'

type TabType = 'bookings' | 'targets' | 'quotas' | 'forecast'
type FormType = 'booking' | 'target' | 'quota' | null

export default function SalesPage() {
  const [activeTab, setActiveTab] = useState<TabType>('bookings')
  const [showForm, setShowForm] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [formType, setFormType] = useState<FormType>(null)

  // Data states
  const [orders, setOrders] = useState<SalesOrder[]>([])
  const [targets, setTargets] = useState<SalesTarget[]>([])
  const [quotas, setQuotas] = useState<SalesQuota[]>([])
  const [metrics, setMetrics] = useState<SalesMetrics | null>(null)
  const [forecast, setForecast] = useState<SalesForecast[]>([])

  // Loading states
  const [ordersLoading, setOrdersLoading] = useState(false)
  const [targetsLoading, setTargetsLoading] = useState(false)
  const [quotasLoading, setQuotasLoading] = useState(false)

  // Load bookings
  const loadOrders = async () => {
    setOrdersLoading(true)
    try {
      const data = await salesService.getOrders()
      setOrders(data)
    } catch (error) {
      toast.error('Failed to load bookings')
    } finally {
      setOrdersLoading(false)
    }
  }

  // Load targets
  const loadTargets = async () => {
    setTargetsLoading(true)
    try {
      const data = await salesService.getTargets()
      setTargets(data)
    } catch (error) {
      toast.error('Failed to load targets')
    } finally {
      setTargetsLoading(false)
    }
  }

  // Load quotas
  const loadQuotas = async () => {
    setQuotasLoading(true)
    try {
      const data = await salesService.getQuotas()
      setQuotas(data)
    } catch (error) {
      toast.error('Failed to load quotas')
    } finally {
      setQuotasLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    try {
      const data = await salesService.getMetrics()
      setMetrics(data)
    } catch (error) {
      toast.error('Failed to load metrics')
    }
  }

  // Load forecast
  const loadForecast = async () => {
    try {
      const data = await salesService.getForecast()
      setForecast(data)
    } catch (error) {
      // Forecast is optional
    }
  }

  // Load data on tab change
  useEffect(() => {
    loadMetrics()
    switch (activeTab) {
      case 'bookings':
        loadOrders()
        break
      case 'targets':
        loadTargets()
        break
      case 'quotas':
        loadQuotas()
        break
      case 'forecast':
        loadForecast()
        break
    }
  }, [activeTab])

  // Booking CRUD
  const handleCreateOrder = () => {
    setEditingItem(null)
    setFormType('booking')
    setShowForm(true)
  }

  const handleEditOrder = (order: SalesOrder) => {
    setEditingItem(order)
    setFormType('booking')
    setShowForm(true)
  }

  const handleDeleteOrder = async (order: SalesOrder) => {
    if (!confirm('Are you sure?')) return
    try {
      await salesService.deleteOrder(order.id || '')
      toast.success('Booking deleted!')
      loadOrders()
    } catch (error) {
      toast.error('Failed to delete booking')
    }
  }

  const handleUpdateOrderStatus = async (order: SalesOrder, status: string) => {
    try {
      await salesService.updateOrderStatus(order.id || '', status)
      toast.success('Booking stage updated!')
      loadOrders()
    } catch (error) {
      toast.error('Failed to update booking stage')
    }
  }

  const handleSubmitOrder = async (data: Partial<SalesOrder>) => {
    try {
      if (editingItem) {
        await salesService.updateOrder(editingItem.id, data)
      } else {
        await salesService.createOrder(data)
      }
      setShowForm(false)
      loadOrders()
    } catch (error) {
      throw error
    }
  }

  // Target CRUD
  const handleCreateTarget = () => {
    setEditingItem(null)
    setFormType('target')
    setShowForm(true)
  }

  const handleEditTarget = (target: SalesTarget) => {
    setEditingItem(target)
    setFormType('target')
    setShowForm(true)
  }

  const handleDeleteTarget = async (target: SalesTarget) => {
    if (!confirm('Are you sure?')) return
    try {
      await salesService.deleteTarget(target.id || '')
      toast.success('Target deleted!')
      loadTargets()
    } catch (error) {
      toast.error('Failed to delete target')
    }
  }

  const handleSubmitTarget = async (data: Partial<SalesTarget>) => {
    try {
      if (editingItem) {
        await salesService.updateTarget(editingItem.id, data)
      } else {
        await salesService.createTarget(data)
      }
      setShowForm(false)
      loadTargets()
    } catch (error) {
      throw error
    }
  }

  // Quota CRUD
  const handleCreateQuota = () => {
    setEditingItem(null)
    setFormType('quota')
    setShowForm(true)
  }

  const handleEditQuota = (quota: SalesQuota) => {
    setEditingItem(quota)
    setFormType('quota')
    setShowForm(true)
  }

  const handleDeleteQuota = async (quota: SalesQuota) => {
    if (!confirm('Are you sure?')) return
    try {
      await salesService.deleteQuota(quota.id || '')
      toast.success('Quota deleted!')
      loadQuotas()
    } catch (error) {
      toast.error('Failed to delete quota')
    }
  }

  const handleSubmitQuota = async (data: Partial<SalesQuota>) => {
    try {
      if (editingItem) {
        await salesService.updateQuota(editingItem.id, data)
      } else {
        await salesService.createQuota(data)
      }
      setShowForm(false)
      loadQuotas()
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
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Real Estate Sales</h1>
          <p className="text-gray-600">Manage property bookings, sales targets, and commission</p>
        </div>

        {/* Real Estate Sales KPIs */}
        {metrics && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Bookings</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.total_bookings}</p>
              <p className="text-xs text-gray-500 mt-1">Properties booked</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Sales Value</p>
              <p className="text-2xl font-bold text-green-600 mt-1">₹{(metrics.total_sales_value / 10000000).toFixed(1)}Cr</p>
              <p className="text-xs text-gray-500 mt-1">Total value</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Avg Booking Value</p>
              <p className="text-2xl font-bold text-purple-600 mt-1">₹{(metrics.average_booking_value / 100000).toFixed(1)}L</p>
              <p className="text-xs text-gray-500 mt-1">Per property</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Target Achievement</p>
              <p className="text-2xl font-bold text-orange-600 mt-1">{metrics.target_achievement_percentage.toFixed(1)}%</p>
              <p className="text-xs text-gray-500 mt-1">This quarter</p>
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="bg-white rounded-lg shadow mb-6 border border-gray-200">
          <div className="flex flex-wrap border-b border-gray-200">
            <button
              onClick={() => {
                setActiveTab('bookings')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'bookings'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Property Bookings
            </button>
            <button
              onClick={() => {
                setActiveTab('targets')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'targets'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Monthly Targets
            </button>
            <button
              onClick={() => {
                setActiveTab('quotas')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'quotas'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Quarterly Quotas
            </button>
            <button
              onClick={() => {
                setActiveTab('forecast')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'forecast'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Forecast
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-6">
            {/* Bookings Tab */}
            {activeTab === 'bookings' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateOrder}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Create Booking
                    </button>
                    <SalesOrderList
                      orders={orders}
                      loading={ordersLoading}
                      onEdit={handleEditOrder}
                      onDelete={handleDeleteOrder}
                      onStatusChange={handleUpdateOrderStatus}
                    />
                  </div>
                ) : (
                  <SalesOrderForm
                    order={editingItem}
                    onSubmit={handleSubmitOrder}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Targets Tab */}
            {activeTab === 'targets' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateTarget}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Set Target
                    </button>
                    <SalesTargetList
                      targets={targets}
                      loading={targetsLoading}
                      onEdit={handleEditTarget}
                      onDelete={handleDeleteTarget}
                    />
                  </div>
                ) : (
                  <SalesTargetForm
                    target={editingItem}
                    onSubmit={handleSubmitTarget}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Quotas Tab */}
            {activeTab === 'quotas' && (
              <div>
                {!showForm ? (
                  <div>
                    <button
                      onClick={handleCreateQuota}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + Set Quota
                    </button>
                    <SalesQuotaList
                      quotas={quotas}
                      loading={quotasLoading}
                      onEdit={handleEditQuota}
                      onDelete={handleDeleteQuota}
                    />
                  </div>
                ) : (
                  <SalesQuotaForm
                    quota={editingItem}
                    onSubmit={handleSubmitQuota}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Forecast Tab */}
            {activeTab === 'forecast' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-gray-900">Project Booking Forecast</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {forecast.map((f, idx) => (
                    <div key={idx} className="bg-gradient-to-br from-cyan-50 to-cyan-100 rounded-lg p-4 border border-cyan-200">
                      <p className="text-sm font-semibold text-cyan-900 mb-3">{f.project_name}</p>
                      <div className="space-y-2">
                        <div>
                          <p className="text-xs text-cyan-700">Period</p>
                          <p className="font-medium text-cyan-900">{f.period}</p>
                        </div>
                        <div>
                          <p className="text-xs text-cyan-700">Forecasted Bookings</p>
                          <p className="text-xl font-bold text-cyan-600">{f.forecasted_bookings}</p>
                        </div>
                        <div>
                          <p className="text-xs text-cyan-700">Forecasted Amount</p>
                          <p className="text-lg font-bold text-cyan-600">₹{(f.forecasted_amount / 10000000).toFixed(1)}Cr</p>
                        </div>
                        <div className="pt-2">
                          <p className="text-xs text-cyan-700">Confidence Level</p>
                          <div className="w-full bg-white rounded-full h-2 overflow-hidden mt-1">
                            <div
                              className="h-full bg-cyan-500"
                              style={{ width: `${f.confidence_level}%` }}
                            />
                          </div>
                          <p className="text-xs text-cyan-700 mt-1">{f.confidence_level}%</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
