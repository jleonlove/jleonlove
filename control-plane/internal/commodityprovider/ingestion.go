package commodityprovider

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"atlas/internal/commoditydata"
)

type QuoteRequest struct {
	Commodity, Metric string
	At                time.Time
}
type Provider interface {
	Name() string
	Fetch(context.Context, QuoteRequest) (commoditydata.Observation, error)
}
type Config struct {
	Priority    int
	Timeout     time.Duration
	MaxFailures int
	Cooldown    time.Duration
	Enabled     bool
}
type state struct {
	failures  int
	openUntil time.Time
}
type Result struct {
	Attempted, Succeeded []string
	Failed               map[string]string
}
type Ingestor struct {
	mu        sync.Mutex
	fabric    *commoditydata.Fabric
	providers map[string]Provider
	cfg       map[string]Config
	state     map[string]state
	now       func() time.Time
}

func New(f *commoditydata.Fabric) *Ingestor {
	return &Ingestor{fabric: f, providers: map[string]Provider{}, cfg: map[string]Config{}, state: map[string]state{}, now: func() time.Time { return time.Now().UTC() }}
}
func (i *Ingestor) Register(p Provider, c Config) error {
	if p == nil || strings.TrimSpace(p.Name()) == "" {
		return errors.New("provider and name required")
	}
	if c.Timeout <= 0 {
		c.Timeout = 5 * time.Second
	}
	if c.MaxFailures <= 0 {
		c.MaxFailures = 3
	}
	if c.Cooldown <= 0 {
		c.Cooldown = time.Minute
	}
	c.Enabled = true
	i.mu.Lock()
	defer i.mu.Unlock()
	i.providers[p.Name()] = p
	i.cfg[p.Name()] = c
	return nil
}
func (i *Ingestor) Fetch(ctx context.Context, q QuoteRequest) (Result, error) {
	if i.fabric == nil {
		return Result{}, errors.New("commodity data fabric required")
	}
	if strings.TrimSpace(q.Commodity) == "" || strings.TrimSpace(q.Metric) == "" {
		return Result{}, errors.New("commodity and metric required")
	}
	i.mu.Lock()
	names := make([]string, 0, len(i.providers))
	for n := range i.providers {
		names = append(names, n)
	}
	sort.Slice(names, func(a, b int) bool {
		ca, cb := i.cfg[names[a]], i.cfg[names[b]]
		if ca.Priority != cb.Priority {
			return ca.Priority < cb.Priority
		}
		return names[a] < names[b]
	})
	i.mu.Unlock()
	r := Result{Failed: map[string]string{}}
	now := i.now()
	for _, n := range names {
		i.mu.Lock()
		p, c, s := i.providers[n], i.cfg[n], i.state[n]
		i.mu.Unlock()
		if !c.Enabled || now.Before(s.openUntil) {
			continue
		}
		r.Attempted = append(r.Attempted, n)
		cctx, cancel := context.WithTimeout(ctx, c.Timeout)
		o, err := p.Fetch(cctx, q)
		cancel()
		if err == nil {
			if o.Source == "" {
				o.Source = n
			}
			if !strings.EqualFold(o.Source, n) {
				err = fmt.Errorf("provider source mismatch: %s", o.Source)
			}
			if o.Commodity == "" {
				o.Commodity = q.Commodity
			}
			if o.Metric == "" {
				o.Metric = q.Metric
			}
			if err == nil {
				err = i.fabric.Ingest(o)
			}
		}
		i.mu.Lock()
		s = i.state[n]
		if err != nil {
			s.failures++
			r.Failed[n] = err.Error()
			if s.failures >= c.MaxFailures {
				s.openUntil = now.Add(c.Cooldown)
				s.failures = 0
			}
		} else {
			s = state{}
			r.Succeeded = append(r.Succeeded, n)
		}
		i.state[n] = s
		i.mu.Unlock()
	}
	if len(r.Succeeded) == 0 {
		return r, errors.New("all commodity providers failed or unavailable")
	}
	return r, nil
}
