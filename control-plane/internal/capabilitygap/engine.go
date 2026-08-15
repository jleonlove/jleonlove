package capabilitygap

import (
	"sort"
	"strings"
	"time"
)

type Signal struct {
	Task, Domain, Failure, MissingCapability string
	Severity, Frequency                      int
	Evidence                                 []string
	At                                       time.Time
}
type Proposal struct {
	Capability, Domain    string
	Priority, Score       int
	Evidence              []string
	RequiresHumanApproval bool
	Status                string
}
type Engine struct{ signals []Signal }

func New() *Engine { return &Engine{} }
func (e *Engine) Observe(s Signal) {
	if s.Severity < 0 {
		s.Severity = 0
	}
	if s.Frequency < 1 {
		s.Frequency = 1
	}
	e.signals = append(e.signals, s)
}
func (e *Engine) Proposals() []Proposal {
	type agg struct {
		p  Proposal
		ev map[string]bool
	}
	m := map[string]*agg{}
	for _, s := range e.signals {
		c := strings.TrimSpace(s.MissingCapability)
		if c == "" {
			continue
		}
		k := strings.ToLower(s.Domain + "|" + c)
		a := m[k]
		if a == nil {
			a = &agg{p: Proposal{Capability: c, Domain: s.Domain, RequiresHumanApproval: true, Status: "proposed"}, ev: map[string]bool{}}
			m[k] = a
		}
		a.p.Score += s.Severity * s.Frequency
		for _, v := range s.Evidence {
			if strings.TrimSpace(v) != "" {
				a.ev[v] = true
			}
		}
	}
	out := make([]Proposal, 0, len(m))
	for _, a := range m {
		if a.p.Score >= 50 {
			a.p.Priority = 0
		} else if a.p.Score >= 20 {
			a.p.Priority = 1
		} else {
			a.p.Priority = 2
		}
		for v := range a.ev {
			a.p.Evidence = append(a.p.Evidence, v)
		}
		sort.Strings(a.p.Evidence)
		out = append(out, a.p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Capability < out[j].Capability
		}
		return out[i].Score > out[j].Score
	})
	return out
}
func (e *Engine) CanSelfDeploy() bool { return false }
