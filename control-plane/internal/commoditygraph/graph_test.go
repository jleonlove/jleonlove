package commoditygraph

import "testing"

func fixture(t *testing.T) *Graph {
	g, e := New([]Node{{ID: "gas", Kind: "commodity", Name: "Natural Gas"}, {ID: "ammonia", Kind: "product", Name: "Ammonia"}, {ID: "urea", Kind: "commodity", Name: "Urea"}, {ID: "wheat", Kind: "commodity", Name: "Wheat"}}, []Edge{{From: "gas", To: "ammonia", Relation: "feedstock_cost", Weight: .8, Evidence: []string{"process"}}, {From: "ammonia", To: "urea", Relation: "input_cost", Weight: .7}, {From: "urea", To: "wheat", Relation: "fertilizer_cost", Weight: .25}})
	if e != nil {
		t.Fatal(e)
	}
	return g
}
func TestPaths(t *testing.T) {
	g := fixture(t)
	p := g.FindPaths("gas", "wheat", 4)
	if len(p) != 1 || len(p[0].Relations) != 3 {
		t.Fatalf("bad path %#v", p)
	}
}
func TestPropagation(t *testing.T) {
	g := fixture(t)
	x, e := g.Propagate(Shock{"gas", 40}, 4, .01)
	if e != nil {
		t.Fatal(e)
	}
	if len(x) != 3 {
		t.Fatalf("want 3 impacts got %d", len(x))
	}
	if x[2].NodeID != "wheat" || x[2].ChangePct <= 0 {
		t.Fatalf("bad impact %#v", x[2])
	}
}
func TestValidation(t *testing.T) {
	_, e := New([]Node{{ID: "a", Kind: "commodity", Name: "A"}}, []Edge{{From: "a", To: "missing", Relation: "x", Weight: 1}})
	if e == nil {
		t.Fatal("expected error")
	}
}
