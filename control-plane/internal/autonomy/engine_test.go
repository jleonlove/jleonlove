package autonomy
import("errors";"testing";"time")
func eng()*Engine{return &Engine{Grant:Grant{Capabilities:map[string]bool{"monitor":true,"notify":true},ExpiresAt:time.Unix(2000,0),MaxActions:2}}}
func ev(id string)Event{return Event{ID:id,Capability:"monitor",At:time.Unix(1000,0)}}
func TestAuthorizedEvent(t *testing.T){if e:=eng().Handle(ev("1"));e!=nil{t.Fatal(e)}}
func TestDuplicateDenied(t *testing.T){e:=eng();_ = e.Handle(ev("1"));if x:=e.Handle(ev("1"));!errors.Is(x,ErrDuplicate){t.Fatal(x)}}
func TestExpiredGrant(t *testing.T){e:=eng();x:=ev("1");x.At=time.Unix(3000,0);if z:=e.Handle(x);!errors.Is(z,ErrExpired){t.Fatal(z)}}
func TestCapabilityDenied(t *testing.T){e:=eng();x:=ev("1");x.Capability="pay";if z:=e.Handle(x);!errors.Is(z,ErrAuthority){t.Fatal(z)}}
func TestBudget(t *testing.T){e:=eng();_ = e.Handle(ev("1"));_ = e.Handle(ev("2"));if z:=e.Handle(ev("3"));!errors.Is(z,ErrBudget){t.Fatal(z)}}
func TestConsequentialNeedsApproval(t *testing.T){e:=eng();x:=ev("1");x.Consequential=true;if z:=e.Handle(x);!errors.Is(z,ErrApproval){t.Fatal(z)}}
func TestApprovedConsequential(t *testing.T){e:=eng();x:=Event{ID:"1",Capability:"notify",At:time.Unix(1000,0),Consequential:true,Approved:true};if z:=e.Handle(x);z!=nil{t.Fatal(z)}}
