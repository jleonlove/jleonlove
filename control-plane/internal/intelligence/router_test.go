package intelligence
import("errors";"testing")
func models()[]Model{return []Model{
{Name:"fast",Quality:6,Reasoning:4,Privacy:2,Qualified:true,Cost:.1,Latency:.2,Tools:true},
{Name:"deep",Quality:9,Reasoning:9,Privacy:4,Qualified:true,Cost:1,Latency:2,Tools:true,Vision:true},
{Name:"unqualified",Quality:10,Reasoning:10,Privacy:5,Qualified:false,Cost:.01,Latency:.1,Tools:true,Vision:true},
}}
func TestRoutesSimpleCheap(t *testing.T){d,e:=Route(Task{Complexity:0,Risk:0,Privacy:1,MaxCost:.5,MaxLatency:1},models());if e!=nil||d.Model!="fast"{t.Fatalf("%+v %v",d,e)}}
func TestRoutesComplexToDeep(t *testing.T){d,e:=Route(Task{Complexity:5,Risk:2,Privacy:3,MaxCost:2,MaxLatency:3,NeedsVision:true},models());if e!=nil||d.Model!="deep"{t.Fatalf("%+v %v",d,e)}}
func TestUnqualifiedNeverSelected(t *testing.T){d,e:=Route(Task{Complexity:9,Risk:5,Privacy:5,MaxCost:10,MaxLatency:10,NeedsVision:true},models());if !errors.Is(e,ErrNoQualifiedModel)||d.Model!=""{t.Fatal("unqualified model selected")}}
func TestBudgetIncreasesWithDifficulty(t *testing.T){a,_:=Route(Task{Complexity:0,Privacy:1,MaxCost:2,MaxLatency:3},models());b,_:=Route(Task{Complexity:5,Privacy:1,MaxCost:2,MaxLatency:3},models());if b.ReasoningBudget<=a.ReasoningBudget{t.Fatal("budget did not increase")}}
func TestPrivacyGate(t *testing.T){_,e:=Route(Task{Privacy:5,MaxCost:2,MaxLatency:3},models());if !errors.Is(e,ErrNoQualifiedModel){t.Fatal("privacy gate bypassed")}}
