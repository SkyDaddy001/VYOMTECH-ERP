# Fix all user_id, created_by, updated_by columns that still use INT format
# These need to be CHAR(36) to match user.id

$migrationDir = "d:\VYOMTECH-ERP\migrations"
$files = Get-ChildItem -Path $migrationDir -Filter "*.sql" -File

Write-Host "Fixing INT user references to CHAR(36)..." -ForegroundColor Green

$patterns = @(
    ('user_id\s+INT\s+NOT\s+NULL', 'user_id CHAR(36) NOT NULL'),
    ('user_id\s+INT\s+DEFAULT', 'user_id CHAR(36) DEFAULT'),
    ('user_id\s+INT\s+UNIQUE', 'user_id CHAR(36) UNIQUE'),
    ('user_id\s+INT,', 'user_id CHAR(36),'),
    ('created_by\s+BIGINT\s+NOT\s+NULL', 'created_by CHAR(36) NOT NULL'),
    ('created_by\s+BIGINT\s+DEFAULT', 'created_by CHAR(36) DEFAULT'),
    ('created_by\s+BIGINT,', 'created_by CHAR(36),'),
    ('updated_by\s+BIGINT\s+NOT\s+NULL', 'updated_by CHAR(36) NOT NULL'),
    ('updated_by\s+BIGINT\s+DEFAULT', 'updated_by CHAR(36) DEFAULT'),
    ('updated_by\s+BIGINT,', 'updated_by CHAR(36),'),
    ('verifier_id\s+INT\s+NOT\s+NULL', 'verifier_id CHAR(36) NOT NULL'),
    ('verifier_id\s+INT\s+DEFAULT', 'verifier_id CHAR(36) DEFAULT'),
    ('approver_id\s+BIGINT\s+NOT\s+NULL', 'approver_id CHAR(36) NOT NULL'),
    ('approver_id\s+BIGINT\s+DEFAULT', 'approver_id CHAR(36) DEFAULT'),
    ('owner_id\s+BIGINT\s+NOT\s+NULL', 'owner_id CHAR(36) NOT NULL'),
    ('owner_id\s+BIGINT\s+DEFAULT', 'owner_id CHAR(36) DEFAULT')
)

foreach ($file in $files) {
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content
    
    foreach ($pattern in $patterns) {
        $content = [regex]::Replace($content, $pattern[0], $pattern[1])
    }
    
    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -Force
        Write-Host "Fixed: $($file.Name)" -ForegroundColor Cyan
    }
}

Write-Host "`nAll INT user references converted to CHAR(36)!" -ForegroundColor Green
