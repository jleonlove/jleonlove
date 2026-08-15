package transactionpassport
import"testing"
func TestSeal(t *testing.T){p:=Passport{"tx",[]string{"b","a"},true,""};if e:=Seal(&p);e!=nil||!Verify(p){t.Fatal(e,p)}}
func TestTamper(t *testing.T){p:=Passport{"tx",[]string{"a"},true,""};_ = Seal(&p);p.Evidence=[]string{"x"};if Verify(p){t.Fatal("tamper")}}