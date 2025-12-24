#!/usr/bin/env pwsh
<#
  Phase 11 Integration Testing Script
  Tests all critical endpoints and functionality
#>

$API_URL = "http://localhost:8080"
$ADMINER_URL = "http://localhost:8081"
$results = @()

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Headers = @{},
        [string]$Body = $null,
        [int]$ExpectedStatus = 200
    )
    
    try {
        $params = @{
            Uri = "$API_URL$Endpoint"
            Method = $Method
            Headers = $Headers
            TimeoutSec = 5
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = "application/json"
        }
        
        $response = Invoke-WebRequest @params -ErrorAction Stop
        $status = $response.StatusCode
        $success = $status -eq $ExpectedStatus
        
        $results += @{
            Name = $Name
            Status = if ($success) { "✅ PASS" } else { "❌ FAIL" }
            Details = "HTTP $status (expected $ExpectedStatus)"
        }
        
        return $success, $response
    }
    catch {
        $results += @{
            Name = $Name
            Status = "❌ FAIL"
            Details = "Error: $($_.Exception.Message)"
        }
        return $false, $null
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "🧪 PHASE 11 INTEGRATION TESTING" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# ============================================================
# TEST 1: HEALTH & BASIC CONNECTIVITY
# ============================================================
Write-Host "1️⃣  HEALTH & BASIC CONNECTIVITY" -ForegroundColor Yellow
Write-Host "================================" -ForegroundColor Yellow

Test-Endpoint -Name "API Health Check" -Method GET -Endpoint "/health" | Out-Null

Write-Host ""

# ============================================================
# TEST 2: MOCK OAUTH ENDPOINTS
# ============================================================
Write-Host "2️⃣  OAUTH ENDPOINTS" -ForegroundColor Yellow
Write-Host "==================" -ForegroundColor Yellow

Test-Endpoint -Name "OAuth Test Data" -Method GET -Endpoint "/mock/oauth/test-data" | Out-Null
Test-Endpoint -Name "OAuth Docs" -Method GET -Endpoint "/mock/oauth/docs" | Out-Null
Test-Endpoint -Name "OAuth Mock Authorize" -Method POST -Endpoint "/mock/oauth/google/authorize" `
    -Body '{"redirect_uri":"http://localhost:3000/callback"}' | Out-Null

Write-Host ""

# ============================================================
# TEST 3: TEST DATA ENDPOINTS
# ============================================================
Write-Host "3️⃣  TEST DATA ENDPOINTS" -ForegroundColor Yellow
Write-Host "=======================" -ForegroundColor Yellow

Test-Endpoint -Name "Test Data List" -Method GET -Endpoint "/test/data/" | Out-Null
Test-Endpoint -Name "Google Ads Test Data" -Method GET -Endpoint "/test/data/google-ads/sample" | Out-Null

Write-Host ""

# ============================================================
# TEST 4: SYNC JOB ENDPOINTS
# ============================================================
Write-Host "4️⃣  SYNC JOB MONITORING" -ForegroundColor Yellow
Write-Host "=======================" -ForegroundColor Yellow

Test-Endpoint -Name "Sync Status" -Method GET -Endpoint "/api/v1/sync/status" | Out-Null
Test-Endpoint -Name "Sync Jobs List" -Method GET -Endpoint "/api/v1/sync/jobs" | Out-Null
Test-Endpoint -Name "Sync Job Details (Google)" -Method GET -Endpoint "/api/v1/sync/jobs/sync-google-ads-metrics" | Out-Null

Write-Host ""

# ============================================================
# TEST 5: SYNC TRIGGERS
# ============================================================
Write-Host "5️⃣  MANUAL SYNC TRIGGERS" -ForegroundColor Yellow
Write-Host "========================" -ForegroundColor Yellow

Test-Endpoint -Name "Trigger Metrics Sync" -Method POST -Endpoint "/api/v1/sync/metrics/sync-now" `
    -ExpectedStatus 200 | Out-Null
    
Start-Sleep -Seconds 1

Test-Endpoint -Name "Trigger Attribution" -Method POST -Endpoint "/api/v1/sync/attribution/process-now" `
    -ExpectedStatus 200 | Out-Null

Test-Endpoint -Name "Trigger ROI Calculation" -Method POST -Endpoint "/api/v1/sync/roi/calculate-now" `
    -ExpectedStatus 200 | Out-Null

Write-Host ""

# ============================================================
# TEST 6: GOOGLE ADS ENDPOINTS
# ============================================================
Write-Host "6️⃣  GOOGLE ADS API ENDPOINTS" -ForegroundColor Yellow
Write-Host "=============================" -ForegroundColor Yellow

Test-Endpoint -Name "Create Google Campaign" -Method POST -Endpoint "/api/v1/google-ads/campaigns" `
    -Body '{"name":"Test Campaign","budget":1000}' -ExpectedStatus 200 | Out-Null

Test-Endpoint -Name "List Google Campaigns" -Method GET -Endpoint "/api/v1/google-ads/campaigns" | Out-Null

Test-Endpoint -Name "Sync Google Metrics" -Method POST -Endpoint "/api/v1/google-ads/sync/metrics" `
    -ExpectedStatus 200 | Out-Null

Write-Host ""

# ============================================================
# TEST 7: META ADS ENDPOINTS
# ============================================================
Write-Host "7️⃣  META ADS API ENDPOINTS" -ForegroundColor Yellow
Write-Host "==========================" -ForegroundColor Yellow

Test-Endpoint -Name "Create Meta Campaign" -Method POST -Endpoint "/api/v1/meta-ads/campaigns" `
    -Body '{"name":"Meta Test Campaign","budget":1000}' -ExpectedStatus 200 | Out-Null

Test-Endpoint -Name "List Meta Campaigns" -Method GET -Endpoint "/api/v1/meta-ads/campaigns" | Out-Null

Test-Endpoint -Name "Sync Meta Metrics" -Method POST -Endpoint "/api/v1/meta-ads/sync/metrics" `
    -ExpectedStatus 200 | Out-Null

Write-Host ""

# ============================================================
# TEST 8: ROI & ANALYTICS
# ============================================================
Write-Host "8️⃣  ROI & ANALYTICS ENDPOINTS" -ForegroundColor Yellow
Write-Host "==============================" -ForegroundColor Yellow

Test-Endpoint -Name "Calculate Daily ROI" -Method POST -Endpoint "/api/v1/roi/daily" `
    -Body '{"start_date":"2025-12-01","end_date":"2025-12-24"}' -ExpectedStatus 200 | Out-Null

Test-Endpoint -Name "Calculate Period ROI" -Method POST -Endpoint "/api/v1/roi/period" `
    -Body '{"platform":"google","days":7}' -ExpectedStatus 200 | Out-Null

Test-Endpoint -Name "ROI Dashboard" -Method GET -Endpoint "/api/v1/roi/dashboard" | Out-Null

Write-Host ""

# ============================================================
# TEST 9: PHASE 8 TEST ENDPOINTS
# ============================================================
Write-Host "9️⃣  PHASE 8 TEST ENDPOINTS" -ForegroundColor Yellow
Write-Host "===========================" -ForegroundColor Yellow

Test-Endpoint -Name "Phase 8 Docs" -Method GET -Endpoint "/test/phase-8/docs" | Out-Null
Test-Endpoint -Name "Phase 8 Quick Start" -Method GET -Endpoint "/test/phase-8/quick-start" | Out-Null
Test-Endpoint -Name "Phase 8 Status" -Method GET -Endpoint "/test/phase-8/status" | Out-Null

Write-Host ""

# ============================================================
# TEST 10: DATABASE CONNECTIVITY
# ============================================================
Write-Host "🔟 DATABASE CONNECTIVITY" -ForegroundColor Yellow
Write-Host "========================" -ForegroundColor Yellow

try {
    $db = Invoke-WebRequest -Uri "http://localhost:8081" -TimeoutSec 5 -ErrorAction Stop
    $results += @{
        Name = "Adminer Database Admin"
        Status = "✅ PASS"
        Details = "Database admin accessible"
    }
}
catch {
    $results += @{
        Name = "Adminer Database Admin"
        Status = "❌ FAIL"
        Details = "Cannot access database admin"
    }
}

Write-Host ""

# ============================================================
# RESULTS SUMMARY
# ============================================================
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "📊 TEST RESULTS SUMMARY" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$passed = ($results | Where-Object { $_.Status -like "*PASS*" }).Count
$failed = ($results | Where-Object { $_.Status -like "*FAIL*" }).Count
$total = $results.Count

$results | Format-Table -AutoSize @(
    @{ Label = "Test Name"; Expression = { $_.Name } },
    @{ Label = "Result"; Expression = { $_.Status } },
    @{ Label = "Details"; Expression = { $_.Details } }
)

Write-Host ""
Write-Host "SUMMARY: $passed/$total tests passed" -ForegroundColor $(if ($passed -eq $total) { "Green" } else { "Red" })
Write-Host ""

if ($passed -eq $total) {
    Write-Host "✅ ALL TESTS PASSED! System is ready for integration." -ForegroundColor Green
    exit 0
}
else {
    Write-Host "⚠️  $failed tests failed. Review logs above." -ForegroundColor Yellow
    exit 1
}
