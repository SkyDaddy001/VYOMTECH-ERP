#!/bin/bash
# Database Migration and Setup Script
# Checks for all migrations and applies sample partner data
# Date: December 3, 2025

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         VYOMTECH-ERP Database Setup & Validation              ║"
echo "╚════════════════════════════════════════════════════════════════╝"

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-callcenter}"
DB_USER="${DB_USER:-callcenter_user}"
DB_PASSWORD="${DB_PASSWORD:-secure_app_pass}"

echo ""
echo "📊 DATABASE CONFIGURATION"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"

# Function to run SQL
run_sql() {
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "$1"
}

# Check connection
echo ""
echo "🔌 Checking database connection..."
if mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" >/dev/null 2>&1; then
    echo "✅ Database connection successful"
else
    echo "❌ Database connection failed"
    exit 1
fi

# Check existing migrations
echo ""
echo "📋 CHECKING MIGRATIONS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

MIGRATION_COUNT=$(ls -1 migrations/*.sql 2>/dev/null | wc -l)
echo "Total migration files found: $MIGRATION_COUNT"

# Check critical tables
echo ""
echo "📊 CHECKING TABLES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

TABLES=$(run_sql "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA='$DB_NAME'")
echo "Total tables in database: $TABLES"

# Check key tables
REQUIRED_TABLES=("partners" "partner_users" "partner_leads" "partner_payouts" "partner_sources" "partner_credit_policies")

for table in "${REQUIRED_TABLES[@]}"; do
    TABLE_EXISTS=$(run_sql "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA='$DB_NAME' AND TABLE_NAME='$table'")
    if [ "$TABLE_EXISTS" -eq 1 ]; then
        COUNT=$(run_sql "SELECT COUNT(*) FROM $table")
        echo "✅ $table: $COUNT records"
    else
        echo "⚠️ $table: NOT FOUND"
    fi
done

# Apply sample partner login data
echo ""
echo "👥 LOADING SAMPLE PARTNER DATA"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -f "migrations/024_sample_partner_logins.sql" ]; then
    echo "Executing sample partner logins migration..."
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < migrations/024_sample_partner_logins.sql
    echo "✅ Sample partner data loaded"
else
    echo "⚠️ Sample login migration not found"
fi

# Display sample credentials
echo ""
echo "🔐 SAMPLE LOGIN CREDENTIALS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

run_sql "
SELECT 
  CASE 
    WHEN p.partner_type = 'portal' THEN '🌐 Property Portal'
    WHEN p.partner_type = 'channel_partner' THEN '🔗 Channel Partner'
    WHEN p.partner_type = 'vendor' THEN '🏭 Vendor'
    WHEN p.partner_type = 'customer' THEN '👤 Customer'
  END as 'Partner Type',
  p.organization_name as 'Organization',
  pu.email as 'Email',
  'password123' as 'Password',
  pu.role as 'Role'
FROM partner_users pu
JOIN partners p ON pu.partner_id = p.id
ORDER BY p.partner_type, pu.email;
"

# Summary
echo ""
echo "✅ DATABASE SETUP COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Sample partners created:"
echo "  1. PropTech Portal (portal)"
echo "  2. BuildTech Solutions (channel_partner)"
echo "  3. Premium Vendors Inc (vendor)"
echo "  4. Happy Customers Ltd (customer)"
echo ""
echo "Each partner has 2 sample users with roles:"
echo "  • admin - Full access"
echo "  • lead_manager - Lead submission & management"
echo "  • viewer - Read-only access"
echo ""
echo "All sample accounts use password: password123"
echo ""
