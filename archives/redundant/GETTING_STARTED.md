# Quick Start - Multi-Tenant AI Call Center (Windows)

## Current Status
✅ Application compiled successfully: `bin/main`  
❌ Database not running (requires MySQL or Docker)  
⏳ Application cannot start without database

---

## Getting Started - Choose Your Setup

### **Setup Option 1: Standalone (No Docker Required)**

#### Step 1: Install MySQL
1. Download MySQL 8.0+ from https://dev.mysql.com/downloads/mysql/
2. Install with default settings
3. MySQL server starts automatically after installation

#### Step 2: Create Database
Open PowerShell or Command Prompt and run:
```powershell
# Open MySQL command line
mysql -u root -p

# Then paste these SQL commands:
CREATE DATABASE callcenter CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'callcenter_user'@'localhost' IDENTIFIED BY 'secure_app_pass';
GRANT ALL PRIVILEGES ON callcenter.* TO 'callcenter_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;

# Import schema
mysql -u callcenter_user -p secure_app_pass callcenter < C:\Users\Skydaddy\Desktop\Developement\migrations\001_initial_schema.sql
```

#### Step 3: Create .env File
Copy `.env.example` to `.env`:
```powershell
Copy-Item -Path ".env.example" -Destination ".env"
```

Edit `.env` with your settings:
```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=callcenter_user
DB_PASSWORD=secure_app_pass
DB_NAME=callcenter
JWT_SECRET=your-secure-jwt-secret-key-change-in-production-minimum-32-chars
DEBUG=true
```

#### Step 4: Run the Application
```powershell
cd C:\Users\Skydaddy\Desktop\Developement
.\bin\main.exe
```

You should see:
```
2025/11/21 14:30:00 [INFO] Starting Multi-Tenant AI Call Center application
2025/11/21 14:30:01 [INFO] Database connection established successfully
2025/11/21 14:30:01 [INFO] Server starting port 8080
```

**✅ Application is running!** Open http://localhost:8080/health in your browser.

---

### **Setup Option 2: Using Docker Desktop**

#### Prerequisites
- Download and install Docker Desktop for Windows: https://www.docker.com/products/docker-desktop
- Launch Docker Desktop and wait for it to start

#### Step 1: Start MySQL Container
```powershell
docker run -d --name mysql-callcenter `
  -e MYSQL_ROOT_PASSWORD=rootpass `
  -e MYSQL_DATABASE=callcenter `
  -e MYSQL_USER=callcenter_user `
  -e MYSQL_PASSWORD=secure_app_pass `
  -p 3306:3306 `
  mysql:8.0
```

#### Step 2: Wait and Initialize Database
```powershell
# Wait for MySQL to be ready
Start-Sleep -Seconds 10

# Import schema
docker exec -i mysql-callcenter mysql -u callcenter_user -psecure_app_pass callcenter < C:\Users\Skydaddy\Desktop\Developement\migrations\001_initial_schema.sql
```

#### Step 3: Create .env File
```powershell
Copy-Item -Path ".env.example" -Destination ".env"
```

Edit `.env`:
```
SERVER_PORT=8080
DB_HOST=host.docker.internal
DB_PORT=3306
DB_USER=callcenter_user
DB_PASSWORD=secure_app_pass
DB_NAME=callcenter
JWT_SECRET=your-secure-jwt-secret-key-change-in-production-minimum-32-chars
DEBUG=true
```

#### Step 4: Run the Application
```powershell
cd C:\Users\Skydaddy\Desktop\Developement
.\bin\main.exe
```

---

## Testing the Application

### Health Check
Open your browser: **http://localhost:8080/health**

Expected response:
```json
{"status":"healthy"}
```

### Register a Test User
```powershell
$body = @{
    email = "admin@example.com"
    password = "Admin123456!"
    role = "admin"
    tenant_id = "tenant-001"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" `
  -Method POST `
  -Body $body `
  -ContentType "application/json" `
  -UseBasicParsing

$response.Content | ConvertFrom-Json | ConvertTo-Json
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "role": "admin",
    "tenant_id": "tenant-001"
  },
  "message": "User registered successfully"
}
```

