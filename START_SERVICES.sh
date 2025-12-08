#!/bin/bash

# VYOMTECH ERP - Complete Startup Script
# This script starts all services for local development

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo "📦 VYOMTECH ERP - Starting All Services"
echo "📍 Project Directory: $PROJECT_DIR"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Docker is running
if ! docker ps > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

echo -e "${BLUE}1. Starting Docker Compose Services${NC}"
echo "   - MySQL Database"
echo "   - Redis Cache"
echo "   - Backend API"
echo ""

cd "$PROJECT_DIR"
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 5

# Check if services are running
echo -e "${BLUE}2. Checking Service Status${NC}"

# Check MySQL
if docker-compose ps mysql | grep -q "Up"; then
    echo -e "${GREEN}✅ MySQL is running on port 3306${NC}"
else
    echo -e "${RED}❌ MySQL failed to start${NC}"
fi

# Check Redis
if docker-compose ps redis | grep -q "Up"; then
    echo -e "${GREEN}✅ Redis is running on port 6379${NC}"
else
    echo -e "${RED}❌ Redis failed to start${NC}"
fi

# Check Backend
if docker-compose ps app | grep -q "Up"; then
    echo -e "${GREEN}✅ Backend API is running on port 8080${NC}"
else
    echo -e "${RED}❌ Backend API failed to start${NC}"
fi

# Check Frontend
if docker-compose ps frontend | grep -q "Up"; then
    echo -e "${GREEN}✅ Frontend is running on port 3000${NC}"
else
    echo -e "${RED}❌ Frontend failed to start${NC}"
fi

echo ""
echo -e "${BLUE}3. Service URLs${NC}"
echo -e "   ${YELLOW}Frontend:${NC} http://localhost:3000"
echo -e "   ${YELLOW}Backend API:${NC} http://localhost:8080"
echo -e "   ${YELLOW}Database:${NC} localhost:3306 (MySQL)"
echo -e "   ${YELLOW}Cache:${NC} localhost:6379 (Redis)"
echo ""

echo -e "${BLUE}4. Test Connectivity${NC}"

# Test Backend Health
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is responding${NC}"
else
    echo -e "${YELLOW}⏳ Backend is starting, trying again in 3 seconds...${NC}"
    sleep 3
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend is now responding${NC}"
    else
        echo -e "${RED}❌ Backend is not responding${NC}"
    fi
fi

# Test Frontend
if curl -s http://localhost:3000 > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend is responding${NC}"
else
    echo -e "${YELLOW}⏳ Frontend is starting, trying again in 3 seconds...${NC}"
    sleep 3
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Frontend is now responding${NC}"
    else
        echo -e "${RED}❌ Frontend is not responding${NC}"
    fi
fi

echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ VYOMTECH ERP is ready!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════${NC}"
echo ""
echo -e "${YELLOW}Demo Credentials:${NC}"
echo "  Email: demo@vyomtech.com"
echo "  Password: DemoPass@123"
echo ""
echo -e "${YELLOW}Quick Commands:${NC}"
echo "  View logs:      docker-compose logs -f app"
echo "  Stop services:  docker-compose down"
echo "  Restart:        docker-compose restart"
echo "  Database shell: docker-compose exec mysql mysql -udemo -pdemo123 vyomtech_demo"
echo ""
echo -e "${BLUE}📖 Documentation:${NC}"
echo "  - Setup Guide: $PROJECT_DIR/frontend/DASHBOARD_SETUP.md"
echo "  - Implementation: $PROJECT_DIR/CUSTOM_ERP_DASHBOARD_COMPLETE.md"
echo "  - Verification: $PROJECT_DIR/FINAL_DASHBOARD_VERIFICATION.md"
echo ""
