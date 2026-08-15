export interface RateSample{ actorId:string; action:string; timestamp:number; }
export interface RateDecision{ allow:boolean; retryAfterMs:number; reason?:string; }
export function rateDecision(samples:RateSample[],actorId:string,action:string,now:number,limit:number,windowMs:number):RateDecision{
 if(limit<1||windowMs<1) return {allow:false,retryAfterMs:windowMs,reason:'INVALID_RATE_POLICY'};
 const recent=samples.filter(x=>x.actorId===actorId&&x.action===action&&x.timestamp>now-windowMs).sort((a,b)=>a.timestamp-b.timestamp);
 if(recent.length<limit) return {allow:true,retryAfterMs:0};
 return {allow:false,retryAfterMs:Math.max(1,recent[0].timestamp+windowMs-now),reason:'RATE_LIMIT_EXCEEDED'};
}
export function detectReplay(nonces:string[]){ const seen=new Set<string>(); const duplicates:string[]=[]; for(const n of nonces){if(seen.has(n))duplicates.push(n);seen.add(n);} return {safe:duplicates.length===0,duplicates}; }
