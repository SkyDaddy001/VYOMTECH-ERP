# Podman Setup - Quick Start Guide

## What is Podman?

Podman is a daemonless container engine for developing, managing, and running Open Container Initiative (OCI) containers on your Linux System. It's a drop-in replacement for Docker.

**Benefits:**
- No daemon required
- More secure (rootless by default)
- Fully compatible with Docker images
- Lighter weight than Docker Desktop

---

## Quick Start (Choose One)

### Option A: PowerShell Script (Easiest for Windows)

```powershell
# Run the startup script
.\startup.ps1 -Action start

# Stop when done
.\startup.ps1 -Action stop

# Check status
.\startup.ps1 -Action status

# Clean up
.\startup.ps1 -Action clean
```

### Option B: Git Bash Script

```bash
# Make script executable
chmod +x startup.sh

# Run the startup script
./startup.sh start

# Stop when done
./startup.sh stop

# Check status
./startup.sh status

# Clean up
./startup.sh clean
```

### Option C: Manual Commands (PowerShell)

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

# 3. Wait for MySQL
Start-Sleep -Seconds 5

# 4. Run migrations
mysql -h localhost -u callcenter_user -psecure_app_pass callcenter < .\migrations\001_initial_schema.sql

# 5. Set environment & run app
$env:DB_HOST = "localhost"
$env:DB_USER = "callcenter_user"
$env:DB_PASSWORD = "secure_app_pass"
$env:JWT_SECRET = "your-secure-jwt-secret-key-minimum-32-characters"
.\bin\main.exe
```

### Option D: Manual Commands (Git Bash)

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

# 3. Wait for MySQL
sleep 5

# 4. Run migrations
mysql -h localhost -u callcenter_user -psecure_app_pass callcenter < ./migrations/001_initial_schema.sql

# 5. Set environment & run app
export DB_HOST=localhost
export DB_USER=callcenter_user
export DB_PASSWORD=secure_app_pass
export JWT_SECRET="your-secure-jwt-secret-key-minimum-32-characters"
./bin/main
```

---

## Testing the Application

Once running, test in a new terminal:

### Health Check
```bash
curl http://localhost:8080/health
```

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!",
    "role": "agent",
    "tenant_id": "tenant-001"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!"
  }'
```

---

## Common Podman Commands

```bash
# Container Management
podman ps                          # List running containers
podman ps -a                       # List all containers
podman logs <container>            # View container logs
podman logs -f <container>         # Follow container logs
podman exec -it <container> bash   # Access container shell

# Image Management
podman images                      # List images
podman pull mysql:8.0              # Pull image
podman rmi <image>                 # Remove image

# Container Lifecycle
podman run -d <image>              # Run container
podman start <container>           # Start stopped container
podman stop <container>            # Stop running container
podman restart <container>         # Restart container
podman rm <container>              # Remove container

# Network
podman network ls                  # List networks
podman network inspect <network>   # Inspect network

# Machine Management (Windows only)
podman machine ls                  # List machines
podman machine start               # Start machine
podman machine stop                # Stop machine
podman machine ssh                 # SSH into machine
```

---

## Environment Variables

The application requires these environment variables to run:

```bash
DB_HOST=localhost              # Database host
DB_PORT=3306                   # Database port
DB_USER=callcenter_user        # Database user
DB_PASSWORD=secure_app_pass    # Database password
DB_NAME=callcenter             # Database name
SERVER_PORT=8080               # Application port
JWT_SECRET=<minimum-32-chars>  # JWT signing key
```

You can set these in:
1. `.env` file (read by the application)
2. Environment variables
3. PowerShell `$env:` variables
4. Container environment with `-e` flag

---

## Troubleshooting

### "podman: command not found"
Podman is not installed or not in PATH.
```bash
# Install Podman for Windows
# Download from https://podman.io/docs/installation/windows

# Add to PATH if needed
$env:PATH = "$env:PATH;C:\Program Files\RedHat\Podman"
```

### "Cannot connect to MySQL on port 3306"
```bash
# Check if MySQL container is running
podman ps | grep mysql

