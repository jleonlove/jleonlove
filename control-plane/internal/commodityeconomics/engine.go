package commodityeconomics

import (
	"errors"
	"math"
	"sort"
	"strings"
)

var ErrInvalidCase = errors.New("invalid commodity economics case")

type Cost struct {
	Name    string
	PerUnit float64
	Fixed   float64
}
type RouteLeg struct {
	Mode, Origin, Destination string
	DistanceKM, CostPerUnit   float64
	Capacity                  float64
}
type Case struct {
	ID, Commodity, Currency                string
	Quantity, PurchasePerUnit, SalePerUnit float64
	Costs                                  []Cost
	Route                                  []RouteLeg
	WorkingCapitalDays                     int
	AnnualFinanceRate                      float64
	LossPct                                float64
}
type Result struct {
	Revenue, PurchaseCost, VariableCosts, FixedCosts, LogisticsCost, FinanceCost, LossCost, LandedCost, GrossProfit float64
	MarginPct, BreakEvenSalePerUnit, ProfitPerUnit                                                                  float64
	RouteFeasible                                                                                                   bool
	Flags                                                                                                           []string
}
type Scenario struct {
	Name                                                       string
	PurchaseShockPct, SaleShockPct, CostShockPct, LossPctDelta float64
}
type ScenarioResult struct {
	Scenario string
	Result   Result
}

func valid(c Case) bool {
	return strings.TrimSpace(c.ID) != "" && strings.TrimSpace(c.Commodity) != "" && strings.TrimSpace(c.Currency) != "" && c.Quantity > 0 && c.PurchasePerUnit >= 0 && c.SalePerUnit >= 0 && c.LossPct >= 0 && c.LossPct < 100 && c.WorkingCapitalDays >= 0 && c.AnnualFinanceRate >= 0
}
func round(v float64) float64 { return math.Round(v*100) / 100 }
func Evaluate(c Case) (Result, error) {
	if !valid(c) {
		return Result{}, ErrInvalidCase
	}
	r := Result{RouteFeasible: true}
	r.Revenue = c.Quantity * c.SalePerUnit
	r.PurchaseCost = c.Quantity * c.PurchasePerUnit
	for _, x := range c.Costs {
		if x.PerUnit < 0 || x.Fixed < 0 {
			return Result{}, ErrInvalidCase
		}
		r.VariableCosts += c.Quantity * x.PerUnit
		r.FixedCosts += x.Fixed
	}
	for _, l := range c.Route {
		if l.CostPerUnit < 0 || l.DistanceKM < 0 || l.Capacity < 0 {
			return Result{}, ErrInvalidCase
		}
		r.LogisticsCost += c.Quantity * l.CostPerUnit
		if l.Capacity > 0 && l.Capacity < c.Quantity {
			r.RouteFeasible = false
			r.Flags = append(r.Flags, "ROUTE_CAPACITY")
		}
		if strings.TrimSpace(l.Origin) == "" || strings.TrimSpace(l.Destination) == "" || strings.TrimSpace(l.Mode) == "" {
			r.RouteFeasible = false
			r.Flags = append(r.Flags, "ROUTE_INCOMPLETE")
		}
	}
	r.LossCost = r.PurchaseCost * (c.LossPct / 100)
	financed := r.PurchaseCost + r.VariableCosts + r.FixedCosts + r.LogisticsCost
	r.FinanceCost = financed * c.AnnualFinanceRate * float64(c.WorkingCapitalDays) / 365
	r.LandedCost = r.PurchaseCost + r.VariableCosts + r.FixedCosts + r.LogisticsCost + r.LossCost + r.FinanceCost
	r.GrossProfit = r.Revenue - r.LandedCost
	if r.Revenue > 0 {
		r.MarginPct = r.GrossProfit / r.Revenue * 100
	}
	r.ProfitPerUnit = r.GrossProfit / c.Quantity
	r.BreakEvenSalePerUnit = r.LandedCost / c.Quantity
	if r.GrossProfit < 0 {
		r.Flags = append(r.Flags, "NEGATIVE_MARGIN")
	}
	if c.WorkingCapitalDays > 180 {
		r.Flags = append(r.Flags, "LONG_WORKING_CAPITAL")
	}
	sort.Strings(r.Flags)
	r.Revenue = round(r.Revenue)
	r.PurchaseCost = round(r.PurchaseCost)
	r.VariableCosts = round(r.VariableCosts)
	r.FixedCosts = round(r.FixedCosts)
	r.LogisticsCost = round(r.LogisticsCost)
	r.FinanceCost = round(r.FinanceCost)
	r.LossCost = round(r.LossCost)
	r.LandedCost = round(r.LandedCost)
	r.GrossProfit = round(r.GrossProfit)
	r.MarginPct = round(r.MarginPct)
	r.ProfitPerUnit = round(r.ProfitPerUnit)
	r.BreakEvenSalePerUnit = round(r.BreakEvenSalePerUnit)
	return r, nil
}
func Stress(c Case, scenarios []Scenario) ([]ScenarioResult, error) {
	out := make([]ScenarioResult, 0, len(scenarios))
	for _, s := range scenarios {
		x := c
		x.PurchasePerUnit *= 1 + s.PurchaseShockPct/100
		x.SalePerUnit *= 1 + s.SaleShockPct/100
		x.LossPct += s.LossPctDelta
		for i := range x.Costs {
			x.Costs[i].PerUnit *= 1 + s.CostShockPct/100
			x.Costs[i].Fixed *= 1 + s.CostShockPct/100
		}
		for i := range x.Route {
			x.Route[i].CostPerUnit *= 1 + s.CostShockPct/100
		}
		r, err := Evaluate(x)
		if err != nil {
			return nil, err
		}
		out = append(out, ScenarioResult{Scenario: s.Name, Result: r})
	}
	return out, nil
}
