@echo off
REM Initialize database with all migrations and sample data
REM Loads all migrations from 001-024 and sets up sample partner logins

setlocal enabledelayedexpansion

REM Configuration
set DB_HOST=%DB_HOST:localhost=localhost%
set DB_PORT=%DB_PORT:3306=3306%
set DB_NAME=%DB_NAME:callcenter=callcenter%
set DB_USER=%DB_USER:callcenter_user=callcenter_user%
set DB_PASSWORD=%DB_PASSWORD:secure_app_pass=secure_app_pass%

echo.
echo ╔════════════════════════════════════════════════════════════════╗
echo ║    Initializing Database with All Migrations ^& Sample Data    ║
echo ╚════════════════════════════════════════════════════════════════╝
echo.

echo 📊 Configuration
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo Host: %DB_HOST%:%DB_PORT%
echo Database: %DB_NAME%

REM Check MySQL availability
echo.
echo 🔌 Checking MySQL connection...
mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% -e "SELECT 1" >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ MySQL not available. Please check connection.
    exit /b 1
)
echo ✅ MySQL is ready

echo.
echo 📋 Loading Migrations
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

setlocal enabledelayedexpansion
set migrations[0]=006_phase2_tasks_notifications.sql
set migrations[1]=006_sample_data.sql
set migrations[2]=007_tenant_customization.sql
set migrations[3]=008_purchase_module_schema.sql
set migrations[4]=009_sales_module_schema.sql
set migrations[5]=010_milestone_tracking_and_reporting.sql
set migrations[6]=011_real_estate_property_management.sql
set migrations[7]=012_hr_payroll_schema.sql
set migrations[8]=013_accounts_gl_schema.sql
set migrations[9]=014_purchase_module_schema.sql
set migrations[10]=015_project_collection_accounts_rera.sql
set migrations[11]=016_hr_compliance_labour_laws.sql
set migrations[12]=017_tax_compliance_income_tax_gst.sql
set migrations[13]=020_comprehensive_test_data.sql
set migrations[14]=021_comprehensive_customization.sql
set migrations[15]=022_external_partner_system.sql
set migrations[16]=023_partner_sources_and_credit_policies.sql
set migrations[17]=024_sample_partner_logins.sql

for /l %%i in (0,1,17) do (
    if exist "migrations\!migrations[%%i]!" (
        echo Loading !migrations[%%i]!...
        mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% < migrations\!migrations[%%i]! 2>nul
        if !errorlevel! equ 0 (
            echo   ✅
        ) else (
            echo   ⚠️ (may have warnings^)
        )
    )
)

echo.
echo 👥 Verifying Sample Partner Data
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% -e "SELECT COUNT(*) as 'Total Partners' FROM partners;" 2>nul

echo.
echo 🔐 Sample Login Credentials
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo.

mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% -e ^
"SELECT ^
  CASE ^
    WHEN p.partner_type = 'portal' THEN '🌐 Portal' ^
    WHEN p.partner_type = 'channel_partner' THEN '🔗 Channel' ^
    WHEN p.partner_type = 'vendor' THEN '🏭 Vendor' ^
    WHEN p.partner_type = 'customer' THEN '👤 Customer' ^
  END as 'Type', ^
  p.organization_name as 'Organization', ^
  pu.email as 'Email', ^
  'password123' as 'Password', ^
  pu.role as 'Role' ^
FROM partner_users pu ^
JOIN partners p ON pu.partner_id = p.id ^
ORDER BY p.partner_type, pu.email;" 2>nul

echo.
echo ✅ DATABASE INITIALIZATION COMPLETE
echo ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
echo.
echo Application is ready at: http://localhost:8080
echo API Documentation: http://localhost:8080/api/v1
echo.
pause
