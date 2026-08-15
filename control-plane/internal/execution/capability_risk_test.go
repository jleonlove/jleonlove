package execution

import (
	"errors"
	"testing"
)

func TestCapabilityIncreaseRequiresRequalification(t *testing.T) {
	old := CapabilityRiskProfile{ReleaseID: "r1", EvaluationDigest: "e1", CyberTier: 1}
	next := CapabilityRiskProfile{ReleaseID: "r2", EvaluationDigest: "e2", CyberTier: 2}
	if !RequiresRequalification(old, next) {
		t.Fatal("capability increase inherited prior qualification")
	}
}

func TestInsufficientSafeguardDenied(t *testing.T) {
	p := CapabilityRiskProfile{RequiredSafeguard: SafeguardHumanGated}
	if err := p.ValidateRuntime(SafeguardRestricted); !errors.Is(err, ErrSafeguardInsufficient) {
		t.Fatalf("err=%v want ErrSafeguardInsufficient", err)
	}
}

func TestEqualOrStrongerSafeguardAllowed(t *testing.T) {
	p := CapabilityRiskProfile{RequiredSafeguard: SafeguardRestricted}
	if err := p.ValidateRuntime(SafeguardHumanGated); err != nil {
		t.Fatal(err)
	}
}
