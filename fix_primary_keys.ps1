# Fix the issue where id BIGINT was changed to id VARCHAR(36)
# Also need to be more careful with entity ID replacements

$migrations_dir = "d:\VYOMTECH-ERP\migrations"

# Get all SQL files
$sql_files = Get-ChildItem -Path $migrations_dir -Filter "*.sql" -File

Write-Host "Reviewing and fixing primary key issues..."

foreach ($file in $sql_files) {
    $content = Get-Content -Path $file.FullName -Raw
    $original_content = $content
    
    # Fix cases where we changed primary key 'id' to VARCHAR when it should stay as BIGINT or INT
    # Pattern: id VARCHAR(36) AUTO_INCREMENT PRIMARY KEY
    # Should be: id BIGINT AUTO_INCREMENT PRIMARY KEY
    $content = $content -replace 'id VARCHAR\(36\) AUTO_INCREMENT PRIMARY KEY', 'id BIGINT AUTO_INCREMENT PRIMARY KEY'
    $content = $content -replace 'id VARCHAR\(36\) PRIMARY KEY AUTO_INCREMENT', 'id BIGINT AUTO_INCREMENT PRIMARY KEY'
    
    # Also fix VARCHAR(36) AUTO_INCREMENT which is also wrong for primary keys
    $content = $content -replace 'id VARCHAR\(36\)\s+AUTO_INCREMENT', 'id BIGINT AUTO_INCREMENT'
    
    if ($content -ne $original_content) {
        Set-Content -Path $file.FullName -Value $content
        Write-Host "[FIXED] $($file.Name)"
    }
}

Write-Host "Done!"
