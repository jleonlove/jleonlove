import type { CryptoAsset } from './crypto-payments';

export type PaymentPurpose = 'ATLAS_SERVICE_INVOICE'|'CUSTOMER_ESCROW'|'CUSTOMER_TRANSFER'|'CUSTOMER_CUSTODY'|'CUSTOMER_EXCHANGE';
export type MerchantPaymentRequest = {
  invoiceId:string; customerId:string; atlasServiceLevel:string; purpose:PaymentPurpose;
  asset:'USD'|CryptoAsset; network?:string; atlasOwnedDestination:boolean;
  destination?:string; memoOrTag?:string;
};

export function enforceMerchantPaymentBoundary(r:MerchantPaymentRequest){
  const reasons:string[]=[];
  if(r.purpose!=='ATLAS_SERVICE_INVOICE') reasons.push('CUSTOMER_FUNDS_ACTIVITY_PROHIBITED');
  if(!r.invoiceId?.trim() || !r.customerId?.trim() || !r.atlasServiceLevel?.trim()) reasons.push('SERVICE_INVOICE_REQUIRED');
  if(!r.atlasOwnedDestination) reasons.push('DESTINATION_MUST_BE_ATLAS_OWNED');
  if(r.asset!=='USD' && !r.destination?.trim()) reasons.push('ATLAS_PAYMENT_ADDRESS_REQUIRED');
  if(r.asset==='XRP' && !r.memoOrTag?.trim()) reasons.push('XRP_DESTINATION_TAG_REQUIRED_BY_POLICY');
  if(r.asset==='XLM' && !r.memoOrTag?.trim()) reasons.push('XLM_MEMO_REQUIRED_BY_POLICY');
  return {allowed:reasons.length===0,reasons};
}

export const PROHIBITED_CUSTOMER_MONEY_CAPABILITIES = Object.freeze([
  'CUSTOMER_BALANCE','CUSTOMER_WITHDRAWAL','CUSTOMER_ESCROW','CUSTOMER_TRANSFER',
  'CUSTOMER_CRYPTO_CUSTODY','CUSTOMER_EXCHANGE','CUSTOMER_PAYMENT_FORWARDING'
] as const);
