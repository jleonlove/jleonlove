package shipmentplausibility

import "testing"

func valid() Shipment { return Shipment{DeclaredVessel:"MV ATLAS",ObservedVessel:"MV ATLAS",DeclaredOrigin:"Houston",ObservedLoadPort:"Houston",DeclaredDestination:"Rotterdam",ObservedDischargePort:"Rotterdam",CustodyEvidenceComplete:true,ObservationEvidencePresent:true} }

func TestValidShipmentPasses(t *testing.T) { d:=Evaluate(valid()); if !d.Plausible || d.Quarantine { t.Fatalf("unexpected: %+v",d) } }
func TestVesselMismatchQuarantines(t *testing.T) { s:=valid(); s.ObservedVessel="MV OTHER"; d:=Evaluate(s); if !d.Quarantine { t.Fatalf("unexpected: %+v",d) } }
func TestUnauthorizedRouteChangeQuarantines(t *testing.T) { s:=valid(); s.RouteChanged=true; d:=Evaluate(s); if !d.Quarantine { t.Fatalf("unexpected: %+v",d) } }
func TestAuthorizedRouteChangeCanPass(t *testing.T) { s:=valid(); s.RouteChanged=true; s.RouteChangeAuthorized=true; d:=Evaluate(s); if !d.Plausible { t.Fatalf("unexpected: %+v",d) } }
func TestUnexpectedTransshipmentQuarantines(t *testing.T) { s:=valid(); s.UnexpectedTransshipments=1; d:=Evaluate(s); if !d.Quarantine { t.Fatalf("unexpected: %+v",d) } }
func TestMissingCustodyEvidenceQuarantines(t *testing.T) { s:=valid(); s.CustodyEvidenceComplete=false; d:=Evaluate(s); if !d.Quarantine { t.Fatalf("unexpected: %+v",d) } }
func TestMissingObservationEvidenceQuarantines(t *testing.T) { s:=valid(); s.ObservationEvidencePresent=false; d:=Evaluate(s); if !d.Quarantine { t.Fatalf("unexpected: %+v",d) } }
