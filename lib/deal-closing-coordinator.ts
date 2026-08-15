import { createHash, randomUUID } from 'crypto';
import { validateExecutionEnvelope } from './deal-execution-envelope';
import type { DealExecutionEnvelope, ExecutionSnapshot } from './deal-execution-envelope';

export type DealPhase='OPEN'|'CLOSING'|'CLOSED'|'CANCELLED'|'EXCEPTION';
export type CommitEvidence={
  evidenceId:string;
  type:string;
  digest:string;
  dealId:string;
  workspaceId:string;
  stateVersion:number;
  verified:boolean;
  revoked?:boolean;
};
export type PreparedClose={
  idempotencyKey:string;
  stateVersion:number;
  envelopeDigest:string;
  evidenceDigest:string;
  preparedAt:number;
};
export type DealRecord={
  dealId:string;
  organizationId:string;
  workspaceId:string;
  version:number;
  phase:DealPhase;
  lease?:{token:string;expiresAt:number};
  executions:Record<string,string>;
  preparedClose?:PreparedClose;
};
export type CloseRequest={
  envelope:DealExecutionEnvelope;
  /** Snapshot submitted by the caller; never treated as authoritative by itself. */
  snapshot:ExecutionSnapshot;
  /** Snapshot read from Atlas's trusted transaction/evidence store at execution time. */
  trustedSnapshot:ExecutionSnapshot;
  expectedVersion:number;
  requiredAuthority:string[];
  /** Current, independently resolved authority grants for the acting principal. */
  trustedAuthorityScope:string[];
  paymentState:'NONE'|'PENDING'|'SETTLED'|'REVERSED';
  cancelRequested?:boolean;
  commitEvidence:CommitEvidence[];
  now?:Date;
};
const digest=(v:unknown)=>createHash('sha256').update(JSON.stringify(v)).digest('hex');
const snapshotDigest=(v:ExecutionSnapshot)=>digest(v);
const evidenceDigest=(e:CommitEvidence[])=>digest([...e].sort((a,b)=>a.evidenceId.localeCompare(b.evidenceId)));

