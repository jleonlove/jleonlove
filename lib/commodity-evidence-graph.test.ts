import {describe,it,expect} from "vitest";
import {analyzeCommodityEvidence,buildClosingReadinessGraph,resolutionPlan} from "./commodity-evidence-graph";
describe("commodity evidence graph",()=>{
 it("detects conflicting quantities",()=>{const x=analyzeCommodityEvidence([{claimId:"a",dealId:"d",type:"quantity",value:100,sourceDocumentId:"bl",status:"CORROBORATED"},{claimId:"b",dealId:"d",type:"quantity",value:90,sourceDocumentId:"invoice",status:"CORROBORATED"}]);expect(x.some(i=>i.code==="CLAIM_CONFLICT_QUANTITY")).toBe(true)});
 it("does not confuse extraction with truth",()=>{const x=analyzeCommodityEvidence([{claimId:"a",dealId:"d",type:"ownership",value:"Seller",sourceDocumentId:"pop",status:"PARSED"}]);expect(x.some(i=>i.code==="CLAIM_NOT_CORROBORATED")).toBe(true)});
 it("blocks only dependent readiness items",()=>{const issue={code:"CLAIM_CONFLICT_QUANTITY",severity:"HIGH" as const,claimIds:["a"],blocks:["CLOSING"],cure:"correct"};const g=buildClosingReadinessGraph([{id:"close",dependsOn:[],issueCodes:[issue.code]},{id:"schedule-inspection",dependsOn:[],issueCodes:[]}],[issue]);expect(g[0].state).toBe("BLOCKED");expect(g[1].state).toBe("READY")});
 it("prioritizes critical cures",()=>{const p=resolutionPlan([{code:"A",severity:"MEDIUM",claimIds:[],blocks:[],cure:"a"},{code:"B",severity:"CRITICAL",claimIds:[],blocks:[],cure:"b"}]);expect(p[0].issue).toBe("B")});
});
