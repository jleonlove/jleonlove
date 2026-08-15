import {describe,it,expect} from "vitest";import {classifyPayment,validateAllocations,refundAvailability,providerFailover,revenueRecognition} from "../lib/financial-loss-prevention";
describe("financial loss prevention",()=>{
it("holds partial and overpayments",()=>{expect(classifyPayment(100n,90n).status).toBe("PARTIAL");expect(classifyPayment(100n,110n).status).toBe("OVERPAYMENT_REVIEW")});
it("prevents over-allocation",()=>expect(validateAllocations(100n,[{invoiceId:"i",paymentId:"p",amountMinor:101n}]).ok).toBe(false));
it("failed refund does not release cash",()=>expect(refundAvailability({refundId:"r",paymentId:"p",amountMinor:10n,state:"FAILED"}).cashReleased).toBe(false));
it("provider disagreement fails closed",()=>expect(providerFailover(true,true,false).ok).toBe(false));
it("revenue requires settlement and performance",()=>expect(revenueRecognition(true,false,false).recognizable).toBe(false));
});
