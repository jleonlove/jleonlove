import { spawnSync } from 'node:child_process';
import fs from 'node:fs';
const locks=['package-lock.json','npm-shrinkwrap.json','pnpm-lock.yaml','yarn.lock'].filter(fs.existsSync);
if(locks.length>1){console.error(`Refusing ambiguous dependency state: ${locks.join(', ')}`);process.exit(1)}
if(locks.length===1){console.log(`Canonical lockfile already exists: ${locks[0]}`);process.exit(0)}
const ping=spawnSync('npm',['ping','--registry=https://registry.npmjs.org/'],{encoding:'utf8',timeout:15000});
if(ping.status!==0){console.error('npm registry unavailable; lockfile was NOT fabricated.');console.error(String(ping.stderr||ping.stdout||'').slice(0,1000));process.exit(2)}
const r=spawnSync('npm',['install','--package-lock-only','--ignore-scripts','--no-audit','--no-fund'],{stdio:'inherit',timeout:180000});
if(r.error?.code==='ETIMEDOUT'){console.error('Lock generation timed out');process.exit(2)}
if(r.status!==0||!fs.existsSync('package-lock.json')){console.error('Lock generation failed');process.exit(1)}
console.log('package-lock.json generated successfully');
