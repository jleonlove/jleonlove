package tradecompliance
import("errors";"testing";"time")
func tm(n int64)time.Time{return time.Unix(n,0)}
func TestTemporalBlock(t *testing.T){r:=Rule{"sanction-1","US","OIL","X","",tm(100),tm(200),Block,false};rs,e:=Compile([]Rule{r});if e!=nil{t.Fatal(e)};_,e=Evaluate(rs,Trade{"OIL","X","Y",tm(150),nil});if !errors.Is(e,ErrBlocked){t.Fatal(e)}}
func TestExpiredRuleIgnored(t *testing.T){rs,_:=Compile([]Rule{{"r","US","OIL","X","",tm(100),tm(200),Block,false}});f,e:=Evaluate(rs,Trade{"OIL","X","Y",tm(250),nil});if e!=nil||len(f)!=0{t.Fatalf("%v %v",f,e)}}
func TestDestinationMatch(t *testing.T){rs,_:=Compile([]Rule{{"tariff","EU","STEEL","","DE",tm(1),time.Time{},Review,false}});f,e:=Evaluate(rs,Trade{"STEEL","BR","DE",tm(2),nil});if e!=nil||len(f)!=1||f[0].Action!=Review{t.Fatal(f,e)}}
func TestEvidenceMissingReview(t *testing.T){rs,_:=Compile([]Rule{{"origin-proof","EU","COFFEE","BR","DE",tm(1),time.Time{},Allow,true}});f,e:=Evaluate(rs,Trade{"COFFEE","BR","DE",tm(2),map[string]bool{}});if e!=nil||len(f)!=1||f[0].Action!=Review{t.Fatal(f,e)}}
func TestEvidenceSatisfied(t *testing.T){rs,_:=Compile([]Rule{{"origin-proof","EU","COFFEE","BR","DE",tm(1),time.Time{},Allow,true}});f,e:=Evaluate(rs,Trade{"COFFEE","BR","DE",tm(2),map[string]bool{"origin-proof":true}});if e!=nil||f[0].Action!=Allow{t.Fatal(f,e)}}
func TestInvalidRule(t *testing.T){if _,e:=Compile([]Rule{{ID:"x"}});!errors.Is(e,ErrRule){t.Fatal(e)}}
