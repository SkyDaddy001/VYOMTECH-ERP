# VYOM ERP - Environment Setup Script (Windows PowerShell)
# This script helps initialize environment variables for development

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "VYOM ERP - Environment Setup" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Check if .env exists
if (Test-Path .env) {
    Write-Host ".env file already exists." -ForegroundColor Yellow
    $response = Read-Host "Do you want to overwrite it? (y/n)"
    if ($response -ne "y" -and $response -ne "Y") {
        Write-Host "Keeping existing .env file" -ForegroundColor Green
        exit 0
    }
}

# Copy template
Write-Host "Creating .env from template..." -ForegroundColor Cyan
Copy-Item .env.example .env
Write-Host "✓ .env created" -ForegroundColor Green
Write-Host ""

# Database Configuration
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "DATABASE CONFIGURATION" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

$db_host = Read-Host "Database Host (localhost)"
if ([string]::IsNullOrWhiteSpace($db_host)) { $db_host = "localhost" }

$db_port = Read-Host "Database Port (5432)"
if ([string]::IsNullOrWhiteSpace($db_port)) { $db_port = "5432" }

$db_user = Read-Host "Database User (vyom_user)"
if ([string]::IsNullOrWhiteSpace($db_user)) { $db_user = "vyom_user" }

$db_password = Read-Host "Database Password" -AsSecureString
$db_password_plain = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToCoTaskMemUnicode($db_password))

$db_name = Read-Host "Database Name (vyomtech_db)"
if ([string]::IsNullOrWhiteSpace($db_name)) { $db_name = "vyomtech_db" }

# Update .env with database config
$env_content = Get-Content .env
$env_content = $env_content -replace "DB_HOST=.*", "DB_HOST=$db_host"
$env_content = $env_content -replace "DB_PORT=.*", "DB_PORT=$db_port"
$env_content = $env_content -replace "DB_USER=.*", "DB_USER=$db_user"
$env_content = $env_content -replace "DB_NAME=.*", "DB_NAME=$db_name"
$env_content = $env_content -replace "DATABASE_URL=.*", "DATABASE_URL=postgresql://$db_user`:$db_password_plain@$db_host`:$db_port/$db_name"
$env_content | Set-Content .env

Write-Host "✓ Database configuration updated" -ForegroundColor Green
Write-Host ""

# Redis Configuration (Optional)
$response = Read-Host "Configure Redis? (y/n)"
if ($response -eq "y" -or $response -eq "Y") {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "REDIS CONFIGURATION" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    
    $redis_host = Read-Host "Redis Host (localhost)"
    if ([string]::IsNullOrWhiteSpace($redis_host)) { $redis_host = "localhost" }
    
    $redis_port = Read-Host "Redis Port (6379)"
    if ([string]::IsNullOrWhiteSpace($redis_port)) { $redis_port = "6379" }
    
    $env_content = Get-Content .env
    $env_content = $env_content -replace "REDIS_HOST=.*", "REDIS_HOST=$redis_host"
    $env_content = $env_content -replace "REDIS_PORT=.*", "REDIS_PORT=$redis_port"
    $env_content = $env_content -replace "REDIS_URL=.*", "REDIS_URL=redis://$redis_host`:$redis_port/0"
    $env_content | Set-Content .env
    
    Write-Host "✓ Redis configuration updated" -ForegroundColor Green
    Write-Host ""
}

# Google OAuth Configuration (Optional)
$response = Read-Host "Configure Google OAuth? (y/n)"
if ($response -eq "y" -or $response -eq "Y") {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "GOOGLE OAUTH 2.0 CONFIGURATION" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    
    $google_client_id = Read-Host "Google OAuth Client ID"
    $google_client_secret = Read-Host "Google OAuth Client Secret"
    
    $env_content = Get-Content .env
    $env_content = $env_content -replace "GOOGLE_OAUTH_CLIENT_ID=.*", "GOOGLE_OAUTH_CLIENT_ID=$google_client_id"
    $env_content = $env_content -replace "GOOGLE_OAUTH_CLIENT_SECRET=.*", "GOOGLE_OAUTH_CLIENT_SECRET=$google_client_secret"
    $env_content | Set-Content .env
    
    Write-Host "✓ Google OAuth configured" -ForegroundColor Green
    Write-Host ""
}

# Razorpay Configuration (Optional)
$response = Read-Host "Configure Razorpay? (y/n)"
if ($response -eq "y" -or $response -eq "Y") {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "RAZORPAY CONFIGURATION" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    
    $razorpay_key_id = Read-Host "Razorpay Key ID"
    $razorpay_key_secret = Read-Host "Razorpay Key Secret"
    $razorpay_webhook = Read-Host "Razorpay Webhook Secret"
    
    $env_content = Get-Content .env
    $env_content = $env_content -replace "RAZORPAY_KEY_ID=.*", "RAZORPAY_KEY_ID=$razorpay_key_id"
    $env_content = $env_content -replace "RAZORPAY_KEY_SECRET=.*", "RAZORPAY_KEY_SECRET=$razorpay_key_secret"
    $env_content = $env_content -replace "RAZORPAY_WEBHOOK_SECRET=.*", "RAZORPAY_WEBHOOK_SECRET=$razorpay_webhook"
    $env_content | Set-Content .env
    
    Write-Host "✓ Razorpay configured" -ForegroundColor Green
    Write-Host ""
}

# Generate JWT secret
$jwt_secret = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((New-Guid).ToString() + (New-Guid).ToString()))
$env_content = Get-Content .env
$env_content = $env_content -replace "JWT_SECRET=.*", "JWT_SECRET=$jwt_secret"
$env_content | Set-Content .env

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Setup Complete!" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "✓ .env file created with your configuration" -ForegroundColor Green
Write-Host "✓ JWT_SECRET auto-generated" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Review and update .env file as needed"
Write-Host "2. Create PostgreSQL database:"
Write-Host "   createuser $db_user"
Write-Host "   createdb -O $db_user $db_name"
Write-Host "3. Setup Prisma:"
Write-Host "   npx prisma generate"
Write-Host "   npx prisma migrate deploy"
Write-Host "4. Start development server:"
Write-Host "   go run ./cmd/api/main.go"
Write-Host ""
Write-Host "Documentation: See ENV_CONFIGURATION.md" -ForegroundColor Yellow
Write-Host ""
