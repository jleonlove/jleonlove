package logisticstwin
import("errors";"sort")
var ErrRoute=errors.New("invalid route")
type Leg struct{From,To string;Hours,Cost,Risk float64}
type Route struct{Legs []Leg}
type Score struct{Hours,Cost,Risk float64}
func Evaluate(r Route)(Score,error){var s Score;if len(r.Legs)==0{return s,ErrRoute};for i,l:=range r.Legs{if l.From==""||l.To==""||l.Hours<0||l.Cost<0||l.Risk<0{return Score{},ErrRoute};if i>0&&r.Legs[i-1].To!=l.From{return Score{},ErrRoute};s.Hours+=l.Hours;s.Cost+=l.Cost;s.Risk+=l.Risk};return s,nil}
func Rank(rs map[string]Route)[]string{type x struct{id string;s float64};a:=[]x{};for id,r:=range rs{v,e:=Evaluate(r);if e==nil{a=append(a,x{id,v.Hours+v.Cost+v.Risk})}};sort.Slice(a,func(i,j int)bool{if a[i].s==a[j].s{return a[i].id<a[j].id};return a[i].s<a[j].s});o:=[]string{};for _,v:=range a{o=append(o,v.id)};return o}