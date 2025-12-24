# More precise fixes for specific patterns that shouldn't have been changed

$migrations_dir = "d:\VYOMTECH-ERP\migrations"
$sql_files = Get-ChildItem -Path $migrations_dir -Filter "*.sql" -File

Write-Host "Fixing overly broad replacements..."

foreach ($file in $sql_files) {
    $content = Get-Content -Path $file.FullName -Raw
    $original_content = $content
    
    # Revert internal reference IDs that shouldn't be VARCHAR(36)
    # These are FOREIGN KEYs to local tables, not to entities
    $content = $content -replace 'category_id VARCHAR\(36\)', 'category_id BIGINT'
    $content = $content -replace 'document_type_id VARCHAR\(36\)', 'document_type_id BIGINT'
    $content = $content -replace 'collection_id VARCHAR\(36\)', 'collection_id BIGINT'
    $content = $content -replace 'document_id VARCHAR\(36\)', 'document_id BIGINT'
    $content = $content -replace 'template_id VARCHAR\(36\)', 'template_id BIGINT'
    $content = $content -replace 'site_id VARCHAR\(36\)', 'site_id BIGINT'
    $content = $content -replace 'project_id VARCHAR\(36\)', 'project_id BIGINT'
    $content = $content -replace 'block_id VARCHAR\(36\)', 'block_id BIGINT'
    $content = $content -replace 'visit_id VARCHAR\(36\)', 'visit_id BIGINT'
    $content = $content -replace 'clearance_id VARCHAR\(36\)', 'clearance_id BIGINT'
    
    # Fix user reference IDs that should be INT, not VARCHAR
    $content = $content -replace 'verifier_id VARCHAR\(36\)', 'verifier_id INT'
    $content = $content -replace 'verified_by VARCHAR\(36\)', 'verified_by INT'
    $content = $content -replace 'created_by VARCHAR\(36\)', 'created_by INT'
    $content = $content -replace 'updated_by VARCHAR\(36\)', 'updated_by INT'
    $content = $content -replace 'approver_id VARCHAR\(36\)', 'approver_id INT'
    $content = $content -replace 'resolver_id VARCHAR\(36\)', 'resolver_id INT'
    $content = $content -replace 'resolved_by VARCHAR\(36\)', 'resolved_by INT'
    $content = $content -replace 'approved_by VARCHAR\(36\)', 'approved_by INT'
    $content = $content -replace 'performed_by VARCHAR\(36\)', 'performed_by INT'
    $content = $content -replace 'action_by VARCHAR\(36\)', 'action_by INT'
    $content = $content -replace 'shared_with_user_id VARCHAR\(36\)', 'shared_with_user_id INT'
    $content = $content -replace 'shared_by VARCHAR\(36\)', 'shared_by INT'
    $content = $content -replace 'review_by_lawyer VARCHAR\(36\)', 'review_by_lawyer INT'
    $content = $content -replace 'user_id VARCHAR\(36\)', 'user_id INT'
    $content = $content -replace 'customer_user_id VARCHAR\(36\)', 'customer_user_id INT'
    $content = $content -replace 'support_user_id VARCHAR\(36\)', 'support_user_id INT'
    $content = $content -replace 'sender_user_id VARCHAR\(36\)', 'sender_user_id INT'
    
    # For entity IDs referenced by text names, keep as VARCHAR(36)
    # But make sure we're not overwriting the explicit INT/BIGINT assignments above
    
    if ($content -ne $original_content) {
        Set-Content -Path $file.FullName -Value $content
        Write-Host "[FIXED] $($file.Name)"
    }
}

Write-Host "Done!"
