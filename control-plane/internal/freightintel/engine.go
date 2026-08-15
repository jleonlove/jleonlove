package freightintel

import (
	"errors"
	"sort"
	"strings"
	"time"
)

var ErrInvalidQuote = errors.New("invalid freight quote")
var ErrNoEligibleQuote = errors.New("no eligible freight quote")

type Quote struct {
	Provider, Origin, Destination, Mode, Currency string
	CostPerUnit, Confidence, Capacity             float64
	ObservedAt                                    time.Time
	MaxAge                                        time.Duration
	DecisionUseAllowed                            bool
}
type Request struct {
	Origin, Destination, Mode string
	Quantity                  float64
	At                        time.Time
	MinConfidence             float64
}
type Decision struct {
	Quote        Quote
	TotalFreight float64
	Rejected     map[string]string
}

func Select(r Request, qs []Quote) (Decision, error) {
	if strings.TrimSpace(r.Origin) == "" || strings.TrimSpace(r.Destination) == "" || r.Quantity <= 0 {
		return Decision{}, ErrInvalidQuote
	}
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	if r.MinConfidence <= 0 {
		r.MinConfidence = .7
	}
	ok := []Quote{}
	bad := map[string]string{}
	for _, q := range qs {
		k := q.Provider
		if k == "" {
			k = "unknown"
		}
		why := ""
		switch {
		case q.Provider == "" || q.Origin == "" || q.Destination == "" || q.Mode == "" || q.Currency == "" || q.CostPerUnit < 0 || q.Confidence < 0 || q.Confidence > 1 || q.ObservedAt.IsZero():
			why = "invalid"
		case !strings.EqualFold(q.Origin, r.Origin) || !strings.EqualFold(q.Destination, r.Destination):
			why = "route_mismatch"
		case r.Mode != "" && !strings.EqualFold(q.Mode, r.Mode):
			why = "mode_mismatch"
		case !q.DecisionUseAllowed:
			why = "license_blocked"
		case q.Confidence < r.MinConfidence:
			why = "low_confidence"
		case q.MaxAge <= 0 || r.At.Sub(q.ObservedAt) > q.MaxAge || q.ObservedAt.After(r.At.Add(time.Minute)):
			why = "stale_or_future"
		case q.Capacity > 0 && q.Capacity < r.Quantity:
			why = "insufficient_capacity"
		}
		if why != "" {
			bad[k] = why
			continue
		}
		ok = append(ok, q)
	}
	if len(ok) == 0 {
		return Decision{Rejected: bad}, ErrNoEligibleQuote
	}
	sort.SliceStable(ok, func(i, j int) bool {
		if ok[i].CostPerUnit != ok[j].CostPerUnit {
			return ok[i].CostPerUnit < ok[j].CostPerUnit
		}
		if ok[i].Confidence != ok[j].Confidence {
			return ok[i].Confidence > ok[j].Confidence
		}
		return ok[i].Provider < ok[j].Provider
	})
	q := ok[0]
	return Decision{Quote: q, TotalFreight: q.CostPerUnit * r.Quantity, Rejected: bad}, nil
}
