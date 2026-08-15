# Atlas RC-000127 — Authority, Evidence, Isolation & Recovery Hardening

RC-000127 strengthens consequential deal closing against authority escalation, evidence poisoning, document/evidence substitution, cross-workspace contamination, and recovery after a mid-closing infrastructure failure.

## Controls
- Trusted execution snapshot is re-read at execution time and compared with caller-submitted state.
- Deal, organization, workspace, and state-version bindings are validated against the authoritative deal record.
- Current authority grants are independently revalidated; an execution envelope cannot escalate beyond those grants.
- Closing evidence must be verified, non-revoked, unique, and bound to the exact deal/workspace/state version.
- Prepared-close recovery is digest-bound to the exact envelope and evidence set; altered recovery inputs are rejected.
- Settlement/cancellation races remain routed to explicit exception state.

These controls are fail-safe by design: a safety gate may block an unsafe action, while the remediation path diagnoses and resolves the underlying condition rather than bypassing the gate.
