export type TaskState="OPEN"|"DONE"|"BLOCKED";
export type DealTask={id:string;label:string;owner:string;dependsOn:string[];state:TaskState;estimatedHours?:number;unlocks?:string[];hardBlocker?:boolean;reason?:string};
export type ClosingPath={readiness:number;blockers:string[];parallelNow:string[];criticalPath:string[];nextAction?:string;deadlocks:string[][];explanations:Record<string,string>};

function detectCycles(tasks:DealTask[]):string[][]{
 const map=new Map(tasks.map(t=>[t.id,t])); const out:string[][]=[]; const seen=new Set<string>(); const stack:string[]=[]; const active=new Set<string>();
 const dfs=(id:string)=>{if(active.has(id)){const i=stack.indexOf(id); if(i>=0) out.push([...stack.slice(i),id]); return;} if(seen.has(id))return; seen.add(id);active.add(id);stack.push(id); for(const d of map.get(id)?.dependsOn??[]) if(map.has(d)) dfs(d); stack.pop();active.delete(id)};
 for(const t of tasks) dfs(t.id); return out;
}
export function fastestLegitimatePathToClose(tasks:DealTask[]):ClosingPath{
 const map=new Map(tasks.map(t=>[t.id,t])); const done=new Set(tasks.filter(t=>t.state==="DONE").map(t=>t.id));
 const actionable=tasks.filter(t=>t.state==="OPEN"&&t.dependsOn.every(d=>done.has(d)));
 const score=(t:DealTask)=>(t.hardBlocker?1000:0)+(t.unlocks?.length??0)*100-(t.estimatedHours??24);
 const ordered=[...actionable].sort((a,b)=>score(b)-score(a)||a.id.localeCompare(b.id));
 const blockers=tasks.filter(t=>t.state==="BLOCKED"||t.hardBlocker&&t.state!=="DONE").map(t=>t.id);
 const remaining=tasks.filter(t=>t.state!=="DONE"); const readiness=tasks.length?Math.round(((tasks.length-remaining.length)/tasks.length)*100):100;
 const criticalPath:string[]=[]; const pending=new Set(remaining.map(t=>t.id)); let guard=0;
 while(pending.size&&guard++<tasks.length*2){const candidates=[...pending].map(id=>map.get(id)!).filter(t=>t.dependsOn.every(d=>done.has(d)||!pending.has(d))).sort((a,b)=>score(b)-score(a)); if(!candidates.length)break; const t=candidates[0];criticalPath.push(t.id);pending.delete(t.id);done.add(t.id)}
 const explanations=Object.fromEntries(ordered.map(t=>[t.id,`${t.label} — owner: ${t.owner}; ${t.reason??"required for closing"}; unlocks ${(t.unlocks??[]).length} downstream task(s).`]));
 return {readiness,blockers,parallelNow:ordered.map(t=>t.id),criticalPath,nextAction:ordered[0]?.id,deadlocks:detectCycles(tasks),explanations};
}
