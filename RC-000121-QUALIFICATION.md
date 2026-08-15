# RC-000121 — Deterministic Qualification Gate

Release advancement is fail-closed. Only `QUALIFIED` is a passing production qualification state.

- `QUALIFIED`: locked dependencies are present and the mandatory local Vitest suite exits 0.
- `FAILED`: a mandatory executable gate ran and failed or timed out.
- `ENVIRONMENT_BLOCKED`: qualification cannot execute because its locked runtime/dependencies are unavailable.

The gate never invokes `npx`, never downloads packages during qualification, uses the repository-local Vitest binary only, applies a hard suite timeout, emits machine-readable evidence, and returns distinct process exit codes.

Offline structural qualification remains supplemental and cannot promote an environment-blocked release.
