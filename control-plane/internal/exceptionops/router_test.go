package exceptionops

import "testing"

func TestLowRiskCanAutoExecute(t *testing.T) {
	d := Route(Request{RiskScore: 10, Confidence: .98, EvidenceComplete: true, PolicyAllowed: true, Reason: "routine read"})
	if d.Disposition != AutoExecute { t.Fatalf("unexpected: %+v", d) }
}

func TestConsequentialRequiresHuman(t *testing.T) {
	d := Route(Request{RiskScore: 20, Confidence: .99, Consequential: true, EvidenceComplete: true, PolicyAllowed: true, Reason: "external effect"})
	if d.Disposition != HumanReview { t.Fatalf("unexpected: %+v", d) }
}

func TestHighImpactRequiresIndependentReview(t *testing.T) {
	d := Route(Request{RiskScore: 85, Confidence: .99, AffectsMoney: true, ComplianceSensitive: true, EvidenceComplete: true, PolicyAllowed: true, Reason: "settlement"})
	if d.Disposition != MultiPartyReview { t.Fatalf("unexpected: %+v", d) }
}

func TestMissingEvidenceQuarantines(t *testing.T) {
	d := Route(Request{RiskScore: 5, Confidence: .99, PolicyAllowed: true, Reason: "missing evidence"})
	if d.Disposition != Quarantine { t.Fatalf("unexpected: %+v", d) }
}

func TestPolicyDenialCannotBeOverriddenByConfidence(t *testing.T) {
	d := Route(Request{RiskScore: 0, Confidence: 1, EvidenceComplete: true, PolicyAllowed: false, Reason: "model confident"})
	if d.Disposition != Deny { t.Fatalf("unexpected: %+v", d) }
}

func TestLowConfidenceEscalates(t *testing.T) {
	d := Route(Request{RiskScore: 10, Confidence: .70, EvidenceComplete: true, PolicyAllowed: true, Reason: "ambiguous"})
	if d.Disposition != HumanReview { t.Fatalf("unexpected: %+v", d) }
}
