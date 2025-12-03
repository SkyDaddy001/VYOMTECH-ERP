# Demo Test Credentials - Setup & Troubleshooting

**Last Updated:** December 3, 2025

---

## Overview

The VYOM ERP system includes demo test credentials for easy testing and evaluation. These credentials are displayed in the login form to help new users quickly access the system.

---

## Available Demo Credentials

| Email | Password | Role | Access |
|-------|----------|------|--------|
| `demo@vyomtech.com` | `DemoPass@123` | Admin | Full system access |
| `agent@vyomtech.com` | `AgentPass@123` | Agent | Call center & leads |
| `manager@vyomtech.com` | `ManagerPass@123` | Manager | Team management |
| `sales@vyomtech.com` | `SalesPass@123` | Sales | Sales pipeline |
| `hr@vyomtech.com` | `HRPass@123` | HR | HR & payroll |

---

## Setup Instructions

### Automatic Seeding (Development Mode)

Demo users are **automatically created** when the application starts in development mode.

**How it works:**
1. Application starts
2. Checks if environment is NOT production (`APP_ENV != production`)
3. Automatically creates/updates demo users in database
4. Logs seeding status

**Requirements:**
- Database migrations must be applied first
- Proper database connection configured
- Application not in production mode

### Manual Setup (If Needed)

#### Option 1: Using Shell Script

```bash
cd /d/VYOMTECH-ERP/scripts

# Set database credentials (optional, defaults to localhost)
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=vyomerp

# Run the seeding script
bash init-demo-users.sh
```

**Output:**
```
========================================
VYOMTECH ERP - Creating Demo Users
========================================
Demo users created successfully!

Available Demo Credentials:
============================================
1. Admin User
   Email: demo@vyomtech.com
   Password: DemoPass@123
   ...
============================================
```

#### Option 2: Using MySQL Client Directly

```sql
-- Create demo tenant
INSERT INTO tenant (id, name, status, max_users, max_concurrent_calls, ai_budget_monthly) 
VALUES ('demo-tenant', 'Demo Tenant', 'active', 100, 50, 1000.00)
ON DUPLICATE KEY UPDATE name='Demo Tenant';

-- Create demo users
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES ('demo@vyomtech.com', '$2a$10$...hash...', 'admin', 'demo-tenant', NOW(), NOW());
```

#### Option 3: Using Go CLI (After Building)

```bash
# Build the seeder binary
cd /d/VYOMTECH-ERP
go build -o ./bin/main ./cmd/main.go

# Run application (seeder runs automatically in dev mode)
./bin/main
```

---

## Troubleshooting

### Problem: "Invalid credentials" Error

**Cause:** Demo users don't exist in database

**Solution:**

1. **Check if demo users exist:**
   ```sql
   SELECT email, role FROM user WHERE tenant_id = 'demo-tenant';
   ```

2. **If no results, seed the users:**
   ```bash
   bash /d/VYOMTECH-ERP/scripts/init-demo-users.sh
   ```

3. **Verify tenant exists:**
   ```sql
   SELECT * FROM tenant WHERE id = 'demo-tenant';
   ```

### Problem: Seeder Doesn't Run Automatically

**Cause:** Application might be in production mode

**Solution:**

1. **Check APP_ENV variable:**
   ```bash
   echo $APP_ENV
   ```

2. **Set to development:**
   ```bash
   export APP_ENV=development
   ```

3. **Restart application:**
   ```bash
   go run cmd/main.go
   ```

### Problem: Password Hash Mismatch

**Cause:** Bcrypt hash algorithm mismatch or corrupted data

**Solution:**

1. **Delete incorrect user:**
   ```sql
   DELETE FROM user WHERE email = 'demo@vyomtech.com';
   ```

2. **Run seeding script to recreate:**
   ```bash
   bash /d/VYOMTECH-ERP/scripts/init-demo-users.sh
   ```

### Problem: Demo Users Show But Can't Login

**Cause:** Tenant doesn't exist or user not linked to tenant

**Solution:**

1. **Verify user-tenant relationship:**
   ```sql
   SELECT u.email, u.role, u.tenant_id, t.name 
   FROM user u 
   LEFT JOIN tenant t ON u.tenant_id = t.id 
   WHERE u.email = 'demo@vyomtech.com';
   ```

2. **If tenant is NULL, update:**
   ```sql
   UPDATE user SET tenant_id = 'demo-tenant' WHERE email = 'demo@vyomtech.com';
   ```

---

## How It Works

### Seeding Process

1. **Application Start:**
   - Main.go initializes database connection
   - Checks environment variable `APP_ENV`
   - If NOT "production", creates DemoUserSeeder

2. **Seeder Execution:**
   ```go
   // Automatically runs in main.go
   if os.Getenv("APP_ENV") != "production" {
       seeder := services.NewDemoUserSeeder(dbConn, log)
       seeder.SeedDemoUsers(ctx)
   }
   ```

3. **User Creation:**
   - Ensures demo tenant exists
   - For each demo user:
     - Checks if user already exists
     - If exists: updates password and role
     - If not: creates new user with hashed password
   - Logs all operations

### Password Hashing

