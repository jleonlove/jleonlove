package productiongate
import("testing";"errors")
func good()Evidence{m:=map[string]bool{};for i:=71;i<=79;i++{m[fmtMilestone(i)]=true};return Evidence{m,true,true,true,true,true,true}}
func TestGatePass(t *testing.T){if e:=Evaluate(good());e!=nil{t.Fatal(e)}}
func TestGateFailClosed(t *testing.T){x:=good();x.Security=false;if e:=Evaluate(x);!errors.Is(e,ErrGate){t.Fatal(e)}}
func TestMissingMilestone(t *testing.T){x:=good();delete(x.Milestones,"000075");if e:=Evaluate(x);!errors.Is(e,ErrGate){t.Fatal(e)}}