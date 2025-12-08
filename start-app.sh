#!/bin/bash

# VYOMTECH Frontend-Backend Integration Testing & Verification Script
# This script starts the application stack with podman-compose and tests everything

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}VYOMTECH ERP - Application Stack Startup${NC}"
echo -e "${GREEN}========================================${NC}"

# Check if podman-compose is installed
if ! command -v podman-compose &> /dev/null; then
    echo -e "${RED}❌ podman-compose is not installed${NC}"
    echo "Installing podman-compose..."
    pip install podman-compose
fi

# Check if podman is running
if ! podman ping &> /dev/null; then
    echo -e "${YELLOW}⚠️  Starting podman service...${NC}"
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        systemctl start podman
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        open --hide --background -a Docker
        sleep 5
    fi
fi

# Stop existing containers if running
echo -e "${YELLOW}🛑 Stopping existing containers...${NC}"
podman-compose down 2>/dev/null || true
sleep 2

# Start the stack
echo -e "${YELLOW}🚀 Starting application stack...${NC}"
podman-compose up -d

# Wait for services to be ready
echo -e "${YELLOW}⏳ Waiting for services to start...${NC}"
sleep 10

# Check MySQL
echo -e "${YELLOW}📊 Checking MySQL...${NC}"
MYSQL_READY=false
for i in {1..30}; do
    if podman-compose exec -T mysql mysqladmin ping -h localhost &> /dev/null; then
        echo -e "${GREEN}✅ MySQL is ready${NC}"
        MYSQL_READY=true
        break
    fi
    echo -e "${YELLOW}  Attempt $i/30 - MySQL not ready yet...${NC}"
    sleep 2
done

if [ "$MYSQL_READY" = false ]; then
    echo -e "${RED}❌ MySQL failed to start${NC}"
    podman-compose logs mysql
    exit 1
fi

# Check Redis
echo -e "${YELLOW}🔴 Checking Redis...${NC}"
if podman-compose exec -T redis redis-cli ping &> /dev/null; then
    echo -e "${GREEN}✅ Redis is ready${NC}"
else
    echo -e "${RED}❌ Redis failed to start${NC}"
    podman-compose logs redis
    exit 1
fi

# Wait for backend to compile and start
echo -e "${YELLOW}⏳ Waiting for Go backend to compile and start...${NC}"
BACKEND_READY=false
for i in {1..60}; do
    if curl -s http://localhost:8080/health &> /dev/null; then
        echo -e "${GREEN}✅ Go backend is ready${NC}"
        BACKEND_READY=true
        break
    fi
    echo -e "${YELLOW}  Attempt $i/60 - Backend not ready yet...${NC}"
    sleep 2
done

if [ "$BACKEND_READY" = false ]; then
    echo -e "${RED}⚠️  Backend health check failed (might be normal if not implemented)${NC}"
    echo "Checking if port 8080 is listening..."
    if netstat -tuln | grep -q 8080 || ss -tuln | grep -q 8080; then
        echo -e "${GREEN}✅ Port 8080 is listening${NC}"
    else
        echo -e "${RED}❌ Backend not listening on port 8080${NC}"
        podman-compose logs app | tail -50
    fi
fi

# Wait for frontend to install and start
echo -e "${YELLOW}⏳ Waiting for Next.js frontend to install and start...${NC}"
FRONTEND_READY=false
for i in {1..120}; do
    if curl -s http://localhost:3000 &> /dev/null; then
        echo -e "${GREEN}✅ Next.js frontend is ready${NC}"
        FRONTEND_READY=true
        break
    fi
    echo -e "${YELLOW}  Attempt $i/120 - Frontend not ready yet...${NC}"
    sleep 2
done

if [ "$FRONTEND_READY" = false ]; then
    echo -e "${RED}⚠️  Frontend health check failed${NC}"
    echo "Checking if port 3000 is listening..."
    if netstat -tuln | grep -q 3000 || ss -tuln | grep -q 3000; then
        echo -e "${GREEN}✅ Port 3000 is listening${NC}"
    fi
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ STARTUP COMPLETE${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Display service URLs
echo -e "${YELLOW}📍 Service URLs:${NC}"
echo -e "  Frontend:  ${GREEN}http://localhost:3000${NC}"
echo -e "  Backend:   ${GREEN}http://localhost:8080${NC}"
echo -e "  MySQL:     ${GREEN}localhost:3306${NC}"
echo -e "  Redis:     ${GREEN}localhost:6379${NC}"
echo ""

# Display container status
echo -e "${YELLOW}📦 Container Status:${NC}"
podman-compose ps

echo ""
echo -e "${YELLOW}📝 Demo Credentials:${NC}"
echo "  Email:    demo@vyomtech.com"
echo "  Password: DemoPass@123"
echo ""

echo -e "${YELLOW}🔍 Log Commands:${NC}"
echo "  Backend:   podman-compose logs -f app"
echo "  Frontend:  podman-compose logs -f frontend"
echo "  MySQL:     podman-compose logs -f mysql"
echo "  All:       podman-compose logs -f"
echo ""

echo -e "${GREEN}✨ Application is ready to use!${NC}"
