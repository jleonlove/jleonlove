package memorygov
import("errors";"testing";"time")
func rec()Record{n:=time.Now();return New("m1","tenant-a","user","internal","remember preference",n,n.Add(time.Hour))}
func TestValidRead(t *testing.T){r:=rec();if _,e:=Read(r,"tenant-a",time.Now());e!=nil{t.Fatal(e)}}
func TestCrossTenantDenied(t *testing.T){r:=rec();if _,e:=Read(r,"tenant-b",time.Now());!errors.Is(e,ErrTenant){t.Fatal(e)}}
func TestRevokedDenied(t *testing.T){r:=Revoke(rec());if _,e:=Read(r,"tenant-a",time.Now());!errors.Is(e,ErrRevoked){t.Fatal(e)}}
func TestExpiredDenied(t *testing.T){r:=rec();if _,e:=Read(r,"tenant-a",r.ExpiresAt);!errors.Is(e,ErrExpired){t.Fatal(e)}}
func TestTamperDenied(t *testing.T){r:=rec();r.Content="poisoned";if _,e:=Read(r,"tenant-a",time.Now());!errors.Is(e,ErrIntegrity){t.Fatal(e)}}
func TestMemoryCannotGrantAuthority(t *testing.T){r:=rec();r.AuthorityClaims=[]string{"prod.deploy"};if _,e:=Read(r,"tenant-a",time.Now());!errors.Is(e,ErrAuthority){t.Fatal(e)}}
