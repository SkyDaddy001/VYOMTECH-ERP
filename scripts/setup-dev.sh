#!/bin/bash

# Multi-Tenant AI Call Center Development Environment Setup
# This script sets up the complete development environment

set -e

echo "🚀 Setting up Multi-Tenant AI Call Center Development Environment"

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | cut -d'o' -f2)
if [[ "$(printf '%s\n' "$GO_VERSION" "1.21" | sort -V | head -n1)" != "1.21" ]]; then
    echo "❌ Go version $GO_VERSION is not supported. Please upgrade to Go 1.21+"
    exit 1
fi

echo "✅ Prerequisites check passed"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p logs
mkdir -p data/mysql
mkdir -p data/redis
mkdir -p .github/workflows
mkdir -p k8s
mkdir -p monitoring
mkdir -p docs

# Start infrastructure services
echo "🐳 Starting infrastructure services..."
docker-compose up -d mysql redis

# Wait for MySQL to be ready
echo "⏳ Waiting for MySQL to be ready..."
sleep 30

# Check if MySQL is ready
if ! docker-compose exec -T mysql mysqladmin ping -h localhost --silent; then
    echo "❌ MySQL is not ready. Please check the logs."
    docker-compose logs mysql
    exit 1
fi

echo "✅ MySQL is ready"

# Initialize database
echo "🗄️ Initializing database..."
if [ -f "migrations/001_initial_schema.sql" ]; then
    docker-compose exec -T mysql mysql -u root -prootpassword callcenter < migrations/001_initial_schema.sql
    echo "✅ Database initialized"
else
    echo "⚠️ Database migration file not found. Please ensure migrations/001_initial_schema.sql exists."
fi

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy
go mod download

# Build the application
echo "🔨 Building application..."
go build -o bin/callcenter ./cmd

# Run tests
echo "🧪 Running tests..."
go test ./... -v

# Create default admin user
echo "👤 Creating default admin user..."
cat << EOF | docker-compose exec -T mysql mysql -u root -prootpassword callcenter
INSERT IGNORE INTO tenant (id, name, domain, status) VALUES
('default-tenant', 'Default Tenant', 'default.callcenter.com', 'active');

INSERT IGNORE INTO user (email, password_hash, role, tenant_id, first_name, last_name) VALUES
('admin@callcenter.com', '\$2a\$10\$8K3VZ6Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Yx5Z8Y', 'admin', 'default-tenant', 'System', 'Administrator');
EOF

echo "✅ Default admin user created (email: admin@callcenter.com, password: admin123)"

# Generate JWT secret if not exists
if [ ! -f ".env" ]; then
    echo "🔐 Generating JWT secret..."
    JWT_SECRET=$(openssl rand -hex 32)
    cat << EOF > .env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_NAME=callcenter
DB_USER=callcenter_user
DB_PASSWORD=callcenter_pass

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT Configuration
JWT_SECRET=$JWT_SECRET

# OpenAI Configuration (set your API key)
OPENAI_API_KEY=your-openai-api-key-here

# Application Configuration
LOG_LEVEL=debug
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=60
MAX_CONCURRENT_CALLS=10
SESSION_TIMEOUT=3600
PASSWORD_RESET_EXPIRY=3600

# Email Configuration (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@yourdomain.com
EOF
    echo "✅ Environment file created. Please update OPENAI_API_KEY and email settings."
fi

# Create development scripts
echo "📜 Creating development scripts..."

cat << 'EOF' > scripts/dev-run.sh
#!/bin/bash
# Development run script

export $(grep -v '^#' .env | xargs)
go run cmd/main.go
EOF

cat << 'EOF' > scripts/test.sh
#!/bin/bash
# Test script

echo "Running unit tests..."
go test ./internal/... -v

echo "Running integration tests..."
go test ./tests/... -v -tags=integration

echo "Generating coverage report..."
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
echo "Coverage report generated: coverage.html"
EOF

cat << 'EOF' > scripts/migrate.sh
#!/bin/bash
# Database migration script

if [ -z "$1" ]; then
    echo "Usage: $0 <migration-file>"
    exit 1
fi

MIGRATION_FILE=$1

if [ ! -f "$MIGRATION_FILE" ]; then
    echo "Migration file $MIGRATION_FILE not found"
    exit 1
fi

echo "Running migration: $MIGRATION_FILE"
docker-compose exec -T mysql mysql -u root -prootpassword callcenter < $MIGRATION_FILE
echo "Migration completed"
EOF

# Make scripts executable
chmod +x scripts/*.sh

echo ""
echo "🎉 Development environment setup complete!"
echo ""
echo "📋 Next steps:"
echo "1. Update .env file with your OpenAI API key and email settings"
echo "2. Start the application: ./scripts/dev-run.sh"
echo "3. Access the API at http://localhost:8080"
echo "4. Default admin login: admin@callcenter.com / admin123"
echo ""
echo "🔧 Useful commands:"
echo "- View logs: docker-compose logs -f"
echo "- Run tests: ./scripts/test.sh"
echo "- Run migrations: ./scripts/migrate.sh <file.sql>"
echo "- Stop services: docker-compose down"
echo ""
echo "📚 Documentation:"
echo "- API docs: Check README.md"
echo "- Deployment: docs/deployment-guide.md"
echo ""
echo "Happy coding! 🚀"
