package execution

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type State string
const (
	Prepared State = "PREPARED"
	Executing State = "EXECUTING"
	Succeeded State = "SUCCEEDED"
	Failed State = "FAILED"
	Uncertain State = "UNCERTAIN"
	Reconciled State = "RECONCILED"
)

type Record struct {
	ExecutionID string
	RequestID string
	OrganizationID string
	AgentID string
	TaskID string
	TrajectoryID string
	ReleaseID string
	State State
	Version uint64
	ActionDigest string
	PolicyDigest string
	ContainmentDigest string
	ProviderIdempotencyKey string
	PreparedAt time.Time
	ResultDigest string
	ErrorCode string
	PreviousHash string
	Hash string
}

func ID(orgID,requestID,actionDigest string) string {
	sum:=sha256.Sum256([]byte(orgID+"\x00"+requestID+"\x00"+actionDigest))
	return hex.EncodeToString(sum[:])
}
