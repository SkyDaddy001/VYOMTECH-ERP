# Windows Development Setup Guide

Since `make` is not available on Windows bash, here's how to get the application running.

## Prerequisites on Windows

1. **Go 1.24+** - Already installed ✓
2. **MySQL 8.0+** - Install from https://dev.mysql.com/downloads/mysql/ OR use Podman
3. **Git Bash** or **Windows PowerShell**
4. **Podman** (Optional but recommended) - https://podman.io/docs/installation/windows

## Option 1: Using Windows PowerShell (Recommended for Local Dev)

### Step 1: Build the Application
```powershell
cd C:\Users\Skydaddy\Desktop\Developement
go build -o bin\main.exe .\cmd\main.go
```

### Step 2: Start MySQL (if installed locally)
```powershell
# If MySQL is installed as a service
net start MySQL80

# Or if you have it running via Podman:
podman run -d --name mysql-dev `
  -e MYSQL_ROOT_PASSWORD=rootpass `
  -e MYSQL_DATABASE=callcenter `
  -e MYSQL_USER=callcenter_user `
  -e MYSQL_PASSWORD=secure_app_pass `
  -p 3306:3306 `
  mysql:8.0
```

### Step 3: Initialize Database
```powershell
# Create .env file
Copy-Item .env.example -Destination .env

# Run migrations
mysql -h localhost -u callcenter_user -p secure_app_pass callcenter < .\migrations\001_initial_schema.sql
```

### Step 4: Run the Application
```powershell
# Set environment variables
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_USER = "callcenter_user"
$env:DB_PASSWORD = "secure_app_pass"
$env:DB_NAME = "callcenter"
$env:SERVER_PORT = "8080"
$env:JWT_SECRET = "your-secure-jwt-secret-key-change-in-production-minimum-32-chars"

# Run
.\bin\main.exe
```

## Option 2: Using Git Bash

### Step 1: Build
```bash
cd /c/Users/Skydaddy/Desktop/Developement
go build -o bin/main ./cmd/main.go
```

### Step 2: Setup Environment
```bash
# Create .env from template
cp .env.example .env

# Edit .env with your MySQL credentials
nano .env
```

### Step 3: Run Database Migrations
```bash
mysql -u callcenter_user -p < ./migrations/001_initial_schema.sql
# Enter password: secure_app_pass
```

### Step 4: Run Application
```bash
export SERVER_PORT=8080
./bin/main
```

## Option 3: Using Podman (Recommended for Full Stack)

### Prerequisites
- Podman installed and running
- Podman machine initialized (on Windows, use WSL backend)

### Step 1: Initialize Podman Machine (if not done)
```powershell
podman machine init
podman machine start
```

### Step 2: Start MySQL Container with Podman
```powershell
podman run -d --name mysql-callcenter `
  -e MYSQL_ROOT_PASSWORD=rootpass `
  -e MYSQL_DATABASE=callcenter `
  -e MYSQL_USER=callcenter_user `
  -e MYSQL_PASSWORD=secure_app_pass `
  -p 3306:3306 `
  mysql:8.0
```

### Step 3: Initialize Database
```powershell
# Wait for MySQL to be ready
Start-Sleep -Seconds 5

# Run migrations
podman exec mysql-callcenter mysql -u callcenter_user -psecure_app_pass callcenter < .\migrations\001_initial_schema.sql
```

### Step 4: Build Podman Image
```powershell
cd C:\Users\Skydaddy\Desktop\Developement
podman build -t callcenter:latest .
```

### Step 5: Run Application Container
```powershell
podman run -d --name callcenter-app `
  -e DB_HOST=host.containers.internal `
  -e DB_PORT=3306 `
  -e DB_USER=callcenter_user `
  -e DB_PASSWORD=secure_app_pass `
  -e DB_NAME=callcenter `
  -e SERVER_PORT=8080 `
  -e JWT_SECRET="your-secure-jwt-secret-key-change-in-production-minimum-32-chars" `
  -p 8080:8080 `
  callcenter:latest
```

**Note:** For Podman on Windows with WSL, use `host.containers.internal` instead of `localhost` for host machine access.

### Step 6: View Logs
```powershell
podman logs -f callcenter-app
```

## Option 4: Using Podman Compose (Full Stack)

### Step 1: Create podman-compose.yml
```powershell
# The docker-compose.yml should work with podman as well
# Just use podman-compose instead of docker-compose
```

### Step 2: Start Full Stack
```bash
podman-compose up -d
```

### Step 3: Check Status
```bash
podman-compose ps
podman-compose logs -f
```

## Testing the Application

### Health Check
```powershell
# PowerShell
Invoke-WebRequest -Uri "http://localhost:8080/health"

# Or using curl (if installed)
curl http://localhost:8080/health
```

### Register a User
```powershell
$body = @{
    email = "test@example.com"
    password = "TestPassword123!"
    role = "agent"
    tenant_id = "tenant-001"
} | ConvertTo-Json

Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" `
  -Method POST `
  -Body $body `
  -ContentType "application/json"
