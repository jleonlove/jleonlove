package app
import("errors";"testing")
func TestWorkspace(t *testing.T){w:=New();if e:=w.Create(Transaction{ID:"tx1",Tenant:"org",Commodity:"Copper"});e!=nil{t.Fatal(e)};if e:=w.AddDocument("org",Document{ID:"d1",TransactionID:"tx1",Name:"BL",Digest:"sha",Evidence:"ev"});e!=nil{t.Fatal(e)};if e:=w.AddFinding("org",Finding{ID:"f1",TransactionID:"tx1",Kind:"DOC",Severity:9,Evidence:"ev"});e!=nil{t.Fatal(e)};_,d,f,e:=w.Snapshot("org","tx1");if e!=nil||len(d)!=1||len(f)!=1{t.Fatal(e)}}
func TestTenantIsolation(t *testing.T){w:=New();_ = w.Create(Transaction{ID:"same",Tenant:"a"});if _,_,_,e:=w.Snapshot("b","same");!errors.Is(e,ErrTransaction){t.Fatal(e)}}
func TestDocumentEvidence(t *testing.T){w:=New();_ = w.Create(Transaction{ID:"x",Tenant:"a"});if e:=w.AddDocument("a",Document{TransactionID:"x"});!errors.Is(e,ErrEvidence){t.Fatal(e)}}
func TestFindingPriority(t *testing.T){w:=New();_ = w.Create(Transaction{ID:"x",Tenant:"a"});_ = w.AddFinding("a",Finding{ID:"low",TransactionID:"x",Severity:1,Evidence:"e"});_ = w.AddFinding("a",Finding{ID:"high",TransactionID:"x",Severity:10,Evidence:"e"});_,_,f,_:=w.Snapshot("a","x");if f[0].ID!="high"{t.Fatal(f)}}
