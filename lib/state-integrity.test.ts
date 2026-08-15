import {describe,it,expect} from 'vitest';
import {validateDealMutation,requireApprovalForMutation,cannotBypassClosingGate} from './state-integrity';
import type {AtlasDeal,DealEvent} from './types';
const deal:AtlasDeal={id:'d1',organizationId:'o1',workspaceId:'w1',createdBy:'u1',title:'Gold',commodity:'gold',stage:'DD',participants:[],evidence:[],requirements:[{id:'authority',title:'Verify authority',status:'BLOCKED',severity:'CRITICAL',ownerRole:'seller',dependsOn:[],evidenceIds:[]}],assumptions:[],createdAt:'2026-08-14T00:00:00Z',updatedAt:'2026-08-14T00:00:00Z'};
const event=(payload:Record<string,unknown>):DealEvent=>({id:'e1',dealId:'d1',type:'CHANGE',occurredAt:'2026-08-14T01:00:00Z',payload,idempotencyKey:'k1'});
describe('state integrity',()=>{
 it('denies tenant mutation',()=>expect(validateDealMutation(deal,event({organizationId:'o2'})).pass).toBe(false));
 it('requires scoped approval for beneficiary changes',()=>expect(requireApprovalForMutation(event({beneficiaryAccount:'x'})).pass).toBe(false));
 it('prevents closing-gate bypass',()=>expect(cannotBypassClosingGate(deal,'execute settlement').allowed).toBe(false));
});
