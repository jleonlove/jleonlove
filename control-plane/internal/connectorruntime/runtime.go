package connectorruntime
import("context";"errors";"sort";"sync";"time")
var(ErrDenied=errors.New("connector invocation denied");ErrMissing=errors.New("connector missing");ErrTimeout=errors.New("connector timeout");ErrEvidence=errors.New("connector evidence missing"))
type Request struct{Tenant,Connector,Action string;Scopes []string;Evidence string;Timeout time.Duration;Input map[string]string}
type Result struct{Connector,Action string;Output map[string]string}
type Connector interface{Invoke(context.Context,Request)(Result,error)}
type Policy func(Request)bool
type Runtime struct{mu sync.RWMutex;connectors map[string]Connector;policy Policy}
func New(p Policy)*Runtime{return &Runtime{connectors:map[string]Connector{},policy:p}}
func(r *Runtime)Register(n string,c Connector){r.mu.Lock();defer r.mu.Unlock();r.connectors[n]=c}
func(r *Runtime)Invoke(ctx context.Context,q Request)(Result,error){
 if q.Tenant==""||q.Connector==""||q.Action==""{return Result{},ErrDenied};if q.Evidence==""{return Result{},ErrEvidence};if r.policy==nil||!r.policy(q){return Result{},ErrDenied}
 r.mu.RLock();c:=r.connectors[q.Connector];r.mu.RUnlock();if c==nil{return Result{},ErrMissing}
 d:=q.Timeout;if d<=0||d>30*time.Second{d=30*time.Second};cx,cancel:=context.WithTimeout(ctx,d);defer cancel()
 out,e:=c.Invoke(cx,q);if errors.Is(cx.Err(),context.DeadlineExceeded){return Result{},ErrTimeout};return out,e
}
func HasScope(q Request,s string)bool{x:=append([]string(nil),q.Scopes...);sort.Strings(x);i:=sort.SearchStrings(x,s);return i<len(x)&&x[i]==s}
