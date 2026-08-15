package freightintel

import (
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	n := time.Now().UTC()
	r := Request{Origin: "Houston", Destination: "Rotterdam", Mode: "ocean", Quantity: 1000, At: n, MinConfidence: .8}
	qs := []Quote{{Provider: "stale", Origin: "Houston", Destination: "Rotterdam", Mode: "ocean", Currency: "USD", CostPerUnit: 10, Confidence: .99, ObservedAt: n.Add(-48 * time.Hour), MaxAge: time.Hour, DecisionUseAllowed: true}, {Provider: "good", Origin: "Houston", Destination: "Rotterdam", Mode: "ocean", Currency: "USD", CostPerUnit: 25, Confidence: .9, Capacity: 2000, ObservedAt: n, MaxAge: time.Hour, DecisionUseAllowed: true}, {Provider: "small", Origin: "Houston", Destination: "Rotterdam", Mode: "ocean", Currency: "USD", CostPerUnit: 20, Confidence: .9, Capacity: 500, ObservedAt: n, MaxAge: time.Hour, DecisionUseAllowed: true}}
	d, e := Select(r, qs)
	if e != nil || d.Quote.Provider != "good" || d.TotalFreight != 25000 || len(d.Rejected) != 2 {
		t.Fatalf("%+v %v", d, e)
	}
}
func TestUnlicensedFailsClosed(t *testing.T) {
	n := time.Now().UTC()
	_, e := Select(Request{Origin: "A", Destination: "B", Quantity: 1, At: n}, []Quote{{Provider: "p", Origin: "A", Destination: "B", Mode: "truck", Currency: "USD", CostPerUnit: 1, Confidence: 1, ObservedAt: n, MaxAge: time.Hour}})
	if e != ErrNoEligibleQuote {
		t.Fatal(e)
	}
}
