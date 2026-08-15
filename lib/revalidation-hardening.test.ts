import {describe,it,expect} from 'vitest';
import type {AtlasDeal} from './types';
import {cascadeEvidenceRevocation,quarantineCounterparty,revalidationGate,type ApprovalSnapshot} from './revalidation-hardening';
const base:AtlasDeal={id:'d',organizationId:'o',workspaceId:'w',createdBy:'u',title:'Gold',commodity:'gold',stage:'closing',participants:[{id:'seller',name:'Seller',role:'seller',authority:'INDEPENDENTLY_VERIFIED'}],evidence:[{id:'e1',claim:'Seller controls product',status:'INDEPENDENTLY_VERIFIED',sourceIds:['s1','s2'],independentSourceGroups:2,observedAt:'2026-08-14'}],requirements:[{id:'r1',title:'Seller authority',status:'SATISFIED',severity:'CRITICAL',ownerRole:'seller',dependsOn:[],evidenceIds:['e1']}],assumptions:[],createdAt:'2026-08-14T00:00:00Z',updatedAt:'2026-08-14T00:00:00Z'};
describe('RC-000112 revalidation',()=>{
 it('cascades revoked evidence into satisfied requirements',()=>{const d=cascadeEvidenceRevocation(base,['e1']);expect(d.evidence[0].status).toBe('CONTRADICTED');expect(d.requirements[0].status).toBe('BLOCKED');});
 it('quarantines compromised counterparty and fails closed',()=>{const d=quarantineCounterparty(base,'seller');expect(d.participants[0].authority).toBe('CONTRADICTED');expect(revalidationGate(d,'COUNTERPARTY_COMPROMISED',[]).pass).toBe(false);});
 it('invalidates approvals after evidence changes',()=>{const a:ApprovalSnapshot={id:'a1',requirementId:'r1',approvedAt:'2026-08-14T00:00:00Z',evidenceIds:['e1'],dealUpdatedAt:'2026-08-14T00:00:00Z',status:'ACTIVE'};const d={...base,updatedAt:'2026-08-14T01:00:00Z'};expect(revalidationGate(d,'MATERIAL_CHANGE',[a],['e1']).staleApprovalIds).toContain('a1');});
});
