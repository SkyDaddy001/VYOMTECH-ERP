export interface Invoice {
  id?: string
  invoice_number: string
  customer_id: string
  customer_name?: string
  invoice_date: string
  due_date: string
  total_amount: number
  tax_amount: number
  discount_amount: number
  net_amount: number
  status: 'draft' | 'sent' | 'paid' | 'overdue' | 'cancelled'
  payment_terms?: string
  notes?: string
  items: InvoiceItem[]
  created_at?: string
  updated_at?: string
}

export interface InvoiceItem {
  id?: string
  invoice_id?: string
  description: string
  quantity: number
  unit_price: number
  line_total: number
  tax_rate?: number
}

export interface Payment {
  id?: string
  invoice_id: string
  invoice_number?: string
  payment_date: string
  amount: number
  payment_method: 'cash' | 'check' | 'card' | 'bank_transfer'
  reference_number?: string
  status: 'pending' | 'processed' | 'failed'
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Expense {
  id?: string
  expense_number: string
  category: string
  description: string
  amount: number
  expense_date: string
  submitted_by?: string
  approved_by?: string
  status: 'draft' | 'submitted' | 'approved' | 'rejected' | 'paid'
  receipt_url?: string
  notes?: string
  created_at?: string
  updated_at?: string
}

export interface Account {
  id?: string
  account_name: string
  account_type: 'checking' | 'savings' | 'credit'
  bank_name: string
  account_number: string
  balance: number
  currency: string
  status: 'active' | 'inactive' | 'frozen'
  opening_balance: number
  opening_date: string
  created_at?: string
  updated_at?: string
}

export interface AccountingMetrics {
  total_invoices: number
  total_revenue: number
  outstanding_amount: number
  total_expenses: number
  net_profit: number
  accounts_count: number
}

export interface FinancialReport {
  period: string
  total_income: number
  total_expenses: number
  net_income: number
  profit_margin: number // percentage
}
