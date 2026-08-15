import {describe,it,expect} from "vitest";
import {validateJournal,reverseEntry,dedupeBankEvents,canRecognizeCashRevenue,paymentAllocation} from "../lib/financial-ledger";
describe("financial ledger hardening",()=>{
 it("rejects unbalanced journals",()=>{const r=validateJournal({id:"j1",organizationId:"o",workspaceId:"w",currency:"USD",effectiveAt:"2026-08-15",sourceEventId:"e1",posted:false,lines:[{accountId:"cash",accountType:"ASSET",debitMinor:100n,creditMinor:0n},{accountId:"rev",accountType:"REVENUE",debitMinor:0n,creditMinor:99n}]});expect(r.ok).toBe(false);expect(r.errors).toContain("UNBALANCED_JOURNAL")});
 it("reverses instead of editing posted entries",()=>{const e={id:"j1",organizationId:"o",workspaceId:"w",currency:"USD" as const,effectiveAt:"2026-08-15",sourceEventId:"e1",posted:true,lines:[{accountId:"cash",accountType:"ASSET" as const,debitMinor:100n,creditMinor:0n},{accountId:"rev",accountType:"REVENUE" as const,debitMinor:0n,creditMinor:100n}]};const r=reverseEntry(e,"j2","return1");expect(r.reversalOf).toBe("j1");expect(r.lines[0].creditMinor).toBe(100n)});
 it("detects duplicate bank events",()=>{const e={eventId:"x",paymentId:"p",invoiceId:"i",organizationId:"o",workspaceId:"w",rail:"ACH" as const,kind:"SETTLED" as const,amountMinor:100n,providerReference:"r"};expect(dedupeBankEvents([e,e]).ok).toBe(false)});
 it("does not recognize initiated cash as settled revenue",()=>{expect(canRecognizeCashRevenue({eventId:"x",paymentId:"p",invoiceId:"i",organizationId:"o",workspaceId:"w",rail:"ACH",kind:"INITIATED",amountMinor:100n,providerReference:"r"})).toBe(false)});
 it("requires exact invoice allocation",()=>expect(paymentAllocation(100n,[{invoiceId:"i",amountMinor:99n}]).ok).toBe(false));
});
