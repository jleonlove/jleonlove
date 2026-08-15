package controltower
import"testing"
func TestPriority(t *testing.T){x:=New();x.Upsert(Alert{ID:"low",TransactionID:"tx",Severity:1});x.Upsert(Alert{ID:"high",TransactionID:"tx",Severity:9});if x.Open("tx")[0].ID!="high"{t.Fatal()}}