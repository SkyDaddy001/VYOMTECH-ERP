export interface Property {
  id?: string
  property_id: string
  project_id?: string
  unit_type: 'residential' | 'commercial' | 'parking' | 'retail'
  wing?: string
  floor?: string | number
  unit_number: string
  property_status: 'available' | 'booked' | 'sold' | 'under_construction' | 'ready'
  super_area: number
  carpet_area: number
  builtup_area?: number
  terrace_area?: number
  parking_count: number
  facing: 'north' | 'south' | 'east' | 'west' | 'northeast' | 'northwest' | 'southeast' | 'southwest'
  configuration: 'studio' | '1bhk' | '2bhk' | '3bhk' | '4bhk' | '5bhk' | 'penthouse'
  base_price: number
  base_price_per_sqft: number
  possession_date_expected?: string
  possession_date_actual?: string
  ownership_type: 'freehold' | 'leasehold'
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface ProjectLayout {
  id?: string
  project_id: string
  layout_name: string
  total_units: number
  total_area: number
  amenities: string[]
  parking_ratio: number
  master_plan_url?: string
  created_at?: string
  updated_at?: string
}

export interface PropertyAvailability {
  property_id: string
  status: 'available' | 'booked' | 'sold' | 'on_hold'
  hold_till_date?: string
  updated_at?: string
}

export interface RealEstateMetrics {
  total_properties: number
  available_properties: number
  booked_properties: number
  sold_properties: number
  total_area: number
  total_value: number
  occupancy_rate: number
}

export interface PropertyHistory {
  id?: string
  property_id: string
  booking_id?: string
  customer_id?: string
  customer_name?: string
  action: 'created' | 'booked' | 'sold' | 'possession' | 'status_change'
  action_date: string
  notes?: string
  created_at?: string
}
