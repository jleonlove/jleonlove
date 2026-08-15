package tradecompliance
import("errors";"sort";"time")
var(ErrRule=errors.New("invalid compliance rule");ErrBlocked=errors.New("trade blocked by compliance rule");ErrEvidence=errors.New("compliance evidence missing"))
type Action string
const(Allow Action="ALLOW";Review Action="REVIEW";Block Action="BLOCK")
type Rule struct{ID,Jurisdiction,Commodity,Origin,Destination string;Effective,Expires time.Time;Action Action;EvidenceRequired bool}
type Trade struct{Commodity,Origin,Destination string;At time.Time;Evidence map[string]bool}
type Finding struct{RuleID string;Action Action;Reason string}
func active(r Rule,t time.Time)bool{return !t.Before(r.Effective)&&(r.Expires.IsZero()||t.Before(r.Expires))}
func Compile(rs []Rule)([]Rule,error){for _,r:=range rs{if r.ID==""||r.Jurisdiction==""||r.Effective.IsZero(){return nil,ErrRule};if r.Action!=Allow&&r.Action!=Review&&r.Action!=Block{return nil,ErrRule}};out:=append([]Rule(nil),rs...);sort.Slice(out,func(i,j int)bool{return out[i].ID<out[j].ID});return out,nil}
func Evaluate(rs []Rule,t Trade)([]Finding,error){
 var out []Finding
 for _,r:=range rs{if !active(r,t.At){continue};if r.Commodity!=""&&r.Commodity!=t.Commodity{continue};if r.Origin!=""&&r.Origin!=t.Origin{continue};if r.Destination!=""&&r.Destination!=t.Destination{continue}
  if r.EvidenceRequired&&!t.Evidence[r.ID]{out=append(out,Finding{r.ID,Review,"required evidence missing"});continue}
  out=append(out,Finding{r.ID,r.Action,"rule matched"})}
 sort.Slice(out,func(i,j int)bool{return out[i].RuleID<out[j].RuleID});for _,f:=range out{if f.Action==Block{return out,ErrBlocked}};return out,nil
}
