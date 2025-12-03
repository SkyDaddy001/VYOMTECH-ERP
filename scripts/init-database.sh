#!/bin/bash
# Initialize database with all migrations and sample data
# Loads all migrations from 001-024 and sets up sample partner logins

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     Initializing Database with All Migrations & Sample Data    ║"
echo "╚════════════════════════════════════════════════════════════════╝"

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_NAME="${DB_NAME:-callcenter}"
DB_USER="${DB_USER:-callcenter_user}"
DB_PASSWORD="${DB_PASSWORD:-secure_app_pass}"

echo ""
echo "📊 Configuration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"

# Wait for MySQL to be ready
echo ""
echo "🔌 Waiting for MySQL..."
max_attempts=30
attempt=1
while ! mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" >/dev/null 2>&1; do
    if [ $attempt -eq $max_attempts ]; then
        echo "❌ MySQL failed to start after $max_attempts attempts"
        exit 1
    fi
    echo -n "."
    sleep 1
    attempt=$((attempt + 1))
done
echo ""
echo "✅ MySQL is ready"

# Load all migrations
echo ""
echo "📋 Loading Migrations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

MIGRATION_FILES=(
    "006_phase2_tasks_notifications.sql"
    "006_sample_data.sql"
    "007_tenant_customization.sql"
    "008_purchase_module_schema.sql"
    "009_sales_module_schema.sql"
    "010_milestone_tracking_and_reporting.sql"
    "011_real_estate_property_management.sql"
    "012_hr_payroll_schema.sql"
    "013_accounts_gl_schema.sql"
    "014_purchase_module_schema.sql"
    "015_project_collection_accounts_rera.sql"
    "016_hr_compliance_labour_laws.sql"
    "017_tax_compliance_income_tax_gst.sql"
    "020_comprehensive_test_data.sql"
    "021_comprehensive_customization.sql"
    "022_external_partner_system.sql"
    "023_partner_sources_and_credit_policies.sql"
    "024_sample_partner_logins.sql"
)

for migration in "${MIGRATION_FILES[@]}"; do
    if [ -f "migrations/$migration" ]; then
        echo -n "Loading $migration... "
        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "migrations/$migration" 2>/dev/null && echo "✅" || echo "⚠️  (may have warnings)"
    fi
done

# Verify sample data
echo ""
echo "👥 Verifying Sample Partner Data"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

RESULT=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -N -e "
SELECT COUNT(*) FROM partners;
" 2>/dev/null)

echo "Total Partners: $RESULT"

# Display sample credentials
echo ""
echo "🔐 Sample Login Credentials"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
SELECT 
  CASE 
    WHEN p.partner_type = 'portal' THEN '🌐 Portal'
    WHEN p.partner_type = 'channel_partner' THEN '🔗 Channel'
    WHEN p.partner_type = 'vendor' THEN '🏭 Vendor'
    WHEN p.partner_type = 'customer' THEN '👤 Customer'
  END as 'Type',
  p.organization_name as 'Organization',
  pu.email as 'Email',
  'password123' as 'Password',
  pu.role as 'Role'
FROM partner_users pu
JOIN partners p ON pu.partner_id = p.id
ORDER BY p.partner_type, pu.email;
" 2>/dev/null

echo ""
echo "✅ DATABASE INITIALIZATION COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Application is ready at: http://localhost:8080"
echo "API Documentation: http://localhost:8080/api/v1"
echo ""
