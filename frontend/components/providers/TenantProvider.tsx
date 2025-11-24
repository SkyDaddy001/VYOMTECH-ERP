'use client'

import React from 'react'
import { TenantProvider } from '@/contexts/TenantContext'

export function TenantWrapper({ children }: { children: React.ReactNode }) {
  return <TenantProvider>{children}</TenantProvider>
}
