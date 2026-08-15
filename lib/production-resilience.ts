export type IncidentSeverity='SEV1'|'SEV2'|'SEV3'|'SEV4';
export type Health='HEALTHY'|'DEGRADED'|'FAILED';
export interface FailoverRequest { primary:Health; replica:Health; replicaLagSeconds:number; maxLagSeconds:number; authorized:boolean; integrityVerified:boolean; tenantIsolationVerified:boolean; }
export interface IncidentTransition { from:'OPEN'|'CONTAINED'|'RECOVERING'|'RESOLVED'; to:'OPEN'|'CONTAINED'|'RECOVERING'|'RESOLVED'; severity:IncidentSeverity; authorized:boolean; evidenceAttached:boolean; }
export interface ReleasePromotion { artifactHash:string; expectedHash:string; sourceCommit:string; attestedCommit:string; testsPassed:boolean; externalGatesPassed:boolean; rollbackArtifactVerified:boolean; }

export function qualifyFailover(r:FailoverRequest){
 const findings:string[]=[];
 if(r.primary!=='FAILED') findings.push('PRIMARY_NOT_FAILED');
 if(r.replica!=='HEALTHY') findings.push('REPLICA_NOT_HEALTHY');
 if(r.replicaLagSeconds>r.maxLagSeconds) findings.push('RPO_EXCEEDED');
 if(!r.authorized) findings.push('FAILOVER_NOT_AUTHORIZED');
 if(!r.integrityVerified) findings.push('FAILOVER_INTEGRITY_UNVERIFIED');
 if(!r.tenantIsolationVerified) findings.push('FAILOVER_TENANT_ISOLATION_UNVERIFIED');
 return {allowed:findings.length===0,findings};
}
export function qualifyIncidentTransition(t:IncidentTransition){
 const findings:string[]=[];
 if(!t.authorized) findings.push('INCIDENT_TRANSITION_UNAUTHORIZED');
 if((t.to==='RECOVERING'||t.to==='RESOLVED')&&!t.evidenceAttached) findings.push('INCIDENT_EVIDENCE_REQUIRED');
 const allowed:Record<string,string[]>= {OPEN:['CONTAINED'],CONTAINED:['RECOVERING'],RECOVERING:['CONTAINED','RESOLVED'],RESOLVED:[]};
 if(!allowed[t.from].includes(t.to)) findings.push('INCIDENT_STATE_TRANSITION_INVALID');
 return {allowed:findings.length===0,findings};
}
export function qualifyReleasePromotion(r:ReleasePromotion){
 const findings:string[]=[];
 if(!r.artifactHash||r.artifactHash!==r.expectedHash) findings.push('ARTIFACT_HASH_MISMATCH');
 if(!r.sourceCommit||r.sourceCommit!==r.attestedCommit) findings.push('SOURCE_ATTESTATION_MISMATCH');
 if(!r.testsPassed) findings.push('TEST_GATE_OPEN');
 if(!r.externalGatesPassed) findings.push('EXTERNAL_GATE_OPEN');
 if(!r.rollbackArtifactVerified) findings.push('ROLLBACK_ARTIFACT_UNVERIFIED');
 return {promotable:findings.length===0,findings};
}
