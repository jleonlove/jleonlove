import fs from 'node:fs';
import { spawnSync } from 'node:child_process';
import crypto from 'node:crypto';
const pkg=JSON.parse(fs.readFileSync('package.json','utf8'));
const RELEASE=pkg.atlasRelease;
if(!/^RC-\d{6}$/.test(RELEASE||'')){console.error('package.json atlasRelease missing or invalid');process.exit(1)}
const out=`QUALIFICATION-${RELEASE}.json`;
const ev={release:RELEASE,started_at:new Date().toISOString(),status:'ENVIRONMENT_BLOCKED',gates:[]};
const gate=(name,status,detail={})=>ev.gates.push({name,status,...detail});
const sha=p=>crypto.createHash('sha256').update(fs.readFileSync(p)).digest('hex');
const finish=(status=ev.status)=>{ev.status=status;ev.finished_at=new Date().toISOString();ev.source={package_json_sha256:sha('package.json'),lockfile:ev.lockfile||null};fs.writeFileSync(out,JSON.stringify(ev,null,2)+'\n');console.log(status);process.exit(status==='QUALIFIED'?0:status==='FAILED'?1:2)};
const expectedNode=pkg.engines?.node, expectedNpm=pkg.engines?.npm;
if(process.version!==`v${expectedNode}`){gate('runtime_pin','FAILED',{expected:expectedNode,actual:process.version});finish('FAILED')}
const npmV=spawnSync('npm',['--version'],{encoding:'utf8'}).stdout.trim();
if(npmV!==expectedNpm){gate('npm_pin','FAILED',{expected:expectedNpm,actual:npmV});finish('FAILED')}
gate('runtime_pin','PASS',{node:process.version,npm:npmV});
const locks=['package-lock.json','npm-shrinkwrap.json','pnpm-lock.yaml','yarn.lock'].filter(fs.existsSync);
if(locks.length!==1){gate('dependency_lock','FAILED',{detail:`expected exactly one lockfile; found ${locks.length}`,locks});finish('FAILED')}
ev.lockfile={path:locks[0],sha256:sha(locks[0])};gate('dependency_lock','PASS',ev.lockfile);
const vitest=process.platform==='win32'?'node_modules/.bin/vitest.cmd':'node_modules/.bin/vitest';
if(!fs.existsSync(vitest)){gate('dependencies_installed','BLOCKED',{detail:'repository-local Vitest missing; run deterministic locked install first'});finish('ENVIRONMENT_BLOCKED')}
const tests=[]; const walk=d=>{for(const e of fs.readdirSync(d,{withFileTypes:true})){const p=`${d}/${e.name}`;if(e.isDirectory()&&!p.includes('node_modules'))walk(p);else if(/\.(test|spec)\.[cm]?[jt]sx?$/.test(e.name))tests.push(p)}}; walk('.');
if(tests.length===0){gate('test_discovery','FAILED',{detail:'zero tests discovered'});finish('FAILED')}
gate('test_discovery','PASS',{discovered_files:tests.length});
const resultFile=`qualification-vitest-${RELEASE}.json`;
const r=spawnSync(vitest,['run','--reporter=json',`--outputFile=${resultFile}`],{encoding:'utf8',timeout:120000,env:{...process.env,CI:'1'}});
if(r.error?.code==='ETIMEDOUT'){gate('vitest','FAILED',{detail:'120s hard timeout'});finish('FAILED')}
if(r.status!==0){gate('vitest','FAILED',{exit_code:r.status,stderr:String(r.stderr||'').slice(0,1500)});finish('FAILED')}
if(!fs.existsSync(resultFile)){gate('vitest_evidence','FAILED',{detail:'runner exited 0 but JSON evidence missing'});finish('FAILED')}
const vr=JSON.parse(fs.readFileSync(resultFile,'utf8')); const total=vr.numTotalTests??0, failed=vr.numFailedTests??0;
if(total===0||failed!==0){gate('vitest_evidence','FAILED',{total_tests:total,failed_tests:failed});finish('FAILED')}
gate('vitest','PASS',{total_tests:total,passed_tests:vr.numPassedTests??null,failed_tests:failed,result_sha256:sha(resultFile)});
const nextBin=process.platform==='win32'?'node_modules/.bin/next.cmd':'node_modules/.bin/next';
if(!fs.existsSync(nextBin)){gate('build','BLOCKED',{detail:'repository-local Next.js binary missing'});finish('ENVIRONMENT_BLOCKED')}
const b=spawnSync(nextBin,['build'],{encoding:'utf8',timeout:180000,env:{...process.env,CI:'1',NEXT_TELEMETRY_DISABLED:'1'}});
if(b.error?.code==='ETIMEDOUT'){gate('build','FAILED',{detail:'180s hard timeout'});finish('FAILED')}
if(b.status!==0){gate('build','FAILED',{exit_code:b.status,stdout:String(b.stdout||'').slice(-2000),stderr:String(b.stderr||'').slice(-2000)});finish('FAILED')}
gate('build','PASS');
// Repository qualification is necessary but not sufficient for production readiness.
const pg=spawnSync(process.execPath,['scripts/production-gate-check.mjs'],{encoding:'utf8'});
if(pg.status!==0){
  const prod=JSON.parse(fs.readFileSync('PRODUCTION-QUALIFICATION.json','utf8'));
  const openProd=Object.entries(prod.gates||{}).filter(([,v])=>v?.status!=='passed').map(([k])=>k);
  gate('production_external_gates','BLOCKED',{open_gates:openProd,validator_exit:pg.status});
  finish('REPOSITORY_QUALIFIED_PRODUCTION_BLOCKED');
}
gate('production_external_gates','PASS');
finish('QUALIFIED');
