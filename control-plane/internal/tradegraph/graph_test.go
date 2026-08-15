package tradegraph
import("errors";"testing";"time")
func graph()*Graph{g:=New();g.AddEntity(Entity{"buyer","PARTY"});g.AddEntity(Entity{"seller","PARTY"});g.AddEntity(Entity{"cargo","COMMODITY"});return g}
func TestRelation(t *testing.T){g:=graph();if e:=g.Link(Relation{"seller","SELLS","cargo"});e!=nil{t.Fatal(e)}}
func TestMissingRelationEntity(t *testing.T){g:=graph();if e:=g.Link(Relation{"ghost","SELLS","cargo"});!errors.Is(e,ErrMissingEntity){t.Fatal(e)}}
func TestAppendTimeline(t *testing.T){g:=graph();a:=Event{"e1","tx1","CONTRACTED",time.Unix(10,0),[]string{"buyer","seller","cargo"},[]string{"sha256:a"}};b:=Event{"e2","tx1","SHIPPED",time.Unix(20,0),[]string{"cargo"},[]string{"sha256:b"}};if e:=g.Append(a);e!=nil{t.Fatal(e)};if e:=g.Append(b);e!=nil{t.Fatal(e)};if x:=g.Timeline("tx1");len(x)!=2||x[0].ID!="e1"{t.Fatal(x)}}
func TestDuplicateEvent(t *testing.T){g:=graph();e:=Event{"e1","tx", "X",time.Unix(1,0),[]string{"cargo"},nil};_ = g.Append(e);if x:=g.Append(e);!errors.Is(x,ErrDuplicateEvent){t.Fatal(x)}}
func TestMissingEventEntity(t *testing.T){g:=graph();if e:=g.Append(Event{"e","tx","X",time.Unix(1,0),[]string{"ghost"},nil});!errors.Is(e,ErrMissingEntity){t.Fatal(e)}}
func TestTimeRegression(t *testing.T){g:=graph();_ = g.Append(Event{"e1","tx","X",time.Unix(20,0),[]string{"cargo"},nil});if e:=g.Append(Event{"e2","tx","Y",time.Unix(10,0),[]string{"cargo"},nil});!errors.Is(e,ErrTimeRegression){t.Fatal(e)}}
