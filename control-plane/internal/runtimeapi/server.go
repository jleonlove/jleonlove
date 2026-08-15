package runtimeapi
import("encoding/json";"net/http";"strings";"sync/atomic";"time")
type Server struct{Token string;Started time.Time;Requests atomic.Uint64}
type response struct{OK bool `json:"ok"`;Status string `json:"status,omitempty"`;RequestID string `json:"request_id,omitempty"`}
func New(token string)*Server{return &Server{Token:token,Started:time.Now().UTC()}}
func(s *Server)Handler()http.Handler{m:=http.NewServeMux()
m.HandleFunc("/healthz",func(w http.ResponseWriter,r *http.Request){write(w,200,response{true,"alive",""})})
m.HandleFunc("/readyz",func(w http.ResponseWriter,r *http.Request){write(w,200,response{true,"ready",""})})
m.HandleFunc("/v1/execute",func(w http.ResponseWriter,r *http.Request){if r.Method!="POST"{write(w,405,response{false,"method_not_allowed",""});return};if s.Token==""||strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer ")!=s.Token{write(w,401,response{false,"unauthorized",""});return};s.Requests.Add(1);id:=r.Header.Get("X-Request-ID");if id==""{id=time.Now().UTC().Format("20060102T150405.000000000Z")};write(w,202,response{true,"accepted",id})});return security(m)}
func security(n http.Handler)http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Header().Set("X-Content-Type-Options","nosniff");w.Header().Set("X-Frame-Options","DENY");w.Header().Set("Cache-Control","no-store");n.ServeHTTP(w,r)})}
func write(w http.ResponseWriter,s int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(s);_ = json.NewEncoder(w).Encode(v)}
