# Fix remaining BIGINT user ID fields that should be INT

$migrations_dir = "d:\VYOMTECH-ERP\migrations"
$sql_files = Get-ChildItem -Path $migrations_dir -Filter "*.sql" -File

Write-Host "Fixing BIGINT user ID fields that should be INT..."

$user_id_fields = @(
    'created_by', 'updated_by', 'verifier_id', 'verified_by',
    'approver_id', 'resolver_id', 'resolved_by', 'approved_by',
    'performed_by', 'action_by', 'shared_with_user_id', 'shared_by',
    'review_by_lawyer', 'user_id', 'customer_user_id', 'support_user_id',
    'sender_user_id'
)

foreach ($file in $sql_files) {
    $content = Get-Content -Path $file.FullName -Raw
    $original_content = $content
    
    # For each user ID field, replace BIGINT with INT
    foreach ($field in $user_id_fields) {
        $pattern = "$field BIGINT"
        $replacement = "$field INT"
        $content = $content -replace [regex]::Escape($pattern), $replacement
    }
    
    if ($content -ne $original_content) {
        Set-Content -Path $file.FullName -Value $content
        Write-Host "[FIXED] $($file.Name)"
    }
}

Write-Host "Done!"
