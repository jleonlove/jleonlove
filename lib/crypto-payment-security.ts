import type { CryptoAsset, CryptoPaymentIntent, PaymentState } from './crypto-payments';

export interface ChainObservation { asset:CryptoAsset; network:string; txHash:string; destination:string; memoOrTag?:string; receivedAtomic:string; confirmations:number; finalized:boolean; observedAt:string; }
export interface PaymentSecurityContext { usedTxHashes:Set<string>; minConfirmations:Partial<Record<CryptoAsset,number>>; now?:Date; }
export interface PaymentDecision { state:PaymentState; accept:boolean; reasons:string[]; }

export function verifyChainPayment(intent:CryptoPaymentIntent, obs:ChainObservation, ctx:PaymentSecurityContext):PaymentDecision {
  const reasons:string[]=[]; const now=ctx.now ?? new Date();
  if (new Date(intent.quoteExpiresAt).getTime() <= now.getTime()) reasons.push('QUOTE_EXPIRED');
  if (obs.asset !== intent.asset) reasons.push('WRONG_ASSET');
  if (obs.network !== intent.network) reasons.push('WRONG_NETWORK');
  if (obs.destination !== intent.destination) reasons.push('WRONG_DESTINATION');
  if ((intent.asset==='XRP'||intent.asset==='XLM') && intent.memoOrTag !== obs.memoOrTag) reasons.push('MEMO_OR_TAG_MISMATCH');
  if (!/^\d+$/.test(obs.receivedAtomic)) reasons.push('INVALID_OBSERVED_AMOUNT');
  if (ctx.usedTxHashes.has(obs.txHash)) reasons.push('TX_REPLAY');
  const minimum=ctx.minConfirmations[intent.asset] ?? 1;
  if (obs.confirmations < minimum || !obs.finalized) reasons.push('NOT_FINAL');
  if (reasons.length) return {state: reasons.includes('QUOTE_EXPIRED')?'EXPIRED':'REVIEW_REQUIRED',accept:false,reasons};
  const got=BigInt(obs.receivedAtomic), want=BigInt(intent.expectedAtomic);
  if(got<want)return {state:'UNDERPAID',accept:false,reasons:['UNDERPAID']};
  if(got>want)return {state:'OVERPAID',accept:false,reasons:['OVERPAID_REVIEW_REQUIRED']};
  return {state:'SETTLED',accept:true,reasons:[]};
}

export type TreasuryAction='REFUND'|'WITHDRAW'|'CONVERT';
export interface TreasuryAuthorization { action:TreasuryAction; organizationId:string; workspaceId:string; asset:CryptoAsset; amountAtomic:string; destination?:string; approvalIds:string[]; requiredApprovals:number; expiresAt:string; nonce:string; }
export function authorizeTreasuryAction(a:TreasuryAuthorization, usedNonces:Set<string>, now=new Date()){
 const reasons:string[]=[];
 if(!/^\d+$/.test(a.amountAtomic)||BigInt(a.amountAtomic)<=0n) reasons.push('INVALID_AMOUNT');
 if(new Date(a.expiresAt).getTime()<=now.getTime()) reasons.push('APPROVAL_EXPIRED');
 if(new Set(a.approvalIds).size<a.requiredApprovals) reasons.push('INSUFFICIENT_DISTINCT_APPROVALS');
 if(usedNonces.has(a.nonce)) reasons.push('NONCE_REPLAY');
 if((a.action==='WITHDRAW'||a.action==='REFUND')&&!a.destination?.trim()) reasons.push('DESTINATION_REQUIRED');
 return {authorized:reasons.length===0,reasons};
}
