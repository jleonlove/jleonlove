package provenance
import("testing";"errors")
func TestChain(t *testing.T){a:=Event{"cargo","mine","CL","ev1",""};b:=Event{"cargo","carrier","sea","ev2",Hash(a)};if e:=Verify([]Event{a,b});e!=nil{t.Fatal(e)}}
func TestTamper(t *testing.T){a:=Event{"cargo","mine","CL","ev1",""};b:=Event{"cargo","carrier","sea","ev2","bad"};if e:=Verify([]Event{a,b});!errors.Is(e,ErrChain){t.Fatal(e)}}