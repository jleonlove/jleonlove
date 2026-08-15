export type EvidenceStatus = "RECEIVED"|"PARSED"|"CONSISTENT"|"ISSUER_VERIFIED"|"CORROBORATED"|"TRANSACTION_RELEVANT";
export type CommodityClaim = { claimId:string; dealId:string; type:string; value:string|number; unit?:string; sourceDocumentId:string; issuer?:string; status:EvidenceStatus; observedAt?:string };
export type EvidenceIssue = { code:string; severity:"MEDIUM"|"HIGH"|"CRITICAL"; claimIds:string[]; blocks:string[]; cure:string };

const norm=(v:unknown)=>String(v??"").trim().toLowerCase();
export function analyzeCommodityEvidence(claims:CommodityClaim[]):EvidenceIssue[]{
  const issues:EvidenceIssue[]=[];
  const byType=new Map<string,CommodityClaim[]>();
  for(const c of claims){ const a=byType.get(c.type)??[]; a.push(c); byType.set(c.type,a); }
  for(const [type, group] of byType){
    const values=new Set(group.map(c=>norm(c.value)));
    if(values.size>1) issues.push({code:`CLAIM_CONFLICT_${type.toUpperCase()}`,severity:"HIGH",claimIds:group.map(c=>c.claimId),blocks:["CLOSING_READINESS"],cure:`Obtain authoritative ${type} evidence and rebind affected documents.`});
  }
  const seen=new Map<string,CommodityClaim>();
  for(const c of claims){
    const key=`${c.sourceDocumentId}:${c.type}:${norm(c.value)}`;
    const prior=seen.get(key); if(prior && prior.dealId!==c.dealId) issues.push({code:"EVIDENCE_REUSED_ACROSS_DEALS",severity:"CRITICAL",claimIds:[prior.claimId,c.claimId],blocks:["TITLE","POP","CLOSING"],cure:"Verify original issuer/provenance and establish deal-specific evidence."}); else seen.set(key,c);
    if(c.status!=="CORROBORATED" && c.status!=="TRANSACTION_RELEVANT") issues.push({code:"CLAIM_NOT_CORROBORATED",severity:"MEDIUM",claimIds:[c.claimId],blocks:["DEPENDENT_EXECUTION"],cure:"Obtain independent corroboration from an authoritative source."});
  }
  return issues;
}

export type ReadinessItem={id:string; dependsOn:string[]; issueCodes:string[]; state:"READY"|"BLOCKED"|"PARALLEL"; cure?:string};
export function buildClosingReadinessGraph(items:Omit<ReadinessItem,"state">[], issues:EvidenceIssue[]):ReadinessItem[]{
  const active=new Set(issues.map(i=>i.code));
  return items.map(i=>{ const hit=i.issueCodes.filter(x=>active.has(x)); return {...i,state:hit.length?"BLOCKED":"READY",cure:hit.length?issues.find(x=>hit.includes(x.code))?.cure:undefined}; });
}

export function resolutionPlan(issues:EvidenceIssue[]){
  const rank={CRITICAL:0,HIGH:1,MEDIUM:2};
  return [...issues].sort((a,b)=>rank[a.severity]-rank[b.severity]).map((i,index)=>({order:index+1,issue:i.code,action:i.cure,blocks:i.blocks}));
}
