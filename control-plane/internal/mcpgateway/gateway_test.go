package mcpgateway
import("context";"errors";"sync/atomic";"testing")
type ft struct{calls atomic.Int64};func(f *ft)Call(context.Context,Request)(Response,error){f.calls.Add(1);return Response{},nil}
type fe struct{};func(*fe)Record(context.Context,string,Request,error){}
func fix()(*Gateway,*ft){t:=&ft{};g:=New(t,&fe{});g.Register(Server{ID:"github",Approved:true,Tools:map[string]Tool{"repo.read":{Name:"repo.read",SchemaDigest:"v1",RequiredCapabilities:[]string{"repository.read"}}}});return g,t}
func req()Request{return Request{ServerID:"github",ToolName:"repo.read",SchemaDigest:"v1",Capabilities:[]string{"repository.read"},Arguments:map[string]any{"repo":"atlas"}}}
func TestApproved(t *testing.T){g,tr:=fix();if _,e:=g.Call(context.Background(),req());e!=nil{t.Fatal(e)};if tr.calls.Load()!=1{t.Fatal("not called")}}
func TestToolDenied(t *testing.T){g,tr:=fix();r:=req();r.ToolName="repo.write";if _,e:=g.Call(context.Background(),r);!errors.Is(e,ErrToolDenied){t.Fatal(e)};if tr.calls.Load()!=0{t.Fatal("transport called")}}
func TestSchemaDenied(t *testing.T){g,tr:=fix();r:=req();r.SchemaDigest="evil";if _,e:=g.Call(context.Background(),r);!errors.Is(e,ErrSchemaMismatch){t.Fatal(e)};if tr.calls.Load()!=0{t.Fatal("transport called")}}
func TestCapabilityDenied(t *testing.T){g,tr:=fix();r:=req();r.Capabilities=nil;if _,e:=g.Call(context.Background(),r);!errors.Is(e,ErrCapabilityDenied){t.Fatal(e)};if tr.calls.Load()!=0{t.Fatal("transport called")}}
func TestRevoked(t *testing.T){g,tr:=fix();g.Revoke("github");if _,e:=g.Call(context.Background(),req());!errors.Is(e,ErrRevokedServer){t.Fatal(e)};if tr.calls.Load()!=0{t.Fatal("transport called")}}
func TestCredentialSmugglingDenied(t *testing.T){g,tr:=fix();r:=req();r.Arguments["api_key"]="x";if _,e:=g.Call(context.Background(),r);!errors.Is(e,ErrCredentialRequest){t.Fatal(e)};if tr.calls.Load()!=0{t.Fatal("transport called")}}
func TestUnapprovedDiscoveryHidden(t *testing.T){tr:=&ft{};g:=New(tr,nil);g.Register(Server{ID:"bad",Tools:map[string]Tool{"x":{Name:"x"}}});if len(g.Discover("bad"))!=0{t.Fatal("disclosed")}}
