package a2a
import("errors";"testing")
func gw()Gateway{return Gateway{AllowedEndpoints:map[string]bool{"https://agent.example":true},MinTrust:70}}
func card()Card{return Card{AgentID:"remote-1",Endpoint:"https://agent.example",Capabilities:[]string{"research","read"},Trust:90,Qualified:true}}
func TestQualifiedDelegation(t *testing.T){d,e:=gw().Delegate(card(),[]string{"research","read"},[]string{"research"});if e!=nil||d.AgentID==""{t.Fatal(e)}}
func TestUnqualifiedDenied(t *testing.T){c:=card();c.Qualified=false;if e:=gw().Qualify(c,[]string{"research","read"});!errors.Is(e,ErrUnqualified){t.Fatal(e)}}
func TestEndpointDenied(t *testing.T){c:=card();c.Endpoint="https://evil";if e:=gw().Qualify(c,[]string{"research","read"});!errors.Is(e,ErrEndpoint){t.Fatal(e)}}
func TestTrustDenied(t *testing.T){c:=card();c.Trust=10;if e:=gw().Qualify(c,[]string{"research","read"});!errors.Is(e,ErrTrust){t.Fatal(e)}}
func TestRemoteCapabilityExpansionDenied(t *testing.T){c:=card();c.Capabilities=append(c.Capabilities,"deploy");if e:=gw().Qualify(c,[]string{"research","read"});!errors.Is(e,ErrCapability){t.Fatal(e)}}
func TestDelegationExpansionDenied(t *testing.T){_,e:=gw().Delegate(card(),[]string{"research","read"},[]string{"deploy"});if !errors.Is(e,ErrAuthority){t.Fatal(e)}}
