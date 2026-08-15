# RC-000120 Problem Elimination

This release removes a qualification-runtime ambiguity and adds financial loss-prevention controls.

## Root cause found
The prior focused `npx vitest` command could not prove a test result because the packaged release intentionally does not contain `node_modules`; `npx` can wait on package resolution/network rather than execute a local runner. A timeout therefore was not evidence of a product-code deadlock.

## Controls added
- Partial payments remain open and cannot masquerade as paid-in-full.
- Overpayments enter review instead of silently changing invoice economics.
- Payment allocations cannot exceed the received amount or duplicate an invoice allocation.
- Failed refunds do not reduce available cash and require fresh authorization before retry.
- Conflicting payment providers fail closed; secondary-provider failover is explicit degraded mode.
- Revenue recognition requires settled cash AND satisfied performance obligations and is blocked by disputes.
- Release qualification distinguishes dependency availability from application-test results.

## Qualification truth rule
No automated suite is called PASS unless its runner and dependencies are locally available and the command exits successfully. Structural/offline qualification is reported separately.
