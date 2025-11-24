import type { Metadata, Viewport } from 'next'
import './globals.css'
import { ToasterProvider } from '@/components/providers/ToasterProvider'
import { AuthProvider } from '@/components/providers/AuthProvider'
import { TenantWrapper } from '@/components/providers/TenantProvider'
import { TenantManagementProvider } from '@/contexts/TenantManagementContext'

// Enable caching for this layout
export const dynamic = 'force-dynamic'
export const revalidate = 3600 // Revalidate every hour

export const metadata: Metadata = {
  title: 'AI Call Center',
  description: 'Multi-Tenant AI Call Center Management System',
}

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <TenantWrapper>
            <TenantManagementProvider>
              <ToasterProvider />
              {children}
            </TenantManagementProvider>
          </TenantWrapper>
        </AuthProvider>
      </body>
    </html>
  )
}
