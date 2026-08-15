package transactionverification

import (
	"testing"
	"time"
)

func ok(id string) Check {
	return Check{ID: id, Required: true, Passed: true, EvidenceRef: "ev:" + id, Source: "verified-source", VerifiedAt: time.Now()}
}
func base() Assessment {
	return Assessment{DealID: "deal-1", SubjectID: "seller-1", KYC: ok("kyc"), KYB: ok("kyb"), UBO: ok("ubo"), Mandate: ok("mandate"), Sanctions: ok("sanctions"), PEP: ok("pep"), AdverseMedia: ok("adverse_media"), ProofOfCommodity: ok("proof_of_commodity"), Facility: ok("facility"), DocumentIntegrity: ok("document_integrity"), BankOwnership: ok("bank_ownership")}
}
func TestPass(t *testing.T) {
	r, e := Evaluate(base(), time.Now())
	if e != nil || !CanExecute(r) || r.Score != 100 {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestMissingEvidenceReviews(t *testing.T) {
	a := base()
	a.Mandate.EvidenceRef = ""
	r, _ := Evaluate(a, time.Now())
	if r.Status != StatusReview || CanExecute(r) {
		t.Fatalf("%+v", r)
	}
}
func TestFailedCheckBlocks(t *testing.T) {
	a := base()
	a.Sanctions.Passed = false
	r, _ := Evaluate(a, time.Now())
	if r.Status != StatusBlock || CanExecute(r) {
		t.Fatalf("%+v", r)
	}
}
func TestStaleEvidenceReviews(t *testing.T) {
	a := base()
	a.KYB.ExpiresAt = time.Now().Add(-time.Minute)
	r, _ := Evaluate(a, time.Now())
	if r.Status != StatusReview {
		t.Fatalf("%+v", r)
	}
}
