# RC-000136 Evidence Integrity Hardening

Production gates can no longer pass from a boolean/status flip alone. A passed gate must carry an evidence artifact, SHA-256 digest, execution timestamp, environment, authority, release binding, and optional expiry that remains valid. The production validator verifies the artifact and digest and rejects stale, missing, mismatched, or cross-release evidence. Deterministic repository qualification delegates the final production decision to this validator.
