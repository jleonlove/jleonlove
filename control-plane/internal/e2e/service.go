package e2e

import (
	"context"
	"crypto/ed25519"
	"errors"
	"sync"
	"time"

	"atlas/internal/containment"
	"atlas/internal/envelope"
	"atlas/internal/execution"
	"atlas/internal/replay"
	"atlas/internal/trajectory"
)

var (
	ErrDenied   = errors.New("governed execution denied")
	ErrAudience = errors.New("audience mismatch")
	ErrExpired  = errors.New("envelope expired")
	ErrRelease  = errors.New("release mismatch")
)

type Runtime interface {
	Execute(context.Context, envelope.Action, string) ([]byte, error)
}

type Request struct {
	Envelope            envelope.Signed
	Action              envelope.Action
	PublicKey           ed25519.PublicKey
	Risk                execution.CapabilityRiskProfile
	Safeguard           execution.SafeguardClass
	Trajectory          trajectory.State
	Limits              trajectory.Limits
	Attestation         containment.Attestation
	ExpectedContainment containment.Expected
}

type Evidence struct {
	ExecutionID string
	Outcome     string
	Reason      string
}

type Service struct {
	Audience string
	Replay   replay.Store
	Runtime  Runtime
	now      func() time.Time
	mu       sync.Mutex
	evidence []Evidence
}

func New(audience string, rs replay.Store, rt Runtime) *Service {
	return &Service{Audience: audience, Replay: rs, Runtime: rt, now: time.Now}
}

func (s *Service) Execute(ctx context.Context, r Request) ([]byte, error) {
	deny := func(reason string, err error) ([]byte, error) {
		s.record(Evidence{Outcome: "DENIED", Reason: reason})
		return nil, errors.Join(ErrDenied, err)
	}
	if err := envelope.Verify(r.PublicKey, r.Envelope); err != nil {
		return deny("signature", err)
	}
	c := r.Envelope.Claims
	now := s.now()
	if c.Audience != s.Audience {
		return deny("audience", ErrAudience)
	}
	if now.Before(c.IssuedAt) || !now.Before(c.ExpiresAt) {
		return deny("time", ErrExpired)
	}
	if err := envelope.VerifyAction(r.Envelope, r.Action); err != nil {
		return deny("action_binding", err)
	}
	if r.Risk.ReleaseID != c.ReleaseID {
		return deny("risk_release", ErrRelease)
	}
	if err := r.Risk.ValidateRuntime(r.Safeguard); err != nil {
		return deny("frontier_gate", err)
	}
	if r.Trajectory.TrajectoryID != c.TrajectoryID {
		return deny("trajectory_binding", ErrRelease)
	}
	if err := trajectory.Evaluate(now, r.Trajectory, trajectory.Action{Name: r.Action.Operation}, r.Limits); err != nil {
		return deny("trajectory", err)
	}
	if r.ExpectedContainment.AgentID != c.AgentID || r.ExpectedContainment.ReleaseID != c.ReleaseID {
		return deny("containment_binding", ErrRelease)
	}
	if err := containment.Verify(now, r.Attestation, r.ExpectedContainment); err != nil {
		return deny("containment", err)
	}
	key := replay.Key{OrganizationID: c.OrganizationID, RequestID: c.RequestID, Nonce: c.Nonce}
	if err := s.Replay.Consume(ctx, key, c.ExpiresAt); err != nil {
		return deny("replay", err)
	}

	executionID := execution.ID(c.OrganizationID, c.RequestID, c.ActionDigest)
	s.record(Evidence{ExecutionID: executionID, Outcome: "PREPARED"})
	out, err := s.Runtime.Execute(ctx, r.Action, executionID)
	if err != nil {
		s.record(Evidence{ExecutionID: executionID, Outcome: "FAILED", Reason: err.Error()})
		return nil, err
	}
	s.record(Evidence{ExecutionID: executionID, Outcome: "SUCCEEDED"})
	return out, nil
}

func (s *Service) record(e Evidence) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evidence = append(s.evidence, e)
}
func (s *Service) Evidence() []Evidence {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Evidence, len(s.evidence))
	copy(out, s.evidence)
	return out
}
