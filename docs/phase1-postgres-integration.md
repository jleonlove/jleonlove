# PostgreSQL evidence isolation: integration checklist

This SQL is a proposal for an isolated test database, not an applied migration.

- The application reader must not own the evidence table and must not have superuser or BYPASSRLS privileges.
- Grant only schema USAGE and table SELECT to the dedicated reader role. Do not grant write access.
- Authenticate the principal and verify organization/workspace membership in a trusted backend before opening a transaction.
- Within that transaction, set both scope variables with transaction-local settings, and execute all reads through a restricted data-access interface.
- Never accept tenant/workspace identifiers directly from untrusted request parameters as authoritative. Custom PostgreSQL session settings can be changed by anyone with unrestricted SQL access; the backend must not expose arbitrary SQL execution.
- Use a separate privileged ingestion identity for evidence writes, with audited authorization.
- Test cross-tenant access, cross-workspace access, missing scope, reused pooled connections, unauthorized writes, and privileged-role bypass.
- Do not migrate legacy records with missing ownership by silently assigning them to JLST.
- A SHA-256 content digest verifies equality against a trusted digest, not independent authenticity of the source.
- Static SQL assertions are not a substitute for live PostgreSQL integration tests.

No production deployment, migration, or test execution is claimed.
