package engineering
import("errors";"sort";"strings")
var(ErrSensitiveReview=errors.New("security-sensitive change requires independent review");ErrSelfApproval=errors.New("author cannot approve own change");ErrGate=errors.New("release gate incomplete");ErrCanonicalWrite=errors.New("agent cannot write canonical workspace directly"))
type Risk string
const(Normal Risk="NORMAL";Sensitive Risk="SECURITY_SENSITIVE")
type Change struct{ID,AuthorID,Workspace string;Files []string;TestsPassed,SecurityReviewPassed,CapabilityDiffPassed bool}
func Classify(files []string)Risk{for _,f:=range files{p:=strings.ToLower(f);for _,x:=range []string{"auth","iam","crypto","policy","secret","containment","release","deploy"}{if strings.Contains(p,x){return Sensitive}}};return Normal}
func ValidateWorkspace(c Change)error{if c.Workspace==""||c.Workspace=="canonical"{return ErrCanonicalWrite};return nil}
func Approve(c Change,reviewerID string)error{if reviewerID==c.AuthorID{return ErrSelfApproval};if Classify(c.Files)==Sensitive&&!c.SecurityReviewPassed{return ErrSensitiveReview};if !c.TestsPassed||!c.CapabilityDiffPassed{return ErrGate};return nil}
func DeterministicFiles(files []string)string{c:=append([]string(nil),files...);sort.Strings(c);return strings.Join(c,"\x00")}
