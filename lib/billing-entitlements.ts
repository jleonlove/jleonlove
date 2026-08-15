import { createHash } from 'node:crypto';

export type InvoiceState='OPEN'|'PAYMENT_DETECTED'|'CONFIRMING'|'PAID'|'EXPIRED'|'REFUNDED'|'DISPUTED';
export type Invoice={invoiceId:string;customerId:string;tier:string;amountMinor:number;currency:string;state:InvoiceState;expiresAt:string;revision:number};
export type PaymentEvidence={paymentId:string;invoiceId:string;amountMinor:number;currency:string;final:boolean;providerVerified:boolean};
export type Entitlement={invoiceId:string;customerId:string;tier:string;active:boolean;evidenceDigest:string};

export function canonicalInvoiceDigest(i:Invoice){
 return createHash('sha256').update(JSON.stringify([i.invoiceId,i.customerId,i.tier,i.amountMinor,i.currency,i.expiresAt,i.revision])).digest('hex');
}

export function qualifyPayment(i:Invoice,p:PaymentEvidence,now=new Date()){
 const reasons:string[]=[];
 if(i.state==='PAID'||i.state==='REFUNDED') reasons.push('INVOICE_NOT_PAYABLE');
 if(new Date(i.expiresAt).getTime()<=now.getTime()) reasons.push('INVOICE_EXPIRED');
 if(p.invoiceId!==i.invoiceId) reasons.push('INVOICE_MISMATCH');
 if(p.amountMinor!==i.amountMinor) reasons.push(p.amountMinor<i.amountMinor?'UNDERPAYMENT':'OVERPAYMENT_REVIEW_REQUIRED');
 if(p.currency!==i.currency) reasons.push('CURRENCY_MISMATCH');
 if(!p.providerVerified) reasons.push('PAYMENT_EVIDENCE_UNVERIFIED');
 if(!p.final) reasons.push('PAYMENT_NOT_FINAL');
 return {qualified:reasons.length===0,reasons};
}

export function activateEntitlement(i:Invoice,p:PaymentEvidence,now=new Date()):Entitlement{
 const q=qualifyPayment(i,p,now); if(!q.qualified) throw new Error(q.reasons.join(','));
 return {invoiceId:i.invoiceId,customerId:i.customerId,tier:i.tier,active:true,evidenceDigest:canonicalInvoiceDigest(i)};
}

export function entitlementAllows(e:Entitlement,requestedTier:string){return e.active&&e.tier===requestedTier;}
