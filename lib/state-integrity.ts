import type { AtlasDeal, DealEvent } from './types';
import { closingGate } from './closing-engine';

const IMMUTABLE = new Set(['id','organizationId','workspaceId','createdBy','createdAt']);
const HIGH_RISK = new Set(['beneficiary','beneficiaryAccount','bankAccount','paymentInstructions','indicativeValue','currency','quantity']);

export interface IntegrityFinding { code:string; severity:'HIGH'|'CRITICAL'; message:string; }
export interface IntegrityReport { pass:boolean; findings:IntegrityFinding[]; }

export function validateDealMutation(current:AtlasDeal,event:DealEvent):IntegrityReport {
  const findings:IntegrityFinding[]=[];
  if(event.dealId!==current.id) findings.push({code:'EVENT_DEAL_MISMATCH',severity:'CRITICAL',message:'Event targets a different deal.'});
  for(const key of Object.keys(event.payload)) {
    if(IMMUTABLE.has(key)) findings.push({code:'IMMUTABLE_FIELD_MUTATION',severity:'CRITICAL',message:`${key} cannot be changed by a deal event.`});
  }
  const org=(event.payload as any).organizationId, ws=(event.payload as any).workspaceId;
  if(org && org!==current.organizationId) findings.push({code:'TENANT_ESCALATION',severity:'CRITICAL',message:'Cross-organization mutation denied.'});
  if(ws && ws!==current.workspaceId) findings.push({code:'WORKSPACE_ESCALATION',severity:'CRITICAL',message:'Cross-workspace mutation denied.'});
  return {pass:findings.length===0,findings};
}

export function requireApprovalForMutation(event:DealEvent,approvedScopes:string[]=[]):IntegrityReport {
  const findings:IntegrityFinding[]=[];
  for(const key of Object.keys(event.payload)) {
    if(HIGH_RISK.has(key) && !approvedScopes.includes(`deal.change.${key}`)) findings.push({code:'APPROVAL_REQUIRED',severity:'CRITICAL',message:`Changing ${key} requires explicit scoped approval.`});
  }
  return {pass:findings.length===0,findings};
}

export function validateStateTransition(before:AtlasDeal,after:AtlasDeal):IntegrityReport {
  const findings:IntegrityFinding[]=[];
  for(const key of IMMUTABLE) if((before as any)[key] !== (after as any)[key]) findings.push({code:'STATE_CORRUPTION',severity:'CRITICAL',message:`Immutable field ${key} changed.`});
  if(after.quantity!==undefined && (!Number.isFinite(after.quantity)||after.quantity<=0)) findings.push({code:'INVALID_QUANTITY',severity:'CRITICAL',message:'Quantity must be finite and positive.'});
  if(after.indicativeValue!==undefined && (!Number.isFinite(after.indicativeValue)||after.indicativeValue<0)) findings.push({code:'INVALID_VALUE',severity:'CRITICAL',message:'Indicative value must be finite and non-negative.'});
  const ids=new Set(after.requirements.map(r=>r.id));
  for(const r of after.requirements) for(const dep of r.dependsOn) if(!ids.has(dep)) findings.push({code:'BROKEN_DEPENDENCY',severity:'HIGH',message:`Requirement ${r.id} depends on missing ${dep}.`});
  return {pass:findings.length===0,findings};
}

export function cannotBypassClosingGate(deal:AtlasDeal,requestedAction:string){
  const gate=closingGate(deal);
  const consequential=/execute|settle|pay|release|bind|close|transfer/i.test(requestedAction);
  return {allowed:!consequential || gate.pass,gate};
}
