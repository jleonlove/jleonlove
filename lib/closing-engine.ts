import type { AtlasDeal, ClosingGate, DealEvent } from './types';
import { canConsequentiallyExecute, diagnoseDeal } from './deal-engine';
import { computeCriticalPath, invalidateAssumptions } from './deal-hardening';

export function closingGate(deal:AtlasDeal):ClosingGate {
 const d=diagnoseDeal(deal), reasons:string[]=[];
 if(!canConsequentiallyExecute(deal)) reasons.push('Deal integrity or critical blockers prevent consequential execution.');
 if(d.staleEvidence.length) reasons.push('Stale or expired evidence requires refresh.');
 if(d.unknowns.length) reasons.push('Material claims remain claimed or unknown.');
 const authority=deal.participants.filter(p=>['seller','buyer','mandate','representative'].includes(p.role.toLowerCase())&&p.authority!=='INDEPENDENTLY_VERIFIED');
 if(authority.length) reasons.push('Material participant authority is not independently verified.');
 const cp=computeCriticalPath(deal); if(cp.unresolvedDependencies.length) reasons.push('Requirement graph contains unresolved dependencies.');
 return {pass:reasons.length===0,reasons,requiredApprovals:deal.requirements.filter(r=>r.ownerRole.toLowerCase().includes('approval')&&r.status!=='SATISFIED').map(r=>r.id)};
}

export function applyDealEvent(deal:AtlasDeal,event:DealEvent,seenKeys:Set<string>){
 if(event.dealId!==deal.id) throw new Error('EVENT_DEAL_MISMATCH');
 if(seenKeys.has(event.idempotencyKey)) return {deal,duplicate:true,changedKeys:[] as string[]};
 const changedKeys=Object.keys(event.payload); const patch:Partial<AtlasDeal>={};
 for(const k of ['quantity','unit','stage','currency','indicativeValue'] as const) if(k in event.payload) (patch as any)[k]=event.payload[k];
 const next={...deal,...patch,assumptions:invalidateAssumptions(deal,changedKeys),updatedAt:event.occurredAt}; seenKeys.add(event.idempotencyKey);
 return {deal:next,duplicate:false,changedKeys};
}

export function nextClosingActions(deal:AtlasDeal,limit=5){
 const cp=computeCriticalPath(deal); return cp.ordered.slice(0,limit).map((r,index)=>({priority:index+1,requirementId:r.id,title:r.title,severity:r.severity,ownerRole:r.ownerRole,status:r.status}));
}
