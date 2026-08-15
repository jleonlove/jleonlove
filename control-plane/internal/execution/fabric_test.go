package execution
import("errors";"testing")
func pol()Policy{return Policy{Capabilities:map[string]bool{"browser.read":true,"computer.click":true,"code.run":true},Destinations:map[string]bool{"https://docs.example":true}}}
func TestBrowserReadAllowed(t *testing.T){e:=Authorize(pol(),Request{Kind:Browser,Capability:"browser.read",Destination:"https://docs.example"});if e!=nil{t.Fatal(e)}}
func TestDestinationDenied(t *testing.T){e:=Authorize(pol(),Request{Kind:Browser,Capability:"browser.read",Destination:"https://evil.example"});if !errors.Is(e,ErrDestination){t.Fatal(e)}}
func TestRawCredentialDenied(t *testing.T){e:=Authorize(pol(),Request{Kind:Browser,Capability:"browser.read",Arguments:map[string]string{"api_key":"x"}});if !errors.Is(e,ErrCredential){t.Fatal(e)}}
func TestConsequentialClickNeedsApproval(t *testing.T){e:=Authorize(pol(),Request{Kind:Computer,Capability:"computer.click",Consequential:true});if !errors.Is(e,ErrApproval){t.Fatal(e)}}
func TestApprovedConsequentialClick(t *testing.T){e:=Authorize(pol(),Request{Kind:Computer,Capability:"computer.click",Consequential:true,Approved:true});if e!=nil{t.Fatal(e)}}
func TestCodeRequiresDisposableSandbox(t *testing.T){e:=Authorize(pol(),Request{Kind:Code,Capability:"code.run"});if !errors.Is(e,ErrSandbox){t.Fatal(e)}}
func TestSandboxedCodeAllowed(t *testing.T){e:=Authorize(pol(),Request{Kind:Code,Capability:"code.run",DisposableSandbox:true});if e!=nil{t.Fatal(e)}}
func TestUnknownCapabilityDenied(t *testing.T){e:=Authorize(pol(),Request{Kind:Browser,Capability:"browser.write"});if !errors.Is(e,ErrCapability){t.Fatal(e)}}
