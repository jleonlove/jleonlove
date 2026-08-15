package constraints
import("errors";"math";"testing")
func TestVerified(t *testing.T){if e:=Verify(Model{"cost":90,"qty":10},[]Constraint{{"cost",LE,100},{"qty",GE,5}});e!=nil{t.Fatal(e)}}
func TestViolation(t *testing.T){if e:=Verify(Model{"cost":110},[]Constraint{{"cost",LE,100}});!errors.Is(e,ErrUnsatisfied){t.Fatal(e)}}
func TestUnknown(t *testing.T){if e:=Verify(Model{},[]Constraint{{"cost",LE,100}});!errors.Is(e,ErrUnknown){t.Fatal(e)}}
func TestNaNRejected(t *testing.T){if e:=Verify(Model{"x":math.NaN()},[]Constraint{{"x",LE,1}});!errors.Is(e,ErrInvalid){t.Fatal(e)}}
func TestFeasible(t *testing.T){if !Feasible(map[string][2]float64{"price":{80,120}},[]Constraint{{"price",GE,90},{"price",LE,100}}){t.Fatal("expected feasible")}}
func TestInfeasible(t *testing.T){if Feasible(map[string][2]float64{"price":{80,85}},[]Constraint{{"price",GE,90}}){t.Fatal("expected infeasible")}}
