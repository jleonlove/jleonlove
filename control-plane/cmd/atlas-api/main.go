package main
import("log";"net/http";"os";"time";"atlas/internal/runtimeapi")
func main(){token:=os.Getenv("ATLAS_API_TOKEN");if token==""{log.Fatal("ATLAS_API_TOKEN is required")};addr:=os.Getenv("ATLAS_ADDR");if addr==""{addr=":8080"};s:=&http.Server{Addr:addr,Handler:runtimeapi.New(token).Handler(),ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:30*time.Second,IdleTimeout:60*time.Second};log.Printf("atlas-api listening on %s",addr);log.Fatal(s.ListenAndServe())}
