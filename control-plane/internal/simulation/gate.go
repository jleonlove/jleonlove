package simulation
import("errors";"math";"sort")
var(ErrModel=errors.New("simulation model invalid");ErrRisk=errors.New("simulated risk exceeds limit");ErrEvidence=errors.New("simulation evidence required");ErrApproval=errors.New("approval required"))
type Scenario struct{Name string;Probability,Loss float64}
type Request struct{Scenarios []Scenario;MaxExpectedLoss float64;Evidence bool;Consequential bool;Approved bool}
type Result struct{ExpectedLoss,WorstLoss float64;Ordered []string}
func Evaluate(r Request)(Result,error){
 if !r.Evidence{return Result{},ErrEvidence};if r.Consequential&&!r.Approved{return Result{},ErrApproval}
 var out Result;var ps float64
 for _,s:=range r.Scenarios{if math.IsNaN(s.Probability)||math.IsNaN(s.Loss)||s.Probability<0||s.Probability>1||s.Loss<0{return Result{},ErrModel};ps+=s.Probability;out.ExpectedLoss+=s.Probability*s.Loss;if s.Loss>out.WorstLoss{out.WorstLoss=s.Loss};out.Ordered=append(out.Ordered,s.Name)}
 if math.Abs(ps-1)>1e-9{return Result{},ErrModel};sort.Strings(out.Ordered);if out.ExpectedLoss>r.MaxExpectedLoss{return out,ErrRisk};return out,nil
}
