# ATLAS-000095 — Live Commodity Data & Provenance Fabric

Status: implemented in RC-000095.

Purpose: provide Atlas with a fail-closed ingestion/resolution layer for live commodity observations before production provider adapters are connected.

Controls:
- mandatory source, timestamp, freshness TTL, confidence and provenance evidence
- decision-use licensing enforcement
- request-specific freshness/confidence constraints
- deterministic provider selection
- material cross-provider conflict detection
- stale/unlicensed/low-confidence data rejected for consequential use

Production provider credentials, commercial licenses, benchmark subscriptions and freight/tariff feeds remain external deployment dependencies.
