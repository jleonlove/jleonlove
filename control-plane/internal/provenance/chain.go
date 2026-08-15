package provenance
import("errors";"crypto/sha256";"fmt")
var ErrChain=errors.New("chain of custody invalid")
type Event struct{Asset,Holder,Location,Evidence,PrevHash string}
func Hash(e Event)string{return fmt.Sprintf("%x",sha256.Sum256([]byte(e.Asset+"\x00"+e.Holder+"\x00"+e.Location+"\x00"+e.Evidence+"\x00"+e.PrevHash)))}
func Verify(es []Event)error{prev:="";for _,e:=range es{if e.Asset==""||e.Holder==""||e.Evidence==""||e.PrevHash!=prev{return ErrChain};prev=Hash(e)};return nil}