package dealexecution

import (
	"atlas/internal/commoditydomain"
	"testing"
)

func engine(t *testing.T) *Engine {
	c := commoditydomain.New()
	if err := commoditydomain.SeedCore(c); err != nil {
		t.Fatal(err)
	}
	return New(c)
}
func deal() Deal {
	return Deal{ID: "D-1", Commodity: "sulphur", Quantity: 25000, Unit: "MT", Origin: "QA", Destination: "IN", Incoterm: "FOB", Transport: "bulk vessel", Payment: "LC"}
}
func TestCompileDomainRequirements(t *testing.T) {
	p, e := engine(t).Compile(deal())
	if e != nil {
		t.Fatal(e)
	}
	if p.CommodityID != "commodity.sulfur" {
		t.Fatal(p.CommodityID)
	}
	found := false
	for _, r := range p.Requirements {
		if r.Name == "certificate of analysis" {
			found = true
		}
	}
	if !found {
		t.Fatal("domain document missing")
	}
}
func TestInvalidDeal(t *testing.T) {
	_, e := engine(t).Compile(Deal{ID: "x", Commodity: "gold"})
	if e != ErrInvalidDeal {
		t.Fatal(e)
	}
}
func TestReadinessAndAdvance(t *testing.T) {
	p, _ := engine(t).Compile(deal())
	r := Evaluate(p, nil)
	if r.Ready || r.Total == 0 {
		t.Fatal(r)
	}
	var ev []Evidence
	for _, q := range p.Requirements {
		if q.Required {
			ev = append(ev, Evidence{q.ID, "evidence:" + q.ID, true})
		}
	}
	r = Evaluate(p, ev)
	if !r.Ready || r.Completed != r.Total {
		t.Fatal(r)
	}
	if e := CanAdvance(p, Settlement, ev); e != nil {
		t.Fatal(e)
	}
}
func TestUnverifiedBlocks(t *testing.T) {
	p, _ := engine(t).Compile(deal())
	var ev []Evidence
	for _, q := range p.Requirements {
		if q.Required {
			ev = append(ev, Evidence{q.ID, "ref", q.ID != "sanctions"})
		}
	}
	if Evaluate(p, ev).Ready {
		t.Fatal("must block")
	}
	if e := CanAdvance(p, Compliance, ev); e == nil {
		t.Fatal("expected block")
	}
}
