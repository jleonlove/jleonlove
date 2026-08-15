package app
import("errors";"sort";"sync";"time")
var(ErrTenant=errors.New("tenant required");ErrTransaction=errors.New("transaction missing");ErrEvidence=errors.New("evidence required"))
type Transaction struct{ID,Tenant,Commodity,Origin,Destination,Status string;UpdatedAt time.Time}
type Document struct{ID,TransactionID,Name,Digest,Evidence string}
type Finding struct{ID,TransactionID,Kind,Message string;Severity int;Evidence string}
type Workspace struct{mu sync.RWMutex;tx map[string]Transaction;docs map[string][]Document;findings map[string][]Finding}
func New()*Workspace{return &Workspace{tx:map[string]Transaction{},docs:map[string][]Document{},findings:map[string][]Finding{}}}
func key(t,id string)string{return t+"\x00"+id}
func(w *Workspace)Create(x Transaction)error{if x.Tenant==""{return ErrTenant};if x.ID==""{return ErrTransaction};w.mu.Lock();defer w.mu.Unlock();x.UpdatedAt=time.Now().UTC();w.tx[key(x.Tenant,x.ID)]=x;return nil}
func(w *Workspace)AddDocument(t string,d Document)error{w.mu.Lock();defer w.mu.Unlock();if _,ok:=w.tx[key(t,d.TransactionID)];!ok{return ErrTransaction};if d.Evidence==""||d.Digest==""{return ErrEvidence};w.docs[key(t,d.TransactionID)]=append(w.docs[key(t,d.TransactionID)],d);return nil}
func(w *Workspace)AddFinding(t string,f Finding)error{w.mu.Lock();defer w.mu.Unlock();if _,ok:=w.tx[key(t,f.TransactionID)];!ok{return ErrTransaction};if f.Evidence==""{return ErrEvidence};w.findings[key(t,f.TransactionID)]=append(w.findings[key(t,f.TransactionID)],f);return nil}
func(w *Workspace)Snapshot(t,id string)(Transaction,[]Document,[]Finding,error){w.mu.RLock();defer w.mu.RUnlock();x,ok:=w.tx[key(t,id)];if !ok{return Transaction{},nil,nil,ErrTransaction};d:=append([]Document(nil),w.docs[key(t,id)]...);f:=append([]Finding(nil),w.findings[key(t,id)]...);sort.Slice(f,func(i,j int)bool{if f[i].Severity==f[j].Severity{return f[i].ID<f[j].ID};return f[i].Severity>f[j].Severity});return x,d,f,nil}
