package credentialbroker
import("errors";"strings";"time")
var(ErrIdentity=errors.New("identity mismatch");ErrScope=errors.New("credential scope denied");ErrAudience=errors.New("credential audience denied");ErrExpired=errors.New("credential expired");ErrSensitive=errors.New("sensitive data policy denied"))
type Grant struct{UserID,AgentID,Audience string;Scopes []string;ExpiresAt time.Time}
type Request struct{UserID,AgentID,Audience,Scope string;Now time.Time}
type Data struct{Classification,Residency string}
type Policy struct{AllowedAudiences map[string]bool;AllowedScopes map[string]bool;AllowedResidencies map[string]bool}
func Issue(p Policy,g Grant,r Request)(string,error){
 if g.UserID!=r.UserID||g.AgentID!=r.AgentID{return "",ErrIdentity}
 if !r.Now.Before(g.ExpiresAt){return "",ErrExpired}
 if g.Audience!=r.Audience||!p.AllowedAudiences[r.Audience]{return "",ErrAudience}
 ok:=false;for _,s:=range g.Scopes{if s==r.Scope{ok=true}};if !ok||!p.AllowedScopes[r.Scope]{return "",ErrScope}
 return "handle:"+g.AgentID+":"+r.Scope,nil
}
func AuthorizeData(p Policy,d Data)error{if d.Classification=="restricted"&&!p.AllowedResidencies[d.Residency]{return ErrSensitive};return nil}
func RawSecret(v string)bool{x:=strings.ToLower(v);return strings.Contains(x,"bearer ")||strings.Contains(x,"api_key=")||strings.Contains(x,"password=")}
