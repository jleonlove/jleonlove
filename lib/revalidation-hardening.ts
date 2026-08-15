import type { AtlasDeal, DealEvidence, DealRequirement } from './types';
import { closingGate } from './closing-engine';

export type RevalidationTrigger = 'EVIDENCE_REVOKED'|'EVIDENCE_STALE'|'COUNTERPARTY_COMPROMISED'|'MATERIAL_CHANGE'|'DOCUMENT_SUPERSEDED';
export interface ApprovalSnapshot { id:string; requirementId:string; approvedAt:string; evidenceIds:string[]; dealUpdatedAt:string; status:'ACTIVE'|'STALE'|'REVOKED'; }
export interface RevalidationReport { pass:boolean; trigger:RevalidationTrigger; invalidatedRequirementIds:string[]; staleApprovalIds:string[]; quarantinedParticipantIds:string[]; reasons:string[]; }

function dependentRequirements(deal:AtlasDeal,evidenceIds:Set<string>){
 return deal.requirements.filter(r=>r.evidenceIds.some(id=>evidenceIds.has(id)));
}

export function cascadeEvidenceRevocation(deal:AtlasDeal, revokedEvidenceIds:string[]):AtlasDeal {
 const revoked=new Set(revokedEvidenceIds);
 const evidence=deal.evidence.map(e=>revoked.has(e.id)?{...e,status:'CONTRADICTED' as const}:e);
 const affected=new Set(dependentRequirements(deal,revoked).map(r=>r.id));
 const requirements=deal.requirements.map(r=>affected.has(r.id)&&r.status==='SATISFIED'?{...r,status:'BLOCKED' as const}:r);
 return {...deal,evidence,requirements,updatedAt:new Date().toISOString()};
}

export function quarantineCounterparty(deal:AtlasDeal, participantId:string):AtlasDeal {
 const participants=deal.participants.map(p=>p.id===participantId?{...p,authority:'CONTRADICTED' as const}:p);
 const requirements=deal.requirements.map(r=>/authority|counterparty|kyc|beneficiar|bank/i.test(r.title)&&r.status==='SATISFIED'?{...r,status:'BLOCKED' as const}:r);
 return {...deal,participants,requirements,updatedAt:new Date().toISOString()};
}

export function revalidateApprovals(deal:AtlasDeal, approvals:ApprovalSnapshot[], changedEvidenceIds:string[]=[]){
 const changed=new Set(changedEvidenceIds);
 return approvals.map(a=>{
  const requirement=deal.requirements.find(r=>r.id===a.requirementId);
  const evidenceChanged=a.evidenceIds.some(id=>changed.has(id));
  const evidenceInvalid=a.evidenceIds.some(id=>{const e=deal.evidence.find(x=>x.id===id);return !e||['CONTRADICTED','STALE','UNKNOWN','CLAIMED'].includes(e.status);});
  const dealChanged=new Date(deal.updatedAt).getTime()>new Date(a.dealUpdatedAt).getTime();
  const invalid=!requirement||requirement.status!=='SATISFIED'||evidenceChanged||evidenceInvalid||dealChanged;
  return invalid?{...a,status:'STALE' as const}:a;
 });
}

export function revalidationGate(deal:AtlasDeal, trigger:RevalidationTrigger, approvals:ApprovalSnapshot[], changedEvidenceIds:string[]=[]):RevalidationReport {
 const checked=revalidateApprovals(deal,approvals,changedEvidenceIds);
 const staleApprovalIds=checked.filter(a=>a.status!=='ACTIVE').map(a=>a.id);
 const invalidatedRequirementIds=deal.requirements.filter(r=>r.status==='BLOCKED').map(r=>r.id);
 const quarantinedParticipantIds=deal.participants.filter(p=>p.authority==='CONTRADICTED').map(p=>p.id);
 const gate=closingGate(deal), reasons=[...gate.reasons];
 if(staleApprovalIds.length) reasons.push('Prior approval is stale after material deal/evidence change and must be re-issued.');
 if(quarantinedParticipantIds.length) reasons.push('Compromised or contradicted counterparty is quarantined from consequential execution.');
 return {pass:gate.pass&&staleApprovalIds.length===0&&quarantinedParticipantIds.length===0,trigger,invalidatedRequirementIds,staleApprovalIds,quarantinedParticipantIds,reasons:[...new Set(reasons)]};
}

export function detectCrossDocumentContradictions(evidence:DealEvidence[]){
 const normalized=(s:string)=>s.toLowerCase().replace(/\s+/g,' ').trim();
 const byClaim=new Map<string,DealEvidence[]>();
 for(const e of evidence){const key=normalized(e.claim).replace(/\b(not|no|never)\b/g,'').trim();const arr=byClaim.get(key)??[];arr.push(e);byClaim.set(key,arr);}
 return [...byClaim.values()].filter(group=>group.length>1&&group.some(e=>e.status==='CONTRADICTED')).map(group=>({evidenceIds:group.map(e=>e.id),claim:group[0].claim,severity:'CRITICAL' as const}));
}
