# Compliance & Security Features Implementation

## Overview
This document outlines the comprehensive Compliance & Security features implemented for the Multi-Tenant AI Call Center system, including RBAC, Audit Logging, Data Encryption, and GDPR Compliance.

## Features Implemented

### 1. Role-Based Access Control (RBAC) Service
**File:** `internal/services/rbac.go`

#### Key Components:
- **Role Management**: Create, read, update, delete roles with permission sets
- **Permission Model**: Granular permission system with codes like "leads.view", "calls.edit"
- **User-Role Assignment**: Map users to roles
- **Permission Verification**: Check if user has specific permission

#### Default Roles:
- **Admin**: Full system access to all operations
- **Manager**: Team and campaign management capabilities
- **Agent**: Day-to-day operations (leads, calls)
- **Supervisor**: Monitoring and QA access

#### Key Methods:
```go
CreateRole(ctx, role) - Create role with permissions
GetRole(ctx, tenantID, roleID) - Retrieve role
ListRoles(ctx, tenantID) - List all tenant roles
UpdateRole(ctx, role) - Update role and permissions
DeleteRole(ctx, tenantID, roleID) - Soft-delete role
GetUserPermissions(ctx, tenantID, userID) - Get all user permissions
HasPermission(ctx, tenantID, userID, permissionCode) - Check specific permission
AssignRoleToUser(ctx, tenantID, userID, roleID) - Assign role
RemoveRoleFromUser(ctx, tenantID, userID, roleID) - Remove role
SetupDefaultRoles(ctx, tenantID) - Initialize standard roles
VerifyPermission(ctx, tenantID, userID, permissionCode) - Verify with error
```

### 2. Audit Logging Service
**File:** `internal/services/audit.go`

#### Key Components:
- **Action Logging**: Log all user actions (CREATE, READ, UPDATE, DELETE, LOGIN, LOGOUT)
- **Audit Trail**: Complete history with timestamps and user information
- **Security Events**: Track security-related incidents
- **Compliance Reports**: Generate compliance and audit summary reports

#### Key Methods:
```go
LogAction(ctx, log) - Log user action
GetAuditLogs(ctx, tenantID, filters, limit, offset) - Retrieve logs with filters
GetAuditLogsByUser(ctx, tenantID, userID, limit, offset) - User-specific logs
GetAuditLogsByResource(ctx, tenantID, resource, limit, offset) - Resource-specific logs
GetAuditLogsByDateRange(ctx, tenantID, startDate, endDate, limit, offset) - Date range filtering
GetAuditSummary(ctx, tenantID, days) - Summary statistics
LogSecurityEvent(ctx, event) - Log security incident
GetSecurityEvents(ctx, tenantID, filters, limit, offset) - Retrieve security events
ResolveSecurityEvent(ctx, tenantID, eventID) - Mark event as resolved
ArchiveOldAuditLogs(ctx, tenantID, retentionDays) - Cleanup old logs
GetComplianceReport(ctx, tenantID, startDate, endDate) - Compliance report
LogUserAction(ctx, tenantID, userID, action, resource, details, ipAddress, userAgent, status) - Detailed logging
```

### 3. Data Encryption Service
**File:** `internal/services/encryption_gdpr.go`

#### Key Components:
- **AES-256-GCM Encryption**: Military-grade encryption for sensitive data
- **Field-Level Encryption**: Encrypt specific sensitive fields
- **Key Rotation**: Support for encryption key rotation
- **Secure Storage**: Nonce-based encryption with authentication

#### Key Methods:
```go
EncryptField(fieldValue) - Encrypt sensitive data
DecryptField(encryptedValue) - Decrypt sensitive data
StoreEncryptedField(ctx, tenantID, fieldName, fieldValue, resourceType, resourceID) - Store encrypted
GetEncryptedField(ctx, tenantID, fieldName, resourceType, resourceID) - Retrieve encrypted
RotateEncryptionKey(ctx, tenantID, newKey) - Rotate encryption key
```

### 4. GDPR Compliance Service
**File:** `internal/services/encryption_gdpr.go`

#### Key Components:
- **Data Access Requests**: Fulfill user right to access data
- **Data Deletion**: Right to be forgotten implementation
- **Data Portability**: Export user data in portable format
- **Consent Management**: Track user consent records
- **Privacy by Design**: GDPR-compliant data handling

