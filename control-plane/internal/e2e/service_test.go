package e2e

import (
	"context"
	"crypto/ed25519"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"atlas/internal/containment"
	"atlas/internal/envelope"
	"atlas/internal/execution"
	"atlas/internal/replay"
	"atlas/internal/trajectory"
)

type fakeRuntime struct{ calls atomic.Int64 }

func (f *fakeRuntime) Execute(_ context.Context, _ envelope.Action, _ string) ([]byte, error) {
	f.calls.Add(1)
	return []byte("ok"), nil
}

func fixture(t *testing.T) (*Service, Request, *fakeRuntime) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	a := envelope.Action{Tool: "ledger", Operation: "read", Arguments: map[string]any{"account": "A"}}
	d, err := envelope.ActionDigest(a)
	if err != nil {
		t.Fatal(err)
	}
	c := envelope.Claims{Version: envelope.Version, RequestID: "req-1", OrganizationID: "org-1", PrincipalID: "p-1", AgentID: "agent-1", TaskID: "task-1", ReleaseID: "rel-1", CapabilityID: "cap-1", ScopeID: "scope-1", TrajectoryID: "traj-1", PolicyDigest: "policy-1", ActionDigest: d, Audience: "atlasd", Nonce: "nonce-1", IssuedAt: now.Add(-time.Second), ExpiresAt: now.Add(time.Minute)}
	signed, err := envelope.Sign("k1", priv, c)
	if err != nil {
		t.Fatal(err)
	}
	att := containment.Attestation{RuntimeID: "rt-1", RuntimeEpoch: 1, AgentID: "agent-1", ReleaseID: "rel-1", ProfileDigest: "p", NetworkDigest: "n", FilesystemDigest: "f", RuntimeDigest: "r", IssuedAt: now.Add(-time.Second), ExpiresAt: now.Add(time.Minute)}
	exp := containment.Expected{AgentID: "agent-1", ReleaseID: "rel-1", ProfileDigest: "p", NetworkDigest: "n", FilesystemDigest: "f", RuntimeDigest: "r"}
	rt := &fakeRuntime{}
	svc := New("atlasd", replay.NewMemoryStore(), rt)
	req := Request{Envelope: signed, Action: a, PublicKey: pub, Risk: execution.CapabilityRiskProfile{ReleaseID: "rel-1", EvaluationDigest: "eval", RequiredSafeguard: execution.SafeguardRestricted}, Safeguard: execution.SafeguardRestricted, Trajectory: trajectory.State{TrajectoryID: "traj-1", Version: 1, Status: trajectory.Active, StartedAt: now.Add(-time.Minute)}, Limits: trajectory.Limits{MaxActions: 10, MaxRetries: 3, MaxSpendCents: 1000, MaxTokens: 1000, MaxRiskScore: 10, MaxDuration: time.Hour}, Attestation: att, ExpectedContainment: exp}
	return svc, req, rt
}

func TestGovernedPathExecutesOnce(t *testing.T) {
	svc, req, rt := fixture(t)
	out, err := svc.Execute(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "ok" {
		t.Fatalf("out=%q", out)
	}
	if rt.calls.Load() != 1 {
		t.Fatalf("calls=%d", rt.calls.Load())
	}
}

func TestModifiedActionNeverReachesRuntime(t *testing.T) {
	svc, req, rt := fixture(t)
	req.Action.Arguments["account"] = "attacker"
	if _, err := svc.Execute(context.Background(), req); err == nil {
		t.Fatal("expected denial")
	}
	if rt.calls.Load() != 0 {
		t.Fatal("modified action reached runtime")
	}
}

func TestBlockedFrontierReleaseNeverReachesRuntime(t *testing.T) {
	svc, req, rt := fixture(t)
	req.Risk.RequiredSafeguard = execution.SafeguardBlocked
	req.Safeguard = execution.SafeguardBlocked
	if _, err := svc.Execute(context.Background(), req); err == nil {
		t.Fatal("expected denial")
	}
	if rt.calls.Load() != 0 {
		t.Fatal("blocked release reached runtime")
	}
}

func TestContainmentDriftNeverReachesRuntime(t *testing.T) {
	svc, req, rt := fixture(t)
	req.Attestation.NetworkDigest = "drifted"
	if _, err := svc.Execute(context.Background(), req); err == nil {
		t.Fatal("expected denial")
	}
	if rt.calls.Load() != 0 {
		t.Fatal("containment drift reached runtime")
	}
}

func TestConcurrentReplayExactlyOneRuntimeCall(t *testing.T) {
	svc, req, rt := fixture(t)
	var wg sync.WaitGroup
	var success atomic.Int64
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := svc.Execute(context.Background(), req); err == nil {
				success.Add(1)
			}
		}()
	}
	wg.Wait()
	if success.Load() != 1 {
		t.Fatalf("success=%d want=1", success.Load())
	}
	if rt.calls.Load() != 1 {
		t.Fatalf("runtime calls=%d want=1", rt.calls.Load())
	}
}
