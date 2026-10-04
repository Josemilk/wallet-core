package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/Josemilk/wallet-core/exchange/backend/internal/ledger"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/matching"
    "github.com/Josemilk/wallet-core/exchange/backend/internal/risk"
)

type server struct {
    ledger *ledger.Ledger
    book   *matching.Book
    risk   *risk.Engine
}

func main() {
    addr := os.Getenv("EXCHANGE_HTTP_ADDR")
    if addr == "" { addr = ":8080" }
    s := &server{ledger: ledger.New(), book: matching.NewBook(), risk: risk.New()}
    mux := http.NewServeMux()
    mux.HandleFunc("GET /healthz", s.health)
    mux.HandleFunc("GET /readyz", s.ready)
    mux.HandleFunc("GET /v1/orderbook", s.orderbook)
    mux.HandleFunc("POST /v1/orders", s.placeOrder)
    mux.HandleFunc("POST /v1/risk/check", s.riskCheck)

    srv := &http.Server{Addr: addr, Handler: securityHeaders(mux), ReadHeaderTimeout: 5*time.Second, ReadTimeout: 10*time.Second, WriteTimeout: 10*time.Second, IdleTimeout: 60*time.Second}
    log.Printf("exchange service listening on %s", addr)
    log.Fatal(srv.ListenAndServe())
}

func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("Referrer-Policy", "no-referrer")
        next.ServeHTTP(w, r)
    })
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) { writeJSON(w, 200, map[string]any{"status":"ok"}) }
func (s *server) ready(w http.ResponseWriter, r *http.Request) {
    if err := s.ledger.Ready(context.Background()); err != nil { writeJSON(w, 503, map[string]any{"status":"not_ready"}); return }
    writeJSON(w, 200, map[string]any{"status":"ready"})
}
func (s *server) orderbook(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, s.book.Snapshot(r.URL.Query().Get("symbol"))) }
func (s *server) placeOrder(w http.ResponseWriter, r *http.Request) { matching.HandleOrderHTTP(w, r, s.book, s.ledger, s.risk) }
func (s *server) riskCheck(w http.ResponseWriter, r *http.Request) { risk.HandleHTTP(w, r, s.risk) }
func writeJSON(w http.ResponseWriter, status int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(status); _ = json.NewEncoder(w).Encode(v) }
