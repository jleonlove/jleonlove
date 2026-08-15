package supplychain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

var (
	ErrUnknownPublisher = errors.New("unknown publisher key")
	ErrRevokedPublisher = errors.New("publisher key revoked")
	ErrInvalidSignature = errors.New("invalid package signature")
	ErrArtifactMismatch = errors.New("artifact digest mismatch")
	ErrSBOMMismatch = errors.New("sbom digest mismatch")
	ErrProvenanceMismatch = errors.New("provenance mismatch")
)

type Provenance struct {
	SourceURI string `json:"source_uri"`
	SourceRevision string `json:"source_revision"`
	BuilderID string `json:"builder_id"`
}

type PackageStatement struct {
	ExtensionID string `json:"extension_id"`
	Version string `json:"version"`
	PublisherID string `json:"publisher_id"`
	KeyID string `json:"key_id"`
	ArtifactDigest string `json:"artifact_digest"`
	SBOMDigest string `json:"sbom_digest"`
	Dependencies []string `json:"dependencies"`
	Provenance Provenance `json:"provenance"`
}

type SignedPackage struct {
	Statement PackageStatement `json:"statement"`
	Signature []byte `json:"signature"`
}

type PublisherKey struct {
	PublisherID string
	KeyID string
	PublicKey ed25519.PublicKey
	Revoked bool
}

type TrustStore struct { keys map[string]PublisherKey }

func NewTrustStore(keys ...PublisherKey) *TrustStore {
	t:=&TrustStore{keys:map[string]PublisherKey{}}
	for _,k:=range keys { t.keys[k.PublisherID+"\x00"+k.KeyID]=k }
	return t
}

func (t *TrustStore) Revoke(publisherID,keyID string) {
	k:=publisherID+"\x00"+keyID
	v,ok:=t.keys[k]; if !ok{return}; v.Revoked=true; t.keys[k]=v
}

func canonical(s PackageStatement)([]byte,error){
	s.Dependencies=append([]string(nil),s.Dependencies...)
	sort.Strings(s.Dependencies)
	return json.Marshal(s)
}

func Sign(priv ed25519.PrivateKey,s PackageStatement)(SignedPackage,error){
	b,err:=canonical(s); if err!=nil{return SignedPackage{},err}
	return SignedPackage{Statement:s,Signature:ed25519.Sign(priv,b)},nil
}

func Digest(b []byte) string {
	sum:=sha256.Sum256(b); return "sha256:"+hex.EncodeToString(sum[:])
}

func Verify(t *TrustStore,p SignedPackage,artifact,sbom []byte,expected Provenance) error {
	k,ok:=t.keys[p.Statement.PublisherID+"\x00"+p.Statement.KeyID]
	if !ok{return ErrUnknownPublisher}
	if k.Revoked{return ErrRevokedPublisher}
	b,err:=canonical(p.Statement); if err!=nil{return err}
	if !ed25519.Verify(k.PublicKey,b,p.Signature){return ErrInvalidSignature}
	if Digest(artifact)!=p.Statement.ArtifactDigest{return ErrArtifactMismatch}
	if Digest(sbom)!=p.Statement.SBOMDigest{return ErrSBOMMismatch}
	if p.Statement.Provenance!=expected{return ErrProvenanceMismatch}
	return nil
}

func DependencyFingerprint(deps []string) string {
	cp:=append([]string(nil),deps...); sort.Strings(cp)
	return Digest([]byte(strings.Join(cp,"\x00")))
}
