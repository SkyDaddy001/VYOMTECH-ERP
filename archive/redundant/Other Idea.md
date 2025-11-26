# Thoughts Schema Documentation

This folder contains the database schema for the **Thoughts Project**, logically grouped into multiple files for better organization and maintainability.

## Schema Files

1. **Shared Utilities**: `shared_utilities.sql`
   - Contains reusable utility functions and enables required PostgreSQL extensions.
   - Includes functions for ULID generation, cleaning old data, and performing maintenance tasks.

2. **Shared Triggers and Functions**: `shared_triggers_and_functions.sql`
   - Contains reusable triggers and functions for common operations.
   - Includes triggers for password hashing and updating `updated_at` columns.

3. **User Management**: `user_management.sql`
   - Defines tables and constraints for managing users and JWT refresh tokens.
   - Includes constraints for email format and password length.

4. **Lead Management**: `lead_management.sql`
   - Manages leads, their activities, statuses, and document collections.

5. **Marketing**: `marketing.sql`
   - Manages vendors, campaigns, sources, budgets, targets, achievements, and retainer fees.

6. **Campaign API**: `campaign_api.sql`
   - Provides APIs for vendors to create campaigns, log leads, and fetch campaign performance.

7. **Lead Pipeline**: `lead_pipeline.sql`
   - Manages lead pipelines and stages for workflows like Pre-Sale, Sales, and QC.

8. **Asterisk Integration**: `asterisk_schema.sql`
   - Supports telephony functionalities such as call logs, voicemail, and IVR.

## Usage Instructions

1. Run the SQL files in the order listed above.
2. Ensure required PostgreSQL extensions are enabled.
3. Use the provided API functions for programmatic interaction.

## Notes

- **Security**:
  - Sensitive data (e.g., passwords, tokens) is encrypted using `pgcrypto`.
  - Access to critical tables is restricted using role-based permissions.
- **Performance**:
  - Leverages PostgreSQL features like `jsonb_path_exists`, `vector`, and `pg_stat_statements`.
  - Uses partitioning and BRIN indexes for large tables.
- **Scalability**:
  - Designed to handle high volumes of data with efficient indexing and partitioning.
