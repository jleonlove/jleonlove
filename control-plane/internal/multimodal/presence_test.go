package multimodal
import("errors";"testing")
func pol()Policy{return Policy{Modalities:map[Modality]bool{Text:true,Audio:true,Image:true,Screen:true}}}
func TestTextAdmitted(t *testing.T){if e:=Admit(pol(),Input{Modality:Text,Source:"user",Session:"s"});e!=nil{t.Fatal(e)}}
func TestDisabledVideoDenied(t *testing.T){if e:=Admit(pol(),Input{Modality:Video,Source:"cam",Session:"s",Consent:true});!errors.Is(e,ErrModality){t.Fatal(e)}}
func TestAudioNeedsConsent(t *testing.T){if e:=Admit(pol(),Input{Modality:Audio,Source:"mic",Session:"s"});!errors.Is(e,ErrConsent){t.Fatal(e)}}
func TestScreenNeedsConsent(t *testing.T){if e:=Admit(pol(),Input{Modality:Screen,Source:"display",Session:"s"});!errors.Is(e,ErrConsent){t.Fatal(e)}}
func TestProvenanceRequired(t *testing.T){if e:=Admit(pol(),Input{Modality:Image,Session:"s"});!errors.Is(e,ErrProvenance){t.Fatal(e)}}
func TestVoiceCannotBypassAuthority(t *testing.T){if e:=Authorize(Action{Capability:"payment.send",Consequential:true});!errors.Is(e,ErrAuthority){t.Fatal(e)}}
func TestAuthorizedConsequentialAction(t *testing.T){if e:=Authorize(Action{Capability:"payment.send",Consequential:true,Authorized:true});e!=nil{t.Fatal(e)}}
