import type { AtlasDeal } from './types'; import { closingGate } from './closing-engine'; import { authorityChain } from './deal-hardening'; import { evidenceLineage } from './evidence-lineage';
export type ConsequentialAction='SEND_CONFIDENTIAL_DOCUMENTS'|'COMMIT_TERMS'|'AUTHORIZE_PAYMENT'|'RELEASE_FUNDS'|'SIGN_OR_BIND'|'CHANGE_BENEFICIARY';
export interface ExecutionDecision { allow:boolean; action:ConsequentialAction; reasons:string[]; requiredApprovals:string[]; }
export function authorizeExecution(deal:AtlasDeal,action:ConsequentialAction,approvals:string[]=[]):ExecutionDecision{
 const gate=closingGate(deal), authority=authorityChain(deal), lineage=evidenceLineage(deal), reasons=[...gate.reasons];
 if(!authority.pass) reasons.push(...authority.unverified.map(x=>`Authority unresolved: ${x}`));
 if(lineage.contradictions.length) reasons.push('Contradictory evidence remains unresolved.');
 if(lineage.expiredClaims.length) reasons.push('Expired evidence requires refresh.');
 const required=action==='SEND_CONFIDENTIAL_DOCUMENTS'?['data_release_approval']:action==='AUTHORIZE_PAYMENT'||action==='RELEASE_FUNDS'?['finance_approval','compliance_approval']:action==='SIGN_OR_BIND'||action==='COMMIT_TERMS'?['authorized_signer_approval']:action==='CHANGE_BENEFICIARY'?['finance_approval','counterparty_reverification']:[];
 const missing=required.filter(x=>!approvals.includes(x)); if(missing.length) reasons.push(`Missing approvals: ${missing.join(', ')}`);
 return {allow:reasons.length===0,action,reasons,requiredApprovals:required};
}
