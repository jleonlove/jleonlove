package tradereconciliation

import "testing"

func base() (Record, Record, Record) {
 c := Record{Source:"contract", Quantity:500, Unit:"MT", Counterparty:"buyer-1", Currency:"USD", Amount:100000, Evidence:"contract-hash"}
 s := Record{Source:"shipment", Quantity:500, Unit:"MT", Counterparty:"buyer-1", Evidence:"bol-hash"}
 f := Record{Source:"financial", Counterparty:"buyer-1", Currency:"USD", Amount:100000, Evidence:"invoice-hash"}
 return c,s,f
}
func TestCleanThreeWayReconciliation(t *testing.T){ c,s,f:=base(); r:=Reconcile(c,s,f,0); if r.Status!="clean" || len(r.Breaks)!=0 { t.Fatalf("%+v",r) } }
func TestQuantityBreak(t *testing.T){ c,s,f:=base(); s.Quantity=485; r:=Reconcile(c,s,f,5); if r.Status=="clean" { t.Fatal("quantity break auto-cleared") } }
func TestMissingEvidenceFailsClosed(t *testing.T){ c,s,f:=base(); s.Evidence=""; r:=Reconcile(c,s,f,0); if r.Status=="clean" { t.Fatal("missing evidence auto-cleared") } }
func TestCounterpartyBreak(t *testing.T){ c,s,f:=base(); f.Counterparty="other"; r:=Reconcile(c,s,f,0); if r.Status=="clean" { t.Fatal("counterparty break auto-cleared") } }
func TestCurrencyBreak(t *testing.T){ c,s,f:=base(); f.Currency="EUR"; r:=Reconcile(c,s,f,0); if r.Status=="clean" { t.Fatal("currency break auto-cleared") } }
