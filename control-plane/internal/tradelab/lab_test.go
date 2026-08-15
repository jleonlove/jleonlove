package tradelab
import"testing"
func TestFraudCases(t *testing.T){r:=Execute([]Case{{"fake mandate",true,func()bool{return true}},{"clean trade",false,func()bool{return false}}});if r.Failed!=0||r.Passed!=2{t.Fatal(r)}}