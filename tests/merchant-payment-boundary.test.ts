import {describe,it,expect} from 'vitest';
import {enforceMerchantPaymentBoundary} from '../lib/merchant-payment-boundary';
const base={invoiceId:'inv1',customerId:'c1',atlasServiceLevel:'PRO',purpose:'ATLAS_SERVICE_INVOICE' as const,asset:'BTC' as const,network:'bitcoin-mainnet',atlasOwnedDestination:true,destination:'bc1atlas'};
describe('Atlas merchant-only payment boundary',()=>{
 it('allows payment owed directly to Atlas for a service',()=>expect(enforceMerchantPaymentBoundary(base).allowed).toBe(true));
 it.each(['CUSTOMER_ESCROW','CUSTOMER_TRANSFER','CUSTOMER_CUSTODY','CUSTOMER_EXCHANGE'] as const)('prohibits %s',purpose=>expect(enforceMerchantPaymentBoundary({...base,purpose}).allowed).toBe(false));
 it('requires Atlas-owned destination',()=>expect(enforceMerchantPaymentBoundary({...base,atlasOwnedDestination:false}).reasons).toContain('DESTINATION_MUST_BE_ATLAS_OWNED'));
 it('requires XRP tag and XLM memo by policy',()=>{expect(enforceMerchantPaymentBoundary({...base,asset:'XRP',network:'xrpl-mainnet'}).allowed).toBe(false);expect(enforceMerchantPaymentBoundary({...base,asset:'XLM',network:'stellar-public'}).allowed).toBe(false)});
});
