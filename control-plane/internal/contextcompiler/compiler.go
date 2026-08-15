package contextcompiler
import("errors";"sort";"strings")
var ErrBudget=errors.New("context budget exceeded")
type Trust int
const(Untrusted Trust=iota;Memory;Evidence;Policy)
type Item struct{ID,Text string;Trust Trust;Relevance int;Tokens int;Instruction bool}
type Compiled struct{Items []Item;Tokens int}
func Compile(items []Item,budget int)(Compiled,error){
 for i:=range items{if items[i].Trust==Untrusted{items[i].Instruction=false;items[i].Text="[UNTRUSTED DATA] "+items[i].Text}}
 sort.SliceStable(items,func(i,j int)bool{if items[i].Trust==items[j].Trust{return items[i].Relevance>items[j].Relevance};return items[i].Trust>items[j].Trust})
 out:=Compiled{}
 for _,x:=range items{if x.Tokens<0{continue};if out.Tokens+x.Tokens>budget{continue};out.Items=append(out.Items,x);out.Tokens+=x.Tokens}
 if out.Tokens>budget{return Compiled{},ErrBudget};return out,nil
}
func HasExecutableUntrusted(c Compiled)bool{for _,x:=range c.Items{if x.Trust==Untrusted&&x.Instruction{return true};if x.Trust==Untrusted&&!strings.HasPrefix(x.Text,"[UNTRUSTED DATA]"){return true}};return false}
