package extensions

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

type State string
const (
	Draft State = "DRAFT"
	Quarantined State = "QUARANTINED"
	Qualified State = "QUALIFIED"
	Approved State = "APPROVED"
	Active State = "ACTIVE"
	Revoked State = "REVOKED"
)

var (
	ErrNotApproved = errors.New("extension not approved")
	ErrRevoked = errors.New("extension revoked")
	ErrCapabilityIncrease = errors.New("capability increase requires requalification")
)

type Manifest struct {
	ID string
	Version string
	Publisher string
	ArtifactDigest string
	Skills []string
	MCPServers []string
	Hooks []string
	Network []string
	Filesystem []string
	Secrets []string
	Capabilities []string
}

type Record struct {
	Manifest Manifest
	State State
	QualificationDigest string
}

type Registry struct{ records map[string]Record }

func NewRegistry() *Registry { return &Registry{records: map[string]Record{}} }
func key(id, version string) string { return id + "@" + version }

func (r *Registry) Put(rec Record) { r.records[key(rec.Manifest.ID, rec.Manifest.Version)] = rec }

func (r *Registry) Activate(id, version string) error {
	k := key(id,version)
	rec, ok := r.records[k]
	if !ok || (rec.State != Approved && rec.State != Active) { return ErrNotApproved }
	if rec.State == Revoked { return ErrRevoked }
	rec.State = Active
	r.records[k] = rec
	return nil
}

func (r *Registry) Revoke(id, version string) {
	k:=key(id,version); rec,ok:=r.records[k]; if !ok{return}; rec.State=Revoked; r.records[k]=rec
}

func (r *Registry) CanInvoke(id, version string) error {
	rec,ok:=r.records[key(id,version)]
	if !ok { return ErrNotApproved }
	if rec.State==Revoked { return ErrRevoked }
	if rec.State!=Active { return ErrNotApproved }
	return nil
}

func normalized(in []string) []string {
	out:=append([]string(nil),in...); sort.Strings(out); return out
}
func set(in []string) map[string]bool { m:=map[string]bool{}; for _,v:=range in{m[v]=true}; return m }
func added(old,new []string) []string {
	o:=set(old); var a []string; for _,v:=range new{if !o[v]{a=append(a,v)}}; sort.Strings(a); return a
}

type CapabilityDiff struct {
	AddedCapabilities []string
	AddedNetwork []string
	AddedFilesystem []string
	AddedSecrets []string
	AddedMCPServers []string
	AddedHooks []string
	RequiresRequalification bool
}

func Diff(oldM,newM Manifest) CapabilityDiff {
	d:=CapabilityDiff{
		AddedCapabilities:added(oldM.Capabilities,newM.Capabilities),
		AddedNetwork:added(oldM.Network,newM.Network),
		AddedFilesystem:added(oldM.Filesystem,newM.Filesystem),
		AddedSecrets:added(oldM.Secrets,newM.Secrets),
		AddedMCPServers:added(oldM.MCPServers,newM.MCPServers),
		AddedHooks:added(oldM.Hooks,newM.Hooks),
	}
	d.RequiresRequalification=len(d.AddedCapabilities)+len(d.AddedNetwork)+len(d.AddedFilesystem)+len(d.AddedSecrets)+len(d.AddedMCPServers)+len(d.AddedHooks)>0
	return d
}

func EffectiveManifestDigest(agentID, releaseID string, manifests []Manifest) string {
	parts:=[]string{agentID,releaseID}
	sort.Slice(manifests,func(i,j int)bool{return key(manifests[i].ID,manifests[i].Version)<key(manifests[j].ID,manifests[j].Version)})
	for _,m:=range manifests{
		parts=append(parts,key(m.ID,m.Version),m.Publisher,m.ArtifactDigest)
		parts=append(parts,normalized(m.Capabilities)...)
	}
	sum:=sha256.Sum256([]byte(strings.Join(parts,"\x00")))
	return hex.EncodeToString(sum[:])
}
