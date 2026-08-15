package modelrouter
import("context";"errors";"testing")
type fake struct{fail bool}
func(f fake)Generate(_ context.Context,m Model,_ Request)(Response,error){if f.fail{return Response{},ErrProvider};return Response{Text:"ok",Cost:m.EstimatedCost},nil}
func TestRoute(t *testing.T){r:=New();r.RegisterProvider("a",fake{});r.SetModels([]Model{{"a","strong",true,true,.5,10}});x,e:=r.Generate(context.Background(),Request{NeedTools:true,MaxCost:1});if e!=nil||x.Model!="strong"{t.Fatal(x,e)}}
func TestCapability(t *testing.T){r:=New();r.RegisterProvider("a",fake{});r.SetModels([]Model{{"a","text",false,true,.1,1}});_,e:=r.Generate(context.Background(),Request{NeedVision:true});if !errors.Is(e,ErrNoModel){t.Fatal(e)}}
func TestBudget(t *testing.T){r:=New();r.RegisterProvider("a",fake{});r.SetModels([]Model{{"a","x",true,true,2,1}});_,e:=r.Generate(context.Background(),Request{MaxCost:1});if !errors.Is(e,ErrBudget){t.Fatal(e)}}
func TestFallback(t *testing.T){r:=New();r.RegisterProvider("bad",fake{true});r.RegisterProvider("good",fake{});r.SetModels([]Model{{"bad","primary",true,true,.2,10},{"good","fallback",true,true,.3,5}});x,e:=r.Generate(context.Background(),Request{});if e!=nil||x.Model!="fallback"{t.Fatal(x,e)}}
