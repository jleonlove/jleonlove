package productiongate
import"errors"
var ErrGate=errors.New("production gate failed")
type Evidence struct{Milestones map[string]bool;Security,Recovery,Authorization,Regression,TradeLab,ExternalFailure bool}
func Evaluate(e Evidence)error{for i:=71;i<=79;i++{k:=fmtMilestone(i);if !e.Milestones[k]{return ErrGate}};if !e.Security||!e.Recovery||!e.Authorization||!e.Regression||!e.TradeLab||!e.ExternalFailure{return ErrGate};return nil}
func fmtMilestone(i int)string{b:=[]byte{'0','0','0','0','0','0'};for p:=5;p>=0;p--{b[p]=byte('0'+i%10);i/=10};return string(b)}