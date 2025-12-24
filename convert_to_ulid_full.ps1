# Full ULID conversion - converts all BIGINT references to CHAR(26)
# Handles: primary keys, foreign keys, and temporary ID columns

$migrationDir = "d:\VYOMTECH-ERP\migrations"
$files = Get-ChildItem -Path $migrationDir -Filter "*.sql" -File

Write-Host "Phase 2: Converting BIGINT foreign key references to CHAR(26)..." -ForegroundColor Green

foreach ($file in $files) {
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content
    
    # Replace BIGINT in foreign key columns (e.g., category_id BIGINT -> category_id CHAR(26))
    # Pattern: column_name BIGINT (when not AUTO_INCREMENT and not already PRIMARY KEY)
    $content = [regex]::Replace($content, '([a-z_]+_id)\s+BIGINT(?!\s+AUTO_INCREMENT)(?!\s+PRIMARY)', '$1 CHAR(26)')
    
    # Replace INT in foreign key columns that are IDs
    $content = [regex]::Replace($content, '([a-z_]+_id)\s+INT(?!\s+DEFAULT)(?!\s+PRIMARY)', '$1 CHAR(26)')
    
    # Handle specific column names that store IDs but might not follow the _id convention
    $content = [regex]::Replace($content, 'parent_id\s+BIGINT', 'parent_id CHAR(26)')
    $content = [regex]::Replace($content, 'parent_id\s+INT', 'parent_id CHAR(26)')
    
    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -Force
        Write-Host "Updated: $($file.Name)" -ForegroundColor Cyan
    }
}

Write-Host "`nPhase 2 complete! All BIGINT references in foreign keys converted to CHAR(26)." -ForegroundColor Green
