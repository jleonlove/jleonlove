import type { AtlasDeal } from './types';
export type IntegrityFinding={code:string;severity:'MEDIUM'|'HIGH'|'CRITICAL';message:string;evidenceIds:string[]};
const num=(s:string,re:RegExp)=>{const m=s.match(re);return m?Number(m[1].replace(/,/g,'')):undefined};
export function inspectCommodityIntegrity(deal:AtlasDeal){
 const findings:IntegrityFinding[]=[];
 const claims=deal.evidence.map(e=>({e,text:e.claim.toLowerCase()}));
 const quantities=claims.map(x=>({id:x.e.id,n:num(x.text,/(\d[\d,.]*)\s*(?:kg|kilograms?|mt|metric\s+tons?)/i)})).filter(x=>x.n!==undefined) as {id:string;n:number}[];
 if(quantities.length>1){const vals=quantities.map(x=>x.n);const hi=Math.max(...vals),lo=Math.min(...vals);if(lo>0&&hi/lo>=1.5)findings.push({code:'MATERIAL_QUANTITY_CONFLICT',severity:'CRITICAL',message:`Material quantity claims conflict (${lo} vs ${hi}).`,evidenceIds:quantities.map(x=>x.id)});}
 const purities=claims.map(x=>({id:x.e.id,n:num(x.text,/(\d{1,3}(?:\.\d+)?)\s*%\s*(?:purity|fine|fineness|au|gold)/i)})).filter(x=>x.n!==undefined) as {id:string;n:number}[];
 if(purities.some(x=>x.n<=0||x.n>100))findings.push({code:'IMPOSSIBLE_PURITY',severity:'CRITICAL',message:'Purity/fineness claim falls outside a physically valid percentage range.',evidenceIds:purities.filter(x=>x.n<=0||x.n>100).map(x=>x.id)});
 if(purities.length>1){const hi=Math.max(...purities.map(x=>x.n)),lo=Math.min(...purities.map(x=>x.n));if(hi-lo>=2)findings.push({code:'MATERIAL_PURITY_CONFLICT',severity:'HIGH',message:`Material purity claims conflict (${lo}% vs ${hi}%).`,evidenceIds:purities.map(x=>x.id)});}
 const bankChange=claims.filter(x=>/beneficiary|bank account|swift|iban/.test(x.text)&&/change|replace|new|updated/.test(x.text));
 if(bankChange.length)findings.push({code:'BANKING_CHANGE_REVERIFY',severity:'CRITICAL',message:'Banking/beneficiary change requires independent re-verification and explicit approval before funds movement.',evidenceIds:bankChange.map(x=>x.e.id)});
 const impossible=claims.filter(x=>/no inspection required|no assay required|skip compliance|skip kyc|pay before verification/.test(x.text));
 if(impossible.length)findings.push({code:'DILIGENCE_BYPASS_ATTEMPT',severity:'CRITICAL',message:'Claim attempts to bypass normal transaction diligence.',evidenceIds:impossible.map(x=>x.e.id)});
 return {pass:!findings.some(f=>f.severity==='CRITICAL'),findings};
}
