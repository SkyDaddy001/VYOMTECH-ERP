import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Plus, Calendar, TrendingUp, DollarSign, Download } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface Milestone {
  id: string
  booking_id: string
  campaign_name: string
  source: string
  subsource: string
  lead_generated_date?: string
  re_engaged_date?: string
  site_visit_date?: string
  revisit_date?: string
  booking_date?: string
  cancelled_date?: string
  status: string
  notes: string
  created_at: string
}

interface Payment {
  id: string
  booking_id: string
  payment_date: string
  payment_mode: string
  amount: number
  receipt_number: string
  towards: string
  status: string
  transaction_id?: string
}

interface Booking {
  id: string
  booking_reference: string
  unit_id: string
}

export function MilestoneAndPaymentTracking() {
  const [bookings, setBookings] = useState<Booking[]>([])
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [payments, setPayments] = useState<Payment[]>([])
  const [selectedBooking, setSelectedBooking] = useState<string | null>(null)
  const [showMilestoneForm, setShowMilestoneForm] = useState(false)
  const [showPaymentForm, setShowPaymentForm] = useState(false)
  const [activeTab, setActiveTab] = useState<'milestones' | 'payments' | 'timeline'>('milestones')

  const [milestoneForm, setMilestoneForm] = useState({
    campaign_name: '',
    source: 'direct',
    subsource: '',
    lead_generated_date: new Date().toISOString().split('T')[0],
    re_engaged_date: '',
    site_visit_date: '',
    revisit_date: '',
    booking_date: '',
    cancelled_date: '',
    notes: ''
  })

  const [paymentForm, setPaymentForm] = useState({
    payment_date: new Date().toISOString().split('T')[0],
    payment_mode: 'bank_transfer',
    amount: 0,
    receipt_number: '',
    towards: 'advance',
    transaction_id: ''
  })

  useEffect(() => {
    fetchBookings()
  }, [])

  useEffect(() => {
    if (selectedBooking) {
      fetchMilestones(selectedBooking)
      fetchPayments(selectedBooking)
    }
  }, [selectedBooking])

  const fetchBookings = async () => {
    try {
      const response = await fetch('/api/v1/real-estate/bookings', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        }
      })
      const data = await response.json()
      setBookings(data || [])
      if (data && data.length > 0 && !selectedBooking) {
        setSelectedBooking(data[0].id)
      }
    } catch (error) {
      console.error('Failed to fetch bookings:', error)
    }
  }

  const fetchMilestones = async (bookingId: string) => {
    try {
      const response = await fetch(`/api/v1/real-estate/milestones/${bookingId}`, {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        }
      })
      const data = await response.json()
      setMilestones(data || [])
    } catch (error) {
      console.error('Failed to fetch milestones:', error)
    }
  }

  const fetchPayments = async (bookingId: string) => {
    try {
      const response = await fetch(`/api/v1/real-estate/bookings/${bookingId}/payments`, {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        }
      })
      const data = await response.json()
      setPayments(data || [])
    } catch (error) {
      console.error('Failed to fetch payments:', error)
    }
  }

  const handleAddMilestone = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedBooking) return

    try {
      const response = await fetch('/api/v1/real-estate/milestones', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: JSON.stringify({
          booking_id: selectedBooking,
          ...milestoneForm
        })
      })

      if (response.ok) {
        fetchMilestones(selectedBooking)
        setShowMilestoneForm(false)
        setMilestoneForm({
          campaign_name: '',
          source: 'direct',
          subsource: '',
          lead_generated_date: new Date().toISOString().split('T')[0],
          re_engaged_date: '',
          site_visit_date: '',
          revisit_date: '',
          booking_date: '',
          cancelled_date: '',
          notes: ''
        })
      }
    } catch (error) {
      console.error('Failed to add milestone:', error)
    }
  }

  const handleAddPayment = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedBooking) return

    try {
      const response = await fetch('/api/v1/real-estate/payments', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: JSON.stringify({
          booking_id: selectedBooking,
          ...paymentForm
        })
      })

      if (response.ok) {
        fetchPayments(selectedBooking)
        setShowPaymentForm(false)
        setPaymentForm({
          payment_date: new Date().toISOString().split('T')[0],
          payment_mode: 'bank_transfer',
          amount: 0,
          receipt_number: '',
          towards: 'advance',
          transaction_id: ''
        })
      }
    } catch (error) {
      console.error('Failed to record payment:', error)
    }
  }

  const milestoneTimeline = milestones
    .filter(m => m.status !== 'cancelled')
    .map(m => ({
      date: m.lead_generated_date || m.booking_date || '',
      label: m.campaign_name || 'Milestone',
      source: m.source,
      type: 'milestone'
    }))
    .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())

  const totalPayments = payments.reduce((sum, p) => sum + p.amount, 0)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold text-gray-900">Milestone & Payment Tracking</h2>
      </div>

      {/* Booking Selector */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center gap-4">
            <label className="text-sm font-medium">Select Booking:</label>
            <Select value={selectedBooking || ''} onValueChange={setSelectedBooking}>
              <SelectTrigger className="w-64">
                <SelectValue placeholder="Choose a booking" />
              </SelectTrigger>
              <SelectContent>
                {bookings.map(booking => (
                  <SelectItem key={booking.id} value={booking.id}>
                    {booking.booking_reference} - {booking.unit_id}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Tabs */}
      <div className="flex gap-4 border-b border-gray-200">
        <button
          onClick={() => setActiveTab('milestones')}
          className={`px-4 py-2 font-medium border-b-2 ${
            activeTab === 'milestones'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'
          }`}
        >
          <Calendar className="w-4 h-4 inline mr-2" /> Milestones
        </button>
        <button
          onClick={() => setActiveTab('payments')}
          className={`px-4 py-2 font-medium border-b-2 ${
            activeTab === 'payments'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'
          }`}
        >
          <DollarSign className="w-4 h-4 inline mr-2" /> Payments
        </button>
        <button
          onClick={() => setActiveTab('timeline')}
          className={`px-4 py-2 font-medium border-b-2 ${
            activeTab === 'timeline'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'
          }`}
        >
          <TrendingUp className="w-4 h-4 inline mr-2" /> Timeline
        </button>
      </div>

      {/* MILESTONES TAB */}
      {activeTab === 'milestones' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h3 className="text-xl font-bold">Campaign Milestones</h3>
            <Button onClick={() => setShowMilestoneForm(true)} className="bg-blue-600 hover:bg-blue-700">
              <Plus className="w-4 h-4 mr-2" /> Add Milestone
            </Button>
          </div>

          {/* Milestone Form */}
          {showMilestoneForm && (
            <Card className="border-blue-300 bg-blue-50">
              <CardHeader>
                <CardTitle>Record Milestone</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleAddMilestone} className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium mb-2">Campaign Name</label>
                    <Input
                      value={milestoneForm.campaign_name}
                      onChange={(e) => setMilestoneForm({...milestoneForm, campaign_name: e.target.value})}
                      placeholder="e.g., Summer Campaign 2025"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Source</label>
                    <Select value={milestoneForm.source} onValueChange={(value) => setMilestoneForm({...milestoneForm, source: value})}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="direct">Direct</SelectItem>
                        <SelectItem value="site_visit">Site Visit</SelectItem>
                        <SelectItem value="broker">Broker</SelectItem>
                        <SelectItem value="referral">Referral</SelectItem>
                        <SelectItem value="digital">Digital</SelectItem>
                        <SelectItem value="exhibition">Exhibition</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Sub Source</label>
                    <Input
                      value={milestoneForm.subsource}
                      onChange={(e) => setMilestoneForm({...milestoneForm, subsource: e.target.value})}
                      placeholder="e.g., Google Ads, Event, Agent Name"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Lead Generated Date</label>
                    <Input
                      type="date"
                      value={milestoneForm.lead_generated_date}
                      onChange={(e) => setMilestoneForm({...milestoneForm, lead_generated_date: e.target.value})}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Re-engaged Date</label>
                    <Input
                      type="date"
                      value={milestoneForm.re_engaged_date}
                      onChange={(e) => setMilestoneForm({...milestoneForm, re_engaged_date: e.target.value})}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Site Visit Date</label>
                    <Input
                      type="date"
                      value={milestoneForm.site_visit_date}
                      onChange={(e) => setMilestoneForm({...milestoneForm, site_visit_date: e.target.value})}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Re-visit Date</label>
                    <Input
                      type="date"
                      value={milestoneForm.revisit_date}
                      onChange={(e) => setMilestoneForm({...milestoneForm, revisit_date: e.target.value})}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Booking Date</label>
                    <Input
                      type="date"
                      value={milestoneForm.booking_date}
                      onChange={(e) => setMilestoneForm({...milestoneForm, booking_date: e.target.value})}
                    />
                  </div>

                  <div className="col-span-2">
                    <label className="block text-sm font-medium mb-2">Notes</label>
                    <Input
                      value={milestoneForm.notes}
                      onChange={(e) => setMilestoneForm({...milestoneForm, notes: e.target.value})}
                      placeholder="Additional notes"
                    />
                  </div>

                  <div className="col-span-2 flex gap-2">
                    <Button type="submit" className="bg-green-600 hover:bg-green-700 flex-1">Save Milestone</Button>
                    <Button type="button" onClick={() => setShowMilestoneForm(false)} className="bg-gray-500 hover:bg-gray-600 flex-1">Cancel</Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}

          {/* Milestones List */}
          <Card>
            <CardHeader>
              <CardTitle>Recorded Milestones</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {milestones.length === 0 ? (
                  <p className="text-gray-500 text-center py-8">No milestones recorded yet</p>
                ) : (
                  milestones.map((milestone, idx) => (
                    <div key={milestone.id} className="border-l-4 border-blue-600 pl-4 py-3 hover:bg-blue-50 rounded">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <p className="font-semibold text-sm">{milestone.campaign_name || 'Campaign'}</p>
                          <div className="grid grid-cols-4 gap-4 mt-2 text-xs text-gray-600">
                            <p><strong>Source:</strong> {milestone.source}</p>
                            <p><strong>Sub-source:</strong> {milestone.subsource || 'N/A'}</p>
                            <p><strong>Status:</strong> {milestone.status}</p>
                            <p><strong>Created:</strong> {formatDateToDDMMMYYYY(milestone.created_at)}</p>
                          </div>
                          {milestone.notes && <p className="mt-2 text-xs text-gray-700">📝 {milestone.notes}</p>}
                        </div>
                      </div>

                      {/* Milestone dates */}
                      <div className="mt-3 flex flex-wrap gap-2">
                        {milestone.lead_generated_date && (
                          <span className="bg-green-100 text-green-800 px-2 py-1 rounded text-xs">
                            Lead: {formatDateToDDMMMYYYY(milestone.lead_generated_date)}
                          </span>
                        )}
                        {milestone.site_visit_date && (
                          <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded text-xs">
                            Visit: {formatDateToDDMMMYYYY(milestone.site_visit_date)}
                          </span>
                        )}
                        {milestone.revisit_date && (
                          <span className="bg-yellow-100 text-yellow-800 px-2 py-1 rounded text-xs">
                            Re-visit: {formatDateToDDMMMYYYY(milestone.revisit_date)}
                          </span>
                        )}
                        {milestone.booking_date && (
                          <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs">
                            Booking: {formatDateToDDMMMYYYY(milestone.booking_date)}
                          </span>
                        )}
                        {milestone.cancelled_date && (
                          <span className="bg-red-100 text-red-800 px-2 py-1 rounded text-xs">
                            Cancelled: {formatDateToDDMMMYYYY(milestone.cancelled_date)}
                          </span>
                        )}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* PAYMENTS TAB */}
      {activeTab === 'payments' && (
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="text-xl font-bold">Payment Records</h3>
              <p className="text-sm text-gray-600">Total Received: ₹{totalPayments.toLocaleString('en-IN')}</p>
            </div>
            <Button onClick={() => setShowPaymentForm(true)} className="bg-green-600 hover:bg-green-700">
              <Plus className="w-4 h-4 mr-2" /> Record Payment
            </Button>
          </div>

          {/* Payment Form */}
          {showPaymentForm && (
            <Card className="border-green-300 bg-green-50">
              <CardHeader>
                <CardTitle>Record Payment</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleAddPayment} className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium mb-2">Payment Date</label>
                    <Input
                      type="date"
                      value={paymentForm.payment_date}
                      onChange={(e) => setPaymentForm({...paymentForm, payment_date: e.target.value})}
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Payment Mode</label>
                    <Select value={paymentForm.payment_mode} onValueChange={(value) => setPaymentForm({...paymentForm, payment_mode: value})}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="cash">Cash</SelectItem>
                        <SelectItem value="cheque">Cheque</SelectItem>
                        <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
                        <SelectItem value="neft">NEFT</SelectItem>
                        <SelectItem value="rtgs">RTGS</SelectItem>
                        <SelectItem value="demand_draft">Demand Draft</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Amount (₹)</label>
                    <Input
                      type="number"
                      step="0.01"
                      value={paymentForm.amount}
                      onChange={(e) => setPaymentForm({...paymentForm, amount: parseFloat(e.target.value)})}
                      placeholder="Amount"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Receipt Number</label>
                    <Input
                      value={paymentForm.receipt_number}
                      onChange={(e) => setPaymentForm({...paymentForm, receipt_number: e.target.value})}
                      placeholder="Receipt number"
                      required
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Towards</label>
                    <Select value={paymentForm.towards} onValueChange={(value) => setPaymentForm({...paymentForm, towards: value})}>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="advance">Advance</SelectItem>
                        <SelectItem value="booking">Booking Amount</SelectItem>
                        <SelectItem value="installment_1">Installment 1</SelectItem>
                        <SelectItem value="installment_2">Installment 2</SelectItem>
                        <SelectItem value="balance">Balance</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium mb-2">Transaction ID (Optional)</label>
                    <Input
                      value={paymentForm.transaction_id}
                      onChange={(e) => setPaymentForm({...paymentForm, transaction_id: e.target.value})}
                      placeholder="Transaction ID"
                    />
                  </div>

                  <div className="col-span-2 flex gap-2">
                    <Button type="submit" className="bg-green-600 hover:bg-green-700 flex-1">Record Payment</Button>
                    <Button type="button" onClick={() => setShowPaymentForm(false)} className="bg-gray-500 hover:bg-gray-600 flex-1">Cancel</Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}

          {/* Payments Table */}
          <Card>
            <CardHeader>
              <CardTitle>Payment History</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Date</TableHead>
                      <TableHead>Receipt No</TableHead>
                      <TableHead>Mode</TableHead>
                      <TableHead>Amount</TableHead>
                      <TableHead>Towards</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Transaction ID</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {payments.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="text-center py-8 text-gray-500">
                          No payments recorded
                        </TableCell>
                      </TableRow>
                    ) : (
                      payments.map(payment => (
                        <TableRow key={payment.id}>
                          <TableCell>{formatDateToDDMMMYYYY(payment.payment_date)}</TableCell>
                          <TableCell className="font-semibold">{payment.receipt_number}</TableCell>
                          <TableCell>
                            <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs font-medium">
                              {payment.payment_mode}
                            </span>
                          </TableCell>
                          <TableCell className="font-semibold text-green-600">₹{payment.amount.toLocaleString('en-IN')}</TableCell>
                          <TableCell>{payment.towards}</TableCell>
                          <TableCell>
                            <span className="px-2 py-1 bg-green-100 text-green-800 rounded text-xs font-medium">
                              {payment.status}
                            </span>
                          </TableCell>
                          <TableCell className="text-xs">{payment.transaction_id || 'N/A'}</TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* TIMELINE TAB */}
      {activeTab === 'timeline' && (
        <Card>
          <CardHeader>
            <CardTitle>Campaign & Payment Timeline</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {milestoneTimeline.length === 0 && payments.length === 0 ? (
                <p className="text-gray-500 text-center py-8">No timeline data available</p>
              ) : (
                <div className="space-y-3">
                  {/* Timeline Entry */}
                  {milestoneTimeline.map((item, idx) => (
                    <div key={idx} className="flex gap-4">
                      <div className="flex flex-col items-center">
                        <div className="w-4 h-4 rounded-full bg-blue-600" />
                        {idx < milestoneTimeline.length - 1 && (
                          <div className="w-0.5 h-12 bg-gray-300" />
                        )}
                      </div>
                      <div className="pb-4">
                        <p className="font-semibold text-sm">{item.label}</p>
                        <p className="text-xs text-gray-600">{formatDateToDDMMMYYYY(item.date)}</p>
                        <p className="text-xs text-gray-500">Source: {item.source}</p>
                      </div>
                    </div>
                  ))}

                  <hr className="my-4" />

                  {/* Payments on Timeline */}
                  {payments
                    .sort((a, b) => new Date(a.payment_date).getTime() - new Date(b.payment_date).getTime())
                    .map((payment, idx) => (
                      <div key={payment.id} className="flex gap-4">
                        <div className="flex flex-col items-center">
                          <div className="w-4 h-4 rounded-full bg-green-600" />
                          {idx < payments.length - 1 && (
                            <div className="w-0.5 h-12 bg-gray-300" />
                          )}
                        </div>
                        <div className="pb-4">
                          <p className="font-semibold text-sm">Payment: ₹{payment.amount.toLocaleString('en-IN')}</p>
                          <p className="text-xs text-gray-600">{formatDateToDDMMMYYYY(payment.payment_date)}</p>
                          <p className="text-xs text-gray-500">{payment.towards} ({payment.receipt_number})</p>
                        </div>
                      </div>
                    ))}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}

export default MilestoneAndPaymentTracking
