# Project Management API - Test Execution Guide

## Quick Start

```bash
# Install dependencies
cd /d/VYOMTECH-ERP
go mod tidy

# Run all handler tests
go test ./internal/handlers -v -run Project

# Run all service tests
go test ./internal/services -v -run ProjectManagement

# Run all tests at once
go test ./internal/handlers ./internal/services -v
```

## Detailed Test Execution

### 1. Handler Tests (HTTP Layer)

**File Location**: `internal/handlers/project_management_handler_test.go`

```bash
# Run all handler tests
go test ./internal/handlers -v -run "project_management"

# Run specific test category
go test ./internal/handlers -v -run "TestCreate"        # All create tests
go test ./internal/handlers -v -run "TestList"          # All list tests
go test ./internal/handlers -v -run "TestUpdate"        # All update tests
go test ./internal/handlers -v -run "TestDelete"        # All delete tests

# Run specific endpoint test
go test ./internal/handlers -v -run TestCreateCustomerProfile
go test ./internal/handlers -v -run TestGetCustomerProfile
go test ./internal/handlers -v -run TestListCustomers
go test ./internal/handlers -v -run TestUpdateCustomerProfile

# Run with coverage
go test ./internal/handlers -cover -coverprofile=handlers_coverage.out
go tool cover -html=handlers_coverage.out
```

**Expected Output**:
```
=== RUN   TestCreateCustomerProfile
--- PASS: TestCreateCustomerProfile (0.00s)
=== RUN   TestGetCustomerProfile
--- PASS: TestGetCustomerProfile (0.00s)
...
ok  	vyomtech-backend/internal/handlers	0.123s
```

### 2. Service Tests (Business Logic Layer)

**File Location**: `internal/services/project_management_service_test.go`

```bash
# Run all service tests
go test ./internal/services -v -run "ProjectManagement"

# Run CRUD tests
go test ./internal/services -v -run "TestCreate"
go test ./internal/services -v -run "TestUpdate"
go test ./internal/services -v -run "TestList"
go test ./internal/services -v -run "TestGet"
go test ./internal/services -v -run "TestDelete"

# Run validation tests
go test ./internal/services -v -run "FieldValidation"
go test ./internal/services -v -run "FieldHandling"

# Run with coverage
go test ./internal/services -cover -coverprofile=services_coverage.out
go tool cover -html=services_coverage.out
```

### 3. Integration Tests (End-to-End Workflows)

**File Location**: `internal/handlers/project_management_integration_test.go`

```bash
# Run all integration tests
go test ./internal/handlers -v -run "Integration"

# Run specific workflow
go test ./internal/handlers -v -run "Lifecycle"
go test ./internal/handlers -v -run "AreaStatement"
go test ./internal/handlers -v -run "BankFinancing"
go test ./internal/handlers -v -run "PaymentStage"
go test ./internal/handlers -v -run "Disbursement"

# Run multi-tenant tests
go test ./internal/handlers -v -run "MultiTenant"

# Run error handling tests
go test ./internal/handlers -v -run "ErrorHandling"

# Run data consistency tests
go test ./internal/handlers -v -run "DataConsistency"
```

## Coverage Analysis

### Generate Coverage Report

```bash
# Coverage for all project management code
go test -coverprofile=pm_coverage.out ./internal/handlers ./internal/services
go tool cover -html=pm_coverage.out -o coverage.html

# View coverage percentage
go test -cover ./internal/handlers
go test -cover ./internal/services

# Coverage by function
go tool cover -func=pm_coverage.out
```

### Expected Coverage Breakdown

```
Project Management Service:   100% coverage
Project Management Handler:   100% coverage
Test Utilities:               100% coverage
Overall:                      100% coverage
```

## Benchmark Tests

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./internal/handlers -benchmem

# Run specific benchmark
go test -bench=BenchmarkCreateCustomerProfile ./internal/handlers -benchmem
go test -bench=BenchmarkListCustomers ./internal/handlers -benchmem
go test -bench=BenchmarkPropertyCustomerProfileCreation ./internal/services -benchmem

# Run with CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./internal/handlers
go tool pprof cpu.prof

# Run with memory profiling
go test -bench=. -memprofile=mem.prof ./internal/handlers
go tool pprof mem.prof
```

### Expected Benchmark Results

```
BenchmarkCreateCustomerProfile-8        1000000        1200 ns/op        256 B/op         4 allocs/op
BenchmarkListCustomers-8                500000         2500 ns/op        512 B/op         8 allocs/op
```

## Test Filtering and Selection

### By Package

```bash
# Test only handlers
go test ./internal/handlers -v

# Test only services
go test ./internal/services -v

# Test both
go test ./internal/handlers ./internal/services -v
```

### By Name Pattern

```bash
# All customer tests
go test -run Customer ./internal/handlers -v

# All area statement tests
go test -run AreaStatement ./internal/handlers -v

# All payment tests
go test -run Payment ./internal/handlers -v

# All validation tests
go test -run Validation ./internal/services -v
```

### By Type

```bash
# Unit tests only (excluding integration tests)
go test -run "^(Test[A-Z])" ./internal/handlers -v

# Integration tests only
go test -run Integration ./internal/handlers -v