### Login and Get JWT Token
```powershell
$loginBody = @{
    email = "admin@example.com"
    password = "Admin123456!"
} | ConvertTo-Json

$loginResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST `
  -Body $loginBody `
  -ContentType "application/json" `
  -UseBasicParsing

$token = ($loginResponse.Content | ConvertFrom-Json).token
Write-Host "Your JWT Token: $token"
```

### Test Protected Endpoint (Get Agents)
```powershell
# Use the token from above
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-WebRequest -Uri "http://localhost:8080/api/v1/agents" `
  -Method GET `
  -Headers $headers `
  -UseBasicParsing | Select-Object -ExpandProperty Content | ConvertFrom-Json
```

---

## Stopping the Application

### If Running Standalone
```powershell
# Press Ctrl+C in the PowerShell window running the app
# Or kill the process:
Stop-Process -Name main -Force
```

### If Using Docker
```powershell
# Stop and remove containers
docker stop mysql-callcenter
docker rm mysql-callcenter

# Verify
docker ps -a
```

---

## Common Issues & Solutions

### Issue: "Could not connect to server" on localhost:8080
**Solution:**
- Check if application is running: `Get-Process | Where-Object {$_.Name -like "*main*"}`
- Check error messages in the terminal
- Ensure database is running

### Issue: "Database error: connection refused"
**Solution:**
- If using local MySQL: Ensure MySQL service is running
- If using Docker: Verify container is running: `docker ps`
- Check DB_HOST in .env file

### Issue: "Error: database error: no such table"
**Solution:**
- Migrations not applied. Run:
  ```powershell
  mysql -u callcenter_user -p secure_app_pass callcenter < .\migrations\001_initial_schema.sql
  ```

### Issue: Port 8080 already in use
**Solution:**
- Find what's using the port:
  ```powershell
  netstat -ano | findstr :8080
  ```
- Stop that process or use different port:
  ```powershell
  $env:SERVER_PORT = 9000
  .\bin\main.exe
  ```

### Issue: "authentication plugin caching_sha2_password cannot be loaded"
**Solution:**
- Create user with mysql_native_password:
  ```sql
  CREATE USER 'callcenter_user'@'localhost' IDENTIFIED WITH mysql_native_password BY 'secure_app_pass';
  ```

---

## Development Workflow

### Building After Changes
```powershell
cd C:\Users\Skydaddy\Desktop\Developement
go build -o bin\main.exe .\cmd\main.go
```

### Code Formatting
```powershell
go fmt .\...
```

### Running Tests
```powershell
go test -v .\...
```

### Viewing Logs
The application logs to console. For persistent logs, redirect output:
```powershell
.\bin\main.exe > app.log 2>&1
```

---

## Next Steps

1. ✅ Choose setup option (Standalone or Docker)
2. ✅ Setup database
3. ✅ Create .env file
4. ✅ Run `.\bin\main.exe`
5. ✅ Test health endpoint in browser
6. ✅ Review `QUICK_REFERENCE.md` for more API examples

---

## Full API Endpoint Reference

### Authentication (Public)
```
POST /api/v1/auth/register          - Register new user
POST /api/v1/auth/login             - Login and get JWT token
POST /api/v1/password-reset/request - Request password reset
POST /api/v1/password-reset/reset   - Complete password reset
```

### Authentication (Protected)
```
GET  /api/v1/auth/validate          - Validate JWT token
POST /api/v1/auth/change-password   - Change password
```

### Agents (Protected, Tenant-Isolated)
```
GET    /api/v1/agents               - List all agents
GET    /api/v1/agents/{id}          - Get agent by ID
PATCH  /api/v1/agents/status        - Update availability
GET    /api/v1/agents/available     - List available agents
GET    /api/v1/agents/stats         - Get agent statistics
```

### Health & Status (Public)
```
GET /health                         - Health check
GET /ready                          - Readiness probe
```

---

## Where to Go From Here

- **DEVELOPMENT.md** - Comprehensive development guide
- **QUICK_REFERENCE.md** - Command examples and API calls
- **AUDIT_COMPLETION_REPORT.md** - Complete project status
- **WINDOWS_SETUP.md** - Detailed Windows instructions

---

**Created:** November 21, 2025  
**Application:** Multi-Tenant AI Call Center  
**Status:** ✅ Ready to run
