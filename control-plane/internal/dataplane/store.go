package dataplane
import("crypto/sha256";"encoding/hex";"errors";"sort";"sync";"time")
var(ErrTenant=errors.New("tenant required");ErrKey=errors.New("key required");ErrConflict=errors.New("version conflict");ErrNotFound=errors.New("record not found");ErrEvidence=errors.New("evidence required"))
type Record struct{Tenant,Key string;Version uint64;Value []byte;Evidence string;UpdatedAt time.Time;Digest string}
type Store struct{mu sync.RWMutex;records map[string]Record}
func New()*Store{return &Store{records:map[string]Record{}}}
func id(t,k string)string{return t+"\x00"+k}
func digest(t,k string,v uint64,b []byte,e string)string{h:=sha256.New();h.Write([]byte(t));h.Write([]byte{0});h.Write([]byte(k));h.Write([]byte{0});h.Write([]byte(e));h.Write(b);return hex.EncodeToString(h.Sum(nil))}
func(s *Store)Put(r Record,expected uint64)(Record,error){
 if r.Tenant==""{return Record{},ErrTenant};if r.Key==""{return Record{},ErrKey};if r.Evidence==""{return Record{},ErrEvidence}
 s.mu.Lock();defer s.mu.Unlock();old,ok:=s.records[id(r.Tenant,r.Key)];if ok&&old.Version!=expected{return Record{},ErrConflict};if !ok&&expected!=0{return Record{},ErrConflict}
 r.Version=expected+1;r.UpdatedAt=time.Now().UTC();r.Value=append([]byte(nil),r.Value...);r.Digest=digest(r.Tenant,r.Key,r.Version,r.Value,r.Evidence);s.records[id(r.Tenant,r.Key)]=r;return r,nil
}
func(s *Store)Get(t,k string)(Record,error){s.mu.RLock();defer s.mu.RUnlock();r,ok:=s.records[id(t,k)];if !ok{return Record{},ErrNotFound};r.Value=append([]byte(nil),r.Value...);return r,nil}
func(s *Store)Keys(t string)[]string{s.mu.RLock();defer s.mu.RUnlock();o:=[]string{};for _,r:=range s.records{if r.Tenant==t{o=append(o,r.Key)}};sort.Strings(o);return o}
