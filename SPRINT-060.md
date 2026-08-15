# Sprint 060 — Organization Context

Atlas now includes a tenant-aware organization control plane.

## Shipped
- Persistent organizations, users, memberships, roles, and workspaces
- Cookie-backed active organization/workspace/user context
- Workspace switcher in the application shell
- Organization control-plane page
- Tenant- and workspace-scoped dashboard metrics
- Tenant-scoped document ingestion and listing
- Backward-compatible state migration for existing Sprint 059 data

## Security architecture
Every document and trust event now carries organization, workspace, actor, and role attribution. The local repository remains a pilot implementation; production should enforce the same boundaries in PostgreSQL row-level security and application authorization.
