# Atlas Supercharge RC-000102 — Closing Engine

Adds fail-closed closing gates, stale/unknown evidence blocking, material authority blocking, unresolved dependency blocking, idempotent deal events, automatic assumption invalidation on material changes, and prioritized closing actions.

This is intentionally conservative: Atlas must not treat a deal as execution-ready while critical truth, authority, evidence freshness, or dependency conditions remain unresolved.
