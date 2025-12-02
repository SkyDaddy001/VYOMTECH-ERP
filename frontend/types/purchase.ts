export interface Vendor {
  id?: string
  vendor_code: string
  vendor_name: string
  contact_person: string
  email: string
  phone: string
  alternate_phone?: string
  address: string
  city: string
  state: string
  zip_code: string
  gst_number?: string
  pan_number?: string
  bank_account_number?: string
  bank_ifsc_code?: string
  bank_name?: string
  account_holder_name?: string
  payment_terms: 'net_15' | 'net_30' | 'net_45' | 'net_60' | 'advance'
  payment_method: 'bank_transfer' | 'cheque' | 'credit_card' | 'cash'
  status: 'active' | 'inactive' | 'suspended'
  rating: number
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface PurchaseOrderItem {
  item_id?: string
  description: string
  quantity: number
  unit_price: number
  amount: number
}

export interface PurchaseOrder {
  id?: string
  po_number: string
  vendor_id?: string
  vendor_name: string
  po_date: string
  due_date?: string
  delivery_date?: string
  po_status: 'draft' | 'sent' | 'confirmed' | 'partial_received' | 'received' | 'cancelled'
  items: PurchaseOrderItem[]
  subtotal_amount: number
  tax_percentage: number
  tax_amount: number
  shipping_amount?: number
  discount_amount?: number
  total_amount: number
  payment_status: 'pending' | 'partial' | 'paid'
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface GRNLineItem {
  po_item_id: string
  description: string
  quantity_ordered: number
  quantity_received: number
  variance?: number
  quality_status: 'accepted' | 'rejected' | 'partial'
  remarks?: string
}

export interface GoodsReceiptNote {
  id?: string
  grn_number: string
  po_number?: string
  vendor_name: string
  grn_date: string
  received_date: string
  grn_status: 'pending_inspection' | 'inspected' | 'accepted' | 'rejected' | 'partial_accepted'
  items: GRNLineItem[]
  warehouse_location?: string
  received_by?: string
  inspected_by?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface PurchaseMetrics {
  total_vendors: number
  active_vendors: number
  total_po_value: number
  pending_deliveries: number
  average_vendor_rating: number
  total_grn_processed: number
}
