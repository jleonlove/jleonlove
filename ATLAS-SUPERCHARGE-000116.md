# Atlas RC-000116 — Custody & Webhook Boundary Hardening

Adds a fail-closed custody request boundary, deterministic request digests, tenant/workspace scoping, expiring nonces, distinct multi-approval requirements, explicit application-state private-key rejection, authenticated webhook freshness checks, and replay/idempotency enforcement.

This does not enable production custody or real-money transfers. External custody/HSM/MPC integration and live-chain qualification remain launch gates.
