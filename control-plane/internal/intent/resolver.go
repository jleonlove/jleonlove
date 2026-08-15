package intent
import("crypto/sha256";"encoding/hex";"encoding/json";"errors";"strings")
var ErrClarificationRequired=errors.New("clarification required")
type Risk string
const(Low Risk="LOW";High Risk="HIGH";Irreversible Risk="IRREVERSIBLE")
type Request struct{Objective string;Risk Risk;Target string;Amount string;Destructive bool}
type Decision struct{Proceed bool;Question string;Options []string;IntentDigest string}
func digest(r Request)string{b,_:=json.Marshal(r);s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
func Resolve(r Request)Decision{
 ambiguous:=strings.TrimSpace(r.Target)=="" || (r.Risk!=Low && strings.TrimSpace(r.Objective)=="")
 if r.Destructive && strings.TrimSpace(r.Target)==""{ambiguous=true}
 if ambiguous{return Decision{Proceed:false,Question:"Confirm the exact target and intended action.",Options:[]string{"Provide details","Cancel"},IntentDigest:digest(r)}}
 return Decision{Proceed:true,IntentDigest:digest(r)}
}
func Confirm(original Request,confirmed Request,expectedDigest string)error{
 if digest(original)!=expectedDigest{return errors.New("intent record changed")}
 if strings.TrimSpace(confirmed.Target)==""{return ErrClarificationRequired}
 return nil
}
