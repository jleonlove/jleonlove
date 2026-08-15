package tradelab
type Case struct{Name string;ExpectedBlock bool;Run func()bool}
type Report struct{Passed,Failed int;Failures []string}
func Execute(cs []Case)Report{r:=Report{};for _,c:=range cs{blocked:=c.Run();if blocked==c.ExpectedBlock{r.Passed++}else{r.Failed++;r.Failures=append(r.Failures,c.Name)}};return r}