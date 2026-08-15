package tradefinance
import("testing";"errors")
func TestDraw(t *testing.T){f:=Facility{1000,0,800,.25};if Available(f)!=600{t.Fatal()};if e:=Draw(&f,500);e!=nil{t.Fatal(e)}}
func TestOverdraw(t *testing.T){f:=Facility{1000,0,100,.5};if e:=Draw(&f,60);!errors.Is(e,ErrCollateral){t.Fatal(e)}}
func TestPaymentGate(t *testing.T){if e:=Release(Payment{100,true,true,true});e!=nil{t.Fatal(e)};if e:=Release(Payment{100,true,false,true});!errors.Is(e,ErrPayment){t.Fatal(e)}}