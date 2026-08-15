package tradegraph
import("errors";"sort";"time")
var(ErrDuplicateEvent=errors.New("duplicate trade event");ErrMissingEntity=errors.New("trade entity missing");ErrTimeRegression=errors.New("event time regression");ErrInvalidRelation=errors.New("invalid trade relation"))
type Entity struct{ID,Kind string}
type Relation struct{From,Type,To string}
type Event struct{ID,TransactionID,Type string;At time.Time;EntityIDs []string;Evidence []string}
type Graph struct{Entities map[string]Entity;Relations []Relation;Events map[string]Event;LastEvent map[string]time.Time}
func New()*Graph{return &Graph{Entities:map[string]Entity{},Events:map[string]Event{},LastEvent:map[string]time.Time{}}}
func(g *Graph)AddEntity(e Entity){g.Entities[e.ID]=e}
func(g *Graph)Link(r Relation)error{if _,ok:=g.Entities[r.From];!ok{return ErrMissingEntity};if _,ok:=g.Entities[r.To];!ok{return ErrMissingEntity};if r.Type==""{return ErrInvalidRelation};g.Relations=append(g.Relations,r);return nil}
func(g *Graph)Append(e Event)error{
 if _,ok:=g.Events[e.ID];ok{return ErrDuplicateEvent}
 for _,id:=range e.EntityIDs{if _,ok:=g.Entities[id];!ok{return ErrMissingEntity}}
 if t,ok:=g.LastEvent[e.TransactionID];ok&&e.At.Before(t){return ErrTimeRegression}
 g.Events[e.ID]=e;g.LastEvent[e.TransactionID]=e.At;return nil
}
func(g *Graph)Timeline(tx string)[]Event{var out []Event;for _,e:=range g.Events{if e.TransactionID==tx{out=append(out,e)}};sort.Slice(out,func(i,j int)bool{if out[i].At.Equal(out[j].At){return out[i].ID<out[j].ID};return out[i].At.Before(out[j].At)});return out}
