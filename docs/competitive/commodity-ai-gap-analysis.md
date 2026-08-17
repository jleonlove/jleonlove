# Atlas Commodity AI Competitive Gap Analysis

Status: implementation input, not a production-readiness claim.

## External benchmark signals

### CommodityAI
Observed strengths worth matching or exceeding:
- omnichannel document ingestion: email, WhatsApp, Teams, SharePoint, voice, PDF, Excel
- extract -> validate -> reconcile -> enrich -> generate -> sync workflows
- three-way reconciliation across trade documents and payments
- live vessel/port/market enrichment
- CTRM/ETRM, ERP/accounting, warehouse and channel connectors
- confidence thresholds with human exception queues
- full agent/human audit trails
- commodity-specific workflows such as demurrage/laytime and assay reconciliation
- enterprise controls: SSO/SAML, scoped API/OAuth, encryption, data residency, SOC 2 evidence

### Atlas Verified
Observed strengths worth matching or exceeding:
- large purpose-built verification tool estate
- authoritative-source supplier verification
- sanctions, import-alert, certification and entity checks
- cross-document consistency and duplicate/hash detection
- route and volume plausibility analysis
- persistent supplier verification histories
- autonomous monitoring agents
- chain-of-custody and shipment tracking
- defensible reports with source attribution
- reusable supplier verification/trust artifacts

## Atlas implementation priorities

P0 — Verification Mesh
Every consequential supplier, shipment, document and counterparty assertion must support independent checks against authoritative sources, freshness, provenance, contradiction detection and evidence retention.

P0 — Trade Reconciliation Engine
Implement deterministic reconciliation contracts for invoice/BOL/payment, quantity/weight, assay/quality, contract/shipment and settlement records. Variances must route to governed exception handling rather than silent correction.

P0 — Document Authenticity Graph
Add duplicate/reuse detection, document hashes, cross-document entity/date/quantity consistency, certificate scope/timeline validation and provenance linkage.

P0 — Shipment Plausibility & Chain of Custody
Bind shipment events to vessel/container/port/route evidence; detect impossible routes, timing anomalies, unexpected custody transitions and material divergence from contract terms.

P0 — Connector Control Plane
Treat CTRM/ETRM, ERP, accounting, warehouse, communications and external data connectors as certified capabilities with scoped OAuth/API credentials, schema contracts, idempotency, replay protection, version fingerprints and revocation.

P0 — Human Exception Operations
Every confidence/risk threshold must have explicit auto-execute, review, escalate, quarantine or deny behavior. Human overrides require identity, reason, exact-effect binding and audit evidence.

P1 — Omnichannel Trade Intake
Normalize approved email, messaging, voice, spreadsheet and document channels into a single provenance-preserving intake envelope. No channel may bypass document or authority controls.

P1 — Persistent Counterparty Intelligence
Maintain evidence-backed supplier/counterparty histories with verification freshness, risk changes, certifications, shipment performance and contradictions. Never convert a historical trust score into permanent authority.

P1 — Defensible Intelligence Reports
Generate shareable reports where every material conclusion resolves to source evidence, methodology, timestamps, model/release identity and confidence/disagreement state.

P1 — Commodity Workflow Packs
Build certified workflow packs for agriculture, energy, metals/minerals and chemicals, then specialize for assay reconciliation, demurrage/laytime, letters of credit, certificates, quality/quantity and settlement.

## Atlas differentiation to preserve
Atlas should not merely copy competitors. Its target advantage is governed end-to-end execution:
principal -> intent -> evidence -> verification -> compliance -> authority -> workflow -> human/agent action -> economic effect -> settlement evidence -> recovery/reconstruction.

No competitor feature is admitted merely because it is commercially useful. It must pass Atlas identity, information-flow, supply-chain, regulatory, authority, provenance, recovery and certification controls first.
