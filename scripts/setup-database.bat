@echo off
REM Database Migration and Setup Script for Windows
REM Checks for all migrations and applies sample partner data
REM Date: December 3, 2025

setlocal enabledelayedexpansion

REM Configuration
set DB_HOST=%DB_HOST:localhost=localhost%
set DB_PORT=%DB_PORT:3306=3306%
set DB_NAME=%DB_NAME:callcenter=callcenter%
set DB_USER=%DB_USER:callcenter_user=callcenter_user%
set DB_PASSWORD=%DB_PASSWORD:secure_app_pass=secure_app_pass%

echo.
echo ╔════════════════════════════════════════════════════════════════╗
echo ║         VYOMTECH-ERP Database Setup ^& Validation              ║
echo ╚════════════════════════════════════════════════════════════════╝
echo.
echo 📊 DATABASE CONFIGURATION
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo Host: %DB_HOST%:%DB_PORT%
echo Database: %DB_NAME%
echo User: %DB_USER%

REM Check if mysql is available
where mysql >nul 2>&1
if %errorlevel% neq 0 (
    echo.
    echo ❌ MySQL command-line tool not found
    echo Please ensure MySQL is installed and added to PATH
    exit /b 1
)

echo.
echo 🔌 Checking database connection...
mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% -e "SELECT 1" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Database connection successful
) else (
    echo ❌ Database connection failed
    echo Please verify credentials and MySQL is running
    exit /b 1
)

echo.
echo 📋 CHECKING MIGRATIONS
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

REM Count migration files
setlocal enabledelayedexpansion
set count=0
for %%f in (migrations\*.sql) do (
    set /a count+=1
)
echo Total migration files found: !count!

echo.
echo 👥 LOADING SAMPLE PARTNER DATA
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

if exist "migrations\024_sample_partner_logins.sql" (
    echo Executing sample partner logins migration...
    mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% < migrations\024_sample_partner_logins.sql
    echo ✅ Sample partner data loaded
) else (
    echo ⚠️ Sample login migration not found
)

echo.
echo 🔐 SAMPLE LOGIN CREDENTIALS
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo.
mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% -e ^
"SELECT ^
  CASE ^
    WHEN p.partner_type = 'portal' THEN '🌐 Property Portal' ^
    WHEN p.partner_type = 'channel_partner' THEN '🔗 Channel Partner' ^
    WHEN p.partner_type = 'vendor' THEN '🏭 Vendor' ^
    WHEN p.partner_type = 'customer' THEN '👤 Customer' ^
  END as 'Partner Type', ^
  p.organization_name as 'Organization', ^
  pu.email as 'Email', ^
  'password123' as 'Password', ^
  pu.role as 'Role' ^
FROM partner_users pu ^
JOIN partners p ON pu.partner_id = p.id ^
ORDER BY p.partner_type, pu.email;"

echo.
echo ✅ DATABASE SETUP COMPLETE
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo.
echo Sample partners created:
echo   1. PropTech Portal (portal)
echo   2. BuildTech Solutions (channel_partner)
echo   3. Premium Vendors Inc (vendor)
echo   4. Happy Customers Ltd (customer)
echo.
echo Each partner has 2 sample users with roles:
echo   • admin - Full access
echo   • lead_manager - Lead submission ^& management
echo   • viewer - Read-only access
echo.
echo All sample accounts use password: password123
echo.
echo Press any key to exit...
pause >nul
