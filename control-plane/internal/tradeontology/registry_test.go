package tradeontology
import("errors";"testing")
func TestCanonicalResolve(t *testing.T){r:=New();c:=Concept{ID:"cmd:cu",Kind:Commodity,Standard:"ATLAS",Code:"CU",Name:"Copper",Unit:"MT",Aliases:[]string{"copper cathode"}};if e:=r.Add(c);e!=nil{t.Fatal(e)};x,ok:=r.Resolve("Copper Cathode");if !ok||x.ID!="cmd:cu"{t.Fatal("resolve")}}
func TestUnknownKind(t *testing.T){r:=New();if e:=r.Add(Concept{ID:"x",Kind:"BAD",Standard:"x",Code:"x"});!errors.Is(e,ErrKind){t.Fatal(e)}}
func TestMissingStandardCode(t *testing.T){r:=New();if e:=r.Add(Concept{ID:"x",Kind:Commodity});!errors.Is(e,ErrCode){t.Fatal(e)}}
func TestDuplicate(t *testing.T){r:=New();c:=Concept{ID:"x",Kind:Commodity,Standard:"A",Code:"1"};_ = r.Add(c);if e:=r.Add(c);!errors.Is(e,ErrDuplicate){t.Fatal(e)}}
func TestUnitMismatch(t *testing.T){if e:=CompatibleUnit(Concept{Unit:"MT"},"BBL");!errors.Is(e,ErrUnit){t.Fatal(e)}}
func TestLocationAndDocument(t *testing.T){r:=New();for _,c:=range []Concept{{ID:"loc:rot",Kind:Location,Standard:"UNLOCODE",Code:"NLRTM",Name:"Rotterdam"},{ID:"doc:bl",Kind:Document,Standard:"DCSA",Code:"BL",Name:"Bill of Lading"}}{if e:=r.Add(c);e!=nil{t.Fatal(e)}};if len(r.IDs())!=2{t.Fatal("registry")}}
