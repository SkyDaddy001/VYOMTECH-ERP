# Final comprehensive scan for any remaining type mismatches
# Find all columns that reference id columns but have wrong types

$migrationDir = "d:\VYOMTECH-ERP\migrations"
$files = Get-ChildItem -Path $migrationDir -Filter "*.sql" -File | Sort-Object Name

Write-Host "Comprehensive scan for type mismatches..." -ForegroundColor Green

$issues = @()

foreach ($file in $files) {
    $content = Get-Content -Path $file.FullName -Raw
    $lines = $content -split "`n"
    
    # Look for BIGINT or INT that should be CHAR(36) or CHAR(26)
    for ($i = 0; $i -lt $lines.Count; $i++) {
        $line = $lines[$i]
        
        # Find columns that end with _id and are BIGINT/INT (when not AUTO_INCREMENT or storage)
        if ($line -match '^\s*`?[a-z_]*_id`?\s+(INT|BIGINT)(?!\s+AUTO_INCREMENT)' -and 
            $line -notmatch 'max_file_size' -and 
            $line -notmatch 'display_order' -and
            $line -notmatch 'DEFAULT 0' -and
            $line -notmatch 'COMMENT') {
            $issues += "$($file.Name):$($i+1) - $($line.Trim())"
        }
    }
}

if ($issues.Count -gt 0) {
    Write-Host "Found $($issues.Count) potential type mismatches:" -ForegroundColor Yellow
    $issues | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
} else {
    Write-Host "No type mismatches found!" -ForegroundColor Green
}
