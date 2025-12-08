@echo off
REM Service verification script for VYOMTECH ERP Stack

cd /d "C:\Users\Skydaddy\Desktop\VYOM - ERP"

echo.
echo ========================================
echo VYOMTECH ERP - SERVICE VERIFICATION
echo ========================================
echo.

echo [1/5] Checking Container Status...
podman-compose ps
echo.

echo [2/5] Checking MySQL...
podman-compose exec mysql mysql -u callcenter_user -psecure_app_pass -e "SELECT 1;" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✓ MySQL is healthy
) else (
    echo ✗ MySQL not responding - checking logs
    podman-compose logs --tail=20 mysql
)
echo.

echo [3/5] Checking Redis...
podman-compose exec redis redis-cli PING >nul 2>&1
if %errorlevel% equ 0 (
    echo ✓ Redis is healthy
) else (
    echo ✗ Redis not responding
    podman-compose logs --tail=20 redis
)
echo.

echo [4/5] Checking Backend API...
curl -s -o /dev/null -w "Backend Status: %%{http_code}\n" http://localhost:8080/api/v1/health
echo.

echo [5/5] Checking Frontend...
curl -s -o /dev/null -w "Frontend Status: %%{http_code}\n" http://localhost:3000
echo.

echo ========================================
echo SERVICES READY!
echo ========================================
echo.
echo Access Points:
echo   Frontend:  http://localhost:3000
echo   Backend:   http://localhost:8080
echo   Database:  localhost:3306
echo   Redis:     localhost:6379
echo.
echo Demo Credentials:
echo   Email:     demo@vyomtech.com
echo   Password:  DemoPass@123
echo.
