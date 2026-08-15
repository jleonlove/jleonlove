import type { AtlasDeal, DealEvent } from './types';
import { stressClosing, type ClosingStressReport } from './closing-intelligence';

export type EventDisposition = 'ACCEPT'|'DUPLICATE'|'LATE'|'QUARANTINE'|'RECONCILE';
export interface AuthoritativeDealEvent extends DealEvent {
  organizationId:string; workspaceId:string; source:string; sourceAuthority:number;
  receivedAt:string; effectiveAt:string; baseVersion:number; evidenceDigest?:string;
}
export interface EventDecision { eventId:string; disposition:EventDisposition; reason:string; }
export interface DealControlState { deal:AtlasDeal; version:number; seen:Set<string>; ledger:AuthoritativeDealEvent[]; quarantined:AuthoritativeDealEvent[]; reconciliation:AuthoritativeDealEvent[]; closing:ClosingStressReport; }

const MATERIAL = new Set(['TERMS_CHANGED','COUNTERPARTY_CHANGED','AUTHORITY_CHANGED','EVIDENCE_CHANGED','DOCUMENT_CHANGED','INSPECTION_CHANGED','LOGISTICS_CHANGED','FINANCING_CHANGED','COMPLIANCE_CHANGED','SETTLEMENT_CHANGED','APPROVAL_CHANGED']);
const CONFLICT_FAMILIES:Record<string,string>={SETTLEMENT_CONFIRMED:'settlement',SETTLEMENT_REJECTED:'settlement',TITLE_CONFIRMED:'title',TITLE_REJECTED:'title',INSPECTION_PASSED:'inspection',INSPECTION_FAILED:'inspection'};

function isoMs(v:string){ const n=Date.parse(v); return Number.isFinite(n)?n:-1; }
function stableEventOrder(a:AuthoritativeDealEvent,b:AuthoritativeDealEvent){return isoMs(a.effectiveAt)-isoMs(b.effectiveAt)||isoMs(a.occurredAt)-isoMs(b.occurredAt)||a.id.localeCompare(b.id);}
function cloneDeal(d:AtlasDeal):AtlasDeal{return JSON.parse(JSON.stringify(d)) as AtlasDeal;}
function materialMutation(deal:AtlasDeal,e:AuthoritativeDealEvent):AtlasDeal{
 const next=cloneDeal(deal); next.updatedAt=e.effectiveAt;
 if(e.type==='TERMS_CHANGED'){
  for(const k of ['commodity','quantity','unit','currency','indicativeValue','stage'] as const){ if(k in e.payload) (next as any)[k]=e.payload[k]; }
 }
 if(e.type==='EVIDENCE_CHANGED' && e.payload.evidence && typeof e.payload.evidence==='object'){
  const ev=e.payload.evidence as any; next.evidence=[...next.evidence.filter(x=>x.id!==ev.id),ev];
 }
 if(e.type==='APPROVAL_CHANGED' && typeof e.payload.requirementId==='string'){
  next.requirements=next.requirements.map(r=>r.id===e.payload.requirementId?{...r,status:String(e.payload.status||r.status) as any}:r);
 }
 return next;
}
function conflicting(state:DealControlState,e:AuthoritativeDealEvent){
 const family=CONFLICT_FAMILIES[e.type]; if(!family) return false;
 return state.ledger.some(x=>CONFLICT_FAMILIES[x.type]===family && x.type!==e.type && x.sourceAuthority>=e.sourceAuthority && Math.abs(isoMs(x.effectiveAt)-isoMs(e.effectiveAt))<86400000);
}
export function createControlState(deal:AtlasDeal,version=1):DealControlState{return {deal:cloneDeal(deal),version,seen:new Set(),ledger:[],quarantined:[],reconciliation:[],closing:stressClosing(deal)};}
export function ingestDealEvent(state:DealControlState,e:AuthoritativeDealEvent):EventDecision{
 if(state.seen.has(e.idempotencyKey)||state.ledger.some(x=>x.id===e.id)){return {eventId:e.id,disposition:'DUPLICATE',reason:'Event already processed.'};}
 if(e.dealId!==state.deal.id||e.organizationId!==state.deal.organizationId||e.workspaceId!==state.deal.workspaceId||!e.source||e.sourceAuthority<0||e.sourceAuthority>100||isoMs(e.effectiveAt)<0){state.quarantined.push(e);return {eventId:e.id,disposition:'QUARANTINE',reason:'Invalid provenance, scope, authority, or time.'};}
 if(e.baseVersion>state.version){state.quarantined.push(e);return {eventId:e.id,disposition:'QUARANTINE',reason:'Event references a future transaction version.'};}
 if(conflicting(state,e)){state.reconciliation.push(e);return {eventId:e.id,disposition:'RECONCILE',reason:'Conflicting authoritative material event requires reconciliation.'};}
 const late=e.baseVersion<state.version;
 state.seen.add(e.idempotencyKey); state.ledger.push(e); state.ledger.sort(stableEventOrder);
 if(MATERIAL.has(e.type) && !late){state.deal=materialMutation(state.deal,e);state.version++;state.closing=stressClosing(state.deal);}
 return {eventId:e.id,disposition:late?'LATE':'ACCEPT',reason:late?'Accepted as historical; current state was not blindly overwritten.':'Accepted and transaction truth recalculated.'};
}
export function replayDealEvents(seed:AtlasDeal,events:AuthoritativeDealEvent[],version=1):DealControlState{
 const s=createControlState(seed,version); for(const e of [...events].sort(stableEventOrder)) ingestDealEvent(s,e); return s;
}
export function explainReadinessChange(before:DealControlState,after:DealControlState){return {from:before.closing.score,to:after.closing.score,delta:after.closing.score-before.closing.score,versionFrom:before.version,versionTo:after.version,reasons:after.closing.reasons};}
