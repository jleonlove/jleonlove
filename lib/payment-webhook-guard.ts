import { verifyWebhookSignature } from "./payment-provider-integrity";

export type WebhookEnvelope={provider:string;eventId:string;timestamp:string;rawBody:string;signatureHex:string};
export function qualifyWebhook(e:WebhookEnvelope, secret:string, seen:Set<string>, now=new Date(), maxSkewMs=5*60*1000){
 const reasons:string[]=[];
 if(!verifyWebhookSignature(e.rawBody,e.signatureHex,secret)) reasons.push("INVALID_SIGNATURE");
 if(!e.eventId || seen.has(`${e.provider}:${e.eventId}`)) reasons.push("DUPLICATE_EVENT");
 const t=new Date(e.timestamp).getTime();
 if(!Number.isFinite(t) || Math.abs(now.getTime()-t)>maxSkewMs) reasons.push("STALE_OR_INVALID_TIMESTAMP");
 return {accepted:reasons.length===0,reasons,key:`${e.provider}:${e.eventId}`};
}
export function commitWebhook(key:string,seen:Set<string>){ if(seen.has(key)) return false; seen.add(key); return true; }
