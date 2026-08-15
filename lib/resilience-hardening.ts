import type { WorkflowStep } from './workflow-recovery';

export interface WorkflowEvent { id:string; workflowId:string; sequence:number; idempotencyKey:string; kind:string; checksum?:string; payload?:unknown; }
export interface ResilienceFinding { severity:'MEDIUM'|'HIGH'|'CRITICAL'; code:string; detail:string; }
export interface ResilienceReport { safe:boolean; findings:ResilienceFinding[]; nextSequence:number; }

export function validateWorkflowEvents(events:WorkflowEvent[]):ResilienceReport {
  const findings:ResilienceFinding[]=[];
  const ids=new Set<string>(), keys=new Set<string>();
  const sorted=[...events].sort((a,b)=>a.sequence-b.sequence);
  let expected=sorted.length ? sorted[0].sequence : 0;
  for(const e of sorted){
    if(ids.has(e.id)) findings.push({severity:'CRITICAL',code:'DUPLICATE_EVENT_ID',detail:e.id});
    if(keys.has(e.idempotencyKey)) findings.push({severity:'CRITICAL',code:'DUPLICATE_IDEMPOTENCY_KEY',detail:e.idempotencyKey});
    if(e.sequence!==expected) findings.push({severity:'CRITICAL',code:'EVENT_SEQUENCE_GAP',detail:`expected ${expected}, received ${e.sequence}`});
    if(!e.workflowId || !e.kind || !e.idempotencyKey) findings.push({severity:'CRITICAL',code:'CORRUPT_EVENT',detail:e.id});
    ids.add(e.id); keys.add(e.idempotencyKey); expected=e.sequence+1;
  }
  return {safe:!findings.some(f=>f.severity==='CRITICAL'),findings,nextSequence:expected};
}

export function validateConcurrentSteps(steps:WorkflowStep[]):ResilienceReport {
  const findings:ResilienceFinding[]=[];
  const running=new Map<string,string>();
  for(const s of steps){
    if(s.status==='RUNNING'){
      const prior=running.get(s.idempotencyKey);
      if(prior && prior!==s.id) findings.push({severity:'CRITICAL',code:'CONCURRENT_SIDE_EFFECT_RACE',detail:`${prior} and ${s.id} share ${s.idempotencyKey}`});
      running.set(s.idempotencyKey,s.id);
    }
    if(s.attempts> s.maxAttempts) findings.push({severity:'CRITICAL',code:'RETRY_BUDGET_BREACH',detail:s.id});
  }
  return {safe:findings.length===0,findings,nextSequence:0};
}

export function requireCompensationAfterPartialFailure(steps:WorkflowStep[]):ResilienceFinding[] {
  const failedIndex=steps.findIndex(s=>s.status==='FAILED');
  if(failedIndex<0) return [];
  return steps.slice(0,failedIndex)
    .filter(s=>s.status==='SUCCEEDED' && s.compensatable)
    .map(s=>({severity:'HIGH' as const,code:'COMPENSATION_REQUIRED',detail:s.id}));
}
