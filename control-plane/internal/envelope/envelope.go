package envelope

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"
)

const Version = "atlas-envelope/v1"

var (
	ErrInvalidSignature = errors.New("invalid envelope signature")
	ErrActionMismatch   = errors.New("action digest mismatch")
)

type Claims struct {
	Version        string    `json:"version"`
	RequestID      string    `json:"request_id"`
	OrganizationID string    `json:"organization_id"`
	PrincipalID    string    `json:"principal_id"`
	AgentID        string    `json:"agent_id"`
	TaskID         string    `json:"task_id"`
	ReleaseID      string    `json:"release_id"`
	CapabilityID   string    `json:"capability_id"`
	ScopeID        string    `json:"scope_id"`
	TrajectoryID   string    `json:"trajectory_id"`
	PolicyDigest   string    `json:"policy_digest"`
	ActionDigest   string    `json:"action_digest"`
	Audience       string    `json:"audience"`
	Nonce          string    `json:"nonce"`
	IssuedAt       time.Time `json:"issued_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type Signed struct {
	KeyID string `json:"kid"`
	Claims Claims `json:"claims"`
	Signature []byte `json:"signature"`
}

type Action struct {
	Tool string `json:"tool"`
	Operation string `json:"operation"`
	Arguments map[string]any `json:"arguments"`
}

func Canonical(c Claims) ([]byte, error) { return json.Marshal(c) }

func Sign(keyID string, key ed25519.PrivateKey, c Claims) (Signed, error) {
	b, err := Canonical(c)
	if err != nil { return Signed{}, err }
	return Signed{KeyID:keyID, Claims:c, Signature:ed25519.Sign(key,b)}, nil
}

func Verify(pub ed25519.PublicKey, s Signed) error {
	b, err := Canonical(s.Claims)
	if err != nil { return err }
	if !ed25519.Verify(pub,b,s.Signature) { return ErrInvalidSignature }
	return nil
}

func ActionDigest(a Action) (string,error) {
	b,err:=json.Marshal(a)
	if err!=nil{return "",err}
	sum:=sha256.Sum256(b)
	return hex.EncodeToString(sum[:]),nil
}

func VerifyAction(s Signed,a Action) error {
	d,err:=ActionDigest(a)
	if err!=nil{return err}
	if d!=s.Claims.ActionDigest{return ErrActionMismatch}
	return nil
}
