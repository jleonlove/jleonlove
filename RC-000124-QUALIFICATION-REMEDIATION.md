# RC-000124 — Qualification Remediation

Fail-closed remains the safety invariant. This release adds an automated remediation path that attempts lockfile generation, deterministic npm ci, and the existing qualification gate. It emits qualification/remediation.json and uses exit 42 only when the execution environment prevents dependency resolution. It never converts an environment block into PASS.
