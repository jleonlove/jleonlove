import {describe,it,expect} from 'vitest';
import {validateExecutionEnvelope,ExecutionReplayGuard,type DealExecutionEnvelope} from '../lib/deal-execution-envelope';
const base={dealId:'d1',organizationId:'o1',workspaceId:'w1',stateVersion:7,termsDigest:'t',documentSetDigest:'d',evidenceDigest:'e',approvalDigest:'a',authorityScope:['close']};
const env:DealExecutionEnvelope={...base,expiresAt:'2099-01-01T00:00:00Z',idempotencyKey:'once'};
describe('deal execution envelope',()=>{
 it('accepts exact approved snapshot',()=>expect(validateExecutionEnvelope(env,base).pass).toBe(true));
 it.each([
  ['stateVersion',8],['termsDigest','changed'],['documentSetDigest','changed'],['evidenceDigest','changed'],['approvalDigest','changed'],['workspaceId','other']
 ] as const)('invalidates stale approval when %s changes',(k,v)=>expect(validateExecutionEnvelope(env,{...base,[k]:v}).pass).toBe(false));
 it('rejects expired authority',()=>expect(validateExecutionEnvelope({...env,expiresAt:'2020-01-01T00:00:00Z'},base).pass).toBe(false));
 it('rejects duplicate consequential execution',()=>{const g=new ExecutionReplayGuard();expect(g.consume('x').pass).toBe(true);expect(g.consume('x').pass).toBe(false)});
});
