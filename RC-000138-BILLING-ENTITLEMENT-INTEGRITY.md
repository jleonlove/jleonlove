# RC-000138 — Billing & Entitlement Integrity

Atlas merchant payments now bind verified, final payment evidence to an immutable invoice identity and activate only the exact purchased service tier. Underpayment, overpayment, invoice mismatch, expired invoices, unverified evidence, non-final payments, and replay against paid/refunded invoices fail closed or route to review. Customer-money custody/transmission remains prohibited.
