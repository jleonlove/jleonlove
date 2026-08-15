const injectionPatterns=[/ignore\s+(all\s+)?previous/i,/system\s+prompt/i,/developer\s+message/i,/reveal\s+(the\s+)?secret/i,/bypass\s+(policy|approval|compliance|kyc)/i,/send\s+funds\s+(now|immediately)/i];
export interface AbuseScan { safe:boolean; findings:string[]; }
export function scanUntrustedText(text:string):AbuseScan { const findings=injectionPatterns.filter(p=>p.test(text)).map(p=>`UNTRUSTED_INSTRUCTION:${p.source}`); return {safe:findings.length===0,findings}; }
export function enforceTenantBoundary(resourceOrg:string,resourceWorkspace:string,activeOrg:string,activeWorkspace:string){ if(resourceOrg!==activeOrg||resourceWorkspace!==activeWorkspace) throw new Error('TENANT_BOUNDARY_VIOLATION'); return true; }
