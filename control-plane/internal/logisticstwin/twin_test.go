package logisticstwin
import("testing";"errors")
func TestRoute(t *testing.T){s,e:=Evaluate(Route{[]Leg{{"A","B",10,20,1},{"B","C",5,10,2}}});if e!=nil||s.Hours!=15||s.Cost!=30||s.Risk!=3{t.Fatal(s,e)}}
func TestBrokenRoute(t *testing.T){_,e:=Evaluate(Route{[]Leg{{"A","B",1,1,1},{"X","C",1,1,1}}});if !errors.Is(e,ErrRoute){t.Fatal(e)}}
func TestRank(t *testing.T){r:=map[string]Route{"b":{[]Leg{{"A","B",2,2,2}}},"a":{[]Leg{{"A","B",1,1,1}}}};if Rank(r)[0]!="a"{t.Fatal("rank")}}