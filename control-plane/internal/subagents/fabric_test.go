package subagents
import("errors";"testing";"time")
func parent()Envelope{return Envelope{AgentID:"parent",Workspace:"canonical",Capabilities:[]string{"repo.read","repo.write"},MaxTokens:1000,Deadline:time.Now().Add(time.Hour)}}
func child()Envelope{return Envelope{AgentID:"child",ParentID:"parent",Workspace:"worktree-child",Capabilities:[]string{"repo.read"},MaxTokens:500,Deadline:time.Now().Add(time.Minute)}}
func TestBoundedDelegation(t *testing.T){if e:=Delegate(parent(),child());e!=nil{t.Fatal(e)}}
func TestAuthorityEscalationDenied(t *testing.T){c:=child();c.Capabilities=append(c.Capabilities,"prod.deploy");if e:=Delegate(parent(),c);!errors.Is(e,ErrAuthorityEscalation){t.Fatal(e)}}
func TestBudgetEscalationDenied(t *testing.T){c:=child();c.MaxTokens=2000;if e:=Delegate(parent(),c);!errors.Is(e,ErrBudget){t.Fatal(e)}}
func TestWorkspaceIsolationRequired(t *testing.T){c:=child();c.Workspace="canonical";if e:=Delegate(parent(),c);!errors.Is(e,ErrWorkspace){t.Fatal(e)}}
func TestSelfMergeDenied(t *testing.T){if e:=Integrate(Proposal{AgentID:"child",Workspace:"w",DiffDigest:"abc",SelfApproved:true});!errors.Is(e,ErrSelfMerge){t.Fatal(e)}}
func TestEvidenceChangesWithAuthority(t *testing.T){a:=child();b:=child();b.Capabilities=[]string{"repo.read","repo.write"};if EvidenceDigest(a)==EvidenceDigest(b){t.Fatal("digest unchanged")}}