#### Key Methods:
```go
CreateDataAccessRequest(ctx, tenantID, userID) - Create access request
CreateDataDeletionRequest(ctx, tenantID, userID, reason) - Create deletion request
GetGDPRRequest(ctx, tenantID, requestID) - Retrieve request
ExportUserData(ctx, tenantID, userID) - Export all user data
DeleteUserData(ctx, tenantID, userID) - Permanently delete user data
RecordConsent(ctx, tenantID, userID, consentType, given) - Record consent
GetUserConsents(ctx, tenantID, userID) - Retrieve consent records
```

### 5. Security & Compliance Middleware
**File:** `internal/middleware/rbac_security.go`

#### Middleware Components:

##### PermissionMiddleware
- Validates user has required permission
- Returns 403 if permission denied
```go
PermissionMiddleware(rbacService, requiredPermission, log)
```

##### PermissionBasedAccessMiddleware
- Restricts access to users with specific roles
- Works with RBAC service for dynamic role checking
```go
PermissionBasedAccessMiddleware(rbacService, allowedRoles, log)
```

##### AuditMiddleware
- Automatically logs all requests
- Captures status codes and user information
- Async-safe audit trail creation
```go
AuditMiddleware(auditService, log)
```

##### SecurityHeadersMiddleware
- Adds security headers to all responses:
  - X-Content-Type-Options: nosniff (MIME sniffing prevention)
  - X-XSS-Protection: 1; mode=block (XSS protection)
  - X-Frame-Options: DENY (Clickjacking protection)
  - Content-Security-Policy (CSP)
  - Strict-Transport-Security (HSTS)
  - Referrer-Policy
```go
SecurityHeadersMiddleware(log)
```

##### RateLimitMiddleware
- Prevents abuse with request rate limiting
- Per-IP address rate limiting
```go
RateLimitMiddleware(requestsPerMinute, log)
```

##### DataMaskingMiddleware
- Masks sensitive parameters in logs
- Prevents credential exposure in audit trails
```go
DataMaskingMiddleware(log)
```

##### EnforceHTTPSMiddleware
- Redirects HTTP to HTTPS in production
```go
EnforceHTTPSMiddleware(log)
```

### 6. Compliance Handler & REST API
**File:** `internal/handlers/compliance.go`

#### REST Endpoints:

##### RBAC Endpoints
```
POST   /api/v1/compliance/roles               - Create role
GET    /api/v1/compliance/roles               - List roles
```

##### Audit Log Endpoints
```
GET    /api/v1/compliance/audit-logs          - Get audit logs with filters
GET    /api/v1/compliance/audit-summary       - Get audit summary
GET    /api/v1/compliance/security-events     - Get security events
GET    /api/v1/compliance/report              - Get compliance report
```

##### GDPR Compliance Endpoints
```
POST   /api/v1/compliance/gdpr/access         - Request data access
GET    /api/v1/compliance/gdpr/export         - Export user data
POST   /api/v1/compliance/gdpr/deletion       - Request data deletion
GET    /api/v1/compliance/gdpr/consent        - Get consent records
POST   /api/v1/compliance/gdpr/consent        - Record consent
```

### 7. Compliance Models
**File:** `internal/models/compliance.go`

#### Data Models:
- **Role**: User roles with permissions
- **Permission**: Granular actions
- **RolePermission**: Role-permission mapping
- **AuditLog**: Complete audit trail
- **DataEncryption**: Encrypted field storage
- **GDPRRequest**: Data access/deletion requests
- **ConsentRecord**: User consent tracking
- **DataClassification**: Data sensitivity levels
- **SecurityEvent**: Security incident tracking

## Security Features

### Authentication & Authorization
- JWT token validation (existing)
- Multi-tenant isolation (existing)
- RBAC with granular permissions (NEW)
- Permission-based access control (NEW)
- Role-based access control (ENHANCED)

### Data Protection
- AES-256-GCM encryption for sensitive fields
- Field-level encryption capability
- Encryption key rotation support
- Secure password hashing (existing)

### Audit & Compliance
- Comprehensive audit logging
- Security event tracking
- Compliance reporting
- GDPR data subject requests
- Consent management
- Data retention policies

### Network Security
- Security headers (CSP, HSTS, X-Frame-Options)
- CORS protection
- HTTPS enforcement option
- Rate limiting framework

## Configuration

