import { createHash } from "crypto";

export type CustodyRequest = {
  tenantId:string; workspaceId:string; action:"SIGN"|"WITHDRAW"|"REFUND"|"CONVERT";
  asset:string; network:string; destination:string; amountAtomic:bigint;
  approvalIds:string[]; nonce:string; expiresAt:string;
};

export function custodyRequestDigest(r:CustodyRequest){
  return createHash("sha256").update(JSON.stringify({...r,amountAtomic:r.amountAtomic.toString(),approvalIds:[...r.approvalIds].sort()})).digest("hex");
}

// Atlas is a merchant for its own services, never a custodian or transmitter of customer funds.
// Generic customer custody requests are therefore structurally prohibited, regardless of approvals.
export function qualifyCustodyRequest(r:CustodyRequest, now=new Date()){
  const reasons:string[]=["CUSTOMER_CUSTODY_PROHIBITED"];
  if(!r.tenantId||!r.workspaceId) reasons.push("MISSING_SCOPE");
  if(!r.asset||!r.network||!r.destination||r.amountAtomic<=0n) reasons.push("INVALID_TRANSFER");
  if(!r.nonce) reasons.push("MISSING_NONCE");
  if(new Date(r.expiresAt).getTime()<=now.getTime()) reasons.push("REQUEST_EXPIRED");
  return {allowed:false,reasons,digest:custodyRequestDigest(r)};
}

export function assertNoPrivateKeyMaterial(value:unknown){
  const text=JSON.stringify(value,(_,v)=>typeof v==="bigint"?v.toString():v).toLowerCase();
  const forbidden=["privatekey","private_key","seedphrase","seed_phrase","mnemonic","secretkey","secret_key"];
  const hit=forbidden.find(k=>text.includes(k));
  return {safe:!hit,reason:hit?"PRIVATE_KEY_MATERIAL_FORBIDDEN" as const:"CLEAR" as const};
}
