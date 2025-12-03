#!/bin/bash

# Renaming map based on dependency analysis
declare -A rename_map=(
    ["001_initial_schema.sql"]="001_initial_schema.sql"
    ["002_civil_schema.sql"]="002_civil_schema.sql"
    ["003_multi_tenant_users.sql"]="003_multi_tenant_users.sql"
    ["004_construction_schema.sql"]="004_construction_schema.sql"
    ["005_create_team_table.sql"]="005_team_table.sql"
    ["006_gamification_schema.sql"]="006_gamification_schema.sql"
    ["007_gamification_system.sql"]="007_gamification_system.sql"
    ["008_modular_monetization_schema.sql"]="008_modular_monetization_schema.sql"
    ["009_phase1_features.sql"]="009_phase1_features.sql"
    ["010_scheduled_tasks_schema.sql"]="010_scheduled_tasks_schema.sql"
    ["011_phase2_tasks_notifications.sql"]="011_phase2_tasks_notifications.sql"
    ["012_sample_data.sql"]="012_sample_data.sql"
    ["013_accounts_gl_schema.sql"]="013_accounts_gl_schema.sql"
    ["013_tenant_customization.sql"]="014_tenant_customization.sql"
    ["014_purchase_module_schema.sql"]="015_purchase_module_schema.sql"
    ["015_sales_module_schema.sql"]="016_sales_module_schema.sql"
    ["015_project_collection_accounts_rera.sql"]="017_project_collection_accounts_rera.sql"
    ["016_milestone_tracking_and_reporting.sql"]="018_milestone_tracking_and_reporting.sql"
    ["016_hr_compliance_labour_laws.sql"]="019_hr_compliance_labour_laws.sql"
    ["017_real_estate_property_management.sql"]="020_real_estate_property_management.sql"
    ["017_tax_compliance_income_tax_gst.sql"]="021_tax_compliance_income_tax_gst.sql"
    ["018_hr_payroll_schema.sql"]="022_hr_payroll_schema.sql"
    ["020_comprehensive_test_data.sql"]="023_comprehensive_test_data.sql"
    ["021_comprehensive_customization.sql"]="024_comprehensive_customization.sql"
    ["021a_roles_permissions.sql"]="025_roles_permissions.sql"
    ["022_external_partner_system.sql"]="026_external_partner_system.sql"
    ["023_partner_sources_and_credit_policies.sql"]="027_partner_sources_and_credit_policies.sql"
    ["024_sample_partner_logins.sql"]="028_sample_partner_logins.sql"
    ["030_phase3_analytics.sql"]="029_phase3_analytics.sql"
    ["031_phase3_workflows.sql"]="030_phase3_workflows.sql"
)

# Execute renames
for old_name in "${!rename_map[@]}"; do
    new_name="${rename_map[$old_name]}"
    if [ -f "$old_name" ]; then
        mv "$old_name" "$new_name"
        echo "Renamed: $old_name → $new_name"
    fi
done

echo "Migration renaming complete!"
