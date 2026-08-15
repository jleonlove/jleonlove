import { checksumPayload } from './data-integrity';

export type KeyState='ACTIVE'|'RETIRED'|'REVOKED';
export interface KeyRecord { keyId:string; purpose:string; version:number; state:KeyState; createdAt:string; }
export interface RetentionRecord { recordId:string; organizationId:string; deleteAfter:string; legalHold:boolean; deletedAt?:string; }
export interface RestoreRequest { organizationId:string; workspaceId:string; snapshotOrganizationId:string; snapshotWorkspaceId:string; snapshotChecksum:string; computedChecksum:string; encryptionKeyState:KeyState; approved:boolean; }

export function qualifyKeyRotation(current:KeyRecord,next:KeyRecord){
 const findings:string[]=[];
 if(current.state!=='ACTIVE') findings.push('CURRENT_KEY_NOT_ACTIVE');
 if(next.state!=='ACTIVE') findings.push('NEXT_KEY_NOT_ACTIVE');
 if(current.keyId===next.keyId) findings.push('KEY_ID_REUSE');
 if(next.version<=current.version) findings.push('KEY_VERSION_NOT_MONOTONIC');
 if(current.purpose!==next.purpose) findings.push('KEY_PURPOSE_MISMATCH');
 return {allowed:findings.length===0,findings};
}

export function qualifyRestoreRequest(r:RestoreRequest){
 const findings:string[]=[];
 if(!r.approved) findings.push('RESTORE_NOT_APPROVED');
 if(r.organizationId!==r.snapshotOrganizationId||r.workspaceId!==r.snapshotWorkspaceId) findings.push('RESTORE_TENANT_SCOPE_MISMATCH');
 if(r.snapshotChecksum!==r.computedChecksum) findings.push('RESTORE_CHECKSUM_MISMATCH');
 if(r.encryptionKeyState!=='ACTIVE') findings.push('RESTORE_KEY_NOT_ACTIVE');
 return {safeToRestore:findings.length===0,findings};
}

export function qualifyDeletion(record:RetentionRecord,now=new Date()){
 const findings:string[]=[];
 if(record.legalHold) findings.push('LEGAL_HOLD_ACTIVE');
 if(Number.isNaN(Date.parse(record.deleteAfter))) findings.push('RETENTION_DATE_INVALID');
 else if(now.getTime()<Date.parse(record.deleteAfter)) findings.push('RETENTION_PERIOD_ACTIVE');
 if(record.deletedAt) findings.push('ALREADY_DELETED');
 return {deletable:findings.length===0,findings};
}

export function deletionEvidence(record:RetentionRecord,deletedAt:string){
 return {recordId:record.recordId,organizationId:record.organizationId,deletedAt,evidenceHash:checksumPayload({recordId:record.recordId,organizationId:record.organizationId,deletedAt})};
}
