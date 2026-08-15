package productionreadiness
import("errors";"sort")
var ErrNotReady=errors.New("production readiness gate failed")
type Check struct{Name string;Pass bool;Evidence string;Critical bool}
type Report struct{Ready bool;Failures []string}
func Evaluate(cs []Check)Report{
 r:=Report{Ready:true}
 required:=map[string]bool{"live_integrations":false,"regulatory_data":false,"observability":false,"load_chaos":false,"security_assessment":false,"disaster_recovery":false,"red_team":false,"end_to_end_trade":false}
 for _,c:=range cs{
  if _,ok:=required[c.Name];ok{required[c.Name]=c.Pass&&c.Evidence!=""}
  if c.Critical&&(!c.Pass||c.Evidence==""){r.Ready=false;r.Failures=append(r.Failures,c.Name)}
 }
 for n,ok:=range required{if !ok{r.Ready=false;r.Failures=append(r.Failures,n)}}
 sort.Strings(r.Failures);r.Failures=dedupe(r.Failures);return r
}
func dedupe(x []string)[]string{if len(x)==0{return x};o:=[]string{x[0]};for _,v:=range x[1:]{if v!=o[len(o)-1]{o=append(o,v)}};return o}
