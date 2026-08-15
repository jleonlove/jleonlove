package traderegulatory

import (
	"errors"
	"testing"
	"time"
)

func TestEvaluateDutyAndDocs(t *testing.T) {
	now := time.Now().UTC()
	d, e := Evaluate(Request{Origin: "US", Destination: "CA", Commodity: "wheat", Quantity: 100, CustomsValue: 10000, At: now}, []Rule{{ID: "r1", Source: "customs", Origin: "US", Destination: "CA", Commodity: "wheat", Currency: "USD", DutyRate: .05, RequiredDocuments: []string{"origin_certificate"}, Confidence: .99, ObservedAt: now, MaxAge: 24 * time.Hour, DecisionUseAllowed: true}})
	if e != nil || d.Duty != 500 || len(d.RequiredDocuments) != 1 {
		t.Fatalf("bad decision %#v %v", d, e)
	}
}
func TestBlockedFailsClosed(t *testing.T) {
	now := time.Now().UTC()
	_, e := Evaluate(Request{Origin: "A", Destination: "B", Commodity: "gold", Quantity: 1, CustomsValue: 1, At: now}, []Rule{{ID: "ban", Source: "reg", Origin: "A", Destination: "B", Commodity: "gold", Prohibited: true, Confidence: 1, ObservedAt: now, MaxAge: time.Hour, DecisionUseAllowed: true}})
	if !errors.Is(e, ErrTradeBlocked) {
		t.Fatalf("expected block: %v", e)
	}
}
func TestStaleRuleRejected(t *testing.T) {
	now := time.Now().UTC()
	d, e := Evaluate(Request{Origin: "A", Destination: "B", Commodity: "x", Quantity: 1, At: now}, []Rule{{ID: "old", Source: "reg", Origin: "A", Destination: "B", Commodity: "x", Confidence: 1, ObservedAt: now.Add(-48 * time.Hour), MaxAge: time.Hour, DecisionUseAllowed: true}})
	if !errors.Is(e, ErrNoApplicableRule) || d.Rejected["old"] != "stale_or_future" {
		t.Fatalf("expected stale rejection %#v %v", d, e)
	}
}
