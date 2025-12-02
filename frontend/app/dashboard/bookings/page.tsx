'use client'

import { useState, useEffect } from 'react'
import toast from 'react-hot-toast'
import BookingList from '@/components/modules/Bookings/BookingList'
import BookingForm from '@/components/modules/Bookings/BookingForm'
import PaymentList from '@/components/modules/Bookings/PaymentList'
import { bookingsService } from '@/services/bookings.service'
import { Booking, BookingPayment, BookingMetrics } from '@/types/bookings'

type TabType = 'bookings' | 'payments'
type FormType = 'booking' | 'payment' | null

export default function BookingsPage() {
  const [activeTab, setActiveTab] = useState<TabType>('bookings')
  const [showForm, setShowForm] = useState(false)
  const [editingItem, setEditingItem] = useState<any>(null)
  const [formType, setFormType] = useState<FormType>(null)

  // Data states
  const [bookings, setBookings] = useState<Booking[]>([])
  const [payments, setPayments] = useState<BookingPayment[]>([])
  const [metrics, setMetrics] = useState<BookingMetrics | null>(null)
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null)

  // Loading states
  const [bookingsLoading, setBookingsLoading] = useState(false)
  const [paymentsLoading, setPaymentsLoading] = useState(false)

  // Load bookings
  const loadBookings = async () => {
    setBookingsLoading(true)
    try {
      const data = await bookingsService.getBookings()
      setBookings(data)
      if (data.length > 0 && !selectedBooking) {
        setSelectedBooking(data[0])
      }
    } catch (error) {
      toast.error('Failed to load bookings')
    } finally {
      setBookingsLoading(false)
    }
  }

  // Load payments for selected booking
  const loadPayments = async () => {
    if (!selectedBooking?.id) return
    setPaymentsLoading(true)
    try {
      const data = await bookingsService.getBookingPayments(selectedBooking.id)
      setPayments(data)
    } catch (error) {
      toast.error('Failed to load payments')
    } finally {
      setPaymentsLoading(false)
    }
  }

  // Load metrics
  const loadMetrics = async () => {
    try {
      const data = await bookingsService.getMetrics()
      setMetrics(data)
    } catch (error) {
      // Metrics are optional
    }
  }

  // Load data on mount and tab change
  useEffect(() => {
    loadMetrics()
    if (activeTab === 'bookings') {
      loadBookings()
    } else if (activeTab === 'payments') {
      loadPayments()
    }
  }, [activeTab])

  // Booking CRUD
  const handleCreateBooking = () => {
    setEditingItem(null)
    setFormType('booking')
    setShowForm(true)
  }

  const handleEditBooking = (booking: Booking) => {
    setEditingItem(booking)
    setFormType('booking')
    setShowForm(true)
  }

  const handleDeleteBooking = async (booking: Booking) => {
    if (!confirm('Are you sure?')) return
    try {
      await bookingsService.deleteBooking(booking.id || '')
      toast.success('Booking deleted!')
      loadBookings()
    } catch (error) {
      toast.error('Failed to delete booking')
    }
  }

  const handleUpdateBookingStatus = async (booking: Booking, status: string) => {
    try {
      await bookingsService.updateBookingStatus(booking.id || '', status)
      toast.success('Status updated!')
      loadBookings()
    } catch (error) {
      toast.error('Failed to update status')
    }
  }

  const handleSubmitBooking = async (data: Partial<Booking>) => {
    try {
      if (editingItem) {
        await bookingsService.updateBooking(editingItem.id, data)
      } else {
        await bookingsService.createBooking(data)
      }
      setShowForm(false)
      loadBookings()
    } catch (error) {
      throw error
    }
  }

  // Payment CRUD
  const handleAddPayment = async (data: Partial<BookingPayment>) => {
    try {
      await bookingsService.addPayment(selectedBooking?.id || '', data)
      toast.success('Payment added!')
      loadPayments()
      setShowForm(false)
    } catch (error) {
      throw error
    }
  }

  const handleUpdatePayment = async (data: Partial<BookingPayment>) => {
    try {
      await bookingsService.updatePayment(selectedBooking?.id || '', editingItem.id, data)
      toast.success('Payment updated!')
      loadPayments()
      setShowForm(false)
    } catch (error) {
      throw error
    }
  }

  const handleDeletePayment = async (payment: BookingPayment) => {
    if (!confirm('Are you sure?')) return
    try {
      await bookingsService.deletePayment(selectedBooking?.id || '', payment.id || '')
      toast.success('Payment deleted!')
      loadPayments()
    } catch (error) {
      toast.error('Failed to delete payment')
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
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Customer Bookings</h1>
          <p className="text-gray-600">Manage customer bookings, payments, and handovers</p>
        </div>

        {/* Metrics */}
        {metrics && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Total Bookings</p>
              <p className="text-2xl font-bold text-blue-600 mt-1">{metrics.total_bookings}</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Booked Value</p>
              <p className="text-2xl font-bold text-green-600 mt-1">₹{(metrics.total_booked_value / 10000000).toFixed(0)}Cr</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Received</p>
              <p className="text-2xl font-bold text-purple-600 mt-1">₹{(metrics.total_received / 10000000).toFixed(0)}Cr</p>
            </div>
            <div className="bg-white rounded-lg shadow p-4 border border-gray-200">
              <p className="text-gray-600 text-xs font-medium">Pending</p>
              <p className="text-2xl font-bold text-orange-600 mt-1">₹{(metrics.pending_realization / 10000000).toFixed(0)}Cr</p>
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
              Bookings
            </button>
            <button
              onClick={() => {
                setActiveTab('payments')
                setShowForm(false)
              }}
              className={`flex-1 py-4 px-6 text-center font-medium transition ${
                activeTab === 'payments'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Payments
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
                      onClick={handleCreateBooking}
                      className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                    >
                      + New Booking
                    </button>
                    <BookingList
                      bookings={bookings}
                      loading={bookingsLoading}
                      onEdit={handleEditBooking}
                      onDelete={handleDeleteBooking}
                      onUpdateStatus={handleUpdateBookingStatus}
                    />
                  </div>
                ) : (
                  <BookingForm
                    booking={editingItem}
                    onSubmit={handleSubmitBooking}
                    onCancel={closeForm}
                  />
                )}
              </div>
            )}

            {/* Payments Tab */}
            {activeTab === 'payments' && (
              <div>
                {selectedBooking ? (
                  <>
                    <div className="mb-4 p-4 bg-blue-50 border border-blue-200 rounded-lg flex justify-between items-center">
                      <div>
                        <p className="text-sm text-gray-600">Booking: <span className="font-semibold text-blue-700">{selectedBooking.customer_name}</span></p>
                        <p className="text-xs text-gray-600 mt-1">Ref: {selectedBooking.booking_reference}</p>
                      </div>
                      <select
                        value={selectedBooking.id || ''}
                        onChange={(e) => {
                          const selected = bookings.find(b => b.id === e.target.value)
                          if (selected) setSelectedBooking(selected)
                        }}
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                      >
                        {bookings.map(b => (
                          <option key={b.id} value={b.id}>{b.customer_name}</option>
                        ))}
                      </select>
                    </div>
                    {!showForm ? (
                      <div>
                        <button
                          onClick={() => {
                            setEditingItem(null)
                            setFormType('payment')
                            setShowForm(true)
                          }}
                          className="mb-4 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 font-medium"
                        >
                          + Add Payment
                        </button>
                        <PaymentList
                          payments={payments}
                          loading={paymentsLoading}
                          onEdit={(payment) => {
                            setEditingItem(payment)
                            setFormType('payment')
                            setShowForm(true)
                          }}
                          onDelete={handleDeletePayment}
                        />
                      </div>
                    ) : (
                      <div className="bg-white rounded-lg p-6 max-w-2xl">
                        <h3 className="text-lg font-semibold mb-4">{editingItem ? 'Edit Payment' : 'Add Payment'}</h3>
                        <p className="text-gray-600 text-sm mb-4">Coming soon: Full payment form implementation</p>
                        <button
                          onClick={closeForm}
                          className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 font-medium"
                        >
                          Cancel
                        </button>
                      </div>
                    )}
                  </>
                ) : (
                  <div className="text-center py-12 text-gray-500">
                    <p>Please create or select a booking first</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