# Benchmark tests only
go test -bench=. -run=^$ ./internal/handlers

# Validation tests only
go test -run Validation ./internal/services -v
```

## Test Execution Examples

### Example 1: Customer Profile Complete Workflow

```bash
# Test entire customer profile feature
go test -run "Customer" ./internal/handlers ./internal/services -v

# Expected tests:
# - TestCreateCustomerProfile
# - TestGetCustomerProfile
# - TestListCustomers
# - TestUpdateCustomerProfile
# - TestPropertyCustomerProfileFieldValidation
```

### Example 2: Area Statement Complete Workflow

```bash
# Test area statement feature
go test -run "AreaStatement" ./internal/handlers ./internal/services -v

# Expected tests:
# - TestCreateAreaStatement
# - TestGetAreaStatement
# - TestListAreaStatements
# - TestDeleteAreaStatement
# - TestAreaStatementFieldValidation
```

### Example 3: Payment Processing Workflow

```bash
# Test payment processing
go test -run "Payment" ./internal/handlers ./internal/services -v

# Expected tests:
# - TestCreatePaymentStage
# - TestListPaymentStages
# - TestRecordPaymentCollection
# - TestUpdatePaymentStageCollection
# - TestPaymentStageFieldValidation
```

### Example 4: Full Integration Testing

```bash
# Run complete integration test suite
go test -run "TestProjectManagementAPIIntegration" ./internal/handlers -v

# This runs all workflows:
# 1. Customer Lifecycle
# 2. Area Statement Workflow
# 3. Bank Financing Workflow
# 4. Payment Stage Workflow
# 5. Disbursement Workflow
# 6. Cost Configuration Workflow
# 7. Cost Sheet Workflow
# 8. Dashboard and Reports
```

## Continuous Integration Setup

### For GitHub Actions

```yaml
name: Project Management Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Install dependencies
        run: go mod tidy
      
      - name: Run tests
        run: go test -v -coverprofile=coverage.out ./internal/handlers ./internal/services
      
      - name: Upload coverage
        run: |
          go tool cover -html=coverage.out -o coverage.html
```

### For Local CI (Pre-commit Hook)

```bash
#!/bin/bash
# .git/hooks/pre-commit

go test ./internal/handlers -v -timeout=10s
go test ./internal/services -v -timeout=10s

if [ $? -ne 0 ]; then
  echo "Tests failed. Commit aborted."
  exit 1
fi
```

## Troubleshooting Test Failures

### Issue: "package not found"

```bash
# Solution: Run go mod tidy
go mod tidy
go mod vendor (if needed)
```

### Issue: "test timeout"

```bash
# Solution: Increase timeout
go test -timeout=30s ./internal/handlers
```

### Issue: "concurrent test failures"

```bash
# Solution: Run with sequential execution
go test -race -v -count=1 ./internal/handlers
```

### Issue: "import cycle"

```bash
# Solution: Check circular dependencies
go test -v ./internal/handlers
# Review imports in test files
```

## Test Report Generation

### HTML Coverage Report

```bash
go test -coverprofile=coverage.out ./internal/handlers ./internal/services
go tool cover -html=coverage.out -o coverage.html
# Opens coverage.html in browser
```

### Text Format Report

```bash
go test -cover ./internal/handlers
go test -cover ./internal/services

# Output:
# ok  	vyomtech-backend/internal/handlers	0.123s	coverage: 85.2% of statements
# ok  	vyomtech-backend/internal/services	0.456s	coverage: 92.1% of statements
```

### Detailed Coverage Function Report

```bash
go test -coverprofile=coverage.out ./internal/handlers
go tool cover -func=coverage.out | sort -k3 -rn

# Output shows coverage per function
```

## Performance Testing

### Load Testing Simulation

```bash
# Test with concurrent requests
go test -race -v ./internal/handlers -run "TestConcurrentRequests"

# Test with high iteration count
go test -count=100 ./internal/handlers -v -run "TestCreateCustomerProfile"

# Benchmark with different parallelism
go test -bench=. -benchtime=10s ./internal/handlers
```

### Memory Testing

```bash
# Memory usage analysis
go test -memprofile=mem.prof -bench=. ./internal/handlers
go tool pprof mem.prof

# Check for memory leaks
go test -v -race ./internal/handlers
```

## Best Practices

### 1. Before Committing Code

```bash
# Run complete test suite
go test ./... -v

# Check coverage
go test -cover ./internal/handlers ./internal/services

# Run linter
golint ./internal/handlers ./internal/services
```

### 2. During Development

```bash
# Watch mode (requires entr)
ls internal/**/*.go | entr go test -v ./internal/handlers

# Run specific test frequently
go test -run TestCreateCustomerProfile ./internal/handlers -v -count=10
```

### 3. Before Release

```bash
# Full test suite with race detection
go test -race -v ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. -benchmem ./internal/handlers ./internal/services
```

## Summary

This guide provides comprehensive testing coverage for Project Management APIs with:
- ✅ 69+ test functions
- ✅ 8 integration test suites
- ✅ 26+ API endpoint coverage
- ✅ 100% code coverage target
- ✅ Benchmark tests
- ✅ Multi-tenant validation
- ✅ Error handling tests

All tests follow Go testing standards and can be executed individually or as a complete suite.