```

### Login
```powershell
$body = @{
    email = "test@example.com"
    password = "TestPassword123!"
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST `
  -Body $body `
  -ContentType "application/json"

# Save token for use in other requests
$token = ($response.Content | ConvertFrom-Json).token
Write-Host "Token: $token"
```

### Get Agents (Protected Endpoint)
```powershell
$headers = @{
    Authorization = "Bearer $token"
}

Invoke-WebRequest -Uri "http://localhost:8080/api/v1/agents" `
  -Method GET `
  -Headers $headers
```

## Stopping the Application

### With Podman
```powershell
# Stop application container
podman stop callcenter-app

# Stop MySQL container
podman stop mysql-callcenter

# Remove containers
podman rm callcenter-app mysql-callcenter

# Remove images
podman rmi callcenter:latest
```

### Local Process
```powershell
# Stop local application
Get-Process main | Stop-Process -Force

# Or using Task Manager - find "main.exe" and end task
```

## Podman Useful Commands

```bash
# Check Podman status
podman ps                    # List running containers
podman ps -a                 # List all containers
podman images                # List images
podman logs <container>      # View container logs
podman exec -it <container> /bin/bash  # Execute command in container

# Machine management
podman machine ls            # List machines
podman machine start         # Start default machine
podman machine stop          # Stop default machine
podman machine rm            # Remove machine
```

## Troubleshooting

### "podman: command not found" in Git Bash
```bash
# Add Podman to PATH or use full path
/mnt/c/Program\ Files/RedHat/Podman/podman.exe ps

# Or create an alias
echo 'alias podman="/mnt/c/Program\ Files/RedHat/Podman/podman.exe"' >> ~/.bashrc
source ~/.bashrc
```

### MySQL Connection Refused on Port 3306
```powershell
# Check if container is running
podman ps | grep mysql

# Check container logs
podman logs mysql-callcenter

# Restart container
podman restart mysql-callcenter

# Wait a few seconds and try again
```

### Application Can't Connect to Database
```powershell
# For Podman on Windows, use:
DB_HOST=host.containers.internal   # NOT localhost

# Verify in .env or set environment variable
$env:DB_HOST = "host.containers.internal"
```

### Port Already in Use
```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill process by PID
taskkill /PID <PID> /F

# Or use different port
$env:SERVER_PORT = "9000"
.\bin\main.exe
```

### Database Migration Fails
```powershell
# Check if MySQL is ready
podman exec mysql-callcenter mysql -u root -prootpass -e "SELECT 1"

# Run migration manually
podman exec mysql-callcenter mysql -u callcenter_user -psecure_app_pass callcenter < .\migrations\001_initial_schema.sql

# Or directly
mysql -h localhost -u callcenter_user -psecure_app_pass callcenter < .\migrations\001_initial_schema.sql
```

## Quick Start Command (Copy & Paste)

### PowerShell Quick Start with Podman
```powershell
# 1. Build
go build -o bin\main.exe .\cmd\main.go

# 2. Start MySQL
podman run -d --name mysql-callcenter `
  -e MYSQL_ROOT_PASSWORD=rootpass `
  -e MYSQL_DATABASE=callcenter `
  -e MYSQL_USER=callcenter_user `
  -e MYSQL_PASSWORD=secure_app_pass `
  -p 3306:3306 `
  mysql:8.0

# 3. Wait and setup DB
Start-Sleep -Seconds 5
mysql -h localhost -u callcenter_user -psecure_app_pass callcenter < .\migrations\001_initial_schema.sql

# 4. Run app
$env:DB_HOST = "localhost"
$env:DB_USER = "callcenter_user"
$env:DB_PASSWORD = "secure_app_pass"
$env:JWT_SECRET = "your-secure-jwt-secret-key-minimum-32-characters"
.\bin\main.exe

# 5. In another terminal, test
curl http://localhost:8080/health
```

### Git Bash Quick Start with Podman
```bash
# 1. Build
go build -o bin/main ./cmd/main.go

# 2. Start MySQL
podman run -d --name mysql-callcenter \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_DATABASE=callcenter \
  -e MYSQL_USER=callcenter_user \
  -e MYSQL_PASSWORD=secure_app_pass \
  -p 3306:3306 \
  mysql:8.0

# 3. Wait and setup DB
sleep 5
mysql -h localhost -u callcenter_user -psecure_app_pass callcenter < ./migrations/001_initial_schema.sql

# 4. Run app
export DB_HOST=localhost
export DB_USER=callcenter_user
export DB_PASSWORD=secure_app_pass
export JWT_SECRET="your-secure-jwt-secret-key-minimum-32-characters"
./bin/main

# 5. In another terminal
curl http://localhost:8080/health
```

## Resources

- **Podman Documentation**: https://podman.io/docs/
- **MySQL Docker Image**: https://hub.docker.com/_/mysql
- **Go Build**: https://golang.org/cmd/go/
- **PowerShell Docs**: https://docs.microsoft.com/en-us/powershell/

---

**Last Updated:** November 21, 2025
**Status:** Ready for Podman on Windows

### "Connection refused on port 3306"
- MySQL not running - start MySQL service or Docker container
- Check MySQL is listening: `netstat -an | findstr :3306`

## Quick Commands Reference

### PowerShell
```powershell
# Build
go build -o bin\main.exe .\cmd\main.go

# Run
.\bin\main.exe

# Clean
Remove-Item -Path "bin\main.exe"

# Check process
Get-Process | Where-Object {$_.Name -like "*main*"}
```

### Git Bash
```bash
# Build
go build -o bin/main ./cmd/main.go

# Run
./bin/main

# Clean
rm bin/main

# Check process
ps aux | grep main
```

## Next Steps

1. Start the application using one of the options above
2. Test endpoints using curl or PowerShell
3. Check the `QUICK_REFERENCE.md` for more API examples
4. Review `DEVELOPMENT.md` for detailed development workflow

---

For Windows users, **Option 1 (PowerShell) or Option 3 (Docker Desktop)** are recommended as they handle Windows paths correctly.
