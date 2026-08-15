import fs from 'node:fs';
import {spawnSync} from 'node:child_process';
const evidence={schema:'atlas.qualification-remediation.v1',at:new Date().toISOString(),attempts:[],status:'ENVIRONMENT_BLOCKED'};
const run=(name,cmd,args,ms)=>{const r=spawnSync(cmd,args,{encoding:'utf8',timeout:ms,env:{...process.env,npm_config_audit:'false',npm_config_fund:'false',npm_config_ignore_scripts:'true'}});evidence.attempts.push({name,status:r.status,signal:r.signal,error:r.error?.code||null,stderr:(r.stderr||'').slice(-2000)});return r.status===0};
if(!fs.existsSync('package-lock.json')){
  const ok=run('generate-lock','npm',['install','--package-lock-only','--ignore-scripts','--no-audit','--no-fund'],120000);
  if(!ok||!fs.existsSync('package-lock.json')){fs.mkdirSync('qualification',{recursive:true});fs.writeFileSync('qualification/remediation.json',JSON.stringify(evidence,null,2));console.error('ENVIRONMENT_BLOCKED: lockfile resolution unavailable');process.exit(42)}
}
if(!run('npm-ci','npm',['ci','--ignore-scripts','--no-audit','--no-fund'],180000)){fs.mkdirSync('qualification',{recursive:true});fs.writeFileSync('qualification/remediation.json',JSON.stringify(evidence,null,2));process.exit(42)}
const q=run('qualification','npm',['run','qualify'],180000); evidence.status=q?'QUALIFIED':'FAILED';fs.mkdirSync('qualification',{recursive:true});fs.writeFileSync('qualification/remediation.json',JSON.stringify(evidence,null,2));process.exit(q?0:1);
