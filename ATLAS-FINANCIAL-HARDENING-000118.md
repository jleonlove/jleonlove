# Atlas Financial Hardening 000118

Adds deterministic double-entry validation, posted-journal immutability, reversal-entry construction, bank-event deduplication, settlement-only cash recognition, and exact invoice allocation checks. ACH returns/reversals must be represented by compensating/reversal events rather than destructive history edits.

Release packaging is deterministic: stage from the confirmed prior RC, add source/tests/release evidence, create ZIP, hash ZIP, and verify the ZIP can be listed before upload.
