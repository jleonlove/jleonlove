# Atlas Revenue Engine — RC-000117
USD is the primary commercial currency. Supported USD rails are ACH, wire, enterprise invoice/payment instructions, and provider-backed card processing. Digital-asset rails remain BTC, XRP, ETH, HBAR, ADA, and XLM.

Core invariants: initiated is not settled; returned/reversed funds are not available cash; every payment is invoice + organization + workspace scoped; reconciliation uses integer minor/atomic units; receipt authority is separate from refund/payout/withdraw/convert authority; live providers and real-money movement remain gated until credentials, provider integration, compliance, and production qualification are complete.
