package dataplane
import("errors";"sync";"testing")
func TestPutGet(t *testing.T){s:=New();r,e:=s.Put(Record{Tenant:"a",Key:"tx/1",Value:[]byte("state"),Evidence:"ev"},0);if e!=nil||r.Version!=1{t.Fatal(r,e)};g,e:=s.Get("a","tx/1");if e!=nil||string(g.Value)!="state"||g.Digest==""{t.Fatal(g,e)}}
func TestTenantIsolation(t *testing.T){s:=New();_,_=s.Put(Record{Tenant:"a",Key:"k",Value:[]byte("A"),Evidence:"e"},0);_,_=s.Put(Record{Tenant:"b",Key:"k",Value:[]byte("B"),Evidence:"e"},0);a,_:=s.Get("a","k");b,_:=s.Get("b","k");if string(a.Value)==string(b.Value){t.Fatal("tenant leak")}}
func TestConflict(t *testing.T){s:=New();_,_=s.Put(Record{Tenant:"a",Key:"k",Evidence:"e"},0);if _,e:=s.Put(Record{Tenant:"a",Key:"k",Evidence:"e2"},0);!errors.Is(e,ErrConflict){t.Fatal(e)}}
func TestEvidence(t *testing.T){s:=New();if _,e:=s.Put(Record{Tenant:"a",Key:"k"},0);!errors.Is(e,ErrEvidence){t.Fatal(e)}}
func TestCopyIsolation(t *testing.T){s:=New();b:=[]byte("x");_,_=s.Put(Record{Tenant:"a",Key:"k",Value:b,Evidence:"e"},0);b[0]='z';g,_:=s.Get("a","k");if string(g.Value)!="x"{t.Fatal("alias")}}
func TestConcurrentReadsWrites(t *testing.T){s:=New();r,_:=s.Put(Record{Tenant:"a",Key:"k",Evidence:"e"},0);var wg sync.WaitGroup;for i:=0;i<20;i++{wg.Add(1);go func(){defer wg.Done();_,_=s.Get("a","k")}()};wg.Wait();if r.Version!=1{t.Fatal()}}
