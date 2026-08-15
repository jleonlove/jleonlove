package intent
import"testing"
func TestLowRiskClearProceeds(t *testing.T){d:=Resolve(Request{Objective:"read docs",Risk:Low,Target:"docs"});if !d.Proceed{t.Fatal("blocked")}}
func TestHighRiskMissingTargetClarifies(t *testing.T){d:=Resolve(Request{Objective:"deploy",Risk:High});if d.Proceed||len(d.Options)==0{t.Fatal("guessed")}}
func TestDestructiveMissingTargetClarifies(t *testing.T){if Resolve(Request{Objective:"delete",Risk:Irreversible,Destructive:true}).Proceed{t.Fatal("guessed destructive target")}}
func TestIntentDigestBindsDecision(t *testing.T){r:=Request{Objective:"deploy",Risk:High,Target:"staging"};d:=Resolve(r);mut:=r;mut.Target="prod";if Confirm(mut,mut,d.IntentDigest)==nil{t.Fatal("mutated intent accepted")}}
