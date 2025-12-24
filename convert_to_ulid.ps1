# Convert all migrations from BIGINT AUTO_INCREMENT to CHAR(26) for ULID
# ULID format: 26 characters (timestamp + randomness)

$migrationDir = "d:\VYOMTECH-ERP\migrations"
$files = Get-ChildItem -Path $migrationDir -Filter "*.sql" -File

Write-Host "Converting all migrations to use ULID (CHAR(26)) for IDs..." -ForegroundColor Green

foreach ($file in $files) {
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content
    
    # Replace all BIGINT AUTO_INCREMENT PRIMARY KEY with CHAR(26) PRIMARY KEY
    $content = [regex]::Replace($content, 'id\s+BIGINT\s+AUTO_INCREMENT\s+PRIMARY\s+KEY', 'id CHAR(26) PRIMARY KEY')
    $content = [regex]::Replace($content, 'id\s+BIGINT\s+PRIMARY\s+KEY\s+AUTO_INCREMENT', 'id CHAR(26) PRIMARY KEY')
    $content = [regex]::Replace($content, 'id\s+BIGINT\s+AUTO_INCREMENT', 'id CHAR(26)')
    $content = [regex]::Replace($content, 'id\s+INT\s+AUTO_INCREMENT\s+PRIMARY\s+KEY', 'id CHAR(26) PRIMARY KEY')
    $content = [regex]::Replace($content, 'id\s+INT\s+AUTO_INCREMENT', 'id CHAR(26)')
    
    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -Force
        Write-Host "OK: $($file.Name)" -ForegroundColor Cyan
    }
}

Write-Host "Conversion complete! All BIGINT AUTO_INCREMENT IDs converted to CHAR(26)." -ForegroundColor Green
