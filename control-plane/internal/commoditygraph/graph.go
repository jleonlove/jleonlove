package commoditygraph

import (
	"errors"
	"sort"
	"strings"
)

var ErrInvalidGraph = errors.New("invalid commodity graph")

type Node struct {
	ID, Kind, Name string
	Attributes     map[string]string
}
type Edge struct {
	From, To, Relation string
	Weight             float64
	Evidence           []string
}
type Graph struct {
	Nodes map[string]Node
	Edges []Edge
}
type Path struct {
	Nodes     []string
	Relations []string
	Score     float64
}
type Shock struct {
	NodeID    string
	ChangePct float64
}
type Impact struct {
	NodeID    string
	ChangePct float64
	Depth     int
	Path      []string
}

func New(nodes []Node, edges []Edge) (*Graph, error) {
	g := &Graph{Nodes: map[string]Node{}, Edges: append([]Edge(nil), edges...)}
	for _, n := range nodes {
		if strings.TrimSpace(n.ID) == "" || strings.TrimSpace(n.Kind) == "" || strings.TrimSpace(n.Name) == "" {
			return nil, ErrInvalidGraph
		}
		if _, ok := g.Nodes[n.ID]; ok {
			return nil, ErrInvalidGraph
		}
		g.Nodes[n.ID] = n
	}
	for _, e := range edges {
		if _, ok := g.Nodes[e.From]; !ok {
			return nil, ErrInvalidGraph
		}
		if _, ok := g.Nodes[e.To]; !ok {
			return nil, ErrInvalidGraph
		}
		if strings.TrimSpace(e.Relation) == "" || e.Weight < -1 || e.Weight > 1 {
			return nil, ErrInvalidGraph
		}
	}
	return g, nil
}
func (g *Graph) Neighbors(id string) []Edge {
	var out []Edge
	for _, e := range g.Edges {
		if e.From == id {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].To == out[j].To {
			return out[i].Relation < out[j].Relation
		}
		return out[i].To < out[j].To
	})
	return out
}
func (g *Graph) FindPaths(from, to string, maxDepth int) []Path {
	if maxDepth < 1 {
		return nil
	}
	var out []Path
	var walk func(string, []string, []string, float64, map[string]bool, int)
	walk = func(cur string, ns, rs []string, score float64, seen map[string]bool, d int) {
		if d > maxDepth {
			return
		}
		if cur == to && len(rs) > 0 {
			out = append(out, Path{append([]string(nil), ns...), append([]string(nil), rs...), score})
			return
		}
		for _, e := range g.Neighbors(cur) {
			if seen[e.To] {
				continue
			}
			cp := map[string]bool{}
			for k, v := range seen {
				cp[k] = v
			}
			cp[e.To] = true
			walk(e.To, append(ns, e.To), append(rs, e.Relation), score*e.Weight, cp, d+1)
		}
	}
	if _, ok := g.Nodes[from]; !ok {
		return nil
	}
	if _, ok := g.Nodes[to]; !ok {
		return nil
	}
	walk(from, []string{from}, nil, 1, map[string]bool{from: true}, 0)
	sort.Slice(out, func(i, j int) bool { return abs(out[i].Score) > abs(out[j].Score) })
	return out
}
func (g *Graph) Propagate(s Shock, maxDepth int, minAbs float64) ([]Impact, error) {
	if _, ok := g.Nodes[s.NodeID]; !ok || maxDepth < 1 || minAbs < 0 {
		return nil, ErrInvalidGraph
	}
	impacts := map[string]Impact{}
	type state struct {
		id    string
		chg   float64
		depth int
		path  []string
	}
	q := []state{{s.NodeID, s.ChangePct, 0, []string{s.NodeID}}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.depth >= maxDepth {
			continue
		}
		for _, e := range g.Neighbors(cur.id) {
			chg := cur.chg * e.Weight
			if abs(chg) < minAbs {
				continue
			}
			p := append(append([]string(nil), cur.path...), e.To)
			old, ok := impacts[e.To]
			if !ok || abs(chg) > abs(old.ChangePct) {
				impacts[e.To] = Impact{e.To, chg, cur.depth + 1, p}
				q = append(q, state{e.To, chg, cur.depth + 1, p})
			}
		}
	}
	out := make([]Impact, 0, len(impacts))
	for _, v := range impacts {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Depth == out[j].Depth {
			return out[i].NodeID < out[j].NodeID
		}
		return out[i].Depth < out[j].Depth
	})
	return out, nil
}
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
