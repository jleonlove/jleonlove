package commoditydata

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type LicensePolicy struct {
	AllowDecisionUse  bool
	AllowRedistribute bool
}

type Observation struct {
	Commodity  string
	Metric     string
	Value      float64
	Unit       string
	Currency   string
	Source     string
	SourceID   string
	ObservedAt time.Time
	IngestedAt time.Time
	MaxAge     time.Duration
	Confidence float64
	License    LicensePolicy
	Evidence   []string
}

type Request struct {
	Commodity           string
	Metric              string
	At                  time.Time
	RequireDecisionUse  bool
	MinConfidence       float64
	MaxAge              time.Duration
	MaxRelativeConflict float64
}

type Resolution struct {
	Observation     Observation
	Candidates      int
	FreshCandidates int
	Conflicting     bool
	Provenance      []string
}

type Fabric struct {
	mu           sync.RWMutex
	observations []Observation
}

func New() *Fabric { return &Fabric{} }

func (f *Fabric) Ingest(o Observation) error {
	o.Commodity = strings.TrimSpace(o.Commodity)
	o.Metric = strings.TrimSpace(o.Metric)
	o.Source = strings.TrimSpace(o.Source)
	if o.Commodity == "" || o.Metric == "" || o.Source == "" {
		return errors.New("commodity, metric, and source are required")
	}
	if o.ObservedAt.IsZero() {
		return errors.New("observed_at is required")
	}
	if o.Confidence < 0 || o.Confidence > 1 {
		return errors.New("confidence must be between 0 and 1")
	}
	if o.MaxAge <= 0 {
		return errors.New("max_age must be positive")
	}
	if len(o.Evidence) == 0 {
		return errors.New("provenance evidence is required")
	}
	if o.IngestedAt.IsZero() {
		o.IngestedAt = time.Now().UTC()
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.observations = append(f.observations, o)
	return nil
}

func (f *Fabric) Resolve(r Request) (Resolution, error) {
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	if r.MinConfidence < 0 || r.MinConfidence > 1 {
		return Resolution{}, errors.New("min confidence must be between 0 and 1")
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	var all, eligible []Observation
	for _, o := range f.observations {
		if !strings.EqualFold(o.Commodity, strings.TrimSpace(r.Commodity)) || !strings.EqualFold(o.Metric, strings.TrimSpace(r.Metric)) {
			continue
		}
		all = append(all, o)
		ageLimit := o.MaxAge
		if r.MaxAge > 0 && r.MaxAge < ageLimit {
			ageLimit = r.MaxAge
		}
		age := r.At.Sub(o.ObservedAt)
		if age < 0 {
			continue
		}
		if age > ageLimit || o.Confidence < r.MinConfidence {
			continue
		}
		if r.RequireDecisionUse && !o.License.AllowDecisionUse {
			continue
		}
		eligible = append(eligible, o)
	}
	if len(all) == 0 {
		return Resolution{}, errors.New("no observations found")
	}
	if len(eligible) == 0 {
		return Resolution{Candidates: len(all)}, errors.New("no fresh, licensed, sufficiently confident observation")
	}
	sort.SliceStable(eligible, func(i, j int) bool {
		if eligible[i].Confidence != eligible[j].Confidence {
			return eligible[i].Confidence > eligible[j].Confidence
		}
		if !eligible[i].ObservedAt.Equal(eligible[j].ObservedAt) {
			return eligible[i].ObservedAt.After(eligible[j].ObservedAt)
		}
		return eligible[i].Source < eligible[j].Source
	})
	threshold := r.MaxRelativeConflict
	if threshold <= 0 {
		threshold = 0.05
	}
	best := eligible[0]
	conflicting := false
	for _, o := range eligible[1:] {
		if best.Value == 0 {
			if o.Value != 0 {
				conflicting = true
			}
			continue
		}
		d := o.Value - best.Value
		if d < 0 {
			d = -d
		}
		if d/abs(best.Value) > threshold {
			conflicting = true
			break
		}
	}
	if conflicting {
		return Resolution{Candidates: len(all), FreshCandidates: len(eligible), Conflicting: true}, fmt.Errorf("conflicting fresh observations exceed %.2f%% threshold", threshold*100)
	}
	prov := append([]string(nil), best.Evidence...)
	sort.Strings(prov)
	return Resolution{Observation: best, Candidates: len(all), FreshCandidates: len(eligible), Provenance: prov}, nil
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