export class DealClosingCoordinator {
 private record:DealRecord;
 constructor(record:DealRecord){this.record=record}
 acquireLease(expectedVersion:number,now=Date.now(),ttlMs=30_000){
  if(this.record.version!==expectedVersion) return {pass:false,error:'VERSION_CONFLICT'} as const;
  if(this.record.lease && this.record.lease.expiresAt>now) return {pass:false,error:'CLOSING_LEASE_HELD'} as const;
  const token=randomUUID(); this.record.lease={token,expiresAt:now+ttlMs}; return {pass:true,token} as const;
 }
 close(req:CloseRequest,leaseToken:string){
  const pre=this.preflight(req,leaseToken); if(!pre.pass) return pre;
  const now=req.now??new Date();
  const prior=this.record.executions[req.envelope.idempotencyKey]; if(prior) return {pass:true,replayed:true,receipt:prior};
  if(req.cancelRequested && req.paymentState==='SETTLED'){ this.record.phase='EXCEPTION'; this.record.version++; delete this.record.lease; return this.fail('SETTLEMENT_CANCELLATION_RACE'); }
  if(req.cancelRequested){ this.record.phase='CANCELLED'; this.record.version++; delete this.record.lease; return this.fail('CLOSING_CANCELLED'); }
  if(req.paymentState==='REVERSED') return this.fail('PAYMENT_REVERSED');
  const ev=this.validateCommitEvidence(req.commitEvidence,req.trustedSnapshot); if(!ev.pass) return ev;

  const envelopeDigest=digest(req.envelope);
  const evDigest=evidenceDigest(req.commitEvidence);
  this.record.phase='CLOSING';
  this.record.preparedClose={idempotencyKey:req.envelope.idempotencyKey,stateVersion:this.record.version,envelopeDigest,evidenceDigest:evDigest,preparedAt:now.getTime()};

  const receipt=digest({dealId:this.record.dealId,version:this.record.version,key:req.envelope.idempotencyKey,envelopeDigest,evidenceDigest:evDigest});
  this.record.phase='CLOSED'; this.record.executions[req.envelope.idempotencyKey]=receipt; this.record.version++;
  delete this.record.preparedClose; delete this.record.lease;
  return {pass:true,replayed:false,receipt};
 }
 /** Recover a close that was durably PREPARED before a worker/process failure. */
 recoverPrepared(req:CloseRequest,leaseToken:string){
  const pre=this.preflight(req,leaseToken); if(!pre.pass) return pre;
  const prepared=this.record.preparedClose; if(!prepared) return this.fail('NO_PREPARED_CLOSE');
  const envDigest=digest(req.envelope); const evDigest=evidenceDigest(req.commitEvidence);
  if(prepared.idempotencyKey!==req.envelope.idempotencyKey) return this.fail('RECOVERY_IDEMPOTENCY_MISMATCH');
  if(prepared.stateVersion!==this.record.version) return this.fail('RECOVERY_VERSION_MISMATCH');
  if(prepared.envelopeDigest!==envDigest) return this.fail('RECOVERY_ENVELOPE_MISMATCH');
  if(prepared.evidenceDigest!==evDigest) return this.fail('RECOVERY_EVIDENCE_MISMATCH');
  const ev=this.validateCommitEvidence(req.commitEvidence,req.trustedSnapshot); if(!ev.pass) return ev;
  const receipt=digest({dealId:this.record.dealId,version:this.record.version,key:req.envelope.idempotencyKey,envelopeDigest:envDigest,evidenceDigest:evDigest});
  this.record.phase='CLOSED'; this.record.executions[req.envelope.idempotencyKey]=receipt; this.record.version++;
  delete this.record.preparedClose; delete this.record.lease;
  return {pass:true,recovered:true,receipt};
 }
 snapshot(){return structuredClone(this.record)}
 private preflight(req:CloseRequest,leaseToken:string){
  const now=req.now??new Date();
  if(this.record.version!==req.expectedVersion) return this.fail('VERSION_CONFLICT');
  if(!this.record.lease||this.record.lease.token!==leaseToken||this.record.lease.expiresAt<=now.getTime()) return this.fail('INVALID_OR_EXPIRED_LEASE');
  if(req.trustedSnapshot.dealId!==this.record.dealId) return this.fail('DEAL_IDENTITY_MISMATCH');
  if(req.trustedSnapshot.organizationId!==this.record.organizationId) return this.fail('ORGANIZATION_BOUNDARY_VIOLATION');
  if(req.trustedSnapshot.workspaceId!==this.record.workspaceId) return this.fail('WORKSPACE_BOUNDARY_VIOLATION');
  if(req.trustedSnapshot.stateVersion!==this.record.version) return this.fail('TRUSTED_STATE_VERSION_MISMATCH');
  if(snapshotDigest(req.snapshot)!==snapshotDigest(req.trustedSnapshot)) return this.fail('UNTRUSTED_SNAPSHOT_MISMATCH');
  const valid=validateExecutionEnvelope(req.envelope,req.trustedSnapshot,now); if(!valid.pass) return {pass:false,errors:valid.errors};
  if(req.requiredAuthority.some(a=>!req.trustedAuthorityScope.includes(a))) return this.fail('AUTHORITY_REVALIDATION_FAILED');
  if(req.envelope.authorityScope.some(a=>!req.trustedAuthorityScope.includes(a))) return this.fail('AUTHORITY_SCOPE_ESCALATION');
  return {pass:true} as const;
 }
 private validateCommitEvidence(items:CommitEvidence[],trusted:ExecutionSnapshot){
  if(items.length===0) return this.fail('MISSING_COMMIT_EVIDENCE');
  const ids=new Set<string>();
  for(const e of items){
   if(!e.evidenceId||!e.digest) return this.fail('MALFORMED_COMMIT_EVIDENCE');
   if(ids.has(e.evidenceId)) return this.fail('DUPLICATE_COMMIT_EVIDENCE'); ids.add(e.evidenceId);
   if(e.dealId!==this.record.dealId) return this.fail('EVIDENCE_DEAL_MISMATCH');
   if(e.workspaceId!==this.record.workspaceId) return this.fail('EVIDENCE_WORKSPACE_MISMATCH');
   if(e.stateVersion!==this.record.version||e.stateVersion!==trusted.stateVersion) return this.fail('STALE_COMMIT_EVIDENCE');
   if(!e.verified) return this.fail('UNVERIFIED_COMMIT_EVIDENCE');
   if(e.revoked) return this.fail('REVOKED_COMMIT_EVIDENCE');
  }
  return {pass:true} as const;
 }
 private fail(error:string){return {pass:false,errors:[error]};}
}

