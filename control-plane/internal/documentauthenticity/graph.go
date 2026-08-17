package documentauthenticity

import (
	"fmt"
	"strings"
)

type Document struct {
	ID string
	Type string
	ContentHash string
	Issuer string
	Counterparty string
	Commodity string
	Vessel string
	Container string
	Origin string
	Destination string
	Quantity string
	IssueDate string
	ExternalVerification bool
}

type Finding struct {
	Code string
	Severity string
	DocumentIDs []string
	Detail string
}

type Report struct {
	Clear bool
	Findings []Finding
}

// Evaluate detects duplicate/recycled documents and cross-document identity inconsistencies.
// It is intentionally conservative: missing hashes and unverifiable issuers cannot clear silently.
func Evaluate(docs []Document) Report {
	findings := []Finding{}
	hashes := map[string]string{}
	containers := map[string]string{}
	for _, d := range docs {
		if strings.TrimSpace(d.ID) == "" || strings.TrimSpace(d.ContentHash) == "" {
			findings = append(findings, Finding{Code:"DOCUMENT_IDENTITY_INCOMPLETE", Severity:"high", DocumentIDs:[]string{d.ID}})
			continue
		}
		if prior, ok := hashes[d.ContentHash]; ok && prior != d.ID {
			findings = append(findings, Finding{Code:"DUPLICATE_DOCUMENT_HASH", Severity:"critical", DocumentIDs:[]string{prior,d.ID}, Detail:"same content presented as distinct documents"})
		} else { hashes[d.ContentHash] = d.ID }
		if d.Container != "" {
			if prior, ok := containers[d.Container]; ok && prior != d.ID {
				findings = append(findings, Finding{Code:"REUSED_CONTAINER_REFERENCE", Severity:"high", DocumentIDs:[]string{prior,d.ID}, Detail:fmt.Sprintf("container %s reused", d.Container)})
			} else { containers[d.Container] = d.ID }
		}
		if strings.TrimSpace(d.Issuer) == "" || !d.ExternalVerification {
			findings = append(findings, Finding{Code:"ISSUER_NOT_VERIFIED", Severity:"high", DocumentIDs:[]string{d.ID}})
		}
	}
	if len(docs) > 1 {
		base := docs[0]
		for _, d := range docs[1:] {
			check := func(field, a, b string) { if a != "" && b != "" && !strings.EqualFold(a,b) { findings = append(findings, Finding{Code:"CROSS_DOCUMENT_"+field+"_MISMATCH", Severity:"high", DocumentIDs:[]string{base.ID,d.ID}, Detail:a+" != "+b}) } }
			check("COUNTERPARTY", base.Counterparty, d.Counterparty)
			check("COMMODITY", base.Commodity, d.Commodity)
			check("VESSEL", base.Vessel, d.Vessel)
			check("ORIGIN", base.Origin, d.Origin)
			check("DESTINATION", base.Destination, d.Destination)
		}
	}
	return Report{Clear: len(findings)==0, Findings:findings}
}
