package tradeontology
import("errors";"sort";"strings")
var(ErrKind=errors.New("unknown trade concept kind");ErrCode=errors.New("invalid trade code");ErrUnit=errors.New("unit incompatible");ErrDuplicate=errors.New("duplicate canonical concept"))
type Kind string
const(Commodity Kind="COMMODITY";Location Kind="LOCATION";Party Kind="PARTY";Document Kind="DOCUMENT";Event Kind="EVENT";Standard Kind="STANDARD")
type Concept struct{ID string;Kind Kind;Standard,Code,Name,Unit string;Aliases []string}
type Registry struct{ByID map[string]Concept;Alias map[string]string}
func New() *Registry{return &Registry{ByID:map[string]Concept{},Alias:map[string]string{}}}
func validKind(k Kind)bool{switch k{case Commodity,Location,Party,Document,Event,Standard:return true};return false}
func(r *Registry)Add(c Concept)error{
 if !validKind(c.Kind){return ErrKind};if strings.TrimSpace(c.ID)==""||strings.TrimSpace(c.Code)==""||strings.TrimSpace(c.Standard)==""{return ErrCode}
 if _,ok:=r.ByID[c.ID];ok{return ErrDuplicate};r.ByID[c.ID]=c
 for _,a:=range append(c.Aliases,c.Name,c.Code){x:=strings.ToLower(strings.TrimSpace(a));if x!=""{r.Alias[x]=c.ID}}
 return nil
}
func(r *Registry)Resolve(x string)(Concept,bool){id,ok:=r.Alias[strings.ToLower(strings.TrimSpace(x))];if !ok{id=x};c,ok:=r.ByID[id];return c,ok}
func CompatibleUnit(c Concept,u string)error{if c.Unit!=""&&u!=c.Unit{return ErrUnit};return nil}
func(r *Registry)IDs()[]string{x:=make([]string,0,len(r.ByID));for id:=range r.ByID{x=append(x,id)};sort.Strings(x);return x}
