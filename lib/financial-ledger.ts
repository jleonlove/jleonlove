export type AccountType = "ASSET"|"LIABILITY"|"EQUITY"|"REVENUE"|"EXPENSE";
export type JournalLine = { accountId:string; accountType:AccountType; debitMinor:bigint; creditMinor:bigint };
export type JournalEntry = { id:string; organizationId:string; workspaceId:string; currency:"USD"; effectiveAt:string; sourceEventId:string; lines:JournalLine[]; reversalOf?:string; posted:boolean };
export function validateJournal(e:JournalEntry){
 const errors:string[]=[]; if(!e.id||!e.organizationId||!e.workspaceId||!e.sourceEventId) errors.push("MISSING_SCOPE_OR_ID");
 if(e.currency!=="USD") errors.push("UNSUPPORTED_LEDGER_CURRENCY"); if(e.lines.length<2) errors.push("INSUFFICIENT_LINES");
 let d=0n,c=0n; for(const l of e.lines){ if(l.debitMinor<0n||l.creditMinor<0n) errors.push("NEGATIVE_LINE"); if(l.debitMinor>0n&&l.creditMinor>0n) errors.push("BOTH_DEBIT_AND_CREDIT"); d+=l.debitMinor;c+=l.creditMinor; }
 if(d!==c) errors.push("UNBALANCED_JOURNAL"); if(d===0n) errors.push("ZERO_VALUE_JOURNAL"); return {ok:errors.length===0,errors,debitsMinor:d,creditsMinor:c};
}
export function mutatePostedEntry(e:JournalEntry){ if(e.posted) throw new Error("POSTED_ENTRY_IMMUTABLE_USE_REVERSAL"); return e; }
export function reverseEntry(original:JournalEntry,id:string,sourceEventId:string):JournalEntry { if(!original.posted) throw new Error("ONLY_POSTED_ENTRY_CAN_REVERSE"); return {...original,id,sourceEventId,reversalOf:original.id,posted:false,lines:original.lines.map(l=>({...l,debitMinor:l.creditMinor,creditMinor:l.debitMinor}))}; }
export type BankEvent={eventId:string; paymentId:string; invoiceId:string; organizationId:string; workspaceId:string; rail:"ACH"|"WIRE"|"CARD"; kind:"INITIATED"|"SETTLED"|"RETURNED"|"REVERSED"; amountMinor:bigint; providerReference:string};
export function dedupeBankEvents(events:BankEvent[]){const seen=new Set<string>();const duplicates:string[]=[];for(const e of events){if(seen.has(e.eventId)) duplicates.push(e.eventId);seen.add(e.eventId);}return {ok:duplicates.length===0,duplicates};}
export function canRecognizeCashRevenue(e:BankEvent){return e.kind==="SETTLED";}
export function paymentAllocation(expectedMinor:bigint, allocations:{invoiceId:string;amountMinor:bigint}[]){const total=allocations.reduce((n,a)=>n+a.amountMinor,0n);return {ok:total===expectedMinor,totalMinor:total,deltaMinor:total-expectedMinor};}
