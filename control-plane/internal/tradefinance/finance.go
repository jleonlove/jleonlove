package tradefinance
import"errors"
var(ErrCollateral=errors.New("insufficient collateral");ErrPayment=errors.New("payment condition unmet"))
type Facility struct{Limit,Drawn,Collateral,Haircut float64}
func Available(f Facility)float64{a:=f.Limit-f.Drawn;c:=f.Collateral*(1-f.Haircut);if c<a{return c};return a}
func Draw(f *Facility,amount float64)error{if amount<=0||amount>Available(*f){return ErrCollateral};f.Drawn+=amount;return nil}
type Payment struct{Amount float64;DocsOK,ComplianceOK,Approved bool}
func Release(p Payment)error{if p.Amount<=0||!p.DocsOK||!p.ComplianceOK||!p.Approved{return ErrPayment};return nil}