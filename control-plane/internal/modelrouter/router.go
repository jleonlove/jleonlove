package modelrouter
import("context";"errors";"sort";"sync")
var(ErrNoModel=errors.New("no eligible model");ErrBudget=errors.New("model budget exceeded");ErrProvider=errors.New("provider failure"))
type Request struct{Prompt string;NeedVision,NeedTools bool;MaxCost float64}
type Response struct{Text,Provider,Model string;Cost float64}
type Model struct{Provider,Name string;Vision,Tools bool;EstimatedCost float64;Priority int}
type Provider interface{Generate(context.Context,Model,Request)(Response,error)}
type Router struct{mu sync.RWMutex;models []Model;providers map[string]Provider}
func New()*Router{return &Router{providers:map[string]Provider{}}}
func(r *Router)RegisterProvider(n string,p Provider){r.mu.Lock();defer r.mu.Unlock();r.providers[n]=p}
func(r *Router)SetModels(ms []Model){r.mu.Lock();defer r.mu.Unlock();r.models=append([]Model(nil),ms...)}
func(r *Router)Generate(ctx context.Context,q Request)(Response,error){r.mu.RLock();ms:=append([]Model(nil),r.models...);ps:=map[string]Provider{};for k,v:=range r.providers{ps[k]=v};r.mu.RUnlock()
 sort.Slice(ms,func(i,j int)bool{if ms[i].Priority==ms[j].Priority{return ms[i].EstimatedCost<ms[j].EstimatedCost};return ms[i].Priority>ms[j].Priority})
 eligible:=0;budgetBlocked:=false
 for _,m:=range ms{if q.NeedVision&&!m.Vision||q.NeedTools&&!m.Tools{continue};eligible++;if q.MaxCost>0&&m.EstimatedCost>q.MaxCost{budgetBlocked=true;continue};p:=ps[m.Provider];if p==nil{continue};out,e:=p.Generate(ctx,m,q);if e==nil{out.Provider=m.Provider;out.Model=m.Name;return out,nil}}
 if eligible>0&&budgetBlocked{return Response{},ErrBudget};return Response{},ErrNoModel}
