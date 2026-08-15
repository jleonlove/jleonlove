# RC-000141 — Identity, Tenant, Recovery & Credential Closure

Adds fail-closed authorization binding identity, session, organization, workspace, role, resource and action. Revoked sessions are denied. Recovery-restricted sessions cannot perform privileged actions. Cross-tenant access and member privilege escalation are denied. Account recovery can revoke all subject sessions within the affected organization. Credential rotation rejects identifier reuse.
