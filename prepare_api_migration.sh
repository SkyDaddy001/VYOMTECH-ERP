#!/bin/bash
# Script to help identify and prepare Go model files for ULID migration

echo "=== VyomTech ERP - Go Models ULID Migration Preparation ==="
echo

# Find all model files with int/int64 ID fields
echo "Files that need updating (have int/int64 ID fields):"
echo "======================================================="
grep -l "ID.*int64.*json.*id" internal/models/*.go | while read file; do
  echo "✓ $file"
  # Count how many int64 ID fields
  count=$(grep -c "ID.*int64" "$file")
  echo "  Found $count fields needing update"
done
echo

# Find specific patterns
echo "Critical Models to Update (highest priority):"
echo "=============================================="
echo "1. internal/models/partner.go - 10+ structs with int64 IDs"
echo "2. internal/models/phase1_models.go - 10+ structs with int64 IDs"
echo "3. internal/models/integration.go - int64 IDs"
echo

# Check for Password field references
echo "Models using Password instead of PasswordHash:"
grep -l "Password.*string" internal/models/*.go

echo

echo "RECOMMENDED APPROACH:"
echo "===================="
echo "1. Run Go compiler to see specific errors: go build ./..."
echo "2. Update one file at a time, using regex replacements"
echo "3. For each model file:"
echo "   a. Replace 'ID.*int64' with 'ID string'"
echo "   b. Replace '*int64' ID fields with '*string'"
echo "   c. Update foreign key references (e.g., UserID, PartnerID)"
echo "4. Verify: go build ./... && go test ./..."
echo

# List all model files
echo "All model files to review:"
ls -1 internal/models/*.go
