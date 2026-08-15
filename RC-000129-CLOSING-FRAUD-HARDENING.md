# RC-000129 — Closing Fraud & Settlement Direction Integrity

Adds verified instruction provenance and immutable execution binding for fiat and supported crypto settlement destinations. Material changes to counterparty identity, mandate, settlement rail, beneficiary/bank/wallet destination, memo/tag, provenance, or verification state invalidate prior authorization. Conflicting settlement evidence is routed to explicit reconciliation rather than auto-close.

Regression targets include compromised-email redirection, beneficiary substitution, wallet-address replacement, stale legitimate instruction replay, mandate/counterparty changes, wrong-deal settlement evidence, and provider conflicts.
