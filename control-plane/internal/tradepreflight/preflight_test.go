package tradepreflight
import("testing";"errors")
func TestPass(t *testing.T){if e:=Run([]Check{{"docs",true,"ev"},{"compliance",true,"ev"}});e!=nil{t.Fatal(e)}}
func TestEvidence(t *testing.T){if e:=Run([]Check{{"docs",true,""}});!errors.Is(e,ErrPreflight){t.Fatal(e)}}