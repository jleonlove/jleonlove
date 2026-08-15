# Milestone 000039 — Extension Supply-Chain Trust

Verified 2026-08-07:
- Ed25519 publisher package signatures.
- Publisher/key identity binding.
- Publisher-key revocation.
- Immutable SHA-256 artifact verification.
- Signed SBOM digest verification.
- Source revision + builder provenance binding.
- Package substitution, SBOM substitution, signed-statement tampering, unknown publisher,
  revoked publisher, and provenance mismatch all fail closed.
- Deterministic dependency fingerprinting.
- go test -race ./...: PASS
- go vet ./...: PASS
- atlasd build/boot: PASS

Boundary:
Production deployment still requires external publisher identity proofing, secure key custody/rotation,
transparency logging, vulnerability intelligence, and stronger standardized provenance attestations.
