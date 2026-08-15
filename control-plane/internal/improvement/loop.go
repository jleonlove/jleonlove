package improvement
import("errors";"sort")
var(ErrEvidence=errors.New("release evidence required");ErrRegression=errors.New("regression detected");ErrApproval=errors.New("release approval required");ErrUnsafe=errors.New("unsafe candidate"))
type Eval struct{Name string;Baseline,Candidate float64;HigherIsBetter bool;Critical bool}
type Candidate struct{ID string;Evals []Eval;SafetyPassed,RedTeamPassed,EvidenceSigned,Approved bool}
type Decision struct{Promote bool;Reasons []string}
func Assess(c Candidate)Decision{
 d:=Decision{}
 if !c.SafetyPassed||!c.RedTeamPassed{d.Reasons=append(d.Reasons,ErrUnsafe.Error())}
 if !c.EvidenceSigned{d.Reasons=append(d.Reasons,ErrEvidence.Error())}
 if !c.Approved{d.Reasons=append(d.Reasons,ErrApproval.Error())}
 for _,e:=range c.Evals{
  reg:=(e.HigherIsBetter&&e.Candidate<e.Baseline)||(!e.HigherIsBetter&&e.Candidate>e.Baseline)
  if reg&&e.Critical{d.Reasons=append(d.Reasons,ErrRegression.Error()+": "+e.Name)}
 }
 sort.Strings(d.Reasons);d.Promote=len(d.Reasons)==0;return d
}
