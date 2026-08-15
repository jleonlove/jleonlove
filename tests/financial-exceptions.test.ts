import {describe,it,expect} from "vitest";import {validateTransition,financialAvailability,validateExceptionEvent,reconcileSettlement,detectProviderReferenceReuse} from "../lib/financial-exceptions";import type {FinancialEvent} from "../lib/financial-exceptions";
const base:FinancialEvent={eventId:"e1",paymentId:"p1",invoiceId:"i1",organizationId:"o1",workspaceId:"w1",rail:"ACH",state:"SETTLED",amountMinor:10000n,providerReference:"bank-1",occurredAt:"2026-08-15T00:00:00Z"};
describe("financial exception hardening",()=>{
it("permits a settled ACH return but requires lineage",()=>{const r={...base,eventId:"e2",state:"RETURNED" as const,originalEventId:"e1",providerReference:"bank-2"};expect(validateExceptionEvent(r,base).ok).toBe(true)});
it("rejects destructive impossible transitions",()=>expect(validateTransition("RETURNED","SETTLED").ok).toBe(false));
it("separates settlement from revenue and treasury authority",()=>expect(financialAvailability("SETTLED")).toEqual({cashAvailable:true,revenueEarned:false,treasuryAuthorized:false}));
it("nets returns out of reconciliation",()=>{const ret={...base,eventId:"e2",state:"RETURNED" as const,originalEventId:"e1",providerReference:"bank-2"};expect(reconcileSettlement(0n,[base,ret]).ok).toBe(true)});
it("detects provider reference reused across payments",()=>{const other={...base,eventId:"e3",paymentId:"p2"};expect(detectProviderReferenceReuse([base,other]).ok).toBe(false)});
it("rejects cross-workspace exception injection",()=>{const r={...base,eventId:"e2",workspaceId:"evil",state:"RETURNED" as const,originalEventId:"e1",providerReference:"bank-2"};expect(validateExceptionEvent(r,base).errors).toContain("SCOPE_OR_PAYMENT_MISMATCH")});
});
