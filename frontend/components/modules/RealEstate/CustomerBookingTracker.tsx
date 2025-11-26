import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Plus, Edit2, Trash2, Eye, Download } from 'lucide-react'
import { formatDateToDDMMMYYYY, formatDateToInput } from '@/lib/dateFormat'

interface Booking {
  id: string
  booking_reference: string
  unit_id: string
  customer_id: string
  booking_date: string
  booking_status: string
  rate_per_sqft: number
  composite_guideline_value: number
  allotment_date?: string
  agreement_date?: string
  registration_date?: string
  handover_date?: string
}

interface Unit {
  id: string
  unit_number: string
  unit_type: string
  carpet_area: number
  plinth_area: number
  sbua: number
  status: string
}

export function CustomerBookingTracker() {
  const [bookings, setBookings] = useState<Booking[]>([])
  const [units, setUnits] = useState<Unit[]>([])
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  
  const [formData, setFormData] = useState({
    unit_id: '',
    customer_id: '',
    booking_date: new Date().toISOString().split('T')[0],
    rate_per_sqft: 0,
    composite_guideline_value: 0,
    car_parking_type: 'covered',
    parking_location: ''
  })

  useEffect(() => {
    fetchBookings()
    fetchUnits()
  }, [])

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
    } catch (error) {
      console.error('Failed to fetch bookings:', error)
    }
  }

  const fetchUnits = async () => {
    try {
      const response = await fetch('/api/v1/real-estate/units', {
        headers: {
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        }
      })
      const data = await response.json()
      setUnits(data || [])
    } catch (error) {
      console.error('Failed to fetch units:', error)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      const url = editingId 
        ? `/api/v1/real-estate/bookings/${editingId}`
        : '/api/v1/real-estate/bookings'
      
      const response = await fetch(url, {
        method: editingId ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
          'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
        },
        body: JSON.stringify(formData)
      })

      if (response.ok) {
        fetchBookings()
        setShowForm(false)
        setEditingId(null)
        setFormData({
          unit_id: '',
          customer_id: '',
          booking_date: new Date().toISOString().split('T')[0],
          rate_per_sqft: 0,
          composite_guideline_value: 0,
          car_parking_type: 'covered',
          parking_location: ''
        })
      }
    } catch (error) {
      console.error('Failed to save booking:', error)
    }
  }

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this booking?')) {
      try {
        await fetch(`/api/v1/real-estate/bookings/${id}`, {
          method: 'DELETE',
          headers: {
            'X-Tenant-ID': localStorage.getItem('tenant_id') || '',
            'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
          }
        })
        fetchBookings()
      } catch (error) {
        console.error('Failed to delete booking:', error)
      }
    }
  }

  const filteredBookings = bookings.filter(booking => {
    const matchesSearch = booking.booking_reference.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesStatus = statusFilter === 'all' || booking.booking_status === statusFilter
    return matchesSearch && matchesStatus
  })

  const getProgressStage = (booking: Booking) => {
    const stages = [
      { date: booking.booking_date, label: 'Booking', icon: '📅' },
      { date: booking.allotment_date, label: 'Allotment', icon: '✓' },
      { date: booking.agreement_date, label: 'Agreement', icon: '📋' },
      { date: booking.registration_date, label: 'Registration', icon: '📝' },
      { date: booking.handover_date, label: 'Handover', icon: '🔑' }
    ]
    
    const completed = stages.filter(s => s.date).length
    return { completed, total: stages.length, stages }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold text-gray-900">Customer Bookings</h2>
        <Button onClick={() => setShowForm(true)} className="bg-blue-600 hover:bg-blue-700">
          <Plus className="w-4 h-4 mr-2" /> New Booking
        </Button>
      </div>

      {/* Filters */}
      <div className="flex gap-4">
        <Input
          placeholder="Search booking reference..."
          value={searchTerm}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchTerm(e.target.value)}
          className="flex-1"
        />
        <Select value={statusFilter} onValueChange={(value: string) => setStatusFilter(value)}>
          <SelectTrigger className="w-48">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Status</SelectItem>
            <SelectItem value="active">Active</SelectItem>
            <SelectItem value="cancelled">Cancelled</SelectItem>
            <SelectItem value="completed">Completed</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Form Modal */}
      {showForm && (
        <Card className="border-blue-300 bg-blue-50">
          <CardHeader>
            <CardTitle>{editingId ? 'Edit Booking' : 'New Customer Booking'}</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-2">Unit</label>
                <Select value={formData.unit_id} onValueChange={(value) => setFormData({...formData, unit_id: value})}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select unit" />
                  </SelectTrigger>
                  <SelectContent>
                    {units.filter(u => u.status === 'available').map(unit => (
                      <SelectItem key={unit.id} value={unit.id}>
                        {unit.unit_number} ({unit.unit_type})
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              
              <div>
                <label className="block text-sm font-medium mb-2">Customer ID</label>
                <Input
                  type="text"
                  value={formData.customer_id}
                  onChange={(e) => setFormData({...formData, customer_id: e.target.value})}
                  placeholder="Customer ID"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Booking Date</label>
                <Input
                  type="date"
                  value={formData.booking_date}
                  onChange={(e) => setFormData({...formData, booking_date: e.target.value})}
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Rate/Sqft</label>
                <Input
                  type="number"
                  step="0.01"
                  value={formData.rate_per_sqft}
                  onChange={(e) => setFormData({...formData, rate_per_sqft: parseFloat(e.target.value)})}
                  placeholder="Rate per sqft"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Composite Guideline Value</label>
                <Input
                  type="number"
                  step="0.01"
                  value={formData.composite_guideline_value}
                  onChange={(e) => setFormData({...formData, composite_guideline_value: parseFloat(e.target.value)})}
                  placeholder="Guideline value"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Parking Type</label>
                <Select value={formData.car_parking_type} onValueChange={(value) => setFormData({...formData, car_parking_type: value})}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="covered">Covered</SelectItem>
                    <SelectItem value="open">Open</SelectItem>
                    <SelectItem value="none">None</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="col-span-2">
                <label className="block text-sm font-medium mb-2">Parking Location</label>
                <Input
                  type="text"
                  value={formData.parking_location}
                  onChange={(e) => setFormData({...formData, parking_location: e.target.value})}
                  placeholder="Parking location"
                />
              </div>

              <div className="col-span-2 flex gap-2">
                <Button type="submit" className="bg-green-600 hover:bg-green-700 flex-1">Save</Button>
                <Button type="button" onClick={() => {setShowForm(false); setEditingId(null)}} className="bg-gray-500 hover:bg-gray-600 flex-1">Cancel</Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {/* Bookings Table */}
      <Card>
        <CardHeader>
          <CardTitle>Bookings List</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Reference</TableHead>
                  <TableHead>Unit</TableHead>
                  <TableHead>Booking Date</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Progress</TableHead>
                  <TableHead>Rate/Sqft</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredBookings.map(booking => {
                  const progress = getProgressStage(booking)
                  return (
                    <TableRow key={booking.id}>
                      <TableCell className="font-semibold">{booking.booking_reference}</TableCell>
                      <TableCell>{booking.unit_id}</TableCell>
                      <TableCell>{formatDateToDDMMMYYYY(booking.booking_date)}</TableCell>
                      <TableCell>
                        <span className={`px-2 py-1 rounded text-xs font-medium ${
                          booking.booking_status === 'active' ? 'bg-green-100 text-green-800' :
                          booking.booking_status === 'cancelled' ? 'bg-red-100 text-red-800' :
                          'bg-blue-100 text-blue-800'
                        }`}>
                          {booking.booking_status}
                        </span>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1 text-xs">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div 
                              className="bg-blue-600 h-2 rounded-full" 
                              style={{width: `${(progress.completed / progress.total) * 100}%`}}
                            />
                          </div>
                          {progress.completed}/{progress.total}
                        </div>
                      </TableCell>
                      <TableCell>₹{booking.rate_per_sqft.toLocaleString()}</TableCell>
                      <TableCell>
                        <div className="flex gap-2">
                          <button className="p-1 hover:bg-blue-50 rounded">
                            <Eye className="w-4 h-4 text-blue-600" />
                          </button>
                          <button 
                            onClick={() => {
                              setEditingId(booking.id)
                              setFormData({
                                unit_id: booking.unit_id,
                                customer_id: booking.customer_id,
                                booking_date: booking.booking_date,
                                rate_per_sqft: booking.rate_per_sqft,
                                composite_guideline_value: booking.composite_guideline_value,
                                car_parking_type: 'covered',
                                parking_location: ''
                              })
                              setShowForm(true)
                            }}
                            className="p-1 hover:bg-yellow-50 rounded"
                          >
                            <Edit2 className="w-4 h-4 text-yellow-600" />
                          </button>
                          <button 
                            onClick={() => handleDelete(booking.id)}
                            className="p-1 hover:bg-red-50 rounded"
                          >
                            <Trash2 className="w-4 h-4 text-red-600" />
                          </button>
                        </div>
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Progress Timeline */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Milestones</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {filteredBookings.slice(0, 5).map(booking => {
              const progress = getProgressStage(booking)
              return (
                <div key={booking.id} className="border-l-4 border-blue-600 pl-4 py-2">
                  <p className="font-semibold text-sm">{booking.booking_reference}</p>
                  <p className="text-xs text-gray-600 mb-2">Progress: {progress.completed}/{progress.total} stages</p>
                  <div className="flex gap-2 flex-wrap">
                    {progress.stages.map((stage, idx) => (
                      <span key={idx} className={`px-2 py-1 rounded text-xs ${
                        stage.date ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-600'
                      }`}>
                        {stage.icon} {stage.label}
                      </span>
                    ))}
                  </div>
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default CustomerBookingTracker
