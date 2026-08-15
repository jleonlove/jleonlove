import { createHash } from 'crypto';

export type DealExecutionEnvelope = {
  dealId:string; organizationId:string; workspaceId:string; stateVersion:number;
  termsDigest:string; documentSetDigest:string; evidenceDigest:string; approvalDigest:string;
  authorityScope:string[]; expiresAt:string; idempotencyKey:string;
};
export type ExecutionSnapshot = Omit<DealExecutionEnvelope,'expiresAt'|'idempotencyKey'>;
const sha=(v:unknown)=>createHash('sha256').update(typeof v==='string'?v:JSON.stringify(v)).digest('hex');
export const digestCanonical=(v:unknown)=>sha(v);
export function validateExecutionEnvelope(env:DealExecutionEnvelope,current:ExecutionSnapshot,now=new Date()){
 const errors:string[]=[];
 if(new Date(env.expiresAt).getTime()<=now.getTime()) errors.push('EXECUTION_ENVELOPE_EXPIRED');
 for(const k of ['dealId','organizationId','workspaceId','stateVersion','termsDigest','documentSetDigest','evidenceDigest','approvalDigest'] as const)
   if(env[k]!==current[k]) errors.push(`STALE_${String(k).toUpperCase()}`);
 if(!env.idempotencyKey) errors.push('MISSING_IDEMPOTENCY_KEY');
 if(!env.authorityScope.length) errors.push('MISSING_AUTHORITY_SCOPE');
 const scopeMissing=current.authorityScope.some(s=>!env.authorityScope.includes(s));
 if(scopeMissing) errors.push('AUTHORITY_SCOPE_CHANGED');
 return {pass:errors.length===0,errors};
}
export class ExecutionReplayGuard {
 private used=new Set<string>();
 consume(key:string){ if(!key) return {pass:false,error:'MISSING_IDEMPOTENCY_KEY'}; if(this.used.has(key)) return {pass:false,error:'DUPLICATE_EXECUTION'}; this.used.add(key); return {pass:true}; }
}
