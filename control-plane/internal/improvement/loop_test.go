package improvement
import"testing"
func good()Candidate{return Candidate{ID:"c1",SafetyPassed:true,RedTeamPassed:true,EvidenceSigned:true,Approved:true,Evals:[]Eval{{Name:"accuracy",Baseline:.8,Candidate:.85,HigherIsBetter:true,Critical:true},{Name:"hallucination",Baseline:.1,Candidate:.08,HigherIsBetter:false,Critical:true}}}}
func TestGoodCandidatePromotes(t *testing.T){if d:=Assess(good());!d.Promote{t.Fatalf("%v",d.Reasons)}}
func TestUnsignedCannotPromote(t *testing.T){c:=good();c.EvidenceSigned=false;if Assess(c).Promote{t.Fatal("unsigned promoted")}}
func TestUnapprovedCannotPromote(t *testing.T){c:=good();c.Approved=false;if Assess(c).Promote{t.Fatal("unapproved promoted")}}
func TestSafetyFailureCannotPromote(t *testing.T){c:=good();c.SafetyPassed=false;if Assess(c).Promote{t.Fatal("unsafe promoted")}}
func TestRedTeamFailureCannotPromote(t *testing.T){c:=good();c.RedTeamPassed=false;if Assess(c).Promote{t.Fatal("red-team failure promoted")}}
func TestCriticalRegressionCannotPromote(t *testing.T){c:=good();c.Evals[0].Candidate=.7;if Assess(c).Promote{t.Fatal("regression promoted")}}
