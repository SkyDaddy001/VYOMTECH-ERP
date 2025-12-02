export interface Ledger {
  id?: string
  booking_id: string
  customer_id: string
  customer_name?: string
  transaction_date: string
  transaction_type: 'debit' | 'credit'
  description: string
  debit_amount: number
  credit_amount: number
  opening_balance: number
  closing_balance: number
  payment_id?: string
  reference_number: string
  created_at?: string
  updated_at?: string
}

export interface LedgerSummary {
  customer_id: string
  customer_name?: string
  opening_balance: number
  total_debit: number
  total_credit: number
  closing_balance: number
  last_transaction_date?: string
}

export interface LedgerMetrics {
  total_transactions: number
  total_debit: number
  total_credit: number
  net_balance: number
  customers_with_outstanding: number
  total_outstanding: number
}

export interface LedgerReport {
  period: string
  transactions: number
  total_amount: number
  balance: number
}
