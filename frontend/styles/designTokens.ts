/**
 * Design Tokens for Call Center Dashboard
 * Based on 12 Design Tips for Dashboard Design
 * Enforces consistency, hierarchy, and data-driven design
 */

// ============================================================================
// COLOR SYSTEM - Semantic meaning (Tip 7: Maintain Consistency)
// ============================================================================

export const COLORS = {
  // Status Colors - Always mean the same thing
  success: '#10b981',      // Green - good, achieved
  warning: '#f59e0b',      // Orange - needs attention
  critical: '#ef4444',     // Red - urgent action needed
  neutral: '#6b7280',      // Gray - informational only

  // Metric Type Colors - Consistent across dashboard
  performance: '#3b82f6',  // Blue - efficiency metrics
  quality: '#8b5cf6',      // Purple - quality metrics
  volume: '#06b6d4',       // Cyan - volume/count metrics
  revenue: '#14b8a6',      // Teal - revenue metrics
  engagement: '#ec4899',   // Pink - engagement/gamification

  // UI Colors
  background: '#ffffff',
  backgroundAlt: '#f9fafb', // Tip 3: Consider Data Ink Ratio
  border: '#e5e7eb',
  borderLight: '#f3f4f6',
  text: '#111827',
  textSecondary: '#6b7280',
  textTertiary: '#9ca3af',
  textInverted: '#ffffff',
} as const

// ============================================================================
// RARITY SYSTEM (Gamification)
// ============================================================================

export const RARITY_COLORS = {
  common: {
    bg: '#f3f4f6',
    text: '#374151',
    border: '#d1d5db',
    accent: '#6b7280',
  },
  uncommon: {
    bg: '#dcfce7',
    text: '#166534',
    border: '#86efac',
    accent: '#22c55e',
  },
  rare: {
    bg: '#dbeafe',
    text: '#0c4a6e',
    border: '#7dd3fc',
    accent: '#3b82f6',
  },
  epic: {
    bg: '#f3e8ff',
    text: '#6b21a8',
    border: '#d8b4fe',
    accent: '#a855f7',
  },
  legendary: {
    bg: '#fef3c7',
    text: '#78350f',
    border: '#fcd34d',
    accent: '#f59e0b',
  },
} as const

// ============================================================================
// TYPOGRAPHY - Consistent Font Scales (Tip 7: Maintain Consistency)
// ============================================================================

export const TYPOGRAPHY = {
  // Size scales (Tip 8: Use Size & Position for Hierarchy)
  sizes: {
    xs: '12px',     // Supporting details
    sm: '14px',     // Labels, captions
    base: '16px',   // Body text
    lg: '18px',     // Card titles
    xl: '24px',     // Section headers
    '2xl': '30px',  // Important metrics (Tier 2)
    '3xl': '36px',  // Critical metrics (Tier 1)
    '4xl': '48px',  // Hero metrics (Tier 1)
  },

  weights: {
    normal: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },

  // Component-specific
  metric: {
    label: { size: '12px', weight: 500, color: COLORS.textSecondary },
    value: { size: '36px', weight: 700, color: COLORS.text },
    subtext: { size: '12px', weight: 400, color: COLORS.textTertiary },
  },

  card: {
    title: { size: '18px', weight: 600, color: COLORS.text },
    text: { size: '14px', weight: 400, color: COLORS.textSecondary },
  },

  header: {
    title: { size: '30px', weight: 700, color: COLORS.text },
    subtitle: { size: '14px', weight: 400, color: COLORS.textSecondary },
  },
} as const

// ============================================================================
// SPACING - Data Ink Ratio Optimization (Tip 3)
// ============================================================================

export const SPACING = {
  xs: '4px',
  sm: '8px',
  md: '12px',
  lg: '16px',
  xl: '24px',
  '2xl': '32px',
  '3xl': '48px',
} as const

// ============================================================================
// LAYOUT HIERARCHY (Tip 8: Use Size & Position for Hierarchy)
// ============================================================================

export const HIERARCHY = {
  tier1: {
    // Critical metrics - Largest, top, bold
    fontSize: '36px',
    fontWeight: 700,
    padding: '24px',
    background: 'from-blue-600 to-purple-600',
    textColor: COLORS.textInverted,
    gridColumns: '2',
    marginBottom: '32px',
  },

  tier2: {
    // Important metrics - Medium, clear
    fontSize: '24px',
    fontWeight: 600,
    padding: '16px',
    background: COLORS.background,
    border: `1px solid ${COLORS.border}`,
    textColor: COLORS.text,
    gridColumns: '4',
    marginBottom: '24px',
  },

  tier3: {
    // Supporting metrics - Smaller, detailed
    fontSize: '16px',
    fontWeight: 500,
    padding: '12px',
    background: COLORS.background,
    border: `1px solid ${COLORS.borderLight}`,
    textColor: COLORS.text,
    gridColumns: '3',
  },
} as const

