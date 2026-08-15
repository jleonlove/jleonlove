package tradeintegrity
import("errors";"testing";"time")
func graph()*Graph{g:=New();g.AddEntity(Entity{"seller","ORG"});g.AddEntity(Entity{"agent","PERSON"});g.AddEntity(Entity{"operator","ORG"});return g}
func TestValidMandate(t *testing.T){g:=graph();m:=Mandate{"seller","agent","SELL",time.Unix(2000,0),"sha256:m"};if e:=g.AddMandate(m,time.Unix(1000,0));e!=nil||!g.Authorized("seller","agent","SELL",time.Unix(1500,0)){t.Fatal(e)}}
func TestExpiredMandate(t *testing.T){g:=graph();m:=Mandate{"seller","agent","SELL",time.Unix(1000,0),"sha256:m"};if e:=g.AddMandate(m,time.Unix(1000,0));!errors.Is(e,ErrExpired){t.Fatal(e)}}
func TestMandateEvidenceRequired(t *testing.T){g:=graph();m:=Mandate{"seller","agent","SELL",time.Unix(2000,0),""};if e:=g.AddMandate(m,time.Unix(1000,0));!errors.Is(e,ErrAuthority){t.Fatal(e)}}
func TestWrongScopeNotAuthorized(t *testing.T){g:=graph();_ = g.AddMandate(Mandate{"seller","agent","SELL",time.Unix(2000,0),"ev"},time.Unix(1000,0));if g.Authorized("seller","agent","BORROW",time.Unix(1500,0)){t.Fatal("scope escalation")}}
func TestFacilityCapability(t *testing.T){g:=graph();f:=Facility{"refinery","operator",map[string]bool{"REFINE_COPPER":true},"sha256:f"};if e:=g.AddFacility(f);e!=nil{t.Fatal(e)};if e:=g.Supports("refinery","REFINE_COPPER");e!=nil{t.Fatal(e)}}
func TestUnsupportedFacilityClaim(t *testing.T){g:=graph();_ = g.AddFacility(Facility{"warehouse","operator",map[string]bool{"STORE":true},"ev"});if e:=g.Supports("warehouse","REFINE");!errors.Is(e,ErrFacility){t.Fatal(e)}}
func TestUnknownOperator(t *testing.T){g:=graph();if e:=g.AddFacility(Facility{"x","ghost",map[string]bool{"STORE":true},"ev"});!errors.Is(e,ErrEntity){t.Fatal(e)}}
