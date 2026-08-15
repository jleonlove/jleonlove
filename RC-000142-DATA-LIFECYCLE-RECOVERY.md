# RC-000142 — Data Lifecycle & Recovery Hardening

Adds fail-closed qualification for encryption-key rotation, restore authorization, tenant-bound restore scope, backup checksum integrity, revoked-key rejection, retention periods, legal holds, duplicate deletion prevention, and content-free deletion evidence.

No restore is considered safe merely because backup bytes exist. Restore requires authorization, exact tenant/workspace binding, integrity verification, and an active encryption key. Deletion cannot bypass retention or legal hold.
