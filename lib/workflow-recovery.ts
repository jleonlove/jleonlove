export type WorkflowStepStatus='PENDING'|'RUNNING'|'SUCCEEDED'|'FAILED'|'COMPENSATED';
export interface WorkflowStep { id:string; status:WorkflowStepStatus; attempts:number; maxAttempts:number; idempotencyKey:string; compensatable:boolean; }
export interface RecoveryPlan { retry:string[]; compensate:string[]; halt:boolean; reasons:string[]; }
export function planRecovery(steps:WorkflowStep[]):RecoveryPlan {
 const retry:string[]=[], compensate:string[]=[], reasons:string[]=[];
 const keys=new Set<string>(); let halt=false;
 for(const s of steps){
  if(keys.has(s.idempotencyKey)){halt=true;reasons.push(`Duplicate idempotency key: ${s.idempotencyKey}`);} keys.add(s.idempotencyKey);
  if(s.status==='FAILED'){
   if(s.attempts<s.maxAttempts) retry.push(s.id);
   else {halt=true;reasons.push(`Retry budget exhausted: ${s.id}`); if(s.compensatable) compensate.push(s.id);}
  }
 }
 return {retry,compensate,halt,reasons};
}
