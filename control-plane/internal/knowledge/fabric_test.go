package knowledge
import("errors";"testing";"time")
func now()time.Time{return time.Unix(1000,0)}
func fact(obj,src string,rel float64)Fact{return Fact{ID:src,Subject:"atlas",Predicate:"status",Object:obj,Source:src,Tenant:"t1",Classification:"internal",ObservedAt:now(),Reliability:rel,Readers:map[string]bool{"agent":true}}}
func TestProvenanceRequired(t *testing.T){f:=fact("ready","",.9);if Ingest(f)==nil{t.Fatal("accepted missing provenance")}}
func TestPermissionBoundary(t *testing.T){_,e:=Retrieve([]Fact{fact("ready","a",.9)},Query{Tenant:"t1",Reader:"other",Subject:"atlas",Predicate:"status",Now:now()});if !errors.Is(e,ErrPermission){t.Fatal(e)}}
func TestTenantBoundary(t *testing.T){_,e:=Retrieve([]Fact{fact("ready","a",.9)},Query{Tenant:"t2",Reader:"agent",Subject:"atlas",Predicate:"status",Now:now()});if !errors.Is(e,ErrPermission){t.Fatal(e)}}
func TestExpiredExcluded(t *testing.T){f:=fact("ready","a",.9);f.ValidUntil=now();_,e:=Retrieve([]Fact{f},Query{Tenant:"t1",Reader:"agent",Subject:"atlas",Predicate:"status",Now:now()});if !errors.Is(e,ErrPermission){t.Fatal(e)}}
func TestConflictSurfaced(t *testing.T){r,e:=Retrieve([]Fact{fact("ready","a",.9),fact("not-ready","b",.8)},Query{Tenant:"t1",Reader:"agent",Subject:"atlas",Predicate:"status",Now:now()});if e!=nil||!r.Conflicted{t.Fatalf("%+v %v",r,e)}}
func TestReliabilityRanks(t *testing.T){r,e:=Retrieve([]Fact{fact("low","a",.2),fact("high","b",.9)},Query{Tenant:"t1",Reader:"agent",Subject:"atlas",Predicate:"status",Now:now()});if e!=nil||r.Facts[0].Object!="high"{t.Fatal("ranking failed")}}
