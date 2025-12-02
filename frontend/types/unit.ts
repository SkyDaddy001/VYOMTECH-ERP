export interface PropertyUnit {
  id: string
  tenant_id: string
  project_id: string
  block_id: string
  unit_number: string
  floor: number
  unit_type: string
  facing: string
  carpet_area: number
  carpet_area_with_balcony: number
  utility_area: number
  plinth_area: number
  sbua: number
  uds_sqft: number
  status: 'available' | 'booked' | 'sold' | 'reserved'
  alloted_to?: string
  allotment_date?: string
  created_at: string
  updated_at: string
}

export interface UnitCostSheet {
  id: string
  tenant_id: string
  unit_id: string
  rate_per_sqft: number
  sbua_rate: number
  base_price: number
  frc: number
  car_parking_cost: number
  plc: number
  statutory_charges: number
  other_charges: number
  legal_charges: number
  apartment_cost_exc_govt: number
  apartment_cost_inc_govt: number
  composite_guideline_value: number
  actual_sold_price: number
  car_parking_type: string
  parking_location: string
  effective_date?: string
  validity_date?: string
  created_at: string
  updated_at: string
}

export interface UnitSummary {
  total_units: number
  available: number
  booked: number
  sold: number
  reserved: number
  total_value: number
}
