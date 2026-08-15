package economicauthority
import("errors";"math";"testing")
func grant()Grant{return Grant{Currency:"USD",MaxSingle:1000,MaxAggregate:1500,Recipients:map[string]bool{"supplier":true},RequireApprovalAbove:500}}
func TestSmallAuthorized(t *testing.T){l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:100,Currency:"USD",Recipient:"supplier"});e!=nil{t.Fatal(e)}}
func TestCurrencyDenied(t *testing.T){l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:1,Currency:"EUR",Recipient:"supplier"});!errors.Is(e,ErrCurrency){t.Fatal(e)}}
func TestRecipientDenied(t *testing.T){l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:1,Currency:"USD",Recipient:"unknown"});!errors.Is(e,ErrRecipient){t.Fatal(e)}}
func TestSingleLimit(t *testing.T){l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:1001,Currency:"USD",Recipient:"supplier",Approved:true});!errors.Is(e,ErrSingle){t.Fatal(e)}}
func TestAggregateLimit(t *testing.T){l:=Ledger{Spent:1000};if e:=Authorize(grant(),&l,Request{Amount:600,Currency:"USD",Recipient:"supplier",Approved:true});!errors.Is(e,ErrAggregate){t.Fatal(e)}}
func TestApprovalThreshold(t *testing.T){l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:600,Currency:"USD",Recipient:"supplier"});!errors.Is(e,ErrApproval){t.Fatal(e)}}
func TestApprovedSpend(t *testing.T){l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:600,Currency:"USD",Recipient:"supplier",Approved:true});e!=nil||l.Spent!=600{t.Fatal(e)}}
func TestInvalidAmounts(t *testing.T){for _,a:=range []float64{0,-1,math.NaN(),math.Inf(1)}{l:=Ledger{};if e:=Authorize(grant(),&l,Request{Amount:a,Currency:"USD",Recipient:"supplier"});!errors.Is(e,ErrAmount){t.Fatalf("%v %v",a,e)}}}
