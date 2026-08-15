package tradeintegrity
import("errors";"sort";"time")
var(ErrEntity=errors.New("integrity entity missing");ErrAuthority=errors.New("representative authority invalid");ErrExpired=errors.New("mandate expired");ErrFacility=errors.New("facility capability unsupported");ErrConflict=errors.New("conflicting integrity claim"))
type Entity struct{ID,Kind string}
type Mandate struct{Principal,Representative,Scope string;ExpiresAt time.Time;Evidence string}
type Facility struct{ID,Operator string;Capabilities map[string]bool;Evidence string}
type Graph struct{Entities map[string]Entity;Mandates []Mandate;Facilities map[string]Facility}
func New()*Graph{return &Graph{Entities:map[string]Entity{},Facilities:map[string]Facility{}}}
func(g *Graph)AddEntity(e Entity){g.Entities[e.ID]=e}
func(g *Graph)AddMandate(m Mandate,now time.Time)error{if _,ok:=g.Entities[m.Principal];!ok{return ErrEntity};if _,ok:=g.Entities[m.Representative];!ok{return ErrEntity};if !now.Before(m.ExpiresAt){return ErrExpired};if m.Scope==""||m.Evidence==""{return ErrAuthority};g.Mandates=append(g.Mandates,m);return nil}
func(g *Graph)Authorized(principal,rep,scope string,now time.Time)bool{for _,m:=range g.Mandates{if m.Principal==principal&&m.Representative==rep&&m.Scope==scope&&now.Before(m.ExpiresAt)&&m.Evidence!=""{return true}};return false}
func(g *Graph)AddFacility(f Facility)error{if _,ok:=g.Entities[f.Operator];!ok{return ErrEntity};if f.Evidence==""{return ErrFacility};g.Facilities[f.ID]=f;return nil}
func(g *Graph)Supports(id,cap string)error{f,ok:=g.Facilities[id];if !ok||!f.Capabilities[cap]{return ErrFacility};return nil}
func(g *Graph)FacilityIDs()[]string{x:=make([]string,0,len(g.Facilities));for id:=range g.Facilities{x=append(x,id)};sort.Strings(x);return x}
