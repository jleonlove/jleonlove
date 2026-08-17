package protocolpolicy

import "fmt"

type Decision struct {
	Allowed           bool
	NegotiatedVersion string
	Reason            string
}

type Policy struct {
	AllowedVersions map[string]bool
	AllowDowngrade  bool
}

// Authorize is deny-by-default. A protocol version must be explicitly admitted.
// Downgrades are never silent and require both policy permission and caller consent.
func (p Policy) Authorize(requested, negotiated string, callerAcceptedDowngrade bool) Decision {
	if requested == "" || negotiated == "" {
		return Decision{Reason: "protocol_version_required"}
	}
	if !p.AllowedVersions[negotiated] {
		return Decision{NegotiatedVersion: negotiated, Reason: "protocol_version_not_allowed"}
	}
	if requested != negotiated && (!p.AllowDowngrade || !callerAcceptedDowngrade) {
		return Decision{NegotiatedVersion: negotiated, Reason: "protocol_downgrade_denied"}
	}
	return Decision{Allowed: true, NegotiatedVersion: negotiated, Reason: "allowed"}
}

func (d Decision) Evidence() string {
	return fmt.Sprintf("allowed=%t version=%s reason=%s", d.Allowed, d.NegotiatedVersion, d.Reason)
}
