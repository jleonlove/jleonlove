package runtimeapi
import("net/http/httptest";"strings";"testing")
func TestHealth(t *testing.T){s:=New("secret");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,httptest.NewRequest("GET","/healthz",nil));if w.Code!=200{t.Fatal(w.Code)}}
func TestUnauthorized(t *testing.T){s:=New("secret");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,httptest.NewRequest("POST","/v1/execute",strings.NewReader(`{}`)));if w.Code!=401{t.Fatal(w.Code)}}
func TestExecute(t *testing.T){s:=New("secret");r:=httptest.NewRequest("POST","/v1/execute",strings.NewReader(`{}`));r.Header.Set("Authorization","Bearer secret");r.Header.Set("X-Request-ID","r1");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=202||s.Requests.Load()!=1{t.Fatal(w.Code)}}
func TestMethod(t *testing.T){s:=New("secret");r:=httptest.NewRequest("GET","/v1/execute",nil);r.Header.Set("Authorization","Bearer secret");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,r);if w.Code!=405{t.Fatal(w.Code)}}
func TestSecurityHeaders(t *testing.T){s:=New("secret");w:=httptest.NewRecorder();s.Handler().ServeHTTP(w,httptest.NewRequest("GET","/healthz",nil));if w.Header().Get("X-Content-Type-Options")!="nosniff"||w.Header().Get("Cache-Control")!="no-store"{t.Fatal(w.Header())}}
