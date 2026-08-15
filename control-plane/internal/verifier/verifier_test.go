package verifier
import("errors";"testing")
func TestGeneratorConfidenceIsNotProof(t *testing.T){r:=Verify(Claim{Text:"x",GeneratorConfidence:.999});if r.Verified||r.Status!=Unverified{t.Fatal("confidence became proof")}}
func TestIndependentEvidenceSupports(t *testing.T){r:=Verify(Claim{Text:"x",Evidence:[]Evidence{{Source:"a",Independent:true,Supports:true,Reliability:.9},{Source:"b",Independent:true,Supports:true,Reliability:.9}}});if !r.Verified||r.Status!=Supported{t.Fatalf("%+v",r)}}
func TestConflictSurfaced(t *testing.T){r:=Verify(Claim{Text:"x",Evidence:[]Evidence{{Source:"a",Independent:true,Supports:true,Reliability:.9},{Source:"b",Independent:true,Supports:false,Reliability:.8}}});if r.Status!=Conflicted{t.Fatalf("%+v",r)};if e:=RequireVerified(r);!errors.Is(e,ErrConflict){t.Fatal(e)}}
func TestWeakEvidenceCannotAuthorize(t *testing.T){r:=Verify(Claim{Text:"x",Evidence:[]Evidence{{Source:"a",Supports:true,Reliability:.4}}});if e:=RequireVerified(r);!errors.Is(e,ErrInsufficientEvidence){t.Fatal(e)}}
func TestKnownNeedsStrongIndependentEvidence(t *testing.T){r:=Verify(Claim{Text:"x",Evidence:[]Evidence{{Source:"a",Independent:true,Supports:true,Reliability:.95},{Source:"b",Independent:true,Supports:true,Reliability:.95},{Source:"c",Independent:true,Supports:true,Reliability:.95}}});if r.Status!=Known||!r.Verified{t.Fatalf("%+v",r)}}
