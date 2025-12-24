# PowerShell script to fix all migration files

$migrations_dir = "d:\VYOMTECH-ERP\migrations"

# Define replacements as an array of tuples
$replacements = @(

    # Table reference changes
    @('REFERENCES tenants(', 'REFERENCES `tenant`('),
    @('REFERENCES users(', 'REFERENCES `user`('),
    @('REFERENCES bookings(', 'REFERENCES booking('),
    @('REFERENCES leads(', 'REFERENCES sales_lead(')
)

# Get all SQL files
$sql_files = Get-ChildItem -Path $migrations_dir -Filter "*.sql" -File

Write-Host "Processing $($sql_files.Count) migration files..."

foreach ($file in $sql_files) {
    Write-Host "Processing: $($file.Name)"
    $content = Get-Content -Path $file.FullName -Raw
    $original_content = $content
    
    # Apply all replacements
    foreach ($replacement in $replacements) {
        $content = $content -replace [regex]::Escape($replacement[0]), $replacement[1]
    }
    
    # Only write if content changed
    if ($content -ne $original_content) {
        Set-Content -Path $file.FullName -Value $content
        Write-Host "  [UPDATED]"
    } else {
        Write-Host "  [OK]"
    }
}

Write-Host "Done! All migrations processed."
