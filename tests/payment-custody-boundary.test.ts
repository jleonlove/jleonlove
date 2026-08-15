import {describe,it,expect} from "vitest";
import {qualifyCustodyRequest,assertNoPrivateKeyMaterial} from "../lib/custody-boundary";
const req={tenantId:"t1",workspaceId:"w1",action:"WITHDRAW" as const,asset:"XLM",network:"stellar-mainnet",destination:"GDEST",amountAtomic:100n,approvalIds:["a1","a2"],nonce:"n1",expiresAt:"2099-01-01T00:00:00Z"};
describe("custody boundary",()=>{
 it("requires two distinct approvals",()=>expect(qualifyCustodyRequest({...req,approvalIds:["a1","a1"]}).allowed).toBe(false));
 it("rejects expired requests",()=>expect(qualifyCustodyRequest({...req,expiresAt:"2020-01-01T00:00:00Z"}).allowed).toBe(false));
 it("forbids private key material in app state",()=>expect(assertNoPrivateKeyMaterial({privateKey:"secret"}).safe).toBe(false));
});
