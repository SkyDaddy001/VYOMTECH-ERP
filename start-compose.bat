@echo off
REM Start application stack with podman-compose

cd /d "C:\Users\Skydaddy\Desktop\VYOM - ERP"

echo ========================================
echo Starting VYOMTECH ERP Stack
echo ========================================
echo.

podman-compose up -d

echo.
echo ========================================
echo Stack started! Waiting for services...
echo ========================================
echo.

timeout /t 5 /nobreak

echo Checking service status...
podman-compose ps

echo.
echo ========================================
echo Services Info:
echo ========================================
echo MySQL:    localhost:3306
echo Redis:    localhost:6379
echo Backend:  http://localhost:8080
echo Frontend: http://localhost:3000
echo ========================================
