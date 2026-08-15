package execution

import (
	"errors"
	"testing"
)

func TestUnknownSafeguardDenied(t *testing.T) {
	p := CapabilityRiskProfile{RequiredSafeguard: SafeguardRestricted}
	if err := p.ValidateRuntime(SafeguardClass("UNKNOWN")); !errors.Is(err, ErrSafeguardInsufficient) {
		t.Fatalf("err=%v want ErrSafeguardInsufficient", err)
	}
}

func TestBlockedReleaseCannotExecuteEvenInBlockedClass(t *testing.T) {
	p := CapabilityRiskProfile{RequiredSafeguard: SafeguardBlocked}
	if err := p.ValidateRuntime(SafeguardBlocked); !errors.Is(err, ErrSafeguardInsufficient) {
		t.Fatalf("err=%v want ErrSafeguardInsufficient", err)
	}
}

func TestEvaluationDigestChangeRequiresRequalification(t *testing.T) {
	old := CapabilityRiskProfile{ReleaseID: "r1", EvaluationDigest: "eval-a", CyberTier: 2}
	next := CapabilityRiskProfile{ReleaseID: "r1", EvaluationDigest: "eval-b", CyberTier: 2}
	if !RequiresRequalification(old, next) {
		t.Fatal("changed evaluation evidence inherited qualification")
	}
}

func TestCapabilityDowngradeWithSameEvidenceDoesNotInventEscalation(t *testing.T) {
	old := CapabilityRiskProfile{ReleaseID: "r1", EvaluationDigest: "eval-a", CyberTier: 3, AutonomyTier: 3}
	next := CapabilityRiskProfile{ReleaseID: "r2", EvaluationDigest: "eval-a", CyberTier: 2, AutonomyTier: 2}
	if RequiresRequalification(old, next) {
		t.Fatal("pure downgrade with identical evaluation evidence incorrectly marked escalation")
	}
}
