package execution
import("errors";"strings")
var(ErrCapability=errors.New("execution capability denied");ErrDestination=errors.New("destination denied");ErrCredential=errors.New("raw credential denied");ErrApproval=errors.New("approval required");ErrSandbox=errors.New("sandbox required"))
type Kind string
const(Browser Kind="BROWSER";Computer Kind="COMPUTER";Code Kind="CODE")
type Request struct{Kind Kind;Capability,Destination string;Arguments map[string]string;Consequential bool;Approved bool;DisposableSandbox bool}
type Policy struct{Capabilities map[string]bool;Destinations map[string]bool}
func Authorize(p Policy,r Request)error{
 if !p.Capabilities[r.Capability]{return ErrCapability}
 if r.Destination!=""&&!p.Destinations[r.Destination]{return ErrDestination}
 for k:=range r.Arguments{q:=strings.ToLower(k);if q=="password"||q=="api_key"||q=="token"||q=="secret"{return ErrCredential}}
 if r.Consequential&&!r.Approved{return ErrApproval}
 if r.Kind==Code&&!r.DisposableSandbox{return ErrSandbox}
 return nil
}
