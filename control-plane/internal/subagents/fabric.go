package subagents
import("crypto/sha256";"encoding/hex";"errors";"sort";"strings";"time")
var(ErrAuthorityEscalation=errors.New("delegated authority exceeds parent");ErrBudget=errors.New("subagent budget invalid");ErrWorkspace=errors.New("workspace isolation required");ErrSelfMerge=errors.New("subagent cannot self-merge"))
type Envelope struct{AgentID,ParentID,Workspace string;Capabilities []string;MaxTokens int;Deadline time.Time}
type Proposal struct{AgentID,Workspace,DiffDigest string;SelfApproved bool}
func subset(child,parent []string)bool{m:=map[string]bool{};for _,x:=range parent{m[x]=true};for _,x:=range child{if !m[x]{return false}};return true}
func Delegate(parent Envelope,child Envelope)error{
 if !subset(child.Capabilities,parent.Capabilities){return ErrAuthorityEscalation}
 if child.MaxTokens<=0||child.MaxTokens>parent.MaxTokens||child.Deadline.After(parent.Deadline){return ErrBudget}
 if child.Workspace==""||child.Workspace==parent.Workspace{return ErrWorkspace}
 return nil
}
func Integrate(p Proposal)error{if p.SelfApproved{return ErrSelfMerge};if p.DiffDigest==""||p.Workspace==""{return ErrWorkspace};return nil}
func EvidenceDigest(e Envelope)string{c:=append([]string(nil),e.Capabilities...);sort.Strings(c);s:=sha256.Sum256([]byte(strings.Join(append([]string{e.AgentID,e.ParentID,e.Workspace},c...),"\x00")));return hex.EncodeToString(s[:])}
