import {describe,it,expect} from 'vitest';
import {ClosingRecoveryGuard} from '../lib/deal-closing-coordinator';
import type {DealRecord} from '../lib/deal-closing-coordinator';
import type {DealExecutionEnvelope,ExecutionSnapshot} from '../lib/deal-execution-envelope';
const record=():DealRecord=>({dealId:'d1',organizationId:'o1',workspaceId:'w1',version:7,phase:'CLOSING',executions:{}});
const req=()=>{const s:ExecutionSnapshot={dealId:'d1',organizationId:'o1',workspaceId:'w1',stateVersion:7,termsDigest:'t',documentSetDigest:'d',evidenceDigest:'e',approvalDigest:'a',authorityScope:['deal.close']};const e:DealExecutionEnvelope={...s,expiresAt:'2099-01-01T00:00:00Z',idempotencyKey:'k'};return {envelope:e,snapshot:s,trustedSnapshot:s,expectedVersion:7,requiredAuthority:['deal.close'],trustedAuthorityScope:['deal.close'],paymentState:'SETTLED' as const,commitEvidence:[{evidenceId:'x',type:'WIRE',digest:'x',dealId:'d1',workspaceId:'w1',stateVersion:7,verified:true}]};};
describe('ClosingRecoveryGuard',()=>{
 it('rejects stale recovery checkpoints',()=>{const r=req();r.now=new Date(1000) as any;const g=new ClosingRecoveryGuard(record());g.prepare(r,10);expect(g.validateResume(r,1011)).toEqual({pass:false,error:'STALE_RECOVERY_CHECKPOINT'});});
 it('detects authority revocation/change during a long workflow',()=>{const r=req();const g=new ClosingRecoveryGuard(record());g.prepare(r);expect(g.validateResume({...r,trustedAuthorityScope:[]})).toEqual({pass:false,error:'AUTHORITY_CHANGED_DURING_WORKFLOW'});});
 it('requires compensation for partial external side effects before resume',()=>{const r=req();const g=new ClosingRecoveryGuard(record());g.prepare(r);g.registerSideEffect({effectId:'pay',service:'payments',operation:'capture',digest:'p'});g.markApplied('pay');expect(g.validateResume(r)).toEqual({pass:false,error:'UNCOMPENSATED_SIDE_EFFECTS'});g.compensate('pay','refund');expect(g.validateResume(r).pass).toBe(true);});
 it('rejects changed recovery inputs',()=>{const r=req();const g=new ClosingRecoveryGuard(record());g.prepare(r);const changed={...r,envelope:{...r.envelope,termsDigest:'changed'}};expect(g.validateResume(changed)).toEqual({pass:false,error:'RECOVERY_INPUT_CHANGED'});});
 it('blocks duplicate cross-service messages',()=>{const g=new ClosingRecoveryGuard(record());expect(g.acceptServiceMessage('payments','m1','p').pass).toBe(true);expect(g.acceptServiceMessage('payments','m1','p')).toEqual({pass:false,error:'CROSS_SERVICE_REPLAY'});});
});
