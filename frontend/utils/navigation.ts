// Navigation Utilities and Helpers
import { ReactNode } from 'react'

export interface NavigationItem {
  href: string
  label: string
  icon?: string
  description?: string
  subItems?: NavigationItem[]
  badge?: string | number
}

export interface NavigationCategory {
  name: string
  items: NavigationItem[]
}

/**
 * Complete navigation structure for the application
 * Used for sidebar, navigation menus, and site map
 */
export const NAVIGATION_STRUCTURE: NavigationCategory[] = [
  {
    name: 'Core Modules',
    items: [
      {
        href: '/dashboard',
        label: 'Dashboard',
        icon: '📊',
        description: 'Overview & KPIs',
      },
      {
        href: '/dashboard/sales',
        label: 'Sales',
        icon: '📈',
        description: 'Sales pipeline & management',
      },
      {
        href: '/dashboard/presales',
        label: 'Pre-Sales',
        icon: '🎯',
        description: 'Pre-sales activities',
      },
      {
        href: '/dashboard/leads',
        label: 'Leads',
        icon: '🔍',
        description: 'Lead management',
      },
      {
        href: '/dashboard/accounts',
        label: 'Finance',
        icon: '💰',
        description: 'GL & accounting',
      },
      {
        href: '/dashboard/ledgers',
        label: 'Ledgers',
        icon: '📑',
        description: 'Account ledgers',
      },
    ],
  },
  {
    name: 'Operations',
    items: [
      {
        href: '/dashboard/purchase',
        label: 'Purchase',
        icon: '📦',
        description: 'Vendor & PO management',
      },
      {
        href: '/dashboard/hr',
        label: 'HR',
        icon: '👨‍💼',
        description: 'Human resources',
      },
      {
        href: '/dashboard/projects',
        label: 'Projects',
        icon: '📌',
        description: 'Project management',
      },
      {
        href: '/dashboard/workflows',
        label: 'Workflows',
        icon: '⚙️',
        description: 'Automation & workflows',
      },
      {
        href: '/dashboard/reports',
        label: 'Reports',
        icon: '📋',
        description: 'Business reports',
      },
    ],
  },
  {
    name: 'Real Estate & Construction',
    items: [
      {
        href: '/dashboard/real-estate',
        label: 'Real Estate',
        icon: '🏢',
        description: 'Properties & bookings',
      },
      {
        href: '/dashboard/construction',
        label: 'Construction',
        icon: '🏗️',
        description: 'Construction tracking',
      },
      {
        href: '/dashboard/civil',
        label: 'Civil',
        icon: '🏛️',
        description: 'Civil engineering',
      },
      {
        href: '/dashboard/units',
        label: 'Units',
        icon: '🏠',
        description: 'Unit management',
      },
    ],
  },
  {
    name: 'Marketing & Communications',
    items: [
      {
        href: '/dashboard/marketing',
        label: 'Marketing',
        icon: '📣',
        description: 'Marketing campaigns',
      },
      {
        href: '/dashboard/campaigns',
        label: 'Campaigns',
        icon: '🎪',
        description: 'Campaign management',
      },
      {
        href: '/dashboard/calls',
        label: 'Calls',
        icon: '☎️',
        description: 'Call management',
      },
      {
        href: '/dashboard/agents',
        label: 'Agents',
        icon: '📞',
        description: 'Agent management',
      },
    ],
  },
  {
    name: 'Administration',
    items: [
      {
        href: '/dashboard/users',
        label: 'Users',
        icon: '👥',
        description: 'User management',
      },
      {
        href: '/dashboard/tenants',
        label: 'Tenants',
        icon: '🏢',
        description: 'Tenant management',
      },
      {
        href: '/dashboard/company',
        label: 'Company',
        icon: '🏛️',
        description: 'Company settings',
      },
      {
        href: '/dashboard/bookings',
        label: 'Bookings',
        icon: '📅',
        description: 'Booking management',
      },
    ],
  },
]

/**
 * Quick access items - frequently used navigation links
 */
export const QUICK_ACCESS_ITEMS: NavigationItem[] = [
  { href: '/dashboard', label: 'Dashboard', icon: '📊' },
  { href: '/dashboard/sales', label: 'Sales', icon: '📈' },
  { href: '/dashboard/leads', label: 'Leads', icon: '🔍' },
  { href: '/dashboard/calls', label: 'Calls', icon: '☎️' },
  { href: '/dashboard/projects', label: 'Projects', icon: '📌' },
  { href: '/dashboard/users', label: 'Users', icon: '👥' },
  { href: '/dashboard/reports', label: 'Reports', icon: '📋' },
  { href: '/dashboard/workflows', label: 'Workflows', icon: '⚙️' },
]

/**
 * Get breadcrumb path from URL pathname
 */
export function getBreadcrumbPath(pathname: string): Array<{ label: string; href?: string }> {
  const items: Array<{ label: string; href?: string }> = [
    { label: '🏠 Dashboard', href: '/dashboard' },
  ]

  if (pathname === '/dashboard') return items

  const pathSegments = pathname.split('/').filter(Boolean)
  let currentPath = ''

  for (let i = 1; i < pathSegments.length; i++) {
    currentPath += `/${pathSegments[i]}`
    const label = formatPathSegment(pathSegments[i])

    if (i === pathSegments.length - 1) {
      items.push({ label })
    } else {
      items.push({ label, href: currentPath })
    }
  }

  return items
}

/**
 * Format URL segment to readable label
 */
function formatPathSegment(segment: string): string {
  return segment
    .split('-')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

/**
 * Find navigation item by href
 */
export function findNavigationItem(href: string): NavigationItem | undefined {
  for (const category of NAVIGATION_STRUCTURE) {
    const found = category.items.find((item) => item.href === href)
    if (found) return found
  }
  return undefined
}

/**
 * Get all navigation items flattened
 */
export function getAllNavigationItems(): NavigationItem[] {
  return NAVIGATION_STRUCTURE.flatMap((category) => category.items)
}

/**
 * Search navigation items by keyword
 */
export function searchNavigationItems(keyword: string): NavigationItem[] {
  const lowerKeyword = keyword.toLowerCase()
  return getAllNavigationItems().filter(
    (item) =>
      item.label.toLowerCase().includes(lowerKeyword) ||
      item.description?.toLowerCase().includes(lowerKeyword)
  )
}
