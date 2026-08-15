package tradeintegrityengine
type Signals struct{Docs,Counterparty,Compliance,Provenance,Economics bool}
type Result struct{Score int;Pass bool;Failures []string}
func Evaluate(s Signals)Result{r:=Result{};xs:=[]struct{n string;v bool}{{"docs",s.Docs},{"counterparty",s.Counterparty},{"compliance",s.Compliance},{"provenance",s.Provenance},{"economics",s.Economics}};for _,x:=range xs{if x.v{r.Score+=20}else{r.Failures=append(r.Failures,x.n)}};r.Pass=r.Score==100;return r}