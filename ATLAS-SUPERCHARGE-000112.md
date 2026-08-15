# Atlas RC-000112 — Revalidation & Compromise Containment

Adds evidence-revocation cascading, stale-approval invalidation, compromised-counterparty quarantine, cross-document contradiction propagation, and mandatory closing-gate revalidation after material changes.

## Fail-closed invariants
- Revoked/contradicted evidence can no longer support a previously satisfied requirement.
- Prior approvals do not survive material deal or evidence changes without revalidation.
- Compromised counterparties are quarantined from consequential execution.
- Closing readiness is recomputed after trust-state changes.
