package documentauthenticity

import "testing"

func good(id, hash string) Document { return Document{ID:id, Type:"BOL", ContentHash:hash, Issuer:"verified-carrier", Counterparty:"buyer-a", Commodity:"wheat", Vessel:"vessel-a", Origin:"US", Destination:"NL", ExternalVerification:true} }

func TestCleanDocumentsClear(t *testing.T) {
	r := Evaluate([]Document{good("d1","h1"), good("d2","h2")})
	if !r.Clear { t.Fatalf("unexpected findings: %+v", r.Findings) }
}

func TestDuplicateHashBlocked(t *testing.T) {
	r := Evaluate([]Document{good("d1","same"), good("d2","same")})
	if r.Clear { t.Fatal("duplicate document cleared") }
}

func TestUnverifiedIssuerBlocked(t *testing.T) {
	d := good("d1","h1"); d.ExternalVerification=false
	if Evaluate([]Document{d}).Clear { t.Fatal("unverified issuer cleared") }
}

func TestCounterpartyMismatchBlocked(t *testing.T) {
	a,b := good("d1","h1"),good("d2","h2"); b.Counterparty="buyer-b"
	if Evaluate([]Document{a,b}).Clear { t.Fatal("counterparty mismatch cleared") }
}

func TestReusedContainerDetected(t *testing.T) {
	a,b := good("d1","h1"),good("d2","h2"); a.Container="C1"; b.Container="C1"
	if Evaluate([]Document{a,b}).Clear { t.Fatal("reused container cleared") }
}
