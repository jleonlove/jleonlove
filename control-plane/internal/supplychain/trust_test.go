package supplychain

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
)

func fixture(t *testing.T)(*TrustStore,SignedPackage,[]byte,[]byte,Provenance){
	t.Helper()
	pub,priv,err:=ed25519.GenerateKey(rand.Reader); if err!=nil{t.Fatal(err)}
	artifact:=[]byte("immutable extension package")
	sbom:=[]byte(`{"bomFormat":"CycloneDX","components":[]}`)
	prov:=Provenance{SourceURI:"https://example.invalid/repo",SourceRevision:"abc123",BuilderID:"atlas-builder-v1"}
	stmt:=PackageStatement{ExtensionID:"browser-review",Version:"1.0.0",PublisherID:"publisher-1",KeyID:"key-1",
		ArtifactDigest:Digest(artifact),SBOMDigest:Digest(sbom),Dependencies:[]string{"dep-b@2","dep-a@1"},Provenance:prov}
	signed,err:=Sign(priv,stmt); if err!=nil{t.Fatal(err)}
	return NewTrustStore(PublisherKey{PublisherID:"publisher-1",KeyID:"key-1",PublicKey:pub}),signed,artifact,sbom,prov
}

func TestValidPackage(t *testing.T){ts,p,a,s,prov:=fixture(t);if err:=Verify(ts,p,a,s,prov);err!=nil{t.Fatal(err)}}
func TestPackageSubstitutionDenied(t *testing.T){ts,p,_,s,prov:=fixture(t);if err:=Verify(ts,p,[]byte("substituted"),s,prov);!errors.Is(err,ErrArtifactMismatch){t.Fatalf("got %v",err)}}
func TestSBOMSubstitutionDenied(t *testing.T){ts,p,a,_,prov:=fixture(t);if err:=Verify(ts,p,a,[]byte("fake sbom"),prov);!errors.Is(err,ErrSBOMMismatch){t.Fatalf("got %v",err)}}
func TestTamperedStatementDenied(t *testing.T){ts,p,a,s,prov:=fixture(t);p.Statement.Version="9.9.9";if err:=Verify(ts,p,a,s,prov);!errors.Is(err,ErrInvalidSignature){t.Fatalf("got %v",err)}}
func TestRevokedPublisherDenied(t *testing.T){ts,p,a,s,prov:=fixture(t);ts.Revoke("publisher-1","key-1");if err:=Verify(ts,p,a,s,prov);!errors.Is(err,ErrRevokedPublisher){t.Fatalf("got %v",err)}}
func TestUnknownPublisherDenied(t *testing.T){_,p,a,s,prov:=fixture(t);if err:=Verify(NewTrustStore(),p,a,s,prov);!errors.Is(err,ErrUnknownPublisher){t.Fatalf("got %v",err)}}
func TestProvenanceMismatchDenied(t *testing.T){ts,p,a,s,_:=fixture(t);bad:=Provenance{SourceURI:"other",SourceRevision:"abc123",BuilderID:"atlas-builder-v1"};if err:=Verify(ts,p,a,s,bad);!errors.Is(err,ErrProvenanceMismatch){t.Fatalf("got %v",err)}}
func TestDependencyFingerprintOrderIndependent(t *testing.T){a:=DependencyFingerprint([]string{"b@2","a@1"});b:=DependencyFingerprint([]string{"a@1","b@2"});if a!=b{t.Fatal("dependency fingerprint is order-dependent")}}
