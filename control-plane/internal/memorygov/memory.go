package memorygov
import("crypto/sha256";"encoding/hex";"errors";"time")
var(ErrTenant=errors.New("memory tenant mismatch");ErrRevoked=errors.New("memory revoked");ErrExpired=errors.New("memory expired");ErrIntegrity=errors.New("memory integrity failure");ErrAuthority=errors.New("memory cannot grant authority"))
type Record struct{ID,Tenant,Source,Classification,Content,Digest string;CreatedAt,ExpiresAt time.Time;Revoked bool;AuthorityClaims []string}
func Digest(content string)string{s:=sha256.Sum256([]byte(content));return hex.EncodeToString(s[:])}
func New(id,tenant,source,class,content string,now,expires time.Time)Record{return Record{ID:id,Tenant:tenant,Source:source,Classification:class,Content:content,Digest:Digest(content),CreatedAt:now,ExpiresAt:expires}}
func Read(r Record,tenant string,now time.Time)(string,error){if r.Tenant!=tenant{return "",ErrTenant};if r.Revoked{return "",ErrRevoked};if !r.ExpiresAt.IsZero()&&!now.Before(r.ExpiresAt){return "",ErrExpired};if Digest(r.Content)!=r.Digest{return "",ErrIntegrity};if len(r.AuthorityClaims)>0{return "",ErrAuthority};return r.Content,nil}
func Revoke(r Record)Record{r.Revoked=true;return r}
