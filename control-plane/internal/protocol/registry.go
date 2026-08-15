package protocol
import("errors";"sort";"strings")
var(ErrProtocol=errors.New("protocol not allowed");ErrVersion=errors.New("protocol version incompatible");ErrCapability=errors.New("plugin capability denied");ErrUnsigned=errors.New("unsigned plugin");ErrSchema=errors.New("invalid plugin schema"))
type Plugin struct{Name,Protocol,Version,Signature string;Capabilities []string;Schema string}
type Policy struct{Protocols map[string][]string;Capabilities map[string]bool;RequireSignature bool}
func Admit(p Policy,x Plugin)error{
 versions,ok:=p.Protocols[x.Protocol];if !ok{return ErrProtocol}
 vok:=false;for _,v:=range versions{if v==x.Version{vok=true}};if !vok{return ErrVersion}
 if p.RequireSignature&&strings.TrimSpace(x.Signature)==""{return ErrUnsigned}
 if strings.TrimSpace(x.Schema)==""{return ErrSchema}
 for _,c:=range x.Capabilities{if !p.Capabilities[c]{return ErrCapability}}
 return nil
}
func Fingerprint(x Plugin)string{c:=append([]string(nil),x.Capabilities...);sort.Strings(c);return strings.Join(append([]string{x.Name,x.Protocol,x.Version,x.Schema},c...),"\x00")}
