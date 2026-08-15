package autonomy
import("errors";"sort";"time")
var(ErrAuthority=errors.New("autonomous authority denied");ErrExpired=errors.New("autonomy grant expired");ErrBudget=errors.New("autonomy budget exceeded");ErrDuplicate=errors.New("duplicate event");ErrApproval=errors.New("human approval required"))
type Grant struct{Capabilities map[string]bool;ExpiresAt time.Time;MaxActions int}
type Event struct{ID,Capability string;At time.Time;Consequential bool;Approved bool}
type Engine struct{Grant Grant;Used int;Seen map[string]bool}
func(e *Engine)Handle(x Event)error{
 if e.Seen==nil{e.Seen=map[string]bool{}}
 if e.Seen[x.ID]{return ErrDuplicate}
 if !x.At.Before(e.Grant.ExpiresAt){return ErrExpired}
 if !e.Grant.Capabilities[x.Capability]{return ErrAuthority}
 if e.Used>=e.Grant.MaxActions{return ErrBudget}
 if x.Consequential&&!x.Approved{return ErrApproval}
 e.Seen[x.ID]=true;e.Used++;return nil
}
func Pending(events []Event)[]string{var x []string;for _,e:=range events{x=append(x,e.ID)};sort.Strings(x);return x}
