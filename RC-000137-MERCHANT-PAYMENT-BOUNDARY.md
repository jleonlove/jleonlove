# RC-000137 — Merchant Payment Boundary

Atlas receives funds only as payment for Atlas products/services. Atlas does not receive, custody, escrow, exchange, forward, withdraw, or transmit customer funds on behalf of customers.

Supported merchant payment assets remain USD, BTC, XRP, ETH, HBAR, ADA, and XLM. Crypto payment instructions must point to an Atlas-owned receiving destination associated with an Atlas service invoice. XRP destination-tag and XLM memo policy remains enforced.

Generic customer custody requests are structurally denied. Private keys, seed phrases, mnemonics, and signing secrets remain forbidden from application/agent state.
