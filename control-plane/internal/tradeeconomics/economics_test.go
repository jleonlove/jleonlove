package tradeeconomics
import("testing";"errors";"math")
func TestCalc(t *testing.T){r,e:=Calculate(Deal{10,100,20,10,5,5,1});if e!=nil||r.Gross!=1000||r.Costs!=40||r.Landed!=1040{t.Fatal(r,e)}}
func TestBad(t *testing.T){_,e:=Calculate(Deal{Quantity:math.NaN(),FX:1});if !errors.Is(e,ErrInput){t.Fatal(e)}}