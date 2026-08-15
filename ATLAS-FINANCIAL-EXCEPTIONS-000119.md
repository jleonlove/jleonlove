# Atlas RC-000119 — Financial Exception & Loss-Prevention Hardening

Adds explicit payment state transitions and permanent regression coverage for ACH returns/reversals, card disputes/refunds, provider-reference reuse, cross-workspace exception injection, settlement netting, and separation of cash availability from revenue recognition and treasury authority.

Financial history is append-only: exception events must preserve lineage to the original event; they do not erase it.
