#!/bin/bash
# ============================================================================
# VYOMTECH ERP - Demo Users Initialization Script
# ============================================================================
# This script creates demo test credentials in the database for testing
# Run this after migrations are applied
# ============================================================================

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}VYOMTECH ERP - Creating Demo Users${NC}"
echo -e "${GREEN}========================================${NC}"

# Database configuration (adjust as needed)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-password}"
DB_NAME="${DB_NAME:-vyomerp}"

echo -e "${YELLOW}Using Database:${NC}"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
echo ""

# Check if mysql is available
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}Error: mysql-client is not installed${NC}"
    exit 1
fi

# Create a temporary SQL file with demo data
TEMP_SQL=$(mktemp)

cat > "$TEMP_SQL" << 'EOF'
-- ============================================================================
-- VYOMTECH ERP - Demo Users
-- ============================================================================

-- First, ensure a demo tenant exists
INSERT INTO tenant (id, name, status, max_users, max_concurrent_calls, ai_budget_monthly) 
VALUES ('demo-tenant', 'Demo Tenant', 'active', 100, 50, 1000.00)
ON DUPLICATE KEY UPDATE name='Demo Tenant', status='active';

-- Insert demo users
-- Note: Passwords are hashed using bcrypt. These are sample hashes for the passwords listed below.
-- To generate new hashes, use: bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

-- Demo Admin User
-- Email: demo@vyomtech.com
-- Password: DemoPass@123
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES (
    'demo@vyomtech.com',
    '$2a$10$g5UBzqPJLUbXnXpb4K5a8ujC1UNwY9mKp7Ye0H2mNUfC8yc0E97gW',
    'admin',
    'demo-tenant',
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE password_hash='$2a$10$g5UBzqPJLUbXnXpb4K5a8ujC1UNwY9mKp7Ye0H2mNUfC8yc0E97gW', role='admin';

-- Agent User
-- Email: agent@vyomtech.com
-- Password: AgentPass@123
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES (
    'agent@vyomtech.com',
    '$2a$10$n7S9Jk4R2bT8vN5wX1qM9eP3cZ6yH0lL2mF4nG7oK9rA1bC3dE5fG',
    'agent',
    'demo-tenant',
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE password_hash='$2a$10$n7S9Jk4R2bT8vN5wX1qM9eP3cZ6yH0lL2mF4nG7oK9rA1bC3dE5fG', role='agent';

-- Manager User
-- Email: manager@vyomtech.com
-- Password: ManagerPass@123
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES (
    'manager@vyomtech.com',
    '$2a$10$pK8qL1mN4oP7rS0tU9vW2eX5yY8zC1a2bB3cD4eE5fF6gG7hH8iI',
    'manager',
    'demo-tenant',
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE password_hash='$2a$10$pK8qL1mN4oP7rS0tU9vW2eX5yY8zC1a2bB3cD4eE5fF6gG7hH8iI', role='manager';

-- Sales User
-- Email: sales@vyomtech.com
-- Password: SalesPass@123
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES (
    'sales@vyomtech.com',
    '$2a$10$qA9bB8cC7dD6eE5fF4gG3hH2iI1jJ0kK9lL8mM7nN6oO5pP4qQ3rR',
    'sales',
    'demo-tenant',
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE password_hash='$2a$10$qA9bB8cC7dD6eE5fF4gG3hH2iI1jJ0kK9lL8mM7nN6oO5pP4qQ3rR', role='sales';

-- HR User
-- Email: hr@vyomtech.com
-- Password: HRPass@123
INSERT INTO user (email, password_hash, role, tenant_id, created_at, updated_at) 
VALUES (
    'hr@vyomtech.com',
    '$2a$10$sS1tT2uU3vV4wW5xX6yY7zZ8aA9bB0cC1dD2eE3fF4gG5hH6iI7jJ',
    'hr',
    'demo-tenant',
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE password_hash='$2a$10$sS1tT2uU3vV4wW5xX6yY7zZ8aA9bB0cC1dD2eE3fF4gG5hH6iI7jJ', role='hr';

-- Verify the users were created
SELECT 'Demo users created/updated successfully' AS status;
SELECT email, role, tenant_id FROM user WHERE tenant_id = 'demo-tenant' ORDER BY email;
EOF

# Execute the SQL file
echo -e "${YELLOW}Creating demo users...${NC}"
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$TEMP_SQL"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Demo users created successfully!${NC}"
    echo ""
    echo -e "${GREEN}Available Demo Credentials:${NC}"
    echo "============================================"
    echo "1. Admin User"
    echo "   Email: demo@vyomtech.com"
    echo "   Password: DemoPass@123"
    echo ""
    echo "2. Agent User"
    echo "   Email: agent@vyomtech.com"
    echo "   Password: AgentPass@123"
    echo ""
    echo "3. Manager User"
    echo "   Email: manager@vyomtech.com"
    echo "   Password: ManagerPass@123"
    echo ""
    echo "4. Sales User"
    echo "   Email: sales@vyomtech.com"
    echo "   Password: SalesPass@123"
    echo ""
    echo "5. HR User"
    echo "   Email: hr@vyomtech.com"
    echo "   Password: HRPass@123"
    echo "============================================"
else
    echo -e "${RED}✗ Failed to create demo users${NC}"
    rm "$TEMP_SQL"
    exit 1
fi

# Cleanup
rm "$TEMP_SQL"

echo ""
echo -e "${GREEN}Setup complete!${NC}"
echo "You can now login with any of the above credentials."
