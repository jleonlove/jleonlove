package traderegulatory

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var ErrInvalidRequest = errors.New("invalid regulatory request")
var ErrNoApplicableRule = errors.New("no applicable regulatory rule")
var ErrTradeBlocked = errors.New("trade blocked by regulatory rule")

type Rule struct {
	ID, Source, Origin, Destination, Commodity, HSCode, Currency string
	EffectiveFrom, EffectiveTo, ObservedAt                       time.Time
	DutyRate, FixedDutyPerUnit, Confidence                       float64
	RequiredDocuments                                            []string
	Prohibited, DecisionUseAllowed                               bool
	MaxAge                                                       time.Duration
}
type Request struct {
	Origin, Destination, Commodity, HSCode string
	Quantity, CustomsValue                 float64
	At                                     time.Time
	MinConfidence                          float64
}
type Decision struct {
	RuleIDs           []string
	Duty              float64
	Currency          string
	RequiredDocuments []string
	Rejected          map[string]string
}

func Evaluate(r Request, rules []Rule) (Decision, error) {
	if strings.TrimSpace(r.Origin) == "" || strings.TrimSpace(r.Destination) == "" || strings.TrimSpace(r.Commodity) == "" || r.Quantity <= 0 || r.CustomsValue < 0 {
		return Decision{}, ErrInvalidRequest
	}
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	if r.MinConfidence <= 0 {
		r.MinConfidence = .8
	}
	d := Decision{Rejected: map[string]string{}}
	docs := map[string]bool{}
	applicable := []Rule{}
	for _, x := range rules {
		key := x.ID
		if key == "" {
			key = x.Source
		}
		why := ""
		switch {
		case x.ID == "" || x.Source == "" || x.ObservedAt.IsZero() || x.Confidence < 0 || x.Confidence > 1 || x.DutyRate < 0 || x.FixedDutyPerUnit < 0:
			why = "invalid"
		case !strings.EqualFold(x.Origin, r.Origin) || !strings.EqualFold(x.Destination, r.Destination) || !strings.EqualFold(x.Commodity, r.Commodity):
			why = "scope_mismatch"
		case r.HSCode != "" && x.HSCode != "" && !strings.EqualFold(x.HSCode, r.HSCode):
			why = "hs_mismatch"
		case !x.DecisionUseAllowed:
			why = "license_blocked"
		case x.Confidence < r.MinConfidence:
			why = "low_confidence"
		case x.MaxAge <= 0 || r.At.Sub(x.ObservedAt) > x.MaxAge || x.ObservedAt.After(r.At.Add(time.Minute)):
			why = "stale_or_future"
		case !x.EffectiveFrom.IsZero() && r.At.Before(x.EffectiveFrom):
			why = "not_yet_effective"
		case !x.EffectiveTo.IsZero() && r.At.After(x.EffectiveTo):
			why = "expired"
		}
		if why != "" {
			d.Rejected[key] = why
			continue
		}
		applicable = append(applicable, x)
	}
	if len(applicable) == 0 {
		return d, ErrNoApplicableRule
	}
	sort.Slice(applicable, func(i, j int) bool { return applicable[i].ID < applicable[j].ID })
	for _, x := range applicable {
		if x.Prohibited {
			d.RuleIDs = append(d.RuleIDs, x.ID)
			return d, ErrTradeBlocked
		}
		d.RuleIDs = append(d.RuleIDs, x.ID)
		d.Duty += r.CustomsValue*x.DutyRate + r.Quantity*x.FixedDutyPerUnit
		if d.Currency == "" {
			d.Currency = x.Currency
		} else if x.Currency != "" && d.Currency != x.Currency {
			return d, ErrInvalidRequest
		}
		for _, doc := range x.RequiredDocuments {
			doc = strings.TrimSpace(doc)
			if doc != "" {
				docs[doc] = true
			}
		}
	}
	for doc := range docs {
		d.RequiredDocuments = append(d.RequiredDocuments, doc)
	}
	sort.Strings(d.RequiredDocuments)
	return d, nil
}
