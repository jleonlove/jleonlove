package runtimeapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"atlas/internal/productionreadiness"
)

type Server struct {
	Token    string
	Started  time.Time
	Requests atomic.Uint64
	Checks   []productionreadiness.Check
}

type response struct {
	OK       bool     `json:"ok"`
	Status   string   `json:"status,omitempty"`
	RequestID string  `json:"request_id,omitempty"`
	Failures []string `json:"failures,omitempty"`
}

func New(token string) *Server {
	return &Server{Token: token, Started: time.Now().UTC()}
}

func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		write(w, http.StatusOK, response{OK: true, Status: "alive"})
	})
	m.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		report := productionreadiness.Evaluate(s.Checks)
		if !report.Ready {
			write(w, http.StatusServiceUnavailable, response{OK: false, Status: "not_ready", Failures: report.Failures})
			return
		}
		write(w, http.StatusOK, response{OK: true, Status: "ready"})
	})
	m.HandleFunc("/v1/execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			write(w, http.StatusMethodNotAllowed, response{OK: false, Status: "method_not_allowed"})
			return
		}
		if s.Token == "" || strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ") != s.Token {
			write(w, http.StatusUnauthorized, response{OK: false, Status: "unauthorized"})
			return
		}
		s.Requests.Add(1)
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = time.Now().UTC().Format("20060102T150405.000000000Z")
		}
		write(w, http.StatusAccepted, response{OK: true, Status: "accepted", RequestID: id})
	})
	return security(m)
}

func security(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Cache-Control", "no-store")
		n.ServeHTTP(w, r)
	})
}

func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
