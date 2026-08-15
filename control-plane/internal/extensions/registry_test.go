package extensions

import (
	"errors"
	"testing"
)

func base() Manifest {
	return Manifest{ID:"browser-review",Version:"1.0.0",Publisher:"verified-publisher",
		ArtifactDigest:"sha256:a",Capabilities:[]string{"browser.read"},
		Network:[]string{"docs.example"},Filesystem:[]string{"/workspace:ro"}}
}

func TestInstallDoesNotGrantAuthority(t *testing.T){
	r:=NewRegistry(); r.Put(Record{Manifest:base(),State:Quarantined})
	if err:=r.Activate("browser-review","1.0.0");!errors.Is(err,ErrNotApproved){t.Fatalf("got %v",err)}
}

func TestApprovedCanActivate(t *testing.T){
	r:=NewRegistry(); r.Put(Record{Manifest:base(),State:Approved})
	if err:=r.Activate("browser-review","1.0.0");err!=nil{t.Fatal(err)}
	if err:=r.CanInvoke("browser-review","1.0.0");err!=nil{t.Fatal(err)}
}

func TestRevocationStopsInvocation(t *testing.T){
	r:=NewRegistry(); r.Put(Record{Manifest:base(),State:Approved})
	if err:=r.Activate("browser-review","1.0.0");err!=nil{t.Fatal(err)}
	r.Revoke("browser-review","1.0.0")
	if err:=r.CanInvoke("browser-review","1.0.0");!errors.Is(err,ErrRevoked){t.Fatalf("got %v",err)}
}

func TestCapabilityDiffForcesRequalification(t *testing.T){
	old:=base(); newer:=base(); newer.Version="1.1.0"
	newer.Secrets=[]string{"GITHUB_TOKEN"}; newer.Capabilities=[]string{"browser.read","repository.write"}
	d:=Diff(old,newer)
	if !d.RequiresRequalification{t.Fatal("authority increase did not require requalification")}
	if len(d.AddedSecrets)!=1 || len(d.AddedCapabilities)!=1{t.Fatalf("bad diff: %#v",d)}
}

func TestNoCapabilityIncreaseDoesNotForceRequalification(t *testing.T){
	old:=base(); newer:=base(); newer.Version="1.0.1"; newer.ArtifactDigest="sha256:b"
	d:=Diff(old,newer)
	if d.RequiresRequalification{t.Fatalf("unexpected capability escalation: %#v",d)}
}

func TestEffectiveManifestDigestChangesWithLoadedExtension(t *testing.T){
	a:=EffectiveManifestDigest("agent","release",[]Manifest{base()})
	m:=base(); m.Capabilities=append(m.Capabilities,"repository.write")
	b:=EffectiveManifestDigest("agent","release",[]Manifest{m})
	if a==b{t.Fatal("effective manifest digest did not change")}
}