export type SideEffectState='PLANNED'|'APPLIED'|'COMPENSATED';
export type SideEffect={effectId:string;service:string;operation:string;digest:string;state:SideEffectState;compensationDigest?:string};
export type RecoveryCheckpoint={checkpointId:string;dealId:string;workspaceId:string;stateVersion:number;envelopeDigest:string;evidenceDigest:string;authorityDigest:string;createdAt:number;expiresAt:number;sideEffects:SideEffect[]};
const authorityDigest=(scope:string[])=>digest([...scope].sort());

/** Durable saga guard for long-running closing work spanning external services. */
export class ClosingRecoveryGuard {
 private seenReplay=new Set<string>();
 private record:DealRecord; private checkpoint?:RecoveryCheckpoint;
 constructor(record:DealRecord, checkpoint?:RecoveryCheckpoint){this.record=record;this.checkpoint=checkpoint}
 prepare(req:CloseRequest,ttlMs=300_000):RecoveryCheckpoint{
  const now=(req.now??new Date()).getTime();
  const cp:RecoveryCheckpoint={checkpointId:randomUUID(),dealId:this.record.dealId,workspaceId:this.record.workspaceId,stateVersion:this.record.version,envelopeDigest:digest(req.envelope),evidenceDigest:evidenceDigest(req.commitEvidence),authorityDigest:authorityDigest(req.trustedAuthorityScope),createdAt:now,expiresAt:now+ttlMs,sideEffects:[]};
  this.checkpoint=cp; return structuredClone(cp);
 }
 registerSideEffect(effect:Omit<SideEffect,'state'>){
  if(!this.checkpoint) return {pass:false,error:'NO_RECOVERY_CHECKPOINT'} as const;
  if(this.checkpoint.sideEffects.some(e=>e.effectId===effect.effectId)) return {pass:false,error:'DUPLICATE_SIDE_EFFECT'} as const;
  this.checkpoint.sideEffects.push({...effect,state:'PLANNED'}); return {pass:true} as const;
 }
 markApplied(effectId:string){const e=this.checkpoint?.sideEffects.find(x=>x.effectId===effectId);if(!e)return {pass:false,error:'UNKNOWN_SIDE_EFFECT'} as const;e.state='APPLIED';return {pass:true} as const;}
 compensate(effectId:string,compensationDigest:string){const e=this.checkpoint?.sideEffects.find(x=>x.effectId===effectId);if(!e)return {pass:false,error:'UNKNOWN_SIDE_EFFECT'} as const;if(e.state!=='APPLIED')return {pass:false,error:'SIDE_EFFECT_NOT_APPLIED'} as const;e.state='COMPENSATED';e.compensationDigest=compensationDigest;return {pass:true} as const;}
 validateResume(req:CloseRequest,now=Date.now()){
  const cp=this.checkpoint;if(!cp)return {pass:false,error:'NO_RECOVERY_CHECKPOINT'} as const;
  if(cp.expiresAt<=now)return {pass:false,error:'STALE_RECOVERY_CHECKPOINT'} as const;
  if(cp.dealId!==this.record.dealId||cp.workspaceId!==this.record.workspaceId)return {pass:false,error:'RECOVERY_BOUNDARY_MISMATCH'} as const;
  if(cp.stateVersion!==this.record.version||cp.stateVersion!==req.trustedSnapshot.stateVersion)return {pass:false,error:'RECOVERY_STATE_CHANGED'} as const;
  if(cp.envelopeDigest!==digest(req.envelope)||cp.evidenceDigest!==evidenceDigest(req.commitEvidence))return {pass:false,error:'RECOVERY_INPUT_CHANGED'} as const;
  if(cp.authorityDigest!==authorityDigest(req.trustedAuthorityScope))return {pass:false,error:'AUTHORITY_CHANGED_DURING_WORKFLOW'} as const;
  if(cp.sideEffects.some(e=>e.state==='APPLIED'))return {pass:false,error:'UNCOMPENSATED_SIDE_EFFECTS'} as const;
  return {pass:true} as const;
 }
 acceptServiceMessage(service:string,messageId:string,payloadDigest:string){
  const key=digest({dealId:this.record.dealId,workspaceId:this.record.workspaceId,service,messageId,payloadDigest});
  if(this.seenReplay.has(key))return {pass:false,error:'CROSS_SERVICE_REPLAY'} as const;
  this.seenReplay.add(key);return {pass:true,replayKey:key} as const;
 }
 snapshot(){return this.checkpoint?structuredClone(this.checkpoint):undefined;}
}
