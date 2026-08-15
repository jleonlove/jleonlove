package verifier
import("errors";"math")
var ErrInsufficientEvidence=errors.New("insufficient evidence")
var ErrConflict=errors.New("conflicting evidence")
type Status string
const(Known Status="KNOWN";Supported Status="SUPPORTED";Inferred Status="INFERRED";Conflicted Status="CONFLICTED";Unverified Status="UNVERIFIED")
type Evidence struct{Source string;Independent bool;Supports bool;Reliability float64}
type Claim struct{Text string;GeneratorConfidence float64;Evidence []Evidence}
type Result struct{Status Status;CalibratedConfidence float64;Verified bool}
func Verify(c Claim)Result{
 if len(c.Evidence)==0{return Result{Status:Unverified}}
 var pos,neg float64;ind:=0
 for _,e:=range c.Evidence{r:=math.Max(0,math.Min(1,e.Reliability));if e.Supports{pos+=r}else{neg+=r};if e.Independent{ind++}}
 if pos>0&&neg>0{return Result{Status:Conflicted,CalibratedConfidence:pos/(pos+neg)}}
 if pos==0{return Result{Status:Unverified}}
 conf:=pos/float64(len(c.Evidence));if conf>1{conf=1}
 st:=Inferred;if ind>=2&&conf>=0.75{st=Supported};if ind>=3&&conf>=0.9{st=Known}
 return Result{Status:st,CalibratedConfidence:conf,Verified:st==Supported||st==Known}
}
func RequireVerified(r Result)error{if r.Status==Conflicted{return ErrConflict};if !r.Verified{return ErrInsufficientEvidence};return nil}
