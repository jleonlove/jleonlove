package corporateauthority

import "strings"

type Action string

const (
	OwnershipChange Action = "ownership_change"
	EquityIssuance Action = "equity_issuance"
	IPTransfer Action = "ip_transfer"
	ProductionRootTransfer Action = "production_root_transfer"
	RecoveryRootTransfer Action = "recovery_root_transfer"
	DomainTransfer Action = "domain_transfer"
	PolicyRootChange Action = "policy_root_change"
	EvidenceDestruction Action = "evidence_destruction"
)

type Request struct {
	Action Action
	AuthorityEvidence string
	IndependentApproval string
	ExactEffectDigest string
}

type Decision struct {
	Allowed bool
	Reason string
}

var protected = map[Action]bool{
	OwnershipChange: true,
	EquityIssuance: true,
	IPTransfer: true,
	ProductionRootTransfer: true,
	RecoveryRootTransfer: true,
	DomainTransfer: true,
	PolicyRootChange: true,
	EvidenceDestruction: true,
}

// Authorize is deliberately fail-closed. Protected corporate-control effects
// require evidence of authority, independent approval, and binding to the exact effect.
func Authorize(r Request) Decision {
	if !protected[r.Action] {
		return Decision{Allowed: false, Reason: "action_not_admitted"}
	}
	if strings.TrimSpace(r.AuthorityEvidence) == "" {
		return Decision{Reason: "authority_evidence_required"}
	}
	if strings.TrimSpace(r.IndependentApproval) == "" {
		return Decision{Reason: "independent_approval_required"}
	}
	if strings.TrimSpace(r.ExactEffectDigest) == "" {
		return Decision{Reason: "exact_effect_binding_required"}
	}
	return Decision{Allowed: true, Reason: "authorized"}
}
