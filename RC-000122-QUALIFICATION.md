# RC-000122 Qualification Hardening

RC-000122 strengthens the deterministic release gate. It pins Node/npm, requires exactly one dependency lockfile, refuses network/package-manager fallback during test execution, verifies non-zero test discovery, requires machine-readable Vitest evidence, records hashes tying evidence to package/lock inputs, and distinguishes QUALIFIED / FAILED / ENVIRONMENT_BLOCKED.

A dependency lockfile could not be truthfully generated in the current constrained build environment because dependency resolution did not complete. RC-000122 therefore remains fail-closed until a deterministic locked install and complete local suite succeed.
