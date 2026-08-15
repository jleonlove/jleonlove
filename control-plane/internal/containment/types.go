package containment

import (
	"errors"
	"time"
)

var ErrMismatch = errors.New("containment attestation mismatch")

type Attestation struct {
	RuntimeID string
	RuntimeEpoch uint64
	AgentID string
	ReleaseID string
	ProfileDigest string
	NetworkDigest string
	FilesystemDigest string
	RuntimeDigest string
	IssuedAt time.Time
	ExpiresAt time.Time
	Signature []byte
}

type Expected struct {
	AgentID string
	ReleaseID string
	ProfileDigest string
	NetworkDigest string
	FilesystemDigest string
	RuntimeDigest string
}

func Verify(now time.Time,a Attestation,e Expected) error {
	if !now.Before(a.ExpiresAt) { return ErrMismatch }
	if a.AgentID != e.AgentID || a.ReleaseID != e.ReleaseID ||
		a.ProfileDigest != e.ProfileDigest || a.NetworkDigest != e.NetworkDigest ||
		a.FilesystemDigest != e.FilesystemDigest || a.RuntimeDigest != e.RuntimeDigest {
		return ErrMismatch
	}
	return nil
}
