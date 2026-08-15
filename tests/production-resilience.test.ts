import {describe,it,expect} from 'vitest';
import {qualifyFailover,qualifyIncidentTransition,qualifyReleasePromotion} from '../lib/production-resilience';
describe('RC-000143 production resilience',()=>{
 it('fails closed on stale disaster-recovery replica',()=>expect(qualifyFailover({primary:'FAILED',replica:'HEALTHY',replicaLagSeconds:91,maxLagSeconds:60,authorized:true,integrityVerified:true,tenantIsolationVerified:true}).allowed).toBe(false));
 it('requires failover authorization and integrity',()=>expect(qualifyFailover({primary:'FAILED',replica:'HEALTHY',replicaLagSeconds:1,maxLagSeconds:60,authorized:false,integrityVerified:false,tenantIsolationVerified:true}).allowed).toBe(false));
 it('blocks resolving incident without evidence',()=>expect(qualifyIncidentTransition({from:'RECOVERING',to:'RESOLVED',severity:'SEV1',authorized:true,evidenceAttached:false}).allowed).toBe(false));
 it('blocks invalid incident state jumps',()=>expect(qualifyIncidentTransition({from:'OPEN',to:'RESOLVED',severity:'SEV2',authorized:true,evidenceAttached:true}).allowed).toBe(false));
 it('blocks release when artifact hash differs',()=>expect(qualifyReleasePromotion({artifactHash:'a',expectedHash:'b',sourceCommit:'c',attestedCommit:'c',testsPassed:true,externalGatesPassed:true,rollbackArtifactVerified:true}).promotable).toBe(false));
 it('blocks release while external gates are open',()=>expect(qualifyReleasePromotion({artifactHash:'a',expectedHash:'a',sourceCommit:'c',attestedCommit:'c',testsPassed:true,externalGatesPassed:false,rollbackArtifactVerified:true}).promotable).toBe(false));
 it('requires verified rollback artifact',()=>expect(qualifyReleasePromotion({artifactHash:'a',expectedHash:'a',sourceCommit:'c',attestedCommit:'c',testsPassed:true,externalGatesPassed:true,rollbackArtifactVerified:false}).promotable).toBe(false));
});
