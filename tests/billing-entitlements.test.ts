import {describe,it,expect} from 'vitest';
import {activateEntitlement,entitlementAllows,qualifyPayment,type Invoice} from '../lib/billing-entitlements';
const inv:Invoice={invoiceId:'inv-138',customerId:'cust',tier:'PRO',amountMinor:10000,currency:'USD',state:'OPEN',expiresAt:'2099-01-01T00:00:00Z',revision:1};
const pay={paymentId:'p1',invoiceId:'inv-138',amountMinor:10000,currency:'USD',final:true,providerVerified:true};
describe('billing entitlement integrity',()=>{
 it('activates exactly the purchased tier only after verified final payment',()=>{const e=activateEntitlement(inv,pay);expect(entitlementAllows(e,'PRO')).toBe(true);expect(entitlementAllows(e,'ENTERPRISE')).toBe(false)});
 it('rejects underpayment',()=>expect(qualifyPayment(inv,{...pay,amountMinor:9999}).reasons).toContain('UNDERPAYMENT'));
 it('routes overpayment to review',()=>expect(qualifyPayment(inv,{...pay,amountMinor:10001}).reasons).toContain('OVERPAYMENT_REVIEW_REQUIRED'));
 it('rejects unverified/non-final evidence',()=>expect(qualifyPayment(inv,{...pay,providerVerified:false,final:false}).qualified).toBe(false));
 it('rejects replay against paid invoice',()=>expect(qualifyPayment({...inv,state:'PAID'},pay).reasons).toContain('INVOICE_NOT_PAYABLE'));
 it('rejects invoice mismatch',()=>expect(qualifyPayment(inv,{...pay,invoiceId:'other'}).reasons).toContain('INVOICE_MISMATCH'));
});
