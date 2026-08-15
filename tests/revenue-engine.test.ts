import {describe,it,expect} from "vitest";
import {validateUsdPayment,settlementDisposition,reconcileRevenue,requiresMoneyMovementApproval} from "../lib/revenue-engine";
const base={id:"p1",invoiceId:"i1",organizationId:"o1",workspaceId:"w1",currency:"USD" as const,rail:"ACH",expectedMinorOrAtomic:"10000",state:"INITIATED" as const};
describe("USD-first revenue engine",()=>{
 it("supports scoped USD ACH",()=>expect(validateUsdPayment(base).ok).toBe(true));
 it("does not confuse initiated ACH with settled cash",()=>expect(settlementDisposition(base).cashAvailable).toBe(false));
 it("treats returned funds as unavailable",()=>expect(settlementDisposition({...base,state:"RETURNED"})).toEqual({cashAvailable:false,final:false}));
 it("reconciles exact USD cents",()=>expect(reconcileRevenue({...base,receivedMinorOrAtomic:"10000"}).status).toBe("MATCH"));
 it("separates receipt from treasury authority",()=>{expect(requiresMoneyMovementApproval("RECEIVE")).toBe(false);expect(requiresMoneyMovementApproval("REFUND")).toBe(true);});
});
