package dataflowvm
import("errors";"sort")
var(ErrCycle=errors.New("dataflow cycle");ErrCapability=errors.New("tool capability denied");ErrDependency=errors.New("dependency unavailable");ErrBudget=errors.New("execution budget exceeded"))
type Node struct{ID,Capability string;DependsOn []string;Cost int}
type Program struct{Nodes map[string]Node;Budget int}
type Plan struct{Order []string;Cost int}
func Compile(p Program,allowed map[string]bool)(Plan,error){
 indeg:=map[string]int{};next:=map[string][]string{}
 for id,n:=range p.Nodes{if !allowed[n.Capability]{return Plan{},ErrCapability};indeg[id]=len(n.DependsOn);for _,d:=range n.DependsOn{if _,ok:=p.Nodes[d];!ok{return Plan{},ErrDependency};next[d]=append(next[d],id)}}
 var ready []string;for id,d:=range indeg{if d==0{ready=append(ready,id)}};sort.Strings(ready)
 out:=Plan{}
 for len(ready)>0{id:=ready[0];ready=ready[1:];n:=p.Nodes[id];out.Cost+=n.Cost;if out.Cost>p.Budget{return Plan{},ErrBudget};out.Order=append(out.Order,id);for _,x:=range next[id]{indeg[x]--;if indeg[x]==0{ready=append(ready,x);sort.Strings(ready)}}}
 if len(out.Order)!=len(p.Nodes){return Plan{},ErrCycle};return out,nil
}
