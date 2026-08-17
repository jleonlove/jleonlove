package shipmentplausibility

import "strings"

type Shipment struct {
	DeclaredVessel string
	ObservedVessel string
	DeclaredOrigin string
	ObservedLoadPort string
	DeclaredDestination string
	ObservedDischargePort string
	RouteChanged bool
	RouteChangeAuthorized bool
	UnexpectedTransshipments int
	CustodyEvidenceComplete bool
	ObservationEvidencePresent bool
}

type Decision struct {
	Plausible bool
	Quarantine bool
	Reasons []string
}

// Evaluate never treats paperwork as sufficient proof of shipment reality.
// Material route/custody contradictions quarantine the shipment for review.
func Evaluate(s Shipment) Decision {
	reasons := []string{}
	if !s.ObservationEvidencePresent { reasons = append(reasons, "observation_evidence_missing") }
	if strings.TrimSpace(s.DeclaredVessel) == "" || strings.TrimSpace(s.ObservedVessel) == "" || !strings.EqualFold(s.DeclaredVessel, s.ObservedVessel) { reasons = append(reasons, "vessel_mismatch") }
	if strings.TrimSpace(s.DeclaredOrigin) == "" || strings.TrimSpace(s.ObservedLoadPort) == "" || !strings.EqualFold(s.DeclaredOrigin, s.ObservedLoadPort) { reasons = append(reasons, "origin_mismatch") }
	if strings.TrimSpace(s.DeclaredDestination) == "" || strings.TrimSpace(s.ObservedDischargePort) == "" || !strings.EqualFold(s.DeclaredDestination, s.ObservedDischargePort) { reasons = append(reasons, "destination_mismatch") }
	if s.RouteChanged && !s.RouteChangeAuthorized { reasons = append(reasons, "unauthorized_route_change") }
	if s.UnexpectedTransshipments > 0 { reasons = append(reasons, "unexpected_transshipment") }
	if !s.CustodyEvidenceComplete { reasons = append(reasons, "custody_evidence_incomplete") }
	return Decision{Plausible: len(reasons) == 0, Quarantine: len(reasons) > 0, Reasons: reasons}
}
