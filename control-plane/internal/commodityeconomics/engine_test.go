package commodityeconomics

import "testing"

func TestEvaluateAndStress(t *testing.T) {
	c := Case{ID: "deal-1", Commodity: "sulfur", Currency: "USD", Quantity: 10000, PurchasePerUnit: 100, SalePerUnit: 145, Costs: []Cost{{Name: "inspection", PerUnit: 1}, {Name: "insurance", Fixed: 10000}}, Route: []RouteLeg{{Mode: "vessel", Origin: "A", Destination: "B", CostPerUnit: 15, Capacity: 20000}}, WorkingCapitalDays: 45, AnnualFinanceRate: .10, LossPct: 1}
	r, err := Evaluate(c)
	if err != nil {
		t.Fatal(err)
	}
	if !r.RouteFeasible || r.GrossProfit <= 0 || r.BreakEvenSalePerUnit <= 100 {
		t.Fatalf("bad result: %+v", r)
	}
	s, err := Stress(c, []Scenario{{Name: "downside", PurchaseShockPct: 10, SaleShockPct: -10, CostShockPct: 20, LossPctDelta: 1}})
	if err != nil || len(s) != 1 {
		t.Fatal(err)
	}
	if s[0].Result.GrossProfit >= r.GrossProfit {
		t.Fatal("stress should reduce profit")
	}
}
func TestCapacityFlag(t *testing.T) {
	c := Case{ID: "x", Commodity: "wheat", Currency: "USD", Quantity: 100, PurchasePerUnit: 1, SalePerUnit: 2, Route: []RouteLeg{{Mode: "truck", Origin: "a", Destination: "b", Capacity: 50}}}
	r, _ := Evaluate(c)
	if r.RouteFeasible {
		t.Fatal("expected infeasible")
	}
}
