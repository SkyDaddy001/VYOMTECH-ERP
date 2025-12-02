export interface Vendor {
  id?: string
  vendor_name: string
  vendor_type: 'material_supplier' | 'contractor' | 'consultant' | 'service_provider'
  contact_person: string
  phone: string
  email: string
  address: string
  city: string
  state: string
  postal_code: string
  gst_number?: string
  bank_account?: string
  ifsc_code?: string
  rating: number // 1-5
  status: 'active' | 'inactive' | 'suspended'
  created_at?: string
  updated_at?: string
}

export interface PurchaseOrder {
  id?: string
  po_number: string
  vendor_id: string
  vendor_name?: string
  project_id?: string
  item_description: string
  quantity: number
  unit: string
  unit_rate: number
  total_amount: number
  po_date: string
  delivery_date: string
  gst_percentage: number
  gst_amount: number
  net_amount: number
  status: 'draft' | 'confirmed' | 'received' | 'invoiced' | 'cancelled'
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface VendorMetrics {
  total_vendors: number
  active_vendors: number
  total_po_value: number
  pending_deliveries: number
  average_vendor_rating: number
}

export interface VendorPayment {
  id?: string
  vendor_id: string
  vendor_name?: string
  po_id?: string
  payment_date: string
  amount: number
  payment_mode: 'check' | 'transfer' | 'neft' | 'cash'
  reference_number: string
  status: 'pending' | 'cleared'
  created_at?: string
  updated_at?: string
}
