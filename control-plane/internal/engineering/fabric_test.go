package engineering
import("errors";"testing")
func good()Change{return Change{ID:"c1",AuthorID:"builder",Workspace:"worktree-c1",Files:[]string{"internal/foo.go"},TestsPassed:true,SecurityReviewPassed:true,CapabilityDiffPassed:true}}
func TestNormalIndependentApproval(t *testing.T){if e:=Approve(good(),"reviewer");e!=nil{t.Fatal(e)}}
func TestSelfApprovalDenied(t *testing.T){if e:=Approve(good(),"builder");!errors.Is(e,ErrSelfApproval){t.Fatal(e)}}
func TestSensitiveRequiresSecurityReview(t *testing.T){c:=good();c.Files=[]string{"internal/auth/jwt.go"};c.SecurityReviewPassed=false;if e:=Approve(c,"reviewer");!errors.Is(e,ErrSensitiveReview){t.Fatal(e)}}
func TestFailedTestsBlock(t *testing.T){c:=good();c.TestsPassed=false;if e:=Approve(c,"reviewer");!errors.Is(e,ErrGate){t.Fatal(e)}}
func TestCapabilityDiffBlock(t *testing.T){c:=good();c.CapabilityDiffPassed=false;if e:=Approve(c,"reviewer");!errors.Is(e,ErrGate){t.Fatal(e)}}
func TestCanonicalWriteDenied(t *testing.T){c:=good();c.Workspace="canonical";if e:=ValidateWorkspace(c);!errors.Is(e,ErrCanonicalWrite){t.Fatal(e)}}
func TestSensitiveClassification(t *testing.T){if Classify([]string{"deploy/policy.yaml"})!=Sensitive{t.Fatal("not sensitive")}}
