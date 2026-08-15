package protocol
import("errors";"testing")
func pol()Policy{return Policy{Protocols:map[string][]string{"mcp":{"1"},"a2a":{"1"}},Capabilities:map[string]bool{"read":true,"search":true},RequireSignature:true}}
func plug()Plugin{return Plugin{Name:"docs",Protocol:"mcp",Version:"1",Signature:"sig",Capabilities:[]string{"read"},Schema:"v1"}}
func TestAdmit(t *testing.T){if e:=Admit(pol(),plug());e!=nil{t.Fatal(e)}}
func TestUnknownProtocol(t *testing.T){x:=plug();x.Protocol="evil";if e:=Admit(pol(),x);!errors.Is(e,ErrProtocol){t.Fatal(e)}}
func TestBadVersion(t *testing.T){x:=plug();x.Version="99";if e:=Admit(pol(),x);!errors.Is(e,ErrVersion){t.Fatal(e)}}
func TestUnsigned(t *testing.T){x:=plug();x.Signature="";if e:=Admit(pol(),x);!errors.Is(e,ErrUnsigned){t.Fatal(e)}}
func TestCapabilityExpansion(t *testing.T){x:=plug();x.Capabilities=[]string{"write"};if e:=Admit(pol(),x);!errors.Is(e,ErrCapability){t.Fatal(e)}}
func TestSchemaRequired(t *testing.T){x:=plug();x.Schema="";if e:=Admit(pol(),x);!errors.Is(e,ErrSchema){t.Fatal(e)}}
func TestFingerprintStable(t *testing.T){a:=plug();a.Capabilities=[]string{"search","read"};b:=plug();b.Capabilities=[]string{"read","search"};if Fingerprint(a)!=Fingerprint(b){t.Fatal("unstable")}}
