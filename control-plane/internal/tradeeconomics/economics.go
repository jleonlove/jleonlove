package tradeeconomics
import("errors";"math")
var ErrInput=errors.New("invalid economics input")
type Deal struct{Quantity,Price,Freight,Insurance,Duty,Finance,FX float64}
type Result struct{Gross,Costs,Landed,Margin float64}
func Calculate(d Deal)(Result,error){v:=[]float64{d.Quantity,d.Price,d.Freight,d.Insurance,d.Duty,d.Finance,d.FX};for _,x:=range v{if math.IsNaN(x)||math.IsInf(x,0)||x<0{return Result{},ErrInput}};if d.FX==0{return Result{},ErrInput};g:=d.Quantity*d.Price*d.FX;c:=(d.Freight+d.Insurance+d.Duty+d.Finance)*d.FX;return Result{g,c,g+c,-c},nil}