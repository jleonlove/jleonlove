import { createHash } from 'node:crypto';

export interface IntegrityEnvelope<T=unknown>{ schemaVersion:number; recordId:string; organizationId:string; workspaceId:string; updatedAt:string; payload:T; checksum:string; }
export interface IntegrityResult{ valid:boolean; findings:string[]; }

function canonical(v:unknown):string{
 if(v===null||typeof v!=='object') return JSON.stringify(v);
 if(Array.isArray(v)) return `[${v.map(canonical).join(',')}]`;
 const o=v as Record<string,unknown>;
 return `{${Object.keys(o).sort().map(k=>`${JSON.stringify(k)}:${canonical(o[k])}`).join(',')}}`;
}
export function checksumPayload(v:unknown){ return createHash('sha256').update(canonical(v)).digest('hex'); }
export function sealRecord<T>(input:Omit<IntegrityEnvelope<T>,'checksum'>):IntegrityEnvelope<T>{ return {...input,checksum:checksumPayload(input)}; }
export function verifyRecord<T>(record:IntegrityEnvelope<T>, expectedOrg?:string, expectedWorkspace?:string):IntegrityResult{
 const findings:string[]=[];
 const {checksum,...unsigned}=record;
 if(checksumPayload(unsigned)!==checksum) findings.push('CHECKSUM_MISMATCH');
 if(!Number.isInteger(record.schemaVersion)||record.schemaVersion<1) findings.push('INVALID_SCHEMA_VERSION');
 if(!record.recordId||!record.organizationId||!record.workspaceId) findings.push('MISSING_IDENTITY_SCOPE');
 if(expectedOrg&&record.organizationId!==expectedOrg) findings.push('ORGANIZATION_SCOPE_MISMATCH');
 if(expectedWorkspace&&record.workspaceId!==expectedWorkspace) findings.push('WORKSPACE_SCOPE_MISMATCH');
 if(Number.isNaN(Date.parse(record.updatedAt))) findings.push('INVALID_UPDATED_AT');
 return {valid:findings.length===0,findings};
}
