import type { AtlasDeal } from './types';
import { diagnoseDeal } from './deal-engine';
import { authorizeExecution, type ConsequentialAction } from './execution-guard';
import { computeCriticalPath } from './deal-hardening';
import { evidenceLineage } from './evidence-lineage';

export type TortureSeverity='LOW'|'MEDIUM'|'HIGH'|'CRITICAL';
export interface TortureFinding { code:string; severity:TortureSeverity; message:string; evidenceIds:string[]; }
export interface TortureReport { dealId:string; pass:boolean; findings:TortureFinding[]; executionLocked:boolean; }

const suspicious=(s:string)=>/guaranteed|no\s+risk|skip\s+(kyc|compliance|verification)|urgent\s+payment|change\s+beneficiary|send\s+funds/i.test(s);

export function runCommodityTortureSuite(deal:AtlasDeal, action:ConsequentialAction='COMMIT_TERMS'):TortureReport {
 const findings:TortureFinding[]=[]; const diagnosis=diagnoseDeal(deal); const lineage=evidenceLineage(deal); const path=computeCriticalPath(deal);
 for(const e of deal.evidence){
  if(suspicious(e.claim)) findings.push({code:'SOCIAL_ENGINEERING_OR_BYPASS_LANGUAGE',severity:'CRITICAL',message:'Evidence/claim contains language attempting to bypass normal diligence or induce urgent consequential action.',evidenceIds:[e.id]});
  if(e.sourceIds.length===0) findings.push({code:'SOURCELESS_EVIDENCE',severity:'HIGH',message:`Evidence ${e.id} has no source lineage.`,evidenceIds:[e.id]});
  if(e.status==='INDEPENDENTLY_VERIFIED' && e.independentSourceGroups<1) findings.push({code:'INVALID_VERIFICATION_STATE',severity:'CRITICAL',message:`Evidence ${e.id} is marked independently verified without an independent source group.`,evidenceIds:[e.id]});
 }
 if(lineage.contradictions.length) findings.push({code:'UNRESOLVED_CONTRADICTIONS',severity:'CRITICAL',message:'Contradictory evidence remains unresolved.',evidenceIds:lineage.contradictions});
 if(lineage.expiredClaims.length) findings.push({code:'EXPIRED_EVIDENCE',severity:'HIGH',message:'Expired evidence requires refresh.',evidenceIds:lineage.expiredClaims});
 if(path.unresolvedDependencies.length) findings.push({code:'BROKEN_REQUIREMENT_GRAPH',severity:'CRITICAL',message:`Missing requirement dependencies: ${path.unresolvedDependencies.join(', ')}`,evidenceIds:[]});
 if(deal.quantity !== undefined && (!Number.isFinite(deal.quantity)||deal.quantity<=0)) findings.push({code:'INVALID_QUANTITY',severity:'CRITICAL',message:'Quantity must be finite and greater than zero.',evidenceIds:[]});
 if(deal.indicativeValue !== undefined && (!Number.isFinite(deal.indicativeValue)||deal.indicativeValue<0)) findings.push({code:'INVALID_VALUE',severity:'CRITICAL',message:'Indicative value must be finite and non-negative.',evidenceIds:[]});
 if(diagnosis.criticalBlockers.length) findings.push({code:'CRITICAL_BLOCKERS',severity:'CRITICAL',message:'Critical transaction blockers remain open.',evidenceIds:diagnosis.criticalBlockers.flatMap(r=>r.evidenceIds)});
 const execution=authorizeExecution(deal,action,[]);
 const executionLocked=!execution.allow || findings.some(f=>f.severity==='CRITICAL');
 return {dealId:deal.id,pass:!findings.some(f=>f.severity==='CRITICAL'),findings,executionLocked};
}
