#!/usr/bin/env pwsh
# VYOM ERP - Docker Startup Script
# Simplified deployment without Dockerfile builds

param(
    [string]$Action = "start",
    [string]$Service = "",
    [switch]$Rebuild = $false,
    [switch]$Clean = $false
)

# Colors for output
$Green = "`e[32m"
$Red = "`e[31m"
$Yellow = "`e[33m"
$Blue = "`e[34m"
$Reset = "`e[0m"

function Write-Success {
    param([string]$Message)
    Write-Host "$Green✓$Reset $Message"
}

function Write-Error {
    param([string]$Message)
    Write-Host "$Red✗$Reset $Message"
}

function Write-Info {
    param([string]$Message)
    Write-Host "$Blue ℹ$Reset $Message"
}

function Write-Warning {
    param([string]$Message)
    Write-Host "$Yellow⚠$Reset $Message"
}

# Check prerequisites
function Check-Prerequisites {
    Write-Info "Checking prerequisites..."
    
    # Check Docker
    try {
        $version = docker --version
        Write-Success "Docker installed: $version"
    } catch {
        Write-Error "Docker not found. Please install Docker Desktop."
        exit 1
    }
    
    # Check Docker Compose
    try {
        $version = docker-compose --version
        Write-Success "Docker Compose installed: $version"
    } catch {
        Write-Error "Docker Compose not found. Please install Docker Compose."
        exit 1
    }
    
    # Check .env file
    if (-not (Test-Path ".env")) {
        Write-Warning ".env file not found. Creating from .env.example..."
        if (Test-Path ".env.example") {
            Copy-Item ".env.example" ".env"
            Write-Success ".env created. Please edit with your configuration."
        } else {
            Write-Error ".env.example not found."
            exit 1
        }
    }
}

# Start services
function Start-Services {
    param([bool]$RebuildImages = $false)
    
    Write-Info "Starting VYOM ERP services..."
    
    # Check if using old or new compose file
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    if ($RebuildImages) {
        Write-Info "Building and starting services..."
        docker-compose -f $ComposeFile up -d --build
    } else {
        Write-Info "Starting services (building if needed)..."
        docker-compose -f $ComposeFile up -d
    }
    
    if ($?) {
        Write-Success "Services started successfully"
        Start-Sleep -Seconds 3
        Show-ServiceStatus
    } else {
        Write-Error "Failed to start services"
        exit 1
    }
}

# Stop services
function Stop-Services {
    Write-Info "Stopping VYOM ERP services..."
    
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    docker-compose -f $ComposeFile stop
    
    if ($?) {
        Write-Success "Services stopped"
    } else {
        Write-Error "Failed to stop services"
        exit 1
    }
}

# Remove services
function Remove-Services {
    param([bool]$RemoveData = $false)
    
    Write-Warning "Removing VYOM ERP services..."
    
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    if ($RemoveData) {
        Write-Warning "Also removing all volumes and data..."
        docker-compose -f $ComposeFile down -v
    } else {
        docker-compose -f $ComposeFile down
    }
    
    if ($?) {
        Write-Success "Services removed"
    } else {
        Write-Error "Failed to remove services"
        exit 1
    }
}

# Show service status
function Show-ServiceStatus {
    Write-Info "Service Status:"
    
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    docker-compose -f $ComposeFile ps
    
    Write-Info ""
    Write-Info "Access Points:"
    Write-Host "  Frontend:   $Blue http://localhost:3000$Reset"
    Write-Host "  API:        $Blue http://localhost:8080$Reset"
    Write-Host "  Adminer:    $Blue http://localhost:8081$Reset"
    Write-Host "  PHPMyAdmin: $Blue http://localhost:8082$Reset"
}

# View logs
function Show-Logs {
    param([string]$Service = "")
    
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    if ($Service) {
        Write-Info "Showing logs for $Service..."
        docker-compose -f $ComposeFile logs -f $Service
    } else {
        Write-Info "Showing logs for all services..."
        docker-compose -f $ComposeFile logs -f
    }
}

# Restart services
function Restart-Services {
    param([string]$Service = "")
    
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    if ($Service) {
        Write-Info "Restarting $Service..."
        docker-compose -f $ComposeFile restart $Service
    } else {
        Write-Info "Restarting all services..."
        docker-compose -f $ComposeFile restart
    }
    
    if ($?) {
        Write-Success "Services restarted"
    } else {
        Write-Error "Failed to restart services"
        exit 1
    }
}

# Build images
function Build-Images {
    Write-Info "Building images with docker-compose..."
    
    $ComposeFile = "docker-compose.new.yml"
    if (-not (Test-Path $ComposeFile)) {
        $ComposeFile = "docker-compose.yml"
    }
    
    docker-compose -f $ComposeFile build
    
    if ($?) {
        Write-Success "Images built successfully"
    } else {
        Write-Error "Failed to build images"
        exit 1
    }
}

# Show help
function Show-Help {
    Write-Host @"
$Blue╔════════════════════════════════════════════════════════╗
║       VYOM ERP - Docker Deployment Script               ║
╚════════════════════════════════════════════════════════╝$Reset

USAGE:
  .\deploy.ps1 [Action] [Options]

ACTIONS:
  start      Start services (default)
  stop       Stop services
  restart    Restart services
  logs       Show logs
  status     Show service status
  build      Build images from source
  clean      Remove services and data
  help       Show this help message

OPTIONS:
  -Service <name>    Specify service (backend, frontend, mysql, etc.)
  -Rebuild           Rebuild images before starting
  -Clean             Remove containers and volumes

EXAMPLES:
  .\deploy.ps1                              # Start services
  .\deploy.ps1 -Rebuild                     # Start with rebuild
  .\deploy.ps1 stop                         # Stop services
  .\deploy.ps1 logs -Service backend        # View backend logs
  .\deploy.ps1 restart -Service frontend    # Restart frontend
  .\deploy.ps1 build                        # Build images
  .\deploy.ps1 clean                        # Remove everything

SERVICES:
  - mysql      Database server
  - backend    Go API server
  - frontend   Next.js application
  - adminer    Database UI
  - phpmyadmin Database UI (optional)

$Reset
"@
}

# Main execution
Write-Host @"
$Blue╔════════════════════════════════════════════════════════╗
║          VYOM ERP - Docker Deployment                  ║
║                                                        ║
║  Deployment Command: $Action $Reset
╚════════════════════════════════════════════════════════╝
$Reset
"@

Check-Prerequisites

switch ($Action.ToLower()) {
    "start" {
        Start-Services -RebuildImages $Rebuild
    }
    "stop" {
        Stop-Services
    }
    "restart" {
        if ($Service) {
            Restart-Services -Service $Service
        } else {
            Stop-Services
            Start-Sleep -Seconds 2
            Start-Services
        }
    }
    "logs" {
        Show-Logs -Service $Service
    }
    "status" {
        Show-ServiceStatus
    }
    "build" {
        Build-Images
    }
    "clean" {
        Remove-Services -RemoveData $true
        if (Test-Path "data") {
            Remove-Item -Recurse -Force "data" -ErrorAction SilentlyContinue
        }
        Write-Success "Cleanup complete"
    }
    "help" {
        Show-Help
    }
    default {
        Write-Error "Unknown action: $Action"
        Show-Help
        exit 1
    }
}

Write-Host ""
Write-Success "Done!"
