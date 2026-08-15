package constraints
import("errors";"math")
var(ErrUnsatisfied=errors.New("constraint unsatisfied");ErrUnknown=errors.New("unknown variable");ErrInvalid=errors.New("invalid numeric value"))
type Op string
const(LE Op="<=";GE Op=">=";EQ Op="==")
type Constraint struct{Variable string;Op Op;Value float64}
type Model map[string]float64
func Verify(m Model,cs []Constraint)error{
 for _,c:=range cs{v,ok:=m[c.Variable];if !ok{return ErrUnknown};if math.IsNaN(v)||math.IsInf(v,0)||math.IsNaN(c.Value)||math.IsInf(c.Value,0){return ErrInvalid}
  good:=false;switch c.Op{case LE:good=v<=c.Value;case GE:good=v>=c.Value;case EQ:good=v==c.Value;default:return ErrInvalid};if !good{return ErrUnsatisfied}}
 return nil
}
func Feasible(bounds map[string][2]float64,cs []Constraint)bool{
 for _,c:=range cs{b,ok:=bounds[c.Variable];if !ok{return false};switch c.Op{case LE:if b[0]>c.Value{return false};case GE:if b[1]<c.Value{return false};case EQ:if c.Value<b[0]||c.Value>b[1]{return false};default:return false}}
 return true
}
