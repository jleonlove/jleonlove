package transactionpassport
import("crypto/sha256";"encoding/hex";"errors";"sort";"strings")
var ErrPassport=errors.New("invalid transaction passport")
type Passport struct{TransactionID string;Evidence []string;Approved bool;Digest string}
func Seal(p *Passport)error{if p.TransactionID==""||len(p.Evidence)==0||!p.Approved{return ErrPassport};x:=append([]string(nil),p.Evidence...);sort.Strings(x);h:=sha256.Sum256([]byte(p.TransactionID+"\x00"+strings.Join(x,"\x00")));p.Digest=hex.EncodeToString(h[:]);return nil}
func Verify(p Passport)bool{d:=p.Digest;p.Digest="";return Seal(&p)==nil&&p.Digest==d}