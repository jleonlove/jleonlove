package instructions
import("errors";"sort";"strings")
var(ErrPolicyOverride=errors.New("lower-level instruction attempts to override higher policy");ErrAuthorityExpansion=errors.New("instruction attempts authority expansion"))
type Level int
const(Organization Level=iota;Workspace;Repository;Directory;Task;Agent)
type Rule struct{Level Level;Key,Value string;Deny bool}
type Effective struct{Rules map[string]Rule}
func Compile(rules []Rule)(Effective,error){
 sort.SliceStable(rules,func(i,j int)bool{return rules[i].Level<rules[j].Level})
 out:=Effective{Rules:map[string]Rule{}}
 for _,r:=range rules{
  if old,ok:=out.Rules[r.Key];ok{
   if old.Deny && !r.Deny{return Effective{},ErrPolicyOverride}
   if strings.HasPrefix(r.Key,"capability.") && old.Value=="deny" && r.Value=="allow"{return Effective{},ErrAuthorityExpansion}
  }
  out.Rules[r.Key]=r
 }
 return out,nil
}
func Allowed(e Effective,key string)bool{r,ok:=e.Rules[key];return ok&&!r.Deny&&r.Value!="deny"}
