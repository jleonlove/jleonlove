package intelligence
import("errors";"math")
var ErrNoQualifiedModel=errors.New("no qualified model")
type Task struct{Complexity,Risk,Privacy int;NeedsVision,NeedsTools bool;MaxCost,MaxLatency float64}
type Model struct{Name string;Quality,Reasoning,Privacy int;Vision,Tools,Qualified bool;Cost,Latency float64}
type Decision struct{Model string;ReasoningBudget int;Score float64}
func Route(t Task,models []Model)(Decision,error){
 best:=Decision{Score:-math.MaxFloat64}
 for _,m:=range models{
  if !m.Qualified||m.Cost>t.MaxCost||m.Latency>t.MaxLatency{continue}
  if t.NeedsVision&&!m.Vision||t.NeedsTools&&!m.Tools{continue}
  if m.Privacy<t.Privacy{continue}
  score:=float64(m.Quality*3+m.Reasoning*(t.Complexity+1)+m.Privacy*t.Risk)-m.Cost*10-m.Latency
  if score>best.Score{best=Decision{Model:m.Name,ReasoningBudget:budget(t,m),Score:score}}
 }
 if best.Model==""{return Decision{},ErrNoQualifiedModel}
 return best,nil
}
func budget(t Task,m Model)int{b:=128+t.Complexity*256+t.Risk*128;if b>4096{b=4096};return b}
