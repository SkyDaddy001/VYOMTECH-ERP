export interface Booking {
  id?: string
  unit_id: string
  unit_number?: string
  customer_id: string
  customer_name?: string
  booking_date: string
  booking_reference: string
  booking_status: 'inquired' | 'interested' | 'booked' | 'registered' | 'handed_over' | 'cancelled'
  welcome_date?: string
  allotment_date?: string
  agreement_date?: string
  registration_date?: string
  handover_date?: string
  rate_per_sqft: number
  booking_amount: number
  total_cost: number
  paid_amount: number
  balance_amount: number
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface BookingPayment {
  id?: string
  booking_id: string
  payment_date: string
  payment_mode: 'cash' | 'cheque' | 'transfer' | 'neft' | 'rtgs' | 'dd'
  amount: number
  receipt_number: string
  towards: string
  cheque_number?: string
  cheque_date?: string
  bank_name?: string
  transaction_id?: string
  status: 'pending' | 'cleared' | 'bounced'
  remarks?: string
  created_at?: string
  updated_at?: string
}

export interface BookingPaymentSchedule {
  id?: string
  booking_id: string
  number_of_installments: number
  payment_stage: 'booking' | 'agreement' | 'foundation' | 'structure' | 'finishing' | 'possession'
  installment_date: string
  installment_amount: number
  status: 'pending' | 'received' | 'overdue' | 'bounced'
  created_at?: string
  updated_at?: string
}

export interface BookingMetrics {
  total_bookings: number
  active_bookings: number
  completed_handovers: number
  cancelled_bookings: number
  total_booked_value: number
  total_received: number
  pending_realization: number
}

export interface BookingTimeline {
  month: string
  new_bookings: number
  handovers: number
  collections: number
}
