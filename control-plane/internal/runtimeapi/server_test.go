package runtimeapi

import (
	"net/http/httptest"
	"strings"
	"testing"

	"atlas/internal/productionreadiness"
)

func TestHealth(t *testing.T) {
	s := New("secret")
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/healthz", nil))
	if w.Code != 200 { t.Fatal(w.Code) }
}

func TestReadyFailsClosedWithoutEvidence(t *testing.T) {
	s := New("secret")
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/readyz", nil))
	if w.Code != 503 { t.Fatalf("expected 503, got %d: %s", w.Code, w.Body.String()) }
	if !strings.Contains(w.Body.String(), "live_integrations") || !strings.Contains(w.Body.String(), "end_to_end_trade") {
		t.Fatalf("missing readiness failures: %s", w.Body.String())
	}
}

func TestReadyOnlyWithAllRequiredEvidence(t *testing.T) {
	s := New("secret")
	for _, name := range []string{"live_integrations", "regulatory_data", "observability", "load_chaos", "security_assessment", "disaster_recovery", "red_team", "end_to_end_trade"} {
		s.Checks = append(s.Checks, productionreadiness.Check{Name: name, Pass: true, Evidence: "signed-evidence", Critical: true})
	}
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/readyz", nil))
	if w.Code != 200 { t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String()) }
}

func TestUnauthorized(t *testing.T) {
	s := New("secret")
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("POST", "/v1/execute", strings.NewReader(`{}`)))
	if w.Code != 401 { t.Fatal(w.Code) }
}

func TestExecute(t *testing.T) {
	s := New("secret")
	r := httptest.NewRequest("POST", "/v1/execute", strings.NewReader(`{}`))
	r.Header.Set("Authorization", "Bearer secret")
	r.Header.Set("X-Request-ID", "r1")
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 202 || s.Requests.Load() != 1 { t.Fatal(w.Code) }
}

func TestMethod(t *testing.T) {
	s := New("secret")
	r := httptest.NewRequest("GET", "/v1/execute", nil)
	r.Header.Set("Authorization", "Bearer secret")
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 405 { t.Fatal(w.Code) }
}

func TestSecurityHeaders(t *testing.T) {
	s := New("secret")
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/healthz", nil))
	if w.Header().Get("X-Content-Type-Options") != "nosniff" || w.Header().Get("Cache-Control") != "no-store" { t.Fatal(w.Header()) }
}