Passwords are hashed using **bcrypt** (Go's crypto package):

```go
hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password), 
    bcrypt.DefaultCost,
)
```

**Security:**
- One-way hashing (cannot be reversed)
- Each hash is unique (includes salt)
- Cost factor = 10 (industry standard)

---

## Frontend Integration

### Login Form Display

The frontend automatically displays demo credentials in a **green card** on the login page:

```tsx
// frontend/components/auth/LoginForm.tsx
const TEST_CREDENTIALS = [
  {
    email: 'demo@vyomtech.com',
    password: 'DemoPass@123',
    role: 'Admin',
    description: 'Full system access',
  },
  // ... other credentials
]
```

### Auto-Fill & Login

Users can:
1. **Click any credential card** → Auto-fills email & password
2. **Click Sign In** → Logs in with that account
3. Or manually type email/password

---

## Database Schema

### User Table
```sql
CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    tenant_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenant(id) ON DELETE CASCADE
);
```

### Tenant Table
```sql
CREATE TABLE tenant (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status ENUM('active', 'inactive'),
    max_users INT DEFAULT 100,
    max_concurrent_calls INT DEFAULT 50,
    ai_budget_monthly DECIMAL(10, 2)
);
```

---

## Backend API

### Login Endpoint

**Request:**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "demo@vyomtech.com",
  "password": "DemoPass@123"
}
```

**Response (Success):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "demo@vyomtech.com",
    "role": "admin",
    "tenant_id": "demo-tenant"
  },
  "message": "Login successful"
}
```

**Response (Failure):**
```json
{
  "error": "invalid credentials"
}
```

### Implementation

**Handler:** `internal/handlers/auth.go`
```go
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    token, err := h.authService.Login(ctx, req.Email, req.Password)
    // Verify password and return token
}
```

**Service:** `internal/services/auth.go`
```go
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
    // Query user by email
    // Verify password with bcrypt.CompareHashAndPassword
    // Generate JWT token
}
```

---

## Testing

### Manual Testing

1. **Start Application:**
   ```bash
   cd /d/VYOMTECH-ERP
   go run cmd/main.go
   ```

2. **Open Login Page:**
   - Browser: `http://localhost:3000/auth/login`
   - See green "Demo Test Credentials" card

3. **Click a Credential:**
   - Any credential card auto-fills form
   - Click "Sign In"
   - Should redirect to dashboard

### Testing with cURL

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@vyomtech.com",
    "password": "DemoPass@123"
  }'
```

**Expected Response:**
```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "demo@vyomtech.com",
    "role": "admin",
    "tenant_id": "demo-tenant"
  },
  "message": "Login successful"
}
```

---

## File Locations

| File | Purpose |
|------|---------|
| `cmd/main.go` | Calls seeder on app startup |
| `internal/services/demo_seeder.go` | Seeding logic |
| `scripts/init-demo-users.sh` | Manual seeding script |
| `internal/handlers/auth.go` | Login endpoint |
| `internal/services/auth.go` | Auth business logic |
| `frontend/components/auth/LoginForm.tsx` | Login UI with credentials |

---

## Production Considerations

### ⚠️ Security Note

Demo credentials should **NOT** be used in production:

1. **Disable in Production:**
   ```bash
   export APP_ENV=production
   # Seeder will NOT run
   ```

2. **Remove Frontend Display:**
   - Remove TEST_CREDENTIALS from LoginForm
   - Hide demo credential card
   - Or move to admin-only area

3. **Strong Credentials:**
   - Create admin account with strong password
   - Use SSO/OAuth when possible
   - Implement 2FA

### Recommendations

1. **Development:** Keep demo users for testing
2. **Staging:** Remove or limit demo users
3. **Production:** 
   - Disable automatic seeding (`APP_ENV=production`)
   - Create real admin account
   - Use enterprise authentication

---

## Modification

### To Change Demo Credentials

**Option 1: Update LoginForm (Frontend)**
```tsx
// frontend/components/auth/LoginForm.tsx
const TEST_CREDENTIALS = [
  {
    email: 'your@email.com',
    password: 'YourPassword123',
    role: 'Admin',
    description: 'Full access',
  },
]
```

**Option 2: Update Seeder (Backend)**
```go
// internal/services/demo_seeder.go
demoUsers := []DemoUser{
    {
        Email:    "your@email.com",
        Password: "YourPassword123",
        Role:     "admin",
        TenantID: demoTenantID,
    },
}
```

**Option 3: Update Script**
```bash
# scripts/init-demo-users.sh
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES (
    'your@email.com',
    '$2a$10$...your_bcrypt_hash...',
    'admin',
    'demo-tenant',
    NOW(),
    NOW()
);
```

---

## Support

If demo credentials are still not working after following this guide:

1. **Check logs:**
   ```bash
   # Backend logs during startup
   go run cmd/main.go 2>&1 | grep -i "demo\|seed"
   ```

2. **Verify database:**
   ```sql
   SELECT COUNT(*) FROM user WHERE tenant_id = 'demo-tenant';
   ```

3. **Check migrations:**
   - Ensure all migrations applied: `001_foundation.sql` through `022_project_management_system.sql`
   - User and tenant tables must exist

4. **Test API directly:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"demo@vyomtech.com","password":"DemoPass@123"}'
   ```

---

**If issues persist, please check:**
- ✅ Database connectivity
- ✅ Migrations applied successfully
- ✅ Application environment set correctly
- ✅ No errors in application logs
