import type { AtlasDeal, DealRequirement } from './types';
import { diagnoseDeal } from './deal-engine';
import { computeCriticalPath, evaluatePhysicalReality, type PhysicalClaim } from './deal-hardening';
import { rescueDeal } from './deal-rescue';

export type ClosingDisposition='ADVANCE'|'HOLD'|'RESCUE'|'STOP';
export interface ClosingStressReport { dealId:string; disposition:ClosingDisposition; score:number; releaseBlocked:boolean; reasons:string[]; nextActions:DealRequirement[]; }
export function stressClosing(deal:AtlasDeal, physicalClaims:PhysicalClaim[]=[], now=new Date()):ClosingStressReport{
 const diagnosis=diagnoseDeal(deal,now), cp=computeCriticalPath(deal), rescue=rescueDeal(deal,now), reality=evaluatePhysicalReality(physicalClaims,deal.evidence), reasons:string[]=[];
 const criticalReality=reality.filter(x=>x.severity==='CRITICAL');
 if(diagnosis.contradictions.length) reasons.push('Unresolved contradictory evidence.');
 if(diagnosis.criticalBlockers.length) reasons.push('Critical deal requirements remain blocked.');
 if(cp.unresolvedDependencies.length) reasons.push('Requirement graph contains unresolved dependencies.');
 if(criticalReality.length) reasons.push('Physical-reality claims materially conflict.');
 if(rescue.killReasons.length) reasons.push(...rescue.killReasons);
 const verified=deal.evidence.filter(e=>e.status==='INDEPENDENTLY_VERIFIED').length;
 const total=Math.max(1,deal.evidence.length);
 let score=Math.round(100*verified/total)-diagnosis.criticalBlockers.length*15-diagnosis.contradictions.length*20-criticalReality.length*25-cp.unresolvedDependencies.length*10;
 score=Math.max(0,Math.min(100,score));
 const releaseBlocked=reasons.length>0||diagnosis.integrity==='LOW';
 const disposition:ClosingDisposition=releaseBlocked?(rescue.killReasons.length||criticalReality.length?'STOP':rescue.status==='STALLED'?'RESCUE':'HOLD'):'ADVANCE';
 return {dealId:deal.id,disposition,score,releaseBlocked,reasons,nextActions:cp.ordered.filter(r=>!['SATISFIED','WAIVED'].includes(r.status)).slice(0,5)};
}