// ============================================================================
// DATA INK RATIO - Minimize decorative elements (Tip 3)
// ============================================================================

export const DATA_INK = {
  // Remove unnecessary borders - use whitespace instead
  borders: {
    subtle: `1px solid ${COLORS.borderLight}`,
    standard: `1px solid ${COLORS.border}`,
    accent: `2px solid ${COLORS.performance}`,
    critical: `2px solid ${COLORS.critical}`,
  },

  // Minimal shadows for depth
  shadows: {
    none: 'none',
    subtle: '0 1px 2px rgba(0,0,0,0.05)',
    standard: '0 1px 3px rgba(0,0,0,0.1)',
  },

  // Remove decorative icons - use icons only for meaning
  icons: {
    show: true,  // Use icons for status/meaning only
    decorative: false, // Avoid emoji/decorative elements
  },
} as const

// ============================================================================
// ROUNDING - Number Formatting Rules (Tip 4)
// ============================================================================

export const NUMBER_FORMAT = {
  duration: {
    format: (seconds: number) => {
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      const secs = Math.floor(seconds % 60)
      
      if (hours > 0) return `${hours}h ${minutes}m`
      if (minutes > 0) return `${minutes}m ${secs}s`
      return `${secs}s`
    },
  },

  largeNumber: {
    format: (num: number) => {
      if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
      if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
      return num.toString()
    },
  },

  percentage: {
    format: (num: number, decimals = 0) => num.toFixed(decimals) + '%',
  },

  rating: {
    format: (num: number) => num.toFixed(1),
  },

  money: {
    format: (num: number) => '$' + num.toLocaleString('en-US', { maximumFractionDigits: 0 }),
  },
} as const

// ============================================================================
// CONTEXT & THRESHOLDS - Give numbers meaning (Tip 9)
// ============================================================================

export const THRESHOLDS = {
  queueWaitTime: {
    target: 180, // seconds (3 minutes)
    warning: 300, // seconds (5 minutes)
    critical: 600, // seconds (10 minutes)
  },

  agentProductivity: {
    target: 30, // calls per day
    warning: 20,
    critical: 10,
  },

  customerSatisfaction: {
    target: 4.5,
    warning: 4.0,
    critical: 3.5,
  },

  callCompletionRate: {
    target: 95,
    warning: 85,
    critical: 75,
  },

  firstCallResolution: {
    target: 85,
    warning: 75,
    critical: 60,
  },
} as const

// ============================================================================
// RESPONSIVE BREAKPOINTS
// ============================================================================

export const BREAKPOINTS = {
  mobile: '640px',    // sm
  tablet: '768px',    // md
  desktop: '1024px',  // lg
  wide: '1280px',     // xl
} as const

// ============================================================================
// ANIMATION & TRANSITIONS
// ============================================================================

export const ANIMATION = {
  fast: '150ms ease-out',
  standard: '300ms ease-out',
  slow: '500ms ease-out',
} as const

// ============================================================================
// GRID LAYOUTS
// ============================================================================

export const GRID = {
  tier1: 'grid grid-cols-1 md:grid-cols-2 gap-6',
  tier2: 'grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4',
  tier3: 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4',
} as const

// ============================================================================
// COMPONENT PRESETS
// ============================================================================

export const COMPONENT_PRESETS = {
  metricCard: {
    padding: SPACING.lg,
    border: DATA_INK.borders.standard,
    borderRadius: '8px',
    background: COLORS.background,
  },

  alertBox: {
    padding: SPACING.lg,
    borderRadius: '8px',
    borderLeft: DATA_INK.borders.critical,
  },

  button: {
    primary: {
      background: COLORS.performance,
      text: COLORS.textInverted,
      padding: `${SPACING.md} ${SPACING.lg}`,
      borderRadius: '6px',
    },
    secondary: {
      background: COLORS.backgroundAlt,
      text: COLORS.text,
      padding: `${SPACING.md} ${SPACING.lg}`,
      borderRadius: '6px',
    },
  },
} as const

export default {
  COLORS,
  RARITY_COLORS,
  TYPOGRAPHY,
  SPACING,
  HIERARCHY,
  DATA_INK,
  NUMBER_FORMAT,
  THRESHOLDS,
  BREAKPOINTS,
  ANIMATION,
  GRID,
  COMPONENT_PRESETS,
}
