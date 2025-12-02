// Real Estate Booking Transaction
export interface SalesOrder {
  id?: string
  booking_number: string
  customer_id: string
  customer_name?: string
  project_id: string
  project_name?: string
  property_id: string
  unit_type: 'residential' | 'commercial' | 'parking' | 'other'
  super_area: number
  carpet_area: number
  base_price: number
  base_price_per_sqft: number
  total_amount: number
  gst_amount: number
  registration_amount: number
  discount_amount: number
  net_amount: number
  booking_stage: 'inquiry' | 'quote' | 'booking' | 'agreement' | 'completed' | 'cancelled'
  booking_date: string
  possession_date_expected?: string
  broker_id?: string
  broker_commission_percentage?: number
  broker_commission_amount?: number
  sales_executive_id: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface SalesOrderPaymentTerm {
  id?: string
  booking_id: string
  stage: 'booking' | 'agreement' | 'foundation' | 'structure' | 'finishing' | 'possession'
  stage_percentage: number
  amount_due: number
  due_date: string
  amount_paid?: number
  payment_date?: string
}

export interface SalesTarget {
  id?: string
  sales_executive_id: string
  sales_executive_name?: string
  period: string // YYYY-MM
  project_id?: string
  target_bookings: number
  target_amount: number
  achieved_bookings: number
  achieved_amount: number
  achievement_percentage: number
  status: 'not_started' | 'in_progress' | 'completed' | 'exceeded'
  created_at?: string
  updated_at?: string
}

export interface SalesQuota {
  id?: string
  sales_executive_id: string
  sales_executive_name?: string
  quarter: string // Q1, Q2, Q3, Q4
  year: number
  project_ids?: string[]
  quota_bookings: number
  quota_amount: number
  commission_rate_per_booking: number // fixed commission per booking
  commission_rate_percentage?: number // percentage of booking amount
  created_at?: string
  updated_at?: string
}

export interface SalesMetrics {
  total_bookings: number
  total_sales_value: number
  average_booking_value: number
  bookings_this_month: number
  bookings_this_quarter: number
  target_achievement_percentage: number
  commission_earned: number
  pending_commission: number
  cancellations: number
}

export interface SalesForecast {
  period: string
  project_id: string
  project_name: string
  forecasted_bookings: number
  forecasted_amount: number
  confidence_level: number // 0-100
  critical_factors: string // Pipeline status, market conditions, etc.
}
