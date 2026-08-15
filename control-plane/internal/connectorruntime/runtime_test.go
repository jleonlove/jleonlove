package connectorruntime
import("context";"errors";"testing";"time")
type fake struct{delay time.Duration}
func(f fake)Invoke(ctx context.Context,q Request)(Result,error){select{case<-time.After(f.delay):return Result{q.Connector,q.Action,map[string]string{"ok":"true"}},nil;case<-ctx.Done():return Result{},ctx.Err()}}
func allow(q Request)bool{return HasScope(q,"trade:read")}
func base()Request{return Request{Tenant:"t",Connector:"docs",Action:"read",Scopes:[]string{"trade:read"},Evidence:"ev",Timeout:time.Second}}
func TestInvoke(t *testing.T){r:=New(allow);r.Register("docs",fake{});x,e:=r.Invoke(context.Background(),base());if e!=nil||x.Output["ok"]!="true"{t.Fatal(x,e)}}
func TestDenied(t *testing.T){r:=New(allow);q:=base();q.Scopes=nil;if _,e:=r.Invoke(context.Background(),q);!errors.Is(e,ErrDenied){t.Fatal(e)}}
func TestEvidence(t *testing.T){r:=New(allow);q:=base();q.Evidence="";if _,e:=r.Invoke(context.Background(),q);!errors.Is(e,ErrEvidence){t.Fatal(e)}}
func TestMissing(t *testing.T){r:=New(allow);if _,e:=r.Invoke(context.Background(),base());!errors.Is(e,ErrMissing){t.Fatal(e)}}
func TestTimeout(t *testing.T){r:=New(allow);r.Register("docs",fake{50*time.Millisecond});q:=base();q.Timeout=time.Millisecond;if _,e:=r.Invoke(context.Background(),q);!errors.Is(e,ErrTimeout){t.Fatal(e)}}
