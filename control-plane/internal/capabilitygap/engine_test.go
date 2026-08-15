package capabilitygap

import "testing"

func TestRanksAndFailClosed(t *testing.T) {
	e := New()
	e.Observe(Signal{Domain: "wheat", MissingCapability: "live freight", Severity: 8, Frequency: 5, Evidence: []string{"deal-1"}})
	e.Observe(Signal{Domain: "wheat", MissingCapability: "live freight", Severity: 5, Frequency: 3, Evidence: []string{"deal-2"}})
	p := e.Proposals()
	if len(p) != 1 || p[0].Priority != 0 || p[0].Score != 55 {
		t.Fatalf("bad proposal %#v", p)
	}
	if !p[0].RequiresHumanApproval || e.CanSelfDeploy() {
		t.Fatal("evolution must fail closed")
	}
}
