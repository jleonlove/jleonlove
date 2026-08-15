import { createHash } from 'crypto';

export type SettlementRail='WIRE'|'ACH'|'BTC'|'XRP'|'ETH'|'HBAR'|'ADA'|'XLM';
export type ClosingInstruction={dealId:string;workspaceId:string;stateVersion:number;counterpartyId:string;counterpartyLegalName:string;mandateId:string;rail:SettlementRail;destination:string;beneficiaryName?:string;bankName?:string;intermediaryBank?:string;memoOrTag?:string;sourcePrincipalId:string;sourceChannel:string;sourceChannelVerified:boolean;authorityVerified:boolean;independentlyVerified:boolean;verifiedAt:string;expiresAt:string};
export type InstructionBinding={instructionDigest:string;dealId:string;workspaceId:string;stateVersion:number;counterpartyId:string;mandateId:string;rail:SettlementRail;createdAt:string};
const stable=(v:unknown)=>JSON.stringify(v,Object.keys(v as object).sort());
export const closingInstructionDigest=(i:ClosingInstruction)=>createHash('sha256').update(stable(i)).digest('hex');
export function bindClosingInstruction(i:ClosingInstruction,now=new Date()):InstructionBinding{
 const v=validateClosingInstruction(i,now); if(!v.pass) throw new Error(v.reasons.join(','));
 return {instructionDigest:closingInstructionDigest(i),dealId:i.dealId,workspaceId:i.workspaceId,stateVersion:i.stateVersion,counterpartyId:i.counterpartyId,mandateId:i.mandateId,rail:i.rail,createdAt:now.toISOString()};
}
export function validateClosingInstruction(i:ClosingInstruction,now=new Date()){
 const reasons:string[]=[];
 if(!i.dealId||!i.workspaceId||!i.counterpartyId||!i.mandateId)reasons.push('INSTRUCTION_IDENTITY_INCOMPLETE');
 if(!i.destination.trim())reasons.push('DESTINATION_REQUIRED');
 if(!i.sourceChannelVerified)reasons.push('UNVERIFIED_SOURCE_CHANNEL');
 if(!i.authorityVerified)reasons.push('UNVERIFIED_CHANGE_AUTHORITY');
 if(!i.independentlyVerified)reasons.push('INDEPENDENT_VERIFICATION_REQUIRED');
 if(new Date(i.expiresAt).getTime()<=now.getTime())reasons.push('INSTRUCTION_EXPIRED');
 if((i.rail==='XRP'||i.rail==='XLM')&&i.memoOrTag!==undefined&&!i.memoOrTag.trim())reasons.push('EMPTY_MEMO_OR_TAG');
 return {pass:reasons.length===0,reasons};
}
export function verifyInstructionAtExecution(binding:InstructionBinding,current:ClosingInstruction,now=new Date()){
 const base=validateClosingInstruction(current,now); const reasons=[...base.reasons];
 if(binding.dealId!==current.dealId)reasons.push('INSTRUCTION_DEAL_CHANGED');
 if(binding.workspaceId!==current.workspaceId)reasons.push('INSTRUCTION_WORKSPACE_CHANGED');
 if(binding.stateVersion!==current.stateVersion)reasons.push('INSTRUCTION_STATE_CHANGED');
 if(binding.counterpartyId!==current.counterpartyId)reasons.push('COUNTERPARTY_IDENTITY_CHANGED');
 if(binding.mandateId!==current.mandateId)reasons.push('MANDATE_CHANGED');
 if(binding.rail!==current.rail)reasons.push('SETTLEMENT_RAIL_CHANGED');
 if(binding.instructionDigest!==closingInstructionDigest(current))reasons.push('CLOSING_INSTRUCTION_SUBSTITUTED');
 return {pass:reasons.length===0,reasons};
}
export type SettlementEvidence={provider:string;dealId:string;workspaceId:string;instructionDigest:string;transactionRef:string;status:'PENDING'|'SETTLED'|'REVERSED';amount:string;currencyOrAsset:string;verified:boolean};
export function reconcileSettlementEvidence(items:SettlementEvidence[],binding:InstructionBinding){
 const reasons:string[]=[];
 if(items.length===0)return {pass:false,state:'RECONCILIATION_REQUIRED' as const,reasons:['MISSING_SETTLEMENT_EVIDENCE']};
 for(const e of items){if(!e.verified)reasons.push('UNVERIFIED_SETTLEMENT_EVIDENCE');if(e.dealId!==binding.dealId||e.workspaceId!==binding.workspaceId)reasons.push('SETTLEMENT_BOUNDARY_MISMATCH');if(e.instructionDigest!==binding.instructionDigest)reasons.push('SETTLEMENT_FOR_DIFFERENT_INSTRUCTION');}
 const statuses=new Set(items.map(x=>x.status)); const refs=new Set(items.map(x=>x.transactionRef)); const amounts=new Set(items.map(x=>`${x.amount}:${x.currencyOrAsset}`));
 if(statuses.size>1)reasons.push('CONFLICTING_SETTLEMENT_STATUS'); if(refs.size>1)reasons.push('CONFLICTING_SETTLEMENT_REFERENCE'); if(amounts.size>1)reasons.push('CONFLICTING_SETTLEMENT_AMOUNT');
 if(reasons.length)return {pass:false,state:'RECONCILIATION_REQUIRED' as const,reasons};
 return {pass:items[0].status==='SETTLED',state:items[0].status,reasons:items[0].status==='SETTLED'?[]:['SETTLEMENT_NOT_FINAL']};
}
