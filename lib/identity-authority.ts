export type AtlasRole='member'|'admin'|'owner'|'security_admin'|'billing_admin';
export type SensitiveAction='READ_RESOURCE'|'MANAGE_USERS'|'TRANSFER_OWNERSHIP'|'CHANGE_MERCHANT_DESTINATION'|'ROTATE_CREDENTIAL'|'EXPORT_ORG'|'DELETE_ORG';
export interface AuthorityContext {subjectId:string; organizationId:string; workspaceId:string; role:AtlasRole; sessionId:string; sessionRevoked?:boolean; recoveryRestricted?:boolean;}
export interface ScopedResource {organizationId:string; workspaceId:string;}
const required:Record<SensitiveAction,AtlasRole[]>={READ_RESOURCE:['member','admin','owner','security_admin','billing_admin'],MANAGE_USERS:['admin','owner'],TRANSFER_OWNERSHIP:['owner'],CHANGE_MERCHANT_DESTINATION:['owner','billing_admin'],ROTATE_CREDENTIAL:['owner','security_admin'],EXPORT_ORG:['owner','admin'],DELETE_ORG:['owner']};
export function authorize(ctx:AuthorityContext,resource:ScopedResource,action:SensitiveAction){
 if(!ctx.subjectId||!ctx.sessionId) throw new Error('IDENTITY_REQUIRED');
 if(ctx.sessionRevoked) throw new Error('SESSION_REVOKED');
 if(ctx.organizationId!==resource.organizationId||ctx.workspaceId!==resource.workspaceId) throw new Error('TENANT_ISOLATION_DENY');
 if(ctx.recoveryRestricted && action!=='READ_RESOURCE') throw new Error('RECOVERY_PRIVILEGE_DENY');
 if(!required[action].includes(ctx.role)) throw new Error('AUTHORITY_DENY');
 return true;
}
export interface SessionRecord {sessionId:string; subjectId:string; organizationId:string; revokedAt?:string; reason?:string;}
export function revokeSubjectSessions(sessions:SessionRecord[],subjectId:string,organizationId:string,reason='ACCOUNT_RECOVERY',now=new Date().toISOString()){
 return sessions.map(s=>s.subjectId===subjectId&&s.organizationId===organizationId?{...s,revokedAt:s.revokedAt??now,reason:s.reason??reason}:s);
}
export function assertCredentialRotation(oldCredentialId:string,newCredentialId:string){
 if(!oldCredentialId||!newCredentialId) throw new Error('CREDENTIAL_ID_REQUIRED');
 if(oldCredentialId===newCredentialId) throw new Error('CREDENTIAL_ROTATION_REUSE_DENY');
 return true;
}
