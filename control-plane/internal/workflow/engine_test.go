package workflow
import("errors";"testing")
func mission()Mission{return Mission{ID:"m",AuthorityEpoch:7,Steps:map[string]Step{"plan":{ID:"plan",State:Pending},"execute":{ID:"execute",DependsOn:[]string{"plan"},State:Pending}}}}
func TestDependencyOrder(t *testing.T){m:=mission();r:=Ready(m);if len(r)!=1||r[0]!="plan"{t.Fatalf("%v",r)};if e:=Start(&m,"execute","w",7);!errors.Is(e,ErrDependency){t.Fatal(e)}}
func TestCheckpointPauseResume(t *testing.T){m:=mission();if e:=Start(&m,"plan","w1",7);e!=nil{t.Fatal(e)};if e:=Checkpoint(&m,"plan","w1","cp1");e!=nil{t.Fatal(e)};if e:=Pause(&m,"plan","w1");e!=nil{t.Fatal(e)};if e:=Resume(&m,"plan","w2",7);e!=nil{t.Fatal(e)};if m.Steps["plan"].Checkpoint!="cp1"{t.Fatal("checkpoint lost")}}
func TestStaleAuthorityCannotResume(t *testing.T){m:=mission();_ = Start(&m,"plan","w",7);_ = Checkpoint(&m,"plan","w","cp");_ = Pause(&m,"plan","w");m.AuthorityEpoch=8;if e:=Resume(&m,"plan","w2",7);!errors.Is(e,ErrAuthorityRefresh){t.Fatal(e)}}
func TestLeaseOwnerEnforced(t *testing.T){m:=mission();_ = Start(&m,"plan","w1",7);if e:=Complete(&m,"plan","w2");!errors.Is(e,ErrLease){t.Fatal(e)}}
func TestCancellationIsTerminal(t *testing.T){m:=mission();Cancel(&m);if len(Ready(m))!=0{t.Fatal("cancelled mission ready")};if e:=Start(&m,"plan","w",7);!errors.Is(e,ErrTerminal){t.Fatal(e)}}
func TestCompletionUnlocksNext(t *testing.T){m:=mission();_ = Start(&m,"plan","w",7);_ = Complete(&m,"plan","w");r:=Ready(m);if len(r)!=1||r[0]!="execute"{t.Fatalf("%v",r)}}
