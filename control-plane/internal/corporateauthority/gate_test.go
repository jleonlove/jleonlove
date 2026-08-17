package corporateauthority

import "testing"

func TestProtectedActionsFailClosed(t *testing.T) {
	for _, action := range []Action{OwnershipChange, EquityIssuance, IPTransfer, ProductionRootTransfer, RecoveryRootTransfer, DomainTransfer, PolicyRootChange, EvidenceDestruction} {
		d := Authorize(Request{Action: action})
		if d.Allowed { t.Fatalf("%s unexpectedly allowed", action) }
	}
}

func TestAuthorityAloneCannotApprove(t *testing.T) {
	d := Authorize(Request{Action: IPTransfer, AuthorityEvidence: "signed-owner-authority", ExactEffectDigest: "sha256:effect"})
	if d.Allowed || d.Reason != "independent_approval_required" { t.Fatalf("unexpected: %+v", d) }
}

func TestApprovalCannotBeReusedWithoutExactEffectBinding(t *testing.T) {
	d := Authorize(Request{Action: DomainTransfer, AuthorityEvidence: "signed-owner-authority", IndependentApproval: "signed-independent-approval"})
	if d.Allowed || d.Reason != "exact_effect_binding_required" { t.Fatalf("unexpected: %+v", d) }
}

func TestProtectedActionRequiresCompleteEvidence(t *testing.T) {
	d := Authorize(Request{Action: PolicyRootChange, AuthorityEvidence: "signed-owner-authority", IndependentApproval: "signed-independent-approval", ExactEffectDigest: "sha256:effect"})
	if !d.Allowed { t.Fatalf("unexpected: %+v", d) }
}

func TestUnknownCorporateActionDenied(t *testing.T) {
	d := Authorize(Request{Action: Action("unknown"), AuthorityEvidence: "x", IndependentApproval: "y", ExactEffectDigest: "z"})
	if d.Allowed || d.Reason != "action_not_admitted" { t.Fatalf("unexpected: %+v", d) }
}
