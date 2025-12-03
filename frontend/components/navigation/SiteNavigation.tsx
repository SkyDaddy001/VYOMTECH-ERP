'use client'

import Link from 'next/link'

export function SiteNavigation() {
  const navigationMap = {
    'Core Modules': [
      { href: '/dashboard', label: '📊 Dashboard', desc: 'Overview & KPIs' },
      { href: '/dashboard/sales', label: '📈 Sales', desc: 'Sales pipeline & management' },
      { href: '/dashboard/presales', label: '🎯 Pre-Sales', desc: 'Pre-sales activities' },
      { href: '/dashboard/leads', label: '🔍 Leads', desc: 'Lead management' },
      { href: '/dashboard/accounts', label: '💰 Finance', desc: 'GL & accounting' },
      { href: '/dashboard/ledgers', label: '📑 Ledgers', desc: 'Account ledgers' },
    ],
    'Operations': [
      { href: '/dashboard/purchase', label: '📦 Purchase', desc: 'Vendor & PO management' },
      { href: '/dashboard/hr', label: '👨‍💼 HR', desc: 'Human resources' },
      { href: '/dashboard/projects', label: '📌 Projects', desc: 'Project management' },
      { href: '/dashboard/workflows', label: '⚙️ Workflows', desc: 'Automation & workflows' },
      { href: '/dashboard/reports', label: '📋 Reports', desc: 'Business reports' },
    ],
    'Real Estate & Construction': [
      { href: '/dashboard/real-estate', label: '🏢 Real Estate', desc: 'Properties & bookings' },
      { href: '/dashboard/construction', label: '🏗️ Construction', desc: 'Construction tracking' },
      { href: '/dashboard/civil', label: '🏛️ Civil', desc: 'Civil engineering' },
      { href: '/dashboard/units', label: '🏠 Units', desc: 'Unit management' },
    ],
    'Marketing & Comms': [
      { href: '/dashboard/marketing', label: '📣 Marketing', desc: 'Marketing campaigns' },
      { href: '/dashboard/campaigns', label: '🎪 Campaigns', desc: 'Campaign management' },
      { href: '/dashboard/calls', label: '☎️ Calls', desc: 'Call management' },
      { href: '/dashboard/agents', label: '📞 Agents', desc: 'Agent management' },
    ],
    'Administration': [
      { href: '/dashboard/users', label: '👥 Users', desc: 'User management' },
      { href: '/dashboard/tenants', label: '🏢 Tenants', desc: 'Tenant management' },
      { href: '/dashboard/company', label: '🏛️ Company', desc: 'Company settings' },
      { href: '/dashboard/bookings', label: '📅 Bookings', desc: 'Booking management' },
    ],
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold mb-6 text-gray-900">System Navigation</h2>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {Object.entries(navigationMap).map(([category, items]) => (
          <div key={category}>
            <h3 className="font-semibold text-gray-900 mb-4 text-lg border-b-2 border-blue-600 pb-2">{category}</h3>
            <div className="space-y-3">
              {items.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="block p-3 rounded-lg hover:bg-blue-50 transition border border-gray-100 hover:border-blue-300"
                >
                  <div className="font-medium text-gray-900">{item.label}</div>
                  <div className="text-xs text-gray-600 mt-1">{item.desc}</div>
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
