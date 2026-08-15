package hooks
import("context";"errors";"sync/atomic";"testing";"time")
type runner struct{calls atomic.Int64;block bool}
func(r *runner)Run(ctx context.Context,_ Request)error{r.calls.Add(1);if r.block{<-ctx.Done();return ctx.Err()};return nil}
func fx()(*Governor,*runner){r:=&runner{};g:=New(r);g.Register(Hook{ID:"audit",Approved:true,Triggers:map[string]bool{"post-tool":true},RequiredCapabilities:[]string{"audit.write"},MaxDepth:1,Timeout:time.Second});return g,r}
func req()Request{return Request{HookID:"audit",Trigger:"post-tool",Capabilities:[]string{"audit.write"},Depth:1}}
func TestApproved(t *testing.T){g,r:=fx();if e:=g.Execute(context.Background(),req());e!=nil{t.Fatal(e)};if r.calls.Load()!=1{t.Fatal("not run")}}
func TestTriggerDenied(t *testing.T){g,r:=fx();q:=req();q.Trigger="pre-auth";if e:=g.Execute(context.Background(),q);!errors.Is(e,ErrTriggerDenied){t.Fatal(e)};if r.calls.Load()!=0{t.Fatal("ran")}}
func TestCapabilityDenied(t *testing.T){g,r:=fx();q:=req();q.Capabilities=nil;if e:=g.Execute(context.Background(),q);!errors.Is(e,ErrCapabilityDenied){t.Fatal(e)};if r.calls.Load()!=0{t.Fatal("ran")}}
func TestRecursionDenied(t *testing.T){g,r:=fx();q:=req();q.Depth=2;if e:=g.Execute(context.Background(),q);!errors.Is(e,ErrRecursion){t.Fatal(e)};if r.calls.Load()!=0{t.Fatal("ran")}}
func TestRevoked(t *testing.T){g,r:=fx();g.Revoke("audit");if e:=g.Execute(context.Background(),req());!errors.Is(e,ErrRevokedHook){t.Fatal(e)};if r.calls.Load()!=0{t.Fatal("ran")}}
func TestTimeout(t *testing.T){r:=&runner{block:true};g:=New(r);g.Register(Hook{ID:"slow",Approved:true,Triggers:map[string]bool{"x":true},MaxDepth:0,Timeout:time.Millisecond});e:=g.Execute(context.Background(),Request{HookID:"slow",Trigger:"x"});if !errors.Is(e,ErrBudget){t.Fatalf("got %v",e)}}
