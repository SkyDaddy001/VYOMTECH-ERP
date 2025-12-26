// Payment utilities and constants

export const PAYMENT_METHODS = {
  NETBANKING: 'netbanking',
  CREDIT_CARD: 'credit_card',
  DEBIT_CARD: 'debit_card',
  UPI: 'upi',
  WALLET: 'wallet',
} as const

export const PAYMENT_PROVIDERS = {
  RAZORPAY: 'razorpay',
  BILLDESK: 'billdesk',
} as const

export const PAYMENT_STATUS = {
  PENDING: 'pending',
  INITIATED: 'initiated',
  PROCESSING: 'processing',
  SUCCESSFUL: 'successful',
  FAILED: 'failed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const

export const CURRENCIES = {
  INR: 'INR',
  USD: 'USD',
  EUR: 'EUR',
  GBP: 'GBP',
} as const

export const BANK_CODES = {
  HDFC: { code: 'HDFC', name: 'HDFC Bank' },
  ICIC: { code: 'ICIC', name: 'ICICI Bank' },
  AXIS: { code: 'AXIS', name: 'Axis Bank' },
  SBIN: { code: 'SBIN', name: 'State Bank of India' },
  UTIB: { code: 'UTIB', name: 'Axis Bank' },
  KOBA: { code: 'KOBA', name: 'Kotak Mahindra Bank' },
  IDFB: { code: 'IDFB', name: 'IDFC Bank' },
  AUBL: { code: 'AUBL', name: 'Aurobindo Bank' },
} as const

export const STATUS_COLORS = {
  [PAYMENT_STATUS.PENDING]: 'bg-yellow-100 text-yellow-800',
  [PAYMENT_STATUS.INITIATED]: 'bg-blue-100 text-blue-800',
  [PAYMENT_STATUS.PROCESSING]: 'bg-blue-100 text-blue-800',
  [PAYMENT_STATUS.SUCCESSFUL]: 'bg-green-100 text-green-800',
  [PAYMENT_STATUS.FAILED]: 'bg-red-100 text-red-800',
  [PAYMENT_STATUS.CANCELLED]: 'bg-gray-100 text-gray-800',
  [PAYMENT_STATUS.REFUNDED]: 'bg-purple-100 text-purple-800',
} as const

export const STATUS_ICONS = {
  [PAYMENT_STATUS.PENDING]: '⏳',
  [PAYMENT_STATUS.INITIATED]: '🔄',
  [PAYMENT_STATUS.PROCESSING]: '⚙️',
  [PAYMENT_STATUS.SUCCESSFUL]: '✅',
  [PAYMENT_STATUS.FAILED]: '❌',
  [PAYMENT_STATUS.CANCELLED]: '⛔',
  [PAYMENT_STATUS.REFUNDED]: '↩️',
} as const

export const METHOD_ICONS = {
  [PAYMENT_METHODS.NETBANKING]: '🏦',
  [PAYMENT_METHODS.CREDIT_CARD]: '💳',
  [PAYMENT_METHODS.DEBIT_CARD]: '💳',
  [PAYMENT_METHODS.UPI]: '📱',
  [PAYMENT_METHODS.WALLET]: '👛',
} as const

export const PROVIDER_COLORS = {
  [PAYMENT_PROVIDERS.RAZORPAY]: 'bg-blue-50 border-blue-200',
  [PAYMENT_PROVIDERS.BILLDESK]: 'bg-orange-50 border-orange-200',
} as const

// Format currency for display
export const formatCurrency = (amount: number, currency: string = 'INR'): string => {
  const formatter = new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
  return formatter.format(amount)
}

// Format date
export const formatDate = (date: string | Date): string => {
  const d = typeof date === 'string' ? new Date(date) : date
  return new Intl.DateTimeFormat('en-IN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d)
}

// Validate payment amount
export const validateAmount = (amount: string): { valid: boolean; error?: string } => {
  const num = parseFloat(amount)
  
  if (!amount || isNaN(num)) {
    return { valid: false, error: 'Please enter a valid amount' }
  }
  
  if (num <= 0) {
    return { valid: false, error: 'Amount must be greater than 0' }
  }
  
  if (num > 999999999) {
    return { valid: false, error: 'Amount is too large' }
  }
  
  return { valid: true }
}

// Validate email
export const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

// Validate phone
export const validatePhone = (phone: string): boolean => {
  const phoneRegex = /^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$/
  return phoneRegex.test(phone)
}

// Get payment method display name
export const getPaymentMethodName = (method: string): string => {
  const names: { [key: string]: string } = {
    netbanking: 'Netbanking',
    credit_card: 'Credit Card',
    debit_card: 'Debit Card',
    upi: 'UPI',
    wallet: 'Digital Wallet',
  }
  return names[method] || method
}

// Get status display name
export const getStatusName = (status: string): string => {
  const names: { [key: string]: string } = {
    pending: 'Pending',
    initiated: 'Initiated',
    processing: 'Processing',
    successful: 'Successful',
    failed: 'Failed',
    cancelled: 'Cancelled',
    refunded: 'Refunded',
  }
  return names[status] || status
}

// Check if payment can be refunded
export const canRefund = (status: string): boolean => {
  return status === PAYMENT_STATUS.SUCCESSFUL
}

// Check if payment is final
export const isFinalStatus = (status: 'successful' | 'failed' | 'cancelled' | 'refunded'): boolean => {
  return [
    PAYMENT_STATUS.SUCCESSFUL,
    PAYMENT_STATUS.FAILED,
    PAYMENT_STATUS.CANCELLED,
    PAYMENT_STATUS.REFUNDED,
  ].includes(status)
}

// Calculate success rate
export const calculateSuccessRate = (successful: number, total: number): number => {
  if (total === 0) return 0
  return Math.round((successful / total) * 100 * 10) / 10
}

// Generate order ID
export const generateOrderID = (): string => {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 8).toUpperCase()
  return `ORDER-${timestamp}-${random}`
}

// Format payment summary
export const formatPaymentSummary = (summary: any) => {
  return {
    totalPayments: summary.total_payments || 0,
    totalSuccessful: summary.total_successful || 0,
    totalFailed: summary.total_failed || 0,
    totalPending: summary.total_pending || 0,
    successfulCount: summary.successful_count || 0,
    failedCount: summary.failed_count || 0,
    pendingCount: summary.pending_count || 0,
    successRate: summary.success_rate || 0,
  }
}

// Error messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  PAYMENT_FAILED: 'Payment failed. Please try again.',
  INVALID_AMOUNT: 'Please enter a valid amount.',
  INVALID_EMAIL: 'Please enter a valid email address.',
  INVALID_PHONE: 'Please enter a valid phone number.',
  GATEWAY_ERROR: 'Payment gateway error. Please try again.',
  SIGNATURE_MISMATCH: 'Payment signature verification failed.',
  PAYMENT_EXPIRED: 'Payment link has expired.',
  ALREADY_PROCESSED: 'Payment has already been processed.',
}

// Success messages
export const SUCCESS_MESSAGES = {
  PAYMENT_INITIATED: 'Payment initiated successfully. Redirecting to payment gateway...',
  PAYMENT_VERIFIED: 'Payment verified successfully.',
  REFUND_INITIATED: 'Refund initiated successfully.',
  GATEWAY_CONFIGURED: 'Payment gateway configured successfully.',
}
