package credentialbroker
import("errors";"testing";"time")
func pol()Policy{return Policy{AllowedAudiences:map[string]bool{"github":true},AllowedScopes:map[string]bool{"repo.read":true},AllowedResidencies:map[string]bool{"us":true}}}
func grant()Grant{return Grant{UserID:"u",AgentID:"a",Audience:"github",Scopes:[]string{"repo.read"},ExpiresAt:time.Unix(2000,0)}}
func req()Request{return Request{UserID:"u",AgentID:"a",Audience:"github",Scope:"repo.read",Now:time.Unix(1000,0)}}
func TestScopedHandleIssued(t *testing.T){h,e:=Issue(pol(),grant(),req());if e!=nil||h==""{t.Fatal(e)}}
func TestAgentSwapDenied(t *testing.T){r:=req();r.AgentID="evil";if _,e:=Issue(pol(),grant(),r);!errors.Is(e,ErrIdentity){t.Fatal(e)}}
func TestAudienceSwapDenied(t *testing.T){r:=req();r.Audience="bank";if _,e:=Issue(pol(),grant(),r);!errors.Is(e,ErrAudience){t.Fatal(e)}}
func TestScopeExpansionDenied(t *testing.T){r:=req();r.Scope="repo.write";if _,e:=Issue(pol(),grant(),r);!errors.Is(e,ErrScope){t.Fatal(e)}}
func TestExpiredDenied(t *testing.T){r:=req();r.Now=time.Unix(3000,0);if _,e:=Issue(pol(),grant(),r);!errors.Is(e,ErrExpired){t.Fatal(e)}}
func TestRestrictedResidencyDenied(t *testing.T){if e:=AuthorizeData(pol(),Data{Classification:"restricted",Residency:"eu"});!errors.Is(e,ErrSensitive){t.Fatal(e)}}
func TestRawSecretDetection(t *testing.T){if !RawSecret("Authorization: Bearer abc"){t.Fatal("missed")}}
