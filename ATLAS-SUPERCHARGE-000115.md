# Atlas Supercharge RC-000115 — Payment Provider & Settlement Integrity

Adds HMAC webhook authenticity verification, independent-provider settlement quorum, provider-disagreement fail-close behavior, chain finality regression/reorg detection, invoice accounting reconciliation, and compliance-hold primitives. A single external provider cannot independently declare a payment settled.

Production activation remains gated on real provider/custody integration, secret management/HSM policy, live-chain qualification, accounting review, compliance configuration, and adversarial testing.
