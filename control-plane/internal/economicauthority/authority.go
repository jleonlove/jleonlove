package economicauthority
import("errors";"math")
var(ErrCurrency=errors.New("currency denied");ErrSingle=errors.New("single transaction limit exceeded");ErrAggregate=errors.New("aggregate spend limit exceeded");ErrApproval=errors.New("economic approval required");ErrRecipient=errors.New("recipient denied");ErrAmount=errors.New("invalid amount"))
type Grant struct{Currency string;MaxSingle,MaxAggregate float64;Recipients map[string]bool;RequireApprovalAbove float64}
type Request struct{Amount float64;Currency,Recipient string;Approved bool}
type Ledger struct{Spent float64}
func Authorize(g Grant,l *Ledger,r Request)error{
 if math.IsNaN(r.Amount)||math.IsInf(r.Amount,0)||r.Amount<=0{return ErrAmount}
 if r.Currency!=g.Currency{return ErrCurrency}
 if !g.Recipients[r.Recipient]{return ErrRecipient}
 if r.Amount>g.MaxSingle{return ErrSingle}
 if l.Spent+r.Amount>g.MaxAggregate{return ErrAggregate}
 if r.Amount>g.RequireApprovalAbove&&!r.Approved{return ErrApproval}
 l.Spent+=r.Amount;return nil
}
