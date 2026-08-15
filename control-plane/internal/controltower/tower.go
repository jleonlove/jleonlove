package controltower
import("sort";"time")
type Alert struct{ID,TransactionID,Kind string;Severity int;At time.Time;Resolved bool}
type Tower struct{Alerts map[string]Alert}
func New()*Tower{return &Tower{Alerts:map[string]Alert{}}}
func(t *Tower)Upsert(a Alert){t.Alerts[a.ID]=a}
func(t *Tower)Open(tx string)[]Alert{var o []Alert;for _,a:=range t.Alerts{if a.TransactionID==tx&&!a.Resolved{o=append(o,a)}};sort.Slice(o,func(i,j int)bool{if o[i].Severity==o[j].Severity{return o[i].ID<o[j].ID};return o[i].Severity>o[j].Severity});return o}