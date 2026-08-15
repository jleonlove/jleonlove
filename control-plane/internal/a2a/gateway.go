package a2a
import("errors";"sort";"strings")
var(ErrUnqualified=errors.New("remote agent unqualified");ErrCapability=errors.New("remote capability denied");ErrTrust=errors.New("remote trust insufficient");ErrAuthority=errors.New("delegation exceeds authority");ErrEndpoint=errors.New("endpoint not allowed"))
type Card struct{AgentID,Endpoint string;Capabilities []string;Trust int;Qualified bool}
type Delegation struct{AgentID string;Capabilities []string;MaxTrust int}
type Gateway struct{AllowedEndpoints map[string]bool;MinTrust int}
func subset(a,b []string)bool{m:=map[string]bool{};for _,x:=range b{m[x]=true};for _,x:=range a{if !m[x]{return false}};return true}
func(g Gateway)Qualify(c Card,allowedCaps []string)error{if !c.Qualified{return ErrUnqualified};if !g.AllowedEndpoints[c.Endpoint]{return ErrEndpoint};if c.Trust<g.MinTrust{return ErrTrust};if !subset(c.Capabilities,allowedCaps){return ErrCapability};return nil}
func(g Gateway)Delegate(c Card,parentCaps,requested []string) (Delegation,error){if e:=g.Qualify(c,parentCaps);e!=nil{return Delegation{},e};if !subset(requested,parentCaps)||!subset(requested,c.Capabilities){return Delegation{},ErrAuthority};return Delegation{AgentID:c.AgentID,Capabilities:requested,MaxTrust:c.Trust},nil}
func Fingerprint(c Card)string{x:=append([]string(nil),c.Capabilities...);sort.Strings(x);return strings.Join(append([]string{c.AgentID,c.Endpoint},x...),"\x00")}
