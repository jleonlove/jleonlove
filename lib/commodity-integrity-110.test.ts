import {describe,it,expect} from 'vitest'; import {inspectCommodityIntegrity} from './commodity-integrity-110';
const base:any={id:'d',organizationId:'o',workspaceId:'w',createdBy:'u',title:'Gold',commodity:'gold',stage:'DD',participants:[],requirements:[],assumptions:[],createdAt:'2026-01-01',updatedAt:'2026-01-01'};
describe('commodity integrity 110',()=>{
 it('locks material quantity conflicts',()=>{const r=inspectCommodityIntegrity({...base,evidence:[{id:'a',claim:'Seller offers 200 kg monthly',status:'CLAIMED',sourceIds:['s'],independentSourceGroups:1,observedAt:'x'},{id:'b',claim:'Facility capacity 500 kg monthly',status:'DOCUMENT_SUPPORTED',sourceIds:['f'],independentSourceGroups:1,observedAt:'x'}]});expect(r.pass).toBe(false);expect(r.findings.some(x=>x.code==='MATERIAL_QUANTITY_CONFLICT')).toBe(true)});
 it('locks beneficiary changes',()=>{const r=inspectCommodityIntegrity({...base,evidence:[{id:'a',claim:'Updated beneficiary bank account - send to new IBAN',status:'CLAIMED',sourceIds:['email'],independentSourceGroups:1,observedAt:'x'}]});expect(r.pass).toBe(false)});
 it('rejects impossible purity',()=>{const r=inspectCommodityIntegrity({...base,evidence:[{id:'a',claim:'Gold 105% purity',status:'CLAIMED',sourceIds:['x'],independentSourceGroups:1,observedAt:'x'}]});expect(r.findings.some(x=>x.code==='IMPOSSIBLE_PURITY')).toBe(true)});
});
