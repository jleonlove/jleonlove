# Sprint 061 — Agent Release Registry and Tool Gateway

## Shipped

- Immutable agent release records with cryptographic digests
- Candidate, approved, and revoked release states
- Minimum evaluation-score release gate
- Owner/admin approval and revocation controls
- Signed capability manifests for tools and memory classifications
- Deny-by-default governed tool gateway
- Tool risk tiers and human-approval requirements
- Trust decisions and evidence for every attempted tool action
- Persistent tool-execution ledger
- Agent Release Registry UI and APIs

## Security model

An agent can execute a tool only when its exact release is approved, the tool appears in its declared manifest, and required human authorization is present. Every attempt creates a Trust Fabric decision and execution record.