# Check container logs
podman logs mysql-callcenter

# Verify port is accessible
telnet localhost 3306

# Restart container
podman restart mysql-callcenter
```

### "Application can't connect to database"
```bash
# On Podman Windows with WSL backend:
# Use: host.containers.internal (NOT localhost)

# In your .env or environment:
DB_HOST=host.containers.internal

# Or use Podman machine IP:
podman machine inspect  # Check IP address
```

### "Port 8080 already in use"
```bash
# Find process using port 8080
netstat -ano | findstr :8080

# Kill process
taskkill /PID <PID> /F

# Or use different port
$env:SERVER_PORT=9000
```

### "Database migration fails"
```bash
# Verify MySQL is ready
podman exec mysql-callcenter mysql -u root -prootpass -e "SELECT 1"

# Run migration manually
podman exec mysql-callcenter mysql -u callcenter_user -psecure_app_pass callcenter < ./migrations/001_initial_schema.sql
```

### "Script execution disabled" (PowerShell)
```powershell
# Check execution policy
Get-ExecutionPolicy

# Allow scripts temporarily
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope Process

# Or run directly
pwsh -ExecutionPolicy RemoteSigned -File .\startup.ps1 -Action start
```

---

## Useful Tips

### View Real-time Application Logs
```bash
podman logs -f mysql-callcenter    # MySQL logs
podman logs -f callcenter-app      # App logs (if containerized)
```

### Connect to MySQL Container
```bash
podman exec -it mysql-callcenter mysql -u callcenter_user -psecure_app_pass callcenter
```

### Check Container Network
```bash
podman network inspect podman  # View default network

# Test connectivity
podman exec mysql-callcenter mysql -u root -prootpass -h host.containers.internal
```

### Increase Resource Limits
```bash
# Check current limits
podman system info | grep -A 10 "limits"

# Increase in Podman settings (Windows)
podman machine set --memory 4096 --cpus 4
```

---

## Integration with Development Workflow

### With VS Code
1. Install "Dev Containers" extension
2. Use provided `.devcontainer` config (optional)
3. Debug directly in container

### With IDE Terminal
Most modern IDEs can run your startup script:
- Set terminal to PowerShell (Windows)
- Set terminal to bash (Linux/Mac)
- Run: `.\startup.ps1 -Action start`

### With Git Hooks
Add to `.git/hooks/pre-commit`:
```bash
#!/bin/bash
./startup.sh status
```

---

## Performance Tuning

### For Better Performance
```powershell
# Increase machine resources
podman machine set --memory 8192 --cpus 8 --disk-size 100

# Reduce log size
podman run --log-driver=journald <image>  # Lighter logging
```

### Monitor Resource Usage
```bash
podman stats                       # Real-time stats
podman inspect <container>         # Detailed info
```

---

## Clean Up & Maintenance

### Prune Unused Resources
```bash
podman system prune           # Remove dangling images
podman system prune -a        # Remove all unused images
podman volume prune           # Remove unused volumes
```

### Backup Database
```bash
# Dump MySQL database
podman exec mysql-callcenter mysqldump -u callcenter_user -psecure_app_pass callcenter > backup.sql

# Restore
mysql -u callcenter_user -psecure_app_pass callcenter < backup.sql
```

---

## Next Steps

1. ✅ Start application using one of the options above
2. ✅ Test endpoints with curl or Postman
3. ✅ Check logs if issues arise
4. ✅ Review `QUICK_REFERENCE.md` for API examples
5. ✅ Read `DEVELOPMENT.md` for full development workflow

---

**Resources:**
- Podman Docs: https://podman.io/docs/
- Podman on Windows: https://podman.io/docs/installation/windows
- MySQL Image: https://hub.docker.com/_/mysql
- Go Build: https://golang.org/cmd/go/

**Last Updated:** November 21, 2025  
**Status:** Ready for Podman on Windows
