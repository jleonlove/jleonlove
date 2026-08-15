package multimodal
import("errors";"sort";"strings")
var(ErrModality=errors.New("modality denied");ErrConsent=errors.New("capture consent required");ErrAuthority=errors.New("action authority denied");ErrProvenance=errors.New("provenance required"))
type Modality string
const(Text Modality="TEXT";Audio Modality="AUDIO";Image Modality="IMAGE";Video Modality="VIDEO";Screen Modality="SCREEN")
type Input struct{Modality Modality;Source,Session string;Consent bool;Trusted bool}
type Action struct{Capability string;Consequential bool;Authorized bool}
type Policy struct{Modalities map[Modality]bool}
func Admit(p Policy,i Input)error{if !p.Modalities[i.Modality]{return ErrModality};if i.Source==""||i.Session==""{return ErrProvenance};if (i.Modality==Audio||i.Modality==Video||i.Modality==Screen)&&!i.Consent{return ErrConsent};return nil}
func Authorize(a Action)error{if a.Consequential&&!a.Authorized{return ErrAuthority};return nil}
func ContextLabel(i Input)string{trust:="UNTRUSTED";if i.Trusted{trust="TRUSTED"};return strings.Join([]string{string(i.Modality),i.Source,i.Session,trust},"|")}
func Fingerprint(inputs []Input)string{x:=make([]string,len(inputs));for n,i:=range inputs{x[n]=ContextLabel(i)};sort.Strings(x);return strings.Join(x,"\x00")}
