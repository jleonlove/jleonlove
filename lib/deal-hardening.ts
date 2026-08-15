import type { AtlasDeal, DealEvidence, DealRequirement, EvidenceStatus } from './types';

export interface AuthorityEdge { id:string; fromParticipantId:string; toParticipantId:string; scope:string[]; status:EvidenceStatus; evidenceIds:string[]; effectiveAt?:string; expiresAt?:string; }
export interface PhysicalClaim { id:string; metric:string; value:number; unit:string; sourceEvidenceIds:string[]; }
export interface RealityFinding { code:string; severity:'LOW'|'MEDIUM'|'HIGH'|'CRITICAL'; message:string; evidenceIds:string[]; }
export interface CriticalPathResult { ordered: DealRequirement[]; blockers: DealRequirement[]; unresolvedDependencies:string[]; }

export function authorityChain(deal:AtlasDeal){
 const unverified=deal.participants.filter(p=>p.authority!=='INDEPENDENTLY_VERIFIED').map(p=>p.id);
 return {pass:unverified.length===0,unverified};
}

export function verifyAuthorityChain(deal:AtlasDeal, edges:AuthorityEdge[], now=new Date()) {
 const participantIds=new Set(deal.participants.map(p=>p.id)); const findings:string[]=[];
 for(const e of edges){
  if(!participantIds.has(e.fromParticipantId)||!participantIds.has(e.toParticipantId)) findings.push(`Authority edge ${e.id} references an unknown participant.`);
  if(e.status!=='INDEPENDENTLY_VERIFIED') findings.push(`Authority edge ${e.id} is ${e.status}.`);
  if(e.expiresAt && new Date(e.expiresAt)<=now) findings.push(`Authority edge ${e.id} is expired.`);
  if(!e.scope.length) findings.push(`Authority edge ${e.id} has no defined scope.`);
 }
 return { valid: findings.length===0, findings };
}

export function evaluatePhysicalReality(claims:PhysicalClaim[], evidence:DealEvidence[]):RealityFinding[]{
 const out:RealityFinding[]=[]; const groups=new Map<string,PhysicalClaim[]>();
 for(const c of claims){const k=`${c.metric}:${c.unit}`;groups.set(k,[...(groups.get(k)||[]),c]);}
 for(const [key,items] of groups){
  if(items.length<2) continue; const values=items.map(i=>i.value); const min=Math.min(...values),max=Math.max(...values);
  if(min>0 && max/min>=1.5) out.push({code:'PHYSICAL_CLAIM_CONFLICT',severity:max/min>=2?'CRITICAL':'HIGH',message:`Material conflict in ${key}: ${min} vs ${max}.`,evidenceIds:[...new Set(items.flatMap(i=>i.sourceEvidenceIds))]});
 }
 for(const c of claims){const support=c.sourceEvidenceIds.map(id=>evidence.find(e=>e.id===id)).filter(Boolean) as DealEvidence[]; if(!support.length) out.push({code:'UNSUPPORTED_PHYSICAL_CLAIM',severity:'HIGH',message:`${c.metric} has no linked evidence.`,evidenceIds:[]});}
 return out;
}

export function computeCriticalPath(deal:AtlasDeal):CriticalPathResult{
 const byId=new Map(deal.requirements.map(r=>[r.id,r])); const unresolved=new Set<string>(); const blockers=deal.requirements.filter(r=>r.status==='BLOCKED'||(r.severity==='CRITICAL'&&r.status!=='SATISFIED'));
 const score=(r:DealRequirement,seen=new Set<string>()):number=>{if(seen.has(r.id))return 1000; seen.add(r.id); let s=({LOW:1,MEDIUM:3,HIGH:7,CRITICAL:15} as const)[r.severity]; for(const id of r.dependsOn){const d=byId.get(id); if(!d){unresolved.add(id);s+=50;} else if(d.status!=='SATISFIED')s+=score(d,new Set(seen));} return s;};
 const ordered=deal.requirements.filter(r=>!['SATISFIED','WAIVED'].includes(r.status)).sort((a,b)=>score(b)-score(a)); return {ordered,blockers,unresolvedDependencies:[...unresolved]};
}

export function invalidateAssumptions(deal:AtlasDeal, changedKeys:string[]):AtlasDeal['assumptions']{
 const changed=new Set(changedKeys); return deal.assumptions.map(a=>a.dependencyKeys.some(k=>changed.has(k))?{...a,valid:false}:a);
}
