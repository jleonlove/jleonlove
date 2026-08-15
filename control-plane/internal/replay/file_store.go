package replay
import("context";"encoding/json";"errors";"os";"path/filepath";"sync";"time")
var ErrStoreUnavailable=errors.New("replay store unavailable")
type FileStore struct{mu sync.Mutex;path string;now func()time.Time}
type diskEntry struct{OrganizationID string `json:"organization_id"`;RequestID string `json:"request_id"`;Nonce string `json:"nonce"`;ExpiresAt time.Time `json:"expires_at"`}
type diskState struct{Used []diskEntry `json:"used"`}
func NewFileStore(path string)*FileStore{return &FileStore{path:path,now:time.Now}}
func (s *FileStore)Consume(_ context.Context,key Key,expiresAt time.Time)error{
 s.mu.Lock();defer s.mu.Unlock();now:=s.now();if !now.Before(expiresAt){return ErrExpired}
 st,err:=s.load();if err!=nil{return errors.Join(ErrStoreUnavailable,err)}
 kept:=make([]diskEntry,0,len(st.Used)+1)
 for _,e:=range st.Used{if !now.Before(e.ExpiresAt){continue};if e.OrganizationID==key.OrganizationID&&e.RequestID==key.RequestID&&e.Nonce==key.Nonce{return ErrReplay};kept=append(kept,e)}
 kept=append(kept,diskEntry{key.OrganizationID,key.RequestID,key.Nonce,expiresAt});st.Used=kept
 if err:=s.persist(st);err!=nil{return errors.Join(ErrStoreUnavailable,err)};return nil
}
func(s *FileStore)load()(diskState,error){var st diskState;b,err:=os.ReadFile(s.path);if errors.Is(err,os.ErrNotExist){return st,nil};if err!=nil{return st,err};if err=json.Unmarshal(b,&st);err!=nil{return st,err};return st,nil}
func(s *FileStore)persist(st diskState)error{
 if err:=os.MkdirAll(filepath.Dir(s.path),0700);err!=nil{return err};b,err:=json.Marshal(st);if err!=nil{return err};tmp:=s.path+".tmp"
 f,err:=os.OpenFile(tmp,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,0600);if err!=nil{return err}
 if _,err=f.Write(b);err!=nil{f.Close();return err};if err=f.Sync();err!=nil{f.Close();return err};if err=f.Close();err!=nil{return err}
 if err=os.Rename(tmp,s.path);err!=nil{return err};d,err:=os.Open(filepath.Dir(s.path));if err!=nil{return err};defer d.Close();return d.Sync()
}
