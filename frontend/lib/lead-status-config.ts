/**
 * Lead Status Management System
 * 
 * This file defines the comprehensive lead lifecycle stages and status mapping
 * for the VYOM ERP sales funnel.
 */

// Detailed Lead Status Types
export type DetailedLeadStatus = 
  // Initial Contact Phase
  | 'Fresh Lead'
  | 'Re Engaged'
  | 'Not Connected'
  | 'Dead'
  | 'Follow Up - Cold'
  | 'Follow Up - Warm'
  | 'Follow Up - Hot'
  | 'Lost'
  | 'Unqualified (Location)'
  | 'Unqualified (Budget)'
  | 'Unqualified - Client Profile'
  // Site Visit Phase
  | 'SV - Scheduled'
  | 'SV - Done'
  | 'SV - Cold'
  | 'SV - Warm'
  | 'SV - Revisit Scheduled'
  | 'SV - Revisit Done'
  | 'SV - Hot'
  | 'SV - Lost (No Response)'
  | 'SV - Lost (Budget)'
  | 'SV - Lost (Plan Dropped)'
  | 'SV - Lost (Location)'
  | 'SV - Lost (Availability)'
  // Face to Face Phase
  | 'F2F - Scheduled'
  | 'F2F - Done'
  | 'F2F - Follow Up'
  | 'F2F - Warm'
  | 'F2F - Hot'
  | 'F2F - Lost (No Response)'
  | 'F2F - Lost (Budget)'
  | 'F2F - Lost (Plan Dropped)'
  | 'F2F - Lost (Location)'
  | 'F2F - Lost (Availability)'
  // Booking Phase
  | 'Booking - In Progress'
  | 'Booking Progress - Lost'
  | 'Booking Done'

// Pipeline Status (Simplified categorization)
export type PipelineStatus = 
  | 'New'
  | 'Connected'
  | 'Interested'
  | 'Analysis'
  | 'Negotiation'
  | 'Pre Booking'
  | 'Booking'

// Next Action
export type NextAction = 
  | 'Capture Date A'
  | 'Capture Date B'
  | 'Capture Date C'
  | 'Capture Date D'
  | 'Transferred to Pre Sales'
  | 'Transferred to Sales'
  | 'Transferred to Post Sales'
  | 'Call back'
  | 'Follow-up'
  | 'Lost'
  | 'RNR' // No Response Required

/**
 * Comprehensive status mapping
 * Maps detailed statuses to pipeline stage and next action
 */
