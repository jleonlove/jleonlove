export type ProtocolKind='mcp'|'a2a'|'trade_document'|'settlement'|'event'|'custom';
export type CertificationState='CLAIMED'|'OBSERVED'|'VERIFIED'|'DEGRADED'|'QUARANTINED'|'REVOKED';
export type TrustTier='sandbox'|'restricted'|'trusted'|'critical';

export interface ProtocolCertification {
  state: CertificationState;
  validUntil: string;
  evidenceDigest: string;
}

export interface ProtocolRegistration {
  id: string;
  protocol: ProtocolKind;
  versions: string[];
  preferredVersion: string;
  schemaDigest: string;
  dependencyDigest: string;
  signature: string;
  issuerPrincipalId: string;
  allowedPrincipals: string[];
  capabilities: string[];
  jurisdictions: string[];
  dataClasses: string[];
  trustTier: TrustTier;
  certification: ProtocolCertification;
  allowDowngrade?: boolean;
  revokedAt?: string;
  revocationReason?: string;
}

export interface ProtocolRequest {
  protocolId: string;
  requestedVersion: string;
  negotiatedVersion: string;
  callerAcceptedDowngrade?: boolean;
  principalId: string;
  capability: string;
  jurisdiction: string;
  dataClass: string;
  schemaDigest: string;
  dependencyDigest: string;
  now?: string;
}

export interface ProtocolDecision {
  allowed: boolean;
  reason: string;
  protocolId: string;
  negotiatedVersion?: string;
  trustTier?: TrustTier;
  evidence: string;
}

const required=(value:string,code:string)=>{if(!value.trim()) throw new Error(code)};
const unique=(values:string[])=>new Set(values).size===values.length;
const contains=(values:string[],value:string)=>values.includes(value);

export function validateProtocolRegistration(r:ProtocolRegistration){
  required(r.id,'PROTOCOL_ID_REQUIRED');
  required(r.schemaDigest,'PROTOCOL_SCHEMA_DIGEST_REQUIRED');
  required(r.dependencyDigest,'PROTOCOL_DEPENDENCY_DIGEST_REQUIRED');
  required(r.signature,'PROTOCOL_SIGNATURE_REQUIRED');
  required(r.issuerPrincipalId,'PROTOCOL_ISSUER_REQUIRED');
  required(r.certification.evidenceDigest,'PROTOCOL_EVIDENCE_REQUIRED');
  if(!r.versions.length||!unique(r.versions)||!contains(r.versions,r.preferredVersion)) throw new Error('PROTOCOL_VERSION_POLICY_INVALID');
  if(!r.allowedPrincipals.length||!unique(r.allowedPrincipals)) throw new Error('PROTOCOL_PRINCIPAL_POLICY_INVALID');
  if(!r.capabilities.length||!unique(r.capabilities)) throw new Error('PROTOCOL_CAPABILITY_POLICY_INVALID');
  if(!r.jurisdictions.length||!unique(r.jurisdictions)) throw new Error('PROTOCOL_JURISDICTION_POLICY_INVALID');
  if(!r.dataClasses.length||!unique(r.dataClasses)) throw new Error('PROTOCOL_DATA_POLICY_INVALID');
  const expiry=Date.parse(r.certification.validUntil);
  if(!Number.isFinite(expiry)) throw new Error('PROTOCOL_CERTIFICATION_EXPIRY_INVALID');
  return true;
}

function evidence(r:ProtocolRegistration|undefined,reason:string,version=''){
  return [r?.id??'unknown',r?.protocol??'unknown',version,reason,r?.schemaDigest??'',r?.dependencyDigest??'',r?.certification.evidenceDigest??''].join('|');
}

export class ProtocolControlPlane {
  private readonly records=new Map<string,ProtocolRegistration>();

  register(registration:ProtocolRegistration){
    validateProtocolRegistration(registration);
    if(this.records.has(registration.id)) throw new Error('PROTOCOL_DUPLICATE_REGISTRATION_DENY');
    this.records.set(registration.id,{...registration,versions:[...registration.versions],allowedPrincipals:[...registration.allowedPrincipals],capabilities:[...registration.capabilities],jurisdictions:[...registration.jurisdictions],dataClasses:[...registration.dataClasses],certification:{...registration.certification}});
    return true;
  }

  revoke(protocolId:string,reason:string,now=new Date().toISOString()){
    const r=this.records.get(protocolId);
    if(!r) throw new Error('PROTOCOL_UNKNOWN_DENY');
    required(reason,'PROTOCOL_REVOCATION_REASON_REQUIRED');
    this.records.set(protocolId,{...r,revokedAt:now,revocationReason:reason,certification:{...r.certification,state:'REVOKED'}});
    return true;
  }

  get(protocolId:string){
    const r=this.records.get(protocolId);
    return r?{...r,versions:[...r.versions],allowedPrincipals:[...r.allowedPrincipals],capabilities:[...r.capabilities],jurisdictions:[...r.jurisdictions],dataClasses:[...r.dataClasses],certification:{...r.certification}}:undefined;
  }

  authorize(request:ProtocolRequest):ProtocolDecision {
    const r=this.records.get(request.protocolId);
    const deny=(reason:string):ProtocolDecision=>({allowed:false,reason,protocolId:request.protocolId,negotiatedVersion:request.negotiatedVersion||undefined,trustTier:r?.trustTier,evidence:evidence(r,reason,request.negotiatedVersion)});
    if(!r) return deny('PROTOCOL_UNKNOWN_DENY');
    if(r.revokedAt||r.certification.state==='REVOKED') return deny('PROTOCOL_REVOKED_DENY');
    if(r.certification.state!=='VERIFIED') return deny('PROTOCOL_NOT_VERIFIED_DENY');
    const now=Date.parse(request.now??new Date().toISOString());
    const expiry=Date.parse(r.certification.validUntil);
    if(!Number.isFinite(now)||!Number.isFinite(expiry)||now>=expiry) return deny('PROTOCOL_CERTIFICATION_EXPIRED_DENY');
    if(!request.principalId||!contains(r.allowedPrincipals,request.principalId)) return deny('PROTOCOL_PRINCIPAL_DENY');
    if(!contains(r.capabilities,request.capability)) return deny('PROTOCOL_CAPABILITY_DENY');
    if(!contains(r.jurisdictions,request.jurisdiction)) return deny('PROTOCOL_JURISDICTION_DENY');
    if(!contains(r.dataClasses,request.dataClass)) return deny('PROTOCOL_DATA_CLASS_DENY');
    if(request.schemaDigest!==r.schemaDigest) return deny('PROTOCOL_SCHEMA_DRIFT_DENY');
    if(request.dependencyDigest!==r.dependencyDigest) return deny('PROTOCOL_DEPENDENCY_DRIFT_DENY');
    if(!request.requestedVersion||!request.negotiatedVersion||!contains(r.versions,request.negotiatedVersion)) return deny('PROTOCOL_VERSION_DENY');
    if(request.requestedVersion!==request.negotiatedVersion&&(!r.allowDowngrade||!request.callerAcceptedDowngrade)) return deny('PROTOCOL_DOWNGRADE_DENY');
    return {allowed:true,reason:'PROTOCOL_ALLOWED',protocolId:r.id,negotiatedVersion:request.negotiatedVersion,trustTier:r.trustTier,evidence:evidence(r,'PROTOCOL_ALLOWED',request.negotiatedVersion)};
  }
}
