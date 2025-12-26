// This is the Prisma schema addition for tenant credentials
// Add this to your prisma/schema.prisma file

// ============= ADD THIS SECTION TO schema.prisma =============

model TenantCredential {
  id                 String    @id @default(cuid())
  tenantId           String
  tenant             Tenant    @relation(fields: [tenantId], references: [id], onDelete: Cascade)
  
  // Type of credential (GOOGLE_OAUTH, META_OAUTH, AWS_S3, EMAIL_SMTP, RAZORPAY, BILLDESK, etc.)
  credentialType     String    @db.VarChar(50)
  
  // Encrypted JSON containing the actual credential data
  // Different types will have different fields
  encryptedValue     String    @db.Text
  
  // Metadata
  description        String?   @db.VarChar(255)
  isActive           Boolean   @default(true)
  lastRotatedAt      DateTime?
  expiresAt          DateTime?
  
  // Audit
  createdAt          DateTime  @default(now())
  updatedAt          DateTime  @updatedAt
  createdBy          String?
  updatedBy          String?
  
  // Ensure one active credential per type per tenant
  @@unique([tenantId, credentialType], where: { isActive: true })
  @@index([tenantId])
  @@index([credentialType])
  @@map("tenant_credentials")
}

// Also add this relation to your existing Tenant model:
// Add to Tenant model in schema.prisma:
// credentials    TenantCredential[]

// ============= END OF ADDITION =============
