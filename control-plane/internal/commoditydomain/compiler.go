package commoditydomain

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

var (
	ErrInvalidPack      = errors.New("invalid commodity domain pack")
	ErrDuplicatePack    = errors.New("duplicate commodity domain pack")
	ErrUnknownCommodity = errors.New("unknown commodity")
)

type Specification struct {
	Name     string
	Unit     string
	Required bool
}
type DocumentRule struct {
	Name     string
	Stage    string
	Required bool
}
type RiskRule struct {
	Code        string
	Description string
	Severity    string
}
type Pack struct {
	ID             string
	Family         string
	Commodity      string
	Aliases        []string
	Grades         []string
	Specifications []Specification
	Benchmarks     []string
	Documents      []DocumentRule
	Risks          []RiskRule
	Version        string
}
type Compiled struct {
	Pack        Pack
	SearchTerms []string
}
type Compiler struct {
	mu      sync.RWMutex
	packs   map[string]Compiled
	aliases map[string]string
}

func New() *Compiler       { return &Compiler{packs: map[string]Compiled{}, aliases: map[string]string{}} }
func norm(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
func validate(p Pack) error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.Family) == "" || strings.TrimSpace(p.Commodity) == "" || strings.TrimSpace(p.Version) == "" {
		return ErrInvalidPack
	}
	for _, s := range p.Specifications {
		if strings.TrimSpace(s.Name) == "" {
			return ErrInvalidPack
		}
	}
	return nil
}
func (c *Compiler) Compile(p Pack) (Compiled, error) {
	if err := validate(p); err != nil {
		return Compiled{}, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.packs[p.ID]; ok {
		return Compiled{}, ErrDuplicatePack
	}
	terms := []string{p.ID, p.Family, p.Commodity}
	terms = append(terms, p.Aliases...)
	terms = append(terms, p.Grades...)
	terms = append(terms, p.Benchmarks...)
	seen := map[string]bool{}
	clean := make([]string, 0, len(terms))
	for _, t := range terms {
		n := norm(t)
		if n != "" && !seen[n] {
			seen[n] = true
			clean = append(clean, n)
		}
	}
	sort.Strings(clean)
	out := Compiled{Pack: p, SearchTerms: clean}
	c.packs[p.ID] = out
	for _, a := range append([]string{p.ID, p.Commodity}, p.Aliases...) {
		if n := norm(a); n != "" {
			c.aliases[n] = p.ID
		}
	}
	return out, nil
}
func (c *Compiler) Resolve(q string) (Compiled, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.aliases[norm(q)]
	if !ok {
		id = q
	}
	p, ok := c.packs[id]
	if !ok {
		return Compiled{}, ErrUnknownCommodity
	}
	return p, nil
}
func (c *Compiler) IDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ids := make([]string, 0, len(c.packs))
	for id := range c.packs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func SeedCore(c *Compiler) error {
	packs := []Pack{
		{ID: "commodity.gold", Family: "Metals/Precious Metals", Commodity: "Gold", Aliases: []string{"au", "bullion", "gold dore"}, Grades: []string{"doré", "999.9 bullion"}, Specifications: []Specification{{"fineness", "ppt", true}, {"weight", "kg", true}}, Benchmarks: []string{"LBMA Gold Price"}, Documents: []DocumentRule{{"assay certificate", "verification", true}, {"certificate of origin", "compliance", true}}, Risks: []RiskRule{{"PROVENANCE", "Unverified origin or chain of custody", "HIGH"}}, Version: "1.0.0"},
		{ID: "commodity.wheat", Family: "Agriculture/Grains", Commodity: "Wheat", Aliases: []string{"milling wheat", "feed wheat"}, Grades: []string{"milling", "feed"}, Specifications: []Specification{{"protein", "%", true}, {"moisture", "%", true}, {"test weight", "kg/hl", false}}, Documents: []DocumentRule{{"phytosanitary certificate", "export", true}, {"certificate of origin", "export", true}}, Risks: []RiskRule{{"QUALITY", "Grade/specification mismatch", "MEDIUM"}}, Version: "1.0.0"},
		{ID: "commodity.sulfur", Family: "Chemicals/Fertilizer Feedstocks", Commodity: "Sulfur", Aliases: []string{"sulphur", "granular sulfur", "prilled sulfur"}, Grades: []string{"granular", "prilled", "lump", "molten"}, Specifications: []Specification{{"purity", "%", true}, {"moisture", "%", false}, {"ash", "%", false}}, Documents: []DocumentRule{{"certificate of analysis", "verification", true}, {"certificate of origin", "export", true}}, Risks: []RiskRule{{"HANDLING", "Storage/transport handling requirements", "MEDIUM"}}, Version: "1.0.0"},
		{ID: "commodity.beef", Family: "Agriculture/Animal Protein", Commodity: "Beef", Aliases: []string{"frozen beef", "chilled beef"}, Grades: []string{"frozen", "chilled"}, Specifications: []Specification{{"cut", "text", true}, {"storage temperature", "C", true}}, Documents: []DocumentRule{{"veterinary health certificate", "export", true}, {"establishment eligibility", "compliance", true}}, Risks: []RiskRule{{"COLD_CHAIN", "Cold-chain integrity failure", "HIGH"}}, Version: "1.0.0"},
		{ID: "commodity.frac_sand", Family: "Industrial Minerals/Oilfield Services", Commodity: "Frac Sand", Aliases: []string{"proppant sand", "hydraulic fracturing sand"}, Grades: []string{"40/70", "100 mesh"}, Specifications: []Specification{{"mesh", "mesh", true}, {"turbidity", "NTU", false}, {"crush resistance", "psi", false}}, Documents: []DocumentRule{{"quality test report", "verification", true}}, Risks: []RiskRule{{"LOGISTICS", "Delivered-cost radius can destroy economics", "HIGH"}}, Version: "1.0.0"},
	}
	for _, p := range packs {
		if _, err := c.Compile(p); err != nil {
			return err
		}
	}
	return nil
}
