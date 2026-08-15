package tradeintegrityengine
import"testing"
func TestPass(t *testing.T){r:=Evaluate(Signals{true,true,true,true,true});if !r.Pass||r.Score!=100{t.Fatal(r)}}
func TestFailClosed(t *testing.T){r:=Evaluate(Signals{true,true,false,true,true});if r.Pass||r.Score!=80{t.Fatal(r)}}