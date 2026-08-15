export type PaymentLedgerEntry = { invoiceId:string; txHash:string; asset:string; amountAtomic:bigint; fiatMinor:number; feeAtomic?:bigint; status:"PENDING"|"SETTLED"|"REVERSED"|"HOLD" };
export function reconcileInvoice(expectedAtomic: bigint, entries: PaymentLedgerEntry[]) {
  const settled = entries.filter(e=>e.status==="SETTLED").reduce((n,e)=>n+e.amountAtomic,0n);
  return { expectedAtomic, settledAtomic:settled, deltaAtomic:settled-expectedAtomic, balanced:settled===expectedAtomic };
}
export function complianceHold(reasons:string[]){ return { allowed: reasons.length===0, status: reasons.length ? "HOLD" as const : "CLEAR" as const, reasons }; }
