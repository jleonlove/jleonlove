import { checksumPayload } from './data-integrity';
export interface BackupSnapshot<T=unknown>{ snapshotId:string; createdAt:string; organizationId:string; workspaceId:string; records:T[]; recordCount:number; checksum:string; }
export interface RestoreQualification{ safeToRestore:boolean; findings:string[]; }
export function createSnapshot<T>(snapshotId:string,organizationId:string,workspaceId:string,records:T[]):BackupSnapshot<T>{
 const base={snapshotId,createdAt:new Date().toISOString(),organizationId,workspaceId,records,recordCount:records.length};
 return {...base,checksum:checksumPayload(base)};
}
export function qualifyRestore<T>(s:BackupSnapshot<T>,organizationId:string,workspaceId:string):RestoreQualification{
 const findings:string[]=[]; const {checksum,...base}=s;
 if(checksumPayload(base)!==checksum) findings.push('BACKUP_CHECKSUM_MISMATCH');
 if(s.organizationId!==organizationId||s.workspaceId!==workspaceId) findings.push('BACKUP_TENANT_SCOPE_MISMATCH');
 if(s.recordCount!==s.records.length) findings.push('BACKUP_RECORD_COUNT_MISMATCH');
 if(!s.snapshotId||Number.isNaN(Date.parse(s.createdAt))) findings.push('BACKUP_METADATA_INVALID');
 return {safeToRestore:findings.length===0,findings};
}
