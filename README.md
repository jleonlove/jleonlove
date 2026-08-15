# Atlas Genesis MVP v0.2

A runnable enterprise-intelligence MVP demonstrating organization context, workspaces, Memory Fabric, Trust Fabric, governed retrieval, and persistent audit evidence.

## Run
```bash
npm install
npm run dev
```
Open `http://localhost:3000`.

## Pilot flow
1. Open Atlas and switch between Executive and Atlas Product workspaces.
2. Upload or paste knowledge in Memory.
3. Ask a governed question.
4. Review Trust evidence.
5. View tenant-scoped dashboard metrics.
6. Open Organization to review users, roles, and workspaces.

## Storage
The pilot stores state atomically in `data/atlas-state.json`. The repository boundary is intentionally isolated so PostgreSQL with row-level security can replace it.


## Sprint 061

Open `/releases` to review immutable agent releases, approve evaluated candidates, and test governed tool requests. Every request is evaluated against the signed release manifest and recorded in the Trust Fabric.

## Milestone 087 — Verified Execution & Settlement Gateway

The control plane now includes a provider-neutral settlement gateway under `control-plane/internal/settlementgateway`. It binds every economic action to principal, agent, trajectory, transaction, upstream economic-authority evidence, compliance evidence, policy evidence, verified input digests, simulation results, bounded slippage, explicit provider/network/asset/currency/counterparty policy, finality proof, and a tamper-evident receipt. Duplicate intent IDs are idempotent and never re-execute a completed settlement; post-execution finality failures enter a reconciliation state without releasing the consumed economic budget.

Live provider adapters remain an external-integration gate. Atlas does not treat installing a payment/blockchain skill as authority to move value.


## RC-000127 hardening
Authority escalation, evidence poisoning, document substitution, tenant/workspace isolation, and prepared-close recovery are hardened in `lib/deal-closing-coordinator.ts`.
