package workflow
import("errors";"sort";"strings")
var(ErrUnknownStep=errors.New("unknown workflow step");ErrDependency=errors.New("dependency incomplete");ErrLease=errors.New("invalid lease");ErrAuthorityRefresh=errors.New("authority refresh required");ErrTerminal=errors.New("workflow terminal"))
type State string
const(Pending State="PENDING";Running State="RUNNING";Paused State="PAUSED";Succeeded State="SUCCEEDED";Failed State="FAILED";Cancelled State="CANCELLED")
type Step struct{ID string;DependsOn []string;State State;LeaseOwner string;Checkpoint string}
type Mission struct{ID string;Steps map[string]Step;AuthorityEpoch uint64;Cancelled bool}
func Ready(m Mission)[]string{var out []string;if m.Cancelled{return out};for id,s:=range m.Steps{if s.State!=Pending{continue};ok:=true;for _,d:=range s.DependsOn{x,exists:=m.Steps[d];if !exists||x.State!=Succeeded{ok=false;break}};if ok{out=append(out,id)}};sort.Strings(out);return out}
func Start(m *Mission,id,owner string,authorityEpoch uint64)error{if m.Cancelled{return ErrTerminal};s,ok:=m.Steps[id];if !ok{return ErrUnknownStep};if authorityEpoch!=m.AuthorityEpoch{return ErrAuthorityRefresh};for _,d:=range s.DependsOn{if m.Steps[d].State!=Succeeded{return ErrDependency}};if owner==""{return ErrLease};s.State=Running;s.LeaseOwner=owner;m.Steps[id]=s;return nil}
func Checkpoint(m *Mission,id,owner,data string)error{s,ok:=m.Steps[id];if !ok{return ErrUnknownStep};if s.State!=Running||s.LeaseOwner!=owner{return ErrLease};s.Checkpoint=data;m.Steps[id]=s;return nil}
func Complete(m *Mission,id,owner string)error{s,ok:=m.Steps[id];if !ok{return ErrUnknownStep};if s.State!=Running||s.LeaseOwner!=owner{return ErrLease};s.State=Succeeded;s.LeaseOwner="";m.Steps[id]=s;return nil}
func Resume(m *Mission,id,owner string,authorityEpoch uint64)error{s,ok:=m.Steps[id];if !ok{return ErrUnknownStep};if authorityEpoch!=m.AuthorityEpoch{return ErrAuthorityRefresh};if s.Checkpoint==""{return errors.New("checkpoint required")};if s.State!=Paused&&s.State!=Running{return ErrTerminal};s.State=Running;s.LeaseOwner=owner;m.Steps[id]=s;return nil}
func Pause(m *Mission,id,owner string)error{s,ok:=m.Steps[id];if !ok{return ErrUnknownStep};if s.State!=Running||s.LeaseOwner!=owner{return ErrLease};s.State=Paused;s.LeaseOwner="";m.Steps[id]=s;return nil}
func Cancel(m *Mission){m.Cancelled=true;for id,s:=range m.Steps{if s.State==Pending||s.State==Running||s.State==Paused{s.State=Cancelled;s.LeaseOwner="";m.Steps[id]=s}}}
func Fingerprint(m Mission)string{ids:=make([]string,0,len(m.Steps));for id:=range m.Steps{ids=append(ids,id)};sort.Strings(ids);var b []string;for _,id:=range ids{s:=m.Steps[id];b=append(b,id,string(s.State),s.Checkpoint)};return strings.Join(b,"\x00")}
