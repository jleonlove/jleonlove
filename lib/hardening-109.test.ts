import {describe,it,expect} from 'vitest';
import {sealRecord,verifyRecord} from './data-integrity';
import {createSnapshot,qualifyRestore} from './backup-recovery';
import {rateDecision,detectReplay} from './abuse-rate-guard';
describe('RC-000109 integrity/recovery/abuse hardening',()=>{
 it('detects tampered tenant-scoped records',()=>{ const r=sealRecord({schemaVersion:1,recordId:'d1',organizationId:'o1',workspaceId:'w1',updatedAt:new Date().toISOString(),payload:{value:10}}); const bad={...r,payload:{value:11}}; expect(verifyRecord(bad,'o1','w1').valid).toBe(false); });
 it('rejects corrupted or cross-tenant backups',()=>{ const s=createSnapshot('s1','o1','w1',[{id:1}]); expect(qualifyRestore(s,'o1','w1').safeToRestore).toBe(true); expect(qualifyRestore({...s,recordCount:2},'o2','w1').safeToRestore).toBe(false); });
 it('rate limits abusive consequential actions',()=>{ const now=100000; const samples=Array.from({length:3},(_,i)=>({actorId:'a',action:'PAY',timestamp:now-100+i})); expect(rateDecision(samples,'a','PAY',now,3,1000).allow).toBe(false); });
 it('detects replayed nonces',()=>{ expect(detectReplay(['n1','n2','n1']).safe).toBe(false); });
});
