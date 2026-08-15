import type { AtlasDeal, DealEvidence } from './types';

export interface EvidenceLineageNode { evidenceId:string; claim:string; sourceIds:string[]; independentSourceGroups:number; status:string; observedAt:string; expiresAt?:string; }
export interface EvidenceLineageReport { dealId:string; nodes:EvidenceLineageNode[]; weakClaims:string[]; expiredClaims:string[]; contradictions:string[]; }
export function evidenceLineage(deal:AtlasDeal, now=new Date()):EvidenceLineageReport {
 const nodes=deal.evidence.map((e:DealEvidence)=>({evidenceId:e.id,claim:e.claim,sourceIds:e.sourceIds,independentSourceGroups:e.independentSourceGroups,status:e.status,observedAt:e.observedAt,expiresAt:e.expiresAt}));
 return {dealId:deal.id,nodes,weakClaims:deal.evidence.filter(e=>e.independentSourceGroups<2&&e.status!=='INDEPENDENTLY_VERIFIED').map(e=>e.id),expiredClaims:deal.evidence.filter(e=>!!e.expiresAt&&new Date(e.expiresAt).getTime()<=now.getTime()).map(e=>e.id),contradictions:deal.evidence.filter(e=>e.status==='CONTRADICTED').map(e=>e.id)};
}
