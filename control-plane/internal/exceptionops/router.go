package exceptionops

import "strings"

type Disposition string

const (
	AutoExecute Disposition = "AUTO_EXECUTE"
	HumanReview Disposition = "HUMAN_REVIEW"
	MultiPartyReview Disposition = "MULTI_PARTY_REVIEW"
	Quarantine Disposition = "QUARANTINE"
	Deny Disposition = "DENY"
)

type Request struct {
	RiskScore int
	Confidence float64
	Consequential bool
	Irreversible bool
	AffectsMoney bool
	ComplianceSensitive bool
	EvidenceComplete bool
	PolicyAllowed bool
	Reason string
}

type Decision struct {
	Disposition Disposition
	Reason string
}

// Route uses deterministic risk policy. The model may supply signals, but it
// cannot grant itself execution authority.
func Route(r Request) Decision {
	if !r.PolicyAllowed { return Decision{Deny, "policy_denied"} }
	if !r.EvidenceComplete { return Decision{Quarantine, "evidence_incomplete"} }
	if r.RiskScore < 0 || r.RiskScore > 100 || r.Confidence < 0 || r.Confidence > 1 {
		return Decision{Quarantine, "invalid_risk_signal"}
	}
	if strings.TrimSpace(r.Reason) == "" { return Decision{Quarantine, "decision_reason_required"} }
	if r.Irreversible || (r.AffectsMoney && r.ComplianceSensitive) || r.RiskScore >= 80 {
		return Decision{MultiPartyReview, "high_impact_independent_approval_required"}
	}
	if r.Consequential || r.AffectsMoney || r.ComplianceSensitive || r.RiskScore >= 40 || r.Confidence < 0.90 {
		return Decision{HumanReview, "human_judgment_required"}
	}
	return Decision{AutoExecute, "bounded_low_risk"}
}
