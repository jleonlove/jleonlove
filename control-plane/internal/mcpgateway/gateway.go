package mcpgateway
import("context";"errors";"sync")
var(ErrUnknownServer=errors.New("unknown mcp server");ErrRevokedServer=errors.New("mcp server revoked");ErrToolDenied=errors.New("mcp tool denied");ErrSchemaMismatch=errors.New("mcp schema mismatch");ErrCapabilityDenied=errors.New("mcp capability denied");ErrCredentialRequest=errors.New("raw credential access denied"))
type Tool struct{Name string;SchemaDigest string;RequiredCapabilities []string}
type Server struct{ID string;Approved bool;Revoked bool;Tools map[string]Tool}
type Request struct{ServerID string;ToolName string;SchemaDigest string;Capabilities []string;Arguments map[string]any}
type Response struct{Data map[string]any}
type Transport interface{Call(context.Context,Request)(Response,error)}
type Evidence interface{Record(context.Context,string,Request,error)}
type Gateway struct{mu sync.RWMutex;servers map[string]Server;transport Transport;evidence Evidence}
func New(t Transport,e Evidence)*Gateway{return &Gateway{servers:map[string]Server{},transport:t,evidence:e}}
func(g *Gateway)Register(s Server){g.mu.Lock();defer g.mu.Unlock();g.servers[s.ID]=s}
func(g *Gateway)Revoke(id string){g.mu.Lock();defer g.mu.Unlock();s,ok:=g.servers[id];if !ok{return};s.Revoked=true;g.servers[id]=s}
func hasAll(got,need []string)bool{m:=map[string]bool{};for _,v:=range got{m[v]=true};for _,v:=range need{if !m[v]{return false}};return true}
func(g *Gateway)Discover(id string)[]Tool{g.mu.RLock();defer g.mu.RUnlock();s,ok:=g.servers[id];if !ok||!s.Approved||s.Revoked{return nil};out:=make([]Tool,0,len(s.Tools));for _,t:=range s.Tools{out=append(out,t)};return out}
func(g *Gateway)Call(ctx context.Context,r Request)(Response,error){g.mu.RLock();s,ok:=g.servers[r.ServerID];g.mu.RUnlock();deny:=func(e error)(Response,error){if g.evidence!=nil{g.evidence.Record(ctx,"DENY",r,e)};return Response{},e};if !ok||!s.Approved{return deny(ErrUnknownServer)};if s.Revoked{return deny(ErrRevokedServer)};t,ok:=s.Tools[r.ToolName];if !ok{return deny(ErrToolDenied)};if t.SchemaDigest!=r.SchemaDigest{return deny(ErrSchemaMismatch)};if !hasAll(r.Capabilities,t.RequiredCapabilities){return deny(ErrCapabilityDenied)};for k:=range r.Arguments{if k=="api_key"||k=="token"||k=="password"||k=="secret"{return deny(ErrCredentialRequest)}};resp,err:=g.transport.Call(ctx,r);if g.evidence!=nil{g.evidence.Record(ctx,"CALL",r,err)};return resp,err}
