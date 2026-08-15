package knowledge
import("errors";"sort";"strings";"time")
var(ErrPermission=errors.New("knowledge permission denied");ErrExpired=errors.New("knowledge expired"))
type Fact struct{ID,Subject,Predicate,Object,Source,Tenant,Classification string;ObservedAt,ValidUntil time.Time;Reliability float64;Readers map[string]bool}
type Query struct{Tenant,Reader,Subject,Predicate string;Now time.Time}
type Result struct{Facts []Fact;Conflicted bool}
func Retrieve(all []Fact,q Query)(Result,error){
 var out []Fact
 for _,f:=range all{
  if f.Tenant!=q.Tenant||f.Subject!=q.Subject||f.Predicate!=q.Predicate{continue}
  if !f.Readers[q.Reader]{continue}
  if !f.ValidUntil.IsZero()&&!q.Now.Before(f.ValidUntil){continue}
  out=append(out,f)
 }
 if len(out)==0{return Result{},ErrPermission}
 sort.SliceStable(out,func(i,j int)bool{if out[i].Reliability==out[j].Reliability{return out[i].ObservedAt.After(out[j].ObservedAt)};return out[i].Reliability>out[j].Reliability})
 vals:=map[string]bool{};for _,f:=range out{vals[strings.TrimSpace(strings.ToLower(f.Object))]=true}
 return Result{Facts:out,Conflicted:len(vals)>1},nil
}
func Ingest(f Fact)error{if f.Source==""||f.Tenant==""||f.Classification==""{return errors.New("provenance required")};if f.Reliability<0||f.Reliability>1{return errors.New("invalid reliability")};return nil}
