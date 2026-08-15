package tradepreflight
import"errors"
var ErrPreflight=errors.New("trade preflight failed")
type Check struct{Name string;Pass bool;Evidence string}
func Run(cs []Check)error{if len(cs)==0{return ErrPreflight};for _,c:=range cs{if !c.Pass||c.Evidence==""{return ErrPreflight}};return nil}