export const STATUS_MAP: Record<DetailedLeadStatus, {
  pipeline: PipelineStatus
  action: NextAction
  color: string
  icon: string
  phase: string
}> = {
  // Initial Contact Phase
  'Fresh Lead': {
    pipeline: 'New',
    action: 'Transferred to Pre Sales',
    color: 'bg-blue-100 text-blue-800',
    icon: '🆕',
    phase: 'Initial Contact'
  },
  'Re Engaged': {
    pipeline: 'Connected',
    action: 'Transferred to Pre Sales',
    color: 'bg-cyan-100 text-cyan-800',
    icon: '🔄',
    phase: 'Initial Contact'
  },
  'Not Connected': {
    pipeline: 'New',
    action: 'RNR',
    color: 'bg-gray-100 text-gray-800',
    icon: '❌',
    phase: 'Initial Contact'
  },
  'Dead': {
    pipeline: 'New',
    action: 'Lost',
    color: 'bg-slate-200 text-slate-800',
    icon: '💀',
    phase: 'Initial Contact'
  },
  'Follow Up - Cold': {
    pipeline: 'Connected',
    action: 'Follow-up',
    color: 'bg-indigo-100 text-indigo-800',
    icon: '❄️',
    phase: 'Initial Contact'
  },
  'Follow Up - Warm': {
    pipeline: 'Connected',
    action: 'Follow-up',
    color: 'bg-yellow-100 text-yellow-800',
    icon: '🔥',
    phase: 'Initial Contact'
  },
  'Follow Up - Hot': {
    pipeline: 'Interested',
    action: 'Follow-up',
    color: 'bg-orange-100 text-orange-800',
    icon: '🔥🔥',
    phase: 'Initial Contact'
  },
  'Lost': {
    pipeline: 'Connected',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '❌',
    phase: 'Initial Contact'
  },
  'Unqualified (Location)': {
    pipeline: 'Connected',
    action: 'Lost',
    color: 'bg-rose-100 text-rose-800',
    icon: '📍',
    phase: 'Initial Contact'
  },
  'Unqualified (Budget)': {
    pipeline: 'Connected',
    action: 'Lost',
    color: 'bg-rose-100 text-rose-800',
    icon: '💰',
    phase: 'Initial Contact'
  },
  'Unqualified - Client Profile': {
    pipeline: 'Connected',
    action: 'Lost',
    color: 'bg-rose-100 text-rose-800',
    icon: '👤',
    phase: 'Initial Contact'
  },

  // Site Visit Phase
  'SV - Scheduled': {
    pipeline: 'Interested',
    action: 'Call back',
    color: 'bg-purple-100 text-purple-800',
    icon: '📅',
    phase: 'Site Visit'
  },
  'SV - Done': {
    pipeline: 'Analysis',
    action: 'Transferred to Sales',
    color: 'bg-green-100 text-green-800',
    icon: '✅',
    phase: 'Site Visit'
  },
  'SV - Cold': {
    pipeline: 'Analysis',
    action: 'Call back',
    color: 'bg-indigo-100 text-indigo-800',
    icon: '❄️',
    phase: 'Site Visit'
  },
  'SV - Warm': {
    pipeline: 'Analysis',
    action: 'Call back',
    color: 'bg-yellow-100 text-yellow-800',
    icon: '🔥',
    phase: 'Site Visit'
  },
  'SV - Revisit Scheduled': {
    pipeline: 'Analysis',
    action: 'Call back',
    color: 'bg-purple-100 text-purple-800',
    icon: '📅',
    phase: 'Site Visit'
  },
  'SV - Revisit Done': {
    pipeline: 'Analysis',
    action: 'Capture Date C',
    color: 'bg-green-100 text-green-800',
    icon: '✅',
    phase: 'Site Visit'
  },
  'SV - Hot': {
    pipeline: 'Negotiation',
    action: 'Call back',
    color: 'bg-orange-100 text-orange-800',
    icon: '🔥🔥',
    phase: 'Site Visit'
  },
  'SV - Lost (No Response)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '❌',
    phase: 'Site Visit'
  },
  'SV - Lost (Budget)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '💰',
    phase: 'Site Visit'
  },
  'SV - Lost (Plan Dropped)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '📉',
    phase: 'Site Visit'
  },
  'SV - Lost (Location)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '📍',
    phase: 'Site Visit'
  },
  'SV - Lost (Availability)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '⏰',
    phase: 'Site Visit'
  },

  // Face to Face Phase
  'F2F - Scheduled': {
    pipeline: 'Interested',
    action: 'Call back',
    color: 'bg-purple-100 text-purple-800',
    icon: '📅',
    phase: 'Face to Face'
  },
  'F2F - Done': {
    pipeline: 'Analysis',
    action: 'Transferred to Sales',
    color: 'bg-green-100 text-green-800',
    icon: '✅',
    phase: 'Face to Face'
  },
  'F2F - Follow Up': {
    pipeline: 'Analysis',
    action: 'Call back',
    color: 'bg-blue-100 text-blue-800',
    icon: '📞',
    phase: 'Face to Face'
  },
  'F2F - Warm': {
    pipeline: 'Analysis',
    action: 'Call back',
    color: 'bg-yellow-100 text-yellow-800',
    icon: '🔥',
    phase: 'Face to Face'
  },
  'F2F - Hot': {
    pipeline: 'Negotiation',
    action: 'Call back',
    color: 'bg-orange-100 text-orange-800',
    icon: '🔥🔥',
    phase: 'Face to Face'
  },
  'F2F - Lost (No Response)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '❌',
    phase: 'Face to Face'
  },
  'F2F - Lost (Budget)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '💰',
    phase: 'Face to Face'
  },
  'F2F - Lost (Plan Dropped)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '📉',
    phase: 'Face to Face'
  },
  'F2F - Lost (Location)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '📍',
    phase: 'Face to Face'
  },
  'F2F - Lost (Availability)': {
    pipeline: 'Analysis',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '⏰',
    phase: 'Face to Face'
  },

  // Booking Phase
  'Booking - In Progress': {
    pipeline: 'Pre Booking',
    action: 'Call back',
    color: 'bg-amber-100 text-amber-800',
    icon: '📋',
    phase: 'Booking'
  },
  'Booking Progress - Lost': {
    pipeline: 'Pre Booking',
    action: 'Lost',
    color: 'bg-red-100 text-red-800',
    icon: '❌',
    phase: 'Booking'
  },
  'Booking Done': {
    pipeline: 'Booking',
    action: 'Transferred to Post Sales',
    color: 'bg-emerald-100 text-emerald-800',
    icon: '🎉',
    phase: 'Booking'
  }
}

/**
 * Get pipeline stage order for funnel view
 */
export const PIPELINE_ORDER: PipelineStatus[] = [
  'New',
  'Connected',
  'Interested',
  'Analysis',
  'Negotiation',
  'Pre Booking',
  'Booking'
]

/**
 * Get all statuses for a specific phase
 */
export const getStatusesByPhase = (phase: string): DetailedLeadStatus[] => {
  return Object.entries(STATUS_MAP)
    .filter(([_, config]) => config.phase === phase)
    .map(([status]) => status as DetailedLeadStatus)
}

/**
 * Get all available phases
 */
export const getPhases = (): string[] => {
  const phases = new Set(Object.values(STATUS_MAP).map(config => config.phase))
  return Array.from(phases)
}

/**
 * Utility function to get status info
 */
export const getStatusInfo = (status: DetailedLeadStatus) => {
  return STATUS_MAP[status] || STATUS_MAP['Fresh Lead']
}

/**
 * Utility function to get color for pipeline status
 */
export const getPipelineColor = (pipeline: PipelineStatus): string => {
  const colors: Record<PipelineStatus, string> = {
    'New': 'bg-blue-50 border-blue-200',
    'Connected': 'bg-cyan-50 border-cyan-200',
    'Interested': 'bg-purple-50 border-purple-200',
    'Analysis': 'bg-yellow-50 border-yellow-200',
    'Negotiation': 'bg-orange-50 border-orange-200',
    'Pre Booking': 'bg-amber-50 border-amber-200',
    'Booking': 'bg-emerald-50 border-emerald-200'
  }
  return colors[pipeline]
}
