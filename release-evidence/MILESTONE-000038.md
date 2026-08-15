# Milestone 000038 — Governed Skill & Extension Fabric

Verified 2026-08-07:
- Extension lifecycle states: DRAFT, QUARANTINED, QUALIFIED, APPROVED, ACTIVE, REVOKED.
- Installation/quarantine does not grant execution authority.
- Only approved extensions can become active.
- Revocation immediately prevents new invocation.
- Capability diffing detects newly requested capabilities, secrets, network, filesystem, MCP, and hooks.
- Capability increases force requalification.
- Effective Agent Manifest digest changes when loaded authority changes.
- go test -race ./...: PASS
- go vet ./...: PASS
- atlasd build/boot: PASS

Next hardening:
publisher signatures, package/SBOM verification, MCP mediation, hook trigger governance,
sandbox qualification, persistent registry, and organization policy bindings.
