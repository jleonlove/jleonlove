package protocolpolicy

import "testing"

func TestDenyUnknownVersion(t *testing.T) {
	p := Policy{AllowedVersions: map[string]bool{"a2a/1.1": true}}
	d := p.Authorize("a2a/1.1", "a2a/9.9", false)
	if d.Allowed || d.Reason != "protocol_version_not_allowed" { t.Fatalf("unexpected decision: %+v", d) }
}

func TestDenySilentDowngrade(t *testing.T) {
	p := Policy{AllowedVersions: map[string]bool{"a2a/1.1": true, "a2a/0.3": true}, AllowDowngrade: true}
	d := p.Authorize("a2a/1.1", "a2a/0.3", false)
	if d.Allowed || d.Reason != "protocol_downgrade_denied" { t.Fatalf("unexpected decision: %+v", d) }
}

func TestAllowExplicitVersion(t *testing.T) {
	p := Policy{AllowedVersions: map[string]bool{"a2a/1.1": true}}
	d := p.Authorize("a2a/1.1", "a2a/1.1", false)
	if !d.Allowed || d.NegotiatedVersion != "a2a/1.1" { t.Fatalf("unexpected decision: %+v", d) }
}

func TestAllowExplicitDowngradeOnlyWithConsent(t *testing.T) {
	p := Policy{AllowedVersions: map[string]bool{"a2a/1.1": true, "a2a/0.3": true}, AllowDowngrade: true}
	d := p.Authorize("a2a/1.1", "a2a/0.3", true)
	if !d.Allowed { t.Fatalf("unexpected decision: %+v", d) }
}
