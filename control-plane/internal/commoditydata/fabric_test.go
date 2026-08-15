package commoditydata

import (
	"testing"
	"time"
)

func base(now time.Time, source string, value float64) Observation {
	return Observation{Commodity: "sulfur", Metric: "FOB", Value: value, Unit: "MT", Currency: "USD", Source: source, SourceID: "q1", ObservedAt: now.Add(-time.Minute), MaxAge: 5 * time.Minute, Confidence: .95, License: LicensePolicy{AllowDecisionUse: true}, Evidence: []string{"signed-feed:" + source}}
}

func TestResolveFreshWithProvenance(t *testing.T) {
	now := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	f := New()
	if err := f.Ingest(base(now, "provider-a", 245)); err != nil {
		t.Fatal(err)
	}
	r, err := f.Resolve(Request{Commodity: "Sulfur", Metric: "fob", At: now, RequireDecisionUse: true, MinConfidence: .9})
	if err != nil || r.Observation.Value != 245 || len(r.Provenance) != 1 {
		t.Fatalf("bad resolution %#v %v", r, err)
	}
}
func TestRejectsStaleAndUnlicensed(t *testing.T) {
	now := time.Now().UTC()
	f := New()
	o := base(now, "a", 245)
	o.ObservedAt = now.Add(-10 * time.Minute)
	if err := f.Ingest(o); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Resolve(Request{Commodity: "sulfur", Metric: "FOB", At: now, RequireDecisionUse: true}); err == nil {
		t.Fatal("stale data must fail closed")
	}
	f = New()
	o = base(now, "a", 245)
	o.License.AllowDecisionUse = false
	_ = f.Ingest(o)
	if _, err := f.Resolve(Request{Commodity: "sulfur", Metric: "FOB", At: now, RequireDecisionUse: true}); err == nil {
		t.Fatal("unlicensed decision use must fail closed")
	}
}
func TestRejectsMaterialProviderConflict(t *testing.T) {
	now := time.Now().UTC()
	f := New()
	_ = f.Ingest(base(now, "a", 245))
	_ = f.Ingest(base(now, "b", 300))
	r, err := f.Resolve(Request{Commodity: "sulfur", Metric: "FOB", At: now, RequireDecisionUse: true, MaxRelativeConflict: .05})
	if err == nil || !r.Conflicting {
		t.Fatalf("expected conflict %#v %v", r, err)
	}
}
func TestValidation(t *testing.T) {
	f := New()
	if err := f.Ingest(Observation{}); err == nil {
		t.Fatal("invalid observation accepted")
	}
}
