package productionreadiness
import"testing"
func good()[]Check{n:=[]string{"live_integrations","regulatory_data","observability","load_chaos","security_assessment","disaster_recovery","red_team","end_to_end_trade"};x:=[]Check{};for _,v:=range n{x=append(x,Check{v,true,"signed-evidence",true})};return x}
func TestReady(t *testing.T){r:=Evaluate(good());if !r.Ready||len(r.Failures)!=0{t.Fatal(r)}}
func TestMissingEvidenceFails(t *testing.T){x:=good();x[0].Evidence="";r:=Evaluate(x);if r.Ready{t.Fatal("missing evidence passed")}}
func TestSecurityFailureFails(t *testing.T){x:=good();x[4].Pass=false;r:=Evaluate(x);if r.Ready{t.Fatal("security failure passed")}}
func TestMissingRequiredGateFails(t *testing.T){x:=good()[1:];r:=Evaluate(x);if r.Ready{t.Fatal("missing required gate passed")}}
