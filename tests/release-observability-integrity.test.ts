import {describe,it,expect} from 'vitest';
import {observationHash,verifyObservations,qualifyImmutableReleaseSet} from '../lib/release-observability-integrity';
describe('RC-000145 release and observability integrity',()=>{
 it('accepts intact observation chain',()=>{const h=observationHash(1,'ALERT','sev1','GENESIS');expect(verifyObservations([{sequence:1,kind:'ALERT',payload:'sev1',previousHash:'GENESIS',hash:h}]).safe).toBe(true)});
 it('rejects tampered observation',()=>{expect(verifyObservations([{sequence:1,kind:'ALERT',payload:'changed',previousHash:'GENESIS',hash:'bad'}]).code).toBe('OBSERVATION_TAMPERED')});
 it('rejects missing release component',()=>{expect(qualifyImmutableReleaseSet({releaseId:'RC-1',artifactHash:'a',checksumHash:'a',manifestHash:'',provenanceHash:'p',retrievedArtifactHash:'a'}).eligible).toBe(false)});
 it('rejects release id reuse',()=>{expect(qualifyImmutableReleaseSet({releaseId:'RC-1',artifactHash:'a',checksumHash:'a',manifestHash:'m',provenanceHash:'p',retrievedArtifactHash:'a',duplicateReleaseId:true}).code).toBe('RELEASE_ID_REUSE')});
 it('rejects retrieved artifact mismatch',()=>{expect(qualifyImmutableReleaseSet({releaseId:'RC-1',artifactHash:'a',checksumHash:'a',manifestHash:'m',provenanceHash:'p',retrievedArtifactHash:'b'}).eligible).toBe(false)});
});
