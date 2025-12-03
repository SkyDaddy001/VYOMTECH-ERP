'use client'

import { usePathname } from 'next/navigation'
import Link from 'next/link'

interface BreadcrumbItem {
  label: string
  href?: string
}

export function Breadcrumbs() {
  const pathname = usePathname()
  
  // Generate breadcrumb items from pathname
  const generateBreadcrumbs = (): BreadcrumbItem[] => {
    const items: BreadcrumbItem[] = [{ label: '🏠 Dashboard', href: '/dashboard' }]
    
    if (pathname === '/dashboard') return items
    
    const pathSegments = pathname.split('/').filter(Boolean)
    let currentPath = ''
    
    for (let i = 1; i < pathSegments.length; i++) {
      currentPath += `/${pathSegments[i]}`
      const label = pathSegments[i]
        .split('-')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ')
      
      if (i === pathSegments.length - 1) {
        items.push({ label })
      } else {
        items.push({ label, href: currentPath })
      }
    }
    
    return items
  }

  const breadcrumbs = generateBreadcrumbs()

  return (
    <nav className="flex items-center gap-2 text-sm mb-6 overflow-x-auto">
      {breadcrumbs.map((item, index) => (
        <div key={index} className="flex items-center gap-2 whitespace-nowrap">
          {item.href ? (
            <Link href={item.href} className="text-blue-600 hover:text-blue-800 transition">
              {item.label}
            </Link>
          ) : (
            <span className="text-gray-600 font-medium">{item.label}</span>
          )}
          {index < breadcrumbs.length - 1 && <span className="text-gray-400">/</span>}
        </div>
      ))}
    </nav>
  )
}
