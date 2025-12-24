#!/bin/bash

# VYOM LMS - Development Startup Script
# This script sets up and runs the entire MVP system

set -e  # Exit on error

echo "🚀 VYOM LMS - MVP Development Startup"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

# Check dependencies
echo -e "${BLUE}1. Checking dependencies...${NC}"
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi
echo "✓ Go $(go version | awk '{print $3}')"

if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed"
    exit 1
fi
echo "✓ Node.js $(node --version)"

if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed"
    exit 1
fi
echo "✓ Docker installed"

echo ""
echo -e "${BLUE}2. Starting database services...${NC}"
docker-compose up -d
echo "✓ Database services started"

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 5

echo ""
echo -e "${BLUE}3. Generating Prisma client...${NC}"
go run github.com/prisma/prisma-client-go generate
echo "✓ Prisma client generated"

echo ""
echo -e "${BLUE}4. Running database migrations...${NC}"
go run github.com/prisma/prisma-client-go migrate deploy
echo "✓ Database migrations completed"

echo ""
echo -e "${BLUE}5. Building backend...${NC}"
go build -o bin/api ./cmd/api
echo "✓ Backend built successfully"

echo ""
echo -e "${YELLOW}========================================${NC}"
echo -e "${GREEN}✅ Setup Complete!${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""
echo "To start the system:"
echo ""
echo "  Terminal 1 - Backend:"
echo "    go run ./cmd/api/main.go"
echo ""
echo "  Terminal 2 - Frontend:"
echo "    cd frontend && npm run dev"
echo ""
echo "Then visit: http://localhost:3000"
echo ""
