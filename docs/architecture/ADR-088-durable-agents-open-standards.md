# Atlas — Durable Agents & Open Standards Integration (ADR-088)

Status: PROPOSED — implementation and external qualification required.
Scope: Atlas control plane; no deployment or financial authority granted.

## 1. Durable transaction agent runtime

Introduce an event-driven durable execution boundary for long-lived trade workflows.

State machine:
CREATED -> AUTHORIZED -> RUNNING -> WAITING -> RESUMING -> RUNNING -> COMPLETED.
Any state may enter QUARANTINED; failures requiring human action enter REVIEW_REQUIRED.
No automatic transition from QUARANTINED to RUNNING.

Each persisted checkpoint MUST contain tenant/workspace, transaction ID, workflow version,
principal and agent identity, signed release digest, intent digest, immutable event cursor,
authorized capability set, policy version, expiry, pending external dependency, and audit correlation ID.
Encrypt at rest and scope access by tenant; never checkpoint raw secrets.

On resume: authenticate event origin; enforce deduplication and replay protection; verify
checkpoint integrity and version compatibility; re-evaluate principal, agent, release,
policy, compliance, and transaction authorization; verify expiration and revocation;
then obtain fresh short-lived credentials. Fail closed on missing or stale evidence.

External side effects use durable outbox/inbox, idempotency keys, and reconciliation
on ambiguous outcomes. Never assume exactly-once delivery. Economic actions must continue
through the existing Verified Execution & Settlement Gateway; agent resume never grants
new authority to transfer or custody customer funds.

## 2. Provider-neutral responses adapter

Define a versioned internal interface for input messages, tool requests/results, streaming
events, structured outputs, usage, safety outcomes, cancellation, and errors.
Implement provider-specific adapters behind the existing Cognitive Execution Router.
Do not presume standards provide identical semantics. Validate tool-call ordering,
schema conformance, interruption/retry behavior, refusal handling, and audit provenance.
Enforce provider terms and customer authorization before enabling routes.

## 3. Certified portable skills

Skills are immutable, signed, versioned packages with declared inputs/outputs, permissions,
data classification, network boundaries, resource budgets, dependencies, and owner.
Admission pipeline: provenance -> malware/dependency scanning -> sandbox tests ->
prompt-injection/red-team tests -> deterministic contract tests -> capability certification ->
human approval -> signed registry activation.

Initial candidates: commodity document checks, counterparty evidence checks,
shipping milestone reconciliation, and trade exception diagnosis.
A skill never inherits settlement privileges merely by installation.

## 4. Isolation and evaluation

Compare sandbox options against Atlas's existing Protected Recovery Plane.
Run OpenEnv-style repeatable evaluation scenarios where useful, without binding
production design to any single vendor.
Required tests:
- 10,000 waiting workflows with bounded memory/cost measurements
- crash/restart and duplicate/out-of-order event delivery
- revoked agent resumed after checkpoint
- cross-tenant checkpoint access attempts
- corrupted checkpoint and incompatible workflow version
- forged skill manifest and dependency substitution
- provider streaming/tool schema mismatches
- delayed settlement callback with ambiguous finality
- incident kill switch, quarantine, and audit reconstruction

## 5. Delivery gates

Phase A: architecture review and threat model.
Phase B: isolated prototype, synthetic events, and benchmark evidence.
Phase C: integration behind disabled-by-default feature flags.
Phase D: external security assessment, observability, load/chaos, disaster recovery,
regulatory data, live integrations, red-team, and end-to-end golden trade qualification.
Production activation is BLOCKED until relevant external gates pass.

## 6. Acceptance criteria

- No cross-tenant access or privilege escalation under test.
- All resume paths reauthorize and record immutable evidence.
- Duplicate events never duplicate economic side effects.
- Failed model/skill compatibility tests prevent activation.
- A kill switch revokes execution authority and preserves audit trails.
- Reproducible performance/cost benchmarks and rollback procedures exist.

References: prior Atlas innovation review (Oct 25, 2026); upstream claims require
independent technical verification before vendor adoption.
