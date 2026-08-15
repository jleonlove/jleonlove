package dealexecution

import (
	"errors"
	"sort"
	"strings"

	"atlas/internal/commoditydomain"
)

var (
	ErrInvalidDeal        = errors.New("invalid deal")
	ErrMissingRequirement = errors.New("required transaction item missing")
	ErrStageBlocked       = errors.New("deal stage blocked")
)

type Stage string

const (
	Qualification Stage = "QUALIFICATION"
	Verification  Stage = "VERIFICATION"
	Commercial    Stage = "COMMERCIAL"
	Compliance    Stage = "COMPLIANCE"
	Finance       Stage = "FINANCE"
	Inspection    Stage = "INSPECTION"
	Logistics     Stage = "LOGISTICS"
	Delivery      Stage = "DELIVERY"
	Settlement    Stage = "SETTLEMENT"
	Closeout      Stage = "CLOSEOUT"
)

type Deal struct {
	ID, Commodity, Origin, Destination, Incoterm, Transport, Payment string
	Quantity                                                         float64
	Unit                                                             string
}
type Requirement struct {
	ID, Name, Kind, Stage, Source string
	Required                      bool
}
type Plan struct {
	Deal         Deal
	CommodityID  string
	Requirements []Requirement
	Stages       []Stage
}
type Evidence struct {
	RequirementID, Reference string
	Verified                 bool
}
type Readiness struct {
	Ready               bool
	Missing, Unverified []string
	Completed           int
	Total               int
}

type Engine struct{ Domains *commoditydomain.Compiler }

func New(domains *commoditydomain.Compiler) *Engine { return &Engine{Domains: domains} }
func valid(d Deal) bool {
	return strings.TrimSpace(d.ID) != "" && strings.TrimSpace(d.Commodity) != "" && d.Quantity > 0 && strings.TrimSpace(d.Unit) != "" && strings.TrimSpace(d.Origin) != "" && strings.TrimSpace(d.Destination) != ""
}
func rid(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	r := strings.NewReplacer(" ", "_", "/", "_", "-", "_")
	return r.Replace(s)
}
func (e *Engine) Compile(d Deal) (Plan, error) {
	if !valid(d) {
		return Plan{}, ErrInvalidDeal
	}
	p, err := e.Domains.Resolve(d.Commodity)
	if err != nil {
		return Plan{}, err
	}
	req := []Requirement{
		{"kyb_buyer", "Buyer KYB/UBO verification", "verification", string(Verification), "atlas", true},
		{"kyb_seller", "Seller KYB/UBO verification", "verification", string(Verification), "atlas", true},
		{"mandate", "Representative mandate/authority verification", "verification", string(Verification), "atlas", true},
		{"commercial_terms", "Agreed commercial terms", "commercial", string(Commercial), "deal", true},
		{"sale_contract", "Sale/purchase agreement", "document", string(Commercial), "deal", true},
		{"sanctions", "Sanctions/compliance clearance", "compliance", string(Compliance), "atlas", true},
		{"payment", "Payment/finance readiness", "finance", string(Finance), "deal", true},
		{"inspection", "Independent quality/quantity inspection", "inspection", string(Inspection), "deal", true},
		{"transport", "Transport/logistics confirmation", "logistics", string(Logistics), "deal", true},
		{"delivery_evidence", "Delivery evidence", "document", string(Delivery), "deal", true},
		{"settlement_authorization", "Settlement authorization", "settlement", string(Settlement), "atlas", true},
	}
	for _, dr := range p.Pack.Documents {
		req = append(req, Requirement{"domain_" + rid(dr.Name), dr.Name, "document", strings.ToUpper(dr.Stage), "domain-pack", dr.Required})
	}
	seen := map[string]bool{}
	out := make([]Requirement, 0, len(req))
	for _, r := range req {
		if !seen[r.ID] {
			seen[r.ID] = true
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Stage == out[j].Stage {
			return out[i].ID < out[j].ID
		}
		return out[i].Stage < out[j].Stage
	})
	return Plan{Deal: d, CommodityID: p.Pack.ID, Requirements: out, Stages: []Stage{Qualification, Verification, Commercial, Compliance, Finance, Inspection, Logistics, Delivery, Settlement, Closeout}}, nil
}
func Evaluate(p Plan, ev []Evidence) Readiness {
	m := map[string]Evidence{}
	for _, x := range ev {
		m[x.RequirementID] = x
	}
	r := Readiness{}
	for _, q := range p.Requirements {
		if !q.Required {
			continue
		}
		r.Total++
		x, ok := m[q.ID]
		if !ok || strings.TrimSpace(x.Reference) == "" {
			r.Missing = append(r.Missing, q.ID)
			continue
		}
		if !x.Verified {
			r.Unverified = append(r.Unverified, q.ID)
			continue
		}
		r.Completed++
	}
	sort.Strings(r.Missing)
	sort.Strings(r.Unverified)
	r.Ready = len(r.Missing) == 0 && len(r.Unverified) == 0
	return r
}
func CanAdvance(p Plan, target Stage, ev []Evidence) error {
	idx := -1
	for i, s := range p.Stages {
		if s == target {
			idx = i
			break
		}
	}
	if idx < 0 {
		return ErrStageBlocked
	}
	allowed := map[string]bool{}
	for i := 0; i <= idx; i++ {
		allowed[string(p.Stages[i])] = true
	}
	m := map[string]Evidence{}
	for _, x := range ev {
		m[x.RequirementID] = x
	}
	for _, q := range p.Requirements {
		if !q.Required || !allowed[q.Stage] {
			continue
		}
		x, ok := m[q.ID]
		if !ok || x.Reference == "" || !x.Verified {
			return ErrMissingRequirement
		}
	}
	return nil
}
