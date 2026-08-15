import fs from 'node:fs';
import crypto from 'node:crypto';
const fail=(m,c=1)=>{console.error(m);process.exit(c)};
const sha=p=>crypto.createHash('sha256').update(fs.readFileSync(p)).digest('hex');
const qPath='PRODUCTION-QUALIFICATION.json';
if(!fs.existsSync(qPath)) fail('Missing PRODUCTION-QUALIFICATION.json');
const q=JSON.parse(fs.readFileSync(qPath,'utf8'));
const pkg=JSON.parse(fs.readFileSync('package.json','utf8'));
if(q.release!==pkg.atlasRelease) fail(`Qualification release mismatch: ${q.release||'missing'} != ${pkg.atlasRelease}`);
const required=['live_integrations','regulatory_data','observability','load_chaos','security_assessment','disaster_recovery','red_team','end_to_end_trade'];
const gates=q.gates||{}; const names=Object.keys(gates);
for(const name of required) if(!Object.hasOwn(gates,name)) fail(`Missing required production gate: ${name}`);
for(const name of names) if(!required.includes(name)) fail(`Unknown production gate: ${name}`);
if(!names.length) fail('No production gates defined');
const open=[];
for(const name of names){
  const g=gates[name];
  if(!g || typeof g!=='object' || Array.isArray(g)) fail(`Invalid gate object: ${name}`);
  if(!['pending','passed','failed','blocked'].includes(g.status)) fail(`Invalid gate status for ${name}: ${g.status}`);
  if(g.release!==q.release) fail(`Gate ${name} is bound to ${g.release||'no release'}, expected ${q.release}`);
  if(g.status!=='passed'){open.push(name);continue;}
  for(const k of ['evidence_path','evidence_sha256','executed_at','environment','authority']) if(!g[k]) fail(`Passed gate ${name} missing ${k}`);
  if(!/^([a-f0-9]{64})$/.test(g.evidence_sha256)) fail(`Passed gate ${name} has invalid evidence_sha256`);
  if(!fs.existsSync(g.evidence_path)) fail(`Passed gate ${name} evidence missing: ${g.evidence_path}`);
  if(sha(g.evidence_path)!==g.evidence_sha256) fail(`Passed gate ${name} evidence digest mismatch`);
  const executed=Date.parse(g.executed_at); if(!Number.isFinite(executed)) fail(`Passed gate ${name} invalid executed_at`);
  if(g.expires_at){const exp=Date.parse(g.expires_at); if(!Number.isFinite(exp)||exp<=Date.now()) fail(`Passed gate ${name} evidence expired/invalid`);}
}
if(open.length){console.error(`PRODUCTION BLOCKED: ${open.join(', ')}`);process.exit(2)}
console.log(`PRODUCTION GATES PASS: ${names.length} evidence-backed gates for ${q.release}`);
