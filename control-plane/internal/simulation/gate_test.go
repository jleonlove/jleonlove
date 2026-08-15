package simulation
import("errors";"testing")
func req()Request{return Request{Evidence:true,MaxExpectedLoss:50,Scenarios:[]Scenario{{"base",.7,10},{"stress",.3,40}}}}
func TestPass(t *testing.T){r,e:=Evaluate(req());if e!=nil||r.ExpectedLoss!=19{t.Fatalf("%+v %v",r,e)}}
func TestProbabilitiesMustSum(t *testing.T){x:=req();x.Scenarios[0].Probability=.5;if _,e:=Evaluate(x);!errors.Is(e,ErrModel){t.Fatal(e)}}
func TestRiskGate(t *testing.T){x:=req();x.MaxExpectedLoss=10;if _,e:=Evaluate(x);!errors.Is(e,ErrRisk){t.Fatal(e)}}
func TestEvidenceGate(t *testing.T){x:=req();x.Evidence=false;if _,e:=Evaluate(x);!errors.Is(e,ErrEvidence){t.Fatal(e)}}
func TestConsequentialApproval(t *testing.T){x:=req();x.Consequential=true;if _,e:=Evaluate(x);!errors.Is(e,ErrApproval){t.Fatal(e)}}
func TestApprovedConsequential(t *testing.T){x:=req();x.Consequential=true;x.Approved=true;if _,e:=Evaluate(x);e!=nil{t.Fatal(e)}}
