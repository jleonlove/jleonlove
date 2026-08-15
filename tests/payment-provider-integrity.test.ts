import { describe, expect, it } from "vitest";
import { reconcileProviderObservations, detectReorg } from "../lib/payment-provider-integrity";
const base={chain:"XLM",txHash:"abc",destination:"GDEST",amountAtomic:1000n,confirmations:1,finalized:true,observedAt:new Date().toISOString()};
describe("payment provider integrity",()=>{
 it("requires independent agreement",()=>{expect(reconcileProviderObservations([{...base,provider:"a"},{...base,provider:"b"}]).settled).toBe(true)});
 it("fails closed on disagreement",()=>{expect(reconcileProviderObservations([{...base,provider:"a"},{...base,provider:"b",amountAtomic:999n}]).settled).toBe(false)});
 it("detects reorg/finality regression",()=>{expect(detectReorg({...base,provider:"a"},{...base,provider:"a",finalized:false})).toBe(true)});
});
