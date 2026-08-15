import type { AtlasDeal, DealRequirement } from './types';
import { diagnoseDeal } from './deal-engine';
import { computeCriticalPath } from './deal-hardening';

export interface DealRescuePlan { dealId:string; status:'HEALTHY'|'STALLED'|'CRITICAL'; causes:string[]; recovery:DealRequirement[]; killReasons:string[]; }
export function rescueDeal(deal:AtlasDeal, now=new Date()):DealRescuePlan {
 const d=diagnoseDeal(deal,now), cp=computeCriticalPath(deal), causes:string[]=[], killReasons:string[]=[];
 if(d.criticalBlockers.length) causes.push(`${d.criticalBlockers.length} critical blocker(s) remain.`);
 if(d.contradictions.length) { causes.push('Material contradictions remain unresolved.'); killReasons.push('Resolve contradictory evidence before further consequential work.'); }
 if(d.staleEvidence.length) causes.push('Required evidence is stale or expired.');
 const overdue=deal.requirements.filter(r=>r.status!=='SATISFIED'&&r.status!=='WAIVED'&&r.dueAt&&new Date(r.dueAt)<now);
 if(overdue.length) causes.push(`${overdue.length} requirement(s) are overdue.`);
 const status=killReasons.length||d.integrity==='LOW'?'CRITICAL':causes.length?'STALLED':'HEALTHY';
 return {dealId:deal.id,status,causes,recovery:cp.ordered.filter(r=>r.status!=='SATISFIED'&&r.status!=='WAIVED').slice(0,5),killReasons};
}
