#!/bin/bash

# Migration Initialization Script - Loads all SQL migrations into MySQL

DB_HOST="${DB_HOST:-mysql}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-callcenter_user}"
DB_PASSWORD="${DB_PASSWORD:-secure_app_pass}"
DB_NAME="${DB_NAME:-callcenter}"
MIGRATIONS_DIR="${MIGRATIONS_DIR:-/migrations}"

echo "=================================================="
echo "Database Migration Initialization"
echo "=================================================="

# Wait for MySQL
echo "Waiting for MySQL..."
max_attempts=30
attempt=0
until MYSQL_PWD="$DB_PASSWORD" mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -e "SELECT 1" &> /dev/null; do
    attempt=$((attempt + 1))
    if [ $attempt -ge $max_attempts ]; then
        echo "ERROR: MySQL not ready after $max_attempts attempts"
        exit 1
    fi
    echo "Attempt $attempt/$max_attempts: Waiting..."
    sleep 2
done
echo "✓ MySQL ready"

# Create tenants table
echo "Setting up tenants..."
MYSQL_PWD="$DB_PASSWORD" mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" << 'SQL'
CREATE TABLE IF NOT EXISTS tenants (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    plan VARCHAR(50) DEFAULT 'free',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT IGNORE INTO tenants (id, name) VALUES ('tenant_001', 'Default Tenant');
SQL
echo "✓ Tenants ready"

# Load migrations
echo "Loading migrations..."
LOADED=0
for f in $MIGRATIONS_DIR/006_*.sql $MIGRATIONS_DIR/007_*.sql $MIGRATIONS_DIR/008_*.sql $MIGRATIONS_DIR/009_*.sql $MIGRATIONS_DIR/010_*.sql $MIGRATIONS_DIR/011_*.sql $MIGRATIONS_DIR/012_*.sql $MIGRATIONS_DIR/013_*.sql $MIGRATIONS_DIR/014_*.sql $MIGRATIONS_DIR/015_*.sql $MIGRATIONS_DIR/016_*.sql $MIGRATIONS_DIR/017_*.sql $MIGRATIONS_DIR/020_*.sql $MIGRATIONS_DIR/021_*.sql $MIGRATIONS_DIR/022_*.sql $MIGRATIONS_DIR/023_*.sql; do
    [ -f "$f" ] || continue
    fname=$(basename "$f")
    echo -n "  $fname... "
    if MYSQL_PWD="$DB_PASSWORD" mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" < "$f" 2>&1 | grep -i "ERROR" > /dev/null; then
        echo "✗ SKIPPED"
    else
        echo "✓"
        LOADED=$((LOADED + 1))
    fi
done

echo "Loaded $LOADED migrations"

# Insert sample partners
echo "Creating sample partners..."
MYSQL_PWD="$DB_PASSWORD" mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" "$DB_NAME" << 'SQL'
INSERT IGNORE INTO partners (tenant_id, partner_code, organization_name, partner_type, status, contact_email, contact_person, created_by) 
VALUES 
('tenant_001', 'PORTAL_001', 'PropTech Portal', 'portal', 'active', 'admin@proptech.com', 'Admin', 1),
('tenant_001', 'CHANNEL_001', 'BuildTech Solutions', 'channel_partner', 'active', 'admin@buildtech.in', 'Manager', 1),
('tenant_001', 'VENDOR_001', 'Premium Vendors Inc', 'vendor', 'active', 'admin@vendors.in', 'Lead', 1),
('tenant_001', 'CUSTOMER_001', 'Happy Customers Ltd', 'customer', 'active', 'admin@customers.in', 'Manager', 1);

INSERT IGNORE INTO partner_users (partner_id, tenant_id, email, first_name, last_name, password_hash, role, is_active) 
VALUES 
(1, 'tenant_001', 'admin@proptech.com', 'Portal', 'Admin', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE),
(2, 'tenant_001', 'admin@buildtech.in', 'Channel', 'Admin', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE),
(3, 'tenant_001', 'admin@vendors.in', 'Vendor', 'Admin', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE),
(4, 'tenant_001', 'admin@customers.in', 'Customer', 'Admin', '$2a$10$slYQmyNdGzin7olVN3DOjeBQfmROe/xhbRj9Lqq8K6gGcRjg7gRZm', 'admin', TRUE);
SQL
echo "✓ Sample partners created"

echo ""
echo "=================================================="
echo "Migration Complete! Test Credentials:"
echo "  email: admin@proptech.com (password: password123)"
echo "  email: admin@buildtech.in (password: password123)"
echo "  email: admin@vendors.in (password: password123)"
echo "  email: admin@customers.in (password: password123)"
echo "=================================================="
