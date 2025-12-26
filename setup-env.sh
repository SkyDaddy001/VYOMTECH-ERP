#!/bin/bash
# VYOM ERP - Environment Setup Script
# This script helps initialize environment variables for development

set -e

echo "=========================================="
echo "VYOM ERP - Environment Setup"
echo "=========================================="
echo ""

# Check if .env exists
if [ -f .env ]; then
    read -p ".env file already exists. Do you want to overwrite it? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Keeping existing .env file"
        exit 0
    fi
fi

# Copy template
echo "Creating .env from template..."
cp .env.example .env
echo "✓ .env created"
echo ""

# Prompt for database configuration
echo "=========================================="
echo "DATABASE CONFIGURATION"
echo "=========================================="
read -p "Database Host (localhost): " db_host
db_host=${db_host:-localhost}

read -p "Database Port (5432): " db_port
db_port=${db_port:-5432}

read -p "Database User (vyom_user): " db_user
db_user=${db_user:-vyom_user}

read -sp "Database Password: " db_password
echo ""

read -p "Database Name (vyomtech_db): " db_name
db_name=${db_name:-vyomtech_db}

# Update .env with database config
sed -i.bak "s|DB_HOST=.*|DB_HOST=$db_host|" .env
sed -i.bak "s|DB_PORT=.*|DB_PORT=$db_port|" .env
sed -i.bak "s|DB_USER=.*|DB_USER=$db_user|" .env
sed -i.bak "s|DB_NAME=.*|DB_NAME=$db_name|" .env
sed -i.bak "s|DATABASE_URL=.*|DATABASE_URL=postgresql://$db_user:$db_password@$db_host:$db_port/$db_name|" .env

echo "✓ Database configuration updated"
echo ""

# Optional: Configure Redis
read -p "Configure Redis? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "=========================================="
    echo "REDIS CONFIGURATION"
    echo "=========================================="
    read -p "Redis Host (localhost): " redis_host
    redis_host=${redis_host:-localhost}
    
    read -p "Redis Port (6379): " redis_port
    redis_port=${redis_port:-6379}
    
    sed -i.bak "s|REDIS_HOST=.*|REDIS_HOST=$redis_host|" .env
    sed -i.bak "s|REDIS_PORT=.*|REDIS_PORT=$redis_port|" .env
    sed -i.bak "s|REDIS_URL=.*|REDIS_URL=redis://$redis_host:$redis_port/0|" .env
    
    echo "✓ Redis configuration updated"
    echo ""
fi

# Optional: Configure OAuth
read -p "Configure Google OAuth? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "=========================================="
    echo "GOOGLE OAUTH 2.0 CONFIGURATION"
    echo "=========================================="
    read -p "Google OAuth Client ID: " google_client_id
    read -p "Google OAuth Client Secret: " google_client_secret
    
    sed -i.bak "s|GOOGLE_OAUTH_CLIENT_ID=.*|GOOGLE_OAUTH_CLIENT_ID=$google_client_id|" .env
    sed -i.bak "s|GOOGLE_OAUTH_CLIENT_SECRET=.*|GOOGLE_OAUTH_CLIENT_SECRET=$google_client_secret|" .env
    
    echo "✓ Google OAuth configured"
    echo ""
fi

# Optional: Configure Razorpay
read -p "Configure Razorpay? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "=========================================="
    echo "RAZORPAY CONFIGURATION"
    echo "=========================================="
    read -p "Razorpay Key ID: " razorpay_key_id
    read -p "Razorpay Key Secret: " razorpay_key_secret
    read -p "Razorpay Webhook Secret: " razorpay_webhook
    
    sed -i.bak "s|RAZORPAY_KEY_ID=.*|RAZORPAY_KEY_ID=$razorpay_key_id|" .env
    sed -i.bak "s|RAZORPAY_KEY_SECRET=.*|RAZORPAY_KEY_SECRET=$razorpay_key_secret|" .env
    sed -i.bak "s|RAZORPAY_WEBHOOK_SECRET=.*|RAZORPAY_WEBHOOK_SECRET=$razorpay_webhook|" .env
    
    echo "✓ Razorpay configured"
    echo ""
fi

# Generate JWT secret if not set
jwt_secret=$(openssl rand -base64 32)
sed -i.bak "s|JWT_SECRET=.*|JWT_SECRET=$jwt_secret|" .env

echo "=========================================="
echo "Setup Complete!"
echo "=========================================="
echo ""
echo "✓ .env file created with your configuration"
echo "✓ JWT_SECRET auto-generated"
echo ""
echo "Next steps:"
echo "1. Review and update .env file as needed"
echo "2. Create PostgreSQL database:"
echo "   createdb -O $db_user $db_name"
echo "3. Setup Prisma:"
echo "   npx prisma generate"
echo "   npx prisma migrate deploy"
echo "4. Start development server:"
echo "   go run ./cmd/api/main.go"
echo ""
echo "Documentation: See ENV_CONFIGURATION.md"
echo ""

# Cleanup backup files
rm -f .env.bak

exit 0
