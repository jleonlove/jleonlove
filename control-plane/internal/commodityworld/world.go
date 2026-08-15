package commodityworld
import("errors";"math";"sort")
var(ErrNode=errors.New("world-model node missing");ErrCycle=errors.New("causal cycle");ErrProbability=errors.New("invalid scenario probability");ErrEvidence=errors.New("observation evidence required"))
type Node struct{ID,Kind string}
type Edge struct{From,To string;Weight float64}
type Observation struct{NodeID,Evidence string;Delta float64}
type Scenario struct{Name string;Probability float64;Shocks map[string]float64}
type Model struct{Nodes map[string]Node;Edges []Edge}
type Forecast struct{Expected map[string]float64;Scenarios []string}
func(m Model)Validate()error{
 indeg:=map[string]int{};next:=map[string][]string{}
 for id:=range m.Nodes{indeg[id]=0}
 for _,e:=range m.Edges{if _,ok:=m.Nodes[e.From];!ok{return ErrNode};if _,ok:=m.Nodes[e.To];!ok{return ErrNode};indeg[e.To]++;next[e.From]=append(next[e.From],e.To)}
 q:=[]string{};for id,d:=range indeg{if d==0{q=append(q,id)}};seen:=0
 for len(q)>0{x:=q[0];q=q[1:];seen++;for _,n:=range next[x]{indeg[n]--;if indeg[n]==0{q=append(q,n)}}}
 if seen!=len(m.Nodes){return ErrCycle};return nil
}
func(m Model)Propagate(base map[string]float64,obs []Observation) (map[string]float64,error){
 if e:=m.Validate();e!=nil{return nil,e};v:=map[string]float64{};for k,x:=range base{v[k]=x}
 for _,o:=range obs{if _,ok:=m.Nodes[o.NodeID];!ok{return nil,ErrNode};if o.Evidence==""{return nil,ErrEvidence};v[o.NodeID]+=o.Delta}
 indeg:=map[string]int{};next:=map[string][]Edge{};for id:=range m.Nodes{indeg[id]=0};for _,e:=range m.Edges{indeg[e.To]++;next[e.From]=append(next[e.From],e)}
 q:=[]string{};for id,d:=range indeg{if d==0{q=append(q,id)}};sort.Strings(q)
 for len(q)>0{x:=q[0];q=q[1:];for _,e:=range next[x]{v[e.To]+=v[x]*e.Weight;indeg[e.To]--;if indeg[e.To]==0{q=append(q,e.To);sort.Strings(q)}}}
 return v,nil
}
func(m Model)Expected(base map[string]float64,ss []Scenario)(Forecast,error){
 var ps float64;out:=Forecast{Expected:map[string]float64{}}
 for _,s:=range ss{if math.IsNaN(s.Probability)||s.Probability<0||s.Probability>1{return out,ErrProbability};ps+=s.Probability;obs:=[]Observation{};for id,d:=range s.Shocks{obs=append(obs,Observation{id,s.Name,d})};v,e:=m.Propagate(base,obs);if e!=nil{return out,e};for id,x:=range v{out.Expected[id]+=s.Probability*x};out.Scenarios=append(out.Scenarios,s.Name)}
 if math.Abs(ps-1)>1e-9{return Forecast{},ErrProbability};sort.Strings(out.Scenarios);return out,nil
}
