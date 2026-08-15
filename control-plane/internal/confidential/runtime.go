package confidential
import("errors";"strings")
var(ErrPlacement=errors.New("execution placement denied");ErrAttestation=errors.New("confidential attestation required");ErrResidency=errors.New("data residency denied");ErrOffline=errors.New("offline capability denied"))
type Mode string
const(Cloud Mode="CLOUD";Edge Mode="EDGE";Confidential Mode="CONFIDENTIAL")
type Workload struct{Mode Mode;Region string;Sensitive,Offline bool;Capability string;Attested bool}
type Policy struct{Regions map[string]bool;OfflineCapabilities map[string]bool}
func Authorize(p Policy,w Workload)error{
 if !p.Regions[w.Region]{return ErrResidency}
 if w.Sensitive&&w.Mode==Cloud{return ErrPlacement}
 if w.Mode==Confidential&&!w.Attested{return ErrAttestation}
 if w.Offline&&!p.OfflineCapabilities[w.Capability]{return ErrOffline}
 return nil
}
func Label(w Workload)string{return strings.Join([]string{string(w.Mode),w.Region,w.Capability},"\x00")}
