# PowerShell script to add SET FOREIGN_KEY_CHECKS statements

$migrations_dir = "d:\VYOMTECH-ERP\migrations"

# Migrations 035-043 that need the statements
$target_migrations = @(
    "035_document_management.sql",
    "036_site_visit_management.sql",
    "037_possession_management.sql",
    "038_title_clearance_management.sql",
    "039_customer_portal.sql",
    "040_advanced_analytics.sql",
    "041_mobile_app_features.sql",
    "042_ai_powered_recommendations.sql",
    "043_integration_hub.sql"
)

foreach ($migration_file in $target_migrations) {
    $file_path = Join-Path $migrations_dir $migration_file
    
    if (-not (Test-Path $file_path)) {
        Write-Host "File not found: $migration_file"
        continue
    }
    
    Write-Host "Processing: $migration_file"
    $content = Get-Content -Path $file_path -Raw
    
    # Check if it already has SET FOREIGN_KEY_CHECKS = 0; near the start
    if ($content -match "SET FOREIGN_KEY_CHECKS\s*=\s*0;") {
        Write-Host "  [SKIP] Already has SET FOREIGN_KEY_CHECKS = 0"
    } else {
        # Find the first CREATE TABLE or comment after initial comments
        # and insert SET FOREIGN_KEY_CHECKS = 0; before it
        $lines = $content -split "`n"
        $insert_at = -1
        
        for ($i = 0; $i -lt $lines.Count; $i++) {
            $line = $lines[$i]
            # Skip initial comments and blank lines
            if ($line -match "^--" -or $line -match "^\s*$") {
                continue
            }
            # Insert before first non-comment, non-blank line
            if ($line -match "^CREATE\s|^INSERT\s|^ALTER\s|^DROP\s") {
                $insert_at = $i
                break
            }
        }
        
        if ($insert_at -gt 0) {
            $lines[$insert_at] = "SET FOREIGN_KEY_CHECKS = 0;`n`n" + $lines[$insert_at]
            $content = $lines -join "`n"
            Set-Content -Path $file_path -Value $content
            Write-Host "  [ADDED] SET FOREIGN_KEY_CHECKS = 0 at line $insert_at"
        }
    }
    
    # Check if it ends with SET FOREIGN_KEY_CHECKS = 1;
    if ($content -match "SET FOREIGN_KEY_CHECKS\s*=\s*1;") {
        Write-Host "  [OK] Already has SET FOREIGN_KEY_CHECKS = 1"
    } else {
        # Append at the end (before any trailing whitespace/backticks)
        $content = $content -replace '```\s*$', ''
        $content = $content.TrimEnd()
        $content = $content + "`n`nSET FOREIGN_KEY_CHECKS = 1;"
        Set-Content -Path $file_path -Value $content
        Write-Host "  [ADDED] SET FOREIGN_KEY_CHECKS = 1 at end"
    }
}

Write-Host "Done!"
