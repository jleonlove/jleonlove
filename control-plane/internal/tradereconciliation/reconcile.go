package tradereconciliation

import (
 "fmt"
 "math"
 "strings"
)

type Record struct {
 Source string
 Quantity float64
 Unit string
 Counterparty string
 Currency string
 Amount float64
 Evidence string
}

type Result struct {
 Status string
 Breaks []string
 Evidence []string
}

// Reconcile compares contract, shipment/document, and financial/system truth.
// Missing evidence never qualifies for straight-through processing.
func Reconcile(contract, shipment, financial Record, quantityTolerance float64) Result {
 r := Result{Status: "review"}
 for _, x := range []Record{contract, shipment, financial} {
  if strings.TrimSpace(x.Evidence) == "" { r.Breaks = append(r.Breaks, fmt.Sprintf("missing_evidence:%s", x.Source)) } else { r.Evidence = append(r.Evidence, x.Evidence) }
 }
 if contract.Unit == "" || shipment.Unit == "" || contract.Unit != shipment.Unit { r.Breaks = append(r.Breaks, "unit_mismatch") }
 if math.Abs(contract.Quantity-shipment.Quantity) > quantityTolerance { r.Breaks = append(r.Breaks, "quantity_break") }
 if contract.Counterparty == "" || shipment.Counterparty == "" || financial.Counterparty == "" || contract.Counterparty != shipment.Counterparty || contract.Counterparty != financial.Counterparty { r.Breaks = append(r.Breaks, "counterparty_break") }
 if contract.Currency == "" || financial.Currency == "" || contract.Currency != financial.Currency { r.Breaks = append(r.Breaks, "currency_break") }
 if contract.Amount > 0 && financial.Amount > 0 && contract.Amount != financial.Amount { r.Breaks = append(r.Breaks, "amount_break") }
 if len(r.Breaks) == 0 { r.Status = "clean" }
 return r
}
