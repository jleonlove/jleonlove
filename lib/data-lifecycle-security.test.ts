import { describe,it,expect } from 'vitest';
import { qualifyKeyRotation,qualifyRestoreRequest,qualifyDeletion,deletionEvidence } from './data-lifecycle-security';

describe('data lifecycle security',()=>{
 it('rejects key reuse and rollback',()=>{ const r=qualifyKeyRotation({keyId:'k1',purpose:'backup',version:2,state:'ACTIVE',createdAt:'2026-01-01T00:00:00Z'},{keyId:'k1',purpose:'backup',version:1,state:'ACTIVE',createdAt:'2026-01-02T00:00:00Z'}); expect(r.allowed).toBe(false); expect(r.findings).toContain('KEY_ID_REUSE'); });
 it('rejects cross-tenant restore',()=>{ expect(qualifyRestoreRequest({organizationId:'a',workspaceId:'w',snapshotOrganizationId:'b',snapshotWorkspaceId:'w',snapshotChecksum:'x',computedChecksum:'x',encryptionKeyState:'ACTIVE',approved:true}).safeToRestore).toBe(false); });
 it('rejects tampered restore',()=>{ expect(qualifyRestoreRequest({organizationId:'a',workspaceId:'w',snapshotOrganizationId:'a',snapshotWorkspaceId:'w',snapshotChecksum:'x',computedChecksum:'y',encryptionKeyState:'ACTIVE',approved:true}).safeToRestore).toBe(false); });
 it('rejects restore using revoked key',()=>{ expect(qualifyRestoreRequest({organizationId:'a',workspaceId:'w',snapshotOrganizationId:'a',snapshotWorkspaceId:'w',snapshotChecksum:'x',computedChecksum:'x',encryptionKeyState:'REVOKED',approved:true}).safeToRestore).toBe(false); });
 it('honors legal hold',()=>{ expect(qualifyDeletion({recordId:'r',organizationId:'a',deleteAfter:'2020-01-01T00:00:00Z',legalHold:true},new Date('2026-01-01')).deletable).toBe(false); });
 it('prevents early deletion',()=>{ expect(qualifyDeletion({recordId:'r',organizationId:'a',deleteAfter:'2030-01-01T00:00:00Z',legalHold:false},new Date('2026-01-01')).deletable).toBe(false); });
 it('creates deletion evidence without content',()=>{ const e=deletionEvidence({recordId:'r',organizationId:'a',deleteAfter:'2020-01-01T00:00:00Z',legalHold:false},'2026-01-01T00:00:00Z'); expect(e.evidenceHash.length).toBeGreaterThan(10); expect((e as any).content).toBeUndefined(); });
});
