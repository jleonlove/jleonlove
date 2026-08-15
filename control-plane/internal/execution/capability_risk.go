package execution

import "errors"

type SafeguardClass string

const (
	SafeguardNormal       SafeguardClass = "NORMAL"
	SafeguardRestricted   SafeguardClass = "RESTRICTED"
	SafeguardHumanGated   SafeguardClass = "HUMAN_GATED"
	SafeguardIsolatedEval SafeguardClass = "ISOLATED_EVAL_ONLY"
	SafeguardBlocked      SafeguardClass = "BLOCKED"
)

type CapabilityRiskProfile struct {
	ReleaseID           string
	EvaluationDigest    string
	AutonomyTier        uint8
	CyberTier           uint8
	SelfImprovementTier uint8
	ConsequentialTier   uint8
	RequiredSafeguard   SafeguardClass
}

var (
	ErrCapabilityEscalation  = errors.New("release capability escalation requires fresh qualification")
	ErrSafeguardInsufficient = errors.New("runtime safeguard class insufficient for release")
)

func (p CapabilityRiskProfile) ValidateRuntime(actual SafeguardClass) error {
	if p.RequiredSafeguard == SafeguardBlocked {
		return ErrSafeguardInsufficient
	}
	if safeguardRank(actual) < safeguardRank(p.RequiredSafeguard) {
		return ErrSafeguardInsufficient
	}
	return nil
}

func RequiresRequalification(previous, next CapabilityRiskProfile) bool {
	return next.AutonomyTier > previous.AutonomyTier ||
		next.CyberTier > previous.CyberTier ||
		next.SelfImprovementTier > previous.SelfImprovementTier ||
		next.ConsequentialTier > previous.ConsequentialTier ||
		next.EvaluationDigest != previous.EvaluationDigest
}

func safeguardRank(s SafeguardClass) int {
	switch s {
	case SafeguardNormal:
		return 1
	case SafeguardRestricted:
		return 2
	case SafeguardHumanGated:
		return 3
	case SafeguardIsolatedEval:
		return 4
	case SafeguardBlocked:
		return 5
	default:
		return 0
	}
}