### Encryption Key Setup
```go
// Generate 32-byte key for AES-256
// Store securely in environment variables
encryptionKey := os.Getenv("ENCRYPTION_KEY") // Must be 32 bytes
encService, _ := services.NewEncryptionService(db, log, encryptionKey)
```

### Default Roles Initialization
```go
rbacService := services.NewRBACService(db, log)
rbacService.SetupDefaultRoles(ctx, tenantID)
```

### Middleware Integration
```go
// Apply audit logging to all requests
r.Use(rbac_security.AuditMiddleware(auditService, log))

// Apply security headers
r.Use(rbac_security.SecurityHeadersMiddleware(log))

// Apply specific permission checks to endpoints
protectedRoute.Use(rbac_security.PermissionMiddleware(rbacService, "leads.view", log))
```

## Database Schema Requirements

The following tables are required (ensure they exist in your database):

```sql
-- Roles and Permissions
CREATE TABLE roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  permissions JSON,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY unique_tenant_role (tenant_id, name)
);

CREATE TABLE permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  category VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE role_permissions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (role_id) REFERENCES roles(id),
  FOREIGN KEY (permission_id) REFERENCES permissions(id),
  UNIQUE KEY unique_role_permission (role_id, permission_id)
);

CREATE TABLE user_roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY unique_user_role (tenant_id, user_id, role_id)
);

-- Audit Logging
CREATE TABLE audit_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  user_id BIGINT,
  action VARCHAR(50) NOT NULL,
  resource VARCHAR(100),
  details JSON,
  ip_address VARCHAR(45),
  user_agent TEXT,
  status VARCHAR(20),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_tenant_user (tenant_id, user_id),
  INDEX idx_created_at (created_at)
);

-- Data Encryption
CREATE TABLE data_encryption (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  field_name VARCHAR(100) NOT NULL,
  field_value LONGBLOB NOT NULL,
  resource_type VARCHAR(50),
  resource_id BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_resource (tenant_id, resource_type, resource_id)
);

-- GDPR Compliance
CREATE TABLE gdpr_requests (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL,
  type VARCHAR(50) NOT NULL,
  status VARCHAR(50),
  reason TEXT,
  expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_tenant_user (tenant_id, user_id)
);

CREATE TABLE consent_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL,
  type VARCHAR(50) NOT NULL,
  given BOOLEAN,
  consented_at TIMESTAMP,
  expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_tenant_user (tenant_id, user_id)
);

-- Security Events
CREATE TABLE security_events (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(255) NOT NULL,
  user_id BIGINT,
  event_type VARCHAR(100) NOT NULL,
  severity VARCHAR(20),
  description TEXT,
  ip_address VARCHAR(45),
  resolved_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_tenant_created (tenant_id, created_at),
  INDEX idx_unresolved (resolved_at)
);
```

## Best Practices

### Encryption
1. Store encryption keys in secure environment variables
2. Rotate encryption keys periodically
3. Use field-level encryption for PII (Personally Identifiable Information)
4. Never log encrypted values

### Audit Logging
1. Enable audit logging for all sensitive operations
2. Archive old audit logs for compliance
3. Monitor security events regularly
4. Generate compliance reports monthly

### GDPR Compliance
1. Obtain explicit user consent before processing data
2. Implement data deletion within 30 days of request
3. Provide data export functionality
4. Maintain consent records with timestamps
5. Document data processing activities

### RBAC Implementation
1. Use principle of least privilege
2. Regularly audit role assignments
3. Separate admin, manager, and agent roles
4. Test permission checks thoroughly
5. Log all role changes

## Security Considerations

- **Principle of Least Privilege**: Users get only necessary permissions
- **Defense in Depth**: Multiple security layers (auth, RBAC, encryption)
- **Audit Trail**: All actions logged for accountability
- **Data Privacy**: Encryption and GDPR compliance
- **Secure Defaults**: Security headers enabled by default

## Testing

All services should be tested with:
1. Unit tests for encryption/decryption
2. Integration tests for RBAC permission checks
3. Audit log verification tests
4. GDPR request processing tests
5. Security header validation tests

## Integration with Existing Features

This Compliance & Security feature integrates seamlessly with:
- Existing JWT authentication
- Multi-tenant isolation middleware
- WebSocket real-time notifications
- Analytics and reporting
- All existing CRUD operations

No breaking changes to existing APIs - all new features are additive.
