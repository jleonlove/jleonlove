package inspectionintel

import (
	"errors"
	"testing"
	"time"
)

func TestEvaluatePassAndFailClosed(t *testing.T) {
	now := time.Now().UTC()
	good := Evidence{ID: "sgs-1", Provider: "SGS", Commodity: "gold", LotID: "lot-7", FacilityID: "ref-1", CertificateHash: "sha256:abc", ObservedAt: now.Add(-time.Hour), Purity: 99.95, Quantity: 200, Confidence: .98, DecisionUseAllowed: true, Independent: true, SignatureVerified: true, MaxAge: 24 * time.Hour}
	r := Request{Commodity: "gold", LotID: "lot-7", FacilityID: "ref-1", MinPurity: 99.5, MinQuantity: 200, MinConfidence: .9, At: now, RequireIndependent: true, RequireSignature: true}
	d, err := Evaluate(r, []Evidence{good})
	if err != nil || d.SelectedID != "sgs-1" {
		t.Fatalf("pass: %#v %v", d, err)
	}
	bad := good
	bad.ID = "bad"
	bad.SignatureVerified = false
	_, err = Evaluate(r, []Evidence{bad})
	if !errors.Is(err, ErrNoEligibleEvidence) {
		t.Fatalf("expected fail closed, got %v", err)
	}
	low := good
	low.ID = "low"
	low.Purity = 90.5
	_, err = Evaluate(r, []Evidence{low})
	if !errors.Is(err, ErrQualityFailed) {
		t.Fatalf("expected quality failure, got %v", err)
	}
}
