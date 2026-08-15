package commodityworld
import("errors";"testing")
func model()Model{return Model{Nodes:map[string]Node{"mine":{"mine","SUPPLY"},"inventory":{"inventory","INVENTORY"},"price":{"price","PRICE"}},Edges:[]Edge{{"mine","inventory",1},{"inventory","price",-0.5}}}}
func TestCausalPropagation(t *testing.T){v,e:=model().Propagate(map[string]float64{},[]Observation{{"mine","source:1",-10}});if e!=nil||v["inventory"]!=-10||v["price"]!=5{t.Fatalf("%v %v",v,e)}}
func TestEvidenceRequired(t *testing.T){_,e:=model().Propagate(nil,[]Observation{{"mine","",-1}});if !errors.Is(e,ErrEvidence){t.Fatal(e)}}
func TestUnknownNode(t *testing.T){_,e:=model().Propagate(nil,[]Observation{{"ghost","ev",-1}});if !errors.Is(e,ErrNode){t.Fatal(e)}}
func TestCycleRejected(t *testing.T){m:=model();m.Edges=append(m.Edges,Edge{"price","mine",1});if e:=m.Validate();!errors.Is(e,ErrCycle){t.Fatal(e)}}
func TestProbabilisticForecast(t *testing.T){f,e:=model().Expected(nil,[]Scenario{{"normal",.7,map[string]float64{"mine":0}},{"disruption",.3,map[string]float64{"mine":-10}}});if e!=nil||f.Expected["price"]!=1.5{t.Fatalf("%v %v",f,e)}}
func TestBadProbability(t *testing.T){_,e:=model().Expected(nil,[]Scenario{{"x",.8,nil}});if !errors.Is(e,ErrProbability){t.Fatal(e)}}
