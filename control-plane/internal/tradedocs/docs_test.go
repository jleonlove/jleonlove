package tradedocs
import("errors";"testing")
func spa()Document{return Document{ID:"spa",Type:"SPA",Source:"upload",Hash:"sha256:a",Fields:map[string]string{"commodity":"Copper","quantity":"25000","unit":"MT","buyer":"B","seller":"S","destination":"Rotterdam"}}}
func bl()Document{return Document{ID:"bl",Type:"BL",Source:"carrier",Hash:"sha256:b",Fields:map[string]string{"commodity":"Copper","quantity":"25000","unit":"MT","origin":"Chile","destination":"Rotterdam"}}}
func TestValid(t *testing.T){if e:=Validate(spa());e!=nil{t.Fatal(e)}}
func TestUnsupported(t *testing.T){x:=spa();x.Type="FAKE";if e:=Validate(x);!errors.Is(e,ErrType){t.Fatal(e)}}
func TestProvenanceRequired(t *testing.T){x:=spa();x.Hash="";if e:=Validate(x);!errors.Is(e,ErrProvenance){t.Fatal(e)}}
func TestRequiredField(t *testing.T){x:=spa();delete(x.Fields,"buyer");if e:=Validate(x);!errors.Is(e,ErrField){t.Fatal(e)}}
func TestNoConflict(t *testing.T){c,e:=Reconcile([]Document{spa(),bl()},[]string{"commodity","quantity","unit","destination"});if e!=nil||len(c)!=0{t.Fatalf("%v %v",c,e)}}
func TestQuantityConflict(t *testing.T){b:=bl();b.Fields["quantity"]="24760";c,e:=Reconcile([]Document{spa(),b},[]string{"quantity"});if e!=nil||len(c)!=1||c[0].Field!="quantity"{t.Fatalf("%v %v",c,e)}}
func TestDestinationConflict(t *testing.T){b:=bl();b.Fields["destination"]="Antwerp";c,e:=Reconcile([]Document{spa(),b},[]string{"destination"});if e!=nil||len(c)!=1{t.Fatalf("%v %v",c,e)}}
