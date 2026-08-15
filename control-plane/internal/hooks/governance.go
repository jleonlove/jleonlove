package hooks
import("context";"errors";"sync";"time")
var(ErrUnknownHook=errors.New("unknown hook");ErrRevokedHook=errors.New("hook revoked");ErrTriggerDenied=errors.New("trigger denied");ErrCapabilityDenied=errors.New("hook capability denied");ErrRecursion=errors.New("hook recursion denied");ErrBudget=errors.New("hook budget exceeded"))
type Hook struct{ID string;Approved bool;Revoked bool;Triggers map[string]bool;RequiredCapabilities []string;MaxDepth int;Timeout time.Duration}
type Request struct{HookID string;Trigger string;Capabilities []string;Depth int}
type Runner interface{Run(context.Context,Request)error}
type Governor struct{mu sync.RWMutex;hooks map[string]Hook;runner Runner}
func New(r Runner)*Governor{return &Governor{hooks:map[string]Hook{},runner:r}}
func(g *Governor)Register(h Hook){g.mu.Lock();defer g.mu.Unlock();g.hooks[h.ID]=h}
func(g *Governor)Revoke(id string){g.mu.Lock();defer g.mu.Unlock();h,ok:=g.hooks[id];if !ok{return};h.Revoked=true;g.hooks[id]=h}
func all(got,need []string)bool{m:=map[string]bool{};for _,x:=range got{m[x]=true};for _,x:=range need{if !m[x]{return false}};return true}
func(g *Governor)Execute(ctx context.Context,r Request)error{g.mu.RLock();h,ok:=g.hooks[r.HookID];g.mu.RUnlock();if !ok||!h.Approved{return ErrUnknownHook};if h.Revoked{return ErrRevokedHook};if !h.Triggers[r.Trigger]{return ErrTriggerDenied};if !all(r.Capabilities,h.RequiredCapabilities){return ErrCapabilityDenied};if r.Depth>h.MaxDepth{return ErrRecursion};if h.Timeout<=0{return ErrBudget};c,cancel:=context.WithTimeout(ctx,h.Timeout);defer cancel();err:=g.runner.Run(c,r);if errors.Is(c.Err(),context.DeadlineExceeded){return ErrBudget};return err}
