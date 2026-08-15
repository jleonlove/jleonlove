# ATLAS-000091 — Counterparty, Commodity & Transaction Verification Fabric

Status: IMPLEMENTED (local qualification)

Adds a fail-closed verification fabric for KYC, KYB, UBO, mandate authority, sanctions, PEP/adverse-media screening, proof-of-commodity, facility verification, document integrity and bank/account ownership evidence. Required checks must carry evidence and source attribution; expired evidence is treated as stale and cannot authorize execution.

Execution invariant: consequential transaction execution requires PASS + score 100 for all required verification checks.
