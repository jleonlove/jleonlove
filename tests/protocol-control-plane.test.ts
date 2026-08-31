import {describe,it,expect} from 'vitest';
import {ProtocolControlPlane, type ProtocolRegistration} from '../lib/protocol-control-plane';

const registration=():ProtocolRegistration=>({
  id:'mcp.trade-docs',
  protocol:'mcp',
  versions:['1.0','1.1'],
  preferredVersion:'1.1',
  schemaDigest:'schema:trade-docs:v1',
  dependencyDigest:'deps:trade-docs:v1',
  signature:'sig:approved',
  issuerPrincipalId:'atlas-security',
  allowedPrincipals:['atlas-agent-trade'],
  capabilities:['document.read','document.verify'],
  jurisdictions:['US','EU'],
  dataClasses:['internal','confidential'],
  trustTier:'trusted',
  certification:{state:'VERIFIED',validUntil:'2030-01-01T00:00:00.000Z',evidenceDigest:'evidence:cert-001'},
  allowDowngrade:false,
});

const request=()=>({
  protocolId:'mcp.trade-docs',
  requestedVersion:'1.1',
  negotiatedVersion:'1.1',
  principalId:'atlas-agent-trade',
  capability:'document.verify',
  jurisdiction:'US',
  dataClass:'confidential',
  schemaDigest:'schema:trade-docs:v1',
  dependencyDigest:'deps:trade-docs:v1',
  now:'2026-08-31T16:00:00.000Z',
});

describe('protocol control plane fail closed',()=>{
  it('authorizes only a fully certified exact runtime',()=>{const c=new ProtocolControlPlane();c.register(registration());expect(c.authorize(request()).allowed).toBe(true)});
  it('denies duplicate registration',()=>{const c=new ProtocolControlPlane();c.register(registration());expect(()=>c.register(registration())).toThrow('PROTOCOL_DUPLICATE_REGISTRATION_DENY')});
  it('denies revoked protocols immediately',()=>{const c=new ProtocolControlPlane();c.register(registration());c.revoke('mcp.trade-docs','incident');expect(c.authorize(request()).reason).toBe('PROTOCOL_REVOKED_DENY')});
  it('denies expired certification',()=>{const c=new ProtocolControlPlane();const r=registration();r.certification.validUntil='2026-08-30T00:00:00.000Z';c.register(r);expect(c.authorize(request()).reason).toBe('PROTOCOL_CERTIFICATION_EXPIRED_DENY')});
  it('denies capability expansion',()=>{const c=new ProtocolControlPlane();c.register(registration());expect(c.authorize({...request(),capability:'document.delete'}).reason).toBe('PROTOCOL_CAPABILITY_DENY')});
  it('denies principal substitution',()=>{const c=new ProtocolControlPlane();c.register(registration());expect(c.authorize({...request(),principalId:'other-agent'}).reason).toBe('PROTOCOL_PRINCIPAL_DENY')});
  it('denies schema and dependency drift',()=>{const c=new ProtocolControlPlane();c.register(registration());expect(c.authorize({...request(),schemaDigest:'tampered'}).reason).toBe('PROTOCOL_SCHEMA_DRIFT_DENY');expect(c.authorize({...request(),dependencyDigest:'tampered'}).reason).toBe('PROTOCOL_DEPENDENCY_DRIFT_DENY')});
  it('denies silent downgrade and permits explicit policy plus consent',()=>{const c=new ProtocolControlPlane();const r=registration();r.allowDowngrade=true;c.register(r);expect(c.authorize({...request(),negotiatedVersion:'1.0'}).reason).toBe('PROTOCOL_DOWNGRADE_DENY');expect(c.authorize({...request(),negotiatedVersion:'1.0',callerAcceptedDowngrade:true}).allowed).toBe(true)});
  it('denies jurisdiction and data-boundary escape',()=>{const c=new ProtocolControlPlane();c.register(registration());expect(c.authorize({...request(),jurisdiction:'XX'}).reason).toBe('PROTOCOL_JURISDICTION_DENY');expect(c.authorize({...request(),dataClass:'secret'}).reason).toBe('PROTOCOL_DATA_CLASS_DENY')});
});
