package containment
import ("crypto/ed25519"; "encoding/json"; "errors"; "time")
var ( ErrInvalidSignature=errors.New("invalid containment signature"); ErrUnknownKey=errors.New("unknown containment signing key") )
type SignedAttestation struct { KeyID string `json:"kid"`; Attestation Attestation `json:"attestation"` }
type Keyring map[string]ed25519.PublicKey
func canonicalUnsigned(a Attestation)([]byte,error){ a.Signature=nil; return json.Marshal(a) }
func Sign(keyID string,key ed25519.PrivateKey,a Attestation)(SignedAttestation,error){ b,e:=canonicalUnsigned(a); if e!=nil{return SignedAttestation{},e}; a.Signature=ed25519.Sign(key,b); return SignedAttestation{KeyID:keyID,Attestation:a},nil }
func VerifySigned(now time.Time,s SignedAttestation,e Expected,keys Keyring) error { pub,ok:=keys[s.KeyID]; if !ok{return ErrUnknownKey}; b,err:=canonicalUnsigned(s.Attestation); if err!=nil{return err}; if !ed25519.Verify(pub,b,s.Attestation.Signature){return ErrInvalidSignature}; return Verify(now,s.Attestation,e) }
