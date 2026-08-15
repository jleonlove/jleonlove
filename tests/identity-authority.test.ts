import {describe,it,expect} from 'vitest';
import {authorize,revokeSubjectSessions,assertCredentialRotation} from '../lib/identity-authority';
const resource={organizationId:'org-a',workspaceId:'ws-a'};
const owner={subjectId:'u1',organizationId:'org-a',workspaceId:'ws-a',role:'owner' as const,sessionId:'s1'};
describe('identity authority fail closed',()=>{
 it('allows owner in exact tenant',()=>expect(authorize(owner,resource,'DELETE_ORG')).toBe(true));
 it('denies cross tenant',()=>expect(()=>authorize(owner,{...resource,organizationId:'org-b'},'READ_RESOURCE')).toThrow('TENANT_ISOLATION_DENY'));
 it('denies revoked session',()=>expect(()=>authorize({...owner,sessionRevoked:true},resource,'READ_RESOURCE')).toThrow('SESSION_REVOKED'));
 it('denies member privilege escalation',()=>expect(()=>authorize({...owner,role:'member'},resource,'TRANSFER_OWNERSHIP')).toThrow('AUTHORITY_DENY'));
 it('restricts recovered sessions from privilege',()=>expect(()=>authorize({...owner,recoveryRestricted:true},resource,'CHANGE_MERCHANT_DESTINATION')).toThrow('RECOVERY_PRIVILEGE_DENY'));
 it('revokes all subject org sessions only',()=>{const r=revokeSubjectSessions([{sessionId:'1',subjectId:'u1',organizationId:'org-a'},{sessionId:'2',subjectId:'u1',organizationId:'org-b'}],'u1','org-a');expect(Boolean(r[0].revokedAt)).toBe(true);expect(r[1].revokedAt).toBe(undefined)});
 it('denies credential reuse during rotation',()=>expect(()=>assertCredentialRotation('key1','key1')).toThrow('CREDENTIAL_ROTATION_REUSE_DENY'));
});
